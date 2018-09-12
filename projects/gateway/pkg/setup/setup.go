package setup

import (
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/reporter"
	"github.com/solo-io/solo-kit/pkg/errors"
	"github.com/solo-io/solo-kit/pkg/utils/contextutils"
	"github.com/solo-io/solo-kit/pkg/utils/errutils"
	"github.com/solo-io/solo-kit/projects/gateway/pkg/api/v1"
	"github.com/solo-io/solo-kit/projects/gateway/pkg/defaults"
	"github.com/solo-io/solo-kit/projects/gateway/pkg/propagator"
	"github.com/solo-io/solo-kit/projects/gateway/pkg/syncer"
	gloov1 "github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1"
	"github.com/solo-io/solo-kit/samples"
)

func Setup(opts Opts) error {
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

func setupForNamespaces(watchNamespaces []string, opts Opts) error {
	opts.WatchOpts = opts.WatchOpts.WithDefaults()

	gatewayClient, err := v1.NewGatewayClient(opts.Gateways)
	if err != nil {
		return err
	}
	if err := gatewayClient.Register(); err != nil {
		return err
	}

	virtualServicesClient, err := v1.NewVirtualServiceClient(opts.VirtualServices)
	if err != nil {
		return err
	}
	if err := virtualServicesClient.Register(); err != nil {
		return err
	}

	proxyClient, err := gloov1.NewProxyClient(opts.Proxies)
	if err != nil {
		return err
	}

	if _, err := gatewayClient.Write(defaults.DefaultGateway(opts.WriteNamespace), clients.WriteOpts{
		Ctx: opts.WatchOpts.Ctx,
	}); err != nil && !errors.IsExist(err) {
		return err
	}

	if opts.SampleData {
		if err := addSampleData(opts, virtualServicesClient); err != nil {
			return err
		}

	}

	emitter := v1.NewApiEmitter(gatewayClient, virtualServicesClient)

	rpt := reporter.NewReporter("gateway", gatewayClient.BaseClient(), virtualServicesClient.BaseClient())
	writeErrs := make(chan error)

	prop := propagator.NewPropagator("gateway", gatewayClient, virtualServicesClient, proxyClient, writeErrs)

	sync := syncer.NewSyncer(opts.WriteNamespace, proxyClient, rpt, prop, writeErrs)

	eventLoop := v1.NewApiEventLoop(emitter, sync)
	eventLoopErrs, err := eventLoop.Run(watchNamespaces, opts.WatchOpts)
	if err != nil {
		return err
	}
	go errutils.AggregateErrs(opts.WatchOpts.Ctx, writeErrs, eventLoopErrs, "event_loop")

	logger := contextutils.LoggerFrom(opts.WatchOpts.Ctx)

	for {
		select {
		case err := <-writeErrs:
			logger.Errorf("error: %v", err)
		case <-opts.WatchOpts.Ctx.Done():
			close(writeErrs)
			return nil
		}
	}
}

func addSampleData(opts Opts, vsClient v1.VirtualServiceClient) error {
	upstreamClient, err := gloov1.NewUpstreamClient(opts.Upstreams)
	if err != nil {
		return err
	}
	secretClient, err := gloov1.NewSecretClient(opts.Secrets)
	if err != nil {
		return err
	}
	virtualServices, upstreams, secrets := samples.VirtualServices(), samples.Upstreams(), samples.Secrets()
	for _, item := range virtualServices {
		if _, err := vsClient.Write(item, clients.WriteOpts{}); err != nil && !errors.IsExist(err) {
			return err
		}
	}
	for _, item := range upstreams {
		if _, err := upstreamClient.Write(item, clients.WriteOpts{}); err != nil && !errors.IsExist(err) {
			return err
		}
	}
	for _, item := range secrets {
		if _, err := secretClient.Write(item, clients.WriteOpts{}); err != nil && !errors.IsExist(err) {
			return err
		}
	}
	return nil
}
