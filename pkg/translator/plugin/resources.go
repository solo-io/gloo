package plugin

import (
	api "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	"github.com/envoyproxy/go-control-plane/envoy/api/v2/route"
	hcm "github.com/envoyproxy/go-control-plane/envoy/config/filter/network/http_connection_manager/v2"
)

type Stage int

const (
	PreInAuth Stage = iota
	InAuth
	PostInAuth
	PreOutAuth
	OutAuth
)

type FilterWrapper struct {
	Filter hcm.HttpFilter
	Stage  Stage
}

type RouteWrapper struct {
	Route  route.Route
	Weight int
	// optional; if not populated, will use 'default'
	VirtualHost string
}

type ClusterWrapper struct {
	Cluster api.Cluster
}

type EnvoyResources struct {
	Filters  []FilterWrapper
	Routes   []RouteWrapper
	Clusters []ClusterWrapper
	//TODO: VirtualHosts []VirtualHostWrapper
}
