package protoc

import (
	"strings"

	"unicode"
	"unicode/utf8"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/plugin"
	"github.com/iancoleman/strcase"
	"github.com/pseudomuto/protokit"
	"github.com/solo-io/solo-kit/pkg/code-generator/typed"
)

// plugin is an implementation of protokit.Plugin
type Plugin struct{}

type resouceTemplateFunc func(params typed.ResourceLevelTemplateParams) (string, error)
type packageTemplateFunc func(params typed.PackageLevelTemplateParams) (string, error)

var resourceFilesToGenerate = map[string]resouceTemplateFunc{
	"_client.go":           typed.GenerateTypedClientCode,
	"_client_kube_test.go": typed.GenerateTypedClientKubeTestCode,
}

var packageFilesToGenerate = map[string]packageTemplateFunc{
	"_suite_test.go": typed.GenerateTestSuiteCode,
	"_cache.go":      typed.GenerateCacheCode,
	"_cache_test.go": typed.GenerateCacheTestCode,
}

func (p *Plugin) Generate(req *plugin_go.CodeGeneratorRequest) (*plugin_go.CodeGeneratorResponse, error) {
	descriptors := protokit.ParseCodeGenRequest(req)

	resp := new(plugin_go.CodeGeneratorResponse)

	for _, d := range descriptors {
		packageName := goPackage(d)
		var resourceTypes []string
		for _, msg := range d.Messages {
			resourceType := msg.GetName()
			var fields []string
			for _, field := range msg.Fields {
				// exclude special fields
				if field.GetName() == "status" || field.GetName() == "metadata" {
					continue
				}
				fields = append(fields, strcase.ToCamel(field.GetName()))
			}
			params := codegenParams(packageName, msg.GetComments(), resourceType, fields)
			if params != nil {
				resourceTypes = append(resourceTypes, resourceType)
				for suffix, genFunc := range resourceFilesToGenerate {
					file, err := generateResourceLevelFile(*params, suffix, genFunc)
					if err != nil {
						return nil, err
					}
					resp.File = append(resp.File, file)
				}
			}
		}
		if len(resourceTypes) > 0 {
			params := typed.PackageLevelTemplateParams{
				PackageName:   packageName,
				ResourceTypes: resourceTypes,
			}
			for suffix, genFunc := range packageFilesToGenerate {
				file, err := generatePackageLevelFile(params, suffix, genFunc)
				if err != nil {
					return nil, err
				}
				resp.File = append(resp.File, file)
			}
		}
	}
	return resp, nil
}

func codegenParams(packageName string, comments *protokit.Comment, resourceType string, fields []string) *typed.ResourceLevelTemplateParams {
	magicComments := strings.Split(comments.Leading, "\n")
	var (
		isResource bool
		shortName  string
		pluralName string
		groupName  string
		version    string
	)
	for _, comment := range magicComments {
		if comment == "@solo-kit:resource" {
			isResource = true
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
	return &typed.ResourceLevelTemplateParams{
		PackageName:           packageName,
		ResourceType:          resourceType,
		ResourceTypeLowerCase: strcase.ToLowerCamel(resourceType),
		ShortName:             shortName,
		PluralName:            pluralName,
		GroupName:             groupName,
		Version:               version,
		Fields:                fields,
	}
}

func generateResourceLevelFile(params typed.ResourceLevelTemplateParams, suffix string, genFunc resouceTemplateFunc) (*plugin_go.CodeGeneratorResponse_File, error) {
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

func generatePackageLevelFile(params typed.PackageLevelTemplateParams, suffix string, genFunc packageTemplateFunc) (*plugin_go.CodeGeneratorResponse_File, error) {
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

// goPackageOption interprets the file's go_package option.
// If there is no go_package, it returns ("", "", false).
// If there's a simple name, it returns ("", pkg, true).
// If the option implies an import path, it returns (impPath, pkg, true).
func goPackage(d *protokit.FileDescriptor) string {
	opt := d.GetOptions().GetGoPackage()
	if opt == "" {
		return ""
	}
	// A semicolon-delimited suffix delimits the import path and package name.
	sc := strings.Index(opt, ";")
	if sc >= 0 {
		return cleanPackageName(opt[sc+1:])
	}
	// The presence of a slash implies there's an import path.
	slash := strings.LastIndex(opt, "/")
	if slash >= 0 {
		return cleanPackageName(opt[slash+1:])
	}
	return cleanPackageName(opt)
}

var isGoKeyword = map[string]bool{
	"break":       true,
	"case":        true,
	"chan":        true,
	"const":       true,
	"continue":    true,
	"default":     true,
	"else":        true,
	"defer":       true,
	"fallthrough": true,
	"for":         true,
	"func":        true,
	"go":          true,
	"goto":        true,
	"if":          true,
	"import":      true,
	"interface":   true,
	"map":         true,
	"package":     true,
	"range":       true,
	"return":      true,
	"select":      true,
	"struct":      true,
	"switch":      true,
	"type":        true,
	"var":         true,
}

func badToUnderscore(r rune) rune {
	if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' {
		return r
	}
	return '_'
}

func cleanPackageName(name string) string {
	name = strings.Map(badToUnderscore, name)
	// Identifier must not be keyword: insert _.
	if isGoKeyword[name] {
		name = "_" + name
	}
	// Identifier must not begin with digit: insert _.
	if r, _ := utf8.DecodeRuneInString(name); unicode.IsDigit(r) {
		name = "_" + name
	}
	return name
}
