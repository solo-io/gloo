// Code generated by MockGen. DO NOT EDIT.
// Source: ./cluster_set.go

// Package mock_multicluster is a generated GoMock package.
package mock_multicluster

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	manager "sigs.k8s.io/controller-runtime/pkg/manager"
)

// MockClusterSet is a mock of ClusterSet interface.
type MockClusterSet struct {
	ctrl     *gomock.Controller
	recorder *MockClusterSetMockRecorder
}

// MockClusterSetMockRecorder is the mock recorder for MockClusterSet.
type MockClusterSetMockRecorder struct {
	mock *MockClusterSet
}

// NewMockClusterSet creates a new mock instance.
func NewMockClusterSet(ctrl *gomock.Controller) *MockClusterSet {
	mock := &MockClusterSet{ctrl: ctrl}
	mock.recorder = &MockClusterSetMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockClusterSet) EXPECT() *MockClusterSetMockRecorder {
	return m.recorder
}

// AddCluster mocks base method.
func (m *MockClusterSet) AddCluster(ctx context.Context, cluster string, mgr manager.Manager) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "AddCluster", ctx, cluster, mgr)
}

// AddCluster indicates an expected call of AddCluster.
func (mr *MockClusterSetMockRecorder) AddCluster(ctx, cluster, mgr interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddCluster", reflect.TypeOf((*MockClusterSet)(nil).AddCluster), ctx, cluster, mgr)
}

// Exists mocks base method.
func (m *MockClusterSet) Exists(cluster string) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Exists", cluster)
	ret0, _ := ret[0].(bool)
	return ret0
}

// Exists indicates an expected call of Exists.
func (mr *MockClusterSetMockRecorder) Exists(cluster interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Exists", reflect.TypeOf((*MockClusterSet)(nil).Exists), cluster)
}

// ListClusters mocks base method.
func (m *MockClusterSet) ListClusters() []string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListClusters")
	ret0, _ := ret[0].([]string)
	return ret0
}

// ListClusters indicates an expected call of ListClusters.
func (mr *MockClusterSetMockRecorder) ListClusters() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListClusters", reflect.TypeOf((*MockClusterSet)(nil).ListClusters))
}

// RemoveCluster mocks base method.
func (m *MockClusterSet) RemoveCluster(cluster string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "RemoveCluster", cluster)
}

// RemoveCluster indicates an expected call of RemoveCluster.
func (mr *MockClusterSetMockRecorder) RemoveCluster(cluster interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveCluster", reflect.TypeOf((*MockClusterSet)(nil).RemoveCluster), cluster)
}