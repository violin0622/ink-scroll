package vo

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrNoContent              = status.New(codes.OutOfRange, `no content`).Err()
	ErrInvalidChapterSequence = status.New(codes.OutOfRange, `chapter sequence out of range`).
					Err()
)

func invalidConvertionError(a, b any) error {
	return status.Errorf(
		codes.FailedPrecondition,
		`invalid status convertion '%s'->'%s'`, a, b)
}
