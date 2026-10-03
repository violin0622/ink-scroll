package novel

import (
	"context"
	"errors"
	"slices"
	"time"

	"github.com/go-logr/logr"

	"github.com/violin0622/ink-scroll/common"
	ev "github.com/violin0622/ink-scroll/contents/domain/novel/events"
	"github.com/violin0622/ink-scroll/contents/domain/novel/vo"
)

var _ Aggregate = (*Novel)(nil)

type Novel struct {
	id            ID
	title         vo.Title
	category      vo.Category
	authorID      ID
	desc          vo.Description
	tags          []vo.Tag
	s             vo.Status
	releaseStatus vo.ReleaseStatus
	cover         vo.Cover
	toc           vo.TOC
	wordCount     int
	updatedAt     time.Time
	createdAt     time.Time

	lg logr.Logger
	er EventRepo
}

func (n Novel) ID() ID                          { return n.id }
func (n Novel) Kind() string                    { return `Novel` }
func (n Novel) WordCount() int                  { return n.wordCount }
func (n Novel) Title() vo.Title                 { return n.title }
func (n Novel) AuthorID() ID                    { return n.authorID }
func (n Novel) Description() vo.Description     { return n.desc }
func (n Novel) Tags() []vo.Tag                  { return n.tags }
func (n Novel) Category() vo.Category           { return n.category }
func (n Novel) Cover() vo.Cover                 { return n.cover }
func (n Novel) Status() vo.Status               { return n.s }
func (n Novel) UpdatedAt() time.Time            { return n.updatedAt }
func (n Novel) CreatedAt() time.Time            { return n.createdAt }
func (n Novel) Chapters() []vo.Chapter          { return n.toc.GetChapters() }
func (n Novel) Volumes() []vo.Volume            { return n.toc.Volumes() }
func (n Novel) ChapterCount() int               { return n.toc.ChapterCount() }
func (n Novel) ReleaseStatus() vo.ReleaseStatus { return n.releaseStatus }

func (n *Novel) ReplayEvents(es ...Event) error {
	if len(es) == 0 {
		return nil
	}

	if _, ok := es[0].(ev.NovelCreated); !ok {
		return common.ErrInitialEvent
	}

	for i := range es {
		if err := n.applyEvent(es[i]); err != nil {
			return NewReplayEventsError(es[i], i, err)
		}
	}

	return n.checkInvariants()
}

func (n *Novel) AppendTags(ctx context.Context, tags ...vo.Tag) error {
	prev := slices.Clone(n.tags)
	for i := range tags {
		n.tags = append(n.tags, vo.Tag(tags[i]))
	}

	return n.finish(ctx, ev.NewTagsUpdated(n.id, prev, n.tags))
}

func (n *Novel) UploadCover(ctx context.Context, cover vo.Cover) error {
	e := ev.NewCoverUploaded(n.id, cover)
	if err := n.imposeCoverUploaded(e); err != nil {
		return err
	}

	return n.finish(ctx, e)
}

func (n *Novel) imposeCoverUploaded(e ev.CoverUploaded) error {
	if n == nil {
		return ErrNilPointer
	}

	n.cover = e.Cover
	return nil
}

func (n *Novel) Unban(ctx context.Context) error {
	e := ev.NewNovelUnbanned(n.ID())
	if err := n.imposeNovelUnbanned(e); err != nil {
		return err
	}

	return n.finish(ctx, e)
}

func (n *Novel) imposeNovelUnbanned(ev.NovelUnbanned) error {
	if n == nil {
		return ErrNilPointer
	}

	var err error
	if n.releaseStatus, err = n.releaseStatus.AppealApproved(); err != nil {
		return err
	}

	return nil
}

// Ban implements novel.
func (n *Novel) Ban(ctx context.Context) error {
	e := ev.NewNovelBanned(n.ID())
	if err := n.imposeNovelBanned(e); err != nil {
		return err
	}
	return n.finish(ctx, e)
}

func (n *Novel) imposeNovelBanned(ev.NovelBanned) error {
	if n == nil {
		return ErrNilPointer
	}

	var err error
	if n.releaseStatus, err = n.releaseStatus.Banned(); err != nil {
		return err
	}

	return nil
}

// ApplyRelease implements novel.
func (n *Novel) ApplyPublish(ctx context.Context) error {
	e := ev.NewPublishReviewSubmitted(n.ID())
	if err := n.imposePublishReviewSubmitted(e); err != nil {
		return err
	}

	return n.finish(ctx, e)
}

