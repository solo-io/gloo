package graphql_handler

import (
	"github.com/rotisserie/eris"
	rpc_edge_v1 "github.com/solo-io/solo-projects/projects/apiserver/pkg/api/rpc.edge.gloo/v1"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/utils/graphql/validation"
)

func GetSchemaDiff(request *rpc_edge_v1.GetSchemaDiffRequest) (*rpc_edge_v1.GetSchemaDiffResponse, error) {
	in := request.GetDiffInput()
	if in == nil {
		return nil, eris.New("must provide input to schema diff")
	}
	out, err := validation.GetSchemaDiff(in)
	if err != nil {
		return nil, err
	}
	return &rpc_edge_v1.GetSchemaDiffResponse{
		DiffOutput: out,
	}, nil
}
