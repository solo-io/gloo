package proxysvc

import (
	errors "github.com/rotisserie/eris"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
)

var (
	FailedToGetProxyError = func(err error, ref *core.ResourceRef) error {
		return errors.Wrapf(err, "Failed to get proxy %v.%v", ref.GetNamespace(), ref.GetName())
	}

	FailedToListProxiesError = func(err error, namespace string) error {
		return errors.Wrapf(err, "Failed to list proxies in %v", namespace)
	}
)
