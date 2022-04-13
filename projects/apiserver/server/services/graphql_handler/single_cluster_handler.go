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
	"github.com/solo-io/solo-projects/projects/apiserver/server/services/glooinstance_handler"
	"go.uber.org/zap"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func NewSingleClusterGraphqlHandler(
	graphqlClientset graphql_v1beta1.Clientset,
	glooInstanceLister glooinstance_handler.SingleClusterGlooInstanceLister,
	settingsClient gloo_v1.SettingsClient,
) rpc_edge_v1.GraphqlConfigApiServer {
	return &singleClusterGraphqlHandler{
		graphqlClientset:   graphqlClientset,
		glooInstanceLister: glooInstanceLister,
		settingsClient:     settingsClient,
	}
}

type singleClusterGraphqlHandler struct {
	graphqlClientset   graphql_v1beta1.Clientset
	glooInstanceLister glooinstance_handler.SingleClusterGlooInstanceLister
	settingsClient     gloo_v1.SettingsClient
}

func (h *singleClusterGraphqlHandler) GetGraphqlApi(ctx context.Context, request *rpc_edge_v1.GetGraphqlApiRequest) (*rpc_edge_v1.GetGraphqlApiResponse, error) {
	if request.GetGraphqlApiRef() == nil {
		return nil, eris.Errorf("graphqlapi ref missing from request: %v", request)
	}

	graphqlApi, err := h.graphqlClientset.GraphQLApis().GetGraphQLApi(ctx, client.ObjectKey{
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
		GraphqlApi: &rpc_edge_v1.GraphqlApi{
			Metadata: apiserverutils.ToMetadata(graphqlApi.ObjectMeta),
			Spec:     &graphqlApi.Spec,
			Status:   &graphqlApi.Status,
			GlooInstance: &skv2_v1.ObjectRef{
				Name:      glooInstance.GetMetadata().GetName(),
				Namespace: glooInstance.GetMetadata().GetNamespace(),
			},
		},
	}, nil
}

// find a gloo instance that is either watching the given namespace or watching all namespaces
func (h *singleClusterGraphqlHandler) getGlooInstanceForNamespace(ctx context.Context, namespace string) (*rpc_edge_v1.GlooInstance, error) {
	glooInstances, err := h.glooInstanceLister.ListGlooInstances(ctx)
	if err != nil {
		return nil, err
	}

	for _, instance := range glooInstances {
		watchedNamespaces := instance.Spec.GetControlPlane().GetWatchedNamespaces()
		if len(watchedNamespaces) == 0 {
			return instance, nil
		}
		for _, ns := range watchedNamespaces {
			if ns == namespace {
				return instance, nil
			}
		}
	}
	return nil, eris.Errorf("could not find a gloo instance with watched namespace %s", namespace)
}

func (h *singleClusterGraphqlHandler) ListGraphqlApis(ctx context.Context, request *rpc_edge_v1.ListGraphqlApisRequest) (*rpc_edge_v1.ListGraphqlApisResponse, error) {
	var rpcGraphqlApis []*rpc_edge_v1.GraphqlApi
	if request.GetGlooInstanceRef() == nil || request.GetGlooInstanceRef().GetName() == "" || request.GetGlooInstanceRef().GetNamespace() == "" {
		// List graphqlapis across all gloo edge instances
		instanceList, err := h.glooInstanceLister.ListGlooInstances(ctx)
		if err != nil {
			wrapped := eris.Wrapf(err, "Failed to list gloo edge instances")
			contextutils.LoggerFrom(ctx).Errorw(wrapped.Error(), zap.Error(err), zap.Any("request", request))
			return nil, wrapped
		}
		for _, instance := range instanceList {
			rpcGraphqlApiList, err := h.listGraphqlApisForGlooInstance(ctx, instance)
			if err != nil {
				wrapped := eris.Wrapf(err, "Failed to list graphqlApis for gloo edge instance %v", instance)
				contextutils.LoggerFrom(ctx).Errorw(wrapped.Error(), zap.Error(err), zap.Any("request", request))
				return nil, wrapped
			}
			rpcGraphqlApis = append(rpcGraphqlApis, rpcGraphqlApiList...)
		}
	} else {
		// List graphqlApis for a specific gloo edge instance
		instance, err := h.glooInstanceLister.GetGlooInstance(ctx, request.GetGlooInstanceRef())
		if err != nil {
			wrapped := eris.Wrap(err, "Failed to get gloo edge instance")
			contextutils.LoggerFrom(ctx).Errorw(wrapped.Error(), zap.Error(err), zap.Any("request", request))
			return nil, wrapped
		}
		rpcGraphqlApis, err = h.listGraphqlApisForGlooInstance(ctx, instance)
		if err != nil {
			wrapped := eris.Wrapf(err, "Failed to list graphqlApis for gloo edge instance %v", instance)
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

func (h *singleClusterGraphqlHandler) listGraphqlApisForGlooInstance(ctx context.Context, instance *rpc_edge_v1.GlooInstance) ([]*rpc_edge_v1.GraphqlApi, error) {
	var graphqlApiList []*graphql_v1beta1.GraphQLApi
	watchedNamespaces := instance.Spec.GetControlPlane().GetWatchedNamespaces()
	if len(watchedNamespaces) > 0 {
		for _, ns := range watchedNamespaces {
			list, err := h.graphqlClientset.GraphQLApis().ListGraphQLApi(ctx, client.InNamespace(ns))
			if err != nil {
				return nil, err
			}
			for _, item := range list.Items {
				item := item
				graphqlApiList = append(graphqlApiList, &item)
			}
		}
	} else {
		list, err := h.graphqlClientset.GraphQLApis().ListGraphQLApi(ctx)
		if err != nil {
			return nil, err
		}
		for _, item := range list.Items {
			item := item
			graphqlApiList = append(graphqlApiList, &item)
		}
	}
	sort.Slice(graphqlApiList, func(i, j int) bool {
		x := graphqlApiList[i]
		y := graphqlApiList[j]
		return x.GetNamespace()+x.GetName() < y.GetNamespace()+y.GetName()
	})

	var rpcGraphqlApis []*rpc_edge_v1.GraphqlApi
	glooInstanceRef := &skv2_v1.ObjectRef{
		Name:      instance.GetMetadata().GetName(),
		Namespace: instance.GetMetadata().GetNamespace(),
	}
	for _, graphqlApi := range graphqlApiList {
		graphqlApi := graphqlApi
		rpcGraphqlApis = append(rpcGraphqlApis, &rpc_edge_v1.GraphqlApi{
			Metadata:     apiserverutils.ToMetadata(graphqlApi.ObjectMeta),
			GlooInstance: glooInstanceRef,
			Spec:         &graphqlApi.Spec,
			Status:       &graphqlApi.Status,
		})
	}
	return rpcGraphqlApis, nil
}

func (h *singleClusterGraphqlHandler) GetGraphqlApiYaml(ctx context.Context, request *rpc_edge_v1.GetGraphqlApiYamlRequest) (*rpc_edge_v1.GetGraphqlApiYamlResponse, error) {
	if request.GetGraphqlApiRef() == nil {
		return nil, eris.Errorf("graphqlapi ref missing from request: %v", request)
	}

	graphqlApi, err := h.graphqlClientset.GraphQLApis().GetGraphQLApi(ctx, client.ObjectKey{
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

func (h *singleClusterGraphqlHandler) CreateGraphqlApi(ctx context.Context, request *rpc_edge_v1.CreateGraphqlApiRequest) (*rpc_edge_v1.CreateGraphqlApiResponse, error) {
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

	graphqlApi := &graphql_v1beta1.GraphQLApi{
		ObjectMeta: apiserverutils.RefToObjectMeta(*request.GetGraphqlApiRef()),
		Spec:       *request.GetSpec(),
	}
	err = h.graphqlClientset.GraphQLApis().CreateGraphQLApi(ctx, graphqlApi)
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
		GraphqlApi: &rpc_edge_v1.GraphqlApi{
			Metadata: apiserverutils.ToMetadata(graphqlApi.ObjectMeta),
			Spec:     &graphqlApi.Spec,
			Status:   &graphqlApi.Status,
			GlooInstance: &skv2_v1.ObjectRef{
				Name:      glooInstance.GetMetadata().GetName(),
				Namespace: glooInstance.GetMetadata().GetNamespace(),
			},
		},
	}, nil
}

func (h *singleClusterGraphqlHandler) UpdateGraphqlApi(ctx context.Context, request *rpc_edge_v1.UpdateGraphqlApiRequest) (*rpc_edge_v1.UpdateGraphqlApiResponse, error) {
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

	// first get the existing graphqlapi
	graphqlApi, err := h.graphqlClientset.GraphQLApis().GetGraphQLApi(ctx, client.ObjectKey{
		Namespace: request.GetGraphqlApiRef().GetNamespace(),
		Name:      request.GetGraphqlApiRef().GetName(),
	})
	if err != nil {
		wrapped := eris.Wrapf(err, "Cannot edit a graphqlapi that does not exist")
		contextutils.LoggerFrom(ctx).Errorw(wrapped.Error(), zap.Error(err), zap.Any("request", request))
		return nil, wrapped
	}
	// apply the changes to its spec
	graphqlApi.Spec = *request.GetSpec()

	// save the updated graphqlapi
	err = h.graphqlClientset.GraphQLApis().UpdateGraphQLApi(ctx, graphqlApi)
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
		GraphqlApi: &rpc_edge_v1.GraphqlApi{
			Metadata: apiserverutils.ToMetadata(graphqlApi.ObjectMeta),
			Spec:     &graphqlApi.Spec,
			Status:   &graphqlApi.Status,
			GlooInstance: &skv2_v1.ObjectRef{
				Name:      glooInstance.GetMetadata().GetName(),
				Namespace: glooInstance.GetMetadata().GetNamespace(),
			},
		},
	}, nil
}

func (h *singleClusterGraphqlHandler) DeleteGraphqlApi(ctx context.Context, request *rpc_edge_v1.DeleteGraphqlApiRequest) (*rpc_edge_v1.DeleteGraphqlApiResponse, error) {
	err := apiserverutils.CheckUpdatesAllowed(ctx, h.settingsClient)
	if err != nil {
		return nil, err
	}
	err = h.checkGraphqlApiRef(request.GetGraphqlApiRef())
	if err != nil {
		return nil, err
	}

	err = h.graphqlClientset.GraphQLApis().DeleteGraphQLApi(ctx, client.ObjectKey{
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

func (h *singleClusterGraphqlHandler) ValidateResolverYaml(_ context.Context, request *rpc_edge_v1.ValidateResolverYamlRequest) (*rpc_edge_v1.ValidateResolverYamlResponse, error) {
	err := ValidateResolverYaml(request.GetYaml(), request.GetResolverType())
	if err != nil {
		return nil, err
	}
	return &rpc_edge_v1.ValidateResolverYamlResponse{}, nil
}

func (h *singleClusterGraphqlHandler) ValidateSchemaDefinition(ctx context.Context, request *rpc_edge_v1.ValidateSchemaDefinitionRequest) (*rpc_edge_v1.ValidateSchemaDefinitionResponse, error) {
	err := ValidateSchemaDefinition(request)
	if err != nil {
		return nil, err
	}
	return &rpc_edge_v1.ValidateSchemaDefinitionResponse{}, nil
}

func (h *singleClusterGraphqlHandler) checkGraphqlApiRef(ref *skv2_v1.ClusterObjectRef) error {
	if ref == nil || ref.GetName() == "" || ref.GetNamespace() == "" {
		return eris.Errorf("request does not contain valid graphqlapi ref: %v", ref)
	}
	return nil
}

func (h *singleClusterGraphqlHandler) GetStitchedSchemaDefinition(ctx context.Context, request *rpc_edge_v1.GetStitchedSchemaDefinitionRequest) (*rpc_edge_v1.GetStitchedSchemaDefinitionResponse, error) {
	return GetStitchedSchemaDefinition(ctx, h.graphqlClientset.GraphQLApis(), request)
}

func (h *singleClusterGraphqlHandler) GetSchemaDiff(ctx context.Context, request *rpc_edge_v1.GetSchemaDiffRequest) (*rpc_edge_v1.GetSchemaDiffResponse, error) {
	return GetSchemaDiff(request)
}
