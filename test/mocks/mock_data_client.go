package mocks

import (
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/factory"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/kube/crd"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// TODO: modify as needed to populate additional fields
func NewMockData(namespace, name string) *MockData {
	return &MockData{
		Metadata: core.Metadata{
			Name:      name,
			Namespace: namespace,
		},
	}
}

func (r *MockData) SetStatus(status core.Status) {
	r.Status = status
}

func (r *MockData) SetMetadata(meta core.Metadata) {
	r.Metadata = meta
}

func (r *MockData) SetData(data map[string]string) {
	r.Data = data
}

var _ resources.Resource = &MockData{}

type MockDataClient interface {
	Register() error
	Read(namespace, name string, opts clients.ReadOpts) (*MockData, error)
	Write(resource *MockData, opts clients.WriteOpts) (*MockData, error)
	Delete(namespace, name string, opts clients.DeleteOpts) error
	List(namespace string, opts clients.ListOpts) ([]*MockData, error)
	Watch(namespace string, opts clients.WatchOpts) (<-chan []*MockData, <-chan error, error)
}

type mockDataClient struct {
	rc clients.ResourceClient
}

func NewMockDataClient(factory *factory.ResourceClientFactory) MockDataClient {
	return &mockDataClient{
		rc: factory.NewResourceClient(&MockData{}),
	}
}

func (client *mockDataClient) Register() error {
	return client.rc.Register()
}

func (client *mockDataClient) Read(namespace, name string, opts clients.ReadOpts) (*MockData, error) {
	resource, err := client.rc.Read(namespace, name, opts)
	if err != nil {
		return nil, err
	}
	return resource.(*MockData), nil
}

func (client *mockDataClient) Write(mockData *MockData, opts clients.WriteOpts) (*MockData, error) {
	resource, err := client.rc.Write(mockData, opts)
	if err != nil {
		return nil, err
	}
	return resource.(*MockData), nil
}

func (client *mockDataClient) Delete(namespace, name string, opts clients.DeleteOpts) error {
	return client.rc.Delete(namespace, name, opts)
}

func (client *mockDataClient) List(namespace string, opts clients.ListOpts) ([]*MockData, error) {
	resourceList, err := client.rc.List(namespace, opts)
	if err != nil {
		return nil, err
	}
	return convertToMockData(resourceList), nil
}

func (client *mockDataClient) Watch(namespace string, opts clients.WatchOpts) (<-chan []*MockData, <-chan error, error) {
	resourcesChan, errs, initErr := client.rc.Watch(namespace, opts)
	if initErr != nil {
		return nil, nil, initErr
	}
	mockDatasChan := make(chan []*MockData)
	go func() {
		for {
			select {
			case resourceList := <-resourcesChan:
				mockDatasChan <- convertToMockData(resourceList)
			}
		}
	}()
	return mockDatasChan, errs, nil
}

func convertToMockData(resources []resources.Resource) []*MockData {
	var mockDataList []*MockData
	for _, resource := range resources {
		mockDataList = append(mockDataList, resource.(*MockData))
	}
	return mockDataList
}

// Kubernetes Adapter for MockData

func (o *MockData) GetObjectKind() schema.ObjectKind {
	t := MockDataCrd.TypeMeta()
	return &t
}

func (o *MockData) DeepCopyObject() runtime.Object {
	return resources.Clone(o).(*MockData)
}

var MockDataCrd = crd.NewCrd("testing.solo.io",
	"fakes",
	"testing.solo.io",
	"v1",
	"MockData",
	"fk",
	&MockData{})
