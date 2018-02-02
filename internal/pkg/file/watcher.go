package file

import (
	"fmt"
	"time"

	"github.com/radovskyb/watcher"
	"k8s.io/apimachinery/pkg/util/runtime"

	"github.com/solo-io/glue/pkg/log"
)

// WatchFile watches for changes
// on a file or directory
// and calls the given handler
// on that file
func WatchFile(path string, handler func(path string), syncFrequency time.Duration) error {
	w := watcher.New()
	w.SetMaxEvents(1)
	// Only notify rename and move events.
	w.FilterOps(watcher.Create, watcher.Move, watcher.Write, watcher.Remove)
	go func() {
		for {
			select {
			case event := <-w.Event:
				log.Debugf("FileWatcher: Watcher received new event: %v %v", event.Op.String(), event.Path)
				if path != event.Path {
					break
				}
				handler(event.Path)
			case err := <-w.Error:
				log.Debugf("FileWatcher: Watcher encountered error: %v", err)
			case <-w.Closed:
				log.Debugf("FileWatcher: Watcher terminated")
				return
			}
		}
	}()

	// Watch this file for changes.
	if err := w.Add(path); err != nil {
		return fmt.Errorf("failed to add watcher to %s: %v", path, err)
	}

	// Print a list of all of the files and folders currently
	// being watched and their paths.
	for path, f := range w.WatchedFiles() {
		log.Debugf("FileWatcher: Watching %s: %s\n", path, f.Name())
	}

	go func() {
		if err := w.Start(syncFrequency); err != nil {
			runtime.HandleError(fmt.Errorf("failed to start watcher to: %v", err))
		}
	}()

	time.Sleep(time.Second)

	return nil
}
