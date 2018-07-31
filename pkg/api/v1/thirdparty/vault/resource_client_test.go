package vault_test

import (
	"github.com/hashicorp/vault/api"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/solo-io/solo-kit/pkg/api/v1/thirdparty"
	. "github.com/solo-io/solo-kit/pkg/api/v1/thirdparty/vault"
	"github.com/solo-io/solo-kit/test/helpers"
)

var _ = Describe("Base", func() {
	var (
		vault              *api.Client
		rootKey            string
		artifacts, secrets thirdparty.ThirdPartyResourceClient
	)
	BeforeEach(func() {
		rootKey = "/secret/" + helpers.RandString(4)
		cfg := api.DefaultConfig()
		cfg.Address = "http://127.0.0.1:8200"
		c, err := api.NewClient(cfg)
		c.SetToken(vaultInstance.Token())
		Expect(err).NotTo(HaveOccurred())
		vault = c
		artifacts = NewThirdPartyResourceClient(vault, rootKey, &thirdparty.Artifact{})
		secrets = NewThirdPartyResourceClient(vault, rootKey, &thirdparty.Secret{})
	})
	AfterEach(func() {
		vault.Logical().Delete(rootKey)
	})
	It("CRUDs secrets", func() {
		helpers.TestThirdPartyClient("", secrets, &thirdparty.Secret{})
	})
	It("CRUDs artifacts", func() {
		helpers.TestThirdPartyClient("", artifacts, &thirdparty.Artifact{})
	})
})
