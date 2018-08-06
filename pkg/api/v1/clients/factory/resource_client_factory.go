package factory

import (
	"reflect"

	"github.com/hashicorp/consul/api"
	vaultapi "github.com/hashicorp/vault/api"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/configmap"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/consul"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/file"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/kube"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/kube/crd"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/kubesecret"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/memory"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/vault"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/solo-io/solo-kit/pkg/errors"
	"k8s.io/client-go/kubernetes"
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
		inputResource, ok := params.ResourceType.(resources.InputResource)
		if !ok {
			return nil, errors.Errorf("the kubernetes crd client can only be used for input resources, received type %v", resources.Kind(resourceType))
		}
		return kube.NewResourceClient(opts.Crd, opts.Cfg, inputResource)
	case *ConsulResourceClientOpts:
		return consul.NewResourceClient(opts.Consul, opts.RootKey, resourceType), nil
	case *FileResourceClientOpts:
		return file.NewResourceClient(opts.RootDir, resourceType), nil
	case *MemoryResourceClientOpts:
		return memory.NewResourceClient(opts.Cache, resourceType), nil
	case *KubeConfigMapClientOpts:
		dataResource, ok := params.ResourceType.(resources.DataResource)
		if !ok {
			return nil, errors.Errorf("the kubernetes configmap client can only be used for data resources, received type %v", resources.Kind(resourceType))
		}
		return configmap.NewResourceClient(opts.Clientset, dataResource)
	case *KubeSecretClientOpts:
		dataResource, ok := params.ResourceType.(resources.DataResource)
		if !ok {
			return nil, errors.Errorf("the kubernetes secret client can only be used for data resources, received type %v", resources.Kind(resourceType))
		}
		return kubesecret.NewResourceClient(opts.Clientset, dataResource)
	case *VaultSecretClientOpts:
		dataResource, ok := params.ResourceType.(resources.DataResource)
		if !ok {
			return nil, errors.Errorf("the vault secret client can only be used for data resources, received type %v", resources.Kind(resourceType))
		}
		return vault.NewResourceClient(opts.Vault, opts.RootKey, dataResource), nil
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

type KubeConfigMapClientOpts struct {
	Clientset kubernetes.Interface
}

func (o *KubeConfigMapClientOpts) isResourceClientOpts() {}

type KubeSecretClientOpts struct {
	Clientset kubernetes.Interface
}

func (o *KubeSecretClientOpts) isResourceClientOpts() {}

type VaultSecretClientOpts struct {
	Vault   *vaultapi.Client
	RootKey string
}

func (o *VaultSecretClientOpts) isResourceClientOpts() {}
