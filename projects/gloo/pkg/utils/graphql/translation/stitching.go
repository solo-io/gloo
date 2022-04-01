package translation

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/graphql-go/graphql/language/ast"
	"github.com/rotisserie/eris"
	v3 "github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/config/core/v3"
	v2 "github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/extensions/filters/http/graphql/v2"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/graphql/v1alpha1"
	"github.com/solo-io/gloo/projects/gloo/pkg/utils"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/utils/graphql/printer"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/utils/graphql/types"
)

var (
	// Env var that has the path to the index.js file that runs the stitching code
	StitchingIndexFilePathEnvVar = "STITCHING_PATH"
	// Env var that has the path to the proto dependencies that the stitching index.js file requires
	StitchingProtoDependenciesPathEnvVar = "STITCHING_PROTO_DIR"

	DefaultStitchingIndexFilePath = "/usr/local/bin/js/index.js"
)

// The resolver info per subschema name
type ResolverInfoPerSubschema map[string]*v2.ResolverInfo

func getStitchingInfo(schema *v1alpha1.StitchedSchema, graphqlApis types.GraphQLApiList) (*v1alpha1.GraphQlToolsStitchingOutput, map[string]ResolverInfoPerSubschema, map[string]string, []*v1alpha1.GraphQLApi, error) {
	schemas := &v1alpha1.GraphQLToolsStitchingInput{}
	// map of types -> map of subschema names -> resolver info
	argMap := map[string]ResolverInfoPerSubschema{}
	// Query field map holds a mapping of Query field -> subschema that the query field is from
	queryFieldMap := map[string]string{}

	var subschemaGraphqlApis []*v1alpha1.GraphQLApi
	for _, subschema := range schema.GetSubschemas() {

		gqlSchema, err := graphqlApis.Find(subschema.GetNamespace(), subschema.GetName())
		if err != nil {
			ret := ""
			for _, s := range graphqlApis.AsResources() {
				ret += " " + s.GetMetadata().GetName()
			}
			return nil, nil, nil, nil, eris.Wrapf(err, "unable to find graphql api resource `%s` in namespace `%s`, here are all graphql apis found: %s", subschema.GetName(), subschema.GetNamespace(), ret)
		}
		subschemaGraphqlApis = append(subschemaGraphqlApis, gqlSchema)

		subschemaName := generateSubschemaName(gqlSchema.GetMetadata().Ref())
		schemaDef, err := getGraphQlApiSchemaDefinition(gqlSchema, graphqlApis)
		if err != nil {
			return nil, nil, nil, nil, err
		}

		doc, err := ParseGraphQLSchema(schemaDef)
		if err != nil {
			// This should never happen as we're already parsing / validating the schema in `translateSchema`
			return nil, nil, nil, nil, err
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

		stitchingScriptSubschema := createStitchingScriptSubschema(subschema, subschemaName, schemaDef, argMap)
		schemas.Subschemas = append(schemas.Subschemas, stitchingScriptSubschema)
	}
	stitchingInfoOut, err := processStitchingInfo(DefaultStitchingIndexFilePath, schemas)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	return stitchingInfoOut, argMap, queryFieldMap, subschemaGraphqlApis, nil
}

func createStitchingScriptSubschema(subschema *v1alpha1.StitchedSchema_SubschemaConfig, subschemaName, schemaDef string, argMap map[string]ResolverInfoPerSubschema) *v1alpha1.GraphQLToolsStitchingInput_Schema {
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

func GetStitchedSchemaDefinition(stitchedSchema *v1alpha1.StitchedSchema, gqlApis types.GraphQLApiList) (string, error) {
	stitchedSchemaOut, _, _, _, err := getStitchingInfo(stitchedSchema, gqlApis)
	if err != nil {
		return "", err
	}
	return stitchedSchemaOut.GetStitchedSchema(), nil
}

// A mock upstream list which will never return an error for the `Find` method.
type MockUpstreamsList struct{}

func (l *MockUpstreamsList) Find(namespace, name string) (*v1.Upstream, error) {
	return &v1.Upstream{
		Metadata: &core.Metadata{
			Name:      "fake-upstream",
			Namespace: "fake-namespace",
		},
	}, nil
}

// Gets only the schema definition without validating upstreams, hence the use of the MockUpstreamList
func getGraphQlApiSchemaDefinition(graphQLApi *v1alpha1.GraphQLApi, gqlApis types.GraphQLApiList) (string, error) {
	v2ApiSchema, err := CreateGraphQlApi(&MockUpstreamsList{}, gqlApis, graphQLApi)
	if err != nil {
		return "", eris.Wrapf(err, "error getting schema definition for GraphQLApi %s.%s", graphQLApi.GetMetadata().GetNamespace(), graphQLApi.GetMetadata().GetName())
	}
	return v2ApiSchema.GetSchemaDefinition().GetInlineString(), nil
}

func processStitchingInfo(pathToStitchingJsFile string, schemas *v1alpha1.GraphQLToolsStitchingInput) (*v1alpha1.GraphQlToolsStitchingOutput, error) {
	schemasBytes, err := proto.Marshal(schemas)
	if err != nil {
		return nil, eris.Wrapf(err, "error marshaling to binary data")
	}
	// This is the default path
	var stitchingPath = pathToStitchingJsFile
	// Used for local testing and unit/e2e tests
	if path := os.Getenv(StitchingIndexFilePathEnvVar); path != "" {
		stitchingPath = path
	}
	cmd := exec.Command("node", stitchingPath, base64.StdEncoding.EncodeToString(schemasBytes))
	// Set the environment variable STITCHING_PROTO_IMPORT_PATH for the node file to know where to import dependencies from.
	protoDirPath := "/usr/local/bin/js/proto/github.com/solo-io/solo-apis/api/gloo/gloo"
	if protoDir := os.Getenv(StitchingProtoDependenciesPathEnvVar); protoDir != "" {
		// JS needs the absolute path to the stitching proto dir, so we join path to current dir + path from repository root
		currentDir, err := os.Getwd()
		if err != nil {
			return nil, eris.Wrap(err, "unable to get current directory path for running stitching script")
		}
		protoDirPath = path.Join(currentDir, os.Getenv(StitchingProtoDependenciesPathEnvVar))
	}
	cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", StitchingProtoDependenciesPathEnvVar, protoDirPath))
	stdOutBuf := bytes.NewBufferString("")
	stdErrBuf := bytes.NewBufferString("")
	cmd.Stdout = stdOutBuf
	cmd.Stderr = stdErrBuf
	err = cmd.Run()
	stdOutString := stdOutBuf.String()
	if err != nil {
		return nil, eris.Wrapf(err, "error running stitching info generation, stdout: %s, \nstderr: %s", stdOutString, stdErrBuf.String())
	}
	if len(stdOutString) == 0 {
		return nil, eris.Errorf("error running stitching info generation, no stitching info generated, \nstderr: %s", stdErrBuf.String())
	}
	decodedStdOutString, err := base64.StdEncoding.DecodeString(stdOutBuf.String())
	if err != nil {
		return nil, eris.Wrapf(err, "error decoding %s from base64 protobuf", stdOutBuf.String())
	}
	stitchingInfoOut := &v1alpha1.GraphQlToolsStitchingOutput{}
	err = proto.Unmarshal(decodedStdOutString, stitchingInfoOut)
	if err != nil {
		return nil, eris.Wrap(err, "unable to unmarshal graphql tools output to Go type")
	}
	return stitchingInfoOut, nil
}

func generateSubschemaName(subschemaResourseRef *core.ResourceRef) string {
	return fmt.Sprintf("Graphqlschema-%s.%s", subschemaResourseRef.GetNamespace(), subschemaResourseRef.GetName())
}

var (
	stitchingExtensionName = "stitching_extension"
)

func addSubschemaNameResolverInfo(mergeTypes map[string]*v2.MergedTypeConfig, argMap map[string]ResolverInfoPerSubschema) {
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

func translateStitchedSchema(upstreams types.UpstreamList, graphqlapis types.GraphQLApiList, schema *v1alpha1.StitchedSchema) (*v2.ExecutableSchema, error) {
	stitchingInfoOut, argMap, queryFieldMap, subschemaGqls, err := getStitchingInfo(schema, graphqlapis)
	if err != nil {
		return nil, err
	}
	gatewaySchema, err := ParseGraphQLSchema(stitchingInfoOut.GetStitchedSchema())
	if err != nil {
		return nil, err
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

	addSubschemaNameResolverInfo(stitchingInfoOut.GetMergedTypes(), argMap)
	subschemaNameToExecutableSchema := map[string]*v2.StitchingInfo_SubschemaConfig{}
	for _, subschemaGql := range subschemaGqls {
		subschemaRef := subschemaGql.GetMetadata().Ref()
		execSchema, err := CreateGraphQlApi(upstreams, graphqlapis, subschemaGql)
		if err != nil {
			return nil, eris.Wrapf(err, "unable to create configuration for subschema %s.%s", subschemaGql.GetMetadata().GetNamespace(), subschemaGql.GetMetadata().GetName())
		}
		subschemaNameToExecutableSchema[generateSubschemaName(subschemaRef)] = &v2.StitchingInfo_SubschemaConfig{
			ExecutableSchema: execSchema,
		}
	}
	stitchingExtension := &v2.StitchingInfo{
		FieldNodesByType:               stitchingInfoOut.GetFieldNodesByType(),
		FieldNodesByField:              stitchingInfoOut.GetFieldNodesByField(),
		MergedTypes:                    stitchingInfoOut.GetMergedTypes(),
		SubschemaNameToSubschemaConfig: subschemaNameToExecutableSchema,
	}

	return &v2.ExecutableSchema{
		SchemaDefinition: &v3.DataSource{
			Specifier: &v3.DataSource_InlineString{
				InlineString: printer.PrettyPrintKubeString(stitchingInfoOut.GetStitchedSchema()),
			},
		},
		Extensions: map[string]*any.Any{
			stitchingExtensionName: utils.MustMessageToAny(stitchingExtension),
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
