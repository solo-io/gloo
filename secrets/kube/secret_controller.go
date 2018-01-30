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
)

type secretController struct {
	secrets       chan secrets.SecretMap
	errors        chan error
	secretsLister v1.SecretLister
	secretRefs    []string
}

func newSecretController(cfg *rest.Config, resyncDuration time.Duration) (*secretController, error) {
	kubeClient, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create kube clientset: %v", err)
	}

	informerFactory := informers.NewSharedInformerFactory(kubeClient, resyncDuration)
	secretInformer := informerFactory.Core().V1().Secrets()

	kubeController := controller.NewController("glue-secrets-controller", kubeClient,
		secretInformer.Informer())

	ctrl := &secretController{
		secrets:       make(chan secrets.SecretMap),
		errors:        make(chan error),
		secretsLister: secretInformer.Lister(),
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
func (c *secretController) UpdateRefs(secretRefs []string) {
	c.secretRefs = secretRefs
	c.syncSecrets()
}

func (c *secretController) Secrets() <-chan secrets.SecretMap {
	return c.secrets
}

func (c *secretController) Error() <-chan error {
	return c.errors
}

// pushes secretmap or error to channel
func (c *secretController) syncSecrets() {
	secretMap, err := c.getUpdatedSecrets()
	if err != nil {
		c.errors <- err
		return
	}
	// ignore empty configs / no secrets to watch
	if len(secretMap) == 0 {
		return
	}
	c.secrets <- secretMap
}

// retrieves secrets from kubernetes
func (c *secretController) getUpdatedSecrets() (secrets.SecretMap, error) {
	secretList, err := c.secretsLister.List(labels.Everything())
	if err != nil {
		return nil, fmt.Errorf("error retrieving routes: %v", err)
	}
	secretMap := make(secrets.SecretMap)
	for _, secret := range secretList {
		for _, ref := range c.secretRefs {
			if secret.Name == ref {
				log.Printf("updated secret %s", ref)
				secretMap[ref] = secret.Data
				break
			}
		}
	}
	return secretMap, nil
}
