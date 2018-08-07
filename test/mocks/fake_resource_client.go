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
func NewFakeResource(namespace, name string) *FakeResource {
	return &FakeResource{
		Metadata: core.Metadata{
			Name:      name,
			Namespace: namespace,
		},
	}
}

func (r *FakeResource) SetStatus(status core.Status) {
	r.Status = status
}

func (r *FakeResource) SetMetadata(meta core.Metadata) {
	r.Metadata = meta
}

var _ resources.Resource = &FakeResource{}

type FakeResourceClient interface {
	Register() error
	Read(namespace, name string, opts clients.ReadOpts) (*FakeResource, error)
	Write(resource *FakeResource, opts clients.WriteOpts) (*FakeResource, error)
	Delete(namespace, name string, opts clients.DeleteOpts) error
	List(namespace string, opts clients.ListOpts) ([]*FakeResource, error)
	Watch(namespace string, opts clients.WatchOpts) (<-chan []*FakeResource, <-chan error, error)
}

type fakeResourceClient struct {
	rc clients.ResourceClient
}

func NewFakeResourceClient(rcFactory factory.ResourceClientFactory) (FakeResourceClient, error) {
	rc, err := rcFactory.NewResourceClient(factory.NewResourceClientParams{
		ResourceType: &FakeResource{},
	})
	if err != nil {
		return nil, errors.Wrapf(err, "creating base FakeResource resource client")
	}
	return &fakeResourceClient{
		rc: rc,
	}, nil
}

func (client *fakeResourceClient) Register() error {
	return client.rc.Register()
}

func (client *fakeResourceClient) Read(namespace, name string, opts clients.ReadOpts) (*FakeResource, error) {
	opts = opts.WithDefaults()
	resource, err := client.rc.Read(namespace, name, opts)
	if err != nil {
		return nil, err
	}
	return resource.(*FakeResource), nil
}

func (client *fakeResourceClient) Write(fakeResource *FakeResource, opts clients.WriteOpts) (*FakeResource, error) {
	opts = opts.WithDefaults()
	resource, err := client.rc.Write(fakeResource, opts)
	if err != nil {
		return nil, err
	}
	return resource.(*FakeResource), nil
}

func (client *fakeResourceClient) Delete(namespace, name string, opts clients.DeleteOpts) error {
	opts = opts.WithDefaults()
	return client.rc.Delete(namespace, name, opts)
}

func (client *fakeResourceClient) List(namespace string, opts clients.ListOpts) ([]*FakeResource, error) {
	opts = opts.WithDefaults()
	resourceList, err := client.rc.List(namespace, opts)
	if err != nil {
		return nil, err
	}
	return convertToFakeResource(resourceList), nil
}

func (client *fakeResourceClient) Watch(namespace string, opts clients.WatchOpts) (<-chan []*FakeResource, <-chan error, error) {
	opts = opts.WithDefaults()
	resourcesChan, errs, initErr := client.rc.Watch(namespace, opts)
	if initErr != nil {
		return nil, nil, initErr
	}
	fakeResourcesChan := make(chan []*FakeResource)
	go func() {
		for {
			select {
			case resourceList := <-resourcesChan:
				fakeResourcesChan <- convertToFakeResource(resourceList)
			case <-opts.Ctx.Done():
				close(fakeResourcesChan)
				return
			}
		}
	}()
	return fakeResourcesChan, errs, nil
}

func convertToFakeResource(resources []resources.Resource) []*FakeResource {
	var fakeResourceList []*FakeResource
	for _, resource := range resources {
		fakeResourceList = append(fakeResourceList, resource.(*FakeResource))
	}
	return fakeResourceList
}

// Kubernetes Adapter for FakeResource

func (o *FakeResource) GetObjectKind() schema.ObjectKind {
	t := FakeResourceCrd.TypeMeta()
	return &t
}

func (o *FakeResource) DeepCopyObject() runtime.Object {
	return resources.Clone(o).(*FakeResource)
}

var FakeResourceCrd = crd.NewCrd("testing.solo.io",
	"fakes",
	"testing.solo.io",
	"v1",
	"FakeResource",
	"fk",
	&FakeResource{})
