package registry

import (
	"github.com/solo-io/solo-kit/projects/gloo/pkg/bootstrap"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/plugins"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/plugins/aws"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/plugins/azure"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/plugins/grpc"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/plugins/kubernetes"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/plugins/rest"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/plugins/static"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/plugins/transformation"
)

type registry struct {
	plugins []plugins.Plugin
}

var globalRegistry = func(opts bootstrap.Opts) *registry {
	transformationPlugin := transformation.NewPlugin()
	return &registry{
		plugins: []plugins.Plugin{
			// plugins should be added here
			aws.NewPlugin(),
			azure.NewPlugin(&transformationPlugin.RequireTransformationFilter),
			kubernetes.NewPlugin(opts),
			rest.NewPlugin(&transformationPlugin.RequireTransformationFilter),
			static.NewPlugin(),
			transformationPlugin,
			grpc.NewPlugin(&transformationPlugin.RequireTransformationFilter),
		},
	}
}

func Plugins(opts bootstrap.Opts) []plugins.Plugin {
	return globalRegistry(opts).plugins
}
