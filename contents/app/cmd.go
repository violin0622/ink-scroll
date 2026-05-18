package app

import (
	"context"
	"errors"
	"fmt"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"play-ddd/contents/app/command"
	"play-ddd/contents/domain/chapter"
	"play-ddd/contents/domain/novel"
	novelvo "play-ddd/contents/domain/novel/vo"
	"play-ddd/utils/retry"

	"github.com/samber/lo"
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
		})

	if retryCount > 0 {
		s.stats.reqRetries[`CreateNovel`] += 1
		s.stats.reqRetryTimes[`CreateNovel`] += uint64(retryCount)
	}

	if err != nil {
		return novel.ZeroID, fmt.Errorf(`create novel: %w`, err)
	}
	return id, err
}

func (s *CommandHandler) createNovel(
	ctx context.Context, cmd command.CreateNovel) (
	id novel.ID, err error,
) {
	err = s.repo.Tx(func(tx Repo) error {
		err := ensureNovelTitleAvailable(ctx, tx.Novel(), cmd.AuthorID, cmd.Title)
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
				cmd.Tags)
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
func (s *CommandHandler) UploadChapter(
	ctx context.Context, cmd command.UploadChapter) (
	id chapter.ID, err error,
) {
	err = s.repo.Tx(func(tx Repo) error {
		novel, err := tx.Novel().Get(ctx, cmd.NovelID)
		if err != nil {
			return err
		}

		c, err := s.cf.
			WithEventRepo(tx.Event()).
			UploadChapter(
				ctx,
				novel.ChapterCount()+1,
				cmd.Title,
				cmd.MainContent,
				cmd.ExtraContent)
		if err != nil {
			return err
		}

		err = novel.UploadNewChapter(ctx, c.ID(), c.WordCount(), time.Now())
		if err != nil {
			return err
		}

		err = tx.Chapter().Save(ctx, c)
		if err != nil {
			return err
		}

		err = tx.Novel().Save(ctx, novel)
		if err != nil {
			return err
		}

		id = c.ID()
		return nil
	})
	if err != nil {
		return chapter.ZeroID, fmt.Errorf(`upload chapter: %w`, err)
	}

	return id, nil
}

func (s *CommandHandler) UpdateNovel(
	ctx context.Context, cmd command.UpdateNovelInfo,
) error {
	retries := -1
	err := retry.New(
		retry.Context(ctx),
		retry.LastErrorOnly(true)).
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
