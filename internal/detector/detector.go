package detector

import (
	"time"

	"sync"

	"github.com/pkg/errors"
	"github.com/solo-io/gloo-api/pkg/api/types/v1"
	"github.com/solo-io/gloo-function-discovery/pkg/resolver"
	"github.com/solo-io/gloo-plugins/kubernetes"
	"github.com/solo-io/gloo/pkg/coreplugins/service"
)

// detectors detect a specific type of functional service
// if they detect the service, they return service info and
// annotations (optional) for the service
type Detector interface {
	// if it detects the upstream is a known functional type, give us the
	// service info and annotations to mark it with
	DetectFunctionalService(addr string) (*v1.ServiceInfo, map[string]string, error)
}

// marker marks the upstream as functional. this modifies the upstream it was received,
// so should not be called concurrently from multiple goroutines
type Marker struct {
	detectors []Detector
	resolver  *resolver.Resolver
}

func NewMarker(detectors []Detector, resolver *resolver.Resolver) *Marker {
	return &Marker{
		detectors: detectors,
		resolver:  resolver,
	}
}

// should only be called for k8s, consul, and service type upstreams
func (m *Marker) MarkFunctionalUpstream(us *v1.Upstream) error {
	if us.Type != kubernetes.UpstreamTypeKube && us.Type != service.UpstreamTypeService {
		// don't run detection for these types of upstreams
		return nil
	}
	if us.ServiceInfo != nil {
		// this upstream has already been marked, skip it
		return nil
	}

	stop := make(chan struct{})
	wg := sync.WaitGroup{}
	// try every possible detector concurrently
	for _, d := range m.detectors {
		addr, err := m.resolver.Resolve(us)
		if err != nil {
			return errors.Wrapf(err, "resolving address for %v", us.Name)
		}
		wg.Add(1)
		go func() {
			withBackoff(func() error {
				serviceInfo, annotations, err := d.DetectFunctionalService(addr)
				if err != nil {
					return err
				}
				// discovered an upstream
				us.ServiceInfo = serviceInfo
				if us.Metadata == nil {
					us.Metadata = &v1.Metadata{}
				}
				us.Metadata.Annotations = mergeAnnotations(us.Metadata.Annotations, annotations)
				// stop the other detectors from running for this upstream
				close(stop)
				return nil
			}, stop)
			wg.Done()
		}()
	}
	wg.Wait()
	return nil
}

// Default values for ExponentialBackOff.
const (
	defaultInitialInterval = 500 * time.Millisecond
	defaultMaxElapsedTime  = 3 * time.Minute
)

func withBackoff(fn func() error, stop chan struct{}) {
	// first try
	if err := fn(); err == nil {
		return
	}
	tilNextRetry := defaultInitialInterval
	for {
		select {
		// stopped by another goroutine
		case <-stop:
			return
		case <-time.After(tilNextRetry):
			tilNextRetry *= 2
			if err := fn(); err == nil || tilNextRetry >= defaultMaxElapsedTime {
				return
			}
		}
	}
}

// get the unique set of funcs between two lists
// if conflict, new wins
func mergeAnnotations(oldAnnotations, newAnnotations map[string]string) map[string]string {
	merged := make(map[string]string)
	for k, v := range oldAnnotations {
		merged[k] = v
	}
	for k, v := range newAnnotations {
		merged[k] = v
	}
	return merged
}
