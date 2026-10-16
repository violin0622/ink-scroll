package app

import (
	"context"
	"errors"
	"fmt"

	"github.com/samber/lo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/violin0622/ink-scroll/contents/app/command"
	"github.com/violin0622/ink-scroll/contents/domain/chapter"
	"github.com/violin0622/ink-scroll/contents/domain/novel"
	novelvo "github.com/violin0622/ink-scroll/contents/domain/novel/vo"
	"github.com/violin0622/ink-scroll/utils/retry"
)

func (s *CommandHandler) CreateNovel(
	ctx context.Context, cmd command.CreateNovel) (
	id novel.ID, err error,
) {
	retryCount := -1
	id, err = retry.NewWithData[novel.ID](
		retry.Context(ctx),
		retry.LastErrorOnly(true),
	).Do(
		func() (novel.ID, error) {
			retryCount += 1
			return s.createNovel(ctx, cmd)
		},
	)

	if retryCount > 0 {
		s.stats.reqRetries[`CreateNovel`] += 1
		s.stats.reqRetryTimes[`CreateNovel`] += uint64(retryCount)
	}

	if err != nil {
		return novel.ZeroID, fmt.Errorf(`create novel: %w`, err)
	}
	return id, err
}

func (s *CommandHandler) LockDraft(
	ctx context.Context, cmd command.LockDraft,
) error {
	return nil
}

func (s *CommandHandler) UnlockDraft(
	ctx context.Context, cmd command.UnlockDraft,
) error {
	return nil
}

func (s *CommandHandler) CreateDraft(
	ctx context.Context, cmd command.CreateDraft,
) (id chapter.ID, err error) {
	d, err := s.df.NewDraft(ctx,
		cmd.NovelID,
		cmd.VolumeID,
		cmd.Title,
		cmd.MainContent,
		cmd.ExtraContent)
	if err != nil {
		return id, err
	}
	if err = s.repo.Draft().Save(ctx, d); err != nil {
		return id, err
	}

	return d.ID(), nil
}

func (s *CommandHandler) DeleteDraft(
	ctx context.Context, cmd command.DeleteDraft,
) error {
	return nil
}

func (s *CommandHandler) DeleteChapter(
	ctx context.Context, cmd command.DeleteChapter,
) error {
	return nil
}

func (s *CommandHandler) createNovel(
	ctx context.Context, cmd command.CreateNovel) (
	id novel.ID, err error,
) {
	err = s.repo.Tx(func(tx Repo) error {
		err := ensureNovelTitleAvailable(
			ctx,
			tx.Novel(),
			cmd.AuthorID,
			cmd.Title,
		)
		if err != nil {
			return err
		}

		n, err := s.novelFact.
			WithEventRepo(tx.Event()).
			Create(
				ctx,
				cmd.AuthorID,
				cmd.Title,
				cmd.Desc,
				cmd.Category,
				cmd.Tags,
			)
		if err != nil {
			return err
		}

		if err = tx.Novel().Save(ctx, n); err != nil {
			return err
		}

		id = n.ID()
		return nil
	})

	return id, err
}

func ensureNovelTitleAvailable(
	ctx context.Context,
	repo NovelRepo,
	authorID novel.ID,
	title string,
) error {
	_, err := repo.ByAuthorAndTitle(ctx, authorID, title)
	var notfound novel.NotfoundError

	if errors.As(err, &notfound) {
		return nil
	}

	if err != nil {
		return status.Errorf(codes.Internal, `check title unique: %s`, err)
	}

	return novel.ErrTitleAlreadyExist
}

// UploadChapter 先后更新 Chapter 和 Novel 聚合。
// func (s *CommandHandler) UploadChapter(
// 	ctx context.Context, cmd command.UploadChapter) (
// 	err error,
// ) {
// 	err = s.repo.Tx(func(tx Repo) error {
// 		n, err := tx.Novel().Get(ctx, cmd.NovelID)
// 		if err != nil {
// 			return err
// 		}
//
// 		c, err := tx.Chapter().Get(ctx, cmd.ChapterID)
// 		if err != nil {
// 			return err
// 		}
//
// 		err = n.UploadChapter(ctx, novel.UpdateArg{
// 			ChapterID:  c.ID(),
// 			VolumeID:   cmd.VolumeID,
// 			WordCount:  c.WordCount(),
// 			ReleasedAt: c.ReleaseAt(),
// 		})
// 		if err != nil {
// 			return err
// 		}
//
// 		err = tx.Novel().Save(ctx, n)
// 		if err != nil {
// 			return err
// 		}
//
// 		return nil
// 	})
// 	if err != nil {
// 		return fmt.Errorf(`upload chapter: %w`, err)
// 	}
//
// 	return nil
// }

func (s *CommandHandler) CreateChapter(
	ctx context.Context, cmd command.CreateChapter) (
	id chapter.ID, err error,
) {
	d, err := s.repo.Draft().Get(ctx, cmd.DraftID)
	if err != nil {
		return id, err
	}

	c, err := s.cf.UploadChapter(
		ctx,
		d.NovelID(),
		d.VolumeID(),
		d.Title(),
		d.MainContent(),
		d.ExtraContent())
	if err != nil {
		return id, err
	}

	if err = s.repo.Chapter().Save(ctx, c); err != nil {
		return id, err
	}

	return c.ID(), nil
}

func (s *CommandHandler) UploadChapter(
	ctx context.Context, cmd command.UploadChapter,
) error {
	n, err := s.repo.Novel().Get(ctx, cmd.NovelID)
	if err != nil {
		return err
	}

	c, err := s.repo.Chapter().Get(ctx, cmd.ChapterID)
	if err != nil {
		return err
	}

	err = n.UploadChapter(ctx, novel.UpdateArg{
		ChapterID:  cmd.ChapterID,
		VolumeID:   cmd.VolumeID,
		WordCount:  c.WordCount(),
		ReleasedAt: c.ReleaseAt(),
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *CommandHandler) UpdateNovel(
	ctx context.Context, cmd command.UpdateNovelInfo,
) error {
	retries := -1
	err := retry.New(
		retry.Context(ctx),
		retry.LastErrorOnly(true),
	).
		Do(func() error {
			retries += 1
			return s.repo.Tx(s.updateNovel(ctx, cmd))
		})

	if retries > 0 {
		s.stats.reqRetries[`UpdateNovel`] += 1
		s.stats.reqRetryTimes[`UpdateNovel`] += uint64(retries)
	}

	if err != nil {
		return fmt.Errorf(`update novel: %w`, err)
	}

	return nil
}

func (s *CommandHandler) updateNovel(
	ctx context.Context,
	cmd command.UpdateNovelInfo,
) func(Repo) error {
	return func(tx Repo) error {
		novel, err := tx.Novel().Get(ctx, cmd.NovelID)
		if err != nil {
			return err
		}

		tags := lo.Map(cmd.Tags, func(t string, _ int) novelvo.Tag {
			return novelvo.Tag(t)
		})
		tobeAdd, tobeRemove := novelvo.DiffTags(tags, novel.Tags())

		if err := novel.AppendTags(ctx, tobeAdd...); err != nil {
			return err
		}

		if err := novel.RemoveTags(ctx, tobeRemove...); err != nil {
			return err
		}

		if len(cmd.Desc) == 0 {
			return nil
		}

		if err := novel.UpdateDescription(ctx, cmd.Desc); err != nil {
			return err
		}

		return tx.Novel().Save(ctx, novel)
	}
}
