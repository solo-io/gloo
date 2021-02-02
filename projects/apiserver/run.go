package apiserver

import (
	"context"

	rpc_v1 "github.com/solo-io/solo-projects/projects/apiserver/pkg/api/fed.rpc/v1"

	"github.com/solo-io/solo-projects/projects/apiserver/internal/settings"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	"github.com/solo-io/go-utils/contextutils"
	"github.com/solo-io/solo-projects/projects/apiserver/server"
	"github.com/solo-io/solo-projects/projects/apiserver/server/health_check"
)

func NewServerRunnable(
	rootCtx context.Context, cfg *settings.ApiServerSettings,
	glooInstanceService rpc_v1.GlooInstanceApiServer,
	failoverSchemeService rpc_v1.FailoverSchemeApiServer,
	routeTableSelectorService rpc_v1.VirtualServiceRoutesApiServer,
	wasmFilterService rpc_v1.WasmFilterApiServer,
	glooResourceService rpc_v1.GlooResourceApiServer,
	gatewayResourceService rpc_v1.GatewayResourceApiServer,
	glooEnterpriseResourceService rpc_v1.EnterpriseGlooResourceApiServer,
	ratelimitResourceService rpc_v1.RatelimitResourceApiServer,
	glooFedResourceService rpc_v1.FederatedGlooResourceApiServer,
	gatewayFedResourceService rpc_v1.FederatedGatewayResourceApiServer,
	glooEnterpriseFedResourceService rpc_v1.FederatedEnterpriseGlooResourceApiServer,
	ratelimitFedResourceService rpc_v1.FederatedRatelimitResourceApiServer,
) manager.Runnable {
	return manager.RunnableFunc(func(ctx context.Context) error {
		ctx = contextutils.WithLogger(rootCtx, "gloofed-apiserver")
		apiServer := server.NewGrpcServer(ctx, glooInstanceService, failoverSchemeService, routeTableSelectorService, wasmFilterService,
			glooResourceService, gatewayResourceService, glooEnterpriseResourceService, ratelimitResourceService, glooFedResourceService,
			gatewayFedResourceService, glooEnterpriseFedResourceService, ratelimitFedResourceService, health_check.NewHealthChecker())

		return apiServer.Run(ctx, cfg)
	})
}
