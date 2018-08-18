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

type UpstreamList []*Upstream

// namespace is optional, if left empty, names can collide if the list contains more than one with the same name
func (list UpstreamList) Find(namespace, name string) (*Upstream, error) {
	for _, upstream := range list {
		if upstream.Metadata.Name == name {
			if namespace == "" || upstream.Metadata.Namespace == namespace {
				return upstream, nil
			}
		}
	}
	return nil, errors.Errorf("list did not find upstream %v.%v", namespace, name)
}

func (list UpstreamList) AsResources() []resources.Resource {
	var ress []resources.Resource
	for _, upstream := range list {
		ress = append(ress, upstream)
	}
	return ress
}

func (list UpstreamList) AsInputResources() []resources.InputResource {
	var ress []resources.InputResource
	for _, upstream := range list {
		ress = append(ress, upstream)
	}
	return ress
}

func (list UpstreamList) Names() []string {
	var names []string
	for _, upstream := range list {
		names = append(names, upstream.Metadata.Name)
	}
	return names
}

func (list UpstreamList) NamespacesDotNames() []string {
	var names []string
	for _, upstream := range list {
		names = append(names, upstream.Metadata.Namespace+"."+upstream.Metadata.Name)
	}
	return names
}

func (list UpstreamList) Sort() {
	sort.SliceStable(list, func(i, j int) bool {
		return list[i].Metadata.Less(list[j].Metadata)
	})
}

var _ resources.Resource = &Upstream{}

type UpstreamClient interface {
	BaseClient() clients.ResourceClient
	Register() error
	Read(namespace, name string, opts clients.ReadOpts) (*Upstream, error)
	Write(resource *Upstream, opts clients.WriteOpts) (*Upstream, error)
	Delete(namespace, name string, opts clients.DeleteOpts) error
	List(namespace string, opts clients.ListOpts) (UpstreamList, error)
	Watch(namespace string, opts clients.WatchOpts) (<-chan UpstreamList, <-chan error, error)
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

func (client *upstreamClient) BaseClient() clients.ResourceClient {
	return client.rc
}

func (client *upstreamClient) Register() error {
	return client.rc.Register()
}

func (client *upstreamClient) Read(namespace, name string, opts clients.ReadOpts) (*Upstream, error) {
	opts = opts.WithDefaults()
	resource, err := client.rc.Read(namespace, name, opts)
	if err != nil {
		return nil, err
	}
	return resource.(*Upstream), nil
}

func (client *upstreamClient) Write(upstream *Upstream, opts clients.WriteOpts) (*Upstream, error) {
	opts = opts.WithDefaults()
	resource, err := client.rc.Write(upstream, opts)
	if err != nil {
		return nil, err
	}
	return resource.(*Upstream), nil
}

func (client *upstreamClient) Delete(namespace, name string, opts clients.DeleteOpts) error {
	opts = opts.WithDefaults()
	return client.rc.Delete(namespace, name, opts)
}

func (client *upstreamClient) List(namespace string, opts clients.ListOpts) (UpstreamList, error) {
	opts = opts.WithDefaults()
	resourceList, err := client.rc.List(namespace, opts)
	if err != nil {
		return nil, err
	}
	return convertToUpstream(resourceList), nil
}

func (client *upstreamClient) Watch(namespace string, opts clients.WatchOpts) (<-chan UpstreamList, <-chan error, error) {
	opts = opts.WithDefaults()
	resourcesChan, errs, initErr := client.rc.Watch(namespace, opts)
	if initErr != nil {
		return nil, nil, initErr
	}
	upstreamsChan := make(chan UpstreamList)
	go func() {
		for {
			select {
			case resourceList := <-resourcesChan:
				upstreamsChan <- convertToUpstream(resourceList)
			case <-opts.Ctx.Done():
				close(upstreamsChan)
				return
			}
		}
	}()
	return upstreamsChan, errs, nil
}

func convertToUpstream(resources []resources.Resource) UpstreamList {
	var upstreamList UpstreamList
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
