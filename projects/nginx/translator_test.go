package nginx_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/solo-io/solo-kit/projects/gateway/pkg/api/v1"
	gloo_solo_io1 "github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1"
	. "github.com/solo-io/solo-kit/projects/nginx"
)

var _ = Describe("Translator", func() {
	Context("when given a single server with a single empty location", func() {
		It("generates the appropriate Nginx config", func() {
			matcher := &gloo_solo_io1.Matcher{
				PathSpecifier: &gloo_solo_io1.Matcher_Prefix{
					Prefix: "/",
				},
			}
			action := &gloo_solo_io1.Route_RouteAction{
				RouteAction: &gloo_solo_io1.RouteAction{
					Destination: &gloo_solo_io1.RouteAction_Single{
						Single: &gloo_solo_io1.Destination{},
					},
				},
			}
			route := &gloo_solo_io1.Route{
				Matcher: matcher,
				Action:  action,
			}
			routes := []*gloo_solo_io1.Route{
				route,
			}
			virtualHost := &gloo_solo_io1.VirtualHost{
				Routes: routes,
			}
			virtualService := v1.VirtualService{
				VirtualHost: virtualHost,
			}
			virtualServices := []v1.VirtualService{
				virtualService,
			}
			gateway := &v1.Gateway{}

			server, err := Translate(gateway, virtualServices)
			Expect(err).NotTo(HaveOccurred())

			serverContext, err := GenerateServerContext(&server)
			Expect(err).NotTo(HaveOccurred())
			expected := `server {
    location / {
    }
}`
			Expect(string(serverContext)).To(Equal(expected))
		})
	})
})
