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

type BaseEvent[AID, EID comparable] struct {
	id            EID
	aid           AID
	at            time.Time
	kind, aggkind string
}

func NewBaseEvent[AID, EID comparable](
	id EID, aid AID, kind, aggKind string,
) BaseEvent[AID, EID] {
	return BaseEvent[AID, EID]{
		id:      id,
		aid:     aid,
		at:      time.Now(),
		kind:    kind,
		aggkind: aggKind,
	}
}

func (e BaseEvent[AID, EID]) ID() EID              { return e.id }
func (e BaseEvent[AID, EID]) AggID() AID           { return e.aid }
func (e BaseEvent[AID, EID]) Kind() string         { return e.kind }
func (e BaseEvent[AID, EID]) AggKind() string      { return e.aggkind }
func (e BaseEvent[AID, EID]) EmittedAt() time.Time { return e.at }
func (e BaseEvent[AID, EID]) String() string       { return FormatEvent(e) }

func (e BaseEvent[AID, EID]) Payload() ([]byte, error) { panic("unimplemented") }

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
