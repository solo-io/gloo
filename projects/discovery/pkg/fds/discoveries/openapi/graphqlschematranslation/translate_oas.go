package graphql

import (
	"fmt"
	"path"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/Masterminds/goutils"
	"github.com/gertd/go-pluralize"
	openapi "github.com/getkin/kin-openapi/openapi3"
	"github.com/go-openapi/inflect"
	"github.com/graphql-go/graphql"
	"github.com/iancoleman/strcase"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/graphql/v1alpha1"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	. "github.com/solo-io/solo-projects/projects/discovery/pkg/fds/discoveries/openapi/graphqlschematranslation/types"
)

type SchemaOptions struct {
	OperationIdFieldNames bool
	FillEmptyResponses    bool
	AddLimitArgument      bool
	GenericPayloadArgName bool
	SimpleNames           bool
	SingularNames         bool
}

type OasToGqlTranslator struct {
	Upstream *v1.Upstream
	Warnings []Warning
	Errors   []Error
}

type Error struct {
	Type       MitigationType
	Message    string
	Mitigation string
	Location   string
}

func (e Error) Error() string {
	return e.Message
}

type Warning struct {
	Type       MitigationType
	Message    string
	Mitigation string
	Location   string
	// Logger logger
}

type Location struct {
	// only provide one of the 3 for warning location
	Path      []string
	Operation *Operation
	Oas       *openapi.T
}

func (t *OasToGqlTranslator) handleWarningf(mitigationType MitigationType, mitigationAddendum string, location Location, messageFormat string, args ...interface{}) {
	t.handleWarning(mitigationType, mitigationAddendum, location, fmt.Sprintf(messageFormat, args...))
}

func (t *OasToGqlTranslator) handleWarning(mitigationType MitigationType, mitigationAddendum string, location Location, message string) {
	mitigation := Mitigations[mitigationType]
	mitigationMessage := mitigation
	if mitigationAddendum != "" {
		mitigationMessage += " " + mitigationAddendum
	}
	var locationString string
	if location.Path != nil {
		locationString = strings.Join(location.Path, ">")
	} else if location.Operation != nil {
		locationString = "Operation " + location.Operation.OperationId
	} else if location.Oas != nil && location.Oas.Info != nil {
		locationString = "OpenApiSpec " + location.Oas.Info.Title
	}
	t.Warnings = append(t.Warnings, Warning{
		Type:       mitigationType,
		Mitigation: mitigationMessage,
		Message:    message,
		Location:   locationString,
	})
}

func (t *OasToGqlTranslator) createErrorf(mitigationType MitigationType, mitigationAddendum string, location Location, messageFormat string, args ...interface{}) Error {
	return t.createError(mitigationType, mitigationAddendum, location, fmt.Sprintf(messageFormat, args...))
}

func (t *OasToGqlTranslator) createError(mitigationType MitigationType, mitigationAddendum string, location Location, message string) Error {
	mitigation := Mitigations[mitigationType]
	mitigationMessage := mitigation
	if mitigationAddendum != "" {
		mitigationMessage += " " + mitigationAddendum
	}
	var locationString string
	if location.Path != nil {
		locationString = strings.Join(location.Path, ">")
	} else if location.Operation != nil {
		locationString = "Operation " + location.Operation.OperationId
	} else if location.Oas != nil && location.Oas.Info != nil {
		locationString = "OpenApiSpec " + location.Oas.Info.Title
	}
	e := Error{
		Type:       mitigationType,
		Mitigation: mitigationMessage,
		Message:    message,
		Location:   locationString,
	}
	t.Errors = append(t.Errors, e)
	return e
}

func NewOasToGqlTranslator(upstream *v1.Upstream) *OasToGqlTranslator {
	return &OasToGqlTranslator{
		Upstream: upstream,
	}
}

func (t *OasToGqlTranslator) CreateGraphqlSchema(oass []*openapi.T) (*graphql.Schema, []*v1alpha1.Resolution) {
	return t.TranslateOpenApiToGraphQL(oass, SchemaOptions{})
}

