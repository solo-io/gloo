package crd

import (
	"fmt"

	"reflect"

	"github.com/pkg/errors"
	"github.com/solo-io/glue-storage"
	crdclientset "github.com/solo-io/glue-storage/crd/client/clientset/versioned"
	crdv1 "github.com/solo-io/glue-storage/crd/solo.io/v1"
	"github.com/solo-io/glue/pkg/api/types/v1"
	"istio.io/istio/pkg/log"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	apiexts "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
)

type Client struct {
	crds    crdclientset.Interface
	apiexts apiexts.Interface
	// write and read objects to this namespace if not specified on the GlueObjects
	defaultNamespace string
}

func NewClient(cfg *rest.Config, defaultNamespace string) (*Client, error) {
	crdClient, err := crdclientset.NewForConfig(cfg)
	if err != nil {
		return nil, err
	}
	apiextClient, err := apiexts.NewForConfig(cfg)
	if err != nil {
		return nil, err
	}
	return &Client{
		crds:             crdClient,
		apiexts:          apiextClient,
		defaultNamespace: defaultNamespace,
	}, nil
}

func (c *Client) Register() error {
	for _, crd := range crdv1.KnownCRDs {
		toRegister := &v1beta1.CustomResourceDefinition{
			ObjectMeta: metav1.ObjectMeta{Name: crd.FullName()},
			Spec: v1beta1.CustomResourceDefinitionSpec{
				Group:   crd.Group,
				Version: crd.Version,
				Scope:   v1beta1.NamespaceScoped,
				Names: v1beta1.CustomResourceDefinitionNames{
					Plural: crd.Plural,
					Kind:   crd.Kind,
				},
			},
		}
		log.Debugf("registering crd %v", crd)
		if _, err := c.apiexts.ApiextensionsV1beta1().CustomResourceDefinitions().Create(toRegister); err != nil && !apierrors.IsAlreadyExists(err) {
			return fmt.Errorf("failed to create crd: %v", err)
		}
	}
	return nil
}

type crudOperation int

const (
	crudOperationCreate crudOperation = iota
	crudOperationUpdate               = iota
	crudOperationDelete               = iota
)

func (c *Client) Create(item storage.Item) (storage.Item, error) {
	switch glueObject := item.(type) {
	case *v1.Upstream:
		return c.createOrUpdateUpstreamCrd(glueObject, crudOperationCreate)
	case *v1.VirtualHost:
		return c.createOrUpdateVirtualHostCrd(glueObject, crudOperationCreate)
	}
	return nil, errors.Errorf("unsupported object type %v", reflect.TypeOf(item).String())
}

func (c *Client) Update(item storage.Item) (storage.Item, error) {
	switch glueObject := item.(type) {
	case *v1.Upstream:
		return c.createOrUpdateUpstreamCrd(glueObject, crudOperationUpdate)
	case *v1.VirtualHost:
		return c.createOrUpdateVirtualHostCrd(glueObject, crudOperationUpdate)
	}
	return nil, errors.Errorf("unsupported object type %v", reflect.TypeOf(item).String())
}

func (c *Client) createOrUpdateUpstreamCrd(upstream *v1.Upstream, op crudOperation) (*v1.Upstream, error) {
	upstreamCrd, err := UpstreamToCrd(c.defaultNamespace, upstream)
	if err != nil {
		return nil, errors.Wrap(err, "converting glue object to crd")
	}
	upstreams := c.crds.GlueV1().Upstreams(upstreamCrd.Namespace)
	var returnedCrd *crdv1.Upstream
	switch op {
	case crudOperationCreate:
		returnedCrd, err = upstreams.Create(upstreamCrd)
		if err != nil {
			return nil, errors.Wrap(err, "kubernetes create api request")
		}
	case crudOperationUpdate:
		returnedCrd, err = upstreams.Update(upstreamCrd)
		if err != nil {
			return nil, errors.Wrap(err, "kubernetes update api request")
		}
	}
	returnedUpstream, err := UpstreamFromCrd(returnedCrd)
	if err != nil {
		return nil, errors.Wrap(err, "converting created crd to upstream")
	}
	return returnedUpstream, nil
}

func (c *Client) createOrUpdateVirtualHostCrd(virtualHost *v1.VirtualHost, op crudOperation) (*v1.VirtualHost, error) {
	vhostCrd, err := VirtualHostToCrd(c.defaultNamespace, virtualHost)
	if err != nil {
		return nil, errors.Wrap(err, "converting glue object to crd")
	}
	vhosts := c.crds.GlueV1().VirtualHosts(vhostCrd.Namespace)
	var returnedCrd *crdv1.VirtualHost
	switch op {
	case crudOperationCreate:
		returnedCrd, err = vhosts.Create(vhostCrd)
		if err != nil {
			return nil, errors.Wrap(err, "kubernetes create api request")
		}
	case crudOperationUpdate:
		returnedCrd, err = vhosts.Update(vhostCrd)
		if err != nil {
			return nil, errors.Wrap(err, "kubernetes update api request")
		}
	}
	returnedVhost, err := VirtualHostFromCrd(returnedCrd)
	if err != nil {
		return nil, errors.Wrap(err, "converting created crd to virtualHost")
	}
	return returnedVhost, nil
}

//func (c *Client) Delete(item Item) error                                   {}
//func (c *Client) Get(item Item, getOptions *GetOptions) (Item, error)      {}
//func (c *Client) List(item Item, listOptions *ListOptions) ([]Item, error) {}
//func (c *Client) Watch(item Item, watchOptions *WatchOptions, callback func(item Item, operation WatchOperation)) error {
//}
//func (c *Client) WatchStop() {}
