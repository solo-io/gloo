package v1

import (
	"sort"

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

type SecretList []*Secret

// namespace is optional, if left empty, names can collide if the list contains more than one with the same name
func (list SecretList) Find(namespace, name string) (*Secret, error) {
	for _, secret := range list {
		if secret.Metadata.Name == name {
			if namespace == "" || secret.Metadata.Namespace == namespace {
				return secret, nil
			}
		}
	}
	return nil, errors.Errorf("list did not find secret %v.%v", namespace, name)
}

func (list SecretList) AsResources() resources.ResourceList {
	var ress resources.ResourceList
	for _, secret := range list {
		ress = append(ress, secret)
	}
	return ress
}

func (list SecretList) Names() []string {
	var names []string
	for _, secret := range list {
		names = append(names, secret.Metadata.Name)
	}
	return names
}

func (list SecretList) NamespacesDotNames() []string {
	var names []string
	for _, secret := range list {
		names = append(names, secret.Metadata.Namespace+"."+secret.Metadata.Name)
	}
	return names
}

func (list SecretList) Sort() {
	sort.SliceStable(list, func(i, j int) bool {
		return list[i].Metadata.Less(list[j].Metadata)
	})
}

var _ resources.Resource = &Secret{}

type SecretClient interface {
	BaseClient() clients.ResourceClient
	Register() error
	Read(namespace, name string, opts clients.ReadOpts) (*Secret, error)
	Write(resource *Secret, opts clients.WriteOpts) (*Secret, error)
	Delete(namespace, name string, opts clients.DeleteOpts) error
	List(namespace string, opts clients.ListOpts) (SecretList, error)
	Watch(namespace string, opts clients.WatchOpts) (<-chan SecretList, <-chan error, error)
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

func (client *secretClient) BaseClient() clients.ResourceClient {
	return client.rc
}

func (client *secretClient) Register() error {
	return client.rc.Register()
}

func (client *secretClient) Read(namespace, name string, opts clients.ReadOpts) (*Secret, error) {
	opts = opts.WithDefaults()
	resource, err := client.rc.Read(namespace, name, opts)
	if err != nil {
		return nil, err
	}
	return resource.(*Secret), nil
}

func (client *secretClient) Write(secret *Secret, opts clients.WriteOpts) (*Secret, error) {
	opts = opts.WithDefaults()
	resource, err := client.rc.Write(secret, opts)
	if err != nil {
		return nil, err
	}
	return resource.(*Secret), nil
}

func (client *secretClient) Delete(namespace, name string, opts clients.DeleteOpts) error {
	opts = opts.WithDefaults()
	return client.rc.Delete(namespace, name, opts)
}

func (client *secretClient) List(namespace string, opts clients.ListOpts) (SecretList, error) {
	opts = opts.WithDefaults()
	resourceList, err := client.rc.List(namespace, opts)
	if err != nil {
		return nil, err
	}
	return convertToSecret(resourceList), nil
}

func (client *secretClient) Watch(namespace string, opts clients.WatchOpts) (<-chan SecretList, <-chan error, error) {
	opts = opts.WithDefaults()
	resourcesChan, errs, initErr := client.rc.Watch(namespace, opts)
	if initErr != nil {
		return nil, nil, initErr
	}
	secretsChan := make(chan SecretList)
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

func convertToSecret(resources resources.ResourceList) SecretList {
	var secretList SecretList
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
