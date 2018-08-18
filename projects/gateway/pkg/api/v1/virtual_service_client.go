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
func NewVirtualService(namespace, name string) *VirtualService {
	return &VirtualService{
		Metadata: core.Metadata{
			Name:      name,
			Namespace: namespace,
		},
	}
}

func (r *VirtualService) SetStatus(status core.Status) {
	r.Status = status
}

func (r *VirtualService) SetMetadata(meta core.Metadata) {
	r.Metadata = meta
}

type VirtualServiceList []*VirtualService

// namespace is optional, if left empty, names can collide if the list contains more than one with the same name
func (list VirtualServiceList) Find(namespace, name string) (*VirtualService, error) {
	for _, virtualService := range list {
		if virtualService.Metadata.Name == name {
			if namespace == "" || virtualService.Metadata.Namespace == namespace {
				return virtualService, nil
			}
		}
	}
	return nil, errors.Errorf("list did not find virtualService %v.%v", namespace, name)
}

func (list VirtualServiceList) AsResources() resources.ResourceList {
	var ress resources.ResourceList
	for _, virtualService := range list {
		ress = append(ress, virtualService)
	}
	return ress
}

func (list VirtualServiceList) AsInputResources() resources.InputResourceList {
	var ress resources.InputResourceList
	for _, virtualService := range list {
		ress = append(ress, virtualService)
	}
	return ress
}

func (list VirtualServiceList) Names() []string {
	var names []string
	for _, virtualService := range list {
		names = append(names, virtualService.Metadata.Name)
	}
	return names
}

func (list VirtualServiceList) NamespacesDotNames() []string {
	var names []string
	for _, virtualService := range list {
		names = append(names, virtualService.Metadata.Namespace+"."+virtualService.Metadata.Name)
	}
	return names
}

func (list VirtualServiceList) Sort() {
	sort.SliceStable(list, func(i, j int) bool {
		return list[i].Metadata.Less(list[j].Metadata)
	})
}

var _ resources.Resource = &VirtualService{}

type VirtualServiceClient interface {
	BaseClient() clients.ResourceClient
	Register() error
	Read(namespace, name string, opts clients.ReadOpts) (*VirtualService, error)
	Write(resource *VirtualService, opts clients.WriteOpts) (*VirtualService, error)
	Delete(namespace, name string, opts clients.DeleteOpts) error
	List(namespace string, opts clients.ListOpts) (VirtualServiceList, error)
	Watch(namespace string, opts clients.WatchOpts) (<-chan VirtualServiceList, <-chan error, error)
}

type virtualServiceClient struct {
	rc clients.ResourceClient
}

func NewVirtualServiceClient(rcFactory factory.ResourceClientFactory) (VirtualServiceClient, error) {
	rc, err := rcFactory.NewResourceClient(factory.NewResourceClientParams{
		ResourceType: &VirtualService{},
	})
	if err != nil {
		return nil, errors.Wrapf(err, "creating base VirtualService resource client")
	}
	return &virtualServiceClient{
		rc: rc,
	}, nil
}

func (client *virtualServiceClient) BaseClient() clients.ResourceClient {
	return client.rc
}

func (client *virtualServiceClient) Register() error {
	return client.rc.Register()
}

func (client *virtualServiceClient) Read(namespace, name string, opts clients.ReadOpts) (*VirtualService, error) {
	opts = opts.WithDefaults()
	resource, err := client.rc.Read(namespace, name, opts)
	if err != nil {
		return nil, err
	}
	return resource.(*VirtualService), nil
}

func (client *virtualServiceClient) Write(virtualService *VirtualService, opts clients.WriteOpts) (*VirtualService, error) {
	opts = opts.WithDefaults()
	resource, err := client.rc.Write(virtualService, opts)
	if err != nil {
		return nil, err
	}
	return resource.(*VirtualService), nil
}

func (client *virtualServiceClient) Delete(namespace, name string, opts clients.DeleteOpts) error {
	opts = opts.WithDefaults()
	return client.rc.Delete(namespace, name, opts)
}

func (client *virtualServiceClient) List(namespace string, opts clients.ListOpts) (VirtualServiceList, error) {
	opts = opts.WithDefaults()
	resourceList, err := client.rc.List(namespace, opts)
	if err != nil {
		return nil, err
	}
	return convertToVirtualService(resourceList), nil
}

func (client *virtualServiceClient) Watch(namespace string, opts clients.WatchOpts) (<-chan VirtualServiceList, <-chan error, error) {
	opts = opts.WithDefaults()
	resourcesChan, errs, initErr := client.rc.Watch(namespace, opts)
	if initErr != nil {
		return nil, nil, initErr
	}
	virtualServicesChan := make(chan VirtualServiceList)
	go func() {
		for {
			select {
			case resourceList := <-resourcesChan:
				virtualServicesChan <- convertToVirtualService(resourceList)
			case <-opts.Ctx.Done():
				close(virtualServicesChan)
				return
			}
		}
	}()
	return virtualServicesChan, errs, nil
}

func convertToVirtualService(resources resources.ResourceList) VirtualServiceList {
	var virtualServiceList VirtualServiceList
	for _, resource := range resources {
		virtualServiceList = append(virtualServiceList, resource.(*VirtualService))
	}
	return virtualServiceList
}

// Kubernetes Adapter for VirtualService

func (o *VirtualService) GetObjectKind() schema.ObjectKind {
	t := VirtualServiceCrd.TypeMeta()
	return &t
}

func (o *VirtualService) DeepCopyObject() runtime.Object {
	return resources.Clone(o).(*VirtualService)
}

var VirtualServiceCrd = crd.NewCrd("gateway.solo.io",
	"virtualservices",
	"gateway.solo.io",
	"v1",
	"VirtualService",
	"vs",
	&VirtualService{})
