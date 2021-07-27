package server

import (
	"context"
	"fmt"
	"net"
	"net/http"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"github.com/rotisserie/eris"
	"github.com/solo-io/go-utils/contextutils"
	"github.com/solo-io/solo-projects/projects/apiserver/internal/settings"
	fed_rpc_v1 "github.com/solo-io/solo-projects/projects/apiserver/pkg/api/fed.rpc/v1"
	edge_rpc_v1 "github.com/solo-io/solo-projects/projects/apiserver/pkg/api/rpc.edge.gloo/v1"
	"github.com/solo-io/solo-projects/projects/apiserver/server/health_check"
	"go.opencensus.io/plugin/ocgrpc"
	"go.opencensus.io/stats/view"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

func init() {
	view.Register(ocgrpc.DefaultServerViews...)
}

type GrpcServer interface {
	Run(ctx context.Context, cfg *settings.ApiServerSettings) error
	Stop()
}

type grpcServer struct {
	server        *grpc.Server
	healthChecker health_check.HealthChecker
}

func NewGlooFedGrpcServer(
	ctx context.Context,
	bootstrapService edge_rpc_v1.BootstrapApiServer,
	glooInstanceService fed_rpc_v1.GlooInstanceApiServer,
	failoverSchemeService fed_rpc_v1.FailoverSchemeApiServer,
	routeTableSelectorService fed_rpc_v1.VirtualServiceRoutesApiServer,
	wasmFilterService fed_rpc_v1.WasmFilterApiServer,
	glooResourceService fed_rpc_v1.GlooResourceApiServer,
	gatewayResourceService fed_rpc_v1.GatewayResourceApiServer,
	glooEnterpriseResourceService fed_rpc_v1.EnterpriseGlooResourceApiServer,
	ratelimitResourceService fed_rpc_v1.RatelimitResourceApiServer,
	glooFedResourceService fed_rpc_v1.FederatedGlooResourceApiServer,
	gatewayFedResourceService fed_rpc_v1.FederatedGatewayResourceApiServer,
	glooEnterpriseFedResourceService fed_rpc_v1.FederatedEnterpriseGlooResourceApiServer,
	ratelimitFedResourceService fed_rpc_v1.FederatedRatelimitResourceApiServer,
	healthChecker health_check.HealthChecker,
) GrpcServer {

	logger := contextutils.LoggerFrom(ctx)
	server := grpc.NewServer(
		grpc.StatsHandler(&ocgrpc.ServerHandler{}),
		grpc_middleware.WithUnaryServerChain(grpc_zap.UnaryServerInterceptor(logger.Desugar())),
	)
	// register grpc health check
	healthpb.RegisterHealthServer(server, healthChecker)
	// register handlers
	edge_rpc_v1.RegisterBootstrapApiServer(server, bootstrapService)
	fed_rpc_v1.RegisterGlooInstanceApiServer(server, glooInstanceService)
	fed_rpc_v1.RegisterFailoverSchemeApiServer(server, failoverSchemeService)
	fed_rpc_v1.RegisterVirtualServiceRoutesApiServer(server, routeTableSelectorService)
	fed_rpc_v1.RegisterWasmFilterApiServer(server, wasmFilterService)
	fed_rpc_v1.RegisterGlooResourceApiServer(server, glooResourceService)
	fed_rpc_v1.RegisterGatewayResourceApiServer(server, gatewayResourceService)
	fed_rpc_v1.RegisterEnterpriseGlooResourceApiServer(server, glooEnterpriseResourceService)
	fed_rpc_v1.RegisterRatelimitResourceApiServer(server, ratelimitResourceService)
	fed_rpc_v1.RegisterFederatedGlooResourceApiServer(server, glooFedResourceService)
	fed_rpc_v1.RegisterFederatedGatewayResourceApiServer(server, gatewayFedResourceService)
	fed_rpc_v1.RegisterFederatedEnterpriseGlooResourceApiServer(server, glooEnterpriseFedResourceService)
	fed_rpc_v1.RegisterFederatedRatelimitResourceApiServer(server, ratelimitFedResourceService)

	return &grpcServer{
		healthChecker: healthChecker,
		server:        server,
	}
}

func NewGlooInstanceGrpcServer(
	ctx context.Context,
	bootstrapService edge_rpc_v1.BootstrapApiServer,
	healthChecker health_check.HealthChecker,
) GrpcServer {
	logger := contextutils.LoggerFrom(ctx)
	server := grpc.NewServer(
		grpc.StatsHandler(&ocgrpc.ServerHandler{}),
		grpc_middleware.WithUnaryServerChain(grpc_zap.UnaryServerInterceptor(logger.Desugar())),
	)
	// register grpc health check
	healthpb.RegisterHealthServer(server, healthChecker)

	// register handlers
	edge_rpc_v1.RegisterBootstrapApiServer(server, bootstrapService)

	return &grpcServer{
		healthChecker: healthChecker,
		server:        server,
	}
}

func (g *grpcServer) Run(ctx context.Context, cfg *settings.ApiServerSettings) error {
	var eg errgroup.Group
	// Start http health check
	healthCheckListener, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.HealthCheckPort))
	if err != nil {
		return eris.Wrapf(err, "failed to setup health check listener")
	}
	contextutils.LoggerFrom(ctx).Infof("Set up health check listener at port %d", cfg.HealthCheckPort)
	eg.Go(func() error {
		return http.Serve(healthCheckListener, g.healthChecker)
	})
	// start grpc listener
	grpcListener, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GrpcPort))
	if err != nil {
		return eris.Wrapf(err, "failed to setup grpc listener")
	}
	contextutils.LoggerFrom(ctx).Infof("Set up grpc listener at port %d", cfg.GrpcPort)
	eg.Go(func() error {
		return g.server.Serve(grpcListener)
	})
	return eg.Wait()
}

func (g *grpcServer) Stop() {
	g.server.GracefulStop()
}
