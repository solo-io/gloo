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
)

type Opts struct {
	writeNamespace  string
	gateways        factory.ResourceClientFactoryOpts
	virtualServices factory.ResourceClientFactoryOpts
	// TODO(ilackarms): remove upstreams here if not needed, right now only used for sample data
	upstreams  factory.ResourceClientFactoryOpts
	proxies    factory.ResourceClientFactoryOpts
	namespacer namespacing.Namespacer
	watchOpts  clients.WatchOpts
}

//  ilackarms: We can just put any hacky stuff we need here

func DefaultKubernetesConstructOpts() (Opts, error) {
	cfg, err := kubeutils.GetConfig("", "")
	if err != nil {
		return Opts{}, err
	}
	ctx := contextutils.WithLogger(context.Background(), "main")
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
		namespacer: static.NewNamespacer([]string{"default", "gloo-system"}),
		watchOpts: clients.WatchOpts{
			Ctx:         ctx,
			RefreshRate: time.Minute,
		},
	}, nil
}
