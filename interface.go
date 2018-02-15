package storage

import "github.com/solo-io/glue/pkg/api/types/v1"

// Storage is interface to the storage backend
type Storage interface {
	V1() V1
}

type V1 interface {
	Register() error
	Upstreams() Upstreams
	VirtualHosts() VirtualHosts
}

type Upstreams interface {
	Create(*v1.Upstream) (*v1.Upstream, error)
	Update(*v1.Upstream) (*v1.Upstream, error)
	Delete(name string) error
	Get(name string) (*v1.Upstream, error)
	List() ([]*v1.Upstream, error)
	Watch(handlers ...UpstreamEventHandler) *Watcher
}

type VirtualHosts interface {
	Create(*v1.VirtualHost) (*v1.VirtualHost, error)
	Update(*v1.VirtualHost) (*v1.VirtualHost, error)
	Delete(name string) error
	Get(name string) (*v1.VirtualHost, error)
	List() ([]*v1.VirtualHost, error)
	Watch(...VirtualHostEventHandler) *Watcher
}

type Watcher struct {
	runFunc func(stop <-chan struct{})
}

func NewWatcher(runFunc func(stop <-chan struct{})) *Watcher {
	return &Watcher{runFunc: runFunc}
}

func (w *Watcher) Run(stop <-chan struct{}) {
	w.runFunc(stop)
}

type UpstreamEventHandler interface {
	OnAdd(updatedList []*v1.Upstream, obj *v1.Upstream)
	OnUpdate(updatedList []*v1.Upstream, oldObj, newObj *v1.Upstream)
	OnDelete(updatedList []*v1.Upstream, obj *v1.Upstream)
}

type VirtualHostEventHandler interface {
	OnAdd(updatedList []*v1.VirtualHost, obj *v1.VirtualHost)
	OnUpdate(updatedList []*v1.VirtualHost, oldObj, newObj *v1.VirtualHost)
	OnDelete(updatedList []*v1.VirtualHost, obj *v1.VirtualHost)
}
