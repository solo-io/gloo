package validation

import (
	"context"
	"log"
	"os"
	"os/exec"

	envoy_config_bootstrap_v3 "github.com/envoyproxy/go-control-plane/envoy/config/bootstrap/v3"
	envoy_config_cluster_v3 "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	envoy_config_core_v3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	envoy_config_endpoint_v3 "github.com/envoyproxy/go-control-plane/envoy/config/endpoint/v3"
	envoy_config_listener_v3 "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	envoy_config_route_v3 "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	envoy_extensions_filters_network_http_connection_manager_v3 "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/network/http_connection_manager/v3"
	"github.com/rotisserie/eris"
	"github.com/solo-io/gloo/projects/gloo/pkg/utils"
	"github.com/solo-io/go-utils/contextutils"
	envoycache "github.com/solo-io/solo-kit/pkg/api/v1/control-plane/cache"
	"github.com/solo-io/solo-kit/pkg/api/v1/control-plane/resource"
	"github.com/solo-io/solo-kit/pkg/api/v1/control-plane/types"
	"google.golang.org/protobuf/encoding/protojson"
)

const defaultEnvoyPath = "/usr/local/bin/envoy"

func getEnvoyPath() string {
	ep := os.Getenv("ENVOY_BINARY_PATH")
	if len(ep) == 0 {
		ep = defaultEnvoyPath
	}
	return ep
}

func ValidateEntireBootstrap(
	ctx context.Context,
	snap envoycache.Snapshot,
) error {
	bootstrapYaml, err := buildEntireBootstrap(ctx, snap)
	if err != nil {
		contextutils.LoggerFrom(ctx).Error(err)
		return err
	}
	log.Println("validating with envoy")
	log.Println(bootstrapYaml)

	envoyPath := getEnvoyPath()
	validateCmd := exec.Command(envoyPath, "--mode", "validate", "--config-yaml", bootstrapYaml, "-l", "critical", "--log-format", "%v")
	if output, err := validateCmd.CombinedOutput(); err != nil {
		if os.IsNotExist(err) {
			// log a warning and return nil; will allow users to continue to run Gloo locally without
			// relying on the Gloo container with Envoy already published to the expected directory
			contextutils.LoggerFrom(ctx).Warnf("Unable to validate envoy configuration using envoy at %v; "+
				"skipping additional validation of Gloo config.", envoyPath)
			return nil
		}
		return eris.Errorf("envoy validation mode output: %v, error: %v", string(output), err)
	}
	return nil
}

