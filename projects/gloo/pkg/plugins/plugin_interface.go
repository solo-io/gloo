package plugins

import (
	"github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1"
)

type Plugin interface {
	Init() error
}

type EdsPlugin interface {
	Plugin
	RunEds(client v1.EndpointClient, upstreams []*v1.Upstream) error
}
