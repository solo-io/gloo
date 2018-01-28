package crd

import "github.com/solo-io/glue/pkg/api/types"

// CrdWatcher uses Kubernetes CRDs
// to watch for config changes
type crdWatcher struct {
	configs chan *types.Config
	errors  chan error
}

func NewCrdWatcher(masterUrl, kubeconfigPath string) (*crdWatcher, error) {

}
