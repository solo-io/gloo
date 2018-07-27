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
func NewMockResource(namespace, name string) *MockResource {
	return &MockResource{
		Metadata: core.Metadata{
			Name:      name,
			Namespace: namespace,
		},
	}
}

func (r *MockResource) SetStatus(status core.Status) {
	r.Status = status
}

func (r *MockResource) SetMetadata(meta core.Metadata) {
	r.Metadata = meta
}

var _ resources.Resource = &MockResource{}

type MockResourceClient interface {
	Register() error
	Read(namespace, name string, opts clients.ReadOpts) (*MockResource, error)
	Write(resource *MockResource, opts clients.WriteOpts) (*MockResource, error)
	Delete(namespace, name string, opts clients.DeleteOpts) error
	List(namespace string, opts clients.ListOpts) ([]*MockResource, error)
	Watch(namespace string, opts clients.WatchOpts) (<-chan []*MockResource, <-chan error, error)
}

type mockResourceClient struct {
	rc clients.ResourceClient
}

func NewMockResourceClient(factory *factory.ResourceClientFactory) MockResourceClient {
	return &mockResourceClient{
		rc: factory.NewResourceClient(&MockResource{}),
	}
}

func (client *mockResourceClient) Register() error {
	return client.rc.Register()
}

func (client *mockResourceClient) Read(namespace, name string, opts clients.ReadOpts) (*MockResource, error) {
	resource, err := client.rc.Read(namespace, name, opts)
	if err != nil {
		return nil, err
	}
	return resource.(*MockResource), nil
}

func (client *mockResourceClient) Write(mockResource *MockResource, opts clients.WriteOpts) (*MockResource, error) {
	resource, err := client.rc.Write(mockResource, opts)
	if err != nil {
		return nil, err
	}
	return resource.(*MockResource), nil
}

func (client *mockResourceClient) Delete(namespace, name string, opts clients.DeleteOpts) error {
	return client.rc.Delete(namespace, name, opts)
}

func (client *mockResourceClient) List(namespace string, opts clients.ListOpts) ([]*MockResource, error) {
	resourceList, err := client.rc.List(namespace, opts)
	if err != nil {
		return nil, err
	}
	return convertToMockResource(resourceList), nil
}

func (client *mockResourceClient) Watch(namespace string, opts clients.WatchOpts) (<-chan []*MockResource, <-chan error, error) {
	resourcesChan, errs, initErr := client.rc.Watch(namespace, opts)
	if initErr != nil {
		return nil, nil, initErr
	}
	mockResourcesChan := make(chan []*MockResource)
	go func() {
		for {
			select {
			case resourceList := <-resourcesChan:
				mockResourcesChan <- convertToMockResource(resourceList)
			}
		}
	}()
	return mockResourcesChan, errs, nil
}

func convertToMockResource(resources []resources.Resource) []*MockResource {
	var mockResourceList []*MockResource
	for _, resource := range resources {
		mockResourceList = append(mockResourceList, resource.(*MockResource))
	}
	return mockResourceList
}

// Kubernetes Adapter for MockResource

func (o *MockResource) GetObjectKind() schema.ObjectKind {
	t := MockResourceCrd.TypeMeta()
	return &t
}

func (o *MockResource) DeepCopyObject() runtime.Object {
	return resources.Clone(o).(*MockResource)
}

var MockResourceCrd = crd.NewCrd("testing.solo.io",
	"mocks",
	"testing.solo.io",
	"v1",
	"MockResource",
	"mk",
	&MockResource{})
