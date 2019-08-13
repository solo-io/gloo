package gatewaysvc

import (
	"github.com/solo-io/go-utils/errors"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
)

var (
	FailedToGetGatewayError = func(err error, ref *core.ResourceRef) error {
		return errors.Wrapf(err, "Failed to get gateway %v.%v", ref.GetNamespace(), ref.GetName())
	}

	FailedToListGatewaysError = func(err error, namespace string) error {
		return errors.Wrapf(err, "Failed to list gateways in %v", namespace)
	}
)
