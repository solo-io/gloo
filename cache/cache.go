package cache

import "github.com/solo-io/glue/pkg/api/types"

// Cache is responsible for reading and
// writing the glue config
type Cache interface {
	Config() <-chan types.Config
	// shows the status of the current config read by the cache
	// should show valid if the most recent update passed, otherwise a useful error
	Error() <-chan error
}
