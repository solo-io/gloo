package main

import (
	"context"
	"log"
	"os"

	"github.com/solo-io/solo-projects/projects/vcs/pkg/constants"

	"github.com/solo-io/solo-projects/projects/vcs/pkg/git"
	"k8s.io/client-go/rest"

	"github.com/solo-io/gloo/projects/gloo/pkg/defaults"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/factory"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/kube"
	"github.com/solo-io/solo-kit/pkg/utils/contextutils"
	"github.com/solo-io/solo-kit/pkg/utils/kubeutils"
	v1 "github.com/solo-io/solo-projects/projects/vcs/pkg/api/v1"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("%v", err)
	}
}

// Start the event loop watching the changeset resources
func run() error {

	ctx := validateEnv()

	// Retrieve kubernetes configuration
	cfg, err := kubeutils.GetConfig("", "")
	if err != nil {
		return err
	}

	// Configure changeset resource client for emitter
	emitterCsClient, syncerCsClient, err := buildChangeSetClients(cfg)
	if err != nil {
		return err
	}

	// Configure changeset snapshot emitter (producer)
	emitter := v1.NewApiEmitter(emitterCsClient)

	// Start event loop
	errs, err := v1.NewApiEventLoop(emitter, &git.RemoteSyncer{CsClient: &syncerCsClient}).Run(
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

// Verifies that the required env variables have been set
func validateEnv() context.Context {
	ctx := contextutils.WithLogger(context.Background(), constants.AppName)

	if os.Getenv(constants.AuthTokenEnvVariableName) == "" {
		contextutils.LoggerFrom(ctx).Panicf("Environment variable %v is not set", constants.AuthTokenEnvVariableName)
	}
	if os.Getenv(constants.RemoteUriEnvVariableName) == "" {
		contextutils.LoggerFrom(ctx).Panicf("Environment variable %v is not set", constants.RemoteUriEnvVariableName)
	}
	return ctx
}

// Creates and registers two clients for the changeset resource
func buildChangeSetClients(config *rest.Config) (v1.ChangeSetClient, v1.ChangeSetClient, error) {
	csClient1, err := v1.NewChangeSetClient(&factory.KubeResourceClientFactory{
		Crd:         v1.ChangeSetCrd,
		Cfg:         config,
		SharedCache: kube.NewKubeCache(),
	})
	if err != nil {
		return nil, nil, err
	}

	err = csClient1.Register()
	if err != nil {
		return nil, nil, err
	}

	csClient2, err := v1.NewChangeSetClient(&factory.KubeResourceClientFactory{
		Crd:         v1.ChangeSetCrd,
		Cfg:         config,
		SharedCache: kube.NewKubeCache(),
	})
	if err != nil {
		return nil, nil, err
	}

	err = csClient1.Register()
	if err != nil {
		return nil, nil, err
	}

	return csClient1, csClient2, nil
}
