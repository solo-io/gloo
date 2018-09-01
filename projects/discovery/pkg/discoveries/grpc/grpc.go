package grpc

import (
	"context"
	"net/url"

	discovery "github.com/solo-io/solo-kit/projects/discovery/pkg"

	"github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1/plugins"
)

func getswagspec(u *v1.Upstream) *rest_plugins.ServiceSpec_GrpcInfo {
	spec, ok := u.UpstreamSpec.UpstreamType.(specable)
	if !ok {
		return nil
	}
	restwrapper, ok := spec.GetServiceSpec().PluginType.(*plugins.ServiceSpec_Rest)
	if !ok {
		return nil
	}
	rest := restwrapper.Rest
	return rest.SwaggerInfo
}

type FunctionDiscoveryFactory struct{}

func NewFunctionDiscovery(u *v1.Upstream) discovery.UpstreamFunctionDiscovery {
	return &UpstreamFunctionDiscovery{
		upstream: u,
	}
}

type UpstreamFunctionDiscovery struct {
	upstream *v1.Upstream
}

func (f *UpstreamFunctionDiscovery) IsFunctional() bool {

}

func (f *UpstreamFunctionDiscovery) DetectType(ctx context.Context, url *url.URL) (*plugins.ServiceSpec, error) {

}

func (f *UpstreamFunctionDiscovery) DetectFunctions(ctx context.Context, secrets func() v1.SecretList, out func(discovery.UpstreamMutator) error) error {

}
