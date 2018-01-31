package configwatcher

import "github.com/solo-io/glue/pkg/api/types/v1"

// ConfigWatcher reports new configs when they are updated externally
type Interface interface {
	// configs are pushed here whenever they are read
	Config() <-chan *v1.Config
	// shows the status of the current config read by the cache
	// should show valid if the most recent update passed, otherwise a useful error
	Error() <-chan error
}
