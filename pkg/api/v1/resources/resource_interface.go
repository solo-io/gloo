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
type InputResourceList []InputResource
type DataResourceList []DataResource

// namespace is optional, if left empty, names can collide if the list contains more than one with the same name
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

func (list DataResourceList) Find(namespace, name string) (DataResource, error) {
	for _, resource := range list {
		if resource.GetMetadata().Name == name {
			if namespace == "" || resource.GetMetadata().Namespace == namespace {
				return resource, nil
			}
		}
	}
	return nil, errors.Errorf("list did not find resource %v.%v", namespace, name)
}

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
