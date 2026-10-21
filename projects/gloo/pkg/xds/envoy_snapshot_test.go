package xds_test

import (
	envoy_config_cluster_v3 "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	envoy_config_core_v3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	envoy_config_endpoint_v3 "github.com/envoyproxy/go-control-plane/envoy/config/endpoint/v3"
	envoy_config_listener_v3 "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	envoy_config_route_v3 "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	envoy_extensions_filters_network_http_connection_manager_v3 "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/network/http_connection_manager/v3"
	"github.com/envoyproxy/go-control-plane/pkg/wellknown"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/solo-io/gloo/projects/gloo/pkg/utils"
	"github.com/solo-io/gloo/projects/gloo/pkg/xds"
	"github.com/solo-io/solo-kit/pkg/api/v1/control-plane/cache"
	"github.com/solo-io/solo-kit/pkg/api/v1/control-plane/resource"
	"github.com/solo-io/solo-kit/pkg/api/v1/control-plane/types"
)

var _ = Describe("EnvoySnapshot", func() {

	It("clones correctly", func() {

		toBeCloned := xds.NewSnapshot("1234",
			[]cache.Resource{resource.NewEnvoyResource(&envoy_config_endpoint_v3.ClusterLoadAssignment{})},
			[]cache.Resource{resource.NewEnvoyResource(&envoy_config_cluster_v3.Cluster{})},
			[]cache.Resource{resource.NewEnvoyResource(&envoy_config_route_v3.Route{})},
			[]cache.Resource{resource.NewEnvoyResource(&envoy_config_listener_v3.Listener{})},
		)

		// Create an identical struct which is guaranteed not to have been touched to compare against
		untouched := xds.NewSnapshot("1234",
			[]cache.Resource{resource.NewEnvoyResource(&envoy_config_endpoint_v3.ClusterLoadAssignment{})},
			[]cache.Resource{resource.NewEnvoyResource(&envoy_config_cluster_v3.Cluster{})},
			[]cache.Resource{resource.NewEnvoyResource(&envoy_config_route_v3.Route{})},
			[]cache.Resource{resource.NewEnvoyResource(&envoy_config_listener_v3.Listener{})},
		)

		clone := toBeCloned.Clone()

		// Verify that original snapshot and clone are identical
		Expect(toBeCloned.Equal(clone.(*xds.EnvoySnapshot))).To(BeTrue())
		Expect(untouched.Equal(clone.(*xds.EnvoySnapshot))).To(BeTrue())

		// Mutate the clone
		clone.GetResources(
			types.EndpointTypeV3,
		).Items[""].(*resource.EnvoyResource).ResourceProto().(*envoy_config_endpoint_v3.ClusterLoadAssignment).ClusterName = "new_endpoint"

		// Verify that original snapshot was not mutated
		Expect(toBeCloned.Equal(clone.(*xds.EnvoySnapshot))).NotTo(BeTrue())
		Expect(toBeCloned.Equal(untouched)).To(BeTrue())
	})

	It("makes an inconsistent snapshot consistent", func() {
		hcm := &envoy_extensions_filters_network_http_connection_manager_v3.HttpConnectionManager{
			StatPrefix: "placeholder",
			RouteSpecifier: &envoy_extensions_filters_network_http_connection_manager_v3.HttpConnectionManager_Rds{
				Rds: &envoy_extensions_filters_network_http_connection_manager_v3.Rds{
					RouteConfigName: "routeName",
				},
			},
		}

		hcmAny, err := utils.MessageToAny(hcm)
		Expect(err).ToNot(HaveOccurred())

		snapshot := xds.NewSnapshot("1234",
			[]cache.Resource{resource.NewEnvoyResource(&envoy_config_endpoint_v3.ClusterLoadAssignment{})},
			[]cache.Resource{resource.NewEnvoyResource(&envoy_config_cluster_v3.Cluster{
				Name:                 "clusterName",
				ClusterDiscoveryType: &envoy_config_cluster_v3.Cluster_Type{Type: envoy_config_cluster_v3.Cluster_EDS},
				EdsClusterConfig: &envoy_config_cluster_v3.Cluster_EdsClusterConfig{
					EdsConfig:   &envoy_config_core_v3.ConfigSource{},
					ServiceName: "edsServiceName",
				},
			})},
			[]cache.Resource{resource.NewEnvoyResource(&envoy_config_route_v3.Route{})},
			[]cache.Resource{resource.NewEnvoyResource(&envoy_config_listener_v3.Listener{
				FilterChains: []*envoy_config_listener_v3.FilterChain{
					{
						Name: "placeholder_filter_chain",
						Filters: []*envoy_config_listener_v3.Filter{
							{
								ConfigType: &envoy_config_listener_v3.Filter_TypedConfig{
									TypedConfig: hcmAny,
								},
								Name: wellknown.HTTPConnectionManager,
							},
						},
					},
				},
			})})

		Expect(snapshot.Consistent()).To(HaveOccurred())
		snapshot.MakeConsistent()
		Expect(snapshot.Consistent()).NotTo(HaveOccurred())
	})

	// An EDS subscription is satisfied by name: a proxy asks for the names its
	// clusters reference, and matches what comes back by the name each resource
	// gives for itself. A response carrying a resource under the right map key
	// but identifying itself as something else satisfies the cache and not the
	// proxy -- it is delivered, ignored, and never retried, because from the
	// cache's side the subscription has been answered.
	//
	// These specs assert that contract rather than the field assignment that
	// implements it. They derive the expected names from the clusters under
	// test, so they cannot drift from the fixture and cannot be satisfied by
	// restating the implementation.
	Describe("MakeConsistent and EDS subscriptions", func() {
		edsCluster := func(name, serviceName string) cache.Resource {
			return resource.NewEnvoyResource(&envoy_config_cluster_v3.Cluster{
				Name:                 name,
				ClusterDiscoveryType: &envoy_config_cluster_v3.Cluster_Type{Type: envoy_config_cluster_v3.Cluster_EDS},
				EdsClusterConfig: &envoy_config_cluster_v3.Cluster_EdsClusterConfig{
					EdsConfig:   &envoy_config_core_v3.ConfigSource{},
					ServiceName: serviceName,
				},
			})
		}
		assignment := func(name string, addresses ...string) cache.Resource {
			lbs := make([]*envoy_config_endpoint_v3.LbEndpoint, 0, len(addresses))
			for _, address := range addresses {
				lbs = append(lbs, &envoy_config_endpoint_v3.LbEndpoint{
					HostIdentifier: &envoy_config_endpoint_v3.LbEndpoint_Endpoint{
						Endpoint: &envoy_config_endpoint_v3.Endpoint{
							Address: &envoy_config_core_v3.Address{
								Address: &envoy_config_core_v3.Address_SocketAddress{
									SocketAddress: &envoy_config_core_v3.SocketAddress{Address: address},
								},
							},
						},
					},
				})
			}
			return resource.NewEnvoyResource(&envoy_config_endpoint_v3.ClusterLoadAssignment{
				ClusterName: name,
				Endpoints:   []*envoy_config_endpoint_v3.LocalityLbEndpoints{{LbEndpoints: lbs}},
			})
		}

		// subscribed returns the EDS names the clusters ask for, and answered
		// returns the names the endpoint resources give for themselves. A proxy
		// resolves its subscriptions against the second set, not the map keys.
		subscribed := func(snapshot *xds.EnvoySnapshot) []string {
			var names []string
			for _, res := range snapshot.GetResources(types.ClusterTypeV3).Items {
				cluster := res.ResourceProto().(*envoy_config_cluster_v3.Cluster)
				if serviceName := cluster.GetEdsClusterConfig().GetServiceName(); serviceName != "" {
					names = append(names, serviceName)
				}
			}
			return names
		}
		answered := func(snapshot *xds.EnvoySnapshot) []string {
			var names []string
			for _, res := range snapshot.GetResources(types.EndpointTypeV3).Items {
				names = append(names, res.Self().Name)
			}
			return names
		}

		It("answers every EDS name its clusters subscribe to", func() {
			// The Gloo shape: a cluster never serves EDS under its own name.
			snapshot := xds.NewSnapshot("1234",
				[]cache.Resource{},
				[]cache.Resource{
					edsCluster("backend-a_ns", "backend-a_ns-9302032437900230978"),
					edsCluster("backend-b_ns", "backend-b_ns-1797782275941751787"),
				},
				[]cache.Resource{},
				[]cache.Resource{})

			snapshot.MakeConsistent()

			Expect(answered(snapshot)).To(ConsistOf(subscribed(snapshot)),
				"a subscription is resolved by the name a resource gives for itself, so every "+
					"name the clusters reference must appear among the names the endpoints claim")
		})

		It("answers by self-reported name even where that equals the cluster name", func() {
			// The case the previous implementation got right by coincidence,
			// because there the two names are the same string. Kept so the fix
			// cannot regress it.
			snapshot := xds.NewSnapshot("1234",
				[]cache.Resource{},
				[]cache.Resource{edsCluster("selfNamed", "selfNamed")},
				[]cache.Resource{},
				[]cache.Resource{})

			snapshot.MakeConsistent()

			Expect(answered(snapshot)).To(ConsistOf(subscribed(snapshot)))
		})

		It("leaves a real assignment alone rather than replacing it with a placeholder", func() {
			// Guards the other direction: the repair must not discard endpoints
			// that were published, which would turn a working cluster into an
			// empty one.
			snapshot := xds.NewSnapshot("1234",
				[]cache.Resource{assignment("backend-a_ns-9302032437900230978", "10.0.0.1")},
				[]cache.Resource{edsCluster("backend-a_ns", "backend-a_ns-9302032437900230978")},
				[]cache.Resource{},
				[]cache.Resource{})

			snapshot.MakeConsistent()

			Expect(answered(snapshot)).To(ConsistOf(subscribed(snapshot)))
			kept := snapshot.GetResources(types.EndpointTypeV3).Items["backend-a_ns-9302032437900230978"]
			Expect(kept.ResourceProto().(*envoy_config_endpoint_v3.ClusterLoadAssignment).
				GetEndpoints()[0].GetLbEndpoints()).To(HaveLen(1),
				"a published assignment must survive the consistency repair")
		})

		It("repairs only the clusters that are missing an assignment", func() {
			// Mixed snapshot: one cluster answered, one not. Both must end up
			// answered, and the published one must keep its endpoints.
			snapshot := xds.NewSnapshot("1234",
				[]cache.Resource{assignment("backend-a_ns-9302032437900230978", "10.0.0.1", "10.0.0.2")},
				[]cache.Resource{
					edsCluster("backend-a_ns", "backend-a_ns-9302032437900230978"),
					edsCluster("backend-b_ns", "backend-b_ns-1797782275941751787"),
				},
				[]cache.Resource{},
				[]cache.Resource{})

			snapshot.MakeConsistent()

			Expect(answered(snapshot)).To(ConsistOf(subscribed(snapshot)))
			published := snapshot.GetResources(types.EndpointTypeV3).Items["backend-a_ns-9302032437900230978"]
			Expect(published.ResourceProto().(*envoy_config_endpoint_v3.ClusterLoadAssignment).
				GetEndpoints()[0].GetLbEndpoints()).To(HaveLen(2))
			placeholder := snapshot.GetResources(types.EndpointTypeV3).Items["backend-b_ns-1797782275941751787"]
			Expect(placeholder.ResourceProto().(*envoy_config_endpoint_v3.ClusterLoadAssignment).
				GetEndpoints()).To(BeEmpty(), "a placeholder carries no endpoints")
		})
	})
})
