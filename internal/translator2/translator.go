package translator

import (
	envoyapi "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	envoycore "github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	envoyendpoints "github.com/envoyproxy/go-control-plane/envoy/api/v2/endpoint"
	envoyroute "github.com/envoyproxy/go-control-plane/envoy/api/v2/route"
	envoycache "github.com/envoyproxy/go-control-plane/pkg/cache"
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
	"github.com/solo-io/glue/internal/pkg/envoy"

	"github.com/solo-io/glue/internal/reporter"
	"github.com/solo-io/glue/pkg/api/types/v1"
	"github.com/solo-io/glue/pkg/endpointdiscovery"
	"github.com/solo-io/glue/pkg/plugin2"
	"github.com/solo-io/glue/pkg/secretwatcher"
)

const rdsName = "glue-rds"

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
		// the function router plugin must be initialized for any function plugins
		// since it operates on both upstreams and routes, it must be added to both
		// groups of plugins
		functionRouter := &functionRouterPlugin{
			functionPlugins: functionPlugins,
		}
		upstreamPlugins = append([]plugin.UpstreamPlugin{functionRouter}, upstreamPlugins...)
		routePlugins = append([]plugin.RoutePlugin{functionRouter}, routePlugins...)
	}
	return &Translator{
		upstreamPlugins: upstreamPlugins,
		routePlugins:    routePlugins,
	}
}

func (t *Translator) Translate(cfg v1.Config,
	secrets secretwatcher.SecretMap,
	endpoints endpointdiscovery.EndpointGroups) (*envoycache.Snapshot, []reporter.ConfigObjectReport) {

	// endpoints
	clusterLoadAssignments := computeClusterEndpoints(cfg.Upstreams, endpoints)

	// clusters
	clusters, clusterReports := t.computeClusters(cfg, secrets, endpoints)

	// virtualhosts
	virtualHosts, virtualHostReports := t.computeVirtualHosts(cfg)

	routeConfig := &envoyapi.RouteConfiguration{
		Name:         rdsName,
		VirtualHosts: virtualHosts,
	}

	//listener :=
}

// Endpoints

