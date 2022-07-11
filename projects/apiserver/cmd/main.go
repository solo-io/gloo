package main

import (
	"context"
	"os"

	apps_v1 "github.com/solo-io/external-apis/pkg/api/k8s/apps/v1"
	core_v1 "github.com/solo-io/external-apis/pkg/api/k8s/core/v1"
	"github.com/solo-io/go-utils/contextutils"
	client "github.com/solo-io/skv2/pkg/multicluster"
	"github.com/solo-io/skv2/pkg/multicluster/watch"
	enterprise_gloo_v1 "github.com/solo-io/solo-apis/pkg/api/enterprise.gloo.solo.io/v1"
	gateway_v1 "github.com/solo-io/solo-apis/pkg/api/gateway.solo.io/v1"
	gloo_v1 "github.com/solo-io/solo-apis/pkg/api/gloo.solo.io/v1"
	graphql_v1beta1 "github.com/solo-io/solo-apis/pkg/api/graphql.gloo.solo.io/v1beta1"
	ratelimit_v1alpha1 "github.com/solo-io/solo-apis/pkg/api/ratelimit.solo.io/v1alpha1"
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
	"github.com/solo-io/solo-projects/projects/apiserver/server/services/glooinstance_handler/envoy_admin"
	"github.com/solo-io/solo-projects/projects/apiserver/server/services/graphql_handler"
	"github.com/solo-io/solo-projects/projects/apiserver/server/services/rt_selector_handler"
	"github.com/solo-io/solo-projects/projects/apiserver/server/services/single_cluster_resource_handler"
	"github.com/solo-io/solo-projects/projects/apiserver/server/services/wasmfilter_handler"
	glooentfedv1 "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.enterprise.gloo.solo.io/v1"
	gatewayfedv1 "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.gateway.solo.io/v1"
	gloofedv1 "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.gloo.solo.io/v1"
	ratelimitfedv1alpha1 "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.ratelimit.solo.io/v1alpha1"
	fedv1 "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.solo.io/v1"
	fed_bootstrap "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/bootstrap"
	"github.com/solo-io/solo-projects/projects/gloo-fed/pkg/fields"
	"github.com/solo-io/solo-projects/projects/gloo-fed/pkg/multicluster"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/bootstrap"
	"go.uber.org/zap"
	"k8s.io/client-go/discovery"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	// required import to enable kube client-go auth plugins
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

func main() {
	rootCtx := bootstrap.CreateRootContext(context.Background(), "gloo-fed-apiserver")
	apiserverSettings := settings.New()

	licensedFeatureProvider := license.NewLicensedFeatureProvider()
	licensedFeatureProvider.ValidateAndSetLicense(apiserverSettings.LicenseKey)

	apiServerFeatureState := licensedFeatureProvider.GetStateForLicensedFeature(license.Enterprise)
	if !apiServerFeatureState.Enabled {
		contextutils.LoggerFrom(rootCtx).Fatalw("ApiServer is disabled", zap.String("reason", apiServerFeatureState.Reason))
	}

	cfg, err := config.GetConfig()
	if err != nil {
		contextutils.LoggerFrom(rootCtx).Fatalw("A fatal error occurred while getting config", zap.Error(err))
	}
	glooFedEnabled, err := apiserverutils.IsGlooFedEnabled(rootCtx, cfg)
	if err != nil {
		contextutils.LoggerFrom(rootCtx).Fatalw("An error occurred when checking if Gloo Fed is enabled", zap.Error(err))
	}

	var mgr manager.Manager
	if glooFedEnabled {
		mgr = fed_bootstrap.MustLocalManagerFromConfig(rootCtx, cfg)
		initializeGlooFed(rootCtx, mgr, apiserverSettings, licensedFeatureProvider)
	} else {
		if os.Getenv("NAMESPACE_RESTRICTED_MODE") == "true" {
			mgr = fed_bootstrap.MustSingleClusterManagerFromConfig(rootCtx, cfg, apiserverutils.GetInstallNamespace())
		} else {
			mgr = fed_bootstrap.MustSingleClusterManagerFromConfig(rootCtx, cfg, "")
		}
		initializeSingleClusterGloo(rootCtx, mgr, apiserverSettings, licensedFeatureProvider)
	}

	err = mgr.Start(rootCtx)
	if err != nil {
		contextutils.LoggerFrom(rootCtx).Errorw("An error occurred", zap.Error(err))
	}
	contextutils.LoggerFrom(rootCtx).Infow("Shutting down, root context cancelled.")
}

func initializeGlooFed(ctx context.Context, mgr manager.Manager, apiserverSettings *settings.ApiServerSettings,
	licensedFeatureProvider *license.LicensedFeatureProvider) {
	if err := fields.AddGlooInstanceIndexer(ctx, mgr); err != nil {
		contextutils.LoggerFrom(ctx).Fatalw("A fatal error occurred while adding cluster indexer to GlooInstance", zap.Error(err))
	}

	glooClientset := gloo_v1.NewClientset(mgr.GetClient())
	glooInstanceClient := fedv1.NewGlooInstanceClient(mgr.GetClient())
	failoverSchemeClient := fedv1.NewFailoverSchemeClient(mgr.GetClient())
	glooFedClient := gloofedv1.NewClientset(mgr.GetClient())
	glooEnterpriseFedClient := glooentfedv1.NewClientset(mgr.GetClient())
	gatewayFedClient := gatewayfedv1.NewClientset(mgr.GetClient())
	ratelimitFedClient := ratelimitfedv1alpha1.NewClientset(mgr.GetClient())
	clusterWatcher := watch.NewClusterWatcher(ctx, manager.Options{
		Scheme: fed_bootstrap.MustRemoteScheme(ctx),
	})
	clusterSet := multicluster.NewClusterSet()
	clusterWatcher.RegisterClusterHandler(clusterSet)
	mcClient := client.NewClient(clusterWatcher)
	glooMCClient := gloo_v1.NewMulticlusterClientset(mcClient)
	gatewayMCClient := gateway_v1.NewMulticlusterClientset(mcClient)
	graphqlMCClient := graphql_v1beta1.NewMulticlusterClientset(mcClient)
	glooEnterpriseMCClient := enterprise_gloo_v1.NewMulticlusterClientset(mcClient)
	ratelimitMCCLient := ratelimit_v1alpha1.NewMulticlusterClientset(mcClient)

	bootstrapService := bootstrap_handler.NewBootstrapHandler(mgr, licensedFeatureProvider)
	glooInstanceService := glooinstance_handler.NewFedGlooInstanceHandler(clusterWatcher, clusterSet, envoy_admin.NewEnvoyAdminClient(), glooInstanceClient)
	failoverSchemeService := failover_scheme_handler.NewFailoverSchemeHandler(failoverSchemeClient)
	routeTableSelectorService := rt_selector_handler.NewFedVirtualServiceRoutesHandler(gatewayMCClient)
	wasmFilterService := wasmfilter_handler.NewFedWasmFilterHandler(glooInstanceClient, gatewayMCClient)
	graphqlService := graphql_handler.NewFedGraphqlHandler(glooInstanceClient, glooClientset.Settings(), graphqlMCClient)
	glooResourceService := gloo_resource_handler.NewFedGlooResourceHandler(glooInstanceClient, glooMCClient)
	glooEnterpriseResourceService := enterprise_gloo_resource_handler.NewFedEnterpriseGlooResourceHandler(glooInstanceClient, glooEnterpriseMCClient)
	ratelimitResourceService := ratelimit_resource_handler.NewFedRatelimitResourceHandler(glooInstanceClient, ratelimitMCCLient)
	gatewayResourceService := gateway_resource_handler.NewFedGatewayResourceHandler(glooInstanceClient, gatewayMCClient)
	glooFedResourceService := federated_gloo_resource_handler.NewFederatedGlooResourceHandler(glooFedClient)
	gatewayFedResourceService := federated_gateway_resource_handler.NewFederatedGatewayResourceHandler(gatewayFedClient)
	glooEnterpriseFedResourceService := federated_enterprise_gloo_resource_handler.NewFederatedEnterpriseGlooResourceHandler(glooEnterpriseFedClient)
	ratelimitFedResourceService := federated_ratelimit_resource_handler.NewFederatedRatelimitResourceHandler(ratelimitFedClient)

	if err := mgr.Add(apiserver.NewGlooFedServerRunnable(ctx, apiserverSettings, bootstrapService, glooInstanceService, failoverSchemeService, routeTableSelectorService,
		wasmFilterService, graphqlService, glooResourceService, gatewayResourceService, glooEnterpriseResourceService, ratelimitResourceService,
		glooFedResourceService, gatewayFedResourceService, glooEnterpriseFedResourceService, ratelimitFedResourceService)); err != nil {
		contextutils.LoggerFrom(ctx).Fatalw("Unable to set up GlooFed apiserver", zap.Error(err))
	}

	if err := clusterWatcher.Run(mgr); err != nil {
		contextutils.LoggerFrom(ctx).Fatalw("A fatal error occurred while starting the cluster watcher", zap.Error(err))
	}
}

func initializeSingleClusterGloo(ctx context.Context, mgr manager.Manager, apiserverSettings *settings.ApiServerSettings,
	licensedFeatureProvider *license.LicensedFeatureProvider) {
	coreClientset := core_v1.NewClientset(mgr.GetClient())
	appsClientset := apps_v1.NewClientset(mgr.GetClient())
	gatewayClientset := gateway_v1.NewClientset(mgr.GetClient())
	glooClientset := gloo_v1.NewClientset(mgr.GetClient())
	graphqlClientset := graphql_v1beta1.NewClientset(mgr.GetClient())
	enterpriseGlooClientset := enterprise_gloo_v1.NewClientset(mgr.GetClient())
	ratelimitClientset := ratelimit_v1alpha1.NewClientset(mgr.GetClient())

	bootstrapService := bootstrap_handler.NewBootstrapHandler(mgr, licensedFeatureProvider)
	glooInstanceLister := glooinstance_handler.NewSingleClusterGlooInstanceLister(coreClientset, appsClientset, gatewayClientset, glooClientset, enterpriseGlooClientset, ratelimitClientset)
	discoveryClient, err := discovery.NewDiscoveryClientForConfig(mgr.GetConfig())
	if err != nil {
		contextutils.LoggerFrom(ctx).Fatalw("Could not create discovery client", zap.Error(err))
	}
	restClient := discoveryClient.RESTClient()
	glooInstanceService := glooinstance_handler.NewSingleClusterGlooInstanceHandler(glooInstanceLister, restClient, envoy_admin.NewEnvoyAdminClient())

	routeTableSelectorService := rt_selector_handler.NewSingleClusterVirtualServiceRoutesHandler()
	wasmFilterService := wasmfilter_handler.NewSingleClusterWasmFilterHandler()
	graphqlService := graphql_handler.NewSingleClusterGraphqlHandler(graphqlClientset, glooInstanceLister, glooClientset.Settings())

	// generated resource handlers
	gatewayResourceService := single_cluster_resource_handler.NewSingleClusterGatewayResourceHandler(gatewayClientset, glooInstanceLister)
	glooResourceService := single_cluster_resource_handler.NewSingleClusterGlooResourceHandler(glooClientset, glooInstanceLister)
	glooEnterpriseResourceService := single_cluster_resource_handler.NewSingleClusterEnterpriseGlooResourceHandler(enterpriseGlooClientset, glooInstanceLister)
	ratelimitResourceService := single_cluster_resource_handler.NewSingleClusterRatelimitResourceHandler(ratelimitClientset, glooInstanceLister)

	if err := mgr.Add(apiserver.NewSingleClusterGlooServerRunnable(ctx, apiserverSettings, bootstrapService, glooInstanceService,
		routeTableSelectorService, wasmFilterService, graphqlService,
		gatewayResourceService, glooResourceService, glooEnterpriseResourceService, ratelimitResourceService)); err != nil {
		contextutils.LoggerFrom(ctx).Fatalw("Unable to set up GlooEE apiserver", zap.Error(err))
	}
}
