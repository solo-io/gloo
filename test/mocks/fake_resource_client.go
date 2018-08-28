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

type FakeResourceList []*FakeResource
type FakeResourceListsByNamespace map[string]FakeResourceList

// namespace is optional, if left empty, names can collide if the list contains more than one with the same name
func (list FakeResourceList) Find(namespace, name string) (*FakeResource, error) {
	for _, fakeResource := range list {
		if fakeResource.Metadata.Name == name {
			if namespace == "" || fakeResource.Metadata.Namespace == namespace {
				return fakeResource, nil
			}
		}
	}
	return nil, errors.Errorf("list did not find fakeResource %v.%v", namespace, name)
}

func (list FakeResourceList) AsResources() resources.ResourceList {
	var ress resources.ResourceList
	for _, fakeResource := range list {
		ress = append(ress, fakeResource)
	}
	return ress
}

func (list FakeResourceList) AsInputResources() resources.InputResourceList {
	var ress resources.InputResourceList
	for _, fakeResource := range list {
		ress = append(ress, fakeResource)
	}
	return ress
}

func (list FakeResourceList) Names() []string {
	var names []string
	for _, fakeResource := range list {
		names = append(names, fakeResource.Metadata.Name)
	}
	return names
}

func (list FakeResourceList) NamespacesDotNames() []string {
	var names []string
	for _, fakeResource := range list {
		names = append(names, fakeResource.Metadata.Namespace+"."+fakeResource.Metadata.Name)
	}
	return names
}

func (list FakeResourceList) Sort() {
	sort.SliceStable(list, func(i, j int) bool {
		return list[i].Metadata.Less(list[j].Metadata)
	})
}

func (list FakeResourceList) Clone() FakeResourceList {
	var fakeResourceList FakeResourceList
	for _, fakeResource := range list {
		fakeResourceList = append(fakeResourceList, proto.Clone(fakeResource).(*FakeResource))
	}
	return fakeResourceList
}

func (list FakeResourceList) ByNamespace() FakeResourceListsByNamespace {
	byNamespace := make(FakeResourceListsByNamespace)
	for _, fakeResource := range list {
		byNamespace.Add(fakeResource)
	}
	return byNamespace
}

func (byNamespace FakeResourceListsByNamespace) Add(fakeResource ...*FakeResource) {
	for _, item := range fakeResource {
		byNamespace[item.Metadata.Namespace] = append(byNamespace[item.Metadata.Namespace], item)
	}
}

func (byNamespace FakeResourceListsByNamespace) Clear(namespace string) {
	delete(byNamespace, namespace)
}

func (byNamespace FakeResourceListsByNamespace) List() FakeResourceList {
	var list FakeResourceList
	for _, fakeResourceList := range byNamespace {
		list = append(list, fakeResourceList...)
	}
	list.Sort()
	return list
}

func (byNamespace FakeResourceListsByNamespace) Clone() FakeResourceListsByNamespace {
	return byNamespace.List().Clone().ByNamespace()
}

var _ resources.Resource = &FakeResource{}

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
