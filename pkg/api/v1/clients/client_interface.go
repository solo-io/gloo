package clients

import (
	"context"
	"time"

	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
)

const DefaultNamespace = "default"

var DefaultRefreshRate = time.Second * 30

type ResourceClient interface {
	Register() error
	Read(name string, opts GetOpts) (resources.Resource, error)
	Write(resource resources.Resource, opts WriteOpts) (resources.Resource, error)
	Delete(name string, opts DeleteOpts) error
	List(opts ListOpts) ([]resources.Resource, error)
	Watch(opts WatchOpts) (<-chan []resources.Resource, error)
}

type GetOpts struct {
	Ctx       context.Context
	Namespace string
}

func (o GetOpts) WithDefaults() GetOpts {
	if o.Ctx == nil {
		o.Ctx = context.TODO()
	}
	if o.Namespace == "" {
		o.Namespace = DefaultNamespace
	}
	return o
}

type ListOpts struct {
	Ctx       context.Context
	Selector  map[string]string
	Namespace string
}

func (o ListOpts) WithDefaults() ListOpts {
	if o.Ctx == nil {
		o.Ctx = context.TODO()
	}
	if o.Namespace == "" {
		o.Namespace = DefaultNamespace
	}
	return o
}

type WatchOpts struct {
	Ctx         context.Context
	Selector    map[string]string
	Namespace   string
	RefreshRate time.Duration
}

func (o WatchOpts) WithDefaults() WatchOpts {
	if o.Ctx == nil {
		o.Ctx = context.TODO()
	}
	if o.Namespace == "" {
		o.Namespace = DefaultNamespace
	}
	if o.RefreshRate == 0 {
		o.RefreshRate = DefaultRefreshRate
	}
	return o
}

type WriteOpts struct {
	Ctx               context.Context
	OverwriteExisting bool
}

func (o WriteOpts) WithDefaults() WriteOpts {
	if o.Ctx == nil {
		o.Ctx = context.TODO()
	}
	return o
}

type DeleteOpts struct {
	Ctx            context.Context
	Namespace      string
	IgnoreNotExist bool
}

func (o DeleteOpts) WithDefaults() DeleteOpts {
	if o.Ctx == nil {
		o.Ctx = context.TODO()
	}
	if o.Namespace == "" {
		o.Namespace = DefaultNamespace
	}
	return o
}
