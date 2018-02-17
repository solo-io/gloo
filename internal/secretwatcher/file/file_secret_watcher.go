package file

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"k8s.io/apimachinery/pkg/util/runtime"

	"github.com/ghodss/yaml"
	"github.com/mitchellh/hashstructure"

	"github.com/solo-io/gloo/internal/pkg/file"
	"github.com/solo-io/gloo/pkg/log"
	"github.com/solo-io/gloo/pkg/secretwatcher"
)

// FileWatcher uses .yml files in a directory
// to watch secrets
type fileWatcher struct {
	dir            string
	secretsToWatch []string
	secrets        chan secretwatcher.SecretMap
	errors         chan error
	lastSeen       uint64
}

func NewSecretWatcher(dir string, syncFrequency time.Duration) (*fileWatcher, error) {
	os.MkdirAll(dir, 0755)
	secrets := make(chan secretwatcher.SecretMap)
	errors := make(chan error)
	fw := &fileWatcher{
		secrets: secrets,
		errors:  errors,
		dir:     dir,
	}
	if err := file.WatchDir(dir, false, func(_ string) {
		fw.updateSecrets()
	}, syncFrequency); err != nil {
		return nil, fmt.Errorf("failed to start filewatcher: %v", err)
	}

	// do one on start
	go fw.updateSecrets()

	return fw, nil
}

func (fw *fileWatcher) updateSecrets() {
	secretMap, err := fw.getSecrets()
	if err != nil {
		fw.errors <- err
		return
	}
	// ignore empty configs / no secrets to watch
	if len(secretMap) == 0 {
		return
	}
	fw.secrets <- secretMap
}

// triggers an update
func (fw *fileWatcher) TrackSecrets(secretRefs []string) {
	fw.secretsToWatch = secretRefs
	fw.updateSecrets()
}

func (fw *fileWatcher) Secrets() <-chan secretwatcher.SecretMap {
	return fw.secrets
}

func (fw *fileWatcher) Error() <-chan error {
	return fw.errors
}

func (fw *fileWatcher) getSecrets() (secretwatcher.SecretMap, error) {
	secretFiles, err := ioutil.ReadDir(fw.dir)
	if err != nil {
		return nil, err
	}
	desiredSecrets := make(secretwatcher.SecretMap)
	for _, secretFile := range secretFiles {
		yml, err := ioutil.ReadFile(filepath.Join(fw.dir, secretFile.Name()))
		if err != nil {
			return nil, err
		}
		var secretMap secretwatcher.SecretMap
		err = yaml.Unmarshal(yml, &secretMap)
		if err != nil {
			return nil, err
		}
		for _, ref := range fw.secretsToWatch {
			data, ok := secretMap[ref]
			if !ok {
				log.Debugf("ref %v not found", ref)
				return nil, fmt.Errorf("secret ref %v not found in dir %v", ref, fw.dir)
			}
			log.Debugf("ref found: %v", ref)
			desiredSecrets[ref] = data
		}
	}

	hash, err := hashstructure.Hash(desiredSecrets, nil)
	if err != nil {
		runtime.HandleError(err)
		return nil, nil
	}
	if fw.lastSeen == hash {
		return nil, nil
	}
	fw.lastSeen = hash
	return desiredSecrets, nil
}
