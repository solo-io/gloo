package main

import (
	"context"

	"github.com/solo-io/go-utils/contextutils"
	client "github.com/solo-io/skv2/pkg/multicluster"
	"github.com/solo-io/skv2/pkg/multicluster/watch"
	enterprisev1 "github.com/solo-io/solo-apis/pkg/api/enterprise.gloo.solo.io/v1"
	gatewayv1 "github.com/solo-io/solo-apis/pkg/api/gateway.solo.io/v1"
	gloov1 "github.com/solo-io/solo-apis/pkg/api/gloo.solo.io/v1"
	ratelimitv1 "github.com/solo-io/solo-apis/pkg/api/ratelimit.solo.io/v1alpha1"
	"github.com/solo-io/solo-projects/pkg/license"
	"github.com/solo-io/solo-projects/projects/apiserver"
	"github.com/solo-io/solo-projects/projects/apiserver/internal/settings"
	enterprise_gloo_resource_handler "github.com/solo-io/solo-projects/projects/apiserver/pkg/api/enterprise.gloo.solo.io/v1/handler"
	federated_enterprise_gloo_resource_handler "github.com/solo-io/solo-projects/projects/apiserver/pkg/api/fed.enterprise.gloo.solo.io/v1/handler"
	federated_gateway_resource_handler "github.com/solo-io/solo-projects/projects/apiserver/pkg/api/fed.gateway.solo.io/v1/handler"
	federated_gloo_resource_handler "github.com/solo-io/solo-projects/projects/apiserver/pkg/api/fed.gloo.solo.io/v1/handler"
	federated_ratelimit_resource_handler "github.com/solo-io/solo-projects/projects/apiserver/pkg/api/fed.ratelimit.solo.io/v1alpha1/handler"
	gateway_resource_handler "github.com/solo-io/solo-projects/projects/apiserver/pkg/api/gateway.solo.io/v1/handler"
	gloo_resource_handler "github.com/solo-io/solo-projects/projects/apiserver/pkg/api/gloo.solo.io/v1/handler"
	ratelimit_resource_handler "github.com/solo-io/solo-projects/projects/apiserver/pkg/api/ratelimit.api.solo.io/v1alpha1/handler"
	"github.com/solo-io/solo-projects/projects/apiserver/server/apiserverutils"
	"github.com/solo-io/solo-projects/projects/apiserver/server/services/bootstrap_handler"
	"github.com/solo-io/solo-projects/projects/apiserver/server/services/failover_scheme_handler"
	"github.com/solo-io/solo-projects/projects/apiserver/server/services/glooinstance_handler"
	"github.com/solo-io/solo-projects/projects/apiserver/server/services/glooinstance_handler/config_getter"
	"github.com/solo-io/solo-projects/projects/apiserver/server/services/rt_selector_handler"
	"github.com/solo-io/solo-projects/projects/apiserver/server/services/wasmfilter_handler"
	glooentfedv1 "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.enterprise.gloo.solo.io/v1"
	gatewayfedv1 "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.gateway.solo.io/v1"
	gloofedv1 "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.gloo.solo.io/v1"
	ratelimitfedv1alpha1 "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.ratelimit.solo.io/v1alpha1"
	fedv1 "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.solo.io/v1"
	"github.com/solo-io/solo-projects/projects/gloo-fed/pkg/bootstrap"
	"github.com/solo-io/solo-projects/projects/gloo-fed/pkg/fields"
	"github.com/solo-io/solo-projects/projects/gloo-fed/pkg/multicluster"
	"go.uber.org/zap"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	// required import to enable kube client-go auth plugins
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

func main() {
	rootCtx := bootstrap.CreateRootContext(context.Background(), "gloo-fed-apiserver")
	cfg := settings.New()

	if err := license.IsLicenseValid(rootCtx, cfg.LicenseKey); err != nil {
		contextutils.LoggerFrom(rootCtx).Fatalw("License is invalid", zap.String("error", err.Error()))
	}

	mgr := bootstrap.MustLocalManager(rootCtx)

	glooFedEnabled, err := apiserverutils.IsGlooFedEnabled(rootCtx, mgr.GetConfig())
	if err != nil {
		contextutils.LoggerFrom(rootCtx).Fatalw("An error occurred when checking if Gloo Fed is enabled", zap.Error(err))
	}

	if glooFedEnabled {
		initializeGlooFed(rootCtx, mgr, cfg)
	} else {
		initializeGlooInstance(rootCtx, mgr, cfg)
	}

	err = mgr.Start(rootCtx)
	if err != nil {
		contextutils.LoggerFrom(rootCtx).Errorw("An error occurred", zap.Error(err))
	}
	contextutils.LoggerFrom(rootCtx).Infow("Shutting down, root context cancelled.")
}

func initializeGlooFed(ctx context.Context, mgr manager.Manager, cfg *settings.ApiServerSettings) {
	if err := fields.AddGlooInstanceIndexer(ctx, mgr); err != nil {
		contextutils.LoggerFrom(ctx).Fatalw("A fatal error occurred while adding cluster indexer to GlooInstance", zap.Error(err))
	}

	glooInstanceClient := fedv1.NewGlooInstanceClient(mgr.GetClient())
	failoverSchemeClient := fedv1.NewFailoverSchemeClient(mgr.GetClient())
	glooFedClient := gloofedv1.NewClientset(mgr.GetClient())
	glooEnterpriseFedClient := glooentfedv1.NewClientset(mgr.GetClient())
	gatewayFedClient := gatewayfedv1.NewClientset(mgr.GetClient())
	ratelimitFedClient := ratelimitfedv1alpha1.NewClientset(mgr.GetClient())
	clusterWatcher := watch.NewClusterWatcher(ctx, manager.Options{
		Scheme: bootstrap.MustRemoteScheme(ctx),
	})
	clusterSet := multicluster.NewClusterSet()
	clusterWatcher.RegisterClusterHandler(clusterSet)
	mcClient := client.NewClient(clusterWatcher)
	glooMCClient := gloov1.NewMulticlusterClientset(mcClient)
	gatewayMClient := gatewayv1.NewMulticlusterClientset(mcClient)
	glooEnterpriseMCClient := enterprisev1.NewMulticlusterClientset(mcClient)
	ratelimitMCCLient := ratelimitv1.NewMulticlusterClientset(mcClient)

	bootstrapService := bootstrap_handler.NewBootstrapHandler(mgr.GetConfig())
	glooInstanceService := glooinstance_handler.NewGlooInstanceHandler(clusterSet, config_getter.NewEnvoyConfigDumpGetter(clusterWatcher), glooInstanceClient)
	failoverSchemeService := failover_scheme_handler.NewFailoverSchemeHandler(failoverSchemeClient)
	routeTableSelectorService := rt_selector_handler.NewVirtualServiceRoutesHandler(gatewayMClient)
	wasmFilterService := wasmfilter_handler.NewWasmFilterHandler(glooInstanceClient, gatewayMClient)
	glooResourceService := gloo_resource_handler.NewGlooResourceHandler(glooInstanceClient, glooMCClient)
	glooEnterpriseResourceService := enterprise_gloo_resource_handler.NewEnterpriseGlooResourceHandler(glooInstanceClient, glooEnterpriseMCClient)
	ratelimitResourceService := ratelimit_resource_handler.NewRatelimitResourceHandler(glooInstanceClient, ratelimitMCCLient)
	gatewayResourceService := gateway_resource_handler.NewGatewayResourceHandler(glooInstanceClient, gatewayMClient)
	glooFedResourceService := federated_gloo_resource_handler.NewFederatedGlooResourceHandler(glooFedClient)
	gatewayFedResourceService := federated_gateway_resource_handler.NewFederatedGatewayResourceHandler(gatewayFedClient)
	glooEnterpriseFedResourceService := federated_enterprise_gloo_resource_handler.NewFederatedEnterpriseGlooResourceHandler(glooEnterpriseFedClient)
	ratelimitFedResourceService := federated_ratelimit_resource_handler.NewFederatedRatelimitResourceHandler(ratelimitFedClient)

	if err := mgr.Add(apiserver.NewGlooFedServerRunnable(ctx, cfg, bootstrapService, glooInstanceService, failoverSchemeService, routeTableSelectorService,
		wasmFilterService, glooResourceService, gatewayResourceService, glooEnterpriseResourceService, ratelimitResourceService,
		glooFedResourceService, gatewayFedResourceService, glooEnterpriseFedResourceService, ratelimitFedResourceService)); err != nil {
		contextutils.LoggerFrom(ctx).Fatalw("Unable to set up GlooFed apiserver", zap.Error(err))
	}

	if err := clusterWatcher.Run(mgr); err != nil {
		contextutils.LoggerFrom(ctx).Fatalw("A fatal error occurred while starting the cluster watcher", zap.Error(err))
	}
}

func initializeGlooInstance(ctx context.Context, mgr manager.Manager, cfg *settings.ApiServerSettings) {
	bootstrapService := bootstrap_handler.NewBootstrapHandler(mgr.GetConfig())

	if err := mgr.Add(apiserver.NewGlooInstanceServerRunnable(ctx, cfg, bootstrapService)); err != nil {
		contextutils.LoggerFrom(ctx).Fatalw("Unable to set up GlooEE apiserver", zap.Error(err))
	}
}
