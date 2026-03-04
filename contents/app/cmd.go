package app

import (
	"context"
	"fmt"

	"play-ddd/common/backoff"
	"play-ddd/contents/app/command"
	"play-ddd/contents/domain/chapter"
	"play-ddd/contents/domain/novel"
)

func (s *CommandHandler) CreateNovel(
	ctx context.Context, cmd command.CreateNovel) (
	id novel.ID, err error,
) {
	retry := -1
	id, err = backoff.RetryWith(
		ctx,
		func() (novel.ID, error) {
			retry += 1
			return s.createNovel(ctx, cmd)
		})

	if retry > 0 {
		s.stats.reqRetries[`CreateNovel`] += 1
		s.stats.reqRetryTimes[`CreateNovel`] += uint64(retry)
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
		r := s.novelFact.
			WithQuery(tx.Novel()).
			WithEventRepo(tx.Event()).
			Create(
				ctx,
				cmd.AuthorID,
				cmd.Title,
				cmd.Desc,
				cmd.Category,
				cmd.Tags)

		if r.IsError() {
			return r.Error()
		}

		if err = s.repo.Novel().Save(ctx, r.MustGet()); err != nil {
			return err
		}

		id = r.MustGet().ID()
		return nil
	})

	return id, err
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

		err = novel.UploadNewChapter(ctx, c.Title(), c.WordCount())
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
	err := backoff.Retry(ctx, func() error {
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

		var curTags []string
		for _, t := range novel.Tags() {
			curTags = append(curTags, string(t))
		}

		tobeAdd, tobeRemove := diffStrings(cmd.Tags, curTags)

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

		return nil
	}
}

func diffStrings(a, b []string) (onlyA, onlyB []string) {
	aSet := make(map[string]struct{})
	bSet := make(map[string]struct{})

	// 将切片元素存入 map
	for _, item := range a {
		aSet[item] = struct{}{}
	}

	for _, item := range b {
		bSet[item] = struct{}{}
	}

	for e := range aSet {
		if _, ok := bSet[e]; !ok {
			onlyA = append(onlyA, e)
		}
	}

	for e := range bSet {
		if _, ok := aSet[e]; !ok {
			onlyB = append(onlyB, e)
		}
	}

	return onlyA, onlyB
}
