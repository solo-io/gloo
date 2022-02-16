package graphql_handler

import (
	"context"
	"sort"

	"github.com/ghodss/yaml"
	"github.com/rotisserie/eris"
	"github.com/solo-io/go-utils/contextutils"
	skv2_v1 "github.com/solo-io/skv2/pkg/api/core.skv2.solo.io/v1"
	graphql_v1alpha1 "github.com/solo-io/solo-apis/pkg/api/graphql.gloo.solo.io/v1alpha1"
	rpc_edge_v1 "github.com/solo-io/solo-projects/projects/apiserver/pkg/api/rpc.edge.gloo/v1"
	"github.com/solo-io/solo-projects/projects/apiserver/server/apiserverutils"
	fed_v1 "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.solo.io/v1"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func NewFedGraphqlHandler(
	glooInstanceClient fed_v1.GlooInstanceClient,
	graphqlMCClientset graphql_v1alpha1.MulticlusterClientset,
) rpc_edge_v1.GraphqlApiServer {
	return &fedGraphqlHandler{
		glooInstanceClient: glooInstanceClient,
		graphqlMCClientset: graphqlMCClientset,
	}
}

type fedGraphqlHandler struct {
	glooInstanceClient fed_v1.GlooInstanceClient
	graphqlMCClientset graphql_v1alpha1.MulticlusterClientset
}

func (h *fedGraphqlHandler) GetGraphqlSchema(ctx context.Context, request *rpc_edge_v1.GetGraphqlSchemaRequest) (*rpc_edge_v1.GetGraphqlSchemaResponse, error) {
	if request.GetGraphqlSchemaRef() == nil {
		return nil, eris.Errorf("graphqlschema ref missing from request: %v", request)
	}

	clientset, err := h.graphqlMCClientset.Cluster(request.GetGraphqlSchemaRef().GetClusterName())
	if err != nil {
		wrapped := eris.Wrapf(err, "Failed to get graphql client set")
		contextutils.LoggerFrom(ctx).Errorw(wrapped.Error(), zap.Error(err), zap.Any("request", request))
		return nil, wrapped
	}

	graphqlSchema, err := clientset.GraphQLSchemas().GetGraphQLSchema(ctx, client.ObjectKey{
		Namespace: request.GetGraphqlSchemaRef().GetNamespace(),
		Name:      request.GetGraphqlSchemaRef().GetName(),
	})
	if err != nil {
		wrapped := eris.Wrapf(err, "Failed to get graphqlschema")
		contextutils.LoggerFrom(ctx).Errorw(wrapped.Error(), zap.Error(err), zap.Any("request", request))
		return nil, wrapped
	}

	// find which gloo instance this graphql schema belongs to, by finding a gloo instance that is watching
	// the graphql schema's namespace
	glooInstance, err := h.getGlooInstanceForNamespace(ctx, request.GetGraphqlSchemaRef().GetNamespace())
	if err != nil {
		return nil, err
	}
	return &rpc_edge_v1.GetGraphqlSchemaResponse{
		GraphqlSchema: &rpc_edge_v1.GraphqlSchema{
			Metadata: apiserverutils.ToMetadata(graphqlSchema.ObjectMeta),
			Spec:     &graphqlSchema.Spec,
			Status:   &graphqlSchema.Status,
			GlooInstance: &skv2_v1.ObjectRef{
				Name:      glooInstance.GetName(),
				Namespace: glooInstance.GetNamespace(),
			},
		},
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

func (h *fedGraphqlHandler) ListGraphqlSchemas(ctx context.Context, request *rpc_edge_v1.ListGraphqlSchemasRequest) (*rpc_edge_v1.ListGraphqlSchemasResponse, error) {
	var rpcGraphqlSchemas []*rpc_edge_v1.GraphqlSchema
	if request.GetGlooInstanceRef() == nil || request.GetGlooInstanceRef().GetName() == "" || request.GetGlooInstanceRef().GetNamespace() == "" {
		// List graphqlSchemas across all gloo edge instances
		instanceList, err := h.glooInstanceClient.ListGlooInstance(ctx)
		if err != nil {
			wrapped := eris.Wrapf(err, "Failed to list gloo edge instances")
			contextutils.LoggerFrom(ctx).Errorw(wrapped.Error(), zap.Error(err), zap.Any("request", request))
			return nil, wrapped
		}
		for _, instance := range instanceList.Items {
			rpcGraphqlSchemaList, err := h.listGraphqlSchemasForGlooInstance(ctx, &instance)
			if err != nil {
				wrapped := eris.Wrapf(err, "Failed to list graphqlSchemas for gloo edge instance %s.%s", instance.GetNamespace(), instance.GetName())
				contextutils.LoggerFrom(ctx).Errorw(wrapped.Error(), zap.Error(err), zap.Any("request", request))
				return nil, wrapped
			}
			rpcGraphqlSchemas = append(rpcGraphqlSchemas, rpcGraphqlSchemaList...)
		}
	} else {
		// List graphqlSchemas for a specific gloo edge instance
		instance, err := h.glooInstanceClient.GetGlooInstance(ctx, types.NamespacedName{
			Name:      request.GetGlooInstanceRef().GetName(),
			Namespace: request.GetGlooInstanceRef().GetNamespace(),
		})
		if err != nil {
			wrapped := eris.Wrapf(err, "Failed to get gloo edge instance %s.%s", instance.GetNamespace(), instance.GetName())
			contextutils.LoggerFrom(ctx).Errorw(wrapped.Error(), zap.Error(err), zap.Any("request", request))
			return nil, wrapped
		}
		rpcGraphqlSchemas, err = h.listGraphqlSchemasForGlooInstance(ctx, instance)
		if err != nil {
			wrapped := eris.Wrapf(err, "Failed to list graphqlSchemas for gloo edge instance %s.%s", instance.GetNamespace(), instance.GetName())
			contextutils.LoggerFrom(ctx).Errorw(wrapped.Error(), zap.Error(err), zap.Any("request", request))
			return nil, wrapped
		}
	}

	return &rpc_edge_v1.ListGraphqlSchemasResponse{
		GraphqlSchemas: rpcGraphqlSchemas,
	}, nil
}

func (h *fedGraphqlHandler) listGraphqlSchemasForGlooInstance(ctx context.Context, instance *fed_v1.GlooInstance) ([]*rpc_edge_v1.GraphqlSchema, error) {
	clientset, err := h.graphqlMCClientset.Cluster(instance.Spec.GetCluster())
	if err != nil {
		return nil, err
	}
	graphqlSchemaClient := clientset.GraphQLSchemas()

	var graphqlSchemaGraphqlSchemaList []*graphql_v1alpha1.GraphQLSchema
	watchedNamespaces := instance.Spec.GetControlPlane().GetWatchedNamespaces()
	if len(watchedNamespaces) != 0 {
		for _, ns := range watchedNamespaces {
			list, err := graphqlSchemaClient.ListGraphQLSchema(ctx, client.InNamespace(ns))
			if err != nil {
				return nil, err
			}
			for i, _ := range list.Items {
				graphqlSchemaGraphqlSchemaList = append(graphqlSchemaGraphqlSchemaList, &list.Items[i])
			}
		}
	} else {
		list, err := graphqlSchemaClient.ListGraphQLSchema(ctx)
		if err != nil {
			return nil, err
		}
		for i, _ := range list.Items {
			graphqlSchemaGraphqlSchemaList = append(graphqlSchemaGraphqlSchemaList, &list.Items[i])
		}
	}
	sort.Slice(graphqlSchemaGraphqlSchemaList, func(i, j int) bool {
		x := graphqlSchemaGraphqlSchemaList[i]
		y := graphqlSchemaGraphqlSchemaList[j]
		return x.GetNamespace()+x.GetName() < y.GetNamespace()+y.GetName()
	})

	var rpcGraphqlSchemas []*rpc_edge_v1.GraphqlSchema
	for _, graphqlSchema := range graphqlSchemaGraphqlSchemaList {
		rpcGraphqlSchemas = append(rpcGraphqlSchemas, convertGraphqlSchema(graphqlSchema, &skv2_v1.ObjectRef{
			Name:      instance.GetName(),
			Namespace: instance.GetNamespace(),
		}, instance.Spec.GetCluster()))
	}
	return rpcGraphqlSchemas, nil
}

func convertGraphqlSchema(graphqlSchema *graphql_v1alpha1.GraphQLSchema, glooInstance *skv2_v1.ObjectRef, cluster string) *rpc_edge_v1.GraphqlSchema {
	m := &rpc_edge_v1.GraphqlSchema{
		Metadata:     apiserverutils.ToMetadata(graphqlSchema.ObjectMeta),
		GlooInstance: glooInstance,
		Spec:         &graphqlSchema.Spec,
		Status:       &graphqlSchema.Status,
	}
	m.Metadata.ClusterName = cluster
	return m
}

func (h *fedGraphqlHandler) GetGraphqlSchemaYaml(ctx context.Context, request *rpc_edge_v1.GetGraphqlSchemaYamlRequest) (*rpc_edge_v1.GetGraphqlSchemaYamlResponse, error) {
	if request.GetGraphqlSchemaRef() == nil {
		return nil, eris.Errorf("graphqlschema ref missing from request: %v", request)
	}

	clientset, err := h.graphqlMCClientset.Cluster(request.GetGraphqlSchemaRef().GetClusterName())
	if err != nil {
		wrapped := eris.Wrapf(err, "Failed to get graphql client set")
		contextutils.LoggerFrom(ctx).Errorw(wrapped.Error(), zap.Error(err), zap.Any("request", request))
		return nil, wrapped
	}
	graphqlSchema, err := clientset.GraphQLSchemas().GetGraphQLSchema(ctx, client.ObjectKey{
		Namespace: request.GetGraphqlSchemaRef().GetNamespace(),
		Name:      request.GetGraphqlSchemaRef().GetName(),
	})
	if err != nil {
		wrapped := eris.Wrapf(err, "Failed to get graphqlschema")
		contextutils.LoggerFrom(ctx).Errorw(wrapped.Error(), zap.Error(err), zap.Any("request", request))
		return nil, wrapped
	}
	content, err := yaml.Marshal(graphqlSchema)
	if err != nil {
		wrapped := eris.Wrapf(err, "Failed to marshal kube resource into yaml")
		contextutils.LoggerFrom(ctx).Errorw(wrapped.Error(), zap.Error(err), zap.Any("request", request))
		return nil, wrapped
	}
	return &rpc_edge_v1.GetGraphqlSchemaYamlResponse{
		YamlData: &rpc_edge_v1.ResourceYaml{
			Yaml: string(content),
		},
	}, nil
}

func (h *fedGraphqlHandler) CreateGraphqlSchema(ctx context.Context, request *rpc_edge_v1.CreateGraphqlSchemaRequest) (*rpc_edge_v1.CreateGraphqlSchemaResponse, error) {
	err := h.checkGraphqlSchemaRef(request.GetGraphqlSchemaRef())
	if err != nil {
		return nil, err
	}
	if request.GetSpec() == nil {
		return nil, eris.Errorf("graphqlschema spec missing from request: %v", request)
	}

	clientset, err := h.graphqlMCClientset.Cluster(request.GetGraphqlSchemaRef().GetClusterName())
	if err != nil {
		wrapped := eris.Wrapf(err, "Failed to get graphql client set")
		contextutils.LoggerFrom(ctx).Errorw(wrapped.Error(), zap.Error(err), zap.Any("request", request))
		return nil, wrapped
	}

	graphqlSchema := &graphql_v1alpha1.GraphQLSchema{
		ObjectMeta: apiserverutils.RefToObjectMeta(*request.GetGraphqlSchemaRef()),
		Spec:       *request.GetSpec(),
	}
	err = clientset.GraphQLSchemas().CreateGraphQLSchema(ctx, graphqlSchema)
	if err != nil {
		wrapped := eris.Wrapf(err, "Failed to create graphqlschema")
		contextutils.LoggerFrom(ctx).Errorw(wrapped.Error(), zap.Error(err), zap.Any("request", request))
		return nil, wrapped
	}

	// find which gloo instance this graphql schema belongs to, by finding a gloo instance that is watching
	// the graphql schema's namespace
	glooInstance, err := h.getGlooInstanceForNamespace(ctx, request.GetGraphqlSchemaRef().GetNamespace())
	if err != nil {
		return nil, err
	}
	return &rpc_edge_v1.CreateGraphqlSchemaResponse{
		GraphqlSchema: &rpc_edge_v1.GraphqlSchema{
			Metadata: apiserverutils.ToMetadata(graphqlSchema.ObjectMeta),
			Spec:     &graphqlSchema.Spec,
			Status:   &graphqlSchema.Status,
			GlooInstance: &skv2_v1.ObjectRef{
				Name:      glooInstance.GetName(),
				Namespace: glooInstance.GetNamespace(),
			},
		},
	}, nil
}

func (h *fedGraphqlHandler) UpdateGraphqlSchema(ctx context.Context, request *rpc_edge_v1.UpdateGraphqlSchemaRequest) (*rpc_edge_v1.UpdateGraphqlSchemaResponse, error) {
	err := h.checkGraphqlSchemaRef(request.GetGraphqlSchemaRef())
	if err != nil {
		return nil, err
	}
	if request.GetSpec() == nil {
		return nil, eris.Errorf("graphqlschema spec missing from request: %v", request)
	}

	clientset, err := h.graphqlMCClientset.Cluster(request.GetGraphqlSchemaRef().GetClusterName())
	if err != nil {
		wrapped := eris.Wrapf(err, "Failed to get graphql client set")
		contextutils.LoggerFrom(ctx).Errorw(wrapped.Error(), zap.Error(err), zap.Any("request", request))
		return nil, wrapped
	}

	// first get the existing graphqlschema
	graphqlSchema, err := clientset.GraphQLSchemas().GetGraphQLSchema(ctx, client.ObjectKey{
		Namespace: request.GetGraphqlSchemaRef().GetNamespace(),
		Name:      request.GetGraphqlSchemaRef().GetName(),
	})
	if err != nil {
		wrapped := eris.Wrapf(err, "Cannot edit a graphqlschema that does not exist")
		contextutils.LoggerFrom(ctx).Errorw(wrapped.Error(), zap.Error(err), zap.Any("request", request))
		return nil, wrapped
	}
	// apply the changes to its spec
	graphqlSchema.Spec = *request.GetSpec()

	// save the updated graphqlschema
	err = clientset.GraphQLSchemas().UpdateGraphQLSchema(ctx, graphqlSchema)
	if err != nil {
		wrapped := eris.Wrapf(err, "Failed to update graphqlschema")
		contextutils.LoggerFrom(ctx).Errorw(wrapped.Error(), zap.Error(err), zap.Any("request", request))
		return nil, wrapped
	}

	// find which gloo instance this graphql schema belongs to, by finding a gloo instance that is watching
	// the graphql schema's namespace
	glooInstance, err := h.getGlooInstanceForNamespace(ctx, request.GetGraphqlSchemaRef().GetNamespace())
	if err != nil {
		return nil, err
	}
	return &rpc_edge_v1.UpdateGraphqlSchemaResponse{
		GraphqlSchema: &rpc_edge_v1.GraphqlSchema{
			Metadata: apiserverutils.ToMetadata(graphqlSchema.ObjectMeta),
			Spec:     &graphqlSchema.Spec,
			Status:   &graphqlSchema.Status,
			GlooInstance: &skv2_v1.ObjectRef{
				Name:      glooInstance.GetName(),
				Namespace: glooInstance.GetNamespace(),
			},
		},
	}, nil
}

func (h *fedGraphqlHandler) DeleteGraphqlSchema(ctx context.Context, request *rpc_edge_v1.DeleteGraphqlSchemaRequest) (*rpc_edge_v1.DeleteGraphqlSchemaResponse, error) {
	err := h.checkGraphqlSchemaRef(request.GetGraphqlSchemaRef())
	if err != nil {
		return nil, err
	}

	clientset, err := h.graphqlMCClientset.Cluster(request.GetGraphqlSchemaRef().GetClusterName())
	if err != nil {
		wrapped := eris.Wrapf(err, "Failed to get graphql client set")
		contextutils.LoggerFrom(ctx).Errorw(wrapped.Error(), zap.Error(err), zap.Any("request", request))
		return nil, wrapped
	}

	err = clientset.GraphQLSchemas().DeleteGraphQLSchema(ctx, client.ObjectKey{
		Name:      request.GetGraphqlSchemaRef().GetName(),
		Namespace: request.GetGraphqlSchemaRef().GetNamespace(),
	})
	if err != nil {
		wrapped := eris.Wrapf(err, "Failed to delete graphqlschema")
		contextutils.LoggerFrom(ctx).Errorw(wrapped.Error(), zap.Error(err), zap.Any("request", request))
		return nil, wrapped
	}

	return &rpc_edge_v1.DeleteGraphqlSchemaResponse{
		GraphqlSchemaRef: request.GetGraphqlSchemaRef(),
	}, nil
}

func (h *fedGraphqlHandler) ValidateResolverYaml(_ context.Context, request *rpc_edge_v1.ValidateResolverYamlRequest) (*rpc_edge_v1.ValidateResolverYamlResponse, error) {
	err := ValidateResolverYaml(request.GetYaml(), request.GetResolverType())
	if err != nil {
		return nil, err
	}
	return &rpc_edge_v1.ValidateResolverYamlResponse{}, nil
}

func (h *fedGraphqlHandler) checkGraphqlSchemaRef(ref *skv2_v1.ClusterObjectRef) error {
	if ref == nil || ref.GetName() == "" || ref.GetNamespace() == "" || ref.GetClusterName() == "" {
		return eris.Errorf("request does not contain valid graphqlschema ref: %v", ref)
	}
	return nil
}