func (t *OasToGqlTranslator) TranslateOpenApiToGraphQL(oass []*openapi.T, options SchemaOptions) (*graphql.Schema, []*v1alpha1.Resolution) {
	data := t.PreprocessOas(oass, options)

	queryFields := map[string]*graphql.Field{}
	mutationFields := map[string]*graphql.Field{}
	//todo-  subscription fields, authQueryfields, authMutationFields, authSubscriptionFields
	for operationId, operation := range data.Operations {
		field := t.GetFieldForOperation(
			operation,
			data,
		)

		saneOperationId := sanitizeString(operationId, CaseStyle_camelCase)
		if operation.OperationType == GraphQlOperationType_Query {
			fieldName := goutils.Uncapitalize(operation.ResponseDefinition.GraphQLTypeName)
			fieldName = StoreSaneName(saneOperationId, operationId, data.SaneMap)
			queryFields[fieldName] = field
		} else {
			saneFieldName := StoreSaneName(saneOperationId,
				operationId,
				data.SaneMap,
			)
			if _, ok := mutationFields[saneFieldName]; ok {
			} else {
				mutationFields[saneFieldName] = field
			}
		}
	}

	schemaConfig := graphql.SchemaConfig{
		Query: graphql.NewObject(graphql.ObjectConfig{
			Name: "Query",
			Fields: graphql.Fields(map[string]*graphql.Field{
				"DoNotUse": {
					Name:              "DoNotUse",
					Type:              graphql.Output(graphql.Boolean),
					Args:              nil,
					Resolve:           nil,
					DeprecationReason: "Do Not query this type, it is a placeholder",
					Description:       "DO NOT USE",
				},
			}),
		}),
	}
	if len(queryFields) > 0 {
		schemaConfig.Query = graphql.NewObject(graphql.ObjectConfig{
			Name:   "Query",
			Fields: graphql.Fields(queryFields),
		})
	}
	if len(mutationFields) > 0 {
		schemaConfig.Mutation = graphql.NewObject(graphql.ObjectConfig{
			Name:   "Mutation",
			Fields: graphql.Fields(mutationFields),
		})
	}
	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		t.createErrorf(GRAPHQL_SCHEMA_CREATION_ERR, "", Location{Oas: oass[0]}, "Unable to create schema with schema config: %+v", schemaConfig)
		return nil, nil
	}

	resolvers := t.CreateResolversForSchema(schema)
	return &schema, resolvers
}

func (t *OasToGqlTranslator) CreateResolversForSchema(schema graphql.Schema) []*v1alpha1.Resolution {
	resolutions := map[string]*v1alpha1.Resolution{}
	t.CreateResolverForField("Query", schema.QueryType(), resolutions)
	t.CreateResolverForField("Mutation", schema.MutationType(), resolutions)

	var r []*v1alpha1.Resolution
	for _, resol := range resolutions {
		r = append(r, resol)
	}
	return r
}

func (o *OasToGqlTranslator) GetFieldForOperation(operation *Operation, data *PreprocessingData) *graphql.Field {
	t := o.GetGraphQlType(CreateOrReuseComplexTypeParams{
		Def:               operation.ResponseDefinition,
		Data:              data,
		Operation:         operation,
		IsInputObjectType: false,
		Iteration:         0,
	})
	t = t.(graphql.Output)

	var payloadSchemaName string
	if operation.PayloadDefinition != nil {
		payloadSchemaName = operation.PayloadDefinition.GraphQLInputObjectTypeName
	}
	args := o.GetArgs(
		GetArgsParams{
			RequestPayloadDef: operation.PayloadDefinition,
			Parameters:        operation.Parameters,
			Operation:         operation,
			Data:              data,
		})

	//todo- support Subscription

	return &graphql.Field{
		Name:        payloadSchemaName,
		Description: operation.Description,
		Type:        t,
		Args:        args,
		Resolve: graphql.FieldResolveFn(func(p graphql.ResolveParams) (interface{}, error) {
			return GetResolverParams{
				Operation: operation,
				Data:      data,
			}, nil
		}),
	}
}

