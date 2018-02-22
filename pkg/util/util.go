package util

import (
	"time"

	storage "github.com/solo-io/gloo-storage"
	"github.com/solo-io/gloo-storage/crd"
	"github.com/spf13/cobra"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func GetStorageClient(c *cobra.Command) (storage.Interface, error) {
	var cfg *rest.Config
	var err error
	flags := c.InheritedFlags()
	kubeConfig, _ := flags.GetString("kubeconfig")
	namespace, _ := flags.GetString("namespace")
	if kubeConfig != "" {
		cfg, err = clientcmd.BuildConfigFromFlags("", kubeConfig)
		if err != nil {
			return nil, err
		}
	} else {
		cfg, err = rest.InClusterConfig()
		if err != nil {
			return nil, err
		}
	}
	period, _ := flags.GetInt("sync-period")
	syncPeriod := time.Duration(period) * time.Second

	return crd.NewStorage(cfg, namespace, syncPeriod)
}
