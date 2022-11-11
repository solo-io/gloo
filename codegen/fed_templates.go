package codegen

import (
	"github.com/gobuffalo/packr"
	"github.com/solo-io/skv2/codegen/model"
)

var fedTemplatesBox = packr.NewBox("./custom_templates/federated_resource_templates")

var FederatedResourceTemplates = func() []model.CustomTemplates {
	fedReconciler, err := fedTemplatesBox.FindString("reconcilers/federated_reconcilers.gotmpl")
	if err != nil {
		panic(err)
	}

	fedReconcilerRunner, err := fedTemplatesBox.FindString("reconcilers/federated_reconciler_runner.gotmpl")
	if err != nil {
		panic(err)
	}

	fedClusterHandler, err := fedTemplatesBox.FindString("reconcilers/federated_resource_cluster_handler.gotmpl")
	if err != nil {
		panic(err)
	}

	needsReconcile, err := fedTemplatesBox.FindString("reconcilers/federated_needs_reconcile.gotmpl")
	if err != nil {
		panic(err)
	}

	fedApiserverHandler, err := fedTemplatesBox.FindString("apiserver/federated_apiserver_handler.gotmpl")
	if err != nil {
		panic(err)
	}
	fedApiserverProtos, err := fedTemplatesBox.FindString("apiserver/federated_apiserver_protos.gotmpl")
	if err != nil {
		panic(err)
	}

	statusExtensions, err := fedTemplatesBox.FindString("types/status_extensions.gotmpl")
	if err != nil {
		panic(err)
	}

	return []model.CustomTemplates{
		{
			Templates: map[string]string{
				"federation/federation_reconcilers.go": fedReconciler,
			},
			Funcs: GetTemplateFuncs(),
		},
		{
			Templates: map[string]string{
				"federation/federation_reconciler_runner.go": fedReconcilerRunner,
			},
			Funcs: GetTemplateFuncs(),
		},
		{
			Templates: map[string]string{
				"federation/federation_cluster_handler.go": fedClusterHandler,
			},
			Funcs: GetTemplateFuncs(),
		},
		{
			Templates: map[string]string{
				"federated_resource_needs_reconcile.go": needsReconcile,
			},
		},
		{
			Templates: map[string]string{
				"handler/fed_handler.go": fedApiserverHandler,
			},
			Funcs: GetTemplateFuncs(),
		},
		{
			Templates: map[string]string{
				"resource_apis.proto": fedApiserverProtos,
			},
			Funcs: GetTemplateFuncs(),
		},

		{
			Templates: map[string]string{
				"types/status_extensions.gen.go": statusExtensions,
			},
			Funcs: GetTemplateFuncs(),
		},
	}
}()
