package translator

import (
	"fmt"
	"sort"

	envoyapi "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	envoycore "github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	envoyendpoints "github.com/envoyproxy/go-control-plane/envoy/api/v2/endpoint"
	envoylistener "github.com/envoyproxy/go-control-plane/envoy/api/v2/listener"
	envoyroute "github.com/envoyproxy/go-control-plane/envoy/api/v2/route"
	envoyhttp "github.com/envoyproxy/go-control-plane/envoy/config/filter/network/http_connection_manager/v2"
	envoycache "github.com/envoyproxy/go-control-plane/pkg/cache"
	envoyutil "github.com/envoyproxy/go-control-plane/pkg/util"

	"github.com/ghodss/yaml"
	"github.com/gogo/protobuf/proto"
	"github.com/hashicorp/go-multierror"
	"github.com/mitchellh/hashstructure"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/util/runtime"

	"github.com/solo-io/gloo/pkg/log"
	"github.com/solo-io/gloo/pkg/protoutil"
	"github.com/solo-io/gloo/internal/reporter"
	"github.com/solo-io/gloo/pkg/api/types/v1"
	"github.com/solo-io/gloo/pkg/endpointdiscovery"
	"github.com/solo-io/gloo/pkg/plugin"
	"github.com/solo-io/gloo/pkg/secretwatcher"
)

const (
	rdsName       = "gloo-rds"
	listenerName  = "listener-" + rdsName
	listenerPort  = uint32(8080)
	connMgrFilter = "envoy.http_connection_manager"
	routerFilter  = "envoy.router"
)

type Translator struct {
	plugins []plugin.TranslatorPlugin
}

func NewTranslator(plugins []plugin.TranslatorPlugin) *Translator {
	// special routing must be done for upstream plugins that support functions
	var functionPlugins []plugin.FunctionPlugin
	for _, plug := range plugins {
		if functionPlugin, ok := plug.(plugin.FunctionPlugin); ok {
			functionPlugins = append(functionPlugins, functionPlugin)
		}
	}
	// the initializer plugin must be initialized with any function plugins
	// it's responsible for setting cluster weights and common route properties
	initPlugin := newInitializerPlugin(functionPlugins)
	plugins = append([]plugin.TranslatorPlugin{initPlugin}, plugins...)
	return &Translator{
		plugins: plugins,
	}
}

func (t *Translator) Translate(cfg *v1.Config,
	secrets secretwatcher.SecretMap,
	endpoints endpointdiscovery.EndpointGroups) (*envoycache.Snapshot, []reporter.ConfigObjectReport, error) {

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

	// listeners
	// TODO: eventaully support multiple listeners (e.g. for TLS)
	listener, err := t.constructHttpListener(listenerName, listenerPort, routeConfig.Name)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "constructing http listener %v", listenerName)
	}

	// proto-ify everything
	var endpointsProto []proto.Message
	for _, cla := range clusterLoadAssignments {
		endpointsProto = append(endpointsProto, cla)
	}

	var clustersProto []proto.Message
	for _, cluster := range clusters {
		clustersProto = append(clustersProto, cluster)
	}

	routesProto := []proto.Message{routeConfig}
	listenersProto := []proto.Message{listener}

	// construct version
	// TODO: investigate whether we need a more sophisticated versionining algorithm
	version, err := hashstructure.Hash([][]proto.Message{
		endpointsProto,
		clustersProto,
		routesProto,
		listenersProto,
	}, nil)
	if err != nil {
		return nil, nil, errors.Wrap(err, "constructing version hash for envoy snapshot components")
	}

	log.Printf("DID ERR? %v", printYaml(map[envoycache.ResponseType][]proto.Message{
		envoycache.EndpointResponse: endpointsProto,
		envoycache.ClusterResponse:  clustersProto,
		envoycache.RouteResponse:    routesProto,
		envoycache.ListenerResponse: listenersProto,
	}))

	// construct snapshot
	snapshot := envoycache.NewSnapshot(fmt.Sprintf("%v", version), endpointsProto, clustersProto, routesProto, listenersProto)

	// aggregate reports
	reports := append(clusterReports, virtualHostReports...)

	return &snapshot, reports, nil
}

func printYaml(snaps map[envoycache.ResponseType][]proto.Message) error {
	for resourceType, snap := range snaps {
		log.GreyPrintf("\n\n%s\n", resourceType)
		for _, pro := range snap {
			jsn, err := protoutil.Marshal(pro)
			if err != nil {
				return err
			}
			yam, err := yaml.JSONToYAML(jsn)
			if err != nil {
				return err
			}
			log.GreyPrintf("%s", yam)
		}
	}
	return nil
}

// Endpoints

