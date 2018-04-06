package service

import (
	"net"

	envoyapi "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	envoycore "github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	"github.com/pkg/errors"

	"github.com/solo-io/gloo/pkg/api/types/v1"
	"github.com/solo-io/gloo/pkg/plugins"
)

type Plugin struct{}

const (
	// define Upstream type name
	UpstreamTypeService = "service"
)

func (p *Plugin) GetDependencies(_ *v1.Config) *plugins.Dependencies {
	return nil
}

func (p *Plugin) ProcessUpstream(_ *plugins.UpstreamPluginParams, in *v1.Upstream, out *envoyapi.Cluster) error {
	if in.Type != UpstreamTypeService {
		return nil
	}
	// decode does validation for us
	spec, err := DecodeUpstreamSpec(in.Spec)
	if err != nil {
		return errors.Wrap(err, "invalid service upstream spec")
	}
	for _, host := range spec.Hosts {
		if host.Addr == "" {
			return errors.New("addr cannot be empty for host")
		}
		if host.Port == 0 {
			return errors.New("port cannot be empty for host")
		}
		ip := net.ParseIP(host.Addr)
		if ip != nil {
			out.Type = envoyapi.Cluster_STATIC
		} else {
			out.Type = envoyapi.Cluster_STRICT_DNS
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
