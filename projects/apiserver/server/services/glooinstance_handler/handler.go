package glooinstance_handler

import (
	"context"
	"sort"

	"github.com/rotisserie/eris"
	"github.com/solo-io/go-utils/contextutils"
	rpc_v1 "github.com/solo-io/solo-projects/projects/apiserver/pkg/api/fed.rpc/v1"
	"github.com/solo-io/solo-projects/projects/apiserver/server/apiserverutils"
	"github.com/solo-io/solo-projects/projects/apiserver/server/services/glooinstance_handler/config_getter"
	fedv1 "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.solo.io/v1"
	"github.com/solo-io/solo-projects/projects/gloo-fed/pkg/multicluster"
	"go.uber.org/zap"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func NewGlooInstanceHandler(
	clusterClient multicluster.ClusterSet,
	configClient config_getter.EnvoyConfigDumpGetter,
	glooInstanceClient fedv1.GlooInstanceClient,
) rpc_v1.GlooInstanceApiServer {
	return &glooInstanceHandler{
		clusterClient:      clusterClient,
		configClient:       configClient,
		glooInstanceClient: glooInstanceClient,
	}
}

type glooInstanceHandler struct {
	clusterClient      multicluster.ClusterSet
	configClient       config_getter.EnvoyConfigDumpGetter
	glooInstanceClient fedv1.GlooInstanceClient
}

func (k *glooInstanceHandler) ListClusterDetails(ctx context.Context, request *rpc_v1.ListClusterDetailsRequest) (*rpc_v1.ListClusterDetailsResponse, error) {
	glooInstancesByCluster := make(map[string][]*rpc_v1.GlooInstance)
	for _, cluster := range k.clusterClient.ListClusters() {
		glooInstancesByCluster[cluster] = []*rpc_v1.GlooInstance{}
	}
	list, err := k.glooInstanceClient.ListGlooInstance(ctx)
	if err != nil {
		wrapped := eris.Wrapf(err, "Failed to get list gloo instances")
		contextutils.LoggerFrom(ctx).Errorw(wrapped.Error(), zap.Error(err), zap.Any("request", request))
		return nil, wrapped
	}
	for _, glooInstance := range list.Items {
		glooInstancesByCluster[glooInstance.Spec.Cluster] = append(glooInstancesByCluster[glooInstance.Spec.Cluster],
			BuildRpcGlooInstance(glooInstance))
	}

	var rpcClusterDetails []*rpc_v1.ClusterDetails
	for cluster, glooInstances := range glooInstancesByCluster {
		sortGlooInstances(glooInstances)
		rpcClusterDetails = append(rpcClusterDetails, &rpc_v1.ClusterDetails{
			Cluster:       cluster,
			GlooInstances: glooInstances,
		})
	}
	sort.Slice(rpcClusterDetails, func(i, j int) bool {
		x := rpcClusterDetails[i]
		y := rpcClusterDetails[j]
		return x.GetCluster() < y.GetCluster()
	})
	return &rpc_v1.ListClusterDetailsResponse{
		ClusterDetails: rpcClusterDetails,
	}, nil
}

func (k *glooInstanceHandler) ListGlooInstances(ctx context.Context, request *rpc_v1.ListGlooInstancesRequest) (*rpc_v1.ListGlooInstancesResponse, error) {
	list, err := k.glooInstanceClient.ListGlooInstance(ctx)
	if err != nil {
		wrapped := eris.Wrapf(err, "Failed to get list gloo instances")
		contextutils.LoggerFrom(ctx).Errorw(wrapped.Error(), zap.Error(err), zap.Any("request", request))
		return nil, wrapped
	}
	var glooInstances []*rpc_v1.GlooInstance
	for _, glooInstance := range list.Items {
		glooInstances = append(glooInstances, BuildRpcGlooInstance(glooInstance))
	}
	sortGlooInstances(glooInstances)

	return &rpc_v1.ListGlooInstancesResponse{
		GlooInstances: glooInstances,
	}, nil
}

func sortGlooInstances(glooInstances []*rpc_v1.GlooInstance) {
	sort.Slice(glooInstances, func(i, j int) bool {
		x := glooInstances[i]
		y := glooInstances[j]
		return x.GetMetadata().GetNamespace()+x.GetMetadata().GetName() < y.GetMetadata().GetNamespace()+y.GetMetadata().GetName()
	})
}

func BuildRpcGlooInstance(glooInstance fedv1.GlooInstance) *rpc_v1.GlooInstance {
	return &rpc_v1.GlooInstance{
		Metadata: apiserverutils.ToMetadata(glooInstance.ObjectMeta),
		Spec:     &glooInstance.Spec,
		Status:   &glooInstance.Status,
	}
}

func (k *glooInstanceHandler) GetConfigDumps(ctx context.Context, request *rpc_v1.GetConfigDumpsRequest) (*rpc_v1.GetConfigDumpsResponse, error) {
	glooInstance, err := k.glooInstanceClient.GetGlooInstance(ctx, client.ObjectKey{
		Name:      request.GlooInstanceRef.GetName(),
		Namespace: request.GlooInstanceRef.GetNamespace(),
	})
	if err != nil {
		return nil, eris.Wrapf(err, "could not find gloo instance %v.%v", request.GlooInstanceRef.GetName(), request.GlooInstanceRef.GetNamespace())
	}
	// Get envoy proxy config dumps for gloo instance
	configDumps, err := k.configClient.GetConfigs(ctx, *glooInstance)
	if err != nil {
		contextutils.LoggerFrom(ctx).Warnf("Unable to get config dump for Gloo Instance %s", glooInstance.GetName())
		return nil, err
	}

	return &rpc_v1.GetConfigDumpsResponse{
		ConfigDumps: configDumps,
	}, nil
}
