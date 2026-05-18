package events

import (
	"encoding/json"
	"time"

	"github.com/oklog/ulid/v2"

	"play-ddd/contents/domain/novel/vo"
)

var _ Event = (*NewChapterUploaded)(nil)

type NewChapterUploaded struct {
	id  ID
	aid ID
	at  time.Time

	WordCount  int
	UploadedAt time.Time
	Chapter    vo.Chapter
}

func NewNewChapterUploaded(
	aid ID,
	c vo.Chapter,
	wordcount int,
	uploadAt time.Time,
) NewChapterUploaded {
	return NewChapterUploaded{
		id:         ulid.Make(),
		at:         time.Now(),
		aid:        aid,
		WordCount:  wordcount,
		UploadedAt: uploadAt,
		Chapter:    c,
	}
}

func (t NewChapterUploaded) AggID() ID                { return t.aid }
func (t NewChapterUploaded) AggKind() string          { return `Novel` }
func (t NewChapterUploaded) EmittedAt() time.Time     { return t.at }
func (t NewChapterUploaded) ID() ID                   { return t.id }
func (t NewChapterUploaded) Payload() ([]byte, error) { return json.Marshal(t) }

func (t NewChapterUploaded) Kind() string   { return `NewChapterUploaded` }
func (t NewChapterUploaded) String() string { return formatEvent(t) }
