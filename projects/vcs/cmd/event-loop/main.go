package main

import (
	"context"
	"flag"
	"fmt"
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

// start an event loop watching the changeset resources in the glo
func run() error {

	configPath := flag.String("kube-config", "", "Path to the KUBECONFIG file. Leave blank for in-cluster configuration")
	flag.Parse()

	csClient, err := buildChangeSetClient(configPath)
	if err != nil {
		return err
	}

	ctx := contextutils.WithLogger(context.Background(), "vcs")
	errs, err := v1.NewApiEventLoop(v1.NewApiEmitter(csClient), &mockApiSyncer{}).Run(
		[]string{defaults.GlooSystem},
		clients.WatchOpts{Ctx: ctx})

	if err != nil {
		return err
	}

	for err := range errs {
		contextutils.LoggerFrom(ctx).Errorf("error in setup: %v", err)
	}

	return nil
}

// Creates and registers a client for the changeset resource
func buildChangeSetClient(configPath *string) (v1.ChangeSetClient, error) {
	cfg, err := kubeutils.GetConfig("", *configPath)
	if err != nil {
		return nil, err
	}

	csClient, err := v1.NewChangeSetClient(&factory.KubeResourceClientFactory{
		Crd:         v1.ChangeSetCrd,
		Cfg:         cfg,
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

type mockApiSyncer struct {
	synced bool
}

func (s *mockApiSyncer) Sync(ctx context.Context, snap *v1.ApiSnapshot) error {

	// iterate over all available changesets
	for _, cs := range snap.Changesets[defaults.GlooSystem] {
		fmt.Printf("Changeset with name [%v] has to be commited? [%v]\n", cs.Metadata.Name, cs.CommitPending.GetValue())

		// TODO: If commit_pending == true -> push changes to github
		if cs.CommitPending.GetValue() {
			fmt.Printf("Committing changeset: [%v]\n", cs.Metadata.Name)
		}
	}
	s.synced = true

	return nil
}
