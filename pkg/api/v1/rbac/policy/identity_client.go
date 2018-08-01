package policy

import (
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/factory"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/kube/crd"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// TODO: modify as needed to populate additional fields
func NewIdentity(namespace, name string) *Identity {
	return &Identity{
		Metadata: core.Metadata{
			Name:      name,
			Namespace: namespace,
		},
	}
}

func (r *Identity) SetStatus(status core.Status) {
	r.Status = status
}

func (r *Identity) SetMetadata(meta core.Metadata) {
	r.Metadata = meta
}

var _ resources.Resource = &Identity{}

type IdentityClient interface {
	Register() error
	Read(namespace, name string, opts clients.ReadOpts) (*Identity, error)
	Write(resource *Identity, opts clients.WriteOpts) (*Identity, error)
	Delete(namespace, name string, opts clients.DeleteOpts) error
	List(namespace string, opts clients.ListOpts) ([]*Identity, error)
	Watch(namespace string, opts clients.WatchOpts) (<-chan []*Identity, <-chan error, error)
}

type identityClient struct {
	rc clients.ResourceClient
}

func NewIdentityClient(factory *factory.ResourceClientFactory) IdentityClient {
	return &identityClient{
		rc: factory.NewResourceClient(&Identity{}),
	}
}

func (client *identityClient) Register() error {
	return client.rc.Register()
}

func (client *identityClient) Read(namespace, name string, opts clients.ReadOpts) (*Identity, error) {
	resource, err := client.rc.Read(namespace, name, opts)
	if err != nil {
		return nil, err
	}
	return resource.(*Identity), nil
}

func (client *identityClient) Write(identity *Identity, opts clients.WriteOpts) (*Identity, error) {
	resource, err := client.rc.Write(identity, opts)
	if err != nil {
		return nil, err
	}
	return resource.(*Identity), nil
}

func (client *identityClient) Delete(namespace, name string, opts clients.DeleteOpts) error {
	return client.rc.Delete(namespace, name, opts)
}

func (client *identityClient) List(namespace string, opts clients.ListOpts) ([]*Identity, error) {
	resourceList, err := client.rc.List(namespace, opts)
	if err != nil {
		return nil, err
	}
	return convertToIdentity(resourceList), nil
}

func (client *identityClient) Watch(namespace string, opts clients.WatchOpts) (<-chan []*Identity, <-chan error, error) {
	resourcesChan, errs, initErr := client.rc.Watch(namespace, opts)
	if initErr != nil {
		return nil, nil, initErr
	}
	identitysChan := make(chan []*Identity)
	go func() {
		for {
			select {
			case resourceList := <-resourcesChan:
				identitysChan <- convertToIdentity(resourceList)
			}
		}
	}()
	return identitysChan, errs, nil
}

func convertToIdentity(resources []resources.Resource) []*Identity {
	var identityList []*Identity
	for _, resource := range resources {
		identityList = append(identityList, resource.(*Identity))
	}
	return identityList
}

// Kubernetes Adapter for Identity

func (o *Identity) GetObjectKind() schema.ObjectKind {
	t := IdentityCrd.TypeMeta()
	return &t
}

func (o *Identity) DeepCopyObject() runtime.Object {
	return resources.Clone(o).(*Identity)
}

var IdentityCrd = crd.NewCrd("authorization.solo.io",
	"identities",
	"authorization.solo.io",
	"v1",
	"Identity",
	"identity",
	&Identity{})
