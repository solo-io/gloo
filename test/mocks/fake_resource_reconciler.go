package mocks

import (
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/reconcile"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/solo-io/solo-kit/pkg/utils/contextutils"
)

// Option to copy anything from the original to the desired before writing
type TransitionFakeResourceFunc func(original, desired *FakeResource)

type FakeResourceReconciler interface {
	Reconcile(namespace string, desiredResources []*FakeResource, opts clients.ListOpts) error
}

func fakeResourcesToResources(list FakeResourceList) []resources.Resource {
	var resourceList []resources.Resource
	for _, fakeResource := range list {
		resourceList = append(resourceList, fakeResource)
	}
	return resourceList
}

func NewFakeResourceReconciler(client FakeResourceClient, transition TransitionFakeResourceFunc) FakeResourceReconciler {
	var transitionResources reconcile.TransitionResourcesFunc
	if transition != nil {
		transitionResources = func(original, desired resources.Resource) {
			transition(original.(*FakeResource), desired.(*FakeResource))
		}
	}
	return &fakeResourceReconciler{
		base: reconcile.NewReconciler(client.BaseClient(), transitionResources),
	}
}

type fakeResourceReconciler struct {
	base reconcile.Reconciler
}

func (r *fakeResourceReconciler) Reconcile(namespace string, desiredResources []*FakeResource, opts clients.ListOpts) error {
	opts = opts.WithDefaults()
	opts.Ctx = contextutils.WithLogger(opts.Ctx, "fakeResource_reconciler")
	return r.base.Reconcile(namespace, fakeResourcesToResources(desiredResources), opts)
}
