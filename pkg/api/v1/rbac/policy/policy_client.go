package policy

import (
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/factory"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/kube/crd"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// TODO: modify as needed to populate additional fields
func NewPolicy(namespace, name string) *Policy {
	return &Policy{
		Metadata: core.Metadata{
			Name:      name,
			Namespace: namespace,
		},
	}
}

func (r *Policy) SetStatus(status core.Status) {
	r.Status = status
}

func (r *Policy) SetMetadata(meta core.Metadata) {
	r.Metadata = meta
}

var _ resources.Resource = &Policy{}

type PolicyClient interface {
	Register() error
	Read(namespace, name string, opts clients.ReadOpts) (*Policy, error)
	Write(resource *Policy, opts clients.WriteOpts) (*Policy, error)
	Delete(namespace, name string, opts clients.DeleteOpts) error
	List(namespace string, opts clients.ListOpts) ([]*Policy, error)
	Watch(namespace string, opts clients.WatchOpts) (<-chan []*Policy, <-chan error, error)
}

type policyClient struct {
	rc clients.ResourceClient
}

func NewPolicyClient(factory *factory.ResourceClientFactory) PolicyClient {
	return &policyClient{
		rc: factory.NewResourceClient(&Policy{}),
	}
}

func (client *policyClient) Register() error {
	return client.rc.Register()
}

func (client *policyClient) Read(namespace, name string, opts clients.ReadOpts) (*Policy, error) {
	resource, err := client.rc.Read(namespace, name, opts)
	if err != nil {
		return nil, err
	}
	return resource.(*Policy), nil
}

func (client *policyClient) Write(policy *Policy, opts clients.WriteOpts) (*Policy, error) {
	resource, err := client.rc.Write(policy, opts)
	if err != nil {
		return nil, err
	}
	return resource.(*Policy), nil
}

func (client *policyClient) Delete(namespace, name string, opts clients.DeleteOpts) error {
	return client.rc.Delete(namespace, name, opts)
}

func (client *policyClient) List(namespace string, opts clients.ListOpts) ([]*Policy, error) {
	resourceList, err := client.rc.List(namespace, opts)
	if err != nil {
		return nil, err
	}
	return convertToPolicy(resourceList), nil
}

func (client *policyClient) Watch(namespace string, opts clients.WatchOpts) (<-chan []*Policy, <-chan error, error) {
	resourcesChan, errs, initErr := client.rc.Watch(namespace, opts)
	if initErr != nil {
		return nil, nil, initErr
	}
	policysChan := make(chan []*Policy)
	go func() {
		for {
			select {
			case resourceList := <-resourcesChan:
				policysChan <- convertToPolicy(resourceList)
			}
		}
	}()
	return policysChan, errs, nil
}

func convertToPolicy(resources []resources.Resource) []*Policy {
	var policyList []*Policy
	for _, resource := range resources {
		policyList = append(policyList, resource.(*Policy))
	}
	return policyList
}

// Kubernetes Adapter for Policy

func (o *Policy) GetObjectKind() schema.ObjectKind {
	t := PolicyCrd.TypeMeta()
	return &t
}

func (o *Policy) DeepCopyObject() runtime.Object {
	return resources.Clone(o).(*Policy)
}

var PolicyCrd = crd.NewCrd("authorization.solo.io",
	"policies",
	"authorization.solo.io",
	"v1",
	"Policy",
	"policy",
	&Policy{})
