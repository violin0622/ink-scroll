package status

import (
	"slices"
	"time"

	epb "google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/protoadapt"
	"google.golang.org/protobuf/types/known/durationpb"
)

func NewWithDetails(
	code codes.Code, msg string, details ...protoadapt.MessageV1,
) error {
	st := status.New(code, msg)
	st, err := st.WithDetails(details...)
	if err != nil {
		return err
	}

	return st.Err()
}

func CanceledWithDetails(msg string, details ...protoadapt.MessageV1) error {
	return NewWithDetails(codes.Canceled, msg, details...)
}

func UnknownWithDetails(msg string, details ...protoadapt.MessageV1) error {
	return NewWithDetails(codes.Unknown, msg, details...)
}

func InvalidArgumentWithDetails(
	msg string,
	details ...protoadapt.MessageV1,
) error {
	return NewWithDetails(codes.InvalidArgument, msg, details...)
}

func DeadlineExceededWithDetails(
	msg string,
	details ...protoadapt.MessageV1,
) error {
	return NewWithDetails(codes.DeadlineExceeded, msg, details...)
}

func NotFoundWithDetails(msg string, details ...protoadapt.MessageV1) error {
	return NewWithDetails(codes.NotFound, msg, details...)
}

func AlreadyExistsWithDetails(
	msg string,
	details ...protoadapt.MessageV1,
) error {
	return NewWithDetails(codes.AlreadyExists, msg, details...)
}

func PermissionDeniedWithDetails(
	msg string,
	details ...protoadapt.MessageV1,
) error {
	return NewWithDetails(codes.PermissionDenied, msg, details...)
}

func ResourceExhaustedWithDetails(
	msg string,
	details ...protoadapt.MessageV1,
) error {
	return NewWithDetails(codes.ResourceExhausted, msg, details...)
}

func FailedPreconditionWithDetails(
	msg string,
	details ...protoadapt.MessageV1,
) error {
	return NewWithDetails(codes.FailedPrecondition, msg, details...)
}

func AbortedWithDetails(msg string, details ...protoadapt.MessageV1) error {
	return NewWithDetails(codes.Aborted, msg, details...)
}

func OutOfRangeWithDetails(msg string, details ...protoadapt.MessageV1) error {
	return NewWithDetails(codes.OutOfRange, msg, details...)
}

func UnimplementedWithDetails(
	msg string,
	details ...protoadapt.MessageV1,
) error {
	return NewWithDetails(codes.Unimplemented, msg, details...)
}

func InternalWithDetails(msg string, details ...protoadapt.MessageV1) error {
	return NewWithDetails(codes.Internal, msg, details...)
}

func UnavailableWithDetails(msg string, details ...protoadapt.MessageV1) error {
	return NewWithDetails(codes.Unavailable, msg, details...)
}

func DataLossWithDetails(msg string, details ...protoadapt.MessageV1) error {
	return NewWithDetails(codes.DataLoss, msg, details...)
}

func UnauthenticatedWithDetails(
	msg string,
	details ...protoadapt.MessageV1,
) error {
	return NewWithDetails(codes.Unauthenticated, msg, details...)
}

type builder struct {
	s       *status.Status
	details []protoadapt.MessageV1
}

func Canceled(msg string, a ...any) *builder {
	return newBuilder(codes.Canceled, msg, a...)
}

func Unknown(msg string, a ...any) *builder {
	return newBuilder(codes.Unknown, msg, a...)
}

func InvalidArgument(msg string, a ...any) *builder {
	return newBuilder(codes.InvalidArgument, msg, a...)
}

func DeadlineExceeded(msg string, a ...any) *builder {
	return newBuilder(codes.DeadlineExceeded, msg, a...)
}

func NotFound(msg string, a ...any) *builder {
	return newBuilder(codes.NotFound, msg, a...)
}

func AlreadyExists(msg string, a ...any) *builder {
	return newBuilder(codes.AlreadyExists, msg, a...)
}

func PermissionDenied(msg string, a ...any) *builder {
	return newBuilder(codes.PermissionDenied, msg, a...)
}

func ResourceExhausted(msg string, a ...any) *builder {
	return newBuilder(codes.ResourceExhausted, msg, a...)
}

func FailedPrecondition(msg string, a ...any) *builder {
	return newBuilder(codes.FailedPrecondition, msg, a...)
}

func Aborted(msg string, a ...any) *builder {
	return newBuilder(codes.Aborted, msg, a...)
}

func OutOfRange(msg string, a ...any) *builder {
	return newBuilder(codes.OutOfRange, msg, a...)
}

func Unimplemented(msg string, a ...any) *builder {
	return newBuilder(codes.Unimplemented, msg, a...)
}

