package namespacing

import (
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
)

type Namespacer interface {
	Namespaces(opts clients.WatchOpts) (<-chan []string, <-chan error, error)
}
