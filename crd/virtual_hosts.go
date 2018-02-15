package crd

import (
	"time"

	"github.com/pkg/errors"
	"github.com/solo-io/glue-storage"
	crdclientset "github.com/solo-io/glue-storage/crd/client/clientset/versioned"
	crdv1 "github.com/solo-io/glue-storage/crd/solo.io/v1"
	"github.com/solo-io/glue/pkg/api/types/v1"
	apiexts "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"

	"github.com/solo-io/glue-storage/internal/crud"
	"k8s.io/client-go/tools/cache"
)

type virtualHostsClient struct {
	crds    crdclientset.Interface
	apiexts apiexts.Interface
	// write and read objects to this namespace if not specified on the GlueObjects
	namespace     string
	syncFrequency time.Duration
}

func (c *virtualHostsClient) Create(item *v1.VirtualHost) (*v1.VirtualHost, error) {
	return c.createOrUpdateVirtualHostCrd(item, crud.OperationCreate)
}

func (c *virtualHostsClient) Update(item *v1.VirtualHost) (*v1.VirtualHost, error) {
	return c.createOrUpdateVirtualHostCrd(item, crud.OperationUpdate)
}

func (c *virtualHostsClient) Delete(name string) error {
	return c.crds.GlueV1().VirtualHosts(c.namespace).Delete(name, nil)
}

func (c *virtualHostsClient) Get(name string) (*v1.VirtualHost, error) {
	crdVh, err := c.crds.GlueV1().VirtualHosts(c.namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "failed performing get api request")
	}
	returnedVirtualHost, err := VirtualHostFromCrd(crdVh)
	if err != nil {
		return nil, errors.Wrap(err, "converting returned crd to virtualHost")
	}
	return returnedVirtualHost, nil
}

func (c *virtualHostsClient) List() ([]*v1.VirtualHost, error) {
	crdList, err := c.crds.GlueV1().VirtualHosts(c.namespace).List(metav1.ListOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "failed performing list api request")
	}
	var returnedVirtualHosts []*v1.VirtualHost
	for _, crdVh := range crdList.Items {
		virtualHost, err := VirtualHostFromCrd(&crdVh)
		if err != nil {
			return nil, errors.Wrap(err, "converting returned crd to virtualHost")
		}
		returnedVirtualHosts = append(returnedVirtualHosts, virtualHost)
	}
	return returnedVirtualHosts, nil
}

func (u *virtualHostsClient) Watch(handlers ...storage.VirtualHostEventHandler) (*storage.Watcher, error) {
	lw := cache.NewListWatchFromClient(u.crds.GlueV1().RESTClient(), crdv1.VirtualHostCRD.Plural, metav1.NamespaceAll, fields.Everything())
	sw := cache.NewSharedInformer(lw, new(crdv1.VirtualHost), u.syncFrequency)
	for _, h := range handlers {
		sw.AddEventHandler(&virtualHostEventHandler{handler: h, store: sw.GetStore()})
	}
	return storage.NewWatcher(func(stop <-chan struct{}) {
		sw.Run(stop)
	}), nil
}

func (c *virtualHostsClient) createOrUpdateVirtualHostCrd(virtualHost *v1.VirtualHost, op crud.Operation) (*v1.VirtualHost, error) {
	vhostCrd, err := VirtualHostToCrd(c.namespace, virtualHost)
	if err != nil {
		return nil, errors.Wrap(err, "converting glue object to crd")
	}
	vhosts := c.crds.GlueV1().VirtualHosts(vhostCrd.Namespace)
	var returnedCrd *crdv1.VirtualHost
	switch op {
	case crud.OperationCreate:
		returnedCrd, err = vhosts.Create(vhostCrd)
		if err != nil {
			return nil, errors.Wrap(err, "kubernetes create api request")
		}
	case crud.OperationUpdate:
		returnedCrd, err = vhosts.Update(vhostCrd)
		if err != nil {
			return nil, errors.Wrap(err, "kubernetes update api request")
		}
	}
	returnedVhost, err := VirtualHostFromCrd(returnedCrd)
	if err != nil {
		return nil, errors.Wrap(err, "converting returned crd to virtualHost")
	}
	return returnedVhost, nil
}

// implements the kubernetes ResourceEventHandler interface
type virtualHostEventHandler struct {
	handler storage.VirtualHostEventHandler
	store   cache.Store
}

func (eh *virtualHostEventHandler) getUpdatedList() []*v1.VirtualHost {
	updatedList := eh.store.List()
	var updatedVirtualHostList []*v1.VirtualHost
	for _, updated := range updatedList {
		usCrd, ok := updated.(*crdv1.VirtualHost)
		if !ok {
			continue
		}
		updatedVirtualHost, err := VirtualHostFromCrd(usCrd)
		if err != nil {
			continue
		}
		updatedVirtualHostList = append(updatedVirtualHostList, updatedVirtualHost)
	}
	return updatedVirtualHostList
}

func convertVh(obj interface{}) (*v1.VirtualHost, bool) {
	usCrd, ok := obj.(*crdv1.VirtualHost)
	if !ok {
		return nil, ok
	}
	us, err := VirtualHostFromCrd(usCrd)
	if err != nil {
		return nil, false
	}
	return us, ok
}

func (eh *virtualHostEventHandler) OnAdd(obj interface{}) {
	us, ok := convertVh(obj)
	if !ok {
		return
	}
	eh.handler.OnAdd(eh.getUpdatedList(), us)
}
func (eh *virtualHostEventHandler) OnUpdate(_, newObj interface{}) {
	newVh, ok := convertVh(newObj)
	if !ok {
		return
	}
	eh.handler.OnUpdate(eh.getUpdatedList(), newVh)
}

func (eh *virtualHostEventHandler) OnDelete(obj interface{}) {
	us, ok := convertVh(obj)
	if !ok {
		return
	}
	eh.handler.OnDelete(eh.getUpdatedList(), us)
}
