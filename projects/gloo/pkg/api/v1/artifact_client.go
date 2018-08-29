package v1

import (
	"sort"

	"github.com/gogo/protobuf/proto"
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

type ArtifactList []*Artifact
type ArtifactsByNamespace map[string]ArtifactList

// namespace is optional, if left empty, names can collide if the list contains more than one with the same name
func (list ArtifactList) Find(namespace, name string) (*Artifact, error) {
	for _, artifact := range list {
		if artifact.Metadata.Name == name {
			if namespace == "" || artifact.Metadata.Namespace == namespace {
				return artifact, nil
			}
		}
	}
	return nil, errors.Errorf("list did not find artifact %v.%v", namespace, name)
}

func (list ArtifactList) AsResources() resources.ResourceList {
	var ress resources.ResourceList
	for _, artifact := range list {
		ress = append(ress, artifact)
	}
	return ress
}

func (list ArtifactList) Names() []string {
	var names []string
	for _, artifact := range list {
		names = append(names, artifact.Metadata.Name)
	}
	return names
}

func (list ArtifactList) NamespacesDotNames() []string {
	var names []string
	for _, artifact := range list {
		names = append(names, artifact.Metadata.Namespace+"."+artifact.Metadata.Name)
	}
	return names
}

func (list ArtifactList) Sort() {
	sort.SliceStable(list, func(i, j int) bool {
		return list[i].Metadata.Less(list[j].Metadata)
	})
}

func (list ArtifactList) Clone() ArtifactList {
	var artifactList ArtifactList
	for _, artifact := range list {
		artifactList = append(artifactList, proto.Clone(artifact).(*Artifact))
	}
	return artifactList
}

func (list ArtifactList) ByNamespace() ArtifactsByNamespace {
	byNamespace := make(ArtifactsByNamespace)
	for _, artifact := range list {
		byNamespace.Add(artifact)
	}
	return byNamespace
}

func (byNamespace ArtifactsByNamespace) Add(artifact ...*Artifact) {
	for _, item := range artifact {
		byNamespace[item.Metadata.Namespace] = append(byNamespace[item.Metadata.Namespace], item)
	}
}

func (byNamespace ArtifactsByNamespace) Clear(namespace string) {
	delete(byNamespace, namespace)
}

func (byNamespace ArtifactsByNamespace) List() ArtifactList {
	var list ArtifactList
	for _, artifactList := range byNamespace {
		list = append(list, artifactList...)
	}
	list.Sort()
	return list
}

func (byNamespace ArtifactsByNamespace) Clone() ArtifactsByNamespace {
	return byNamespace.List().Clone().ByNamespace()
}

var _ resources.Resource = &Artifact{}

type ArtifactClient interface {
	BaseClient() clients.ResourceClient
	Register() error
	Read(namespace, name string, opts clients.ReadOpts) (*Artifact, error)
	Write(resource *Artifact, opts clients.WriteOpts) (*Artifact, error)
	Delete(namespace, name string, opts clients.DeleteOpts) error
	List(namespace string, opts clients.ListOpts) (ArtifactList, error)
	Watch(namespace string, opts clients.WatchOpts) (<-chan ArtifactList, <-chan error, error)
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

func (client *artifactClient) BaseClient() clients.ResourceClient {
	return client.rc
}

func (client *artifactClient) Register() error {
	return client.rc.Register()
}

func (client *artifactClient) Read(namespace, name string, opts clients.ReadOpts) (*Artifact, error) {
	opts = opts.WithDefaults()
	resource, err := client.rc.Read(namespace, name, opts)
	if err != nil {
		return nil, err
	}
	return resource.(*Artifact), nil
}

func (client *artifactClient) Write(artifact *Artifact, opts clients.WriteOpts) (*Artifact, error) {
	opts = opts.WithDefaults()
	resource, err := client.rc.Write(artifact, opts)
	if err != nil {
		return nil, err
	}
	return resource.(*Artifact), nil
}

func (client *artifactClient) Delete(namespace, name string, opts clients.DeleteOpts) error {
	opts = opts.WithDefaults()
	return client.rc.Delete(namespace, name, opts)
}

func (client *artifactClient) List(namespace string, opts clients.ListOpts) (ArtifactList, error) {
	opts = opts.WithDefaults()
	resourceList, err := client.rc.List(namespace, opts)
	if err != nil {
		return nil, err
	}
	return convertToArtifact(resourceList), nil
}

func (client *artifactClient) Watch(namespace string, opts clients.WatchOpts) (<-chan ArtifactList, <-chan error, error) {
	opts = opts.WithDefaults()
	resourcesChan, errs, initErr := client.rc.Watch(namespace, opts)
	if initErr != nil {
		return nil, nil, initErr
	}
	artifactsChan := make(chan ArtifactList)
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

func convertToArtifact(resources resources.ResourceList) ArtifactList {
	var artifactList ArtifactList
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
