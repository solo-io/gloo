package glooinstance_handler

import (
	"context"

	"github.com/rotisserie/eris"
	"github.com/solo-io/go-utils/contextutils"
	rpc_edge_v1 "github.com/solo-io/solo-projects/projects/apiserver/pkg/api/rpc.edge.gloo/v1"
	"github.com/solo-io/solo-projects/projects/apiserver/server/services/glooinstance_handler/config_getter"
	"k8s.io/client-go/discovery"
)

func NewSingleClusterGlooInstanceHandler(
	glooInstanceLister SingleClusterGlooInstanceLister,
	discoveryClient discovery.DiscoveryClient,
	envoyConfigClient config_getter.EnvoyConfigDumpGetter,
) rpc_edge_v1.GlooInstanceApiServer {
	return &singleClusterGlooInstanceHandler{
		glooInstanceLister: glooInstanceLister,
		discoveryClient:    discoveryClient,
		envoyConfigClient:  envoyConfigClient,
	}
}

type singleClusterGlooInstanceHandler struct {
	glooInstanceLister SingleClusterGlooInstanceLister
	discoveryClient    discovery.DiscoveryClient
	envoyConfigClient  config_getter.EnvoyConfigDumpGetter
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
	glooInstance, err := h.glooInstanceLister.GetGlooInstance(ctx, request.GetGlooInstanceRef())
	if err != nil {
		return nil, eris.Wrapf(err, "could not find gloo instance %v", request.GetGlooInstanceRef())
	}

	// Get envoy proxy config dumps for gloo instance
	configDumps, err := h.envoyConfigClient.GetConfigs(ctx, glooInstance, h.discoveryClient)
	if err != nil {
		contextutils.LoggerFrom(ctx).Warnf("Unable to get config dump for Gloo Instance %v", glooInstance)
		return nil, err
	}

	return &rpc_edge_v1.GetConfigDumpsResponse{
		ConfigDumps: configDumps,
	}, nil
}