// buildEntireBootstrap queries the gloo xds dump using cli code and converts the output
// into valid bootstrap json.
func buildEntireBootstrap(
	ctx context.Context,
	snap envoycache.Snapshot,
) (string, error) {

	// get the resources we're going to need
	listeners := snap.GetResources(types.ListenerTypeV3).Items
	clusters := snap.GetResources(types.ClusterTypeV3).Items
	routes := snap.GetResources(types.RouteTypeV3).Items
	endpoints := snap.GetResources(types.EndpointTypeV3).Items
	routedCluster := map[string]struct{}{}
	for _, v := range listeners {
		l := v.ResourceProto().(*envoy_config_listener_v3.Listener)
		for _, fc := range l.FilterChains {
			for _, f := range fc.Filters {

				if f.GetTypedConfig().TypeUrl == "type.googleapis.com/envoy.extensions.filters.network.http_connection_manager.v3.HttpConnectionManager" {
					hcmAny, err := utils.AnyToMessage(f.GetTypedConfig())
					if err != nil {
						return "", err
					}
					if hcm, ok := hcmAny.(*envoy_extensions_filters_network_http_connection_manager_v3.HttpConnectionManager); ok {
						if n := hcm.GetRds().RouteConfigName; n != "" {
							// find route
							for _, rt := range routes {
								r := rt.ResourceProto().(*envoy_config_route_v3.RouteConfiguration)
								if r.Name == n {
									hcm.RouteSpecifier = &envoy_extensions_filters_network_http_connection_manager_v3.HttpConnectionManager_RouteConfig{
										RouteConfig: r,
									}
									for _, v := range r.VirtualHosts {
										for _, r := range v.Routes {
											if r.GetRoute() != nil {
												if c := r.GetRoute().GetCluster(); c != "" {
													routedCluster[c] = struct{}{}
												}
												if c := r.GetRoute().GetWeightedClusters().GetClusters(); len(c) != 0 {
													for _, c := range c {
														routedCluster[c.Name] = struct{}{}
													}
												}
											}
										}
									}

									hcmAny, err := utils.MessageToAny(hcm)
									if err != nil {
										return "", err
									}

									f.ConfigType = &envoy_config_listener_v3.Filter_TypedConfig{
										TypedConfig: hcmAny,
									}
								}
							}
						}

					}
				}

			}
		}
	}

	for _, cl := range clusters {
		c := cl.ResourceProto().(*envoy_config_cluster_v3.Cluster)
		// remove existing clusters
		delete(routedCluster, c.Name)
		if c.GetEdsClusterConfig() != nil {
			name := c.Name
			if n2 := c.GetEdsClusterConfig().GetServiceName(); n2 != "" {
				name = n2
			}

			// find endpoints
			for _, en := range endpoints {
				e := en.ResourceProto().(*envoy_config_endpoint_v3.ClusterLoadAssignment)
				if e.ClusterName == name {
					c.LoadAssignment = e
					c.EdsClusterConfig = nil
					c.ClusterDiscoveryType = &envoy_config_cluster_v3.Cluster_Type{
						Type: envoy_config_cluster_v3.Cluster_STRICT_DNS,
					}
				}
			}

		}
	}
	// routedClusters now contains clusters that have a route but no cluster.
	// these are effectively blackhole clusters. in static mode, envoy won't start without them
	// so just add them as blackhole clusters
	for c := range routedCluster {
		clusters[c] = resource.NewEnvoyResource(&envoy_config_cluster_v3.Cluster{
			Name: c,
			ClusterDiscoveryType: &envoy_config_cluster_v3.Cluster_Type{
				Type: envoy_config_cluster_v3.Cluster_STATIC,
			},
			LoadAssignment: &envoy_config_endpoint_v3.ClusterLoadAssignment{
				ClusterName: c,
				Endpoints:   []*envoy_config_endpoint_v3.LocalityLbEndpoints{},
			},
		})
	}

	var concreteListeners []*envoy_config_listener_v3.Listener
	for _, v := range listeners {
		l := v.ResourceProto().(*envoy_config_listener_v3.Listener)
		concreteListeners = append(concreteListeners, l)
	}
	var concreteClusters []*envoy_config_cluster_v3.Cluster
	for _, v := range clusters {
		c := v.ResourceProto().(*envoy_config_cluster_v3.Cluster)
		concreteClusters = append(concreteClusters, c)
	}
	bs := envoy_config_bootstrap_v3.Bootstrap{
		Node: &envoy_config_core_v3.Node{
			Id:      "test-id",
			Cluster: "test-cluster",
		},
		Admin: &envoy_config_bootstrap_v3.Admin{
			Address: &envoy_config_core_v3.Address{
				Address: &envoy_config_core_v3.Address_SocketAddress{
					SocketAddress: &envoy_config_core_v3.SocketAddress{
						Address: "127.0.0.1",
						PortSpecifier: &envoy_config_core_v3.SocketAddress_PortValue{
							PortValue: 19000,
						},
					},
				},
			},
		},
		StaticResources: &envoy_config_bootstrap_v3.Bootstrap_StaticResources{
			Listeners: concreteListeners,
			Clusters:  concreteClusters,
		},
	}

	// jsonpb marshal
	j, err := protojson.MarshalOptions{Multiline: true, Indent: "  "}.Marshal(&bs)
	if err != nil {
		return "", err
	}

	return string(j), nil

}
