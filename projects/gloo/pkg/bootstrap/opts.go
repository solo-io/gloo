package bootstrap

import (
	"net"

	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/factory"
	"github.com/solo-io/solo-kit/pkg/namespacing"
	"google.golang.org/grpc"
	"k8s.io/client-go/kubernetes"
)

type Opts struct {
	WriteNamespace string
	Upstreams      factory.ResourceClientFactoryOpts
	Proxies        factory.ResourceClientFactoryOpts
	Secrets        factory.ResourceClientFactoryOpts
	Artifacts      factory.ResourceClientFactoryOpts
	Namespacer     namespacing.Namespacer
	BindAddr       net.Addr
	GrpcServer     *grpc.Server
	KubeClient     kubernetes.Interface
	WatchOpts      clients.WatchOpts
	DevMode        bool
}
