package main

import (
	"context"
	"flag"
	"github.com/solo-io/solo-projects/projects/vcs/pkg"
	"github.com/solo-io/solo-projects/projects/vcs/pkg/git"
	"k8s.io/client-go/rest"
	"log"

	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/factory"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/kube"
	"github.com/solo-io/solo-kit/pkg/utils/contextutils"
	"github.com/solo-io/solo-kit/pkg/utils/kubeutils"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/defaults"
	"github.com/solo-io/solo-projects/projects/vcs/pkg/api/v1"
)

// TODO(marco): This is just a very simple and temporary way of bootstrapping the main loop
func main() {
	if err := run(); err != nil {
		log.Fatalf("%v", err)
	}
}

// Start the event loop watching the changeset resources
func run() error {

	ctx := contextutils.WithLogger(context.Background(), pkg.AppName)

	// TODO: remove
	configPath := flag.String("kube-config", "", "Path to the KUBECONFIG file. Leave blank for in-cluster configuration")
	flag.Parse()

	// Retrieve kubernetes configuration
	cfg, err := kubeutils.GetConfig("", *configPath)
	if err != nil {
		return err
	}

	// Configure changeset resource client
	csClient, err := buildChangeSetClient(cfg)
	if err != nil {
		return err
	}

	// Configure changeset snapshot emitter (producer)
	emitter := v1.NewApiEmitter(csClient)

	// Start event loop
	errs, err := v1.NewApiEventLoop(emitter, &git.RemoteSyncer{}).Run(
		[]string{defaults.GlooSystem},
		clients.WatchOpts{Ctx: ctx})

	if err != nil {
		return err
	}

	// Receive from error channel until it closes
	for err := range errs {
		contextutils.LoggerFrom(ctx).Errorf("error in setup: %v", err)
	}

	return nil
}

// Creates and registers a client for the changeset resource
func buildChangeSetClient(config *rest.Config) (v1.ChangeSetClient, error) {
	csClient, err := v1.NewChangeSetClient(&factory.KubeResourceClientFactory{
		Crd:         v1.ChangeSetCrd,
		Cfg:         config,
		SharedCache: kube.NewKubeCache(),
	})
	if err != nil {
		return nil, err
	}

	err = csClient.Register()
	if err != nil {
		return nil, err
	}

	return csClient, nil
}
