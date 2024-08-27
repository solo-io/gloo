package serviceentry

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/kube/controller"
	istioclient "istio.io/client-go/pkg/clientset/versioned"
	"k8s.io/client-go/kubernetes"

	networkingclient "istio.io/client-go/pkg/apis/networking/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
)

// copied from kube plugin
const resyncPeriod = 12 * time.Hour

type informerLister[T any] struct {
	Informer cache.SharedIndexInformer
}

func watchInformers(updates chan struct{}, stop <-chan struct{}, informers ...cache.SharedIndexInformer) error {
	kubeController := controller.NewController("serviceentry-plugin-controller",
		controller.NewLockingSyncHandler(func() {
			updates <- struct{}{}
		}),
		informers...)

	err := kubeController.Run(2, stop)
	if err != nil {
		return err
	}

	var syncFuncs []cache.InformerSynced
	for _, informer := range informers {
		syncFuncs = append(syncFuncs, informer.HasSynced)
	}

	ok := cache.WaitForCacheSync(stop, syncFuncs...)
	if !ok {
		return errors.Errorf("waiting for kube pod, endpoints, services cache sync failed")
	}

	return nil
}

func (n *informerLister[T]) List(namespace string, selector labels.Selector) ([]T, error) {
	var res []T
	err := cache.ListAllByNamespace(n.Informer.GetIndexer(), namespace, selector, func(i any) {
		cast := i.(T)
		res = append(res, cast)
	})
	return res, err
}

func serviceEntryInformer(client istioclient.Interface, ns string) informerLister[*networkingclient.ServiceEntry] {
	l := func(options metav1.ListOptions) (runtime.Object, error) {
		return client.NetworkingV1beta1().ServiceEntries(ns).List(context.Background(), options)
	}
	w := func(options metav1.ListOptions) (watch.Interface, error) {
		return client.NetworkingV1beta1().ServiceEntries(ns).Watch(context.Background(), options)
	}
	return informerLister[*networkingclient.ServiceEntry]{
		Informer: cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc:  l,
				WatchFunc: w,
			},
			&networkingclient.ServiceEntry{},
			resyncPeriod,
			cache.Indexers{
				cache.NamespaceIndex: cache.MetaNamespaceIndexFunc,
			},
		),
	}
}

func podInformer(client kubernetes.Interface, ns string) informerLister[*corev1.Pod] {
	l := func(options metav1.ListOptions) (runtime.Object, error) {
		return client.CoreV1().Pods(ns).List(context.Background(), options)
	}
	w := func(options metav1.ListOptions) (watch.Interface, error) {
		return client.CoreV1().Pods(ns).Watch(context.Background(), options)
	}
	return informerLister[*corev1.Pod]{
		Informer: cache.NewSharedIndexInformer(
			&cache.ListWatch{
				ListFunc:  l,
				WatchFunc: w,
			},
			&corev1.Pod{},
			resyncPeriod,
			cache.Indexers{
				cache.NamespaceIndex: cache.MetaNamespaceIndexFunc,
			},
		),
	}
}
