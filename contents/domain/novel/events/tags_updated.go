package events

import (
	"encoding/json"
	"time"

	"github.com/oklog/ulid/v2"

	"github.com/violin0622/ink-scroll/contents/domain/novel/vo"
)

var _ Event = (*TagsUpdated)(nil)

type TagsUpdated struct {
	id  ID
	aid ID
	at  time.Time

	PrevTags, CurrentTags []vo.Tag
}

func NewTagsUpdated(aid ID, prev, cur []vo.Tag) TagsUpdated {
	return TagsUpdated{
		id:          ulid.Make(),
		at:          time.Now(),
		aid:         aid,
		PrevTags:    prev,
		CurrentTags: cur,
	}
}

func (t TagsUpdated) AggID() ID                { return t.aid }
func (t TagsUpdated) AggKind() string          { return `Novel` }
func (t TagsUpdated) EmittedAt() time.Time     { return t.at }
func (t TagsUpdated) ID() ID                   { return t.id }
func (t TagsUpdated) Kind() string             { return `TagsUpdated` }
func (t TagsUpdated) String() string           { return formatEvent(t) }
func (t TagsUpdated) Payload() ([]byte, error) { return json.Marshal(t) }
