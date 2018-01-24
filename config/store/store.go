package store

import "github.com/solo-io/glue/config"

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
