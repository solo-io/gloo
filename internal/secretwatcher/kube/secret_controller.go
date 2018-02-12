package kube

import (
	"encoding/base64"
	"fmt"
	"time"

	"github.com/mitchellh/hashstructure"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/rest"

	"github.com/solo-io/glue/internal/pkg/kube/controller"
	"github.com/solo-io/glue/pkg/log"
	"github.com/solo-io/glue/pkg/secretwatcher"
)

type secretController struct {
	secrets       chan secretwatcher.SecretMap
	errors        chan error
	secretsLister v1.SecretLister
	secretRefs    []string
	lastSeen      uint64
}

func newSecretController(cfg *rest.Config, resyncDuration time.Duration, stopCh <-chan struct{}) (*secretController, error) {
	kubeClient, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create kube clientset: %v", err)
	}

	informerFactory := informers.NewSharedInformerFactory(kubeClient, resyncDuration)
	secretInformer := informerFactory.Core().V1().Secrets()

	ctrl := &secretController{
		secrets:       make(chan secretwatcher.SecretMap),
		errors:        make(chan error),
		secretsLister: secretInformer.Lister(),
	}

	kubeController := controller.NewController(
		"glue-secrets-controller", kubeClient,
		func(_, _ string, _ interface{}) {
			ctrl.getUpdatedSecrets()
		},
		secretInformer.Informer())

	go informerFactory.Start(stopCh)
	go func() {
		kubeController.Run(2, stopCh)
	}()

	return ctrl, nil
}

// triggers an update
func (c *secretController) TrackSecrets(secretRefs []string) {
	c.secretRefs = secretRefs
	c.syncSecrets()
}

func (c *secretController) Secrets() <-chan secretwatcher.SecretMap {
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
func (c *secretController) getUpdatedSecrets() (secretwatcher.SecretMap, error) {
	secretList, err := c.secretsLister.List(labels.Everything())
	if err != nil {
		return nil, fmt.Errorf("error retrieving secrets: %v", err)
	}
	secrets := make(secretwatcher.SecretMap)
	for _, secret := range secretList {
		for _, ref := range c.secretRefs {
			if secret.Name == ref {
				log.Debugf("updated secret %s", ref)
				secrets[ref] = make(map[string]string)
				for key, encodedValue := range secret.Data {
					var decodedValue []byte
					if _, err := base64.StdEncoding.Decode(decodedValue, encodedValue); err != nil {
						runtime.HandleError(errors.Wrap(err, "failed to decode base64 encoding on secret"))
					}
					secrets[ref][key] = string(encodedValue)
				}
				break
			}
		}
	}
	hash, err := hashstructure.Hash(secrets, nil)
	if err != nil {
		runtime.HandleError(err)
		return nil, nil
	}
	if c.lastSeen == hash {
		return nil, nil
	}
	c.lastSeen = hash
	return secrets, nil
}
