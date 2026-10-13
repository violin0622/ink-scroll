package chapter

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/violin0622/ink-scroll/common"
)

var _ Aggregate = Chapter{}

type Chapter struct {
	id, nid, vid ID

	title        string
	mainContent  string
	extraContent string
	releaseAt    time.Time

	er EventRepo
}

func (c Chapter) ID() ID               { return c.id }
func (c Chapter) NovelID() ID          { return c.nid }
func (c Chapter) VolumeID() ID         { return c.vid }
func (c Chapter) Kind() string         { return `Chapter` }
func (c Chapter) Title() string        { return c.title }
func (c Chapter) MainContent() string  { return c.mainContent }
func (c Chapter) ExtraContent() string { return c.extraContent }
func (c Chapter) WordCount() int       { return countWords(c.mainContent) }
func (c Chapter) ReleaseAt() time.Time { return c.releaseAt }

func New(er EventRepo) Chapter { return Chapter{er: er} }

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

// Revise modifies chapter contents. Either title, mainContent, exContent can
// be empty, which means no change.
func (c *Chapter) Revise(
	ctx context.Context,
	title, mainContent, exContent string,
) error {
	e := NewChapterRevised(c.ID(), title, mainContent, exContent)
	if err := c.imposeReviseEvent(e); err != nil {
		return fmt.Errorf(`revise chapter: %w`, err)
	}

	return c.finish(ctx, e)
}

func (c *Chapter) imposeEvent(e Event) error {
	switch e := e.(type) {
	case ChapterRevised:
		return c.imposeReviseEvent(e)
	case ChapterUploaded:
		return c.imposeUploadedEvent(e)
	default:
		return common.ErrUnknownEventKind(e.Kind())
	}
}

func (c *Chapter) imposeUploadedEvent(e ChapterUploaded) error {
	c.id = e.AggID()
	c.nid = e.NovelID
	c.vid = e.VolumeID
	c.mainContent = e.MainContent
	c.extraContent = e.ExtraContent
	c.releaseAt = e.EmittedAt()

	return nil
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

var Now = time.Now

func (n *Chapter) finish(ctx context.Context, es ...Event) error {
	if err := n.checkInvariants(); err != nil {
		return common.NewInvariantsBrokenError(err)
	}

	return n.er.Append(ctx, es...)
}

func (c Chapter) checkInvariants() error {
	return nil
}

func countWords(s string) int {
	return len(strings.Split(s, ` `))
}
