package setup

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"slices"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/solo-io/gloo/pkg/utils/settingsutil"
	"github.com/solo-io/gloo/pkg/utils/statsutils/metrics"
	"github.com/solo-io/gloo/projects/gloo/pkg/debug"
	"github.com/solo-io/gloo/projects/gloo/pkg/servers/iosnapshot"
	"github.com/solo-io/gloo/projects/gloo/pkg/utils"
	"github.com/solo-io/solo-kit/pkg/api/v1/control-plane/cache"

	"github.com/golang/protobuf/ptypes/duration"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	consulapi "github.com/hashicorp/consul/api"
	vaultapi "github.com/hashicorp/vault/api"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/proto"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"github.com/solo-io/go-utils/contextutils"
	"github.com/solo-io/go-utils/errutils"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/factory"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/kube"
	corecache "github.com/solo-io/solo-kit/pkg/api/v1/clients/kube/cache"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/memory"
	"github.com/solo-io/solo-kit/pkg/api/v1/control-plane/server"
	xdsserver "github.com/solo-io/solo-kit/pkg/api/v1/control-plane/server"
	"github.com/solo-io/solo-kit/pkg/api/v1/control-plane/types"
	"github.com/solo-io/solo-kit/pkg/api/v2/reporter"
	"github.com/solo-io/solo-kit/pkg/errors"
	"github.com/solo-io/solo-kit/pkg/utils/prototime"
	"github.com/solo-io/solo-kit/pkg/utils/statusutils"

	"github.com/solo-io/gloo/pkg/bootstrap/leaderelector"
	"github.com/solo-io/gloo/pkg/utils/channelutils"
	"github.com/solo-io/gloo/pkg/utils/envutils"
	"github.com/solo-io/gloo/pkg/utils/namespaces"
	"github.com/solo-io/gloo/pkg/utils/setuputils"
	gloostatusutils "github.com/solo-io/gloo/pkg/utils/statusutils"
	gateway "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	gwdefaults "github.com/solo-io/gloo/projects/gateway/pkg/defaults"
	gwreconciler "github.com/solo-io/gloo/projects/gateway/pkg/reconciler"
	"github.com/solo-io/gloo/projects/gateway/pkg/services/k8sadmission"
	gwsyncer "github.com/solo-io/gloo/projects/gateway/pkg/syncer"
	gwtranslator "github.com/solo-io/gloo/projects/gateway/pkg/translator"
	gwvalidation "github.com/solo-io/gloo/projects/gateway/pkg/validation"
	"github.com/solo-io/gloo/projects/gateway2/extensions"
	ggv2utils "github.com/solo-io/gloo/projects/gateway2/utils"
	"github.com/solo-io/gloo/projects/gloo/constants"
	rlv1alpha1 "github.com/solo-io/gloo/projects/gloo/pkg/api/external/solo/ratelimit"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	extauth "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/extauth/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/graphql/v1beta1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/ratelimit"
	v1snap "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/gloosnapshot"
	"github.com/solo-io/gloo/projects/gloo/pkg/bootstrap"
	bootstrap_clients "github.com/solo-io/gloo/projects/gloo/pkg/bootstrap/clients"
	"github.com/solo-io/gloo/projects/gloo/pkg/bootstrap/clients/vault"
	"github.com/solo-io/gloo/projects/gloo/pkg/defaults"
	"github.com/solo-io/gloo/projects/gloo/pkg/discovery"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins"
	consulplugin "github.com/solo-io/gloo/projects/gloo/pkg/plugins/consul"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins/registry"
	"github.com/solo-io/gloo/projects/gloo/pkg/syncer"
	extauthExt "github.com/solo-io/gloo/projects/gloo/pkg/syncer/extauth"
	ratelimitExt "github.com/solo-io/gloo/projects/gloo/pkg/syncer/ratelimit"
	"github.com/solo-io/gloo/projects/gloo/pkg/syncer/sanitizer"
	syncerValidation "github.com/solo-io/gloo/projects/gloo/pkg/syncer/validation"
	"github.com/solo-io/gloo/projects/gloo/pkg/translator"
	"github.com/solo-io/gloo/projects/gloo/pkg/upstreams"
	"github.com/solo-io/gloo/projects/gloo/pkg/upstreams/consul"
	"github.com/solo-io/gloo/projects/gloo/pkg/validation"
	"github.com/solo-io/gloo/projects/gloo/pkg/xds"
)

// TODO: (copied from gateway) switch AcceptAllResourcesByDefault to false after validation has been tested in user environments
var AcceptAllResourcesByDefault = true

var AllowWarnings = true

type RunFunc func(opts bootstrap.Opts) error

func NewSetupFunc(setupOpts *bootstrap.SetupOpts) setuputils.SetupFunc {
	return NewSetupFuncWithRunAndExtensions(RunGloo, setupOpts, nil)
}

// used outside of this repo
//
//goland:noinspection GoUnusedExportedFunction
func NewSetupFuncWithExtensions(setupOpts *bootstrap.SetupOpts, extensions Extensions) setuputils.SetupFunc {
	runWithExtensions := func(opts bootstrap.Opts) error {
		return RunGlooWithExtensions(opts, extensions)
	}
	return NewSetupFuncWithRunAndExtensions(runWithExtensions, setupOpts, &extensions)
}

// for use by UDS, FDS, other v1.SetupSyncers
func NewSetupFuncWithRun(runFunc RunFunc) setuputils.SetupFunc {
	return NewSetupFuncWithRunAndExtensions(runFunc, &bootstrap.SetupOpts{}, nil)
}

// Called directly by GlooEE
func NewSetupFuncWithRunAndExtensions(runFunc RunFunc, setupOpts *bootstrap.SetupOpts, extensions *Extensions) setuputils.SetupFunc {
	if setupOpts == nil {
		setupOpts = &bootstrap.SetupOpts{}
	}

	s := &setupSyncer{
		extensions: extensions,
		makeGrpcServer: func(ctx context.Context, options ...grpc.ServerOption) *grpc.Server {
			serverOpts := []grpc.ServerOption{
				grpc.StreamInterceptor(
					grpc_middleware.ChainStreamServer(
						grpc_ctxtags.StreamServerInterceptor(),
						grpc_zap.StreamServerInterceptor(zap.NewNop()),
						func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
							contextutils.LoggerFrom(ctx).Debugf("gRPC call: %v", info.FullMethod)
							return handler(srv, ss)
						},
					)),
			}
			serverOpts = append(serverOpts, options...)
			return grpc.NewServer(serverOpts...)
		},
		runFunc:   runFunc,
		setupOpts: setupOpts,
	}
	return s.Setup
}

// grpcServer contains grpc server configuration fields we will need to persist after starting a server
// to later check if they changed and we need to trigger a server restart
type grpcServer struct {
	addr            string
	maxGrpcRecvSize int
	cancel          context.CancelFunc
}

type setupSyncer struct {
	extensions               *Extensions
	runFunc                  RunFunc
	setupOpts                *bootstrap.SetupOpts
	makeGrpcServer           func(ctx context.Context, options ...grpc.ServerOption) *grpc.Server
	previousXdsServer        grpcServer
	previousControlPlane     bootstrap.ControlPlane
	previousValidationServer grpcServer
	previousProxyDebugServer grpcServer
	controlPlane             bootstrap.ControlPlane
	validationServer         bootstrap.ValidationServer
	proxyDebugServer         bootstrap.ProxyDebugServer
	callbacks                xdsserver.Callbacks
}

