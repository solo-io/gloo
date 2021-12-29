package graphql

import (
	"fmt"
	"strconv"

	openapi "github.com/getkin/kin-openapi/openapi3"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
	"github.com/solo-io/go-utils/cliutils"
	. "github.com/solo-io/solo-projects/projects/discovery/pkg/fds/discoveries/openapi/graphqlschematranslation/types"
)

type GetArgsParams struct {
	RequestPayloadDef *DataDefinition
	Parameters        openapi.Parameters
	Operation         *Operation
	Data              *PreprocessingData
}

type CreateOrReuseComplexTypeParams struct {
	Def               *DataDefinition
	Operation         *Operation
	Iteration         int64              // Count of recursions used to create type
	IsInputObjectType bool               // Does not require isInputObjectType because unions must be composed of objects
	Data              *PreprocessingData // Data produced by preprocessing
}

type CreateOrReuseSimpleTypeParams struct {
	Def  *DataDefinition
	Data *PreprocessingData
}

type CreateFieldsParams struct {
	Def               *DataDefinition
	Links             map[string]*openapi.Link
	Operation         *Operation
	Iteration         int64
	IsInputObjectType bool
	Data              *PreprocessingData
}

type LinkOpRefToOpIdParams struct {
	Links     map[string]*openapi.Link
	LinkKey   string
	Operation Operation
	data      *PreprocessingData
}

var JsonScalar = graphql.NewScalar(graphql.ScalarConfig{
	Name:         "JSON",
	Description:  "The `JSON` scalar type represents JSON values as specified by [ECMA-404](http://www.ecma-international.org/publications/files/ECMA-ST/ECMA-404.pdf).",
	Serialize:    Identity,
	ParseValue:   Identity,
	ParseLiteral: ParseLiteral,
})

func (t *OasToGqlTranslator) GetGraphQlType(params CreateOrReuseComplexTypeParams) graphql.Type {
	def := params.Def
	name := def.GraphQLTypeName
	if params.IsInputObjectType {
		name = def.GraphQLInputObjectTypeName
	}

	if params.Iteration >= 50 {
		t.handleWarningf(EXCESSIVE_NESTING, "The child type will have type JSON", Location{Operation: params.Operation}, "Type %s has excessive nesting", name)
		def.GraphQlType = JsonScalar
		return def.GraphQlType
	}

	switch def.TargetGraphQLType {
	case "object":
		return t.CreateOrReuseOt(params)
	case "union":
		return t.CreateOrReuseUnion(params)
	case "list":
		return t.CreateOrReuseList(params)
	case "enum":
		return t.CreateOrReuseEnum(CreateOrReuseSimpleTypeParams{
			Def:  params.Def,
			Data: params.Data,
		})
	default:
		return t.GetScalarType(CreateOrReuseComplexTypeParams{
			Def:  params.Def,
			Data: params.Data,
		})
	}
}

func (t *OasToGqlTranslator) GetScalarType(params CreateOrReuseComplexTypeParams) *graphql.Scalar {
	def := params.Def
	switch def.TargetGraphQLType {
	case "id":
		def.GraphQlType = graphql.ID
	case "string":
		def.GraphQlType = graphql.String
	case "integer":
		def.GraphQlType = graphql.Int
	case "int64":
		def.GraphQlType = &graphql.Scalar{}
	case "number":
		fallthrough
	case "float":
		def.GraphQlType = graphql.Float
	case "boolean":
		def.GraphQlType = graphql.Boolean
	case "json":
		def.GraphQlType = JsonScalar
	}

	return def.GraphQlType.(*graphql.Scalar)
}

func ParseLiteral(ast.Value) interface{} {
	return ""
}

func Identity(val interface{}) interface{} {
	return val
}

