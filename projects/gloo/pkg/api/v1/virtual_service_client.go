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

var _ resources.Resource = &VirtualService{}

type VirtualServiceClient interface {
	Register() error
	Read(namespace, name string, opts clients.ReadOpts) (*VirtualService, error)
	Write(resource *VirtualService, opts clients.WriteOpts) (*VirtualService, error)
	Delete(namespace, name string, opts clients.DeleteOpts) error
	List(namespace string, opts clients.ListOpts) ([]*VirtualService, error)
	Watch(namespace string, opts clients.WatchOpts) (<-chan []*VirtualService, <-chan error, error)
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

func (client *virtualServiceClient) Register() error {
	return client.rc.Register()
}

func (client *virtualServiceClient) Read(namespace, name string, opts clients.ReadOpts) (*VirtualService, error) {
	resource, err := client.rc.Read(namespace, name, opts)
	if err != nil {
		return nil, err
	}
	return resource.(*VirtualService), nil
}

func (client *virtualServiceClient) Write(virtualService *VirtualService, opts clients.WriteOpts) (*VirtualService, error) {
	resource, err := client.rc.Write(virtualService, opts)
	if err != nil {
		return nil, err
	}
	return resource.(*VirtualService), nil
}

func (client *virtualServiceClient) Delete(namespace, name string, opts clients.DeleteOpts) error {
	return client.rc.Delete(namespace, name, opts)
}

func (client *virtualServiceClient) List(namespace string, opts clients.ListOpts) ([]*VirtualService, error) {
	resourceList, err := client.rc.List(namespace, opts)
	if err != nil {
		return nil, err
	}
	return convertToVirtualService(resourceList), nil
}

func (client *virtualServiceClient) Watch(namespace string, opts clients.WatchOpts) (<-chan []*VirtualService, <-chan error, error) {
	resourcesChan, errs, initErr := client.rc.Watch(namespace, opts)
	if initErr != nil {
		return nil, nil, initErr
	}
	virtualServicesChan := make(chan []*VirtualService)
	go func() {
		for {
			select {
			case resourceList := <-resourcesChan:
				virtualServicesChan <- convertToVirtualService(resourceList)
			}
		}
	}()
	return virtualServicesChan, errs, nil
}

func convertToVirtualService(resources []resources.Resource) []*VirtualService {
	var virtualServiceList []*VirtualService
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

var VirtualServiceCrd = crd.NewCrd("gloo.solo.io",
	"virtualservices",
	"gloo.solo.io",
	"v1",
	"VirtualService",
	"vs",
	&VirtualService{})
