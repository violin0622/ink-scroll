package events

import (
	"github.com/oklog/ulid/v2"

	"github.com/violin0622/ink-scroll/common"
)

type (
	ID    = ulid.ULID
	Event = common.Event[ID, ID]
)

var formatEvent = common.FormatEvent[ID, ID]

func emptyPayload() ([]byte, error) { return []byte(`{}`), nil }
