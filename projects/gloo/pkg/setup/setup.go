package setup

import (
	"net"

	"github.com/solo-io/solo-kit/pkg/api/v1/clients/factory"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/memory"
	"github.com/solo-io/solo-kit/pkg/api/v1/reporter"
	"github.com/solo-io/solo-kit/pkg/utils/contextutils"
	"github.com/solo-io/solo-kit/pkg/utils/errutils"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/bootstrap"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/discovery"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/syncer"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/translator"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/xds"
)

func Setup(opts bootstrap.Opts) error {
	// TODO: Ilackarms: move this to multi-eventloop
	namespaces, errs, err := opts.Namespacer.Namespaces(opts.WatchOpts)
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

func setupForNamespaces(discoveredNamespaces []string, opts bootstrap.Opts) error {
	watchOpts := opts.WatchOpts.WithDefaults()

	watchOpts.Ctx = contextutils.WithLogger(watchOpts.Ctx, "setup")
	upstreamFactory := factory.NewResourceClientFactory(opts.Upstreams)
	proxyFactory := factory.NewResourceClientFactory(opts.Proxies)
	secretFactory := factory.NewResourceClientFactory(opts.Secrets)
	artifactFactory := factory.NewResourceClientFactory(opts.Artifacts)
	endpointsFactory := factory.NewResourceClientFactory(&factory.MemoryResourceClientOpts{
		Cache: memory.NewInMemoryResourceCache(),
	})

	upstreamClient, err := v1.NewUpstreamClient(upstreamFactory)
	if err != nil {
		return err
	}
	if err := upstreamClient.Register(); err != nil {
		return err
	}

	proxyClient, err := v1.NewProxyClient(proxyFactory)
	if err != nil {
		return err
	}
	if err := proxyClient.Register(); err != nil {
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

	cache := v1.NewApiEmitter(artifactClient, endpointClient, proxyClient, secretClient, upstreamClient)

	xdsHasher, xdsCache := xds.SetupEnvoyXds(opts.WatchOpts.Ctx, opts.GrpcServer, nil)

	rpt := reporter.NewReporter("gloo", upstreamClient.BaseClient(), proxyClient.BaseClient())

	disc := discovery.NewDiscovery(opts.WriteNamespace, upstreamClient, endpointClient)

	sync := syncer.NewSyncer(translator.NewTranslator(opts), xdsCache, xdsHasher, rpt, opts.DevMode)
	eventLoop := v1.NewApiEventLoop(cache, sync)

	errs := make(chan error)

	udsErrs, err := discovery.RunUds(disc, watchOpts, discovery.Opts{
	// TODO(ilackarms)
	})
	if err != nil {
		return err
	}
	go errutils.AggregateErrs(watchOpts.Ctx, errs, udsErrs, "uds.gloo")

	edsErrs, err := discovery.RunEds(upstreamClient, disc, opts.WriteNamespace, watchOpts)
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

	go func() {

		for {
			select {
			case err := <-errs:
				logger.Errorf("error: %v", err)
			case <-watchOpts.Ctx.Done():
				return
			}
		}
	}()

	lis, err := net.Listen(opts.BindAddr.Network(), opts.BindAddr.String())
	if err != nil {
		return err
	}
	return opts.GrpcServer.Serve(lis)
}
