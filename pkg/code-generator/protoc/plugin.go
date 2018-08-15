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

const (
	statusTypeName = ".core.solo.io.Status"
)

// plugin is an implementation of protokit.Plugin
type Plugin struct{}

type resouceTemplateFunc func(params typed.ResourceLevelTemplateParams) (string, error)
type packageTemplateFunc func(params typed.PackageLevelTemplateParams) (string, error)

var resourceFilesToGenerate = map[string]resouceTemplateFunc{
	"_client.go":      typed.GenerateResourceClientCode,
	"_client_test.go": typed.GenerateResourceClientTestCode,
	"_reconciler.go":  typed.GenerateReconcilerCode,
}

var packageFilesToGenerate = map[string]packageTemplateFunc{
	"_suite_test.go":      typed.GenerateTestSuiteCode,
	"_cache.go":           typed.GenerateCacheCode,
	"_cache_test.go":      typed.GenerateCacheTestCode,
	"_event_loop.go":      typed.GenerateEventLoopCode,
	"_event_loop_test.go": typed.GenerateEventLoopTestCode,
}

var descriptors []*protokit.FileDescriptor

func (p *Plugin) Generate(req *plugin_go.CodeGeneratorRequest) (*plugin_go.CodeGeneratorResponse, error) {
	descriptors = protokit.ParseCodeGenRequest(req)

	resp := new(plugin_go.CodeGeneratorResponse)

	resourcesByPackage := make(map[string][]*typed.ResourceLevelTemplateParams)

	resourceDescriptors := make(map[string]*protokit.Descriptor)

	for _, d := range descriptors {
		packageName := goPackage(d)
		var resourceTypes []*typed.ResourceLevelTemplateParams
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
			resourceLevelParams = make(map[string]*typed.ResourceLevelTemplateParams)
		)
		for _, resource := range resourceParams {
			resourceTypes = append(resourceTypes, resource.ResourceType)
			resourceLevelParams[resource.ResourceType] = resource
		}
		params := typed.PackageLevelTemplateParams{
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

//type message struct {
//	name   string
//	fields []field
//}
//
//type field struct {
//	name  string
//	ftype string
//	// nil if scalar
//	messagetype *message
//}
//
//func findMessage(typeName string) *protokit.Descriptor {
//
//}
//
//func constructTree(topLevel *protokit.Descriptor) *message {
//	msg := &message{
//		name: topLevel.GetName(),
//	}
//	for _, f := range topLevel.GetField() {
//		var fieldMsg *message
//		if f.GetType() == descriptor.FieldDescriptorProto_TYPE_MESSAGE {
//			fieldMsg = constructTree(findMessage(f.GetTypeName()))
//		}
//		var repeated string
//		if f.GetLabel() == descriptor.FieldDescriptorProto_LABEL_REPEATED {
//			repeated = "repeated "
//		}
//		field := field{
//			name:        f.GetName(),
//			ftype:       repeated + f.GetTypeName(),
//			messagetype: fieldMsg,
//		}
//	}
//	return msg
//}
//
//func findProtoMessage(msgTypeName string) (*protokit.Descriptor, error) {
//	packageName := strings.TrimPrefix(msgTypeName, ".")
//	for _, d := range descriptors {
//		if d.GetPackage()
//		log.DefaultOut = os.Stderr
//		log.Printf("%v", d)
//		time.Sleep(time.Minute)
//	}
//	return nil, nil
//}
//
//func genCmd(resourceDescriptors map[string]*protokit.Descriptor) ([]*plugin_go.CodeGeneratorResponse_File, error) {
//	var files []*plugin_go.CodeGeneratorResponse_File
//	for resourceName, descriptor := range resourceDescriptors {
//		cmdFile, err := genCmdFile(descriptor)
//		if err != nil {
//			return nil, err
//		}
//		files = append(files, &plugin_go.CodeGeneratorResponse_File{
//			Name:    proto.String("cmd/" + resourceName + ".json"),
//			Content: proto.String(cmdFile),
//		})
//	}
//	return files, nil
//}
//
//func toString(fields []*protokit.FieldDescriptor) string {
//	var msgs []string
//	for _, f := range fields {
//		msgs = append(msgs, fmt.Sprintf("- Name: %v", f.GetName()))
//		if f.GetLabel() == descriptor.FieldDescriptorProto_LABEL_REPEATED {
//			msgs = append(msgs, fmt.Sprintf("  Type: []%v", f.GetTypeName()))
//		} else {
//			msgs = append(msgs, fmt.Sprintf("  Type: %v", f.GetTypeName()))
//		}
//		msgs = append(msgs, fmt.Sprintf("  Field Type: %v", f.GetType()))
//	}
//	return strings.Join(msgs, "\n")
//}
//
//func genCmdFile(message *protokit.Descriptor) (string, error) {
//	var msgs []string
//	for _, v := range [][]interface{}{
//		{"message.GetName:\n\t", message.GetName()},
//		{"message.GetComments:\n\t", message.GetComments()},
//		{"message.GetOptions:\n\t", message.GetOptions()},
//		{"message.GetOneofDecl:\n\t", message.GetOneofDecl()},
//		{"message.GetMessageFields:\n", toString(message.GetMessageFields())},
//	} {
//		msgs = append(msgs, fmt.Sprintf("%v%v", v[0].(string), v[1]))
//	}
//	return strings.Join(msgs, "\n\n"), nil
//}

func codegenParams(packageName string, msg *protokit.Descriptor, resourceType string, fields []string) *typed.ResourceLevelTemplateParams {
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
	return &typed.ResourceLevelTemplateParams{
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
