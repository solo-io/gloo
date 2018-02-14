package crd

import (
	"fmt"

	crdclientset "github.com/solo-io/glue-storage/crd/client/clientset/versioned"
	crdv1 "github.com/solo-io/glue-storage/crd/solo.io/v1"
	"istio.io/istio/pkg/log"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	apiexts "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
)

type Client struct {
	crds      crdclientset.Interface
	apiexts   apiexts.Interface
	namespace string
}

func NewClient(cfg *rest.Config, namespace string) (*Client, error) {
	crdClient, err := crdclientset.NewForConfig(cfg)
	if err != nil {
		return nil, err
	}
	apiextClient, err := apiexts.NewForConfig(cfg)
	if err != nil {
		return nil, err
	}
	return &Client{
		crds:      crdClient,
		apiexts:   apiextClient,
		namespace: namespace,
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

//func (c *Client) Create(item Item) (Item, error)                           {}
//func (c *Client) Update(item Item) (Item, error)                           {}
//func (c *Client) Delete(item Item) error                                   {}
//func (c *Client) Get(item Item, getOptions *GetOptions) (Item, error)      {}
//func (c *Client) List(item Item, listOptions *ListOptions) ([]Item, error) {}
//func (c *Client) Watch(item Item, watchOptions *WatchOptions, callback func(item Item, operation WatchOperation)) error {
//}
//func (c *Client) WatchStop() {}
