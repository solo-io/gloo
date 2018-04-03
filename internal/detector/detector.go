package detector

import (
	"time"

	"sync"

	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
	"github.com/solo-io/gloo-api/pkg/api/types/v1"
	"github.com/solo-io/gloo-function-discovery/pkg/resolver"
	"github.com/solo-io/gloo-plugins/kubernetes"
	"github.com/solo-io/gloo/pkg/coreplugins/service"
)

// detectors detect a specific type of functional service
// if they detect the service, they return service info and
// annotations (optional) for the service
type Interface interface {
	// if it detects the upstream is a known functional type, give us the
	// service info and annotations to mark it with
	DetectFunctionalService(addr string) (*v1.ServiceInfo, map[string]string, error)
}

// marker marks the upstream as functional. this modifies the upstream it was received,
// so should not be called concurrently from multiple goroutines
type Marker struct {
	detectors []Interface
	resolver  *resolver.Resolver

	markedOrFailed map[string]bool
	m              sync.RWMutex
}

func NewMarker(detectors []Interface, resolver *resolver.Resolver) *Marker {
	return &Marker{
		detectors:      detectors,
		resolver:       resolver,
		markedOrFailed: make(map[string]bool),
	}
}

// should only be called for k8s, consul, and service type upstreams
func (m *Marker) DetectFunctionalUpstream(us *v1.Upstream) (*v1.ServiceInfo, map[string]string, error) {
	if us.Type != kubernetes.UpstreamTypeKube && us.Type != service.UpstreamTypeService {
		// don't run detection for these types of upstreams
		return nil, nil, nil
	}
	if us.ServiceInfo != nil {
		return nil, nil, nil
		// this upstream has already been marked, skip it
	}

	m.m.RLock()
	// already tried this upstream
	already := m.markedOrFailed[us.Name]
	m.m.RUnlock()
	if already {
		return nil, nil, nil
	}

	m.m.Lock()
	m.markedOrFailed[us.Name] = true
	m.m.Unlock()

	addr, err := m.resolver.Resolve(us)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "resolving address for %v", us.Name)
	}

	stop := make(chan struct{})
	failed := make(chan error)

	serviceInfoC := make(chan *v1.ServiceInfo)
	annotationsC := make(chan map[string]string)

	// try every possible detector concurrently
	for _, d := range m.detectors {
		go func(d Interface) {
			err := withBackoff(func() error {
				serviceInfo, annotations, err := d.DetectFunctionalService(addr)
				if err != nil {
					return err
				}
				// success
				close(stop)
				serviceInfoC <- serviceInfo
				annotationsC <- annotations
				return nil
			}, stop)
			if err != nil {
				failed <- err
			}
		}(d)
	}

	var totalFailed int

	var errs error
	for {
		select {
		case <-stop:
			return <-serviceInfoC, <-annotationsC, nil
		case err := <-failed:
			errs = multierror.Append(errs, err)
			totalFailed++
			if totalFailed >= len(m.detectors) {
				return nil, nil, errors.Errorf("service type detection failed for %s: %v", us.Name, errs)
			}
		}
	}
}

// Default values for ExponentialBackOff.
const (
	defaultInitialInterval = 500 * time.Millisecond
	defaultMaxElapsedTime  = 60 * time.Second
)

func withBackoff(fn func() error, stop chan struct{}) error {
	// first try
	if err := fn(); err == nil {
		return nil
	}
	tilNextRetry := defaultInitialInterval
	for {
		select {
		// stopped by another goroutine
		case <-stop:
			return nil
		case <-time.After(tilNextRetry):
			tilNextRetry *= 2
			err := fn()
			if err == nil {
				return nil
			}
			if tilNextRetry >= defaultMaxElapsedTime {
				return err
			}
		}
	}
}
