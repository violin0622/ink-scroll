package events

import (
	"encoding/json"
	"time"

	"github.com/oklog/ulid/v2"
)

var _ Event = (*Completed)(nil)

type Completed struct {
	id  ID
	aid ID
	at  time.Time
}

func NewCompleted(aid ID) Completed {
	return Completed{
		id:  ulid.Make(),
		at:  time.Now(),
		aid: aid,
	}
}

func (t Completed) AggID() ID            { return t.aid }
func (t Completed) AggKind() string      { return `Novel` }
func (t Completed) EmittedAt() time.Time { return t.at }
func (t Completed) ID() ID               { return t.id }
func (t Completed) Kind() string         { return `Completed` }
func (t Completed) String() string       { return formatEvent(t) }

func (t Completed) Payload() ([]byte, error) { return json.Marshal(map[string]any{}) }
