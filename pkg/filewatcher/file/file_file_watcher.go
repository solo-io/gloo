package file

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"k8s.io/apimachinery/pkg/util/runtime"

	"github.com/mitchellh/hashstructure"

	"github.com/pkg/errors"
	"github.com/solo-io/gloo/internal/pkg/file"
	"github.com/solo-io/gloo/pkg/filewatcher"
)

// FileWatcher uses .yml files in a directory
// to watch artifacts
type fileWatcher struct {
	dir              string
	artifactsToWatch []string
	artifacts        chan filewatcher.Files
	errors           chan error
	lastSeen         uint64
}

func NewArtifactWatcher(dir string, syncFrequency time.Duration) (*fileWatcher, error) {
	os.MkdirAll(dir, 0755)
	artifacts := make(chan filewatcher.Files)
	errs := make(chan error)
	fw := &fileWatcher{
		artifacts: artifacts,
		errors:    errs,
		dir:       dir,
	}
	if err := file.WatchDir(dir, false, func(_ string) {
		fw.updateArtifacts()
	}, syncFrequency); err != nil {
		return nil, fmt.Errorf("failed to start filewatcher: %v", err)
	}

	// do one on start
	go fw.updateArtifacts()

	return fw, nil
}

func (fw *fileWatcher) updateArtifacts() {
	Artifacts, err := fw.getArtifacts()
	if err != nil {
		fw.errors <- err
		return
	}
	// ignore empty configs / no artifacts to watch
	if len(Artifacts) == 0 {
		return
	}
	fw.artifacts <- Artifacts
}

// triggers an update
func (fw *fileWatcher) TrackArtifacts(artifactRefs []string) {
	fw.artifactsToWatch = artifactRefs
	fw.updateArtifacts()
}

func (fw *fileWatcher) Artifacts() <-chan filewatcher.Files {
	return fw.artifacts
}

func (fw *fileWatcher) Error() <-chan error {
	return fw.errors
}

func (fw *fileWatcher) getArtifacts() (filewatcher.Files, error) {
	desiredArtifacts := make(filewatcher.Files)
	// ref should be the filename
	for _, ref := range fw.artifactsToWatch {
		data, err := ioutil.ReadFile(filepath.Join(fw.dir, ref))
		if err != nil {
			return nil, errors.Wrapf(err, "reading file: %v", filepath.Join(fw.dir, ref))
		}
		// because, on the filesystem,
		desiredArtifacts[ref] = map[string][]byte{ref: data}
	}

	hash, err := hashstructure.Hash(desiredArtifacts, nil)
	if err != nil {
		runtime.HandleError(err)
		return nil, nil
	}
	if fw.lastSeen == hash {
		return nil, nil
	}
	fw.lastSeen = hash
	return desiredArtifacts, nil
}