// WithdrawRelease implements novel.
func (n *Novel) WithdrawRelease(ctx context.Context) error {
	e := ev.NewWithdrawed(n.ID())
	if err := n.imposeWithdrawed(e); err != nil {
		return err
	}

	return n.finish(ctx, e)
}

func (n *Novel) imposePublishReviewSubmitted(ev.PublishReviewSubmitted) error {
	if n == nil {
		return ErrNilPointer
	}

	var err error
	if n.releaseStatus, err = n.releaseStatus.PublishReviewSubmitted(); err != nil {
		return err
	}

	return nil
}

func (n *Novel) imposeWithdrawed(ev.Withdrawed) error {
	if n == nil {
		return ErrNilPointer
	}

	var err error
	if n.releaseStatus, err = n.releaseStatus.PublishReviewRejected(); err != nil {
		return err
	}

	return nil
}

// Release implements novel.
func (n *Novel) Publish(ctx context.Context) error {
	e := ev.NewNovelPublished(n.ID())
	if err := n.imposeNovelPublished(e); err != nil {
		return err
	}

	return n.finish(ctx, e)
}

func (n *Novel) imposeNovelPublished(ev.NovelPublished) error {
	if n == nil {
		return ErrNilPointer
	}

	var err error
	if n.releaseStatus, err = n.releaseStatus.PublishReviewApproved(); err != nil {
		return err
	}

	return nil
}

func (n *Novel) RemoveTags(ctx context.Context, tags ...vo.Tag) error {
	prev := slices.Clone(n.tags)
	for i := range tags {
		n.tags = slices.DeleteFunc(
			n.tags,
			func(t vo.Tag) bool { return t == vo.Tag(tags[i]) },
		)
	}

	return n.finish(ctx, ev.NewTagsUpdated(n.id, prev, n.tags))
}

func (n *Novel) UpdateDescription(ctx context.Context, d string) error {
	desc, err := vo.NewDescription(d)
	if err != nil {
		return err
	}

	n.desc = desc
	return n.finish(ctx, ev.NewDescUpdated(n.id, n.desc))
}

type UpdateArg struct {
	ChapterID, VolumeID ID
	WordCount           int
	ReleasedAt          time.Time
}

func (a *UpdateArg) validate() error {
	if a.ChapterID == ZeroID {
		return errors.New(`chapter id must not be zero`)
	}

	if a.VolumeID == ZeroID {
		return errors.New(`volume id must not be zero`)
	}

	if a.WordCount == 0 {
		return errors.New(`word count must not be zero`)
	}

	if a.ReleasedAt.IsZero() {
		a.ReleasedAt = time.Now()
	}

	return nil
}

func (n *Novel) UploadChapter(
	ctx context.Context, arg UpdateArg,
) error {
	if err := arg.validate(); err != nil {
		return VolumeError{
			NovelID:  n.ID(),
			VolumeID: arg.VolumeID,
			Cause:    err,
		}
	}

	c := vo.Chapter{ID: arg.ChapterID, VolumeID: arg.VolumeID}
	e := ev.NewNewChapterUploaded(n.id, c, arg.WordCount, arg.ReleasedAt)
	if err := n.imposeNewChapterUploaded(e); err != nil {
		return err
	}

	return n.finish(ctx, e)
}

func (n *Novel) WithdrawChapter(
	ctx context.Context, arg UpdateArg,
) error {
	if err := arg.validate(); err != nil {
		return VolumeError{
			NovelID:  n.ID(),
			VolumeID: arg.VolumeID,
			Cause:    err,
		}
	}

	e := ev.NewChapterWithdrawed(n.ID(), vo.Chapter{
		ID:       arg.ChapterID,
		VolumeID: arg.VolumeID,
	})
	if err := n.imposeChapterWithdrawed(e); err != nil {
		return err
	}

	return n.finish(ctx, e)
}

func (n *Novel) DropVolume(vid ID) error {
	e := ev.NewVolumeDropped(n.ID(), vid)
	if err := n.imposeVolumeDropped(e); err != nil {
		return err
	}

	return nil
}

func (n *Novel) imposeVolumeDropped(e ev.VolumeDropped) error {
	toc, err := n.toc.DropVolume(e.VID)
	if err != nil {
		return VolumeError{
			VolumeID: e.VID,
			NovelID:  n.ID(),
			Cause:    err,
		}
	}

	n.toc = toc
	return nil
}

func (n *Novel) ReorderVolume(vid ID, seq int) error {
	e := ev.NewVolumeReordered(n.ID(), vid, seq)
	return n.imposeVolumeReordered(e)
}

