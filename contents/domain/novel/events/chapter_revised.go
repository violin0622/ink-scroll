package events

import (
	"encoding/json"
	"time"

	"github.com/oklog/ulid/v2"

	"github.com/violin0622/ink-scroll/contents/domain/novel/vo"
)

var _ Event = (*ChapterRevised)(nil)

type ChapterRevised struct {
	id  ID
	aid ID
	at  time.Time

	RevisedChapter vo.Chapter
}

func NewChapterRevised(aid ID, c vo.Chapter) ChapterRevised {
	return ChapterRevised{
		id:  ulid.Make(),
		at:  time.Now(),
		aid: aid,

		RevisedChapter: c,
	}
}

func (t ChapterRevised) AggID() ID            { return t.aid }
func (t ChapterRevised) AggKind() string      { return `Novel` }
func (t ChapterRevised) EmittedAt() time.Time { return t.at }
func (t ChapterRevised) ID() ID               { return t.id }
func (t ChapterRevised) Kind() string         { return `ChapterRevised` }
func (t ChapterRevised) String() string       { return formatEvent(t) }

func (t ChapterRevised) Payload() ([]byte, error) { return json.Marshal(t) }
