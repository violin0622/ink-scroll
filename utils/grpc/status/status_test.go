package status_test

import (
	"testing"
	"time"

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
		{
			name: "Canceled",
			code: codes.Canceled,
			fn:   statusx.CanceledWithDetails,
		},
		{name: "Unknown", code: codes.Unknown, fn: statusx.UnknownWithDetails},
		{
			name: "InvalidArgument",
			code: codes.InvalidArgument,
			fn:   statusx.InvalidArgumentWithDetails,
		},
		{
			name: "DeadlineExceeded",
			code: codes.DeadlineExceeded,
			fn:   statusx.DeadlineExceededWithDetails,
		},
		{
			name: "NotFound",
			code: codes.NotFound,
			fn:   statusx.NotFoundWithDetails,
		},
		{
			name: "AlreadyExists",
			code: codes.AlreadyExists,
			fn:   statusx.AlreadyExistsWithDetails,
		},
		{
			name: "PermissionDenied",
			code: codes.PermissionDenied,
			fn:   statusx.PermissionDeniedWithDetails,
		},
		{
			name: "ResourceExhausted",
			code: codes.ResourceExhausted,
			fn:   statusx.ResourceExhaustedWithDetails,
		},
		{
			name: "FailedPrecondition",
			code: codes.FailedPrecondition,
			fn:   statusx.FailedPreconditionWithDetails,
		},
		{name: "Aborted", code: codes.Aborted, fn: statusx.AbortedWithDetails},
		{
			name: "OutOfRange",
			code: codes.OutOfRange,
			fn:   statusx.OutOfRangeWithDetails,
		},
		{
			name: "Unimplemented",
			code: codes.Unimplemented,
			fn:   statusx.UnimplementedWithDetails,
		},
		{
			name: "Internal",
			code: codes.Internal,
			fn:   statusx.InternalWithDetails,
		},
		{
			name: "Unavailable",
			code: codes.Unavailable,
			fn:   statusx.UnavailableWithDetails,
		},
		{
			name: "DataLoss",
			code: codes.DataLoss,
			fn:   statusx.DataLossWithDetails,
		},
		{
			name: "Unauthenticated",
			code: codes.Unauthenticated,
			fn:   statusx.UnauthenticatedWithDetails,
		},
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

func TestBuilderCodeHelpers(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		code  codes.Code
		build func() error
	}{
		{
			name: "Canceled", code: codes.Canceled,
			build: func() error {
				return statusx.Canceled("message %d", 7).
					WithErrorInfo("REASON", "domain", nil).
					Err()
			},
		},
		{
			name: "Unknown", code: codes.Unknown,
			build: func() error {
				return statusx.Unknown("message %d", 7).
					WithErrorInfo("REASON", "domain", nil).
					Err()
			},
		},
		{
			name: "InvalidArgument", code: codes.InvalidArgument,
			build: func() error {
				return statusx.InvalidArgument("message %d", 7).
					WithErrorInfo("REASON", "domain", nil).
					Err()
			},
		},
		{
			name: "DeadlineExceeded", code: codes.DeadlineExceeded,
			build: func() error {
				return statusx.DeadlineExceeded("message %d", 7).
					WithErrorInfo("REASON", "domain", nil).
					Err()
			},
		},
		{
			name: "NotFound", code: codes.NotFound,
			build: func() error {
				return statusx.NotFound("message %d", 7).
					WithErrorInfo("REASON", "domain", nil).
					Err()
			},
		},
		{
			name: "AlreadyExists", code: codes.AlreadyExists,
			build: func() error {
				return statusx.AlreadyExists("message %d", 7).
					WithErrorInfo("REASON", "domain", nil).
					Err()
			},
		},
		{
			name: "PermissionDenied", code: codes.PermissionDenied,
			build: func() error {
				return statusx.PermissionDenied("message %d", 7).
					WithErrorInfo("REASON", "domain", nil).
					Err()
			},
		},
		{
			name: "ResourceExhausted", code: codes.ResourceExhausted,
			build: func() error {
				return statusx.ResourceExhausted("message %d", 7).
					WithErrorInfo("REASON", "domain", nil).
					Err()
			},
		},
		{
			name: "FailedPrecondition", code: codes.FailedPrecondition,
			build: func() error {
				return statusx.FailedPrecondition("message %d", 7).
					WithErrorInfo("REASON", "domain", nil).
					Err()
			},
		},
		{
			name: "Aborted", code: codes.Aborted,
			build: func() error {
				return statusx.Aborted("message %d", 7).
					WithErrorInfo("REASON", "domain", nil).
					Err()
			},
		},
		{
			name: "OutOfRange", code: codes.OutOfRange,
			build: func() error {
				return statusx.OutOfRange("message %d", 7).
					WithErrorInfo("REASON", "domain", nil).
					Err()
			},
		},
		{
			name: "Unimplemented", code: codes.Unimplemented,
			build: func() error {
				return statusx.Unimplemented("message %d", 7).
					WithErrorInfo("REASON", "domain", nil).
					Err()
			},
		},
		{
			name: "Internal", code: codes.Internal,
			build: func() error {
				return statusx.Internal("message %d", 7).
					WithErrorInfo("REASON", "domain", nil).
					Err()
			},
		},
		{
			name: "Unavailable", code: codes.Unavailable,
			build: func() error {
				return statusx.Unavailable("message %d", 7).
					WithErrorInfo("REASON", "domain", nil).
					Err()
			},
		},
		{
			name: "DataLoss", code: codes.DataLoss,
			build: func() error {
				return statusx.DataLoss("message %d", 7).
					WithErrorInfo("REASON", "domain", nil).
					Err()
			},
		},
		{
			name: "Unauthenticated", code: codes.Unauthenticated,
			build: func() error {
				return statusx.Unauthenticated("message %d", 7).
					WithErrorInfo("REASON", "domain", nil).
					Err()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assertStatus(t, tt.build(), tt.code, "message 7")
		})
	}
}

