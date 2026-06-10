package events_test

import (
	"github.com/oklog/ulid/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	ne "play-ddd/contents/domain/novel/events"
)

var _ = Describe(`publish review events`, func() {
	It(`emits empty payloads with stable kinds`, func() {
		submitted := ne.NewPublishReviewSubmitted(ulid.Make())
		raw, err := submitted.Payload()
		Ω(err).ShouldNot(HaveOccurred())
		Ω(string(raw)).Should(Equal(`{}`))
		Ω(submitted.Kind()).Should(Equal(`PublishReviewSubmitted`))

		withdrawed := ne.NewWithdrawed(ulid.Make())
		raw, err = withdrawed.Payload()
		Ω(err).ShouldNot(HaveOccurred())
		Ω(string(raw)).Should(Equal(`{}`))
		Ω(withdrawed.Kind()).Should(Equal(`PublishReviewWithdrawed`))
	})
})
