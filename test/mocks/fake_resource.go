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
type FakesByNamespace map[string]FakeResourceList

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

func (list FakeResourceList) Sort() FakeResourceList {
	sort.SliceStable(list, func(i, j int) bool {
		return list[i].Metadata.Less(list[j].Metadata)
	})
	return list
}

func (list FakeResourceList) Clone() FakeResourceList {
	var fakeResourceList FakeResourceList
	for _, fakeResource := range list {
		fakeResourceList = append(fakeResourceList, proto.Clone(fakeResource).(*FakeResource))
	}
	return fakeResourceList
}

func (list FakeResourceList) ByNamespace() FakesByNamespace {
	byNamespace := make(FakesByNamespace)
	for _, fakeResource := range list {
		byNamespace.Add(fakeResource)
	}
	return byNamespace
}

func (byNamespace FakesByNamespace) Add(fakeResource ...*FakeResource) {
	for _, item := range fakeResource {
		byNamespace[item.Metadata.Namespace] = append(byNamespace[item.Metadata.Namespace], item)
	}
}

func (byNamespace FakesByNamespace) Clear(namespace string) {
	delete(byNamespace, namespace)
}

func (byNamespace FakesByNamespace) List() FakeResourceList {
	var list FakeResourceList
	for _, fakeResourceList := range byNamespace {
		list = append(list, fakeResourceList...)
	}
	return list.Sort()
}

func (byNamespace FakesByNamespace) Clone() FakesByNamespace {
	return byNamespace.List().Clone().ByNamespace()
}

var _ resources.Resource = &FakeResource{}

// Kubernetes Adapter for FakeResource

func (o *FakeResource) GetObjectKind() schema.ObjectKind {
	t := FakeResourceCrd.TypeMeta()
	return &t
}

func (o *FakeResource) DeepCopyObject() runtime.Object {
	return resources.Clone(o).(*FakeResource)
}

var FakeResourceCrd = crd.NewCrd("mocks.api.v1",
	"fakes",
	"mocks.api.v1",
	"v1",
	"FakeResource",
	"fk",
	&FakeResource{})
