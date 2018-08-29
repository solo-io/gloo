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
func NewResolverMap(namespace, name string) *ResolverMap {
	return &ResolverMap{
		Metadata: core.Metadata{
			Name:      name,
			Namespace: namespace,
		},
	}
}

func (r *ResolverMap) SetStatus(status core.Status) {
	r.Status = status
}

func (r *ResolverMap) SetMetadata(meta core.Metadata) {
	r.Metadata = meta
}

type ResolverMapList []*ResolverMap
type ResolverMapsByNamespace map[string]ResolverMapList

// namespace is optional, if left empty, names can collide if the list contains more than one with the same name
func (list ResolverMapList) Find(namespace, name string) (*ResolverMap, error) {
	for _, resolverMap := range list {
		if resolverMap.Metadata.Name == name {
			if namespace == "" || resolverMap.Metadata.Namespace == namespace {
				return resolverMap, nil
			}
		}
	}
	return nil, errors.Errorf("list did not find resolverMap %v.%v", namespace, name)
}

func (list ResolverMapList) AsResources() resources.ResourceList {
	var ress resources.ResourceList
	for _, resolverMap := range list {
		ress = append(ress, resolverMap)
	}
	return ress
}

func (list ResolverMapList) AsInputResources() resources.InputResourceList {
	var ress resources.InputResourceList
	for _, resolverMap := range list {
		ress = append(ress, resolverMap)
	}
	return ress
}

func (list ResolverMapList) Names() []string {
	var names []string
	for _, resolverMap := range list {
		names = append(names, resolverMap.Metadata.Name)
	}
	return names
}

func (list ResolverMapList) NamespacesDotNames() []string {
	var names []string
	for _, resolverMap := range list {
		names = append(names, resolverMap.Metadata.Namespace+"."+resolverMap.Metadata.Name)
	}
	return names
}

func (list ResolverMapList) Sort() {
	sort.SliceStable(list, func(i, j int) bool {
		return list[i].Metadata.Less(list[j].Metadata)
	})
}

func (list ResolverMapList) Clone() ResolverMapList {
	var resolverMapList ResolverMapList
	for _, resolverMap := range list {
		resolverMapList = append(resolverMapList, proto.Clone(resolverMap).(*ResolverMap))
	}
	return resolverMapList
}

func (list ResolverMapList) ByNamespace() ResolverMapsByNamespace {
	byNamespace := make(ResolverMapsByNamespace)
	for _, resolverMap := range list {
		byNamespace.Add(resolverMap)
	}
	return byNamespace
}

func (byNamespace ResolverMapsByNamespace) Add(resolverMap ...*ResolverMap) {
	for _, item := range resolverMap {
		byNamespace[item.Metadata.Namespace] = append(byNamespace[item.Metadata.Namespace], item)
	}
}

func (byNamespace ResolverMapsByNamespace) Clear(namespace string) {
	delete(byNamespace, namespace)
}

func (byNamespace ResolverMapsByNamespace) List() ResolverMapList {
	var list ResolverMapList
	for _, resolverMapList := range byNamespace {
		list = append(list, resolverMapList...)
	}
	list.Sort()
	return list
}

func (byNamespace ResolverMapsByNamespace) Clone() ResolverMapsByNamespace {
	return byNamespace.List().Clone().ByNamespace()
}

var _ resources.Resource = &ResolverMap{}

type ResolverMapClient interface {
	BaseClient() clients.ResourceClient
	Register() error
	Read(namespace, name string, opts clients.ReadOpts) (*ResolverMap, error)
	Write(resource *ResolverMap, opts clients.WriteOpts) (*ResolverMap, error)
	Delete(namespace, name string, opts clients.DeleteOpts) error
	List(namespace string, opts clients.ListOpts) (ResolverMapList, error)
	Watch(namespace string, opts clients.WatchOpts) (<-chan ResolverMapList, <-chan error, error)
}

type resolverMapClient struct {
	rc clients.ResourceClient
}

func NewResolverMapClient(rcFactory factory.ResourceClientFactory) (ResolverMapClient, error) {
	rc, err := rcFactory.NewResourceClient(factory.NewResourceClientParams{
		ResourceType: &ResolverMap{},
	})
	if err != nil {
		return nil, errors.Wrapf(err, "creating base ResolverMap resource client")
	}
	return &resolverMapClient{
		rc: rc,
	}, nil
}

func (client *resolverMapClient) BaseClient() clients.ResourceClient {
	return client.rc
}

func (client *resolverMapClient) Register() error {
	return client.rc.Register()
}

func (client *resolverMapClient) Read(namespace, name string, opts clients.ReadOpts) (*ResolverMap, error) {
	opts = opts.WithDefaults()
	resource, err := client.rc.Read(namespace, name, opts)
	if err != nil {
		return nil, err
	}
	return resource.(*ResolverMap), nil
}

func (client *resolverMapClient) Write(resolverMap *ResolverMap, opts clients.WriteOpts) (*ResolverMap, error) {
	opts = opts.WithDefaults()
	resource, err := client.rc.Write(resolverMap, opts)
	if err != nil {
		return nil, err
	}
	return resource.(*ResolverMap), nil
}

func (client *resolverMapClient) Delete(namespace, name string, opts clients.DeleteOpts) error {
	opts = opts.WithDefaults()
	return client.rc.Delete(namespace, name, opts)
}

func (client *resolverMapClient) List(namespace string, opts clients.ListOpts) (ResolverMapList, error) {
	opts = opts.WithDefaults()
	resourceList, err := client.rc.List(namespace, opts)
	if err != nil {
		return nil, err
	}
	return convertToResolverMap(resourceList), nil
}

func (client *resolverMapClient) Watch(namespace string, opts clients.WatchOpts) (<-chan ResolverMapList, <-chan error, error) {
	opts = opts.WithDefaults()
	resourcesChan, errs, initErr := client.rc.Watch(namespace, opts)
	if initErr != nil {
		return nil, nil, initErr
	}
	resolverMapsChan := make(chan ResolverMapList)
	go func() {
		for {
			select {
			case resourceList := <-resourcesChan:
				resolverMapsChan <- convertToResolverMap(resourceList)
			case <-opts.Ctx.Done():
				close(resolverMapsChan)
				return
			}
		}
	}()
	return resolverMapsChan, errs, nil
}

func convertToResolverMap(resources resources.ResourceList) ResolverMapList {
	var resolverMapList ResolverMapList
	for _, resource := range resources {
		resolverMapList = append(resolverMapList, resource.(*ResolverMap))
	}
	return resolverMapList
}

// Kubernetes Adapter for ResolverMap

func (o *ResolverMap) GetObjectKind() schema.ObjectKind {
	t := ResolverMapCrd.TypeMeta()
	return &t
}

func (o *ResolverMap) DeepCopyObject() runtime.Object {
	return resources.Clone(o).(*ResolverMap)
}

var ResolverMapCrd = crd.NewCrd("sqoop.solo.io",
	"resolvermaps",
	"sqoop.solo.io",
	"v1",
	"ResolverMap",
	"rm",
	&ResolverMap{})
