package failover_test

import (
	"context"
	"fmt"
	"net"

	envoy_config_cluster_v3 "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	envoy_config_core_v3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	envoy_config_endpoint_v3 "github.com/envoyproxy/go-control-plane/envoy/config/endpoint/v3"
	envoytls "github.com/envoyproxy/go-control-plane/envoy/extensions/transport_sockets/tls/v3"
	"github.com/envoyproxy/go-control-plane/pkg/wellknown"
	"github.com/golang/mock/gomock"
	structpb "github.com/golang/protobuf/ptypes/struct"
	"github.com/golang/protobuf/ptypes/wrappers"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rotisserie/eris"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/api/v2/cluster"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/api/v2/core"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	gloov1snap "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/gloosnapshot"
	gloossl "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/ssl"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins"
	mock_consul "github.com/solo-io/gloo/projects/gloo/pkg/plugins/consul/mocks"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins/static"
	"github.com/solo-io/gloo/projects/gloo/pkg/utils"
	"github.com/solo-io/go-utils/testutils"
	"github.com/solo-io/solo-kit/test/matchers"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/failover"
	mock_utils "github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/failover/mocks"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

var _ = Describe("Failover", func() {

	var (
		ctrl          *gomock.Controller
		ctx           context.Context
		cancel        context.CancelFunc
		sslTranslator *mock_utils.MockSslConfigTranslator
		dnsResolver   *mock_consul.MockDnsResolver

		apiEmitNotificationChan chan struct{}

		tlsContext = &envoytls.UpstreamTlsContext{
			Sni: "test",
		}
		ipAddr1 = net.IPAddr{
			IP: net.IPv4(10, 0, 0, 1),
		}
		ipAddr2 = net.IPAddr{
			IP: net.IPv4(10, 0, 0, 2),
		}
		upstream     *gloov1.Upstream
		httpEndpoint *gloov1.LbEndpoint
		sslEndpoint  *gloov1.LbEndpoint

		uniqueName    string
		metadataMatch *structpb.Struct

		buildExpectedCluster = func() *envoy_config_cluster_v3.Cluster {
			anyCfg, err := utils.MessageToAny(tlsContext)
			Expect(err).NotTo(HaveOccurred())

			return &envoy_config_cluster_v3.Cluster{
				TransportSocketMatches: []*envoy_config_cluster_v3.Cluster_TransportSocketMatch{
					{
						Name:  uniqueName,
						Match: metadataMatch,
						TransportSocket: &envoy_config_core_v3.TransportSocket{
							Name: wellknown.TransportSocketTls,
							ConfigType: &envoy_config_core_v3.TransportSocket_TypedConfig{
								TypedConfig: anyCfg,
							},
						},
					},
				},
			}
		}

		runPlugin = func(
			plugin plugins.Plugin,
			params plugins.Params,
			upstream *gloov1.Upstream,
			cluster *envoy_config_cluster_v3.Cluster,
			endpoints *envoy_config_endpoint_v3.ClusterLoadAssignment,
		) error {
			plugin.Init(plugins.InitParams{Ctx: ctx, Settings: &gloov1.Settings{Gloo: &gloov1.GlooOptions{
				FailoverUpstreamDnsPollingInterval: &durationpb.Duration{
					Seconds: 1,
				},
			}}})
			ups, ok := plugin.(plugins.UpstreamPlugin)
			Expect(ok).To(BeTrue())
			err := ups.ProcessUpstream(params, upstream, cluster)
			if err != nil {
				return err
			}
			eps, ok := plugin.(plugins.EndpointPlugin)
			Expect(ok).To(BeTrue())
			return eps.ProcessEndpoints(params, upstream, endpoints)
		}
	)

	BeforeEach(func() {
		ctrl, ctx = gomock.WithContext(context.TODO(), GinkgoT())
		ctx, cancel = context.WithCancel(ctx)
		sslTranslator = mock_utils.NewMockSslConfigTranslator(ctrl)
		dnsResolver = mock_consul.NewMockDnsResolver(ctrl)
		apiEmitNotificationChan = make(chan struct{})

		httpEndpoint = &gloov1.LbEndpoint{
			Address: "127.0.0.1",
			Port:    10101,
			HealthCheckConfig: &gloov1.LbEndpoint_HealthCheckConfig{
				PortValue: 9090,
				Hostname:  "new.host.who.dis",
			},
			LoadBalancingWeight: &wrappers.UInt32Value{
				Value: 9999,
			},
		}

		sslEndpoint = &gloov1.LbEndpoint{
			Address: "ssl.address.who.dis",
			Port:    10101,
			UpstreamSslConfig: &gloossl.UpstreamSslConfig{
				Sni: "test",
			},
		}

		upstream = &gloov1.Upstream{
			HealthChecks: []*core.HealthCheck{{}},
			Failover: &gloov1.Failover{
				Policy: &gloov1.Failover_Policy{
					OverprovisioningFactor: wrapperspb.UInt32(123),
				},
				PrioritizedLocalities: []*gloov1.Failover_PrioritizedLocality{
					{
						LocalityEndpoints: []*gloov1.LocalityLbEndpoints{
							{
								Locality: &gloov1.Locality{
									Region:  "p1_region",
									Zone:    "p1_zone",
									SubZone: "p1_sub_zone",
								},
								LbEndpoints: []*gloov1.LbEndpoint{
									sslEndpoint,
								},
								LoadBalancingWeight: &wrappers.UInt32Value{
									Value: 8888,
								},
							},
						},
					},
					{
						LocalityEndpoints: []*gloov1.LocalityLbEndpoints{
							{
								Locality: &gloov1.Locality{
									Region:  "p2_region",
									Zone:    "p2_zone",
									SubZone: "p2_sub_zone",
								},
								LbEndpoints: []*gloov1.LbEndpoint{
									httpEndpoint,
								},
								LoadBalancingWeight: &wrappers.UInt32Value{
									Value: 7777,
								},
							},
						},
					},
				},
			},
		}
		uniqueName = failover.PrioritizedEndpointName(sslEndpoint.GetAddress(), sslEndpoint.GetPort(), 1, 0)
		metadataMatch = &structpb.Struct{
			Fields: map[string]*structpb.Value{
				uniqueName: {
					Kind: &structpb.Value_BoolValue{
						BoolValue: true,
					},
				},
			},
		}
	})

	AfterEach(func() {
		close(apiEmitNotificationChan)
		ctrl.Finish()
		cancel()
	})

	It("will return nil if failover cfg is nil", func() {
		plugin := failover.NewFailoverPlugin(sslTranslator, dnsResolver, apiEmitNotificationChan)
		err := runPlugin(plugin, plugins.Params{}, &gloov1.Upstream{}, nil, nil)
		Expect(err).NotTo(HaveOccurred())
	})

	It("will fail if no healthchecks are present", func() {
		plugin := failover.NewFailoverPlugin(sslTranslator, dnsResolver, apiEmitNotificationChan)
		err := runPlugin(plugin, plugins.Params{}, &gloov1.Upstream{Failover: &gloov1.Failover{}}, nil, nil)
		Expect(err).To(HaveOccurred())
		Expect(err).To(testutils.HaveInErrorChain(failover.NoHealthCheckError))
	})

	It("will fail if a DNS endpoint is specified with weights", func() {
		plugin := failover.NewFailoverPlugin(sslTranslator, dnsResolver, apiEmitNotificationChan)
		err := runPlugin(plugin, plugins.Params{}, &gloov1.Upstream{
			OutlierDetection: &cluster.OutlierDetection{},
			Failover: &gloov1.Failover{
				PrioritizedLocalities: []*gloov1.Failover_PrioritizedLocality{
					{
						LocalityEndpoints: []*gloov1.LocalityLbEndpoints{
							{
								LbEndpoints: []*gloov1.LbEndpoint{
									{
										Address: "dns.name",
										LoadBalancingWeight: &wrappers.UInt32Value{
											Value: 9999,
										},
									},
								},
							},
						},
					},
				},
			}}, nil, nil)
		Expect(err).To(HaveOccurred())
		Expect(err).To(testutils.HaveInErrorChain(failover.WeightedDnsError))
	})

	It("will successfully return failover endpoints in the Cluster.ClusterLoadAssignment", func() {

		expected := &envoy_config_endpoint_v3.ClusterLoadAssignment{
			Policy: &envoy_config_endpoint_v3.ClusterLoadAssignment_Policy{
				OverprovisioningFactor: &wrappers.UInt32Value{
					Value: 123,
				},
			},
			Endpoints: []*envoy_config_endpoint_v3.LocalityLbEndpoints{
				{
					Locality: &envoy_config_core_v3.Locality{
						Region:  "p1_region",
						Zone:    "p1_zone",
						SubZone: "p1_sub_zone",
					},
					LbEndpoints: []*envoy_config_endpoint_v3.LbEndpoint{
						{
							HostIdentifier: &envoy_config_endpoint_v3.LbEndpoint_Endpoint{
								Endpoint: &envoy_config_endpoint_v3.Endpoint{
									Hostname: sslEndpoint.GetAddress(),
									HealthCheckConfig: &envoy_config_endpoint_v3.Endpoint_HealthCheckConfig{
										Hostname: sslEndpoint.GetAddress(),
									},
									Address: &envoy_config_core_v3.Address{
										Address: &envoy_config_core_v3.Address_SocketAddress{
											SocketAddress: &envoy_config_core_v3.SocketAddress{
												Address: sslEndpoint.GetAddress(),
												PortSpecifier: &envoy_config_core_v3.SocketAddress_PortValue{
													PortValue: sslEndpoint.GetPort(),
												},
											},
										},
									},
								},
							},
							Metadata: &envoy_config_core_v3.Metadata{
								FilterMetadata: map[string]*structpb.Struct{
									static.TransportSocketMatchKey: metadataMatch,
								},
							},
						},
					},
					LoadBalancingWeight: &wrappers.UInt32Value{
						Value: 8888,
					},
					Priority: 1,
				},
				{
					Locality: &envoy_config_core_v3.Locality{
						Region:  "p2_region",
						Zone:    "p2_zone",
						SubZone: "p2_sub_zone",
					},
					LbEndpoints: []*envoy_config_endpoint_v3.LbEndpoint{
						{
							HostIdentifier: &envoy_config_endpoint_v3.LbEndpoint_Endpoint{
								Endpoint: &envoy_config_endpoint_v3.Endpoint{
									Address: &envoy_config_core_v3.Address{
										Address: &envoy_config_core_v3.Address_SocketAddress{
											SocketAddress: &envoy_config_core_v3.SocketAddress{
												Address: httpEndpoint.GetAddress(),
												PortSpecifier: &envoy_config_core_v3.SocketAddress_PortValue{
													PortValue: httpEndpoint.GetPort(),
												},
											},
										},
									},
									HealthCheckConfig: &envoy_config_endpoint_v3.Endpoint_HealthCheckConfig{
										PortValue: httpEndpoint.GetHealthCheckConfig().GetPortValue(),
										Hostname:  httpEndpoint.GetHealthCheckConfig().GetHostname(),
									},
									Hostname: httpEndpoint.GetAddress(),
								},
							},
							LoadBalancingWeight: &wrappers.UInt32Value{
								Value: httpEndpoint.GetLoadBalancingWeight().GetValue(),
							},
						},
					},
					LoadBalancingWeight: &wrappers.UInt32Value{
						Value: 7777,
					},
					Priority: 2,
				},
			},
		}

		secretList := gloov1.SecretList{{}}
		sslTranslator.EXPECT().
			ResolveUpstreamSslConfig(secretList, sslEndpoint.GetUpstreamSslConfig()).
			Return(tlsContext, nil)

		cluster := &envoy_config_cluster_v3.Cluster{
			ClusterDiscoveryType: &envoy_config_cluster_v3.Cluster_Type{
				Type: envoy_config_cluster_v3.Cluster_STRICT_DNS,
			},
		}
		cluster.LoadAssignment = &envoy_config_endpoint_v3.ClusterLoadAssignment{
			Policy: &envoy_config_endpoint_v3.ClusterLoadAssignment_Policy{
				OverprovisioningFactor: &wrappers.UInt32Value{Value: 123},
			},
		}
		expectedCluster := buildExpectedCluster()
		expectedCluster.LoadAssignment = expected
		expectedCluster.ClusterDiscoveryType = &envoy_config_cluster_v3.Cluster_Type{
			Type: envoy_config_cluster_v3.Cluster_STRICT_DNS,
		}

		plugin := failover.NewFailoverPlugin(sslTranslator, dnsResolver, apiEmitNotificationChan)
		params := plugins.Params{
			Ctx: ctx,
			Snapshot: &gloov1snap.ApiSnapshot{
				Secrets: secretList,
			},
		}
		err := runPlugin(plugin, params, upstream, cluster, &envoy_config_endpoint_v3.ClusterLoadAssignment{})
		Expect(err).NotTo(HaveOccurred())
		Expect(cluster).To(matchers.MatchProto(expectedCluster))
	})

	It("will successfully return failover endpoints in the EDS ClusterLoadAssignment", func() {
		expected := &envoy_config_endpoint_v3.ClusterLoadAssignment{
			Policy: &envoy_config_endpoint_v3.ClusterLoadAssignment_Policy{
				OverprovisioningFactor: &wrappers.UInt32Value{
					Value: 123,
				},
			},
			Endpoints: []*envoy_config_endpoint_v3.LocalityLbEndpoints{
				{
					Locality: &envoy_config_core_v3.Locality{
						Region:  "p1_region",
						Zone:    "p1_zone",
						SubZone: "p1_sub_zone",
					},
					LbEndpoints: []*envoy_config_endpoint_v3.LbEndpoint{
						{
							HostIdentifier: &envoy_config_endpoint_v3.LbEndpoint_Endpoint{
								Endpoint: &envoy_config_endpoint_v3.Endpoint{
									Hostname: sslEndpoint.GetAddress(),
									HealthCheckConfig: &envoy_config_endpoint_v3.Endpoint_HealthCheckConfig{
										Hostname: sslEndpoint.GetAddress(),
									},
									Address: &envoy_config_core_v3.Address{
										Address: &envoy_config_core_v3.Address_SocketAddress{
											SocketAddress: &envoy_config_core_v3.SocketAddress{
												Address: ipAddr1.String(),
												PortSpecifier: &envoy_config_core_v3.SocketAddress_PortValue{
													PortValue: sslEndpoint.GetPort(),
												},
											},
										},
									},
								},
							},
							Metadata: &envoy_config_core_v3.Metadata{
								FilterMetadata: map[string]*structpb.Struct{
									static.TransportSocketMatchKey: metadataMatch,
								},
							},
						},
						{
							HostIdentifier: &envoy_config_endpoint_v3.LbEndpoint_Endpoint{
								Endpoint: &envoy_config_endpoint_v3.Endpoint{
									Hostname: sslEndpoint.GetAddress(),
									HealthCheckConfig: &envoy_config_endpoint_v3.Endpoint_HealthCheckConfig{
										Hostname: sslEndpoint.GetAddress(),
									},
									Address: &envoy_config_core_v3.Address{
										Address: &envoy_config_core_v3.Address_SocketAddress{
											SocketAddress: &envoy_config_core_v3.SocketAddress{
												Address: ipAddr2.String(),
												PortSpecifier: &envoy_config_core_v3.SocketAddress_PortValue{
													PortValue: sslEndpoint.GetPort(),
												},
											},
										},
									},
								},
							},
							Metadata: &envoy_config_core_v3.Metadata{
								FilterMetadata: map[string]*structpb.Struct{
									static.TransportSocketMatchKey: metadataMatch,
								},
							},
						},
					},
					LoadBalancingWeight: &wrappers.UInt32Value{
						Value: 8888,
					},
					Priority: 1,
				},
				{
					Locality: &envoy_config_core_v3.Locality{
						Region:  "p2_region",
						Zone:    "p2_zone",
						SubZone: "p2_sub_zone",
					},
					LbEndpoints: []*envoy_config_endpoint_v3.LbEndpoint{
						{
							HostIdentifier: &envoy_config_endpoint_v3.LbEndpoint_Endpoint{
								Endpoint: &envoy_config_endpoint_v3.Endpoint{
									Address: &envoy_config_core_v3.Address{
										Address: &envoy_config_core_v3.Address_SocketAddress{
											SocketAddress: &envoy_config_core_v3.SocketAddress{
												Address: httpEndpoint.GetAddress(),
												PortSpecifier: &envoy_config_core_v3.SocketAddress_PortValue{
													PortValue: httpEndpoint.GetPort(),
												},
											},
										},
									},
									Hostname: httpEndpoint.GetAddress(),
									HealthCheckConfig: &envoy_config_endpoint_v3.Endpoint_HealthCheckConfig{
										PortValue: httpEndpoint.GetHealthCheckConfig().GetPortValue(),
										Hostname:  httpEndpoint.GetHealthCheckConfig().GetHostname(),
									},
								},
							},
							LoadBalancingWeight: &wrappers.UInt32Value{
								Value: httpEndpoint.GetLoadBalancingWeight().GetValue(),
							},
						},
					},
					LoadBalancingWeight: &wrappers.UInt32Value{
						Value: 7777,
					},
					Priority: 2,
				},
			},
		}

		secretList := gloov1.SecretList{{}}
		sslTranslator.EXPECT().
			ResolveUpstreamSslConfig(secretList, sslEndpoint.GetUpstreamSslConfig()).
			Return(tlsContext, nil)

		dnsResolver.EXPECT().Resolve(gomock.Any(), sslEndpoint.GetAddress()).Return([]net.IPAddr{ipAddr1, ipAddr2}, nil)

		cluster := &envoy_config_cluster_v3.Cluster{}
		cluster.ClusterDiscoveryType = &envoy_config_cluster_v3.Cluster_Type{
			Type: envoy_config_cluster_v3.Cluster_EDS,
		}
		expectedCluster := buildExpectedCluster()
		expectedCluster.ClusterDiscoveryType = &envoy_config_cluster_v3.Cluster_Type{
			Type: envoy_config_cluster_v3.Cluster_EDS,
		}
		expectedCluster.EdsClusterConfig = &envoy_config_cluster_v3.Cluster_EdsClusterConfig{
			EdsConfig: &envoy_config_core_v3.ConfigSource{
				ResourceApiVersion: envoy_config_core_v3.ApiVersion_V3,
				ConfigSourceSpecifier: &envoy_config_core_v3.ConfigSource_Ads{
					Ads: &envoy_config_core_v3.AggregatedConfigSource{},
				},
			},
		}

		plugin := failover.NewFailoverPlugin(sslTranslator, dnsResolver, apiEmitNotificationChan)
		params := plugins.Params{
			Ctx: ctx,
			Snapshot: &gloov1snap.ApiSnapshot{
				Secrets: secretList,
			},
		}
		endpoints := &envoy_config_endpoint_v3.ClusterLoadAssignment{}
		err := runPlugin(plugin, params, upstream, cluster, endpoints)
		Expect(err).NotTo(HaveOccurred())
		Expect(cluster).To(matchers.MatchProto(expectedCluster))
		Expect(endpoints).To(matchers.MatchProto(expected))
	})

	It("will set endpoint metadata for advanced http check if path and/or method is set on upstream failover healthcheck", func() {

		secretList := gloov1.SecretList{{}}
		sslTranslator.EXPECT().
			ResolveUpstreamSslConfig(secretList, sslEndpoint.GetUpstreamSslConfig()).
			Return(tlsContext, nil)

		dnsResolver.EXPECT().Resolve(gomock.Any(), sslEndpoint.GetAddress()).Return([]net.IPAddr{ipAddr1, ipAddr2}, nil)

		cluster := &envoy_config_cluster_v3.Cluster{}
		cluster.ClusterDiscoveryType = &envoy_config_cluster_v3.Cluster_Type{
			Type: envoy_config_cluster_v3.Cluster_EDS,
		}
		expectedCluster := buildExpectedCluster()
		expectedCluster.ClusterDiscoveryType = &envoy_config_cluster_v3.Cluster_Type{
			Type: envoy_config_cluster_v3.Cluster_EDS,
		}
		expectedCluster.EdsClusterConfig = &envoy_config_cluster_v3.Cluster_EdsClusterConfig{
			EdsConfig: &envoy_config_core_v3.ConfigSource{
				ResourceApiVersion: envoy_config_core_v3.ApiVersion_V3,
				ConfigSourceSpecifier: &envoy_config_core_v3.ConfigSource_Ads{
					Ads: &envoy_config_core_v3.AggregatedConfigSource{},
				},
			},
		}

		plugin := failover.NewFailoverPlugin(sslTranslator, dnsResolver, apiEmitNotificationChan)
		params := plugins.Params{
			Ctx: ctx,
			Snapshot: &gloov1snap.ApiSnapshot{
				Secrets: secretList,
			},
		}
		endpoints := &envoy_config_endpoint_v3.ClusterLoadAssignment{}
		upstream.Failover.PrioritizedLocalities[0].LocalityEndpoints[0].LbEndpoints[0].HealthCheckConfig = &gloov1.LbEndpoint_HealthCheckConfig{
			Path:   "some/path/1",
			Method: "POST",
		}
		err := runPlugin(plugin, params, upstream, cluster, endpoints)

		Expect(err).NotTo(HaveOccurred())
		Expect(cluster).To(matchers.MatchProto(expectedCluster))
		filterMetadata := endpoints.GetEndpoints()[0].GetLbEndpoints()[0].GetMetadata().GetFilterMetadata()
		Expect(filterMetadata).To(HaveKey(static.AdvancedHttpCheckerName))
		fields := filterMetadata[static.AdvancedHttpCheckerName].GetFields()
		Expect(fields).To(HaveKey(static.PathFieldName))
		Expect(fields[static.PathFieldName].GetStringValue()).To(Equal("some/path/1"))
		Expect(fields).To(HaveKey(static.MethodFieldName))
		Expect(fields[static.MethodFieldName].GetStringValue()).To(Equal("POST"))
	})

	It("force emits when a DNS resolution changes", func() {
		secretList := gloov1.SecretList{{}}
		sslTranslator.EXPECT().
			ResolveUpstreamSslConfig(secretList, sslEndpoint.GetUpstreamSslConfig()).
			Return(tlsContext, nil)

		initialIps := []net.IPAddr{ipAddr1, ipAddr2}
		updatedIps := []net.IPAddr{{IP: net.IPv4(127, 0, 0, 1)}, {IP: net.IPv4(127, 0, 0, 100)}}

		// 2 times for the initial resolved that is called during the main call in ProcessUpstreams
		// and then the initial resolve that is called in the go routine
		dnsResolver.EXPECT().Resolve(gomock.Any(), gomock.Any()).Do(func(context.Context, string) {
			fmt.Fprint(GinkgoWriter, "Initial resolve called for endpoint.")
		}).Return(initialIps, nil).Times(2)

		// Once we see an updated address, all dns resolution go routines will be cancelled after the
		// notifcation channel is notified to emit
		dnsResolver.EXPECT().Resolve(gomock.Any(), gomock.Any()).Do(func(context.Context, string) {
			fmt.Fprint(GinkgoWriter, "Updated resolve called for endpoint.")
		}).Return(updatedIps, nil).Times(1)

		cluster := &envoy_config_cluster_v3.Cluster{}
		cluster.ClusterDiscoveryType = &envoy_config_cluster_v3.Cluster_Type{
			Type: envoy_config_cluster_v3.Cluster_EDS,
		}

		endpoints := &envoy_config_endpoint_v3.ClusterLoadAssignment{}

		plugin := failover.NewFailoverPlugin(sslTranslator, dnsResolver, apiEmitNotificationChan)
		params := plugins.Params{
			Ctx: ctx,
			Snapshot: &gloov1snap.ApiSnapshot{
				Secrets: secretList,
			},
		}

		err := runPlugin(plugin, params, upstream, cluster, endpoints)
		Expect(err).NotTo(HaveOccurred())

		Eventually(apiEmitNotificationChan, "5s", "1s").Should(Receive(BeEquivalentTo(struct{}{})))
		Consistently(apiEmitNotificationChan, "5s", "1s").ShouldNot(Receive())
	})

	It("uses the previous dns resolution when there is an error", func() {
		secretList := gloov1.SecretList{{}}
		sslTranslator.EXPECT().
			ResolveUpstreamSslConfig(secretList, sslEndpoint.GetUpstreamSslConfig()).
			Return(tlsContext, nil)

		initialIps := []net.IPAddr{ipAddr1, ipAddr2}
		resolutionError := eris.New("DNS resolution error")

		// 2 times for the initial resolved that is called during the main call in ProcessUpstreams
		// and then the initial resolve that is called in the go routine
		dnsResolver.EXPECT().Resolve(gomock.Any(), gomock.Any()).Do(func(context.Context, string) {
			fmt.Fprint(GinkgoWriter, "Initial resolve called for endpoint.")
		}).Return(initialIps, nil).Times(2)

		// We should use the previous IPs when we encounter this error. Given that, there will be no
		// DNS changes so expect that a notification will not be sent
		dnsResolver.EXPECT().Resolve(gomock.Any(), gomock.Any()).Do(func(context.Context, string) {
			fmt.Fprint(GinkgoWriter, "Errored resolve called for endpoint.")
		}).Return([]net.IPAddr{}, resolutionError).AnyTimes()

		cluster := &envoy_config_cluster_v3.Cluster{}
		cluster.ClusterDiscoveryType = &envoy_config_cluster_v3.Cluster_Type{
			Type: envoy_config_cluster_v3.Cluster_EDS,
		}

		endpoints := &envoy_config_endpoint_v3.ClusterLoadAssignment{}
		expected := &envoy_config_endpoint_v3.ClusterLoadAssignment{
			Policy: &envoy_config_endpoint_v3.ClusterLoadAssignment_Policy{
				OverprovisioningFactor: &wrappers.UInt32Value{
					Value: 123,
				},
			},
			Endpoints: []*envoy_config_endpoint_v3.LocalityLbEndpoints{
				{
					Locality: &envoy_config_core_v3.Locality{
						Region:  "p1_region",
						Zone:    "p1_zone",
						SubZone: "p1_sub_zone",
					},
					LbEndpoints: []*envoy_config_endpoint_v3.LbEndpoint{
						{
							HostIdentifier: &envoy_config_endpoint_v3.LbEndpoint_Endpoint{
								Endpoint: &envoy_config_endpoint_v3.Endpoint{
									Hostname: sslEndpoint.GetAddress(),
									HealthCheckConfig: &envoy_config_endpoint_v3.Endpoint_HealthCheckConfig{
										Hostname: sslEndpoint.GetAddress(),
									},
									Address: &envoy_config_core_v3.Address{
										Address: &envoy_config_core_v3.Address_SocketAddress{
											SocketAddress: &envoy_config_core_v3.SocketAddress{
												Address: ipAddr1.String(),
												PortSpecifier: &envoy_config_core_v3.SocketAddress_PortValue{
													PortValue: sslEndpoint.GetPort(),
												},
											},
										},
									},
								},
							},
							Metadata: &envoy_config_core_v3.Metadata{
								FilterMetadata: map[string]*structpb.Struct{
									static.TransportSocketMatchKey: metadataMatch,
								},
							},
						},
						{
							HostIdentifier: &envoy_config_endpoint_v3.LbEndpoint_Endpoint{
								Endpoint: &envoy_config_endpoint_v3.Endpoint{
									Hostname: sslEndpoint.GetAddress(),
									HealthCheckConfig: &envoy_config_endpoint_v3.Endpoint_HealthCheckConfig{
										Hostname: sslEndpoint.GetAddress(),
									},
									Address: &envoy_config_core_v3.Address{
										Address: &envoy_config_core_v3.Address_SocketAddress{
											SocketAddress: &envoy_config_core_v3.SocketAddress{
												Address: ipAddr2.String(),
												PortSpecifier: &envoy_config_core_v3.SocketAddress_PortValue{
													PortValue: sslEndpoint.GetPort(),
												},
											},
										},
									},
								},
							},
							Metadata: &envoy_config_core_v3.Metadata{
								FilterMetadata: map[string]*structpb.Struct{
									static.TransportSocketMatchKey: metadataMatch,
								},
							},
						},
					},
					LoadBalancingWeight: &wrappers.UInt32Value{
						Value: 8888,
					},
					Priority: 1,
				},
				{
					Locality: &envoy_config_core_v3.Locality{
						Region:  "p2_region",
						Zone:    "p2_zone",
						SubZone: "p2_sub_zone",
					},
					LbEndpoints: []*envoy_config_endpoint_v3.LbEndpoint{
						{
							HostIdentifier: &envoy_config_endpoint_v3.LbEndpoint_Endpoint{
								Endpoint: &envoy_config_endpoint_v3.Endpoint{
									Address: &envoy_config_core_v3.Address{
										Address: &envoy_config_core_v3.Address_SocketAddress{
											SocketAddress: &envoy_config_core_v3.SocketAddress{
												Address: httpEndpoint.GetAddress(),
												PortSpecifier: &envoy_config_core_v3.SocketAddress_PortValue{
													PortValue: httpEndpoint.GetPort(),
												},
											},
										},
									},
									Hostname: httpEndpoint.GetAddress(),
									HealthCheckConfig: &envoy_config_endpoint_v3.Endpoint_HealthCheckConfig{
										PortValue: httpEndpoint.GetHealthCheckConfig().GetPortValue(),
										Hostname:  httpEndpoint.GetHealthCheckConfig().GetHostname(),
									},
								},
							},
							LoadBalancingWeight: &wrappers.UInt32Value{
								Value: httpEndpoint.GetLoadBalancingWeight().GetValue(),
							},
						},
					},
					LoadBalancingWeight: &wrappers.UInt32Value{
						Value: 7777,
					},
					Priority: 2,
				},
			},
		}

		plugin := failover.NewFailoverPlugin(sslTranslator, dnsResolver, apiEmitNotificationChan)
		params := plugins.Params{
			Ctx: ctx,
			Snapshot: &gloov1snap.ApiSnapshot{
				Secrets: secretList,
			},
		}

		err := runPlugin(plugin, params, upstream, cluster, endpoints)
		Expect(err).NotTo(HaveOccurred())
		Expect(endpoints).To(matchers.MatchProto(expected))

		Consistently(apiEmitNotificationChan, "5s", "1s").ShouldNot(Receive())
	})

})