type FieldResolverWrapper struct {
	Resolvers []*v1alpha1.RESTResolver
	*graphql.Field
}

func (t *OasToGqlTranslator) GetArgs(params GetArgsParams) graphql.FieldConfigArgument {
	args := graphql.FieldConfigArgument{}
	parameters := params.Parameters
	if parameters == nil {
		return nil
	}
	for _, param := range parameters {
		var schema *openapi.Schema
		if param.Value != nil {
			p := param.Value
			if p.Schema != nil {
				schema = p.Schema.Value
			} else if p.Content != nil {
				jsonContentType := p.Content.Get("application/json")
				if jsonContentType == nil {
					t.handleWarningf(NON_APPLICATION_JSON_SCHEMA, "", Location{Operation: params.Operation},
						"The operation %s contains a parameter %s that has a content property but no schemas in application/json format. The parameter will not be created",
						params.Operation.OperationString, p.Name)
					return nil
				}
				schema = jsonContentType.Schema.Value
			} else {
				t.handleWarningf(INVALID_OAS, "", Location{Operation: params.Operation},
					"The operation %s contains a parameter %s with no 'schema' or 'content' property.", params.Operation.OperationString, p.Name)
				return nil
			}
			names := Names{
				FromSchema: p.Name,
			}
			paramDef := t.CreateDataDef(
				names,
				schema,
				true,
				params.Data,
				nil,
				params.Operation.Oas)

			graphqlType := t.GetGraphQlType(CreateOrReuseComplexTypeParams{
				Def:               paramDef,
				Operation:         params.Operation,
				Data:              params.Data,
				Iteration:         0,
				IsInputObjectType: true,
			})

			saneName := sanitizeString(p.Name, CaseStyle_camelCase)

			hasDefault := p.Schema.Value.Default != nil
			paramRequired := p.Required && !hasDefault
			var argType graphql.Input
			argType = graphql.NewNonNull(graphqlType)
			if !paramRequired {
				argType = graphqlType
			}
			args[saneName] = &graphql.ArgumentConfig{
				Type:        argType,
				Description: p.Description,
			}
		}
	}

	//todo- add limit argument (gql-mesh:schema_builder.ts:1202)
	if params.RequestPayloadDef != nil {
		reqObjectType := t.GetGraphQlType(CreateOrReuseComplexTypeParams{
			Def:               params.RequestPayloadDef,
			Data:              params.Data,
			Operation:         params.Operation,
			IsInputObjectType: true,
		})
		saneName := goutils.Uncapitalize(params.RequestPayloadDef.GraphQLInputObjectTypeName)
		reqRequired := params.Operation.PayloadRequired
		var argType graphql.Input
		argType = graphql.NewNonNull(reqObjectType)
		if !reqRequired {
			argType = reqObjectType
		}
		args[saneName] = &graphql.ArgumentConfig{
			Type:        argType,
			Description: params.RequestPayloadDef.Schema.Description,
		}
	}
	return args
}

func FormatOperationString(method string, path string, oas *openapi.T) string {
	result := strings.ToUpper(method) + " " + path
	if oas.Info != nil && oas.Info.Title != "" {
		result = oas.Info.Title + " " + result
	}
	return result
}

func GetParameters(path string, method string, operation *openapi.Operation, item *openapi.PathItem, oas *openapi.T) openapi.Parameters {
	parameters := openapi.Parameters{}
	for _, param := range item.Parameters {
		parameters = append(parameters, param)
	}
	for _, param := range operation.Parameters {
		parameters = append(parameters, param)
	}
	return parameters
}

func GetRefName(ref string) string {
	fromRef := strings.Split(ref, "/")
	if len(fromRef) == 0 {
		return ""
	}
	return fromRef[len(fromRef)-1]
}

