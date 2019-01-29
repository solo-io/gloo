package secret_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/solo-io/gloo/projects/gloo/cli/pkg/helpers"
	pluginutils "github.com/solo-io/gloo/projects/gloo/pkg/plugins/utils"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-projects/projects/gloo/cli/pkg/testutils"
	extauthpb "github.com/solo-io/solo-projects/projects/gloo/pkg/api/v1/plugins/extauth"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/extauth"
)

var _ = Describe("ExtauthOauth", func() {

	BeforeEach(func() {
		helpers.UseMemoryClients()
	})

	It("should create secret", func() {
		err := testutils.GlooctlEE("create secret oauth --name oauth --namespace gloo-system --client-secret 123")
		Expect(err).NotTo(HaveOccurred())

		secret, err := helpers.MustSecretClient().Read("gloo-system", "oauth", clients.ReadOpts{})
		Expect(err).NotTo(HaveOccurred())

		var extension extauthpb.OauthSecret
		err = pluginutils.ExtensionToProto(secret.GetExtension(), extauth.ExtensionName, &extension)
		Expect(err).NotTo(HaveOccurred())

		Expect(extension).To(Equal(extauthpb.OauthSecret{ClientSecret: "123"}))
	})

	It("should error when no client secret provided", func() {
		err := testutils.GlooctlEE("create secret oauth --name oauth --namespace gloo-system")
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("client-secret not provided"))
	})

})
