package envoydetails

import (
	"fmt"

	"github.com/solo-io/go-utils/errors"
)

// Go errors
var (
	FailedToListPodsError = func(err error, namespace, selector string) error {
		return errors.Wrapf(err, "Failed to list pods in %v using LabelSelector %v", namespace, selector)
	}
)

// String error messages
var (
	FailedToGetEnvoyConfig = func(namespace, name string) string {
		return fmt.Sprintf("Failed to get envoy config from pod %v.%v", namespace, name)
	}
)
