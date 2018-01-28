package watcher

import "github.com/solo-io/glue/pkg/api/types"

// Watcher is responsible for reading the glue config
type Watcher interface {
	// configs are pushed here whenever they are read
	Config() <-chan *types.Config
	// shows the status of the current config read by the cache
	// should show valid if the most recent update passed, otherwise a useful error
	Error() <-chan error
}
