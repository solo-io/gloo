package apiserver

import (
	"context"

	"github.com/solo-io/go-utils/contextutils"
	"github.com/solo-io/solo-projects/projects/apiserver/internal/settings"
	rpc_fed_v1 "github.com/solo-io/solo-projects/projects/apiserver/pkg/api/fed.rpc/v1"
	rpc_edge_v1 "github.com/solo-io/solo-projects/projects/apiserver/pkg/api/rpc.edge.gloo/v1"
	"github.com/solo-io/solo-projects/projects/apiserver/server"
	"github.com/solo-io/solo-projects/projects/apiserver/server/health_check"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

func NewGlooFedServerRunnable(
	rootCtx context.Context,
	cfg *settings.ApiServerSettings,
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
) manager.Runnable {
	return manager.RunnableFunc(func(ctx context.Context) error {
		ctx = contextutils.WithLogger(rootCtx, "gloo-fed-apiserver")
		apiServer := server.NewGlooFedGrpcServer(ctx, bootstrapService, glooInstanceService, failoverSchemeService, routeTableSelectorService, wasmFilterService,
			glooResourceService, gatewayResourceService, glooEnterpriseResourceService, ratelimitResourceService, glooFedResourceService,
			gatewayFedResourceService, glooEnterpriseFedResourceService, ratelimitFedResourceService, health_check.NewHealthChecker())

		return apiServer.Run(ctx, cfg)
	})
}

func NewSingleClusterGlooServerRunnable(
	rootCtx context.Context,
	cfg *settings.ApiServerSettings,
	bootstrapService rpc_edge_v1.BootstrapApiServer,
	glooInstanceService rpc_edge_v1.GlooInstanceApiServer,
	gatewayResourceService rpc_edge_v1.GatewayResourceApiServer,
	glooResourceService rpc_edge_v1.GlooResourceApiServer,
	glooEnterpriseResourceService rpc_edge_v1.EnterpriseGlooResourceApiServer,
	ratelimitResourceService rpc_edge_v1.RatelimitResourceApiServer,
) manager.Runnable {
	return manager.RunnableFunc(func(ctx context.Context) error {
		ctx = contextutils.WithLogger(rootCtx, "gloo-ee-apiserver")
		apiServer := server.NewSingleClusterGlooGrpcServer(ctx, bootstrapService, glooInstanceService,
			gatewayResourceService, glooResourceService, glooEnterpriseResourceService, ratelimitResourceService,
			health_check.NewHealthChecker())

		return apiServer.Run(ctx, cfg)
	})
}
