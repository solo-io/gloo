package mocks

import (
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/factory"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/kube/crd"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-kit/pkg/errors"
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

type MockDataList []*MockData

// namespace is optional, if left empty, names can collide if the list contains more than one with the same name
func (list MockDataList) Find(namespace, name string) (*MockData, error) {
	for _, mockData := range list {
		if mockData.Metadata.Name == name {
			if namespace == "" || mockData.Metadata.Namespace == namespace {
				return mockData, nil
			}
		}
	}
	return nil, errors.Errorf("list did not find mockData %v.%v", namespace, name)
}

var _ resources.Resource = &MockData{}

type MockDataClient interface {
	Register() error
	Read(namespace, name string, opts clients.ReadOpts) (*MockData, error)
	Write(resource *MockData, opts clients.WriteOpts) (*MockData, error)
	Delete(namespace, name string, opts clients.DeleteOpts) error
	List(namespace string, opts clients.ListOpts) (MockDataList, error)
	Watch(namespace string, opts clients.WatchOpts) (<-chan MockDataList, <-chan error, error)
}

type mockDataClient struct {
	rc clients.ResourceClient
}

func NewMockDataClient(rcFactory factory.ResourceClientFactory) (MockDataClient, error) {
	rc, err := rcFactory.NewResourceClient(factory.NewResourceClientParams{
		ResourceType: &MockData{},
	})
	if err != nil {
		return nil, errors.Wrapf(err, "creating base MockData resource client")
	}
	return &mockDataClient{
		rc: rc,
	}, nil
}

func (client *mockDataClient) Register() error {
	return client.rc.Register()
}

func (client *mockDataClient) Read(namespace, name string, opts clients.ReadOpts) (*MockData, error) {
	opts = opts.WithDefaults()
	resource, err := client.rc.Read(namespace, name, opts)
	if err != nil {
		return nil, err
	}
	return resource.(*MockData), nil
}

func (client *mockDataClient) Write(mockData *MockData, opts clients.WriteOpts) (*MockData, error) {
	opts = opts.WithDefaults()
	resource, err := client.rc.Write(mockData, opts)
	if err != nil {
		return nil, err
	}
	return resource.(*MockData), nil
}

func (client *mockDataClient) Delete(namespace, name string, opts clients.DeleteOpts) error {
	opts = opts.WithDefaults()
	return client.rc.Delete(namespace, name, opts)
}

func (client *mockDataClient) List(namespace string, opts clients.ListOpts) (MockDataList, error) {
	opts = opts.WithDefaults()
	resourceList, err := client.rc.List(namespace, opts)
	if err != nil {
		return nil, err
	}
	return convertToMockData(resourceList), nil
}

func (client *mockDataClient) Watch(namespace string, opts clients.WatchOpts) (<-chan MockDataList, <-chan error, error) {
	opts = opts.WithDefaults()
	resourcesChan, errs, initErr := client.rc.Watch(namespace, opts)
	if initErr != nil {
		return nil, nil, initErr
	}
	mockDatasChan := make(chan MockDataList)
	go func() {
		for {
			select {
			case resourceList := <-resourcesChan:
				mockDatasChan <- convertToMockData(resourceList)
			case <-opts.Ctx.Done():
				close(mockDatasChan)
				return
			}
		}
	}()
	return mockDatasChan, errs, nil
}

func convertToMockData(resources []resources.Resource) MockDataList {
	var mockDataList MockDataList
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
	"mockdatas",
	"testing.solo.io",
	"v1",
	"MockData",
	"mkd",
	&MockData{})
