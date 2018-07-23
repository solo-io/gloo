package mocks

import (
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/factory"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/kube/crd"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime"
)

func (r *MockResource) SetStatus(status core.Status) {
	r.Status = status
}

func (r *MockResource) SetMetadata(meta core.Metadata) {
	r.Metadata = meta
}

var _ resources.Resource = &MockResource{}

type MockResourceClient interface {
	Register() error
	Read(name string, opts clients.ReadOpts) (*MockResource, error)
	Write(resource *MockResource, opts clients.WriteOpts) (*MockResource, error)
	Delete(name string, opts clients.DeleteOpts) error
	List(opts clients.ListOpts) ([]*MockResource, error)
	Watch(opts clients.WatchOpts) (<-chan []*MockResource, <-chan error, error)
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

func (client *mockResourceClient) Read(name string, opts clients.ReadOpts) (*MockResource, error) {
	resource, err := client.rc.Read(name, opts)
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

func (client *mockResourceClient) Delete(name string, opts clients.DeleteOpts) error {
	return client.rc.Delete(name, opts)
}

func (client *mockResourceClient) List(opts clients.ListOpts) ([]*MockResource, error) {
	resourceList, err := client.rc.List(opts)
	if err != nil {
		return nil, err
	}
	return convertResources(resourceList), nil
}

func (client *mockResourceClient) Watch(opts clients.WatchOpts) (<-chan []*MockResource, <-chan error, error) {
	resourcesChan, errs, initErr := client.rc.Watch(opts)
	if initErr != nil {
		return nil, nil, initErr
	}
	mockResourcesChan := make(chan []*MockResource)
	go func() {
		for {
			select {
			case resourceList := <-resourcesChan:
				mockResourcesChan <- convertResources(resourceList)
			}
		}
	}()
	return mockResourcesChan, errs, nil
}

func convertResources(resources []resources.Resource) []*MockResource {
	var mockResourceList []*MockResource
	for _, resource := range resources {
		mockResourceList = append(mockResourceList, resource.(*MockResource))
	}
	return mockResourceList
}

// Kubernetes Adapter for MockResource

type MockResourceCrd struct {
	resources.Resource
}

func (m *MockResourceCrd) GetObjectKind() schema.ObjectKind {
	t := MockResourceCrdDefinition.TypeMeta()
	return &t
}

func (m *MockResourceCrd) DeepCopyObject() runtime.Object {
	return &MockResourceCrd{
		Resource: resources.Clone(m.Resource),
	}
}

var MockResourceCrdDefinition = crd.NewCrd("testing.solo.io",
	"mocks",
	"testing.solo.io",
	"v1",
	"MockResource",
	"mk",
	&MockResourceCrd{})
