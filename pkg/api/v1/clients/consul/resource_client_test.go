package consul_test

import (
	"time"

	"github.com/hashicorp/consul/api"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/solo-io/solo-kit/pkg/api/v1/clients/consul"
	"github.com/solo-io/solo-kit/test/helpers"
	"github.com/solo-io/solo-kit/test/mocks"
	"github.com/solo-io/solo-kit/test/tests/generic"
)

var _ = Describe("Base", func() {
	var (
		consul  *api.Client
		client  *ResourceClient
		rootKey string
	)
	BeforeEach(func() {
		rootKey = helpers.RandString(4)
		c, err := api.NewClient(api.DefaultConfig())
		Expect(err).NotTo(HaveOccurred())
		consul = c
		client = NewResourceClient(consul, rootKey, &mocks.MockResource{})
	})
	AfterEach(func() {
		consul.KV().DeleteTree(rootKey, nil)
	})
	It("CRUDs resources", func() {
		generic.TestCrudClient("", client, time.Minute)
	})
})
