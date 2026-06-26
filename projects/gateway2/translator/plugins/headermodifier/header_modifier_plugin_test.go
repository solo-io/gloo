package headermodifier_test

import (
	"context"

	"github.com/golang/protobuf/ptypes/wrappers"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"google.golang.org/protobuf/proto"

	"github.com/solo-io/gloo/projects/gateway2/translator/plugins"
	"github.com/solo-io/gloo/projects/gateway2/translator/plugins/filtertests"
	"github.com/solo-io/gloo/projects/gateway2/translator/plugins/headermodifier"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/headers"
	"github.com/solo-io/solo-kit/pkg/api/external/envoy/api/v2/core"
	gwv1 "sigs.k8s.io/gateway-api/apis/v1"
)

var _ = DescribeTable(
	"HeaderModifierPlugin",
	func(
		plugin plugins.RoutePlugin,
		filter gwv1.HTTPRouteFilter,
		expectedRoute *v1.Route,
	) {
		filtertests.AssertExpectedRoute(
			plugin,
			expectedRoute,
			true,
			filter,
		)
	},
	Entry(
		"applies request header modifier filter",
		headermodifier.NewPlugin(),
		gwv1.HTTPRouteFilter{
			Type: gwv1.HTTPRouteFilterRequestHeaderModifier,
			RequestHeaderModifier: &gwv1.HTTPHeaderFilter{
				Add: []gwv1.HTTPHeader{
					{
						Name:  "foo",
						Value: "bar",
					},
				},
				Set: []gwv1.HTTPHeader{
					{
						Name:  "foo",
						Value: "bar",
					},
				},
				Remove: []string{"foo"},
			},
		},
		&v1.Route{
			Options: &v1.RouteOptions{
				HeaderManipulation: &headers.HeaderManipulation{
					RequestHeadersToAdd: []*core.HeaderValueOption{
						{
							HeaderOption: &core.HeaderValueOption_Header{
								Header: &core.HeaderValue{
									Key:   "foo",
									Value: "bar",
								},
							},
							Append: &wrappers.BoolValue{Value: true},
						},
						{
							HeaderOption: &core.HeaderValueOption_Header{
								Header: &core.HeaderValue{
									Key:   "foo",
									Value: "bar",
								},
							},
							Append: &wrappers.BoolValue{Value: false},
						},
					},
					RequestHeadersToRemove: []string{"foo"},
				},
			},
		},
	),
	Entry(
		"applies response header modifier filter",
		headermodifier.NewPlugin(),
		gwv1.HTTPRouteFilter{
			Type: gwv1.HTTPRouteFilterResponseHeaderModifier,
			ResponseHeaderModifier: &gwv1.HTTPHeaderFilter{
				Add: []gwv1.HTTPHeader{
					{
						Name:  "foo",
						Value: "bar",
					},
				},
				Set: []gwv1.HTTPHeader{
					{
						Name:  "foo",
						Value: "bar",
					},
				},
				Remove: []string{"foo"},
			},
		},
		&v1.Route{
			Options: &v1.RouteOptions{
				HeaderManipulation: &headers.HeaderManipulation{
					ResponseHeadersToAdd: []*headers.HeaderValueOption{
						{
							Header: &headers.HeaderValue{
								Key:   "foo",
								Value: "bar",
							},
							Append: &wrappers.BoolValue{Value: true},
						},
						{
							Header: &headers.HeaderValue{
								Key:   "foo",
								Value: "bar",
							},
							Append: &wrappers.BoolValue{Value: false},
						},
					},
					ResponseHeadersToRemove: []string{"foo"},
				},
			},
		},
	),
)

