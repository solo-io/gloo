package setup

import (
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/factory"
	"github.com/solo-io/solo-kit/pkg/api/v1/reporter"
	"github.com/solo-io/solo-kit/pkg/utils/contextutils"
	"github.com/solo-io/solo-kit/projects/gateway/pkg/api/v1"
	"github.com/solo-io/solo-kit/projects/gateway/pkg/propagator"
	"github.com/solo-io/solo-kit/projects/gateway/pkg/syncer"
	gloov1 "github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1"
)

func Run(namespace string, inputResourceOpts factory.ResourceClientFactoryOpts, opts clients.WatchOpts) error {
	opts = opts.WithDefaults()
	opts.Ctx = contextutils.WithLogger(opts.Ctx, "setup")
	inputFactory := factory.NewResourceClientFactory(inputResourceOpts)

	gatewayClient, err := v1.NewGatewayClient(inputFactory)
	if err != nil {
		return err
	}

	virtualServicesClient, err := v1.NewVirtualServiceClient(inputFactory)
	if err != nil {
		return err
	}

	proxyClient, err := gloov1.NewProxyClient(inputFactory)
	if err != nil {
		return err
	}

	cache := v1.NewCache(gatewayClient, virtualServicesClient)

	rpt := reporter.NewReporter("gateway", gatewayClient.BaseClient(), virtualServicesClient.BaseClient())
	writeErrs := make(chan error)

	prop := propagator.NewPropagator("gateway", gatewayClient, virtualServicesClient, proxyClient, writeErrs)

	sync := syncer.NewSyncer(namespace, proxyClient, rpt, prop, writeErrs)

	eventLoop := v1.NewEventLoop(cache, sync)
	eventLoop.Run(namespace, opts)

	logger := contextutils.LoggerFrom(opts.Ctx)

	for {
		select {
		case err := <-writeErrs:
			logger.Errorf("error: %v", err)
		case <-opts.Ctx.Done():
			close(writeErrs)
			return nil
		}
	}
}
