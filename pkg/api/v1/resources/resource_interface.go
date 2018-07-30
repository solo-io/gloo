package resources

import (
	"reflect"

	"github.com/gogo/protobuf/proto"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-kit/pkg/errors"
	"k8s.io/apimachinery/pkg/util/validation"
)

type BaseResource interface {
	GetMetadata() core.Metadata
	SetMetadata(meta core.Metadata)
}

type Resource interface {
	proto.Message
	BaseResource
	Equal(that interface{}) bool
	GetStatus() core.Status
	SetStatus(status core.Status)
}

func Clone(resource Resource) Resource {
	return proto.Clone(resource).(Resource)
}

func Kind(resource Resource) string {
	return reflect.TypeOf(resource).String()
}

func UpdateMetadata(resource BaseResource, updateFunc func(meta *core.Metadata)) {
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
