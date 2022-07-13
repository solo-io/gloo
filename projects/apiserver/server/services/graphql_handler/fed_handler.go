package graphql_handler

import (
	"context"
	"sort"

	"github.com/ghodss/yaml"
	"github.com/rotisserie/eris"
	"github.com/solo-io/go-utils/contextutils"
	skv2_v1 "github.com/solo-io/skv2/pkg/api/core.skv2.solo.io/v1"
	gloo_v1 "github.com/solo-io/solo-apis/pkg/api/gloo.solo.io/v1"
	graphql_v1beta1 "github.com/solo-io/solo-apis/pkg/api/graphql.gloo.solo.io/v1beta1"
	rpc_edge_v1 "github.com/solo-io/solo-projects/projects/apiserver/pkg/api/rpc.edge.gloo/v1"
	"github.com/solo-io/solo-projects/projects/apiserver/server/apiserverutils"
	fed_v1 "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.solo.io/v1"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func NewFedGraphqlHandler(
	glooInstanceClient fed_v1.GlooInstanceClient,
	settingsClient gloo_v1.SettingsClient,
	graphqlMCClientset graphql_v1beta1.MulticlusterClientset,
) rpc_edge_v1.GraphqlConfigApiServer {
	return &fedGraphqlHandler{
		glooInstanceClient: glooInstanceClient,
		settingsClient:     settingsClient,
		graphqlMCClientset: graphqlMCClientset,
	}
}

type fedGraphqlHandler struct {
	glooInstanceClient fed_v1.GlooInstanceClient
	settingsClient     gloo_v1.SettingsClient
	graphqlMCClientset graphql_v1beta1.MulticlusterClientset
}

func (h *fedGraphqlHandler) GetGraphqlApi(ctx context.Context, request *rpc_edge_v1.GetGraphqlApiRequest) (*rpc_edge_v1.GetGraphqlApiResponse, error) {
	if request.GetGraphqlApiRef() == nil {
		return nil, eris.Errorf("graphqlapi ref missing from request: %v", request)
	}

	clientset, err := h.graphqlMCClientset.Cluster(request.GetGraphqlApiRef().GetClusterName())
	if err != nil {
		wrapped := eris.Wrapf(err, "Failed to get graphql client set")
		contextutils.LoggerFrom(ctx).Errorw(wrapped.Error(), zap.Error(err), zap.Any("request", request))
		return nil, wrapped
	}

	graphqlApi, err := clientset.GraphQLApis().GetGraphQLApi(ctx, client.ObjectKey{
		Namespace: request.GetGraphqlApiRef().GetNamespace(),
		Name:      request.GetGraphqlApiRef().GetName(),
	})
	if err != nil {
		wrapped := eris.Wrapf(err, "Failed to get graphqlapi")
		contextutils.LoggerFrom(ctx).Errorw(wrapped.Error(), zap.Error(err), zap.Any("request", request))
		return nil, wrapped
	}

	// find which gloo instance this graphql api belongs to, by finding a gloo instance that is watching
	// the graphql api's namespace
	glooInstance, err := h.getGlooInstanceForNamespace(ctx, request.GetGraphqlApiRef().GetNamespace())
	if err != nil {
		return nil, err
	}
	return &rpc_edge_v1.GetGraphqlApiResponse{
		GraphqlApi: convertGraphqlApi(graphqlApi, &skv2_v1.ObjectRef{
			Name:      glooInstance.GetName(),
			Namespace: glooInstance.GetNamespace(),
		}, glooInstance.Spec.GetCluster()),
	}, nil
}

// find a gloo instance that is either watching the given namespace or watching all namespaces
func (h *fedGraphqlHandler) getGlooInstanceForNamespace(ctx context.Context, namespace string) (*fed_v1.GlooInstance, error) {
	glooInstances, err := h.glooInstanceClient.ListGlooInstance(ctx)
	if err != nil {
		return nil, err
	}

	for _, instance := range glooInstances.Items {
		instance := instance
		watchedNamespaces := instance.Spec.GetControlPlane().GetWatchedNamespaces()
		if len(watchedNamespaces) == 0 {
			return &instance, nil
		}
		for _, ns := range watchedNamespaces {
			if ns == namespace {
				return &instance, nil
			}
		}
	}
	return nil, eris.Errorf("could not find a gloo instance with watched namespace %s", namespace)
}

func (h *fedGraphqlHandler) ListGraphqlApis(ctx context.Context, request *rpc_edge_v1.ListGraphqlApisRequest) (*rpc_edge_v1.ListGraphqlApisResponse, error) {
	var rpcGraphqlApis []*rpc_edge_v1.GraphqlApi
	if request.GetGlooInstanceRef() == nil || request.GetGlooInstanceRef().GetName() == "" || request.GetGlooInstanceRef().GetNamespace() == "" {
		// List graphqlApis across all gloo edge instances
		instanceList, err := h.glooInstanceClient.ListGlooInstance(ctx)
		if err != nil {
			wrapped := eris.Wrapf(err, "Failed to list gloo edge instances")
			contextutils.LoggerFrom(ctx).Errorw(wrapped.Error(), zap.Error(err), zap.Any("request", request))
			return nil, wrapped
		}
		for _, instance := range instanceList.Items {
			rpcGraphqlApiList, err := h.listGraphqlApisForGlooInstance(ctx, &instance)
			if err != nil {
				wrapped := eris.Wrapf(err, "Failed to list graphqlApis for gloo edge instance %s.%s", instance.GetNamespace(), instance.GetName())
				contextutils.LoggerFrom(ctx).Errorw(wrapped.Error(), zap.Error(err), zap.Any("request", request))
				return nil, wrapped
			}
			rpcGraphqlApis = append(rpcGraphqlApis, rpcGraphqlApiList...)
		}
	} else {
		// List graphqlApis for a specific gloo edge instance
		instance, err := h.glooInstanceClient.GetGlooInstance(ctx, types.NamespacedName{
			Name:      request.GetGlooInstanceRef().GetName(),
			Namespace: request.GetGlooInstanceRef().GetNamespace(),
		})
		if err != nil {
			wrapped := eris.Wrapf(err, "Failed to get gloo edge instance %s.%s", instance.GetNamespace(), instance.GetName())
			contextutils.LoggerFrom(ctx).Errorw(wrapped.Error(), zap.Error(err), zap.Any("request", request))
			return nil, wrapped
		}
		rpcGraphqlApis, err = h.listGraphqlApisForGlooInstance(ctx, instance)
		if err != nil {
			wrapped := eris.Wrapf(err, "Failed to list graphqlApis for gloo edge instance %s.%s", instance.GetNamespace(), instance.GetName())
			contextutils.LoggerFrom(ctx).Errorw(wrapped.Error(), zap.Error(err), zap.Any("request", request))
			return nil, wrapped
		}
	}

	summaries, err := getGraphqlApiSummaries(rpcGraphqlApis)
	if err != nil {
		return nil, err
	}
	return &rpc_edge_v1.ListGraphqlApisResponse{
		GraphqlApis: summaries,
	}, nil
}

func (h *fedGraphqlHandler) listGraphqlApisForGlooInstance(ctx context.Context, instance *fed_v1.GlooInstance) ([]*rpc_edge_v1.GraphqlApi, error) {
	clientset, err := h.graphqlMCClientset.Cluster(instance.Spec.GetCluster())
	if err != nil {
		return nil, err
	}
	graphqlApiClient := clientset.GraphQLApis()

	var graphqlApiGraphqlApiList []*graphql_v1beta1.GraphQLApi
	watchedNamespaces := instance.Spec.GetControlPlane().GetWatchedNamespaces()
	if len(watchedNamespaces) != 0 {
		for _, ns := range watchedNamespaces {
			list, err := graphqlApiClient.ListGraphQLApi(ctx, client.InNamespace(ns))
			if err != nil {
				return nil, err
			}
			for i, _ := range list.Items {
				graphqlApiGraphqlApiList = append(graphqlApiGraphqlApiList, &list.Items[i])
			}
		}
	} else {
		list, err := graphqlApiClient.ListGraphQLApi(ctx)
		if err != nil {
			return nil, err
		}
		for i, _ := range list.Items {
			graphqlApiGraphqlApiList = append(graphqlApiGraphqlApiList, &list.Items[i])
		}
	}
	sort.Slice(graphqlApiGraphqlApiList, func(i, j int) bool {
		x := graphqlApiGraphqlApiList[i]
		y := graphqlApiGraphqlApiList[j]
		return x.GetNamespace()+x.GetName() < y.GetNamespace()+y.GetName()
	})

	var rpcGraphqlApis []*rpc_edge_v1.GraphqlApi
	for _, graphqlApi := range graphqlApiGraphqlApiList {
		rpcGraphqlApis = append(rpcGraphqlApis, convertGraphqlApi(graphqlApi, &skv2_v1.ObjectRef{
			Name:      instance.GetName(),
			Namespace: instance.GetNamespace(),
		}, instance.Spec.GetCluster()))
	}
	return rpcGraphqlApis, nil
}

func convertGraphqlApi(graphqlApi *graphql_v1beta1.GraphQLApi, glooInstance *skv2_v1.ObjectRef, cluster string) *rpc_edge_v1.GraphqlApi {
	m := &rpc_edge_v1.GraphqlApi{
		Metadata:     apiserverutils.ToMetadata(graphqlApi.ObjectMeta),
		GlooInstance: glooInstance,
		Spec:         &graphqlApi.Spec,
		Status:       &graphqlApi.Status,
	}
	m.Metadata.ClusterName = cluster
	return m
}

func (h *fedGraphqlHandler) GetGraphqlApiYaml(ctx context.Context, request *rpc_edge_v1.GetGraphqlApiYamlRequest) (*rpc_edge_v1.GetGraphqlApiYamlResponse, error) {
	if request.GetGraphqlApiRef() == nil {
		return nil, eris.Errorf("graphqlapi ref missing from request: %v", request)
	}

	clientset, err := h.graphqlMCClientset.Cluster(request.GetGraphqlApiRef().GetClusterName())
	if err != nil {
		wrapped := eris.Wrapf(err, "Failed to get graphql client set")
		contextutils.LoggerFrom(ctx).Errorw(wrapped.Error(), zap.Error(err), zap.Any("request", request))
		return nil, wrapped
	}
	graphqlApi, err := clientset.GraphQLApis().GetGraphQLApi(ctx, client.ObjectKey{
		Namespace: request.GetGraphqlApiRef().GetNamespace(),
		Name:      request.GetGraphqlApiRef().GetName(),
	})
	if err != nil {
		wrapped := eris.Wrapf(err, "Failed to get graphqlapi")
		contextutils.LoggerFrom(ctx).Errorw(wrapped.Error(), zap.Error(err), zap.Any("request", request))
		return nil, wrapped
	}
	content, err := yaml.Marshal(graphqlApi)
	if err != nil {
		wrapped := eris.Wrapf(err, "Failed to marshal kube resource into yaml")
		contextutils.LoggerFrom(ctx).Errorw(wrapped.Error(), zap.Error(err), zap.Any("request", request))
		return nil, wrapped
	}
	return &rpc_edge_v1.GetGraphqlApiYamlResponse{
		YamlData: &rpc_edge_v1.ResourceYaml{
			Yaml: string(content),
		},
	}, nil
}

func (h *fedGraphqlHandler) CreateGraphqlApi(ctx context.Context, request *rpc_edge_v1.CreateGraphqlApiRequest) (*rpc_edge_v1.CreateGraphqlApiResponse, error) {
	err := apiserverutils.CheckUpdatesAllowed(ctx, h.settingsClient)
	if err != nil {
		return nil, err
	}
	err = h.checkGraphqlApiRef(request.GetGraphqlApiRef())
	if err != nil {
		return nil, err
	}
	if request.GetSpec() == nil {
		return nil, eris.Errorf("graphqlapi spec missing from request: %v", request)
	}

	// make sure the schema definition is valid before we save anything
	err = ValidateSchemaDefinition(&rpc_edge_v1.ValidateSchemaDefinitionRequest{
		Input: &rpc_edge_v1.ValidateSchemaDefinitionRequest_SchemaDefinition{
			SchemaDefinition: request.GetSpec().GetExecutableSchema().GetSchemaDefinition(),
		},
	})
	if err != nil {
		wrapped := eris.Wrapf(err, "Rejected graphqlapi creation due to validation error")
		contextutils.LoggerFrom(ctx).Errorw(wrapped.Error(), zap.Error(err), zap.Any("request", request))
		return nil, wrapped
	}

	clientset, err := h.graphqlMCClientset.Cluster(request.GetGraphqlApiRef().GetClusterName())
	if err != nil {
		wrapped := eris.Wrapf(err, "Failed to get graphql client set")
		contextutils.LoggerFrom(ctx).Errorw(wrapped.Error(), zap.Error(err), zap.Any("request", request))
		return nil, wrapped
	}

	graphqlApi := &graphql_v1beta1.GraphQLApi{
		ObjectMeta: apiserverutils.RefToObjectMeta(*request.GetGraphqlApiRef()),
		Spec:       *request.GetSpec(),
	}
	err = clientset.GraphQLApis().CreateGraphQLApi(ctx, graphqlApi)
	if err != nil {
		wrapped := eris.Wrapf(err, "Failed to create graphqlapi")
		contextutils.LoggerFrom(ctx).Errorw(wrapped.Error(), zap.Error(err), zap.Any("request", request))
		return nil, wrapped
	}

	// find which gloo instance this graphql api belongs to, by finding a gloo instance that is watching
	// the graphql api's namespace
	glooInstance, err := h.getGlooInstanceForNamespace(ctx, request.GetGraphqlApiRef().GetNamespace())
	if err != nil {
		return nil, err
	}
	return &rpc_edge_v1.CreateGraphqlApiResponse{
		GraphqlApi: convertGraphqlApi(graphqlApi, &skv2_v1.ObjectRef{
			Name:      glooInstance.GetName(),
			Namespace: glooInstance.GetNamespace(),
		}, glooInstance.Spec.GetCluster()),
	}, nil
}

func (h *fedGraphqlHandler) UpdateGraphqlApi(ctx context.Context, request *rpc_edge_v1.UpdateGraphqlApiRequest) (*rpc_edge_v1.UpdateGraphqlApiResponse, error) {
	err := apiserverutils.CheckUpdatesAllowed(ctx, h.settingsClient)
	if err != nil {
		return nil, err
	}
	err = h.checkGraphqlApiRef(request.GetGraphqlApiRef())
	if err != nil {
		return nil, err
	}
	if request.GetSpec() == nil {
		return nil, eris.Errorf("graphqlapi spec missing from request: %v", request)
	}

	clientset, err := h.graphqlMCClientset.Cluster(request.GetGraphqlApiRef().GetClusterName())
	if err != nil {
		wrapped := eris.Wrapf(err, "Failed to get graphql client set")
		contextutils.LoggerFrom(ctx).Errorw(wrapped.Error(), zap.Error(err), zap.Any("request", request))
		return nil, wrapped
	}

	// first get the existing graphqlapi
	graphqlApi, err := clientset.GraphQLApis().GetGraphQLApi(ctx, client.ObjectKey{
		Namespace: request.GetGraphqlApiRef().GetNamespace(),
		Name:      request.GetGraphqlApiRef().GetName(),
	})
	if err != nil {
		wrapped := eris.Wrapf(err, "Cannot edit a graphqlapi that does not exist")
		contextutils.LoggerFrom(ctx).Errorw(wrapped.Error(), zap.Error(err), zap.Any("request", request))
		return nil, wrapped
	}

	// make sure the new schema definition is valid before we save anything
	err = ValidateGraphqlApiUpdate(ctx, h.settingsClient, graphqlApi, request.GetSpec())
	if err != nil {
		wrapped := eris.Wrapf(err, "Rejected graphqlapi update due to validation error")
		contextutils.LoggerFrom(ctx).Errorw(wrapped.Error(), zap.Error(err), zap.Any("request", request))
		return nil, wrapped
	}

	// apply the changes to its spec
	graphqlApi.Spec = *request.GetSpec()

	// save the updated graphqlapi
	err = clientset.GraphQLApis().UpdateGraphQLApi(ctx, graphqlApi)
	if err != nil {
		wrapped := eris.Wrapf(err, "Failed to update graphqlapi")
		contextutils.LoggerFrom(ctx).Errorw(wrapped.Error(), zap.Error(err), zap.Any("request", request))
		return nil, wrapped
	}

	// find which gloo instance this graphql api belongs to, by finding a gloo instance that is watching
	// the graphql api's namespace
	glooInstance, err := h.getGlooInstanceForNamespace(ctx, request.GetGraphqlApiRef().GetNamespace())
	if err != nil {
		return nil, err
	}
	return &rpc_edge_v1.UpdateGraphqlApiResponse{
		GraphqlApi: convertGraphqlApi(graphqlApi, &skv2_v1.ObjectRef{
			Name:      glooInstance.GetName(),
			Namespace: glooInstance.GetNamespace(),
		}, glooInstance.Spec.GetCluster()),
	}, nil
}

func (h *fedGraphqlHandler) DeleteGraphqlApi(ctx context.Context, request *rpc_edge_v1.DeleteGraphqlApiRequest) (*rpc_edge_v1.DeleteGraphqlApiResponse, error) {
	err := apiserverutils.CheckUpdatesAllowed(ctx, h.settingsClient)
	if err != nil {
		return nil, err
	}
	err = h.checkGraphqlApiRef(request.GetGraphqlApiRef())
	if err != nil {
		return nil, err
	}

	clientset, err := h.graphqlMCClientset.Cluster(request.GetGraphqlApiRef().GetClusterName())
	if err != nil {
		wrapped := eris.Wrapf(err, "Failed to get graphql client set")
		contextutils.LoggerFrom(ctx).Errorw(wrapped.Error(), zap.Error(err), zap.Any("request", request))
		return nil, wrapped
	}

	err = clientset.GraphQLApis().DeleteGraphQLApi(ctx, client.ObjectKey{
		Name:      request.GetGraphqlApiRef().GetName(),
		Namespace: request.GetGraphqlApiRef().GetNamespace(),
	})
	if err != nil {
		wrapped := eris.Wrapf(err, "Failed to delete graphqlapi")
		contextutils.LoggerFrom(ctx).Errorw(wrapped.Error(), zap.Error(err), zap.Any("request", request))
		return nil, wrapped
	}

	return &rpc_edge_v1.DeleteGraphqlApiResponse{
		GraphqlApiRef: request.GetGraphqlApiRef(),
	}, nil
}

func (h *fedGraphqlHandler) ValidateResolverYaml(_ context.Context, request *rpc_edge_v1.ValidateResolverYamlRequest) (*rpc_edge_v1.ValidateResolverYamlResponse, error) {
	err := ValidateResolverYaml(request.GetYaml(), request.GetResolverType())
	if err != nil {
		return nil, err
	}
	return &rpc_edge_v1.ValidateResolverYamlResponse{}, nil
}

func (h *fedGraphqlHandler) ValidateSchemaDefinition(_ context.Context, request *rpc_edge_v1.ValidateSchemaDefinitionRequest) (*rpc_edge_v1.ValidateSchemaDefinitionResponse, error) {
	err := ValidateSchemaDefinition(request)
	if err != nil {
		return nil, err
	}
	return &rpc_edge_v1.ValidateSchemaDefinitionResponse{}, nil
}

func (h *fedGraphqlHandler) checkGraphqlApiRef(ref *skv2_v1.ClusterObjectRef) error {
	if ref == nil || ref.GetName() == "" || ref.GetNamespace() == "" || ref.GetClusterName() == "" {
		return eris.Errorf("request does not contain valid graphqlapi ref: %v", ref)
	}
	return nil
}

func (h *fedGraphqlHandler) GetStitchedSchemaDefinition(ctx context.Context, request *rpc_edge_v1.GetStitchedSchemaDefinitionRequest) (*rpc_edge_v1.GetStitchedSchemaDefinitionResponse, error) {
	clusterName := request.GetStitchedSchemaApiRef().GetClusterName()
	gqlClient, err := h.graphqlMCClientset.Cluster(clusterName)
	if err != nil {
		return nil, eris.Wrapf(err, "unable to get GraphQLApi client for cluster %s", clusterName)
	}
	return GetStitchedSchemaDefinition(ctx, gqlClient.GraphQLApis(), request)
}

func (h *fedGraphqlHandler) GetSchemaDiff(_ context.Context, request *rpc_edge_v1.GetSchemaDiffRequest) (*rpc_edge_v1.GetSchemaDiffResponse, error) {
	return GetSchemaDiff(request)
}
