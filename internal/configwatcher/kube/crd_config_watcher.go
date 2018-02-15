package kube

import (
	"fmt"
	"sync"
	"time"

	"github.com/pkg/errors"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/solo-io/glue-storage"
	"github.com/solo-io/glue-storage/crd"
	"github.com/solo-io/glue/pkg/api/types/v1"
	"github.com/solo-io/glue/pkg/configwatcher"
)

func NewCrdWatcher(masterUrl, kubeconfigPath, namespace string, resyncDuration time.Duration) (configwatcher.Interface, error) {
	cfg, err := clientcmd.BuildConfigFromFlags(masterUrl, kubeconfigPath)
	if err != nil {
		return nil, fmt.Errorf("failed to build rest config: %v", err)
	}
	storageClient, err := crd.NewStorage(cfg, namespace, resyncDuration)
	if err != nil {
		return nil, fmt.Errorf("create storage client: %v", err)
	}
	err = storageClient.V1().Register()
	if err != nil {
		return nil, fmt.Errorf("failed to register crds: %v", err)
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
	}, nil
}

type configWatcher struct {
	watchers []*storage.Watcher
	configs  chan *v1.Config
}

func (w *configWatcher) Run(stop <-chan struct{}) {
	done := &sync.WaitGroup{}
	for _, watcher := range w.watchers {
		done.Add(1)
		go func(watcher *storage.Watcher) {
			watcher.Run(stop)
			done.Done()
		}(watcher)
	}
	done.Wait()
}
func (w *configWatcher) Config() <-chan *v1.Config {
	return w.configs
}

// implemented only for the interface. errors will not be emitted on this chan
func (w *configWatcher) Error() <-chan error {
	return make(chan error)
}
