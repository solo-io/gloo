package setup

import (
	"context"

	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/factory"
	"github.com/solo-io/solo-kit/pkg/namespacing"
	"github.com/solo-io/solo-kit/pkg/namespacing/static"
	"github.com/solo-io/solo-kit/pkg/utils/contextutils"
	"github.com/solo-io/solo-kit/pkg/utils/kubeutils"
	gatewayv1 "github.com/solo-io/solo-kit/projects/gateway/pkg/api/v1"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/defaults"
	"k8s.io/client-go/kubernetes"
)

type Opts struct {
	WriteNamespace  string
	Gateways        factory.ResourceClientFactory
	VirtualServices factory.ResourceClientFactory
	Proxies         factory.ResourceClientFactory
	Namespacer      namespacing.Namespacer
	WatchOpts       clients.WatchOpts
	DevMode         bool
}

func NewOpts(
	writeNamespace string,
	gateways,
	virtualServices,
	proxies factory.ResourceClientFactory,
	namespacer namespacing.Namespacer,
	watchOpts clients.WatchOpts,
	devMode bool,
) Opts {
	return Opts{
		WriteNamespace:  writeNamespace,
		Gateways:        gateways,
		VirtualServices: virtualServices,
		Proxies:         proxies,
		Namespacer:      namespacer,
		WatchOpts:       watchOpts,
		DevMode:         devMode,
	}
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
	ctx := contextutils.WithLogger(context.Background(), "gateway")
	return Opts{
		WriteNamespace: defaults.GlooSystem,
		Gateways: &factory.KubeResourceClientFactory{
			Crd: gatewayv1.GatewayCrd,
			Cfg: cfg,
		},
		VirtualServices: &factory.KubeResourceClientFactory{
			Crd: gatewayv1.VirtualServiceCrd,
			Cfg: cfg,
		},
		Proxies: &factory.KubeResourceClientFactory{
			Crd: v1.ProxyCrd,
			Cfg: cfg,
		},
		Namespacer: static.NewNamespacer([]string{"default", defaults.GlooSystem}),
		WatchOpts: clients.WatchOpts{
			Ctx:         ctx,
			RefreshRate: defaults.RefreshRate,
		},
		DevMode: true,
	}, nil
}
