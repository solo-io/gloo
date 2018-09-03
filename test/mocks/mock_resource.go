package mocks

import (
	"sort"

	"github.com/gogo/protobuf/proto"
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
type MocksByNamespace map[string]MockResourceList

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

func (list MockResourceList) Sort() MockResourceList {
	sort.SliceStable(list, func(i, j int) bool {
		return list[i].Metadata.Less(list[j].Metadata)
	})
	return list
}

func (list MockResourceList) Clone() MockResourceList {
	var mockResourceList MockResourceList
	for _, mockResource := range list {
		mockResourceList = append(mockResourceList, proto.Clone(mockResource).(*MockResource))
	}
	return mockResourceList
}

func (list MockResourceList) ByNamespace() MocksByNamespace {
	byNamespace := make(MocksByNamespace)
	for _, mockResource := range list {
		byNamespace.Add(mockResource)
	}
	return byNamespace
}

func (byNamespace MocksByNamespace) Add(mockResource ...*MockResource) {
	for _, item := range mockResource {
		byNamespace[item.Metadata.Namespace] = append(byNamespace[item.Metadata.Namespace], item)
	}
}

func (byNamespace MocksByNamespace) Clear(namespace string) {
	delete(byNamespace, namespace)
}

func (byNamespace MocksByNamespace) List() MockResourceList {
	var list MockResourceList
	for _, mockResourceList := range byNamespace {
		list = append(list, mockResourceList...)
	}
	return list.Sort()
}

func (byNamespace MocksByNamespace) Clone() MocksByNamespace {
	return byNamespace.List().Clone().ByNamespace()
}

var _ resources.Resource = &MockResource{}

// Kubernetes Adapter for MockResource

func (o *MockResource) GetObjectKind() schema.ObjectKind {
	t := MockResourceCrd.TypeMeta()
	return &t
}

func (o *MockResource) DeepCopyObject() runtime.Object {
	return resources.Clone(o).(*MockResource)
}

var MockResourceCrd = crd.NewCrd("mocks.api.v1",
	"mocks",
	"mocks.api.v1",
	"v1",
	"MockResource",
	"mk",
	&MockResource{})
