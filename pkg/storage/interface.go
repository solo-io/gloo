package storage

import "github.com/solo-io/gloo/pkg/api/types/v1"

// Interface is interface to the storage backend
type Interface interface {
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
	Watch(handlers ...UpstreamEventHandler) (*Watcher, error)
}

type VirtualHosts interface {
	Create(*v1.VirtualHost) (*v1.VirtualHost, error)
	Update(*v1.VirtualHost) (*v1.VirtualHost, error)
	Delete(name string) error
	Get(name string) (*v1.VirtualHost, error)
	List() ([]*v1.VirtualHost, error)
	Watch(...VirtualHostEventHandler) (*Watcher, error)
}
