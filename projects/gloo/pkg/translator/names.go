package translator

import (
	"github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1"
)

func routeConfigName(listener *v1.Listener) string {
	return listener.Name + "-routes"
}
