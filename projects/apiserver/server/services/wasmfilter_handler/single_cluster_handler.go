package wasmfilter_handler

import (
	"context"

	rpc_edge_v1 "github.com/solo-io/solo-projects/projects/apiserver/pkg/api/rpc.edge.gloo/v1"
)

func NewSingleClusterWasmFilterHandler() rpc_edge_v1.WasmFilterApiServer {
	return &singleClusterWamFilterHandler{}
}

type singleClusterWamFilterHandler struct {
}

func (h *singleClusterWamFilterHandler) DescribeWasmFilter(ctx context.Context, request *rpc_edge_v1.DescribeWasmFilterRequest) (*rpc_edge_v1.DescribeWasmFilterResponse, error) {
	// TODO implement
	return &rpc_edge_v1.DescribeWasmFilterResponse{
		WasmFilter: &rpc_edge_v1.WasmFilter{},
	}, nil
}

func (h *singleClusterWamFilterHandler) ListWasmFilters(ctx context.Context, request *rpc_edge_v1.ListWasmFiltersRequest) (*rpc_edge_v1.ListWasmFiltersResponse, error) {
	// TODO implement
	return &rpc_edge_v1.ListWasmFiltersResponse{
		WasmFilters: []*rpc_edge_v1.WasmFilter{},
	}, nil
}
