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
)

type Opts struct {
	writeNamespace  string
	gateways        factory.ResourceClientFactoryOpts
	virtualServices factory.ResourceClientFactoryOpts
	// TODO(ilackarms): remove upstreams here if not needed, right now only used for sample data
	upstreams  factory.ResourceClientFactoryOpts
	secrets    factory.ResourceClientFactoryOpts
	proxies    factory.ResourceClientFactoryOpts
	namespacer namespacing.Namespacer
	watchOpts  clients.WatchOpts
	devMode    bool
}

func NewOpts(
	writeNamespace string,
	gateways,
	virtualServices,
	upstreams,
	secrets,
	proxies factory.ResourceClientFactoryOpts,
	namespacer namespacing.Namespacer,
	watchOpts clients.WatchOpts,
	devMode bool,
) Opts {
	return Opts{
		writeNamespace:  writeNamespace,
		gateways:        gateways,
		virtualServices: virtualServices,
		upstreams:       upstreams,
		secrets:         secrets,
		proxies:         proxies,
		namespacer:      namespacer,
		watchOpts:       watchOpts,
	}
}

//  ilackarms: We can just put any hacky stuff we need here

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
		writeNamespace: "gloo-system",
		gateways: &factory.KubeResourceClientOpts{
			Crd: gatewayv1.GatewayCrd,
			Cfg: cfg,
		},
		virtualServices: &factory.KubeResourceClientOpts{
			Crd: gatewayv1.VirtualServiceCrd,
			Cfg: cfg,
		},
		proxies: &factory.KubeResourceClientOpts{
			Crd: v1.ProxyCrd,
			Cfg: cfg,
		},
		upstreams: &factory.KubeResourceClientOpts{
			Crd: v1.UpstreamCrd,
			Cfg: cfg,
		},
		secrets: &factory.KubeSecretClientOpts{
			Clientset: clientset,
		},
		namespacer: static.NewNamespacer([]string{"default", "gloo-system"}),
		watchOpts: clients.WatchOpts{
			Ctx:         ctx,
			RefreshRate: time.Minute,
		},
	}, nil
}
