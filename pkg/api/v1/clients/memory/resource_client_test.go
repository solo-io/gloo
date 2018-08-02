package memory_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/solo-io/solo-kit/pkg/api/v1/clients/memory"
	"github.com/solo-io/solo-kit/test/mocks"
	"github.com/solo-io/solo-kit/test/tests"
)

var _ = Describe("Base", func() {
	var (
		client *ResourceClient
	)
	BeforeEach(func() {
		client = NewResourceClient(NewInMemoryResourceCache(), &mocks.MockData{})
	})
	AfterEach(func() {
	})
	It("CRUDs resources", func() {
		tests.TestCrudClient("", client)
	})
})
