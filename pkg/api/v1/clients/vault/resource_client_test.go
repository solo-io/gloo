package vault_test

import (
	"time"

	"github.com/hashicorp/vault/api"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	. "github.com/solo-io/solo-kit/pkg/api/v1/clients/vault"
	"github.com/solo-io/solo-kit/test/helpers"
	"github.com/solo-io/solo-kit/test/mocks"
	"github.com/solo-io/solo-kit/test/tests/generic"
)

var _ = Describe("Base", func() {
	var (
		vault   *api.Client
		rootKey string
		secrets clients.ResourceClient
	)
	BeforeEach(func() {
		rootKey = "/secret/" + helpers.RandString(4)
		cfg := api.DefaultConfig()
		cfg.Address = "http://127.0.0.1:8200"
		c, err := api.NewClient(cfg)
		c.SetToken(vaultInstance.Token())
		Expect(err).NotTo(HaveOccurred())
		vault = c
		secrets = NewResourceClient(vault, rootKey, &mocks.MockResource{})
	})
	AfterEach(func() {
		vault.Logical().Delete(rootKey)
	})
	It("CRUDs secrets", func() {
		generic.TestCrudClient("", secrets, time.Second/8)
	})
})
