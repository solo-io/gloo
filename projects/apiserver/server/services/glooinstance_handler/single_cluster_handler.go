package glooinstance_handler

import (
	"context"

	rpc_edge_v1 "github.com/solo-io/solo-projects/projects/apiserver/pkg/api/rpc.edge.gloo/v1"
)

func NewSingleClusterGlooInstanceHandler(
	glooInstanceLister SingleClusterGlooInstanceLister,
) rpc_edge_v1.GlooInstanceApiServer {
	return &singleClusterGlooInstanceHandler{
		glooInstanceLister: glooInstanceLister,
	}
}

type singleClusterGlooInstanceHandler struct {
	glooInstanceLister SingleClusterGlooInstanceLister
}

func (h *singleClusterGlooInstanceHandler) ListClusterDetails(ctx context.Context, _ *rpc_edge_v1.ListClusterDetailsRequest) (*rpc_edge_v1.ListClusterDetailsResponse, error) {
	glooInstances, err := h.glooInstanceLister.ListGlooInstances(ctx)
	if err != nil {
		return nil, err
	}

	return &rpc_edge_v1.ListClusterDetailsResponse{
		ClusterDetails: []*rpc_edge_v1.ClusterDetails{
			{
				Cluster:       ClusterName,
				GlooInstances: glooInstances,
			},
		},
	}, nil
}

func (h *singleClusterGlooInstanceHandler) ListGlooInstances(ctx context.Context, _ *rpc_edge_v1.ListGlooInstancesRequest) (*rpc_edge_v1.ListGlooInstancesResponse, error) {
	glooInstances, err := h.glooInstanceLister.ListGlooInstances(ctx)
	if err != nil {
		return nil, err
	}

	return &rpc_edge_v1.ListGlooInstancesResponse{
		GlooInstances: glooInstances,
	}, nil
}

func (h *singleClusterGlooInstanceHandler) GetConfigDumps(ctx context.Context, request *rpc_edge_v1.GetConfigDumpsRequest) (*rpc_edge_v1.GetConfigDumpsResponse, error) {
	// TODO implement
	return &rpc_edge_v1.GetConfigDumpsResponse{
		ConfigDumps: []*rpc_edge_v1.ConfigDump{},
	}, nil
}
