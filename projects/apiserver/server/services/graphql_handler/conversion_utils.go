package graphql_handler

import (
	"github.com/rotisserie/eris"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/graphql/v1beta1"
	graphql_v1beta1 "github.com/solo-io/solo-apis/pkg/api/graphql.gloo.solo.io/v1beta1"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	rpc_edge_v1 "github.com/solo-io/solo-projects/projects/apiserver/pkg/api/rpc.edge.gloo/v1"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/utils/graphql/types"
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
		case *graphql_v1beta1.GraphQLApiSpec_ExecutableSchema:
			summary.ApiTypeSummary = &rpc_edge_v1.GraphqlApiSummary_Executable{
				Executable: &rpc_edge_v1.GraphqlApiSummary_ExecutableSchemaSummary{
					// TODO when we add other executor types (besides local), handle them here
					NumResolvers: uint32(len(graphqlApi.GetSpec().GetExecutableSchema().GetExecutor().GetLocal().GetResolutions())),
				},
			}
		case *graphql_v1beta1.GraphQLApiSpec_StitchedSchema:
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

func NewGraphqlApiList(graphqlapis *graphql_v1beta1.GraphQLApiList) (RequestGraphQlApiList, error) {
	subschemaApiList := map[ResourceRefGolangType]*v1beta1.GraphQLApi{}
	for _, api := range graphqlapis.Items {
		ret := &v1beta1.GraphQLApi{}
		err := types.ConvertGoProtoTypes(&api.Spec, ret)
		if err != nil {
			return nil, eris.Wrapf(err, "unable to convert solo-apis GraphQLApi spec to enterprise gloo GraphQlApi type")
		}
		ret.SetMetadata(&core.Metadata{
			Name:      api.GetName(),
			Namespace: api.GetNamespace(),
		})
		subschemaApiList[ResourceRefGolangType{
			Name:      api.GetName(),
			Namespace: api.GetNamespace(),
		}] = ret
	}
	return subschemaApiList, nil
}

type RequestGraphQlApiList map[ResourceRefGolangType]*v1beta1.GraphQLApi

func (list RequestGraphQlApiList) Find(namespace, name string) (*v1beta1.GraphQLApi, error) {
	if api, ok := list[ResourceRefGolangType{Namespace: namespace, Name: name}]; ok {
		return api, nil
	}
	return nil, eris.Errorf("did not find graphQLApi %v.%v", namespace, name)
}

func (list RequestGraphQlApiList) AsResources() resources.ResourceList {
	var ret resources.ResourceList
	for resourceRef, _ := range list {

		ret = append(ret, &v1beta1.GraphQLApi{
			Metadata: &core.Metadata{
				Name:      resourceRef.Name,
				Namespace: resourceRef.Namespace,
			},
		})
	}
	return ret
}
