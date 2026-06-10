package events

import (
	"time"

	"github.com/oklog/ulid/v2"
)

var _ Event = (*NovelBanned)(nil)

type NovelBanned struct {
	id  ID        `json:"-"`
	aid ID        `json:"-"`
	at  time.Time `json:"-"`
}

func NewNovelBanned(aid ID) NovelBanned {
	return NovelBanned{
		id:  ulid.Make(),
		at:  time.Now(),
		aid: aid,
	}
}

func (t NovelBanned) AggID() ID                { return t.aid }
func (t NovelBanned) AggKind() string          { return `Novel` }
func (t NovelBanned) EmittedAt() time.Time     { return t.at }
func (t NovelBanned) ID() ID                   { return t.id }
func (t NovelBanned) Kind() string             { return `NovelBanned` }
func (t NovelBanned) String() string           { return formatEvent(t) }
func (t NovelBanned) Payload() ([]byte, error) { return emptyPayload() }
