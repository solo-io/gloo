package graphql_handler

import (
	"github.com/rotisserie/eris"
	rpc_edge_v1 "github.com/solo-io/solo-projects/projects/apiserver/pkg/api/rpc.edge.gloo/v1"
	sticthing "github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/graphql/v8go"
)

func GetSchemaDiff(request *rpc_edge_v1.GetSchemaDiffRequest) (*rpc_edge_v1.GetSchemaDiffResponse, error) {
	in := request.GetDiffInput()
	if in == nil {
		return nil, eris.New("must provide input to schema diff")
	}
	runner, err := sticthing.GetStitchingScriptRunner()
	if err != nil {
		return nil, eris.Wrap(err, "could not create the stitching script runner for GetSchemaDiff")
	}
	out, err := runner.RunSchemaDiff(in)
	if err != nil {
		return nil, err
	}
	return &rpc_edge_v1.GetSchemaDiffResponse{
		DiffOutput: out,
	}, nil
}
