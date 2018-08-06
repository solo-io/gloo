package plugins

import (
	"github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/bootstrap"
)

type Plugin interface {
	Init(options bootstrap.Options) error
}

type EdsPlugin interface {
	Plugin
	RunEds(client v1.EndpointClient, upstreams []*v1.Upstream) error
}
