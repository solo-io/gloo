package graphql_test

import (
	"encoding/base64"
	"fmt"

	"github.com/rotisserie/eris"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/utils/graphql/translation"
	"google.golang.org/protobuf/types/known/structpb"
	"k8s.io/utils/lru"

	openapi "github.com/getkin/kin-openapi/openapi3"
	"github.com/ghodss/yaml"
	"github.com/golang/protobuf/ptypes"
	gographql "github.com/graphql-go/graphql"
	. "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/graphql/v1beta1"
	"github.com/solo-io/gloo/projects/gloo/pkg/utils"
	v2 "github.com/solo-io/solo-apis/pkg/api/gloo.solo.io/external/envoy/extensions/filters/http/graphql/v2"
	"github.com/solo-io/solo-kit/test/matchers"
	graphql2 "github.com/solo-io/solo-projects/projects/discovery/pkg/fds/discoveries/openapi-graphql/graphqlschematranslation"
	"github.com/solo-io/solo-projects/projects/discovery/pkg/fds/discoveries/openapi-graphql/printer"
	schemas "github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/graphql"

	envoy_config_route_v3 "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	v1snap "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/gloosnapshot"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
)

type MockArtifactsList struct{}

func (m *MockArtifactsList) Find(namespace string, name string) (*v1.Artifact, error) {
	var proto string

	if name == "fake-artifact-one" || name == "fake-artifact-two" {
		proto = `ClEKC21vdmllLnByb3RvEgVtb3ZpZSIzCgVNb3ZpZRISCgRuYW1lGAEgASgJUgRuYW1lEhYKBnJhdGluZxgCIAEoBVIGcmF0aW5nYgZwcm90bzM=`
	} else if name == "illegal-artifact" {
		proto = `&$(*@#$&@! A bunch of illegal bytes here +_!(#_+#`
	} else if name == "malformed-artifact" {
		proto = `ClEKC21vdmllLnByb3RvEgVtb3ZpZSIzCgVNb3ISCgRuYW1lGAEgASgJUgRuYW1lEhYKBnJhdGluZxgCIAEoBVIGcmF0aW5nYgZwcm90bzM=`
	} else {
		return nil, eris.Errorf("Some could not find error")
	}

	return &v1.Artifact{
		Metadata: &core.Metadata{
			Name:      name,
			Namespace: namespace,
		},
		Data: map[string]string{
			"proto": proto,
		},
	}, nil
}

type EmptyMockArtifactsList struct{}

func (m *EmptyMockArtifactsList) Find(namespace string, name string) (*v1.Artifact, error) {
	return &v1.Artifact{
		Metadata: &core.Metadata{
			Name:      name,
			Namespace: namespace,
		},
		Data: map[string]string{},
	}, nil
}

