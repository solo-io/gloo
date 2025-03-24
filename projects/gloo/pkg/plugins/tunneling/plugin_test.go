package tunneling_test

import (
	"fmt"

	envoy_config_cluster_v3 "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	envoy_config_core_v3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	envoy_config_endpoint_v3 "github.com/envoyproxy/go-control-plane/envoy/config/endpoint/v3"
	envoy_config_listener_v3 "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	envoy_config_route_v3 "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	envoytcp "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/network/tcp_proxy/v3"
	envoyauth "github.com/envoyproxy/go-control-plane/envoy/extensions/transport_sockets/tls/v3"
	"github.com/golang/protobuf/ptypes/wrappers"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	v1snap "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/gloosnapshot"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/ssl"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins/tunneling"
	"github.com/solo-io/gloo/projects/gloo/pkg/translator"
	"github.com/solo-io/gloo/projects/gloo/pkg/utils"
	"github.com/solo-io/skv2/test/matchers"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-kit/pkg/api/v2/reporter"
	"google.golang.org/protobuf/proto"
)

const (
	httpProxyHostname string = "host.com:443"
)

var _ = Describe("Plugin", func() {

	var (
		proxy                 *v1.Proxy
		params                plugins.Params
		inRouteConfigurations []*envoy_config_route_v3.RouteConfiguration
		inClusters            []*envoy_config_cluster_v3.Cluster

		us          *v1.Upstream
		clusterName string

		reports reporter.ResourceReports
	)

	BeforeEach(func() {
		proxy = &v1.Proxy{
			Metadata: &core.Metadata{
				Name:      "test",
				Namespace: "gloo-system",
			},
		}
		us = &v1.Upstream{
			Metadata: &core.Metadata{
				Name:      "http-proxy-upstream",
				Namespace: "gloo-system",
			},
			SslConfig:         nil,
			HttpProxyHostname: &wrappers.StringValue{Value: httpProxyHostname},
		}

		// use UpstreamToClusterName to emulate a real translation loop.
		clusterName = translator.UpstreamToClusterName(us.Metadata.Ref())

		params = plugins.Params{
			Snapshot: &v1snap.ApiSnapshot{
				Upstreams: []*v1.Upstream{us},
			},
			Settings: &v1.Settings{
				Gateway: &v1.GatewayOptions{},
			},
		}

		inRouteConfigurations = []*envoy_config_route_v3.RouteConfiguration{
			{
				Name: "listener-::-11082-routes",
				VirtualHosts: []*envoy_config_route_v3.VirtualHost{
					{
						Name:    "gloo-system_vs",
						Domains: []string{"*"},
						Routes: []*envoy_config_route_v3.Route{
							{
								Name: "testroute",
								Match: &envoy_config_route_v3.RouteMatch{
									PathSpecifier: &envoy_config_route_v3.RouteMatch_Prefix{
										Prefix: "/",
									},
								},
								Action: &envoy_config_route_v3.Route_Route{
									Route: &envoy_config_route_v3.RouteAction{
										ClusterSpecifier: &envoy_config_route_v3.RouteAction_Cluster{
											Cluster: translator.UpstreamToClusterName(us.Metadata.Ref()),
										},
									},
								},
							},
						},
					},
				},
			},
		}

		inClusters = []*envoy_config_cluster_v3.Cluster{
			{
				Name: clusterName,
				LoadAssignment: &envoy_config_endpoint_v3.ClusterLoadAssignment{
					ClusterName: clusterName,
					Endpoints: []*envoy_config_endpoint_v3.LocalityLbEndpoints{
						{
							LbEndpoints: []*envoy_config_endpoint_v3.LbEndpoint{
								{
									HostIdentifier: &envoy_config_endpoint_v3.LbEndpoint_Endpoint{
										Endpoint: &envoy_config_endpoint_v3.Endpoint{
											Address: &envoy_config_core_v3.Address{
												Address: &envoy_config_core_v3.Address_SocketAddress{
													SocketAddress: &envoy_config_core_v3.SocketAddress{
														Address: "192.168.0.1",
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		}

		reports = reporter.ResourceReports{}
	})

	It("should update resources properly", func() {
		p := tunneling.NewPlugin()

		generatedClusters, _, _, generatedListeners := p.GeneratedResources(params,
			proxy, inClusters, nil, inRouteConfigurations, nil, reports)
		Expect(reports).To(BeEmpty())
		Expect(generatedClusters).ToNot(BeNil())
		Expect(generatedListeners).ToNot(BeNil())

		validateForwarding(clusterName, inClusters, generatedClusters, generatedListeners)
	})

	Context("UpstreamTlsContext", func() {
		BeforeEach(func() {
			// add an UpstreamTlsContext
			cfg, err := utils.MessageToAny(&envoyauth.UpstreamTlsContext{
				CommonTlsContext: &envoyauth.CommonTlsContext{},
				Sni:              httpProxyHostname,
			})
			Expect(err).ToNot(HaveOccurred())
			inClusters[0].TransportSocket = &envoy_config_core_v3.TransportSocket{
				Name: "",
				ConfigType: &envoy_config_core_v3.TransportSocket_TypedConfig{
					TypedConfig: cfg,
				},
			}

			inRoute := inRouteConfigurations[0].VirtualHosts[0].Routes[0]

			// update route input with duplicate route, the duplicate points to the same upstream as existing
			dupRoute := proto.Clone(inRoute).(*envoy_config_route_v3.Route)
			dupRoute.Name = dupRoute.Name + "-duplicate"
			dupRoute.Action = &envoy_config_route_v3.Route_Route{
				Route: &envoy_config_route_v3.RouteAction{
					ClusterSpecifier: &envoy_config_route_v3.RouteAction_Cluster{
						Cluster: translator.UpstreamToClusterName(us.Metadata.Ref()),
					},
				},
			}
			inRouteConfigurations[0].VirtualHosts[0].Routes = append(inRouteConfigurations[0].VirtualHosts[0].Routes, dupRoute)
		})

		It("should allow multiple routes to same upstream", func() {
			p := tunneling.NewPlugin()
			generatedClusters, _, _, _ := p.GeneratedResources(params, proxy, inClusters,
				nil, inRouteConfigurations, nil, reports)
			Expect(reports).To(BeEmpty())
			Expect(generatedClusters).To(HaveLen(1), "should generate a single cluster for the upstream")
			Expect(generatedClusters[0].GetTransportSocket()).ToNot(BeNil())
		})
	})
	Context("multiple routes and clusters", func() {

		BeforeEach(func() {

			usCopy := &v1.Upstream{}
			us.DeepCopyInto(usCopy)
			usCopy.Metadata.Name = usCopy.Metadata.Name + "-copy"

			// update snapshot with copied upstream, pointing to same HTTP proxy. copied upstream only has different name
			params.Snapshot.Upstreams = append(params.Snapshot.Upstreams, usCopy)

			// update route input with duplicate route, the copy points to the cluster correlating to the copied upstream
			inRoute := inRouteConfigurations[0].VirtualHosts[0].Routes[0]
			cpRoute := proto.Clone(inRoute).(*envoy_config_route_v3.Route)
			cpRoute.Name = cpRoute.Name + "-copy"
			cpRoute.Action = &envoy_config_route_v3.Route_Route{
				Route: &envoy_config_route_v3.RouteAction{
					ClusterSpecifier: &envoy_config_route_v3.RouteAction_Cluster{
						Cluster: translator.UpstreamToClusterName(usCopy.Metadata.Ref()),
					},
				},
			}
			inRouteConfigurations[0].VirtualHosts[0].Routes = append(inRouteConfigurations[0].VirtualHosts[0].Routes, cpRoute)

			// update cluster input such that copied cluster's name matches the copied upstream
			inCluster := inClusters[0]
			cpCluster := proto.Clone(inCluster).(*envoy_config_cluster_v3.Cluster)
			cpCluster.Name = translator.UpstreamToClusterName(usCopy.Metadata.Ref())
			inClusters = append(inClusters, cpCluster)
		})

		It("should namespace generated clusters, avoiding duplicates", func() {
			p := tunneling.NewPlugin()

			generatedClusters, _, _, _ := p.GeneratedResources(params, proxy, inClusters,
				nil, inRouteConfigurations, nil, reports)
			Expect(reports).To(BeEmpty())
			Expect(generatedClusters).To(HaveLen(2), "should generate a cluster for each input route")

			Expect(generatedClusters[0].Name).ToNot(Equal(generatedClusters[1].Name), "should not have name collisions")
			Expect(generatedClusters[0].LoadAssignment.ClusterName).ToNot(Equal(generatedClusters[1].LoadAssignment.ClusterName), "should not route to same cluster")
			cluster1Pipe := generatedClusters[0].LoadAssignment.Endpoints[0].LbEndpoints[0].GetEndpoint().GetAddress().GetPipe().GetPath()
			cluster2Pipe := generatedClusters[1].LoadAssignment.Endpoints[0].LbEndpoints[0].GetEndpoint().GetAddress().GetPipe().GetPath()
			Expect(cluster1Pipe).ToNot(BeEmpty(), "generated cluster should be in-memory pipe to self")
			Expect(cluster2Pipe).ToNot(BeEmpty(), "generated cluster should be in-memory pipe to self")
			Expect(cluster1Pipe).ToNot(Equal(cluster2Pipe), "should not reuse the same pipe to same generated tcp listener")

			// wipe the fields we expect to be different
			generatedClusters[0].Name = ""
			generatedClusters[0].LoadAssignment.ClusterName = ""
			generatedClusters[0].LoadAssignment.Endpoints[0].LbEndpoints[0] = nil
			generatedClusters[1].Name = ""
			generatedClusters[1].LoadAssignment.ClusterName = ""
			generatedClusters[1].LoadAssignment.Endpoints[0].LbEndpoints[0] = nil

			Expect(generatedClusters[0]).To(matchers.MatchProto(generatedClusters[1]), "generated clusters should be identical, barring name, clustername, endpoints")
		})

		It("should namespace generated listeners, avoiding duplicates", func() {
			p := tunneling.NewPlugin()

			_, _, _, generatedListeners := p.GeneratedResources(params, proxy, inClusters,
				nil, inRouteConfigurations, nil, reports)
			Expect(reports).To(BeEmpty())
			Expect(generatedListeners).To(HaveLen(2), "should generate a listener for each input route")

			listener1Pipe := generatedListeners[0].GetAddress().GetPipe().GetPath()
			listener2Pipe := generatedListeners[1].GetAddress().GetPipe().GetPath()
			Expect(listener1Pipe).ToNot(BeEmpty(), "generated listener should be in-memory pipe to self")
			Expect(listener2Pipe).ToNot(BeEmpty(), "generated listener should be in-memory pipe to self")

			Expect(listener1Pipe).ToNot(Equal(listener2Pipe), "should not reuse the same pipe to same generated tcp listener")
			Expect(generatedListeners[0].Name).ToNot(Equal(generatedListeners[1].Name), "should not have name collisions")

			listener1TcpProto, err := utils.AnyToMessage(generatedListeners[0].FilterChains[0].Filters[0].GetTypedConfig())
			Expect(err).ToNot(HaveOccurred())
			Expect(listener1TcpProto).To(BeAssignableToTypeOf(&envoytcp.TcpProxy{}))
			listener1TcpCfg := listener1TcpProto.(*envoytcp.TcpProxy)

			listener2TcpProto, err := utils.AnyToMessage(generatedListeners[1].FilterChains[0].Filters[0].GetTypedConfig())
			Expect(err).To(Not(HaveOccurred()))
			Expect(listener2TcpProto).To(BeAssignableToTypeOf(&envoytcp.TcpProxy{}))
			listener2TcpCfg := listener2TcpProto.(*envoytcp.TcpProxy)

			Expect(listener1TcpCfg.GetStatPrefix()).ToNot(Equal(listener2TcpCfg.GetStatPrefix()), "should not reuse the same stats prefix on generated tcp listeners")

			// wipe the fields we expect to be different
			generatedListeners[0].Name = ""
			generatedListeners[0].Address = nil
			generatedListeners[0].FilterChains[0].Filters[0].ConfigType = nil
			generatedListeners[1].Name = ""
			generatedListeners[1].Address = nil
			generatedListeners[1].FilterChains[0].Filters[0].ConfigType = nil

			Expect(generatedListeners[0]).To(matchers.MatchProto(generatedListeners[1]), "generated listeners should be identical, barring name, address, and tcp stats prefix")
		})
	})

	Context("upstream without a route", func() {
		BeforeEach(func() {
			inRouteConfigurations = []*envoy_config_route_v3.RouteConfiguration{}
		})

		It("should generate resources", func() {
			p := tunneling.NewPlugin()

			newClusters, _, _, newListeners := p.GeneratedResources(params, proxy, inClusters,
				nil, inRouteConfigurations, nil, reports)
			Expect(reports).To(BeEmpty())
			Expect(newClusters).To(HaveLen(1), "should generate a cluster for the upstream")
			Expect(newListeners).To(HaveLen(1), "should generate a listener for the upstream")

			validateForwarding(clusterName, inClusters, newClusters, newListeners)
		})

		It("should not mutate the same cluster multiple times", func() {
			// add the same upstream to the snapshot so we have Upstreams pointing same cluster
			params.Snapshot.Upstreams = append(params.Snapshot.Upstreams, us)

			p := tunneling.NewPlugin()

			newClusters, _, _, newListeners := p.GeneratedResources(params, proxy, inClusters,
				nil, inRouteConfigurations, nil, reports)
			Expect(reports).To(BeEmpty())
			Expect(newClusters).To(HaveLen(1), "should generate a single cluster for the upstream")
			Expect(newListeners).To(HaveLen(1), "should generate a single listener for the upstream")
		})
	})

	Context("upstream without a tunneling configuration", func() {
		BeforeEach(func() {
			us.HttpProxyHostname = nil
		})

		It("should not generate resources", func() {
			p := tunneling.NewPlugin()

			newClusters, _, _, newListeners := p.GeneratedResources(params, proxy, inClusters,
				nil, inRouteConfigurations, nil, reports)
			Expect(reports).To(BeEmpty())
			Expect(newClusters).To(BeNil(), "should not generate a cluster for the upstream")
			Expect(newListeners).To(BeNil(), "should not generate a listener for the upstream")

			Expect(inClusters[0].Name).ToNot(ContainSubstring(tunneling.OriginalClusterSuffix),
				"should not have updated the original cluster name")
			Expect(newClusters).To(BeEmpty(), "should not generate a cluster for the upstream")
			Expect(newListeners).To(BeEmpty(), "should not generate a listener for the upstream")
		})
	})

	Context("tunnel with an ssl configuration", func() {
		BeforeEach(func() {
			us.HttpConnectSslConfig = &ssl.UpstreamSslConfig{
				Sni: "ansni.example.com",
				SslSecrets: &ssl.UpstreamSslConfig_SecretRef{
					SecretRef: &core.ResourceRef{
						Name:      "secret",
						Namespace: "gloo-system",
					},
				},
			}

			params.Snapshot.Secrets = []*v1.Secret{
				{
					Metadata: &core.Metadata{
						Name:      "secret",
						Namespace: "gloo-system",
					},
					Kind: &v1.Secret_Tls{
						Tls: &v1.TlsSecret{
							CertChain:  certChain,
							PrivateKey: privateKey,
							RootCa:     rootCa,
						},
					},
				},
			}
		})

		It("should generate resources", func() {
			p := tunneling.NewPlugin()

			fmt.Println("params: ", params.Snapshot.Upstreams)

			newClusters, _, _, newListeners := p.GeneratedResources(params, proxy, inClusters,
				nil, inRouteConfigurations, nil, reports)
			Expect(reports).To(BeEmpty())
			Expect(newClusters).To(HaveLen(1), "should generate a cluster for the upstream")
			Expect(newListeners).To(HaveLen(1), "should generate a listener for the upstream")

			validateForwarding(clusterName, inClusters, newClusters, newListeners)

			// validate the generated cluster has the correct transport socket
			transportSocket := inClusters[0].GetTransportSocket()
			Expect(transportSocket).ToNot(BeNil(), "should have a transport socket")
			Expect(transportSocket.GetName()).To(Equal("envoy.transport_sockets.tls"))
			Expect(transportSocket.GetTypedConfig()).ToNot(BeNil(), "should have a typed config")

			anyPb := transportSocket.GetTypedConfig()
			msg, err := utils.AnyToMessage(anyPb)
			Expect(err).ToNot(HaveOccurred(), "should be able to convert the any to a message")
			Expect(msg).To(BeAssignableToTypeOf(&envoyauth.UpstreamTlsContext{}),
				"should have an upstream TLS context")

			upstreamTlsContext := msg.(*envoyauth.UpstreamTlsContext)
			Expect(upstreamTlsContext.GetCommonTlsContext()).ToNot(BeNil(), "should have a common TLS context")
			Expect(upstreamTlsContext.GetSni()).To(Equal("ansni.example.com"),
				"should have the correct SNI on the original cluster")
		})

		It("should error if secret not found", func() {
			p := tunneling.NewPlugin()

			params.Snapshot.Secrets = nil

			newClusters, _, _, newListeners := p.GeneratedResources(params, proxy, inClusters,
				nil, inRouteConfigurations, nil, reports)
			Expect(reports).To(HaveLen(1), "should have a report for the upstream")
			Expect(newClusters).To(BeEmpty(), "should not generate a cluster for the upstream")
			Expect(newListeners).To(BeEmpty(), "should not generate a listener for the upstream")

			report, ok := reports[us]
			Expect(ok).To(BeTrue(), "should have a report for the upstream")
			Expect(report.Errors.Error()).
				To(ContainSubstring("SSL secret not found: list did not find secret gloo-system.secret"),
					"should have an error for the missing secret")
		})

		It("should warn if secret not found", func() {
			p := tunneling.NewPlugin()

			params.Snapshot.Secrets = nil
			params.Settings.Gateway.Validation = &v1.GatewayOptions_ValidationOptions{
				WarnMissingTlsSecret: &wrappers.BoolValue{Value: true},
			}

			newClusters, _, _, newListeners := p.GeneratedResources(params, proxy, inClusters,
				nil, inRouteConfigurations, nil, reports)
			Expect(reports).To(HaveLen(1), "should have a report for the upstream")
			Expect(newClusters).To(HaveLen(1), "should generate a cluster for the upstream")
			Expect(newListeners).To(HaveLen(1), "should generate a listener for the upstream")

			report, ok := reports[us]
			Expect(ok).To(BeTrue(), "should have a report for the upstream")
			Expect(report.Warnings).To(HaveLen(1), "should have a warning")
			Expect(report.Warnings[0]).
				To(ContainSubstring("SSL secret not found: list did not find secret gloo-system.secret"),
					"should have a warning for the missing secret")
		})
	})
})

func validateForwarding(
	clusterName string,
	inClusters []*envoy_config_cluster_v3.Cluster,
	generatedClusters []*envoy_config_cluster_v3.Cluster,
	generatedListeners []*envoy_config_listener_v3.Listener,
) {
	adjustedClusterName := clusterName + tunneling.OriginalClusterSuffix

	// generated self tcp cluster should pipe to in memory tcp listener
	selfClusterPipe := generatedClusters[0].GetLoadAssignment().GetEndpoints()[0].GetLbEndpoints()[0].
		GetEndpoint().GetAddress().GetPipe()
	selfListenerPipe := generatedListeners[0].GetAddress().GetPipe()
	Expect(selfClusterPipe).To(Equal(selfListenerPipe), "we should be routing to ourselves")

	// generated listener encapsulates tcp data in HTTP CONNECT and sends to the original upstream w/ new name
	generatedTcpConfig := generatedListeners[0].GetFilterChains()[0].GetFilters()[0].GetTypedConfig()
	typedTcpConfig, err := utils.AnyToMessage(generatedTcpConfig)
	Expect(err).ToNot(HaveOccurred())

	Expect(typedTcpConfig).To(BeAssignableToTypeOf(&envoytcp.TcpProxy{}))
	tcpConfig := typedTcpConfig.(*envoytcp.TcpProxy)

	Expect(tcpConfig.GetCluster()).To(Equal(adjustedClusterName), "should forward to original destination")
	Expect(inClusters[0].Name).To(Equal(adjustedClusterName), "should have updated the original cluster name")
}

// example values copied from projects/gloo/cli/pkg/cmd/create/secret/secret_test.go
var (
	rootCa     = "foo"
	privateKey = `
-----BEGIN PRIVATE KEY-----
MIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQDKtfsD/lq/htGu
28ivzLzzmZe8fOiqBn2eNUj8QVOlHLnaxTZvrqlkTAnPsDvdHCSTQrZDGI5am1/c
tMKodDbflybgskKvZvDi6BF1w1sbZ3kPGZEt0EnykoiKDvc++XuiNUvMQG02B03v
4iurrLa6TLapB8NWtmRJ3n/Pba500zu+uDF6WFkNTEawzuMMl7jLX5gmkDQeBVvX
12IRY0E1zyq7oTNxsyppxvknhAVKdLYsKRVCtwruC0ACjdwYgq47Nn1S7d7b7M8T
cRkcsHqpy4fUvski1ykfX45kt3QcH8XEpCLPA6B74/1YDAGKfdfBKTsZbeuUt5Q3
v8u2/VJDAgMBAAECggEBALVb/0phZLt6WVtCE9kFKgAf6JugfWCxEe5b6gKWR8mv
W7CZRMzCzZaIWTbRi52VMjv2Ma7x51q1LB0ANDAWWYnNZ+EcW4Embln0GrrrnzVz
a+HQlA0TDzXRWAt8vERBXRWQ7VG+SndOLbJy/XNIwOsI+1tbMK8B29QjFuJ0VOL7
CUoG1/0BXw6sloh9UcoqM4nVcNaqocti4K7d3Z0/T1zJzx6o4o89DqBRnzTj0b9O
ge+sRqGFxEejr2hqMJ0jF5BPYUOyt/4kz1gJcRst0NI9dttwTPWWodu2a3eQT1SF
GmYaJpZQAhknQ6LbFMIdGxLUjt9Ej0D/YRjAvjrs5SECgYEA8o5MJoq8gIMbYFuS
UpR9/qhOh1EaR1D6SEKBpStvVXkRR0JHU77h5gmZZNjzeoGzR5rON5iq2kuKNomW
jt4S6BdEPoR+Kqxu/55uVhmdrSmP/tMXN0XTQa63HscmWUsadBm/E1xAaivYG7W9
GyzsVW+HApm6mTmL2S+YQ1aB1esCgYEA1fJPTiFifjAuhPKOLHQZ6WmgdAWNAsDi
nIoM9BSaxZYNo7o2nes2FIa88hIAtlyCzLMNI8r7xkD3XJEsWoe+VX1gbeTjARpm
W6yLMkYvQsF6Urt4Jw1m2+mCbbh/aNfdaL3HrHevU8mMUaP9iRekXYNh24CaNK+v
t2YRuWsCJwkCgYBgP+Mr8CW5AU2dwPihWFde9D6lJ6O75QBMKEf12PSHAFHA6yYO
r1JIzEpYYFbNqCYSJfXqzeQOV6dy2MoryyfJfWIRRNYj7OTm/mFePS/6hOGlBvLR
dh3MlJ4J0pD/IfRPWeAeuJ6/AsLwy/9Mh1kI1gbHG2WWY+WAu4g6QFupHQKBgQCK
ODWMIH1lUPN86Md5aLik15zV2BA1yy+cOoQL3JPxOvQs5s0KUT9rG3FOYtsa9cF7
ReIjUaw/dRFaOGATTMdmq810sf8GY2vlph93p2g5FI5WjM8fS8U8JiwhfqSxs2RT
mug5QEmBNCD3TZ8qxp9l2s+J5Be8GhTHw6WHyN5nIQKBgHseKAYNH0SKMTBmD9tC
+DMhw6Ypxe4VsDBFoDr1Wxpt6SmrZcy7JcBO/jmXY/xwsnGyehdbpsMxX03c7QSf
AmoJCgOtm0FUXc+eybFzgjM9dvB/ZaKRk7LtA2KJjFtMPGwKmLwjC4+cD8xL57Ej
ZEhjfeyucd48M+JNbyMuE2ZC
-----END PRIVATE KEY-----
`
	certChain = `
-----BEGIN CERTIFICATE-----
MIICvDCCAaQCCQDro6ZXybhlYDANBgkqhkiG9w0BAQsFADAgMR4wHAYDVQQDDBVw
ZXRzdG9yZTEuZXhhbXBsZS5jb20wHhcNMTkwNDA1MTkyODQ2WhcNMjAwNDA0MTky
ODQ2WjAgMR4wHAYDVQQDDBVwZXRzdG9yZTEuZXhhbXBsZS5jb20wggEiMA0GCSqG
SIb3DQEBAQUAA4IBDwAwggEKAoIBAQDKtfsD/lq/htGu28ivzLzzmZe8fOiqBn2e
NUj8QVOlHLnaxTZvrqlkTAnPsDvdHCSTQrZDGI5am1/ctMKodDbflybgskKvZvDi
6BF1w1sbZ3kPGZEt0EnykoiKDvc++XuiNUvMQG02B03v4iurrLa6TLapB8NWtmRJ
3n/Pba500zu+uDF6WFkNTEawzuMMl7jLX5gmkDQeBVvX12IRY0E1zyq7oTNxsypp
xvknhAVKdLYsKRVCtwruC0ACjdwYgq47Nn1S7d7b7M8TcRkcsHqpy4fUvski1ykf
X45kt3QcH8XEpCLPA6B74/1YDAGKfdfBKTsZbeuUt5Q3v8u2/VJDAgMBAAEwDQYJ
KoZIhvcNAQELBQADggEBAL8m5TjFEb58MEXKGdbyGZEdS0FpNI+fYBYyxkpU5/z3
06hV2ajisgvHGyGun/HLBDXtWnbNWKpSjiJiS9Kpkv6X73hba6Q3p3prjgdXkpSU
ONozwlMM1SM0dj/O5VULkcW4uhSQJEyIRRPiA8fslqWuIlr5KWWPbdIkDex/9Ddf
oC7D1exclZNVDVmJzYFSxb1js/rSsln11VJ7uyozpk23lrAVGIrtg5Xr4vxqUZHU
TOeFSVH6LMC5j/Fff+bEBhbPxJAI0P7VXaphYh/dMyAEq+xRxm6ssuccgCyvttmz
+6sUivvxaDhUCAzAoLSa5Xgn5eNdsePz6PQ5Vy/Yidg=
-----END CERTIFICATE-----
`
)
