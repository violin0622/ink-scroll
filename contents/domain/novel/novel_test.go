package novel_test

import (
	"context"
	"time"

	"github.com/go-logr/logr"
	"github.com/oklog/ulid/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"play-ddd/common"
	"play-ddd/contents/domain/novel"
	"play-ddd/contents/domain/novel/events"
	"play-ddd/contents/domain/novel/vo"
)

var _ novel.EventRepo = (*nopEventRepo)(nil)

type nopEventRepo struct{}

func (nopEventRepo) Append(context.Context, ...novel.Event) error { return nil }
func (nopEventRepo) Fetch(context.Context, novel.ID) ([]novel.Event, error) {
	return nil, nil
}

var _ = Describe(`NovelReplayEvents`, func() {
	var (
		n            novel.Novel
		err          error
		id, autherID = ulid.Make(), ulid.Make()
		cover, _     = vo.NewCover(`https://example.com/cover.png`)
		fac          = novel.NewFactory(logr.Discard()).WithEventRepo(nopEventRepo{})
	)

	// events
	var (
		novelCreated = events.NewNovelCreated(
			id,
			autherID,
			vo.Title(`测试小说`),
			vo.NoCover,
			vo.Description(`初始描述`),
			vo.Category(`玄幻`),
			[]vo.Tag{vo.Tag(`修仙`), vo.Tag(`穿越`)})

		reviewSubmitted  = events.NewPublishReviewSubmitted(id)
		reviewWithdrawed = events.NewWithdrawed(id)
		coverUploaded    = events.NewCoverUploaded(id, cover)
		descUpdated      = events.NewDescUpdated(id, vo.Description(`更新后的描述`))
		tagsUpdated      = events.NewTagsUpdated(id,
			[]vo.Tag{vo.Tag(`修仙`), vo.Tag(`穿越`)},
			[]vo.Tag{vo.Tag(`修仙`), vo.Tag(`穿越`), vo.Tag(`系统`)})
		chapterUploaded1 = events.NewNewChapterUploaded(id, vo.Chapter{
			ID: ulid.Make(),
		}, 1500, time.Now().Add(-time.Minute))
		chapterUploaded2 = events.NewNewChapterUploaded(id, vo.Chapter{
			ID: ulid.Make(),
		}, 1500, time.Now())
		chapterRevised = events.NewChapterRevised(id, vo.Chapter{
			ID: ulid.Make(),
		})
		novelBanned    = events.NewNovelBanned(id)
		novelUnbanned  = events.NewNovelUnbanned(id)
		nolongerUpdate = events.NewNolongerUpdate(id)
		completed      = events.NewCompleted(id)
	)
	eventsHist := []events.Event{
		novelCreated,
		reviewSubmitted,
		reviewWithdrawed,
		coverUploaded,
		descUpdated,
		tagsUpdated,
		chapterUploaded1,
		chapterUploaded2,
		nolongerUpdate,
		completed,
	}

	BeforeEach(func() {
		n = fac.New()
	})

	When(`first event is NovelCreated`, func() {
		It(`should accept`, func() {
			err = n.ReplayEvents(novelCreated)
			Ω(err).ShouldNot(HaveOccurred())
		})
	})

	When(`first event is not NovelCreated`, func() {
		DescribeTable(`should reject invalid initial events`,
			func(e novel.Event) {
				err = n.ReplayEvents(e)
				Ω(err).Should(MatchError(common.ErrInitialEvent))
			},
			Entry(`PublishReviewSubmitted`, reviewSubmitted),
			Entry(`PublishReviewWithdrawed`, reviewWithdrawed),
			Entry(`DescUpdated`, descUpdated),
			Entry(`TagsUpdated`, tagsUpdated),
			Entry(`CoverUploaded`, coverUploaded),
			Entry(`NovelBanned`, novelBanned),
			Entry(`NovelUnbanned`, novelUnbanned),
			Entry(`NewChapterUploaded`, chapterUploaded2),
			Entry(`ChapterRevised`, chapterRevised),
			Entry(`Completed`, completed),
			Entry(`NolongerUpdate`, nolongerUpdate),
		)
	})

	When(`replaying multiple events starting with NovelCreated`, func() {
		BeforeEach(func() {
			err = n.ReplayEvents(eventsHist...)
			Ω(err).ShouldNot(HaveOccurred())
		})

		It(`should maintain correct basic properties`, func() {
			Ω(n.ID()).Should(Equal(id))
			Ω(n.AuthorID()).Should(Equal(autherID))
			Ω(n.Title()).Should(Equal(vo.Title(`测试小说`)))
			Ω(n.Description()).Should(Equal(vo.Description(`更新后的描述`)))
			Ω(n.Category()).Should(Equal(vo.Category(`玄幻`)))
			Ω(n.Cover()).Should(Equal(cover))
		})

		It(`should maintain correct status and tags`, func() {
			Ω(n.Status()).Should(Equal(vo.Completed))
			expectedTags := []vo.Tag{vo.Tag(`修仙`), vo.Tag(`穿越`), vo.Tag(`系统`)}
			Ω(n.Tags()).Should(Equal(expectedTags))
		})

		// It(`should maintain correct chapter information`, func() {
		// 	chapters := n.Chapters()
		// 	Ω(n.ChapterCount()).Should(Equal(2))
		// 	Ω(chapters[0].Title).Should(Equal("第一章（修订版）"))
		// 	Ω(chapters[0].Sequence).Should(Equal(1))
		// 	Ω(chapters[0].WordCount).Should(Equal(1200))
		// 	Ω(chapters[1].Title).Should(Equal("第二章"))
		// 	Ω(chapters[1].Sequence).Should(Equal(2))
		// 	Ω(chapters[1].WordCount).Should(Equal(1500))
		// })

		It(`should maintain correct word count and timestamps`, func() {
			Ω(n.WordCount()).Should(Equal(3000))
			Ω(n.CreatedAt()).ShouldNot(BeZero())
			Ω(n.UpdatedAt()).ShouldNot(BeZero())
		})
	})
	It(`uploads cover and appends an event`, func() {
		err = n.ReplayEvents(novelCreated)
		Ω(err).ShouldNot(HaveOccurred())

		uploadedCover, err := vo.NewCover(`https://example.com/uploaded-cover.png`)
		Ω(err).ShouldNot(HaveOccurred())

		err = n.UploadCover(context.Background(), uploadedCover)
		Ω(err).ShouldNot(HaveOccurred())
		Ω(n.Cover()).Should(Equal(uploadedCover))
	})
	It(`submits and withdraws publish review`, func() {
		err = n.ReplayEvents(novelCreated)
		Ω(err).ShouldNot(HaveOccurred())

		err = n.ApplyPublish(context.Background())
		Ω(err).ShouldNot(HaveOccurred())
		Ω(n.ReleaseStatus()).Should(Equal(vo.Reviewing))

		err = n.WithdrawRelease(context.Background())
		Ω(err).ShouldNot(HaveOccurred())
		Ω(n.ReleaseStatus()).Should(Equal(vo.Draft))
	})
})

// var _ = Describe(``, func() {
// 	var _ novel.EventRepo = fake.New[novel.ID, novel.ID]()
// })
