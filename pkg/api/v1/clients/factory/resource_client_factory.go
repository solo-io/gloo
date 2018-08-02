package factory

import (
	"reflect"

	"github.com/hashicorp/consul/api"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/consul"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/file"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/kube"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/kube/crd"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/memory"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"k8s.io/client-go/rest"
)

type ResourceClientFactory interface {
	NewResourceClient(params NewResourceClientParams) (clients.ResourceClient, error)
}

type resourceClientFactory struct {
	opts ResourceClientFactoryOpts
}

func NewResourceClientFactory(opts ResourceClientFactoryOpts) ResourceClientFactory {
	if opts == nil {
		panic("resource client factory opts cannot be nil")
	}
	return &resourceClientFactory{
		opts: opts,
	}
}

type NewResourceClientParams struct {
	ResourceType resources.Resource
	Token        string
}

func (factory *resourceClientFactory) NewResourceClient(params NewResourceClientParams) (clients.ResourceClient, error) {
	resourceType := params.ResourceType
	switch opts := factory.opts.(type) {
	case *KubeResourceClientOpts:
		if params.Token != "" {
			opts.Cfg.BearerToken = params.Token
		}
		return kube.NewResourceClient(opts.Crd, opts.Cfg, resourceType)
	case *ConsulResourceClientOpts:
		return consul.NewResourceClient(opts.Consul, opts.RootKey, resourceType), nil
	case *FileResourceClientOpts:
		return file.NewResourceClient(opts.RootDir, resourceType), nil
	case *MemoryResourceClientOpts:
		return memory.NewResourceClient(opts.Cache, resourceType), nil
	}
	panic("unsupported type " + reflect.TypeOf(factory.opts).Name())
}

// https://golang.org/doc/faq#generics
type ResourceClientFactoryOpts interface {
	isResourceClientOpts()
}

type KubeResourceClientOpts struct {
	Crd crd.Crd
	Cfg *rest.Config
}

func (o *KubeResourceClientOpts) isResourceClientOpts() {}

type ConsulResourceClientOpts struct {
	Consul  *api.Client
	RootKey string
}

func (o *ConsulResourceClientOpts) isResourceClientOpts() {}

type FileResourceClientOpts struct {
	RootDir string
}

func (o *FileResourceClientOpts) isResourceClientOpts() {}

type MemoryResourceClientOpts struct {
	Cache memory.InMemoryResourceCache
}

func (o *MemoryResourceClientOpts) isResourceClientOpts() {}
