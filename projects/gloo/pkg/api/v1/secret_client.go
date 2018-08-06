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
func NewSecret(namespace, name string) *Secret {
	return &Secret{
		Metadata: core.Metadata{
			Name:      name,
			Namespace: namespace,
		},
	}
}

func (r *Secret) SetMetadata(meta core.Metadata) {
	r.Metadata = meta
}

func (r *Secret) SetData(data map[string]string) {
	r.Data = data
}

var _ resources.Resource = &Secret{}

type SecretClient interface {
	Register() error
	Read(namespace, name string, opts clients.ReadOpts) (*Secret, error)
	Write(resource *Secret, opts clients.WriteOpts) (*Secret, error)
	Delete(namespace, name string, opts clients.DeleteOpts) error
	List(namespace string, opts clients.ListOpts) ([]*Secret, error)
	Watch(namespace string, opts clients.WatchOpts) (<-chan []*Secret, <-chan error, error)
}

type secretClient struct {
	rc clients.ResourceClient
}

func NewSecretClient(rcFactory factory.ResourceClientFactory) (SecretClient, error) {
	rc, err := rcFactory.NewResourceClient(factory.NewResourceClientParams{
		ResourceType: &Secret{},
	})
	if err != nil {
		return nil, errors.Wrapf(err, "creating base Secret resource client")
	}
	return &secretClient{
		rc: rc,
	}, nil
}

func (client *secretClient) Register() error {
	return client.rc.Register()
}

func (client *secretClient) Read(namespace, name string, opts clients.ReadOpts) (*Secret, error) {
	resource, err := client.rc.Read(namespace, name, opts)
	if err != nil {
		return nil, err
	}
	return resource.(*Secret), nil
}

func (client *secretClient) Write(secret *Secret, opts clients.WriteOpts) (*Secret, error) {
	resource, err := client.rc.Write(secret, opts)
	if err != nil {
		return nil, err
	}
	return resource.(*Secret), nil
}

func (client *secretClient) Delete(namespace, name string, opts clients.DeleteOpts) error {
	return client.rc.Delete(namespace, name, opts)
}

func (client *secretClient) List(namespace string, opts clients.ListOpts) ([]*Secret, error) {
	resourceList, err := client.rc.List(namespace, opts)
	if err != nil {
		return nil, err
	}
	return convertToSecret(resourceList), nil
}

func (client *secretClient) Watch(namespace string, opts clients.WatchOpts) (<-chan []*Secret, <-chan error, error) {
	resourcesChan, errs, initErr := client.rc.Watch(namespace, opts)
	if initErr != nil {
		return nil, nil, initErr
	}
	secretsChan := make(chan []*Secret)
	go func() {
		for {
			select {
			case resourceList := <-resourcesChan:
				secretsChan <- convertToSecret(resourceList)
			case <-opts.Ctx.Done():
				close(secretsChan)
				return
			}
		}
	}()
	return secretsChan, errs, nil
}

func convertToSecret(resources []resources.Resource) []*Secret {
	var secretList []*Secret
	for _, resource := range resources {
		secretList = append(secretList, resource.(*Secret))
	}
	return secretList
}

// Kubernetes Adapter for Secret

func (o *Secret) GetObjectKind() schema.ObjectKind {
	t := SecretCrd.TypeMeta()
	return &t
}

func (o *Secret) DeepCopyObject() runtime.Object {
	return resources.Clone(o).(*Secret)
}

var SecretCrd = crd.NewCrd("gloo.solo.io",
	"secrets",
	"gloo.solo.io",
	"v1",
	"Secret",
	"sec",
	&Secret{})
