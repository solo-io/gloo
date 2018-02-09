package translator

import (
	"github.com/solo-io/glue/pkg/api/types/v1"
	"github.com/solo-io/glue/pkg/plugin2"
)

type functionRouterPlugin struct {
	functionPlugins []plugin.FunctionPlugin
}

func (p *functionRouterPlugin) GetDependencies(_ v1.Config) *plugin.Dependencies {
	return nil
}
