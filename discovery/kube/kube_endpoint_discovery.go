package kube

import (
	"fmt"
	"time"

	"k8s.io/client-go/tools/clientcmd"

	"github.com/solo-io/glue/secrets/watcher"
)

func NewSecretWatcher(masterUrl, kubeconfigPath string, resyncDuration time.Duration) (watcher.Watcher, error) {
	cfg, err := clientcmd.BuildConfigFromFlags(masterUrl, kubeconfigPath)
	if err != nil {
		return nil, fmt.Errorf("failed to build rest config: %v", err)
	}

	ctl, err := newDiscoveryController(cfg, resyncDuration)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize secret controller: %v", err)
	}

	return ctl, nil
}
