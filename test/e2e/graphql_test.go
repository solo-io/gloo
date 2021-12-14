package e2e_test

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/fgrosse/zaptest"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/core/matchers"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/graphql/v1alpha1"
	"github.com/solo-io/gloo/projects/gloo/pkg/defaults"
	"github.com/solo-io/gloo/test/helpers"
	"github.com/solo-io/gloo/test/v1helpers"
	"github.com/solo-io/go-utils/contextutils"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/memory"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-projects/test/services"
)

var _ = Describe("graphql", func() {

	var (
		ctx         context.Context
		cancel      context.CancelFunc
		testClients services.TestClients
	)

	var getGraphQLSchema = func(us1Ref *core.ResourceRef) *v1alpha1.GraphQLSchema {
		schema := `
		      schema { query: Query }
		      input Map {
		        a: Int!
		      }
		      type Query {
		        field1(intArg: Int!, boolArg: Boolean!, floatArg: Float!, stringArg: String!, mapArg: Map!, listArg: [Int!]!): SimpleType
		      }
		      type SimpleType {
		        simple: String
		        child: String
		      }
`

		resolutions := []*v1alpha1.Resolution{
			{
				Matcher: &v1alpha1.QueryMatcher{
					Match: &v1alpha1.QueryMatcher_FieldMatcher_{
						FieldMatcher: &v1alpha1.QueryMatcher_FieldMatcher{
							Type:  "Query",
							Field: "field1",
						},
					},
				},
				Resolver: &v1alpha1.Resolution_RestResolver{
					RestResolver: &v1alpha1.RESTResolver{
						UpstreamRef: us1Ref,
						RequestTransform: &v1alpha1.RequestTemplate{
							Headers:      nil, // configured and tested later on
							QueryParams:  nil, // configured and tested later on
							OutgoingBody: nil, // configured and tested later on
						},
						SpanName: "",
					},
				},
			},
		}

		return &v1alpha1.GraphQLSchema{
			Metadata: &core.Metadata{
				Name:      "gql",
				Namespace: "gloo-system",
			},
			Schema:              schema,
			EnableIntrospection: false,
			Resolutions:         resolutions,
		}
	}

	var getProxy = func(envoyPort uint32) *gloov1.Proxy {

		var vhosts []*gloov1.VirtualHost

		vhost := &gloov1.VirtualHost{
			Name:    "gloo-system.virt1",
			Domains: []string{"*"},
			Routes: []*gloov1.Route{
				{
					Matchers: []*matchers.Matcher{{
						PathSpecifier: &matchers.Matcher_Prefix{
							Prefix: "/testroute",
						},
					}},
					Action: &gloov1.Route_GraphqlSchemaRef{
						GraphqlSchemaRef: &core.ResourceRef{
							Name:      "gql",
							Namespace: "gloo-system",
						},
					},
				},
			},
		}

		vhosts = append(vhosts, vhost)

		p := &gloov1.Proxy{
			Metadata: &core.Metadata{
				Name:      "proxy",
				Namespace: "default",
			},
			Listeners: []*gloov1.Listener{{
				Name:        "listener",
				BindAddress: "0.0.0.0",
				BindPort:    envoyPort,
				ListenerType: &gloov1.Listener_HttpListener{
					HttpListener: &gloov1.HttpListener{
						VirtualHosts: vhosts,
					},
				},
			}},
		}

		return p
	}

	BeforeEach(func() {

		logger := zaptest.LoggerWriter(GinkgoWriter)
		contextutils.SetFallbackLogger(logger.Sugar())

		ctx, cancel = context.WithCancel(context.Background())
		cache := memory.NewInMemoryResourceCache()

		testClients = services.GetTestClients(ctx, cache)
		testClients.GlooPort = int(services.AllocateGlooPort())

		what := services.What{
			DisableGateway: true,
			DisableUds:     true,
			DisableFds:     true,
		}

		services.RunGlooGatewayUdsFdsOnPort(ctx, cache, int32(testClients.GlooPort), what, defaults.GlooSystem, nil, nil, nil)
	})

	AfterEach(func() {
		cancel()
	})
	Context("With envoy", func() {
		var (
			envoyInstance *services.EnvoyInstance
			testUpstream1 *v1helpers.TestUpstream
			envoyPort     = uint32(8080)

			proxy         *gloov1.Proxy
			graphQlSchema *v1alpha1.GraphQLSchema
		)

		var testRequest = func(result string) {
			var resp *http.Response
			Eventually(func() (int, error) {
				query := `
{
  "query":"{f:field1(intArg: 2, boolArg: true, floatArg: 9.99993, stringArg: \"this is a string arg\", mapArg: {a: 9}, listArg: [21,22,23]){simple}}"
}
`
				client := http.DefaultClient
				reqUrl, err := url.Parse(fmt.Sprintf("http://%s:%d/testroute", "localhost", envoyPort))
				Expect(err).NotTo(HaveOccurred())
				resp, err = client.Do(&http.Request{
					Method: http.MethodPost,
					URL:    reqUrl,
					Body:   ioutil.NopCloser(strings.NewReader(query)),
				})
				if resp == nil {
					return 0, nil
				}
				return resp.StatusCode, nil
			}, "5s", "0.5s").Should(Equal(http.StatusOK))
			bodyStr, err := ioutil.ReadAll(resp.Body)
			Expect(err).NotTo(HaveOccurred())
			Expect(bodyStr).To(ContainSubstring(result))
		}

		var configureProxy = func() {
			Expect(proxy).NotTo(BeNil())
			_, err := testClients.ProxyClient.Write(proxy, clients.WriteOpts{})
			Expect(err).NotTo(HaveOccurred())

			helpers.EventuallyResourceAccepted(func() (resources.InputResource, error) {
				return testClients.ProxyClient.Read(proxy.Metadata.Namespace, proxy.Metadata.Name, clients.ReadOpts{})
			})
		}

		BeforeEach(func() {
			var err error
			envoyInstance, err = envoyFactory.NewEnvoyInstance()
			Expect(err).NotTo(HaveOccurred())

			err = envoyInstance.Run(testClients.GlooPort)
			Expect(err).NotTo(HaveOccurred())

			testUpstream1 = v1helpers.NewTestHttpUpstreamWithReply(ctx, envoyInstance.LocalAddr(), "{\"simple\":\"foo\"}")

			_, err = testClients.UpstreamClient.Write(testUpstream1.Upstream, clients.WriteOpts{})
			Expect(err).NotTo(HaveOccurred())

			helpers.EventuallyResourceAccepted(func() (resources.InputResource, error) {
				return testClients.UpstreamClient.Read(testUpstream1.Upstream.Metadata.Namespace,
					testUpstream1.Upstream.Metadata.Name, clients.ReadOpts{})
			})
			graphQlSchema = getGraphQLSchema(testUpstream1.Upstream.Metadata.Ref())
		})
		JustBeforeEach(func() {
			_, err := testClients.GraphQLSchemaClient.Write(graphQlSchema, clients.WriteOpts{})
			Expect(err).NotTo(HaveOccurred())

			proxy = getProxy(envoyPort)
			configureProxy()
		})

		AfterEach(func() {
			if envoyInstance != nil {
				envoyInstance.Clean()
			}
		})

		Context("route rules", func() {

			It("resolves graphql queries to REST upstreams", func() {
				testRequest("{\"data\":{\"f\":{\"simple\":\"foo\"}}}")
				Eventually(testUpstream1.C).Should(Receive(PointTo(MatchFields(IgnoreExtras, Fields{
					"URL": PointTo(Equal(url.URL{
						Path: "/",
					})),
				}))))
			})

			Context("with body to upstream", func() {

				BeforeEach(func() {
					graphQlSchema.Resolutions[0].GetRestResolver().RequestTransform.OutgoingBody = &v1alpha1.JsonValue{
						JsonVal: &v1alpha1.JsonValue_Node{
							Node: &v1alpha1.JsonNode{
								KeyValues: []*v1alpha1.JsonKeyValue{
									{
										Key: "key1",
										Value: &v1alpha1.JsonValue{
											JsonVal: &v1alpha1.JsonValue_ValueProvider{
												ValueProvider: &v1alpha1.ValueProvider{
													Providers: map[string]*v1alpha1.ValueProvider_Provider{
														"namedProvider": {
															Provider: &v1alpha1.ValueProvider_Provider_TypedProvider{
																TypedProvider: &v1alpha1.ValueProvider_TypedValueProvider{
																	Type: v1alpha1.ValueProvider_TypedValueProvider_STRING,
																	ValProvider: &v1alpha1.ValueProvider_TypedValueProvider_Value{
																		Value: "value1",
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
					}
				})

				It("resolves graphql queries to REST upstreams with body", func() {
					testRequest("{\"data\":{\"f\":{\"simple\":\"foo\"}}}")
					Eventually(testUpstream1.C).Should(Receive(PointTo(MatchFields(IgnoreExtras, Fields{
						"Body": Equal([]byte("{\"key1\":\"value1\"}")),
						"URL": PointTo(Equal(url.URL{
							Path: "/",
						})),
					}))))
				})
			})

			Context("with query params", func() {

				BeforeEach(func() {
					graphQlSchema.Resolutions[0].GetRestResolver().RequestTransform.QueryParams = map[string]*v1alpha1.ValueProvider{
						"queryparam": {
							Providers: map[string]*v1alpha1.ValueProvider_Provider{
								"namedProvider": {
									Provider: &v1alpha1.ValueProvider_Provider_TypedProvider{
										TypedProvider: &v1alpha1.ValueProvider_TypedValueProvider{
											Type: v1alpha1.ValueProvider_TypedValueProvider_STRING,
											ValProvider: &v1alpha1.ValueProvider_TypedValueProvider_Value{
												Value: "queryparamval",
											},
										},
									},
								},
							},
						},
					}
				})

				It("resolves graphql queries to REST upstreams with query params", func() {
					testRequest("{\"data\":{\"f\":{\"simple\":\"foo\"}}}")
					Eventually(testUpstream1.C).Should(Receive(PointTo(MatchFields(IgnoreExtras, Fields{
						"URL": PointTo(Equal(url.URL{
							Path:     "/",
							RawQuery: "queryparam=queryparamval",
						})),
					}))))
				})
			})

			Context("with headers", func() {

				BeforeEach(func() {
					graphQlSchema.Resolutions[0].GetRestResolver().RequestTransform.Headers = map[string]*v1alpha1.ValueProvider{
						"header": {
							Providers: map[string]*v1alpha1.ValueProvider_Provider{
								"namedProvider": {
									Provider: &v1alpha1.ValueProvider_Provider_TypedProvider{
										TypedProvider: &v1alpha1.ValueProvider_TypedValueProvider{
											Type: v1alpha1.ValueProvider_TypedValueProvider_STRING,
											ValProvider: &v1alpha1.ValueProvider_TypedValueProvider_Value{
												Value: "headerval",
											},
										},
									},
								},
							},
						},
					}
				})

				It("resolves graphql queries to REST upstreams with headers", func() {
					testRequest("{\"data\":{\"f\":{\"simple\":\"foo\"}}}")
					Eventually(testUpstream1.C).Should(Receive(PointTo(MatchFields(IgnoreExtras, Fields{
						"URL": PointTo(Equal(url.URL{
							Path: "/",
						})),
						"Headers": HaveKeyWithValue("Header", []string{"headerval"}),
					}))))
				})
			})

		})
	})
})
