package v1

import (
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/reconcile"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/solo-io/solo-kit/pkg/utils/contextutils"
)

// Option to copy anything from the original to the desired before writing
type TransitionVirtualServiceFunc func(original, desired *VirtualService) error

type VirtualServiceReconciler interface {
	Reconcile(namespace string, desiredResources []*VirtualService, transition TransitionVirtualServiceFunc, opts clients.ListOpts) error
}

func virtualServicesToResources(list VirtualServiceList) []resources.Resource {
	var resourceList []resources.Resource
	for _, virtualService := range list {
		resourceList = append(resourceList, virtualService)
	}
	return resourceList
}

func NewVirtualServiceReconciler(client VirtualServiceClient) VirtualServiceReconciler {
	return &virtualServiceReconciler{
		base: reconcile.NewReconciler(client.BaseClient()),
	}
}

type virtualServiceReconciler struct {
	base reconcile.Reconciler
}

func (r *virtualServiceReconciler) Reconcile(namespace string, desiredResources []*VirtualService, transition TransitionVirtualServiceFunc, opts clients.ListOpts) error {
	opts = opts.WithDefaults()
	opts.Ctx = contextutils.WithLogger(opts.Ctx, "virtualService_reconciler")
	var transitionResources reconcile.TransitionResourcesFunc
	if transition != nil {
		transitionResources = func(original, desired resources.Resource) error {
			return transition(original.(*VirtualService), desired.(*VirtualService))
		}
	}
	return r.base.Reconcile(namespace, virtualServicesToResources(desiredResources), transitionResources, opts)
}
