package mocks

import (
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/factory"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/solo-io/solo-kit/pkg/errors"
)

type FakeResourceClient interface {
	BaseClient() clients.ResourceClient
	Register() error
	Read(namespace, name string, opts clients.ReadOpts) (*FakeResource, error)
	Write(resource *FakeResource, opts clients.WriteOpts) (*FakeResource, error)
	Delete(namespace, name string, opts clients.DeleteOpts) error
	List(namespace string, opts clients.ListOpts) (FakeResourceList, error)
	Watch(namespace string, opts clients.WatchOpts) (<-chan FakeResourceList, <-chan error, error)
}

type fakeResourceClient struct {
	rc clients.ResourceClient
}

func NewFakeResourceClient(rcFactory factory.ResourceClientFactory) (FakeResourceClient, error) {
	return NewFakeResourceClientWithToken(rcFactory, "")
}

func NewFakeResourceClientWithToken(rcFactory factory.ResourceClientFactory, token string) (FakeResourceClient, error) {
	rc, err := rcFactory.NewResourceClient(factory.NewResourceClientParams{
		ResourceType: &FakeResource{},
		Token:        token,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "creating base FakeResource resource client")
	}
	return &fakeResourceClient{
		rc: rc,
	}, nil
}

func (client *fakeResourceClient) BaseClient() clients.ResourceClient {
	return client.rc
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

func (client *fakeResourceClient) List(namespace string, opts clients.ListOpts) (FakeResourceList, error) {
	opts = opts.WithDefaults()
	resourceList, err := client.rc.List(namespace, opts)
	if err != nil {
		return nil, err
	}
	return convertToFakeResource(resourceList), nil
}

func (client *fakeResourceClient) Watch(namespace string, opts clients.WatchOpts) (<-chan FakeResourceList, <-chan error, error) {
	opts = opts.WithDefaults()
	resourcesChan, errs, initErr := client.rc.Watch(namespace, opts)
	if initErr != nil {
		return nil, nil, initErr
	}
	fakesChan := make(chan FakeResourceList)
	go func() {
		for {
			select {
			case resourceList := <-resourcesChan:
				fakesChan <- convertToFakeResource(resourceList)
			case <-opts.Ctx.Done():
				close(fakesChan)
				return
			}
		}
	}()
	return fakesChan, errs, nil
}

func convertToFakeResource(resources resources.ResourceList) FakeResourceList {
	var fakeResourceList FakeResourceList
	for _, resource := range resources {
		fakeResourceList = append(fakeResourceList, resource.(*FakeResource))
	}
	return fakeResourceList
}
