package internal

import (
	"context"

	skv2v1 "github.com/solo-io/skv2/pkg/api/core.skv2.solo.io/v1"
	"github.com/solo-io/skv2/pkg/reconcile"
	gloov1 "github.com/solo-io/solo-apis/pkg/api/gloo.solo.io/v1"
	"github.com/solo-io/solo-kit/pkg/utils/prototime"
	fedv1 "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.solo.io/v1"
	fed_types "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.solo.io/v1/types"
	"k8s.io/apimachinery/pkg/api/errors"
)

func IgnoreIsConflictError(err error) error {
	if errors.IsConflict(err) {
		return nil
	}
	return err
}

const (
	DependentUpdateMessage = "dependent has been updated"
)

func NewFailoverDependentReconciler(
	ctx context.Context,
	depCalc FailoverDependencyCalculator,
	failoverSchemeClient fedv1.FailoverSchemeClient,
) FailoverDependentReconciler {
	return &failoverDependentReconcilerImpl{
		ctx:                  ctx,
		depCalc:              depCalc,
		failoverSchemeClient: failoverSchemeClient,
	}
}

type failoverDependentReconcilerImpl struct {
	ctx                  context.Context
	depCalc              FailoverDependencyCalculator
	failoverSchemeClient fedv1.FailoverSchemeClient
}

func (f *failoverDependentReconcilerImpl) ReconcileGlooInstance(obj *fedv1.GlooInstance) (reconcile.Result, error) {
	failoverSchemes, err := f.depCalc.ForGlooInstance(f.ctx, &skv2v1.ObjectRef{
		Name:      obj.GetName(),
		Namespace: obj.GetNamespace(),
	})
	if err != nil {
		return reconcile.Result{}, err
	}
	return reconcile.Result{}, f.updateMultiStatus(failoverSchemes)
}

func (f *failoverDependentReconcilerImpl) ReconcileGlooInstanceDeletion(req reconcile.Request) error {
	failoverSchemes, err := f.depCalc.ForGlooInstance(f.ctx, &skv2v1.ObjectRef{
		Name:      req.Name,
		Namespace: req.Namespace,
	})
	if err != nil {
		return err
	}
	return f.updateMultiStatus(failoverSchemes)
}

func (f *failoverDependentReconcilerImpl) ReconcileUpstream(
	clusterName string,
	obj *gloov1.Upstream,
) (reconcile.Result, error) {
	failoverSchemes, err := f.depCalc.ForUpstream(f.ctx, &skv2v1.ClusterObjectRef{
		Name:        obj.GetName(),
		Namespace:   obj.GetNamespace(),
		ClusterName: clusterName,
	})
	if err != nil {
		return reconcile.Result{}, err
	}
	return reconcile.Result{}, f.updateMultiStatus(failoverSchemes)
}

func (f *failoverDependentReconcilerImpl) ReconcileUpstreamDeletion(
	clusterName string,
	req reconcile.Request,
) error {
	failoverSchemes, err := f.depCalc.ForUpstream(f.ctx, &skv2v1.ClusterObjectRef{
		Name:        req.Name,
		Namespace:   req.Namespace,
		ClusterName: clusterName,
	})
	if err != nil {
		return err
	}
	return f.updateMultiStatus(failoverSchemes)
}

// updateMultiStatus attempts to update the status on all passed in failover schemes. Each Failover Scheme passed in
// represents a failover scheme which has been calculated to depend on either a Gloo Instance or an Upstream which
// has changed. Will return on first failed update.
func (f *failoverDependentReconcilerImpl) updateMultiStatus(failoverSchemes []*fedv1.FailoverScheme) error {
	for _, failoverScheme := range failoverSchemes {
		failoverScheme.Status = fed_types.FailoverSchemeStatus{
			State:              fed_types.FailoverSchemeStatus_PENDING,
			Message:            DependentUpdateMessage,
			ObservedGeneration: failoverScheme.GetGeneration(),
			ProcessingTime:     prototime.Now(),
		}
		// Only return an error if there is no update conflict, If the CRD has been updated than we can most
		// likely skip reprocessing this event
		if err := IgnoreIsConflictError(f.failoverSchemeClient.UpdateFailoverSchemeStatus(f.ctx, failoverScheme)); err != nil {
			return err
		}
	}
	return nil
}