var _ = Describe("Graphql plugin", func() {
	var (
		plugin       plugins.Plugin
		params       plugins.Params
		vhostParams  plugins.VirtualHostParams
		virtualHost  *v1.VirtualHost
		route        *v1.Route
		routeAction  *v1.Route_GraphqlApiRef
		httpListener *v1.HttpListener
		gqlApiSpec   *GraphQLApi
	)

	BeforeEach(func() {
		routeAction = &v1.Route_GraphqlApiRef{
			GraphqlApiRef: &core.ResourceRef{
				Name:      "gql",
				Namespace: "gloo-system",
			},
		}
		route = &v1.Route{
			Action: routeAction,
		}

		gqlApiSpec = &GraphQLApi{
			Metadata: &core.Metadata{
				Name:      "gql",
				Namespace: "gloo-system",
			},
			Schema: &GraphQLApi_ExecutableSchema{
				ExecutableSchema: &ExecutableSchema{
					Executor: &Executor{
						Executor: &Executor_Local_{
							Local: &Executor_Local{},
						},
					},
				},
			},
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

		params.Snapshot = &v1snap.ApiSnapshot{
			Proxies: v1.ProxyList{proxy},
			GraphqlApis: GraphQLApiList{
				gqlApiSpec,
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
		messages := (make(map[*core.ResourceRef][]string))
		params.Messages = messages
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
			goTpfc := outRoute.TypedPerFilterConfig[schemas.FilterName]
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
			err := plugin.(plugins.RoutePlugin).ProcessRoute(routesParams, route, &outRoute)
			Expect(err).NotTo(HaveOccurred())
			outFilters, err = plugin.(plugins.HttpFilterPlugin).HttpFilters(params, httpListener)
			Expect(err).NotTo(HaveOccurred())
		})

		BeforeEach(func() {
			plugin = schemas.NewPlugin()
			plugin.Init(plugins.InitParams{})
		})

		Context("empty extensions", func() {
			It("can create the proper filters", func() {
				Expect(outFilters).To(HaveLen(1))
				gqlFilter := outFilters[0]
				Expect(gqlFilter.HttpFilter.Name).To(Equal(schemas.FilterName))
				Expect(gqlFilter.Stage).To(Equal(plugins.BeforeStage(plugins.RouteStage)))
				tc := gqlFilter.HttpFilter.GetTypedConfig()
				// graphql is always added to HCM, only routes with graphql route config will use the empty config
				Expect(tc).NotTo(BeNil())
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
					pfc := outRoute.TypedPerFilterConfig[schemas.FilterName]
					Expect(pfc).To(BeNil())
				})
			})

			Context("translates gRPC config", func() {
				It("properly accumulates and returns referenced configmaps", func() {
					output, err := translation.TranslateExtensions(&MockArtifactsList{}, &GraphQLApi{
						Schema: &GraphQLApi_ExecutableSchema{
							ExecutableSchema: &ExecutableSchema{
								GrpcDescriptorRegistry: &GrpcDescriptorRegistry{
									DescriptorSet: &GrpcDescriptorRegistry_ProtoRefsList{
										ProtoRefsList: &GrpcDescriptorRegistry_ProtoRefs{
											ConfigMapRefs: []*core.ResourceRef{
												&core.ResourceRef{
													Name:      "fake-artifact-one",
													Namespace: "fake-namespace",
												},
											},
										},
									},
								},
							},
						},
					})
					var expected = "ClEKC21vdmllLnByb3RvEgVtb3ZpZSIzCgVNb3ZpZRISCgRuYW1lGAEgASgJUgRuYW1lEhYKBnJhdGluZxgCIAEoBVIGcmF0aW5nYgZwcm90bzM="
					expectedBytes, _ := base64.StdEncoding.DecodeString(expected)
					message, _ := utils.AnyToMessage(output["grpc_extension"])
					bytes := message.(*v2.GrpcDescriptorRegistry).ProtoDescriptors.GetInlineBytes()
					Expect(bytes).To(Equal(expectedBytes))
					Expect(err).ToNot(HaveOccurred())
				})
				It("properly deduplictes protos", func() {
					output, err := translation.TranslateExtensions(&MockArtifactsList{}, &GraphQLApi{
						Schema: &GraphQLApi_ExecutableSchema{
							ExecutableSchema: &ExecutableSchema{
								GrpcDescriptorRegistry: &GrpcDescriptorRegistry{
									DescriptorSet: &GrpcDescriptorRegistry_ProtoRefsList{
										ProtoRefsList: &GrpcDescriptorRegistry_ProtoRefs{
											ConfigMapRefs: []*core.ResourceRef{
												&core.ResourceRef{
													Name:      "fake-artifact-one",
													Namespace: "fake-namespace",
												}, &core.ResourceRef{
													Name:      "fake-artifact-two",
													Namespace: "fake-namespace",
												},
											},
										},
									},
								},
							},
						},
					})

					var expected = "ClEKC21vdmllLnByb3RvEgVtb3ZpZSIzCgVNb3ZpZRISCgRuYW1lGAEgASgJUgRuYW1lEhYKBnJhdGluZxgCIAEoBVIGcmF0aW5nYgZwcm90bzM="
					expectedBytes, _ := base64.StdEncoding.DecodeString(expected)
					message, _ := utils.AnyToMessage(output["grpc_extension"])
					bytes := message.(*v2.GrpcDescriptorRegistry).ProtoDescriptors.GetInlineBytes()
					Expect(bytes).To(Equal(expectedBytes))
					Expect(err).ToNot(HaveOccurred())
				})
				It("fails if no matching configmap exists in artifacts", func() {
					_, err := translation.TranslateExtensions(&MockArtifactsList{}, &GraphQLApi{
						Schema: &GraphQLApi_ExecutableSchema{
							ExecutableSchema: &ExecutableSchema{
								GrpcDescriptorRegistry: &GrpcDescriptorRegistry{
									DescriptorSet: &GrpcDescriptorRegistry_ProtoRefsList{
										ProtoRefsList: &GrpcDescriptorRegistry_ProtoRefs{
											ConfigMapRefs: []*core.ResourceRef{
												&core.ResourceRef{
													Name:      "missing-artifact",
													Namespace: "fake-namespace",
												},
											},
										},
									},
								},
							},
						},
					})
					Expect(err).To(HaveOccurred())
					Expect(err).To(MatchError(ContainSubstring("Could not find ConfigMap with ref fake-namespace.missing-artifact to use a gRPC proto registry source")))
				})
				It("fails if no keys exist on configmaps", func() {
					_, err := translation.TranslateExtensions(&EmptyMockArtifactsList{}, &GraphQLApi{
						Schema: &GraphQLApi_ExecutableSchema{
							ExecutableSchema: &ExecutableSchema{
								GrpcDescriptorRegistry: &GrpcDescriptorRegistry{
									DescriptorSet: &GrpcDescriptorRegistry_ProtoRefsList{
										ProtoRefsList: &GrpcDescriptorRegistry_ProtoRefs{
											ConfigMapRefs: []*core.ResourceRef{
												&core.ResourceRef{
													Name:      "fake-artifact",
													Namespace: "fake-namespace",
												},
											},
										},
									},
								},
							},
						},
					})
					Expect(err).To(HaveOccurred())
					Expect(err).To(MatchError(ContainSubstring("No keys exist in fake-namespace.fake-artifact")))
				})
				It("fails if invalid value exists on configmap", func() {
					_, err := translation.TranslateExtensions(&MockArtifactsList{}, &GraphQLApi{
						Schema: &GraphQLApi_ExecutableSchema{
							ExecutableSchema: &ExecutableSchema{
								GrpcDescriptorRegistry: &GrpcDescriptorRegistry{
									DescriptorSet: &GrpcDescriptorRegistry_ProtoRefsList{
										ProtoRefsList: &GrpcDescriptorRegistry_ProtoRefs{
											ConfigMapRefs: []*core.ResourceRef{
												&core.ResourceRef{
													Name:      "illegal-artifact",
													Namespace: "fake-namespace",
												},
											},
										},
									},
								},
							},
						},
					})
					Expect(err).To(HaveOccurred())
					Expect(err).To(MatchError(ContainSubstring("Error decoding proto data in proto: &$(*@#$&@! A bunch of illegal bytes here +_!(#_+#")))
				})
				It("fails if concatenated protos arent valid", func() {
					_, err := translation.TranslateExtensions(&MockArtifactsList{}, &GraphQLApi{
						Schema: &GraphQLApi_ExecutableSchema{
							ExecutableSchema: &ExecutableSchema{
								GrpcDescriptorRegistry: &GrpcDescriptorRegistry{
									DescriptorSet: &GrpcDescriptorRegistry_ProtoRefsList{
										ProtoRefsList: &GrpcDescriptorRegistry_ProtoRefs{
											ConfigMapRefs: []*core.ResourceRef{
												&core.ResourceRef{
													Name:      "fake-artifact-one",
													Namespace: "fake-namespace",
												},
												&core.ResourceRef{
													Name:      "malformed-artifact",
													Namespace: "fake-namespace",
												},
											},
										},
									},
								},
							},
						},
					})
					Expect(err).To(HaveOccurred())
					Expect(err).To(MatchError(ContainSubstring("key proto in configMap fake-namespace.malformed-artifact does not contain valid proto bytes")))
				})
			})
			Context("translate route config", func() {
				BeforeEach(func() {
					gqlApiSpec.GetExecutableSchema().GetExecutor().GetLocal().EnableIntrospection = true
				})

				It("sets enable introspection", func() {
					perRouteGql := translateRoute()
					Expect(perRouteGql.GetExecutableSchema().GetExecutor().GetLocal().GetEnableIntrospection()).To(BeTrue())
				})

				Context("translate resolutions", func() {

					BeforeEach(func() {
						body := `{"k1": {"k2": "val"}}`
						bodyStruct := &structpb.Value{}
						err := yaml.Unmarshal([]byte(body), bodyStruct)
						Expect(err).NotTo(HaveOccurred())
						// Only resolvers referenced in the schema are translated.
						gqlApiSpec.GetExecutableSchema().SchemaDefinition = `type Query {
	field1: String @resolve(name: "resolver1") @cacheControl(maxAge: 60, inheritMaxAge: false, scope: private)
}`
						gqlApiSpec.GetExecutableSchema().Executor.GetLocal().Resolutions = map[string]*Resolution{
							"resolver1": {
								Resolver: &Resolution_RestResolver{
									RestResolver: &RESTResolver{
										UpstreamRef: &core.ResourceRef{
											Name:      "us",
											Namespace: "gloo-system",
										},
										Request: &RequestTemplate{
											Headers: map[string]string{
												"header": "7.89",
											},
											QueryParams: map[string]string{
												"qp": "true",
											},
											Body: bodyStruct,
										},
										SpanName: "span",
									},
								},
							},
						}
					})

					Context("with type-level directive", func() {

						Context("type-level directives only", func() {
							BeforeEach(func() {
								gqlApiSpec.GetExecutableSchema().SchemaDefinition = `type Query @cacheControl(maxAge: 70, inheritMaxAge: false, scope: private) {
	field1: String @resolve(name: "resolver1")
}`
							})
							It("sets resolvers and cache control defaults", func() {
								perRouteGql := translateRoute()
								resolutions := perRouteGql.GetExecutableSchema().GetExecutor().GetLocal().GetResolutions()
								Expect(resolutions[0].Matcher.GetFieldMatcher().GetType()).To(Equal("Query"))
								Expect(resolutions[0].Matcher.GetFieldMatcher().GetField()).To(Equal("field1"))

								Expect(resolutions[0].GetCacheControl().GetMaxAge().GetValue()).To(Equal(uint32(70)))
								Expect(resolutions[0].GetCacheControl().GetInheritMaxAge()).To(Equal(false))
								Expect(resolutions[0].GetCacheControl().GetScope().String()).To(Equal("PRIVATE"))

								any := resolutions[0].GetResolver()
								Expect(any).NotTo(BeNil())
								msg, err := utils.AnyToMessage(any.TypedConfig)
								Expect(err).NotTo(HaveOccurred())
								restResolver := msg.(*v2.RESTResolver)

								Expect(restResolver.GetSpanName()).To(Equal("span"))
								Expect(restResolver.GetRequestTransform().GetHeaders()["header"].GetProviders()["ARBITRARY_PROVIDER_NAME"].GetTypedProvider().GetValue()).To(Equal("7.89"))
								Expect(restResolver.GetRequestTransform().GetQueryParams()["qp"].GetProviders()["ARBITRARY_PROVIDER_NAME"].GetTypedProvider().GetValue()).To(Equal("true"))

								// testing the recursive translation
								Expect(restResolver.GetRequestTransform().GetOutgoingBody().GetNode().GetKeyValues()[0].GetKey()).To(Equal("k1"))
								Expect(restResolver.GetRequestTransform().GetOutgoingBody().GetNode().GetKeyValues()[0].GetValue().GetNode().GetKeyValues()[0].GetKey()).To(Equal("k2"))
								Expect(restResolver.GetRequestTransform().GetOutgoingBody().GetNode().GetKeyValues()[0].GetValue().GetNode().GetKeyValues()[0].GetValue().GetValueProvider().GetProviders()["ARBITRARY_PROVIDER_NAME"].GetTypedProvider().GetValue()).To(Equal("val"))
							})
						})

						Context("type-level and field-level directives", func() {

							BeforeEach(func() {
								gqlApiSpec.GetExecutableSchema().SchemaDefinition = `type Query @cacheControl(maxAge: 70, inheritMaxAge: false, scope: private) {
	field1: String @resolve(name: "resolver1") @cacheControl(maxAge: 90, inheritMaxAge: false, scope: public)
}`
							})

							It("sets resolvers and cache control defaults -- field-level cache control overrides type-level configuration", func() {
								perRouteGql := translateRoute()
								resolutions := perRouteGql.GetExecutableSchema().GetExecutor().GetLocal().GetResolutions()
								Expect(resolutions[0].Matcher.GetFieldMatcher().GetType()).To(Equal("Query"))
								Expect(resolutions[0].Matcher.GetFieldMatcher().GetField()).To(Equal("field1"))

								Expect(resolutions[0].GetCacheControl().GetMaxAge().GetValue()).To(Equal(uint32(90)))
								Expect(resolutions[0].GetCacheControl().GetInheritMaxAge()).To(Equal(false))
								Expect(resolutions[0].GetCacheControl().GetScope().String()).To(Equal("PUBLIC"))

								any := resolutions[0].GetResolver()
								Expect(any).NotTo(BeNil())
								msg, err := utils.AnyToMessage(any.TypedConfig)
								Expect(err).NotTo(HaveOccurred())
								restResolver := msg.(*v2.RESTResolver)

								Expect(restResolver.GetSpanName()).To(Equal("span"))
								Expect(restResolver.GetRequestTransform().GetHeaders()["header"].GetProviders()["ARBITRARY_PROVIDER_NAME"].GetTypedProvider().GetValue()).To(Equal("7.89"))
								Expect(restResolver.GetRequestTransform().GetQueryParams()["qp"].GetProviders()["ARBITRARY_PROVIDER_NAME"].GetTypedProvider().GetValue()).To(Equal("true"))

								// testing the recursive translation
								Expect(restResolver.GetRequestTransform().GetOutgoingBody().GetNode().GetKeyValues()[0].GetKey()).To(Equal("k1"))
								Expect(restResolver.GetRequestTransform().GetOutgoingBody().GetNode().GetKeyValues()[0].GetValue().GetNode().GetKeyValues()[0].GetKey()).To(Equal("k2"))
								Expect(restResolver.GetRequestTransform().GetOutgoingBody().GetNode().GetKeyValues()[0].GetValue().GetNode().GetKeyValues()[0].GetValue().GetValueProvider().GetProviders()["ARBITRARY_PROVIDER_NAME"].GetTypedProvider().GetValue()).To(Equal("val"))
							})
						})
					})

					It("sets resolvers and cache control defaults", func() {
						perRouteGql := translateRoute()
						resolutions := perRouteGql.GetExecutableSchema().GetExecutor().GetLocal().GetResolutions()
						Expect(resolutions[0].Matcher.GetFieldMatcher().GetType()).To(Equal("Query"))
						Expect(resolutions[0].Matcher.GetFieldMatcher().GetField()).To(Equal("field1"))

						Expect(resolutions[0].GetCacheControl().GetMaxAge().GetValue()).To(Equal(uint32(60)))
						Expect(resolutions[0].GetCacheControl().GetInheritMaxAge()).To(Equal(false))
						Expect(resolutions[0].GetCacheControl().GetScope().String()).To(Equal("PRIVATE"))

						any := resolutions[0].GetResolver()
						Expect(any).NotTo(BeNil())
						msg, err := utils.AnyToMessage(any.TypedConfig)
						Expect(err).NotTo(HaveOccurred())
						restResolver := msg.(*v2.RESTResolver)

						Expect(restResolver.GetSpanName()).To(Equal("span"))
						Expect(restResolver.GetRequestTransform().GetHeaders()["header"].GetProviders()["ARBITRARY_PROVIDER_NAME"].GetTypedProvider().GetValue()).To(Equal("7.89"))
						Expect(restResolver.GetRequestTransform().GetQueryParams()["qp"].GetProviders()["ARBITRARY_PROVIDER_NAME"].GetTypedProvider().GetValue()).To(Equal("true"))

						// testing the recursive translation
						Expect(restResolver.GetRequestTransform().GetOutgoingBody().GetNode().GetKeyValues()[0].GetKey()).To(Equal("k1"))
						Expect(restResolver.GetRequestTransform().GetOutgoingBody().GetNode().GetKeyValues()[0].GetValue().GetNode().GetKeyValues()[0].GetKey()).To(Equal("k2"))
						Expect(restResolver.GetRequestTransform().GetOutgoingBody().GetNode().GetKeyValues()[0].GetValue().GetNode().GetKeyValues()[0].GetValue().GetValueProvider().GetProviders()["ARBITRARY_PROVIDER_NAME"].GetTypedProvider().GetValue()).To(Equal("val"))
					})
				})
			})
			Context("Translate schema config for remote executor", func() {
				BeforeEach(func() {
					gqlApiSpec = &GraphQLApi{
						Metadata: &core.Metadata{
							Name:      "gql",
							Namespace: "gloo-system",
						},
						Schema: &GraphQLApi_ExecutableSchema{
							ExecutableSchema: &ExecutableSchema{
								Executor: &Executor{
									Executor: &Executor_Remote_{
										Remote: &Executor_Remote{
											UpstreamRef: &core.ResourceRef{
												Name:      "us",
												Namespace: "gloo-system",
											},
											Headers: map[string]string{
												"foo": "far",
												"boo": "{$headers.bar}",
												"zoo": "{$metadata.io.solo.transformation:endpoint_url}",
											},
											QueryParams: map[string]string{
												"moo": "mar",
												"noo": "{$headers.nar}",
											},
											SpanName: "TestSpanName",
										},
									},
								},
							},
						},
					}
				})
				It("Translates user facing api to envoy api for remote executors", func() {
					upstreams := v1.UpstreamList{
						{
							Metadata: &core.Metadata{
								Name:      "us",
								Namespace: "gloo-system",
							},
						},
					}
					api, err := translation.CreateGraphQlApi(translation.CreateGraphQLApiParams{&MockArtifactsList{}, upstreams, nil, gqlApiSpec, lru.New(1024), lru.New(1024)})
					Expect(err).ToNot(HaveOccurred())
					translatedExecutor := api.GetExecutor().GetRemote()
					Expect(translatedExecutor.GetSpanName()).To(Equal("TestSpanName"))
					Expect(translatedExecutor.GetRequest().GetHeaders()["foo"].GetValue()).To(Equal("far"))
					Expect(translatedExecutor.GetRequest().GetHeaders()["boo"].GetHeader()).To(Equal("bar"))
					Expect(translatedExecutor.GetRequest().GetHeaders()["zoo"].GetDynamicMetadata().GetMetadataNamespace()).To(Equal("io.solo.transformation"))
					Expect(translatedExecutor.GetRequest().GetHeaders()["zoo"].GetDynamicMetadata().GetKey()).To(Equal("endpoint_url"))
					Expect(translatedExecutor.GetRequest().GetQueryParams()["moo"].GetValue()).To(Equal("mar"))
					Expect(translatedExecutor.GetRequest().GetQueryParams()["noo"].GetHeader()).To(Equal("nar"))
				})
			})
			Context("Incorrectly translates schema config for remote executor", func() {
				JustBeforeEach(func() {
					gqlApiSpec = &GraphQLApi{
						Metadata: &core.Metadata{
							Name:      "gql",
							Namespace: "gloo-system",
						},
						Schema: &GraphQLApi_ExecutableSchema{
							ExecutableSchema: &ExecutableSchema{
								Executor: &Executor{
									Executor: &Executor_Remote_{
										Remote: &Executor_Remote{
											UpstreamRef: &core.ResourceRef{
												Name:      "us",
												Namespace: "gloo-system",
											},
											Headers: map[string]string{
												"zoo": "{$metadata.}",
											},
											SpanName: "TestSpanName",
										},
									},
								},
							},
						},
					}
				})
				It("Incorrectly translates user facing api to envoy api for remote executors", func() {
					upstreams := v1.UpstreamList{
						{
							Metadata: &core.Metadata{
								Name:      "us",
								Namespace: "gloo-system",
							},
						},
					}
					_, err := translation.CreateGraphQlApi(translation.CreateGraphQLApiParams{&MockArtifactsList{}, upstreams, nil, gqlApiSpec, nil, nil})
					Expect(err).To(MatchError(ContainSubstring("Malformed value for dynamic metadata zoo: {$metadata.}")))

				})
			})
		})
		Context("graphql translation", func() {

			var (
				graphqlSchema *gographql.Schema
				resolutions   map[string]*Resolution
			)
			AfterEach(func() {
				graphqlSchema = nil
				resolutions = nil
			})
			translateToSchema := func(openapiSchema string) {
				t := graphql2.NewOasToGqlTranslator(&v1.Upstream{
					Metadata: &core.Metadata{
						Name:      "hi",
						Namespace: "default",
					},
				})
				l := openapi.NewLoader()
				l.IsExternalRefsAllowed = true
				spec, err := l.LoadFromData([]byte(openapiSchema))
				ExpectWithOffset(1, err).NotTo(HaveOccurred())
				oass := []*openapi.T{spec}
				_, graphqlSchema, resolutions, err = t.CreateGraphqlSchema(oass)
				ExpectWithOffset(1, err).NotTo(HaveOccurred())
				fmt.Println(printer.PrintFilteredSchema(graphqlSchema))
				crd := &GraphQLApi{
					Schema: &GraphQLApi_ExecutableSchema{
						ExecutableSchema: &ExecutableSchema{
							Executor: &Executor{
								Executor: &Executor_Local_{
									Local: &Executor_Local{
										Resolutions: resolutions,
									},
								},
							},
							SchemaDefinition: printer.PrintFilteredSchema(graphqlSchema),
						},
					},
				}
				b, err := yaml.Marshal(crd)
				ExpectWithOffset(1, err).NotTo(HaveOccurred())
				fmt.Println(string(b))
			}

			It("Handles links", func() {
				schemaToTranslate := schemas.GetSimpleLinkSchema()
				translateToSchema(schemaToTranslate)

				fields := graphqlSchema.QueryType().Fields()
				// Expect getEmployeeById query field to be created
				Expect(fields).To(HaveKey("getEmployeeById"))
				employeeByIdQueryField := fields["getEmployeeById"]
				Expect(employeeByIdQueryField.Args).To(HaveLen(1))
				Expect(employeeByIdQueryField.Args[0].Type).To(Equal(gographql.NewNonNull(gographql.String)))
				Expect(employeeByIdQueryField.Type.String()).To(Equal("Employee"))

				// Employee type should have fields, and userManager link should be resolved.
				employeeType := graphqlSchema.Type("Employee").(*gographql.Object)
				Expect(employeeType).To(Not(BeNil()))
				Expect(employeeType.Fields()).To(HaveLen(4))
				Expect(employeeType.Fields()).To(HaveKey("userManager"))
				userManagerField := employeeType.Fields()["userManager"]
				Expect(userManagerField.Type.String()).To(Equal("Employee"))

				// Resolvers should exist for Query.getEmployeeById and Employee.userManager
				Expect(resolutions).To(HaveLen(2))
				Expect(resolutions).To(ContainElement(matchers.MatchProto(&Resolution{
					Resolver: &Resolution_RestResolver{
						RestResolver: &RESTResolver{
							UpstreamRef: &core.ResourceRef{
								Name:      "hi",
								Namespace: "default",
							},
							Request: &RequestTemplate{
								Headers: map[string]string{
									":method": "GET",
									":path":   "/employees/{$args.id}",
								},
							},
						},
					},
				})))
				Expect(resolutions).To(ContainElement(matchers.MatchProto(&Resolution{
					Resolver: &Resolution_RestResolver{
						RestResolver: &RESTResolver{
							UpstreamRef: &core.ResourceRef{
								Name:      "hi",
								Namespace: "default",
							},
							Request: &RequestTemplate{
								Headers: map[string]string{
									":method": "GET",
									":path":   "/employees/{$parent.managerId}",
								},
							},
						},
					},
				})))
			})
		})
	})
})