func NewControlPlane(ctx context.Context, snapshotCache cache.SnapshotCache, grpcServer *grpc.Server, bindAddr net.Addr, kubeControlPlaneCfg bootstrap.KubernetesControlPlaneConfig,
	callbacks xdsserver.Callbacks, start bool) bootstrap.ControlPlane {
	xdsServer := server.NewServer(ctx, snapshotCache, callbacks)
	reflection.Register(grpcServer)

	return bootstrap.ControlPlane{
		GrpcService: &bootstrap.GrpcService{
			GrpcServer:      grpcServer,
			StartGrpcServer: start,
			BindAddr:        bindAddr,
			Ctx:             ctx,
		},
		SnapshotCache: snapshotCache,
		XDSServer:     xdsServer,
		Kube:          kubeControlPlaneCfg,
	}
}

func NewValidationServer(ctx context.Context, grpcServer *grpc.Server, bindAddr net.Addr, start bool) bootstrap.ValidationServer {
	return bootstrap.ValidationServer{
		GrpcService: &bootstrap.GrpcService{
			GrpcServer:      grpcServer,
			StartGrpcServer: start,
			BindAddr:        bindAddr,
			Ctx:             ctx,
		},
		Server: validation.NewValidationServer(),
	}
}

func NewProxyDebugServer(ctx context.Context, grpcServer *grpc.Server, bindAddr net.Addr, start bool) bootstrap.ProxyDebugServer {
	return bootstrap.ProxyDebugServer{
		GrpcService: &bootstrap.GrpcService{
			Ctx:             ctx,
			BindAddr:        bindAddr,
			GrpcServer:      grpcServer,
			StartGrpcServer: start,
		},
		Server: debug.NewProxyEndpointServer(),
	}
}

var (
	DefaultXdsBindAddr        = fmt.Sprintf("0.0.0.0:%v", defaults.GlooXdsPort)
	DefaultValidationBindAddr = fmt.Sprintf("0.0.0.0:%v", defaults.GlooValidationPort)
	DefaultRestXdsBindAddr    = fmt.Sprintf("0.0.0.0:%v", defaults.GlooRestXdsPort)
	DefaultProxyDebugAddr     = fmt.Sprintf("0.0.0.0:%v", defaults.GlooProxyDebugPort)
)

func getAddr(addr string) (*net.TCPAddr, error) {
	addrParts := strings.Split(addr, ":")
	if len(addrParts) != 2 {
		return nil, errors.Errorf("invalid bind addr: %v", addr)
	}
	ip := net.ParseIP(addrParts[0])

	port, err := strconv.Atoi(addrParts[1])
	if err != nil {
		return nil, errors.Wrapf(err, "invalid bind addr: %v", addr)
	}

	return &net.TCPAddr{IP: ip, Port: port}, nil
}
func GetWriteNamespace(settings *v1.Settings) string {
	writeNamespace := settings.GetDiscoveryNamespace()
	if writeNamespace == "" {
		writeNamespace = defaults.GlooSystem
	}
	return writeNamespace
}

