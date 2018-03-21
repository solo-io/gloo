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
// to watch files
type fileWatcher struct {
	dir          string
	filesToWatch []string
	files        chan filewatcher.Files
	errors       chan error
	lastSeen     uint64
}

func NewFileWatcher(dir string, syncFrequency time.Duration) (*fileWatcher, error) {
	os.MkdirAll(dir, 0755)
	files := make(chan filewatcher.Files)
	errs := make(chan error)
	fw := &fileWatcher{
		files:  files,
		errors: errs,
		dir:    dir,
	}
	if err := file.WatchDir(dir, false, func(_ string) {
		fw.updateFiles()
	}, syncFrequency); err != nil {
		return nil, fmt.Errorf("failed to start filewatcher: %v", err)
	}

	// do one on start
	go fw.updateFiles()

	return fw, nil
}

func (fw *fileWatcher) updateFiles() {
	Files, err := fw.getFiles()
	if err != nil {
		fw.errors <- err
		return
	}
	// ignore empty configs / no files to watch
	if len(Files) == 0 {
		return
	}
	fw.files <- Files
}

// triggers an update
func (fw *fileWatcher) TrackFiles(fileRefs []string) {
	fw.filesToWatch = fileRefs
	fw.updateFiles()
}

func (fw *fileWatcher) Files() <-chan filewatcher.Files {
	return fw.files
}

func (fw *fileWatcher) Error() <-chan error {
	return fw.errors
}

func (fw *fileWatcher) getFiles() (filewatcher.Files, error) {
	desiredFiles := make(filewatcher.Files)
	// ref should be the filename
	for _, ref := range fw.filesToWatch {
		data, err := ioutil.ReadFile(filepath.Join(fw.dir, ref))
		if err != nil {
			return nil, errors.Wrapf(err, "reading file: %v", filepath.Join(fw.dir, ref))
		}
		// because, on the filesystem,
		desiredFiles[ref] = data
	}

	hash, err := hashstructure.Hash(desiredFiles, nil)
	if err != nil {
		runtime.HandleError(err)
		return nil, nil
	}
	if fw.lastSeen == hash {
		return nil, nil
	}
	fw.lastSeen = hash
	return desiredFiles, nil
}
