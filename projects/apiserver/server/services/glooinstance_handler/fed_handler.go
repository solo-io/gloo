package glooinstance_handler

import (
	"context"
	"sort"

	"github.com/rotisserie/eris"
	"github.com/solo-io/go-utils/contextutils"
	skv2_multicluster "github.com/solo-io/skv2/pkg/multicluster"
	rpc_edge_v1 "github.com/solo-io/solo-projects/projects/apiserver/pkg/api/rpc.edge.gloo/v1"
	"github.com/solo-io/solo-projects/projects/apiserver/server/apiserverutils"
	"github.com/solo-io/solo-projects/projects/apiserver/server/services/glooinstance_handler/config_getter"
	fedv1 "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.solo.io/v1"
	"github.com/solo-io/solo-projects/projects/gloo-fed/pkg/multicluster"
	"go.uber.org/zap"
	"k8s.io/client-go/discovery"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func NewFedGlooInstanceHandler(
	managerSet skv2_multicluster.ManagerSet,
	clusterSet multicluster.ClusterSet,
	envoyConfigClient config_getter.EnvoyConfigDumpGetter,
	glooInstanceClient fedv1.GlooInstanceClient,
) rpc_edge_v1.GlooInstanceApiServer {
	return &fedGlooInstanceHandler{
		managerSet:         managerSet,
		clusterSet:         clusterSet,
		envoyConfigClient:  envoyConfigClient,
		glooInstanceClient: glooInstanceClient,
	}
}

type fedGlooInstanceHandler struct {
	managerSet         skv2_multicluster.ManagerSet
	clusterSet         multicluster.ClusterSet
	envoyConfigClient  config_getter.EnvoyConfigDumpGetter
	glooInstanceClient fedv1.GlooInstanceClient
}

func (h *fedGlooInstanceHandler) ListClusterDetails(ctx context.Context, request *rpc_edge_v1.ListClusterDetailsRequest) (*rpc_edge_v1.ListClusterDetailsResponse, error) {
	glooInstancesByCluster := make(map[string][]*rpc_edge_v1.GlooInstance)
	for _, cluster := range h.clusterSet.ListClusters() {
		glooInstancesByCluster[cluster] = []*rpc_edge_v1.GlooInstance{}
	}
	list, err := h.glooInstanceClient.ListGlooInstance(ctx)
	if err != nil {
		wrapped := eris.Wrapf(err, "Failed to list gloo instances")
		contextutils.LoggerFrom(ctx).Errorw(wrapped.Error(), zap.Error(err), zap.Any("request", request))
		return nil, wrapped
	}
	for _, glooInstance := range list.Items {
		glooInstance := glooInstance
		glooInstancesByCluster[glooInstance.Spec.Cluster] = append(glooInstancesByCluster[glooInstance.Spec.Cluster],
			apiserverutils.ConvertToRpcGlooInstance(&glooInstance))
	}

	var rpcClusterDetails []*rpc_edge_v1.ClusterDetails
	for cluster, glooInstances := range glooInstancesByCluster {
		sortGlooInstances(glooInstances)
		rpcClusterDetails = append(rpcClusterDetails, &rpc_edge_v1.ClusterDetails{
			Cluster:       cluster,
			GlooInstances: glooInstances,
		})
	}
	sort.Slice(rpcClusterDetails, func(i, j int) bool {
		x := rpcClusterDetails[i]
		y := rpcClusterDetails[j]
		return x.GetCluster() < y.GetCluster()
	})
	return &rpc_edge_v1.ListClusterDetailsResponse{
		ClusterDetails: rpcClusterDetails,
	}, nil
}

func (h *fedGlooInstanceHandler) ListGlooInstances(ctx context.Context, request *rpc_edge_v1.ListGlooInstancesRequest) (*rpc_edge_v1.ListGlooInstancesResponse, error) {
	list, err := h.glooInstanceClient.ListGlooInstance(ctx)
	if err != nil {
		wrapped := eris.Wrapf(err, "Failed to list gloo instances")
		contextutils.LoggerFrom(ctx).Errorw(wrapped.Error(), zap.Error(err), zap.Any("request", request))
		return nil, wrapped
	}
	var glooInstances []*rpc_edge_v1.GlooInstance
	for _, glooInstance := range list.Items {
		glooInstance := glooInstance
		glooInstances = append(glooInstances, apiserverutils.ConvertToRpcGlooInstance(&glooInstance))
	}
	sortGlooInstances(glooInstances)

	return &rpc_edge_v1.ListGlooInstancesResponse{
		GlooInstances: glooInstances,
	}, nil
}

func (h *fedGlooInstanceHandler) GetConfigDumps(ctx context.Context, request *rpc_edge_v1.GetConfigDumpsRequest) (*rpc_edge_v1.GetConfigDumpsResponse, error) {
	glooInstance, err := h.glooInstanceClient.GetGlooInstance(ctx, client.ObjectKey{
		Name:      request.GlooInstanceRef.GetName(),
		Namespace: request.GlooInstanceRef.GetNamespace(),
	})
	if err != nil {
		return nil, eris.Wrapf(err, "could not find gloo instance %v", request.GetGlooInstanceRef())
	}

	rpcGlooInstance := apiserverutils.ConvertToRpcGlooInstance(glooInstance)
	mgr, err := h.managerSet.Cluster(glooInstance.Spec.GetCluster())
	if err != nil {
		return nil, err
	}
	discoveryClient, err := discovery.NewDiscoveryClientForConfig(mgr.GetConfig())
	if err != nil {
		return nil, err
	}

	// Get envoy proxy config dumps for gloo instance
	configDumps, err := h.envoyConfigClient.GetConfigs(ctx, rpcGlooInstance, *discoveryClient)
	if err != nil {
		contextutils.LoggerFrom(ctx).Warnf("Unable to get config dump for Gloo Instance %v", glooInstance)
		return nil, err
	}

	return &rpc_edge_v1.GetConfigDumpsResponse{
		ConfigDumps: configDumps,
	}, nil
}
