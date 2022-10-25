package syncer

import (
	"github.com/solo-io/gloo/projects/discovery/pkg/fds"
	"github.com/solo-io/gloo/projects/discovery/pkg/fds/syncer"
	"github.com/solo-io/gloo/projects/gloo/pkg/bootstrap"
	"github.com/solo-io/solo-projects/projects/discovery/pkg/fds/discoveries/grpc-graphql"
	openApi "github.com/solo-io/solo-projects/projects/discovery/pkg/fds/discoveries/openapi-graphql"
)

func GetFDSEnterpriseExtensions() syncer.Extensions {
	return syncer.Extensions{
		DiscoveryFactoryFuncs: []func(opts bootstrap.Opts) fds.FunctionDiscoveryFactory{
			openApi.NewFunctionDiscoveryFactory,
			grpc.NewFunctionDiscoveryFactory,
		},
	}
}
