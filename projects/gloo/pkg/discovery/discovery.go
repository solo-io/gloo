package discovery

import (
	"reflect"
	"sort"
	"sync"

	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/errors"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/plugins"
)

type DiscoveryPlugin interface {
	plugins.Plugin

	// UDS API
	// send us an updated list of upstreams on every change
	// namespace is for writing to, not necessarily reading from
	WatchUpstreams(namespace string, opts clients.WatchOpts, discOpts Opts) (chan v1.UpstreamList, chan error, error)
	// finalize any changes to the desired upstream before it gets written
	// for example, copying the functions from the old upstream to the new.
	// a value of false indicates that the resource does not need to be updated
	UpdateUpstream(original, desired *v1.Upstream) (bool, error)

	// EDS API
	// start the EDS watch which sends a new list of endpoints on any change
	// will send only endpoints for upstreams configured with TrackUpstreams
	WatchEndpoints(writeNamespace string, upstreamsToTrack v1.UpstreamList, opts clients.WatchOpts) (<-chan v1.EndpointList, <-chan error, error)
}

type Discovery struct {
	writeNamespace     string
	upstreamReconciler v1.UpstreamReconciler
	endpointReconciler v1.EndpointReconciler
	discoveryPlugins   []DiscoveryPlugin
}

func NewDiscovery(writeNamespace string,
	upstreamClient v1.UpstreamClient,
	endpointsClient v1.EndpointClient,
	discoveryPlugins []DiscoveryPlugin) *Discovery {
	return &Discovery{
		writeNamespace:     writeNamespace,
		upstreamReconciler: v1.NewUpstreamReconciler(upstreamClient),
		endpointReconciler: v1.NewEndpointReconciler(endpointsClient),
		discoveryPlugins:   discoveryPlugins,
	}
}

// launch a goroutine for all the UDS plugins
func (d *Discovery) StartUds(opts clients.WatchOpts, discOpts Opts) (chan error, error) {
	aggregatedErrs := make(chan error)
	upstreamsByUds := make(map[DiscoveryPlugin]v1.UpstreamList)
	lock := sync.Mutex{}
	for _, uds := range d.discoveryPlugins {
		upstreams, errs, err := uds.WatchUpstreams(d.writeNamespace, opts, discOpts)
		if err != nil {
			return nil, errors.Wrapf(err, "initializing UDS for %v", reflect.TypeOf(uds).Name())
		}
		go func(uds DiscoveryPlugin) {
			for {
				select {
				case upstreamList := <-upstreams:
					lock.Lock()
					upstreamsByUds[uds] = upstreamList
					desiredUpstreams := aggregateUpstreams(upstreamsByUds)
					if err := d.upstreamReconciler.Reconcile(d.writeNamespace, desiredUpstreams, uds.UpdateUpstream, clients.ListOpts{
						Ctx:      opts.Ctx,
						Selector: opts.Selector,
					}); err != nil {
						aggregatedErrs <- err
					}
					lock.Unlock()
				case err := <-errs:
					aggregatedErrs <- errors.Wrapf(err, "error in uds plugin %v", reflect.TypeOf(uds).Name())
				case <-opts.Ctx.Done():
					return
				}
			}
		}(uds)
	}
	return aggregatedErrs, nil
}

func aggregateUpstreams(endpointsByUds map[DiscoveryPlugin]v1.UpstreamList) v1.UpstreamList {
	var upstreams v1.UpstreamList
	for _, upstreamList := range endpointsByUds {
		upstreams = append(upstreams, upstreamList...)
	}
	sort.SliceStable(upstreams, func(i, j int) bool {
		return upstreams[i].Metadata.Less(upstreams[j].Metadata)
	})
	return upstreams
}

// launch a goroutine for all the UDS plugins with a single cancel to close them all
func (d *Discovery) StartEds(upstreamsToTrack v1.UpstreamList, opts clients.WatchOpts) (chan error, error) {
	aggregatedErrs := make(chan error)
	endpointsByUds := make(map[DiscoveryPlugin]v1.EndpointList)
	lock := sync.Mutex{}
	for _, eds := range d.discoveryPlugins {
		endpoints, errs, err := eds.WatchEndpoints(d.writeNamespace, upstreamsToTrack, opts)
		if err != nil {
			return nil, errors.Wrapf(err, "initializing UDS for %v", reflect.TypeOf(eds).Name())
		}

		go func(eds DiscoveryPlugin) {
			for {
				select {
				case endpointList := <-endpoints:
					lock.Lock()
					endpointsByUds[eds] = endpointList
					desiredEndpoints := aggregateEndpoints(endpointsByUds)
					if err := d.endpointReconciler.Reconcile(d.writeNamespace, desiredEndpoints, nil, clients.ListOpts{
						Ctx:      opts.Ctx,
						Selector: opts.Selector,
					}); err != nil {
						aggregatedErrs <- err
					}
					lock.Unlock()
				case err := <-errs:
					aggregatedErrs <- errors.Wrapf(err, "error in eds plugin %v", reflect.TypeOf(eds).Name())
				case <-opts.Ctx.Done():
					return
				}
			}
		}(eds)
	}
	return aggregatedErrs, nil
}

func aggregateEndpoints(endpointsByUds map[DiscoveryPlugin]v1.EndpointList) v1.EndpointList {
	var endpoints v1.EndpointList
	for _, endpointList := range endpointsByUds {
		endpoints = append(endpoints, endpointList...)
	}
	sort.SliceStable(endpoints, func(i, j int) bool {
		return endpoints[i].Metadata.Less(endpoints[j].Metadata)
	})
	return endpoints
}