func TestBuilderAggregatesRepeatedDetails(t *testing.T) {
	t.Parallel()

	err := statusx.InvalidArgument("invalid request").
		WithBadRequestFieldViolation("title", "must not be empty").
		WithBadRequestFieldViolation("author", "must not be empty").
		WithPreconditionViolation("state", "/novels/1", "must be draft").
		WithPreconditionViolation("state", "/novels/2", "must be released").
		WithQuotaViolation("user:1", "too many requests").
		WithQuotaViolation("project:2", "daily limit exceeded").
		WithHelpLink("docs", "https://example.test/docs").
		WithHelpLink("support", "https://example.test/support").
		Err()

	st := grpcstatus.Convert(err)
	if st.Code() != codes.InvalidArgument {
		t.Fatalf("code = %s, want %s", st.Code(), codes.InvalidArgument)
	}

	details := st.Details()
	badRequest := findDetail[*epb.BadRequest](t, details)
	if got := len(badRequest.GetFieldViolations()); got != 2 {
		t.Fatalf("bad request violations = %d, want 2", got)
	}
	if got := badRequest.GetFieldViolations()[0].GetField(); got != "title" {
		t.Fatalf("first bad request field = %q, want %q", got, "title")
	}

	precondition := findDetail[*epb.PreconditionFailure](t, details)
	if got := len(precondition.GetViolations()); got != 2 {
		t.Fatalf("precondition violations = %d, want 2", got)
	}
	if got := precondition.GetViolations()[1].GetSubject(); got != "/novels/2" {
		t.Fatalf("second precondition subject = %q, want %q", got, "/novels/2")
	}

	quota := findDetail[*epb.QuotaFailure](t, details)
	if got := len(quota.GetViolations()); got != 2 {
		t.Fatalf("quota violations = %d, want 2", got)
	}
	if got := quota.GetViolations()[1].GetDescription(); got != "daily limit exceeded" {
		t.Fatalf(
			"second quota description = %q, want %q",
			got,
			"daily limit exceeded",
		)
	}

	help := findDetail[*epb.Help](t, details)
	if got := len(help.GetLinks()); got != 2 {
		t.Fatalf("help links = %d, want 2", got)
	}
	if got := help.GetLinks()[0].GetUrl(); got != "https://example.test/docs" {
		t.Fatalf(
			"first help link URL = %q, want %q",
			got,
			"https://example.test/docs",
		)
	}
}

func TestBuilderAddsStandaloneDetails(t *testing.T) {
	t.Parallel()

	err := statusx.Unavailable("temporarily unavailable").
		WithDetails(&epb.ErrorInfo{
			Reason: "TEMPORARY",
			Domain: "domain",
		}).
		WithRetryInfo(1500*time.Millisecond).
		WithDebugInfo([]string{"frame1", "frame2"}, "trace-id=abc").
		WithRequestInfo("request-1", "shard-a").
		WithResourceInfo("novel", "novels/1", "user:42", "locked").
		WithLocalizedMessage("en-US", "Retry later.").
		Err()

	st := grpcstatus.Convert(err)
	if st.Code() != codes.Unavailable {
		t.Fatalf("code = %s, want %s", st.Code(), codes.Unavailable)
	}

	details := st.Details()
	info := findDetail[*epb.ErrorInfo](t, details)
	if info.GetReason() != "TEMPORARY" {
		t.Fatalf(
			"error info reason = %q, want %q",
			info.GetReason(),
			"TEMPORARY",
		)
	}

	retry := findDetail[*epb.RetryInfo](t, details)
	if got := retry.GetRetryDelay().AsDuration(); got != 1500*time.Millisecond {
		t.Fatalf("retry delay = %s, want %s", got, 1500*time.Millisecond)
	}

	debug := findDetail[*epb.DebugInfo](t, details)
	if got := len(debug.GetStackEntries()); got != 2 {
		t.Fatalf("debug stack entries = %d, want 2", got)
	}
	if got := debug.GetDetail(); got != "trace-id=abc" {
		t.Fatalf("debug detail = %q, want %q", got, "trace-id=abc")
	}

	request := findDetail[*epb.RequestInfo](t, details)
	if got := request.GetRequestId(); got != "request-1" {
		t.Fatalf("request id = %q, want %q", got, "request-1")
	}

	resource := findDetail[*epb.ResourceInfo](t, details)
	if got := resource.GetResourceName(); got != "novels/1" {
		t.Fatalf("resource name = %q, want %q", got, "novels/1")
	}

	localized := findDetail[*epb.LocalizedMessage](t, details)
	if got := localized.GetMessage(); got != "Retry later." {
		t.Fatalf("localized message = %q, want %q", got, "Retry later.")
	}
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

func findDetail[T any](t *testing.T, details []any) T {
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
