package chapter

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	"github.com/oklog/ulid/v2"
)

type Factory struct {
	er  EventRepo
	log logr.Logger
}

func NewFactory(log logr.Logger) Factory {
	return Factory{log: log}
}

func (f Factory) WithEventRepo(er EventRepo) Factory {
	return Factory{
		log: f.log,
		er:  er,
	}
}

func (cf Factory) UploadChapter(
	ctx context.Context,
	nid, vid ID,
	title, mainContent, exContent string) (
	c Chapter, err error,
) {
	c.er = cf.er

	e := NewChapterUploaded(ulid.Make(),
		title,
		mainContent,
		exContent,
		nid, vid,
	)

	if err = c.imposeUploadedEvent(e); err != nil {
		return Chapter{}, fmt.Errorf(`upload chapter: %w`, err)
	}

	if err = c.finish(ctx, e); err != nil {
		return Chapter{}, fmt.Errorf(`upload chapter: %w`, err)
	}

	return c, nil
}
