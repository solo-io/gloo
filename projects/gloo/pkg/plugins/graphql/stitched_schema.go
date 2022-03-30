package graphql

import (
	"fmt"
	"strings"

	"github.com/solo-io/solo-projects/projects/gloo/utils/graphql/stitching"

	"github.com/golang/protobuf/ptypes/any"
	"github.com/graphql-go/graphql/language/ast"
	"github.com/graphql-go/graphql/language/parser"
	"github.com/rotisserie/eris"
	v3 "github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/config/core/v3"
	v2 "github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/extensions/filters/http/graphql/v2"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/graphql/v1alpha1"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins"
	"github.com/solo-io/gloo/projects/gloo/pkg/utils"
)

var (
	StitchingExtensionName = "stitching_extension"
)

// The resolver info per subschema name
type ResolverInfoPerSubschema map[string]*v2.ResolverInfo

func translateStitchedSchema(params plugins.RouteParams, schema *v1alpha1.StitchedSchema) (*v2.ExecutableSchema, error) {
	schemas := &v1alpha1.GraphQLToolsStitchingInput{}
	subschemaNameToExecutableSchema := map[string]*v2.StitchingInfo_SubschemaConfig{}
	// map of types -> map of subschema names -> resolver info
	argMap := map[string]ResolverInfoPerSubschema{}
	// Query field map holds a mapping of Query field -> subschema that the query field is from
	queryFieldMap := map[string]string{}
	for _, subschema := range schema.GetSubschemas() {

		gqlSchema, err := params.Snapshot.GraphqlApis.Find(subschema.GetNamespace(), subschema.GetName())
		if err != nil {
			ret := ""
			for _, s := range params.Snapshot.GraphqlApis {
				ret += " " + s.Metadata.Name
			}
			return nil, eris.Wrapf(err, "unable to find graphql api resource `%s` in namespace `%s`, here are all graphql apis found: %s", subschema.GetName(), subschema.GetNamespace(), ret)
		}

		subschemaName := GenerateSubschemaName(gqlSchema)
		execSchema, err := createGraphQlApi(params, gqlSchema)
		if err != nil {
			return nil, eris.Wrapf(err, "unable to create executable schema for subschema %s in namespace %s", subschema.GetName(), subschema.GetNamespace())
		}
		subschemaNameToExecutableSchema[subschemaName] = &v2.StitchingInfo_SubschemaConfig{
			ExecutableSchema: execSchema,
		}
		doc, err := parseGraphQLSchema(execSchema.GetSchemaDefinition().GetInlineString())
		if err != nil {
			// This should never happen as we're already parsing / validating the schema in `createGraphQlApi`
			return nil, err
		}
		// Populate queryFieldMap with the query fields that this subschema has
		for _, def := range doc.Definitions {
			// todo - when we do a second pass of type merging,
			// we want to also support Mutation field merging here.
			if obj, ok := def.(*ast.ObjectDefinition); ok && obj.Name.Value == "Query" {
				for _, field := range obj.Fields {
					queryFieldMap[field.Name.Value] = subschemaName
				}
			}
		}

		stitchingScriptSubschema := CreateStitchingScriptSubschema(subschema, subschemaName, execSchema.GetSchemaDefinition().GetInlineString(), argMap)
		schemas.Subschemas = append(schemas.Subschemas, stitchingScriptSubschema)
	}

	stitchingInfoOut, err := stitching.ProcessStitchingInfo(DefaultStitchingIndexFilePath, schemas)
	if err != nil {
		return nil, err
	}
	gatewaySchema, err := parser.Parse(parser.ParseParams{Source: stitchingInfoOut.GetStitchedSchema()})
	if err != nil {
		return nil, eris.Wrapf(err, "error parsing stitched schema %s", stitchingInfoOut.GetStitchedSchema())
	}
	var resolutions []*v2.Resolution
	for _, def := range gatewaySchema.Definitions {
		if objDef, ok := def.(*ast.ObjectDefinition); ok && objDef.Name.Value == "Query" {
			for _, field := range objDef.Fields {
				resolver := &v2.StitchingResolver{
					SubschemaName: queryFieldMap[field.Name.Value],
				}
				r := &v2.Resolution{
					Matcher: &v2.QueryMatcher{
						Match: &v2.QueryMatcher_FieldMatcher_{
							FieldMatcher: &v2.QueryMatcher_FieldMatcher{
								Type:  "Query",
								Field: field.Name.Value,
							},
						},
					},
					Resolver: &v3.TypedExtensionConfig{
						Name:        "io.solo.graphql.resolver.stitching",
						TypedConfig: utils.MustMessageToAny(resolver),
					},
				}
				resolutions = append(resolutions, r)
			}
		}
	}

	AddSubschemaNameResolverInfo(stitchingInfoOut.GetMergedTypes(), argMap)
	stitchingExtension := &v2.StitchingInfo{
		FieldNodesByType:               stitchingInfoOut.GetFieldNodesByType(),
		FieldNodesByField:              stitchingInfoOut.GetFieldNodesByField(),
		MergedTypes:                    stitchingInfoOut.GetMergedTypes(),
		SubschemaNameToSubschemaConfig: subschemaNameToExecutableSchema,
	}

	return &v2.ExecutableSchema{
		SchemaDefinition: &v3.DataSource{
			Specifier: &v3.DataSource_InlineString{
				InlineString: PrettyPrintKubeString(stitchingInfoOut.GetStitchedSchema()),
			},
		},
		Extensions: map[string]*any.Any{
			StitchingExtensionName: utils.MustMessageToAny(stitchingExtension),
		},
		Executor: &v2.Executor{
			Executor: &v2.Executor_Local_{
				Local: &v2.Executor_Local{
					Resolutions:         resolutions,
					EnableIntrospection: true,
				},
			},
		},
	}, nil
}

