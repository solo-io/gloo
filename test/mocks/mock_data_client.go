package mocks

import (
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/factory"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/solo-io/solo-kit/pkg/errors"
)

type MockDataClient interface {
	BaseClient() clients.ResourceClient
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
	return NewMockDataClientWithToken(rcFactory, "")
}

func NewMockDataClientWithToken(rcFactory factory.ResourceClientFactory, token string) (MockDataClient, error) {
	rc, err := rcFactory.NewResourceClient(factory.NewResourceClientParams{
		ResourceType: &MockData{},
		Token:        token,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "creating base MockData resource client")
	}
	return &mockDataClient{
		rc: rc,
	}, nil
}

func (client *mockDataClient) BaseClient() clients.ResourceClient {
	return client.rc
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

func convertToMockData(resources resources.ResourceList) MockDataList {
	var mockDataList MockDataList
	for _, resource := range resources {
		mockDataList = append(mockDataList, resource.(*MockData))
	}
	return mockDataList
}
