package syncer

import (
	"github.com/solo-io/gloo/projects/discovery/pkg/fds"
	"github.com/solo-io/gloo/projects/discovery/pkg/fds/syncer"
	openApi "github.com/solo-io/solo-projects/projects/discovery/pkg/fds/discoveries/openapi"
)

func GetFDSEnterpriseExtensions() syncer.Extensions {
	return syncer.Extensions{
		DiscoveryFactoryFuncs: []func() fds.FunctionDiscoveryFactory{
			openApi.NewFunctionDiscoveryFactory,
		},
	}
}