// Setup constructs bootstrap options based on settings and other input, and calls the runFunc with these options.
func (s *setupSyncer) Setup(ctx context.Context, kubeCache kube.SharedCache, memCache memory.InMemoryResourceCache, settings *v1.Settings, identity leaderelector.Identity) error {
	xdsBindAddr := settings.GetGloo().GetXdsBindAddr()
	if xdsBindAddr == "" {
		xdsBindAddr = DefaultXdsBindAddr
	}
	xdsTcpAddress, err := getAddr(xdsBindAddr)
	if err != nil {
		return errors.Wrapf(err, "parsing xds addr")
	}

	validationAddr := settings.GetGloo().GetValidationBindAddr()
	if validationAddr == "" {
		validationAddr = DefaultValidationBindAddr
	}
	validationTcpAddress, err := getAddr(validationAddr)
	if err != nil {
		return errors.Wrapf(err, "parsing validation addr")
	}

	proxyDebugAddr := settings.GetGloo().GetProxyDebugBindAddr()
	if proxyDebugAddr == "" {
		proxyDebugAddr = DefaultProxyDebugAddr
	}
	proxyDebugTcpAddress, err := getAddr(proxyDebugAddr)
	if err != nil {
		return errors.Wrapf(err, "parsing proxy debug endpoint addr")
	}
	refreshRate := time.Minute
	if settings.GetRefreshRate() != nil {
		refreshRate = prototime.DurationFromProto(settings.GetRefreshRate())
	}

	writeNamespace := GetWriteNamespace(settings)
	watchNamespaces := namespaces.ProcessWatchNamespaces(settingsutil.GetNamespacesToWatch(settings), writeNamespace)

	consulClient, err := bootstrap_clients.ConsulClientForSettings(ctx, settings)
	if err != nil {
		return err
	}

	getVaultInit := func(vaultSettings *v1.Settings_VaultSecrets) vault.VaultClientInitFunc {
		return func(initCtx context.Context) *vaultapi.Client {
			c, err := vault.VaultClientForSettings(initCtx, vaultSettings)
			if err != nil {
				// We log this error here, but we do not have a feasible way to raise
				// it when this function is called in NewVaultSecretClientFactory.
				// The error is handled after we create the factory in getFactoryForSource
				// for the multi client, and NewSecretResourceClientFactory for the
				// traditional single client.
				contextutils.LoggerFrom(initCtx).Error(err)
			}
			return c
		}
	}
	vaultInitMap := make(map[int]vault.VaultClientInitFunc)
	vaultSettings := settings.GetVaultSecretSource()
	if vaultSettings != nil {
		vaultInitMap[bootstrap_clients.SecretSourceAPIVaultClientInitIndex] = getVaultInit(vaultSettings)
	}
	if secretSources := settings.GetSecretOptions().GetSources(); secretSources != nil {
		sourceList := bootstrap_clients.SourceList(secretSources)
		sort.Stable(sourceList)
		for i := range sourceList {
			switch src := sourceList[i].GetSource().(type) {
			case *v1.Settings_SecretOptions_Source_Vault:
				vaultInitMap[i] = getVaultInit(src.Vault)
			default:
				// We are not concerned with source types that are not Vault since
				// we are only getting the indices of Vault sources to accurately
				// map the vault source to its API client setup func
				continue
			}
		}
		settings.GetSecretOptions().Sources = sourceList
	}

	var clientset kubernetes.Interface
	opts, err := constructOpts(ctx, s.setupOpts,
		constructOptsParams{
			clientset:          &clientset,
			kubeCache:          kubeCache,
			consulClient:       consulClient,
			vaultClientInitMap: vaultInitMap,
			memCache:           memCache,
			settings:           settings,
			writeNamespace:     writeNamespace,
		},
	)
	if err != nil {
		return err
	}

	// get xds host/port to pass to ControlPlane below, only if we're running in k8s
	var xdsHost string
	var xdsPort int32
	switch settings.GetConfigSource().(type) {
	case *v1.Settings_KubernetesConfigSource:
		// Since the installation namespace and the write namespace can be different, we pass the installation namespace
		// to fetch the service
		glooService, err := GetControlPlaneService(ctx, namespaces.GetPodNamespace(), opts.KubeServiceClient)
		if err != nil {
			return err
		}

		xdsHost = GetControlPlaneXdsHost(glooService)
		xdsPort, err = GetControlPlaneXdsPort(glooService)
		if err != nil {
			return err
		}

		logger := contextutils.LoggerFrom(ctx)
		logger.Infof("using xds host %v and xds port %v", xdsHost, xdsPort)
	}

	// process grpcserver options to understand if any servers will need a restart

	maxGrpcRecvSize := -1
	// Use the same maxGrpcMsgSize for both validation server and proxy debug server as the message size is determined by the size of proxies.
	if maxGrpcMsgSize := settings.GetGateway().GetValidation().GetValidationServerGrpcMaxSizeBytes(); maxGrpcMsgSize != nil {
		if maxGrpcMsgSize.GetValue() < 0 {
			return errors.Errorf("validationServerGrpcMaxSizeBytes in settings CRD must be non-negative, current value: %v", maxGrpcMsgSize.GetValue())
		}
		maxGrpcRecvSize = int(maxGrpcMsgSize.GetValue())
	}

	emptyControlPlane := bootstrap.ControlPlane{}
	emptyValidationServer := bootstrap.ValidationServer{}
	emptyProxyDebugServer := bootstrap.ProxyDebugServer{}

	// check if we need to restart the control plane
	if xdsBindAddr != s.previousXdsServer.addr ||
		xdsHost != s.previousControlPlane.Kube.XdsHost ||
		xdsPort != s.previousControlPlane.Kube.XdsPort {
		if s.previousXdsServer.cancel != nil {
			s.previousXdsServer.cancel()
			s.previousXdsServer.cancel = nil
		}
		s.controlPlane = emptyControlPlane
	}

	// check if we need to restart the validation server
	if validationAddr != s.previousValidationServer.addr || maxGrpcRecvSize != s.previousValidationServer.maxGrpcRecvSize {
		if s.previousValidationServer.cancel != nil {
			s.previousValidationServer.cancel()
			s.previousValidationServer.cancel = nil
		}
		s.validationServer = emptyValidationServer
	}

	// check if we need to restart the proxy debug server
	if proxyDebugAddr != s.previousProxyDebugServer.addr || maxGrpcRecvSize != s.previousProxyDebugServer.maxGrpcRecvSize {
		if s.previousProxyDebugServer.cancel != nil {
			s.previousProxyDebugServer.cancel()
			s.previousProxyDebugServer.cancel = nil
		}
		s.proxyDebugServer = emptyProxyDebugServer
	}

	// initialize the control plane context in this block either on the first loop, or if bind addr changed
	if s.controlPlane == emptyControlPlane {
		// create new context as the grpc server might survive multiple iterations of this loop.
		ctx, cancel := context.WithCancel(context.Background())
		var callbacks xdsserver.Callbacks
		if s.extensions != nil {
			callbacks = s.extensions.XdsCallbacks
		}
		s.controlPlane = NewControlPlane(ctx, s.setupOpts.Cache, s.makeGrpcServer(ctx), xdsTcpAddress,
			bootstrap.KubernetesControlPlaneConfig{XdsHost: xdsHost, XdsPort: xdsPort}, multiCallbacks(s.setupOpts.ExtraCallbacks, callbacks), true)

		s.setupOpts.SetXdsAddress(xdsHost, xdsPort)

		s.previousXdsServer.cancel = cancel
		s.previousXdsServer.addr = xdsBindAddr
		s.previousControlPlane.Kube.XdsHost = xdsHost
		s.previousControlPlane.Kube.XdsPort = xdsPort
	}

	// initialize the validation server context in this block either on the first loop, or if bind addr changed
	if s.validationServer == emptyValidationServer {
		// create new context as the grpc server might survive multiple iterations of this loop.
		ctx, cancel := context.WithCancel(context.Background())
		var validationGrpcServerOpts []grpc.ServerOption
		// if validationServerGrpcMaxSizeBytes was set this will be non-negative, otherwise use gRPC default
		if maxGrpcRecvSize >= 0 {
			validationGrpcServerOpts = append(validationGrpcServerOpts, grpc.MaxRecvMsgSize(maxGrpcRecvSize))
		}
		s.validationServer = NewValidationServer(ctx, s.makeGrpcServer(ctx, validationGrpcServerOpts...), validationTcpAddress, true)
		s.previousValidationServer.cancel = cancel
		s.previousValidationServer.addr = validationAddr
		s.previousValidationServer.maxGrpcRecvSize = maxGrpcRecvSize
	}
	// initialize the proxy debug server context in this block either on the first loop, or if bind addr changed
	if s.proxyDebugServer == emptyProxyDebugServer {
		// create new context as the grpc server might survive multiple iterations of this loop.
		ctx, cancel := context.WithCancel(context.Background())

		proxyGrpcServerOpts := []grpc.ServerOption{grpc.MaxRecvMsgSize(maxGrpcRecvSize)}
		s.proxyDebugServer = NewProxyDebugServer(ctx, s.makeGrpcServer(ctx, proxyGrpcServerOpts...), proxyDebugTcpAddress, true)
		s.previousProxyDebugServer.cancel = cancel
		s.previousProxyDebugServer.addr = proxyDebugAddr
		s.previousProxyDebugServer.maxGrpcRecvSize = maxGrpcRecvSize
	}

	// populate rest of opts
	opts.Identity = identity
	opts.WriteNamespace = writeNamespace
	opts.StatusReporterNamespace = gloostatusutils.GetStatusReporterNamespaceOrDefault(writeNamespace)
	opts.WatchNamespaces = watchNamespaces
	opts.WatchOpts = clients.WatchOpts{
		Ctx:         ctx,
		RefreshRate: refreshRate,
	}
	opts.ControlPlane = s.controlPlane
	opts.ValidationServer = s.validationServer
	opts.ProxyDebugServer = s.proxyDebugServer
	// if nil, kube plugin disabled
	opts.KubeClient = clientset
	opts.DevMode = settings.GetDevMode()
	opts.Settings = settings
	opts.KrtDebugger = s.setupOpts.KrtDebugger

	opts.Consul.DnsServer = settings.GetConsul().GetDnsAddress()
	if len(opts.Consul.DnsServer) == 0 {
		opts.Consul.DnsServer = consulplugin.DefaultDnsAddress
	}
	opts.Consul.DnsPollingInterval = settings.GetConsul().GetDnsPollingInterval()

	// if consul service discovery specified, initialize consul watcher
	if consulServiceDiscovery := settings.GetConsul().GetServiceDiscovery(); consulServiceDiscovery != nil {
		consulClientWrapper, err := consul.NewConsulWatcher(consulClient, consulServiceDiscovery.GetDataCenters(), settings.GetConsulDiscovery().GetServiceTagsAllowlist())
		if err != nil {
			return err
		}
		opts.Consul.ConsulWatcher = consulClientWrapper
	}

	err = s.runFunc(opts)

	s.validationServer.StartGrpcServer = opts.ValidationServer.StartGrpcServer
	s.controlPlane.StartGrpcServer = opts.ControlPlane.StartGrpcServer

	return err
}

