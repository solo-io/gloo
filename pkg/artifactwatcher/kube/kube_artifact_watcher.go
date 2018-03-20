package kube

import (
	"fmt"
	"time"

	"k8s.io/client-go/tools/clientcmd"

	"github.com/solo-io/gloo/pkg/artifactwatcher"
)

func NewArtifactWatcher(masterUrl, kubeconfigPath string, resyncDuration time.Duration, stopCh <-chan struct{}) (artifactwatcher.Interface, error) {
	cfg, err := clientcmd.BuildConfigFromFlags(masterUrl, kubeconfigPath)
	if err != nil {
		return nil, fmt.Errorf("failed to build rest config: %v", err)
	}

	ctl, err := newConfigmapController(cfg, resyncDuration, stopCh)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize artifact controller: %v", err)
	}

	return ctl, nil
}
