package fds

import (
	"context"

	"github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1"
)

type FunctionDiscovery struct {
	updater       *Updater
	prevupstreams v1.UpstreamList
}

func NewFunctionDiscovery(updater *Updater) *FunctionDiscovery {
	return &FunctionDiscovery{
		updater: updater,
	}
}

func (d *FunctionDiscovery) Sync(ctx context.Context, snap *v1.DiscoverySnapshot) error {
	// update the upstream and secrets
	//	this will feed the update loop via a channel
	return nil
}

func (d *FunctionDiscovery) Setup(context.Context) error {
	/*
		setup an event loop and an update loop
	*/
	return nil
}

func (d *FunctionDiscovery) Update(upstreams v1.UpstreamList, secrets v1.SecretList) error {
	d.updater.SetSecrets(secrets)
	// get new snapshot from sync and update the upstreams and secrets in the updater
	old := d.prevupstreams
	d.prevupstreams = upstreams

	// find one the ones that were removed.
	removed := diff(old, upstreams)
	added := diff(upstreams, old)
	// find the once that are left, and update them.
	potentiallyUpdated := diff(upstreams, added)

	for _, u := range removed {
		d.updater.UpstreamRemoved(u)
	}
	for _, u := range added {
		d.updater.UpstreamAdded(u)
	}
	for _, u := range potentiallyUpdated {
		// TODO: TEST IF THEY WERE REALLY CHANGED, perhaps by comparing the resource version?
		d.updater.UpstreamUpdated(u)
	}

	return nil
}

func diff(one, two v1.UpstreamList) v1.UpstreamList {
	newlist := make([]*v1.Upstream, 0, len(one))

	for _, up := range two {
		meta := up.Metadata
		if _, err := one.Find(meta.Namespace, meta.Name); err != nil {
			// upstream from two is not present in one. add it to result list
			newlist = append(newlist, up)
		}
	}
	return newlist
}
