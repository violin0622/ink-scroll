package events

import (
	"encoding/json"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/samber/lo"

	"github.com/violin0622/ink-scroll/contents/domain/novel/vo"
	novelv1 "github.com/violin0622/ink-scroll/proto/gen/go/contents/novel/v1"
	eventv1 "github.com/violin0622/ink-scroll/proto/gen/go/event/v1"
)

// var _ fromPB[novelv1.NovelCreated] = NovelCreated{}

// NovelCreated is init event of novel agg.
type NovelCreated struct {
	id  ID        `json:"-"`
	aid ID        `json:"-"`
	at  time.Time `json:"-"`

	AuthorID ID
	Title    vo.Title
	Desc     vo.Description
	Tags     []vo.Tag
	Category vo.Category
	Cover    vo.Cover
}

func NewNovelCreated(
	aid ID,
	autherID ID,
	title vo.Title,
	cover vo.Cover,
	desc vo.Description,
	category vo.Category,
	tags []vo.Tag,
	// c vo.Chapter,
) NovelCreated {
	return NovelCreated{
		id:  ulid.Make(),
		at:  time.Now(),
		aid: aid,

		AuthorID: autherID,
		Title:    title,
		Desc:     desc,
		Tags:     tags,
		Category: category,
		Cover:    cover,
		// FirstChapter: c,
	}
}

func (t NovelCreated) AggID() ID                { return t.aid }
func (t NovelCreated) AggKind() string          { return `Novel` }
func (t NovelCreated) EmittedAt() time.Time     { return t.at }
func (t NovelCreated) ID() ID                   { return t.id }
func (t NovelCreated) Kind() string             { return `NovelCreated` }
func (t NovelCreated) String() string           { return formatEvent(t) }
func (t NovelCreated) Payload() ([]byte, error) { return json.Marshal(t) }
func (t *NovelCreated) FromPB(e *eventv1.Event) error {
	t.id = e.GetId().Into()
	t.aid = e.GetAggregateId().Into()
	t.at = e.GetEmitAt().AsTime()

	var payload novelv1.NovelCreatedPayload
	if err := e.GetPayload().UnmarshalTo(&payload); err != nil {
		return err
	}
	t.Title = vo.Title(payload.GetTitle())
	t.AuthorID = payload.GetAuthorId().Into()
	t.Desc = vo.Description(payload.GetDescription())
	t.Tags = lo.Map(payload.GetTags(), func(t string, _ int) vo.Tag {
		return vo.Tag(t)
	})
	t.Category = vo.Category(payload.GetCategory())
	return nil
}
