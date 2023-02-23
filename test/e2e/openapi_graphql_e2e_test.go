package e2e_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	. "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/graphql/v1beta1"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/graphql"
	"github.com/solo-io/solo-projects/test/v1helpers"

	"github.com/fgrosse/zaptest"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/core/matchers"
	"github.com/solo-io/gloo/test/helpers"
	"github.com/solo-io/go-utils/contextutils"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/memory"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-projects/test/services"
)

var _ = Describe("Graphql E2E test", func() {

	var (
		ctx         context.Context
		cancel      context.CancelFunc
		testClients services.TestClients
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
	Context("With envoy", func() {
		var (
			envoyInstance *services.EnvoyInstance
			testUpstream1 *v1helpers.TestUpstream
			envoyPort     = uint32(8080)

			proxy *gloov1.Proxy
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
													Name:      testUpstream1.Upstream.Metadata.Name,
													Namespace: testUpstream1.Upstream.Metadata.Namespace,
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
			router := NewOpenapiBackend()
			testUpstream1 = v1helpers.NewTestHttpUpstreamWithHandler(ctx, envoyInstance.LocalAddr(), func(w http.ResponseWriter, r *http.Request) bool {
				router.ServeHTTP(w, r)
				return false
			})
			_, err = testClients.UpstreamClient.Write(testUpstream1.Upstream, clients.WriteOpts{})
			Expect(err).NotTo(HaveOccurred())

			helpers.EventuallyResourceAccepted(func() (resources.InputResource, error) {
				return testClients.UpstreamClient.Read(testUpstream1.Upstream.Metadata.Namespace,
					testUpstream1.Upstream.Metadata.Name, clients.ReadOpts{})
			})

		})

		AfterEach(func() {
			if envoyInstance != nil {
				envoyInstance.Clean()
			}
		})

		Context("test", func() {

			doQuery := func(query string) string {
				graphqlAddress := fmt.Sprintf("http://localhost:%d/graphql", envoyPort)
				req, err := http.NewRequest(http.MethodPost, graphqlAddress, strings.NewReader(query))
				ExpectWithOffset(1, err).NotTo(HaveOccurred())
				var bodyStr string
				EventuallyWithOffset(1, func() (int, error) {
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

			It("discovers openapi schema and translates to graphql schema", func() {
				By("Ensure that discovery discovers openapi", func() {
					Eventually(func() (string, error) {
						m := testUpstream1.Upstream.Metadata
						u, err := testClients.UpstreamClient.Read(m.GetNamespace(), m.GetName(), clients.ReadOpts{})
						if err != nil {
							return "", err
						}
						if u.GetStatic().GetServiceSpec() == nil {
							return "", fmt.Errorf("no service spec found on upstream %s.%s, %s", u.Metadata.Namespace, u.Metadata.Name, u.String())
						}

						return u.GetStatic().GetServiceSpec().GetRest().GetSwaggerInfo().GetUrl(), nil

					}, 20*time.Second, 1*time.Second).Should(Equal(fmt.Sprintf("http://%s/openapi.json", testUpstream1.Address)))

					By("creates graphql schema from discovered openapi", func() {

						m := testUpstream1.Upstream.GetMetadata()
						Eventually(func() (*GraphQLApi, error) {
							a, err := testClients.
								GraphQLApiClient.
								Read(m.Namespace, m.Name, clients.ReadOpts{})
							if err != nil {
								return nil, err
							}

							//resolutions, _ := yaml.Marshal(a.GetResolutions())
							//fmt.Printf("schema:\n %s\n resolvers:\n %s\n", a.GetSchema(), resolutions)

							return a, err

						}, 15*time.Second, 1*time.Second).ShouldNot(BeNil())

					})

					By("proxy gets accepted with graphql route referencing schema", func() {
						proxy = getProxy(envoyPort)
						configureProxy()
					})

					By("Basic Query", func() {
						// Queries existing pet
						basicQuery := `{"query":"query {  getPetById(petId: 1) {tags {name } name}}","variables":{}}`
						resBody := doQuery(basicQuery)
						Expect(resBody).To(Equal(`{"data":{"getPetById":{"tags":[{"name":"Tag 2"},{"name":"Tag 2"}],"name":"Pet 2"}}}`))
					})

					By("Basic Mutation", func() {
						// Creates Pet 3 via mutation and reads created pet's name
						graphqlQuery := `mutation Mutation {  addPet(petInput: {name: \"Pet 3\", id: 3, photoUrls:[\"url1\", \"url2\"], tags: []}){ name id }}`
						basicQuery := fmt.Sprintf(`{"query":"%s","variables":{},"operationName":"Mutation"}`, graphqlQuery)
						resBody := doQuery(basicQuery)
						Expect(resBody).To(Equal(`{"data":{"addPet":{"name":"Pet 3","id":3}}}`))
					})

					By("Update Pet with mutation", func() {
						updateMutation := `{"query":"mutation Mutation { updatePet(petInput: {id: 3, name: \"Cat 3\", photoUrls:[\"hi\"], tags: []})}"}`
						doQuery(updateMutation)
						basicQuery := `{"query":"query {  getPetById(petId: 3) { name photoUrls }}","variables":{}}`
						resBody := doQuery(basicQuery)
						Expect(resBody).To(Equal(`{"data":{"getPetById":{"name":"Cat 3","photoUrls":["hi"]}}}`))
					})

				})
			})

		})

	})
})

// Setup backend for testing
func NewOpenapiBackend() http.Handler {
	router := gin.Default()
	data := NewData()
	router.GET("/openapi.json", func(c *gin.Context) {
		c.String(http.StatusOK, graphql.GetFullJsonSchema())
	})

	v3 := router.Group("/v3")
	{
		v3.GET("/pet/:id", func(c *gin.Context) {
			idParam, _ := c.Params.Get("id")
			petToReturn, ok := data.Pets[idParam]
			if !ok {
				c.String(405, "Pet with id %s not found", idParam)
			}
			c.JSON(200, petToReturn)
		})
		v3.POST("/pet", func(c *gin.Context) {
			body, _ := ioutil.ReadAll(c.Request.Body)
			p := NewPet()
			err := json.Unmarshal(body, p)
			if err != nil {
				c.String(405, "Invalid Pet JSON: %s", err)
				return
			}
			data.Pets[strconv.Itoa(p.ID)] = p
			c.JSON(200, p)
		})
		v3.PUT("/pet", func(c *gin.Context) {
			body, _ := ioutil.ReadAll(c.Request.Body)
			petUpdateMsg := NewPet()
			err := json.Unmarshal(body, petUpdateMsg)
			if err != nil {
				c.String(405, "Invalid Pet JSON input")
			}
			p, ok := data.Pets[strconv.Itoa(petUpdateMsg.ID)]
			if !ok {
				c.String(405, "Cannot update pet with id %d because it does not exist", petUpdateMsg.ID)
			}
			if petUpdateMsg.Name != "" {
				p.Name = petUpdateMsg.Name
			}
			if petUpdateMsg.PhotoURLs != nil {
				p.PhotoURLs = petUpdateMsg.PhotoURLs
			}

		})
	}
	return router
}

type Data struct {
	Pets map[string]*Pet
}

func NewPet() *Pet {
	return &Pet{
		PhotoURLs: []string{""},
	}
}

func NewData() *Data {
	return &Data{
		Pets: map[string]*Pet{
			"0": {
				ID:        0,
				Name:      "Pet 0",
				PhotoURLs: []string{"pet 0 url"},
				Tags: []Tag{{
					ID:   0,
					Name: "Tag 1",
				},
					{
						ID:   1,
						Name: "Tag 2",
					}},
			},
			"1": {
				ID:        1,
				Name:      "Pet 2",
				PhotoURLs: []string{"pet 1 url"},
				Tags: []Tag{{
					ID:   0,
					Name: "Tag 2",
				},
					{
						ID:   1,
						Name: "Tag 2",
					}},
			},
		},
	}
}

// Tag the tag model
type Tag struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Pet the pet model
type Pet struct {
	ID        int      `json:"id"`
	Name      string   `json:"name"`
	PhotoURLs []string `json:"photoUrls,omitempty"`
	Tags      []Tag    `json:"tags,omitempty"`
}
