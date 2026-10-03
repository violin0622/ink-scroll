package events_test

import (
	"github.com/oklog/ulid/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/violin0622/ink-scroll/contents/domain/novel/events"
	"github.com/violin0622/ink-scroll/contents/domain/novel/vo"
)

var _ = Describe(`desc updated event`, func() {
	e := events.NewDescUpdated(
		ulid.Make(),
		vo.Description(`updated description`),
	)

	It(`marshal to json`, func() {
		raw, err := e.Payload()
		Ω(err).ShouldNot(HaveOccurred())
		Ω(string(raw)).Should(Equal(`{"desc":"updated description"}`))
	})

	// It(`unmarshal from json`, func() {
	// 	ev, err := events.RestoreEvent(
	// 		e.ID(),
	// 		e.AggID(),
	// 		e.EmittedAt(),
	// 		e.Kind(),
	// 		`Novel`,
	// 		[]byte(`{"desc":"updated description"}`),
	// 	)
	// 	Ω(err).ShouldNot(HaveOccurred())
	// 	Ω(ev.String()).Should(Equal(e.String()))
	// })

	// It(`not unmarshal id/aid/at`, func() {
	// 	id2, aid2, at2 := ulid.Make(), ulid.Make(), time.Now()
	// 	raw := fmt.Appendf(nil, `{
	// 		"id":"%s",
	// 		"aid":"%s",
	// 		"at":"%s",
	// 		"desc":"another description"
	// 	}`, id2, aid2, at2)

	// 	ev, err := events.RestoreEvent(
	// 		e.ID(),
	// 		e.AggID(),
	// 		e.EmittedAt(),
	// 		e.Kind(),
	// 		e.AggKind(),
	// 		raw,
	// 	)
	// 	Ω(err).ShouldNot(HaveOccurred())
	// 	Ω(ev.ID()).Should(Equal(e.ID()))
	// 	Ω(ev.AggID()).Should(Equal(e.AggID()))
	// 	Ω(ev.EmittedAt()).Should(Equal(e.EmittedAt()))
	// 	Ω(
	// 		ev.(events.DescUpdated).Payload(),
	// 	).Should(Equal([]byte(`{"desc":"another description"}`)))
	// })
})
