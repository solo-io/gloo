package codegen

import (
	"encoding/json"
	"io/ioutil"
	"strings"

	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/golang/protobuf/protoc-gen-go/plugin"
	"github.com/iancoleman/strcase"
	"github.com/pseudomuto/protokit"
	"github.com/solo-io/solo-kit/pkg/errors"
	"github.com/solo-io/solo-kit/pkg/utils/log"
)

func ParseRequest(req *plugin_go.CodeGeneratorRequest) (*Project, error) {
	log.Printf("parsing request %v", req.FileToGenerate, req.GetParameter())
	params := req.GetParameter()
	if params == "" {
		return nil, errors.Errorf("must provide path to project.json file with --solo-kit_out=${PWD}/project.json:${OUTDIR}")
	}

	log.Printf("got cli param from protoc invoke: %v", params)
	projectConfig, err := loadProjectConfig(params)
	if err != nil {
		return nil, err
	}

	descriptors := protokit.ParseCodeGenRequest(req)
	var messages []*protokit.Descriptor
	for _, file := range descriptors {
		messages = append(messages, file.GetMessages()...)
	}

	var services []*protokit.ServiceDescriptor
	for _, file := range descriptors {
		services = append(services, file.GetServices()...)
	}

	var groupName string
	for _, desc := range descriptors {
		if groupName == "" {
			groupName = desc.GetPackage()
		}
		if groupName != desc.GetPackage() {
			return nil, errors.Errorf("package conflict: %v must match %v", groupName, desc.GetPackage())
		}
	}

	project := &Project{
		ProjectConfig: projectConfig,
		GroupName:     groupName,
	}
	resources, resourceGroups, err := getResources(project, messages)
	if err != nil {
		return nil, err
	}

	xdsResources, err := getXdsResources(project, messages, services)
	if err != nil {
		return nil, err
	}

	project.Resources = resources
	project.ResourceGroups = resourceGroups
	project.XDSResources = xdsResources

	return project, nil
}

func loadProjectConfig(path string) (ProjectConfig, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return ProjectConfig{}, err
	}
	var pc ProjectConfig
	err = json.Unmarshal(b, &pc)
	return pc, err
}

func goName(n string) string {
	return strcase.ToCamel(strings.Split(n, ".")[0])
}

func collectFields(msg *protokit.Descriptor) []*Field {
	var fields []*Field
	for _, f := range msg.GetField() {
		fields = append(fields, &Field{
			Name:     f.GetName(),
			TypeName: f.GetTypeName(),
		})
	}
	log.Printf("%v", fields)
	return fields
}

func hasField(msg *protokit.Descriptor, fieldName, fieldType string) bool {
	for _, field := range msg.Fields {
		if field.GetName() == fieldName && field.GetTypeName() == fieldType {
			return true
		}
	}
	return false
}

func hasPrimitiveField(msg *protokit.Descriptor, fieldName string, fieldType descriptor.FieldDescriptorProto_Type) bool {
	for _, field := range msg.Fields {
		if field.GetName() == fieldName && field.GetType() == fieldType {
			return true
		}
	}
	return false
}

func getCommentValue(comments []string, key string) (string, bool) {
	for _, c := range comments {
		if strings.HasPrefix(c, key) {
			return strings.TrimPrefix(c, key), true
		}
	}
	return "", false
}
