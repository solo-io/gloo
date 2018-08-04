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
func NewRole(namespace, name string) *Role {
	return &Role{
		Metadata: core.Metadata{
			Name:      name,
			Namespace: namespace,
		},
	}
}

func (r *Role) SetStatus(status core.Status) {
	r.Status = status
}

func (r *Role) SetMetadata(meta core.Metadata) {
	r.Metadata = meta
}

var _ resources.Resource = &Role{}

type RoleClient interface {
	Register() error
	Read(namespace, name string, opts clients.ReadOpts) (*Role, error)
	Write(resource *Role, opts clients.WriteOpts) (*Role, error)
	Delete(namespace, name string, opts clients.DeleteOpts) error
	List(namespace string, opts clients.ListOpts) ([]*Role, error)
	Watch(namespace string, opts clients.WatchOpts) (<-chan []*Role, <-chan error, error)
}

type roleClient struct {
	rc clients.ResourceClient
}

func NewRoleClient(rcFactory factory.ResourceClientFactory) (RoleClient, error) {
	rc, err := rcFactory.NewResourceClient(factory.NewResourceClientParams{
		ResourceType: &Role{},
	})
	if err != nil {
		return nil, errors.Wrapf(err, "creating base Role resource client")
	}
	return &roleClient{
		rc: rc,
	}, nil
}

func (client *roleClient) Register() error {
	return client.rc.Register()
}

func (client *roleClient) Read(namespace, name string, opts clients.ReadOpts) (*Role, error) {
	resource, err := client.rc.Read(namespace, name, opts)
	if err != nil {
		return nil, err
	}
	return resource.(*Role), nil
}

func (client *roleClient) Write(role *Role, opts clients.WriteOpts) (*Role, error) {
	resource, err := client.rc.Write(role, opts)
	if err != nil {
		return nil, err
	}
	return resource.(*Role), nil
}

func (client *roleClient) Delete(namespace, name string, opts clients.DeleteOpts) error {
	return client.rc.Delete(namespace, name, opts)
}

func (client *roleClient) List(namespace string, opts clients.ListOpts) ([]*Role, error) {
	resourceList, err := client.rc.List(namespace, opts)
	if err != nil {
		return nil, err
	}
	return convertToRole(resourceList), nil
}

func (client *roleClient) Watch(namespace string, opts clients.WatchOpts) (<-chan []*Role, <-chan error, error) {
	resourcesChan, errs, initErr := client.rc.Watch(namespace, opts)
	if initErr != nil {
		return nil, nil, initErr
	}
	rolesChan := make(chan []*Role)
	go func() {
		for {
			select {
			case resourceList := <-resourcesChan:
				rolesChan <- convertToRole(resourceList)
			}
		}
	}()
	return rolesChan, errs, nil
}

func convertToRole(resources []resources.Resource) []*Role {
	var roleList []*Role
	for _, resource := range resources {
		roleList = append(roleList, resource.(*Role))
	}
	return roleList
}

// Kubernetes Adapter for Role

func (o *Role) GetObjectKind() schema.ObjectKind {
	t := RoleCrd.TypeMeta()
	return &t
}

func (o *Role) DeepCopyObject() runtime.Object {
	return resources.Clone(o).(*Role)
}

var RoleCrd = crd.NewCrd("gloo.solo.io",
	"mocks",
	"gloo.solo.io",
	"v1",
	"Role",
	"mk",
	&Role{})
