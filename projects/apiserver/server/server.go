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
	rpc_fed_v1 "github.com/solo-io/solo-projects/projects/apiserver/pkg/api/fed.rpc/v1"
	rpc_edge_v1 "github.com/solo-io/solo-projects/projects/apiserver/pkg/api/rpc.edge.gloo/v1"
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
	bootstrapService rpc_edge_v1.BootstrapApiServer,
	glooInstanceService rpc_edge_v1.GlooInstanceApiServer,
	failoverSchemeService rpc_fed_v1.FailoverSchemeApiServer,
	routeTableSelectorService rpc_edge_v1.VirtualServiceRoutesApiServer,
	wasmFilterService rpc_edge_v1.WasmFilterApiServer,
	glooResourceService rpc_edge_v1.GlooResourceApiServer,
	gatewayResourceService rpc_edge_v1.GatewayResourceApiServer,
	glooEnterpriseResourceService rpc_edge_v1.EnterpriseGlooResourceApiServer,
	ratelimitResourceService rpc_edge_v1.RatelimitResourceApiServer,
	glooFedResourceService rpc_fed_v1.FederatedGlooResourceApiServer,
	gatewayFedResourceService rpc_fed_v1.FederatedGatewayResourceApiServer,
	glooEnterpriseFedResourceService rpc_fed_v1.FederatedEnterpriseGlooResourceApiServer,
	ratelimitFedResourceService rpc_fed_v1.FederatedRatelimitResourceApiServer,
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
	rpc_edge_v1.RegisterBootstrapApiServer(server, bootstrapService)
	rpc_edge_v1.RegisterGlooInstanceApiServer(server, glooInstanceService)
	rpc_fed_v1.RegisterFailoverSchemeApiServer(server, failoverSchemeService)
	rpc_edge_v1.RegisterVirtualServiceRoutesApiServer(server, routeTableSelectorService)
	rpc_edge_v1.RegisterWasmFilterApiServer(server, wasmFilterService)
	rpc_edge_v1.RegisterGlooResourceApiServer(server, glooResourceService)
	rpc_edge_v1.RegisterGatewayResourceApiServer(server, gatewayResourceService)
	rpc_edge_v1.RegisterEnterpriseGlooResourceApiServer(server, glooEnterpriseResourceService)
	rpc_edge_v1.RegisterRatelimitResourceApiServer(server, ratelimitResourceService)
	rpc_fed_v1.RegisterFederatedGlooResourceApiServer(server, glooFedResourceService)
	rpc_fed_v1.RegisterFederatedGatewayResourceApiServer(server, gatewayFedResourceService)
	rpc_fed_v1.RegisterFederatedEnterpriseGlooResourceApiServer(server, glooEnterpriseFedResourceService)
	rpc_fed_v1.RegisterFederatedRatelimitResourceApiServer(server, ratelimitFedResourceService)

	return &grpcServer{
		healthChecker: healthChecker,
		server:        server,
	}
}

func NewSingleClusterGlooGrpcServer(
	ctx context.Context,
	bootstrapService rpc_edge_v1.BootstrapApiServer,
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
	rpc_edge_v1.RegisterBootstrapApiServer(server, bootstrapService)

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
