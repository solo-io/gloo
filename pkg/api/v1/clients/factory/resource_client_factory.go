package factory

import (
	"reflect"

	"github.com/hashicorp/consul/api"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/consul"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/file"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/kube"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/kube/crd"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/kube/crd/client/clientset/versioned"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	apiexts "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
)

type ResourceClientFactory struct {
	opts resourceClientOpts
}

func NewResourceClientFactory(opts resourceClientOpts) *ResourceClientFactory {
	return &ResourceClientFactory{
		opts: opts,
	}
}

func (factory *ResourceClientFactory) NewResourceClient(resourceType resources.Resource) clients.ResourceClient {
	if factory.opts == nil {
		panic("resource client factory opts cannot be nil")
	}
	switch opts := factory.opts.(type) {
	case *KubeResourceClientOpts:
		return kube.NewResourceClient(opts.Crd, opts.ApiExts, opts.Kube, resourceType)
	case *ConsulResourceClientOpts:
		return consul.NewResourceClient(opts.Consul, opts.RootKey, resourceType)
	case *FileResourceClientOpts:
		return file.NewResourceClient(opts.RootDir, resourceType)
	}
	panic("unsupported type " + reflect.TypeOf(factory.opts).Name())
}

// https://golang.org/doc/faq#generics
type resourceClientOpts interface {
	isResourceClientOpts()
}

type KubeResourceClientOpts struct {
	Crd     crd.Crd
	ApiExts apiexts.Interface
	Kube    versioned.Interface
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
