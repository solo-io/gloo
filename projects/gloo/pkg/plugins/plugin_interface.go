package plugins

import (
	envoyapi "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1"
)

type Plugin interface {
	Init() error
}

type PluginParams struct {
	Secrets     []*v1.Secret
	Artifacts   []*v1.Artifact
}

type UpstreamPlugin interface {
	Plugin
	ProcessUpstream(params PluginParams, in *v1.Upstream, out *envoyapi.Cluster) error
}

type EdsPlugin interface {
	Plugin
	RunEds(client v1.EndpointClient, upstreams []*v1.Upstream) error
}
