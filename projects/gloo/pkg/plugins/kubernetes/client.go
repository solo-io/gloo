package kubernetes

import (
	"context"
	"sync"
	"time"

	errors "github.com/rotisserie/eris"
	"k8s.io/client-go/tools/cache"

	"github.com/solo-io/solo-kit/pkg/api/v1/clients/kube/controller"
	kubeinformers "k8s.io/client-go/informers"
	corev1listers "k8s.io/client-go/listers/core/v1"
	discoveryv1listers "k8s.io/client-go/listers/discovery/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

//go:generate mockgen -destination ./mocks/kubesharedfactory_mock.go github.com/solo-io/gloo/projects/gloo/pkg/plugins/kubernetes KubePluginSharedFactory

type KubePluginSharedFactory interface {
	EndpointsLister(ns string) corev1listers.EndpointsLister
	EndpointSlicesLister(ns string) discoveryv1listers.EndpointSliceLister
	Subscribe() <-chan struct{}
	Unsubscribe(<-chan struct{})
}

type KubePluginListers struct {
	initError error

	endpointsLister      map[string]corev1listers.EndpointsLister
	endpointSlicesLister map[string]discoveryv1listers.EndpointSliceLister

	cacheUpdatedWatchers      []chan struct{}
	cacheUpdatedWatchersMutex sync.Mutex
}

func getInformerFactory(ctx context.Context, client kubernetes.Interface, watchNamespaces []string) *KubePluginListers {
	if len(watchNamespaces) == 0 {
		watchNamespaces = []string{metav1.NamespaceAll}
	}
	kubePluginSharedFactory := startInformerFactory(ctx, client, watchNamespaces)
	if kubePluginSharedFactory.initError != nil {
		// This is an unrecoverable error (no shared informer factory means all of kube EDS won't work, which is
		// probably the most valuable / important role for gloo) and  users know immediately about e.g. any rbac errors
		// preventing this from working rather than this hiding in the logs
		panic(kubePluginSharedFactory.initError)
	}
	return kubePluginSharedFactory
}

func startInformerFactory(ctx context.Context, client kubernetes.Interface, watchNamespaces []string) *KubePluginListers {
	const resyncDuration = 12 * time.Hour

	var informers []cache.SharedIndexInformer
	k := &KubePluginListers{
		endpointsLister:      map[string]corev1listers.EndpointsLister{},
		endpointSlicesLister: map[string]discoveryv1listers.EndpointSliceLister{},
	}
	for _, nsToWatch := range watchNamespaces {
		kubeInformerFactory := kubeinformers.NewSharedInformerFactoryWithOptions(client, resyncDuration, kubeinformers.WithNamespace(nsToWatch))
		endpointInformer := kubeInformerFactory.Core().V1().Endpoints()
		endpointSliceInformer := kubeInformerFactory.Discovery().V1().EndpointSlices()
		k.endpointsLister[nsToWatch] = endpointInformer.Lister()
		k.endpointSlicesLister[nsToWatch] = endpointSliceInformer.Lister()
		informers = append(informers, endpointInformer.Informer(), endpointSliceInformer.Informer())
	}

	kubeController := controller.NewController(
		"kube-plugin-controller",
		controller.NewLockingSyncHandler(k.updatedOccurred),
		informers...,
	)

	stop := ctx.Done()
	if err := kubeController.Run(2, stop); err != nil && ctx.Err() == nil {
		k.initError = errors.Wrapf(err, "could not start shared informer factory")
		return k
	}

	syncFuncs := make([]cache.InformerSynced, 0, len(informers))
	for _, informer := range informers {
		syncFuncs = append(syncFuncs, informer.HasSynced)
	}

	ok := cache.WaitForCacheSync(stop, syncFuncs...)
	if !ok && ctx.Err() == nil {
		// if initError is non-nil, the kube resource client will panic
		k.initError = errors.Errorf("waiting for kube pod, endpoints, services cache sync failed")
	}

	return k
}

func (k *KubePluginListers) EndpointsLister(ns string) corev1listers.EndpointsLister {
	return k.endpointsLister[ns]
}

func (k *KubePluginListers) EndpointSlicesLister(ns string) discoveryv1listers.EndpointSliceLister {
	return k.endpointSlicesLister[ns]
}

func (k *KubePluginListers) Subscribe() <-chan struct{} {
	k.cacheUpdatedWatchersMutex.Lock()
	defer k.cacheUpdatedWatchersMutex.Unlock()
	c := make(chan struct{}, 1)
	k.cacheUpdatedWatchers = append(k.cacheUpdatedWatchers, c)
	return c
}

func (k *KubePluginListers) Unsubscribe(c <-chan struct{}) {
	k.cacheUpdatedWatchersMutex.Lock()
	defer k.cacheUpdatedWatchersMutex.Unlock()
	for i, cacheUpdated := range k.cacheUpdatedWatchers {
		if cacheUpdated == c {
			k.cacheUpdatedWatchers = append(k.cacheUpdatedWatchers[:i], k.cacheUpdatedWatchers[i+1:]...)
			return
		}
	}
}

func (k *KubePluginListers) updatedOccurred() {
	k.cacheUpdatedWatchersMutex.Lock()
	defer k.cacheUpdatedWatchersMutex.Unlock()
	for _, cacheUpdated := range k.cacheUpdatedWatchers {
		select {
		case cacheUpdated <- struct{}{}:
		default:
		}
	}
}
