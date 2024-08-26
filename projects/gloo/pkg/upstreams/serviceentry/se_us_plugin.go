package serviceentry

import (
	"github.com/solo-io/gloo/pkg/utils/kubeutils"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins"
	"github.com/solo-io/gloo/projects/gloo/pkg/upstreams"

	"istio.io/istio/pkg/kube"
)

var _ upstreams.ClientPlugin = &seUsClientPlugin{}

type seUsClientPlugin struct{}

// Client implements upstreams.ClientPlugin.
func (s *seUsClientPlugin) Client() v1.UpstreamClient {
	c, err := kubeutils.GetRestConfigWithKubeContext("")
	if err != nil {
		panic(err)
	}
	clientCfg := kube.NewClientConfigForRestConfig(cfg)
	client, _ := kube.NewClient(clientCfg, "")
	return NewServiceEntryUpstreamClient(c)
}

// Init implements upstreams.ClientPlugin.
func (s *seUsClientPlugin) Init(params plugins.InitParams) {
}

// Name implements upstreams.ClientPlugin.
func (s *seUsClientPlugin) Name() string {
	panic("unimplemented")
}

// SourceName implements upstreams.ClientPlugin.
func (s *seUsClientPlugin) SourceName() string {
	panic("unimplemented")
}
