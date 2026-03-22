package outbox_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gleak"
)

var _ = AfterEach(func() {
	Eventually(Goroutines).ShouldNot(HaveLeaked())
})

func TestControllers(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecs(t, "Novel Suite")
}
