package e2e_test

import (
	"context"
	json2 "encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/solo-io/solo-projects/projects/gloo/pkg/utils/graphql/translation"

	"github.com/fgrosse/zaptest"
	structpb "github.com/golang/protobuf/ptypes/struct"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/core/matchers"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/graphql/v1beta1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/cors"
	"github.com/solo-io/gloo/projects/gloo/pkg/defaults"
	"github.com/solo-io/gloo/test/helpers"
	"github.com/solo-io/go-utils/contextutils"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/memory"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-projects/test/services"
)

var _ = Describe("graphql stitching", func() {

	var (
		ctx         context.Context
		cancel      context.CancelFunc
		testClients services.TestClients
	)

	productSchemaDef := `
type User {
	username: String
}

type Product {
	id: Int
	seller: User
}

type Query {
  products: Product @resolve(name: "product_resolver")
}
`

	productGqlApi := &v1beta1.GraphQLApi{
		Metadata: &core.Metadata{
			Name:      "product-gql",
			Namespace: "gloo-system",
		},
		Schema: &v1beta1.GraphQLApi_ExecutableSchema{
			ExecutableSchema: &v1beta1.ExecutableSchema{
				SchemaDefinition: productSchemaDef,
				Executor: &v1beta1.Executor{
					Executor: &v1beta1.Executor_Local_{
						Local: &v1beta1.Executor_Local{
							Resolutions: map[string]*v1beta1.Resolution{
								"product_resolver": {
									Resolver: &v1beta1.Resolution_MockResolver{
										MockResolver: &v1beta1.MockResolver{
											Response: &v1beta1.MockResolver_SyncResponse{
												SyncResponse: JsonToStructPbValue(`{"id": 1, "seller": {"username": "user1"}}`),
											},
										},
									},
								},
							},
							EnableIntrospection: true,
						},
					},
				},
			},
		},
	}

	userSchemaDef := `
type User {
	username: String
	firstName: String
	lastName: String
}

type Query {
  user: User @resolve(name: "user_resolver")
}
`

	userGqlApi := &v1beta1.GraphQLApi{
		Metadata: &core.Metadata{
			Name:      "users-gql",
			Namespace: "gloo-system",
		},
		Schema: &v1beta1.GraphQLApi_ExecutableSchema{
			ExecutableSchema: &v1beta1.ExecutableSchema{
				SchemaDefinition: userSchemaDef,
				Executor: &v1beta1.Executor{
					Executor: &v1beta1.Executor_Local_{
						Local: &v1beta1.Executor_Local{
							Resolutions: map[string]*v1beta1.Resolution{
								"user_resolver": {
									Resolver: &v1beta1.Resolution_MockResolver{
										MockResolver: &v1beta1.MockResolver{
											Response: &v1beta1.MockResolver_SyncResponse{
												SyncResponse: JsonToStructPbValue(`{"username": "user1", "firstName": "User", "lastName": "One"}`),
											},
										},
									},
								},
							},
							EnableIntrospection: true,
						},
					},
				},
			},
		},
	}

	stitchedGqlApi := &v1beta1.GraphQLApi{
		Metadata: &core.Metadata{
			Name:      "stitched-gql",
			Namespace: "gloo-system",
		},
		Schema: &v1beta1.GraphQLApi_StitchedSchema{
			StitchedSchema: &v1beta1.StitchedSchema{
				Subschemas: []*v1beta1.StitchedSchema_SubschemaConfig{
					// Stitch product api
					{
						Name:      "product-gql",
						Namespace: "gloo-system",
					},
					// Stitch Users api with type merge configuration for Users type
					{
						Name:      "users-gql",
						Namespace: "gloo-system",
						TypeMerge: map[string]*v1beta1.StitchedSchema_SubschemaConfig_TypeMergeConfig{
							"User": {
								QueryName:    "user",
								SelectionSet: "{ username }",
								Args: map[string]string{
									"username": "username",
								},
							},
						},
					},
				},
			},
		},
	}

	var getProxy = func(envoyPort uint32) *gloov1.Proxy {

		var vhosts []*gloov1.VirtualHost

		vhost := &gloov1.VirtualHost{
			Name:    "gloo-system.virt1",
			Domains: []string{"*"},
			Options: &gloov1.VirtualHostOptions{
				// This isn't needed for this end-to-end test to work, but is useful when
				// debugging using a graphql explorer IDE like apollo sandbox or graphql playground
				Cors: &cors.CorsPolicy{
					AllowCredentials: true,
					AllowHeaders:     []string{"content-type", "x-apollo-tracing"},
					AllowMethods:     []string{"POST"},
					AllowOriginRegex: []string{"\\*"},
				},
			},
			Routes: []*gloov1.Route{
				{
					Matchers: []*matchers.Matcher{{
						PathSpecifier: &matchers.Matcher_Prefix{
							Prefix: "/testroute",
						},
					}},
					Action: &gloov1.Route_GraphqlApiRef{
						GraphqlApiRef: &core.ResourceRef{
							Name:      "stitched-gql",
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
			envoyInstance *services.EnvoyInstance
			envoyPort     = uint32(8080)
			query         string

			proxy *gloov1.Proxy
		)

		var testRequestWithRespAssertions = func(result string, f func(resp *http.Response)) {
			var resp *http.Response
			EventuallyWithOffset(1, func() (int, error) {
				client := http.DefaultClient
				reqUrl, err := url.Parse(fmt.Sprintf("http://%s:%d/testroute", "localhost", envoyPort))
				ExpectWithOffset(1, err).NotTo(HaveOccurred())
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
			if f != nil {
				f(resp)
			}
			ExpectWithOffset(1, err).NotTo(HaveOccurred())
			ExpectWithOffset(1, bodyStr).To(ContainSubstring(result))
		}

		var configureProxy = func() {
			ExpectWithOffset(1, proxy).NotTo(BeNil())
			_, err := testClients.ProxyClient.Write(proxy, clients.WriteOpts{})
			ExpectWithOffset(1, err).NotTo(HaveOccurred())

			helpers.EventuallyResourceAccepted(func() (resources.InputResource, error) {
				test, err := testClients.ProxyClient.Read(proxy.Metadata.Namespace, proxy.Metadata.Name, clients.ReadOpts{})
				return test, err
			})
		}

		BeforeEach(func() {
			var err error
			envoyInstance, err = envoyFactory.NewEnvoyInstance()
			Expect(err).NotTo(HaveOccurred())

			err = envoyInstance.Run(testClients.GlooPort)
			Expect(err).NotTo(HaveOccurred())
			err = os.Setenv(translation.GraphqlJsRootEnvVar, "../../projects/gloo/pkg/plugins/graphql/js/")
			Expect(err).NotTo(HaveOccurred())
			err = os.Setenv(translation.GraphqlProtoRootEnvVar, "../../projects/ui/src/proto/")
			Expect(err).NotTo(HaveOccurred())

			query = `
{"query":" {products { id seller { username firstName lastName}}}"}`

		})
		JustBeforeEach(func() {
			_, err := testClients.GraphQLApiClient.Write(userGqlApi, clients.WriteOpts{})
			Expect(err).NotTo(HaveOccurred())

			_, err = testClients.GraphQLApiClient.Write(productGqlApi, clients.WriteOpts{})
			Expect(err).NotTo(HaveOccurred())

			_, err = testClients.GraphQLApiClient.Write(stitchedGqlApi, clients.WriteOpts{})
			Expect(err).NotTo(HaveOccurred())

			proxy = getProxy(envoyPort)
			configureProxy()
		})

		AfterEach(func() {
			if envoyInstance != nil {
				envoyInstance.Clean()
			}
			Expect(os.Unsetenv(translation.GraphqlProtoRootEnvVar)).NotTo(HaveOccurred())
			Expect(os.Unsetenv(translation.GraphqlJsRootEnvVar)).NotTo(HaveOccurred())
		})

		Context("request to stitched schema", func() {

			It("stitches delegated responses from subschemas to a stitched response", func() {
				testRequestWithRespAssertions(`{"data":{"products":{"id":1,"seller":{"username":"user1","firstName":"User","lastName":"One"}}}}`, nil)
			})
		})
	})
})

func JsonToStructPbValue(js string) *structpb.Value {
	ret := &structpb.Value{}
	err := json2.Unmarshal([]byte(js), ret)
	ExpectWithOffset(1, err).NotTo(HaveOccurred())
	return ret
}
