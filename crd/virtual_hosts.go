package crd

import (
	"time"

	"github.com/pkg/errors"
	"github.com/solo-io/gloo-api/pkg/api/types/v1"
	"github.com/solo-io/gloo-storage"
	crdclientset "github.com/solo-io/gloo-storage/crd/client/clientset/versioned"
	crdv1 "github.com/solo-io/gloo-storage/crd/solo.io/v1"
	apiexts "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"

	"github.com/solo-io/gloo-storage/internal/crud"
	"k8s.io/client-go/tools/cache"
)

type virtualHostsClient struct {
	crds    crdclientset.Interface
	apiexts apiexts.Interface
	// write and read objects to this namespace if not specified on the GlooObjects
	namespace     string
	syncFrequency time.Duration
}

func (v *virtualHostsClient) Create(item *v1.VirtualHost) (*v1.VirtualHost, error) {
	return v.createOrUpdateVirtualHostCrd(item, crud.OperationCreate)
}

func (v *virtualHostsClient) Update(item *v1.VirtualHost) (*v1.VirtualHost, error) {
	return v.createOrUpdateVirtualHostCrd(item, crud.OperationUpdate)
}

func (v *virtualHostsClient) Delete(name string) error {
	return v.crds.GlooV1().VirtualHosts(v.namespace).Delete(name, nil)
}

func (v *virtualHostsClient) Get(name string) (*v1.VirtualHost, error) {
	crdVh, err := v.crds.GlooV1().VirtualHosts(v.namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "failed performing get api request")
	}
	returnedVirtualHost, err := VirtualHostFromCrd(crdVh)
	if err != nil {
		return nil, errors.Wrap(err, "converting returned crd to virtualHost")
	}
	return returnedVirtualHost, nil
}

func (v *virtualHostsClient) List() ([]*v1.VirtualHost, error) {
	crdList, err := v.crds.GlooV1().VirtualHosts(v.namespace).List(metav1.ListOptions{})
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

func (v *virtualHostsClient) Watch(handlers ...storage.VirtualHostEventHandler) (*storage.Watcher, error) {
	lw := cache.NewListWatchFromClient(v.crds.GlooV1().RESTClient(), crdv1.VirtualHostCRD.Plural, v.namespace, fields.Everything())
	sw := cache.NewSharedInformer(lw, new(crdv1.VirtualHost), v.syncFrequency)
	for _, h := range handlers {
		sw.AddEventHandler(&virtualHostEventHandler{handler: h, store: sw.GetStore()})
	}
	return storage.NewWatcher(func(stop <-chan struct{}, _ chan error) {
		sw.Run(stop)
	}), nil
}

func (v *virtualHostsClient) createOrUpdateVirtualHostCrd(virtualHost *v1.VirtualHost, op crud.Operation) (*v1.VirtualHost, error) {
	vhostCrd, err := VirtualHostToCrd(v.namespace, virtualHost)
	if err != nil {
		return nil, errors.Wrap(err, "converting gloo object to crd")
	}
	vhosts := v.crds.GlooV1().VirtualHosts(vhostCrd.Namespace)
	var returnedCrd *crdv1.VirtualHost
	switch op {
	case crud.OperationCreate:
		returnedCrd, err = vhosts.Create(vhostCrd)
		if err != nil {
			return nil, errors.Wrap(err, "kubernetes create api request")
		}
	case crud.OperationUpdate:
		// need to make sure we preserve labels and annotations
		currentCrd, err := vhosts.Get(vhostCrd.Name, metav1.GetOptions{ResourceVersion: vhostCrd.ResourceVersion})
		if err != nil {
			return nil, errors.Wrap(err, "kubernetes get api request")
		}
		// copy labels
		vhostCrd.Labels = currentCrd.Labels
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
		vhCrd, ok := updated.(*crdv1.VirtualHost)
		if !ok {
			continue
		}
		updatedVirtualHost, err := VirtualHostFromCrd(vhCrd)
		if err != nil {
			continue
		}
		updatedVirtualHostList = append(updatedVirtualHostList, updatedVirtualHost)
	}
	return updatedVirtualHostList
}

func convertVh(obj interface{}) (*v1.VirtualHost, bool) {
	vhCrd, ok := obj.(*crdv1.VirtualHost)
	if !ok {
		return nil, false
	}
	vh, err := VirtualHostFromCrd(vhCrd)
	if err != nil {
		return nil, false
	}
	return vh, ok
}

func (eh *virtualHostEventHandler) OnAdd(obj interface{}) {
	vh, ok := convertVh(obj)
	if !ok {
		return
	}
	eh.handler.OnAdd(eh.getUpdatedList(), vh)
}
func (eh *virtualHostEventHandler) OnUpdate(_, newObj interface{}) {
	newVh, ok := convertVh(newObj)
	if !ok {
		return
	}
	eh.handler.OnUpdate(eh.getUpdatedList(), newVh)
}

func (eh *virtualHostEventHandler) OnDelete(obj interface{}) {
	vh, ok := convertVh(obj)
	if !ok {
		return
	}
	eh.handler.OnDelete(eh.getUpdatedList(), vh)
}
