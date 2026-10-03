package chapter

import (
	"github.com/oklog/ulid/v2"

	"github.com/violin0622/ink-scroll/common"
	"github.com/violin0622/ink-scroll/contents/domain"
)

type (
	ID        = ulid.ULID
	Event     = common.Event[ID, ID]
	EventRepo = domain.EventRepo[ID, ID]
	Aggregate = common.Aggregate[ID]
	Repo      = common.AggregateRepo[ID, Chapter]
)

var ZeroID = ulid.Zero
