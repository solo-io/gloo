package defaults

import (
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-kit/projects/gateway/pkg/api/v1"
	gloov1 "github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1"
)

func DefaultGateway(writeNamespace string) *v1.Gateway {
	return &v1.Gateway{
		Metadata: core.Metadata{
			Name:      "gateway",
			Namespace: writeNamespace,
		},
		BindAddress:     "::",
		BindPort:        80,
		VirtualServices: []string{}, // all virtualservices
	}
}

func DefaultVirtualService(writeNamespace string) *v1.VirtualService {
	return &v1.VirtualService{
		Metadata: core.Metadata{
			Name:      "routeman",
			Namespace: writeNamespace,
		},
		VirtualHost: &gloov1.VirtualHost{
			Name:    "routes",
			Domains: []string{"*"},
			Routes:  []*gloov1.Route{},
		},
	}
}

func LocalUpstream(writeNamespace string) *gloov1.Upstream {
	return &gloov1.Upstream{
		Metadata: core.Metadata{
			Name:      "upstream",
			Namespace: writeNamespace,
		},
		UpstreamSpec: &gloov1.UpstreamSpec{},
	}
}
