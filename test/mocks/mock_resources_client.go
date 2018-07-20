package mocks

import (
	"time"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/factory"
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

type typedResourceClient struct {
	rc clients.ResourceClient
}

func NewMockResourceClient(factory *factory.ResourceClientFactory) MockResourceClient {
	return &typedResourceClient{
		rc: factory.NewResourceClient(&MockResource{}),
	}
}

func (client *typedResourceClient) Register() error {
	return client.rc.Register()
}

func (client *typedResourceClient) Read(name string, opts clients.ReadOpts) (*MockResource, error) {
	resource, err := client.rc.Read(name, opts)
	if err != nil {
		return nil, err
	}
	return resource.(*MockResource), nil
}

func (client *typedResourceClient) Write(typedResource *MockResource, opts clients.WriteOpts) (*MockResource, error) {
	resource, err := client.rc.Write(typedResource, opts)
	if err != nil {
		return nil, err
	}
	return resource.(*MockResource), nil
}

func (client *typedResourceClient) Delete(name string, opts clients.DeleteOpts) error {
	return client.rc.Delete(name, opts)
}

func (client *typedResourceClient) List(opts clients.ListOpts) ([]*MockResource, error) {
	resourceList, err := client.rc.List(opts)
	if err != nil {
		return nil, err
	}
	return convertResources(resourceList), nil
}

func (client *typedResourceClient) Watch(opts clients.WatchOpts) (<-chan []*MockResource, <-chan error, error) {
	resourcesChan, errs, initErr := client.rc.Watch(opts)
	if initErr != nil {
		return nil, nil, initErr
	}
	typedResourcesChan := make(chan []*MockResource)
	go func(){
		for {
			select {
			case resourceList := <- resourcesChan:
				typedResourcesChan <- convertResources(resourceList)
			}
		}
	}()
	return typedResourcesChan, errs, nil
}

func convertResources(resources []resources.Resource) []*MockResource {
	var typedResourceList []*MockResource
	for _, resource := range resources {
		typedResourceList = append(typedResourceList, resource.(*MockResource))
	}
	return typedResourceList
}