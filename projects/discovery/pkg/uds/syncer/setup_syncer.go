package setup

import (
	"github.com/solo-io/solo-kit/pkg/utils/contextutils"
	"github.com/solo-io/solo-kit/pkg/utils/errutils"
	"github.com/solo-io/solo-kit/projects/discovery/pkg/uds/syncer"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/bootstrap"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/discovery"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/plugins/registry"
	"k8s.io/client-go/rest"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/memory"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"context"
	"github.com/solo-io/solo-kit/pkg/errors"
	"github.com/gogo/protobuf/types"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/defaults"
)

func NewSetupSyncer() v1.SetupSyncer {
	return &settingsSyncer{}
}

type settingsSyncer struct{}

func (s *settingsSyncer) Sync(ctx context.Context, snap *v1.SetupSnapshot) error {
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
	cache := memory.NewInMemoryResourceCache()

	proxyFactory, err := bootstrap.ConfigFactoryForSettings(
		settings,
		cache,
		v1.ProxyCrd,
		&cfg,
	)
	if err != nil {
		return err
	}

	virtualServiceFactory, err := bootstrap.ConfigFactoryForSettings(
		settings,
		cache,
		v1.VirtualServiceCrd,
		&cfg,
	)
	if err != nil {
		return err
	}

	gatewayFactory, err := bootstrap.ConfigFactoryForSettings(
		settings,
		cache,
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
		writeNamespace = defaults.GlooSystem
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
	opts := bootstrap.Opts{
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

func RunUDS(opts bootstrap.Opts) error {
	watchOpts := opts.WatchOpts.WithDefaults()
	watchOpts.Ctx = contextutils.WithLogger(watchOpts.Ctx, "uds")

	upstreamClient, err := v1.NewUpstreamClient(opts.Upstreams)
	if err != nil {
		return err
	}
	if err := upstreamClient.Register(); err != nil {
		return err
	}

	secretClient, err := v1.NewSecretClient(opts.Secrets)
	if err != nil {
		return err
	}
	if err := secretClient.Register(); err != nil {
		return err
	}

	emit := make(chan struct{})
	emitter := v1.NewDiscoveryEmitterWithEmit(secretClient, upstreamClient, emit)

	// jump start all the watches
	go func() {
		emit <- struct{}{}
	}()

	plugins := registry.Plugins(opts)

	var discoveryPlugins []discovery.DiscoveryPlugin
	for _, plug := range plugins {
		disc, ok := plug.(discovery.DiscoveryPlugin)
		if ok {
			discoveryPlugins = append(discoveryPlugins, disc)
		}
	}
	disc := discovery.NewUpstreamDiscovery(opts.WatchNamespaces, opts.WriteNamespace, upstreamClient, discoveryPlugins)

	sync := syncer.NewSyncer(disc,
		discovery.Opts{}, // TODO(ilackarms)
		watchOpts.RefreshRate)
	eventLoop := v1.NewDiscoveryEventLoop(emitter, sync)

	errs := make(chan error)

	eventLoopErrs, err := eventLoop.Run(opts.WatchNamespaces, watchOpts)
	if err != nil {
		return err
	}
	go errutils.AggregateErrs(watchOpts.Ctx, errs, eventLoopErrs, "event_loop.gloo")

	logger := contextutils.LoggerFrom(watchOpts.Ctx)

	for {
		select {
		case err := <-errs:
			logger.Errorf("error: %v", err)
		case <-watchOpts.Ctx.Done():
			return nil
		}
	}
}