func RunGloo(opts bootstrap.Opts) error {
	glooExtensions := Extensions{
		K8sGatewayExtensionsFactory: extensions.NewK8sGatewayExtensions,
		PluginRegistryFactory:       registry.GetPluginRegistryFactory(registry.FromBootstrap(opts)),
		SyncerExtensions: []syncer.TranslatorSyncerExtensionFactory{
			ratelimitExt.NewTranslatorSyncerExtension,
			extauthExt.NewTranslatorSyncerExtension,
		},
		ApiEmitterChannel:      make(chan struct{}),
		XdsCallbacks:           nil,
		SnapshotHistoryFactory: iosnapshot.GetHistoryFactory(),
	}

	return RunGlooWithExtensions(opts, glooExtensions)
}

// RunGlooWithExtensions is the core entrypoint to the Gloo components.
// THIS FUNCTION MUST NOT BLOCK:
//
//	It is invoked by an outer control loop (SetupFunc) which monitors
//	the current Settings resource, and re-runs this function each time the global Settings change
//
// This function is called directly by GlooEE
func RunGlooWithExtensions(opts bootstrap.Opts, extensions Extensions) error {
	if err := extensions.Validate(); err != nil {
		return err
	}

	watchOpts := opts.WatchOpts.WithDefaults()
	watchOpts.Ctx = contextutils.WithLogger(watchOpts.Ctx, "setup")
	opts.WatchOpts.Ctx = contextutils.WithLogger(opts.WatchOpts.Ctx, "gloo")

	runErrorGroup, _ := errgroup.WithContext(watchOpts.Ctx)
	logger := contextutils.LoggerFrom(watchOpts.Ctx)

	// MARK: build resource clients
	upstreamClient, err := v1.NewUpstreamClient(watchOpts.Ctx, opts.Upstreams)
	if err != nil {
		return err
	}
	if err := upstreamClient.Register(); err != nil {
		return err
	}

	kubeServiceClient := opts.KubeServiceClient
	if opts.Settings.GetGloo().GetDisableKubernetesDestinations() {
		kubeServiceClient = nil
	}
	hybridUsClient, err := upstreams.NewHybridUpstreamClient(upstreamClient, kubeServiceClient, opts.Consul.ConsulWatcher, opts.Settings)
	if err != nil {
		return err
	}

	proxyClient, err := v1.NewProxyClient(watchOpts.Ctx, opts.Proxies)
	if err != nil {
		return err
	}
	if err := proxyClient.Register(); err != nil {
		return err
	}

	upstreamGroupClient, err := v1.NewUpstreamGroupClient(watchOpts.Ctx, opts.UpstreamGroups)
	if err != nil {
		return err
	}
	if err := upstreamGroupClient.Register(); err != nil {
		return err
	}

	// create in-memory cache for endpoints
	// see (https://github.com/solo-io/gloo/blob/main/devel/architecture/endpoint-discovery.md) for more info
	endpointsFactory := &factory.MemoryResourceClientFactory{
		Cache: memory.NewInMemoryResourceCache(),
	}
	endpointClient, err := v1.NewEndpointClient(watchOpts.Ctx, endpointsFactory)
	if err != nil {
		return err
	}

	secretClient, err := v1.NewSecretClient(watchOpts.Ctx, opts.Secrets)
	if err != nil {
		return err
	}

	artifactClient, err := v1.NewArtifactClient(watchOpts.Ctx, opts.Artifacts)
	if err != nil {
		return err
	}

	authConfigClient, err := extauth.NewAuthConfigClient(watchOpts.Ctx, opts.AuthConfigs)
	if err != nil {
		return err
	}
	if err := authConfigClient.Register(); err != nil {
		return err
	}

	graphqlApiClient, err := v1beta1.NewGraphQLApiClient(watchOpts.Ctx, opts.GraphQLApis)
	if err != nil {
		return err
	}
	if err := graphqlApiClient.Register(); err != nil {
		return err
	}

	rlClient, rlReporterClient, err := rlv1alpha1.NewRateLimitClients(watchOpts.Ctx, opts.RateLimitConfigs)
	if err != nil {
		return err
	}
	if err := rlClient.Register(); err != nil {
		return err
	}

	virtualServiceClient, err := gateway.NewVirtualServiceClient(watchOpts.Ctx, opts.VirtualServices)
	if err != nil {
		return err
	}
	if err := virtualServiceClient.Register(); err != nil {
		return err
	}

	rtClient, err := gateway.NewRouteTableClient(watchOpts.Ctx, opts.RouteTables)
	if err != nil {
		return err
	}
	if err := rtClient.Register(); err != nil {
		return err
	}

	gatewayClient, err := gateway.NewGatewayClient(watchOpts.Ctx, opts.Gateways)
	if err != nil {
		return err
	}
	if err := gatewayClient.Register(); err != nil {
		return err
	}

	matchableHttpGatewayClient, err := gateway.NewMatchableHttpGatewayClient(watchOpts.Ctx, opts.MatchableHttpGateways)
	if err != nil {
		return err
	}
	if err := matchableHttpGatewayClient.Register(); err != nil {
		return err
	}

	matchableTcpGatewayClient, err := gateway.NewMatchableTcpGatewayClient(watchOpts.Ctx, opts.MatchableTcpGateways)
	if err != nil {
		return err
	}
	if err := matchableTcpGatewayClient.Register(); err != nil {
		return err
	}

	virtualHostOptionClient, err := gateway.NewVirtualHostOptionClient(watchOpts.Ctx, opts.VirtualHostOptions)
	if err != nil {
		return err
	}
	if err := virtualHostOptionClient.Register(); err != nil {
		return err
	}

	routeOptionClient, err := gateway.NewRouteOptionClient(watchOpts.Ctx, opts.RouteOptions)
	if err != nil {
		return err
	}
	if err := routeOptionClient.Register(); err != nil {
		return err
	}
	if opts.ProxyCleanup != nil {
		opts.ProxyCleanup()
	}

	statusClient := gloostatusutils.GetStatusClientForNamespace(opts.StatusReporterNamespace)

	// Register grpc endpoints to the grpc server
	xds.SetupEnvoyXds(opts.ControlPlane.GrpcServer, opts.ControlPlane.XDSServer, opts.ControlPlane.SnapshotCache)

	pluginRegistry := extensions.PluginRegistryFactory(watchOpts.Ctx)
	var discoveryPlugins []discovery.DiscoveryPlugin
	for _, plug := range pluginRegistry.GetPlugins() {
		disc, ok := plug.(discovery.DiscoveryPlugin)
		if ok {
			disc.Init(plugins.InitParams{Ctx: watchOpts.Ctx, Settings: opts.Settings})
			discoveryPlugins = append(discoveryPlugins, disc)
		}
	}

	startRestXdsServer(opts)

	errs := make(chan error)

	// MARK: build and run EDS loop
	disc := discovery.NewEndpointDiscovery(opts.WatchNamespaces, opts.WriteNamespace, endpointClient, statusClient, discoveryPlugins)
	edsSync := discovery.NewEdsSyncer(disc, discovery.Opts{}, watchOpts.RefreshRate)
	edsEmitter := v1.NewEdsEmitter(hybridUsClient)
	edsEventLoop := v1.NewEdsEventLoop(edsEmitter, edsSync)
	edsErrs, err := edsEventLoop.Run(opts.WatchNamespaces, watchOpts)
	if err != nil {
		return err
	}

	warmTimeout := opts.Settings.GetGloo().GetEndpointsWarmingTimeout()

	if warmTimeout == nil {
		warmTimeout = &duration.Duration{
			Seconds: 5 * 60,
		}
	}
	if warmTimeout.GetSeconds() != 0 || warmTimeout.GetNanos() != 0 {
		warmTimeoutDuration := prototime.DurationFromProto(warmTimeout)
		ctx := watchOpts.Ctx
		err = channelutils.WaitForReady(ctx, warmTimeoutDuration, edsEventLoop.Ready(), disc.Ready())
		if err != nil {
			// make sure that the reason we got here is not context cancellation
			if ctx.Err() != nil {
				return ctx.Err()
			}
			logger.Panicw("failed warming up endpoints - consider adjusting endpointsWarmingTimeout", "warmTimeoutDuration", warmTimeoutDuration)
		}
	}

	// We are ready!

	go errutils.AggregateErrs(watchOpts.Ctx, errs, edsErrs, "eds.gloo")

	rpt := reporter.NewReporter(defaults.GlooReporter,
		statusClient,
		hybridUsClient.BaseClient(),
		proxyClient.BaseClient(),
		upstreamGroupClient.BaseClient(),
		authConfigClient.BaseClient(),
		gatewayClient.BaseClient(),
		matchableHttpGatewayClient.BaseClient(),
		matchableTcpGatewayClient.BaseClient(),
		virtualServiceClient.BaseClient(),
		rtClient.BaseClient(),
		virtualHostOptionClient.BaseClient(),
		routeOptionClient.BaseClient(),
		rlReporterClient,
	)
	statusMetrics, err := metrics.NewConfigStatusMetrics(opts.Settings.GetObservabilityOptions().GetConfigStatusMetricLabels())
	if err != nil {
		return err
	}
	//The validation grpc server is available for custom controllers
	if opts.ValidationServer.StartGrpcServer {
		validationServer := opts.ValidationServer
		lis, err := net.Listen(validationServer.BindAddr.Network(), validationServer.BindAddr.String())
		if err != nil {
			return err
		}
		validationServer.Server.Register(validationServer.GrpcServer)

		go func() {
			<-validationServer.Ctx.Done()
			validationServer.GrpcServer.Stop()
		}()

		go func() {
			if err := validationServer.GrpcServer.Serve(lis); err != nil {
				logger.Errorf("validation grpc server failed to start")
			}
		}()
		opts.ValidationServer.StartGrpcServer = false
	}

	if opts.ControlPlane.StartGrpcServer {
		// copy for the go-routines
		controlPlane := opts.ControlPlane
		lis, err := net.Listen(opts.ControlPlane.BindAddr.Network(), opts.ControlPlane.BindAddr.String())
		if err != nil {
			return err
		}
		go func() {
			<-controlPlane.GrpcService.Ctx.Done()
			controlPlane.GrpcServer.Stop()
		}()

		go func() {
			if err := controlPlane.GrpcServer.Serve(lis); err != nil {
				logger.Errorf("xds grpc server failed to start")
			}
		}()
		opts.ControlPlane.StartGrpcServer = false
	}
	if opts.ProxyDebugServer.StartGrpcServer {
		proxyDebugServer := opts.ProxyDebugServer
		proxyDebugServer.Server.RegisterProxyReader(proxyClient)
		proxyDebugServer.Server.Register(proxyDebugServer.GrpcServer)
		lis, err := net.Listen(opts.ProxyDebugServer.BindAddr.Network(), opts.ProxyDebugServer.BindAddr.String())
		if err != nil {
			return err
		}
		go func() {
			<-proxyDebugServer.GrpcService.Ctx.Done()
			proxyDebugServer.GrpcServer.Stop()
		}()

		go func() {
			if err := proxyDebugServer.GrpcServer.Serve(lis); err != nil {
				logger.Errorf("Proxy debug grpc server failed to start")
			}
		}()
		opts.ProxyDebugServer.StartGrpcServer = false
	}

	resourceHasher := translator.EnvoyCacheResourcesListToFnvHash

	// Set up the syncer extensions
	syncerExtensionParams := syncer.TranslatorSyncerExtensionParams{
		RateLimitServiceSettings: &ratelimit.ServiceSettings{
			Descriptors:    opts.Settings.GetRatelimit().GetDescriptors(),
			SetDescriptors: opts.Settings.GetRatelimit().GetSetDescriptors(),
		},
		Hasher: resourceHasher,
	}
	var syncerExtensions []syncer.TranslatorSyncerExtension
	for _, syncerExtensionFactory := range extensions.SyncerExtensions {
		syncerExtension := syncerExtensionFactory(watchOpts.Ctx, syncerExtensionParams)
		syncerExtensions = append(syncerExtensions, syncerExtension)
	}

	// MARK: build gloo translator
	sharedTranslator := TranslatorFactory{PluginRegistry: extensions.PluginRegistryFactory}.NewTranslator(watchOpts.Ctx,
		opts.Settings)
	routeReplacingSanitizer, err := sanitizer.NewRouteReplacingSanitizer(opts.Settings.GetGloo().GetInvalidConfigPolicy())
	if err != nil {
		return err
	}

	xdsSanitizers := sanitizer.XdsSanitizers{
		sanitizer.NewUpstreamRemovingSanitizer(),
		routeReplacingSanitizer,
	}

	vc := validation.ValidatorConfig{
		Ctx: watchOpts.Ctx,
		GlooValidatorConfig: validation.GlooValidatorConfig{
			XdsSanitizer: xdsSanitizers,
			Translator:   sharedTranslator,
			Settings:     opts.Settings,
		},
	}
	validator := validation.NewValidator(vc)
	if opts.ValidationServer.Server != nil {
		opts.ValidationServer.Server.SetValidator(validator)
	}

	// MARK: build gateway translator
	gwOpts := gwtranslator.Opts{
		GlooNamespace:                  opts.WriteNamespace,
		WriteNamespace:                 opts.WriteNamespace,
		StatusReporterNamespace:        opts.StatusReporterNamespace,
		WatchNamespaces:                opts.WatchNamespaces,
		Gateways:                       opts.Gateways,
		VirtualServices:                opts.VirtualServices,
		RouteTables:                    opts.RouteTables,
		Proxies:                        opts.Proxies,
		RouteOptions:                   opts.RouteOptions,
		VirtualHostOptions:             opts.VirtualHostOptions,
		WatchOpts:                      opts.WatchOpts,
		DevMode:                        opts.DevMode,
		ReadGatewaysFromAllNamespaces:  opts.ReadGatwaysFromAllNamespaces,
		Validation:                     opts.ValidationOpts,
		ConfigStatusMetricOpts:         nil,
		IsolateVirtualHostsBySslConfig: opts.Settings.GetGateway().GetIsolateVirtualHostsBySslConfig().GetValue(),
	}
	var (
		gwTranslatorSyncer *gwsyncer.TranslatorSyncer
		gatewayTranslator  *gwtranslator.GwTranslator
	)
	if opts.GatewayControllerEnabled {
		logger.Debugf("Setting up gateway translator")
		gatewayTranslator = gwtranslator.NewDefaultTranslator(gwOpts)
		proxyReconciler := gwreconciler.NewProxyReconciler(validator.Validate, proxyClient, statusClient)
		gwTranslatorSyncer = gwsyncer.NewTranslatorSyncer(
			watchOpts.Ctx,
			opts.WriteNamespace,
			proxyClient,
			proxyReconciler,
			rpt,
			gatewayTranslator,
			statusClient,
			statusMetrics,
			opts.Identity)
	} else {
		logger.Debugf("Gateway translation is disabled. Proxies are provided from another source")
	}

	// filter the list of extensions to only include the rate limit extension for validation
	syncerValidatorExtensions := []syncer.TranslatorSyncerExtension{}
	for _, ext := range syncerExtensions {
		// currently only supporting ratelimit extension in validation
		if ext.ID() == ratelimitExt.ServerRole {
			syncerValidatorExtensions = append(syncerValidatorExtensions, ext)
		}
	}
	// create a validator to validate extensions
	extensionValidator := syncerValidation.NewValidator(syncerValidatorExtensions, opts.Settings)

	// allow by default
	disableValidationAgainstPreviousState := os.Getenv("DISABLE_VALIDATION_AGAINST_PREVIOUS_STATE") == "true"
	validationConfig := gwvalidation.ValidatorConfig{
		Translator:                            gatewayTranslator,
		GlooValidator:                         validator.ValidateGloo,
		ExtensionValidator:                    extensionValidator,
		DisableValidationAgainstPreviousState: disableValidationAgainstPreviousState,
	}
	if gwOpts.Validation != nil {
		valOpts := gwOpts.Validation
		if opts.GatewayControllerEnabled {
			validationConfig.AllowWarnings = valOpts.AllowWarnings
		}
	}
	gwValidationSyncer := gwvalidation.NewValidator(validationConfig)

	// startFuncs represents the set of StartFunc that should be executed at startup
	// At the moment, the functionality is used minimally.
	// Overtime, we should break up this large function into smaller StartFunc
	startFuncs := map[string]StartFunc{}

	// snapshotHistory is a utility for managing the state of the input/output snapshots that the Control Plane
	// consumes and produces. This object is then used by our Admin Server, to provide this data on demand
	snapshotHistory := extensions.SnapshotHistoryFactory(iosnapshot.HistoryFactoryParameters{
		Settings:                    opts.Settings,
		Cache:                       opts.ControlPlane.SnapshotCache,
		EnableK8sGatewayIntegration: opts.GlooGateway.EnableK8sGatewayController,
	})

	startFuncs["admin-server"] = AdminServerStartFunc(snapshotHistory, opts.KrtDebugger)

	if opts.ProxyReconcileQueue != nil {
		go runQueue(watchOpts.Ctx, opts.ProxyReconcileQueue, opts.WriteNamespace, proxyClient)
	}

	// MARK: build translator syncer
	translationSync := syncer.NewTranslatorSyncer(
		watchOpts.Ctx,
		sharedTranslator,
		opts.ControlPlane.SnapshotCache,
		xdsSanitizers,
		rpt,
		// opts.DevMode should be deprecated
		// https://github.com/solo-io/gloo/issues/6494
		// We are starting to build out a true Admin Server, and enhancements should be added to that server
		opts.DevMode,
		syncerExtensions,
		opts.Settings,
		statusMetrics,
		gwTranslatorSyncer,
		proxyClient,
		opts.WriteNamespace,
		opts.Identity,
		snapshotHistory,
	)

	// MARK: build & run api snap loop
	apiEmitter := v1snap.NewApiEmitterWithEmit(
		artifactClient,
		endpointClient,
		proxyClient,
		upstreamGroupClient,
		secretClient,
		hybridUsClient,
		authConfigClient,
		rlClient,
		virtualServiceClient,
		rtClient,
		gatewayClient,
		virtualHostOptionClient,
		routeOptionClient,
		matchableHttpGatewayClient,
		matchableTcpGatewayClient,
		graphqlApiClient,
		extensions.ApiEmitterChannel,
	)

	syncers := v1snap.ApiSyncers{
		validator,
		translationSync,
	}
	if opts.GatewayControllerEnabled {
		syncers = append(syncers, gwValidationSyncer)
	}

	apiEventLoop := v1snap.NewApiEventLoop(apiEmitter, syncers)
	apiEventLoopErrs, err := apiEventLoop.Run(opts.WatchNamespaces, watchOpts)
	if err != nil {
		return err
	}
	go errutils.AggregateErrs(watchOpts.Ctx, errs, apiEventLoopErrs, "event_loop.gloo")

	go func() {
		for {
			select {
			case <-watchOpts.Ctx.Done():
				logger.Debugf("context cancelled")
				return
			}
		}
	}()

	// MARK: start validation server
	validationMustStart := os.Getenv("VALIDATION_MUST_START")
	// only starting validation server if the env var is true or empty (previously, it always started, so this avoids causing unwanted changes for users)
	if validationMustStart == "true" || validationMustStart == "" {
		// Start the validation webhook
		validationServerErr := make(chan error, 1)
		if gwOpts.Validation != nil {
			// make sure non-empty WatchNamespaces contains the gloo instance's own namespace if
			// ReadGatewaysFromAllNamespaces is false
			if !gwOpts.ReadGatewaysFromAllNamespaces && !namespaces.AllNamespaces(opts.WatchNamespaces) {
				foundSelf := false
				for _, namespace := range opts.WatchNamespaces {
					if gwOpts.GlooNamespace == namespace {
						foundSelf = true
						break
					}
				}
				if !foundSelf {
					return errors.Errorf("The gateway configuration value readGatewaysFromAllNamespaces was set "+
						"to false, but the non-empty settings.watchNamespaces "+
						"list (%s) did not contain this gloo instance's own namespace: %s.",
						strings.Join(opts.WatchNamespaces, ", "), gwOpts.GlooNamespace)
				}
			}

			validationWebhook, err := k8sadmission.NewGatewayValidatingWebhook(
				k8sadmission.NewWebhookConfig(
					watchOpts.Ctx,
					gwValidationSyncer,
					gwOpts.WatchNamespaces,
					gwOpts.Validation.ValidatingWebhookPort,
					gwOpts.Validation.ValidatingWebhookCertPath,
					gwOpts.Validation.ValidatingWebhookKeyPath,
					gwOpts.Validation.AlwaysAcceptResources,
					gwOpts.ReadGatewaysFromAllNamespaces,
					gwOpts.GlooNamespace,
					opts.GlooGateway.EnableK8sGatewayController, // controls validation of KubeGateway policies (e.g. RouteOption, VirtualHostOption)
				),
			)
			if err != nil {
				return errors.Wrapf(err, "creating validating webhook")
			}

			go func() {
				// close out validation server when context is cancelled
				<-watchOpts.Ctx.Done()
				validationWebhook.Close()
			}()
			go func() {
				contextutils.LoggerFrom(watchOpts.Ctx).Infow("starting gateway validation server",
					zap.Int("port", gwOpts.Validation.ValidatingWebhookPort),
					zap.String("cert", gwOpts.Validation.ValidatingWebhookCertPath),
					zap.String("key", gwOpts.Validation.ValidatingWebhookKeyPath),
				)
				if err := validationWebhook.ListenAndServeTLS("", ""); err != nil && err != http.ErrServerClosed {
					select {
					case validationServerErr <- err:
					default:
						logger.DPanicw("failed to start validation webhook server", zap.Error(err))
					}
				}
			}()
		}

		// give the validation server 100ms to start
		select {
		case err := <-validationServerErr:
			return errors.Wrapf(err, "failed to start validation webhook server")
		case <-time.After(time.Millisecond * 100):
		}
	}

	ExecuteAsynchronousStartFuncs(
		watchOpts.Ctx,
		opts,
		extensions,
		startFuncs,
		runErrorGroup,
	)

	go func() {
		// It is critical that the RunGlooWithExtensions function does not block.
		// As a result, we monitor the runErrorGroup and just drop errors on the shared "errs" channel if one occurs
		runErr := runErrorGroup.Wait()
		if runErr != nil {
			errs <- runErr
		}
	}()

	go func() {
		for {
			select {
			case err, ok := <-errs:
				if !ok {
					return
				}
				logger.Errorw("gloo main event loop", zap.Error(err))
			case <-watchOpts.Ctx.Done():
				// think about closing this channel
				// close(errs)
				return
			}
		}
	}()

	logger.Infof("Gloo setup completed successfully")
	return nil
}

