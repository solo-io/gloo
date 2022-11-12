package failover

import (
	"context"

	"github.com/solo-io/skv2/pkg/reconcile"
	gloov1 "github.com/solo-io/solo-apis/pkg/api/gloo.solo.io/v1"
	fedv1 "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.solo.io/v1"
	"github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.solo.io/v1/controller"
	fed_types "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.solo.io/v1/types"
	"github.com/solo-io/solo-projects/projects/gloo-fed/pkg/federation"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Internal interface representing all functions linked to this reconciler
type failoverSchemeReconciler interface {
	controller.FailoverSchemeDeletionReconciler
	controller.FailoverSchemeFinalizer
}

func NewFailoverSchemeReconciler(
	ctx context.Context,
	processor FailoverProcessor,
	failoverSchemeClient fedv1.FailoverSchemeClient,
	glooClientset gloov1.MulticlusterClientset,
	statusManager *StatusManager,
) failoverSchemeReconciler {
	return &failoverSchemeReconcilerImpl{
		ctx:                  ctx,
		processor:            processor,
		failoverSchemeClient: failoverSchemeClient,
		glooClientset:        glooClientset,
		statusManager:        statusManager,
	}
}

type failoverSchemeReconcilerImpl struct {
	ctx                  context.Context
	processor            FailoverProcessor
	failoverSchemeClient fedv1.FailoverSchemeClient
	glooClientset        gloov1.MulticlusterClientset
	statusManager        *StatusManager
}

func (f *failoverSchemeReconcilerImpl) ReconcileFailoverScheme(obj *fedv1.FailoverScheme) (reconcile.Result, error) {
	currentStatus := f.statusManager.GetStatus(obj)
	if obj.GetGeneration() == currentStatus.GetObservedGeneration() &&
		(currentStatus.GetState() == fed_types.FailoverSchemeStatus_INVALID ||
			currentStatus.GetState() == fed_types.FailoverSchemeStatus_ACCEPTED) {
		return reconcile.Result{}, nil
	}
	us, statusBuilder := f.processor.ProcessFailoverUpdate(f.ctx, obj)
	if statusBuilder != nil {
		err := f.statusManager.UpdateStatus(f.ctx, statusBuilder)
		return reconcile.Result{}, err
	}

	statusBuilder = f.statusManager.NewStatusBuilder(obj)
	glooClusterClientset, err := f.glooClientset.Cluster(obj.Spec.GetPrimary().GetClusterName())
	if err != nil {
		return reconcile.Result{}, f.statusManager.UpdateStatus(f.ctx, statusBuilder.Fail(err))
	}

	if err = glooClusterClientset.Upstreams().UpsertUpstream(f.ctx, us); err != nil {
		return reconcile.Result{}, f.statusManager.UpdateStatus(f.ctx, statusBuilder.Fail(err))
	}

	return reconcile.Result{}, f.statusManager.UpdateStatus(f.ctx, statusBuilder.Accept())
}

func (f *failoverSchemeReconcilerImpl) ReconcileFailoverSchemeDeletion(req reconcile.Request) error {
	return nil
}

func (f *failoverSchemeReconcilerImpl) FinalizeFailoverScheme(obj *fedv1.FailoverScheme) error {
	upstream, err := f.processor.ProcessFailoverDelete(f.ctx, obj)
	if err != nil {
		return client.IgnoreNotFound(err)
	} else if upstream == nil {
		return nil
	}

	glooClusterClientset, err := f.glooClientset.Cluster(obj.Spec.GetPrimary().GetClusterName())
	if err != nil {
		return err
	}

	return client.IgnoreNotFound(glooClusterClientset.Upstreams().UpsertUpstream(f.ctx, upstream))
}

func (f *failoverSchemeReconcilerImpl) FailoverSchemeFinalizerName() string {
	return federation.HubFinalizer
}
