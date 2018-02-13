package crd

import (
	"fmt"

	gluev1 "github.com/solo-io/glue/pkg/api/types/v1"
	crdclient "github.com/solo-io/glue/pkg/platform/kube/crd/client/clientset/versioned"
	crd "github.com/solo-io/glue/pkg/platform/kube/crd/solo.io/v1"
	"github.com/solo-io/gluectl/pkg/storage"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	apiexts "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/rest"
)

var evtTypeToOperation = map[watch.EventType]storage.WatchOperation{
	watch.Added:    storage.Create,
	watch.Modified: storage.Update,
	watch.Deleted:  storage.Delete,
	watch.Error:    storage.Error,
}

type CrdStorage struct {
	clientset       *crdclient.Clientset
	apiextset       *apiexts.Clientset
	namespace       string
	upstreamWatcher watch.Interface
	vhostWatcher    watch.Interface
}

func NewCrdStorage(cfg *rest.Config, namespace string) (*CrdStorage, error) {
	cl, err := crdclient.NewForConfig(cfg)
	if err != nil {
		return nil, err
	}
	client, err := apiexts.NewForConfig(cfg)
	if err != nil {
		return nil, err
	}
	return &CrdStorage{
		clientset: cl,
		namespace: namespace,
		apiextset: client,
	}, nil
}

func (c *CrdStorage) Register(item storage.Item) error {

	switch item.(type) {

	case *gluev1.Upstream:
		upstream := &v1beta1.CustomResourceDefinition{
			ObjectMeta: metav1.ObjectMeta{Name: "upstreams.glue.solo.io"},
			Spec: v1beta1.CustomResourceDefinitionSpec{
				Group:   crd.GroupName,
				Version: crd.Version,
				Scope:   v1beta1.NamespaceScoped,
				Names: v1beta1.CustomResourceDefinitionNames{
					Plural: "upstreams",
					Kind:   "Upstream",
				},
			},
		}
		if _, err := c.apiextset.ApiextensionsV1beta1().CustomResourceDefinitions().Create(upstream); err != nil && !apierrors.IsAlreadyExists(err) {
			return err
		}
	case *gluev1.VirtualHost:
		vhost := &v1beta1.CustomResourceDefinition{
			ObjectMeta: metav1.ObjectMeta{Name: "virtualhosts.glue.solo.io"},
			Spec: v1beta1.CustomResourceDefinitionSpec{
				Group:   crd.GroupName,
				Version: crd.Version,
				Scope:   v1beta1.NamespaceScoped,
				Names: v1beta1.CustomResourceDefinitionNames{
					Plural: "virtualhosts",
					Kind:   "VirtualHost",
				},
			},
		}
		if _, err := c.apiextset.ApiextensionsV1beta1().CustomResourceDefinitions().Create(vhost); err != nil && !apierrors.IsAlreadyExists(err) {
			return err
		}
	default:
	}
	return nil
}

func (c *CrdStorage) Create(item storage.Item) (storage.Item, error) {
	if obj, ok := item.(*gluev1.Upstream); ok {
		cobj := crd.UpstreamToCRD(metav1.ObjectMeta{Name: obj.Name}, *obj)
		return upstream(c.clientset.GlueV1().Upstreams(c.namespace).Create(cobj))
	} else if obj, ok := item.(*gluev1.VirtualHost); ok {
		cobj := crd.VirtualHostToCRD(metav1.ObjectMeta{Name: obj.Name}, *obj)
		return vhost(c.clientset.GlueV1().VirtualHosts(c.namespace).Create(cobj))
	}
	return nil, fmt.Errorf("Unknown Item Type: %t", item)
}

func (c *CrdStorage) Update(item storage.Item) (storage.Item, error) {
	if obj, ok := item.(*gluev1.Upstream); ok {
		xobj, err := c.clientset.GlueV1().Upstreams(c.namespace).Get(obj.Name, metav1.GetOptions{})
		if err != nil {
			return nil, err
		}
		cobj := crd.UpstreamToCRD(xobj.ObjectMeta, *obj)
		return upstream(c.clientset.GlueV1().Upstreams(c.namespace).Update(cobj))
	} else if obj, ok := item.(*gluev1.VirtualHost); ok {
		xobj, err := c.clientset.GlueV1().VirtualHosts(c.namespace).Get(obj.Name, metav1.GetOptions{})
		if err != nil {
			return nil, err
		}
		cobj := crd.VirtualHostToCRD(xobj.ObjectMeta, *obj)
		return vhost(c.clientset.GlueV1().VirtualHosts(c.namespace).Update(cobj))
	}
	return nil, fmt.Errorf("Unknown Item Type: %t", item)
}

