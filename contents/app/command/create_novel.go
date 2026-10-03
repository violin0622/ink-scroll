package command

import (
	"github.com/oklog/ulid/v2"

	"github.com/violin0622/ink-scroll/contents/domain/chapter"
	"github.com/violin0622/ink-scroll/contents/domain/novel"
)

type CreateNovel struct {
	AuthorID ulid.ULID
	Title    string
	Desc     string
	Category string
	Tags     []string
}

type UploadChapter struct {
	ChapterID chapter.ID
	NovelID   novel.ID
	VolumeID  novel.ID
}

type UpdateNovelInfo struct {
	NovelID novel.ID
	Desc    string
	Tags    []string
}
