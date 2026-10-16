package draft

import (
	"context"

	"github.com/go-logr/logr"
	"github.com/oklog/ulid/v2"

	"github.com/violin0622/ink-scroll/contents/domain"
	"github.com/violin0622/ink-scroll/contents/domain/novel"
)

type (
	EventRepo = domain.EventRepo[ID, ID]
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

func (f Factory) NewDraft(
	ctx context.Context,
	nid novel.ID, vid novel.ID,
	title, main, ex string,
) (d Draft, err error) {
	e := NewDraftCreated(ulid.Make(), nid, vid, title, main, ex)
	if err = d.imposeCreatedEvent(e); err != nil {
		return d, err
	}

	err = f.er.Append(ctx, e)
	return d, err
}
