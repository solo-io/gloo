package kube

import (
	"fmt"
	"time"

	"github.com/mitchellh/hashstructure"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/rest"

	"github.com/solo-io/gloo/pkg/artifactwatcher"
	"github.com/solo-io/gloo/pkg/log"
	"github.com/solo-io/kubecontroller"
)

type artifactController struct {
	artifacts       chan artifactwatcher.Artifacts
	errors          chan error
	configmapLister v1.ConfigMapLister
	configmapRefs   []string
	lastSeen        uint64
}

func newConfigmapController(cfg *rest.Config, resyncDuration time.Duration, stopCh <-chan struct{}) (*artifactController, error) {
	kubeClient, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create kube clientset: %v", err)
	}

	informerFactory := informers.NewSharedInformerFactory(kubeClient, resyncDuration)
	configmapInformer := informerFactory.Core().V1().ConfigMaps()

	c := &artifactController{
		artifacts:       make(chan artifactwatcher.Artifacts),
		errors:          make(chan error),
		configmapLister: configmapInformer.Lister(),
	}

	kubeController := kubecontroller.NewController(
		"gloo-artifacts-controller", kubeClient,
		kubecontroller.NewSyncHandler(c.syncArtifacts),
		configmapInformer.Informer())

	go informerFactory.Start(stopCh)
	go func() {
		tick := time.Tick(time.Minute)
		go func() {
			for {
				select {
				case <-tick:
					c.syncArtifacts()
				case <-stopCh:
					return
				}
			}
		}()
		kubeController.Run(2, stopCh)
	}()

	return c, nil
}

// triggers an update
func (c *artifactController) TrackArtifacts(artifactRefs []string) {
	c.configmapRefs = artifactRefs
	c.syncArtifacts()
}

func (c *artifactController) Artifacts() <-chan artifactwatcher.Artifacts {
	return c.artifacts
}

func (c *artifactController) Error() <-chan error {
	return c.errors
}

// pushes Artifacts or error to channel
func (c *artifactController) syncArtifacts() {
	Artifacts, err := c.getUpdatedArtifacts()
	if err != nil {
		c.errors <- err
		return
	}
	// ignore empty configs / no artifacts to watch
	if len(Artifacts) == 0 {
		return
	}
	c.artifacts <- Artifacts
}

// retrieves artifacts from kubernetes
func (c *artifactController) getUpdatedArtifacts() (artifactwatcher.Artifacts, error) {
	configmapList, err := c.configmapLister.List(labels.Everything())
	if err != nil {
		return nil, fmt.Errorf("error retrieving artifacts: %v", err)
	}
	artifacts := make(artifactwatcher.Artifacts)
	for _, configmap := range configmapList {
		for _, ref := range c.configmapRefs {
			if configmap.Name == ref {
				log.Debugf("updated artifact %s", ref)
				artifacts[ref] = make(map[string][]byte)
				for key, data := range configmap.Data {
					artifacts[ref][key] = []byte(data)
				}
				for key, data := range configmap.BinaryData {
					artifacts[ref][key] = data
				}
				break
			}
		}
	}
	hash, err := hashstructure.Hash(artifacts, nil)
	if err != nil {
		runtime.HandleError(err)
		return nil, nil
	}
	if c.lastSeen == hash {
		return nil, nil
	}
	c.lastSeen = hash
	return artifacts, nil
}
