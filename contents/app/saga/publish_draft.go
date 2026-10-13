package saga

import (
	"context"

	"github.com/violin0622/ink-scroll/contents/app"
	"github.com/violin0622/ink-scroll/contents/app/command"
	"github.com/violin0622/ink-scroll/contents/domain/chapter"
	"github.com/violin0622/ink-scroll/contents/domain/draft"
)

var ach app.CommandHandler

// 1. lock draft
// 2. create chapter from draft
// 3. update novel
// 4. destroy draft
func PublishDraft(ctx context.Context, id draft.ID) error {
	// step 1
	ach.LockDraft()

	// step 2
	cid, err := ach.CreateChapterFromDraft(ctx, command.CreateChapterFromDraft{
		DraftID: id,
	})
	if err != nil {
		return err
	}

	// step 3
	var chapterRepo app.AggRepo[chapter.Chapter, chapter.ID]
	c, err := chapterRepo.Get(ctx, cid)
	if err != nil {
		return err
	}

	err = ach.UploadChapter(ctx, command.UploadChapter{
		ChapterID: cid,
		NovelID:   c.NovelID(),
		VolumeID:  c.VolumeID(),
	})
	if err != nil {
		return err
	}

	return nil
}
