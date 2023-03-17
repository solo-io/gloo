package services

import (
	"net"
	"time"

	"github.com/solo-io/gloo/pkg/bootstrap/leaderelector/singlereplica"

	license2 "github.com/solo-io/solo-projects/pkg/license"

	"github.com/solo-io/licensing/pkg/model"

	"github.com/solo-io/solo-projects/projects/discovery/pkg/fds/syncer"

	graphqlv1beta1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/graphql/v1beta1"

	"github.com/solo-io/gloo/pkg/utils/statusutils"

	v1alpha1 "github.com/solo-io/gloo/projects/gloo/pkg/api/external/solo/ratelimit"

	extauthv1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/extauth/v1"

	"github.com/solo-io/solo-projects/projects/gloo/pkg/setup"

	"context"
	"sync/atomic"

	"github.com/solo-io/gloo/pkg/utils/settingsutil"

	"github.com/solo-io/go-utils/contextutils"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/factory"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/memory"

	gatewayv1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gateway/pkg/translator"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/bootstrap"
	skkube "github.com/solo-io/solo-kit/pkg/api/v1/resources/common/kubernetes"
	"google.golang.org/grpc"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"go.uber.org/zap"

	. "github.com/onsi/gomega"
	fds_syncer "github.com/solo-io/gloo/projects/discovery/pkg/fds/syncer"
	uds_syncer "github.com/solo-io/gloo/projects/discovery/pkg/uds/syncer"
	"github.com/solo-io/gloo/projects/gloo/pkg/defaults"
	syncer_setup "github.com/solo-io/gloo/projects/gloo/pkg/syncer/setup"
	"k8s.io/client-go/kubernetes"
)

type TestClients struct {
	GatewayClient         gatewayv1.GatewayClient
	VirtualServiceClient  gatewayv1.VirtualServiceClient
	ProxyClient           gloov1.ProxyClient
	UpstreamClient        gloov1.UpstreamClient
	SecretClient          gloov1.SecretClient
	ArtifactClient        gloov1.ArtifactClient
	ServiceClient         skkube.ServiceClient
	AuthConfigClient      extauthv1.AuthConfigClient
	RateLimitConfigClient v1alpha1.RateLimitConfigClient
	GraphQLApiClient      graphqlv1beta1.GraphQLApiClient
	GlooPort              int
}

var glooPort int32 = 8100

func RunGateway(ctx context.Context, justgloo bool) TestClients {
	return RunGatewayWithNamespaceAndKubeClient(ctx, justgloo, defaults.GlooSystem, nil)
}

func RunGatewayWithNamespaceAndKubeClient(ctx context.Context, justgloo bool, ns string, kubeclient kubernetes.Interface) TestClients {
	return RunGatewayWithKubeClientAndSettings(ctx, justgloo, ns, kubeclient, nil)
}

func RunGatewayWithSettings(ctx context.Context, justgloo bool, extensions *v1.Extensions) TestClients {
	return RunGatewayWithKubeClientAndSettings(ctx, justgloo, defaults.GlooSystem, nil, extensions)
}

func RunGatewayWithKubeClientAndSettings(ctx context.Context, justgloo bool, ns string, kubeclient kubernetes.Interface, extensions *v1.Extensions) TestClients {
	cache := memory.NewInMemoryResourceCache()

	testclients := GetTestClients(ctx, cache)
	testclients.GlooPort = RunGlooGatewayUdsFds(RunGlooGatewayOpts{
		Ctx:        ctx,
		Cache:      cache,
		What:       What{DisableGateway: justgloo},
		Namespace:  ns,
		KubeClient: kubeclient,
		Extensions: extensions,
	})
	return testclients
}

