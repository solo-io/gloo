package failover_test

import (
	"context"

	v2 "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	envoy_api_v2_auth "github.com/envoyproxy/go-control-plane/envoy/api/v2/auth"
	envoy_api_v2_core "github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	envoy_api_v2_endpoint "github.com/envoyproxy/go-control-plane/envoy/api/v2/endpoint"
	"github.com/envoyproxy/go-control-plane/pkg/wellknown"
	"github.com/gogo/protobuf/types"
	"github.com/golang/mock/gomock"
	structpb "github.com/golang/protobuf/ptypes/struct"
	"github.com/golang/protobuf/ptypes/wrappers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/api/v2/core"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins/pluginutils"
	"github.com/solo-io/go-utils/testutils"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/failover"
	mock_utils "github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/failover/mocks"
)

var _ = Describe("Failover", func() {

	var (
		ctrl          *gomock.Controller
		ctx           context.Context
		sslTranslator *mock_utils.MockSslConfigTranslator

		sslEndpoint = &gloov1.LbEndpoint{
			Address: "ssl.address.who.dis",
			Port:    10101,
			UpstreamSslConfig: &gloov1.UpstreamSslConfig{
				Sni: "test",
			},
			LoadBalancingWeight: &types.UInt32Value{
				Value: 9999,
			},
		}
		tlsContext = &envoy_api_v2_auth.UpstreamTlsContext{
			Sni: "test",
		}
		httpEndpoint = &gloov1.LbEndpoint{
			Address: "http.address.who.dis",
			Port:    10101,
			HealthCheckConfig: &gloov1.LbEndpoint_HealthCheckConfig{
				PortValue: 9090,
				Hostname:  "new.host.who.dis",
			},
			LoadBalancingWeight: &types.UInt32Value{
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
								LoadBalancingWeight: &types.UInt32Value{
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
								LoadBalancingWeight: &types.UInt32Value{
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

		expected = &v2.ClusterLoadAssignment{
			Endpoints: []*envoy_api_v2_endpoint.LocalityLbEndpoints{
				{
					Locality: &envoy_api_v2_core.Locality{
						Region:  "p1_region",
						Zone:    "p1_zone",
						SubZone: "p1_sub_zone",
					},
					LbEndpoints: []*envoy_api_v2_endpoint.LbEndpoint{
						{
							HostIdentifier: &envoy_api_v2_endpoint.LbEndpoint_Endpoint{
								Endpoint: &envoy_api_v2_endpoint.Endpoint{
									Address: &envoy_api_v2_core.Address{
										Address: &envoy_api_v2_core.Address_SocketAddress{
											SocketAddress: &envoy_api_v2_core.SocketAddress{
												Address: sslEndpoint.GetAddress(),
												PortSpecifier: &envoy_api_v2_core.SocketAddress_PortValue{
													PortValue: sslEndpoint.GetPort(),
												},
											},
										},
									},
								},
							},
							Metadata: &envoy_api_v2_core.Metadata{
								FilterMetadata: map[string]*structpb.Struct{
									failover.TransportSocketMatchKey: metadataMatch,
								},
							},
							LoadBalancingWeight: &wrappers.UInt32Value{
								Value: sslEndpoint.GetLoadBalancingWeight().GetValue(),
							},
						},
					},
					LoadBalancingWeight: &wrappers.UInt32Value{
						Value: 8888,
					},
					Priority: 1,
				},
				{
					Locality: &envoy_api_v2_core.Locality{
						Region:  "p2_region",
						Zone:    "p2_zone",
						SubZone: "p2_sub_zone",
					},
					LbEndpoints: []*envoy_api_v2_endpoint.LbEndpoint{
						{
							HostIdentifier: &envoy_api_v2_endpoint.LbEndpoint_Endpoint{
								Endpoint: &envoy_api_v2_endpoint.Endpoint{
									Address: &envoy_api_v2_core.Address{
										Address: &envoy_api_v2_core.Address_SocketAddress{
											SocketAddress: &envoy_api_v2_core.SocketAddress{
												Address: httpEndpoint.GetAddress(),
												PortSpecifier: &envoy_api_v2_core.SocketAddress_PortValue{
													PortValue: httpEndpoint.GetPort(),
												},
											},
										},
									},
									HealthCheckConfig: &envoy_api_v2_endpoint.Endpoint_HealthCheckConfig{
										PortValue: httpEndpoint.GetHealthCheckConfig().GetPortValue(),
										Hostname:  httpEndpoint.GetHealthCheckConfig().GetHostname(),
									},
								},
							},
							LoadBalancingWeight: &wrappers.UInt32Value{
								Value: sslEndpoint.GetLoadBalancingWeight().GetValue(),
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

		buildExpectedCluster = func() *v2.Cluster {
			anyCfg, err := pluginutils.MessageToAny(tlsContext)
			Expect(err).NotTo(HaveOccurred())

			return &v2.Cluster{
				TransportSocketMatches: []*v2.Cluster_TransportSocketMatch{
					{
						Name:  uniqueName,
						Match: metadataMatch,
						TransportSocket: &envoy_api_v2_core.TransportSocket{
							Name: wellknown.TransportSocketTls,
							ConfigType: &envoy_api_v2_core.TransportSocket_TypedConfig{
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
			cluster *v2.Cluster,
			endpoints *v2.ClusterLoadAssignment,
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
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	It("will return nil if failover cfg is nil", func() {
		plugin := failover.NewFailoverPlugin(sslTranslator)
		err := runPlugin(plugin, plugins.Params{}, &gloov1.Upstream{}, nil, nil)
		Expect(err).NotTo(HaveOccurred())
	})

	It("will fail if no healthchecks are present", func() {
		plugin := failover.NewFailoverPlugin(sslTranslator)
		err := runPlugin(plugin, plugins.Params{}, &gloov1.Upstream{Failover: &gloov1.Failover{}}, nil, nil)
		Expect(err).To(HaveOccurred())
		Expect(err).To(testutils.HaveInErrorChain(failover.NoHealthCheckError))
	})

	It("will successfully return failover endpoints in the Cluster.ClusterLoadAssignment", func() {
		secretList := gloov1.SecretList{{}}
		sslTranslator.EXPECT().
			ResolveUpstreamSslConfig(secretList, sslEndpoint.GetUpstreamSslConfig()).
			Return(tlsContext, nil)

		cluster := &v2.Cluster{}
		cluster.LoadAssignment = &v2.ClusterLoadAssignment{}
		expectedCluster := buildExpectedCluster()
		expectedCluster.LoadAssignment = expected

		plugin := failover.NewFailoverPlugin(sslTranslator)
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

		cluster := &v2.Cluster{}
		cluster.ClusterDiscoveryType = &v2.Cluster_Type{
			Type: v2.Cluster_EDS,
		}
		expectedCluster := buildExpectedCluster()
		expectedCluster.ClusterDiscoveryType = &v2.Cluster_Type{
			Type: v2.Cluster_EDS,
		}

		plugin := failover.NewFailoverPlugin(sslTranslator)
		params := plugins.Params{
			Ctx: ctx,
			Snapshot: &gloov1.ApiSnapshot{
				Secrets: secretList,
			},
		}
		endpoints := &v2.ClusterLoadAssignment{}
		err := runPlugin(plugin, params, upstream, cluster, endpoints)
		Expect(err).NotTo(HaveOccurred())
		Expect(cluster).To(Equal(expectedCluster))
		Expect(endpoints).To(Equal(expected))
	})

})
