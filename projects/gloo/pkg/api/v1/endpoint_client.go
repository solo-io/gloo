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
func NewEndpoint(namespace, name string) *Endpoint {
	return &Endpoint{
		Metadata: core.Metadata{
			Name:      name,
			Namespace: namespace,
		},
	}
}

func (r *Endpoint) SetMetadata(meta core.Metadata) {
	r.Metadata = meta
}

var _ resources.Resource = &Endpoint{}

type EndpointClient interface {
	Register() error
	Read(namespace, name string, opts clients.ReadOpts) (*Endpoint, error)
	Write(resource *Endpoint, opts clients.WriteOpts) (*Endpoint, error)
	Delete(namespace, name string, opts clients.DeleteOpts) error
	List(namespace string, opts clients.ListOpts) ([]*Endpoint, error)
	Watch(namespace string, opts clients.WatchOpts) (<-chan []*Endpoint, <-chan error, error)
}

type endpointClient struct {
	rc clients.ResourceClient
}

func NewEndpointClient(rcFactory factory.ResourceClientFactory) (EndpointClient, error) {
	rc, err := rcFactory.NewResourceClient(factory.NewResourceClientParams{
		ResourceType: &Endpoint{},
	})
	if err != nil {
		return nil, errors.Wrapf(err, "creating base Endpoint resource client")
	}
	return &endpointClient{
		rc: rc,
	}, nil
}

func (client *endpointClient) Register() error {
	return client.rc.Register()
}

func (client *endpointClient) Read(namespace, name string, opts clients.ReadOpts) (*Endpoint, error) {
	opts = opts.WithDefaults()
	resource, err := client.rc.Read(namespace, name, opts)
	if err != nil {
		return nil, err
	}
	return resource.(*Endpoint), nil
}

func (client *endpointClient) Write(endpoint *Endpoint, opts clients.WriteOpts) (*Endpoint, error) {
	opts = opts.WithDefaults()
	resource, err := client.rc.Write(endpoint, opts)
	if err != nil {
		return nil, err
	}
	return resource.(*Endpoint), nil
}

func (client *endpointClient) Delete(namespace, name string, opts clients.DeleteOpts) error {
	opts = opts.WithDefaults()
	return client.rc.Delete(namespace, name, opts)
}

func (client *endpointClient) List(namespace string, opts clients.ListOpts) ([]*Endpoint, error) {
	opts = opts.WithDefaults()
	resourceList, err := client.rc.List(namespace, opts)
	if err != nil {
		return nil, err
	}
	return convertToEndpoint(resourceList), nil
}

func (client *endpointClient) Watch(namespace string, opts clients.WatchOpts) (<-chan []*Endpoint, <-chan error, error) {
	opts = opts.WithDefaults()
	resourcesChan, errs, initErr := client.rc.Watch(namespace, opts)
	if initErr != nil {
		return nil, nil, initErr
	}
	endpointsChan := make(chan []*Endpoint)
	go func() {
		for {
			select {
			case resourceList := <-resourcesChan:
				endpointsChan <- convertToEndpoint(resourceList)
			case <-opts.Ctx.Done():
				close(endpointsChan)
				return
			}
		}
	}()
	return endpointsChan, errs, nil
}

func convertToEndpoint(resources []resources.Resource) []*Endpoint {
	var endpointList []*Endpoint
	for _, resource := range resources {
		endpointList = append(endpointList, resource.(*Endpoint))
	}
	return endpointList
}

// Kubernetes Adapter for Endpoint

func (o *Endpoint) GetObjectKind() schema.ObjectKind {
	t := EndpointCrd.TypeMeta()
	return &t
}

func (o *Endpoint) DeepCopyObject() runtime.Object {
	return resources.Clone(o).(*Endpoint)
}

var EndpointCrd = crd.NewCrd("gloo.solo.io",
	"endpoints",
	"gloo.solo.io",
	"v1",
	"Endpoint",
	"ep",
	&Endpoint{})
