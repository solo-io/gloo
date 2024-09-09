package bootstrap

import (
	"bytes"
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
	"github.com/envoyproxy/go-control-plane/pkg/wellknown"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	anypb "github.com/golang/protobuf/ptypes/any"
	"github.com/rotisserie/eris"
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/xdsinspection"
	"github.com/solo-io/gloo/projects/gloo/pkg/utils"
	"github.com/solo-io/go-utils/contextutils"
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

func ValidateBootstrap(
	ctx context.Context,
	filterName string,
	msg proto.Message,
) error {
	bootstrapYaml, err := buildPerFilterBootstrapYaml(filterName, msg)
	if err != nil {
		return err
	}

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

func buildPerFilterBootstrapYaml(filterName string, msg proto.Message) (string, error) {

	typedFilter, err := utils.MessageToAny(msg)
	if err != nil {
		return "", err
	}
	vhosts := []*envoy_config_route_v3.VirtualHost{
		{
			Name:    "placeholder_host",
			Domains: []string{"*"},
			TypedPerFilterConfig: map[string]*anypb.Any{
				filterName: {
					TypeUrl: typedFilter.GetTypeUrl(),
					Value:   typedFilter.GetValue(),
				},
			},
		},
	}

	rc := &envoy_config_route_v3.RouteConfiguration{VirtualHosts: vhosts}

	hcm := &envoy_extensions_filters_network_http_connection_manager_v3.HttpConnectionManager{
		StatPrefix:     "placeholder",
		RouteSpecifier: &envoy_extensions_filters_network_http_connection_manager_v3.HttpConnectionManager_RouteConfig{RouteConfig: rc},
	}

	hcmAny, err := utils.MessageToAny(hcm)
	if err != nil {
		return "", err
	}
	bootstrap := &envoy_config_bootstrap_v3.Bootstrap{
		Node: &envoy_config_core_v3.Node{
			Id:      "imspecial",
			Cluster: "doesntmatter",
		},
		StaticResources: &envoy_config_bootstrap_v3.Bootstrap_StaticResources{
			Listeners: []*envoy_config_listener_v3.Listener{
				{
					Name: "placeholder_listener",
					Address: &envoy_config_core_v3.Address{
						Address: &envoy_config_core_v3.Address_SocketAddress{SocketAddress: &envoy_config_core_v3.SocketAddress{
							Address:       "0.0.0.0",
							PortSpecifier: &envoy_config_core_v3.SocketAddress_PortValue{PortValue: 8081},
						}},
					},
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
				},
			},
		},
	}

	buf := &bytes.Buffer{}
	marshaler := &jsonpb.Marshaler{
		OrigName: true,
	}
	marshaler.Marshal(buf, bootstrap)
	json := string(buf.Bytes())
	return json, nil // returns a json, but json is valid yaml
}

func ValidateEntireBootstrap(ctx context.Context, port int, ns string, proxyName string) error {
	bootstrapYaml, err := buildEntireBootstrap(port, ns, proxyName)
	if err != nil {
		return err
	}

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
func buildEntireBootstrap(port int, ns, proxyName string) (string, error) {

	dump, err := xdsinspection.GetGlooXdsDump(context.Background(), proxyName, ns, true)
	if err != nil {
		log.Fatal(err)
	}

	// get listeners and clusters
	clusters := dump.Clusters
	listeners := dump.Listeners
	routedCluster := map[string]struct{}{}
	for i := range listeners {
		l := &listeners[i]
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
							for j := range dump.Routes {
								r := &dump.Routes[j]
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

	for i := range clusters {
		c := &clusters[i]
		// remove existing clusters
		delete(routedCluster, c.Name)
		if c.GetEdsClusterConfig() != nil {
			name := c.Name
			if n2 := c.GetEdsClusterConfig().GetServiceName(); n2 != "" {
				name = n2
			}

			// find endpoints
			for j := range dump.Endpoints {
				e := &dump.Endpoints[j]
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
		clusters = append(clusters, envoy_config_cluster_v3.Cluster{
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
			Listeners: toPtr(listeners),
			Clusters:  toPtr(clusters),
		},
	}

	// jsonpb marshal
	j, err := protojson.MarshalOptions{Multiline: true, Indent: "  "}.Marshal(&bs)
	if err != nil {
		return "", err
	}

	return string(j), nil

}

func toPtr[T any](s []T) []*T {
	ptrs := make([]*T, len(s))
	for i, v := range s {
		ptrs[i] = &v
	}
	return ptrs
}