func GetResponseObject(operation *openapi.Operation, statusCode string, oas *openapi.T) (string, *openapi.Response) {
	// this check is probably not needed since we got a valid status code from a response
	if len(operation.Responses) > 0 {
		if responseObject := operation.Responses[statusCode].Value; responseObject != nil {
			if len(responseObject.Content) > 0 {
				if isJson := responseObject.Content.Get("application/json"); isJson != nil {
					return "application/json", responseObject
				} else {
					// pick first random content type that's not json
					for contentType, _ := range responseObject.Content {
						return contentType, responseObject
					}
				}
			}
		}
	}
	return "", nil
}

/**
 * Determines name to use for schema from previously determined schemaNames and
 * considering not reusing existing names.
 */
func GetSchemaName(names Names, usedNames map[string]bool, caseStyle CaseStyle) string {
	schemaName := ""

	// Case - name from reference
	if len(names.FromRef) > 0 {
		saneName := sanitizeString(names.FromRef, caseStyle)
		if !usedNames[saneName] {
			schemaName = names.FromRef
		}
	}
	// case - name from title property in schema
	if schemaName == "" && len(names.FromSchema) > 0 {
		saneName := sanitizeString(names.FromSchema, caseStyle)
		if !usedNames[saneName] {
			schemaName = names.FromSchema
		}
	}
	// case - name from path
	if schemaName == "" && len(names.FromPath) > 0 {
		saneName := sanitizeString(names.FromPath, caseStyle)
		if !usedNames[saneName] {
			schemaName = names.FromPath
		}
	}

	// case -- all names already used, create appropriate name
	if schemaName == "" {
		if names.FromRef != "" {
			schemaName = names.FromRef
		} else if names.FromSchema != "" {
			schemaName = names.FromSchema
		} else if names.FromPath != "" {
			schemaName = names.FromPath
		} else {
			schemaName = "PlaceholderName"
		}
		schemaName = sanitizeString(schemaName, caseStyle)
	}

	if usedNames[schemaName] {
		appendix := 2
		newName := schemaName + strconv.Itoa(appendix)
		for ; usedNames[newName]; appendix++ {
			newName = schemaName + strconv.Itoa(appendix)
		}
		schemaName = newName
	}
	return schemaName
}

func GetExistingDataDef(name string, schema *openapi.Schema, defs []*DataDefinition) *DataDefinition {
	for _, def := range defs {
		if name == def.PreferredName && reflect.DeepEqual(schema, def.Schema) {
			return def
		}
	}
	return nil
}

func GetPreferredName(names Names) string {
	if names.FromRef != "" {
		return sanitizeString(names.FromRef, CaseStyle_PascalCase)
	}
	if names.FromSchema != "" {
		return sanitizeString(names.FromSchema, CaseStyle_PascalCase)
	}
	if names.FromPath != "" {
		return sanitizeString(names.FromPath, CaseStyle_PascalCase)
	}
	return "PlaceholderName"
}

func (t *OasToGqlTranslator) addObjectPropertiesToDataDef(def *DataDefinition, schema *openapi.Schema, isInputObjectType bool, data *PreprocessingData, oas *openapi.T) {
	for _, requiredProperty := range schema.Required {
		def.Required = append(def.Required, requiredProperty)
	}

	for propertyKey, property := range schema.Properties {
		propSchemaName := propertyKey
		propSchema := property.Value
		propSchemaNames := Names{
			FromRef: propSchemaName,
		}

		_, ok := def.SubDefinitions.InputObjectType[propertyKey]
		if !ok {
			subDef := t.CreateDataDef(
				propSchemaNames,
				propSchema,
				isInputObjectType,
				data,
				nil,
				oas)
			def.SubDefinitions.InputObjectType[propertyKey] = subDef
		}
	}
}

