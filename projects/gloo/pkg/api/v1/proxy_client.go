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

var _ resources.Resource = &Proxy{}

type ProxyClient interface {
	Register() error
	Read(namespace, name string, opts clients.ReadOpts) (*Proxy, error)
	Write(resource *Proxy, opts clients.WriteOpts) (*Proxy, error)
	Delete(namespace, name string, opts clients.DeleteOpts) error
	List(namespace string, opts clients.ListOpts) ([]*Proxy, error)
	Watch(namespace string, opts clients.WatchOpts) (<-chan []*Proxy, <-chan error, error)
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

func (client *proxyClient) Register() error {
	return client.rc.Register()
}

func (client *proxyClient) Read(namespace, name string, opts clients.ReadOpts) (*Proxy, error) {
	resource, err := client.rc.Read(namespace, name, opts)
	if err != nil {
		return nil, err
	}
	return resource.(*Proxy), nil
}

func (client *proxyClient) Write(proxy *Proxy, opts clients.WriteOpts) (*Proxy, error) {
	resource, err := client.rc.Write(proxy, opts)
	if err != nil {
		return nil, err
	}
	return resource.(*Proxy), nil
}

func (client *proxyClient) Delete(namespace, name string, opts clients.DeleteOpts) error {
	return client.rc.Delete(namespace, name, opts)
}

func (client *proxyClient) List(namespace string, opts clients.ListOpts) ([]*Proxy, error) {
	resourceList, err := client.rc.List(namespace, opts)
	if err != nil {
		return nil, err
	}
	return convertToProxy(resourceList), nil
}

func (client *proxyClient) Watch(namespace string, opts clients.WatchOpts) (<-chan []*Proxy, <-chan error, error) {
	resourcesChan, errs, initErr := client.rc.Watch(namespace, opts)
	if initErr != nil {
		return nil, nil, initErr
	}
	proxysChan := make(chan []*Proxy)
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

func convertToProxy(resources []resources.Resource) []*Proxy {
	var proxyList []*Proxy
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
