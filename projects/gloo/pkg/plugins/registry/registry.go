package plugins

import (
	"github.com/solo-io/solo-kit/projects/gloo/pkg/plugins/aws"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/plugins/azure"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/plugins/transformation"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/plugins/kubernetes"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/plugins/rest"
)

type registry struct {
	plugins []Plugin
}

var globalRegistry = func() *registry{
	transformationPlugin := transformation.NewPlugin()
	return &registry{
		plugins: []Plugin{
			// plugins should be added here
			aws.NewPlugin(),
			azure.NewPlugin(&transformationPlugin.RequireTransformationFilter),
			kubernetes.NewPlugin(),
			rest.NewPlugin(&transformationPlugin.RequireTransformationFilter),
			transformationPlugin,
		},
	}
}()

func Plugins() []Plugin {
	return globalRegistry.plugins
}
