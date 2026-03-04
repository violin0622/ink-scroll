package event

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"play-ddd/common"
	"play-ddd/contents/domain"
	"play-ddd/contents/domain/novel"
)

var _ domain.EventRepo[novel.ID, novel.ID] = eventRepo{}

type eventRepo struct {
	db *gorm.DB
}

func New(db *gorm.DB) eventRepo { return eventRepo{db} }

// Append implements domain.EventRepo.
func (e eventRepo) Append(
	ctx context.Context, es ...common.Event[novel.ID, novel.ID],
) error {
	return e.db.Transaction(func(tx *gorm.DB) error {
		for i := range es {
			var ev Event
			var ver uint64
			var err error
			if err := ev.FromDomain(es[i]); err != nil {
				return fmt.Errorf(`append events: %w`, err)
			}

			err = tx.Model(Event{}).
				WithContext(ctx).
				Select(`COALESCE(MAX(aggregate_seq),0)`).
				Where(`"aggregate_id"=?`, ev.AggregateID).
				Take(&ver).
				Error
			if err != nil {
				return fmt.Errorf(`append events: %w`, err)
			}

			ev.AggregateSeq = ver + 1

			err = tx.WithContext(ctx).
				Omit(`status`, `reason`, `seq`).
				Create(&ev).Error
			if err != nil {
				return fmt.Errorf(`append events: %w`, err)
			}
		}

		return nil
	})
}
