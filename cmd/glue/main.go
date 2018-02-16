package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/solo-io/glue/internal/eventloop"
	"github.com/solo-io/glue/pkg/bootstrap"
	"github.com/solo-io/glue/pkg/signals"

	//register plugins
	_ "github.com/solo-io/glue/internal/install"
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var opts bootstrap.Options

var rootCmd = &cobra.Command{
	Use:   "glue-server",
	Short: "Glue Server runs the glue service to turn Envoy into a Function Gateway",
	RunE: func(cmd *cobra.Command, args []string) error {
		stopCh := signals.SetupSignalHandler()
		eventLoop, err := eventloop.Setup(opts, stopCh)
		if err != nil {
			return errors.Wrap(err, "setting up event loop")
		}
		if err := eventLoop.Run(); err != nil {
			return errors.Wrap(err, "failed running event loop")
		}
		return nil
	},
}

func init() {
	// config watcher
	rootCmd.PersistentFlags().StringVar(&opts.ConfigWatcherOptions.Type, "storage.type", bootstrap.WatcherTypeFile, fmt.Sprintf("storage backend for config objects. supported: [%s]", strings.Join(bootstrap.SupportedCwTypes, " | ")))
	rootCmd.PersistentFlags().DurationVar(&opts.ConfigWatcherOptions.SyncFrequency, "storage.refreshrate", time.Second, "refresh rate for polling config")

	// storage watcher
	rootCmd.PersistentFlags().StringVar(&opts.SecretWatcherOptions.Type, "secrets.type", bootstrap.WatcherTypeFile, fmt.Sprintf("storage backend for secrets. supported: [%s]", strings.Join(bootstrap.SupportedSwTypes, " | ")))
	rootCmd.PersistentFlags().DurationVar(&opts.SecretWatcherOptions.SyncFrequency, "secrets.refreshrate", time.Second, "refresh rate for polling secrets")

	// xds port
	rootCmd.PersistentFlags().IntVar(&opts.XdsOptions.Port, "xds.port", 8081, "port to serve envoy xDS services. this port should be specified in your envoy's static config")

	// file
	rootCmd.PersistentFlags().StringVar(&opts.FileOptions.ConfigDir, "file.config.dir", "_glue_config", "root directory to use for storing glue config files")
	rootCmd.PersistentFlags().StringVar(&opts.FileOptions.SecretDir, "file.secret.dir", "_glue_secrets", "root directory to use for storing glue secret files")

	// kube
	rootCmd.PersistentFlags().StringVar(&opts.KubeOptions.MasterURL, "kube.master", "", "url of the kubernetes apiserver. not needed if running in-cluster")
	rootCmd.PersistentFlags().StringVar(&opts.KubeOptions.KubeConfig, "kube.config", "", "path to kubeconfig file. not needed if running in-cluster")

	// vault
	rootCmd.PersistentFlags().StringVar(&opts.VaultOptions.VaultAddr, "vault.addr", "", "url for vault server")
	rootCmd.PersistentFlags().StringVar(&opts.VaultOptions.AuthToken, "vault.token", "", "auth token for reading vault secrets")
	rootCmd.PersistentFlags().IntVar(&opts.VaultOptions.Retries, "vault.retries", 3, "number of times to retry failed requests to vault")
}
