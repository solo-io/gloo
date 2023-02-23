package e2e_test

import (
	"context"
	"fmt"
	"net"
	"strings"

	"github.com/solo-io/gloo/test/helpers"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"

	_ "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/router/v3"

	envoy_admin_v3 "github.com/envoyproxy/go-control-plane/envoy/admin/v3"
	envoy_listener "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	_ "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/router/v3"
	envoy_hcm "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/network/http_connection_manager/v3"
	"github.com/fgrosse/zaptest"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/ptypes"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/extensions/proxylatency"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/defaults"
	"github.com/solo-io/go-utils/contextutils"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/memory"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	extauthrunner "github.com/solo-io/solo-projects/projects/extauth/pkg/runner"
	"github.com/solo-io/solo-projects/test/services"
	"github.com/solo-io/solo-projects/test/v1helpers"
)

var _ = Describe("Proxy latency", func() {

	var (
		ctx         context.Context
		cancel      context.CancelFunc
		testClients services.TestClients
		settings    extauthrunner.Settings
		cache       memory.InMemoryResourceCache
	)

	BeforeEach(func() {

		logger := zaptest.LoggerWriter(GinkgoWriter)
		contextutils.SetFallbackLogger(logger.Sugar())

		ctx, cancel = context.WithCancel(context.Background())
		cache = memory.NewInMemoryResourceCache()

		testClients = services.GetTestClients(ctx, cache)
		testClients.GlooPort = int(services.AllocateGlooPort())

		// Initialize settings for extauthrunner
		settings.GlooAddress = fmt.Sprintf("localhost:%d", testClients.GlooPort)
		settings.ExtAuthSettings.HealthCheckHttpPath = "/healthcheck"
		settings.ExtAuthSettings.HealthCheckHttpPort = int(services.AllocateGlooPort())

		what := services.What{
			DisableGateway: true,
			DisableUds:     true,
			DisableFds:     true,
		}

		services.RunGlooGatewayUdsFdsOnPort(services.RunGlooGatewayOpts{Ctx: ctx, Cache: cache, LocalGlooPort: int32(testClients.GlooPort), What: what, Namespace: defaults.GlooSystem})
		go func(testCtx context.Context) {
			defer GinkgoRecover()
			err := extauthrunner.RunWithSettings(testCtx, settings)
			if testCtx.Err() == nil {
				Expect(err).NotTo(HaveOccurred())
			}
		}(ctx)
	})

	AfterEach(func() {
		cancel()
	})

	Context("With envoy", func() {

		var (
			envoyInstance *services.EnvoyInstance
			testUpstream  *v1helpers.TestUpstream
			envoyPort     = uint32(8080)
		)

		BeforeEach(func() {
			var err error
			envoyInstance, err = envoyFactory.NewEnvoyInstance()
			Expect(err).NotTo(HaveOccurred())

			err = envoyInstance.Run(testClients.GlooPort)
			Expect(err).NotTo(HaveOccurred())

			testUpstream = v1helpers.NewTestHttpUpstream(ctx, envoyInstance.LocalAddr())

			var opts clients.WriteOpts
			up := testUpstream.Upstream
			_, err = testClients.UpstreamClient.Write(up, opts)
			Expect(err).NotTo(HaveOccurred())
		})

		AfterEach(func() {
			_ = envoyInstance.Clean()
		})

		Context("proxy latency", func() {

			BeforeEach(func() {
				proxy := getProxyLatencyProxy(envoyPort, testUpstream.Upstream.Metadata.Ref())

				_, err := testClients.ProxyClient.Write(proxy, clients.WriteOpts{Ctx: ctx})
				Expect(err).NotTo(HaveOccurred())

				helpers.EventuallyResourceAccepted(func() (resources.InputResource, error) {
					return testClients.ProxyClient.Read(proxy.Metadata.Namespace, proxy.Metadata.Name, clients.ReadOpts{})
				})
			})

			It("should be last", func() {
				Eventually(func() (string, error) {
					filters, err := getFilters(envoyInstance)
					if err != nil {
						return "", err
					}
					// last filter is router, so get the one before that:
					if len(filters) >= 2 {
						return filters[len(filters)-2], nil
					}
					return "", nil
				}, "5s", "0.5s").Should(ContainSubstring("proxy_latency"))
			})
		})
	})

})

func getFilters(envoyInstance *services.EnvoyInstance) ([]string, error) {
	resp, err := envoyInstance.GetConfigDump()
	if err != nil {
		return nil, err
	}

	jsonpbMarshaler := &jsonpb.Unmarshaler{
		// Ever since upgrading the go-control-plane to v0.10.1 this test fails with the following error:
		// unknown field \"hidden_envoy_deprecated_build_version\" in envoy.config.core.v3.Node"
		// Set AllowUnknownFields to true to get around this
		AllowUnknownFields: true,
	}

	var cfgDump envoy_admin_v3.ConfigDump
	if err = jsonpbMarshaler.Unmarshal(resp.Body, &cfgDump); err != nil {
		return nil, err
	}
	if err = resp.Body.Close(); err != nil {
		return nil, err
	}

	for _, cfg := range cfgDump.GetConfigs() {
		if strings.Contains(cfg.GetTypeUrl(), "ListenersConfigDump") {
			var listeners envoy_admin_v3.ListenersConfigDump
			err = ptypes.UnmarshalAny(cfg, &listeners)
			if err != nil {
				return nil, err
			}
			if len(listeners.GetDynamicListeners()) != 1 {
				return nil, fmt.Errorf("listener not found")
			}
			anyListener := listeners.GetDynamicListeners()[0].GetActiveState().GetListener()
			if anyListener == nil {
				return nil, fmt.Errorf("listener not active")
			}

			var listener envoy_listener.Listener
			err = ptypes.UnmarshalAny(anyListener, &listener)
			if err != nil {
				return nil, err
			}

			anyHcm := listener.GetFilterChains()[0].GetFilters()[0]
			var hcm envoy_hcm.HttpConnectionManager
			err = ptypes.UnmarshalAny(anyHcm.GetTypedConfig(), &hcm)
			if err != nil {
				return nil, err
			}
			httpFilters := hcm.GetHttpFilters()
			var names []string
			for _, filter := range httpFilters {
				names = append(names, filter.GetName())
			}
			return names, nil
		}
	}
	return nil, fmt.Errorf("config not found")
}

func getProxyLatencyProxy(envoyPort uint32, upstream *core.ResourceRef) *gloov1.Proxy {
	var vhosts []*gloov1.VirtualHost

	vhost := &gloov1.VirtualHost{
		Name:    "gloo-system.virt1",
		Domains: []string{"*"},
		Routes: []*gloov1.Route{{
			Action: &gloov1.Route_RouteAction{
				RouteAction: &gloov1.RouteAction{
					Destination: &gloov1.RouteAction_Single{
						Single: &gloov1.Destination{
							DestinationType: &gloov1.Destination_Upstream{
								Upstream: upstream,
							},
						},
					},
				},
			},
		}},
	}

	vhosts = append(vhosts, vhost)

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
					Options: &gloov1.HttpListenerOptions{
						ProxyLatency: &proxylatency.ProxyLatency{
							MeasureRequestInternally: true,
						},
					},
					VirtualHosts: vhosts,
				},
			},
		}},
	}

	return p
}
