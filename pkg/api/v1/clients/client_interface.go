package clients

import (
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"context"
)

type ResourceClient interface {
	Register() error
	Get(name string, opts *GetOptions) (resources.Resource, error)
	Create(resource resources.Resource, opts *CreateOptions) (resources.Resource, error)
	Update(resource resources.Resource, opts *UpdateOptions) (resources.Resource, error)
	Delete(name string, opts *DeleteOptions) error
	List(opts *ListOptions) ([]resources.Resource, error)
	Watch(opts *WatchOptions) (<-chan []resources.Resource, error)
}

type GetOptions struct {
	Ctx  context.Context
	Selector map[string]string
	Namespace string
}

type ListOptions struct {
	Ctx  context.Context
	Selector map[string]string
	Namespace string
}

type WatchOptions struct {
	Ctx  context.Context
	Selector map[string]string
	Namespace string
}

type CreateOptions struct {
	Ctx  context.Context
}

type UpdateOptions struct {
	Ctx  context.Context
}

type DeleteOptions struct {
	Ctx  context.Context
}
