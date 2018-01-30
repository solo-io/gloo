package configwatcher

import (
	"fmt"
	"io/ioutil"
	"time"

	"github.com/ghodss/yaml"

	"github.com/solo-io/glue/implemented_modules/file/pkg/watcher"
	"github.com/solo-io/glue/pkg/api/types/v1"
)

// FileWatcher uses .yml files in a directory
// to watch for config changes
type fileWatcher struct {
	configs chan *v1.Config
	errors  chan error
}

func NewFileConfigWatcher(file string, syncFrequency time.Duration) (*fileWatcher, error) {
	configs := make(chan *v1.Config)
	errors := make(chan error)
	if err := watcher.WatchFile(file, func(path string) {
		cfg, err := parseConfig(path)
		if err != nil {
			errors <- err
			return
		}
		configs <- &cfg
	}, syncFrequency); err != nil {
		return nil, fmt.Errorf("failed to start filewatcher: %v", err)
	}

	return &fileWatcher{
		configs: configs,
		errors:  errors,
	}, nil
}

func (fc *fileWatcher) Config() <-chan *v1.Config {
	return fc.configs
}

func (fc *fileWatcher) Error() <-chan error {
	return fc.errors
}

func parseConfig(path string) (v1.Config, error) {
	var cfg v1.Config
	yml, err := ioutil.ReadFile(path)
	if err != nil {
		return cfg, err
	}
	err = yaml.Unmarshal(yml, &cfg)

	return cfg, err
}
