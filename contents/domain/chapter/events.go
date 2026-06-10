package chapter

import (
	"encoding/json"
	"time"

	"github.com/oklog/ulid/v2"

	"play-ddd/common"
	"play-ddd/contents/domain/novel/vo"
)

var formatEvent = common.FormatEvent[ID, ID]

var _ Event = (*ChapterRevised)(nil)

type ChapterRevised struct {
	id           ID
	aid          ID
	at           time.Time
	Seq          int
	Title        string
	MainContent  string
	ExtraContent string
}

func NewChapterRevised(aid ID, c vo.Chapter) ChapterRevised {
	return ChapterRevised{
		id:  ulid.Make(),
		at:  time.Now(),
		aid: aid,
	}
}

func (t ChapterRevised) AggID() ID                { return t.aid }
func (t ChapterRevised) AggKind() string          { return `Chapter` }
func (t ChapterRevised) EmittedAt() time.Time     { return t.at }
func (t ChapterRevised) ID() ID                   { return t.id }
func (t ChapterRevised) Kind() string             { return `ChapterRevised` }
func (t ChapterRevised) String() string           { return formatEvent(t) }
func (t ChapterRevised) Payload() ([]byte, error) { return json.Marshal(t) }

var _ Event = ChapterUploaded{}

type ChapterUploaded struct {
	id  ID
	aid ID
	at  time.Time

	// Seq               int
	Title             string
	MainContent       string
	ExtraContent      string
	WordCount         int
	NovelID, VolumeID ID
}

func (t ChapterUploaded) AggID() ID                { return t.aid }
func (t ChapterUploaded) AggKind() string          { return `Chapter` }
func (t ChapterUploaded) EmittedAt() time.Time     { return t.at }
func (t ChapterUploaded) ID() ID                   { return t.id }
func (t ChapterUploaded) Kind() string             { return `ChapterUploaded` }
func (t ChapterUploaded) String() string           { return formatEvent(t) }
func (t ChapterUploaded) Payload() ([]byte, error) { return json.Marshal(t) }
