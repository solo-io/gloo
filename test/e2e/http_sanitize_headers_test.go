package e2e_test

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	"github.com/rotisserie/eris"

	"github.com/fgrosse/zaptest"
	"github.com/golang/protobuf/ptypes/wrappers"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	extauthv1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/extauth/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/defaults"
	"github.com/solo-io/gloo/projects/gloo/pkg/translator"
	"github.com/solo-io/go-utils/contextutils"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/memory"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
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

		testClients = services.GetTestClients(ctx, cache)
		testClients.GlooPort = int(services.AllocateGlooPort())

		what := services.What{
			DisableGateway: true,
			DisableFds:     true,
			DisableUds:     true,
		}

		services.RunGlooGatewayUdsFdsOnPort(services.RunGlooGatewayOpts{Ctx: ctx, Cache: cache, LocalGlooPort: int32(testClients.GlooPort), What: what, Namespace: "gloo-system", Settings: &gloov1.Settings{}})
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
			if err != nil {
				return 0, err
			}
			defer response.Body.Close()
			_, _ = io.ReadAll(response.Body)
			return response.StatusCode, err
		}, 20*time.Second, 1*time.Second).Should(Equal(200))

	}

	AfterEach(func() {
		if envoyInstance != nil {
			_ = envoyInstance.Clean()
		}
		cancel()
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
				response, err := http.DefaultClient.Do(request)
				if err != nil {
					return 0
				}
				defer response.Body.Close()
				_, _ = io.ReadAll(response.Body)
				return response.StatusCode
			}, 10*time.Second, 1*time.Second).Should(Equal(200))

			select {
			case received := <-testUpstream.C:
				Expect(received.Headers.Get("cluster-header-name")).To(Equal(upstreamName))
			case <-time.After(time.Second * 5):
				Fail("request didn't make it upstream")
			}

		})

		It("sanitizes downstream http header for cluster_header when told to do so", func() {
			setupProxy(true)
			us := testUpstream.Upstream
			upstreamName := translator.UpstreamToClusterName(us.Metadata.Ref())

			// Create a regular request
			request, err := http.NewRequest("GET", fmt.Sprintf("http://localhost:%d/", envoyPort), nil)
			Expect(err).NotTo(HaveOccurred())

			// TODO(kdorosh) it appears sanitizer filter has been broken a while. to simulate it working, I'm setting the wrong header here
			// see https://github.com/solo-io/solo-projects/issues/3507
			// appears this test was working because listener wasn't served (bug) but this was masking another bug in cluster header routing...
			request.Header.Add("cluster-header-name-FIXME", upstreamName)

			// Check that the request is not sent upstream because the cluster_header is sanitized, and the request has no upstream destination
			Eventually(func() error {
				response, err := http.DefaultClient.Do(request)
				if err != nil {
					return err
				}
				defer response.Body.Close()
				body, err := ioutil.ReadAll(response.Body)
				if err != nil {
					return err
				}
				if response.StatusCode != 503 {
					// we want 503 from envoy, sample reply:
					//
					//[2022-04-11 21:59:37.552][23][debug][http] [external/envoy/source/common/http/filter_manager.cc:953] [C2][S1853545153833015884] Sending local reply with details cluster_not_found
					//[2022-04-11 21:59:37.552][23][debug][http] [external/envoy/source/common/http/conn_manager_impl.cc:1467] [C2][S1853545153833015884] encoding headers via codec (end_stream=true):
					//':status', '503'
					//'date', 'Mon, 11 Apr 2022 21:59:37 GMT'
					//'server', 'envoy'
					return eris.Errorf("bad status code: %v (%v)", response.StatusCode, string(body))
				}
				return nil
			}, 10*time.Second, 1*time.Second).ShouldNot(HaveOccurred())

			select {
			case received := <-testUpstream.C:
				Fail(fmt.Sprintf("request received upstream %v", received))
			case <-time.After(time.Second * 5):
				// hooray! request did not make it upstream
			}

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
				Extauth: &extauthv1.Settings{ClearRouteCache: true,
					TransportApiVersion: extauthv1.Settings_V3,
				},
				SanitizeClusterHeader: &wrappers.BoolValue{
					Value: true,
				},
			}
		}
		return &gloov1.HttpListenerOptions{
			SanitizeClusterHeader: &wrappers.BoolValue{
				Value: false,
			},
		}
	}

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
					VirtualHosts: vhosts,
					Options:      getHttpListenerOptions(),
				},
			},
		}},
	}

	return p
}
