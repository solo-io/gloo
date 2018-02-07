package translator

import (
	"sort"
	"strings"

	envoycache "github.com/envoyproxy/go-control-plane/pkg/cache"
	"github.com/gogo/protobuf/proto"
	"github.com/gogo/protobuf/types"

	api "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	envoy_api_v2_core "github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	apiep "github.com/envoyproxy/go-control-plane/envoy/api/v2/endpoint"
	envoy_api_v2_listener "github.com/envoyproxy/go-control-plane/envoy/api/v2/listener"
	apiroute "github.com/envoyproxy/go-control-plane/envoy/api/v2/route"
	hcm "github.com/envoyproxy/go-control-plane/envoy/config/filter/network/http_connection_manager/v2"

	"github.com/envoyproxy/go-control-plane/pkg/util"
	"github.com/solo-io/glue/pkg/api/types/v1"
	translatoriface "github.com/solo-io/glue/pkg/translator"
	"github.com/solo-io/glue/pkg/translator/plugin"

	"github.com/solo-io/glue/pkg/endpointdiscovery"
	"github.com/solo-io/glue/pkg/secretwatcher"

	"github.com/hashicorp/go-multierror"
)

type dependencies struct {
	secrets secretwatcher.SecretMap
}

func (d *dependencies) Secrets() secretwatcher.SecretMap {
	return d.secrets
}

type Translator struct {
	plugins        []plugin.Plugin
	nameTranslator translatoriface.NameTranslator
}

func NewTranslator(plugins []plugin.Plugin, nameTranslator translatoriface.NameTranslator) *Translator {

	var functionPlugins []plugin.FunctionalPlugin
	for _, p := range plugins {
		if fp, ok := p.(plugin.FunctionalPlugin); ok {
			functionPlugins = append(functionPlugins, fp)
		}
	}

	plugins = append([]plugin.Plugin{NewInitPlugin(functionPlugins)}, plugins...)

	return &Translator{plugins: plugins, nameTranslator: nameTranslator}
}

func constructMatch(in *v1.Matcher) apiroute.RouteMatch {
	var out apiroute.RouteMatch
	if in.Path.Exact != "" {
		out.PathSpecifier = &apiroute.RouteMatch_Path{Path: in.Path.Exact}
	} else if in.Path.Prefix != "" {
		out.PathSpecifier = &apiroute.RouteMatch_Prefix{Prefix: in.Path.Prefix}
	} else if in.Path.Regex != "" {
		out.PathSpecifier = &apiroute.RouteMatch_Regex{Regex: in.Path.Regex}
	}

	if len(in.Verbs) == 1 {
		out.Headers = append(out.Headers, &apiroute.HeaderMatcher{Name: ":method", Value: in.Verbs[0]})
	} else if len(in.Verbs) >= 1 {
		out.Headers = append(out.Headers, &apiroute.HeaderMatcher{Name: ":method", Value: strings.Join(in.Verbs, "|"), Regex: &types.BoolValue{Value: true}})
	}

	for k, v := range in.Headers {
		out.Headers = append(out.Headers, &apiroute.HeaderMatcher{Name: k, Value: v})
	}

	return out
}

func constructRoute(in *v1.Route) *apiroute.Route {
	var out apiroute.Route
	out.Match = constructMatch(&in.Matcher)
	out.Action = &apiroute.Route_Route{
		Route: &apiroute.RouteAction{
			PrefixRewrite: in.RewritePrefix,
		},
	}

	return &out
}

func (t *Translator) constructUpstream(in *v1.Upstream) *api.Cluster {
	var out api.Cluster

	out.Name = t.nameTranslator.UpstreamToClusterName(in.Name)
	return &out
}

