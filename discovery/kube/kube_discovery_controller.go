package kube

import (
	"fmt"
	"time"

	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/sample-controller/pkg/signals"

	"github.com/solo-io/glue/adapters/kube/controller"
	"github.com/solo-io/glue/pkg/log"
	"github.com/solo-io/glue/secrets"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/listers/core/v1"
	"github.com/solo-io/glue/discovery"
)

type discoveryController struct {
	clusters        chan discovery.Clusters
	errors          chan error
	serviceLister   v1.ServiceLister
	endpointsLister v1.EndpointsLister
	serviceRefs    []serviceRef
}

type serviceRef struct {
	namespace, name string
}

func newDiscoveryController(cfg *rest.Config, resyncDuration time.Duration) (*discoveryController, error) {
	kubeClient, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create kube clientset: %v", err)
	}

	informerFactory := informers.NewSharedInformerFactory(kubeClient, resyncDuration)
	serviceInformer := informerFactory.Core().V1().Services()
	endpointsLister := informerFactory.Core().V1().Endpoints()

	kubeController := controller.NewController("glue-secrets-controller", kubeClient,
		serviceInformer.Informer(), endpointsLister.Informer())

	ctrl := &discoveryController{
		clusters:      make(chan discovery.Clusters),
		errors:        make(chan error),
		serviceLister: serviceInformer.Lister(),
		endpointsLister: endpointsLister.Lister(),
	}

	kubeController.AddEventHandler(controller.Added, func(_, _ string, _ interface{}) {
		ctrl.getUpdatedSecrets()
	})
	kubeController.AddEventHandler(controller.Updated, func(namespace, name string, _ interface{}) {
		ctrl.getUpdatedSecrets()
	})
	kubeController.AddEventHandler(controller.Deleted, func(namespace, name string, _ interface{}) {
		ctrl.getUpdatedSecrets()
	})

	// set up signals so we handle the first shutdown signal gracefully
	stopCh := signals.SetupSignalHandler()

	go informerFactory.Start(stopCh)
	go func() {
		kubeController.Run(2, stopCh)
	}()

	return ctrl, nil
}

// triggers an update
func (c *discoveryController) UpdateClusterRefs(clusterRefs []string) {
	var serviceRefs []serviceRef
	for _, cluster := range clusterRefs {

	}
	c.syncSecrets()
}

func (c *discoveryController) Ckusters() <-chan secrets.SecretMap {
	return c.clusters
}

func (c *discoveryController) Error() <-chan error {
	return c.errors
}

// pushes secretmap or error to channel
func (c *discoveryController) syncSecrets() {
	secretMap, err := c.getUpdatedSecrets()
	if err != nil {
		c.errors <- err
		return
	}
	// ignore empty configs / no secrets to watch
	if len(secretMap) == 0 {
		return
	}
	c.clusters <- secretMap
}

// retrieves secrets from kubernetes
func (c *discoveryController) getUpdatedSecrets() (secrets.SecretMap, error) {
	secretList, err := c.serviceLister.List(labels.Everything())
	if err != nil {
		return nil, fmt.Errorf("error retrieving routes: %v", err)
	}
	secretMap := make(secrets.SecretMap)
	for _, secret := range secretList {
		for _, ref := range c.serviceNames {
			if secret.Name == ref {
				log.Printf("updated secret %s", ref)
				secretMap[ref] = secret.Data
				break
			}
		}
	}
	return secretMap, nil
}
