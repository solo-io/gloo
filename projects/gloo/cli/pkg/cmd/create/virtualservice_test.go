package create_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	"github.com/solo-io/gloo/projects/gloo/cli/pkg/helpers"
	pluginutils "github.com/solo-io/gloo/projects/gloo/pkg/plugins/utils"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-projects/projects/gloo/cli/pkg/testutils"
	extauthpb "github.com/solo-io/solo-projects/projects/gloo/pkg/api/v1/plugins/extauth"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/extauth"
)

var _ = Describe("Virtualservice", func() {
	BeforeEach(func() {
		helpers.UseMemoryClients()
	})

	DescribeTable("should create vhost",
		func(cmd string, expected extauthpb.OAuth) {
			err := testutils.GlooctlEE(cmd)
			Expect(err).NotTo(HaveOccurred())

			vs, err := helpers.MustVirtualServiceClient().Read("gloo-system", "vs1", clients.ReadOpts{})
			Expect(err).NotTo(HaveOccurred())
			Expect(vs.Metadata.Name).To(Equal("vs1"))

			var extension extauthpb.VhostExtension
			err = pluginutils.UnmarshalExtension(vs.GetVirtualHost().GetVirtualHostPlugins(), extauth.ExtensionName, &extension)
			Expect(err).NotTo(HaveOccurred())

			oidc := extension.AuthConfig.(*extauthpb.VhostExtension_Oauth).Oauth
			Expect(*oidc).To(Equal(expected))
		},
		Entry("with oid config", "create vs --name vs1 --enable-oidc-auth --oidc-auth-client-id "+
			"1 --oidc-auth-app-url http://app.example.com --oidc-auth-client-name fake "+
			"--oidc-auth-client-namespace fakens --oidc-auth-issuer-url http://issuer.example.com "+
			"--oidc-auth-callback-path /cb",
			extauthpb.OAuth{
				ClientId: "1",
				ClientSecretRef: &core.ResourceRef{
					Name:      "fake",
					Namespace: "fakens",
				},
				CallbackPath: "/cb",
				IssuerUrl:    "http://issuer.example.com",
				AppUrl:       "http://app.example.com",
			}),
		Entry("with default callback", "create vs --name vs1 --enable-oidc-auth --oidc-auth-client-id "+
			"1 --oidc-auth-app-url http://app.example.com --oidc-auth-client-name fake "+
			"--oidc-auth-client-namespace fakens --oidc-auth-issuer-url http://issuer.example.com",
			extauthpb.OAuth{
				ClientId: "1",
				ClientSecretRef: &core.ResourceRef{
					Name:      "fake",
					Namespace: "fakens",
				},
				CallbackPath: "/oidc-gloo-callback",
				IssuerUrl:    "http://issuer.example.com",
				AppUrl:       "http://app.example.com",
			}),
	)
})
