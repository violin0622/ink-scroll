package events

import (
	"time"

	"github.com/oklog/ulid/v2"
)

var _ Event = (*PublishReviewSubmitted)(nil)

type PublishReviewSubmitted struct {
	id  ID        `json:"-"`
	aid ID        `json:"-"`
	at  time.Time `json:"-"`
}

func NewPublishReviewSubmitted(aid ID) PublishReviewSubmitted {
	return PublishReviewSubmitted{
		id:  ulid.Make(),
		at:  time.Now(),
		aid: aid,
	}
}

func (t PublishReviewSubmitted) AggID() ID            { return t.aid }
func (t PublishReviewSubmitted) AggKind() string      { return `Novel` }
func (t PublishReviewSubmitted) EmittedAt() time.Time { return t.at }
func (t PublishReviewSubmitted) ID() ID               { return t.id }

func (t PublishReviewSubmitted) Kind() string { return `PublishReviewSubmitted` }

func (t PublishReviewSubmitted) String() string { return formatEvent(t) }

func (t PublishReviewSubmitted) Payload() ([]byte, error) { return emptyPayload() }
