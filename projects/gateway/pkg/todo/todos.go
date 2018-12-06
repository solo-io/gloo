package TODO

import (
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
)

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
