package failover

import (
	"context"

	skv2v1 "github.com/solo-io/skv2/pkg/api/core.skv2.solo.io/v1"
	gloo_solo_io_v1_controllers "github.com/solo-io/solo-apis/pkg/api/gloo.solo.io/v1/controller"
	"github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.solo.io/v1/controller"

	gloov1 "github.com/solo-io/solo-apis/pkg/api/gloo.solo.io/v1"
	fedv1 "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.solo.io/v1"
)

const (
	PortNumber           = 15443
	PortName             = "failover"
	DownstreamSecretName = "failover-downstream"
	UpstreamSecretName   = "failover-upstream"
)

//go:generate mockgen -source ./interfaces.go -destination ./mocks/mock_interfaces.go

// FailoverProcessor calculates the output Upstream for a given change/event of a FailoverProcessor.
type FailoverProcessor interface {
	// ProcessFailoverUpdate returns a status similar to a typical error return value, if nil it can be ignored,
	// otherwise it should be saved back to the parent object.
	ProcessFailoverUpdate(
		ctx context.Context,
		obj *fedv1.FailoverScheme,
	) (*gloov1.Upstream, StatusBuilder)
	// ProcessFailoverDelete does not return a status as the object is being deleted.
	// May return nil, nil if primary is nil, as nothing can be done, and there is no "error"
	ProcessFailoverDelete(ctx context.Context, obj *fedv1.FailoverScheme) (*gloov1.Upstream, error)
}

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