func (n *Novel) imposeVolumeReordered(e ev.VolumeReordered) error {
	toc, err := n.toc.ReorderVolume(e.VID, e.Seq)
	if err != nil {
		return VolumeError{
			NovelID:  n.ID(),
			VolumeID: e.VID,
			Cause:    err,
		}
	}
	n.toc = toc
	return nil
}

// func (n *Novel) WithdrawChapter(ctx context.Context) error {
// 	last := n.toc.ChapterCount()
// 	c, err := n.toc.Get(last)
// 	if err != nil {
// 		return fmt.Errorf(`withdraw chapter: %w`, err)
// 	}
//
// 	if err = n.withdrawChapter(); err != nil {
// 		return fmt.Errorf(`withdraw chapter: %w`, err)
// 	}
//
// 	return n.finish(ctx, ev.NewChapterWithdrawed(n.id, c))
// }

// func (n *Novel) ReviseChapter(
// 	ctx context.Context, title string, wc, seq int,
// ) error {
// 	c, err := n.toc.Get(seq)
// 	if err != nil {
// 		return fmt.Errorf(`revise chapter: %w`, err)
// 	}
//
// 	// c.Title = title
// 	// c.UpdatedAt = Now()
// 	// c.WordCount = wc
//
// 	if c, err = n.reviseChapter(c); err != nil {
// 		return fmt.Errorf(`revise chapter: %w`, err)
// 	}
//
// 	return n.finish(ctx, ev.NewChapterRevised(n.id, c))
// }

func (n *Novel) Complete(ctx context.Context) error {
	n.s = vo.Completed
	return n.finish(ctx, ev.NewCompleted(n.id))
}

func (n *Novel) NolongerUpdate(ctx context.Context) error {
	n.s = vo.Desolated
	return n.finish(ctx, ev.NewNolongerUpdate(n.id))
}

// finish chech invariants and emit events.
func (n Novel) finish(ctx context.Context, e ...Event) error {
	if err := n.checkInvariants(); err != nil {
		return newInvariantsBrokenError(err)
	}

	return n.er.Append(ctx, e...)
}

func (n Novel) checkInvariants() error {
	invarants := []func() error{
		n.idMustExist,
		n.statusNotUnknown,
		// n.draftHasNoChapter,
		n.authorMustExist,
		n.atMost10Tags,
	}

	for i := range invarants {
		if err := invarants[i](); err != nil {
			return err
		}
	}

	return nil
}

// func (n Novel) draftHasNoChapter() error {
// 	if n.s == vo.Draft && n.toc.ChapterCount() != 0 {
// 		return status.New(codes.Internal, `Draft novel can't has chapter`).Err()
// 	}
// 	return nil
// }

func (n Novel) statusNotUnknown() error {
	if n.s.String() == `<unknown>` {
		return errors.New(`unknown status`)
	}
	return nil
}

func (n Novel) idMustExist() error {
	if n.id.IsZero() {
		return errors.New(`id must exist`)
	}
	return nil
}

func (n Novel) authorMustExist() error {
	if n.authorID.IsZero() {
		return errors.New(`authorID must exist`)
	}
	return nil
}

// func (n Novel) updatedAtNewstChapterUpload() error {
// 	if n.s == vo.Draft {
// 		return nil
// 	}
//
// 	c, err := n.toc.LatestChapter()
// 	if err != nil {
// 		return err
// 	}
//
// 	if n.updatedAt.Equal(c.UploadedAt) {
// 		return nil
// 	}
//
// 	return status.New(codes.Internal,
// 		`novel's update time should be identical to last chapter's upload time`,
// 	).Err()
// }

// func (n Novel) totalWordCountIsSummitOfChapterWordCount() error {
// 	sum := 0
// 	for _, c := range n.toc.GetChapters() {
// 		sum += c.WordCount
// 	}
//
// 	if sum != n.wordCount {
// 		return fmt.Errorf(
// 			`novel's word count %d should be summit of every chapter's word count %d`,
// 			n.wordCount,
// 			sum,
// 		)
// 	}
//
// 	return nil
// }

func (n Novel) atMost10Tags() error {
	if len(n.tags) > 10 {
		return errors.New(`at most 10 tags allowed`)
	}

	return nil
}

