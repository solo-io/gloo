package kubernetes

import (
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/plugins"
	"k8s.io/client-go/kubernetes"
)

type KubePlugin struct {
	kube kubernetes.Interface
}

func (p *KubePlugin) Init(params plugins.InitParams) error {
	p.kube = params.Bootstrap.KubeClient()
	return nil
}

func (p *KubePlugin) RunEds(opts clients.WatchOpts) error {
	panic("implement me")
}

func (p *KubePlugin) SubscribeUpstream(upstream *v1.Upstream) (<-chan []*v1.Endpoint, error) {
	panic("implement me")
}
