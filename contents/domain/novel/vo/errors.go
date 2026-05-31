package vo

import (
	"errors"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrNoContent              = errors.New(`no content`)
	ErrInvalidChapterSequence = errors.New(`chapter sequence out of range`)
	ErrInvalidVolumeSequence  = errors.New(`volume sequence out of range`)
	ErrVolumeNotFound         = errors.New(`volume not found`)
	ErrChapterNotFound        = errors.New(`chapter not found`)
	ErrDropLastVolume         = errors.New(`last volume`)
)

type VolumeHasChapterError struct{ Cnt int }

func (e VolumeHasChapterError) Error() string {
	return fmt.Sprintf(`remains %d chapters in the volume`, e.Cnt)
}

type VolumeTitleAlreadyExistError struct {
	ID    ID
	Title string
}

func (e VolumeTitleAlreadyExistError) Error() string {
	return fmt.Sprintf(`title '%s' already been used by volume %s`, e.Title, e.ID)
}

func invalidConvertionError(a, b any) error {
	return status.Errorf(
		codes.FailedPrecondition,
		`invalid status convertion '%s'->'%s'`, a, b)
}

type InvariantBrokenError struct {
	Msg string
	Num string
}

func (e InvariantBrokenError) Error() string {
	return fmt.Sprintf(`invariance[%s] broken: %s`, e.Num, e.Msg)
}
