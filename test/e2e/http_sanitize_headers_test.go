package e2e_test

import (
	"context"
	"fmt"
	"net/http"
	"time"

	extauthv1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/extauth/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/defaults"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"

	"github.com/fgrosse/zaptest"
	"github.com/gogo/protobuf/types"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/translator"
	"github.com/solo-io/go-utils/contextutils"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/memory"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/solo-io/solo-projects/test/services"
	"github.com/solo-io/solo-projects/test/v1helpers"
)

var _ = Describe("Http Sanitize Headers Local E2E", func() {

	var (
		ctx           context.Context
		cancel        context.CancelFunc
		testClients   services.TestClients
		envoyInstance *services.EnvoyInstance
		testUpstream  *v1helpers.TestUpstream
		envoyPort     uint32
	)

	BeforeEach(func() {

		logger := zaptest.LoggerWriter(GinkgoWriter)
		contextutils.SetFallbackLogger(logger.Sugar())

		ctx, cancel = context.WithCancel(context.Background())
		cache := memory.NewInMemoryResourceCache()

		testClients = services.GetTestClients(cache)
		testClients.GlooPort = int(services.AllocateGlooPort())

		what := services.What{
			DisableFds: true,
			DisableUds: true,
		}

		services.RunGlooGatewayUdsFdsOnPort(ctx, cache, int32(testClients.GlooPort), what, "gloo-system", nil, nil, &gloov1.Settings{})
	})

	setupProxy := func(headerSanitation bool) {
		var err error
		envoyInstance, err = envoyFactory.NewEnvoyInstance()
		Expect(err).NotTo(HaveOccurred())
		err = envoyInstance.Run(testClients.GlooPort)
		Expect(err).NotTo(HaveOccurred())

		testUpstream = v1helpers.NewTestHttpUpstream(ctx, envoyInstance.LocalAddr())

		envoyPort = defaults.HttpPort
		proxy := getProxyWithHeaderSanitation(envoyPort, headerSanitation)

		_, err = testClients.ProxyClient.Write(proxy, clients.WriteOpts{})
		Expect(err).NotTo(HaveOccurred())

		_, err = testClients.UpstreamClient.Write(testUpstream.Upstream, clients.WriteOpts{})
		Expect(err).NotTo(HaveOccurred())

		request, err := http.NewRequest("GET", fmt.Sprintf("http://localhost:%d", envoyInstance.AdminPort), nil)
		Expect(err).NotTo(HaveOccurred())
		client := &http.Client{}
		Eventually(func() (int, error) {
			response, err := client.Do(request)
			if response == nil {
				return 0, err
			}
			return response.StatusCode, err
		}, 20*time.Second, 1*time.Second).Should(Equal(200))

	}

	AfterEach(func() {
		cancel()
		if envoyInstance != nil {
			_ = envoyInstance.Clean()
		}
		// Wait till envoy is completely cleaned up
		request, err := http.NewRequest("GET", fmt.Sprintf("http://localhost:%d", envoyInstance.AdminPort), nil)
		Expect(err).NotTo(HaveOccurred())
		client := &http.Client{}
		Eventually(func() error {
			_, err := client.Do(request)
			return err
		}, 5*time.Second, 1*time.Second).Should(HaveOccurred())
	})

	Context("With envoy", func() {

		It("sanity check for cluster header destination", func() {
			setupProxy(false)
			us := testUpstream.Upstream
			upstreamName := translator.UpstreamToClusterName(us.Metadata.Ref())

			request, err := http.NewRequest("GET", fmt.Sprintf("http://localhost:%d/", envoyPort), nil)
			Expect(err).NotTo(HaveOccurred())
			request.Header.Add("cluster-header-name", upstreamName)

			Eventually(func() int {
				response, _ := http.DefaultClient.Do(request)
				if response == nil {
					return 0
				}
				return response.StatusCode
			}, 10*time.Second, 1*time.Second).Should(Equal(200))

		})

		It("sanitizes downstream http header for cluster_header when told to do so", func() {
			setupProxy(true)
			us := testUpstream.Upstream
			upstreamName := translator.UpstreamToClusterName(us.Metadata.Ref())

			// Create a regular request
			request, err := http.NewRequest("GET", fmt.Sprintf("http://localhost:%d/", envoyPort), nil)
			Expect(err).NotTo(HaveOccurred())
			request.Header.Add("cluster-header-name", upstreamName)

			// Check that the request times out because the cluster_header is sanitized, and request has no upstream destination
			Eventually(func() error {
				_, err := http.DefaultClient.Do(request)
				return err
			}, 5*time.Second, 1*time.Second).Should(HaveOccurred())

		})

	})

})

func getProxyWithHeaderSanitation(envoyPort uint32, headerSanitation bool) *gloov1.Proxy {
	var vhosts []*gloov1.VirtualHost

	vhost := &gloov1.VirtualHost{
		Name:    "virt1",
		Domains: []string{"*"},
		Routes: []*gloov1.Route{{
			Action: &gloov1.Route_RouteAction{
				RouteAction: &gloov1.RouteAction{
					Destination: &gloov1.RouteAction_ClusterHeader{
						ClusterHeader: "cluster-header-name",
					},
				},
			}}},
	}

	vhosts = append(vhosts, vhost)

	getHttpListenerOptions := func() *gloov1.HttpListenerOptions {
		if headerSanitation {
			return &gloov1.HttpListenerOptions{
				Extauth: &extauthv1.Settings{ClearRouteCache: true},
				SanitizeClusterHeader: &types.BoolValue{
					Value: true,
				},
			}
		}
		return &gloov1.HttpListenerOptions{
			SanitizeClusterHeader: &types.BoolValue{
				Value: false,
			},
		}
	}

	p := &gloov1.Proxy{
		Metadata: core.Metadata{
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
					Options:      getHttpListenerOptions(),
				},
			},
		}},
	}

	return p
}
