package clients

import (
	"context"
	"time"

	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
)

const DefaultNamespace = "default"

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

type ListOpts struct {
	Ctx       context.Context
	Selector  map[string]string
	Namespace string
}

type WatchOpts struct {
	Ctx         context.Context
	Selector    map[string]string
	Namespace   string
	RefreshRate time.Duration
}

type WriteOpts struct {
	Ctx               context.Context
	OverwriteExisting bool
}

type UpdateOpts struct {
	Ctx context.Context
}

type DeleteOpts struct {
	Ctx            context.Context
	Namespace      string
	IgnoreNotExist bool
}
