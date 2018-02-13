package file

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ghodss/yaml"
	"github.com/gogo/protobuf/proto"
	"github.com/pkg/errors"
	"github.com/solo-io/glue/pkg/log"

	"github.com/solo-io/glue/internal/pkg/file"
	"github.com/solo-io/glue/pkg/api/types/v1"
	"github.com/solo-io/glue/pkg/protoutil"
)

var (
	upstreamDir     = "upstreams"
	virtualhostsDir = "virtualhosts"
	subdirs         = []string{upstreamDir, virtualhostsDir}
)

// FileWatcher uses .yml files in a directory
// to watch for config changes
type fileWatcher struct {
	configs chan *v1.Config
	errors  chan error
}

func NewFileConfigWatcher(dir string, syncFrequency time.Duration) (*fileWatcher, error) {
	configs := make(chan *v1.Config)
	errs := make(chan error)
	for _, subdir := range subdirs {
		os.MkdirAll(filepath.Join(dir, subdir), 0755)
	}
	if err := file.WatchDir(dir, true, func(string) {
		cfg, err := refreshConfig(dir)
		if err != nil {
			errs <- err
			return
		}
		configs <- cfg
	}, syncFrequency); err != nil {
		return nil, fmt.Errorf("failed to start filewatcher: %v", err)
	}

	// do one on start
	go func() {
		cfg, err := refreshConfig(dir)
		if err != nil {
			errs <- err
			return
		}
		configs <- cfg
	}()

	return &fileWatcher{
		configs: configs,
		errors:  errs,
	}, nil
}

func (fc *fileWatcher) Config() <-chan *v1.Config {
	return fc.configs
}

func (fc *fileWatcher) Error() <-chan error {
	return fc.errors
}

func refreshConfig(configDir string) (*v1.Config, error) {
	var (
		upstreams    []*v1.Upstream
		virtualHosts []*v1.VirtualHost
	)
	fullUpstreamDir := filepath.Join(configDir, upstreamDir)
	upstreamFiles, err := ioutil.ReadDir(fullUpstreamDir)
	if err != nil {
		return nil, errors.New("failed to read directory " + fullUpstreamDir)
	}
	fullVirtualhostDir := filepath.Join(configDir, virtualhostsDir)
	virtualhostFiles, err := ioutil.ReadDir(fullVirtualhostDir)
	if err != nil {
		return nil, errors.New("failed to read directory " + fullVirtualhostDir)
	}

	for _, f := range upstreamFiles {
		if f.IsDir() {
			continue
		}
		if !strings.HasSuffix(f.Name(), ".yml") && !strings.HasSuffix(f.Name(), ".yaml") {
			continue
		}
		var u v1.Upstream
		path := filepath.Join(fullUpstreamDir, f.Name())
		if err := readFileInto(path, &u); err != nil {
			return nil, errors.Errorf("failed to read file into upstream: %v", err)
		}
		upstreams = append(upstreams, &u)
	}

	for _, f := range virtualhostFiles {
		if f.IsDir() {
			continue
		}
		if !strings.HasSuffix(f.Name(), ".yml") && !strings.HasSuffix(f.Name(), ".yaml") {
			continue
		}
		var vh v1.VirtualHost
		path := filepath.Join(fullVirtualhostDir, f.Name())
		if err := readFileInto(path, &vh); err != nil {
			return nil, errors.Errorf("failed to read file into virtualhost: %v", err)
		}
		virtualHosts = append(virtualHosts, &vh)
	}

	return &v1.Config{
		VirtualHosts: virtualHosts,
		Upstreams:    upstreams,
	}, err
}

func readFileInto(f string, v proto.Message) error {
	data, err := ioutil.ReadFile(f)
	if err != nil {
		return errors.Errorf("error reading file: %v", err)
	}
	log.GreyPrintf("%s", string(data))
	jsn, err := yaml.YAMLToJSON(data)
	if err != nil {
		return err
	}
	return protoutil.Unmarshal(jsn, v)
}
