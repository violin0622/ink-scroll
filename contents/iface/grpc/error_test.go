package grpc

import (
	"errors"
	"strings"
	"testing"

	epb "google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	grpcstatus "google.golang.org/grpc/status"
)

func TestUnknownErrorWithDebugStacktraceSkipsHelper(t *testing.T) {
	t.Parallel()

	err := unknownErrorForStacktraceTest()

	st := grpcstatus.Convert(err)
	if got := st.Code(); got != codes.Unknown {
		t.Fatalf("code = %s, want %s", got, codes.Unknown)
	}
	if got := st.Message(); got != "Unknown error: boom" {
		t.Fatalf("message = %q, want %q", got, "Unknown error: boom")
	}

	debug := findStatusDetail[*epb.DebugInfo](t, st.Details())
	entries := debug.GetStackEntries()
	if len(entries) == 0 {
		t.Fatal("debug stack entries are empty")
	}
	if strings.Contains(entries[0], "unknownErrorWithDebugStacktrace") {
		t.Fatalf(
			"first stack entry = %q, want helper to be skipped",
			entries[0],
		)
	}
	if !strings.Contains(entries[0], "unknownErrorForStacktraceTest") {
		t.Fatalf("first stack entry = %q, want caller frame", entries[0])
	}
}

func unknownErrorForStacktraceTest() error {
	return unknownErrorWithDebugStacktrace(errors.New("boom"), 0)
}

func findStatusDetail[T any](t *testing.T, details []any) T {
	t.Helper()

	var zero T
	for _, detail := range details {
		if value, ok := detail.(T); ok {
			return value
		}
	}

	t.Fatalf("detail %T not found in %v", zero, details)
	return zero
}
