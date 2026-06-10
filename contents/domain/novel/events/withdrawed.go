package events

import (
	"time"

	"github.com/oklog/ulid/v2"
)

var _ Event = (*Withdrawed)(nil)

type Withdrawed struct {
	id  ID        `json:"-"`
	aid ID        `json:"-"`
	at  time.Time `json:"-"`
}

func NewWithdrawed(aid ID) Withdrawed {
	return Withdrawed{
		id:  ulid.Make(),
		at:  time.Now(),
		aid: aid,
	}
}

func (t Withdrawed) AggID() ID                { return t.aid }
func (t Withdrawed) AggKind() string          { return `Novel` }
func (t Withdrawed) EmittedAt() time.Time     { return t.at }
func (t Withdrawed) ID() ID                   { return t.id }
func (t Withdrawed) Kind() string             { return `Withdrawed` }
func (t Withdrawed) String() string           { return formatEvent(t) }
func (t Withdrawed) Payload() ([]byte, error) { return emptyPayload() }
