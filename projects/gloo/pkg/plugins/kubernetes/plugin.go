package kubernetes

import (
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

