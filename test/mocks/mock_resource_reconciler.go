package mocks

import (
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/reconcile"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/solo-io/solo-kit/pkg/utils/contextutils"
)

// Option to copy anything from the original to the desired before writing
type TransitionMockResourceFunc func(original, desired *MockResource) error

type MockResourceReconciler interface {
	Reconcile(namespace string, desiredResources []*MockResource, transition TransitionMockResourceFunc, opts clients.ListOpts) error
}

func mockResourcesToResources(list MockResourceList) []resources.Resource {
	var resourceList []resources.Resource
	for _, mockResource := range list {
		resourceList = append(resourceList, mockResource)
	}
	return resourceList
}

func NewMockResourceReconciler(client MockResourceClient) MockResourceReconciler {
	return &mockResourceReconciler{
		base: reconcile.NewReconciler(client.BaseClient()),
	}
}

type mockResourceReconciler struct {
	base reconcile.Reconciler
}

func (r *mockResourceReconciler) Reconcile(namespace string, desiredResources []*MockResource, transition TransitionMockResourceFunc, opts clients.ListOpts) error {
	opts = opts.WithDefaults()
	opts.Ctx = contextutils.WithLogger(opts.Ctx, "mockResource_reconciler")
	var transitionResources reconcile.TransitionResourcesFunc
	if transition != nil {
		transitionResources = func(original, desired resources.Resource) error {
			return transition(original.(*MockResource), desired.(*MockResource))
		}
	}
	return r.base.Reconcile(namespace, mockResourcesToResources(desiredResources), transitionResources, opts)
}
