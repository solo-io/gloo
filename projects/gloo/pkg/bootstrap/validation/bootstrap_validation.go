package validation

import (
	"context"
	"errors"
	"fmt"
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

var (
	errNoHcm = eris.New("no HttpConnectionManager found")
)

func getEnvoyPath() string {
	ep := os.Getenv("ENVOY_BINARY_PATH")
	if len(ep) == 0 {
		ep = defaultEnvoyPath
	}
	return ep
}

func ValidateBootstrap(ctx context.Context, bootstrap string) error {
	logger := contextutils.LoggerFrom(ctx)
	fmt.Println("validating with envoy:")
	fmt.Println(bootstrap)

	envoyPath := getEnvoyPath()
	validateCmd := exec.Command(envoyPath, "--mode", "validate", "--config-yaml", bootstrap, "-l", "critical", "--log-format", "%v")
	output, err := validateCmd.CombinedOutput()
	if err != nil {
		if os.IsNotExist(err) {
			// log a warning and return nil; will allow users to continue to run Gloo locally without
			// relying on the Gloo container with Envoy already published to the expected directory
			logger.Warnf("Unable to validate envoy configuration using envoy at %v; "+
				"skipping additional validation of Gloo config.", envoyPath)
			return nil
		}
		return eris.Errorf("envoy validation mode output: %v, error: %v", string(output), err)
	}
	fmt.Println("envoy validation output:")
	fmt.Println(string(output))
	return nil
}

// ValidateSnapshot accepts an xDS snapshot, clones it, and does the necessary
// conversions to imitate the same config being provided as static bootsrap config to
// Envoy, then executes Envoy in validate mode to ensure the config is valid.
func ValidateSnapshot(
	ctx context.Context,
	snap envoycache.Snapshot,
) error {
	// THIS IS CRITICAL SO WE DO NOT INTERFERE WITH THE CONTROL PLANE.
	snap = snap.Clone()

	logger := contextutils.LoggerFrom(ctx)

	bootstrap, err := BootstrapFromSnapshot(ctx, snap)
	if err != nil {
		logger.Error(err)
		return err
	}

	// jsonpb marshal
	bootstrapBytes, err := protojson.MarshalOptions{Multiline: true, Indent: "  "}.Marshal(bootstrap)
	if err != nil {
		logger.Error(err)
		return err
	}

	bootstrapJson := string(bootstrapBytes)

	return ValidateBootstrap(ctx, bootstrapJson)

}

// BootstrapFromSnapshot accepts an xds Snapshot and converts it into valid bootstrap json.
func BootstrapFromSnapshot(
	ctx context.Context,
	snap envoycache.Snapshot,
) (*envoy_config_bootstrap_v3.Bootstrap, error) {

	// Get the resources we're going to need.
	listeners := snap.GetResources(types.ListenerTypeV3).Items
	clusters := snap.GetResources(types.ClusterTypeV3).Items
	routes := snap.GetResources(types.RouteTypeV3).Items
	endpoints := snap.GetResources(types.EndpointTypeV3).Items

	// This map will hold the aggregate of all cluster names that are routed to
	// by a FilterChain.
	routedCluster := map[string]struct{}{}

	// Gather up all of the clusters that we target with RouteConfigs that are associated with a FilterChain.
	for _, v := range listeners {
		l := v.ResourceProto().(*envoy_config_listener_v3.Listener)
		for _, fc := range l.GetFilterChains() {
			// Get the HttpConnectionManager for this FilterChain if it exists.
			hcm, f, err := getHcmForFilterChain(fc)
			if err != nil {
				// If we just don't have an hcm on this filter chain, skip to the next one.
				if errors.Is(err, errNoHcm) {
					continue
				}
				// If we encountered any other error, fail loudly.
				return nil, err
			}

			// We use Route Discovery Service (RDS) in lieu of static route table config, so we
			// need to get the RouteConfiguration name to lookup in our Snapshot-provided routes,
			// which contain what we serve over RDS.
			routeConfigName := hcm.GetRds().GetRouteConfigName()
			if routeConfigName == "" {
				continue
			}
			// Find matching route config from snapshot.
			for _, rt := range routes {
				r, ok := rt.ResourceProto().(*envoy_config_route_v3.RouteConfiguration)
				if !ok {
					return nil, eris.New("found route with wrong type")
				}

				if r.GetName() != routeConfigName {
					// These aren't the routes you're looking for.
					continue
				}

				// Add clusters targeted by routes on this config to our aggregate list of all targeted clusters
				findTargetedClusters(r, routedCluster)

				// We need to add our route table as a static config to this hcm instead of
				// relying on RDS, the we pack it back up and set it back on the filter chain.
				if err = setStaticRouteConfig(f, hcm, r); err != nil {
					return nil, err
				}
			}

		}
	}

	// Next, we will look through our Snapshot's clusters and delete the ones which are
	// already routed to.
	for _, cl := range clusters {
		c, ok := cl.ResourceProto().(*envoy_config_cluster_v3.Cluster)
		if !ok {
			return nil, eris.New("found cluster with wrong type")
		}

		delete(routedCluster, c.GetName())

		// We use Endpoint Discovery Service (EDS) in lieu of static endpoint config, so we
		// need to get the EDS ServiceName name to lookup in our Snapshot-provided endpoints,
		// which contain what we serve over EDS.
		if c.GetEdsClusterConfig() != nil {
			clusterName := c.GetName()
			if edsServiceName := c.GetEdsClusterConfig().GetServiceName(); edsServiceName != "" {
				clusterName = edsServiceName
			}

			// Find endpoints matching our EDS config and convert the cluster to use
			// static endpoint config matching that which would have been served over EDS.
			for _, en := range endpoints {
				e, ok := en.ResourceProto().(*envoy_config_endpoint_v3.ClusterLoadAssignment)
				if !ok {
					return nil, eris.New("found endpoint with wrong type")
				}
				if e.GetClusterName() == clusterName {
					c.LoadAssignment = e
					c.EdsClusterConfig = nil
					c.ClusterDiscoveryType = &envoy_config_cluster_v3.Cluster_Type{
						Type: envoy_config_cluster_v3.Cluster_STRICT_DNS,
					}
				}
			}

		}
	}
	// We now need to find clusters which do not exist, even though they are targeted by
	// a route. In static mode, envoy won't start without these. At this point in the
	// processing, routedClusters holds this list, so we range over the map and create
	// blackhole clusters for these routes to target.
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

	// Because our snapshot provided us with abstractions over our resources,
	// we need to convert them to slices of pointers to their concrete types
	// in order to include them in the bootstrap struct.
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

	// Finally, we build the static bootstrap config holding all of our converted xDS config.
	bs := envoy_config_bootstrap_v3.Bootstrap{
		Node: &envoy_config_core_v3.Node{
			Id:      "validation-id",
			Cluster: "validation-cluster",
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

	return &bs, nil

}

func getHcmForFilterChain(fc *envoy_config_listener_v3.FilterChain) (
	*envoy_extensions_filters_network_http_connection_manager_v3.HttpConnectionManager,
	*envoy_config_listener_v3.Filter,
	error,
) {

	for _, f := range fc.GetFilters() {

		if f.GetTypedConfig().GetTypeUrl() == "type.googleapis.com/envoy.extensions.filters.network.http_connection_manager.v3.HttpConnectionManager" {
			hcmAny, err := utils.AnyToMessage(f.GetTypedConfig())
			if err != nil {
				return nil, nil, err
			}
			if hcm, ok := hcmAny.(*envoy_extensions_filters_network_http_connection_manager_v3.HttpConnectionManager); ok {
				return hcm, f, nil
			} else {
				return nil, nil, eris.New("hcm config casting to concrete failed; likely wrong type")
			}
		}
	}
	return nil, nil, errNoHcm
}

func findTargetedClusters(r *envoy_config_route_v3.RouteConfiguration, routedCluster map[string]struct{}) {
	for _, v := range r.GetVirtualHosts() {
		for _, r := range v.GetRoutes() {
			if r.GetRoute() == nil {
				continue
			}

			if c := r.GetRoute().GetCluster(); c != "" {
				routedCluster[c] = struct{}{}
			}
			if wc := r.GetRoute().GetWeightedClusters().GetClusters(); len(wc) != 0 {
				for _, c := range wc {
					routedCluster[c.GetName()] = struct{}{}
				}
			}
		}
	}
}

func setStaticRouteConfig(
	f *envoy_config_listener_v3.Filter,
	hcm *envoy_extensions_filters_network_http_connection_manager_v3.HttpConnectionManager,
	r *envoy_config_route_v3.RouteConfiguration,
) error {
	hcm.RouteSpecifier = &envoy_extensions_filters_network_http_connection_manager_v3.HttpConnectionManager_RouteConfig{
		RouteConfig: r,
	}

	hcmAny, err := utils.MessageToAny(hcm)
	if err != nil {
		return err
	}

	f.ConfigType = &envoy_config_listener_v3.Filter_TypedConfig{
		TypedConfig: hcmAny,
	}
	return nil
}
