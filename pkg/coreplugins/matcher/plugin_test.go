package matcher

import (
	envoyroute "github.com/envoyproxy/go-control-plane/envoy/api/v2/route"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/solo-io/gloo-testing/helpers"
	. "github.com/solo-io/gloo/pkg/coreplugins/matcher"
)

var _ = Describe("Plugin", func() {
	Describe("ProcessRoute", func() {
		It("takes an event matcher and creates a route match for envoy", func() {
			plug := &Plugin{}
			route := NewTestRoute1()
			out := &envoyroute.Route{}
			err := plug.ProcessRoute(nil, route, out)
			Expect(err).NotTo(HaveOccurred())
			Expect(out.Match.PathSpecifier).To(Equal(&envoyroute.RouteMatch_Prefix{Prefix: "/foo"}))
		})
	})
})
