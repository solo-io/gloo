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
	"github.com/solo-io/solo-projects/projects/apiserver/server/services/glooinstance_handler"
	"go.uber.org/zap"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func NewSingleClusterGraphqlHandler(
	graphqlClientset graphql_v1alpha1.Clientset,
	glooInstanceLister glooinstance_handler.SingleClusterGlooInstanceLister,
) rpc_edge_v1.GraphqlApiServer {
	return &singleClusterGraphqlHandler{
		graphqlClientset:   graphqlClientset,
		glooInstanceLister: glooInstanceLister,
	}
}

type singleClusterGraphqlHandler struct {
	graphqlClientset   graphql_v1alpha1.Clientset
	glooInstanceLister glooinstance_handler.SingleClusterGlooInstanceLister
}

func (h *singleClusterGraphqlHandler) GetGraphqlSchema(ctx context.Context, request *rpc_edge_v1.GetGraphqlSchemaRequest) (*rpc_edge_v1.GetGraphqlSchemaResponse, error) {
	if request.GetGraphqlSchemaRef() == nil {
		return nil, eris.Errorf("graphqlschema ref missing from request: %v", request)
	}

	graphqlSchema, err := h.graphqlClientset.GraphQLSchemas().GetGraphQLSchema(ctx, client.ObjectKey{
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

func (h *singleClusterGraphqlHandler) ListGraphqlSchemas(ctx context.Context, request *rpc_edge_v1.ListGraphqlSchemasRequest) (*rpc_edge_v1.ListGraphqlSchemasResponse, error) {
	var rpcGraphqlSchemas []*rpc_edge_v1.GraphqlSchema
	if request.GetGlooInstanceRef() == nil || request.GetGlooInstanceRef().GetName() == "" || request.GetGlooInstanceRef().GetNamespace() == "" {
		// List graphqlschemas across all gloo edge instances
		instanceList, err := h.glooInstanceLister.ListGlooInstances(ctx)
		if err != nil {
			wrapped := eris.Wrapf(err, "Failed to list gloo edge instances")
			contextutils.LoggerFrom(ctx).Errorw(wrapped.Error(), zap.Error(err), zap.Any("request", request))
			return nil, wrapped
		}
		for _, instance := range instanceList {
			rpcGraphqlSchemaList, err := h.listGraphqlSchemasForGlooInstance(ctx, instance)
			if err != nil {
				wrapped := eris.Wrapf(err, "Failed to list graphqlSchemas for gloo edge instance %v", instance)
				contextutils.LoggerFrom(ctx).Errorw(wrapped.Error(), zap.Error(err), zap.Any("request", request))
				return nil, wrapped
			}
			rpcGraphqlSchemas = append(rpcGraphqlSchemas, rpcGraphqlSchemaList...)
		}
	} else {
		// List graphqlSchemas for a specific gloo edge instance
		instance, err := h.glooInstanceLister.GetGlooInstance(ctx, request.GetGlooInstanceRef())
		if err != nil {
			wrapped := eris.Wrap(err, "Failed to get gloo edge instance")
			contextutils.LoggerFrom(ctx).Errorw(wrapped.Error(), zap.Error(err), zap.Any("request", request))
			return nil, wrapped
		}
		rpcGraphqlSchemas, err = h.listGraphqlSchemasForGlooInstance(ctx, instance)
		if err != nil {
			wrapped := eris.Wrapf(err, "Failed to list graphqlSchemas for gloo edge instance %v", instance)
			contextutils.LoggerFrom(ctx).Errorw(wrapped.Error(), zap.Error(err), zap.Any("request", request))
			return nil, wrapped
		}
	}

	return &rpc_edge_v1.ListGraphqlSchemasResponse{
		GraphqlSchemas: rpcGraphqlSchemas,
	}, nil
}

func (h *singleClusterGraphqlHandler) listGraphqlSchemasForGlooInstance(ctx context.Context, instance *rpc_edge_v1.GlooInstance) ([]*rpc_edge_v1.GraphqlSchema, error) {
	var graphqlSchemaList []*graphql_v1alpha1.GraphQLSchema
	watchedNamespaces := instance.Spec.GetControlPlane().GetWatchedNamespaces()
	if len(watchedNamespaces) > 0 {
		for _, ns := range watchedNamespaces {
			list, err := h.graphqlClientset.GraphQLSchemas().ListGraphQLSchema(ctx, client.InNamespace(ns))
			if err != nil {
				return nil, err
			}
			for _, item := range list.Items {
				item := item
				graphqlSchemaList = append(graphqlSchemaList, &item)
			}
		}
	} else {
		list, err := h.graphqlClientset.GraphQLSchemas().ListGraphQLSchema(ctx)
		if err != nil {
			return nil, err
		}
		for _, item := range list.Items {
			item := item
			graphqlSchemaList = append(graphqlSchemaList, &item)
		}
	}
	sort.Slice(graphqlSchemaList, func(i, j int) bool {
		x := graphqlSchemaList[i]
		y := graphqlSchemaList[j]
		return x.GetNamespace()+x.GetName() < y.GetNamespace()+y.GetName()
	})

	var rpcGraphqlSchemas []*rpc_edge_v1.GraphqlSchema
	glooInstanceRef := &skv2_v1.ObjectRef{
		Name:      instance.GetMetadata().GetName(),
		Namespace: instance.GetMetadata().GetNamespace(),
	}
	for _, graphqlSchema := range graphqlSchemaList {
		graphqlSchema := graphqlSchema
		rpcGraphqlSchemas = append(rpcGraphqlSchemas, &rpc_edge_v1.GraphqlSchema{
			Metadata:     apiserverutils.ToMetadata(graphqlSchema.ObjectMeta),
			GlooInstance: glooInstanceRef,
			Spec:         &graphqlSchema.Spec,
			Status:       &graphqlSchema.Status,
		})
	}
	return rpcGraphqlSchemas, nil
}

func (h *singleClusterGraphqlHandler) GetGraphqlSchemaYaml(ctx context.Context, request *rpc_edge_v1.GetGraphqlSchemaYamlRequest) (*rpc_edge_v1.GetGraphqlSchemaYamlResponse, error) {
	if request.GetGraphqlSchemaRef() == nil {
		return nil, eris.Errorf("graphqlschema ref missing from request: %v", request)
	}

	graphqlSchema, err := h.graphqlClientset.GraphQLSchemas().GetGraphQLSchema(ctx, client.ObjectKey{
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

func (h *singleClusterGraphqlHandler) CreateGraphqlSchema(ctx context.Context, request *rpc_edge_v1.CreateGraphqlSchemaRequest) (*rpc_edge_v1.CreateGraphqlSchemaResponse, error) {
	err := h.checkGraphqlSchemaRef(request.GetGraphqlSchemaRef())
	if err != nil {
		return nil, err
	}
	if request.GetSpec() == nil {
		return nil, eris.Errorf("graphqlschema spec missing from request: %v", request)
	}

	graphqlSchema := &graphql_v1alpha1.GraphQLSchema{
		ObjectMeta: apiserverutils.RefToObjectMeta(*request.GetGraphqlSchemaRef()),
		Spec:       *request.GetSpec(),
	}
	err = h.graphqlClientset.GraphQLSchemas().CreateGraphQLSchema(ctx, graphqlSchema)
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
				Name:      glooInstance.GetMetadata().GetName(),
				Namespace: glooInstance.GetMetadata().GetNamespace(),
			},
		},
	}, nil
}

func (h *singleClusterGraphqlHandler) UpdateGraphqlSchema(ctx context.Context, request *rpc_edge_v1.UpdateGraphqlSchemaRequest) (*rpc_edge_v1.UpdateGraphqlSchemaResponse, error) {
	err := h.checkGraphqlSchemaRef(request.GetGraphqlSchemaRef())
	if err != nil {
		return nil, err
	}
	if request.GetSpec() == nil {
		return nil, eris.Errorf("graphqlschema spec missing from request: %v", request)
	}

	// first get the existing graphqlschema
	graphqlSchema, err := h.graphqlClientset.GraphQLSchemas().GetGraphQLSchema(ctx, client.ObjectKey{
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
	err = h.graphqlClientset.GraphQLSchemas().UpdateGraphQLSchema(ctx, graphqlSchema)
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
				Name:      glooInstance.GetMetadata().GetName(),
				Namespace: glooInstance.GetMetadata().GetNamespace(),
			},
		},
	}, nil
}

func (h *singleClusterGraphqlHandler) DeleteGraphqlSchema(ctx context.Context, request *rpc_edge_v1.DeleteGraphqlSchemaRequest) (*rpc_edge_v1.DeleteGraphqlSchemaResponse, error) {
	err := h.checkGraphqlSchemaRef(request.GetGraphqlSchemaRef())
	if err != nil {
		return nil, err
	}

	err = h.graphqlClientset.GraphQLSchemas().DeleteGraphQLSchema(ctx, client.ObjectKey{
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

func (h *singleClusterGraphqlHandler) ValidateResolverYaml(ctx context.Context, request *rpc_edge_v1.ValidateResolverYamlRequest) (*rpc_edge_v1.ValidateResolverYamlResponse, error) {
	// TODO implement
	return &rpc_edge_v1.ValidateResolverYamlResponse{}, nil
}

func (h *singleClusterGraphqlHandler) checkGraphqlSchemaRef(ref *skv2_v1.ClusterObjectRef) error {
	if ref == nil || ref.GetName() == "" || ref.GetNamespace() == "" {
		return eris.Errorf("request does not contain valid graphqlschema ref: %v", ref)
	}
	return nil
}
