package static

import (
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
)

type Namespacer struct {
	namespaces []string
}

func (n *Namespacer) Namespaces(opts clients.WatchOpts) (<-chan []string, <-chan error, error) {
	namespacesC := make(chan []string)
	errC := make(chan error)
	go func() {
		namespacesC <- n.namespaces
	}()
	return namespacesC, errC, nil
}

func NewNamespacer(namespaces []string) *Namespacer {
	return &Namespacer{
		namespaces: namespaces,
	}
}
