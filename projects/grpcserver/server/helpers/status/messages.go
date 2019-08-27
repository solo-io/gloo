package status

import (
	"fmt"

	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
)

var (
	ResourcePending = func(namespace, name string) string {
		return fmt.Sprintf("Resource %v.%v is pending", namespace, name)
	}

	ResourceRejected = func(namespace, name, reason string) string {
		return fmt.Sprintf("Resource %v.%v is rejected with reason: %v", namespace, name, reason)
	}

	UnknownFailure = func(namespace, name string, status core.Status_State) string {
		return fmt.Sprintf("Resource %v.%v has an unknown status: %v", namespace, name, status)
	}
)
