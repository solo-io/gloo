package service

import (
	"net"

	envoyapi "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	envoycore "github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	"github.com/pkg/errors"

	"github.com/solo-io/glue/internal/plugins"
	"github.com/solo-io/glue/pkg/api/types/v1"
	"github.com/solo-io/glue/pkg/plugin"
	"github.com/solo-io/glue/pkg/secretwatcher"
)

func init() {
	plugins.Register(&Plugin{}, nil)
}

type Plugin struct{}

const (
	// define Upstream type name
	UpstreamTypeService = "service"
)

func (p *Plugin) GetDependencies(_ *v1.Config) *plugin.Dependencies {
	return nil
}

func (p *Plugin) ProcessUpstream(in *v1.Upstream, _ secretwatcher.SecretMap, out *envoyapi.Cluster) error {
	if in.Type != UpstreamTypeService {
		return nil
	}
	// decode does validation for us
	spec, err := DecodeUpstreamSpec(in.Spec)
	if err != nil {
		return errors.Wrap(err, "invalid service upstream spec")
	}
	for _, host := range spec.Hosts {
		ip := net.ParseIP(host.Addr)
		if ip != nil {
			out.Type = envoyapi.Cluster_STATIC
		} else {
			out.Type = envoyapi.Cluster_LOGICAL_DNS
		}
		out.Hosts = append(out.Hosts, &envoycore.Address{
			Address: &envoycore.Address_SocketAddress{
				SocketAddress: &envoycore.SocketAddress{
					Protocol: envoycore.TCP,
					Address:  host.Addr,
					PortSpecifier: &envoycore.SocketAddress_PortValue{
						PortValue: host.Port,
					},
				},
			},
		})
	}
	return nil
}

//func (p *Plugin) ProcessRoute(in v1.Route, out *envoyroute.Route) error {
//	upstreamNames := destinationUpstreams(in)
//	switch {
//	case len(upstreamNames) == 0:
//		return errors.New("no upstreams found for route")
//	case len(upstreamNames) == 1:
//		out.Action = &envoyroute.Route_Route{
//			Route: &envoyroute.RouteAction{
//				ClusterSpecifier: &envoyroute.RouteAction_Cluster{
//					Cluster: envoy.ClusterName(upstreamNames[0]),
//				},
//			},
//		}
//	case len(upstreamNames) > 1:
//		out.Action = &envoyroute.Route_Route{
//			Route: &envoyroute.RouteAction{
//				ClusterSpecifier: &envoyroute.RouteAction_WeightedClusters{
//					Cluster: envoy.ClusterName(upstreamNames[0]),
//				},
//			},
//		}
//	}
//	return nil
//}
//
//func destinationUpstreams(in v1.Route) []string {
//	var destinationUpstreams []string
//	dests := []v1.SingleDestination{in.Destination.SingleDestination}
//	for _, dest := range in.Destination.Destinations {
//		dests = append(dests, dest.SingleDestination)
//	}
//	for _, dest := range dests {
//		if upstreamName := getUpstreamName(dest); upstreamName != "" {
//			destinationUpstreams = append(destinationUpstreams, upstreamName)
//		}
//	}
//	return destinationUpstreams
//}
//
//func getUpstreamName(dest v1.SingleDestination) string {
//	switch {
//	case dest.UpstreamDestination != nil:
//		return dest.UpstreamDestination.UpstreamName
//	case dest.FunctionDestination != nil:
//		return dest.FunctionDestination.UpstreamName
//	}
//	return ""
//}
