package resources

import (
	"reflect"
	"sort"

	"github.com/gogo/protobuf/proto"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-kit/pkg/errors"
	"k8s.io/apimachinery/pkg/util/validation"
)

type Resource interface {
	proto.Message
	GetMetadata() core.Metadata
	SetMetadata(meta core.Metadata)
	Equal(that interface{}) bool
}

type ResourcesByType map[string]ResourceList

// mixed type resource list
func (m ResourcesByType) List() ResourceList {
	var all ResourceList
	for _, list := range m {
		all = append(all, list...)
	}
	// sort by type
	sort.SliceStable(all, func(i, j int) bool {
		if Kind(all[i]) < Kind(all[j]) {
			return true
		}
		return all[i].GetMetadata().Less(all[j].GetMetadata())
	})
	return all
}

type InputResource interface {
	Resource
	GetStatus() core.Status
	SetStatus(status core.Status)
}

type DataResource interface {
	Resource
	GetData() map[string]string
	SetData(map[string]string)
}

type ResourceList []Resource
type ResourcesByKind map[string]ResourceList

func (m ResourcesByKind) Add(resource Resource) {
	m[Kind(resource)] = append(m[Kind(resource)], resource)
}
func (m ResourcesByKind) Get(resource Resource) []Resource {
	return m[Kind(resource)]
}
func (list ResourceList) Copy() ResourceList {
	var cpy ResourceList
	for _, res := range list {
		cpy = append(cpy, Clone(res))
	}
	return cpy
}
func (list ResourceList) Contains(list2 ResourceList) bool {
	for _, res := range list {
		var found bool
		for _, res2 := range list {
			if res.GetMetadata().Name == res2.GetMetadata().Name && res.GetMetadata().Namespace == res2.GetMetadata().Namespace {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}
func (list ResourceList) Find(namespace, name string) (Resource, error) {
	for _, resource := range list {
		if resource.GetMetadata().Name == name {
			if namespace == "" || resource.GetMetadata().Namespace == namespace {
				return resource, nil
			}
		}
	}
	return nil, errors.Errorf("list did not find resource %v.%v", namespace, name)
}
func (list ResourceList) Names() []string {
	var names []string
	for _, resource := range list {
		names = append(names, resource.GetMetadata().Name)
	}
	return names
}
func (list ResourceList) Namespaces() []string {
	var namespaces []string
	for _, resource := range list {
		namespaces = append(namespaces, resource.GetMetadata().Namespace)
	}
	return namespaces
}
func (list ResourceList) FilterByNames(names []string) ResourceList {
	var filtered ResourceList
	for _, resource := range list {
		for _, name := range names {
			if name == resource.GetMetadata().Name {
				filtered = append(filtered, resource)
				break
			}
		}
	}
	return filtered
}
func (list ResourceList) FilterByNamespaces(namespaces []string) ResourceList {
	var filtered ResourceList
	for _, resource := range list {
		for _, namespace := range namespaces {
			if namespace == resource.GetMetadata().Namespace {
				filtered = append(filtered, resource)
				break
			}
		}
	}
	return filtered
}
func (list ResourceList) FilterByList(list2 ResourceList) ResourceList {
	return list.FilterByNamespaces(list2.Namespaces()).FilterByNames(list.Names())
}
func (list ResourceList) AsInputResourceList() InputResourceList {
	var inputs InputResourceList
	for _, res := range list {
		inputRes, ok := res.(InputResource)
		if !ok {
			continue
		}
		inputs = append(inputs, inputRes)
	}
	return inputs
}

type InputResourceList []InputResource
type InputResourcesByKind map[string]InputResourceList

func (m InputResourcesByKind) Add(resource InputResource) {
	m[Kind(resource)] = append(m[Kind(resource)], resource)
}
func (m InputResourcesByKind) Get(resource InputResource) []InputResource {
	return m[Kind(resource)]
}
func (list InputResourceList) Copy() InputResourceList {
	var cpy InputResourceList
	for _, res := range list {
		cpy = append(cpy, Clone(res).(InputResource))
	}
	return cpy
}
func (list InputResourceList) Contains(list2 InputResourceList) bool {
	for _, res := range list {
		var found bool
		for _, res2 := range list {
			if res.GetMetadata().Name == res2.GetMetadata().Name && res.GetMetadata().Namespace == res2.GetMetadata().Namespace {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}
func (list InputResourceList) Find(namespace, name string) (InputResource, error) {
	for _, resource := range list {
		if resource.GetMetadata().Name == name {
			if namespace == "" || resource.GetMetadata().Namespace == namespace {
				return resource, nil
			}
		}
	}
	return nil, errors.Errorf("list did not find resource %v.%v", namespace, name)
}
func (list InputResourceList) Names() []string {
	var names []string
	for _, resource := range list {
		names = append(names, resource.GetMetadata().Name)
	}
	return names
}
func (list InputResourceList) Namespaces() []string {
	var namespaces []string
	for _, resource := range list {
		namespaces = append(namespaces, resource.GetMetadata().Namespace)
	}
	return namespaces
}
func (list InputResourceList) FilterByNames(names []string) InputResourceList {
	var filtered InputResourceList
	for _, resource := range list {
		for _, name := range names {
			if name == resource.GetMetadata().Name {
				filtered = append(filtered, resource)
				break
			}
		}
	}
	return filtered
}
func (list InputResourceList) FilterByNamespaces(namespaces []string) InputResourceList {
	var filtered InputResourceList
	for _, resource := range list {
		for _, namespace := range namespaces {
			if namespace == resource.GetMetadata().Namespace {
				filtered = append(filtered, resource)
				break
			}
		}
	}
	return filtered
}
func (list InputResourceList) FilterByList(list2 InputResourceList) InputResourceList {
	return list.FilterByNamespaces(list2.Namespaces()).FilterByNames(list.Names())
}
func (list InputResourceList) AsInputResourceList() ResourceList {
	var resources ResourceList
	for _, res := range list {
		resources = append(resources, res)
	}
	return resources
}

type DataResourceList []DataResource

func Clone(resource Resource) Resource {
	return proto.Clone(resource).(Resource)
}

func Kind(resource Resource) string {
	return reflect.TypeOf(resource).String()
}

func UpdateMetadata(resource Resource, updateFunc func(meta *core.Metadata)) {
	meta := resource.GetMetadata()
	updateFunc(&meta)
	resource.SetMetadata(meta)
}

func Validate(resource Resource) error {
	return ValidateName(resource.GetMetadata().Name)
}

func ValidateName(name string) error {
	errs := validation.IsDNS1035Label(name)
	if len(name) < 1 {
		errs = append(errs, "name cannot be empty")
	}
	if len(name) > 253 {
		errs = append(errs, "name has a max length of 253 characters")
	}
	if len(errs) > 0 {
		return errors.Errors(errs)
	}
	return nil
}
