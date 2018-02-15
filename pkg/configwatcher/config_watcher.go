package configwatcher

import "github.com/solo-io/glue/pkg/api/types/v1"

// ConfigWatcher reports new configs when they are updated externally
type Interface interface {
	// starts the ConfigWatcher
	Run(stop <-chan struct{})

	// configs are pushed here whenever the watcher detects an update
	Config() <-chan *v1.Config

	// if an error occurs, it should be pushed to this channel by the watcher
	Error() <-chan error
}
