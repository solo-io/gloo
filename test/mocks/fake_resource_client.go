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

func (r *FakeResource) SetStatus(status core.Status) {
	r.Status = status
}

func (r *FakeResource) SetMetadata(meta core.Metadata) {
	r.Metadata = meta
}

var _ resources.Resource = &FakeResource{}

type FakeResourceClient interface {
	Register() error
	Read(name string, opts clients.ReadOpts) (*FakeResource, error)
	Write(resource *FakeResource, opts clients.WriteOpts) (*FakeResource, error)
	Delete(name string, opts clients.DeleteOpts) error
	List(opts clients.ListOpts) ([]*FakeResource, error)
	Watch(opts clients.WatchOpts) (<-chan []*FakeResource, <-chan error, error)
}

type fakeResourceClient struct {
	rc clients.ResourceClient
}

func NewFakeResourceClient(factory *factory.ResourceClientFactory) FakeResourceClient {
	return &fakeResourceClient{
		rc: factory.NewResourceClient(&FakeResource{}),
	}
}

func (client *fakeResourceClient) Register() error {
	return client.rc.Register()
}

func (client *fakeResourceClient) Read(name string, opts clients.ReadOpts) (*FakeResource, error) {
	resource, err := client.rc.Read(name, opts)
	if err != nil {
		return nil, err
	}
	return resource.(*FakeResource), nil
}

func (client *fakeResourceClient) Write(fakeResource *FakeResource, opts clients.WriteOpts) (*FakeResource, error) {
	resource, err := client.rc.Write(fakeResource, opts)
	if err != nil {
		return nil, err
	}
	return resource.(*FakeResource), nil
}

func (client *fakeResourceClient) Delete(name string, opts clients.DeleteOpts) error {
	return client.rc.Delete(name, opts)
}

func (client *fakeResourceClient) List(opts clients.ListOpts) ([]*FakeResource, error) {
	resourceList, err := client.rc.List(opts)
	if err != nil {
		return nil, err
	}
	return convertResources(resourceList), nil
}

func (client *fakeResourceClient) Watch(opts clients.WatchOpts) (<-chan []*FakeResource, <-chan error, error) {
	resourcesChan, errs, initErr := client.rc.Watch(opts)
	if initErr != nil {
		return nil, nil, initErr
	}
	fakeResourcesChan := make(chan []*FakeResource)
	go func() {
		for {
			select {
			case resourceList := <-resourcesChan:
				fakeResourcesChan <- convertResources(resourceList)
			}
		}
	}()
	return fakeResourcesChan, errs, nil
}

func convertResources(resources []resources.Resource) []*FakeResource {
	var fakeResourceList []*FakeResource
	for _, resource := range resources {
		fakeResourceList = append(fakeResourceList, resource.(*FakeResource))
	}
	return fakeResourceList
}

// Kubernetes Adapter for FakeResource

type FakeResourceCrd struct {
	resources.Resource
}

func (m *FakeResourceCrd) GetObjectKind() schema.ObjectKind {
	t := FakeResourceCrdDefinition.TypeMeta()
	return &t
}

func (m *FakeResourceCrd) DeepCopyObject() runtime.Object {
	return &FakeResourceCrd{
		Resource: resources.Clone(m.Resource),
	}
}

var FakeResourceCrdDefinition = crd.NewCrd("testing.solo.io",
	"fakes",
	"testing.solo.io",
	"v1",
	"FakeResource",
	"fk",
	&FakeResourceCrd{})
