package failover

import (
	"context"

	"github.com/solo-io/skv2/pkg/multicluster"
	"github.com/solo-io/skv2/pkg/reconcile"
	gloov1 "github.com/solo-io/solo-apis/pkg/api/gloo.solo.io/v1"
	controller2 "github.com/solo-io/solo-apis/pkg/api/gloo.solo.io/v1/controller"
	fedv1 "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.solo.io/v1"
	"github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.solo.io/v1/controller"
	"github.com/solo-io/solo-projects/projects/gloo-fed/pkg/routing/failover/internal"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

func NewFailoverDependentReconciler(
	ctx context.Context,
	failoverClient fedv1.FailoverSchemeClient,
	glooInstanceClient fedv1.GlooInstanceClient,
) internal.FailoverDependentReconciler {
	return internal.NewFailoverDependentReconciler(
		ctx,
		internal.NewFailoverDependencyCalculator(failoverClient, glooInstanceClient),
		failoverClient,
	)
}

func InitializeFailover(
	ctx context.Context,
	localManager manager.Manager,
	mcClient multicluster.Client,
	cw multicluster.ClusterWatcher,
) error {

	fedv1Clientset := fedv1.NewClientset(localManager.GetClient())
	failoverProcessor := NewFailoverProcessor(
		gloov1.NewMulticlusterClientset(mcClient),
		fedv1Clientset.GlooInstances(),
		fedv1Clientset.FailoverSchemes(),
	)
	failoverReconciler := NewFailoverSchemeReconciler(
		ctx,
		failoverProcessor,
		fedv1Clientset.FailoverSchemes(),
		gloov1.NewMulticlusterClientset(mcClient),
	)

	failoverDep := NewFailoverDependentReconciler(
		ctx,
		fedv1Clientset.FailoverSchemes(),
		fedv1Clientset.GlooInstances(),
	)
	controller2.NewMulticlusterUpstreamReconcileLoop("failover-dep", cw, reconcile.Options{}).
		AddMulticlusterUpstreamReconciler(ctx, failoverDep)
	if err := controller.NewGlooInstanceReconcileLoop("failover-dep", localManager, reconcile.Options{}).
		RunGlooInstanceReconciler(ctx, failoverDep); err != nil {
		return err
	}

	return controller.NewFailoverSchemeReconcileLoop("failover-scheme", localManager, reconcile.Options{}).
		RunFailoverSchemeReconciler(ctx, failoverReconciler)
}
