package translation

import (
	"encoding/base64"
	"hash/fnv"
	"strings"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/golang/protobuf/ptypes/duration"
	"github.com/graphql-go/graphql/language/ast"
	"github.com/graphql-go/graphql/language/parser"
	"github.com/graphql-go/graphql/language/printer"
	"github.com/pkg/errors"
	"github.com/rotisserie/eris"
	v3 "github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/config/core/v3"
	v2 "github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/extensions/filters/http/graphql/v2"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/graphql/v1beta1"
	"github.com/solo-io/gloo/projects/gloo/pkg/translator"
	"github.com/solo-io/gloo/projects/gloo/pkg/utils"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/graphql/dot_notation"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/graphql/resolvers/grpc"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/graphql/resolvers/mock"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/graphql/resolvers/rest"
	jsonUtils "github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/graphql/resolvers/utils"
	resolver_utils "github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/graphql/resolvers/utils"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/utils/graphql/directives"
	printer2 "github.com/solo-io/solo-projects/projects/gloo/pkg/utils/graphql/printer"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/utils/graphql/types"
	"k8s.io/utils/lru"
)

type CreateGraphQLApiParams struct {
	Artifacts   types.ArtifactList
	Upstreams   types.UpstreamList
	Graphqlapis types.GraphQLApiList
	Graphqlapi  *v1beta1.GraphQLApi

	ProcessedSchemaCache *lru.Cache
	StitchingInfoCache   *lru.Cache
}

func CreateGraphQlApi(params CreateGraphQLApiParams) (*v2.ExecutableSchema, error) {
	executableSchema, err := translateSchema(params)
	if err != nil {
		return nil, err
	}
	executableSchema.LogRequestResponseInfo = params.Graphqlapi.GetOptions().GetLogSensitiveInfo()
	return executableSchema, nil
}

func translateSchema(params CreateGraphQLApiParams) (*v2.ExecutableSchema, error) {
	switch schema := params.Graphqlapi.GetSchema().(type) {
	case *v1beta1.GraphQLApi_StitchedSchema:
		{
			return translateStitchedSchema(params, schema.StitchedSchema)
		}
	case *v1beta1.GraphQLApi_ExecutableSchema:
		{
			return translateExecutableSchema(params)
		}
	default:
		{
			return nil, eris.Errorf("unknown schema type %T", params.Graphqlapi.GetSchema())
		}
	}
}
func glooToEnvoyTranslation(namedExtractions map[string]string) (map[string]*v2.Executor_Remote_Extraction, error) {
	output := make(map[string]*v2.Executor_Remote_Extraction)
	for key, val := range namedExtractions {
		submatches := resolver_utils.ProviderTemplateRegex.FindAllStringSubmatch(val, -1)
		out := ""
		if len(submatches) > 1 {
			return nil, eris.Errorf("'%s' is a templated extraction, which we not currently support with remote executor extractions", val)
		}
		if len(submatches) == 0 {
			output[key] = &v2.Executor_Remote_Extraction{
				ExtractionType: &v2.Executor_Remote_Extraction_Value{
					Value: val,
				},
			}
			continue
		}
		if len(submatches[0]) < 2 {
			return nil, eris.Errorf("Malformed value for dynamic metadata %s: %s", key, val)
		}
		dotNot, err := dot_notation.DotNotationToPathSegments(submatches[0][1])
		if err != nil || dotNot == nil || dotNot[0] == nil {
			return nil, eris.Errorf("Malformed value for dynamic metadata %s: %s", key, val)
		}
		switch dotNot[0].GetKey() {
		case jsonUtils.HEADERS:
			for _, segment := range dotNot[1:] {
				out += segment.GetKey()
			}
			output[key] = &v2.Executor_Remote_Extraction{
				ExtractionType: &v2.Executor_Remote_Extraction_Header{
					Header: out}}
		case jsonUtils.METADATA:
			var segmentStrings []string
			for _, segment := range dotNot[1:len(dotNot)] {
				segmentStrings = append(segmentStrings, segment.GetKey())
			}
			if len(segmentStrings) == 0 {
				return nil, eris.Errorf("No name specified for dynamic metadata %s: %s", key, val)
			}
			joinedSegments := strings.Join(segmentStrings, ".")
			if strings.Index(joinedSegments, ":") == -1 {
				return nil, eris.Errorf("No namespace specified for dynamic metadata %s: %s", key, val)
			}
			output[key] = &v2.Executor_Remote_Extraction{
				ExtractionType: &v2.Executor_Remote_Extraction_DynamicMetadata{
					DynamicMetadata: &v2.Executor_Remote_Extraction_DynamicMetadataExtraction{
						MetadataNamespace: joinedSegments[:strings.Index(joinedSegments, ":")],
						Key:               joinedSegments[strings.Index(joinedSegments, ":")+1:],
					},
				},
			}
		}
	}
	return output, nil
}

