package runner

import (
	"context"

	errors "github.com/rotisserie/eris"
	client "github.com/solo-io/skv2/pkg/multicluster"
	"github.com/solo-io/skv2/pkg/multicluster/watch"
	enterprisev1 "github.com/solo-io/solo-apis/pkg/api/enterprise.gloo.solo.io/v1"
	gatewayv1 "github.com/solo-io/solo-apis/pkg/api/gateway.solo.io/v1"
	gloov1 "github.com/solo-io/solo-apis/pkg/api/gloo.solo.io/v1"
	ratelimitv1alpha1 "github.com/solo-io/solo-apis/pkg/api/ratelimit.solo.io/v1alpha1"
	"github.com/solo-io/solo-projects/pkg/license"
	fedenterprisev1 "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.enterprise.gloo.solo.io/v1"
	enterprisefed "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.enterprise.gloo.solo.io/v1/federation"
	fedgatewayv1 "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.gateway.solo.io/v1"
	gatewayfed "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.gateway.solo.io/v1/federation"
	fedgloov1 "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.gloo.solo.io/v1"
	gloofed "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.gloo.solo.io/v1/federation"
	fedratelimitv1alpha1 "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.ratelimit.solo.io/v1alpha1"
	ratelimitfed "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.ratelimit.solo.io/v1alpha1/federation"
	fed_bootstrap "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/bootstrap"
	"github.com/solo-io/solo-projects/projects/gloo-fed/pkg/discovery"
	"github.com/solo-io/solo-projects/projects/gloo-fed/pkg/federation/placement"
	"github.com/solo-io/solo-projects/projects/gloo-fed/pkg/fields"
	"github.com/solo-io/solo-projects/projects/gloo-fed/pkg/multicluster"
	"github.com/solo-io/solo-projects/projects/gloo-fed/pkg/routing/failover"
	"go.uber.org/zap"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

// Run starts running the Gloo Federation component.
// Run blocks until the context is closed or an error is return
func Run(runCtx context.Context, settings *Settings) error {
	// Validate that the provided license supports Federation
	licensedFeatureProvider := license.NewLicensedFeatureProvider()
	licensedFeatureProvider.ValidateAndSetLicense(settings.LicenseKey)
	federationFeatureState := licensedFeatureProvider.GetStateForLicensedFeature(license.Enterprise)
	if !federationFeatureState.Enabled {
		return errors.Errorf("Federation is disabled", zap.String("reason", federationFeatureState.Reason))
	}

	// Initialize the NewManager
	// This is the component responsible for managing clients and caches and supplying them to controllers
	mgr := fed_bootstrap.MustLocalManager(runCtx)

	if err := fields.AddGlooInstanceIndexer(runCtx, mgr); err != nil {
		return errors.Errorf("A fatal error occurred while adding cluster indexer to GlooInstance", zap.Error(err))
	}

	clusterWatcher := watch.NewClusterWatcher(runCtx, manager.Options{
		Scheme: fed_bootstrap.MustRemoteScheme(runCtx),
	}, []string{fed_bootstrap.GetInstallNamespace()})

	// Generate the set of cluster (names) that the cluster watcher will be registered with
	clusterSet := multicluster.NewClusterSet()
	clusterWatcher.RegisterClusterHandler(clusterSet)

	placementManager := placement.NewManager(settings.PodNamespace, settings.PodName)

	// Register a ClusterHandler for each type of Federated Edge resource
	registerFederatedResourceClusterHandlers(runCtx, clusterWatcher, mgr, placementManager)

	multiclusterClient := client.NewClient(clusterWatcher)
	discovery.InitializeDiscovery(runCtx, settings.WriteNamespace, mgr, multiclusterClient, clusterWatcher)

	if err := failover.InitializeFailover(runCtx, mgr, multiclusterClient, clusterWatcher, settings.PodNamespace); err != nil {
		return errors.Errorf("A fatal error occurred while setting up failover reconciler", zap.Error(err))
	}

	if err := runFederatedResourceReconcilers(runCtx, mgr, multiclusterClient, clusterSet, placementManager); err != nil {
		return errors.Errorf("A fatal error occurred while federated resource reconcilers", zap.Error(err))
	}

	if err := clusterWatcher.Run(mgr); err != nil {
		return errors.Errorf("A fatal error occurred while starting the cluster watcher", zap.Error(err))
	}

	// Start starts all registered Controllers and blocks until the context is cancelled.
	// Returns an error if there is an error starting any controller.
	return mgr.Start(runCtx)
}

func registerFederatedResourceClusterHandlers(
	ctx context.Context,
	clusterWatcher client.ClusterWatcher,
	mgr manager.Manager,
	placementManager placement.Manager,
) {
	fedGlooClusterHandler := gloofed.NewClusterHandler(ctx, fedgloov1.NewClientset(mgr.GetClient()), placementManager)
	clusterWatcher.RegisterClusterHandler(fedGlooClusterHandler)

	fedGatewayClusterHandler := gatewayfed.NewClusterHandler(ctx, fedgatewayv1.NewClientset(mgr.GetClient()), placementManager)
	clusterWatcher.RegisterClusterHandler(fedGatewayClusterHandler)

	fedEnterpriseGlooClusterHandler := enterprisefed.NewClusterHandler(ctx, fedenterprisev1.NewClientset(mgr.GetClient()), placementManager)
	clusterWatcher.RegisterClusterHandler(fedEnterpriseGlooClusterHandler)

	fedRatelimitClusterHandler := ratelimitfed.NewClusterHandler(ctx, fedratelimitv1alpha1.NewClientset(mgr.GetClient()), placementManager)
	clusterWatcher.RegisterClusterHandler(fedRatelimitClusterHandler)
}

func runFederatedResourceReconcilers(
	ctx context.Context,
	mgr manager.Manager,
	multiclusterClient client.Client,
	clusterSet multicluster.ClusterSet,
	placementManager placement.Manager,
) error {
	if err := gatewayfed.Initialize(ctx, mgr, gatewayv1.NewMulticlusterClientset(multiclusterClient), clusterSet, placementManager); err != nil {
		return err
	}

	if err := gloofed.Initialize(ctx, mgr, gloov1.NewMulticlusterClientset(multiclusterClient), clusterSet, placementManager); err != nil {
		return err
	}

	if err := enterprisefed.Initialize(ctx, mgr, enterprisev1.NewMulticlusterClientset(multiclusterClient), clusterSet, placementManager); err != nil {
		return err
	}

	if err := ratelimitfed.Initialize(ctx, mgr, ratelimitv1alpha1.NewMulticlusterClientset(multiclusterClient), clusterSet, placementManager); err != nil {
		return err
	}

	// all reconcilers started successfully
	return nil
}
