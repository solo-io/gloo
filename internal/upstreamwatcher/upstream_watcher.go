package upstreamwatcher

import (
	"github.com/solo-io/gloo-api/pkg/api/types/v1"
	"github.com/solo-io/gloo-storage"
)

// returns a channel of upstreams that need function discovery
// upstream selector tells us whether the upstream should be selected for discovery
func WatchUpstreams(gloo storage.Interface,
	selectUpstream func(*v1.Upstream) bool,
	stop <-chan struct{},
	errs chan error) (<-chan []*v1.Upstream, error) {

	upstreams := make(chan []*v1.Upstream)
	syncFunc := func(newList []*v1.Upstream, _ *v1.Upstream) {
		selected := selectUpstreams(newList, selectUpstream)
		if len(selected) > 0 {
			upstreams <- selected
		}
	}
	w, err := gloo.V1().Upstreams().Watch(storage.UpstreamEventHandlerFuncs{
		AddFunc:    syncFunc,
		UpdateFunc: syncFunc,
	})
	if err != nil {
		return nil, err
	}
	go w.Run(stop, errs)
	return upstreams, nil
}

func selectUpstreams(list []*v1.Upstream, selectUpstream func(*v1.Upstream) bool) []*v1.Upstream {
	var selected []*v1.Upstream
	for _, us := range list {
		if selectUpstream(us) {
			selected = append(selected, us)
		}
	}
	return selected
}
