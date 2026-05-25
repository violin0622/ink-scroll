package outbox_test

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"play-ddd/contents/infra/outbox"
	ulidpb "play-ddd/proto/gen/go/ulid"
)

var _ ulidpb.ULID

var _ = Describe(`Relay`, func() {
	eb := &mockEventBus{}
	er := newRepo()
	relay := outbox.NewRelay(eb, er)

	It(`Start Ideoponent`, func() {
		By(`start first time`, func() {
			Expect(relay.Start()).Should(Succeed())
			Expect(relay.Ready()).Should(BeTrue())
		})
		By(`start second time`, func() {
			Expect(relay.Start()).Should(Succeed())
			Expect(relay.Ready()).Should(BeTrue())
		})
		By(`stop`, func() {
			Expect(relay.Stop()).Should(Succeed())
			Expect(relay.Ready()).Should(BeFalse())
		})
	})
})

var _ = Describe(`Relay`, func() {
	DescribeTable(`supports different max-pub and max-fetch like`,
		func(pub, fetch int) {
			eb := &mockEventBus{}
			er := newRepo()
			relay := outbox.NewRelay(eb, er,
				outbox.WithLogr(GinkgoLogr),
				outbox.WithInterval(time.Millisecond),
				outbox.WithMaxPub(pub),
				outbox.WithMaxFetch(fetch))

			Expect(relay.Start()).Should(Succeed())
			defer func() { Expect(relay.Stop()).Should(Succeed()) }()

			Eventually(func() int {
				return eb.published
			}).
				WithTimeout(time.Second).
				WithPolling(10 * time.Millisecond).
				Should(Equal(len(events)))
		},
		Entry(`pub=1 fetch=1`, 1, 1),
		Entry(`pub=2 fetch=3`, 2, 3),
		Entry(`pub=3 fetch=2`, 3, 2),
		Entry(`pub=5 fetch=5`, 5, 5),
		Entry(`pub=7 fetch=5`, 7, 5),
		Entry(`pub=5 fetch=7`, 5, 7),
	)
})

var _ = PDescribe(`Relay`, func() {
	Describe(
		`mark all subsequent events in same aggregate as blocked when pub failed`,
		func() {},
	)
})
