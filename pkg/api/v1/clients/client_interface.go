package clients

import (
	"context"
	"time"

	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
)

const DefaultNamespace = "default"

type ResourceClient interface {
	Register() error
	Read(name string, opts GetOptions) (resources.Resource, error)
	Write(resource resources.Resource, opts WriteOptions) (resources.Resource, error)
	Delete(name string, opts DeleteOptions) error
	List(opts ListOptions) ([]resources.Resource, error)
	Watch(opts WatchOptions) (<-chan []resources.Resource, error)
}

type GetOptions struct {
	Ctx       context.Context
	Namespace string
}

type ListOptions struct {
	Ctx       context.Context
	Selector  map[string]string
	Namespace string
}

type WatchOptions struct {
	Ctx         context.Context
	Selector    map[string]string
	Namespace   string
	RefreshRate time.Duration
}

type WriteOptions struct {
	Ctx               context.Context
	OverwriteExisting bool
}

type UpdateOptions struct {
	Ctx context.Context
}

type DeleteOptions struct {
	Ctx context.Context
}
