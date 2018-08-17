package v1

import (
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/factory"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/kube/crd"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-kit/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// TODO: modify as needed to populate additional fields
func NewProxy(namespace, name string) *Proxy {
	return &Proxy{
		Metadata: core.Metadata{
			Name:      name,
			Namespace: namespace,
		},
	}
}

func (r *Proxy) SetStatus(status core.Status) {
	r.Status = status
}

func (r *Proxy) SetMetadata(meta core.Metadata) {
	r.Metadata = meta
}

type ProxyList []*Proxy

// namespace is optional, if left empty, names can collide if the list contains more than one with the same name
func (list ProxyList) Find(namespace, name string) (*Proxy, error) {
	for _, proxy := range list {
		if proxy.Metadata.Name == name {
			if namespace == "" || proxy.Metadata.Namespace == namespace {
				return proxy, nil
			}
		}
	}
	return nil, errors.Errorf("list did not find proxy %v.%v", namespace, name)
}

func (list *ProxyList) AsResources() []resources.Resource {
	var ress []resources.Resource
	for _, proxy := range list {
		ress = append(ress, proxy)
	}
	return ress
}

func (list *ProxyList) AsInputResources() []resources.InputResource {
	var ress []resources.InputResource
	for _, proxy := range list {
		ress = append(ress, proxy)
	}
	return ress
}

var _ resources.Resource = &Proxy{}

type ProxyClient interface {
	BaseClient() clients.ResourceClient
	Register() error
	Read(namespace, name string, opts clients.ReadOpts) (*Proxy, error)
	Write(resource *Proxy, opts clients.WriteOpts) (*Proxy, error)
	Delete(namespace, name string, opts clients.DeleteOpts) error
	List(namespace string, opts clients.ListOpts) (ProxyList, error)
	Watch(namespace string, opts clients.WatchOpts) (<-chan ProxyList, <-chan error, error)
}

type proxyClient struct {
	rc clients.ResourceClient
}

func NewProxyClient(rcFactory factory.ResourceClientFactory) (ProxyClient, error) {
	rc, err := rcFactory.NewResourceClient(factory.NewResourceClientParams{
		ResourceType: &Proxy{},
	})
	if err != nil {
		return nil, errors.Wrapf(err, "creating base Proxy resource client")
	}
	return &proxyClient{
		rc: rc,
	}, nil
}

func (client *proxyClient) BaseClient() clients.ResourceClient {
	return client.rc
}

func (client *proxyClient) Register() error {
	return client.rc.Register()
}

func (client *proxyClient) Read(namespace, name string, opts clients.ReadOpts) (*Proxy, error) {
	opts = opts.WithDefaults()
	resource, err := client.rc.Read(namespace, name, opts)
	if err != nil {
		return nil, err
	}
	return resource.(*Proxy), nil
}

func (client *proxyClient) Write(proxy *Proxy, opts clients.WriteOpts) (*Proxy, error) {
	opts = opts.WithDefaults()
	resource, err := client.rc.Write(proxy, opts)
	if err != nil {
		return nil, err
	}
	return resource.(*Proxy), nil
}

func (client *proxyClient) Delete(namespace, name string, opts clients.DeleteOpts) error {
	opts = opts.WithDefaults()
	return client.rc.Delete(namespace, name, opts)
}

func (client *proxyClient) List(namespace string, opts clients.ListOpts) (ProxyList, error) {
	opts = opts.WithDefaults()
	resourceList, err := client.rc.List(namespace, opts)
	if err != nil {
		return nil, err
	}
	return convertToProxy(resourceList), nil
}

func (client *proxyClient) Watch(namespace string, opts clients.WatchOpts) (<-chan ProxyList, <-chan error, error) {
	opts = opts.WithDefaults()
	resourcesChan, errs, initErr := client.rc.Watch(namespace, opts)
	if initErr != nil {
		return nil, nil, initErr
	}
	proxysChan := make(chan ProxyList)
	go func() {
		for {
			select {
			case resourceList := <-resourcesChan:
				proxysChan <- convertToProxy(resourceList)
			case <-opts.Ctx.Done():
				close(proxysChan)
				return
			}
		}
	}()
	return proxysChan, errs, nil
}

func convertToProxy(resources []resources.Resource) ProxyList {
	var proxyList ProxyList
	for _, resource := range resources {
		proxyList = append(proxyList, resource.(*Proxy))
	}
	return proxyList
}

// Kubernetes Adapter for Proxy

func (o *Proxy) GetObjectKind() schema.ObjectKind {
	t := ProxyCrd.TypeMeta()
	return &t
}

func (o *Proxy) DeepCopyObject() runtime.Object {
	return resources.Clone(o).(*Proxy)
}

var ProxyCrd = crd.NewCrd("gloo.solo.io",
	"proxies",
	"gloo.solo.io",
	"v1",
	"Proxy",
	"px",
	&Proxy{})
