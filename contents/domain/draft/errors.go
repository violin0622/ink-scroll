package draft

import "fmt"

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
