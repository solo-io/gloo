package codegen

import (
	"github.com/gobuffalo/packr"
	"github.com/solo-io/skv2/codegen/model"
)

var baseTemplateBox = packr.NewBox("./custom_templates/base_gloo_resource_templates")

var BaseGlooResourceTemplates = func() []model.CustomTemplates {
	// This template generates the apiserver protobuf messages and grpc service interfaces used by the UI
	// to get and list Gloo resources.
	apiserverProtos, err := baseTemplateBox.FindString("apiserver/apiserver_protos.gotmpl")
	if err != nil {
		panic(err)
	}

	// These templates generate the handlers that implement the service interfaces defined in the protos.
	// We need a separate Gloo Fed and single-cluster implementation of each service, because clients need
	// to be initialized differently in each case (e.g. the Fed version needs a multi-cluster clientset
	// and a cluster name to access resources, which may be on remote clusters).
	// Note that for some resource types (e.g. GraphQLSchemas), we manually write the protos and handlers
	// instead of generating them from these templates, because they may have unique needs, such as being
	// able to return a different subset of data in the List APIs, needing to aggregate data from multiple
	// resources, or supporting create/update/delete functionality.
	fedApiserverHandler, err := baseTemplateBox.FindString("apiserver/fed_apiserver_handler.gotmpl")
	if err != nil {
		panic(err)
	}
	singleClusterApiserverHandler, err := baseTemplateBox.FindString("apiserver/single_cluster_apiserver_handler.gotmpl")
	if err != nil {
		panic(err)
	}

	// Templates for glooctl fed cli commands
	apiserverCliClient, err := baseTemplateBox.FindString("apiserver/apiserver_cli_client.gotmpl")
	if err != nil {
		panic(err)
	}

	// these are the templates that auto-generate the code that returns summaries (counts, errors, warnings) for the various resources.
	// the fed checker is used during gloo fed discovery and the single cluster checker is used for the single-cluster apiserver/UI
	fedChecker, err := baseTemplateBox.FindString("discovery_check/fed_check.gotmpl")
	if err != nil {
		panic(err)
	}
	// since the single-cluster checker is only needed for the apiserver/UI (rather than discovery, etc), we put the template
	// in the apiserver dir
	singleClusterChecker, err := baseTemplateBox.FindString("apiserver/single_cluster_check.gotmpl")
	if err != nil {
		panic(err)
	}

	return []model.CustomTemplates{
		{
			Templates: map[string]string{
				"check/discovery_check.go": fedChecker,
			},
			Funcs: GetTemplateFuncs(),
		},
		{
			Templates: map[string]string{
				"check/single_cluster_check.go": singleClusterChecker,
			},
			Funcs: GetTemplateFuncs(),
		},
		{
			Templates: map[string]string{
				"handler/fed_handler.go": fedApiserverHandler,
			},
			Funcs: GetTemplateFuncs(),
		},
		{
			Templates: map[string]string{
				"handler/single_cluster_handler.go": singleClusterApiserverHandler,
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
