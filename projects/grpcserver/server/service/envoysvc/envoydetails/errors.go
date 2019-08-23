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

	GatewayProxyPodIsNotRunning = func(namespace, name, phase string) string {
		return fmt.Sprintf("Gateway proxy pod %v.%v is not running. Current phase: %v", namespace, name, phase)
	}

	ProxyResourceNotFound = func(name string) string {
		return fmt.Sprintf("Could not find gloo proxy resource for gateway-proxy %v", name)
	}

	ProxyResourceRejected = func(namespace, name, reason string) string {
		return fmt.Sprintf("Proxy resource %v.%v is rejected with reason: %v", namespace, name, reason)
	}

	ProxyResourcePending = func(namespace, name string) string {
		return fmt.Sprintf("Proxy resource %v.%v is pending", namespace, name)
	}
)
