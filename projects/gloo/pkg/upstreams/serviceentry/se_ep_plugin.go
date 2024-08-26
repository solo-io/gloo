package serviceentry

import (
	"github.com/solo-io/gloo/projects/gloo/pkg/discovery"

	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
)

var _ discovery.DiscoveryPlugin = &sePlugin{}

func (s *sePlugin) WatchEndpoints(writeNamespace string, upstreamsToTrack v1.UpstreamList, opts clients.WatchOpts) (<-chan v1.EndpointList, <-chan error, error) {
	// kubeUpstreams := slices.Filter(upstreamsToTrack, func(e *v1.Upstream) bool {
	// 	return e.GetKube() != nil
	// })
	panic("todo")
}

// satisfy interface but don't implement; only implement hybrid client for ServiceEntry Upstreams.
func (s *sePlugin) DiscoverUpstreams(watchNamespaces []string, writeNamespace string, opts clients.WatchOpts, discOpts discovery.Opts) (chan v1.UpstreamList, chan error, error) {
	us, errs := make(chan v1.UpstreamList), make(chan error)
	close(us)
	close(errs)
	return us, errs, nil
}

// satisfy interface but don't implement; only implement hybrid client for ServiceEntry Upstreams.
func (s *sePlugin) UpdateUpstream(original *v1.Upstream, desired *v1.Upstream) (bool, error) {
	return false, nil
}
