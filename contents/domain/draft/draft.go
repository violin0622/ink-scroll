package draft

import (
	"context"
	"fmt"
	"time"

	"github.com/oklog/ulid/v2"

	"github.com/violin0622/ink-scroll/common"
)

type (
	ID = ulid.ULID
)

type Draft struct {
	id, vid, nid ID
	title        string
	mainContent  string
	extraContent string
	releaseAt    time.Time
}

func (c Draft) ID() ID               { return c.id }
func (c Draft) VolumeID() ID         { return c.vid }
func (c Draft) NovelID() ID          { return c.nid }
func (c Draft) Title() string        { return c.title }
func (c Draft) MainContent() string  { return c.mainContent }
func (c Draft) ExtraContent() string { return c.extraContent }

func (c *Draft) Revise(
	ctx context.Context,
	title, mainContent, exContent string,
) error {
	e := NewDraftRevised(c.id, title, mainContent, exContent)

	if title != `` {
		e.Title = title
	}
	if mainContent != `` {
		e.MainContent = mainContent
	}
	if exContent != `` {
		e.ExtraContent = exContent
	}

	if err := c.imposeRevisedEvent(e); err != nil {
		return fmt.Errorf(`revise chapter: %w`, err)
	}

	return c.finish(ctx, e)
}

func (c *Draft) ReplayEvents(es ...Event) error {
	if len(es) == 0 {
		return nil
	}

	if _, ok := es[0].(DraftCreated); !ok {
		return common.ErrInitialEvent
	}

	for i := range es {
		if err := c.imposeEvent(es[i]); err != nil {
			return NewReplayEventsError(es[i], i, err)
		}
	}

	return c.checkInvariants()
}

func (c *Draft) imposeEvent(e Event) error {
	switch e := e.(type) {
	case DraftCreated:
		return c.imposeCreatedEvent(e)
	case DraftRevised:
		return c.imposeRevisedEvent(e)
	case DraftDestoryed:
		return c.imposeDestroyedEvent(e)
	default:
		return nil
	}
}

func (c *Draft) imposeRevisedEvent(e DraftRevised) error {
	if e.Title != `` {
		c.title = e.Title
	}
	if e.MainContent != `` {
		c.mainContent = e.MainContent
	}
	if e.ExtraContent != `` {
		c.extraContent = e.ExtraContent
	}
	return nil
}

func (c *Draft) imposeCreatedEvent(e DraftCreated) error {
	c.id = e.AggID()
	c.mainContent = e.MainContent
	c.extraContent = e.MainContent
	return nil
}

func (c *Draft) imposeDestroyedEvent(e DraftDestoryed) error { return nil }

func (c Draft) checkInvariants() error { return nil }

func (c *Draft) finish(context.Context, ...Event) error { return nil }
