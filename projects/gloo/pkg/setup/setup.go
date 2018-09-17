package setup

import (
	"net"

	"context"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/golang/protobuf/ptypes"
	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/factory"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/memory"
	"github.com/solo-io/solo-kit/pkg/api/v1/reporter"
	"github.com/solo-io/solo-kit/pkg/errors"
	"github.com/solo-io/solo-kit/pkg/namespacing/static"
	"github.com/solo-io/solo-kit/pkg/utils/contextutils"
	"github.com/solo-io/solo-kit/pkg/utils/errutils"
	"github.com/solo-io/solo-kit/pkg/utils/kubeutils"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/bootstrap"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/defaults"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/discovery"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/plugins/registry"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/syncer"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/translator"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/xds"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func SetupRoot(configdir string) error {
	el := 
}

func NewSettingsSync() v1.SetupSyncer {
	return &settingsSyncer{
		grpcServer: func(ctx context.Context) *grpc.Server {
			return grpc.NewServer(grpc.StreamInterceptor(
				grpc_middleware.ChainStreamServer(
					grpc_ctxtags.StreamServerInterceptor(),
					grpc_zap.StreamServerInterceptor(zap.NewNop()),
					func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
						contextutils.LoggerFrom(ctx).Infof("gRPC call: %v", info.FullMethod)
						return handler(srv, ss)
					},
				)),
			)
		}
	}
}

type settingsSyncer struct{
	grpcServer func(ctx context.Context) *grpc.Server
}

func (s *settingsSyncer) Sync(ctx context.Context, snap *v1.SetupSnapshot) error {
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
		artifactFactory factory.ResourceClientFactory
	)
	var cfg *rest.Config
	var clientset kubernetes.Interface
	cache := memory.NewInMemoryResourceCache()

	if settings.ConfigSource == nil {
		upstreamFactory = &factory.MemoryResourceClientFactory{
			Cache: cache,
		}
		proxyFactory = &factory.MemoryResourceClientFactory{
			Cache: cache,
		}
	} else {
		switch source := settings.ConfigSource.(type) {
		case *v1.Settings_KubernetesConfigSource:
			var err error
			if cfg == nil {
				cfg, err = kubeutils.GetConfig("", "")
				if err != nil {
					return err
				}
			}
			upstreamFactory = &factory.KubeResourceClientFactory{
				Crd: v1.UpstreamCrd,
				Cfg: cfg,
			}
			proxyFactory = &factory.KubeResourceClientFactory{
				Crd: v1.ProxyCrd,
				Cfg: cfg,
			}
		case *v1.Settings_DirectoryConfigSource:
			upstreamFactory = &factory.FileResourceClientFactory{
				RootDir: filepath.Join(source.DirectoryConfigSource.Directory + "upstreams"),
			}
			proxyFactory = &factory.FileResourceClientFactory{
				RootDir: filepath.Join(source.DirectoryConfigSource.Directory + "proxies"),
			}
		default:
			return errors.Errorf("invalid config source type")
		}
	}

	if settings.SecretSource == nil {
		secretFactory = &factory.MemoryResourceClientFactory{
			Cache: cache,
		}
	} else {
		switch source := settings.SecretSource.(type) {
		case *v1.Settings_KubernetesSecretSource:
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
		case *v1.Settings_VaultSecretSource:
			return errors.Errorf("vault configuration not implemented")
		case *v1.Settings_DirectorySecretSource:
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
		case *v1.Settings_KubernetesArtifactSource:
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
		case *v1.Settings_DirectoryArtifactSource:
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
	refreshRate, err := ptypes.Duration(settings.RefreshRate)
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
		WriteNamespace: writeNamespace,
		Namespacer:     static.NewNamespacer(watchNamespaces),
		Upstreams:      upstreamFactory,
		Proxies:        proxyFactory,
		Secrets:        secretFactory,
		Artifacts:      artifactFactory,
		WatchOpts: clients.WatchOpts{
			Ctx:         ctx,
			RefreshRate: refreshRate,
		},
		BindAddr: &net.TCPAddr{
			IP:   net.ParseIP(ipPort[0]),
			Port: port,
		},
		GrpcServer: s.grpcServer(ctx),
		// if nil, kube plugin disabled
		KubeClient: clientset,
		DevMode:    true,
	}

	return Setup(opts)
}

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

func setupForNamespaces(watchNamespaces []string, opts bootstrap.Opts) error {
	watchOpts := opts.WatchOpts.WithDefaults()
	opts.WatchOpts.Ctx = contextutils.WithLogger(opts.WatchOpts.Ctx, "gloo")

	watchOpts.Ctx = contextutils.WithLogger(watchOpts.Ctx, "setup")
	endpointsFactory := &factory.MemoryResourceClientFactory{
		Cache: memory.NewInMemoryResourceCache(),
	}

	upstreamClient, err := v1.NewUpstreamClient(opts.Upstreams)
	if err != nil {
		return err
	}
	if err := upstreamClient.Register(); err != nil {
		return err
	}

	proxyClient, err := v1.NewProxyClient(opts.Proxies)
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
	if err := endpointClient.Register(); err != nil {
		return err
	}

	secretClient, err := v1.NewSecretClient(opts.Secrets)
	if err != nil {
		return err
	}
	if err := secretClient.Register(); err != nil {
		return err
	}

	artifactClient, err := v1.NewArtifactClient(opts.Artifacts)
	if err != nil {
		return err
	}
	if err := artifactClient.Register(); err != nil {
		return err
	}

	cache := v1.NewApiEmitter(artifactClient, endpointClient, proxyClient, secretClient, upstreamClient)

	xdsHasher, xdsCache := xds.SetupEnvoyXds(opts.WatchOpts.Ctx, opts.GrpcServer, nil)

	rpt := reporter.NewReporter("gloo", upstreamClient.BaseClient(), proxyClient.BaseClient())

	plugins := registry.Plugins(opts)

	var discoveryPlugins []discovery.DiscoveryPlugin
	for _, plug := range plugins {
		disc, ok := plug.(discovery.DiscoveryPlugin)
		if ok {
			discoveryPlugins = append(discoveryPlugins, disc)
		}
	}

	sync := syncer.NewSyncer(translator.NewTranslator(plugins), xdsCache, xdsHasher, rpt, opts.DevMode)
	eventLoop := v1.NewApiEventLoop(cache, sync)

	errs := make(chan error)

	eds := discovery.NewEndpointDiscovery(opts.WriteNamespace, endpointClient, discoveryPlugins)
	edsErrs, err := discovery.RunEds(upstreamClient, eds, opts.WriteNamespace, watchOpts)
	if err != nil {
		return err
	}
	go errutils.AggregateErrs(watchOpts.Ctx, errs, edsErrs, "eds.gloo")

	eventLoopErrs, err := eventLoop.Run(watchNamespaces, watchOpts)
	if err != nil {
		return err
	}
	go errutils.AggregateErrs(watchOpts.Ctx, errs, eventLoopErrs, "event_loop.gloo")

	logger := contextutils.LoggerFrom(watchOpts.Ctx)

	go func() {

		for {
			select {
			case err, ok := <-errs:
				if !ok {
					return
				}
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
