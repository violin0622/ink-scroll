package saga

import (
	"context"
	"fmt"
	"time"

	"github.com/oklog/ulid/v2"

	"github.com/violin0622/ink-scroll/common"
	"github.com/violin0622/ink-scroll/contents/domain/chapter"
	"github.com/violin0622/ink-scroll/contents/domain/novel"
	"github.com/violin0622/ink-scroll/utils/xerr"
)

type (
	ID   = ulid.ULID
	SAGA struct {
		id        ID
		novelRepo NovelRepo
		steps     []Step
		nextStep  int
		direction direction

		common.EventList[ID, ID]
	}
)

type direction int

const (
	_ direction = iota
	forward
	backward
)

func OnEvent(ctx context.Context, e common.Event[ID, ID]) (err error) {
	defer xerr.Expect(&err, `Handle event`)

	var saga *SAGA
	// TODO: load saga by e.AggID()

	switch e := e.(type) {
	case StepDone:
		saga.nextStep = e.NextStep
		saga.run(ctx)
		saga.Events()
	}

	return err
}

func (s *SAGA) run(ctx context.Context) (err error) {
	defer xerr.Expect(&err, `run`)

	switch s.direction {
	case forward:
		return s.runBackward(ctx)
	case backward:
		return s.runForward(ctx)
	default:
		return fmt.Errorf(
			`unknown direction %v, raise a bug issue`,
			s.direction,
		)
	}
}

func (s *SAGA) runForward(ctx context.Context) (err error) {
	defer xerr.Expect(&err, `run forward`)

	// we're done.
	if s.nextStep >= len(s.steps) {
		return err
	}

	ok, err := s.steps[s.nextStep].Execute(ctx)
	if ok && err == nil {
		s.nextStep += 1
		return err
	}

	return nil
}

func (s *SAGA) runBackward(ctx context.Context) (err error) {
	defer xerr.Expect(&err, `run backward`)

	return nil
}

type Step interface {
	Execute(context.Context) (bool, error)
	Compensate(context.Context) (bool, error)
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

var _ common.Event[ID, ID] = StepDone{}

type StepDone struct {
	id  ID
	aid ID
	at  time.Time

	SAGAID   ID
	NextStep int
}

// AggID implements [common.Event].
func (n StepDone) AggID() ID { return n.aid }

// AggKind implements [common.Event].
func (n StepDone) AggKind() string { return `SAGA` }

// EmittedAt implements [common.Event].
func (n StepDone) EmittedAt() time.Time { return n.at }

// ID implements [common.Event].
func (n StepDone) ID() ID { return n.id }

// Kind implements [common.Event].
func (n StepDone) Kind() string { return `Next` }

// Payload implements [common.Event].
func (n StepDone) Payload() ([]byte, error) { return nil, nil }

// String implements [common.Event].
func (n StepDone) String() string { return common.FormatEvent[ID, ID](n) }
