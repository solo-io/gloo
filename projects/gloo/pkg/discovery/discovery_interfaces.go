package discovery

import (
	"reflect"
	"sort"
	"sync"

	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/bootstrap"
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
	UpdateUpstream(original, desired *v1.Upstream) error

	// EDS API
	// start the EDS watch which sends a new list of endpoints on any change
	// will send only endpoints for upstreams configured with TrackUpstreams
	WatchEndpoints(writeNamespace string, upstreamsToTrack v1.UpstreamList, opts clients.WatchOpts) (<-chan v1.EndpointList, <-chan error, error)
}

type Discovery struct {
	bootstrap          bootstrap.Config
	upstreamReconciler v1.UpstreamReconciler
	endpointReconciler v1.EndpointReconciler
	discoveryPlugins   []DiscoveryPlugin
}

func NewDiscovery(upstreamClient v1.UpstreamClient,
	endpointsClient v1.EndpointClient,
	discoveryPlugins ...DiscoveryPlugin) *Discovery {
	return &Discovery{
		upstreamReconciler: v1.NewUpstreamReconciler(upstreamClient),
		endpointReconciler: v1.NewEndpointReconciler(endpointsClient),
		discoveryPlugins:   discoveryPlugins,
	}
}

// launch a goroutine for all the UDS plugins
func (d *Discovery) StartUds(writeNamespace string, opts clients.WatchOpts, discOpts Opts) (chan error, error) {
	aggregatedErrs := make(chan error)
	upstreamsByUds := make(map[DiscoveryPlugin]v1.UpstreamList)
	lock := sync.Mutex{}
	for _, uds := range d.discoveryPlugins {
		upstreams, errs, err := uds.WatchUpstreams(writeNamespace, opts, discOpts)
		if err != nil {
			return nil, errors.Wrapf(err, "initializing UDS for %v", reflect.TypeOf(uds).Name())
		}
		syncFunc := func(uds DiscoveryPlugin, upstreamList v1.UpstreamList) {
			lock.Lock()
			upstreamsByUds[uds] = upstreamList
			desiredUpstreams := aggregateUpstreams(upstreamsByUds)
			lock.Unlock()
			if err := d.upstreamReconciler.Reconcile(writeNamespace, desiredUpstreams, uds.UpdateUpstream, clients.ListOpts{
				Ctx:      opts.Ctx,
				Selector: opts.Selector,
			}); err != nil {
				aggregatedErrs <- err
			}
		}

		go func(uds DiscoveryPlugin) {
			for {
				select {
				case upstreamList := <-upstreams:
					syncFunc(uds, upstreamList)
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
	var endpoints v1.UpstreamList
	for _, endpointList := range endpointsByUds {
		endpoints = append(endpoints, endpointList...)
	}
	sort.SliceStable(endpoints, func(i, j int) bool {
		return endpoints[i].Metadata.Less(endpoints[j].Metadata)
	})
	return endpoints
}

// launch a goroutine for all the UDS plugins with a single cancel to close them all
func (d *Discovery) StartEds(namespace string, upstreamsToTrack v1.UpstreamList, opts clients.WatchOpts) (chan error, error) {
	aggregatedErrs := make(chan error)
	endpointsByUds := make(map[DiscoveryPlugin]v1.EndpointList)
	lock := sync.Mutex{}
	for _, eds := range d.discoveryPlugins {
		endpoints, errs, err := eds.WatchEndpoints(namespace, upstreamsToTrack, opts)
		if err != nil {
			return nil, errors.Wrapf(err, "initializing UDS for %v", reflect.TypeOf(eds).Name())
		}
		syncFunc := func(uds DiscoveryPlugin, endpointList v1.EndpointList) {
			lock.Lock()
			endpointsByUds[uds] = endpointList
			desiredEndpoints := aggregateEndpoints(endpointsByUds)
			lock.Unlock()
			if err := d.endpointReconciler.Reconcile(namespace, desiredEndpoints, nil, clients.ListOpts{
				Ctx:      opts.Ctx,
				Selector: opts.Selector,
			}); err != nil {
				aggregatedErrs <- err
			}
		}

		go func(eds DiscoveryPlugin) {
			for {
				select {
				case endpointList := <-endpoints:
					syncFunc(eds, endpointList)
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
