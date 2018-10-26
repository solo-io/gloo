package registry

import (
	"github.com/solo-io/solo-kit/projects/gloo/pkg/bootstrap"
	"reflect"
	"testing"
)

func TestPlugins(t *testing.T) {
	opts := bootstrap.Opts{}
	plugins := Plugins(opts)
	pluginTypes := make(map[reflect.Type]int)
	for index,plugin := range plugins {
		pluginType := reflect.TypeOf(plugin)
		pluginTypes[pluginType] = index
	}
	if len(plugins) > len(pluginTypes) {
		t.Errorf("Multiple plugins with the same type.")
	}
}