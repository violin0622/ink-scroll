package event_test

import (
	"testing"

	"github.com/oklog/ulid/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/anypb"

	ne "github.com/violin0622/ink-scroll/contents/domain/novel/events"
	"github.com/violin0622/ink-scroll/contents/domain/novel/vo"
	"github.com/violin0622/ink-scroll/contents/infra/repository/pg/event"
	novelv1 "github.com/violin0622/ink-scroll/proto/gen/go/contents/novel/v1"
	ulidpb "github.com/violin0622/ink-scroll/proto/gen/go/ulid"
)

func TestControllers(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecs(t, "Event Suite")
}

// var db *gorm.DB

// var _ = BeforeSuite(func() {
// 	var err error
// 	db, err = pg.InitDB(pg.DSN)
// 	Ω(err).ShouldNot(HaveOccurred())
// })

var _ = Describe(`Event`, func() {
	var (
		id, authID = ulid.Make(), ulid.Make()
		desc       = vo.Description(`A New Description.`)
		title      = vo.Title(`My Title`)
		category   = vo.Category(`Story`)
		tags       = []vo.Tag{`Hot`}
		strTags    = []string{`Hot`}
	)

	When(`convert from domain`, func() {
		It(`can convert from novel.Event`, func() {
			novelCreatedEvent := ne.NewNovelCreated(
				id,
				authID,
				title,
				vo.Cover{},
				desc,
				category,
				tags,
			)
			a, _ := anypb.New(&novelv1.NovelCreatedPayload{
				AuthorId:    ulidpb.From(authID),
				Title:       string(title),
				Description: string(desc),
				Tags:        strTags,
				Category:    string(category),
			})
			expectPayload, _ := protojson.Marshal(a)

			var event event.Event
			err := event.FromDomain(novelCreatedEvent)

			Ω(err).ShouldNot(HaveOccurred())
			Ω(event.Payload).Should(BeEquivalentTo(expectPayload))
		})

		It(`can convert from empty payload`, func() {
			var event event.Event
			nbe := ne.NewNovelBanned(id)
			err := event.FromDomain(nbe)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(event.Payload).Should(BeEquivalentTo(`{}`))
		})
	})
})
