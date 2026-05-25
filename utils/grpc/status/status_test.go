package status_test

import (
	"testing"

	epb "google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	grpcstatus "google.golang.org/grpc/status"
	"google.golang.org/protobuf/protoadapt"

	statusx "play-ddd/utils/grpc/status"
)

func TestCodeHelpersWithDetails(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		code codes.Code
		fn   func(string, ...protoadapt.MessageV1) error
	}{
		{name: "Canceled", code: codes.Canceled, fn: statusx.CanceledWithDetails},
		{name: "Unknown", code: codes.Unknown, fn: statusx.UnknownWithDetails},
		{name: "InvalidArgument", code: codes.InvalidArgument, fn: statusx.InvalidArgumentWithDetails},
		{name: "DeadlineExceeded", code: codes.DeadlineExceeded, fn: statusx.DeadlineExceededWithDetails},
		{name: "NotFound", code: codes.NotFound, fn: statusx.NotFoundWithDetails},
		{name: "AlreadyExists", code: codes.AlreadyExists, fn: statusx.AlreadyExistsWithDetails},
		{name: "PermissionDenied", code: codes.PermissionDenied, fn: statusx.PermissionDeniedWithDetails},
		{name: "ResourceExhausted", code: codes.ResourceExhausted, fn: statusx.ResourceExhaustedWithDetails},
		{name: "FailedPrecondition", code: codes.FailedPrecondition, fn: statusx.FailedPreconditionWithDetails},
		{name: "Aborted", code: codes.Aborted, fn: statusx.AbortedWithDetails},
		{name: "OutOfRange", code: codes.OutOfRange, fn: statusx.OutOfRangeWithDetails},
		{name: "Unimplemented", code: codes.Unimplemented, fn: statusx.UnimplementedWithDetails},
		{name: "Internal", code: codes.Internal, fn: statusx.InternalWithDetails},
		{name: "Unavailable", code: codes.Unavailable, fn: statusx.UnavailableWithDetails},
		{name: "DataLoss", code: codes.DataLoss, fn: statusx.DataLossWithDetails},
		{name: "Unauthenticated", code: codes.Unauthenticated, fn: statusx.UnauthenticatedWithDetails},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := tt.fn("message", &epb.ErrorInfo{
				Reason: "REASON",
				Domain: "domain",
			})

			assertStatus(t, err, tt.code, "message")
		})
	}
}

func TestFaildPreconditionWithDetailsKeepsLegacySpelling(t *testing.T) {
	t.Parallel()

	err := statusx.FaildPreconditionWithDetails("message", &epb.ErrorInfo{
		Reason: "REASON",
		Domain: "domain",
	})

	assertStatus(t, err, codes.FailedPrecondition, "message")
}

func assertStatus(t *testing.T, err error, code codes.Code, msg string) {
	t.Helper()

	st := grpcstatus.Convert(err)
	if st.Code() != code {
		t.Fatalf("code = %s, want %s", st.Code(), code)
	}
	if st.Message() != msg {
		t.Fatalf("message = %q, want %q", st.Message(), msg)
	}

	details := st.Details()
	if len(details) != 1 {
		t.Fatalf("details length = %d, want 1", len(details))
	}

	info, ok := details[0].(*epb.ErrorInfo)
	if !ok {
		t.Fatalf("detail = %T, want *errdetails.ErrorInfo", details[0])
	}
	if info.GetReason() != "REASON" {
		t.Fatalf("detail reason = %q, want %q", info.GetReason(), "REASON")
	}
	if info.GetDomain() != "domain" {
		t.Fatalf("detail domain = %q, want %q", info.GetDomain(), "domain")
	}
}
