package events

import (
	"time"

	"github.com/oklog/ulid/v2"
)

var _ Event = (*NovelUnbanned)(nil)

type NovelUnbanned struct {
	id  ID        `json:"-"`
	aid ID        `json:"-"`
	at  time.Time `json:"-"`
}

func NewNovelUnbanned(aid ID) NovelUnbanned {
	return NovelUnbanned{
		id:  ulid.Make(),
		at:  time.Now(),
		aid: aid,
	}
}

func (t NovelUnbanned) AggID() ID                { return t.aid }
func (t NovelUnbanned) AggKind() string          { return `Novel` }
func (t NovelUnbanned) EmittedAt() time.Time     { return t.at }
func (t NovelUnbanned) ID() ID                   { return t.id }
func (t NovelUnbanned) Kind() string             { return `NovelUnbanned` }
func (t NovelUnbanned) String() string           { return formatEvent(t) }
func (t NovelUnbanned) Payload() ([]byte, error) { return emptyPayload() }
