package thirdparty

import (
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/factory"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/kube/crd"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

type ThirdPartyResource interface {
	GetData() string
	IsSecret() bool
}

// TODO: modify as needed to populate additional fields
func NewArtifact(namespace, name, data string) *Artifact {
	return &Artifact{
		Metadata: core.Metadata{
			Name:      name,
			Namespace: namespace,
		},
		Data: data,
	}
}

func (r *Artifact) SetMetadata(meta core.Metadata) {
	r.Metadata = meta
}

func (r *Artifact) IsSecret() bool {
	return false
}

func (r *Secret) SetMetadata(meta core.Metadata) {
	r.Metadata = meta
}

func (r *Secret) IsSecret() bool {
	return true
}

var _ ThirdPartyResource = &Artifact{}
var _ ThirdPartyResource = &Secret{}

type ThirdPartyResourceClient interface {
	Register() error
	Read(namespace, name string, opts clients.ReadOpts) (*Artifact, error)
	Write(resource *Artifact, opts clients.WriteOpts) (*Artifact, error)
	Delete(namespace, name string, opts clients.DeleteOpts) error
	List(namespace string, opts clients.ListOpts) ([]*Artifact, error)
	Watch(namespace string, opts clients.WatchOpts) (<-chan []*Artifact, <-chan error, error)
}

type thirdPartyResourceClient struct {
	rc clients.ResourceClient
}

func NewArtifactClient(factory *factory.ResourceClientFactory) ThirdPartyResourceClient {
	return &thirdPartyResourceClient{
		rc: factory.NewResourceClient(&Artifact{}),
	}
}

func (client *thirdPartyResourceClient) Register() error {
	return client.rc.Register()
}

func (client *thirdPartyResourceClient) Read(namespace, name string, opts clients.ReadOpts) (*Artifact, error) {
	resource, err := client.rc.Read(namespace, name, opts)
	if err != nil {
		return nil, err
	}
	return resource.(*Artifact), nil
}

func (client *thirdPartyResourceClient) Write(artifact *Artifact, opts clients.WriteOpts) (*Artifact, error) {
	resource, err := client.rc.Write(artifact, opts)
	if err != nil {
		return nil, err
	}
	return resource.(*Artifact), nil
}

func (client *thirdPartyResourceClient) Delete(namespace, name string, opts clients.DeleteOpts) error {
	return client.rc.Delete(namespace, name, opts)
}

func (client *thirdPartyResourceClient) List(namespace string, opts clients.ListOpts) ([]*Artifact, error) {
	resourceList, err := client.rc.List(namespace, opts)
	if err != nil {
		return nil, err
	}
	return convertToArtifact(resourceList), nil
}

func (client *thirdPartyResourceClient) Watch(namespace string, opts clients.WatchOpts) (<-chan []*Artifact, <-chan error, error) {
	resourcesChan, errs, initErr := client.rc.Watch(namespace, opts)
	if initErr != nil {
		return nil, nil, initErr
	}
	artifactsChan := make(chan []*Artifact)
	go func() {
		for {
			select {
			case resourceList := <-resourcesChan:
				artifactsChan <- convertToArtifact(resourceList)
			}
		}
	}()
	return artifactsChan, errs, nil
}

func convertToArtifact(resources []ThirdPartyResource) []*Artifact {
	var artifactList []*Artifact
	for _, resource := range resources {
		artifactList = append(artifactList, resource.(*Artifact))
	}
	return artifactList
}

// Kubernetes Adapter for Artifact

func (o *Artifact) GetObjectKind() schema.ObjectKind {
	t := ArtifactCrd.TypeMeta()
	return &t
}

func (o *Artifact) DeepCopyObject() runtime.Object {
	return resources.Clone(o).(*Artifact)
}

var ArtifactCrd = crd.NewCrd("testing.solo.io",
	"fakes",
	"testing.solo.io",
	"v1",
	"Artifact",
	"fk",
	&Artifact{})
