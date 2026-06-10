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
	seq int,
	title, mainContent, exContent string) (
	Chapter, error,
) {
	if seq < 1 {
		return Chapter{}, fmt.Errorf(`invalid sequence`)
	}

	c := Chapter{
		id:           ulid.Make(),
		seq:          seq,
		title:        title,
		mainContent:  mainContent,
		extraContent: exContent,
		er:           cf.er,
	}

	err := cf.er.Append(ctx, ChapterUploaded{
		id:  ulid.Make(),
		aid: c.id,
		at:  Now(),
		// Seq:          seq,
		Title:        c.title,
		MainContent:  c.mainContent,
		ExtraContent: c.extraContent,
		WordCount:    countWords(c.mainContent),
	})
	if err != nil {
		return Chapter{}, fmt.Errorf(`upload chapter: %w`, err)
	}
	return c, nil
}
