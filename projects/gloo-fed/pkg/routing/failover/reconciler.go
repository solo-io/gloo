package failover

import (
	"context"
	"fmt"

	"github.com/solo-io/k8s-utils/kubeutils"
	"github.com/solo-io/skv2/pkg/reconcile"
	gloov1 "github.com/solo-io/solo-apis/pkg/api/gloo.solo.io/v1"
	"github.com/solo-io/solo-kit/pkg/utils/prototime"
	fedv1 "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.solo.io/v1"
	"github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.solo.io/v1/controller"
	fed_types "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.solo.io/v1/types"
	"github.com/solo-io/solo-projects/projects/gloo-fed/pkg/federation"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
) failoverSchemeReconciler {
	return &failoverSchemeReconcilerImpl{
		ctx:                  ctx,
		processor:            processor,
		failoverSchemeClient: failoverSchemeClient,
		glooClientset:        glooClientset,
	}
}

type failoverSchemeReconcilerImpl struct {
	ctx                  context.Context
	processor            FailoverProcessor
	failoverSchemeClient fedv1.FailoverSchemeClient
	glooClientset        gloov1.MulticlusterClientset
}

func (f *failoverSchemeReconcilerImpl) ReconcileFailoverScheme(obj *fedv1.FailoverScheme) (reconcile.Result, error) {
	if obj.GetGeneration() == obj.Status.GetObservedGeneration() &&
		(obj.Status.GetState() == fed_types.FailoverSchemeStatus_INVALID ||
			obj.Status.GetState() == fed_types.FailoverSchemeStatus_ACCEPTED) {
		return reconcile.Result{}, nil
	}
	us, status := f.processor.ProcessFailoverUpdate(f.ctx, obj)
	if status != nil {
		obj.Status = *status
		if err := f.failoverSchemeClient.UpdateFailoverSchemeStatus(f.ctx, obj); err != nil {
			return reconcile.Result{}, err
		}
		return reconcile.Result{}, nil
	}

	statusBuilder := NewFailoverStatusBuilder(obj)
	glooClusterClientset, err := f.glooClientset.Cluster(obj.Spec.GetPrimary().GetClusterName())
	if err != nil {
		return reconcile.Result{}, statusBuilder.Fail(err).Update(f.ctx, f.failoverSchemeClient)
	}

	if err = glooClusterClientset.Upstreams().UpsertUpstream(f.ctx, us); err != nil {
		return reconcile.Result{}, statusBuilder.Fail(err).Update(f.ctx, f.failoverSchemeClient)
	}

	return reconcile.Result{}, statusBuilder.Accept().Update(f.ctx, f.failoverSchemeClient)
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

func UpstreamName(object metav1.Object, servicePort int32) string {
	return kubeutils.SanitizeNameV2(fmt.Sprintf("%s-%s-%v", object.GetNamespace(), object.GetName(), servicePort))
}

func NewFailoverStatusBuilder(
	obj *fedv1.FailoverScheme,
) *failoverStatusBuilder {
	return &failoverStatusBuilder{
		obj: obj,
	}
}

type failoverStatusBuilder struct {
	obj    *fedv1.FailoverScheme
	status fed_types.FailoverSchemeStatus
}

func (f *failoverStatusBuilder) Accept() *failoverStatusBuilder {

	f.status = fed_types.FailoverSchemeStatus{
		State:              fed_types.FailoverSchemeStatus_ACCEPTED,
		ObservedGeneration: f.obj.GetGeneration(),
		ProcessingTime:     prototime.Now(),
	}
	return f
}

func (f *failoverStatusBuilder) Fail(
	err error,
) *failoverStatusBuilder {
	f.status = fed_types.FailoverSchemeStatus{
		State:              fed_types.FailoverSchemeStatus_FAILED,
		Message:            err.Error(),
		ObservedGeneration: f.obj.GetGeneration(),
		ProcessingTime:     prototime.Now(),
	}
	return f

}

func (f *failoverStatusBuilder) Invalidate(
	err error,
) *failoverStatusBuilder {
	f.status = fed_types.FailoverSchemeStatus{
		State:              fed_types.FailoverSchemeStatus_INVALID,
		Message:            err.Error(),
		ObservedGeneration: f.obj.GetGeneration(),
		ProcessingTime:     prototime.Now(),
	}
	return f
}

func (f *failoverStatusBuilder) Update(ctx context.Context, failoverSchemeClient fedv1.FailoverSchemeClient) error {
	f.obj.Status = f.status
	return failoverSchemeClient.UpdateFailoverSchemeStatus(ctx, f.obj)
}

func (f *failoverStatusBuilder) Build() *fed_types.FailoverSchemeStatus {
	return &f.status
}