var _ = Describe("HeaderModifierPlugin mutation safety", func() {
	It("does not mutate a HeaderManipulation already present on the route options", func() {
		// A route's options sub-messages can be shared with every other route referencing the
		// same RouteOption (solo-io/solo-projects#8802). Parent filters re-applied to delegated
		// child routes reach this plugin with options already populated from RouteOptions, so it
		// must never write into the existing HeaderManipulation in place.
		shared := &headers.HeaderManipulation{
			RequestHeadersToAdd: []*core.HeaderValueOption{
				{
					HeaderOption: &core.HeaderValueOption_Header{
						Header: &core.HeaderValue{Key: "from-route-option", Value: "original"},
					},
					Append: &wrappers.BoolValue{Value: true},
				},
			},
			ResponseHeadersToRemove: []string{"x-strip-response"},
		}
		snapshot := proto.Clone(shared).(*headers.HeaderManipulation)

		outputRoute := &v1.Route{Options: &v1.RouteOptions{HeaderManipulation: shared}}
		rtCtx := &plugins.RouteContext{
			HTTPRoute: &gwv1.HTTPRoute{},
			Rule: &gwv1.HTTPRouteRule{
				Filters: []gwv1.HTTPRouteFilter{{
					Type: gwv1.HTTPRouteFilterRequestHeaderModifier,
					RequestHeaderModifier: &gwv1.HTTPHeaderFilter{
						Add: []gwv1.HTTPHeader{{Name: "foo", Value: "bar"}},
					},
				}},
			},
		}

		err := headermodifier.NewPlugin().ApplyRoutePlugin(context.Background(), rtCtx, outputRoute)
		Expect(err).NotTo(HaveOccurred())

		// the filter's headers land on the route...
		result := outputRoute.GetOptions().GetHeaderManipulation()
		Expect(result.GetRequestHeadersToAdd()).To(HaveLen(1))
		Expect(result.GetRequestHeadersToAdd()[0].GetHeader().GetKey()).To(Equal("foo"))
		// ...fields the filter does not touch survive from the original message...
		Expect(result.GetResponseHeadersToRemove()).To(ConsistOf("x-strip-response"))
		// ...and the original stays untouched.
		Expect(proto.Equal(shared, snapshot)).To(BeTrue(),
			"ApplyRoutePlugin mutated a HeaderManipulation shared with other routes")
	})

	It("does not mutate a HeaderManipulation already present on the route options (response filter)", func() {
		// applyResponseFilter is an independent write path from applyRequestFilter; it needs its
		// own pin so neither can regress to writing into the shared message in place.
		shared := &headers.HeaderManipulation{
			RequestHeadersToAdd: []*core.HeaderValueOption{
				{
					HeaderOption: &core.HeaderValueOption_Header{
						Header: &core.HeaderValue{Key: "from-route-option", Value: "original"},
					},
					Append: &wrappers.BoolValue{Value: true},
				},
			},
			ResponseHeadersToRemove: []string{"x-strip-response"},
		}
		snapshot := proto.Clone(shared).(*headers.HeaderManipulation)

		outputRoute := &v1.Route{Options: &v1.RouteOptions{HeaderManipulation: shared}}
		rtCtx := &plugins.RouteContext{
			HTTPRoute: &gwv1.HTTPRoute{},
			Rule: &gwv1.HTTPRouteRule{
				Filters: []gwv1.HTTPRouteFilter{{
					Type: gwv1.HTTPRouteFilterResponseHeaderModifier,
					ResponseHeaderModifier: &gwv1.HTTPHeaderFilter{
						Add: []gwv1.HTTPHeader{{Name: "foo", Value: "bar"}},
					},
				}},
			},
		}

		err := headermodifier.NewPlugin().ApplyRoutePlugin(context.Background(), rtCtx, outputRoute)
		Expect(err).NotTo(HaveOccurred())

		// the filter's headers land on the route...
		result := outputRoute.GetOptions().GetHeaderManipulation()
		Expect(result.GetResponseHeadersToAdd()).To(HaveLen(1))
		Expect(result.GetResponseHeadersToAdd()[0].GetHeader().GetKey()).To(Equal("foo"))
		// ...fields the filter does not touch survive from the original message (the request
		// side here: the response filter legitimately overwrites the response-side fields)...
		Expect(result.GetRequestHeadersToAdd()).To(HaveLen(1))
		Expect(result.GetRequestHeadersToAdd()[0].GetHeader().GetKey()).To(Equal("from-route-option"))
		// ...and the original stays untouched.
		Expect(proto.Equal(shared, snapshot)).To(BeTrue(),
			"ApplyRoutePlugin mutated a HeaderManipulation shared with other routes")
	})
})
