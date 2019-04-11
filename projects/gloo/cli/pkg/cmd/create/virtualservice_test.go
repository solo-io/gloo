package create_test

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	"github.com/solo-io/gloo/pkg/cliutil/testutil"
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

	getOIDCConfig := func() *extauthpb.OAuth {

		vs, err := helpers.MustVirtualServiceClient().Read("gloo-system", "vs1", clients.ReadOpts{})
		Expect(err).NotTo(HaveOccurred())
		Expect(vs.Metadata.Name).To(Equal("vs1"))

		var extension extauthpb.VhostExtension
		err = pluginutils.UnmarshalExtension(vs.GetVirtualHost().GetVirtualHostPlugins(), extauth.ExtensionName, &extension)
		Expect(err).NotTo(HaveOccurred())

		return extension.AuthConfig.(*extauthpb.VhostExtension_Oauth).Oauth
	}

	DescribeTable("should create vhost",
		func(cmd string, expected extauthpb.OAuth) {
			err := testutils.GlooctlEE(cmd)
			Expect(err).NotTo(HaveOccurred())
			oidc := getOIDCConfig()
			Expect(*oidc).To(Equal(expected))
		},
		Entry("with oid config", "create vs --name vs1 --enable-oidc-auth --oidc-auth-client-id "+
			"1 --oidc-auth-app-url http://app.example.com --oidc-auth-client-secret-name fake "+
			"--oidc-auth-client-secret-namespace fakens --oidc-auth-issuer-url http://issuer.example.com "+
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
			"1 --oidc-auth-app-url http://app.example.com --oidc-auth-client-secret-name fake "+
			"--oidc-auth-client-secret-namespace fakens --oidc-auth-issuer-url http://issuer.example.com",
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

	Context("Interactive tests", func() {

		It("should create vs with no rate limits and auth", func() {
			testutil.ExpectInteractive(func(c *testutil.Console) {
				c.ExpectString("Add another domain for this virtual service")
				c.SendLine("")
				c.ExpectString("do you wish to add rate limiting to the virtual service")
				c.SendLine("n")
				c.ExpectString("do you wish to add oidc auth to the virtual service")
				c.SendLine("n")
				c.ExpectString("Use default namespace (gloo-system)?")
				c.SendLine("")
				c.ExpectString("name of the resource:")
				c.SendLine("default")
				c.ExpectEOF()
			}, func() {
				err := testutils.GlooctlEE("create vs -i")
				Expect(err).NotTo(HaveOccurred())
				_, err = helpers.MustVirtualServiceClient().Read("gloo-system", "default", clients.ReadOpts{})
				Expect(err).NotTo(HaveOccurred())
			})
		})

		It("should create vs with auth", func() {
			testutil.ExpectInteractive(func(c *testutil.Console) {
				c.ExpectString("Add another domain for this virtual service")
				c.SendLine("")
				c.ExpectString("do you wish to add rate limiting to the virtual service")
				c.SendLine("n")
				c.ExpectString("do you wish to add oidc auth to the virtual service")
				c.SendLine("y")
				c.ExpectString("What is your app url?")
				c.SendLine("http://app.example.com")
				c.ExpectString("What is your issuer url?")
				c.SendLine("https://accounts.google.com")
				c.ExpectString("What path (relative to your app url) should we use as a callback from the issuer?")
				c.SendLine("/auth-callback")
				c.ExpectString("What is your client id?")
				c.SendLine("me")
				c.ExpectString("What is your client secret name?")
				c.SendLine("secret-name")
				c.ExpectString("What is your client secret namespace?")
				c.SendLine("gloo-system")

				c.ExpectString("Use default namespace (gloo-system)?")
				c.SendLine("")
				c.ExpectString("name of the resource:")
				c.SendLine("vs1")
				c.ExpectEOF()
			}, func() {
				err := testutils.GlooctlEE("create vs -i")
				Expect(err).NotTo(HaveOccurred())

				oidc := getOIDCConfig()
				expected := extauthpb.OAuth{
					ClientId: "me",
					ClientSecretRef: &core.ResourceRef{
						Name:      "secret-name",
						Namespace: "gloo-system",
					},
					CallbackPath: "/auth-callback",
					IssuerUrl:    "https://accounts.google.com",
					AppUrl:       "http://app.example.com",
				}
				Expect(*oidc).To(Equal(expected))

			})
		})

	})

	var _ = Describe("dry-run", func() {
		It("can print as kube yaml in dry run", func() {
			out, err := testutils.GlooctlEEOut("create virtualservice kube --dry-run --name vs --domains foo.bar,baz.qux")
			Expect(err).NotTo(HaveOccurred())
			fmt.Print(out)
			Expect(out).To(Equal(`apiVersion: gateway.solo.io/v1
kind: VirtualService
metadata:
  creationTimestamp: null
  name: vs
  namespace: gloo-system
spec:
  displayName: vs
  virtualHost:
    domains:
    - foo.bar
    - baz.qux
status: {}
`))
		})
	})

})
