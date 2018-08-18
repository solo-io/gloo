package mocks

import (
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/reconcile"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/solo-io/solo-kit/pkg/utils/contextutils"
)

// Option to copy anything from the original to the desired before writing
type TransitionFakeResourceFunc func(original, desired *FakeResource) error

type FakeResourceReconciler interface {
	Reconcile(namespace string, desiredResources []*FakeResource, transition TransitionFakeResourceFunc, opts clients.ListOpts) error
}

func fakeResourcesToResources(list FakeResourceList) resources.ResourceList {
	var resourceList resources.ResourceList
	for _, fakeResource := range list {
		resourceList = append(resourceList, fakeResource)
	}
	return resourceList
}

func NewFakeResourceReconciler(client FakeResourceClient) FakeResourceReconciler {
	return &fakeResourceReconciler{
		base: reconcile.NewReconciler(client.BaseClient()),
	}
}

type fakeResourceReconciler struct {
	base reconcile.Reconciler
}

func (r *fakeResourceReconciler) Reconcile(namespace string, desiredResources []*FakeResource, transition TransitionFakeResourceFunc, opts clients.ListOpts) error {
	opts = opts.WithDefaults()
	opts.Ctx = contextutils.WithLogger(opts.Ctx, "fakeResource_reconciler")
	var transitionResources reconcile.TransitionResourcesFunc
	if transition != nil {
		transitionResources = func(original, desired resources.Resource) error {
			return transition(original.(*FakeResource), desired.(*FakeResource))
		}
	}
	return r.base.Reconcile(namespace, fakeResourcesToResources(desiredResources), transitionResources, opts)
}
