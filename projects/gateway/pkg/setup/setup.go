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
	"github.com/solo-io/solo-kit/samples"
)

func Setup(opts Opts) error {
	// TODO: Ilackarms: move this to multi-eventloop
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

func setupForNamespaces(watchNamespaces []string, opts Opts) error {
	opts.watchOpts = opts.watchOpts.WithDefaults()
	gatewayFactory := factory.NewResourceClientFactory(opts.gateways)
	virtualServiceFactory := factory.NewResourceClientFactory(opts.virtualServices)
	proxyFactory := factory.NewResourceClientFactory(opts.proxies)

	gatewayClient, err := v1.NewGatewayClient(gatewayFactory)
	if err != nil {
		return err
	}
	if err := gatewayClient.Register(); err != nil {
		return err
	}

	virtualServicesClient, err := v1.NewVirtualServiceClient(virtualServiceFactory)
	if err != nil {
		return err
	}
	if err := virtualServicesClient.Register(); err != nil {
		return err
	}

	proxyClient, err := gloov1.NewProxyClient(proxyFactory)
	if err != nil {
		return err
	}

	// TODO(ilackarms): make this part of -dev mode
	if err := addSampleData(opts, virtualServicesClient); err != nil {
		return err
	}

	cache := v1.NewCache(gatewayClient, virtualServicesClient)

	rpt := reporter.NewReporter("gateway", gatewayClient.BaseClient(), virtualServicesClient.BaseClient())
	writeErrs := make(chan error)

	prop := propagator.NewPropagator("gateway", gatewayClient, virtualServicesClient, proxyClient, writeErrs)

	sync := syncer.NewSyncer(opts.writeNamespace, proxyClient, rpt, prop, writeErrs)

	eventLoop := v1.NewEventLoop(cache, sync)
	eventLoop.Run(watchNamespaces, opts.watchOpts)

	logger := contextutils.LoggerFrom(opts.watchOpts.Ctx)

	for {
		select {
		case err := <-writeErrs:
			logger.Errorf("error: %v", err)
		case <-opts.watchOpts.Ctx.Done():
			close(writeErrs)
			return nil
		}
	}
}

func addSampleData(opts Opts, vsClient v1.VirtualServiceClient) error {
	upstreamClient, err := gloov1.NewUpstreamClient(factory.NewResourceClientFactory(opts.upstreams))
	if err != nil {
		return err
	}
	virtualServices, upstreams := samples.VirtualServices(), samples.Upstreams()
	for _, item := range virtualServices {
		if _, err := vsClient.Write(item, clients.WriteOpts{}); err != nil {
			return err
		}
	}
	for _, item := range upstreams {
		if _, err := upstreamClient.Write(item, clients.WriteOpts{}); err != nil {
			return err
		}
	}
	return nil
}
