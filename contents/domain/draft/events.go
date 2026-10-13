package draft

import (
	"encoding/json"

	"github.com/oklog/ulid/v2"

	"github.com/violin0622/ink-scroll/common"
)

type Event = common.Event[ID, ID]

var (
	_ Event = (*DraftRevised)(nil)
	_ Event = (*DraftCreated)(nil)
	_ Event = (*DraftDestoryed)(nil)
)

type DraftCreated struct {
	common.BaseEvent[ID, ID]

	NovelID, VolumeID ID
	Title             string
	MainContent       string
	ExtraContent      string
}

func NewDraftCreated(aid, nid, vid ID, title, main, extra string) DraftCreated {
	return DraftCreated{
		BaseEvent:    common.NewBaseEvent(ulid.Make(), aid, `Draft`, `DraftCreated`),
		NovelID:      nid,
		VolumeID:     vid,
		Title:        title,
		MainContent:  main,
		ExtraContent: extra,
	}
}

type DraftRevised struct {
	common.BaseEvent[ID, ID]

	Title        string
	MainContent  string
	ExtraContent string
}

func NewDraftRevised(aid ID, title, main, extra string) DraftRevised {
	return DraftRevised{
		BaseEvent: common.NewBaseEvent(
			ulid.Make(),
			aid,
			`Draft`,
			`DraftRevised`,
		),
		Title:        title,
		MainContent:  main,
		ExtraContent: extra,
	}
}

func (e DraftRevised) Payload() ([]byte, error) { return json.Marshal(e) }

type DraftDestoryed struct {
	common.BaseEvent[ID, ID]
}

func NewDraftDestroyed(aid ID) DraftDestoryed {
	return DraftDestoryed{
		BaseEvent: common.NewBaseEvent(
			ulid.Make(),
			aid,
			`Draft`,
			`DraftDestoryed`,
		),
	}
}
