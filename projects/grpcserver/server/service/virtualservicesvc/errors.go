package virtualservicesvc

import (
	"github.com/pkg/errors"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
)

var (
	FailedToReadVirtualServiceError = func(err error, ref *core.ResourceRef) error {
		return errors.Wrapf(err, "Failed to read virtual service %v.%v", ref.GetNamespace(), ref.GetName())
	}

	FailedToListVirtualServicesError = func(err error, namespace string) error {
		return errors.Wrapf(err, "Failed to list virtual services in %v", namespace)
	}

	InvalidInputError = errors.New("No input provided")

	FailedToCreateVirtualServiceError = func(err error, ref *core.ResourceRef) error {
		return errors.Wrapf(err, "Failed to create virtual service %v.%v", ref.GetNamespace(), ref.GetName())
	}

	FailedToUpdateVirtualServiceError = func(err error, ref *core.ResourceRef) error {
		return errors.Wrapf(err, "Failed to update virtual service %v.%v", ref.GetNamespace(), ref.GetName())
	}

	FailedToDeleteVirtualServiceError = func(err error, ref *core.ResourceRef) error {
		return errors.Wrapf(err, "Failed to delete virtual service %v.%v", ref.GetNamespace(), ref.GetName())
	}

	FailedToCreateRouteError = func(err error, ref *core.ResourceRef) error {
		return errors.Wrapf(err, "Failed to create route on virtual service %v.%v", ref.GetNamespace(), ref.GetName())
	}

	FailedToUpdateRouteError = func(err error, ref *core.ResourceRef) error {
		return errors.Wrapf(err, "Failed to update route on virtual service %v.%v", ref.GetNamespace(), ref.GetName())
	}

	FailedToDeleteRouteError = func(err error, ref *core.ResourceRef) error {
		return errors.Wrapf(err, "Failed to delete route on virtual service %v.%v", ref.GetNamespace(), ref.GetName())
	}

	FailedToSwapRoutesError = func(err error, ref *core.ResourceRef) error {
		return errors.Wrapf(err, "Failed to swap routes on virtual service %v.%v", ref.GetNamespace(), ref.GetName())
	}

	FailedToShiftRoutesError = func(err error, ref *core.ResourceRef) error {
		return errors.Wrapf(err, "Failed to shift routes on virtual service %v.%v", ref.GetNamespace(), ref.GetName())
	}
)
