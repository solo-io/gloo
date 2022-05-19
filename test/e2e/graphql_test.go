package e2e_test

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"

	"github.com/solo-io/licensing/pkg/model"

	structpb "github.com/golang/protobuf/ptypes/struct"

	"github.com/fgrosse/zaptest"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/core/matchers"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/graphql/v1beta1"
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
		ctx                    context.Context
		cancel                 context.CancelFunc
		testClients            services.TestClients
		executor               v1beta1.Executor
		grpcDescriptorRegistry v1beta1.GrpcDescriptorRegistry
	)

	var getGraphQLApi = func(restUsRef, grpcUsRef *core.ResourceRef) *v1beta1.GraphQLApi {
		schema := `
		      schema { query: Query }
		      input Map {
		        a: Int!
		      }
		      type Query {
		        field1(intArg: Int!, boolArg: Boolean!, floatArg: Float!, stringArg: String!, mapArg: Map!, listArg: [Int!]!): SimpleType
							@resolve(name: "simple_resolver")

		        field2: TestResponse @resolve(name: "field2_resolver")

		        field3: SimpleType @resolve(name: "simple_resolver") @cacheControl(maxAge: 60, scope: private)
		      }
		      type SimpleType {
		        simple: String
		        setme: String
		      }
		      type TestResponse {
		        str: String
		      }`

		resolutions := map[string]*v1beta1.Resolution{
			"simple_resolver": {
				Resolver: &v1beta1.Resolution_RestResolver{
					RestResolver: &v1beta1.RESTResolver{
						UpstreamRef: restUsRef,
						Request: &v1beta1.RequestTemplate{
							Headers:     nil, // configured and tested later on
							QueryParams: nil, // configured and tested later on
							Body:        nil, // configured and tested later on
						},
						SpanName: "",
					},
				},
			},
			"field2_resolver": {
				Resolver: &v1beta1.Resolution_GrpcResolver{
					GrpcResolver: &v1beta1.GrpcResolver{
						UpstreamRef: grpcUsRef,
						RequestTransform: &v1beta1.GrpcRequestTemplate{
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

		grpcDescriptorRegistry = v1beta1.GrpcDescriptorRegistry{
			DescriptorSet: &v1beta1.GrpcDescriptorRegistry_ProtoDescriptorBin{
				ProtoDescriptorBin: glootestpb.ProtoBytes,
			},
		}
		executor = v1beta1.Executor{
			Executor: &v1beta1.Executor_Local_{
				Local: &v1beta1.Executor_Local{
					Resolutions:         resolutions,
					EnableIntrospection: false,
				},
			},
		}

		return &v1beta1.GraphQLApi{
			Metadata: &core.Metadata{
				Name:      "gql",
				Namespace: "gloo-system",
			},
			Schema: &v1beta1.GraphQLApi_ExecutableSchema{
				ExecutableSchema: &v1beta1.ExecutableSchema{
					SchemaDefinition:       schema,
					Executor:               &executor,
					GrpcDescriptorRegistry: &grpcDescriptorRegistry,
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
					Action: &gloov1.Route_GraphqlApiRef{
						GraphqlApiRef: &core.ResourceRef{
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

		services.RunGlooGatewayUdsFdsOnPort(services.RunGlooGatewayOpts{Ctx: ctx, Cache: cache, LocalGlooPort: int32(testClients.GlooPort), What: what, Namespace: defaults.GlooSystem, License: GetValidGraphqlLicense()})
	})

	AfterEach(func() {
		cancel()
	})
	Context("With envoy", func() {
		var (
			envoyInstance                             *services.EnvoyInstance
			restUpstream, grpcUpstream, graphqlServer *v1helpers.TestUpstream
			envoyPort                                 = uint32(8080)
			query                                     string

			proxy      *gloov1.Proxy
			graphqlApi *v1beta1.GraphQLApi
		)

		var testRequestWithHeaders = func(result string, headers map[string]string) {
			var bodyStr string
			var respHeaders http.Header
			Eventually(func() (int, error) {
				client := http.DefaultClient
				reqUrl, err := url.Parse(fmt.Sprintf("http://%s:%d/testroute", "localhost", envoyPort))
				Expect(err).NotTo(HaveOccurred())
				resp, err := client.Do(&http.Request{
					Method: http.MethodPost,
					URL:    reqUrl,
					Body:   ioutil.NopCloser(strings.NewReader(query)),
					Header: map[string][]string{
						"bar":           []string{"bam"},
						"queryparamkey": []string{"queryparamvalue"},
					},
				})
				if err != nil {
					return 0, err
				}
				defer resp.Body.Close()
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					return 0, err
				}
				bodyStr = string(body)
				respHeaders = resp.Header
				return resp.StatusCode, nil
			}, "5s", "0.5s").Should(Equal(http.StatusOK))
			Expect(bodyStr).To(ContainSubstring(result))
			for k, v := range headers {
				Expect(respHeaders.Get(k)).To(Equal(v))
			}
		}
		var testRequest = func(result string) {
			testRequestWithHeaders(result, nil)
		}
		var testGetRequest = func(result string, includeQuery bool) {
			var bodyStr string
			Eventually(func() (int, error) {
				client := http.DefaultClient
				reqUrl, err := url.Parse(fmt.Sprintf("http://%s:%d/testroute", "localhost", envoyPort))
				Expect(err).NotTo(HaveOccurred())
				values := reqUrl.Query()
				if includeQuery {
					values.Add("query", query)
				}
				sum := sha256.Sum256([]byte(query))
				values.Add("extensions", fmt.Sprintf(`{"persistedQuery":{"version":1,"sha256Hash":"%x"}}`, sum))
				reqUrl.RawQuery = values.Encode()
				resp, err := client.Do(&http.Request{
					Method: http.MethodGet,
					URL:    reqUrl,
				})
				if err != nil {
					return 0, err
				}
				defer resp.Body.Close()
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					return 0, err
				}
				bodyStr = string(body)
				return resp.StatusCode, nil
			}, "5s", "0.5s").Should(Equal(http.StatusOK))
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

			query = `{"query":"{f:field1(intArg: 2, boolArg: true, floatArg: 9.99993, stringArg: \"this is a string arg\", mapArg: {a: 9}, listArg: [21,22,23]){simple}}","variables":{}}`

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

			graphqlServer = v1helpers.NewTestHttpUpstreamWithReply(ctx, envoyInstance.LocalAddr(), "{\"data\":{\"user\":\"JohnDoe\"}}")
			_, err = testClients.UpstreamClient.Write(graphqlServer.Upstream, clients.WriteOpts{})
			Expect(err).NotTo(HaveOccurred())
			helpers.EventuallyResourceAccepted(func() (resources.InputResource, error) {
				return testClients.UpstreamClient.Read(graphqlServer.Upstream.Metadata.Namespace,
					graphqlServer.Upstream.Metadata.Name, clients.ReadOpts{})
			})

			graphqlApi = getGraphQLApi(restUpstream.Upstream.Metadata.Ref(), grpcUpstream.Upstream.Metadata.Ref())
		})
		JustBeforeEach(func() {
			_, err := testClients.GraphQLApiClient.Write(graphqlApi, clients.WriteOpts{})
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
				testRequest(`{"data":{"f":{"simple":"foo"}}}`)
				Eventually(restUpstream.C).Should(Receive(PointTo(MatchFields(IgnoreExtras, Fields{
					"URL": PointTo(Equal(url.URL{
						Path: "/",
					})),
				}))))
			})
			Context("grpc resolver", func() {
				BeforeEach(func() {
					query = `{"query":"{f:field2{str}}"}`
				})
				Context("With artifact list", func() {
					BeforeEach(func() {
						grpcDescriptorRegistry = v1beta1.GrpcDescriptorRegistry{
							DescriptorSet: &v1beta1.GrpcDescriptorRegistry_ProtoRefsList{
								ProtoRefsList: &v1beta1.GrpcDescriptorRegistry_ProtoRefs{
									ConfigMapRefs: []*core.ResourceRef{
										&core.ResourceRef{
											Name:      "fake-artifact-one",
											Namespace: "gloo-system",
										},
									},
								},
							},
						}
						//create artifacts
						artifactOne := v1.Artifact{
							Metadata: &core.Metadata{
								Name:      "fake-artifact-one",
								Namespace: "gloo-system",
							},
							Data: map[string]string{
								"proto": base64.StdEncoding.EncodeToString(glootestpb.ProtoBytes),
							},
						}
						_, err := testClients.ArtifactClient.Write(&artifactOne, clients.WriteOpts{})
						Expect(err).ToNot(HaveOccurred())

					})
					AfterEach(func() {
						//delete artifacts
						err := testClients.ArtifactClient.Delete("gloo-system", "fake-artifact-one", clients.DeleteOpts{})
						Expect(err).ToNot(HaveOccurred())

					})

					It("resolves graphql queries to GRPC upstreams with artifacts", func() {
						testRequest(`{"data":{"f":{"str":"foo"}}}`)
						Eventually(grpcUpstream.C).Should(Receive(PointTo(MatchFields(IgnoreExtras, Fields{
							"GRPCRequest": PointTo(Equal(glootest.TestRequest{Str: "foo"})),
						}))))
					})
				})
				It("resolves graphql queries to GRPC upstreams", func() {
					testRequest(`{"data":{"f":{"str":"foo"}}}`)
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
					graphqlApi.GetExecutableSchema().GetExecutor().GetLocal().GetResolutions()["simple_resolver"].GetRestResolver().Request.Body = body
				})
				It("resolves graphql queries to REST upstreams with body", func() {
					testRequest(`{"data":{"f":{"simple":"foo"}}}`)
					Eventually(restUpstream.C).Should(Receive(PointTo(MatchFields(IgnoreExtras, Fields{
						"Body": Equal([]byte(`{"key1":"value1"}`)),
						"URL": PointTo(Equal(url.URL{
							Path: "/",
						})),
					}))))
				})
				Context("With remote executor", func() {
					BeforeEach(func() {
						graphqlApi.GetExecutableSchema().Executor = &v1beta1.Executor{
							Executor: &v1beta1.Executor_Remote_{
								Remote: &v1beta1.Executor_Remote{
									UpstreamRef: graphqlServer.Upstream.Metadata.Ref(),
								},
							},
						}
					})

					It("resolves graphql queries to REST upstreams with body via remote executor", func() {
						testRequest(`{"data":{"user":"JohnDoe"}}`)
						Eventually(graphqlServer.C).Should(Receive(PointTo(MatchFields(IgnoreExtras, Fields{
							"Body": Equal([]byte(query)),
							"URL": PointTo(Equal(url.URL{
								Path: "/",
							})),
						}))))
					})
				})
			})
			Context("with query params", func() {
				BeforeEach(func() {
					graphqlApi.GetExecutableSchema().GetExecutor().GetLocal().GetResolutions()["simple_resolver"].GetRestResolver().Request.QueryParams = map[string]string{
						"queryparam": "queryparamval",
					}
				})
				It("resolves graphql queries to REST upstreams with query params", func() {
					testRequest(`{"data":{"f":{"simple":"foo"}}}`)
					Eventually(restUpstream.C).Should(Receive(PointTo(MatchFields(IgnoreExtras, Fields{
						"URL": PointTo(Equal(url.URL{
							Path:     "/",
							RawQuery: "queryparam=queryparamval",
						})),
					}))))
				})
				Context("With remote executor and query params", func() {
					BeforeEach(func() {
						graphqlApi.GetExecutableSchema().Executor = &v1beta1.Executor{
							Executor: &v1beta1.Executor_Remote_{
								Remote: &v1beta1.Executor_Remote{
									UpstreamRef: graphqlServer.Upstream.Metadata.Ref(),
									QueryParams: map[string]string{
										"queryparam":    "queryparamval",
										"queryparamtwo": "{$headers.queryparamkey}",
									},
								},
							},
						}
					})

					It("resolves graphql queries to REST upstreams with query param via remote executor", func() {
						testRequest(`{"data":{"user":"JohnDoe"}}`)
						Eventually(graphqlServer.C).Should(Receive(PointTo(MatchFields(IgnoreExtras, Fields{
							"Body": Equal([]byte(query)),
							"URL": PointTo(Equal(url.URL{
								Path:     "/",
								RawQuery: "queryparam=queryparamval&queryparamtwo=queryparamvalue",
							})),
						}))))
					})
				})
			})
			Context("with headers", func() {
				BeforeEach(func() {
					graphqlApi.GetExecutableSchema().GetExecutor().GetLocal().GetResolutions()["simple_resolver"].GetRestResolver().Request.Headers = map[string]string{
						"header": "headerval",
					}
				})

				It("resolves graphql queries to REST upstreams with headers", func() {
					testRequest(`{"data":{"f":{"simple":"foo"}}}`)
					Eventually(restUpstream.C).Should(Receive(PointTo(MatchFields(IgnoreExtras, Fields{
						"URL": PointTo(Equal(url.URL{
							Path: "/",
						})),
						"Headers": HaveKeyWithValue("Header", []string{"headerval"}),
					}))))
				})
				Context("With remote executor and headers", func() {
					BeforeEach(func() {
						graphqlApi.GetExecutableSchema().Executor = &v1beta1.Executor{
							Executor: &v1beta1.Executor_Remote_{
								Remote: &v1beta1.Executor_Remote{
									UpstreamRef: graphqlServer.Upstream.Metadata.Ref(),
									Headers: map[string]string{
										"foo": "far",
										"boo": "{$headers.bar}",
									},
								},
							},
						}
					})

					It("resolves graphql queries to REST upstreams with headers via remote executor", func() {
						testRequest(`{"data":{"user":"JohnDoe"}}`)
						Eventually(graphqlServer.C).Should(Receive(PointTo(MatchFields(IgnoreExtras, Fields{
							"Body": Equal([]byte(query)),
							"URL": PointTo(Equal(url.URL{
								Path: "/",
							})),
							"Headers": And(
								HaveKeyWithValue("Foo", []string{"far"}),
								HaveKeyWithValue("Boo", []string{"bam"})),
						}))))
					})
				})
			})
			Context("allowlist", func() {
				Context("allowed", func() {
					BeforeEach(func() {
						graphqlApi.AllowedQueryHashes = []string{"075f4c9392a098f9b6d4e45fa87551d461edc7eedbc67b604bedc1cb9c854692"}
					})

					It("resolves allowed graphql queries", func() {
						testRequest(`{"data":{"f":{"simple":"foo"}}}`)
						Eventually(restUpstream.C).Should(Receive(PointTo(MatchFields(IgnoreExtras, Fields{
							"URL": PointTo(Equal(url.URL{
								Path: "/",
							})),
						}))))
					})
				})
				Context("disallowed", func() {
					BeforeEach(func() {
						graphqlApi.AllowedQueryHashes = []string{"hashnotfound"}
					})
					It("denies disallowed query hashes", func() {
						testRequest(`{"errors":[{"message":"hash 075f4c9392a098f9b6d4e45fa87551d461edc7eedbc67b604bedc1cb9c854692 not found in allowlist for query: '{f:field1(intArg: 2, boolArg: true, floatArg: 9.99993, stringArg: \"this is a string arg\", mapArg: {a: 9}, listArg: [21,22,23]){simple}}'"}]}`)
					})
				})
			})
			Context("persisted queries", func() {
				BeforeEach(func() {
					query = `{__typename}`
					graphqlApi.PersistedQueryCacheConfig = &v1beta1.PersistedQueryCacheConfig{CacheSize: 10}
				})

				It("happy path", func() {
					testGetRequest(`{"errors":[{"message":"persisted query not found: sha256 ecf4edb46db40b5132295c0291d62fb65d6759a9eedfa4d5d612dd5ec54a6b38"}]}`, false)

					testGetRequest(`{"data":{"__typename":"Query"}}`, true)

					// make same initial request, should now be cache hit in automatic persisted query cache
					testGetRequest(`{"data":{"__typename":"Query"}}`, false)
				})
			})
			Context("response setters and cache control", func() {

				BeforeEach(func() {
					query = `
{
  "query":"{f:field3{simple setme}}"
}`
				})
				Context("cache control", func() {
					BeforeEach(func() {
						graphqlApi.GetExecutableSchema().GetExecutor().GetLocal().GetResolutions()["simple_resolver"].GetRestResolver().Response = &v1beta1.ResponseTemplate{
							Setters: map[string]string{
								"setme": "{$body.simple}",
							},
						}
					})

					It("sets cache control header and simple field on response", func() {
						testRequestWithHeaders(`{"data":{"f":{"simple":"foo","setme":"foo"}}}`,
							map[string]string{"Cache-Control": "private, max-age=60"})
						Eventually(restUpstream.C).Should(Receive(PointTo(MatchFields(IgnoreExtras, Fields{
							"URL": PointTo(Equal(url.URL{
								Path: "/",
							})),
						}))))
					})
				})
				Context("response template", func() {
					BeforeEach(func() {
						graphqlApi.GetExecutableSchema().GetExecutor().GetLocal().GetResolutions()["simple_resolver"].GetRestResolver().Response = &v1beta1.ResponseTemplate{
							Setters: map[string]string{
								"setme": "abc {$body.simple} 123",
							},
						}
					})

					It("sets fields on response", func() {
						testRequest(`{"data":{"f":{"simple":"foo","setme":"abc foo 123"}}}`)
						Eventually(restUpstream.C).Should(Receive(PointTo(MatchFields(IgnoreExtras, Fields{
							"URL": PointTo(Equal(url.URL{
								Path: "/",
							})),
						}))))
					})
				})
			})
		})
	})
})

func GetValidGraphqlLicense() *model.License {
	return &model.License{
		IssuedAt:    time.Now(),
		ExpiresAt:   time.Now().Add(time.Minute * 100),
		LicenseType: model.LicenseType_Trial,
		Product:     model.Product_Gloo,
		AddOns: model.AddOns{
			GraphQL: true,
		},
	}
}
