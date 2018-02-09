package cmd

import (
	"github.com/spf13/cobra"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	kubeConfig string
	namespace  string
)

func RootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:   "glue-discovery",
		Short: "Glue Function Discovery service",
	}
	root.PersistentFlags().StringVar(&kubeConfig, "kubeconf", "", "Path to K8S config. Needed for out-of-cluster")
	root.PersistentFlags().StringVarP(&namespace, "namespace", "n", "default", "K8S namespace to use")
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
