package events

import (
	"encoding/json"
	"time"

	"github.com/oklog/ulid/v2"

	"play-ddd/contents/domain/novel/vo"
)

var _ Event = (*DescUpdated)(nil)

type DescUpdated struct {
	id  ID        `json:"-"`
	aid ID        `json:"-"`
	at  time.Time `json:"-"`

	Desc vo.Description `json:"desc"`
}

func NewDescUpdated(aid ID, desc vo.Description) DescUpdated {
	return DescUpdated{
		id:   ulid.Make(),
		at:   time.Now(),
		aid:  aid,
		Desc: desc,
	}
}

func (t DescUpdated) AggID() ID                { return t.aid }
func (t DescUpdated) AggKind() string          { return `Novel` }
func (t DescUpdated) EmittedAt() time.Time     { return t.at }
func (t DescUpdated) ID() ID                   { return t.id }
func (t DescUpdated) Kind() string             { return `DescUpdated` }
func (t DescUpdated) String() string           { return formatEvent(t) }
func (t DescUpdated) Payload() ([]byte, error) { return json.Marshal(t) }

// func (t DescUpdated) restore(
// 	id, aid ID,
// 	at time.Time,
// 	payload []byte,
// ) (Event, error) {
// 	t.id, t.aid, t.at = id, aid, at
// 	return t, json.Unmarshal(payload, &t)
// }