func (t *OasToGqlTranslator) CreateOrReuseEnum(params CreateOrReuseSimpleTypeParams) *graphql.Enum {
	def := params.Def
	if def.GraphQlType != nil {
		return def.GraphQlType.(*graphql.Enum)
	} else {
		values := graphql.EnumValueConfigMap{}
		for _, e := range def.Schema.Enum {
			sanitized := fmt.Sprintf("%s", e)
			values[sanitized] = &graphql.EnumValueConfig{
				Value: e,
			}
		}

		def.GraphQlType = graphql.NewEnum(graphql.EnumConfig{
			Name:   def.GraphQLTypeName,
			Values: values,
		})
		return def.GraphQlType.(*graphql.Enum)
	}
}

func (t *OasToGqlTranslator) CreateOrReuseList(params CreateOrReuseComplexTypeParams) *graphql.List {
	def := params.Def

	name := def.GraphQLInputObjectTypeName
	if !params.IsInputObjectType {
		name = def.GraphQLTypeName
	}
	if !params.IsInputObjectType && def.GraphQlType != nil {
		return def.GraphQlType.(*graphql.List)
	} else if params.IsInputObjectType && def.GraphQLInputObjectType != nil {
		return def.GraphQLInputObjectType.(*graphql.List)
	}

	itemDef := def.SubDefinitions.ListType
	itemsName := itemDef.GraphQLTypeName
	itemsType := t.GetGraphQlType(CreateOrReuseComplexTypeParams{
		Def:               itemDef,
		Data:              params.Data,
		Operation:         params.Operation,
		Iteration:         params.Iteration + 1,
		IsInputObjectType: params.IsInputObjectType,
	})
	if itemsType != nil {
		listObjectType := graphql.NewList(itemsType)
		if !params.IsInputObjectType {
			def.GraphQlType = listObjectType
		} else {
			def.GraphQLInputObjectType = listObjectType
		}
		return listObjectType
	}
	t.handleWarningf(INVALID_LIST_TYPE, "", Location{Operation: params.Operation}, "Cannot create list item object type %s in list %s", itemsName, name)
	return nil
}

func (t *OasToGqlTranslator) CreateOrReuseUnion(params CreateOrReuseComplexTypeParams) *graphql.Union {
	def := params.Def
	if def.GraphQlType != nil {
		return def.GraphQlType.(*graphql.Union)
	}
	schema := def.Schema
	description := schema.Description
	memberTypeDefinitions := def.SubDefinitions.UnionType

	var types []*graphql.Object
	for _, memberTypeDef := range memberTypeDefinitions {
		types = append(types, t.GetGraphQlType(CreateOrReuseComplexTypeParams{
			Def:               memberTypeDef,
			Operation:         params.Operation,
			Data:              params.Data,
			Iteration:         params.Iteration + 1,
			IsInputObjectType: false,
		}).(*graphql.Object))
	}

	resolveTypeFunc := func(p graphql.ResolveTypeParams) *graphql.Object {
		valMap := p.Value.(map[string]interface{})
		// remove custom _openAPIToGraphQL property used to pass data
		delete(valMap, "_openAPIToGraphQL")

		/**
		 * Find appropriate member type
		 *
		 * currently, the check is performed by only checking the property
		 * names. In the future, we should also check the types of those
		 * properties.
		 *
		 *  there is a chance a that an intended member type cannot be
		 * identified if, for whatever reason, the return data is a superset
		 * of the fields specified in the OAS
		 */
		for _, t := range types {
			typeFields := t.Fields()
			every := true
			if len(valMap) <= len(typeFields) {
				for fieldName := range typeFields {
					_, ok := t.Fields()[fieldName]
					if !ok {
						every = false
						break
					}
				}
			}
			if every {
				return t
			}
		}
		return nil
	}

	// check ambiguous member types
	def.GraphQlType = graphql.NewUnion(graphql.UnionConfig{
		Name:        def.GraphQLTypeName,
		Types:       types,
		Description: description,
		ResolveType: resolveTypeFunc,
	})

	return def.GraphQlType.(*graphql.Union)
}

