package graphql_handler

import (
	"context"

	"github.com/golang/protobuf/proto"
	"github.com/rotisserie/eris"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/graphql/v1alpha1"
	graphql_v1alpha1 "github.com/solo-io/solo-apis/pkg/api/graphql.gloo.solo.io/v1alpha1"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	rpc_edge_v1 "github.com/solo-io/solo-projects/projects/apiserver/pkg/api/rpc.edge.gloo/v1"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/utils/graphql/translation"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Converts a list of GraphqlApi into a list of GraphqlApiSummary
func getGraphqlApiSummaries(graphqlApis []*rpc_edge_v1.GraphqlApi) ([]*rpc_edge_v1.GraphqlApiSummary, error) {
	summaries := make([]*rpc_edge_v1.GraphqlApiSummary, 0, len(graphqlApis))
	for _, graphqlApi := range graphqlApis {
		summary := &rpc_edge_v1.GraphqlApiSummary{
			Metadata:     graphqlApi.GetMetadata(),
			Status:       graphqlApi.GetStatus(),
			GlooInstance: graphqlApi.GetGlooInstance(),
		}

		switch graphqlApi.GetSpec().GetSchema().(type) {
		case *graphql_v1alpha1.GraphQLApiSpec_ExecutableSchema:
			summary.ApiTypeSummary = &rpc_edge_v1.GraphqlApiSummary_Executable{
				Executable: &rpc_edge_v1.GraphqlApiSummary_ExecutableSchemaSummary{
					// TODO when we add other executor types (besides local), handle them here
					NumResolvers: uint32(len(graphqlApi.GetSpec().GetExecutableSchema().GetExecutor().GetLocal().GetResolutions())),
				},
			}
		case *graphql_v1alpha1.GraphQLApiSpec_StitchedSchema:
			summary.ApiTypeSummary = &rpc_edge_v1.GraphqlApiSummary_Stitched{
				Stitched: &rpc_edge_v1.GraphqlApiSummary_StitchedSchemaSummary{
					NumSubschemas: uint32(len(graphqlApi.GetSpec().GetStitchedSchema().GetSubschemas())),
				},
			}
		default:
			return nil, eris.Errorf("Unexpected GraphQLApi schema type: %v", graphqlApi.GetSpec().GetSchema())
		}

		summaries = append(summaries, summary)
	}
	return summaries, nil
}

type ResourceRefGolangType struct {
	Name      string
	Namespace string
}

func NewGraphqlApiList(graphqlapis *graphql_v1alpha1.GraphQLApiList, subschemas []*graphql_v1alpha1.StitchedSchema_SubschemaConfig) (RequestGraphQlApiList, error) {
	graphqlapiMap := map[ResourceRefGolangType]*graphql_v1alpha1.StitchedSchema_SubschemaConfig{}
	for _, subschema := range subschemas {
		graphqlapiMap[ResourceRefGolangType{
			Name:      subschema.GetName(),
			Namespace: subschema.GetNamespace(),
		}] = subschema
	}
	var subschemaApiList map[ResourceRefGolangType]*v1alpha1.GraphQLApi
	for _, api := range graphqlapis.Items {
		if subschema, ok := graphqlapiMap[ResourceRefGolangType{Name: api.Name, Namespace: api.Namespace}]; ok {
			ret := &v1alpha1.GraphQLApi{}
			err := ConvertGoProtoTypes(&api.Spec, ret)
			if err != nil {
				return nil, eris.Wrapf(err, "unable to convert solo-apis GraphQLApi spec to enterprise gloo GraphQlApi type")
			}
			ret.SetMetadata(&core.Metadata{
				Name:      api.GetName(),
				Namespace: api.GetNamespace(),
			})
			subschemaApiList[ResourceRefGolangType{
				Name:      subschema.GetName(),
				Namespace: subschema.GetNamespace(),
			}] = ret
			break
		} else {
			return nil, eris.Errorf("did not find GraphQLApi subschema %s.%s", subschema.GetNamespace(), subschema.GetName())
		}
	}
	return subschemaApiList, nil
}

func GetStitchedSchemaDefinition(ctx context.Context, gqlClient graphql_v1alpha1.GraphQLApiClient, request *rpc_edge_v1.GetStitchedSchemaDefinitionRequest) (*rpc_edge_v1.GetStitchedSchemaDefinitionResponse, error) {
	stitchedSchema, err := gqlClient.GetGraphQLApi(ctx, client.ObjectKey{
		Namespace: request.GetStitchedSchemaApiRef().GetNamespace(),
		Name:      request.GetStitchedSchemaApiRef().GetName(),
	})
	if err != nil {
		return nil, eris.Wrapf(err, "unable to find stitched schema GraphQLApi %s.%s", request.GetStitchedSchemaApiRef().GetNamespace(),
			request.GetStitchedSchemaApiRef().GetName())
	}
	gqlApiList, err := gqlClient.ListGraphQLApi(ctx)
	if err != nil {
		return nil, eris.Wrap(err, "unable to list all graphql apis")
	}
	stitchedSchemaCfg := stitchedSchema.Spec.GetStitchedSchema()
	if stitchedSchemaCfg == nil {
		return nil, eris.Errorf("GraphQLApi %s.%s does not have a stitched schema definition",
			request.GetStitchedSchemaApiRef().GetNamespace(), request.GetStitchedSchemaApiRef().GetName())
	}
	glooStitchedSchemaConfig := &v1alpha1.StitchedSchema{}
	err = ConvertGoProtoTypes(stitchedSchemaCfg, glooStitchedSchemaConfig)
	if err != nil {
		return nil, err
	}
	graphqlApiList, err := NewGraphqlApiList(gqlApiList, stitchedSchemaCfg.GetSubschemas())
	if err != nil {
		return nil, err
	}
	stitchedSchemaDef, err := translation.GetStitchedSchemaDefinition(glooStitchedSchemaConfig, graphqlApiList)
	if err != nil {
		return nil, err
	}
	return &rpc_edge_v1.GetStitchedSchemaDefinitionResponse{
		StitchedSchemaSdl: stitchedSchemaDef,
	}, nil
}

func ConvertGoProtoTypes(inputMessage proto.Message, outputProtoMessage proto.Message) error {
	protoIntermediateBytes, err := proto.Marshal(inputMessage)
	if err != nil {
		return eris.Wrapf(err, "proto message %s cannot be marshalled", inputMessage.String())
	}
	err = proto.Unmarshal(protoIntermediateBytes, outputProtoMessage)
	if err != nil {
		return eris.Wrapf(err, "proto message %s cannot be unmarshalled into proto message %s", inputMessage.String(), outputProtoMessage.String())
	}
	return nil
}

type RequestGraphQlApiList map[ResourceRefGolangType]*v1alpha1.GraphQLApi

func (list RequestGraphQlApiList) Find(namespace, name string) (*v1alpha1.GraphQLApi, error) {
	if api, ok := list[ResourceRefGolangType{Namespace: namespace, Name: name}]; ok {
		return api, nil
	}
	return nil, eris.Errorf("did not find graphQLApi %v.%v", namespace, name)
}

func (list RequestGraphQlApiList) AsResources() resources.ResourceList {
	var ret resources.ResourceList
	for resourceRef, _ := range list {

		ret = append(ret, &v1alpha1.GraphQLApi{
			Metadata: &core.Metadata{
				Name:      resourceRef.Name,
				Namespace: resourceRef.Namespace,
			},
		})
	}
	return ret
}
