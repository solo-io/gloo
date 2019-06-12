package utils

import (
	"github.com/gogo/protobuf/proto"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins/utils"
)

type tmpPluginContainer struct {
	params plugins.InitParams
}

func (t *tmpPluginContainer) GetExtensions() *v1.Extensions {
	return t.params.ExtensionsSettings
}

func GetSettings(params plugins.InitParams, name string, settings proto.Message) (bool, error) {
	err := utils.UnmarshalExtension(&tmpPluginContainer{params}, name, settings)
	if err != nil {
		if err == utils.NotFoundError {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