func (t *OasToGqlTranslator) CreateOrReuseOt(params CreateOrReuseComplexTypeParams) graphql.Type {
	// try to reuse a prexisting (input) object type
	// Case: query - reuse object type
	def := params.Def
	if !params.IsInputObjectType {
		if def.GraphQlType != nil {
			return def.GraphQlType
		}
	} else {
		if def.GraphQLInputObjectType != nil {
			return def.GraphQLInputObjectType
		}
	}

	schema := def.Schema
	description := schema.Description

	// Case: query - create object type
	if !params.IsInputObjectType {
		createFieldsThunk := func() graphql.Fields {
			fields := t.CreateFields(CreateFieldsParams{
				Def:               def,
				Links:             def.Links,
				Operation:         params.Operation,
				Data:              params.Data,
				Iteration:         params.Iteration,
				IsInputObjectType: false,
			})
			return fields
		}
		obj := graphql.NewObject(graphql.ObjectConfig{
			Name:        def.GraphQLTypeName,
			Description: description,
			Fields:      graphql.FieldsThunk(createFieldsThunk),
		})
		def.GraphQlType = graphql.Type(obj)
		return def.GraphQlType
	} else /*Case - mutation, create input object type */ {
		def.GraphQLInputObjectType = graphql.NewInputObject(graphql.InputObjectConfig{
			Name:        def.GraphQLInputObjectTypeName,
			Description: description,
			Fields: t.CreateInputFields(CreateFieldsParams{
				Def:               def,
				Links:             map[string]*openapi.Link{},
				Operation:         params.Operation,
				Iteration:         params.Iteration,
				IsInputObjectType: true,
				Data:              params.Data,
			}),
		})

		return def.GraphQLInputObjectType
	}
}

func (t *OasToGqlTranslator) CreateInputFields(params CreateFieldsParams) graphql.InputObjectConfigFieldMap {
	fields := graphql.InputObjectConfigFieldMap{}

	def := params.Def
	fieldTypeDefinitions := def.SubDefinitions.InputObjectType
	for fieldTypeKey, fieldTypeDefinition := range fieldTypeDefinitions {
		fieldSchema := fieldTypeDefinition.Schema
		objectType := t.GetGraphQlType(CreateOrReuseComplexTypeParams{
			Def:               fieldTypeDefinition,
			Operation:         params.Operation,
			Data:              params.Data,
			Iteration:         params.Iteration + 1,
			IsInputObjectType: params.IsInputObjectType,
		})

		requiredProperty := cliutils.Contains(def.Required, fieldTypeKey) && !fieldTypeDefinition.Schema.Nullable

		if objectType != nil {
			saneFieldTypeKey := sanitizeString(fieldTypeKey, CaseStyle_camelCase)
			sanePropName := StoreSaneName(
				saneFieldTypeKey,
				fieldTypeKey,
				params.Data.SaneMap)
			fieldType := graphql.Output(objectType)
			if requiredProperty {
				fieldType = graphql.Output(&graphql.NonNull{
					OfType: objectType,
				})
			}
			fields[sanePropName] = &graphql.InputObjectFieldConfig{
				Type:        fieldType,
				Description: fieldSchema.Description,
			}
		} else {
			t.handleWarningf(CANNOT_GET_FIELD_TYPE, "", Location{Operation: params.Operation}, "Cannot obtain GraphQL type for field '%s'", fieldTypeKey)
		}
	}

	return fields
}

