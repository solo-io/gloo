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

	"strings"

	"github.com/pkg/errors"
	"github.com/solo-io/gloo/pkg/filewatcher"
	"github.com/solo-io/gloo/pkg/log"
	"github.com/solo-io/kubecontroller"
)

/*
 * NOTE: the Structure for FileRefs in kubernetes follows the following format:
 *               "configmap_name/key_name"
 * Config maps may contain more than one file, but a ref should refer to only a
 * single file
 */

type artifactController struct {
	files           chan filewatcher.Files
	errors          chan error
	configmapLister v1.ConfigMapLister
	fileRefs        []string
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
		files:           make(chan filewatcher.Files),
		errors:          make(chan error),
		configmapLister: configmapInformer.Lister(),
	}

	kubeController := kubecontroller.NewController(
		"gloo-files-controller", kubeClient,
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
func (c *artifactController) TrackFiles(fileRefs []string) {
	c.fileRefs = fileRefs
	c.syncArtifacts()
}

func (c *artifactController) Files() <-chan filewatcher.Files {
	return c.files
}

func (c *artifactController) Error() <-chan error {
	return c.errors
}

// pushes Files or error to channel
func (c *artifactController) syncArtifacts() {
	Artifacts, err := c.getUpdatedArtifacts()
	if err != nil {
		c.errors <- err
		return
	}
	// ignore empty configs / no files to watch
	if len(Artifacts) == 0 {
		return
	}
	c.files <- Artifacts
}

// retrieves files from kubernetes
func (c *artifactController) getUpdatedArtifacts() (filewatcher.Files, error) {
	configmapList, err := c.configmapLister.List(labels.Everything())
	if err != nil {
		return nil, fmt.Errorf("error retrieving files: %v", err)
	}
	artifacts := make(filewatcher.Files)
	for _, ref := range c.fileRefs {
		configmapName, configmapKey, err := getConfigmapRef(ref)
		if err != nil {
			log.Warnf("ignoring file ref: %v", err)
		}
	find:
		for _, configmap := range configmapList {
			if configmap.Name == configmapName {
				log.Debugf("updated artifact %s", ref)
				for key, data := range configmap.Data {
					if key == configmapKey {
						artifacts[ref] = []byte(data)
						break find
					}
				}
				for key, data := range configmap.Data {
					if key == configmapKey {
						artifacts[ref] = []byte(data)
						break find
					}
				}
			}
			log.Warnf("config map %v or key %v not found", configmapName, configmapKey)
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

func getConfigmapRef(fileRef string) (string, string, error) {
	parts := strings.Split(fileRef, "/")
	if len(parts) != 2 {
		return "", "", errors.Errorf("invalid file ref for kubernetes: %v. file refs for "+
			"kubernetes must follow the format <configmap_name>/<key_name>", fileRef)
	}
	return parts[0], parts[1], nil
}
