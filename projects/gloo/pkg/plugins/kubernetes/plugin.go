package kubernetes

import (
	"github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/plugins"
	"k8s.io/client-go/kubernetes"
)

type KubePlugin struct {
	kube kubernetes.Interface

	// indicates a resync is required
	resync chan struct{}

	// track upstreams for eds
	trackUpstreams func(list v1.UpstreamList)
}

func (p *KubePlugin) Init(params plugins.InitParams) error {
	p.kube = params.Bootstrap.KubeClient()
	return nil
}

