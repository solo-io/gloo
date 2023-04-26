package e2e_test

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	glooV1helpers "github.com/solo-io/gloo/test/v1helpers"
	"github.com/solo-io/solo-projects/test/v1helpers"

	"github.com/fgrosse/zaptest"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	todo "github.com/solo-io/gloo-graphql-example/code/todo-app/server"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/core/matchers"
	. "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/graphql/v1beta1"
	plugins "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options"
	"github.com/solo-io/gloo/test/helpers"
	"github.com/solo-io/go-utils/contextutils"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/memory"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-projects/test/services"
)

var _ = Describe("Graphql Remote and gRPC E2E test", func() {

	const graphqlPort = "8280"

	var (
		ctx         context.Context
		cancel      context.CancelFunc
		testClients services.TestClients
		server      *todo.TodoApp
	)

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
			DisableFds:     false,
		}

		services.RunGlooGatewayUdsFdsOnPort(services.RunGlooGatewayOpts{Ctx: ctx, Cache: cache, LocalGlooPort: int32(testClients.GlooPort), What: what, Namespace: "gloo-system", Settings: &gloov1.Settings{
			Discovery: &gloov1.Settings_DiscoveryOptions{
				FdsMode: gloov1.Settings_DiscoveryOptions_BLACKLIST,
			},
		},
			License: GetValidGraphqlLicense(),
		})
	})

	AfterEach(func() {
		cancel()
	})
	Context("finding a gRPC service", func() {
		var (
			envoyInstance *services.EnvoyInstance
			grpcUpstream  *glooV1helpers.TestUpstream
		)
		BeforeEach(func() {
			var err error
			envoyInstance, err = envoyFactory.NewEnvoyInstance()
			Expect(err).NotTo(HaveOccurred())

			err = envoyInstance.Run(testClients.GlooPort)
			Expect(err).NotTo(HaveOccurred())

			grpcUpstream = glooV1helpers.NewTestGRPCUpstream(ctx, envoyInstance.LocalAddr(), 1)
			_, err = testClients.UpstreamClient.Write(grpcUpstream.Upstream, clients.WriteOpts{})
			Expect(err).NotTo(HaveOccurred())
			helpers.EventuallyResourceAccepted(func() (resources.InputResource, error) {
				return testClients.UpstreamClient.Read(grpcUpstream.Upstream.Metadata.Namespace,
					grpcUpstream.Upstream.Metadata.Name, clients.ReadOpts{})
			})
		})

		AfterEach(func() {
			if envoyInstance != nil {
				envoyInstance.Clean()
			}
		})

		It("should discover the gRPC service", func() {
			Eventually(func() ([]byte, error) {
				m := grpcUpstream.Upstream.Metadata
				u, err := testClients.UpstreamClient.Read(m.GetNamespace(), m.GetName(), clients.ReadOpts{})
				if err != nil {
					return nil, err
				}
				grpcServiceSpec, ok := u.GetStatic().GetServiceSpec().GetPluginType().(*plugins.ServiceSpec_Grpc)
				if grpcServiceSpec == nil && !ok {
					return nil, fmt.Errorf("no grpc service spec found on upstream %s.%s, %s", u.Metadata.Namespace, u.Metadata.Name, u.String())
				}
				graphqlResource, err := testClients.
					GraphQLApiClient.
					Read(m.Namespace, m.Name, clients.ReadOpts{})
				if err != nil {
					return nil, err
				}
				return graphqlResource.GetExecutableSchema().GetGrpcDescriptorRegistry().GetProtoDescriptorBin(), nil
			}, 1*time.Minute, 1*time.Second).ShouldNot(BeNil())
		})
	})

	Context("finding a graphql remote service", func() {
		var (
			envoyInstance *services.EnvoyInstance
			testUpstream  *v1helpers.TestUpstream
			envoyPort     = uint32(8080)
			proxy         *gloov1.Proxy
		)

		var getProxy = func(envoyPort uint32) *gloov1.Proxy {

			p := &gloov1.Proxy{
				Metadata: &core.Metadata{
					Name:      "proxy",
					Namespace: "default",
				},
				Listeners: []*gloov1.Listener{{
					Name:        "listener",
					BindAddress: net.IPv4zero.String(),
					BindPort:    envoyPort,
					ListenerType: &gloov1.Listener_HttpListener{
						HttpListener: &gloov1.HttpListener{
							VirtualHosts: []*gloov1.VirtualHost{
								{
									Name:    "gloo-system.virt1",
									Domains: []string{"*"},
									Routes: []*gloov1.Route{
										{
											Matchers: []*matchers.Matcher{{
												PathSpecifier: &matchers.Matcher_Prefix{
													Prefix: "/graphql",
												},
											}},
											Action: &gloov1.Route_GraphqlApiRef{
												GraphqlApiRef: &core.ResourceRef{
													Name:      testUpstream.Upstream.Metadata.Name,
													Namespace: testUpstream.Upstream.Metadata.Namespace,
												},
											},
										},
									},
								},
							},
						},
					},
				}},
			}

			return p
		}

		var configureProxy = func() {
			ExpectWithOffset(1, proxy).NotTo(BeNil())
			_, err := testClients.ProxyClient.Write(proxy, clients.WriteOpts{})
			ExpectWithOffset(1, err).NotTo(HaveOccurred())

			helpers.EventuallyResourceAccepted(func() (resources.InputResource, error) {
				p, err := testClients.ProxyClient.Read(proxy.Metadata.Namespace, proxy.Metadata.Name, clients.ReadOpts{})
				return p, err
			})
			EventuallyWithOffset(1, func() (int, error) {
				resp, err := http.Get(fmt.Sprintf("http://%s:%d", "localhost", envoyInstance.AdminPort))
				if err != nil {
					return 0, err
				}
				defer resp.Body.Close()
				_, _ = io.ReadAll(resp.Body)
				return resp.StatusCode, nil
			}, "5s", "0.5s").Should(Equal(http.StatusOK))
		}

		BeforeEach(func() {
			var err error
			envoyInstance, err = envoyFactory.NewEnvoyInstance()
			Expect(err).NotTo(HaveOccurred())

			err = envoyInstance.Run(testClients.GlooPort)
			Expect(err).NotTo(HaveOccurred())

			server = todo.NewTodoServer(graphqlPort)
			errs, err := server.Start(ctx)
			Expect(err).NotTo(HaveOccurred())
			if errs != nil {
				go func() {
					select {
					case er := <-errs:
						fmt.Println(er.Error())
						break
					case <-ctx.Done():
						break
					}
				}()
			}

			port, err := strconv.Atoi(graphqlPort)
			Expect(err).ToNot(HaveOccurred())
			testUpstream = v1helpers.NewTestHttpUpstreamWithAddress(envoyInstance.LocalAddr(), uint32(port))

			_, err = testClients.UpstreamClient.Write(testUpstream.Upstream, clients.WriteOpts{})
			Expect(err).NotTo(HaveOccurred())

			helpers.EventuallyResourceAccepted(func() (resources.InputResource, error) {
				return testClients.UpstreamClient.Read(testUpstream.Upstream.Metadata.Namespace,
					testUpstream.Upstream.Metadata.Name, clients.ReadOpts{})
			})

		})

		AfterEach(func() {
			if envoyInstance != nil {
				envoyInstance.Clean()
			}
			err := server.Kill(ctx)
			Expect(err).ToNot(HaveOccurred())
		})

		Context("GraphQL Discovery Function", func() {

			doQuery := func(query string) string {
				var bodyStr string
				EventuallyWithOffset(1, func() (int, error) {
					graphqlAddress := fmt.Sprintf("http://localhost:%d/graphql", envoyPort)
					req, err := http.NewRequest(http.MethodPost, graphqlAddress, strings.NewReader(query))
					ExpectWithOffset(1, err).NotTo(HaveOccurred())
					res, err := http.DefaultClient.Do(req)
					if err != nil {
						return 0, err
					}
					defer res.Body.Close()
					body, err := io.ReadAll(res.Body)
					if err != nil {
						return 0, err
					}
					bodyStr = string(body)
					return res.StatusCode, err
				}, "10s", "1s").Should(Equal(http.StatusOK))
				return bodyStr
			}

			It("discovers graphql schema", func() {
				By("Ensure that discovery discovers graphql", func() {
					Eventually(func() (string, error) {
						m := testUpstream.Upstream.Metadata
						u, err := testClients.UpstreamClient.Read(m.GetNamespace(), m.GetName(), clients.ReadOpts{})
						if err != nil {
							return "", err
						}
						if u.GetStatic().GetServiceSpec() == nil {
							return "", fmt.Errorf("no service spec found on upstream %s.%s, %s", u.Metadata.Namespace, u.Metadata.Name, u.String())
						}

						url := u.GetStatic().GetServiceSpec().GetGraphql().GetEndpoint().GetUrl()
						return url, nil
					}, 1*time.Minute, 1*time.Second).Should(Equal(fmt.Sprintf("http://%s/graphql", testUpstream.Address)))

					By("creates graphql schema from discovered graphql server", func() {
						m := testUpstream.Upstream.GetMetadata()
						var graphqlResource *GraphQLApi
						Eventually(func() (*GraphQLApi, error) {
							graphqlResource, err := testClients.
								GraphQLApiClient.
								Read(m.Namespace, m.Name, clients.ReadOpts{})
							if err != nil {
								return nil, err
							}
							return graphqlResource, err

						}, 15*time.Second, 1*time.Second).ShouldNot(BeNil())
						schema := graphqlResource.GetExecutableSchema().GetSchemaDefinition()
						Expect(schema).ToNot(BeNil())
					})

					By("proxy gets accepted with graphql route referencing schema", func() {
						proxy = getProxy(envoyPort)
						configureProxy()
					})

					By("Basic Query", func() {
						// Queries existing todo
						basicQuery := `{"query":"{todo(id:\"b\"){id,text,done}}"}`
						resBody := doQuery(basicQuery)
						Expect(resBody).To(Equal(`{"data":{"todo":{"done":false,"id":"b","text":"This is the most important"}}}`))
					})

					By("Basic Mutation", func() {
						basicQuery := `{"query":"mutation _{createTodo(text:\"My new todo\"){id,text,done}}","variables":{},"operationName":"_"}`
						resBody := doQuery(basicQuery)
						Expect(resBody).To(ContainSubstring(`{"data":{"createTodo":{"done":false,"id":"`))
						Expect(resBody).To(ContainSubstring(`","text":"My new todo"}}}`))
					})

					By("Update Todo with mutation", func() {
						updateMutation := `{"query":"mutation _{updateTodo(id:\"b\",done:true){id,text,done}}"}`
						resBody := doQuery(updateMutation)
						Expect(resBody).To(Equal(`{"data":{"updateTodo":{"done":true,"id":"b","text":"This is the most important"}}}`))
					})

				})
			})
		})
	})
})
