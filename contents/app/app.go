package app

import (
	"context"

	"github.com/go-logr/logr"

	"github.com/violin0622/ink-scroll/contents/domain/chapter"
	"github.com/violin0622/ink-scroll/contents/domain/draft"
	"github.com/violin0622/ink-scroll/contents/domain/novel"
)

type CommandHandler struct {
	novelFact novel.Factory
	cf        chapter.Factory
	df        draft.Factory
	repo      Repo
	log       logr.Logger

	stats Stats
}

type Stats struct {
	reqRetryTimes map[string]uint64
	reqRetries    map[string]uint64
}

func NewCommandHandler(
	repo Repo,
	log logr.Logger,
) CommandHandler {
	return CommandHandler{
		repo: repo,
		log:  log,
		stats: Stats{
			reqRetryTimes: map[string]uint64{},
			reqRetries:    map[string]uint64{},
		},
	}
}

func NewQueryHandler() QueryHandler {
	return QueryHandler{}
}

type QueryHandler struct {
	// repo Repo
	// log logr.Logger
}

type Repo interface {
	Draft() AggRepo[draft.Draft, draft.ID]
	Chapter() ChapterRepo
	Novel() NovelRepo

	Event() novel.EventRepo
	Tx(func(Repo) error) error
}

type AggRepo[A any, B comparable] interface {
	Get(context.Context, B) (A, error)
	Save(context.Context, A) error
	Drop(context.Context, A) error
}

type ChapterRepo interface {
	// chapter.EventRepo
	Get(context.Context, chapter.ID) (chapter.Chapter, error)
	Save(context.Context, chapter.Chapter) error
}

type NovelRepo interface {
	ByAuthorAndTitle(context.Context, novel.ID, string) (novel.Novel, error)
	Get(context.Context, novel.ID) (novel.Novel, error)
	Save(context.Context, novel.Novel) error
}
