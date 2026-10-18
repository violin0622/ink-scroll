package saga

import (
	"context"
	"fmt"

	"github.com/oklog/ulid/v2"

	"github.com/violin0622/ink-scroll/contents/app"
	"github.com/violin0622/ink-scroll/contents/app/command"
	"github.com/violin0622/ink-scroll/contents/domain/chapter"
	"github.com/violin0622/ink-scroll/contents/domain/draft"
	"github.com/violin0622/ink-scroll/contents/domain/novel"
)

var ach app.CommandHandler

// PublishChapterSAGA workflows:
// 1. lock draft
// 2. create chapter from draft
// 3. update novel
// 4. destroy draft
// PublishChapterSAGA is a special aggregate, which also produce and source from
// events.
type PublishChapterSAGA struct {
	id  ulid.ULID
	did draft.ID
	cid chapter.ID
	nid novel.ID
	vid novel.ID
	s   Status
	ach app.CommandHandler
	dr  app.AggRepo[draft.Draft, draft.ID]
}

type Status int

const (
	// begin status
	Pending Status = iota
	// terminate status
	Done
	Failed

	// transient status
	LockingDraft
	CreatingChapter
	UpdatingNovel
	UnlockingDraft
	DroppingChapter
	DroppingDraft
)

func NewSAGA(
	did draft.ID,
	dr app.AggRepo[draft.Draft, draft.ID],
) *PublishChapterSAGA {
	return &PublishChapterSAGA{
		id:  ulid.Make(),
		did: did,
		s:   Pending,
		ach: ach,
		dr:  dr,
	}
}

type SAGARepo = app.AggRepo[PublishChapterSAGA, ID]

func DoSAGA(ctx context.Context, sagaID ID, repo SAGARepo) error {
	saga, err := repo.Get(ctx, sagaID)
	if err != nil {
		return fmt.Errorf(`load saga: %w`, err)
	}
	if saga.Terminated() {
		return nil
	}

	saga.Execute(ctx)
	if err = repo.Save(ctx, saga); err != nil {
		return fmt.Errorf(`save saga: %w`, err)
	}

	return nil
}

func (pcs *PublishChapterSAGA) Terminated() bool {
	return pcs.s == Done || pcs.s == Failed
}

func (pcs *PublishChapterSAGA) Execute(ctx context.Context) {
	switch pcs.s {
	case Pending:
		d, err := pcs.dr.Get(ctx, pcs.did)
		if err != nil {
			pcs.s = Failed
			return
		}
		pcs.nid = d.NovelID()
		pcs.vid = d.VolumeID()
		pcs.s = LockingDraft
	case LockingDraft:
		if err := pcs.ach.LockDraft(
			ctx,
			command.LockDraft{ID: pcs.did},
		); err != nil {
			pcs.s = Failed
			return
		}
		pcs.s = CreatingChapter
	case CreatingChapter:
		cid, err := pcs.ach.CreateChapter(
			ctx,
			command.CreateChapter{DraftID: pcs.did},
		)
		if err != nil {
			pcs.s = UnlockingDraft
			return
		}
		pcs.cid = cid
		pcs.s = UpdatingNovel
	case UpdatingNovel:
		err := pcs.ach.UploadChapter(ctx, command.UploadChapter{
			ChapterID: pcs.cid,
			NovelID:   pcs.nid,
			VolumeID:  pcs.vid,
		})
		if err != nil {
			pcs.s = DroppingChapter
			return
		}
		pcs.s = DroppingDraft
	case DroppingDraft:
		if err := pcs.ach.DeleteDraft(
			ctx,
			command.DeleteDraft{ID: pcs.did},
		); err != nil {
			// retry
			return
		}
		pcs.s = Done
	case DroppingChapter:
		if err := pcs.ach.DeleteChapter(
			ctx,
			command.DeleteChapter{ID: pcs.cid},
		); err != nil {
			// retry
			return
		}
		pcs.s = UnlockingDraft
	case UnlockingDraft:
		if err := pcs.ach.UnlockDraft(
			ctx,
			command.UnlockDraft{ID: pcs.did},
		); err != nil {
			return
		}
		pcs.s = Failed
	default:
		panic(fmt.Sprintf("unexpected saga.Status: %#v", pcs.s))
	}
}
