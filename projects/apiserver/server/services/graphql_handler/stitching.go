package graphql_handler

import (
	"context"

	"github.com/rotisserie/eris"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/graphql/v1beta1"
	graphql_v1beta1 "github.com/solo-io/solo-apis/pkg/api/graphql.gloo.solo.io/v1beta1"
	rpc_edge_v1 "github.com/solo-io/solo-projects/projects/apiserver/pkg/api/rpc.edge.gloo/v1"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/utils/graphql/translation"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/utils/graphql/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func GetStitchedSchemaDefinition(ctx context.Context, gqlClient graphql_v1beta1.GraphQLApiClient, request *rpc_edge_v1.GetStitchedSchemaDefinitionRequest) (*rpc_edge_v1.GetStitchedSchemaDefinitionResponse, error) {
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
	glooStitchedSchemaConfig := &v1beta1.StitchedSchema{}
	err = types.ConvertGoProtoTypes(stitchedSchemaCfg, glooStitchedSchemaConfig)
	if err != nil {
		return nil, err
	}
	graphqlApiList, err := NewGraphqlApiList(gqlApiList)
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
