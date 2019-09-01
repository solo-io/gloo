package utils

import (
	"github.com/gogo/protobuf/proto"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins/utils"
)

type tmpPluginContainer struct {
	extensions *v1.Extensions
}

func (t *tmpPluginContainer) GetExtensions() *v1.Extensions {
	return t.extensions
}

func GetSettings(params plugins.InitParams, name string, settings proto.Message) (bool, error) {
	return UnmarshalExtension(params.ExtensionsSettings, name, settings)
}

func UnmarshalExtension(ext *v1.Extensions, name string, settings proto.Message) (bool, error) {
	err := utils.UnmarshalExtension(&tmpPluginContainer{extensions: ext}, name, settings)
	if err != nil {
		if err == utils.NotFoundError {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
