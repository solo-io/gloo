package discovery

import (
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/bootstrap"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1"
	"github.com/solo-io/solo-kit/pkg/errors"
	"reflect"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"context"
)

type UdsPlugin interface {
	// send us an updated list of upstreams on every change
	WatchUpstreams(namespace string, opts clients.WatchOpts, discOpts Opts) (chan v1.UpstreamList, chan error, error)
	// tell us how to update from an existing upstream
	UpdateUpstream(original, desired *v1.Upstream)
}

type Discovery struct {
	uds []UdsPlugin
}

func NewDiscovery(uds []UdsPlugin) *Discovery {
	return &Discovery{uds: uds}
}

func (d *Discovery) Start(bstrp bootstrap.Config, namespace string, opts clients.WatchOpts, discOpts Opts, client v1.UpstreamClient) (chan error, error) {
	aggregatedErrs := make(chan error)
	for _, uds := range d.uds {
		upstreams, errs, err := uds.WatchUpstreams(namespace, opts, discOpts)
		if err != nil {
			return nil, errors.Wrapf(err, "initializing UDS for %v", reflect.TypeOf(uds).Name())
		}
		go func(uds UdsPlugin) {
			for {
				select {
				case upstreamList := <-upstreams:
					if err := reconcile(namespace, upstreamList, client, opts, uds); err != nil {
						aggregatedErrs <- err
					}
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

func reconcile(namespace string, desiredResources v1.UpstreamList, client v1.UpstreamClient, opts clients.WatchOpts, uds UdsPlugin) error {
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

func syncResource(ctx context.Context, client v1.UpstreamClient, desired *v1.Upstream, originalResources v1.UpstreamList, uds UdsPlugin) error {
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

		desired.SetStatus(core.Status{})

		uds.UpdateUpstream(original, desired)
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