func translateExecutableSchema(params CreateGraphQLApiParams) (*v2.ExecutableSchema, error) {
	graphqlApi := params.Graphqlapi
	extensions, err := TranslateExtensions(params.Artifacts, graphqlApi)
	if err != nil {
		return nil, err
	}
	switch typedExecutor := params.Graphqlapi.GetExecutableSchema().Executor.Executor.(type) {
	case *v1beta1.Executor_Local_:
		schemaStr := graphqlApi.GetExecutableSchema().GetSchemaDefinition()
		_, resolutions, processedSchema, err := processGraphqlSchema(params.Upstreams, schemaStr, typedExecutor.Local.GetResolutions(), params.ProcessedSchemaCache)
		if err != nil {
			return nil, err
		}
		return &v2.ExecutableSchema{
			Executor: &v2.Executor{
				Executor: &v2.Executor_Local_{
					Local: &v2.Executor_Local{
						Resolutions:         resolutions,
						EnableIntrospection: typedExecutor.Local.GetEnableIntrospection(),
					},
				},
			},
			SchemaDefinition: &v3.DataSource{
				Specifier: &v3.DataSource_InlineString{InlineString: printer2.PrettyPrintKubeString(processedSchema)},
			},
			Extensions: extensions,
		}, nil
	case *v1beta1.Executor_Remote_:
		remoteExecutor := typedExecutor.Remote
		headers, err := glooToEnvoyTranslation(remoteExecutor.GetHeaders())
		if err != nil {
			return nil, err
		}
		queryParams, err := glooToEnvoyTranslation(remoteExecutor.GetQueryParams())
		if err != nil {
			return nil, err
		}
		upstream, err := params.Upstreams.Find(remoteExecutor.GetUpstreamRef().GetNamespace(), remoteExecutor.GetUpstreamRef().GetName())
		if err != nil {
			return nil, eris.Errorf(
				"No upstream found on cluster with namespace.name: %s.%s",
				remoteExecutor.GetUpstreamRef().GetNamespace(),
				remoteExecutor.GetUpstreamRef().GetName())
		}
		return &v2.ExecutableSchema{
			Executor: &v2.Executor{
				Executor: &v2.Executor_Remote_{
					Remote: &v2.Executor_Remote{
						ServerUri: &v3.HttpUri{
							Uri: "ignored", // ignored by graphql filter,
							HttpUpstreamType: &v3.HttpUri_Cluster{
								Cluster: translator.UpstreamToClusterName(upstream.GetMetadata().Ref()),
							},
							Timeout: &duration.Duration{
								Seconds: upstream.GetConnectionConfig().GetConnectTimeout().GetSeconds(),
								Nanos:   upstream.GetConnectionConfig().GetConnectTimeout().GetNanos(),
							},
						},
						Request: &v2.Executor_Remote_RemoteSchemaRequest{
							Headers:     headers,
							QueryParams: queryParams,
						},
						SpanName: typedExecutor.Remote.GetSpanName(),
					},
				},
			},
			SchemaDefinition: &v3.DataSource{
				Specifier: &v3.DataSource_InlineString{
					InlineString: params.Graphqlapi.GetExecutableSchema().GetSchemaDefinition(),
				},
			},
		}, nil
	default:
		return nil, eris.Errorf("unsupported executor type %T", typedExecutor)
	}
}

