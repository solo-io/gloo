package plugins

import (
	"log"
)

var defaultRegistry = &registry{}

// TODO(yuval-k): delete this when we can
func Register(plugin Plugin) {
	RegisterFunc(func() Plugin { return plugin })
}

func RegisterFunc(plugin func() Plugin) {
	if plugin == nil {
		log.Fatalf("plugin can not be nil")
	}
	defaultRegistry.plugins = append(defaultRegistry.plugins, plugin)
}

func RegisteredPlugins(initparams InitParams) func() []Plugin {
	return func() []Plugin {
		var plugins []Plugin
		for _, pc := range defaultRegistry.plugins {
			p := pc()
			// TODO(yuval-k): not ignore error
			p.Init(initparams)
			plugins = append(plugins, p)
		}
		return plugins
	}
}

type registry struct {
	plugins []func() Plugin
}
