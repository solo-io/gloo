package grpc

import (
	"fmt"

	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	structpb "github.com/golang/protobuf/ptypes/struct"
	"github.com/graphql-go/graphql/language/ast"
	"github.com/jhump/protoreflect/desc"
	errors "github.com/rotisserie/eris"
)

type SchemaBuilder struct {
	QueryType     *ast.ObjectDefinition
	InputTypeDefs map[string]*ast.InputObjectDefinition
	TypeDefs      map[string]*ast.ObjectDefinition
	EnumDefs      map[string]*ast.EnumDefinition
}

func NewSchemaBuilder() *SchemaBuilder {
	return &SchemaBuilder{
		QueryType:     ast.NewObjectDefinition(&ast.ObjectDefinition{Name: CreateNameType("Query")}),
		InputTypeDefs: map[string]*ast.InputObjectDefinition{},
		TypeDefs:      map[string]*ast.ObjectDefinition{},
		EnumDefs:      map[string]*ast.EnumDefinition{},
	}
}

func (sb *SchemaBuilder) Build() *ast.Document {
	doc := ast.NewDocument(&ast.Document{})
	for _, t := range sb.InputTypeDefs {
		doc.Definitions = append(doc.Definitions, t)
	}
	for _, t := range sb.TypeDefs {
		doc.Definitions = append(doc.Definitions, t)
	}
	for _, t := range sb.EnumDefs {
		doc.Definitions = append(doc.Definitions, t)
	}
	doc.Definitions = append(doc.Definitions, sb.QueryType)
	return doc
}

type CreateTypeFunc func(t *desc.MessageDescriptor) (ast.Definition, string, error)

func CreateNamedType(name string) *ast.Named {
	return ast.NewNamed(&ast.Named{
		Name: CreateNameType(name),
	})
}

func CreateNameType(name string) *ast.Name {
	return ast.NewName(&ast.Name{Value: name})
}

func (sb *SchemaBuilder) CreateGraphqlType(t *desc.FieldDescriptor, CreateObjTypeFunc CreateTypeFunc) (ast.Type, string, error) {
	named := func(name string) (ast.Type, string, error) {
		namedNode := CreateNamedType(name)
		if t.GetLabel() == descriptor.FieldDescriptorProto_LABEL_REPEATED {
			return ast.NewList(&ast.List{
				Type: namedNode,
			}), name, nil
		}
		return namedNode, name, nil
	}
	switch t.GetType() {
	case descriptor.FieldDescriptorProto_TYPE_BOOL:
		return named("Boolean")
	case descriptor.FieldDescriptorProto_TYPE_DOUBLE:
		fallthrough
	case descriptor.FieldDescriptorProto_TYPE_FLOAT:
		return named("Float")
	case descriptor.FieldDescriptorProto_TYPE_INT64:
		fallthrough
	case descriptor.FieldDescriptorProto_TYPE_UINT64:
		fallthrough
	case descriptor.FieldDescriptorProto_TYPE_UINT32:
		fallthrough
	case descriptor.FieldDescriptorProto_TYPE_INT32:
		fallthrough
	case descriptor.FieldDescriptorProto_TYPE_SFIXED32:
		fallthrough
	case descriptor.FieldDescriptorProto_TYPE_SFIXED64:
		fallthrough
	case descriptor.FieldDescriptorProto_TYPE_SINT32:
		fallthrough
	case descriptor.FieldDescriptorProto_TYPE_SINT64:
		fallthrough
	case descriptor.FieldDescriptorProto_TYPE_FIXED64:
		fallthrough
	case descriptor.FieldDescriptorProto_TYPE_FIXED32:
		return named("Int")
	case descriptor.FieldDescriptorProto_TYPE_BYTES:
		fallthrough
	case descriptor.FieldDescriptorProto_TYPE_STRING:
		return named("String")
	case descriptor.FieldDescriptorProto_TYPE_GROUP:
		fallthrough
	case descriptor.FieldDescriptorProto_TYPE_MESSAGE:
		_, name, err := CreateObjTypeFunc(t.GetMessageType())
		if err != nil {
			return nil, "", err
		}
		return named(name)
	case descriptor.FieldDescriptorProto_TYPE_ENUM:
		enumName := sb.CreateEnumType(t.GetEnumType())
		return named(enumName)
	default:
		return nil, "", errors.New(fmt.Sprintf("Unable to translate protobuf type %s", t.GetType().String()))
	}
}

func GetMessageName(t desc.Descriptor) string {
	if parentMsg, ok := t.GetParent().(*desc.MessageDescriptor); ok {
		return GetMessageName(parentMsg) + "_" + t.GetName()
	}
	return t.GetName()
}

