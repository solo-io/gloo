package graphql

import (
	"fmt"
	"path"
	"regexp"
	"strings"

	"github.com/solo-io/gloo/projects/gloo/pkg/translator"

	openapi "github.com/getkin/kin-openapi/openapi3"
	structpb "github.com/golang/protobuf/ptypes/struct"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
	. "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/graphql/v1beta1"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-projects/projects/discovery/pkg/fds/discoveries/openapi-graphql/graphqlschematranslation/types"
)

type GetResolverParams struct {
	Operation    *types.Operation
	ArgsFromLink map[string]interface{}
	PayloadName  string
	ResponseName string
	Data         *types.PreprocessingData
	BaseUrl      string
}

// During the schema building, we stored the operation metadata (GetResolverParams) from the openapi spec in the
// `ResolveFn` thunk. We now get use the result of calling the thunk to build the resolver for each
// field.
func (t *OasToGqlTranslator) CreateResolverForField(
	parentTypeName string,
	gqlObj *graphql.Object,
	gqlType *ast.ObjectDefinition,
	resolutions map[string]*Resolution,
	typeDefs map[string]*ast.ObjectDefinition) {
	if gqlType == nil {
		return
	}
	// Iterate over fields of this type
	for _, field := range gqlType.Fields {
		fieldName := field.Name.Value
		resolutionName := CreateBaseResolverName(t.Upstream.GetMetadata().Ref(), parentTypeName, fieldName)
		if _, ok := resolutions[resolutionName]; ok {
			/* This terminates the recursion where we have a type that has a field that is of the same type
			   so we are not infinitely creating resolvers for the same type. For example:
			   type Employee {
			     userManager: Employee
			   }
			   The resolution name guarantees resolver uniqueness for a type-field pair per schema.
			*/
			continue
		}
		fieldObj, ok := gqlObj.Fields()[fieldName]
		if !ok || fieldObj.Resolve == nil {
			continue
		}
		getResolverParams, _ := fieldObj.Resolve(graphql.ResolveParams{})
		gRParams := getResolverParams.(GetResolverParams)
		resolvers := t.GetResolver(gRParams)
		for _, r := range resolvers {
			field.Directives = append(field.Directives,
				ast.NewDirective(&ast.Directive{
					Name: ast.NewName(&ast.Name{Value: "resolve"}),
					Arguments: []*ast.Argument{
						ast.NewArgument(&ast.Argument{
							Name:  ast.NewName(&ast.Name{Value: "name"}),
							Value: ast.NewStringValue(&ast.StringValue{Value: resolutionName}),
						}),
					},
				},
				))
			resolutions[resolutionName] = &Resolution{
				Resolver: &Resolution_RestResolver{
					RestResolver: r,
				},
			}
		}
		if childField, ok := fieldObj.Type.(*graphql.Object); ok {
			fieldTypeName := fieldObj.Type.String()
			t.CreateResolverForField(fieldTypeName, childField, typeDefs[fieldTypeName], resolutions, typeDefs)
		}
	}
}

// Creates a unique name for a field resolver given the field name, upstream metadata ref, and type name
func CreateBaseResolverName(ref *core.ResourceRef, typeName, fieldName string) string {
	return translator.UpstreamToClusterName(ref) + "|" + typeName + "|" + fieldName
}

