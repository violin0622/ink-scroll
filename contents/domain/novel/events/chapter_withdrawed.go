package events

import (
	"encoding/json"
	"time"

	"github.com/oklog/ulid/v2"

	"github.com/violin0622/ink-scroll/contents/domain/novel/vo"
)

var _ Event = (*ChapterWithdrawed)(nil)

type ChapterWithdrawed struct {
	id  ID
	aid ID
	at  time.Time

	Chapter vo.Chapter
}

func NewChapterWithdrawed(aid ID, c vo.Chapter) ChapterWithdrawed {
	return ChapterWithdrawed{
		id:  ulid.Make(),
		at:  time.Now(),
		aid: aid,

		Chapter: c,
	}
}

func (t ChapterWithdrawed) AggID() ID            { return t.aid }
func (t ChapterWithdrawed) AggKind() string      { return `Novel` }
func (t ChapterWithdrawed) EmittedAt() time.Time { return t.at }
func (t ChapterWithdrawed) ID() ID               { return t.id }

func (t ChapterWithdrawed) Kind() string             { return `ChapterWithdrawed` }
func (t ChapterWithdrawed) String() string           { return formatEvent(t) }
func (t ChapterWithdrawed) Payload() ([]byte, error) { return json.Marshal(t) }