func (t *OasToGqlTranslator) LinkOpRefToOpId(links map[string]*openapi.Link, linkKey string, operation *Operation, data *PreprocessingData) (string, bool) {
	link := links[linkKey]
	var linkLocation, linkRelativePathAndMethod string

	/**
	 * Example relative path: '#/paths/~12.0~1repositories~1{username}/get'
	 * Example absolute path: 'https://na2.gigantic-server.com/#/paths/~12.0~1repositories~1{username}/get'
	 * Extract relative path from path
	 */
	if opRef := link.OperationRef; opRef != "" {
		if opRef[:8] == "#/paths/" {
			linkRelativePathAndMethod = opRef
		} else {
			if firstPathIndex := strings.Index(opRef, "#/paths/"); firstPathIndex != -1 {
				lastPathIndex := strings.LastIndex(opRef, "#/paths/")
				if firstPathIndex != lastPathIndex {
					t.handleWarningf(AMBIGUOUS_LINK, "", Location{Operation: operation},
						"The link %s in operation %s contains an ambiguous operationRef %s, meaning it has multiple instances of the string '#/paths'",
						linkKey, operation.OperationString, opRef)
					return "", false
				}

				linkLocation = opRef[:firstPathIndex]
				linkRelativePathAndMethod = opRef[firstPathIndex:]
			} else {
				t.handleWarningf(UNRESOLVABLE_LINK, "", Location{Operation: operation},
					`The link %s in operation %s `+
						`does not contain a valid path in operationRef '%s', `+
						`meaning it does not contain a string '#/paths/'`, linkKey, operation.OperationString, opRef)
				return "", false
			}
		}
		/**
		 * NOTE: I wish we could extract the linkedOpId by matching the
		 * linkedOpObject with an operation in data and extracting the operationId
		 * there but that does not seem to be possible especiially because you
		 * need to know the operationId just to access the operations so what I
		 * have to do is reconstruct the operationId the same way preprocessing
		 * does it
		 */

		/**
		 * linkPath should be the path followed by the method
		 *
		 * Find the slash that divides the path from the method
		 */
		if len(linkRelativePathAndMethod) > 0 {
			var linkPath, linkMethod string
			pivotSlashIndex := strings.LastIndex(linkRelativePathAndMethod, "/")
			// Check if there are any '/' in the linkPath
			if pivotSlashIndex != -1 {
				// Check if there is a method at the end of the linkPath
				if pivotSlashIndex != len(linkRelativePathAndMethod)-1 {
					linkMethod = stringToHttpMethod(linkRelativePathAndMethod[pivotSlashIndex+1:])
					if linkMethod == "" {
						t.handleWarningf(UNRESOLVABLE_LINK, "", Location{Operation: operation},
							"The operationRef %s contains an invalid HTTP method %s", opRef, linkRelativePathAndMethod[pivotSlashIndex+1:])
						return "", false
					}
				} else {
					t.handleWarningf(UNRESOLVABLE_LINK, "", Location{Operation: operation},
						"The operationRef %s does not contain an HTTP method", opRef)
					return "", false
				}
				/**
				 * Get path
				 *
				 * Substring starts at index 8 and ends at pivotSlashIndex to exclude
				 * the '/'s at the ends of the path
				 *
				 */
				linkPath = linkRelativePathAndMethod[8:pivotSlashIndex]
				/**
				 * linkPath is currently a JSON Pointer
				 *
				 * Revert the escaped '/', represented by '~1', to form intended path
				 */
				linkPath = strings.ReplaceAll(linkPath, "~1", "/")
				// Find the right oas
				if linkLocation != "" {
					t.handleWarningf(UNRESOLVABLE_LINK, "", Location{Operation: operation},
						"Operation '%s' contains a link which an external link location, which is not currently supported. Skipping link.")
					return "", false
				}
				oas := operation.Oas
				if linkMethod != "" && linkPath != "" {
					var linkedOpId string
					if pathObj, linkPathInPaths := oas.Paths[linkPath]; linkPathInPaths {
						if op := pathObj.GetOperation(linkMethod); op != nil {
							linkedOpObject := op
							linkedOpId = linkedOpObject.OperationID

							if linkedOpId == "" {
								linkedOpId = GenerateOperationId(linkMethod, linkPath)
							}
							if _, ok := data.Operations[linkedOpId]; ok {
								return linkedOpId, true
							} else {
								t.handleWarningf(UNRESOLVABLE_LINK, "", Location{Operation: operation},
									"The link '%s' references an operation with operationId '%s' but no such operation exists."+
										"Note that the operationId may be autogenerated but regardless, the link could not be matched to an operation.", linkKey, linkedOpId)
								return "", false
							}
						} else {
							// Path and method could not be found
							t.handleWarningf(UNRESOLVABLE_LINK, "", Location{Operation: operation},
								"Cannot identify path and/or method, '%s' and '%s' respectively, from operationRef '%s' in link '%s'",
								linkPath, linkMethod, opRef, linkKey)
						}
					}
				}
			} else {
				// Cannot split relative path into path and method sections
				t.handleWarningf(UNRESOLVABLE_LINK, "", Location{Operation: operation},
					"Cannot extract path and/or method from operationRef '%s' in link '%s'", opRef, linkKey)
			}
		} else {
			// Cannot extract relative path from absolute path
			t.handleWarningf(UNRESOLVABLE_LINK, "", Location{Operation: operation},
				"Cannot extract path and/or method from operationRef '%s' in link '%s'", opRef, linkKey)
		}
	}
	return "", false
}

