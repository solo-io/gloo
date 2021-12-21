package e2e_test

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	structpb "github.com/golang/protobuf/ptypes/struct"

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
	glootest "github.com/solo-io/gloo/test/v1helpers/test_grpc_service/glootest/protos"
	"github.com/solo-io/go-utils/contextutils"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/memory"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-projects/test/services"
	glootestpb "github.com/solo-io/solo-projects/test/v1helpers/test_grpc_service"
)

var _ = Describe("graphql", func() {

	var (
		ctx         context.Context
		cancel      context.CancelFunc
		testClients services.TestClients
	)

	var getGraphQLSchema = func(restUsRef, grpcUsRef *core.ResourceRef) *v1alpha1.GraphQLSchema {
		schema := `
		      schema { query: Query }
		      input Map {
		        a: Int!
		      }
		      type Query {
		        field1(intArg: Int!, boolArg: Boolean!, floatArg: Float!, stringArg: String!, mapArg: Map!, listArg: [Int!]!): SimpleType
		        field2: TestResponse
		      }
		      type SimpleType {
		        simple: String
		        child: String
		      }
		      type TestResponse {
		        str: String
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
						UpstreamRef: restUsRef,
						Request: &v1alpha1.RequestTemplate{
							Headers:     nil, // configured and tested later on
							QueryParams: nil, // configured and tested later on
							Body:        nil, // configured and tested later on
						},
						SpanName: "",
					},
				},
			},
			{
				Matcher: &v1alpha1.QueryMatcher{
					Match: &v1alpha1.QueryMatcher_FieldMatcher_{
						FieldMatcher: &v1alpha1.QueryMatcher_FieldMatcher{
							Type:  "Query",
							Field: "field2",
						},
					},
				},
				Resolver: &v1alpha1.Resolution_GrpcResolver{
					GrpcResolver: &v1alpha1.GrpcResolver{
						UpstreamRef: grpcUsRef,
						RequestTransform: &v1alpha1.GrpcRequestTemplate{
							OutgoingMessageJson: &structpb.Value{
								Kind: &structpb.Value_StructValue{
									StructValue: &structpb.Struct{
										Fields: map[string]*structpb.Value{
											"str": {Kind: &structpb.Value_StringValue{StringValue: "foo"}},
										},
									},
								},
							},
							ServiceName: "glootest.TestService",
							MethodName:  "TestMethod",
						},
					},
				},
			},
		}

		return &v1alpha1.GraphQLSchema{
			Metadata: &core.Metadata{
				Name:      "gql",
				Namespace: "gloo-system",
			},
			ExecutableSchema: &v1alpha1.ExecutableSchema{
				SchemaDefinition: schema,
				Executor: &v1alpha1.Executor{
					Executor: &v1alpha1.Executor_Local_{
						Local: &v1alpha1.Executor_Local{
							Resolutions:         resolutions,
							EnableIntrospection: false,
						},
					},
				},
				GrpcDescriptorRegistry: &v1alpha1.GrpcDescriptorRegistry{
					DescriptorSet: &v1alpha1.GrpcDescriptorRegistry_ProtoDescriptorBin{
						ProtoDescriptorBin: glootestpb.ProtoBytes,
					},
				},
			},
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
			envoyInstance              *services.EnvoyInstance
			restUpstream, grpcUpstream *v1helpers.TestUpstream
			envoyPort                  = uint32(8080)
			query                      string

			proxy         *gloov1.Proxy
			graphQlSchema *v1alpha1.GraphQLSchema
		)

		var testRequest = func(result string) {
			var resp *http.Response
			Eventually(func() (int, error) {
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

			query = `
{
  "query":"{f:field1(intArg: 2, boolArg: true, floatArg: 9.99993, stringArg: \"this is a string arg\", mapArg: {a: 9}, listArg: [21,22,23]){simple}}"
}`
			restUpstream = v1helpers.NewTestHttpUpstreamWithReply(ctx, envoyInstance.LocalAddr(), "{\"simple\":\"foo\"}")
			_, err = testClients.UpstreamClient.Write(restUpstream.Upstream, clients.WriteOpts{})
			Expect(err).NotTo(HaveOccurred())
			helpers.EventuallyResourceAccepted(func() (resources.InputResource, error) {
				return testClients.UpstreamClient.Read(restUpstream.Upstream.Metadata.Namespace,
					restUpstream.Upstream.Metadata.Name, clients.ReadOpts{})
			})

			grpcUpstream = v1helpers.NewTestGRPCUpstream(ctx, envoyInstance.LocalAddr(), 1)
			_, err = testClients.UpstreamClient.Write(grpcUpstream.Upstream, clients.WriteOpts{})
			Expect(err).NotTo(HaveOccurred())
			helpers.EventuallyResourceAccepted(func() (resources.InputResource, error) {
				return testClients.UpstreamClient.Read(grpcUpstream.Upstream.Metadata.Namespace,
					grpcUpstream.Upstream.Metadata.Name, clients.ReadOpts{})
			})

			graphQlSchema = getGraphQLSchema(restUpstream.Upstream.Metadata.Ref(), grpcUpstream.Upstream.Metadata.Ref())
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
				Eventually(restUpstream.C).Should(Receive(PointTo(MatchFields(IgnoreExtras, Fields{
					"URL": PointTo(Equal(url.URL{
						Path: "/",
					})),
				}))))
			})

			Context("grpc resolver", func() {
				BeforeEach(func() {
					query = `
{
  "query":"{f:field2{str}}"
}`
				})

				It("resolves graphql queries to REST upstreams", func() {
					testRequest("{\"data\":{\"f\":{\"str\":\"foo\"}}}")
					Eventually(grpcUpstream.C).Should(Receive(PointTo(MatchFields(IgnoreExtras, Fields{
						"GRPCRequest": PointTo(Equal(glootest.TestRequest{Str: "foo"})),
					}))))
				})
			})

			Context("with body to upstream", func() {

				BeforeEach(func() {
					body := &structpb.Value{
						Kind: &structpb.Value_StructValue{
							StructValue: &structpb.Struct{
								Fields: map[string]*structpb.Value{
									"key1": {Kind: &structpb.Value_StringValue{StringValue: "value1"}},
								},
							},
						},
					}
					graphQlSchema.ExecutableSchema.GetExecutor().GetLocal().GetResolutions()[0].GetRestResolver().Request.Body = body
				})

				It("resolves graphql queries to REST upstreams with body", func() {
					testRequest("{\"data\":{\"f\":{\"simple\":\"foo\"}}}")
					Eventually(restUpstream.C).Should(Receive(PointTo(MatchFields(IgnoreExtras, Fields{
						"Body": Equal([]byte("{\"key1\":\"value1\"}")),
						"URL": PointTo(Equal(url.URL{
							Path: "/",
						})),
					}))))
				})
			})

			Context("with query params", func() {

				BeforeEach(func() {
					graphQlSchema.ExecutableSchema.GetExecutor().GetLocal().GetResolutions()[0].GetRestResolver().Request.QueryParams = map[string]string{
						"queryparam": "queryparamval",
					}
				})

				It("resolves graphql queries to REST upstreams with query params", func() {
					testRequest("{\"data\":{\"f\":{\"simple\":\"foo\"}}}")
					Eventually(restUpstream.C).Should(Receive(PointTo(MatchFields(IgnoreExtras, Fields{
						"URL": PointTo(Equal(url.URL{
							Path:     "/",
							RawQuery: "queryparam=queryparamval",
						})),
					}))))
				})
			})

			Context("with headers", func() {

				BeforeEach(func() {
					graphQlSchema.ExecutableSchema.GetExecutor().GetLocal().GetResolutions()[0].GetRestResolver().Request.Headers = map[string]string{
						"header": "headerval",
					}
				})

				It("resolves graphql queries to REST upstreams with headers", func() {
					testRequest("{\"data\":{\"f\":{\"simple\":\"foo\"}}}")
					Eventually(restUpstream.C).Should(Receive(PointTo(MatchFields(IgnoreExtras, Fields{
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
