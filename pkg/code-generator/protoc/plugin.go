package protoc

import (
	"strings"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/plugin"
	"github.com/iancoleman/strcase"
	"github.com/pseudomuto/protokit"
	"github.com/solo-io/solo-kit/pkg/code-generator/templates"
)

const (
	statusTypeName = ".core.solo.io.Status"
)

// plugin is an implementation of protokit.Plugin
type Plugin struct{}

type resouceTemplateFunc func(params templates.ResourceLevelTemplateParams) (string, error)
type packageTemplateFunc func(params templates.PackageLevelTemplateParams) (string, error)

var resourceFilesToGenerate = map[string]resouceTemplateFunc{
	"_client.go":      templates.GenerateResourceClientCode,
	"_client_test.go": templates.GenerateResourceClientTestCode,
	"_reconciler.go":  templates.GenerateReconcilerCode,
}

var packageFilesToGenerate = map[string]packageTemplateFunc{
	"_suite_test.go":      templates.GenerateTestSuiteCode,
	"_cache.go":           templates.GenerateCacheCode,
	"_cache_test.go":      templates.GenerateCacheTestCode,
	"_event_loop.go":      templates.GenerateEventLoopCode,
	"_event_loop_test.go": templates.GenerateEventLoopTestCode,
}

func (p *Plugin) Generate(req *plugin_go.CodeGeneratorRequest) (*plugin_go.CodeGeneratorResponse, error) {
	descriptors := protokit.ParseCodeGenRequest(req)

	resp := new(plugin_go.CodeGeneratorResponse)

	resourcesByPackage := make(map[string][]*templates.ResourceLevelTemplateParams)

	resourceDescriptors := make(map[string]*protokit.Descriptor)

	for _, d := range descriptors {
		packageName := goPackage(d)
		var resourceTypes []*templates.ResourceLevelTemplateParams
		for _, msg := range d.Messages {
			resourceType := msg.GetName()
			var fields []string
			for _, field := range msg.Fields {
				if field.GetComments().GetLeading() == "@solo-kit:ignore_field=true" {
					continue
				}
				// exclude special fields
				if field.GetName() == "status" || field.GetName() == "metadata" {
					continue
				}
				fields = append(fields, strcase.ToCamel(field.GetName()))
			}
			params := codegenParams(packageName, msg, resourceType, fields)
			if params != nil {
				resourceTypes = append(resourceTypes, params)
				for suffix, genFunc := range resourceFilesToGenerate {
					file, err := generateResourceLevelFile(*params, suffix, genFunc)
					if err != nil {
						return nil, err
					}
					resp.File = append(resp.File, file)
				}
				resourceDescriptors[msg.GetName()] = msg
			}
		}
		if len(resourceTypes) > 0 {
			resourcesByPackage[packageName] = append(resourcesByPackage[packageName], resourceTypes...)
		}
	}

	for packageName, resourceParams := range resourcesByPackage {
		var (
			resourceTypes       []string
			resourceLevelParams = make(map[string]*templates.ResourceLevelTemplateParams)
		)
		for _, resource := range resourceParams {
			resourceTypes = append(resourceTypes, resource.ResourceType)
			resourceLevelParams[resource.ResourceType] = resource
		}
		params := templates.PackageLevelTemplateParams{
			PackageName:         packageName,
			ResourceTypes:       resourceTypes,
			ResourceLevelParams: resourceLevelParams,
		}
		for suffix, genFunc := range packageFilesToGenerate {
			file, err := generatePackageLevelFile(params, suffix, genFunc)
			if err != nil {
				return nil, err
			}
			resp.File = append(resp.File, file)
		}
	}

	//cmdFiles, err := genCmd(resourceDescriptors)
	//if err != nil {
	//	return nil, err
	//}
	//resp.File = append(resp.File, cmdFiles...)

	return resp, nil
}

func codegenParams(packageName string, msg *protokit.Descriptor, resourceType string, fields []string) *templates.ResourceLevelTemplateParams {
	magicComments := strings.Split(msg.GetComments().Leading, "\n")
	var (
		isResource  bool
		isDataType  bool
		isInputType bool
		shortName   string
		pluralName  string
		groupName   string
		version     string
	)
	for _, field := range msg.Fields {
		if field.GetTypeName() == statusTypeName {
			isInputType = true
			break
		}
	}
	for _, comment := range magicComments {
		if comment == "@solo-kit:resource" {
			isResource = true
			continue
		}
		if comment == "@solo-kit:resource.data_type" {
			isDataType = true
			continue
		}
		if strings.HasPrefix(comment, "@solo-kit:resource.short_name=") {
			shortName = strings.TrimPrefix(comment, "@solo-kit:resource.short_name=")
			continue
		}
		if strings.HasPrefix(comment, "@solo-kit:resource.plural_name=") {
			pluralName = strings.TrimPrefix(comment, "@solo-kit:resource.plural_name=")
			continue
		}
		if strings.HasPrefix(comment, "@solo-kit:resource.group_name=") {
			groupName = strings.TrimPrefix(comment, "@solo-kit:resource.group_name=")
			continue
		}
		if strings.HasPrefix(comment, "@solo-kit:resource.version=") {
			version = strings.TrimPrefix(comment, "@solo-kit:resource.version=")
			continue
		}
	}
	if !isResource {
		return nil
	}
	return &templates.ResourceLevelTemplateParams{
		PackageName:           packageName,
		ResourceType:          resourceType,
		IsDataType:            isDataType,
		IsInputType:           isInputType,
		ResourceTypeLowerCase: strcase.ToLowerCamel(resourceType),
		ShortName:             shortName,
		PluralName:            pluralName,
		GroupName:             groupName,
		Version:               version,
		Fields:                fields,
	}
}

func generateResourceLevelFile(params templates.ResourceLevelTemplateParams, suffix string, genFunc resouceTemplateFunc) (*plugin_go.CodeGeneratorResponse_File, error) {
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

func generatePackageLevelFile(params templates.PackageLevelTemplateParams, suffix string, genFunc packageTemplateFunc) (*plugin_go.CodeGeneratorResponse_File, error) {
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