func (sb *SchemaBuilder) CreateOutputMessageType(t *desc.MessageDescriptor) (ast.Definition, string, error) {
	// special case for google.protobuf.WrapperValue types
	if substitutionName := TranslateGoogleProtobufWrapperTypes(t); substitutionName != "" {
		return nil, substitutionName, nil
	}
	typeName := GetMessageName(t)
	if objDef, ok := sb.TypeDefs[typeName]; ok {
		return objDef, typeName, nil
	}
	obj := ast.NewObjectDefinition(&ast.ObjectDefinition{})
	sb.TypeDefs[typeName] = obj
	obj.Name = CreateNameType(typeName)
	obj.Description = ast.NewStringValue(&ast.StringValue{Value: "Created from protobuf type " + t.GetFullyQualifiedName()})

	for _, field := range t.GetFields() {
		t, typeName, err := sb.CreateGraphqlType(field, sb.CreateOutputMessageType)
		if err != nil {
			return nil, "", err
		}
		newValDef := ast.NewFieldDefinition(&ast.FieldDefinition{
			Name: CreateNameType(field.GetName()),
		})

		if t != nil {
			newValDef.Type = t
		} else {
			// If the type t was not created and no err is thrown,
			// this is a case where a Message maps to a GraphQL primitive value
			// For example, google.protobuf.Timestamp is encoded in JSON as a String.
			newValDef.Type = CreateNamedType(typeName)
		}

		obj.Fields = append(obj.Fields, newValDef)
	}
	// we can not have empty graphql types, so we need to add a field definition which will not be used.
	// If queried, this field will always return false.
	if len(t.GetFields()) == 0 {
		obj.Fields = []*ast.FieldDefinition{
			ast.NewFieldDefinition(&ast.FieldDefinition{
				Name:        CreateNameType("_"),
				Description: ast.NewStringValue(&ast.StringValue{Value: "This GraphQL type was generated from an empty proto message. This empty field exists to keep the schema GraphQL spec compliant. If queried, this field will always return false."}),
				Type:        CreateNamedType("Boolean"),
			}),
		}
	}
	return obj, typeName, nil
}

func (sb *SchemaBuilder) CreateInputMessageType(inputType *desc.MessageDescriptor) (ast.Definition, string, error) {
	typeName := GetMessageName(inputType) + "Input"
	if def := sb.InputTypeDefs[typeName]; def != nil {
		return def, typeName, nil
	}
	inputObj := ast.NewInputObjectDefinition(&ast.InputObjectDefinition{})
	sb.InputTypeDefs[typeName] = inputObj
	inputObj.Description = ast.NewStringValue(&ast.StringValue{Value: "Created from protobuf type " + inputType.GetFullyQualifiedName()})
	inputObj.Name = CreateNameType(typeName)
	for _, field := range inputType.GetFields() {
		t, _, err := sb.CreateGraphqlType(field, sb.CreateInputMessageType)
		newInputValDef := ast.NewInputValueDefinition(&ast.InputValueDefinition{
			Name: CreateNameType(field.GetName()),
			Type: t,
		})
		if err != nil {
			return nil, "", err
		}
		inputObj.Fields = append(inputObj.Fields, newInputValDef)
	}
	// we can not have empty graphql types, so we need to add a field definition which will not be used.
	// If queried, this field will always return false.
	if len(inputType.GetFields()) == 0 {
		inputObj.Fields = []*ast.InputValueDefinition{
			ast.NewInputValueDefinition(&ast.InputValueDefinition{
				Name:        CreateNameType("_"),
				Description: ast.NewStringValue(&ast.StringValue{Value: "This GraphQL type was generated from an empty proto message. This empty field exists to keep the schema GraphQL spec compliant. If queried, this field will always return false."}),
				Type:        CreateNamedType("Boolean"),
			}),
		}
	}
	return inputObj, typeName, nil
}

func (sb *SchemaBuilder) CreateEnumType(enumType *desc.EnumDescriptor) string {
	typeName := GetMessageName(enumType)
	if sb.EnumDefs[typeName] != nil {
		return typeName
	}
	enumDef := ast.NewEnumDefinition(&ast.EnumDefinition{})
	enumDef.Name = CreateNameType(typeName)
	sb.EnumDefs[typeName] = enumDef
	for _, enumVal := range enumType.GetValues() {
		enumDef.Values = append(enumDef.Values, ast.NewEnumValueDefinition(&ast.EnumValueDefinition{
			Name: CreateNameType(enumVal.GetName()),
		}))
	}
	return typeName
}

func (sb *SchemaBuilder) AddQueryField(name string, inputType *desc.MessageDescriptor, inputTypeName string, outputType *desc.MessageDescriptor, resolverName string) {
	fieldDef := ast.NewFieldDefinition(&ast.FieldDefinition{})
	fieldDef.Name = CreateNameType(name)
	fieldDef.Type = CreateNamedType(outputType.GetName())
	fieldDef.Arguments = []*ast.InputValueDefinition{
		ast.NewInputValueDefinition(&ast.InputValueDefinition{
			Name: CreateNameType(inputType.GetName()),
			Type: CreateNamedType(inputTypeName),
		}),
	}
	fieldDef.Directives = append(fieldDef.Directives, ast.NewDirective(&ast.Directive{
		Name: CreateNameType("resolve"),
		Arguments: []*ast.Argument{ast.NewArgument(&ast.Argument{
			Name:  CreateNameType("name"),
			Value: ast.NewStringValue(&ast.StringValue{Value: resolverName}),
		})},
	}))

	sb.QueryType.Fields = append(sb.QueryType.Fields, fieldDef)
}

// Generates "OutgoingJsonBody" part of the resolver configuration for the gRPC type
// once it's translated to a graphql type
func GenerateOutgoingJsonBodyForInputType(inputType *desc.MessageDescriptor, argsPath string) *structpb.Value {
	val := &structpb.Struct{
		Fields: map[string]*structpb.Value{},
	}
	for _, f := range inputType.GetFields() {
		path := argsPath

		var newVal *structpb.Value
		path += "." + f.GetName()
		if f.GetLabel() == descriptor.FieldDescriptorProto_LABEL_REPEATED {
			path += "[*]"
		}

		newVal = &structpb.Value{
			Kind: &structpb.Value_StringValue{
				StringValue: path + "}",
			},
		}

		val.Fields[f.GetName()] = newVal
	}
	return &structpb.Value{
		Kind: &structpb.Value_StructValue{
			StructValue: val,
		},
	}
}
