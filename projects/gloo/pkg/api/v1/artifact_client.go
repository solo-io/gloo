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
func NewArtifact(namespace, name string) *Artifact {
	return &Artifact{
		Metadata: core.Metadata{
			Name:      name,
			Namespace: namespace,
		},
	}
}

func (r *Artifact) SetMetadata(meta core.Metadata) {
	r.Metadata = meta
}

func (r *Artifact) SetData(data map[string]string) {
	r.Data = data
}

var _ resources.Resource = &Artifact{}

type ArtifactClient interface {
	Register() error
	Read(namespace, name string, opts clients.ReadOpts) (*Artifact, error)
	Write(resource *Artifact, opts clients.WriteOpts) (*Artifact, error)
	Delete(namespace, name string, opts clients.DeleteOpts) error
	List(namespace string, opts clients.ListOpts) ([]*Artifact, error)
	Watch(namespace string, opts clients.WatchOpts) (<-chan []*Artifact, <-chan error, error)
}

type artifactClient struct {
	rc clients.ResourceClient
}

func NewArtifactClient(rcFactory factory.ResourceClientFactory) (ArtifactClient, error) {
	rc, err := rcFactory.NewResourceClient(factory.NewResourceClientParams{
		ResourceType: &Artifact{},
	})
	if err != nil {
		return nil, errors.Wrapf(err, "creating base Artifact resource client")
	}
	return &artifactClient{
		rc: rc,
	}, nil
}

func (client *artifactClient) Register() error {
	return client.rc.Register()
}

func (client *artifactClient) Read(namespace, name string, opts clients.ReadOpts) (*Artifact, error) {
	resource, err := client.rc.Read(namespace, name, opts)
	if err != nil {
		return nil, err
	}
	return resource.(*Artifact), nil
}

func (client *artifactClient) Write(artifact *Artifact, opts clients.WriteOpts) (*Artifact, error) {
	resource, err := client.rc.Write(artifact, opts)
	if err != nil {
		return nil, err
	}
	return resource.(*Artifact), nil
}

func (client *artifactClient) Delete(namespace, name string, opts clients.DeleteOpts) error {
	return client.rc.Delete(namespace, name, opts)
}

func (client *artifactClient) List(namespace string, opts clients.ListOpts) ([]*Artifact, error) {
	resourceList, err := client.rc.List(namespace, opts)
	if err != nil {
		return nil, err
	}
	return convertToArtifact(resourceList), nil
}

func (client *artifactClient) Watch(namespace string, opts clients.WatchOpts) (<-chan []*Artifact, <-chan error, error) {
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
			case <-opts.Ctx.Done():
				close(artifactsChan)
				return
			}
		}
	}()
	return artifactsChan, errs, nil
}

func convertToArtifact(resources []resources.Resource) []*Artifact {
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

var ArtifactCrd = crd.NewCrd("gloo.solo.io",
	"artifacts",
	"gloo.solo.io",
	"v1",
	"Artifact",
	"art",
	&Artifact{})
