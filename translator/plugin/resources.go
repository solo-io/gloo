package plugin

import 	"github.com/envoyproxy/go-control-plane/api/filter/network"

type Stage int

const (
	PreAuth Stage = iota
	Auth
	PostAuth
)

type EnvoyResources struct {
	Filters  []FilterWrapper
	Routes   []RouteWrapper
	Clusters []ClusterWrapper
	//TODO: VirtualHosts []VirtualHostWrapper
}

type FilterWrapper struct {
	Filter network.HttpFilter
	Stage  Stage
}

type RouteWrapper struct {
	Route  api.Route
	Weight int
	// optional; if not populated, will use 'default'
	VirtualHost string
}

type ClusterWrapper struct {
	Cluster api.Cluster
}