func (c *CrdStorage) Delete(item storage.Item) error {
	if obj, ok := item.(*gluev1.Upstream); ok {
		return c.clientset.GlueV1().Upstreams(c.namespace).Delete(obj.Name, &metav1.DeleteOptions{})
	} else if obj, ok := item.(*gluev1.VirtualHost); ok {
		return c.clientset.GlueV1().VirtualHosts(c.namespace).Delete(obj.Name, &metav1.DeleteOptions{})
	}
	return fmt.Errorf("Unknown Item Type: %t", item)
}

func (c *CrdStorage) Get(item storage.Item, getOptions *storage.GetOptions) (storage.Item, error) {
	if obj, ok := item.(*gluev1.Upstream); ok {
		return upstream(c.clientset.GlueV1().Upstreams(c.namespace).Get(obj.Name, metav1.GetOptions{}))
	} else if obj, ok := item.(*gluev1.VirtualHost); ok {
		return vhost(c.clientset.GlueV1().VirtualHosts(c.namespace).Get(obj.Name, metav1.GetOptions{}))
	}
	return nil, fmt.Errorf("Unknown Item Type: %t", item)
}

func (c *CrdStorage) List(item storage.Item, listOptions *storage.ListOptions) ([]storage.Item, error) {
	switch item.(type) {
	case *gluev1.Upstream:
		list, err := c.clientset.GlueV1().Upstreams(c.namespace).List(metav1.ListOptions{})
		if err != nil {
			return nil, err
		}
		l := make([]storage.Item, 0, len(list.Items))
		for _, i := range list.Items {
			x := crd.UpstreamFromCRD(&i)
			l = append(l, &x)
		}
		return l, nil
	case *gluev1.VirtualHost:
		list, err := c.clientset.GlueV1().VirtualHosts(c.namespace).List(metav1.ListOptions{})
		if err != nil {
			return nil, err
		}
		l := make([]storage.Item, 0, len(list.Items))
		for _, i := range list.Items {
			x := crd.VirtualHostFromCRD(&i)
			l = append(l, &x)
		}
		return l, nil
	default:
		return nil, fmt.Errorf("Unknown Item Type: %t", item)
	}
}

func (c *CrdStorage) Watch(item storage.Item, watchOptions *storage.WatchOptions, callback func(item storage.Item, operation storage.WatchOperation)) error {

	var err error
	switch item.(type) {
	case *gluev1.Upstream:
		c.upstreamWatcher, err = c.clientset.GlueV1().Upstreams(c.namespace).Watch(metav1.ListOptions{})
		if err != nil {
			return err
		}
		go func() {
			for {
				evt, ok := <-c.upstreamWatcher.ResultChan()
				if !ok {
					return
				}
				obj, ok := evt.Object.(*crd.Upstream)
				if ok {
					callback(crd.UpstreamFromCRD(obj), evtTypeToOperation[evt.Type])
				}
			}
		}()
		return nil

	case *gluev1.VirtualHost:
		c.vhostWatcher, err = c.clientset.GlueV1().VirtualHosts(c.namespace).Watch(metav1.ListOptions{})
		if err != nil {
			return err
		}
		go func() {
			for {
				evt, ok := <-c.vhostWatcher.ResultChan()
				if !ok {
					return
				}
				obj, ok := evt.Object.(*crd.VirtualHost)
				if ok {
					callback(crd.VirtualHostFromCRD(obj), evtTypeToOperation[evt.Type])
				}
			}
		}()
		return nil

	default:
		return fmt.Errorf("Unknown Item Type: %t", item)
	}
}

func (c *CrdStorage) WatchStop() {
	if c.upstreamWatcher != nil {
		c.upstreamWatcher.Stop()
	}

	if c.vhostWatcher != nil {
		c.vhostWatcher.Stop()
	}
}

func upstream(u *crd.Upstream, err error) (*gluev1.Upstream, error) {
	if err != nil {
		return nil, err
	}
	us := crd.UpstreamFromCRD(u)
	return &us, nil
}

func vhost(u *crd.VirtualHost, err error) (*gluev1.VirtualHost, error) {
	if err != nil {
		return nil, err
	}
	us := crd.VirtualHostFromCRD(u)
	return &us, nil
}