func computeClusterEndpoints(upstreams []*v1.Upstream, endpoints endpointdiscovery.EndpointGroups) []*envoyapi.ClusterLoadAssignment {
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

func (t *Translator) computeClusters(cfg *v1.Config, secrets secretwatcher.SecretMap, endpoints endpointdiscovery.EndpointGroups) ([]*envoyapi.Cluster, []reporter.ConfigObjectReport) {
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

func (t *Translator) computeCluster(cfg *v1.Config, secrets secretwatcher.SecretMap, upstream *v1.Upstream, edsCluster bool) (*envoyapi.Cluster, error) {
	out := &envoyapi.Cluster{
		Name:     upstream.Name,
		Metadata: new(envoycore.Metadata),
	}
	if edsCluster {
		out.Type = envoyapi.Cluster_EDS
	}
	var upstreamErrors *multierror.Error
	for _, plug := range t.plugins {
		upstreamPlugin, ok := plug.(plugin.UpstreamPlugin)
		if !ok {
			continue
		}
		pluginSecrets := secretsForPlugin(cfg, upstreamPlugin, secrets)
		params := &plugin.UpstreamPluginParams{
			Secrets:              pluginSecrets,
			EnvoyNameForUpstream: clusterName,
		}
		if err := upstreamPlugin.ProcessUpstream(params, upstream, out); err != nil {
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

func createUpstreamReport(upstream *v1.Upstream, err error) reporter.ConfigObjectReport {
	return reporter.ConfigObjectReport{
		CfgObject: upstream,
		Err:       err,
	}
}

func createVirtualHostReport(virtualHost *v1.VirtualHost, err error) reporter.ConfigObjectReport {
	return reporter.ConfigObjectReport{
		CfgObject: virtualHost,
		Err:       err,
	}
}

func secretsForPlugin(cfg *v1.Config, plug plugin.TranslatorPlugin, secrets secretwatcher.SecretMap) secretwatcher.SecretMap {
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

func (t *Translator) computeVirtualHosts(cfg *v1.Config) ([]envoyroute.VirtualHost, []reporter.ConfigObjectReport) {
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

func (t *Translator) computeVirtualHost(upstreams []*v1.Upstream, virtualHost *v1.VirtualHost) (envoyroute.VirtualHost, error) {
	var envoyRoutes []envoyroute.Route
	var routeErrors *multierror.Error
	for _, route := range virtualHost.Routes {
		if err := validateRoute(upstreams, route); err != nil {
			routeErrors = multierror.Append(routeErrors, err)
		}
		out := newBaseEnvoyRoute(route)
		for _, plug := range t.plugins {
			routePlugin, ok := plug.(plugin.RoutePlugin)
			if !ok {
				continue
			}
			params := &plugin.RoutePluginParams{}
			if err := routePlugin.ProcessRoute(params, route, &out); err != nil {
				routeErrors = multierror.Append(routeErrors, err)
			}
		}
		envoyRoutes = append(envoyRoutes, out)
	}
	domains := virtualHost.Domains
	if len(domains) == 0 || (len(domains) == 1 && domains[0] == "") {
		domains = []string{"*"}
	}

	// TODO: handle default virtualhost
	// TODO: handle ssl
	return envoyroute.VirtualHost{
		Name:    virtualHostName(virtualHost.Name),
		Domains: domains,
		Routes:  envoyRoutes,
	}, routeErrors
}

func newBaseEnvoyRoute(route *v1.Route) envoyroute.Route {
	match := envoyroute.RouteMatch{}
	switch path := route.Matcher.Path.(type) {
	case *v1.Matcher_PathRegex:
		match.PathSpecifier = &envoyroute.RouteMatch_Regex{
			Regex: path.PathRegex,
		}
	case *v1.Matcher_PathPrefix:
		match.PathSpecifier = &envoyroute.RouteMatch_Prefix{
			Prefix: path.PathPrefix,
		}
	case *v1.Matcher_PathExact:
		match.PathSpecifier = &envoyroute.RouteMatch_Path{
			Path: path.PathExact,
		}
	}
	for headerName, headerValue := range route.Matcher.Headers {
		match.Headers = append(match.Headers, &envoyroute.HeaderMatcher{
			Name:  headerName,
			Value: headerValue,
		})
	}
	for paramName, paramValue := range route.Matcher.QueryParams {
		match.QueryParameters = append(match.QueryParameters, &envoyroute.QueryParameterMatcher{
			Name:  paramName,
			Value: paramValue,
		})
	}
	return envoyroute.Route{
		Metadata: new(envoycore.Metadata),
		Match:    match,
	}
}

func validateRoute(upstreams []*v1.Upstream, route *v1.Route) error {
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
	switch {
	case route.SingleDestination != nil && len(route.MultipleDestinations) == 0:
		return validateSingleDestination(upstreamsAndTheirFunctions, route.SingleDestination)
	case route.SingleDestination == nil && len(route.MultipleDestinations) > 0:
		return validateMultiDestination(upstreamsAndTheirFunctions, route.MultipleDestinations)
	}
	return errors.Errorf("must specify either 'single_destination' or 'multiple_destinations' for route")
}

func validateMultiDestination(upstreamsAndTheirFunctions map[string][]string, destinations []*v1.WeightedDestination) error {
	for _, dest := range destinations {
		if err := validateSingleDestination(upstreamsAndTheirFunctions, dest.Destination); err != nil {
			return errors.Wrap(err, "invalid destination in weighted destination list")
		}
	}
	return nil
}

func validateSingleDestination(upstreamsAndTheirFunctions map[string][]string, destination *v1.Destination) error {
	switch dest := destination.DestinationType.(type) {
	case *v1.Destination_Upstream:
		return validateUpstreamDestination(upstreamsAndTheirFunctions, dest)
	case *v1.Destination_Function:
		return validateFunctionDestination(upstreamsAndTheirFunctions, dest)
	}
	return errors.New("must specify either a function or upstream on a single destination")
}

func validateUpstreamDestination(upstreamsAndTheirFunctions map[string][]string, upstreamDestination *v1.Destination_Upstream) error {
	upstreamName := upstreamDestination.Upstream.Name
	if _, ok := upstreamsAndTheirFunctions[upstreamName]; !ok {
		return errors.Errorf("upstream %v was not found for function destination", upstreamName)
	}
	return nil
}

func validateFunctionDestination(upstreamsAndTheirFunctions map[string][]string, functionDestination *v1.Destination_Function) error {
	upstreamName := functionDestination.Function.UpstreamName
	upstreamFuncs, ok := upstreamsAndTheirFunctions[upstreamName]
	if !ok {
		return errors.Errorf("upstream %v was not found for function destination", upstreamName)
	}
	functionName := functionDestination.Function.FunctionName
	if !stringInSlice(upstreamFuncs, functionName) {
		return errors.Errorf("function %v/%v was not found for function destination", upstreamName, functionName)
	}
	return nil
}

func stringInSlice(slice []string, s string) bool {
	for _, el := range slice {
		if el == s {
			return true
		}
	}
	return false
}

// Listener

type stagedFilter struct {
	filter *envoyhttp.HttpFilter
	stage  plugin.Stage
}

func (t *Translator) constructHttpListener(name string, port uint32, routeConfigName string) (*envoyapi.Listener, error) {
	var filtersByStage []stagedFilter
	for _, plug := range t.plugins {
		filterPlugin, ok := plug.(plugin.FilterPlugin)
		if !ok {
			continue
		}
		params := &plugin.FilterPluginParams{}
		httpFilter, stage := filterPlugin.HttpFilter(params)
		if httpFilter == nil {
			runtime.HandleError(errors.New("plugin implements HttpFilter() but returned nil"))
			continue
		}
		filtersByStage = append(filtersByStage, stagedFilter{
			filter: httpFilter,
			stage:  stage,
		})
	}

	// sort filters by stage
	httpFilters := sortFilters(filtersByStage)
	httpFilters = append(httpFilters, &envoyhttp.HttpFilter{Name: routerFilter})

	httpConnMgr := &envoyhttp.HttpConnectionManager{
		CodecType:  envoyhttp.AUTO,
		StatPrefix: "http",
		RouteSpecifier: &envoyhttp.HttpConnectionManager_Rds{
			Rds: &envoyhttp.Rds{
				ConfigSource: envoycore.ConfigSource{
					ConfigSourceSpecifier: &envoycore.ConfigSource_Ads{
						Ads: &envoycore.AggregatedConfigSource{},
					},
				},
				RouteConfigName: routeConfigName,
			},
		},
		HttpFilters: httpFilters,
	}

	httpConnMgrCfg, err := envoyutil.MessageToStruct(httpConnMgr)
	if err != nil {
		return nil, errors.Wrap(err, "failed to convert proto message to struct")
	}
	return &envoyapi.Listener{
		Name: name,
		Address: envoycore.Address{
			Address: &envoycore.Address_SocketAddress{
				SocketAddress: &envoycore.SocketAddress{
					Protocol: envoycore.TCP,
					Address:  "0.0.0.0", // bind all
					PortSpecifier: &envoycore.SocketAddress_PortValue{
						PortValue: port,
					},
				},
			},
		},
		FilterChains: []envoylistener.FilterChain{{
			Filters: []envoylistener.Filter{{
				Name:   connMgrFilter,
				Config: httpConnMgrCfg,
			}},
		}},
	}, nil
}

func sortFilters(filters []stagedFilter) []*envoyhttp.HttpFilter {
	// sort them first by stage, then by name.
	less := func(i, j int) bool {
		filteri := filters[i]
		filterj := filters[j]
		if filteri.stage != filterj.stage {
			return filteri.stage < filterj.stage
		}
		return filteri.filter.Name < filterj.filter.Name
	}
	sort.SliceStable(filters, less)

	var sortedFilters []*envoyhttp.HttpFilter
	for _, filter := range filters {
		sortedFilters = append(sortedFilters, filter.filter)
	}

	return sortedFilters
}

// for future-proofing possible safety issues with bad upstream names
func clusterName(upstreamName string) string {
	return upstreamName
}

// for future-proofing possible safety issues with bad virtualhost names
func virtualHostName(virtualHostName string) string {
	return virtualHostName
}
