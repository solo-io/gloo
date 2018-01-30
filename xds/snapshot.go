package xds

import (
	"fmt"
	"sort"

	"github.com/envoyproxy/go-control-plane/api"
	"github.com/envoyproxy/go-control-plane/api/filter/network"
	"github.com/envoyproxy/go-control-plane/pkg/cache"
	"github.com/envoyproxy/go-control-plane/pkg/util"
	"github.com/gogo/protobuf/proto"
	"github.com/mitchellh/hashstructure"
	"github.com/solo-io/glue/config"
	"github.com/solo-io/glue/pkg/log"
)

func createSnapshot(resources []config.EnvoyResources) (cache.Snapshot, error) {
	var (
		filters  byStage
		routes   byWeight
		clusters byName
	)

	for _, resources := range resources {
		for _, filter := range resources.Filters {
			filters = append(filters, filter)
		}
		for _, route := range resources.Routes {
			routes = append(routes, route)
		}
		for _, cluster := range resources.Clusters {
			clusters = append(clusters, cluster)
		}
	}
	api.Cluster{}
	sort.Sort(filters)
	sort.Sort(routes)
	sort.Sort(clusters)
	version, err := hashstructure.Hash(config.EnvoyResources{
		Filters:  filters,
		Routes:   routes,
		Clusters: clusters,
	}, nil)
	if err != nil {
		return cache.Snapshot{}, err
	}

	routeConfig := makeRouteConfig(routes)
	listener := makeListener(filters)

	log.Printf("adding %v routes", len(routes))
	log.Printf("adding %v filters", len(filters))
	log.Printf("adding %v clusters", len(clusters))

	var protoClusters []proto.Message

	for _, cluster := range clusters {
		envoyCluster := cluster.Cluster
		protoClusters = append(protoClusters, &envoyCluster)
	}

	return cache.NewSnapshot(fmt.Sprintf("%d", version),
		nil,
		protoClusters,
		[]proto.Message{routeConfig},
		[]proto.Message{listener}), nil
}

/*
########################
##########Listeners#####
########################
*/
// MakeListener creates a listener.
const (
	//router     = "envoy.router"
	httpFilter         = "envoy.http_connection_manager"
	routeConfigName    = "my_route_config"
	xdsCluster         = "xds_cluster"
	port               = 8080
	listener           = "any-name-is-fine"
	defaultVirtualHost = "default"
)

func makeEndpointConfig(cluster, adress string, port int) *api.ClusterLoadAssignment {
	return &api.ClusterLoadAssignment{
		ClusterName: cluster,
		Endpoints: []*api.LocalityLbEndpoints{{
			LbEndpoints: []*api.LbEndpoint{{
				Endpoint: &api.Endpoint{
					Address: &api.Address{
						Address: &api.Address_SocketAddress{
							SocketAddress: &api.SocketAddress{
								Protocol: api.SocketAddress_TCP,
								Address:  localhost,
								PortSpecifier: &api.SocketAddress_PortValue{
									PortValue: port,
								},
							},
						},
					},
				},
			}},
		}},
	}
}

func makeRouteConfig(routes []config.RouteWrapper) *api.RouteConfiguration {
	routesByVirtualHost := make(map[string][]config.RouteWrapper)
	for _, route := range routes {
		vHostName := route.VirtualHost
		if vHostName == "" {
			vHostName = defaultVirtualHost
		}
		routesByVirtualHost[vHostName] = append(routesByVirtualHost[vHostName], route)
	}
	var vHosts []*api.VirtualHost
	for vHostName, routes := range routesByVirtualHost {
		var envoyRoutes []*api.Route
		for _, route := range routes {
			envoyRoute := route.Route
			envoyRoutes = append(envoyRoutes, &envoyRoute)
		}
		vHost := &api.VirtualHost{
			Name:    vHostName,
			Domains: []string{"*"},
			Routes:  envoyRoutes,
		}
		vHosts = append(vHosts, vHost)
	}
	return &api.RouteConfiguration{
		Name:         routeConfigName,
		VirtualHosts: vHosts,
	}
}

func makeListener(filters []config.FilterWrapper) *api.Listener {
	rdsSource := api.ConfigSource{
		ConfigSourceSpecifier: &api.ConfigSource_ApiConfigSource{
			ApiConfigSource: &api.ApiConfigSource{
				ApiType:      api.ApiConfigSource_GRPC,
				ClusterNames: []string{xdsCluster},
				GrpcServices: []*api.GrpcService{
					{
						TargetSpecifier: &api.GrpcService_EnvoyGrpc_{
							EnvoyGrpc: &api.GrpcService_EnvoyGrpc{
								ClusterName: xdsCluster,
							},
						},
					},
				},
			},
		},
	}
	var httpFilters []*network.HttpFilter
	for _, filter := range filters {
		envoyFilter := filter.Filter
		httpFilters = append(httpFilters, &envoyFilter)
	}
	manager := &network.HttpConnectionManager{
		CodecType:  network.HttpConnectionManager_AUTO,
		StatPrefix: "http",
		RouteSpecifier: &network.HttpConnectionManager_Rds{
			Rds: &network.Rds{
				ConfigSource:    rdsSource,
				RouteConfigName: routeConfigName,
			},
		},
		HttpFilters: httpFilters,
	}
	pbst, err := util.MessageToStruct(manager)
	if err != nil {
		panic(err)
	}

	return &api.Listener{
		Name: listener,
		Address: &api.Address{
			Address: &api.Address_SocketAddress{
				SocketAddress: &api.SocketAddress{
					Protocol: api.SocketAddress_TCP,
					Address:  "0.0.0.0",
					PortSpecifier: &api.SocketAddress_PortValue{
						PortValue: port,
					},
				},
			},
		},
		FilterChains: []*api.FilterChain{{
			Filters: []*api.Filter{{
				Name:   httpFilter,
				Config: pbst,
			}},
		}},
	}
}

/*
########################
############Sorting#####
########################
*/
type byStage []config.FilterWrapper

func (s byStage) Len() int {
	return len(s)
}
func (s byStage) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s byStage) Less(i, j int) bool {
	if s[i].Stage == s[j].Stage {
		return s[i].Filter.Name < s[j].Filter.Name
	}
	return s[i].Stage < s[j].Stage
}

type byWeight []config.RouteWrapper

func (s byWeight) Len() int {
	return len(s)
}
func (s byWeight) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s byWeight) Less(i, j int) bool {
	return s[i].Weight < s[j].Weight
}

type byName []config.ClusterWrapper

func (s byName) Len() int {
	return len(s)
}
func (s byName) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s byName) Less(i, j int) bool {
	return s[i].Cluster.Name < s[j].Cluster.Name
}