func Internal(msg string, a ...any) *builder {
	return newBuilder(codes.Internal, msg, a...)
}

func Unavailable(msg string, a ...any) *builder {
	return newBuilder(codes.Unavailable, msg, a...)
}

func DataLoss(msg string, a ...any) *builder {
	return newBuilder(codes.DataLoss, msg, a...)
}

func Unauthenticated(msg string, a ...any) *builder {
	return newBuilder(codes.Unauthenticated, msg, a...)
}

func newBuilder(code codes.Code, msg string, a ...any) *builder {
	return &builder{s: status.Newf(code, msg, a...)}
}

type Args = map[string]string

func (b *builder) WithDetails(details ...protoadapt.MessageV1) *builder {
	b.details = append(b.details, details...)
	return b
}

func (b *builder) WithErrorInfo(
	reason, domain string,
	md Args,
) *builder {
	return b.WithDetails(&epb.ErrorInfo{
		Reason:   reason,
		Domain:   domain,
		Metadata: md,
	})
}

func (b *builder) WithBadRequestFieldViolation(field, desc string) *builder {
	br := b.badRequest()
	br.FieldViolations = append(
		br.FieldViolations,
		&epb.BadRequest_FieldViolation{
			Field:       field,
			Description: desc,
		},
	)
	return b
}

func (b *builder) WithPreconditionViolation(
	typ, subject, desc string,
) *builder {
	pf := b.preconditionFailure()
	pf.Violations = append(pf.Violations, &epb.PreconditionFailure_Violation{
		Type:        typ,
		Subject:     subject,
		Description: desc,
	})

	return b
}

func (b *builder) WithQuotaViolation(subject, desc string) *builder {
	qf := b.quotaFailure()
	qf.Violations = append(qf.Violations, &epb.QuotaFailure_Violation{
		Subject:     subject,
		Description: desc,
	})
	return b
}

func (b *builder) WithRetryInfo(delay time.Duration) *builder {
	return b.WithDetails(&epb.RetryInfo{RetryDelay: durationpb.New(delay)})
}

func (b *builder) WithDebugInfo(stackEntries []string, detail string) *builder {
	return b.WithDetails(&epb.DebugInfo{
		StackEntries: slices.Clone(stackEntries),
		Detail:       detail,
	})
}

func (b *builder) WithHelpLink(desc, url string) *builder {
	help := b.help()
	help.Links = append(help.Links, &epb.Help_Link{
		Description: desc,
		Url:         url,
	})
	return b
}

func (b *builder) WithLocalizedMessage(locale, msg string) *builder {
	return b.WithDetails(&epb.LocalizedMessage{
		Locale:  locale,
		Message: msg,
	})
}

func (b *builder) WithRequestInfo(requestID, servingData string) *builder {
	return b.WithDetails(&epb.RequestInfo{
		RequestId:   requestID,
		ServingData: servingData,
	})
}

func (b *builder) WithResourceInfo(typ, name, owner, desc string) *builder {
	return b.WithDetails(&epb.ResourceInfo{
		ResourceType: typ,
		ResourceName: name,
		Owner:        owner,
		Description:  desc,
	})
}

func (b *builder) Err() error {
	if len(b.details) == 0 {
		return b.s.Err()
	}

	st, err := b.s.WithDetails(b.details...)
	if err != nil {
		return err
	}

	return st.Err()
}

func (b *builder) Status() (*status.Status, error) {
	st, err := b.s.WithDetails(b.details...)
	if err != nil {
		return nil, err
	}

	return st, nil
}

func (b *builder) badRequest() *epb.BadRequest {
	for _, d := range b.details {
		if br, ok := d.(*epb.BadRequest); ok {
			return br
		}
	}

	br := &epb.BadRequest{}
	b.details = append(b.details, br)
	return br
}

func (b *builder) preconditionFailure() *epb.PreconditionFailure {
	for _, d := range b.details {
		if pf, ok := d.(*epb.PreconditionFailure); ok {
			return pf
		}
	}

	pf := &epb.PreconditionFailure{}
	b.details = append(b.details, pf)
	return pf
}

func (b *builder) quotaFailure() *epb.QuotaFailure {
	for _, d := range b.details {
		if qf, ok := d.(*epb.QuotaFailure); ok {
			return qf
		}
	}

	qf := &epb.QuotaFailure{}
	b.details = append(b.details, qf)
	return qf
}

func (b *builder) help() *epb.Help {
	for _, d := range b.details {
		if help, ok := d.(*epb.Help); ok {
			return help
		}
	}

	help := &epb.Help{}
	b.details = append(b.details, help)
	return help
}
