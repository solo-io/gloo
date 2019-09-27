package routetablesvc

import (
	"github.com/pkg/errors"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
)

var (
	FailedToReadRouteTableError = func(err error, ref *core.ResourceRef) error {
		return errors.Wrapf(err, "Failed to read route table %v.%v", ref.GetNamespace(), ref.GetName())
	}

	FailedToListRouteTablesError = func(err error, namespace string) error {
		return errors.Wrapf(err, "Failed to list route tables in %v", namespace)
	}

	InvalidInputError = errors.New("No input provided")

	FailedToCreateRouteTableError = func(err error, namespace, name string) error {
		return errors.Wrapf(err, "Failed to create route table %v.%v", namespace, name)
	}

	FailedToUpdateRouteTableError = func(err error, namespace, name string) error {
		return errors.Wrapf(err, "Failed to update route table %v.%v", namespace, name)
	}

	FailedToParseRouteTableFromYamlError = func(err error, ref *core.ResourceRef) error {
		return errors.Wrapf(err, "Failed to parse route table %s.%s from YAML", ref.Namespace, ref.Name)
	}

	FailedToDeleteRouteTableError = func(err error, ref *core.ResourceRef) error {
		return errors.Wrapf(err, "Failed to delete route table %v.%v", ref.GetNamespace(), ref.GetName())
	}
)
