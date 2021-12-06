package graphql_test

import (
	"github.com/golang/protobuf/ptypes"
	v2 "github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/extensions/filters/http/graphql/v2"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/graphql/v1alpha1"
	"github.com/solo-io/gloo/projects/gloo/pkg/utils"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/graphql"

	envoy_config_route_v3 "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
)

var _ = Describe("Graphql plugin", func() {
	var (
		plugin        *graphql.Plugin
		params        plugins.Params
		vhostParams   plugins.VirtualHostParams
		virtualHost   *v1.VirtualHost
		route         *v1.Route
		routeAction   *v1.Route_GraphqlSchemaRef
		httpListener  *v1.HttpListener
		gqlSchemaSpec *v1alpha1.GraphQLSchema
	)

	BeforeEach(func() {
		routeAction = &v1.Route_GraphqlSchemaRef{
			GraphqlSchemaRef: &core.ResourceRef{
				Name:      "gql",
				Namespace: "gloo-system",
			},
		}
		route = &v1.Route{
			Action: routeAction,
		}

		gqlSchemaSpec = &v1alpha1.GraphQLSchema{
			Metadata: &core.Metadata{
				Name:      "gql",
				Namespace: "gloo-system",
			},
			Schema:              "",    // customized later
			EnableIntrospection: false, // customized later
			Resolutions:         nil,   // customized later
		}
	})

	JustBeforeEach(func() {
		virtualHost = &v1.VirtualHost{
			Name:    "virt1",
			Domains: []string{"*"},
			Routes:  []*v1.Route{route},
		}

		httpListener = &v1.HttpListener{
			VirtualHosts: []*v1.VirtualHost{virtualHost},
		}
		proxy := &v1.Proxy{
			Metadata: &core.Metadata{
				Name:      "http-proxy",
				Namespace: "default",
			},
			Listeners: []*v1.Listener{{
				Name: "default",
				ListenerType: &v1.Listener_HttpListener{
					HttpListener: httpListener,
				},
			}},
		}

		params.Snapshot = &v1.ApiSnapshot{
			Proxies: v1.ProxyList{proxy},
			GraphqlSchemas: v1alpha1.GraphQLSchemaList{
				gqlSchemaSpec,
			},
			Upstreams: v1.UpstreamList{
				{
					Metadata: &core.Metadata{
						Name:      "us",
						Namespace: "gloo-system",
					},
				},
			},
		}
		vhostParams = plugins.VirtualHostParams{
			Params:   params,
			Proxy:    proxy,
			Listener: proxy.Listeners[0],
		}
	})

	Context("process snapshot", func() {
		var (
			outRoute   envoy_config_route_v3.Route
			outFilters []plugins.StagedHttpFilter
		)

		var translateRoute = func() *v2.GraphQLRouteConfig {
			goTpfc := outRoute.TypedPerFilterConfig[graphql.FilterName]
			Expect(goTpfc).NotTo(BeNil())
			var perRouteGql v2.GraphQLRouteConfig
			err := ptypes.UnmarshalAny(goTpfc, &perRouteGql)
			Expect(err).NotTo(HaveOccurred())
			return &perRouteGql
		}

		JustBeforeEach(func() {
			outRoute = envoy_config_route_v3.Route{}
			routesParams := plugins.RouteParams{
				VirtualHostParams: vhostParams,
				VirtualHost:       virtualHost,
			}
			// run it like the translator:
			err := plugin.ProcessRoute(routesParams, route, &outRoute)
			Expect(err).NotTo(HaveOccurred())
			outFilters, err = plugin.HttpFilters(params, httpListener)
			Expect(err).NotTo(HaveOccurred())
		})

		BeforeEach(func() {
			plugin = graphql.NewPlugin()
			err := plugin.Init(plugins.InitParams{})
			Expect(err).NotTo(HaveOccurred())
		})

		Context("empty extensions", func() {
			It("can create the proper filters", func() {
				Expect(outFilters).To(HaveLen(1))
				gqlFilter := outFilters[0]
				Expect(gqlFilter.HttpFilter.Name).To(Equal(graphql.FilterName))
				Expect(gqlFilter.Stage).To(Equal(plugins.BeforeStage(plugins.RouteStage)))
				st := gqlFilter.HttpFilter.GetTypedConfig()
				// graphql is always added to HCM, only routes with graphql route config will use the empty config
				Expect(st).NotTo(BeNil())
			})
		})

		Context("per route/vhost", func() {

			Context("nil", func() {
				BeforeEach(func() {
					route.Action = &v1.Route_RouteAction{
						RouteAction: nil,
					}
				})

				It("is disabled on routes by default", func() {
					pfc := outRoute.TypedPerFilterConfig[graphql.FilterName]
					Expect(pfc).To(BeNil())
				})
			})

			Context("translate route config", func() {
				BeforeEach(func() {
					gqlSchemaSpec.EnableIntrospection = true
				})

				It("sets enable introspection", func() {
					perRouteGql := translateRoute()
					Expect(perRouteGql.EnableIntrospection).To(BeTrue())
				})

				Context("translate resolutions", func() {

					BeforeEach(func() {
						gqlSchemaSpec.Resolutions = []*v1alpha1.Resolution{
							{
								Matcher: &v1alpha1.QueryMatcher{
									Match: &v1alpha1.QueryMatcher_FieldMatcher_{
										FieldMatcher: &v1alpha1.QueryMatcher_FieldMatcher{
											Type:  "type",
											Field: "field",
										},
									},
								},
								Resolver: &v1alpha1.Resolution_RestResolver{
									RestResolver: &v1alpha1.RESTResolver{
										UpstreamRef: &core.ResourceRef{
											Name:      "us",
											Namespace: "gloo-system",
										},
										RequestTransform: &v1alpha1.RequestTemplate{
											Headers: map[string]*v1alpha1.ValueProvider{
												"header": {
													Provider: &v1alpha1.ValueProvider_TypedProvider{
														TypedProvider: &v1alpha1.ValueProvider_TypedValueProvider{
															ValProvider: &v1alpha1.ValueProvider_TypedValueProvider_Value{Value: "7.89"},
															Type:        v1alpha1.ValueProvider_TypedValueProvider_FLOAT,
														},
													},
												},
											},
											QueryParams: map[string]*v1alpha1.ValueProvider{
												"qp": {
													Provider: &v1alpha1.ValueProvider_TypedProvider{
														TypedProvider: &v1alpha1.ValueProvider_TypedValueProvider{
															ValProvider: &v1alpha1.ValueProvider_TypedValueProvider_Value{Value: "true"},
															Type:        v1alpha1.ValueProvider_TypedValueProvider_BOOLEAN,
														},
													},
												},
											},
											OutgoingBody: &v1alpha1.JsonValue{
												// let's test the recursive translation
												JsonVal: &v1alpha1.JsonValue_Node{
													Node: &v1alpha1.JsonNode{
														KeyValues: []*v1alpha1.JsonKeyValue{
															{
																Key: "k1",
																Value: &v1alpha1.JsonValue{
																	JsonVal: &v1alpha1.JsonValue_Node{
																		Node: &v1alpha1.JsonNode{
																			KeyValues: []*v1alpha1.JsonKeyValue{
																				{
																					Key: "k2",
																					Value: &v1alpha1.JsonValue{
																						JsonVal: &v1alpha1.JsonValue_ValueProvider{
																							ValueProvider: &v1alpha1.ValueProvider{
																								Provider: &v1alpha1.ValueProvider_TypedProvider{
																									TypedProvider: &v1alpha1.ValueProvider_TypedValueProvider{
																										Type:        v1alpha1.ValueProvider_TypedValueProvider_STRING,
																										ValProvider: &v1alpha1.ValueProvider_TypedValueProvider_Value{Value: "val"},
																									},
																								},
																							},
																						},
																					},
																				},
																			},
																		},
																	},
																},
															},
														},
													},
												},
											},
										},
										SpanName: "span",
									},
								},
							},
						}
					})

					It("sets resolvers", func() {
						perRouteGql := translateRoute()
						Expect(perRouteGql.Resolutions[0].Matcher.GetFieldMatcher().GetType()).To(Equal("type"))
						Expect(perRouteGql.Resolutions[0].Matcher.GetFieldMatcher().GetField()).To(Equal("field"))

						any := perRouteGql.GetResolutions()[0].GetResolver()
						Expect(any).NotTo(BeNil())
						msg, err := utils.AnyToMessage(any.TypedConfig)
						Expect(err).NotTo(HaveOccurred())
						restResolver := msg.(*v2.RESTResolver)

						Expect(restResolver.GetSpanName()).To(Equal("span"))
						Expect(restResolver.GetRequestTransform().GetHeaders()["header"].GetTypedProvider().GetType()).To(Equal(v2.ValueProvider_TypedValueProvider_FLOAT))
						Expect(restResolver.GetRequestTransform().GetHeaders()["header"].GetTypedProvider().GetValue()).To(Equal("7.89"))
						Expect(restResolver.GetRequestTransform().GetQueryParams()["qp"].GetTypedProvider().GetType()).To(Equal(v2.ValueProvider_TypedValueProvider_BOOLEAN))
						Expect(restResolver.GetRequestTransform().GetQueryParams()["qp"].GetTypedProvider().GetValue()).To(Equal("true"))

						// testing the recursive translation
						Expect(restResolver.GetRequestTransform().GetOutgoingBody().GetNode().GetKeyValues()[0].GetKey()).To(Equal("k1"))
						Expect(restResolver.GetRequestTransform().GetOutgoingBody().GetNode().GetKeyValues()[0].GetValue().GetNode().GetKeyValues()[0].GetKey()).To(Equal("k2"))
						Expect(restResolver.GetRequestTransform().GetOutgoingBody().GetNode().GetKeyValues()[0].GetValue().GetNode().GetKeyValues()[0].GetValue().GetValueProvider().GetTypedProvider().GetValue()).To(Equal("val"))
					})

				})

			})

		})
	})
})
