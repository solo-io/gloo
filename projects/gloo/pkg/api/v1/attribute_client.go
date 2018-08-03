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
func NewAttribute(namespace, name string) *Attribute {
	return &Attribute{
		Metadata: core.Metadata{
			Name:      name,
			Namespace: namespace,
		},
	}
}

func (r *Attribute) SetStatus(status core.Status) {
	r.Status = status
}

func (r *Attribute) SetMetadata(meta core.Metadata) {
	r.Metadata = meta
}

var _ resources.Resource = &Attribute{}

type AttributeClient interface {
	Register() error
	Read(namespace, name string, opts clients.ReadOpts) (*Attribute, error)
	Write(resource *Attribute, opts clients.WriteOpts) (*Attribute, error)
	Delete(namespace, name string, opts clients.DeleteOpts) error
	List(namespace string, opts clients.ListOpts) ([]*Attribute, error)
	Watch(namespace string, opts clients.WatchOpts) (<-chan []*Attribute, <-chan error, error)
}

type attributeClient struct {
	rc clients.ResourceClient
}

func NewAttributeClient(rcFactory factory.ResourceClientFactory) (AttributeClient, error) {
	rc, err := rcFactory.NewResourceClient(factory.NewResourceClientParams{
		ResourceType: &Attribute{},
	})
	if err != nil {
		return nil, errors.Wrapf(err, "creating base Attribute resource client")
	}
	return &attributeClient{
		rc: rc,
	}, nil
}

func (client *attributeClient) Register() error {
	return client.rc.Register()
}

func (client *attributeClient) Read(namespace, name string, opts clients.ReadOpts) (*Attribute, error) {
	resource, err := client.rc.Read(namespace, name, opts)
	if err != nil {
		return nil, err
	}
	return resource.(*Attribute), nil
}

func (client *attributeClient) Write(attribute *Attribute, opts clients.WriteOpts) (*Attribute, error) {
	resource, err := client.rc.Write(attribute, opts)
	if err != nil {
		return nil, err
	}
	return resource.(*Attribute), nil
}

func (client *attributeClient) Delete(namespace, name string, opts clients.DeleteOpts) error {
	return client.rc.Delete(namespace, name, opts)
}

func (client *attributeClient) List(namespace string, opts clients.ListOpts) ([]*Attribute, error) {
	resourceList, err := client.rc.List(namespace, opts)
	if err != nil {
		return nil, err
	}
	return convertToAttribute(resourceList), nil
}

func (client *attributeClient) Watch(namespace string, opts clients.WatchOpts) (<-chan []*Attribute, <-chan error, error) {
	resourcesChan, errs, initErr := client.rc.Watch(namespace, opts)
	if initErr != nil {
		return nil, nil, initErr
	}
	attributesChan := make(chan []*Attribute)
	go func() {
		for {
			select {
			case resourceList := <-resourcesChan:
				attributesChan <- convertToAttribute(resourceList)
			}
		}
	}()
	return attributesChan, errs, nil
}

func convertToAttribute(resources []resources.Resource) []*Attribute {
	var attributeList []*Attribute
	for _, resource := range resources {
		attributeList = append(attributeList, resource.(*Attribute))
	}
	return attributeList
}

// Kubernetes Adapter for Attribute

func (o *Attribute) GetObjectKind() schema.ObjectKind {
	t := AttributeCrd.TypeMeta()
	return &t
}

func (o *Attribute) DeepCopyObject() runtime.Object {
	return resources.Clone(o).(*Attribute)
}

var AttributeCrd = crd.NewCrd("testing.solo.io",
	"mocks",
	"testing.solo.io",
	"v1",
	"Attribute",
	"mk",
	&Attribute{})
