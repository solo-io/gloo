package v1

import (
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/reconcile"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/solo-io/solo-kit/pkg/utils/contextutils"
)

// Option to copy anything from the original to the desired before writing
type TransitionEndpointFunc func(original, desired *Endpoint) error

type EndpointReconciler interface {
	Reconcile(namespace string, desiredResources []*Endpoint, transition TransitionEndpointFunc, opts clients.ListOpts) error
}

func endpointsToResources(list EndpointList) []resources.Resource {
	var resourceList []resources.Resource
	for _, endpoint := range list {
		resourceList = append(resourceList, endpoint)
	}
	return resourceList
}

func NewEndpointReconciler(client EndpointClient) EndpointReconciler {
	return &endpointReconciler{
		base: reconcile.NewReconciler(client.BaseClient()),
	}
}

type endpointReconciler struct {
	base reconcile.Reconciler
}

func (r *endpointReconciler) Reconcile(namespace string, desiredResources []*Endpoint, transition TransitionEndpointFunc, opts clients.ListOpts) error {
	opts = opts.WithDefaults()
	opts.Ctx = contextutils.WithLogger(opts.Ctx, "endpoint_reconciler")
	var transitionResources reconcile.TransitionResourcesFunc
	if transition != nil {
		transitionResources = func(original, desired resources.Resource) error {
			return transition(original.(*Endpoint), desired.(*Endpoint))
		}
	}
	return r.base.Reconcile(namespace, endpointsToResources(desiredResources), transitionResources, opts)
}