func TranslateExtensions(artifacts types.ArtifactList, api *v1beta1.GraphQLApi) (map[string]*any.Any, error) {
	extensions := map[string]*any.Any{}
	reg := api.GetExecutableSchema().GetGrpcDescriptorRegistry()
	if reg != nil {

		grpcDescRegistry := &v2.GrpcDescriptorRegistry{
			ProtoDescriptors: &v3.DataSource{
				Specifier: nil, // filled in later
			},
		}

		switch regType := reg.DescriptorSet.(type) {
		case *v1beta1.GrpcDescriptorRegistry_ProtoDescriptor:
			grpcDescRegistry.ProtoDescriptors.Specifier = &v3.DataSource_Filename{
				Filename: reg.GetProtoDescriptor(),
			}
		case *v1beta1.GrpcDescriptorRegistry_ProtoDescriptorBin:
			grpcDescRegistry.ProtoDescriptors.Specifier = &v3.DataSource_InlineBytes{
				InlineBytes: reg.GetProtoDescriptorBin(),
			}
		case *v1beta1.GrpcDescriptorRegistry_ProtoRefsList:
			var accumulator []byte
			table := make(map[string]struct{})
			for _, protoRef := range reg.GetProtoRefsList().GetConfigMapRefs() {
				configMap, err := artifacts.Find(protoRef.GetNamespace(), protoRef.GetName())
				if err != nil {
					return nil, eris.Errorf("Could not find ConfigMap with ref %s.%s to use a gRPC proto registry source", protoRef.GetNamespace(), protoRef.GetName())
				}
				if len(configMap.GetData()) == 0 {
					return nil, eris.Errorf("No keys exist in %s.%s", protoRef.GetNamespace(), protoRef.GetName())
				}
				for key, protoData := range configMap.GetData() {
					if protoData == "" {
						return nil, eris.Errorf("Expecting value in configmap %s.%s associated with key %s", protoRef.GetNamespace(), protoRef.GetName(), key)
					}

					bytes, err := base64.StdEncoding.DecodeString(protoData)
					if err != nil {
						return nil, eris.Errorf("Error decoding proto data in %s: %s", key, protoData)
					}

					//Validate the proto
					addr := &descriptor.FileDescriptorProto{}
					err = proto.Unmarshal(bytes, addr)
					if err != nil {
						return nil, eris.Errorf("key %s in configMap %s.%s does not contain valid proto bytes", key, protoRef.GetNamespace(), protoRef.GetName())
					}

					if _, found := table[string(bytes)]; !found {
						accumulator = append(accumulator, bytes...)
						table[string(bytes)] = struct{}{}
					}
				}
			}

			grpcDescRegistry.ProtoDescriptors.Specifier = &v3.DataSource_InlineBytes{
				InlineBytes: accumulator,
			}
		default:
			return nil, eris.Errorf("unimplemented type %T for grpc resolver proto descriptor translation", regType)
		}
		extensions[grpc.GrpcRegistryExtensionName] = utils.MustMessageToAny(grpcDescRegistry)
	}

	if len(extensions) == 0 {
		return nil, nil
	}
	return extensions, nil
}

func processGraphqlSchema(upstreams types.UpstreamList, schema string, resolutions map[string]*v1beta1.Resolution, schemaCache *lru.Cache) (*ast.Document, []*v2.Resolution, string, error) {
	// NOTE:
	// If we add any new inputs to `ParseGraphQLSchema` that affect the output `doc` ast,
	// they must be cached in the schemaCache. Else we risk having false cache hits.
	doc, err := ParseGraphQLSchema(schema)
	if err != nil {
		return nil, nil, "", err
	}
	visitor := directives.NewGraphqlASTVisitor()
	var result []*v2.Resolution
	// Adds a directive visitor to the ast visitor which looks for `@resolve` directives in the schema. For each
	// resolve directive, it translates the resolver with the given name in the control plane resolutions map into
	// a data plane resolver.
	AddResolveDirectiveVisitor(visitor, upstreams, resolutions, &result)
	// Adds a directive visitor to the ast visitor which looks for `@cacheControl` directives in the schema. For each
	// cacheControl directive, it translates it to a data plane cache control and attaches it to the associated resolver.
	AddCacheControlDirectiveVisitor(visitor, &result)
	err = visitor.Visit(doc)
	if err != nil {
		return nil, nil, "", err
	}
	algorithm := fnv.New64a()
	_, err = algorithm.Write([]byte(schema))
	if err != nil {
		return nil, nil, "", err
	}
	if schemaCache == nil {
		// if no schema cache, print out the ast and return
		printedSchema := printer.Print(doc)
		printedSchemaStr, ok := printedSchema.(string)
		if !ok {
			return nil, nil, "", eris.Errorf("cannot convert %v of type %T to string", printedSchemaStr, printedSchemaStr)
		}
		return doc, result, printedSchemaStr, nil
	} else {
		schemaHash := algorithm.Sum64()
		printedSchema, ok := schemaCache.Get(schemaHash)
		if !ok {
			printedSchema = printer.Print(doc)
			schemaCache.Add(schemaHash, printedSchema)
		}
		printedSchemaStr, ok := printedSchema.(string)
		if !ok {
			return nil, nil, "", eris.Errorf("cannot convert %v of type %T to string", printedSchemaStr, printedSchemaStr)
		}
		return doc, result, printedSchemaStr, nil
	}

}

func ParseGraphQLSchema(schema string) (*ast.Document, error) {
	doc, err := parser.Parse(parser.ParseParams{Source: schema})
	if err != nil {
		return nil, eris.Wrapf(err, "unable to parse graphql schema %s", schema)
	}
	return doc, nil
}

