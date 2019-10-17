package services

import (
	"net"
	"time"

	extauthv1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/plugins/extauth/v1"

	"github.com/solo-io/solo-projects/projects/gloo/pkg/setup"

	"github.com/solo-io/gloo/pkg/utils/settingsutil"
	gatewaysyncer "github.com/solo-io/gloo/projects/gateway/pkg/syncer"

	"context"
	"sync/atomic"

	"github.com/solo-io/go-utils/contextutils"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/factory"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/memory"

	gatewayv1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
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
	"github.com/solo-io/gloo/projects/gloo/pkg/syncer"

	"k8s.io/client-go/kubernetes"
)

type TestClients struct {
	GatewayClient        gatewayv1.GatewayClient
	VirtualServiceClient gatewayv1.VirtualServiceClient
	ProxyClient          gloov1.ProxyClient
	UpstreamClient       gloov1.UpstreamClient
	SecretClient         gloov1.SecretClient
	ArtifactClient       gloov1.ArtifactClient
	ServiceClient        skkube.ServiceClient
	AuthConfigClient     extauthv1.AuthConfigClient
	GlooPort             int
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

	testclients := GetTestClients(cache)
	testclients.GlooPort = RunGlooGatewayUdsFds(ctx, cache, What{DisableGateway: justgloo}, ns, kubeclient, extensions)
	return testclients
}

func GetTestClients(cache memory.InMemoryResourceCache) TestClients {

	// construct our own resources:
	rcFactory := &factory.MemoryResourceClientFactory{
		Cache: cache,
	}

	gatewayClient, err := gatewayv1.NewGatewayClient(rcFactory)
	Expect(err).NotTo(HaveOccurred())
	virtualServiceClient, err := gatewayv1.NewVirtualServiceClient(rcFactory)
	Expect(err).NotTo(HaveOccurred())
	upstreamClient, err := gloov1.NewUpstreamClient(rcFactory)
	Expect(err).NotTo(HaveOccurred())
	secretClient, err := gloov1.NewSecretClient(rcFactory)
	Expect(err).NotTo(HaveOccurred())
	artifactClient, err := gloov1.NewArtifactClient(rcFactory)
	Expect(err).NotTo(HaveOccurred())
	proxyClient, err := gloov1.NewProxyClient(rcFactory)
	Expect(err).NotTo(HaveOccurred())
	authConfigClient, err := extauthv1.NewAuthConfigClient(rcFactory)
	Expect(err).NotTo(HaveOccurred())

	return TestClients{
		GatewayClient:        gatewayClient,
		VirtualServiceClient: virtualServiceClient,
		UpstreamClient:       upstreamClient,
		SecretClient:         secretClient,
		ArtifactClient:       artifactClient,
		ProxyClient:          proxyClient,
		AuthConfigClient:     authConfigClient,
	}
}

type What struct {
	DisableGateway bool
	DisableUds     bool
	DisableFds     bool
}

func RunGlooGatewayUdsFds(ctx context.Context, cache memory.InMemoryResourceCache, what What, ns string, kubeclient kubernetes.Interface, extensions *v1.Extensions) int {
	port := AllocateGlooPort()
	RunGlooGatewayUdsFdsOnPort(ctx, cache, port, what, ns, kubeclient, extensions, nil)
	return int(port)
}

func AllocateGlooPort() int32 {
	return atomic.AddInt32(&glooPort, 1)

}

func RunGlooGatewayUdsFdsOnPort(ctx context.Context, cache memory.InMemoryResourceCache, localglooPort int32, what What,
	ns string, kubeclient kubernetes.Interface, extensions *v1.Extensions, s *gloov1.Settings) {

	// no gateway for now
	opts := DefaultTestConstructOpts(ctx, cache, ns)
	if !what.DisableGateway {
		go gatewaysyncer.RunGateway(opts)
	}
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
	go syncer.RunGlooWithExtensions(glooOpts, setup.GetGlooEeExtensions())
	if !what.DisableFds {
		go fds_syncer.RunFDS(glooOpts)
	}
	if !what.DisableUds {
		go uds_syncer.RunUDS(glooOpts)
	}

}

func DefaultTestConstructOpts(ctx context.Context, cache memory.InMemoryResourceCache, ns string) gatewaysyncer.Opts {
	ctx = contextutils.WithLogger(ctx, "gateway")
	ctx = contextutils.SilenceLogger(ctx)
	f := &factory.MemoryResourceClientFactory{
		Cache: cache,
	}

	return gatewaysyncer.Opts{
		WriteNamespace:  ns,
		WatchNamespaces: []string{"default", ns},
		Gateways:        f,
		VirtualServices: f,
		Proxies:         f,
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
		WriteNamespace:  ns,
		Upstreams:       f,
		UpstreamGroups:  f,
		Proxies:         f,
		Secrets:         f,
		Artifacts:       f,
		AuthConfigs:     f,
		WatchNamespaces: []string{"default", ns},
		WatchOpts: clients.WatchOpts{
			Ctx:         ctx,
			RefreshRate: time.Second / 10,
		},
		ControlPlane: syncer.NewControlPlane(ctx, grpcServer, &net.TCPAddr{
			IP:   net.ParseIP("0.0.0.0"),
			Port: 8081,
		}, nil, true),
		ValidationServer: syncer.NewValidationServer(ctx, grpcServerValidation, &net.TCPAddr{
			IP:   net.ParseIP("0.0.0.0"),
			Port: 8081,
		}, true),
		KubeClient: kubeclient,
		DevMode:    true,
	}
}
