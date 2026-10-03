package novel

import (
	"github.com/oklog/ulid/v2"

	"github.com/violin0622/ink-scroll/common"
)

type (
	ID        = ulid.ULID
	AuthorID  = ulid.ULID
	Event     = common.Event[ID, ID]
	EventRepo = common.EventRepo[ID, ID]
	Aggregate = common.Aggregate[ID]
	Repo      = common.AggregateRepo[ID, Novel]
)

var ZeroID = ulid.Zero

type entity interface {
	checkInvariants() error
}