func GetTestClients(ctx context.Context, cache memory.InMemoryResourceCache) TestClients {

	// construct our own resources:
	rcFactory := &factory.MemoryResourceClientFactory{
		Cache: cache,
	}

	gatewayClient, err := gatewayv1.NewGatewayClient(ctx, rcFactory)
	Expect(err).NotTo(HaveOccurred())
	virtualServiceClient, err := gatewayv1.NewVirtualServiceClient(ctx, rcFactory)
	Expect(err).NotTo(HaveOccurred())
	upstreamClient, err := gloov1.NewUpstreamClient(ctx, rcFactory)
	Expect(err).NotTo(HaveOccurred())
	secretClient, err := gloov1.NewSecretClient(ctx, rcFactory)
	Expect(err).NotTo(HaveOccurred())
	artifactClient, err := gloov1.NewArtifactClient(ctx, rcFactory)
	Expect(err).NotTo(HaveOccurred())
	proxyClient, err := gloov1.NewProxyClient(ctx, rcFactory)
	Expect(err).NotTo(HaveOccurred())
	authConfigClient, err := extauthv1.NewAuthConfigClient(ctx, rcFactory)
	Expect(err).NotTo(HaveOccurred())
	rlcClient, err := v1alpha1.NewRateLimitConfigClient(ctx, rcFactory)
	Expect(err).NotTo(HaveOccurred())
	gqlClient, err := graphqlv1beta1.NewGraphQLApiClient(ctx, rcFactory)
	Expect(err).NotTo(HaveOccurred())

	return TestClients{
		GatewayClient:         gatewayClient,
		VirtualServiceClient:  virtualServiceClient,
		UpstreamClient:        upstreamClient,
		SecretClient:          secretClient,
		ArtifactClient:        artifactClient,
		ProxyClient:           proxyClient,
		AuthConfigClient:      authConfigClient,
		RateLimitConfigClient: rlcClient,
		GraphQLApiClient:      gqlClient,
	}
}

type What struct {
	DisableGateway bool
	DisableUds     bool
	DisableFds     bool
}

func RunGlooGatewayUdsFds(opts RunGlooGatewayOpts) int {
	port := AllocateGlooPort()
	opts.LocalGlooPort = port
	RunGlooGatewayUdsFdsOnPort(opts)
	return int(port)
}

func AllocateGlooPort() int32 {
	return atomic.AddInt32(&glooPort, 1)

}

type RunGlooGatewayOpts struct {
	Ctx           context.Context
	Cache         memory.InMemoryResourceCache
	LocalGlooPort int32
	What          What
	Namespace     string
	KubeClient    kubernetes.Interface
	Extensions    *v1.Extensions
	Settings      *gloov1.Settings
	License       *model.License
}

func RunGlooGatewayUdsFdsOnPort(runOpts RunGlooGatewayOpts) {
	ctx, cache, ns, what, s, extensions, kubeclient, localglooPort, licenseState := runOpts.Ctx, runOpts.Cache, runOpts.Namespace, runOpts.What, runOpts.Settings, runOpts.Extensions, runOpts.KubeClient, runOpts.LocalGlooPort, runOpts.License
	// no gateway for now
	opts := DefaultTestConstructOpts(ctx, cache, ns)
	settings := v1.Settings{}
	if s != nil {
		settings = *s
	}
	settings.Extensions = extensions
	settings.WatchNamespaces = opts.WatchNamespaces
	settings.DiscoveryNamespace = opts.WriteNamespace

	ctx = settingsutil.WithSettings(ctx, &settings)

	glooOpts := defaultGlooOpts(ctx, cache, ns, kubeclient)
	glooOpts.ControlPlane.BindAddr.(*net.TCPAddr).Port = int(localglooPort)
	glooOpts.Settings = &settings
	glooOpts.ControlPlane.StartGrpcServer = true
	glooOpts.GatewayControllerEnabled = !what.DisableGateway
	apiEmitterChan := make(chan struct{})
	// For testing purposes, load the LicensedFeatureProvider with the injected license
	licensedFeatureProvider := setupLicensedFeatureProvider(licenseState)

	// Initialize the Gloo Enterprise setup extensions
	setupExtensions := setup.GetGlooEExtensions(ctx, licensedFeatureProvider, apiEmitterChan)
	setupExtensions.PluginRegistryFactory = setup.GetPluginRegistryFactory(glooOpts, apiEmitterChan, licensedFeatureProvider)
	go syncer_setup.RunGlooWithExtensions(glooOpts, setupExtensions)
	if !what.DisableFds {
		go fds_syncer.RunFDSWithExtensions(glooOpts, syncer.GetFDSEnterpriseExtensions())
	}
	if !what.DisableUds {
		go uds_syncer.RunUDS(glooOpts)
	}

}

