package crd

import (
	"fmt"
	"time"

	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/sample-controller/pkg/signals"

	clientset "github.com/solo-io/glue/config/watcher/crd/client/clientset/versioned"
	informers "github.com/solo-io/glue/config/watcher/crd/client/informers/externalversions"
	"github.com/solo-io/glue/config/watcher/crd/controller"
)

func NewCrdWatcher(masterUrl, kubeconfigPath string, resyncDuration time.Duration) (*controller.Controller, error) {
	cfg, err := clientcmd.BuildConfigFromFlags(masterUrl, kubeconfigPath)
	if err != nil {
		return nil, fmt.Errorf("failed to build rest config: %v", err)
	}
	err = controller.RegisterCrds(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to register crds: %v", err)
	}
	kubeClient, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create kube clientset: %v", err)
	}

	glueClient, err := clientset.NewForConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create glue clientset: %v", err)
	}

	glueInformerFactory := informers.NewSharedInformerFactory(glueClient, resyncDuration)

	ctl := controller.NewController(kubeClient, glueInformerFactory)

	// set up signals so we handle the first shutdown signal gracefully
	stopCh := signals.SetupSignalHandler()

	go glueInformerFactory.Start(stopCh)

	go func() {
		err = ctl.Run(2, stopCh)
		runtime.HandleError(err)
	}()

	return ctl, nil
}
