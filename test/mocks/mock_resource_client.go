package mocks

import (
	"sort"

	"github.com/gogo/protobuf/proto"
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

type MockResourceList []*MockResource
type MockResourceListsByNamespace map[string]MockResourceList

// namespace is optional, if left empty, names can collide if the list contains more than one with the same name
func (list MockResourceList) Find(namespace, name string) (*MockResource, error) {
	for _, mockResource := range list {
		if mockResource.Metadata.Name == name {
			if namespace == "" || mockResource.Metadata.Namespace == namespace {
				return mockResource, nil
			}
		}
	}
	return nil, errors.Errorf("list did not find mockResource %v.%v", namespace, name)
}

func (list MockResourceList) AsResources() resources.ResourceList {
	var ress resources.ResourceList
	for _, mockResource := range list {
		ress = append(ress, mockResource)
	}
	return ress
}

func (list MockResourceList) AsInputResources() resources.InputResourceList {
	var ress resources.InputResourceList
	for _, mockResource := range list {
		ress = append(ress, mockResource)
	}
	return ress
}

func (list MockResourceList) Names() []string {
	var names []string
	for _, mockResource := range list {
		names = append(names, mockResource.Metadata.Name)
	}
	return names
}

func (list MockResourceList) NamespacesDotNames() []string {
	var names []string
	for _, mockResource := range list {
		names = append(names, mockResource.Metadata.Namespace+"."+mockResource.Metadata.Name)
	}
	return names
}

func (list MockResourceList) Sort() {
	sort.SliceStable(list, func(i, j int) bool {
		return list[i].Metadata.Less(list[j].Metadata)
	})
}

func (list MockResourceList) Clone() MockResourceList {
	var mockResourceList MockResourceList
	for _, mockResource := range list {
		mockResourceList = append(mockResourceList, proto.Clone(mockResource).(*MockResource))
	}
	return mockResourceList
}

func (list MockResourceList) ByNamespace() MockResourceListsByNamespace {
	byNamespace := make(MockResourceListsByNamespace)
	for _, mockResource := range list {
		byNamespace.Add(mockResource)
	}
	return byNamespace
}

func (byNamespace MockResourceListsByNamespace) Add(mockResource ...*MockResource) {
	for _, item := range mockResource {
		byNamespace[item.Metadata.Namespace] = append(byNamespace[item.Metadata.Namespace], item)
	}
}

func (byNamespace MockResourceListsByNamespace) Clear(namespace string) {
	delete(byNamespace, namespace)
}

func (byNamespace MockResourceListsByNamespace) List() MockResourceList {
	var list MockResourceList
	for _, mockResourceList := range byNamespace {
		list = append(list, mockResourceList...)
	}
	list.Sort()
	return list
}

func (byNamespace MockResourceListsByNamespace) Clone() MockResourceListsByNamespace {
	return byNamespace.List().Clone().ByNamespace()
}

var _ resources.Resource = &MockResource{}

type MockResourceClient interface {
	BaseClient() clients.ResourceClient
	Register() error
	Read(namespace, name string, opts clients.ReadOpts) (*MockResource, error)
	Write(resource *MockResource, opts clients.WriteOpts) (*MockResource, error)
	Delete(namespace, name string, opts clients.DeleteOpts) error
	List(namespace string, opts clients.ListOpts) (MockResourceList, error)
	Watch(namespace string, opts clients.WatchOpts) (<-chan MockResourceList, <-chan error, error)
}

type mockResourceClient struct {
	rc clients.ResourceClient
}

func NewMockResourceClient(rcFactory factory.ResourceClientFactory) (MockResourceClient, error) {
	rc, err := rcFactory.NewResourceClient(factory.NewResourceClientParams{
		ResourceType: &MockResource{},
	})
	if err != nil {
		return nil, errors.Wrapf(err, "creating base MockResource resource client")
	}
	return &mockResourceClient{
		rc: rc,
	}, nil
}

func (client *mockResourceClient) BaseClient() clients.ResourceClient {
	return client.rc
}

func (client *mockResourceClient) Register() error {
	return client.rc.Register()
}

func (client *mockResourceClient) Read(namespace, name string, opts clients.ReadOpts) (*MockResource, error) {
	opts = opts.WithDefaults()
	resource, err := client.rc.Read(namespace, name, opts)
	if err != nil {
		return nil, err
	}
	return resource.(*MockResource), nil
}

func (client *mockResourceClient) Write(mockResource *MockResource, opts clients.WriteOpts) (*MockResource, error) {
	opts = opts.WithDefaults()
	resource, err := client.rc.Write(mockResource, opts)
	if err != nil {
		return nil, err
	}
	return resource.(*MockResource), nil
}

func (client *mockResourceClient) Delete(namespace, name string, opts clients.DeleteOpts) error {
	opts = opts.WithDefaults()
	return client.rc.Delete(namespace, name, opts)
}

func (client *mockResourceClient) List(namespace string, opts clients.ListOpts) (MockResourceList, error) {
	opts = opts.WithDefaults()
	resourceList, err := client.rc.List(namespace, opts)
	if err != nil {
		return nil, err
	}
	return convertToMockResource(resourceList), nil
}

func (client *mockResourceClient) Watch(namespace string, opts clients.WatchOpts) (<-chan MockResourceList, <-chan error, error) {
	opts = opts.WithDefaults()
	resourcesChan, errs, initErr := client.rc.Watch(namespace, opts)
	if initErr != nil {
		return nil, nil, initErr
	}
	mocksChan := make(chan MockResourceList)
	go func() {
		for {
			select {
			case resourceList := <-resourcesChan:
				mocksChan <- convertToMockResource(resourceList)
			case <-opts.Ctx.Done():
				close(mocksChan)
				return
			}
		}
	}()
	return mocksChan, errs, nil
}

func convertToMockResource(resources resources.ResourceList) MockResourceList {
	var mockResourceList MockResourceList
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
