package setup

import (
	"context"
	"time"

	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/factory"
	"github.com/solo-io/solo-kit/pkg/namespacing"
	"github.com/solo-io/solo-kit/pkg/namespacing/static"
	"github.com/solo-io/solo-kit/pkg/utils/contextutils"
	"github.com/solo-io/solo-kit/pkg/utils/kubeutils"
	gatewayv1 "github.com/solo-io/solo-kit/projects/gateway/pkg/api/v1"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1"
	"k8s.io/client-go/kubernetes"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/defaults"
)

type Opts struct {
	WriteNamespace  string
	Gateways        factory.ResourceClientFactory
	VirtualServices factory.ResourceClientFactory
	// TODO(ilackarms): remove Upstreams here if not needed, right now only used for sample data
	Upstreams  factory.ResourceClientFactory
	Secrets    factory.ResourceClientFactory
	Proxies    factory.ResourceClientFactory
	Namespacer namespacing.Namespacer
	WatchOpts  clients.WatchOpts
	DevMode    bool
	SampleData    bool
}

func NewOpts(
	writeNamespace string,
	gateways,
	virtualServices,
	upstreams,
	secrets,
	proxies factory.ResourceClientFactory,
	namespacer namespacing.Namespacer,
	watchOpts clients.WatchOpts,
	devMode bool,
) Opts {
	return Opts{
		WriteNamespace:  writeNamespace,
		Gateways:        gateways,
		VirtualServices: virtualServices,
		Upstreams:       upstreams,
		Secrets:         secrets,
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
	ctx = contextutils.SilenceLogger(ctx)
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
		Upstreams: &factory.KubeResourceClientFactory{
			Crd: v1.UpstreamCrd,
			Cfg: cfg,
		},
		Secrets: &factory.KubeSecretClientFactory{
			Clientset: clientset,
		},
		Namespacer: static.NewNamespacer([]string{"default", defaults.GlooSystem}),
		WatchOpts: clients.WatchOpts{
			Ctx:         ctx,
			RefreshRate: time.Minute,
		},
		DevMode: false,
	}, nil
}
