package failover_test

import (
	"context"
	"net"
	"time"

	envoy_config_cluster_v3 "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	envoy_config_core_v3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	envoy_config_endpoint_v3 "github.com/envoyproxy/go-control-plane/envoy/config/endpoint/v3"
	envoytls "github.com/envoyproxy/go-control-plane/envoy/extensions/transport_sockets/tls/v3"
	"github.com/envoyproxy/go-control-plane/pkg/wellknown"
	"github.com/golang/mock/gomock"
	"github.com/golang/protobuf/ptypes"
	structpb "github.com/golang/protobuf/ptypes/struct"
	"github.com/golang/protobuf/ptypes/wrappers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/api/v2/cluster"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/api/v2/core"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins"
	mock_consul "github.com/solo-io/gloo/projects/gloo/pkg/plugins/consul/mocks"
	"github.com/solo-io/gloo/projects/gloo/pkg/utils"
	"github.com/solo-io/go-utils/testutils"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/failover"
	mock_utils "github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/failover/mocks"
)

var _ = Describe("Failover", func() {

	var (
		ctrl          *gomock.Controller
		ctx           context.Context
		sslTranslator *mock_utils.MockSslConfigTranslator
		dnsResolver   *mock_consul.MockDnsResolver

		sslEndpoint = &gloov1.LbEndpoint{
			Address: "ssl.address.who.dis",
			Port:    10101,
			UpstreamSslConfig: &gloov1.UpstreamSslConfig{
				Sni: "test",
			},
		}
		tlsContext = &envoytls.UpstreamTlsContext{
			Sni: "test",
		}
		ipAddr1 = net.IPAddr{
			IP: []byte("10.0.0.1"),
		}
		ipAddr2 = net.IPAddr{
			IP: []byte("10.0.0.2"),
		}
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

		upstream = &gloov1.Upstream{
			HealthChecks: []*core.HealthCheck{{}},
			Failover: &gloov1.Failover{
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

		uniqueName    = failover.PrioritizedEndpointName(sslEndpoint.GetAddress(), sslEndpoint.GetPort(), 1, 0)
		metadataMatch = &structpb.Struct{
			Fields: map[string]*structpb.Value{
				uniqueName: {
					Kind: &structpb.Value_BoolValue{
						BoolValue: true,
					},
				},
			},
		}

		expected = &envoy_config_endpoint_v3.ClusterLoadAssignment{
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
									Address: &envoy_config_core_v3.Address{
										Address: &envoy_config_core_v3.Address_SocketAddress{
											SocketAddress: &envoy_config_core_v3.SocketAddress{
												Address: ipAddr1.IP.String(),
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
									failover.TransportSocketMatchKey: metadataMatch,
								},
							},
						},
						{
							HostIdentifier: &envoy_config_endpoint_v3.LbEndpoint_Endpoint{
								Endpoint: &envoy_config_endpoint_v3.Endpoint{
									Address: &envoy_config_core_v3.Address{
										Address: &envoy_config_core_v3.Address_SocketAddress{
											SocketAddress: &envoy_config_core_v3.SocketAddress{
												Address: ipAddr2.IP.String(),
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
									failover.TransportSocketMatchKey: metadataMatch,
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
			err := plugin.Init(plugins.InitParams{Ctx: ctx})
			Expect(err).NotTo(HaveOccurred())
			ups, ok := plugin.(plugins.UpstreamPlugin)
			Expect(ok).To(BeTrue())
			err = ups.ProcessUpstream(params, upstream, cluster)
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
		sslTranslator = mock_utils.NewMockSslConfigTranslator(ctrl)
		dnsResolver = mock_consul.NewMockDnsResolver(ctrl)
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	It("will return nil if failover cfg is nil", func() {
		plugin := failover.NewFailoverPlugin(sslTranslator, dnsResolver)
		err := runPlugin(plugin, plugins.Params{}, &gloov1.Upstream{}, nil, nil)
		Expect(err).NotTo(HaveOccurred())
	})

	It("will fail if no healthchecks are present", func() {
		plugin := failover.NewFailoverPlugin(sslTranslator, dnsResolver)
		err := runPlugin(plugin, plugins.Params{}, &gloov1.Upstream{Failover: &gloov1.Failover{}}, nil, nil)
		Expect(err).To(HaveOccurred())
		Expect(err).To(testutils.HaveInErrorChain(failover.NoHealthCheckError))
	})

	It("will fail if a DNS endpoint is specified with weights", func() {
		plugin := failover.NewFailoverPlugin(sslTranslator, dnsResolver)
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
		secretList := gloov1.SecretList{{}}
		sslTranslator.EXPECT().
			ResolveUpstreamSslConfig(secretList, sslEndpoint.GetUpstreamSslConfig()).
			Return(tlsContext, nil)

		dnsResolver.EXPECT().Resolve(ctx, sslEndpoint.GetAddress()).Return([]net.IPAddr{ipAddr1, ipAddr2}, nil)

		cluster := &envoy_config_cluster_v3.Cluster{}
		cluster.LoadAssignment = &envoy_config_endpoint_v3.ClusterLoadAssignment{}
		expectedCluster := buildExpectedCluster()
		expectedCluster.LoadAssignment = expected

		plugin := failover.NewFailoverPlugin(sslTranslator, dnsResolver)
		params := plugins.Params{
			Ctx: ctx,
			Snapshot: &gloov1.ApiSnapshot{
				Secrets: secretList,
			},
		}
		err := runPlugin(plugin, params, upstream, cluster, nil)
		Expect(err).NotTo(HaveOccurred())
		Expect(cluster).To(Equal(expectedCluster))
	})

	It("will successfully return failover endpoints in the EDS ClusterLoadAssignment", func() {

		secretList := gloov1.SecretList{{}}
		sslTranslator.EXPECT().
			ResolveUpstreamSslConfig(secretList, sslEndpoint.GetUpstreamSslConfig()).
			Return(tlsContext, nil)

		dnsResolver.EXPECT().Resolve(ctx, sslEndpoint.GetAddress()).Return([]net.IPAddr{ipAddr1, ipAddr2}, nil)

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
				ConfigSourceSpecifier: &envoy_config_core_v3.ConfigSource_ApiConfigSource{
					ApiConfigSource: &envoy_config_core_v3.ApiConfigSource{
						TransportApiVersion:       envoy_config_core_v3.ApiVersion_V3,
						ApiType:                   envoy_config_core_v3.ApiConfigSource_REST,
						ClusterNames:              []string{failover.RestXdsCluster},
						RefreshDelay:              ptypes.DurationProto(time.Second * 5),
						RequestTimeout:            ptypes.DurationProto(time.Second * 5),
						RateLimitSettings:         nil,
						SetNodeOnFirstMessageOnly: false,
					},
				},
			},
		}

		plugin := failover.NewFailoverPlugin(sslTranslator, dnsResolver)
		params := plugins.Params{
			Ctx: ctx,
			Snapshot: &gloov1.ApiSnapshot{
				Secrets: secretList,
			},
		}
		endpoints := &envoy_config_endpoint_v3.ClusterLoadAssignment{}
		err := runPlugin(plugin, params, upstream, cluster, endpoints)
		Expect(err).NotTo(HaveOccurred())
		Expect(cluster).To(Equal(expectedCluster))
		Expect(endpoints).To(Equal(expected))
	})

})
