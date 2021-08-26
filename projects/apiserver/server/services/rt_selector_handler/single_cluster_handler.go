package rt_selector_handler

import (
	"context"

	rpc_edge_v1 "github.com/solo-io/solo-projects/projects/apiserver/pkg/api/rpc.edge.gloo/v1"
)

func NewSingleClusterVirtualServiceRoutesHandler() rpc_edge_v1.VirtualServiceRoutesApiServer {
	return &singleClusterVsRoutesHandler{}
}

type singleClusterVsRoutesHandler struct {
}

func (h *singleClusterVsRoutesHandler) GetVirtualServiceRoutes(ctx context.Context, request *rpc_edge_v1.GetVirtualServiceRoutesRequest) (*rpc_edge_v1.GetVirtualServiceRoutesResponse, error) {
	// TODO implement
	return &rpc_edge_v1.GetVirtualServiceRoutesResponse{
		VsRoutes: []*rpc_edge_v1.SubRouteTableRow{},
	}, nil
}