func startRestXdsServer(opts bootstrap.Opts) {
	restClient := server.NewHTTPGateway(
		contextutils.LoggerFrom(opts.WatchOpts.Ctx),
		opts.ControlPlane.XDSServer,
		map[string]string{
			types.FetchEndpointsV3: types.EndpointTypeV3,
		},
	)
	restXdsAddr := opts.Settings.GetGloo().GetRestXdsBindAddr()
	if restXdsAddr == "" {
		restXdsAddr = DefaultRestXdsBindAddr
	}
	srv := &http.Server{
		Addr:    restXdsAddr,
		Handler: restClient,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			// TODO: Add metrics for rest xds server
			contextutils.LoggerFrom(opts.WatchOpts.Ctx).Warnf("error while running REST xDS server", zap.Error(err))
		}
	}()
	go func() {
		<-opts.WatchOpts.Ctx.Done()
		if err := srv.Close(); err != nil {
			contextutils.LoggerFrom(opts.WatchOpts.Ctx).Warnf("error while shutting down REST xDS server", zap.Error(err))
		}
	}()
}

type constructOptsParams struct {
	clientset          *kubernetes.Interface
	kubeCache          kube.SharedCache
	consulClient       *consulapi.Client
	vaultClientInitMap map[int]vault.VaultClientInitFunc
	memCache           memory.InMemoryResourceCache
	settings           *v1.Settings
	writeNamespace     string
}

