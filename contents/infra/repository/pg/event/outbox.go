package event

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"play-ddd/contents/infra/outbox"
	novelv1 "play-ddd/proto/gen/go/contents/novel/v1"
	"play-ddd/utils/ptr"
	"play-ddd/utils/xerr"
)

var (
	_ outbox.EventRepo = eventRepo{}
	// _ outbox.Event     = outboxEvent{}
)

//	type outboxEvent struct {
//		id            outbox.ID
//		aggregateID   outbox.ID
//		aggregateKind string
//		kind          string
//		emitAt        time.Time
//		payload       []byte
//	}
//
// func (o outboxEvent) ID() outbox.ID          { return o.id }
// func (o outboxEvent) AggregateID() outbox.ID { return o.aggregateID }
// func (o outboxEvent) AggregateKind() string  { return o.aggregateKind }
// func (o outboxEvent) Kind() string           { return o.kind }
// func (o outboxEvent) EmitAt() time.Time      { return o.emitAt }
// func (o outboxEvent) Payload() []byte        { return o.payload }
//
// Process implements outbox.EventRepo.
func (e eventRepo) Process(
	ctx context.Context,
	fn func(outbox.EventsBatch) error,
) (err error) {
	defer xerr.Expect(&err, "process")

	tx := e.db.WithContext(ctx).Begin()
	defer e.tryCommit(&err, tx)

	var cursor EventCursor
	err = tx.
		Model(EventCursor{}).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Take(&cursor).
		Error
	if err != nil {
		return err
	}

	eb := &eventBatch{tx: tx, ec: cursor}
	if err = fn(eb); err != nil {
		return err
	}

	err = tx.
		Model(EventCursor{}).
		Where("target_table = ?", eb.ec.TargetTable).
		Updates(eb.ec).
		Error
	if err != nil {
		return err
	}

	return err
}

func (e eventRepo) tryCommit(err *error, tx *gorm.DB) {
	if err == nil || *err == nil {
		*err = tx.Commit().Error
	}

	if *err != nil {
		tx.Rollback()
	}
}

var _ outbox.EventsBatch = (*eventBatch)(nil)

type EventCursor struct {
	TargetTable string
	RelayCursor uint64
}

type eventBatch struct {
	tx *gorm.DB
	ec EventCursor
}

// AdvanceCursor implements outbox.EventsBatch.
func (e *eventBatch) AdvanceCursor(results ...outbox.Result) (err error) {
	defer xerr.Expect(&err, `advance cursor`)

	e.ec.RelayCursor += uint64(len(results))
	for i := range results {
		ev := fromOutbox(results[i])
		err = e.tx.
			Model(&ev).
			Updates(map[string]any{
				"status": ev.Status,
				"reason": ev.Reason,
			}).
			Error
		if err != nil {
			return err
		}
	}

	return err
}

// Poll implements outbox.EventsBatch.
func (e *eventBatch) Poll(arg outbox.Arg) (es []*novelv1.Event, err error) {
	defer xerr.Expect(&err, `poll events`)

	var events []Event
	err = e.tx.
		Model(Event{}).
		Where(`"seq">?`, e.ec.RelayCursor).
		Limit(arg.Limit).
		Order(`"seq"`).
		Find(&events).
		Error
	if err != nil {
		return es, err
	}

	es = make([]*novelv1.Event, len(events))
	for i := range events {
		if es[i], err = events[i].IntoPB(); err != nil {
			return nil, err
		}
	}

	return es, err
}

func fromOutbox(o outbox.Result) (e Event) {
	e.ID = ID(o.ID)
	e.Status = ptr.To(string(o.Status))
	e.Reason = ptr.To(o.Reason)
	return e
}
