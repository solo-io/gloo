package mocks

import (
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/reconcile"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/solo-io/solo-kit/pkg/utils/contextutils"
)

// Option to copy anything from the original to the desired before writing
type TransitionMockResourceFunc func(original, desired *MockResource)

type MockResourceReconciler interface {
	Reconcile(namespace string, desiredResources []*MockResource, opts clients.ListOpts) error
}

func mockResourcesToResources(list MockResourceList) []resources.Resource {
	var resourceList []resources.Resource
	for _, mockResource := range list {
		resourceList = append(resourceList, mockResource)
	}
	return resourceList
}

func NewMockResourceReconciler(client MockResourceClient, transition TransitionMockResourceFunc) MockResourceReconciler {
	var transitionResources reconcile.TransitionResourcesFunc
	if transition != nil {
		transitionResources = func(original, desired resources.Resource) {
			transition(original.(*MockResource), desired.(*MockResource))
		}
	}
	return &mockResourceReconciler{
		base: reconcile.NewReconciler(client.BaseClient(), transitionResources),
	}
}

type mockResourceReconciler struct {
	base reconcile.Reconciler
}

func (r *mockResourceReconciler) Reconcile(namespace string, desiredResources []*MockResource, opts clients.ListOpts) error {
	opts = opts.WithDefaults()
	opts.Ctx = contextutils.WithLogger(opts.Ctx, "mockResource_reconciler")
	return r.base.Reconcile(namespace, mockResourcesToResources(desiredResources), opts)
}
