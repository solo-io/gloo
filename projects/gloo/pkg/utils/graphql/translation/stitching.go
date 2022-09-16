package translation

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"os/exec"
	"strings"

	"k8s.io/utils/lru"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/graphql-go/graphql/language/ast"
	"github.com/rotisserie/eris"
	gloov3 "github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/config/core/v3"
	gloov2 "github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/extensions/filters/http/graphql/v2"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	gloov1beta1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/graphql/v1beta1"
	"github.com/solo-io/gloo/projects/gloo/pkg/utils"
	v2 "github.com/solo-io/solo-apis/pkg/api/gloo.solo.io/external/envoy/extensions/filters/http/graphql/v2"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	enterprisev1 "github.com/solo-io/solo-projects/projects/gloo/pkg/api/enterprise/graphql/v1"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/utils/graphql/printer"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/utils/graphql/types"
)

// The resolver info per subschema name
type ResolverInfoPerSubschema map[string]*gloov2.ResolverInfo

func getStitchingInfo(schema *gloov1beta1.StitchedSchema, graphqlApis types.GraphQLApiList, schemaCache *lru.Cache, stitchingInfoCache *lru.Cache) (*enterprisev1.GraphQLToolsStitchingOutput, map[string]ResolverInfoPerSubschema, map[string]string, []*gloov1beta1.GraphQLApi, error) {
	schemas := &enterprisev1.GraphQLToolsStitchingInput{}
	// map of types -> map of subschema names -> resolver info
	argMap := map[string]ResolverInfoPerSubschema{}
	// Query field map holds a mapping of Query field -> subschema that the query field is from
	queryFieldMap := map[string]string{}

	var subschemaGraphqlApis []*gloov1beta1.GraphQLApi
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
		schemaDef, err := getGraphQlApiSchemaDefinition(gqlSchema, graphqlApis, schemaCache, stitchingInfoCache)
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
	stitchingInfoOut, err := processStitchingInfo(schemas, stitchingInfoCache)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	return stitchingInfoOut, argMap, queryFieldMap, subschemaGraphqlApis, nil
}

func createStitchingScriptSubschema(subschema *gloov1beta1.StitchedSchema_SubschemaConfig, subschemaName, schemaDef string, argMap map[string]ResolverInfoPerSubschema) *enterprisev1.GraphQLToolsStitchingInput_Schema {
	stitchingScriptSchema := &enterprisev1.GraphQLToolsStitchingInput_Schema{
		Name:            subschemaName,
		Schema:          schemaDef,
		TypeMergeConfig: map[string]*enterprisev1.GraphQLToolsStitchingInput_Schema_TypeMergeConfig{},
	}

	for typeName, mergeCfg := range subschema.GetTypeMerge() {
		stitchingScriptSchema.TypeMergeConfig[typeName] = &enterprisev1.GraphQLToolsStitchingInput_Schema_TypeMergeConfig{
			SelectionSet: mergeCfg.GetSelectionSet(),
			FieldName:    mergeCfg.GetQueryName(),
		}
		if _, ok := argMap[typeName]; !ok {
			argMap[typeName] = map[string]*gloov2.ResolverInfo{}
		}

		ri := &gloov2.ResolverInfo{
			FieldName: mergeCfg.GetQueryName(),
		}
		for setter, extraction := range mergeCfg.GetArgs() {
			ri.Args = append(ri.Args, &gloov2.ArgPath{
				ExtractionPath: strings.Split(extraction, "."),
				SetterPath:     strings.Split(setter, "."),
			})
		}
		argMap[typeName][subschemaName] = ri
	}
	return stitchingScriptSchema
}

// only used by apiserver, so cache is not available here as we don't have any way of persisting state in the apiserver.
func GetStitchedSchemaDefinition(stitchedSchema *gloov1beta1.StitchedSchema, gqlApis types.GraphQLApiList) (string, error) {
	stitchedSchemaOut, _, _, _, err := getStitchingInfo(stitchedSchema, gqlApis, nil, nil)
	if err != nil {
		return "", err
	}
	return stitchedSchemaOut.GetStitchedSchema(), nil
}

// A mock upstream list which will never return an error for the `Find` method.
type MockUpstreamsList struct{}
type MockArtifactsList struct{}

func (l *MockUpstreamsList) Find(namespace, name string) (*gloov1.Upstream, error) {
	return &gloov1.Upstream{
		Metadata: &core.Metadata{
			Name:      "fake-upstream",
			Namespace: "fake-namespace",
		},
	}, nil
}

func (l *MockArtifactsList) Find(namespace, name string) (*gloov1.Artifact, error) {
	return &gloov1.Artifact{
		Metadata: &core.Metadata{
			Name:      "fake-artifact",
			Namespace: "fake-namespace",
		},
	}, nil
}

