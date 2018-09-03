package registry

import (
	"github.com/solo-io/solo-kit/projects/gloo/pkg/plugins"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/plugins/aws"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/plugins/azure"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/plugins/kubernetes"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/plugins/rest"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/plugins/transformation"
)

type registry struct {
	plugins []plugins.Plugin
}

var globalRegistry = func() *registry {
	transformationPlugin := transformation.NewPlugin()
	return &registry{
		plugins: []plugins.Plugin{
			// plugins should be added here
			aws.NewPlugin(),
			azure.NewPlugin(&transformationPlugin.RequireTransformationFilter),
			kubernetes.NewPlugin(),
			rest.NewPlugin(&transformationPlugin.RequireTransformationFilter),
			transformationPlugin,
		},
	}
}()

func Plugins() []plugins.Plugin {
	return globalRegistry.plugins
}
