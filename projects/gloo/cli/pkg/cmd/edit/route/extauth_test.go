package route_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	"github.com/solo-io/gloo/pkg/cliutil/testutil"
	gatewayv1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/helpers"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins/utils"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-projects/projects/gloo/cli/pkg/testutils"
	extauthpb "github.com/solo-io/solo-projects/projects/gloo/pkg/api/v1/plugins/extauth"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/extauth"
)

var _ = Describe("Extauth", func() {
	var (
		vsvc     *gatewayv1.VirtualService
		vsClient gatewayv1.VirtualServiceClient
	)
	BeforeEach(func() {
		helpers.UseMemoryClients()
		// create a settings object
		vsClient = helpers.MustVirtualServiceClient()
		vsvc = &gatewayv1.VirtualService{
			Metadata: core.Metadata{
				Name:      "vs",
				Namespace: "gloo-system",
			},
			VirtualHost: &v1.VirtualHost{
				Name: "vs",
				Routes: []*v1.Route{{
					Matcher: &v1.Matcher{
						PathSpecifier: &v1.Matcher_Prefix{Prefix: "/"},
					}}, {
					Matcher: &v1.Matcher{
						PathSpecifier: &v1.Matcher_Prefix{Prefix: "/r"},
					}},
				},
			},
		}
		var err error
		vsvc, err = vsClient.Write(vsvc, clients.WriteOpts{})
		Expect(err).NotTo(HaveOccurred())
	})

	extAuthExtension := func(index int) *extauthpb.RouteExtension {
		var extAuthRouteExt extauthpb.RouteExtension
		var err error
		vsvc, err = vsClient.Read(vsvc.Metadata.Namespace, vsvc.Metadata.Name, clients.ReadOpts{})
		Expect(err).NotTo(HaveOccurred())

		err = utils.UnmarshalExtension(vsvc.VirtualHost.Routes[index].RoutePlugins, extauth.ExtensionName, &extAuthRouteExt)
		if err != nil {
			if err == utils.NotFoundError {
				return nil
			}
			Expect(err).NotTo(HaveOccurred())
		}
		return &extAuthRouteExt
	}
	Context("Non-interactive tests", func() {

		DescribeTable("should edit extauth config",
			func(cmd string, index int, expected *extauthpb.RouteExtension) {

				err := testutils.GlooctlEE(cmd)
				Expect(err).NotTo(HaveOccurred())

				extension := extAuthExtension(index)
				Expect(extension).To(Equal(expected))
			},
			Entry("edit route 0 doesnt impact route one", "edit route externalauth --name vs --namespace gloo-system --index 0 --disable=true",
				1,
				nil),
			Entry("edit route and disable it", "edit route externalauth --name vs --namespace gloo-system --index 1 --disable=true",
				1,
				&extauthpb.RouteExtension{
					Disable: true,
				}),
			Entry("edit route and un-disable it", "edit route externalauth --name vs --namespace gloo-system --index 1 --disable=false",
				1,
				&extauthpb.RouteExtension{
					Disable: false,
				}),
		)
	})

	Context("Interactive tests", func() {

		It("should enabled auth on route", func() {
			testutil.ExpectInteractive(func(c *testutil.Console) {
				c.ExpectString("Choose a Virtual Service:")
				c.SendLine("")
				c.ExpectString("Choose the route you wish to change:")
				c.PressDown()
				c.SendLine("")
				c.ExpectString("Disable auth on this route?")
				c.SendLine("")
				c.ExpectEOF()
			}, func() {
				err := testutils.GlooctlEE("edit route externalauth -i")
				Expect(err).NotTo(HaveOccurred())
				extension := extAuthExtension(1)
				Expect(extension).To(Equal(&extauthpb.RouteExtension{
					Disable: true,
				}))
			})
		})

	})

	Context("Errors", func() {
		It("should error with wrong resource version", func() {
			badrv := vsvc.Metadata.ResourceVersion + "a"
			err := testutils.GlooctlEE("edit route externalauth --name vs --namespace gloo-system --index 1 --disable=true --resource-version=" + badrv)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("conflict - resource version does not match"))
		})

	})
})