func setupLicensedFeatureProvider(licenseState *model.License) *license2.LicensedFeatureProvider {
	licensedFeatureProvider := license2.NewLicensedFeatureProvider()

	if licenseState == nil {
		// Default to a valid, enterprise license
		licenseState = &model.License{
			IssuedAt:      time.Now(),
			ExpiresAt:     time.Now(),
			RandomPayload: "",
			LicenseType:   model.LicenseType_Enterprise,
			Product:       model.Product_Gloo,
			AddOns:        nil,
		}
	}

	licensedFeatureProvider.SetValidatedLicense(&license2.ValidatedLicense{
		License: licenseState,
		Warn:    nil,
		Err:     nil,
	})
	return licensedFeatureProvider
}

func DefaultTestConstructOpts(ctx context.Context, cache memory.InMemoryResourceCache, ns string) translator.Opts {
	ctx = contextutils.WithLogger(ctx, "gateway")
	ctx = contextutils.SilenceLogger(ctx)
	f := &factory.MemoryResourceClientFactory{
		Cache: cache,
	}

	return translator.Opts{
		WriteNamespace:          ns,
		StatusReporterNamespace: statusutils.GetStatusReporterNamespaceOrDefault(defaults.GlooSystem),
		WatchNamespaces:         []string{"default", ns},
		Gateways:                f,
		MatchableHttpGateways:   f,
		VirtualServices:         f,
		Proxies:                 f,
		RouteTables:             f,
		RouteOptions:            f,
		VirtualHostOptions:      f,
		WatchOpts: clients.WatchOpts{
			Ctx:         ctx,
			RefreshRate: time.Minute,
		},
		DevMode: false,
	}
}

func defaultGlooOpts(ctx context.Context, cache memory.InMemoryResourceCache, ns string, kubeclient kubernetes.Interface) bootstrap.Opts {
	ctx = contextutils.WithLogger(ctx, "gloo")
	logger := contextutils.LoggerFrom(ctx)
	grpcServer := grpc.NewServer(grpc.StreamInterceptor(
		grpc_middleware.ChainStreamServer(
			grpc_ctxtags.StreamServerInterceptor(),
			grpc_zap.StreamServerInterceptor(zap.NewNop()),
			func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
				logger.Infof("gRPC call: %v", info.FullMethod)
				return handler(srv, ss)
			},
		)),
	)
	grpcServerValidation := grpc.NewServer(grpc.StreamInterceptor(
		grpc_middleware.ChainStreamServer(
			grpc_ctxtags.StreamServerInterceptor(),
			grpc_zap.StreamServerInterceptor(zap.NewNop()),
			func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
				logger.Infof("gRPC call: %v", info.FullMethod)
				return handler(srv, ss)
			},
		)),
	)
	f := &factory.MemoryResourceClientFactory{
		Cache: cache,
	}
	return bootstrap.Opts{
		WriteNamespace:          ns,
		StatusReporterNamespace: statusutils.GetStatusReporterNamespaceOrDefault(defaults.GlooSystem),
		Upstreams:               f,
		UpstreamGroups:          f,
		Proxies:                 f,
		Secrets:                 f,
		Artifacts:               f,
		AuthConfigs:             f,
		RateLimitConfigs:        f,
		GraphQLApis:             f,
		VirtualServices:         f,
		Gateways:                f,
		MatchableHttpGateways:   f,
		RouteTables:             f,
		RouteOptions:            f,
		VirtualHostOptions:      f,
		WatchNamespaces:         []string{"default", ns},
		WatchOpts: clients.WatchOpts{
			Ctx:         ctx,
			RefreshRate: time.Second / 10,
		},
		ControlPlane: syncer_setup.NewControlPlane(ctx, grpcServer, &net.TCPAddr{
			IP:   net.IPv4zero,
			Port: 8081,
		}, nil, true),
		ValidationServer: syncer_setup.NewValidationServer(ctx, grpcServerValidation, &net.TCPAddr{
			IP:   net.IPv4zero,
			Port: 8081,
		}, true),
		ProxyDebugServer: syncer_setup.NewProxyDebugServer(ctx, grpcServerValidation, &net.TCPAddr{
			IP:   net.IPv4zero,
			Port: 8081,
		}, false),
		KubeClient: kubeclient,
		DevMode:    true,
		Identity:   singlereplica.Identity(),
	}
}
