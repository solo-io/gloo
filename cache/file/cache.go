package file

import (
	"fmt"
	"time"

	"github.com/radovskyb/watcher"
	"github.com/solo-io/glue/cache"
	"github.com/solo-io/glue/pkg/api/types"
	"github.com/solo-io/glue/pkg/log"
	"gopkg.in/yaml.v2"
)

const (
	routeDirName       = "routes"
	upstreamDirName    = "upstreams"
	virtualhostDirName = "virtualhosts"
)

// FileCache uses .yml files in a directory
// to write and read configs
type fileCache struct {
	dir           string
	syncFrequency time.Duration
	handler       cache.ConfigChangeHandler
	err           error
}

func NewFileCache(dir string, syncFrequency time.Duration) (*fileCache, error) {
	fc := &fileCache{
		dir:           dir,
		syncFrequency: syncFrequency,
	}
	w := watcher.New()
	w.SetMaxEvents(1)
	// Only notify rename and move events.
	w.FilterOps(watcher.Create, watcher.Write, watcher.Remove)
	go func(fc *fileCache) {
		for {
			select {
			case event := <-w.Event:
				log.Printf("FileCache: Watcher received new event: %v", event)
			case err := <-w.Error:
				log.Printf("FileCache: Watcher encountered error: %v", err)
			case <-w.Closed:
				return
			}
		}
	}(fc)
	return fc, nil
}

func (fc *fileCache) Config() <-chan types.Config {
	return fc.
}

func (fc *fileCache) Error() <-chan error {}

func parseConfig(raw []byte) (types.Config, error) {
	var cfg types.Config
	err := yaml.Unmarshal(raw, &cfg)
	return cfg, err
}