func (t *OasToGqlTranslator) GetResolver(params GetResolverParams) []*RESTResolver {
	var resolvers []*RESTResolver
	baseUrl := params.BaseUrl
	if baseUrl == "" {
		var err error
		baseUrl, err = t.GetBaseUrlPath(params.Operation)
		if err != nil {
			t.handleWarningf(INVALID_SERVER_URL, "", Location{Operation: params.Operation},
				"Invalid server url in operation, skipping building resolver for operation")
			return nil
		}
	}
	resolver := &RESTResolver{
		UpstreamRef: &core.ResourceRef{
			Name:      t.Upstream.GetMetadata().GetName(),
			Namespace: t.Upstream.GetMetadata().GetNamespace(),
		},
	}
	operation := params.Operation
	// arg doesn't exit
	if operation != nil {

		// handle arguments provided by links
		for paramName, value := range params.ArgsFromLink {
			// value is an interface, cast it to string
			valueStr := fmt.Sprintf("%s", value)
			/**
			 * see if the link parameter contains constants that are appended to the link parameter
			 *
			 * e.g. instead of:
			 * $response.body#/employerId
			 *
			 * it could be:
			 * abc_{$response.body#/employerId}
			 */
			var valueProvider string
			if matches := regexp.MustCompile(`{([^}]*)}`).FindAllStringSubmatch(valueStr, 1); len(matches) > 0 {
				t.handleWarningf(LINK_PARAM_TEMPLATE_PRESENT, "", Location{Operation: operation},
					"Templating values with link parameters is not currently supported. Using only the first extraction and ignoring all else.")
				// replace link parameters with appropriate values

				valueStr = matches[0][1]
			}
			valueProvider = t.ResolveLinkParameter(valueStr)
			if valueProvider == "" {
				t.handleWarningf(LINK_UNSUPPORTED_EXTRACTION, "", Location{Operation: operation}, "Link %s uses an unsupported extraction %s. only $response.body extractions are currently supported",
					paramName, valueStr)
				continue
			}

			resolvers = append(resolvers, t.ExtractRequestDataFromParent(baseUrl, operation, paramName, valueProvider))
		}

		// Swallowing error here as it will never happen. URLJoin only returns error when baseUrl is invalid, which we confirmed to be valid in GetBaseUrlPath.
		extendedUrl := path.Join(baseUrl, operation.Path)
		resolverForArgs := ExtractRequestDataFrom(extendedUrl, operation, operation.Parameters, "", params.ArgsFromLink)
		if resolverForArgs == nil {
			return resolvers
		}
		resolver.Request = resolverForArgs

		/**
		  Determine the possible payload type
		*/
		if params.Operation.PayloadDefinition != nil {
			if jsonVal := TraverseGraphqlSchema(params.Operation, params.Operation.PayloadDefinition.GraphQLInputObjectTypeName, params.Operation.PayloadDefinition.Schema, params.Data.SaneMap); jsonVal != nil {
				if resolver.Request == nil {
					resolver.Request = &RequestTemplate{}
				}
				resolver.Request.Body = jsonVal
			}
		}
	}
	resolvers = append(resolvers, resolver)
	return resolvers
}

func TraverseGraphqlSchema(operation *types.Operation, inputTypeName string, schema *openapi.Schema, saneMap map[string]string) *structpb.Value {

	inputSaneName := sanitizeString(inputTypeName, CaseStyle_camelCase)
	if schema.Type == "object" {
		ret := &structpb.Struct{
			Fields: map[string]*structpb.Value{},
		}
		for propName, _ := range schema.Properties {
			p := sanitizeString(propName, CaseStyle_camelCase)
			ret.Fields[p] = &structpb.Value{
				Kind: &structpb.Value_StringValue{
					StringValue: fmt.Sprintf("{$args.%s.%s}", inputSaneName, saneMap[p]),
				},
			}
		}
		return &structpb.Value{
			Kind: &structpb.Value_StructValue{
				StructValue: ret,
			},
		}
	} else if schema.Type == "array" {
		return &structpb.Value{Kind: &structpb.Value_StringValue{StringValue: fmt.Sprintf("{$args.%s}", inputSaneName)}}
	} else {
		return nil
	}
}

func ExtractRequestDataFrom(path string, operation *types.Operation, parameters openapi.Parameters, providerString string, argsFromLink map[string]interface{}) *RequestTemplate {
	method := operation.Method
	requestTemplate := &RequestTemplate{
		Headers: map[string]string{
			":path":   path,
			":method": method,
		},
	}
	if len(parameters) == 0 {
		return requestTemplate
	}
	var parametersNotFromLink openapi.Parameters
	for _, param := range parameters {
		if _, ok := argsFromLink[param.Value.Name]; !ok {
			parametersNotFromLink = append(parametersNotFromLink, param)
		}
	}

	if len(parametersNotFromLink) == 0 {
		// all parameters were provided by link and we should not create a separate resolver
		return nil
	}
	for _, param := range parameters {

		p := param.Value
		sanitizedParamName := sanitizeString(p.Name, CaseStyle_camelCase)
		if providerString == "" {
			providerString = fmt.Sprintf("{$args.%s}", sanitizedParamName)
		}
		switch p.In {
		case "path":
			// replace /pet/{petid} from openapispec to /pet/{$args.ARG_NAME} template string
			// for extraction
			// todo - support multiple name parameters here.
			pString := "{" + p.Name + "}"
			requestTemplate.Headers[":path"] = strings.ReplaceAll(path, pString, providerString)

		case "query":
			if requestTemplate.QueryParams == nil {
				requestTemplate.QueryParams = map[string]string{}
			}
			requestTemplate.QueryParams[p.Name] = providerString
		case "header":
			requestTemplate.Headers[p.Name] = providerString
		case "cookie":
			requestTemplate.Headers["cookie"] = providerString
		}
	}
	return requestTemplate
}
