package errors

import (
	"fmt"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
)

type alreadyExistsErr struct {
	resource resources.Resource
}

func (err *alreadyExistsErr) Error() string {
	return fmt.Sprintf("already exists: %v", err.resource.GetMetadata())
}

func NewAlreadyExistsErr(resource resources.Resource) *alreadyExistsErr {
	return &alreadyExistsErr{resource: resource}
}

func IsAlreadyExists(err error) bool {
	switch err.(type) {
	case *alreadyExistsErr:
		return true
	}
	return false
}
