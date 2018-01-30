package secretwatcher

import (
	"fmt"
	"time"

	"k8s.io/client-go/tools/clientcmd"

	"github.com/solo-io/glue/pkg/module"
)

func NewSecretWatcher(masterUrl, kubeconfigPath string, resyncDuration time.Duration) (module.SecretWatcher, error) {
	cfg, err := clientcmd.BuildConfigFromFlags(masterUrl, kubeconfigPath)
	if err != nil {
		return nil, fmt.Errorf("failed to build rest config: %v", err)
	}

	ctl, err := newSecretController(cfg, resyncDuration)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize secret controller: %v", err)
	}

	return ctl, nil
}
