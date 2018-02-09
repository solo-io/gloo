package translator

import (
	envoyapi "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	envoycore "github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	envoyendpoints "github.com/envoyproxy/go-control-plane/envoy/api/v2/endpoint"
	envoycache "github.com/envoyproxy/go-control-plane/pkg/cache"

	"github.com/solo-io/glue/internal/reporter"
	"github.com/solo-io/glue/pkg/api/types/v1"
	"github.com/solo-io/glue/pkg/endpointdiscovery"
	"github.com/solo-io/glue/pkg/plugin2"
	"github.com/solo-io/glue/pkg/secretwatcher"
)

type Translator struct {
	upstreamPlugins []plugin.UpstreamPlugin
	routePlugins    []plugin.RoutePlugin
}

func NewTranslator(upstreamPlugins []plugin.UpstreamPlugin, routePlugins []plugin.RoutePlugin) *Translator {
	// special routing must be done for upstream plugins that support functions
	var functionPlugins []plugin.FunctionPlugin
	for _, upstreamPlugin := range upstreamPlugins {
		if functionPlugin, ok := upstreamPlugin.(plugin.FunctionPlugin); ok {
			functionPlugins = append(functionPlugins, functionPlugin)
		}
	}
	if len(functionPlugins) > 0 {
		upstreamPlugins = append([]plugin.UpstreamPlugin{&functionRouterPlugin{
			functionPlugins: functionPlugins,
		}}, upstreamPlugins...)
	}
	return &Translator{
		upstreamPlugins: upstreamPlugins,
		routePlugins:    routePlugins,
	}
}

func (t *Translator) Translate(cfg v1.Config,
	secretMap secretwatcher.SecretMap,
	endpoints endpointdiscovery.EndpointGroups) (*envoycache.Snapshot, []reporter.ConfigObjectReport) {

	// endpoints
	clusterLoadAssignments := computeClusterEndpoints(cfg.Upstreams, endpoints)

	// clusters

}

func computeClusterEndpoints(upstreams []v1.Upstream, endpoints endpointdiscovery.EndpointGroups) []*envoyapi.ClusterLoadAssignment {
	var clusterEndpointAssignments []*envoyapi.ClusterLoadAssignment
	for _, us := range upstreams {
		// if there is an endpoint group for this upstream,
		// it's using eds and we need to create a load assignment for it
		if endpointGroup, ok := endpoints[us.Name]; ok {
			loadAssignment := loadAssignmentForCluster(us.Name, endpointGroup)
			clusterEndpointAssignments = append(clusterEndpointAssignments, loadAssignment)
		}
	}
	return clusterEndpointAssignments
}

func loadAssignmentForCluster(clusterName string, addresses []endpointdiscovery.Endpoint) *envoyapi.ClusterLoadAssignment {
	var endpoints []envoyendpoints.LbEndpoint
	for _, addr := range addresses {
		lbEndpoint := envoyendpoints.LbEndpoint{
			Endpoint: &envoyendpoints.Endpoint{
				Address: &envoycore.Address{
					Address: &envoycore.Address_SocketAddress{
						SocketAddress: &envoycore.SocketAddress{
							Protocol: envoycore.TCP,
							Address:  addr.Address,
							PortSpecifier: &envoycore.SocketAddress_PortValue{
								PortValue: uint32(addr.Port),
							},
						},
					},
				},
			},
		}
		endpoints = append(endpoints, lbEndpoint)
	}

	return &envoyapi.ClusterLoadAssignment{
		ClusterName: clusterName,
		Endpoints: []envoyendpoints.LocalityLbEndpoints{{
			LbEndpoints: endpoints,
		}},
	}
}

func (t *Translator) computeClusters(inputs plugin.Inputs, upstreams []v1.Upstream, endpoints endpointdiscovery.EndpointGroups) {

}

func (t *Translator) computeCluster(inputs plugin.Inputs, upstream v1.Upstream, edsCluster bool) *envoyapi.Cluster {
	outputCluster := &envoyapi.Cluster{
		Name: upstream.Name,
	}
	if edsCluster {
		outputCluster.Type = envoyapi.Cluster_EDS
	}
	for _, upstreamPlugin := range t.upstreamPlugins {
		if upstreamPlugin.Type() == upstream.Type
	}
}