// constructs bootstrap opts from settings
func constructOpts(ctx context.Context, setup *bootstrap.SetupOpts, params constructOptsParams) (bootstrap.Opts, error) {

	var (
		cfg           *rest.Config
		kubeCoreCache corecache.KubeCoreCache
	)

	// MARK: build client factories
	factoryParams := bootstrap_clients.NewConfigFactoryParams(
		params.settings,
		params.memCache,
		params.kubeCache,
		&cfg,
		params.consulClient,
	)

	upstreamFactory, err := bootstrap_clients.ConfigFactoryForSettings(factoryParams, v1.UpstreamCrd)
	if err != nil {
		return bootstrap.Opts{}, errors.Wrapf(err, "creating config source from settings")
	}

	kubeServiceClient, err := bootstrap_clients.KubeServiceClientForSettings(
		ctx,
		params.settings,
		params.memCache,
		&cfg,
		params.clientset,
		&kubeCoreCache,
	)
	if err != nil {
		return bootstrap.Opts{}, err
	}

	var proxyFactory factory.ResourceClientFactory
	// Delete proxies that may have been left from prior to an upgrade or from previously having set persistProxySpec
	// Ignore errors because gloo will still work with stray proxies.
	proxyCleanup := func() {
		doProxyCleanup(ctx, factoryParams, params.settings, params.writeNamespace)
	}
	if params.settings.GetGateway().GetPersistProxySpec().GetValue() {
		proxyFactory, err = bootstrap_clients.ConfigFactoryForSettings(factoryParams, v1.ProxyCrd)
		if err != nil {
			return bootstrap.Opts{}, err
		}
	} else {
		proxyFactory = &factory.MemoryResourceClientFactory{
			Cache: memory.NewInMemoryResourceCache(),
		}
	}

	secretFactory, err := bootstrap_clients.SecretFactoryForSettings(ctx,
		bootstrap_clients.SecretFactoryParams{
			Settings:           params.settings,
			SharedCache:        params.memCache,
			Cfg:                &cfg,
			Clientset:          params.clientset,
			KubeCoreCache:      &kubeCoreCache,
			VaultClientInitMap: params.vaultClientInitMap,
			PluralName:         v1.SecretCrd.Plural,
		},
	)
	if err != nil {
		return bootstrap.Opts{}, err
	}

	upstreamGroupFactory, err := bootstrap_clients.ConfigFactoryForSettings(factoryParams, v1.UpstreamGroupCrd)
	if err != nil {
		return bootstrap.Opts{}, err
	}

	artifactFactory, err := bootstrap_clients.ArtifactFactoryForSettings(
		ctx,
		params.settings,
		params.memCache,
		&cfg,
		params.clientset,
		&kubeCoreCache,
		params.consulClient,
		v1.ArtifactCrd.Plural,
	)
	if err != nil {
		return bootstrap.Opts{}, err
	}

	authConfigFactory, err := bootstrap_clients.ConfigFactoryForSettings(factoryParams, extauth.AuthConfigCrd)
	if err != nil {
		return bootstrap.Opts{}, err
	}

	rateLimitConfigFactory, err := bootstrap_clients.ConfigFactoryForSettings(factoryParams, rlv1alpha1.RateLimitConfigCrd)
	if err != nil {
		return bootstrap.Opts{}, err
	}

	graphqlApiFactory, err := bootstrap_clients.ConfigFactoryForSettings(factoryParams, v1beta1.GraphQLApiCrd)
	if err != nil {
		return bootstrap.Opts{}, err
	}

	virtualServiceFactory, err := bootstrap_clients.ConfigFactoryForSettings(factoryParams, gateway.VirtualServiceCrd)
	if err != nil {
		return bootstrap.Opts{}, err
	}

	routeTableFactory, err := bootstrap_clients.ConfigFactoryForSettings(factoryParams, gateway.RouteTableCrd)
	if err != nil {
		return bootstrap.Opts{}, err
	}

	virtualHostOptionFactory, err := bootstrap_clients.ConfigFactoryForSettings(factoryParams, gateway.VirtualHostOptionCrd)
	if err != nil {
		return bootstrap.Opts{}, err
	}

	routeOptionFactory, err := bootstrap_clients.ConfigFactoryForSettings(factoryParams, gateway.RouteOptionCrd)
	if err != nil {
		return bootstrap.Opts{}, err
	}

	gatewayFactory, err := bootstrap_clients.ConfigFactoryForSettings(factoryParams, gateway.GatewayCrd)
	if err != nil {
		return bootstrap.Opts{}, err
	}

	matchableHttpGatewayFactory, err := bootstrap_clients.ConfigFactoryForSettings(factoryParams, gateway.MatchableHttpGatewayCrd)
	if err != nil {
		return bootstrap.Opts{}, err
	}

	matchableTcpGatewayFactory, err := bootstrap_clients.ConfigFactoryForSettings(factoryParams, gateway.MatchableTcpGatewayCrd)
	if err != nil {
		return bootstrap.Opts{}, err
	}

	var validation *gwtranslator.ValidationOpts
	validationCfg := params.settings.GetGateway().GetValidation()

	validationServerEnabled := validationCfg != nil // default to true if validation top level field is set
	if validationCfg.GetServerEnabled() != nil {
		// allow user to explicitly disable validation server
		validationServerEnabled = validationCfg.GetServerEnabled().GetValue()
	}

	var gatewayMode bool
	if params.settings.GetGateway().GetEnableGatewayController() != nil {
		gatewayMode = params.settings.GetGateway().GetEnableGatewayController().GetValue()
	} else {
		gatewayMode = true
	}
	if validationServerEnabled && gatewayMode {
		alwaysAcceptResources := AcceptAllResourcesByDefault

		if alwaysAccept := validationCfg.GetAlwaysAccept(); alwaysAccept != nil {
			alwaysAcceptResources = alwaysAccept.GetValue()
		}

		allowWarnings := AllowWarnings

		if allowWarning := validationCfg.GetAllowWarnings(); allowWarning != nil {
			allowWarnings = allowWarning.GetValue()
		}

		validation = &gwtranslator.ValidationOpts{
			ProxyValidationServerAddress: validationCfg.GetProxyValidationServerAddr(),
			ValidatingWebhookPort:        gwdefaults.ValidationWebhookBindPort,
			ValidatingWebhookCertPath:    validationCfg.GetValidationWebhookTlsCert(),
			ValidatingWebhookKeyPath:     validationCfg.GetValidationWebhookTlsKey(),
			AlwaysAcceptResources:        alwaysAcceptResources,
			AllowWarnings:                allowWarnings,
			WarnOnRouteShortCircuiting:   validationCfg.GetWarnRouteShortCircuiting().GetValue(),
		}
		if validation.ProxyValidationServerAddress == "" {
			validation.ProxyValidationServerAddress = gwdefaults.GlooProxyValidationServerAddr
		}
		if overrideAddr := os.Getenv("PROXY_VALIDATION_ADDR"); overrideAddr != "" {
			validation.ProxyValidationServerAddress = overrideAddr
		}
		if validation.ValidatingWebhookCertPath == "" {
			validation.ValidatingWebhookCertPath = gwdefaults.ValidationWebhookTlsCertPath
		}
		if validation.ValidatingWebhookKeyPath == "" {
			validation.ValidatingWebhookKeyPath = gwdefaults.ValidationWebhookTlsKeyPath
		}
	} else {
		// This will stop users from setting failurePolicy=fail and then removing the webhook configuration
		if validationMustStart := os.Getenv("VALIDATION_MUST_START"); validationMustStart != "" && validationMustStart != "false" && gatewayMode {
			return bootstrap.Opts{}, errors.Errorf("A validation webhook was configured, but no validation configuration was provided in the settings. "+
				"Ensure the v1.Settings %v contains the spec.gateway.validation config."+
				"If you have removed the webhook configuration from K8s since installing and want to disable validation, "+
				"set the environment variable VALIDATION_MUST_START=false",
				params.settings.GetMetadata().Ref())
		}
	}
	readGatewaysFromAllNamespaces := params.settings.GetGateway().GetReadGatewaysFromAllNamespaces()

	return bootstrap.Opts{
		Upstreams:                    upstreamFactory,
		KubeServiceClient:            kubeServiceClient,
		Proxies:                      proxyFactory,
		UpstreamGroups:               upstreamGroupFactory,
		Secrets:                      secretFactory,
		Artifacts:                    artifactFactory,
		AuthConfigs:                  authConfigFactory,
		RateLimitConfigs:             rateLimitConfigFactory,
		GraphQLApis:                  graphqlApiFactory,
		VirtualServices:              virtualServiceFactory,
		RouteTables:                  routeTableFactory,
		VirtualHostOptions:           virtualHostOptionFactory,
		RouteOptions:                 routeOptionFactory,
		Gateways:                     gatewayFactory,
		MatchableHttpGateways:        matchableHttpGatewayFactory,
		MatchableTcpGateways:         matchableTcpGatewayFactory,
		KubeCoreCache:                kubeCoreCache,
		ValidationOpts:               validation,
		ReadGatwaysFromAllNamespaces: readGatewaysFromAllNamespaces,
		GatewayControllerEnabled:     gatewayMode,
		ProxyCleanup:                 proxyCleanup,
		GlooGateway:                  constructGlooGatewayBootstrapOpts(params.settings),
		ProxyReconcileQueue:          setup.ProxyReconcileQueue,
	}, nil
}

