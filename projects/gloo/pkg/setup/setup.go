package setup

import (
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/factory"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/memory"
	"github.com/solo-io/solo-kit/pkg/api/v1/reporter"
	"github.com/solo-io/solo-kit/pkg/utils/contextutils"
	"github.com/solo-io/solo-kit/pkg/utils/errutils"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/discovery"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/syncer"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/translator"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/xds"
)

func Setup(opts Opts) error {
	namespaces, errs, err := opts.namespacer.Namespaces(opts.watchOpts)
	if err != nil {
		return err
	}
	for {
		select {
		case err := <-errs:
			return err
		case watchNamespaces := <-namespaces:
			err := setupForNamespaces(watchNamespaces, opts)
			if err != nil {
				return err
			}
		}
	}
}

func setupForNamespaces(discoveredNamespaces []string, opts Opts) error {
	watchOpts := opts.watchOpts.WithDefaults()
	secretOpts := opts.secretBackend
	artifactOpts := opts.artifactBackend

	watchOpts.Ctx = contextutils.WithLogger(watchOpts.Ctx, "setup")
	upstreamFactory := factory.NewResourceClientFactory(secretOpts)
	proxyFactory := factory.NewResourceClientFactory(secretOpts)
	secretFactory := factory.NewResourceClientFactory(secretOpts)
	artifactFactory := factory.NewResourceClientFactory(artifactOpts)
	endpointsFactory := factory.NewResourceClientFactory(&factory.MemoryResourceClientOpts{
		Cache: memory.NewInMemoryResourceCache(),
	})

	upstreamClient, err := v1.NewUpstreamClient(inputFactory)
	if err != nil {
		return err
	}

	proxyClient, err := v1.NewProxyClient(inputFactory)
	if err != nil {
		return err
	}

	endpointClient, err := v1.NewEndpointClient(endpointsFactory)
	if err != nil {
		return err
	}

	secretClient, err := v1.NewSecretClient(secretFactory)
	if err != nil {
		return err
	}

	artifactClient, err := v1.NewArtifactClient(artifactFactory)
	if err != nil {
		return err
	}

	cache := v1.NewCache(artifactClient, endpointClient, proxyClient, secretClient, upstreamClient)

	xdsHasher, xdsCache := xds.SetupEnvoyXds(opts.watchOpts.Ctx, opts.grpcServer, nil)

	rpt := reporter.NewReporter("gloo", upstreamClient.BaseClient(), proxyClient.BaseClient())

	disc := discovery.NewDiscovery(opts.writeNamespace, upstreamClient, endpointClient)

	sync := syncer.NewSyncer(translator.NewTranslator(), xdsCache, xdsHasher, rpt)
	eventLoop := v1.NewEventLoop(cache, sync)

	errs := make(chan error)

	udsErrs, err := discovery.RunUds(disc, watchOpts, discovery.Opts{
		// TODO(ilackarms)
	})
	if err != nil {
		return err
	}
	go errutils.AggregateErrs(watchOpts.Ctx, errs, udsErrs, "uds.gloo")

	edsErrs, err := discovery.RunEds(upstreamClient, disc, opts.writeNamespace, watchOpts)
	if err != nil {
		return err
	}
	go errutils.AggregateErrs(watchOpts.Ctx, errs, edsErrs, "eds.gloo")

	eventLoopErrs, err := eventLoop.Run(discoveredNamespaces, watchOpts)
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
			close(errs)
			return nil
		}
	}
}
