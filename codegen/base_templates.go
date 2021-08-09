package codegen

import (
	"github.com/gobuffalo/packr"
	"github.com/solo-io/skv2/codegen/model"
)

var baseTemplateBox = packr.NewBox("./custom_templates/base_gloo_resource_templates")

var BaseGlooResourceTemplates = func() []model.CustomTemplates {
	fedApiserverHandler, err := baseTemplateBox.FindString("apiserver/fed_apiserver_handler.gotmpl")
	if err != nil {
		panic(err)
	}
	// TODO add edge_apiserver_handler here
	apiserverProtos, err := baseTemplateBox.FindString("apiserver/apiserver_protos.gotmpl")
	if err != nil {
		panic(err)
	}
	apiserverCliClient, err := baseTemplateBox.FindString("apiserver/apiserver_cli_client.gotmpl")
	if err != nil {
		panic(err)
	}
	checker, err := baseTemplateBox.FindString("discovery_check/check.gotmpl")
	if err != nil {
		panic(err)
	}

	return []model.CustomTemplates{
		{
			Templates: map[string]string{
				"check/discovery_check.go": checker,
			},
			Funcs: GetTemplateFuncs(),
		},
		{
			Templates: map[string]string{
				"handler/handler.go": fedApiserverHandler,
			},
			Funcs: GetTemplateFuncs(),
		},
		{
			Templates: map[string]string{
				"resource_apis.proto": apiserverProtos,
			},
			Funcs: GetTemplateFuncs(),
		},
		{
			Templates: map[string]string{
				"cli/get.go": apiserverCliClient,
			},
			Funcs: GetTemplateFuncs(),
		},
	}
}()
