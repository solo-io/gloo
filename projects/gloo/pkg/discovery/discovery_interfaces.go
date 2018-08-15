package discovery

import (
	"context"
	"reflect"
	"sort"
	"sync"

	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-kit/pkg/bootstrap"
	"github.com/solo-io/solo-kit/pkg/errors"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1"
)

type DiscoveryPlugin interface {
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
	WatchEndpoints(namespace string, opts clients.WatchOpts) (<-chan v1.EndpointList, <-chan error, error)
	// Updates the EDS to track only these upstreams
	TrackUpstreams(list v1.UpstreamList)
}

type Discovery struct {
	upstreamReconciler  v1.UpstreamReconciler
	endpointsReconciler v1.EndpointReconciler
	discoveryPlugins    []DiscoveryPlugin
}

func NewDiscovery(upstreamClient v1.UpstreamClient,
	endpointsClient v1.EndpointClient,
	discoveryPlugins ...DiscoveryPlugin) *Discovery {
	return &Discovery{
		upstreamReconciler:  v1.NewUpstreamReconciler(upstreamClient),
		endpointsReconciler: v1.NewEndpointReconciler(endpointsClient),
		discoveryPlugins:    discoveryPlugins,
	}
}

func aggregateUpstreams(upstreamsByUds map[DiscoveryPlugin]v1.UpstreamList) v1.UpstreamList {
	var upstreams v1.UpstreamList
	for _, upstreamList := range upstreamsByUds {
		upstreams = append(upstreams, upstreamList...)
	}
	sort.SliceStable(upstreams, func(i, j int) bool {
		return upstreams[i].Metadata.Less(upstreams[j].Metadata)
	})
	return upstreams
}

// launch a goroutine for all the UDS plugins
func (d *Discovery) StartUds(bstrp bootstrap.Config, namespace string, opts clients.WatchOpts, discOpts Opts) (chan error, error) {
	aggregatedErrs := make(chan error)
	upstreamsByUds := make(map[DiscoveryPlugin]v1.UpstreamList)
	lock := sync.Mutex{}
	for _, uds := range d.discoveryPlugins {
		upstreams, errs, err := uds.WatchUpstreams(namespace, opts, discOpts)
		if err != nil {
			return nil, errors.Wrapf(err, "initializing UDS for %v", reflect.TypeOf(uds).Name())
		}
		reconcileUpstreams := func(uds DiscoveryPlugin, upstreamList v1.UpstreamList) {
			lock.Lock()
			upstreamsByUds[uds] = upstreamList
			desiredUpstreams := aggregateUpstreams(upstreamsByUds)
			lock.Unlock()
			if err := d.upstreamReconciler.Reconcile(namespace, desiredUpstreams, opts); err != nil {
				aggregatedErrs <- err
			}
		}

		go func(uds DiscoveryPlugin) {
			for {
				select {
				case upstreamList := <-upstreams:
					reconcileUpstreams(uds, upstreamList)
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

func reconcile(namespace string, desiredResources v1.UpstreamList, client v1.UpstreamClient, opts clients.WatchOpts, uds DiscoveryPlugin) error {
	originalResources, err := client.List(namespace, clients.ListOpts{
		Ctx:      opts.Ctx,
		Selector: opts.Selector,
	})
	if err != nil {
		return err
	}
	for _, desired := range desiredResources {
		if err := syncResource(opts.Ctx, client, desired, originalResources, uds); err != nil {
			return errors.Wrapf(err, "reconciling resource %v", desired.GetMetadata().Name)
		}
	}
	// delete unused
	for _, original := range originalResources {
		unused := findResource(original.GetMetadata().Namespace, original.GetMetadata().Name, desiredResources) == nil
		if unused {
			if err := deleteStaleResource(opts.Ctx, client, original); err != nil {
				return errors.Wrapf(err, "deleting stale resource %v", original.GetMetadata().Name)
			}
		}
	}
	return nil
}

func syncResource(ctx context.Context, client v1.UpstreamClient, desired *v1.Upstream, originalResources v1.UpstreamList, uds DiscoveryPlugin) error {
	var overwriteExisting bool
	original := findResource(desired.GetMetadata().Namespace, desired.GetMetadata().Name, originalResources)
	if original != nil {
		// if this is an update,
		// update resource version
		// set status to 0, needs to be re-processed
		overwriteExisting = true
		resources.UpdateMetadata(desired, func(meta *core.Metadata) {
			meta.ResourceVersion = original.GetMetadata().ResourceVersion
		})

		// reset the status
		desired.SetStatus(core.Status{})

		// call update upstream to allow any old context to be copied to the new object before writing
		if err := uds.UpdateUpstream(original, desired); err != nil {
			return err
		}
	}
	_, err := client.Write(desired, clients.WriteOpts{Ctx: ctx, OverwriteExisting: overwriteExisting})
	return err
}

func deleteStaleResource(ctx context.Context, client v1.UpstreamClient, original resources.Resource) error {
	return client.Delete(original.GetMetadata().Namespace, original.GetMetadata().Name, clients.DeleteOpts{
		Ctx:            ctx,
		IgnoreNotExist: true,
	})
}

func findResource(namespace, name string, rss v1.UpstreamList) *v1.Upstream {
	for _, resource := range rss {
		if resource.GetMetadata().Namespace == namespace && resource.GetMetadata().Name == name {
			return resource
		}
	}
	return nil
}

// launch a goroutine for all the EDS plugins
func (d *Discovery) StartEds(bstrp bootstrap.Config, namespace string, opts clients.WatchOpts, client v1.EndpointClient) (chan error, error) {
	aggregatedErrs := make(chan error)
	endpointsByEds := make(map[DiscoveryPlugin]v1.EndpointList)
	lock := sync.Mutex{}
	for _, eds := range d.discoveryPlugins {
		endpoints, errs, err := eds.WatchEndpoints(namespace, opts)
		if err != nil {
			return nil, errors.Wrapf(err, "initializing UDS for %v", reflect.TypeOf(eds).Name())
		}
		reconcileEndpoints := func(eds DiscoveryPlugin, endpointList v1.EndpointList) {
			lock.Lock()
			endpointsByEds[eds] = endpointList
			desiredEndpoints := endpointList
			for ds, upstreams := range endpointsByEds {
				if ds == eds {
					continue
				}
				desiredEndpoints = append(desiredEndpoints, upstreams...)
			}
			lock.Unlock()
			if err := reconcileEndpoints(namespace, desiredEndpoints, client, opts, eds); err != nil {
				aggregatedErrs <- err
			}
		}

		go func(uds DiscoveryPlugin) {
			for {
				select {
				case upstreamList := <-endpoints:
					reconcileEndpoints(uds, upstreamList)
				case err := <-errs:
					aggregatedErrs <- errors.Wrapf(err, "error in eds plugin %v", reflect.TypeOf(uds).Name())
				case <-opts.Ctx.Done():
					return
				}
			}
		}(eds)
	}
	return aggregatedErrs, nil
}
