package kubernetes

import (
	"github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1"
	"k8s.io/client-go/kubernetes"
)

type KubePlugin struct {
	kube kubernetes.Interface
}

func (p *KubePlugin) Init(kube kubernetes.Interface) error {
	p.kube = kube
	return nil
}

func (p *KubePlugin) RunEds(client v1.EndpointClient, upstreams []*v1.Upstream) error {
	panic("implement me")
}
