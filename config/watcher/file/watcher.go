package file

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/ghodss/yaml"
	"github.com/golang/protobuf/jsonpb"
	"github.com/radovskyb/watcher"
	"github.com/solo-io/glue/pkg/api/types"
	"github.com/solo-io/glue/pkg/log"
)

// FileWatcher uses .yml files in a directory
// to watch for config changes
type fileWatcher struct {
	configs chan *types.Config
	errors  chan error
}

func NewFileWatcher(dir string, syncFrequency time.Duration) (*fileWatcher, error) {
	configs := make(chan *types.Config)
	errors := make(chan error)
	w := watcher.New()
	w.SetMaxEvents(1)
	// Only notify rename and move events.
	w.FilterOps(watcher.Create, watcher.Write, watcher.Remove)
	go func() {
		for {
			select {
			case event := <-w.Event:
				log.Debugf("FileCache: Watcher received new event: %v", event)
				if event.IsDir() {
					break
				}
				cfg, err := parseConfig(event.Path)
				if err != nil {
					errors <- err
					break
				}
				configs <- &cfg
			case err := <-w.Error:
				log.Printf("FileCache: Watcher encountered error: %v", err)
			case <-w.Closed:
				log.Printf("FileCache: Watcher terminated")
				return
			}
		}
	}()

	// Watch this folder for changes.
	if err := w.AddRecursive(dir); err != nil {
		return nil, fmt.Errorf("failed to add watcher to %s: %v", dir, err)
	}

	// Print a list of all of the files and folders currently
	// being watched and their paths.
	for path, f := range w.WatchedFiles() {
		log.Printf("FileCache: Watching %s: %s\n", path, f.Name())
	}

	go func() {
		if err := w.Start(syncFrequency); err != nil {
			errors <- fmt.Errorf("failed to start watcher to: %v", err)
		}
	}()

	return &fileWatcher{
		configs: configs,
		errors:  errors,
	}, nil
}

func (fc *fileWatcher) Config() <-chan *types.Config {
	return fc.configs
}

func (fc *fileWatcher) Error() <-chan error {
	return fc.errors
}

func parseConfig(path string) (types.Config, error) {
	var cfg types.Config
	yml, err := ioutil.ReadFile(path)
	if err != nil {
		return cfg, err
	}
	jsn, err := yaml.YAMLToJSON(yml)
	if err != nil {
		return cfg, err
	}
	err = jsonpb.Unmarshal(bytes.NewBuffer(jsn), &cfg)
	if err != nil {
		log.GreyPrintf("WHY!\n%s\n%v", jsn, err)
	}

	return cfg, err
}
