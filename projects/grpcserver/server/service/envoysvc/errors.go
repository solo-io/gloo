package envoysvc

import "github.com/solo-io/go-utils/errors"

var (
	FailedToListEnvoyDetailsError = func(err error) error {
		return errors.Wrapf(err, "Failed to list envoy details")
	}
)
