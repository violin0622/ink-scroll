package chapter

import (
	"encoding/json"

	"github.com/oklog/ulid/v2"

	"github.com/violin0622/ink-scroll/common"
)

var formatEvent = common.FormatEvent[ID, ID]

var (
	_ Event = (*ChapterRevised)(nil)
	_ Event = (*ChapterUploaded)(nil)
)

type ChapterRevised struct {
	common.BaseEvent[ID, ID]

	Title        string
	MainContent  string
	ExtraContent string
}

func NewChapterRevised(
	aid ID,
	title, main, extra string,
) ChapterRevised {
	return ChapterRevised{
		BaseEvent: common.NewBaseEvent(ulid.Make(), aid, `Chapter`, `ChapterRevised`),

		Title:        title,
		MainContent:  main,
		ExtraContent: extra,
	}
}

func (t ChapterRevised) Payload() ([]byte, error) { return json.Marshal(t) }

type ChapterUploaded struct {
	common.BaseEvent[ID, ID]

	Title             string
	MainContent       string
	ExtraContent      string
	WordCount         int
	NovelID, VolumeID ID
}

func NewChapterUploaded(aid ID, title, main, extra string, nid, vid ID) ChapterUploaded {
	return ChapterUploaded{
		BaseEvent: common.NewBaseEvent(ulid.Make(), aid, `Chapter`, `ChapterUploaded`),

		Title:        title,
		MainContent:  main,
		ExtraContent: extra,
		NovelID:      nid,
		WordCount:    countWords(main),
		VolumeID:     vid,
	}
}

func (c ChapterUploaded) Payload() ([]byte, error) { return json.Marshal(c) }
