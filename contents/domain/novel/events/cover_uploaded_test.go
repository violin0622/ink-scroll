package events_test

import (
	"encoding/json"

	"github.com/oklog/ulid/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	ene "github.com/violin0622/ink-scroll/contents/domain/novel/events"
	"github.com/violin0622/ink-scroll/contents/domain/novel/vo"
)

var _ = Describe(`cover uploaded event`, func() {
	It(`marshals the cover payload`, func() {
		cover, err := vo.NewCover(`https://example.com/cover.png`)
		Ω(err).ShouldNot(HaveOccurred())

		e := ene.NewCoverUploaded(ulid.Make(), cover)
		raw, err := e.Payload()
		Ω(err).ShouldNot(HaveOccurred())

		var got ene.CoverUploaded
		Ω(json.Unmarshal(raw, &got)).Should(Succeed())
		Ω(got.Cover).Should(Equal(cover))
	})
})

var _ = Describe(`novel ban events`, func() {
	It(`marshals empty payloads`, func() {
		banned := ene.NewNovelBanned(ulid.Make())
		raw, err := banned.Payload()
		Ω(err).ShouldNot(HaveOccurred())
		Ω(string(raw)).Should(Equal(`{}`))

		unbanned := ene.NewNovelUnbanned(ulid.Make())
		raw, err = unbanned.Payload()
		Ω(err).ShouldNot(HaveOccurred())
		Ω(string(raw)).Should(Equal(`{}`))
	})
})
