package pkg

import (
	"context"

	"github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1"
)

// func Run(opts bootstrap.Options, discoveryOpts options.DiscoveryOptions) error {
// }
//
type DiscoverySyncer struct {
	updater       *Updater
	prevupstreams v1.UpstreamList
}

func (d *DiscoverySyncer) Sync(ctx context.Context, snap *v1.Snapshot) error {
	// update the upstream and secrets
	//	newupstreams = snap.UpstreamList
	return nil
}

func (d *DiscoverySyncer) Setup(context.Context) error {
	return nil

}

func (d *DiscoverySyncer) Update(upstreams v1.UpstreamList, secrets v1.SecretList) error {

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
		if _, err := one.Find(up.Metadata.Namespace, up.Metadata.Name); err != nil {
			// upstream from two is not present in one. add it to result list
			newlist = append(newlist, up)
		}
	}
	return newlist
}
