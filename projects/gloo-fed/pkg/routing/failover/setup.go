package failover

import (
	"context"

	"github.com/solo-io/skv2/pkg/multicluster"
	"github.com/solo-io/skv2/pkg/reconcile"
	gloov1 "github.com/solo-io/solo-apis/pkg/api/gloo.solo.io/v1"
	controller2 "github.com/solo-io/solo-apis/pkg/api/gloo.solo.io/v1/controller"
	fedv1 "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.solo.io/v1"
	"github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.solo.io/v1/controller"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

func InitializeFailover(
	ctx context.Context,
	localManager manager.Manager,
	mcClient multicluster.Client,
	cw multicluster.ClusterWatcher,
	reportingNamespace string,
) error {
	fedv1Clientset := fedv1.NewClientset(localManager.GetClient())

	failoverSchemeStatusManager := NewStatusManager(fedv1Clientset.FailoverSchemes(), reportingNamespace)

	failoverProcessor := NewFailoverProcessor(
		gloov1.NewMulticlusterClientset(mcClient),
		fedv1Clientset.GlooInstances(),
		fedv1Clientset.FailoverSchemes(),
		failoverSchemeStatusManager,
	)
	failoverReconciler := NewFailoverSchemeReconciler(
		ctx,
		failoverProcessor,
		fedv1Clientset.FailoverSchemes(),
		gloov1.NewMulticlusterClientset(mcClient),
		failoverSchemeStatusManager,
	)

	failoverDependentReconciler := NewFailoverDependentReconciler(
		ctx,
		NewFailoverDependencyCalculator(fedv1Clientset.FailoverSchemes(), fedv1Clientset.GlooInstances()),
		fedv1Clientset.FailoverSchemes(),
		failoverSchemeStatusManager,
	)

	controller2.NewMulticlusterUpstreamReconcileLoop("failover-dep", cw, reconcile.Options{}).
		AddMulticlusterUpstreamReconciler(ctx, failoverDependentReconciler)
	if err := controller.NewGlooInstanceReconcileLoop("failover-dep", localManager, reconcile.Options{}).
		RunGlooInstanceReconciler(ctx, failoverDependentReconciler); err != nil {
		return err
	}

	return controller.NewFailoverSchemeReconcileLoop("failover-scheme", localManager, reconcile.Options{}).
		RunFailoverSchemeReconciler(ctx, failoverReconciler)
}
