package TODO

import (
	"github.com/solo-io/solo-kit/projects/gloo/pkg/defaults"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1"
)

const SqoopSidecarBindAddr = "127.0.0.1"
const SqoopSidecarBindPort = 9090
const SqoopSidecarName = "sqoop-sidecar"
const SqoopNamespace = defaults.GlooSystem

// TODO(ilackarms): make this a global function
func TransitionFunction(original, desired *v1.Proxy) (bool, error) {
	if len(original.Listeners) != len(desired.Listeners) {
		return true, nil
	}
	for i := range original.Listeners {
		if !original.Listeners[i].Equal(desired.Listeners[i]) {
			return true, nil
		}
	}
	return false, nil
}
