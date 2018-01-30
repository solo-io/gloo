package discovery

import (
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/ghodss/yaml"

	"github.com/solo-io/glue/implemented_modules/file/pkg/watcher"
	"github.com/solo-io/glue/module"
	"github.com/solo-io/glue/pkg/api/types/v1"
)

// FileWatcher uses .yml files in a directory
// to watch secrets
type fileDiscovery struct {
	file      string
	upstreams []v1.Upstream
	endpoints chan module.EndpointGroups
	errors    chan error
}

func NewServiceDiscovery(file string, syncFrequency time.Duration) (*fileDiscovery, error) {
	endpoints := make(chan module.EndpointGroups)
	errors := make(chan error)
	fw := &fileDiscovery{
		endpoints: endpoints,
		errors:    errors,
		file:      file,
	}
	if err := watcher.WatchFile(file, func(_ string) {
		fw.updateEndpoints()
	}, syncFrequency); err != nil {
		return nil, fmt.Errorf("failed to start filewatcher: %v", err)
	}

	return fw, nil
}

func (fw *fileDiscovery) updateEndpoints() {
	endpointGroups, err := fw.getSecrets()
	if err != nil {
		fw.errors <- err
		return
	}
	// ignore empty groups / no upstreams to watch
	if len(endpointGroups) == 0 {
		return
	}
	fw.endpoints <- endpointGroups
}

// triggers an update
func (fw *fileDiscovery) TrackUpstreams(upstreams []v1.Upstream) {
	fw.upstreams = upstreams
	fw.updateEndpoints()
}

func (fw *fileDiscovery) Endpoints() <-chan module.EndpointGroups {
	return fw.endpoints
}

func (fw *fileDiscovery) Error() <-chan error {
	return fw.errors
}

func (fw *fileDiscovery) getSecrets() (module.EndpointGroups, error) {
	yml, err := ioutil.ReadFile(fw.file)
	if err != nil {
		return nil, err
	}
	var endpointGroups module.EndpointGroups
	err = yaml.Unmarshal(yml, &endpointGroups)
	if err != nil {
		return nil, err
	}
	out := make(module.EndpointGroups)
	// only add endpoints for upstreams we are supposed to watch
	for _, upstream := range fw.upstreams {
		endpoints, ok := endpointGroups[upstream.Name]
		if !ok {
			log.Printf("ref %v not found", upstream)
			continue
		}
		log.Printf("ref found: %v", upstream)
		out[upstream.Name] = endpoints
	}

	return out, err
}
