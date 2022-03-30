package graphql_handler

import (
	"github.com/rotisserie/eris"
	graphql_v1alpha1 "github.com/solo-io/solo-apis/pkg/api/graphql.gloo.solo.io/v1alpha1"
	rpc_edge_v1 "github.com/solo-io/solo-projects/projects/apiserver/pkg/api/rpc.edge.gloo/v1"
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