func computeClusterEndpoints(upstreams []v1.Upstream, endpoints endpointdiscovery.EndpointGroups) []*envoyapi.ClusterLoadAssignment {
	var clusterEndpointAssignments []*envoyapi.ClusterLoadAssignment
	for _, upstream := range upstreams {
		// if there is an endpoint group for this upstream,
		// it's using eds and we need to create a load assignment for it
		if endpointGroup, ok := endpoints[upstream.Name]; ok {
			loadAssignment := loadAssignmentForCluster(upstream.Name, endpointGroup)
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

// Clusters

func (t *Translator) computeClusters(cfg v1.Config, secrets secretwatcher.SecretMap, endpoints endpointdiscovery.EndpointGroups) ([]*envoyapi.Cluster, []reporter.ConfigObjectReport) {
	var (
		reports  []reporter.ConfigObjectReport
		clusters []*envoyapi.Cluster
	)
	for _, upstream := range cfg.Upstreams {
		_, edsCluster := endpoints[upstream.Name]
		cluster, err := t.computeCluster(cfg, secrets, upstream, edsCluster)
		clusters = append(clusters, cluster)
		reports = append(reports, createUpstreamReport(upstream, err))
	}
	return clusters, reports
}

func (t *Translator) computeCluster(cfg v1.Config, secrets secretwatcher.SecretMap, upstream v1.Upstream, edsCluster bool) (*envoyapi.Cluster, error) {
	out := &envoyapi.Cluster{
		Name: upstream.Name,
	}
	if edsCluster {
		out.Type = envoyapi.Cluster_EDS
	}
	var upstreamErrors *multierror.Error
	for _, upstreamPlugin := range t.upstreamPlugins {
		pluginSecrets := secretsForPlugin(cfg, upstreamPlugin, secrets)
		if err := upstreamPlugin.ProcessUpstream(upstream, pluginSecrets, out); err != nil {
			upstreamErrors = multierror.Append(upstreamErrors, err)
		}
	}
	if err := validateCluster(out); err != nil {
		upstreamErrors = multierror.Append(upstreamErrors, err)
	}
	return out, upstreamErrors
}

// TODO: add more validation here
func validateCluster(c *envoyapi.Cluster) error {
	if c.Type == envoyapi.Cluster_STATIC || c.Type == envoyapi.Cluster_STRICT_DNS || c.Type == envoyapi.Cluster_LOGICAL_DNS {
		if len(c.Hosts) < 1 {
			return errors.Errorf("cluster type %v specified but hosts were empty", c.Type.String())
		}
	}
	return nil
}

func createUpstreamReport(upstream v1.Upstream, err error) reporter.ConfigObjectReport {
	return reporter.ConfigObjectReport{
		CfgObject: &upstream,
		Err:       err,
	}
}

func createVirtualHostReport(virtualHost v1.VirtualHost, err error) reporter.ConfigObjectReport {
	return reporter.ConfigObjectReport{
		CfgObject: &virtualHost,
		Err:       err,
	}
}

func secretsForPlugin(cfg v1.Config, plug plugin.TranslatorPlugin, secrets secretwatcher.SecretMap) secretwatcher.SecretMap {
	deps := plug.GetDependencies(cfg)
	if deps == nil || len(deps.SecretRefs) == 0 {
		return nil
	}
	pluginSecrets := make(secretwatcher.SecretMap)
	for _, ref := range deps.SecretRefs {
		pluginSecrets[ref] = secrets[ref]
	}
	return pluginSecrets
}

// VirtualHosts

func (t *Translator) computeVirtualHosts(cfg v1.Config) ([]envoyroute.VirtualHost, []reporter.ConfigObjectReport) {
	var (
		reports      []reporter.ConfigObjectReport
		virtualHosts []envoyroute.VirtualHost
	)
	for _, virtualHost := range cfg.VirtualHosts {
		envoyVirtualHost, err := t.computeVirtualHost(cfg.Upstreams, virtualHost)
		virtualHosts = append(virtualHosts, envoyVirtualHost)
		reports = append(reports, createVirtualHostReport(virtualHost, err))
	}
	return virtualHosts, reports
}

func (t *Translator) computeVirtualHost(upstreams []v1.Upstream, virtualHost v1.VirtualHost) (envoyroute.VirtualHost, error) {
	var envoyRoutes []envoyroute.Route
	var routeErrors *multierror.Error
	for _, route := range virtualHost.Routes {
		if err := validateRoute(upstreams, route); err != nil {
			routeErrors = multierror.Append(routeErrors, err)
		}
		envoyRoute := envoyroute.Route{}
		for _, routePlugin := range t.routePlugins {
			if err := routePlugin.ProcessRoute(route, &envoyRoute); err != nil {
				routeErrors = multierror.Append(routeErrors, err)
			}
		}
		envoyRoutes = append(envoyRoutes, envoyRoute)
	}
	domains := virtualHost.Domains
	if len(domains) == 0 || (len(domains) == 1 && domains[0] == "") {
		domains = []string{"*"}
	}

	// TODO: handle default virtualhost
	// TODO: handle ssl
	return envoyroute.VirtualHost{
		Name:    envoy.VirtualHostName(virtualHost.Name),
		Domains: domains,
		Routes:  envoyRoutes,
	}, routeErrors
}

func validateRoute(upstreams []v1.Upstream, route v1.Route) error {
	// collect existing upstreams/functions for matching
	upstreamsAndTheirFunctions := make(map[string][]string)
	for _, upstream := range upstreams {
		var funcsForUpstream []string
		for _, fn := range upstream.Functions {
			funcsForUpstream = append(funcsForUpstream, fn.Name)
		}
		upstreamsAndTheirFunctions[upstream.Name] = funcsForUpstream
	}

	// make sure the destination itself has the right structure
	if len(route.Destination.Destinations) > 0 {
		return validateMultiDestination(upstreamsAndTheirFunctions, route.Destination.Destinations)
	}
	return validateSingleDestination(upstreamsAndTheirFunctions, route.Destination.SingleDestination)
}

func validateMultiDestination(upstreamsAndTheirFunctions map[string][]string, destinations []v1.WeightedDestination) error {
	for _, dest := range destinations {
		if err := validateSingleDestination(upstreamsAndTheirFunctions, dest.SingleDestination); err != nil {
			return errors.Wrap(err, "invalid destination in weighted destination list")
		}
	}
	return nil
}

func validateSingleDestination(upstreamsAndTheirFunctions map[string][]string, destination v1.SingleDestination) error {
	if destination.FunctionDestination != nil && destination.UpstreamDestination != nil {
		return errors.New("only one of function_destination and upstream_destination can be set on a single destination")
	}
	if destination.UpstreamDestination != nil {
		return validateUpstreamDestination(upstreamsAndTheirFunctions, destination.UpstreamDestination)
	}
	if destination.FunctionDestination != nil {
		return validateFunctionDestination(upstreamsAndTheirFunctions, destination.FunctionDestination)
	}
	return errors.New("must specify either a function or upstream on a single destination")
}

func validateUpstreamDestination(upstreamsAndTheirFunctions map[string][]string, upstreamDestination *v1.UpstreamDestination) error {
	upstreamName := upstreamDestination.UpstreamName
	if _, ok := upstreamsAndTheirFunctions[upstreamName]; !ok {
		return errors.Errorf("upstream %v was not found for function destination", upstreamName)
	}
	return nil
}

func validateFunctionDestination(upstreamsAndTheirFunctions map[string][]string, functionDestination *v1.FunctionDestination) error {
	upstreamName := functionDestination.UpstreamName
	upstreamFuncs, ok := upstreamsAndTheirFunctions[upstreamName]
	if !ok {
		return errors.Errorf("upstream %v was not found for function destination", upstreamName)
	}
	functionName := functionDestination.FunctionName
	if !stringInSlice(upstreamFuncs, functionName) {
		return errors.Errorf("function %v/%v was not found for function destination", upstreamName, functionName)
	}
	return nil
}
