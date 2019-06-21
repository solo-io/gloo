package secretsvc

import (
	"github.com/pkg/errors"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
)

var (
	FailedToReadSecretError = func(err error, ref core.ResourceRef) error {
		return errors.Wrapf(err, "Failed to read newSecret %v.%v", ref.GetNamespace(), ref.GetName())
	}

	FailedToListSecretsError = func(err error, namespace string) error {
		return errors.Wrapf(err, "Failed to list secrets in %v", namespace)
	}

	FailedToCreateSecretError = func(err error, ref core.ResourceRef) error {
		return errors.Wrapf(err, "Failed to create newSecret %v.%v", ref.GetNamespace(), ref.GetName())
	}

	FailedToUpdateSecretError = func(err error, ref core.ResourceRef) error {
		return errors.Wrapf(err, "Failed to update newSecret %v.%v", ref.GetNamespace(), ref.GetName())
	}

	FailedToDeleteSecretError = func(err error, ref core.ResourceRef) error {
		return errors.Wrapf(err, "Failed to delete newSecret %v.%v", ref.GetNamespace(), ref.GetName())
	}
)