func AddSubschemaNameResolverInfo(mergeTypes map[string]*v2.MergedTypeConfig, argMap map[string]ResolverInfoPerSubschema) {
	for typeName, cfg := range mergeTypes {
		if argMap[typeName] == nil || cfg == nil {
			continue
		}
		cfg.SubschemaNameToResolverInfo = map[string]*v2.ResolverInfo{}
		for subschemaName, subschemaMergeCfg := range argMap[typeName] {
			cfg.SubschemaNameToResolverInfo[subschemaName] = subschemaMergeCfg
		}
	}

}

func CreateStitchingScriptSubschema(subschema *v1alpha1.StitchedSchema_SubschemaConfig, subschemaName, schemaDef string, argMap map[string]ResolverInfoPerSubschema) *v1alpha1.GraphQLToolsStitchingInput_Schema {
	stitchingScriptSchema := &v1alpha1.GraphQLToolsStitchingInput_Schema{
		Name:            subschemaName,
		Schema:          schemaDef,
		TypeMergeConfig: map[string]*v1alpha1.GraphQLToolsStitchingInput_Schema_TypeMergeConfig{},
	}

	for typeName, mergeCfg := range subschema.GetTypeMerge() {
		stitchingScriptSchema.TypeMergeConfig[typeName] = &v1alpha1.GraphQLToolsStitchingInput_Schema_TypeMergeConfig{
			SelectionSet: mergeCfg.GetSelectionSet(),
			FieldName:    mergeCfg.GetQueryName(),
		}
		if _, ok := argMap[typeName]; !ok {
			argMap[typeName] = map[string]*v2.ResolverInfo{}
		}

		ri := &v2.ResolverInfo{
			FieldName: mergeCfg.GetQueryName(),
		}
		for setter, extraction := range mergeCfg.GetArgs() {
			ri.Args = append(ri.Args, &v2.ArgPath{
				ExtractionPath: strings.Split(extraction, "."),
				SetterPath:     strings.Split(setter, "."),
			})
		}
		argMap[typeName][subschemaName] = ri
	}
	return stitchingScriptSchema
}

func GenerateSubschemaName(api *v1alpha1.GraphQLApi) string {
	return fmt.Sprintf("Graphqlschema-%s.%s", api.GetMetadata().GetNamespace(), api.GetMetadata().GetName())
}
