package setup

import (
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/factory"
	"github.com/solo-io/solo-kit/pkg/namespacing"
	"google.golang.org/grpc"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1"
	"github.com/solo-io/solo-kit/pkg/utils/kubeutils"
	"k8s.io/client-go/kubernetes"
)

type Opts struct {
	writeNamespace  string
	upstreamBackend factory.ResourceClientFactoryOpts
	proxyBackend    factory.ResourceClientFactoryOpts
	secretBackend   factory.ResourceClientFactoryOpts
	artifactBackend factory.ResourceClientFactoryOpts
	namespacer      namespacing.Namespacer
	grpcServer      *grpc.Server
	watchOpts       clients.WatchOpts
}

func DefaultKubernetesConstructOpts() (Opts, error) {
	cfg, err := kubeutils.GetConfig("", "")
	if err != nil {
		return Opts{}, err
	}
	clientset, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return Opts{}, err
	}
	return Opts{
		writeNamespace: "gloo-system",
		upstreamBackend: &factory.KubeResourceClientOpts{
			Crd: v1.UpstreamCrd,
			Cfg: cfg,
		},
		proxyBackend: &factory.KubeResourceClientOpts{
			Crd: v1.ProxyCrd,
			Cfg: cfg,
		},
		secretBackend: &factory.KubeSecretClientOpts{
			Clientset: clientset,
		},
	}, nil
}
