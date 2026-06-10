package events

import (
	"encoding/json"
	"time"

	"github.com/oklog/ulid/v2"

	"play-ddd/contents/domain/novel/vo"
)

var _ Event = (*CoverUploaded)(nil)

type CoverUploaded struct {
	id  ID        `json:"-"`
	aid ID        `json:"-"`
	at  time.Time `json:"-"`

	Cover vo.Cover `json:"cover"`
}

func NewCoverUploaded(aid ID, cover vo.Cover) CoverUploaded {
	return CoverUploaded{
		id:    ulid.Make(),
		at:    time.Now(),
		aid:   aid,
		Cover: cover,
	}
}

func (t CoverUploaded) AggID() ID                { return t.aid }
func (t CoverUploaded) AggKind() string          { return `Novel` }
func (t CoverUploaded) EmittedAt() time.Time     { return t.at }
func (t CoverUploaded) ID() ID                   { return t.id }
func (t CoverUploaded) Kind() string             { return `CoverUploaded` }
func (t CoverUploaded) String() string           { return formatEvent(t) }
func (t CoverUploaded) Payload() ([]byte, error) { return json.Marshal(t) }