// Gets only the schema definition without validating Upstreams, hence the use of the MockUpstreamList
func getGraphQlApiSchemaDefinition(graphQLApi *gloov1beta1.GraphQLApi, gqlApis types.GraphQLApiList, schemaCache, stitchingInfoCache *lru.Cache) (string, error) {
	v2ApiSchema, err := CreateGraphQlApi(
		CreateGraphQLApiParams{&MockArtifactsList{}, &MockUpstreamsList{}, gqlApis, graphQLApi, schemaCache, stitchingInfoCache},
	)
	if err != nil {
		return "", eris.Wrapf(err, "error getting schema definition for GraphQLApi %s.%s", graphQLApi.GetMetadata().GetNamespace(), graphQLApi.GetMetadata().GetName())
	}
	return v2ApiSchema.GetSchemaDefinition().GetInlineString(), nil
}

func processStitchingInfo(schemas *enterprisev1.GraphQLToolsStitchingInput, stitchingInfoCache *lru.Cache) (*enterprisev1.GraphQLToolsStitchingOutput, error) {

	schemasBytes, err := proto.Marshal(schemas)
	if err != nil {
		return nil, eris.Wrapf(err, "error marshaling to binary data")
	}
	// hash schemaBytes to get a unique key for the cache
	schemasBytesHash := sha256.Sum256(schemasBytes)
	// We cache the stitching info so that we don't have to recompute it every time
	// the key is the base64 encoded protobuf of the stitching info input
	// the value is the protobuf message of the stitching info output
	if stitchingInfoCache != nil {
		stitchingInfo, ok := stitchingInfoCache.Get(schemasBytesHash)
		if ok {
			return stitchingInfo.(*enterprisev1.GraphQLToolsStitchingOutput), nil
		}
	}

	stitchingPath := GetGraphqlJsRoot()
	cmd := exec.Command("node", stitchingPath+"stitching.js", base64.StdEncoding.EncodeToString(schemasBytes))

	protoDirPath, err := GetGraphqlProtoRoot()
	if err != nil {
		return nil, err
	}
	cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", GraphqlProtoRootEnvVar, protoDirPath))
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
	stitchingInfoOut := &enterprisev1.GraphQLToolsStitchingOutput{}
	err = proto.Unmarshal(decodedStdOutString, stitchingInfoOut)
	if err != nil {
		return nil, eris.Wrap(err, "unable to unmarshal graphql tools output to Go type")
	}
	if stitchingInfoCache != nil {
		stitchingInfoCache.Add(schemasBytesHash, stitchingInfoOut)
	}
	return stitchingInfoOut, nil
}

func generateSubschemaName(subschemaResourseRef *core.ResourceRef) string {
	return fmt.Sprintf("Graphqlschema-%s.%s", subschemaResourseRef.GetNamespace(), subschemaResourseRef.GetName())
}

var (
	stitchingExtensionName = "stitching_extension"
)

func addSubschemaNameResolverInfo(mergeTypes map[string]*gloov2.MergedTypeConfig, argMap map[string]ResolverInfoPerSubschema) {
	for typeName, cfg := range mergeTypes {
		if argMap[typeName] == nil || cfg == nil {
			continue
		}
		cfg.SubschemaNameToResolverInfo = map[string]*gloov2.ResolverInfo{}
		for subschemaName, subschemaMergeCfg := range argMap[typeName] {
			cfg.SubschemaNameToResolverInfo[subschemaName] = subschemaMergeCfg
		}
	}
}

