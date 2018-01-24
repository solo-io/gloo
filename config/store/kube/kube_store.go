package kube

import (
	"fmt"

	"github.com/solo-io/glue/config"
)

type ConfigCache interface {
	// SetUpdateHandler sets the handler function
	// to be called when the cached Config changes
	SetUpdateHandler(f func(raw []byte) error)

	// ApplyConfig is called when an API event occurs
	// which updates the config
	// SetUpdateHandler should naturally  be invoked
	// as a consequence (asychronously)
	ApplyConfig(cfg *config.Config) error
}

type kubeCache struct{}

func NewKubeCache() *kubeCache {
	return &kubeCache{}
}

// initialize crd & controller
func (cache *kubeCache) Init() error {
	return nil
}

func (cache *kubeCache) SetUpdateHandler(f func(raw []byte) error) {}

func (cache *kubeCache) ApplyConfig(cfg *config.Config) error {
	return fmt.Errorf("not implemented")
}
