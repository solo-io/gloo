package cmd

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/pkg/errors"
	"github.com/solo-io/glue-discovery/pkg/secret"
	"github.com/solo-io/glue-discovery/pkg/server"
	"github.com/spf13/cobra"
	"k8s.io/client-go/rest"
)

func startCmd() *cobra.Command {
	var port int
	cmd := &cobra.Command{
		Use:   "start",
		Short: "Start Glue Function Discovery service",
		RunE: func(c *cobra.Command, args []string) error {
			cfg, err := getClientConfig()
			if err != nil {
				return errors.Wrap(err, "unable to get client configuration")
			}
			start(cfg, port, namespace)
			return nil
		},
	}
	cmd.LocalFlags().IntVarP(&port, "port", "p", 8080, "Port. If not set tries PORT environment variable before defaulting to 8080")
	return cmd
}

func start(cfg *rest.Config, port int, namespace string) {
	upstreamInterface, err := server.UpstreamInterface(cfg, namespace)
	if err != nil {
		log.Fatalf("Unable to get client to K8S for monitoring upstreams %q\n", err)
	}

	secretRepo, err := secret.NewSecretRepo(cfg)
	if err != nil {
		log.Fatalf("Unable to setup monitoring of secrets %q\n", err)
	}
	server := &server.Server{
		UpstreamRepo: upstreamInterface,
		SecretRepo:   secretRepo,
		Port:         port,
	}
	log.Println("Listening on ", port)
	stop := make(chan struct{})
	server.Start(stop)
	waitSignal(stop)
}

func waitSignal(stop chan struct{}) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
	close(stop)
}