func translateResolver(upstreams types.UpstreamList, resolver *v1beta1.Resolution) (*v3.TypedExtensionConfig, error) {
	switch r := resolver.Resolver.(type) {
	case *v1beta1.Resolution_RestResolver:
		return rest.TranslateRestResolver(upstreams, r.RestResolver)
	case *v1beta1.Resolution_GrpcResolver:
		return grpc.TranslateGrpcResolver(upstreams, r.GrpcResolver)
	case *v1beta1.Resolution_MockResolver:
		return mock.TranslateMockResolver(r.MockResolver)
	default:
		return nil, errors.Errorf("unimplemented resolver type: %T", r)
	}
}

func AddResolveDirectiveVisitor(visitor *directives.GraphqlASTVisitor, upstreams types.UpstreamList, resolutions map[string]*v1beta1.Resolution, result *[]*v2.Resolution) {
	visitor.AddDirectiveVisitor(directives.RESOLVER_DIRECTIVE, func(directiveVisitorParams directives.DirectiveVisitorParams) (bool, error) {
		// validate correct usage of the resolve directive
		resolveDirective := directives.NewResolveDirective()
		err := resolveDirective.Validate(directiveVisitorParams)
		if err != nil {
			return false, err
		}

		// check if the resolver referenced here even exists
		resolution := resolutions[resolveDirective.ResolverName]
		if resolution == nil {
			return false, directives.NewGraphqlSchemaError(resolveDirective.ResolverNameAstValue, "resolver %s is not defined",
				resolveDirective.ResolverName)
		}

		queryMatch := &v2.QueryMatcher{
			Match: &v2.QueryMatcher_FieldMatcher_{
				FieldMatcher: &v2.QueryMatcher_FieldMatcher{
					Type:  directiveVisitorParams.Type.Name.Value,
					Field: directiveVisitorParams.DirectiveField.Name.Value,
				},
			},
		}
		res, err := translateResolver(upstreams, resolution)
		if err != nil {
			return false, err
		}
		statsPrefix := resolveDirective.ResolverName
		if sp := resolution.StatPrefix; sp != nil {
			statsPrefix = sp.Value
		}
		statsPrefix = strings.TrimSuffix(statsPrefix, ".") + "."
		if result != nil {
			for i, resul := range *result {
				if proto.Equal(resul.Matcher, queryMatch) {
					(*result)[i].Resolver = res
					(*result)[i].StatPrefix = statsPrefix
					return true, nil
				}
			}
		}
		*result = append(*result, &v2.Resolution{
			Matcher:    queryMatch,
			Resolver:   res,
			StatPrefix: statsPrefix,
		})
		return true, nil
	})
}

func AddCacheControlDirectiveVisitor(visitor *directives.GraphqlASTVisitor, result *[]*v2.Resolution) {
	visitor.AddDirectiveVisitor(directives.CACHE_CONTROL_DIRECTIVE, func(directiveVisitorParams directives.DirectiveVisitorParams) (bool, error) {
		// validate correct usage of the cacheControl directive
		cacheControlDirective := directives.NewCacheControlDirective()
		_, err := cacheControlDirective.Validate(directiveVisitorParams)
		if err != nil {
			return false, err
		}

		cacheControl := cacheControlDirective.CacheControl

		var queryMatchList []*v2.QueryMatcher
		for _, df := range directiveVisitorParams.DirectiveFields {
			// this is a type-level directive. we will set on all fields
			queryMatch := &v2.QueryMatcher{
				Match: &v2.QueryMatcher_FieldMatcher_{
					FieldMatcher: &v2.QueryMatcher_FieldMatcher{
						Type:  directiveVisitorParams.Type.Name.Value,
						Field: df.Name.Value,
					},
				},
			}
			queryMatchList = append(queryMatchList, queryMatch)
		}

		if directiveVisitorParams.DirectiveField != nil {
			// this is a field-directive
			queryMatch := &v2.QueryMatcher{
				Match: &v2.QueryMatcher_FieldMatcher_{
					FieldMatcher: &v2.QueryMatcher_FieldMatcher{
						Type:  directiveVisitorParams.Type.Name.Value,
						Field: directiveVisitorParams.DirectiveField.Name.Value,
					},
				},
			}
			queryMatchList = append(queryMatchList, queryMatch)
		}

		if len(queryMatchList) == 0 {
			return false, eris.Errorf("logic error: no query match generated but `@cacheControl` directive was found")
		}

		for _, queryMatch := range queryMatchList {
			found := false
			if result != nil {
				for i, res := range *result {
					if proto.Equal(res.Matcher, queryMatch) {
						(*result)[i].CacheControl = cacheControl
						found = true
						break
					}
				}
			}
			if !found {
				*result = append(*result, &v2.Resolution{
					Matcher:      queryMatch,
					CacheControl: cacheControl,
				})
			}
		}
		return true, nil
	})
}
