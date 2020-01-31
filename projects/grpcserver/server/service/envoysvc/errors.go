package envoysvc

import errors "github.com/rotisserie/eris"

var (
	FailedToListEnvoyDetailsError = func(err error) error {
		return errors.Wrapf(err, "Failed to list envoy details")
	}
)
