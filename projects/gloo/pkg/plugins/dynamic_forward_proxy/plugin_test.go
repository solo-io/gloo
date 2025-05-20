package dynamic_forward_proxy_test

import (
	envoy_config_cluster_v3 "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	envoy_config_core_v3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	envoy_extensions_filters_http_dynamic_forward_proxy_v3 "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/dynamic_forward_proxy/v3"
	envoy_extensions_network_dns_resolver_cares_v3 "github.com/envoyproxy/go-control-plane/envoy/extensions/network/dns_resolver/cares/v3"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/duration"
	"github.com/golang/protobuf/ptypes/wrappers"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	v32 "github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/config/core/v3"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	v1_circuitbreaker "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/circuit_breaker"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/gloosnapshot"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/dynamic_forward_proxy"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/ssl"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins"
	. "github.com/solo-io/gloo/projects/gloo/pkg/plugins/dynamic_forward_proxy"
	"github.com/solo-io/gloo/projects/gloo/pkg/utils"
	"github.com/solo-io/solo-kit/pkg/api/v2/reporter"
	"github.com/solo-io/solo-kit/test/matchers"
)

var _ = Describe("dynamic forward proxy plugin", func() {

	var (
		params     plugins.Params
		initParams plugins.InitParams
		listener   *v1.HttpListener
	)

	BeforeEach(func() {
		params = plugins.Params{}
		initParams = plugins.InitParams{}
		listener = &v1.HttpListener{}
	})

	It("does not configure DFP filter if not needed", func() {
		p := NewPlugin()
		filters, err := p.HttpFilters(params, listener)
		Expect(err).NotTo(HaveOccurred())
		Expect(filters).To(BeEmpty())
	})

	Context("sane defaults", func() {

		BeforeEach(func() {
			listener.Options = &v1.HttpListenerOptions{
				DynamicForwardProxy: &dynamic_forward_proxy.FilterConfig{}, // pick up system defaults to resolve DNS
			}
		})

		It("uses sane defaults with empty http filter", func() {
			p := NewPlugin()
			p.Init(initParams)

			filters, err := p.HttpFilters(params, listener)
			Expect(err).NotTo(HaveOccurred())
			Expect(filters).To(HaveLen(1))

			filterCfg := &envoy_extensions_filters_http_dynamic_forward_proxy_v3.FilterConfig{}
			goTypedConfig := filters[0].Filter.GetTypedConfig()
			err = ptypes.UnmarshalAny(goTypedConfig, filterCfg)
			Expect(err).NotTo(HaveOccurred())

			Expect(filterCfg.GetDnsCacheConfig().GetDnsLookupFamily()).To(Equal(envoy_config_cluster_v3.Cluster_V4_PREFERRED))
		})
	})

	Context("translates provided config", func() {

		BeforeEach(func() {
			listener.Options = &v1.HttpListenerOptions{
				DynamicForwardProxy: &dynamic_forward_proxy.FilterConfig{
					DnsCacheConfig: &dynamic_forward_proxy.DnsCacheConfig{
						DnsLookupFamily: dynamic_forward_proxy.DnsLookupFamily_V4_ONLY,
						DnsRefreshRate: &duration.Duration{
							Seconds: 10,
							Nanos:   20,
						},
						HostTtl: &duration.Duration{
							Seconds: 30,
							Nanos:   40,
						},
						MaxHosts: &wrappers.UInt32Value{
							Value: 10,
						},
						DnsCacheType: &dynamic_forward_proxy.DnsCacheConfig_CaresDns{
							CaresDns: &dynamic_forward_proxy.CaresDnsResolverConfig{
								Resolvers: []*v32.Address{
									{
										Address: &v32.Address_SocketAddress{
											SocketAddress: &v32.SocketAddress{
												Protocol: v32.SocketAddress_UDP,
												Address:  "127.0.0.1",
												PortSpecifier: &v32.SocketAddress_PortValue{
													PortValue: 80,
												},
												ResolverName: "resolverName",
												Ipv4Compat:   true,
											},
										},
									},
								},
								DnsResolverOptions: &dynamic_forward_proxy.DnsResolverOptions{
									UseTcpForDnsLookups:   true,
									NoDefaultSearchDomain: true,
								},
							},
						},
					},
					SaveUpstreamAddress: true,
				},
			}
		})

		It("translates cares config and top level fields", func() {
			p := NewPlugin()
			p.Init(initParams)

			filters, err := p.HttpFilters(params, listener)
			Expect(err).NotTo(HaveOccurred())
			Expect(filters).To(HaveLen(1))

			filterCfg := &envoy_extensions_filters_http_dynamic_forward_proxy_v3.FilterConfig{}
			goTypedConfig := filters[0].Filter.GetTypedConfig()
			err = ptypes.UnmarshalAny(goTypedConfig, filterCfg)
			Expect(err).NotTo(HaveOccurred())

			Expect(filterCfg.GetDnsCacheConfig().GetDnsLookupFamily()).To(Equal(envoy_config_cluster_v3.Cluster_V4_ONLY))
			Expect(filterCfg.GetDnsCacheConfig().GetDnsRefreshRate()).To(Equal(&duration.Duration{
				Seconds: 10,
				Nanos:   20,
			}))
			Expect(filterCfg.GetDnsCacheConfig().GetHostTtl()).To(Equal(&duration.Duration{
				Seconds: 30,
				Nanos:   40,
			}))
			Expect(filterCfg.GetDnsCacheConfig().GetMaxHosts()).To(Equal(&wrappers.UInt32Value{
				Value: 10,
			}))

			Expect(filterCfg.GetDnsCacheConfig().GetTypedDnsResolverConfig().Name).To(Equal("envoy.network.dns_resolver.cares"))
			caresCfg := utils.MustAnyToMessage(filterCfg.GetDnsCacheConfig().GetTypedDnsResolverConfig().TypedConfig).(*envoy_extensions_network_dns_resolver_cares_v3.CaresDnsResolverConfig)
			Expect(caresCfg.DnsResolverOptions.UseTcpForDnsLookups).To(BeTrue())
			Expect(caresCfg.DnsResolverOptions.NoDefaultSearchDomain).To(BeTrue())

			Expect(caresCfg.Resolvers).To(HaveLen(1))
			Expect(caresCfg.Resolvers[0].GetSocketAddress()).To(matchers.MatchProto(&envoy_config_core_v3.SocketAddress{
				Protocol: envoy_config_core_v3.SocketAddress_UDP,
				Address:  "127.0.0.1",
				PortSpecifier: &envoy_config_core_v3.SocketAddress_PortValue{
					PortValue: 80,
				},
				ResolverName: "resolverName",
				Ipv4Compat:   true,
			}))
		})

		It("Translates SslConfig", func() {
			// create dummy snapshot
			params.Snapshot = &gloosnapshot.ApiSnapshot{}

			// create plugin
			p := NewPlugin()
			p.Init(initParams)

			// inform plugin of listener
			listener1 := listener.Clone().(*v1.HttpListener)
			listener1.Options.DynamicForwardProxy.SslConfig = &ssl.UpstreamSslConfig{}
			_, err := p.HttpFilters(params, listener1)
			Expect(err).NotTo(HaveOccurred())

			// use plugin to compute expected envoy cluster
			clusters, _, _, _, _ := p.GeneratedResources(params, nil, nil, nil, nil)
			Expect(clusters).To(HaveLen(1))

			// evaluate contents of generated cluster
			ts := clusters[0].GetTransportSocket()
			Expect(ts).NotTo(BeNil())
			Expect(ts.GetName()).To(Equal("envoy.transport_sockets.tls"))
			Expect(ts.GetTypedConfig().GetTypeUrl()).To(Equal("type.googleapis.com/envoy.extensions.transport_sockets.tls.v3.UpstreamTlsContext"))
		})

		It("cache config name is per dns cache config", func() {
			p := NewPlugin()
			p.Init(initParams)

			listener1 := listener.Clone().(*v1.HttpListener)
			listener2 := listener.Clone().(*v1.HttpListener)
			// make a change to listener2 so DNS cache config is different
			listener1.Options.DynamicForwardProxy.DnsCacheConfig = &dynamic_forward_proxy.DnsCacheConfig{DnsLookupFamily: dynamic_forward_proxy.DnsLookupFamily_V4_ONLY}
			listener2.Options.DynamicForwardProxy.DnsCacheConfig = &dynamic_forward_proxy.DnsCacheConfig{DnsLookupFamily: dynamic_forward_proxy.DnsLookupFamily_V6_ONLY}

			filters1, err := p.HttpFilters(params, listener1)
			Expect(err).NotTo(HaveOccurred())
			Expect(filters1).To(HaveLen(1))
			filterCfg := &envoy_extensions_filters_http_dynamic_forward_proxy_v3.FilterConfig{}
			goTypedConfig := filters1[0].Filter.GetTypedConfig()
			err = ptypes.UnmarshalAny(goTypedConfig, filterCfg)
			Expect(err).NotTo(HaveOccurred())
			expectedName1 := "solo_io_generated_dfp:16275453913408509128"
			Expect(filterCfg.GetDnsCacheConfig().GetName()).To(Equal(expectedName1))

			filters2, err := p.HttpFilters(params, listener2)
			Expect(err).NotTo(HaveOccurred())
			Expect(filters2).To(HaveLen(1))
			filterCfg2 := &envoy_extensions_filters_http_dynamic_forward_proxy_v3.FilterConfig{}
			goTypedConfig2 := filters2[0].Filter.GetTypedConfig()
			err = ptypes.UnmarshalAny(goTypedConfig2, filterCfg2)
			Expect(err).NotTo(HaveOccurred())
			expectedName2 := "solo_io_generated_dfp:13361132607922819491"
			Expect(filterCfg2.GetDnsCacheConfig().GetName()).To(Equal(expectedName2))

			// different DNS caches should have different cache names
			Expect(expectedName1).NotTo(Equal(expectedName2))
		})

		var (
			settingsCbc = &v1.Settings{
				Gloo: &v1.GlooOptions{
					CircuitBreakers: &v1_circuitbreaker.CircuitBreakerConfig{
						MaxConnections:     &wrappers.UInt32Value{Value: 1000},
						MaxPendingRequests: &wrappers.UInt32Value{Value: 1000},
						MaxRequests:        &wrappers.UInt32Value{Value: 1000},
						MaxRetries:         &wrappers.UInt32Value{Value: 1000},
						TrackRemaining:     true,
					},
				},
			}

			dfpCbc = &v1_circuitbreaker.CircuitBreakerConfig{
				MaxConnections:     &wrappers.UInt32Value{Value: 2000},
				MaxPendingRequests: &wrappers.UInt32Value{Value: 2000},
				MaxRequests:        &wrappers.UInt32Value{Value: 2000},
				MaxRetries:         &wrappers.UInt32Value{Value: 2000},
				TrackRemaining:     true,
			}

			cbFromSettings = &envoy_config_cluster_v3.CircuitBreakers{
				Thresholds: []*envoy_config_cluster_v3.CircuitBreakers_Thresholds{
					{
						MaxConnections:     &wrappers.UInt32Value{Value: 1000},
						MaxPendingRequests: &wrappers.UInt32Value{Value: 1000},
						MaxRequests:        &wrappers.UInt32Value{Value: 1000},
						MaxRetries:         &wrappers.UInt32Value{Value: 1000},
						TrackRemaining:     true,
					},
				},
			}

			cbFromDfpcb = &envoy_config_cluster_v3.CircuitBreakers{
				Thresholds: []*envoy_config_cluster_v3.CircuitBreakers_Thresholds{
					{
						MaxConnections:     &wrappers.UInt32Value{Value: 2000},
						MaxPendingRequests: &wrappers.UInt32Value{Value: 2000},
						MaxRequests:        &wrappers.UInt32Value{Value: 2000},
						MaxRetries:         &wrappers.UInt32Value{Value: 2000},
						TrackRemaining:     true,
					},
				},
			}
		)

		DescribeTable("Translates CircuitBreakers", func(
			settings *v1.Settings,
			dfpConfig *v1_circuitbreaker.CircuitBreakerConfig,
			expectedCbc *envoy_config_cluster_v3.CircuitBreakers) {

			//create dummy snapshot
			params.Snapshot = &gloosnapshot.ApiSnapshot{}
			if settings != nil {
				params.Settings = settings
			}

			reports := reporter.ResourceReports{}

			// create plugin
			p := NewPlugin()
			p.Init(initParams)

			// inform plugin of listener
			listener1 := listener.Clone().(*v1.HttpListener)
			if dfpConfig != nil {
				listener1.Options.DynamicForwardProxy.CircuitBreakers = dfpConfig
			}
			_, err := p.HttpFilters(params, listener1)
			Expect(err).NotTo(HaveOccurred())

			// use plugin to compute expected envoy cluster
			clusters, _, _, _, _ := p.GeneratedResources(params, nil, nil, nil, nil)
			Expect(reports).To(BeEmpty())
			Expect(clusters).To(HaveLen(1))

			circuitBreaker := clusters[0].GetCircuitBreakers()
			Expect(circuitBreaker).To(BeEquivalentTo(expectedCbc))

		},
			Entry("no config", nil, nil, nil),
			Entry("only defined in settings", settingsCbc, nil, cbFromSettings),
			Entry("only defined in dfp config", nil, dfpCbc, cbFromDfpcb),
			Entry("defined in both", settingsCbc, dfpCbc, cbFromDfpcb),
		)
	})
})
