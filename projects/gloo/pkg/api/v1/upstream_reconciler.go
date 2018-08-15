package v1

import (
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/reconcile"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/solo-io/solo-kit/pkg/utils/contextutils"
)

// Option to copy anything from the original to the desired before writing
type TransitionUpstreamFunc func(original, desired *Upstream)

type UpstreamReconciler interface {
	Reconcile(namespace string, desiredResources []*Upstream, opts clients.ListOpts) error
}

func upstreamsToResources(list UpstreamList) []resources.Resource {
	var resourceList []resources.Resource
	for _, upstream := range list {
		resourceList = append(resourceList, upstream)
	}
	return resourceList
}

func NewUpstreamReconciler(client UpstreamClient, transition TransitionUpstreamFunc) UpstreamReconciler {
	var transitionResources reconcile.TransitionResourcesFunc
	if transition != nil {
		transitionResources = func(original, desired resources.Resource) {
			transition(original.(*Upstream), desired.(*Upstream))
		}
	}
	return &upstreamReconciler{
		base: reconcile.NewReconciler(client.BaseClient(), transitionResources),
	}
}

type upstreamReconciler struct {
	base reconcile.Reconciler
}

func (r *upstreamReconciler) Reconcile(namespace string, desiredResources []*Upstream, opts clients.ListOpts) error {
	opts = opts.WithDefaults()
	opts.Ctx = contextutils.WithLogger(opts.Ctx, "upstream_reconciler")
	return r.base.Reconcile(namespace, upstreamsToResources(desiredResources), opts)
}
