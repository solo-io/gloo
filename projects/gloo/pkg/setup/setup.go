package setup

import (
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/factory"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/memory"
)

func Setup(inputResourceOpts factory.ResourceClientFactoryOpts, secretOpts factory.ResourceClientFactoryOpts, artifactOpts factory.ResourceClientFactoryOpts) error {
	inputFactory := factory.NewResourceClientFactory(inputResourceOpts)
	secretFactory := factory.NewResourceClientFactory(secretOpts)
	artifactFactory := factory.NewResourceClientFactory(artifactOpts)
	// endpoints are internal-only, therefore use the in-memory client
	endpointsFactory := factory.NewResourceClientFactory(&factory.MemoryResourceClientOpts{
		Cache: memory.NewInMemoryResourceCache(),
	})

	upstreamClient, err := v1.NewUpstreamClient(inputFactory)
	if err != nil {
		return err
	}

	proxyClient, err := v1.NewProxyClient(inputFactory)
	if err != nil {
		return err
	}

	endpointClient, err := v1.NewEndpointClient(endpointsFactory)
	if err != nil {
		return err
	}

	secretClient, err := v1.NewSecretClient(secretFactory)
	if err != nil {
		return err
	}

	artifactClient, err := v1.NewArtifactClient(artifactFactory)
	if err != nil {
		return err
	}

	// TODO: initialize endpointClient using EDS plugins
	cache := v1.NewCache(artifactClient, endpointClient, proxyClient, secretClient, upstreamClient)
	el := v1.NewEventLoop(cache, &syncer{})
}

type syncer struct{}

func (syncer) Sync(snap *v1.Snapshot) error {

}