func (t *OasToGqlTranslator) ResolveLinkParameter(param string) string {
	// The only link parameter type we support currently is from the parent request's body

	// CASE: parameter is parent body
	if param == "$response.body" {
		return "{$parent}"
	} else if strings.HasPrefix(param, "$response.body#") {
		// CASE: parameter in parent body

		// e.g. param = $response.body#/Components/Pets
		// extract out the path (/Components/Pets), so we remove $response.body# here
		newVal := param[15:]
		var segments []string
		pathArr := strings.Split(newVal, "/")
		for _, segment := range pathArr {
			if segment == "" {
				continue
			}
			segments = append(segments, segment)
		}
		return fmt.Sprintf("{$parent.%s}", strings.Join(segments, "."))
	}
	return ""
}

func (t *OasToGqlTranslator) ExtractRequestDataFromParent(baseUrl string, operation *Operation, linkParam string, providerString string) *v1alpha1.RESTResolver {
	var targetParam openapi.Parameters
	for _, param := range operation.Parameters {
		if linkParam == param.Value.Name {
			targetParam = append(targetParam, param)
			break
		}
	}

	extendedUrl := path.Join(baseUrl, operation.Path)
	requestTemplate := ExtractRequestDataFrom(extendedUrl, operation, targetParam, providerString, nil)
	resolver := &v1alpha1.RESTResolver{
		UpstreamRef: &core.ResourceRef{
			Name:      t.Upstream.GetMetadata().GetName(),
			Namespace: t.Upstream.GetMetadata().GetNamespace(),
		},
		Request: requestTemplate,
	}
	return resolver
}

func GetSchemaTargetGraphQlType(schema *openapi.Schema, data *PreprocessingData) string {
	if schema.Type == "object" || len(schema.Properties) > 0 {
		if schema.AdditionalProperties != nil && schema.AdditionalProperties.Value != nil {
			return "json"
		}
		return "object"
	}

	if schema.Type == "array" || schema.Items != nil {
		return "list"
	}
	if len(schema.Enum) > 0 {
		return "enum"
	}

	if len(schema.Type) > 0 {
		if len(schema.Format) > 0 {
			if schema.Type == "integer" && schema.Format == "int64" {
				return "integer"

			} else if schema.Type == "string" && schema.Format == "uuid" {
				return "id"
			}
		}
	}
	return schema.Type
}

