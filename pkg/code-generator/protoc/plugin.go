package protoc

import (
	"strings"

	"unicode"
	"unicode/utf8"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/plugin"
	"github.com/pseudomuto/protokit"
	"github.com/solo-io/solo-kit/pkg/code-generator/typed"
)

// plugin is an implementation of protokit.Plugin
type Plugin struct{}

type generateCodeFunc func(params typed.CodeGeneratorParams) (string, error)

var filesToGenerate = map[string]generateCodeFunc{
	"_client.go":           typed.GenerateTypedClientCode,
	"_client_kube_test.go": typed.GenerateTypedClientKubeTestCode,
	"_suite_test.go":       typed.GenerateTypedClientTestSuiteCode,
}

func (p *Plugin) Generate(req *plugin_go.CodeGeneratorRequest) (*plugin_go.CodeGeneratorResponse, error) {
	descriptors := protokit.ParseCodeGenRequest(req)

	resp := new(plugin_go.CodeGeneratorResponse)

	for _, d := range descriptors {
		for _, msg := range d.Messages {
			params := codegenParams(d, msg.GetComments(), msg.GetName())
			if params != nil {
				for suffix, genFunc := range filesToGenerate {
					file, err := generateFile(d, *params, suffix, genFunc)
					if err != nil {
						return nil, err
					}
					resp.File = append(resp.File, file)
				}
			}
		}
	}

	return resp, nil
}

func codegenParams(d *protokit.FileDescriptor, comments *protokit.Comment, messageName string) *typed.CodeGeneratorParams {
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
	return &typed.CodeGeneratorParams{
		PackageName:  goPackage(d),
		ResourceType: messageName,
		ShortName:    shortName,
		PluralName:   pluralName,
		GroupName:    groupName,
		Version:      version,
	}
}

func generateFile(d *protokit.FileDescriptor, params typed.CodeGeneratorParams, suffix string, genFunc generateCodeFunc) (*plugin_go.CodeGeneratorResponse_File, error) {
	fileName := strings.TrimSuffix(d.GetName(), ".proto") + suffix
	content, err := genFunc(params)
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
