package plugins

import (
	"log"
)

var defaultRegistry = &registry{}

func Register(plugin Plugin) {
	if plugin == nil {
		log.Fatalf("plugin can not be nil")
	}
	defaultRegistry.plugins = append(defaultRegistry.plugins, plugin)
}

func RegisteredPlugins() []Plugin {
	return defaultRegistry.plugins
}

type registry struct {
	plugins []Plugin
}
