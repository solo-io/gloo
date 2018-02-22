package cmd

import (
	"github.com/spf13/cobra"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var kubeConfig string

func RootCmd() *cobra.Command {
	var namespace string
	var syncPeriod int

	root := &cobra.Command{
		Use:   "gloo-function-discovery",
		Short: "Gloo Function Discovery service",
	}
	pflags := root.PersistentFlags()
	pflags.StringVar(&kubeConfig, "kubeconfig", "", "Path to K8S config. Needed for out-of-cluster")
	pflags.StringVarP(&namespace, "namespace", "n", "default", "K8S namespace to use")
	pflags.IntVar(&syncPeriod, "sync-period", 300, "sync period (seconds) for resources")
	root.MarkFlagFilename("kubeconfig")
	root.AddCommand(registerCmd())
	root.AddCommand(startCmd())
	return root
}

// GetClientConfig returns the REST configuration necessary to connect to k8s
func getClientConfig() (*rest.Config, error) {
	if kubeConfig != "" {
		return clientcmd.BuildConfigFromFlags("", kubeConfig)
	}
	return rest.InClusterConfig()
}
