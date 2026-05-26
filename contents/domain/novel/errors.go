package novel

import (
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"play-ddd/common"
)

var newInvariantsBrokenError = common.NewInvariantsBrokenError

type NotfoundError struct {
	e error
}

func (n NotfoundError) Error() string {
	return fmt.Sprintf(`not found: %s`, n.e)
}

func (n NotfoundError) Unwarp() error { return n.e }

func NewNotfoundError(e error) error { return NotfoundError{e} }

type ReplayEventsError struct {
	ev  Event
	idx int
	e   error
}

func NewReplayEventsError(ev Event, idx int, e error) error {
	return ReplayEventsError{ev: ev, idx: idx, e: e}
}

func (n ReplayEventsError) Unwarp() error { return n.e }
func (n ReplayEventsError) Error() string {
	return fmt.Sprintf(`replay event index[%d], kind[%s], id[%s]: %s`,
		n.idx, n.ev.Kind(), n.ev.ID(), n.e)
}

var (
	ErrTitleAlreadyExist = status.New(
		codes.AlreadyExists,
		`the author has already created a novel of same title before`).Err()

	ErrNilPointer = status.New(codes.Internal, `nil pointer`).Err()
)

type VolumeError struct {
	NovelID  ID
	VolumeID ID
	Cause    error
}

func (e VolumeError) Error() string {
	return fmt.Sprintf(`novel %s volume %s: %s`, e.NovelID, e.VolumeID, e.Cause)
}