func (t *Translator) constructEds(clustername string, addresses []endpointdiscovery.Endpoint) *api.ClusterLoadAssignment {
	var out api.ClusterLoadAssignment

	var endpoints []apiep.LbEndpoint
	for _, adr := range addresses {
		l := apiep.LbEndpoint{
			Endpoint: &apiep.Endpoint{
				Address: &envoy_api_v2_core.Address{
					Address: &envoy_api_v2_core.Address_SocketAddress{
						SocketAddress: &envoy_api_v2_core.SocketAddress{
							Protocol: envoy_api_v2_core.TCP,
							Address:  adr.Address,
							PortSpecifier: &envoy_api_v2_core.SocketAddress_PortValue{
								PortValue: uint32(adr.Port),
							},
						},
					},
				},
			},
		}
		endpoints = append(endpoints, l)
	}

	out = api.ClusterLoadAssignment{
		ClusterName: clustername,
		Endpoints: []apiep.LocalityLbEndpoints{{
			LbEndpoints: endpoints,
		}},
	}

	return &out
}

func (t *Translator) constructListener(pi *plugin.PluginInputs, listener, route string) *api.Listener {
	router := "envoy.router"
	httpFilter := "envoy.http_connection_manager"
	port := uint32(80)

	rdsSource := envoy_api_v2_core.ConfigSource{}
	rdsSource.ConfigSourceSpecifier = &envoy_api_v2_core.ConfigSource_Ads{
		Ads: &envoy_api_v2_core.AggregatedConfigSource{},
	}

	var whttpfilters []plugin.FilterWrapper
	for _, plgin := range t.plugins {
		whttpfilters = append(whttpfilters, plgin.EnvoyFilters(pi)...)
	}
	httpfilters := sortFilters(whttpfilters)

	httpfilters = append(httpfilters, &hcm.HttpFilter{Name: router})

	manager := &hcm.HttpConnectionManager{
		CodecType:  hcm.AUTO,
		StatPrefix: "http",
		RouteSpecifier: &hcm.HttpConnectionManager_Rds{
			Rds: &hcm.Rds{
				ConfigSource:    rdsSource,
				RouteConfigName: route,
			},
		},
		HttpFilters: httpfilters,
	}
	pbst, err := util.MessageToStruct(manager)
	if err != nil {
		panic("should never happen")
	}

	return &api.Listener{
		Name: listener,
		Address: envoy_api_v2_core.Address{
			Address: &envoy_api_v2_core.Address_SocketAddress{
				SocketAddress: &envoy_api_v2_core.SocketAddress{
					Protocol: envoy_api_v2_core.TCP,
					Address:  "::", //bind all
					PortSpecifier: &envoy_api_v2_core.SocketAddress_PortValue{
						PortValue: port,
					},
				},
			},
		},
		FilterChains: []envoy_api_v2_listener.FilterChain{{
			Filters: []envoy_api_v2_listener.Filter{{
				Name:   httpFilter,
				Config: pbst,
			}},
		}},
	}
}

