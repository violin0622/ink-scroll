// Package nsq is interface layer of DDD.
package nsq

import (
	"context"
	"fmt"

	"github.com/violin0622/ink-scroll/contents/app"
	"github.com/violin0622/ink-scroll/contents/app/command"
	"github.com/violin0622/ink-scroll/contents/domain/chapter"
	eventv1 "github.com/violin0622/ink-scroll/proto/gen/go/event/v1"
)

type EventBus interface {
	Sub(context.Context, string) <-chan *eventv1.Event
}

type consumer struct {
	eb   EventBus
	quit context.Context
	ah   app.CommandHandler
}

func (c *consumer) Run() {
NEXT:
	select {
	case pbe, ok := <-c.eb.Sub(context.TODO(), `chapter.uploaded`):
		if !ok {
			return
		}

		if !ok ||
			pbe.GetAggregateKind() != `chapter` ||
			pbe.GetKind() != `chapter.uploaded` {
			return
		}

		var e chapter.ChapterUploadedV2
		if err := e.FromPB(pbe); err != nil {
			goto NEXT
		}

		c.ah.UploadChapter(
			context.TODO(),
			command.UploadChapter{
				ChapterID: e.AggID(),
				NovelID:   e.Payload().NovelId.Into(),
				VolumeID:  e.Payload().VolumeId.Into(),
			},
		)

	case <-c.quit.Done():
		return
	}

	goto NEXT
}

func OnEvent(ctx context.Context, event *eventv1.Event) error {
	switch event.GetKind() {
	case `chapter.uploaded`:
		var e chapter.ChapterUploadedV2
		if err := e.FromPB(event); err != nil {
			return err
		}
		return onChapterUploaded(ctx, e)
	default:
		return fmt.Errorf(`Unsupported event kind %s.`, event.GetKind())
	}
}

func onChapterUploaded(ctx context.Context, e chapter.ChapterUploadedV2) error {
	return nil
}
