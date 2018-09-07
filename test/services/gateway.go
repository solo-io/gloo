package services

import (
	"net"
	"time"

	"github.com/solo-io/solo-kit/projects/gateway/pkg/setup"

	"context"
	"sync/atomic"

	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/factory"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/memory"
	"github.com/solo-io/solo-kit/pkg/namespacing/static"
	"github.com/solo-io/solo-kit/pkg/utils/contextutils"

	gatewayv1 "github.com/solo-io/solo-kit/projects/gateway/pkg/api/v1"
	gloov1 "github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1"
	gloosetup "github.com/solo-io/solo-kit/projects/gloo/pkg/setup"

	"github.com/solo-io/solo-kit/projects/gloo/pkg/bootstrap"
	"google.golang.org/grpc"

	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"go.uber.org/zap"

	. "github.com/onsi/gomega"
)

type TestClients struct {
	GatewayClient        gatewayv1.GatewayClient
	VirtualServiceClient gatewayv1.VirtualServiceClient
	ProxyClient          gloov1.ProxyClient
	UpstreamClient       gloov1.UpstreamClient
	SecretClient         gloov1.SecretClient
	GlooPort             int
}

var glooPort int32 = 8100

func RunGateway(ctx context.Context, justgloo bool) TestClients {
	localglooPort := atomic.AddInt32(&glooPort, 1)

	cache := memory.NewInMemoryResourceCache()

	glooopts := DefaultGlooOpts(ctx, cache)
	glooopts.BindAddr.(*net.TCPAddr).Port = int(localglooPort)
	// no gateway for now
	if !justgloo {
		opts := DefaultTestConstructOpts(ctx, cache)
		go setup.Setup(opts)
	}
	go gloosetup.Setup(glooopts)

	// construct our own resources:
	a := &factory.MemoryResourceClientOpts{
		Cache: cache,
	}
	factory := factory.NewResourceClientFactory(a)

	gatewayClient, err := gatewayv1.NewGatewayClient(factory)
	Expect(err).NotTo(HaveOccurred())
	virtualServiceClient, err := gatewayv1.NewVirtualServiceClient(factory)
	Expect(err).NotTo(HaveOccurred())
	upstreamClient, err := gloov1.NewUpstreamClient(factory)
	Expect(err).NotTo(HaveOccurred())
	secretClient, err := gloov1.NewSecretClient(factory)
	Expect(err).NotTo(HaveOccurred())
	proxyClient, err := gloov1.NewProxyClient(factory)
	Expect(err).NotTo(HaveOccurred())

	return TestClients{
		GatewayClient:        gatewayClient,
		VirtualServiceClient: virtualServiceClient,
		UpstreamClient:       upstreamClient,
		SecretClient:         secretClient,
		ProxyClient:          proxyClient,
		GlooPort:             int(localglooPort),
	}
}

func DefaultTestConstructOpts(ctx context.Context, cache memory.InMemoryResourceCache) setup.Opts {
	ctx = contextutils.WithLogger(ctx, "gateway")
	ctx = contextutils.SilenceLogger(ctx)
	f := &factory.MemoryResourceClientOpts{
		Cache: cache,
	}
	return setup.NewOpts(
		"gloo-system",
		f,
		f,
		f,
		f,
		f,
		static.NewNamespacer([]string{"default", "gloo-system"}),
		clients.WatchOpts{
			Ctx:         ctx,
			RefreshRate: time.Minute,
		},
		false,
	)
}

func DefaultGlooOpts(ctx context.Context, cache memory.InMemoryResourceCache) bootstrap.Opts {
	ctx = contextutils.WithLogger(ctx, "gloo")
	logger := contextutils.LoggerFrom(ctx)
	grpcServer := grpc.NewServer(grpc.StreamInterceptor(
		grpc_middleware.ChainStreamServer(
			grpc_ctxtags.StreamServerInterceptor(),
			grpc_zap.StreamServerInterceptor(zap.NewNop()),
			func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
				logger.Debugf("gRPC call: %v", info.FullMethod)
				return handler(srv, ss)
			},
		)),
	)
	f := &factory.MemoryResourceClientOpts{
		Cache: cache,
	}
	return bootstrap.Opts{
		WriteNamespace: "gloo-system",
		Upstreams:      f,
		Proxies:        f,
		Secrets:        f,
		Artifacts:      f,
		Namespacer:     static.NewNamespacer([]string{"default", "gloo-system"}),
		WatchOpts: clients.WatchOpts{
			Ctx:         ctx,
			RefreshRate: time.Second / 10,
		},
		BindAddr: &net.TCPAddr{
			IP:   net.ParseIP("0.0.0.0"),
			Port: 8081,
		},
		GrpcServer: grpcServer,
	}
}
