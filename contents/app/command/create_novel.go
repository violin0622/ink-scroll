package command

import (
	"github.com/oklog/ulid/v2"

	"github.com/violin0622/ink-scroll/contents/domain/chapter"
	"github.com/violin0622/ink-scroll/contents/domain/draft"
	"github.com/violin0622/ink-scroll/contents/domain/novel"
)

type CreateNovel struct {
	AuthorID ulid.ULID
	Title    string
	Desc     string
	Category string
	Tags     []string
}

type CreateDraft struct {
	NovelID, VolumeID ulid.ULID
	Title             string
	MainContent       string
	ExtraContent      string
}

type CreateChapterFromDraft struct {
	DraftID draft.ID
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
