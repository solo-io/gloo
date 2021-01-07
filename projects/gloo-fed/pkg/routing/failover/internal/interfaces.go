package internal

import (
	"context"

	skv2v1 "github.com/solo-io/skv2/pkg/api/core.skv2.solo.io/v1"
	gloo_solo_io_v1_controllers "github.com/solo-io/solo-apis/pkg/api/gloo.solo.io/v1/controller"
	fedv1 "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.solo.io/v1"
	"github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.solo.io/v1/controller"
)

//go:generate mockgen -source ./interfaces.go -destination ./mocks/mock_interfaces.go

// FailoverDependencyCalculator calculates all of the Failover Schemes which depend on the incoming resource. In this
// case the 2 possible resources are Upstreams, and GlooInstances. The arguments are refs as opposed to objects, so
// deletes, which may not have access to the object, can be handled as well.
type FailoverDependencyCalculator interface {
	ForUpstream(ctx context.Context, upstream *skv2v1.ClusterObjectRef) ([]*fedv1.FailoverScheme, error)
	ForGlooInstance(ctx context.Context, glooInstance *skv2v1.ObjectRef) ([]*fedv1.FailoverScheme, error)
}

// Internal interface representing all functions linked to this reconciler
type FailoverDependentReconciler interface {
	controller.GlooInstanceReconciler
	controller.GlooInstanceDeletionReconciler
	gloo_solo_io_v1_controllers.MulticlusterUpstreamReconciler
	gloo_solo_io_v1_controllers.MulticlusterUpstreamDeletionReconciler
}
