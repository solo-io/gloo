package mocks

import (
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/reconcile"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/solo-io/solo-kit/pkg/utils/contextutils"
)

// Option to copy anything from the original to the desired before writing
type TransitionMockDataFunc func(original, desired *MockData) error

type MockDataReconciler interface {
	Reconcile(namespace string, desiredResources []*MockData, transition TransitionMockDataFunc, opts clients.ListOpts) error
}

func mockDatasToResources(list MockDataList) resources.ResourceList {
	var resourceList resources.ResourceList
	for _, mockData := range list {
		resourceList = append(resourceList, mockData)
	}
	return resourceList
}

func NewMockDataReconciler(client MockDataClient) MockDataReconciler {
	return &mockDataReconciler{
		base: reconcile.NewReconciler(client.BaseClient()),
	}
}

type mockDataReconciler struct {
	base reconcile.Reconciler
}

func (r *mockDataReconciler) Reconcile(namespace string, desiredResources []*MockData, transition TransitionMockDataFunc, opts clients.ListOpts) error {
	opts = opts.WithDefaults()
	opts.Ctx = contextutils.WithLogger(opts.Ctx, "mockData_reconciler")
	var transitionResources reconcile.TransitionResourcesFunc
	if transition != nil {
		transitionResources = func(original, desired resources.Resource) error {
			return transition(original.(*MockData), desired.(*MockData))
		}
	}
	return r.base.Reconcile(namespace, mockDatasToResources(desiredResources), transitionResources, opts)
}
