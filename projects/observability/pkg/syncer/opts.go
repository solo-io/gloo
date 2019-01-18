package syncer

import (
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/factory"
)

type Opts struct {
	WriteNamespace  string
	WatchNamespaces []string
	Upstreams       factory.ResourceClientFactory
	WatchOpts       clients.WatchOpts
	DevMode         bool
}
