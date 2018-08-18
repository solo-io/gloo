package v1

import (
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/reconcile"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/solo-io/solo-kit/pkg/utils/contextutils"
)

// Option to copy anything from the original to the desired before writing
type TransitionUpstreamFunc func(original, desired *Upstream) error

type UpstreamReconciler interface {
	Reconcile(namespace string, desiredResources []*Upstream, transition TransitionUpstreamFunc, opts clients.ListOpts) error
}

func upstreamsToResources(list UpstreamList) resources.ResourceList {
	var resourceList resources.ResourceList
	for _, upstream := range list {
		resourceList = append(resourceList, upstream)
	}
	return resourceList
}

func NewUpstreamReconciler(client UpstreamClient) UpstreamReconciler {
	return &upstreamReconciler{
		base: reconcile.NewReconciler(client.BaseClient()),
	}
}

type upstreamReconciler struct {
	base reconcile.Reconciler
}

func (r *upstreamReconciler) Reconcile(namespace string, desiredResources []*Upstream, transition TransitionUpstreamFunc, opts clients.ListOpts) error {
	opts = opts.WithDefaults()
	opts.Ctx = contextutils.WithLogger(opts.Ctx, "upstream_reconciler")
	var transitionResources reconcile.TransitionResourcesFunc
	if transition != nil {
		transitionResources = func(original, desired resources.Resource) error {
			return transition(original.(*Upstream), desired.(*Upstream))
		}
	}
	return r.base.Reconcile(namespace, upstreamsToResources(desiredResources), transitionResources, opts)
}
