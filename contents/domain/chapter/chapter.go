package chapter

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/oklog/ulid/v2"

	"github.com/violin0622/ink-scroll/common"
)

var _ Aggregate = Chapter{}

type Chapter struct {
	id           ID
	seq          int
	title        string
	mainContent  string
	extraContent string
	releaseAt    time.Time

	er EventRepo
}

func (c Chapter) ID() ID               { return c.id }
func (c Chapter) Kind() string         { return `Chapter` }
func (c Chapter) Title() string        { return c.title }
func (c Chapter) MainContent() string  { return c.mainContent }
func (c Chapter) ExtraContent() string { return c.extraContent }
func (c Chapter) WordCount() int       { return countWords(c.mainContent) }
func (c Chapter) ReleaseAt() time.Time { return c.releaseAt }

func New(er EventRepo) Chapter { return Chapter{er: er} }

// either title, mainContent, exContent can be empty, which means no change.
func (c *Chapter) Revise(
	ctx context.Context,
	title, mainContent, exContent string,
) error {
	e := ChapterRevised{
		id:  ulid.Make(),
		aid: c.id,
		Seq: c.seq,
	}

	if title != `` {
		e.Title = title
	}
	if mainContent != `` {
		e.MainContent = mainContent
	}
	if exContent != `` {
		e.ExtraContent = exContent
	}

	if err := c.imposeReviseEvent(e); err != nil {
		return fmt.Errorf(`revise chapter: %w`, err)
	}

	return c.finish(ctx, e)
}

func (c *Chapter) imposeReviseEvent(e ChapterRevised) error {
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

func (c Chapter) checkInvariants() error {
	return nil
}

var Now = time.Now

func (c *Chapter) ReplayEvents(es ...Event) error {
	if len(es) == 0 {
		return nil
	}

	if _, ok := es[0].(ChapterUploaded); !ok {
		return common.ErrInitialEvent
	}

	for i := range es {
		if err := c.imposeEvent(es[i]); err != nil {
			return fmt.Errorf(`replay events: %w`, err)
		}
	}

	return c.checkInvariants()
}

func (c *Chapter) imposeEvent(e Event) error {
	switch e := e.(type) {
	case ChapterRevised:
		return c.imposeReviseEvent(e)
	default:
		return common.ErrUnknownEventKind(e.Kind())
	}
}

func (n *Chapter) finish(ctx context.Context, es ...Event) error {
	if err := n.checkInvariants(); err != nil {
		return common.NewInvariantsBrokenError(err)
	}

	return n.er.Append(ctx, es...)
}

func countWords(s string) int {
	return len(strings.Split(s, ` `))
}
