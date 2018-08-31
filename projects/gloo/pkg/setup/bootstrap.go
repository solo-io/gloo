package setup

import (
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/factory"
	"github.com/solo-io/solo-kit/pkg/namespacing"
	"google.golang.org/grpc"
)

type Opts struct {
	writeNamespace  string
	configBackend   factory.ResourceClientFactoryOpts
	secretBackend   factory.ResourceClientFactoryOpts
	artifactBackend factory.ResourceClientFactoryOpts
	namespacer      namespacing.Namespacer
	grpcServer      *grpc.Server
	watchOpts       clients.WatchOpts
}
