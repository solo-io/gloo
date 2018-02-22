package cmd

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pkg/errors"
	"github.com/solo-io/gloo-function-discovery/pkg/secret"
	"github.com/solo-io/gloo-function-discovery/pkg/server"
	"github.com/solo-io/gloo-storage/crd"
	"github.com/spf13/cobra"
	"k8s.io/client-go/rest"
)

func startCmd() *cobra.Command {
	var port int
	cmd := &cobra.Command{
		Use:     "start",
		Aliases: []string{"run"},
		Short:   "Start Gloo Function Discovery service",
		RunE: func(c *cobra.Command, args []string) error {
			cfg, err := getClientConfig()
			if err != nil {
				return errors.Wrap(err, "unable to get client configuration")
			}
			resyncPeriod, _ := c.InheritedFlags().GetInt("sync-period")
			namespace, _ := c.InheritedFlags().GetString("namespace")
			start(cfg, port, time.Duration(resyncPeriod)*time.Second, namespace)
			return nil
		},
	}
	cmd.Flags().IntVarP(&port, "port", "p", 8080, "Port. If not set tries PORT environment variable before defaulting to 8080")
	return cmd
}

func start(cfg *rest.Config, port int, resyncPeriod time.Duration, namespace string) {
	sc, err := crd.NewStorage(cfg, namespace, resyncPeriod)
	if err != nil {
		log.Fatalf("Unable to get client to K8S for monitoring upstreams %q\n", err)
	}
	if namespace == "" {
		namespace = crd.GlooDefaultNamespace
	}
	secretRepo, err := secret.NewSecretRepo(cfg, namespace)
	if err != nil {
		log.Fatalf("Unable to setup monitoring of secrets %q\n", err)
	}
	server := &server.Server{
		Upstreams:  sc.V1().Upstreams(),
		SecretRepo: secretRepo,
		Port:       port,
	}
	log.Println("Listening on ", port)
	stop := make(chan struct{})
	server.Start(resyncPeriod, stop)
	waitSignal(stop)
}

func waitSignal(stop chan struct{}) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
	close(stop)
}
