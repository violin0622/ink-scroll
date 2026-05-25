package status

import (
	"google.golang.org/grpc/codes"
	grpcstatus "google.golang.org/grpc/status"
	"google.golang.org/protobuf/protoadapt"
)

func NewWithDetails(
	code codes.Code, msg string, details ...protoadapt.MessageV1,
) error {
	st := grpcstatus.New(code, msg)
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

func InvalidArgumentWithDetails(msg string, details ...protoadapt.MessageV1) error {
	return NewWithDetails(codes.InvalidArgument, msg, details...)
}

func DeadlineExceededWithDetails(msg string, details ...protoadapt.MessageV1) error {
	return NewWithDetails(codes.DeadlineExceeded, msg, details...)
}

func NotFoundWithDetails(msg string, details ...protoadapt.MessageV1) error {
	return NewWithDetails(codes.NotFound, msg, details...)
}

func AlreadyExistsWithDetails(msg string, details ...protoadapt.MessageV1) error {
	return NewWithDetails(codes.AlreadyExists, msg, details...)
}

func PermissionDeniedWithDetails(msg string, details ...protoadapt.MessageV1) error {
	return NewWithDetails(codes.PermissionDenied, msg, details...)
}

func ResourceExhaustedWithDetails(msg string, details ...protoadapt.MessageV1) error {
	return NewWithDetails(codes.ResourceExhausted, msg, details...)
}

func FailedPreconditionWithDetails(msg string, details ...protoadapt.MessageV1) error {
	return NewWithDetails(codes.FailedPrecondition, msg, details...)
}

func AbortedWithDetails(msg string, details ...protoadapt.MessageV1) error {
	return NewWithDetails(codes.Aborted, msg, details...)
}

func OutOfRangeWithDetails(msg string, details ...protoadapt.MessageV1) error {
	return NewWithDetails(codes.OutOfRange, msg, details...)
}

func UnimplementedWithDetails(msg string, details ...protoadapt.MessageV1) error {
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

func UnauthenticatedWithDetails(msg string, details ...protoadapt.MessageV1) error {
	return NewWithDetails(codes.Unauthenticated, msg, details...)
}
