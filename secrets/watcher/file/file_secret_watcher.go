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
	file           string
	secretsToWatch []string
	secrets        chan watcher.SecretMap
	errors         chan error
}

func NewFileSecretWatcher(file string, syncFrequency time.Duration) (*fileWatcher, error) {
	secrets := make(chan watcher.SecretMap)
	errors := make(chan error)
	fw := &fileWatcher{
		secrets: secrets,
		errors:  errors,
		file:    file,
	}
	if err := filewatch.WatchFile(file, func(_ string) {
		secretMap, err := fw.getSecrets()
		if err != nil {
			errors <- err
			return
		}
		secrets <- secretMap
	}, syncFrequency); err != nil {
		return nil, fmt.Errorf("failed to start filewatcher: %v", err)
	}

	return fw, nil
}

// triggers an update
func (fw *fileWatcher) UpdateRefs(secretRefs []string) {
	fw.secretsToWatch = secretRefs
	secretMap, err := fw.getSecrets()
	if err != nil {
		fw.errors <- err
		return
	}
	fw.secrets <- secretMap
}

func (fw *fileWatcher) Secrets() <-chan watcher.SecretMap {
	return fw.secrets
}

func (fw *fileWatcher) Error() <-chan error {
	return fw.errors
}

func (fw *fileWatcher) getSecrets() (watcher.SecretMap, error) {
	yml, err := ioutil.ReadFile(fw.file)
	if err != nil {
		return nil, err
	}
	var secretMap watcher.SecretMap
	err = yaml.Unmarshal(yml, &secretMap)
	if err != nil {
		return nil, err
	}
	desiredSecrets := make(watcher.SecretMap)
	for _, ref := range fw.secretsToWatch {
		data, ok := secretMap[ref]
		if !ok {
			return nil, fmt.Errorf("secret ref %v not found in file %v", ref, fw.file)
		}
		desiredSecrets[ref] = data
	}

	return desiredSecrets, err
}
