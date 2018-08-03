package setup

import (
	"github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/factory"
)

func Setup() error {
	fact := factory.NewResourceClientFactory()
	attributeClient, err := v1.NewAttributeClient(fact)
	roleClient, err := v1.NewRoleClient(fact)
	if err != nil {
		return err
	}
	upstreamClient, err := v1.NewUpstreamClient(fact)
	if err != nil {
		return err
	}
	virtualServiceClient, err := v1.NewVirtualServiceClient(fact)
	if err != nil {
		return err
	}
	cache := v1.NewCache(attributeClient, roleClient, upstreamClient, virtualServiceClient)
	el := v1.NewEventLoop(cache, &syncer{})
}

type syncer struct {}

func (syncer) Sync(snap *v1.Snapshot) error {

}
