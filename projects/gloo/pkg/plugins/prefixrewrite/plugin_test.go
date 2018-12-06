package prefixrewrite_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	envoyroute "github.com/envoyproxy/go-control-plane/envoy/api/v2/route"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/plugins/transformation"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins"
	. "github.com/solo-io/gloo/projects/gloo/pkg/plugins/prefixrewrite"
)

var _ = Describe("Plugin", func() {
	It("works", func() {

		p := NewPlugin()
		routeAction := &envoyroute.RouteAction{
			PrefixRewrite: "/",
		}
		out := &envoyroute.Route{
			Action: &envoyroute.Route_Route{
				Route: routeAction,
			},
		}
		err := p.ProcessRoute(plugins.Params{}, &v1.Route{
			RoutePlugins: &v1.RoutePlugins{
				PrefixRewrite: &transformation.PrefixRewrite{
					PrefixRewrite: "/foo",
				},
			},
		}, out)
		Expect(err).NotTo(HaveOccurred())
		Expect(routeAction.PrefixRewrite).To(Equal("/foo"))
	})
})
