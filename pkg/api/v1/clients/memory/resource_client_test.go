package memory_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/solo-io/solo-kit/pkg/api/v1/clients/memory"
	"github.com/solo-io/solo-kit/test/helpers"
	"github.com/solo-io/solo-kit/test/mocks"
)

var _ = Describe("Base", func() {
	var (
		client *ResourceClient
	)
	BeforeEach(func() {
		client = NewResourceClient(&mocks.MockResource{})
	})
	AfterEach(func() {
	})
	It("CRUDs resources", func() {
		helpers.TestCrudClient("", client)
	})
})
