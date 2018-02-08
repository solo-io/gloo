package file

import (
	"fmt"
	"time"

	"github.com/radovskyb/watcher"

	"github.com/solo-io/glue/pkg/log"
)

// WatchFile watches for changes
// on a file or directory
// and calls the given handler
// on that file
func WatchDir(path string, recursive bool, handler func(path string), syncFrequency time.Duration) error {
	w := watcher.New()
	w.SetMaxEvents(1)
	// Only notify rename and move events.
	w.FilterOps(watcher.Create, watcher.Move, watcher.Write, watcher.Remove)
	go func() {
		for {
			select {
			case event := <-w.Event:
				log.Debugf("FileWatcher: Watcher received new event: %v %v", event.Op.String(), event.Path)
				//if path != event.Path {
				//	break
				//}
				handler(event.Path)
			case err := <-w.Error:
				log.Debugf("FileWatcher: Watcher encountered error: %v", err)
			case <-w.Closed:
				log.Debugf("FileWatcher: Watcher terminated")
				return
			}
		}
	}()

	// Watch this file or directory for changes.
	if recursive {
		if err := w.AddRecursive(path); err != nil {
			return fmt.Errorf("failed to add watcher to %s: %v", path, err)
		}
	} else {
		if err := w.Add(path); err != nil {
			return fmt.Errorf("failed to add watcher to %s: %v", path, err)
		}
	}

	// Print a list of all of the files and folders currently
	// being watched and their paths.
	for path, f := range w.WatchedFiles() {
		log.Debugf("FileWatcher: Watching %s: %s\n", path, f.Name())
	}

	errC := make(chan error)
	go func() {
		if err := w.Start(syncFrequency); err != nil {
			errC <- fmt.Errorf("failed to start watcher to: %v", err)
		}
	}()

	select {
	case <-time.After(time.Second):
		return nil
	case err := <-errC:
		return err
	}
}