func (n *Novel) applyEvent(e ev.Event) error {
	switch e := e.(type) {
	case ev.NovelCreated:
		return n.imposeNovelCreated(e)
	case ev.NovelPublished:
		return nil
	case ev.PublishReviewSubmitted:
		return n.imposePublishReviewSubmitted(e)
	case ev.Withdrawed:
		return n.imposeWithdrawed(e)
	case ev.CoverUploaded:
		return n.imposeCoverUploaded(e)
	case ev.NovelBanned:
		return n.imposeNovelBanned(e)
	case ev.NovelUnbanned:
		return n.imposeNovelUnbanned(e)
	case ev.TagsUpdated:
		return n.imposeTagsUpdated(e)
	case ev.DescUpdated:
		return n.imposeDescUpdated(e)
	case ev.NewChapterUploaded:
		return n.imposeNewChapterUploaded(e)
	// case ev.ChapterRevised:
	// 	return n.imposeChapterRevised(e)
	// case ev.ChapterWithdrawed:
	// 	return n.imposeChapterWithdrawed(e)
	case ev.NolongerUpdate:
		return n.imposeNolongerUpdate(e)
	case ev.Completed:
		return n.imposeCompleted(e)
	case ev.VolumeDropped:
		return n.imposeVolumeDropped(e)
	case ev.VolumeReordered:
		return n.imposeVolumeReordered(e)
	default:
		return common.ErrUnknownEventKind(e.Kind())
	}
}

func (n *Novel) imposeNovelCreated(e ev.NovelCreated) error {
	if n == nil {
		return errors.New(`nil pointer`)
	}

	n.id = e.AggID()
	n.title = e.Title
	n.authorID = e.AuthorID
	n.tags = e.Tags
	n.cover = e.Cover
	n.s = vo.Serial
	n.releaseStatus = vo.Draft
	n.category = e.Category
	n.createdAt = e.EmittedAt()

	return nil
}

func (n *Novel) imposeTagsUpdated(e ev.TagsUpdated) error {
	n.tags = e.CurrentTags
	return nil
}

func (n *Novel) imposeDescUpdated(e ev.DescUpdated) error {
	n.desc = e.Desc
	return nil
}

func (n *Novel) imposeCompleted(ev.Completed) error {
	var err error
	if n.s, err = n.s.MarkCompleted(); err != nil {
		return err
	}
	return nil
}

func (n *Novel) imposeNolongerUpdate(ev.NolongerUpdate) error {
	var err error
	if n.s, err = n.s.MarkDesolated(); err != nil {
		return err
	}
	return nil
}

func (n *Novel) imposeChapterWithdrawed(e ev.ChapterWithdrawed) (err error) {
	if n == nil {
		return ErrNilPointer
	}

	if n.toc, err = n.toc.DropChapter(e.Chapter.ID); err != nil {
		return err
	}

	return nil
}

func (n *Novel) imposeNewChapterUploaded(e ev.NewChapterUploaded) (err error) {
	if n == nil {
		return ErrNilPointer
	}

	if n.toc, err = n.toc.AppendChapter(e.Chapter); err != nil {
		return err
	}

	if n.s, err = n.s.ChapterPublished(); err != nil {
		return err
	}

	n.wordCount += e.WordCount
	n.updatedAt = e.UploadedAt

	return nil
}

// func (n *Novel) imposeChapterRevised(e ev.ChapterRevised) error {
// 	if n == nil {
// 		return ErrNilPointer
// 	}
//
// 	_, err := n.reviseChapter(e.RevisedChapter)
// 	return err
// }

// func (n *Novel) reviseChapter(c vo.Chapter) (origin vo.Chapter, err error) {
// 	if n.toc, origin, err = n.toc.Revise(c); err != nil {
// 		return origin, err
// 	}
//
// 	delta := origin.WordCount - c.WordCount
// 	n.wordCount -= delta
//
// 	return origin, err
// }

// func (n *Novel) imposeChapterWithdrawed(ev.ChapterWithdrawed) error {
// 	if n == nil {
// 		return ErrNilPointer
// 	}
//
// 	return n.withdrawChapter()
// }

// func (n *Novel) withdrawChapter() error {
// 	if n.s == vo.Completed {
// 		return ErrMutateCompletedNovel
// 	}
//
// 	toc, c, err := n.toc.Pop()
// 	if err != nil {
// 		return err
// 	}
//
// 	n.toc = toc
// 	n.wordCount -= c.WordCount
// 	last := n.toc.ChapterCount()
// 	lastC, err := n.toc.Get(last)
// 	if errors.Is(err, vo.ErrInvalidChapterSequence) {
// 		n.updatedAt = n.createdAt
// 	} else {
// 		n.updatedAt = lastC.UploadedAt
// 	}
//
// 	return nil
// }

// Now is exported only for test, do not modify elsewhere.
var Now = time.Now