func (t *OasToGqlTranslator) CreateFields(params CreateFieldsParams) graphql.Fields {
	fields := graphql.Fields{}

	def := params.Def
	fieldTypeDefinitions := def.SubDefinitions.InputObjectType
	for fieldTypeKey, fieldTypeDefinition := range fieldTypeDefinitions {
		fieldSchema := fieldTypeDefinition.Schema
		objectType := t.GetGraphQlType(CreateOrReuseComplexTypeParams{
			Def:               fieldTypeDefinition,
			Operation:         params.Operation,
			Data:              params.Data,
			Iteration:         params.Iteration + 1,
			IsInputObjectType: params.IsInputObjectType,
		})

		requiredProperty := cliutils.Contains(def.Required, fieldTypeKey) && !fieldTypeDefinition.Schema.Nullable

		if objectType != nil {
			saneFieldTypeKey := sanitizeString(fieldTypeKey, CaseStyle_camelCase)
			sanePropName := StoreSaneName(
				saneFieldTypeKey,
				fieldTypeKey,
				params.Data.SaneMap)
			fieldType := graphql.Output(objectType)
			if requiredProperty {
				fieldType = graphql.Output(graphql.NewNonNull(objectType))
			}
			fields[sanePropName] = &graphql.Field{
				Type:        fieldType,
				Description: fieldSchema.Description,
			}
		} else {
			t.handleWarningf(CANNOT_GET_FIELD_TYPE, "", Location{Operation: params.Operation}, "Cannot obtain GraphQL type for field '%s'", fieldTypeKey)
		}
	}

	if len(params.Links) > 0 && !params.IsInputObjectType {
		for saneLinkKey := range params.Links {
			if _, ok := fields[saneLinkKey]; ok {
				t.handleWarningf(LINK_NAME_COLLISION, "", Location{Operation: params.Operation},
					"Cannot create link %s because parent object type already contains a field with the same sanitized name.",
					saneLinkKey)
			} else {
				link := params.Links[saneLinkKey]
				var (
					linkedOpId string
					found      bool
				)
				if link.OperationID != "" {
					linkedOpId = link.OperationID
					found = true
				} else if link.OperationRef != "" {
					linkedOpId, found = t.LinkOpRefToOpId(params.Links, saneLinkKey, params.Operation, params.Data)
				}
				/**
				 * linkedOpId may not be initialized because operationRef may lead to an
				 * operation object that does not have an operationId
				 */
				if linkedOp, ok := params.Data.Operations[linkedOpId]; ok && found {
					// Determine parameters provided via link
					argsFromLink := link.Parameters
					// Get arguments that are not provided by the linked operation
					var dynamicParams openapi.Parameters
					for _, param := range linkedOp.Parameters {
						if _, ok := argsFromLink[param.Value.Name]; !ok {
							dynamicParams = append(dynamicParams, param)
						}
					}

					// Get arguments for link
					args := t.GetArgs(GetArgsParams{
						Parameters: dynamicParams,
						Operation:  linkedOp,
						Data:       params.Data,
					})

					// Get response object type for link
					var resObjectType graphql.Type
					if linkedOp.ResponseDefinition != nil && linkedOp.ResponseDefinition.GraphQlType != nil {
						resObjectType = linkedOp.ResponseDefinition.GraphQlType
					} else {
						resObjectType = t.GetGraphQlType(CreateOrReuseComplexTypeParams{
							Def:               linkedOp.ResponseDefinition,
							Operation:         params.Operation,
							Data:              params.Data,
							Iteration:         params.Iteration + 1,
							IsInputObjectType: false,
						})
					}

					description := link.Description + "\n\nEquivalent to " + linkedOp.OperationString
					// Finally, add the object type to the fields (using sanitized field name)
					fields[saneLinkKey] = &graphql.Field{
						Type:        graphql.Output(resObjectType),
						Args:        args,
						Description: description,
						Resolve: graphql.FieldResolveFn(func(resolveParams graphql.ResolveParams) (interface{}, error) {
							return GetResolverParams{
								Operation:    linkedOp,
								ArgsFromLink: argsFromLink,
								PayloadName:  params.Def.GraphQLTypeName,
								ResponseName: saneLinkKey,
								Data:         params.Data,
							}, nil
						}),
					}
				} else {
					t.handleWarningf(UNRESOLVABLE_LINK, "", Location{Operation: params.Operation},
						"Cannot resolve target of link '%s'", saneLinkKey)
				}
			}
		}
	}
	return fields
}

func StoreSaneName(saneStr string, str string, mapping map[string]string) string {
	if val, ok := mapping[saneStr]; ok && val != str {
		appendix := 2
		strToTest := saneStr + strconv.Itoa(appendix)
		for {
			strToTest = saneStr + strconv.Itoa(appendix)
			_, ok := mapping[strToTest]
			if ok {
				appendix += 1
			} else {
				break
			}
		}
		return StoreSaneName(
			strToTest,
			str,
			mapping,
		)
	}
	mapping[saneStr] = str
	return saneStr
}
