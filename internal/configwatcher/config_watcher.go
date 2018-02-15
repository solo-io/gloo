package configwatcher

import (
	"fmt"
	"sync"

	"github.com/pkg/errors"

	"github.com/solo-io/glue-storage"
	"github.com/solo-io/glue/pkg/api/types/v1"
	"github.com/solo-io/glue/pkg/configwatcher"
)

func NewConfigWatcher(storageClient storage.Storage) (configwatcher.Interface, error) {
	if err := storageClient.V1().Register(); err != nil {
		return nil, fmt.Errorf("failed to register to storage backend: %v", err)
	}
	cache := v1.Config{}
	configs := make(chan *v1.Config)
	syncUpstreams := func(updatedList []*v1.Upstream, _ *v1.Upstream) {
		cache.Upstreams = updatedList
		configs <- &cache
	}
	upstreamWatcher, err := storageClient.V1().Upstreams().Watch(&storage.UpstreamEventHandlerFuncs{
		AddFunc:    syncUpstreams,
		UpdateFunc: syncUpstreams,
		DeleteFunc: syncUpstreams,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to create watcher for upstreams")
	}
	syncVhosts := func(updatedList []*v1.VirtualHost, _ *v1.VirtualHost) {
		cache.VirtualHosts = updatedList
		configs <- &cache
	}
	vhostWatcher, err := storageClient.V1().VirtualHosts().Watch(&storage.VirtualHostEventHandlerFuncs{
		AddFunc:    syncVhosts,
		UpdateFunc: syncVhosts,
		DeleteFunc: syncVhosts,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to create watcher for virtualhosts")
	}

	return &configWatcher{
		watchers: []*storage.Watcher{vhostWatcher, upstreamWatcher},
		configs:  configs,
		errs:     make(chan error),
	}, nil
}

type configWatcher struct {
	watchers []*storage.Watcher
	configs  chan *v1.Config
	errs     chan error
}

func (w *configWatcher) Run(stop <-chan struct{}) {
	done := &sync.WaitGroup{}
	for _, watcher := range w.watchers {
		done.Add(1)
		go func(watcher *storage.Watcher, stop <-chan struct{}, errs chan error) {
			watcher.Run(stop, errs)
			done.Done()
		}(watcher, stop, w.errs)
	}
	done.Wait()
}

func (w *configWatcher) Config() <-chan *v1.Config {
	return w.configs
}

func (w *configWatcher) Error() <-chan error {
	return w.errs
}
