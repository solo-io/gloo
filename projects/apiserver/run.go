package apiserver

import (
	"context"

	"github.com/solo-io/go-utils/contextutils"
	"github.com/solo-io/solo-projects/projects/apiserver/internal/settings"
	fed_rpc_v1 "github.com/solo-io/solo-projects/projects/apiserver/pkg/api/fed.rpc/v1"
	edge_rpc_v1 "github.com/solo-io/solo-projects/projects/apiserver/pkg/api/rpc.edge.gloo/v1"
	"github.com/solo-io/solo-projects/projects/apiserver/server"
	"github.com/solo-io/solo-projects/projects/apiserver/server/health_check"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

func NewGlooFedServerRunnable(
	rootCtx context.Context,
	cfg *settings.ApiServerSettings,
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
) manager.Runnable {
	return manager.RunnableFunc(func(ctx context.Context) error {
		ctx = contextutils.WithLogger(rootCtx, "gloo-fed-apiserver")
		apiServer := server.NewGlooFedGrpcServer(ctx, bootstrapService, glooInstanceService, failoverSchemeService, routeTableSelectorService, wasmFilterService,
			glooResourceService, gatewayResourceService, glooEnterpriseResourceService, ratelimitResourceService, glooFedResourceService,
			gatewayFedResourceService, glooEnterpriseFedResourceService, ratelimitFedResourceService, health_check.NewHealthChecker())

		return apiServer.Run(ctx, cfg)
	})
}

func NewGlooInstanceServerRunnable(
	rootCtx context.Context,
	cfg *settings.ApiServerSettings,
	bootstrapService edge_rpc_v1.BootstrapApiServer,
) manager.Runnable {
	return manager.RunnableFunc(func(ctx context.Context) error {
		ctx = contextutils.WithLogger(rootCtx, "gloo-ee-apiserver")
		apiServer := server.NewGlooInstanceGrpcServer(ctx, bootstrapService, health_check.NewHealthChecker())

		return apiServer.Run(ctx, cfg)
	})
}
