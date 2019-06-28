package artifactsvc

import (
	"github.com/pkg/errors"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
)

var (
	FailedToReadArtifactError = func(err error, ref *core.ResourceRef) error {
		return errors.Wrapf(err, "Failed to read artifact %v.%v", ref.GetNamespace(), ref.GetName())
	}

	FailedToListArtifactsError = func(err error, namespace string) error {
		return errors.Wrapf(err, "Failed to list artifacts in %v", namespace)
	}

	FailedToCreateArtifactError = func(err error, ref *core.ResourceRef) error {
		return errors.Wrapf(err, "Failed to create artifact %v.%v", ref.GetNamespace(), ref.GetName())
	}

	FailedToUpdateArtifactError = func(err error, ref *core.ResourceRef) error {
		return errors.Wrapf(err, "Failed to update artifact %v.%v", ref.GetNamespace(), ref.GetName())
	}

	FailedToDeleteArtifactError = func(err error, ref *core.ResourceRef) error {
		return errors.Wrapf(err, "Failed to delete artifact %v.%v", ref.GetNamespace(), ref.GetName())
	}
)