func constructGlooGatewayBootstrapOpts(settings *v1.Settings) bootstrap.GlooGateway {
	return bootstrap.GlooGateway{
		// TODO: This value should be inherited at installation time, to determine if the k8s controller is enabled
		// In the interim, we use an env variable to control the value
		EnableK8sGatewayController: envutils.IsEnvTruthy(constants.GlooGatewayEnableK8sGwControllerEnv),
		IstioValues:                constructIstioBootstrapOpts(settings),
	}
}

func constructIstioBootstrapOpts(settings *v1.Settings) bootstrap.IstioValues {
	istioValues := bootstrap.IstioValues{
		IntegrationEnabled: settings.GetGloo().GetIstioOptions().GetEnableIntegration().GetValue(),

		// TODO: enableIstioSidecarOnGateway should be removed as part of: https://github.com/solo-io/solo-projects/issues/5743
		SidecarOnGatewayEnabled: envutils.IsEnvTruthy(constants.IstioInjectionEnabled),
	}

	return istioValues
}

func runQueue(ctx context.Context, proxyReconcileQueue ggv2utils.AsyncQueue[gloov1.ProxyList], writeNamespace string, proxyClient gloov1.ProxyClient) {
	// labels used to uniquely identify Proxies that are managed by the kube gateway controller
	var kubeGatewayProxyLabels = map[string]string{
		// the proxy type key/value must stay in sync with the one defined in projects/gateway2/translator/gateway_translator.go
		utils.ProxyTypeKey: utils.GatewayApiProxyValue,
	}
	ctx = contextutils.WithLogger(ctx, "proxyCache")
	logger := contextutils.LoggerFrom(ctx)

	proxyReconciler := gloov1.NewProxyReconciler(proxyClient, statusutils.NewNoOpStatusClient())
	for {
		proxyList, err := proxyReconcileQueue.Dequeue(ctx)
		if err != nil {
			return
		}
		// Proxy CR is located in the writeNamespace, which may be different from the originating Gateway CR
		err = proxyReconciler.Reconcile(
			writeNamespace,
			proxyList,
			func(original, desired *gloov1.Proxy) (bool, error) {
				// only reconcile if proxies are not equal
				// we reconcile so ggv2 proxies can be used in extension syncing and debug snap storage
				return !proto.Equal(original, desired), nil
			},
			clients.ListOpts{
				Ctx:      ctx,
				Selector: kubeGatewayProxyLabels,
			})
		if err != nil {
			// A write error to our cache should not impact translation
			// We will emit a message, and continue
			logger.Error(err)
		}

	}

}

func multiCallbacks(cb ...xdsserver.Callbacks) xdsserver.Callbacks {
	cb = slices.DeleteFunc(cb, func(x xdsserver.Callbacks) bool {
		return x == nil
	})
	if len(cb) == 0 {
		return nil
	}
	if len(cb) == 1 {
		return cb[0]
	}
	return MutltiCallbacks(cb)
}