func InferResourceNameFromPath(Path string) string {
	parts := strings.Split(Path, "/")
	reTest := regexp.MustCompile(`{`)
	pathNoParams := ReduceStringArray(parts, func(path string, part string, i int) string {
		if !reTest.Match([]byte(part)) {
			if i < len(parts)-1 && (IsIdParam(parts[i+1]) || IsSingularParam(part, parts[i+1])) {
				// pluralize.singular?

				return path + inflect.Singularize(strings.Title(part))
			} else {
				return path + strings.Title(part)
			}
		} else {
			return path
		}
	}, "")
	return pathNoParams
}

func IsSingularParam(part string, nextPart string) bool {
	return "{"+pluralize.NewClient().Singular(part)+"}" == nextPart
}

func IsIdParam(s string) bool {
	idParam := regexp.MustCompile(`^{.*(id|name|key).*}$`)
	return idParam.Match([]byte(s))
}

type Reducer func(accumulator string, part string, index int) string

func ReduceStringArray(arr []string, reducer Reducer, initialVal string) string {
	accumulator := initialVal
	for i, elem := range arr {
		accumulator = reducer(accumulator, elem, i)
	}
	if len(arr) == 1 {
		return arr[0]
	}
	return accumulator
}

func GetRequestBodyObject(operation *openapi.Operation, oas *openapi.T) (string, *openapi.RequestBody) {
	reqBody := operation.RequestBody
	if reqBody == nil || reqBody.Value == nil {
		return "", nil
	}
	reqBodyObj := reqBody.Value

	content := reqBodyObj.Content
	if jsonContentType := content.Get("application/json"); jsonContentType != nil {
		return "application/json", reqBodyObj
	} else if formDataContentType := content.Get("application/x-www-form-urlencoded"); formDataContentType != nil {
		return "application/x-www-form-urlencoded", reqBodyObj
	} else {
		for contentType, mediaType := range content {
			if mediaType != nil {
				return contentType, reqBodyObj
			}
		}
	}

	return "", nil

}

type CaseStyle int64

const (
	CaseStyle_simple CaseStyle = iota
	CaseStyle_ALL_CAPS
	CaseStyle_camelCase
	CaseStyle_PascalCase
)

func sanitizeString(str string, caseStyle CaseStyle) string {
	if caseStyle == CaseStyle_simple {
		r := regexp.MustCompile("[^a-zA-Z0-9_]")
		return r.ReplaceAllString(str, "")
	}

	var removeUnsafeCharRegex *regexp.Regexp
	if caseStyle == CaseStyle_ALL_CAPS {
		// ALL_CAPS has underscores
		removeUnsafeCharRegex = regexp.MustCompile("[^a-zA-Z0-9_]")
	} else {
		removeUnsafeCharRegex = regexp.MustCompile("[^a-zA-Z0-9]")
	}

	operators := map[string]string{
		"<=": "less_than_or_equal_to",
		">=": "greater_than_or_equal_to",
		"<":  "less_than",
		">":  "greater_than",
		"=":  "equal_to",
	}

	for _, op := range operators {
		str = strings.ReplaceAll(str, op, operators[op])
	}

	splitStr := removeUnsafeCharRegex.Split(str, -1)
	sanitized := ReduceStringArray(splitStr, func(path string, part string, _ int) string {
		if caseStyle == CaseStyle_ALL_CAPS {
			return path + "_" + part
		} else {
			return path + goutils.Capitalize(part)
		}
	}, "")

	switch caseStyle {
	case CaseStyle_PascalCase:
		{
			sanitized = strcase.ToCamel(sanitized)
		}
	case CaseStyle_camelCase:
		{
			sanitized = strcase.ToLowerCamel(sanitized)
		}
	case CaseStyle_ALL_CAPS:
		{
			sanitized = strcase.ToScreamingSnake(sanitized)
		}
	}

	startWithNumberRe := regexp.MustCompile(`^[0-9]`)
	if len(sanitized) < 1 || startWithNumberRe.MatchString(sanitized) {
		sanitized = "_" + sanitized
	}
	return sanitized

}
