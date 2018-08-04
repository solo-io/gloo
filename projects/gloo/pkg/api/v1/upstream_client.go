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
func NewUpstream(namespace, name string) *Upstream {
	return &Upstream{
		Metadata: core.Metadata{
			Name:      name,
			Namespace: namespace,
		},
	}
}

func (r *Upstream) SetStatus(status core.Status) {
	r.Status = status
}

func (r *Upstream) SetMetadata(meta core.Metadata) {
	r.Metadata = meta
}

var _ resources.Resource = &Upstream{}

type UpstreamClient interface {
	Register() error
	Read(namespace, name string, opts clients.ReadOpts) (*Upstream, error)
	Write(resource *Upstream, opts clients.WriteOpts) (*Upstream, error)
	Delete(namespace, name string, opts clients.DeleteOpts) error
	List(namespace string, opts clients.ListOpts) ([]*Upstream, error)
	Watch(namespace string, opts clients.WatchOpts) (<-chan []*Upstream, <-chan error, error)
}

type upstreamClient struct {
	rc clients.ResourceClient
}

func NewUpstreamClient(rcFactory factory.ResourceClientFactory) (UpstreamClient, error) {
	rc, err := rcFactory.NewResourceClient(factory.NewResourceClientParams{
		ResourceType: &Upstream{},
	})
	if err != nil {
		return nil, errors.Wrapf(err, "creating base Upstream resource client")
	}
	return &upstreamClient{
		rc: rc,
	}, nil
}

func (client *upstreamClient) Register() error {
	return client.rc.Register()
}

func (client *upstreamClient) Read(namespace, name string, opts clients.ReadOpts) (*Upstream, error) {
	resource, err := client.rc.Read(namespace, name, opts)
	if err != nil {
		return nil, err
	}
	return resource.(*Upstream), nil
}

func (client *upstreamClient) Write(upstream *Upstream, opts clients.WriteOpts) (*Upstream, error) {
	resource, err := client.rc.Write(upstream, opts)
	if err != nil {
		return nil, err
	}
	return resource.(*Upstream), nil
}

func (client *upstreamClient) Delete(namespace, name string, opts clients.DeleteOpts) error {
	return client.rc.Delete(namespace, name, opts)
}

func (client *upstreamClient) List(namespace string, opts clients.ListOpts) ([]*Upstream, error) {
	resourceList, err := client.rc.List(namespace, opts)
	if err != nil {
		return nil, err
	}
	return convertToUpstream(resourceList), nil
}

func (client *upstreamClient) Watch(namespace string, opts clients.WatchOpts) (<-chan []*Upstream, <-chan error, error) {
	resourcesChan, errs, initErr := client.rc.Watch(namespace, opts)
	if initErr != nil {
		return nil, nil, initErr
	}
	upstreamsChan := make(chan []*Upstream)
	go func() {
		for {
			select {
			case resourceList := <-resourcesChan:
				upstreamsChan <- convertToUpstream(resourceList)
			}
		}
	}()
	return upstreamsChan, errs, nil
}

func convertToUpstream(resources []resources.Resource) []*Upstream {
	var upstreamList []*Upstream
	for _, resource := range resources {
		upstreamList = append(upstreamList, resource.(*Upstream))
	}
	return upstreamList
}

// Kubernetes Adapter for Upstream

func (o *Upstream) GetObjectKind() schema.ObjectKind {
	t := UpstreamCrd.TypeMeta()
	return &t
}

func (o *Upstream) DeepCopyObject() runtime.Object {
	return resources.Clone(o).(*Upstream)
}

var UpstreamCrd = crd.NewCrd("gloo.solo.io",
	"upstreams",
	"gloo.solo.io",
	"v1",
	"Upstream",
	"us",
	&Upstream{})
