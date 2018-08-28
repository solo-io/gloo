package setup

import (
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/factory"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/memory"
	"github.com/solo-io/solo-kit/pkg/api/v1/reporter"
	"github.com/solo-io/solo-kit/pkg/utils/contextutils"
	"github.com/solo-io/solo-kit/pkg/utils/errutils"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/discovery"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/syncer"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/translator"
)

type Opts struct {
	namespace string
	inputResourceOpts factory.ResourceClientFactoryOpts
	secretOpts factory.ResourceClientFactoryOpts
	artifactOpts factory.ResourceClientFactoryOpts
	opts clients.WatchOpts
}

func Setup(namespace string, inputResourceOpts factory.ResourceClientFactoryOpts, secretOpts factory.ResourceClientFactoryOpts, artifactOpts factory.ResourceClientFactoryOpts, opts clients.WatchOpts) error {
	opts = opts.WithDefaults()
	opts.Ctx = contextutils.WithLogger(opts.Ctx, "setup")
	inputFactory := factory.NewResourceClientFactory(inputResourceOpts)
	secretFactory := factory.NewResourceClientFactory(secretOpts)
	artifactFactory := factory.NewResourceClientFactory(artifactOpts)
	// endpoints are internal-only, therefore use the in-memory client
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

	disc := discovery.NewDiscovery(namespace, upstreamClient, endpointClient)

	xdsCache, xdsServer, err := newXds()
	if err != nil {
		return err
	}

	go xdsServer.Run(ctx)

	rpt := reporter.NewReporter("gloo", upstreamClient.BaseClient(), proxyClient.BaseClient())

	sync := syncer.NewSyncer(translator.NewTranslator(), xdsCache, rpt)

	eventLoop := v1.NewEventLoop(cache, sync)

	errs := make(chan error)

	udsErrs, err := discovery.RunUds(disc, opts, discovery.Opts{
	// TODO(ilackarms)
	})
	if err != nil {
		return err
	}
	go errutils.AggregateErrs(opts.Ctx, errs, udsErrs, "uds.gloo")

	edsErrs, err := discovery.RunEds(upstreamClient, disc, namespace, opts)
	if err != nil {
		return err
	}
	go errutils.AggregateErrs(opts.Ctx, errs, edsErrs, "eds.gloo")

	eventLoopErrs, err := eventLoop.Run(namespace, opts)
	if err != nil {
		return err
	}
	go errutils.AggregateErrs(opts.Ctx, errs, eventLoopErrs, "event_loop.gloo")

	logger := contextutils.LoggerFrom(opts.Ctx)

	for {
		select {
		case err := <-errs:
			logger.Errorf("error: %v", err)
		case <-opts.Ctx.Done():
			close(errs)
			return nil
		}
	}
}
