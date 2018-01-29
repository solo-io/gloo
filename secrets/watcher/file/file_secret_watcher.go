package file

import (
	"fmt"
	"io/ioutil"
	"time"

	"github.com/ghodss/yaml"

	filewatch "github.com/solo-io/glue/adapters/file/watcher"
	"github.com/solo-io/glue/secrets/watcher"
)

// FileWatcher uses .yml files in a directory
// to watch secrets
type fileWatcher struct {
	secrets chan watcher.SecretMap
	errors  chan error
}

func NewFileWatcher(file string, syncFrequency time.Duration) (*fileWatcher, error) {
	secrets := make(chan watcher.SecretMap)
	errors := make(chan error)
	if err := filewatch.WatchFile(file, func(path string) {
		secretMap, err := getSecrets(path)
		if err != nil {
			errors <- err
			return
		}
		secrets <- secretMap
	}, syncFrequency); err != nil {
		return nil, fmt.Errorf("failed to start filewatcher: %v", err)
	}

	return &fileWatcher{
		secrets: secrets,
		errors:  errors,
	}, nil
}

func (fc *fileWatcher) Config() <-chan watcher.SecretMap {
	return fc.secrets
}

func (fc *fileWatcher) Error() <-chan error {
	return fc.errors
}

func getSecrets(path string) (watcher.SecretMap, error) {
	var secretMap watcher.SecretMap
	yml, err := ioutil.ReadFile(path)
	if err != nil {
		return secretMap, err
	}
	err = yaml.Unmarshal(yml, &secretMap)

	return secretMap, err
}
