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
	"context"
	"path/filepath"
	"strconv"
	"strings"
	"github.com/gogo/protobuf/types"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/factory"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/memory"
	"github.com/solo-io/solo-kit/pkg/namespacing/static"
	"github.com/solo-io/solo-kit/pkg/utils/kubeutils"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	gloodefaults "github.com/solo-io/solo-kit/projects/gloo/pkg/defaults"
)

func NewSetupSyncer() gloov1.SetupSyncer {
	return &settingsSyncer{}
}

type settingsSyncer struct{}

func (s *settingsSyncer) Sync(ctx context.Context, snap *gloov1.SetupSnapshot) error {
	switch {
	case len(snap.Settings.List()) == 0:
		return errors.Errorf("no settings files found")
	case len(snap.Settings.List()) > 1:
		return errors.Errorf("multiple settings files found")
	}
	settings := snap.Settings.List()[0]

	var (
		upstreamFactory factory.ResourceClientFactory
		proxyFactory    factory.ResourceClientFactory
		secretFactory   factory.ResourceClientFactory
	)
	var cfg *rest.Config
	var clientset kubernetes.Interface

	if settings.SecretSource == nil {
		secretFactory = &factory.MemoryResourceClientFactory{
			Cache: cache,
		}
	} else {
		switch source := settings.SecretSource.(type) {
		case *gloov1.Settings_KubernetesSecretSource:
			var err error
			if cfg == nil {
				cfg, err = kubeutils.GetConfig("", "")
				if err != nil {
					return err
				}
			}
			if clientset == nil {
				clientset, err = kubernetes.NewForConfig(cfg)
				if err != nil {
					return err
				}
			}
			secretFactory = &factory.KubeSecretClientFactory{
				Clientset: clientset,
			}
		case *gloov1.Settings_VaultSecretSource:
			return errors.Errorf("vault configuration not implemented")
		case *gloov1.Settings_DirectorySecretSource:
			secretFactory = &factory.FileResourceClientFactory{
				RootDir: filepath.Join(source.DirectorySecretSource.Directory + "secrets"),
			}
		default:
			return errors.Errorf("invalid config source type")
		}
	}

	if settings.ArtifactSource == nil {
		artifactFactory = &factory.MemoryResourceClientFactory{
			Cache: cache,
		}
		switch source := settings.ArtifactSource.(type) {
		case *gloov1.Settings_KubernetesArtifactSource:
			var err error
			if cfg == nil {
				cfg, err = kubeutils.GetConfig("", "")
				if err != nil {
					return err
				}
			}
			if clientset == nil {
				clientset, err = kubernetes.NewForConfig(cfg)
				if err != nil {
					return err
				}
			}
			artifactFactory = &factory.KubeSecretClientFactory{
				Clientset: clientset,
			}
		case *gloov1.Settings_DirectoryArtifactSource:
			artifactFactory = &factory.FileResourceClientFactory{
				RootDir: filepath.Join(source.DirectoryArtifactSource.Directory + "artifacts"),
			}
		default:
			return errors.Errorf("invalid config source type")
		}
	}

	ipPort := strings.Split(settings.BindAddr, ":")
	if len(ipPort) != 2 {
		return errors.Errorf("invalid bind addr: %v", settings.BindAddr)
	}
	port, err := strconv.Atoi(ipPort[1])
	if err != nil {
		return errors.Wrapf(err, "invalid bind addr: %v", settings.BindAddr)
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
		Namespacer:      static.NewNamespacer(watchNamespaces),
		Gateways:        upstreamFactory,
		VirtualServices: upstreamFactory,
		Upstreams:       upstreamFactory,
		Proxies:         proxyFactory,
		Secrets:         secretFactory,
		WatchOpts: clients.WatchOpts{
			Ctx:         ctx,
			RefreshRate: refreshRate,
		},
		DevMode: true,
	}

	return RunGateway(opts)
}

func RunGateway(opts Opts) error {
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
	opts.WatchOpts.Ctx = contextutils.WithLogger(opts.WatchOpts.Ctx, "gateway")

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

	emitter := v1.NewApiEmitter(gatewayClient, virtualServicesClient)

	rpt := reporter.NewReporter("gateway", gatewayClient.BaseClient(), virtualServicesClient.BaseClient())
	writeErrs := make(chan error)

	prop := propagator.NewPropagator("gateway", gatewayClient, virtualServicesClient, proxyClient, writeErrs)

	sync := syncer.NewTranslatorSyncer(opts.WriteNamespace, proxyClient, rpt, prop, writeErrs)

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

func unused(opts Opts, vsClient v1.VirtualServiceClient) error {
	if opts.SampleData {
		if err := addSampleData(opts, vsClient); err != nil {
			return err
		}
	}
	return nil
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