func (t *Translator) Translate(cfg *v1.Config, secretMap secretwatcher.SecretMap, endpoints endpointdiscovery.EndpointGroups) (*envoycache.Snapshot, error) {
	var statues []translatoriface.ConfigStatus

	dependencies := dependencies{secrets: secretMap}
	state := &plugin.State{
		Dependencies: &dependencies,
		Config:       cfg,
	}

	var endpointsproto []proto.Message

	for _, u := range cfg.Upstreams {
		if group, ok := endpoints[u.Name]; ok {
			cla := t.constructEds(t.nameTranslator.UpstreamToClusterName(u.Name), group)
			endpointsproto = append(endpointsproto, cla)
		}
	}

	pi := &plugin.PluginInputs{
		NameTranslator: t.nameTranslator, // TODO
		State:          state,
	}

	var clustersproto []proto.Message

	upstreams := cfg.Upstreams
	for _, upstream := range upstreams {
		var clustererrors *multierror.Error

		envoycluster := t.constructUpstream(&upstream)
		if _, ok := endpoints[upstream.Name]; ok {
			// if we have EDS!
			envoycluster.Type = api.Cluster_EDS
		}

		for _, p := range t.plugins {
			err := p.UpdateEnvoyCluster(pi, &upstream, envoycluster)
			if err != nil {
				clustererrors = multierror.Append(clustererrors, err)
			}
		}

		// TODO: make sure all clusters have lb type.
		// if not mark cluster as invalid.

		// make sure upstream is health
		if clustererrors == nil {
			clustersproto = append(clustersproto, envoycluster)
			statues = append(statues, translatoriface.NewConfigOk(&upstream))

			// now, process functions
			for _, function := range upstream.Functions {
				var functionerrors *multierror.Error
				for _, p := range t.plugins {
					err := p.UpdateFunctionToEnvoyCluster(pi, &upstream, &function, envoycluster)
					if err != nil {
						functionerrors = multierror.Append(functionerrors, err)
					}
				}

				if functionerrors == nil {
					statues = append(statues, translatoriface.NewConfigOk(&function))
				} else {
					statues = append(statues, translatoriface.NewConfigMultiError(&function, functionerrors))
				}
			}

		} else {
			statues = append(statues, translatoriface.NewConfigMultiError(&upstream, clustererrors))

		}

	}

	rdsname := "routes-80"

	var envoyvhosts []apiroute.VirtualHost
	for _, vhost := range cfg.VirtualHosts {

		var routes []apiroute.Route
		for _, route := range vhost.Routes {
			var routeerrors *multierror.Error
			envoyroute := constructRoute(&route)

			// TODO: make sure all clusters that the route points to exist and valid.
			// if not mark route as invalid.

			for _, p := range t.plugins {
				err := p.UpdateEnvoyRoute(pi, &route, envoyroute)
				if err != nil {
					routeerrors = multierror.Append(routeerrors, err)
				}
			}

			if routeerrors == nil {
				routes = append(routes, *envoyroute)
				statues = append(statues, translatoriface.NewConfigOk(&route))
			} else {
				statues = append(statues, translatoriface.NewConfigMultiError(&route, routeerrors))
			}
		}

		envoyvhost := &apiroute.VirtualHost{
			Name:    t.nameTranslator.ToEnvoyVhostName(&vhost),
			Domains: ifEmpty(vhost.Domains, []string{"*"}),
			Routes:  routes,
		}
		statues = append(statues, translatoriface.NewConfigOk(&vhost))

		// if we have ssl certificates, add them to the ssl filter chain.
		// TODO: Create filter chain for listener
		envoyvhosts = append(envoyvhosts, *envoyvhost)
	}
	routeConfig := &api.RouteConfiguration{
		Name:         rdsname,
		VirtualHosts: envoyvhosts,
	}

	var routessproto []proto.Message
	routessproto = append(routessproto, routeConfig)

	listener := t.constructListener(pi, "listener-"+rdsname, rdsname)
	var listenerproto []proto.Message
	listenerproto = append(listenerproto, listener)

	version := "TODO"

	snapshot := envoycache.NewSnapshot(version,
		endpointsproto,
		clustersproto,
		routessproto,
		listenerproto)

	// create the routes

	/*
		create all clusters, and run the filters on all clusters.
		if from some reason a cluster has errored, send it back to user. and remove it
		from the list
	*/

	/*
		Create virtual hosts and ssl certificates and the such.
		for each virtual host, go over it's routes and:
			Create all routes inline, and then send them to be augmented by all filters
	*/

	// runTranslation

	// combine with cluster + endpoints
	// stable sort

	// computer snapshort version
	return &snapshot, nil
}

func ifEmpty(l []string, def []string) []string {
	if len(l) != 0 {
		return l
	}
	return def
}

func sortFilters(filters []plugin.FilterWrapper) []*hcm.HttpFilter {
	// sort them accoirding to stage and then according to the name.
	less := func(i, j int) bool {
		filteri := filters[i]
		filterj := filters[j]
		if filteri.Stage != filterj.Stage {
			return filteri.Stage < filterj.Stage
		}
		return filteri.Filter.Name < filterj.Filter.Name
	}
	sort.Slice(filters, less)

	var sortedFilters []*hcm.HttpFilter
	for _, filter := range filters {
		sortedFilters = append(sortedFilters, &filter.Filter)
	}

	return sortedFilters
}
