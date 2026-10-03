package novel_test

import (
	"net/url"
	"testing"
	"time"

	"github.com/go-logr/logr"
	"github.com/oklog/ulid/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gorm.io/gorm"

	"github.com/violin0622/ink-scroll/contents/domain/novel"
	"github.com/violin0622/ink-scroll/contents/domain/novel/vo"
	"github.com/violin0622/ink-scroll/contents/infra/repository/pg"
	novelrepo "github.com/violin0622/ink-scroll/contents/infra/repository/pg/novel"
)

func TestControllers(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecs(t, "Novel Suite")
}

var db *gorm.DB

var _ = BeforeSuite(func() {
	var err error
	db, err = pg.InitDB(pg.DSN)
	Ω(err).ShouldNot(HaveOccurred())
})

var _ = Describe(`NoveRepo`, func() {
	var (
		now      = time.Now().Truncate(time.Millisecond)
		authorID = ulid.Make()
		novelID  = ulid.Make()
		title    = `An Awesome Story`
		desc     = `The story of the brave one.`
		category = `fantacy`
		tags     = []string{`DnD`, `Sole Famale`, `Legacy`}
		status   = vo.Serial
		chapters = []vo.Chapter{{
			ID: ulid.Make(),
		}, {
			ID: ulid.Make(),
		}}
		wordCount = 30
		createdAt = now
		updatedAt = now

		fact = novel.NewFactory(logr.Discard())
	)

	It(`save novel`, func(ctx SpecContext) {
		repo := novelrepo.New(db, novel.NewFactory(logr.Discard()))
		novel, err := fact.Restore(
			novelID,
			authorID,
			title,
			category,
			desc,
			tags,
			int(status),
			url.URL{},
			nil,
			chapters,
			wordCount,
			updatedAt,
			createdAt,
		)

		Ω(err).ShouldNot(HaveOccurred())
		Ω(repo.Save(ctx, novel)).Should(Succeed())
	})

	It(`load novel`, func(ctx SpecContext) {
		repo := novelrepo.New(db, novel.Factory{})
		novel, err := repo.Get(ctx, novelID)
		Ω(err).ShouldNot(HaveOccurred())
		Ω(novel.ID()).Should(BeEquivalentTo(novelID))
		Ω(novel.Title()).Should(BeEquivalentTo(title))
		Ω(novel.AuthorID()).Should(BeEquivalentTo(authorID))
		Ω(novel.Description()).Should(BeEquivalentTo(desc))
		Ω(novel.Category()).Should(BeEquivalentTo(category))
		Ω(novel.Tags()).Should(And(HaveLen(len(tags))))
		Ω(novel.Status()).Should(BeEquivalentTo(vo.Serial))
		Ω(novel.ChapterCount()).Should(Equal(len(chapters)))
		Ω(novel.WordCount()).Should(Equal(wordCount))
		Ω(novel.UpdatedAt()).Should(BeEquivalentTo(updatedAt))
		Ω(novel.CreatedAt()).Should(BeEquivalentTo(createdAt))

		for i := range tags {
			Ω(novel.Tags()[i]).Should(BeEquivalentTo(tags[i]))
		}

		for i, c := range chapters {
			nc := novel.Chapters()[i]
			Ω(nc.ID).Should(Equal(c.ID))
		}
	})
})
