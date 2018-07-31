package consul_test

import (
	"github.com/hashicorp/consul/api"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/solo-io/solo-kit/pkg/api/v1/thirdparty"
	. "github.com/solo-io/solo-kit/pkg/api/v1/thirdparty/consul"
	"github.com/solo-io/solo-kit/test/helpers"
)

var _ = Describe("Base", func() {
	var (
		consul             *api.Client
		rootKey            string
		artifacts, secrets thirdparty.ThirdPartyResourceClient
	)
	BeforeEach(func() {
		rootKey = helpers.RandString(4)
		c, err := api.NewClient(api.DefaultConfig())
		Expect(err).NotTo(HaveOccurred())
		consul = c
		artifacts = NewThirdPartyResourceClient(consul, rootKey, &thirdparty.Artifact{})
		secrets = NewThirdPartyResourceClient(consul, rootKey, &thirdparty.Secret{})
	})
	AfterEach(func() {
		consul.KV().DeleteTree(rootKey, nil)
	})
	It("CRUDs secrets", func() {
		helpers.TestThirdPartyClient("", secrets, &thirdparty.Secret{})
	})
	It("CRUDs artifacts", func() {
		helpers.TestThirdPartyClient("", artifacts, &thirdparty.Artifact{})
	})
})
