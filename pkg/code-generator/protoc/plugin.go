package protoc

import (
	"bytes"
	"os"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/plugin"
	"github.com/iancoleman/strcase"
	"github.com/solo-io/solo-kit/pkg/code-generator/codegen"
	"github.com/solo-io/solo-kit/pkg/utils/log"
)

// plugin is an implementation of protokit.Plugin
type Plugin struct{}

type resouceTemplateFunc func(params ResourceLevelTemplateParams) (string, error)
type packageTemplateFunc func(params PackageLevelTemplateParams) (string, error)

var resourceFilesToGenerate = map[string]resouceTemplateFunc{
//"_client.go":      templates.GenerateResourceClientCode,
//"_client_test.go": templates.GenerateResourceClientTestCode,
//"_reconciler.go":  templates.GenerateReconcilerCode,
}

var packageFilesToGenerate = map[string]packageTemplateFunc{
//"_suite_test.go":      templates.GenerateTestSuiteCode,
//"_cache.go":           templates.GenerateCacheCode,
//"_cache_test.go":      templates.GenerateCacheTestCode,
//"_event_loop.go":      templates.GenerateEventLoopCode,
//"_event_loop_test.go": templates.GenerateEventLoopTestCode,
}

func (p *Plugin) Generate(req *plugin_go.CodeGeneratorRequest) (*plugin_go.CodeGeneratorResponse, error) {
	log.DefaultOut = &bytes.Buffer{}
	if false {
		log.DefaultOut = os.Stderr
	}
	project, err := codegen.ParseRequest(req)
	if err != nil {
		return nil, err
	}

	code, err := codegen.GenerateFiles(project)
	if err != nil {
		return nil, err
	}

	log.Printf("%v", project)
	log.Printf("%v", code)

	resp := new(plugin_go.CodeGeneratorResponse)

	for _, file := range code {
		resp.File = append(resp.File, &plugin_go.CodeGeneratorResponse_File{
			Name:    proto.String(file.Filename),
			Content: proto.String(file.Content),
		})
	}

	return resp, nil
}

func generateResourceLevelFile(params ResourceLevelTemplateParams, suffix string, genFunc resouceTemplateFunc) (*plugin_go.CodeGeneratorResponse_File, error) {
	fileName := strcase.ToSnake(params.ResourceType) + suffix
	content, err := genFunc(params)
	if err != nil {
		return nil, err
	}
	return &plugin_go.CodeGeneratorResponse_File{
		Name:    proto.String(fileName),
		Content: proto.String(content),
	}, nil
}

func generatePackageLevelFile(params PackageLevelTemplateParams, suffix string, genFunc packageTemplateFunc) (*plugin_go.CodeGeneratorResponse_File, error) {
	fileName := params.PackageName + suffix
	content, err := genFunc(params)
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	return &plugin_go.CodeGeneratorResponse_File{
		Name:    proto.String(fileName),
		Content: proto.String(content),
	}, nil
}
