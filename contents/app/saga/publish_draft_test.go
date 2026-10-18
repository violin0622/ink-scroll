package saga_test

import (
	"context"

	"github.com/oklog/ulid/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/violin0622/ink-scroll/contents/app/saga"
)

var _ = Describe(`Publish draft`, func() {
	var (
		id  = ulid.Make()
		dr  = nil
		pds = saga.NewSAGA(id, dr)
	)

	It(`ends in terminate status`, func() {
		err := saga.DoSAGA(context.TODO(), id, repo)
		Expect(err).ShouldNot(HaveOccurred())
	})
})
