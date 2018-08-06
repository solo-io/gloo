package setup

import (
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/factory"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1"
)

func Setup(inputResourceOpts factory.ResourceClientFactoryOpts, secretOpts factory.ResourceClientFactoryOpts, artifactOpts factory.ResourceClientFactoryOpts) error {
	inputResourceClientFactory := factory.NewResourceClientFactory(inputResourceOpts)
	attributeClient, err := v1.NewAttributeClient(inputResourceClientFactory)
	roleClient, err := v1.NewRoleClient(inputResourceClientFactory)
	if err != nil {
		return err
	}
	upstreamClient, err := v1.NewUpstreamClient(inputResourceClientFactory)
	if err != nil {
		return err
	}
	virtualServiceClient, err := v1.NewVirtualServiceClient(inputResourceClientFactory)
	if err != nil {
		return err
	}

	secretClientFactory := factory.NewResourceClientFactory(secretOpts)
	secretClient, err := v1.NewSecretClient(secretClientFactory)
	if err != nil {
		return err
	}

	artifactClientFactory := factory.NewResourceClientFactory(artifactOpts)
	artifactClient, err := v1.NewArtifactClient(artifactClientFactory)
	if err != nil {
		return err
	}

	// TODO: initialize endpointClient using EDS plugins

	cache := v1.NewCache(artifactClient, attributeClient, endpointClient, roleClient, secretClient, upstreamClient, virtualServiceClient)
	el := v1.NewEventLoop(cache, &syncer{})
}

type syncer struct{}

func (syncer) Sync(snap *v1.Snapshot) error {

}
