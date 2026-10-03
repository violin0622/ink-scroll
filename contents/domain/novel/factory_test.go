package novel_test

import (
	"context"

	"github.com/go-logr/logr"
	"github.com/oklog/ulid/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/violin0622/ink-scroll/contents/domain/novel"
	"github.com/violin0622/ink-scroll/contents/domain/novel/vo"
)

var _ = Describe(`novel factory`, func() {
	var n novel.Novel
	var f novel.Factory
	var (
		_, authorID = ulid.Make(), ulid.Make()
		title       = `An awesome novel`
		desc        = `从前有座山山上有座庙庙里有个老和尚...`
		category    = `志异`
		tags        = []string{}
	)

	BeforeEach(func() {
		f = novel.NewFactory(logr.Discard()).
			WithEventRepo(nopEventRepo{})

		var err error
		n, err = f.Create(context.TODO(), authorID, title, desc, category, tags)
		Expect(err).NotTo(HaveOccurred())
	})

	It(`新建作品发布状态为草稿`, func() {
		Expect(n.ReleaseStatus()).To(Equal(vo.ReleaseStatus(vo.Draft)))
	})

	It(`新建作品连载状态为连载`, func() {
		Expect(n.Status()).To(Equal(vo.Serial))
	})
})
