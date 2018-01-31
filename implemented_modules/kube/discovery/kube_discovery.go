package secretwatcher

import (
	"fmt"
	"time"

	"k8s.io/client-go/tools/clientcmd"

	"github.com/solo-io/glue/pkg/module"
)

func NewEndpointDiscovery(masterUrl, kubeconfigPath string, resyncDuration time.Duration, stopCh <-chan struct{}) (module.Discovery, error) {
	cfg, err := clientcmd.BuildConfigFromFlags(masterUrl, kubeconfigPath)
	if err != nil {
		return nil, fmt.Errorf("failed to build rest config: %v", err)
	}

	ctl, err := newEndpointController(cfg, resyncDuration, stopCh)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize endpoint controller: %v", err)
	}

	return ctl, nil
}
