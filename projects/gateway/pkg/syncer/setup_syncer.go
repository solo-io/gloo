package syncer

import (
	"context"

	"github.com/gogo/protobuf/types"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/kube"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/memory"
	"github.com/solo-io/solo-kit/pkg/api/v1/reporter"
	"github.com/solo-io/solo-kit/pkg/errors"
	"github.com/solo-io/solo-kit/pkg/utils/contextutils"
	"github.com/solo-io/solo-kit/pkg/utils/errutils"
	"github.com/solo-io/solo-kit/projects/gateway/pkg/api/v1"
	"github.com/solo-io/solo-kit/projects/gateway/pkg/defaults"
	"github.com/solo-io/solo-kit/projects/gateway/pkg/propagator"
	gloov1 "github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/bootstrap"
	gloodefaults "github.com/solo-io/solo-kit/projects/gloo/pkg/defaults"
	"k8s.io/client-go/rest"
)

func NewSetupSyncer(inMemoryCache memory.InMemoryResourceCache, kubeCache *kube.KubeCache) gloov1.SetupSyncer {
	return &setupSyncer{
		inMemoryCache: inMemoryCache,
		kubeCache:     kubeCache,
	}
}

type setupSyncer struct {
	kubeCache     *kube.KubeCache
	inMemoryCache memory.InMemoryResourceCache
}

func (s *setupSyncer) Sync(ctx context.Context, snap *gloov1.SetupSnapshot) error {
	switch {
	case len(snap.Settings.List()) == 0:
		return errors.Errorf("no settings files found")
	case len(snap.Settings.List()) > 1:
		return errors.Errorf("multiple settings files found")
	}
	settings := snap.Settings.List()[0]

	var (
		cfg *rest.Config
	)
	cache := s.inMemoryCache
	kubeCache := s.kubeCache

	proxyFactory, err := bootstrap.ConfigFactoryForSettings(
		settings,
		cache,
		kubeCache,
		gloov1.ProxyCrd,
		&cfg,
	)
	if err != nil {
		return err
	}

	virtualServiceFactory, err := bootstrap.ConfigFactoryForSettings(
		settings,
		cache,
		kubeCache,
		v1.VirtualServiceCrd,
		&cfg,
	)
	if err != nil {
		return err
	}

	gatewayFactory, err := bootstrap.ConfigFactoryForSettings(
		settings,
		cache,
		kubeCache,
		v1.GatewayCrd,
		&cfg,
	)
	if err != nil {
		return err
	}

	refreshRate, err := types.DurationFromProto(settings.RefreshRate)
	if err != nil {
		return err
	}

	writeNamespace := settings.DiscoveryNamespace
	if writeNamespace == "" {
		writeNamespace = gloodefaults.GlooSystem
	}
	watchNamespaces := settings.WatchNamespaces
	var writeNamespaceProvided bool
	for _, ns := range watchNamespaces {
		if ns == writeNamespace {
			writeNamespaceProvided = true
			break
		}
	}
	if !writeNamespaceProvided {
		watchNamespaces = append(watchNamespaces, writeNamespace)
	}
	opts := Opts{
		WriteNamespace:  writeNamespace,
		WatchNamespaces: watchNamespaces,
		Gateways:        gatewayFactory,
		VirtualServices: virtualServiceFactory,
		Proxies:         proxyFactory,
		WatchOpts: clients.WatchOpts{
			Ctx:         ctx,
			RefreshRate: refreshRate,
		},
		DevMode: true,
	}

	return RunGateway(opts)
}

func RunGateway(opts Opts) error {
	opts.WatchOpts = opts.WatchOpts.WithDefaults()
	opts.WatchOpts.Ctx = contextutils.WithLogger(opts.WatchOpts.Ctx, "gateway")

	gatewayClient, err := v1.NewGatewayClient(opts.Gateways)
	if err != nil {
		return err
	}
	if err := gatewayClient.Register(); err != nil {
		return err
	}

	virtualServiceClient, err := v1.NewVirtualServiceClient(opts.VirtualServices)
	if err != nil {
		return err
	}
	if err := virtualServiceClient.Register(); err != nil {
		return err
	}

	proxyClient, err := gloov1.NewProxyClient(opts.Proxies)
	if err != nil {
		return err
	}
	if err := proxyClient.Register(); err != nil {
		return err
	}

	if _, err := gatewayClient.Write(defaults.DefaultGateway(opts.WriteNamespace), clients.WriteOpts{
		Ctx: opts.WatchOpts.Ctx,
	}); err != nil && !errors.IsExist(err) {
		return err
	}

	emitter := v1.NewApiEmitter(gatewayClient, virtualServiceClient)

	rpt := reporter.NewReporter("gateway", gatewayClient.BaseClient(), virtualServiceClient.BaseClient())
	writeErrs := make(chan error)

	prop := propagator.NewPropagator("gateway", gatewayClient, virtualServiceClient, proxyClient, writeErrs)

	sync := NewTranslatorSyncer(opts.WriteNamespace, proxyClient, gatewayClient, virtualServiceClient, rpt, prop, writeErrs)

	eventLoop := v1.NewApiEventLoop(emitter, sync)
	eventLoopErrs, err := eventLoop.Run(opts.WatchNamespaces, opts.WatchOpts)
	if err != nil {
		return err
	}
	go errutils.AggregateErrs(opts.WatchOpts.Ctx, writeErrs, eventLoopErrs, "event_loop")

	logger := contextutils.LoggerFrom(opts.WatchOpts.Ctx)

	go func() {
		for {
			select {
			case err := <-writeErrs:
				logger.Errorf("error: %v", err)
			case <-opts.WatchOpts.Ctx.Done():
				close(writeErrs)
				return
			}
		}
	}()
	return nil
}
