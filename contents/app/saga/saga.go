package saga

import (
	"context"

	"github.com/oklog/ulid/v2"

	"play-ddd/contents/domain/chapter"
	"play-ddd/contents/domain/novel"
)

type ID = ulid.ULID
type SAGA struct {
	id        ID
	novelRepo NovelRepo
	steps     []SAGAStep
	curStep   int
	direction string
}

func (s *SAGA) foo() {
	step := s.steps[s.curStep]
	if s.direction == `backward` {
		if err := step.Compensate(); err != nil {
			return
		}
	}

	if err := step.Exec(); err != nil {
		s.direction = `backward`
	}
}

type SAGAStep interface {
	Exec() error
	Compensate() error
}

type NovelRepo interface {
	Get(context.Context, ID) (novel.Novel, error)
	Save(context.Context, novel.Novel) error
}

type OnChapterPublished struct {
	novelRepo NovelRepo
}

func (s OnChapterPublished) Exec(
	ctx context.Context,
	e chapter.ChapterUploaded,
) error {
	// update novel

	n, err := s.novelRepo.Get(ctx, e.NovelID)
	if err != nil {
		return err
	}

	err = n.UploadChapter(ctx, novel.UpdateArg{
		ChapterID:  e.AggID(),
		VolumeID:   e.VolumeID,
		WordCount:  e.WordCount,
		ReleasedAt: e.EmittedAt(),
	})
	if err != nil {
		return err
	}

	if err = s.novelRepo.Save(ctx, n); err != nil {
		return err
	}

	return nil
}
