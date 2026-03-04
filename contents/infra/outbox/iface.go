package outbox

import (
	"context"

	"github.com/oklog/ulid/v2"

	novelv1 "play-ddd/proto/gen/go/contents/novel/v1"
)

// EventRepo is interface between outbox processor and event persistent storage.
// (e.g. Postgres, MySQL or WAL)
// It is assumed that events are arranged in an orderly manner by time overall,
// and individual aggregation events are strictly sorted in the order they
// occur.
type EventRepo interface {
	// Process starts a transaction, executes an arbitrary number of commands
	// in fn within the transaction, commit the changes if no error, otherwise
	// rollback.
	Process(ctx context.Context, fn func(EventsBatch) error) error
}

// EventsBatch fetch and processes a batch of events, record process results by
// AdvanceCursor. PollEvents and AdvanceCursor may be called arbitrary times.
// The amount of Events and Result in each call are not guaranteed to be the
// same, but the order must be the same.
type EventsBatch interface {
	// Poll polls events limited by arg, in order. Events don't have to
	// be in same aggregation.
	Poll(Arg) ([]*novelv1.Event, error)

	// AdvanceCursor pass results, expected them to be persistent.
	AdvanceCursor(...Result) error
}

// EventBus is the message bus interface for publishing events.
type EventBus interface {
	Pub(context.Context, []*novelv1.Event) error
}

// Arg contains arguments for polling events.
type Arg struct {
	Instance ID
	Limit    int
}

// Result contains the result of processing an event.
type Result struct {
	ID     ID
	Status Status
	Reason string
}

// Status represents the publishing status of an event.
type Status string

const (
	// StatusNone means event was just inserted and have not processed by
	// outbox.
	StatusNone Status = ""

	// StatusPending means event has not been processed by outbox,
	// typically because oneof previous events of same aggregate processed
	// failed.
	StatusPending Status = "pending"

	StatusPublished Status = "published"
	StatusFailed    Status = "failed"

	StatusDead Status = "dead" // exceeded max retries
)

type (
	ID = ulid.ULID
)
