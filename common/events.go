package common

import (
	"context"
	"fmt"
	"time"
)

type Event[AID, EID comparable] interface {
	fmt.Stringer

	ID() EID
	Kind() string
	AggID() AID
	AggKind() string
	EmittedAt() time.Time
	Payload() ([]byte, error)
	// Restore([]byte) error
}

type EventRepo[EID, AID comparable] interface {
	// Fetch(context.Context, AID) ([]Event[EID, AID], error)
	Append(context.Context, ...Event[EID, AID]) error
}

func FormatEvent[AID, EID comparable](e Event[AID, EID]) string {
	return fmt.Sprintf(`%s[%v]: %s[%v] @%s`,
		e.AggKind(), e.AggID(), e.Kind(), e.ID(), e.EmittedAt())
}

func Fullname[AID, EID comparable](d string, e Event[AID, EID]) string {
	return fmt.Sprintf(`%s.%s`, e.AggKind(), e.Kind())
}

type EventList[AID, EID comparable] struct {
	events []Event[AID, EID]
}

func (el EventList[AID, EID]) Events() []Event[AID, EID] { return el.events }
func (el *EventList[AID, EID]) Append(e Event[AID, EID]) {
	el.events = append(el.events, e)
}