func translateStitchedSchema(params CreateGraphQLApiParams, schema *gloov1beta1.StitchedSchema) (*gloov2.ExecutableSchema, error) {
	stitchingInfoOut, argMap, queryFieldMap, subschemaGqls, err := getStitchingInfo(schema, params.Graphqlapis, params.ProcessedSchemaCache, params.StitchingInfoCache)
	if err != nil {
		return nil, err
	}
	gatewaySchema, err := ParseGraphQLSchema(stitchingInfoOut.GetStitchedSchema())

	if err != nil {
		return nil, err
	}

	//resolutions generated here
	var resolutions []*gloov2.Resolution
	for _, def := range gatewaySchema.Definitions {
		if objDef, ok := def.(*ast.ObjectDefinition); ok && objDef.Name.Value == "Query" {
			for _, field := range objDef.Fields {
				resolver := &gloov2.StitchingResolver{
					SubschemaName: queryFieldMap[field.Name.Value],
				}
				r := &gloov2.Resolution{
					Matcher: &gloov2.QueryMatcher{
						Match: &gloov2.QueryMatcher_FieldMatcher_{
							FieldMatcher: &gloov2.QueryMatcher_FieldMatcher{
								Type:  "Query",
								Field: field.Name.Value,
							},
						},
					},
					Resolver: &gloov3.TypedExtensionConfig{
						Name:        "io.solo.graphql.resolver.stitching",
						TypedConfig: utils.MustMessageToAny(resolver),
					},
				}
				resolutions = append(resolutions, r)
			}
		}
	}
	// GraphQLToolsStitchingOutput uses the solo-apis version of the envoy apis. we need to convert the fields back to
	// the gloo version of the envoy apis here
	glooFieldNodesByType, err := toGlooFieldNodesByType(stitchingInfoOut.GetFieldNodesByType())
	if err != nil {
		return nil, eris.Wrap(err, "error converting FieldNodes")
	}
	glooFieldNodesByField, err := toGlooFieldNodesByField(stitchingInfoOut.GetFieldNodesByField())
	if err != nil {
		return nil, eris.Wrap(err, "error converting FieldNodeMap")
	}
	glooMergedTypes, err := toGlooMergedTypes(stitchingInfoOut.GetMergedTypes())
	if err != nil {
		return nil, eris.Wrap(err, "error converting MergedTypeConfig")
	}

	addSubschemaNameResolverInfo(glooMergedTypes, argMap)
	subschemaNameToExecutableSchema := map[string]*gloov2.StitchingInfo_SubschemaConfig{}
	for _, subschemaGql := range subschemaGqls {
		subschemaRef := subschemaGql.GetMetadata().Ref()
		subschemaParams := params
		subschemaParams.Graphqlapi = subschemaGql
		execSchema, err := CreateGraphQlApi(subschemaParams)
		if err != nil {
			return nil, eris.Wrapf(err, "unable to create configuration for subschema %s.%s", subschemaGql.GetMetadata().GetNamespace(), subschemaGql.GetMetadata().GetName())
		}
		subschemaNameToExecutableSchema[generateSubschemaName(subschemaRef)] = &gloov2.StitchingInfo_SubschemaConfig{
			ExecutableSchema: execSchema,
		}
	}

	stitchingExtension := &gloov2.StitchingInfo{
		FieldNodesByType:               glooFieldNodesByType,
		FieldNodesByField:              glooFieldNodesByField,
		MergedTypes:                    glooMergedTypes,
		SubschemaNameToSubschemaConfig: subschemaNameToExecutableSchema,
	}
	return &gloov2.ExecutableSchema{
		SchemaDefinition: &gloov3.DataSource{
			Specifier: &gloov3.DataSource_InlineString{
				InlineString: printer.PrettyPrintKubeString(stitchingInfoOut.GetStitchedSchema()),
			},
		},
		Extensions: map[string]*any.Any{
			stitchingExtensionName: utils.MustMessageToAny(stitchingExtension),
		},
		Executor: &gloov2.Executor{
			Executor: &gloov2.Executor_Local_{
				Local: &gloov2.Executor_Local{
					Resolutions:         resolutions,
					EnableIntrospection: true,
				},
			},
		},
	}, nil
}

// solo-apis to gloo conversion functions
func toGlooFieldNodesByType(fieldNodesByType map[string]*v2.FieldNodes) (map[string]*gloov2.FieldNodes, error) {
	glooFieldNodesByType := make(map[string]*gloov2.FieldNodes)
	for k, v := range fieldNodesByType {
		glooVal := &gloov2.FieldNodes{}
		err := types.ConvertGoProtoTypes(v, glooVal)
		if err != nil {
			return nil, err
		}
		glooFieldNodesByType[k] = glooVal
	}
	return glooFieldNodesByType, nil
}

func toGlooFieldNodesByField(fieldNodesByField map[string]*v2.FieldNodeMap) (map[string]*gloov2.FieldNodeMap, error) {
	glooFieldNodesByField := make(map[string]*gloov2.FieldNodeMap)
	for k, v := range fieldNodesByField {
		glooVal := &gloov2.FieldNodeMap{}
		err := types.ConvertGoProtoTypes(v, glooVal)
		if err != nil {
			return nil, err
		}
		glooFieldNodesByField[k] = glooVal
	}
	return glooFieldNodesByField, nil
}

func toGlooMergedTypes(mergedTypes map[string]*v2.MergedTypeConfig) (map[string]*gloov2.MergedTypeConfig, error) {
	glooMergedTypes := make(map[string]*gloov2.MergedTypeConfig)
	for k, v := range mergedTypes {
		glooVal := &gloov2.MergedTypeConfig{}
		err := types.ConvertGoProtoTypes(v, glooVal)
		if err != nil {
			return nil, err
		}
		glooMergedTypes[k] = glooVal
	}
	return glooMergedTypes, nil
}
