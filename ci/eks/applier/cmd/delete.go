/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("delete called")

		userSpecifiedContext, err := cmd.Flags().GetString("context")
		if err != nil {
			return err
		}

		restConfig, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
			&clientcmd.ClientConfigLoadingRules{ExplicitPath: configFlags.ToRawKubeConfigLoader().ConfigAccess().GetDefaultFilename()},
			&clientcmd.ConfigOverrides{
				CurrentContext: userSpecifiedContext,
			}).ClientConfig()
		if err != nil {
			return err
		}

		clientset, err := kubernetes.NewForConfig(restConfig)
		if err != nil {
			return err
		}
		factory := cmdutil.NewFactory(matchVersionKubeConfigFlags)

		dynamicClient, err := factory.DynamicClient()
		if err != nil {
			return err
		}

		deploymentGVR := schema.GroupVersionResource{
			Group:    "apps",
			Version:  "v1",
			Resource: "deployments",
		}

		if false {
			return dynamicClient.Resource(deploymentGVR).Namespace(*configFlags.Namespace).DeleteCollection(context.Background(), v1.DeleteOptions{}, v1.ListOptions{})
		}

		return clientset.CoreV1().Pods("gloo-system").DeleteCollection(context.Background(), v1.DeleteOptions{}, v1.ListOptions{LabelSelector: "test=test1"})
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deleteCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deleteCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
