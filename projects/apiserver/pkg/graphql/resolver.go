package graphql

import (
	"context"

	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/errors"
	"github.com/solo-io/solo-kit/projects/apiserver/pkg/graphql/customtypes"
	"github.com/solo-io/solo-kit/projects/apiserver/pkg/graphql/graph"
	"github.com/solo-io/solo-kit/projects/apiserver/pkg/graphql/models"
	gatewayv1 "github.com/solo-io/solo-kit/projects/gateway/pkg/api/v1"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1"
	sqoopv1 "github.com/solo-io/solo-kit/projects/sqoop/pkg/api/v1"
)

type ApiResolver struct {
	Upstreams       v1.UpstreamClient
	VirtualServices gatewayv1.VirtualServiceClient
	ResolverMaps    sqoopv1.ResolverMapClient
	Converter       *Converter
}

func NewResolvers(upstreams v1.UpstreamClient,
	virtualServices gatewayv1.VirtualServiceClient,
	resolverMaps sqoopv1.ResolverMapClient) *ApiResolver {
	return &ApiResolver{
		Upstreams:       upstreams,
		VirtualServices: virtualServices,
		ResolverMaps:    resolverMaps,
		// TODO(ilackarms): just make these private functions, remove converter
		Converter: &Converter{},
	}
}

func (r *ApiResolver) Mutation() graph.MutationResolver {
	return &mutationResolver{r}
}
func (r *ApiResolver) Query() graph.QueryResolver {
	return &queryResolver{r}
}
func (r *ApiResolver) UpstreamMutation() graph.UpstreamMutationResolver {
	return &upstreamMutationResolver{r}
}
func (r *ApiResolver) UpstreamQuery() graph.UpstreamQueryResolver {
	return &upstreamQueryResolver{r}
}
func (r *ApiResolver) VirtualServiceMutation() graph.VirtualServiceMutationResolver {
	return &virtualServiceMutationResolver{r}
}
func (r *ApiResolver) VirtualServiceQuery() graph.VirtualServiceQueryResolver {
	return &virtualServiceQueryResolver{r}
}
func (r *ApiResolver) ResolverMapMutation() graph.ResolverMapMutationResolver {
	return &resolverMapMutationResolver{r}
}
func (r *ApiResolver) ResolverMapQuery() graph.ResolverMapQueryResolver {
	return &resolverMapQueryResolver{r}
}

type mutationResolver struct{ *ApiResolver }

func (r *mutationResolver) Upstreams(ctx context.Context, namespace string) (customtypes.UpstreamMutation, error) {
	return customtypes.UpstreamMutation{Namespace: namespace}, nil
}
func (r *mutationResolver) VirtualServices(ctx context.Context, namespace string) (customtypes.VirtualServiceMutation, error) {
	return customtypes.VirtualServiceMutation{Namespace: namespace}, nil
}
func (r *mutationResolver) ResolverMaps(ctx context.Context, namespace string) (customtypes.ResolverMapMutation, error) {
	return customtypes.ResolverMapMutation{Namespace: namespace}, nil
}

type queryResolver struct{ *ApiResolver }

func (r *queryResolver) Upstreams(ctx context.Context, namespace string) (customtypes.UpstreamQuery, error) {
	return customtypes.UpstreamQuery{Namespace: namespace}, nil
}
func (r *queryResolver) VirtualServices(ctx context.Context, namespace string) (customtypes.VirtualServiceQuery, error) {
	return customtypes.VirtualServiceQuery{Namespace: namespace}, nil
}
func (r *queryResolver) ResolverMaps(ctx context.Context, namespace string) (customtypes.ResolverMapQuery, error) {
	return customtypes.ResolverMapQuery{Namespace: namespace}, nil
}

type upstreamMutationResolver struct{ *ApiResolver }

func (r *upstreamMutationResolver) write(overwrite bool, ctx context.Context, obj *customtypes.UpstreamMutation, upstream models.InputUpstream) (*models.Upstream, error) {
	ups, err := r.Converter.ConvertInputUpstream(upstream)
	if err != nil {
		return nil, err
	}
	out, err := r.Upstreams.Write(ups, clients.WriteOpts{
		Ctx:               ctx,
		OverwriteExisting: overwrite,
	})
	if err != nil {
		return nil, err
	}
	return r.Converter.ConvertOutputUpstream(out), nil
}

func (r *upstreamMutationResolver) Create(ctx context.Context, obj *customtypes.UpstreamMutation, upstream models.InputUpstream) (*models.Upstream, error) {
	return r.write(false, ctx, obj, upstream)
}
func (r *upstreamMutationResolver) Update(ctx context.Context, obj *customtypes.UpstreamMutation, upstream models.InputUpstream) (*models.Upstream, error) {
	return r.write(true, ctx, obj, upstream)
}
func (r *upstreamMutationResolver) Delete(ctx context.Context, obj *customtypes.UpstreamMutation, name string) (*models.Upstream, error) {
	upstream, err := r.Upstreams.Read(obj.Namespace, name, clients.ReadOpts{
		Ctx: ctx,
	})
	if err != nil {
		if errors.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	err = r.Upstreams.Delete(obj.Namespace, name, clients.DeleteOpts{Ctx: ctx})
	if err != nil {
		return nil, err
	}
	return r.Converter.ConvertOutputUpstream(upstream), nil
}

type upstreamQueryResolver struct{ *ApiResolver }

func (r *upstreamQueryResolver) List(ctx context.Context, obj *customtypes.UpstreamQuery, selector *models.InputMapStringString) ([]*models.Upstream, error) {
	var convertedSelector map[string]string
	if selector != nil {
		convertedSelector = selector.GoType()
	}
	list, err := r.Upstreams.List(obj.Namespace, clients.ListOpts{
		Ctx:      ctx,
		Selector: convertedSelector,
	})
	if err != nil {
		return nil, err
	}
	return r.Converter.ConvertOutputUpstreams(list), nil
}

func (r *upstreamQueryResolver) Get(ctx context.Context, obj *customtypes.UpstreamQuery, name string) (*models.Upstream, error) {
	upstream, err := r.Upstreams.Read(obj.Namespace, name, clients.ReadOpts{
		Ctx: ctx,
	})
	if err != nil {
		return nil, err
	}
	return r.Converter.ConvertOutputUpstream(upstream), nil
}

type virtualServiceMutationResolver struct{ *ApiResolver }

func (r *virtualServiceMutationResolver) write(overwrite bool, ctx context.Context, obj *customtypes.VirtualServiceMutation, virtualService models.InputVirtualService) (*models.VirtualService, error) {
	v1VirtualService, err := r.Converter.ConvertInputVirtualService(virtualService)
	if err != nil {
		return nil, err
	}
	out, err := r.VirtualServices.Write(v1VirtualService, clients.WriteOpts{
		Ctx:               ctx,
		OverwriteExisting: overwrite,
	})
	if err != nil {
		return nil, err
	}
	return r.Converter.ConvertOutputVirtualService(out), nil
}

func (r *virtualServiceMutationResolver) Create(ctx context.Context, obj *customtypes.VirtualServiceMutation, virtualService models.InputVirtualService) (*models.VirtualService, error) {
	return r.write(false, ctx, obj, virtualService)
}
func (r *virtualServiceMutationResolver) Update(ctx context.Context, obj *customtypes.VirtualServiceMutation, virtualService models.InputVirtualService) (*models.VirtualService, error) {
	return r.write(true, ctx, obj, virtualService)
}
func (r *virtualServiceMutationResolver) Delete(ctx context.Context, obj *customtypes.VirtualServiceMutation, name string) (*models.VirtualService, error) {
	virtualService, err := r.VirtualServices.Read(obj.Namespace, name, clients.ReadOpts{
		Ctx: ctx,
	})
	if err != nil {
		if errors.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	err = r.VirtualServices.Delete(obj.Namespace, name, clients.DeleteOpts{Ctx: ctx})
	if err != nil {
		return nil, err
	}
	return r.Converter.ConvertOutputVirtualService(virtualService), nil
}
func (r *virtualServiceMutationResolver) AddRoute(ctx context.Context, obj *customtypes.VirtualServiceMutation, virtualServiceName string, resourceVersion string, index int, route models.InputRoute) (*models.VirtualService, error) {
	v1Route, err := r.Converter.ConvertInputRoute(route)
	if err != nil {
		return nil, err
	}

	virtualService, err := r.VirtualServices.Read(obj.Namespace, virtualServiceName, clients.ReadOpts{Ctx: ctx})
	if err != nil {
		return nil, err
	}
	if virtualService.Metadata.ResourceVersion != resourceVersion {
		return nil, errors.Errorf("resource version mismatch. received %v, want %v", resourceVersion, virtualService.Metadata.ResourceVersion)
	}

	if index > len(virtualService.VirtualHost.Routes) {
		index = len(virtualService.VirtualHost.Routes)
	}
	virtualService.VirtualHost.Routes = append(virtualService.VirtualHost.Routes, nil)
	copy(virtualService.VirtualHost.Routes[index+1:], virtualService.VirtualHost.Routes[index:])
	virtualService.VirtualHost.Routes[index] = v1Route

	out, err := r.VirtualServices.Write(virtualService, clients.WriteOpts{
		Ctx:               ctx,
		OverwriteExisting: true,
	})
	if err != nil {
		return nil, err
	}
	return r.Converter.ConvertOutputVirtualService(out), nil
}

func (r *virtualServiceMutationResolver) UpdateRoute(ctx context.Context, obj *customtypes.VirtualServiceMutation, virtualServiceName string, resourceVersion string, index int, route models.InputRoute) (*models.VirtualService, error) {
	v1Route, err := r.Converter.ConvertInputRoute(route)
	if err != nil {
		return nil, err
	}

	virtualService, err := r.VirtualServices.Read(obj.Namespace, virtualServiceName, clients.ReadOpts{Ctx: ctx})
	if err != nil {
		return nil, err
	}
	if virtualService.Metadata.ResourceVersion != resourceVersion {
		return nil, errors.Errorf("resource version mismatch. received %v, want %v", resourceVersion, virtualService.Metadata.ResourceVersion)
	}

	if index > len(virtualService.VirtualHost.Routes) {
		return nil, errors.Errorf("index out of bounds")
	}

	virtualService.VirtualHost.Routes[index] = v1Route

	out, err := r.VirtualServices.Write(virtualService, clients.WriteOpts{
		Ctx:               ctx,
		OverwriteExisting: true,
	})
	if err != nil {
		return nil, err
	}
	return r.Converter.ConvertOutputVirtualService(out), nil
}

func (r *virtualServiceMutationResolver) DeleteRoute(ctx context.Context, obj *customtypes.VirtualServiceMutation, virtualServiceName string, resourceVersion string, index int) (*models.VirtualService, error) {
	virtualService, err := r.VirtualServices.Read(obj.Namespace, virtualServiceName, clients.ReadOpts{Ctx: ctx})
	if err != nil {
		return nil, err
	}
	if virtualService.Metadata.ResourceVersion != resourceVersion {
		return nil, errors.Errorf("resource version mismatch. received %v, want %v", resourceVersion, virtualService.Metadata.ResourceVersion)
	}

	if index > len(virtualService.VirtualHost.Routes) {
		return nil, errors.Errorf("index out of bounds")
	}

	virtualService.VirtualHost.Routes = append(virtualService.VirtualHost.Routes[:index], virtualService.VirtualHost.Routes[index+1:]...)

	out, err := r.VirtualServices.Write(virtualService, clients.WriteOpts{
		Ctx:               ctx,
		OverwriteExisting: true,
	})
	if err != nil {
		return nil, err
	}
	return r.Converter.ConvertOutputVirtualService(out), nil
}

func (r *virtualServiceMutationResolver) SwapRoutes(ctx context.Context, obj *customtypes.VirtualServiceMutation, virtualServiceName string, resourceVersion string, index1 int, index2 int) (*models.VirtualService, error) {
	virtualService, err := r.VirtualServices.Read(obj.Namespace, virtualServiceName, clients.ReadOpts{Ctx: ctx})
	if err != nil {
		return nil, err
	}
	if virtualService.Metadata.ResourceVersion != resourceVersion {
		return nil, errors.Errorf("resource version mismatch. received %v, want %v", resourceVersion, virtualService.Metadata.ResourceVersion)
	}

	if index1 > len(virtualService.VirtualHost.Routes) || index2 > len(virtualService.VirtualHost.Routes) {
		return nil, errors.Errorf("index out of bounds")
	}

	virtualService.VirtualHost.Routes[index1], virtualService.VirtualHost.Routes[index2] = virtualService.VirtualHost.Routes[index2], virtualService.VirtualHost.Routes[index1]

	out, err := r.VirtualServices.Write(virtualService, clients.WriteOpts{
		Ctx:               ctx,
		OverwriteExisting: true,
	})
	if err != nil {
		return nil, err
	}
	return r.Converter.ConvertOutputVirtualService(out), nil

}

type virtualServiceQueryResolver struct{ *ApiResolver }

func (r *virtualServiceQueryResolver) List(ctx context.Context, obj *customtypes.VirtualServiceQuery, selector *models.InputMapStringString) ([]*models.VirtualService, error) {
	var convertedSelector map[string]string
	if selector != nil {
		convertedSelector = selector.GoType()
	}
	list, err := r.VirtualServices.List(obj.Namespace, clients.ListOpts{
		Ctx:      ctx,
		Selector: convertedSelector,
	})
	if err != nil {
		return nil, err
	}
	return r.Converter.ConvertOutputVirtualServices(list), nil
}

func (r *virtualServiceQueryResolver) Get(ctx context.Context, obj *customtypes.VirtualServiceQuery, name string) (*models.VirtualService, error) {
	virtualService, err := r.VirtualServices.Read(obj.Namespace, name, clients.ReadOpts{
		Ctx: ctx,
	})
	if err != nil {
		return nil, err
	}
	return r.Converter.ConvertOutputVirtualService(virtualService), nil
}

type resolverMapMutationResolver struct{ *ApiResolver }

func (r *resolverMapMutationResolver) write(overwrite bool, ctx context.Context, obj *customtypes.ResolverMapMutation, resolverMap models.InputResolverMap) (*models.ResolverMap, error) {
	ups, err := r.Converter.ConvertInputResolverMap(resolverMap)
	if err != nil {
		return nil, err
	}
	out, err := r.ResolverMaps.Write(ups, clients.WriteOpts{
		Ctx:               ctx,
		OverwriteExisting: overwrite,
	})
	if err != nil {
		return nil, err
	}
	return r.Converter.ConvertOutputResolverMap(out), nil
}

func (r *resolverMapMutationResolver) Create(ctx context.Context, obj *customtypes.ResolverMapMutation, resolverMap models.InputResolverMap) (*models.ResolverMap, error) {
	return r.write(false, ctx, obj, resolverMap)
}
func (r *resolverMapMutationResolver) Update(ctx context.Context, obj *customtypes.ResolverMapMutation, resolverMap models.InputResolverMap) (*models.ResolverMap, error) {
	return r.write(true, ctx, obj, resolverMap)
}
func (r *resolverMapMutationResolver) Delete(ctx context.Context, obj *customtypes.ResolverMapMutation, name string) (*models.ResolverMap, error) {
	resolverMap, err := r.ResolverMaps.Read(obj.Namespace, name, clients.ReadOpts{
		Ctx: ctx,
	})
	if err != nil {
		if errors.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	err = r.ResolverMaps.Delete(obj.Namespace, name, clients.DeleteOpts{Ctx: ctx})
	if err != nil {
		return nil, err
	}
	return r.Converter.ConvertOutputResolverMap(resolverMap), nil
}

type resolverMapQueryResolver struct{ *ApiResolver }

func (r *resolverMapQueryResolver) List(ctx context.Context, obj *customtypes.ResolverMapQuery, selector *models.InputMapStringString) ([]*models.ResolverMap, error) {
	var convertedSelector map[string]string
	if selector != nil {
		convertedSelector = selector.GoType()
	}
	list, err := r.ResolverMaps.List(obj.Namespace, clients.ListOpts{
		Ctx:      ctx,
		Selector: convertedSelector,
	})
	if err != nil {
		return nil, err
	}
	return r.Converter.ConvertOutputResolverMaps(list), nil
}

func (r *resolverMapQueryResolver) Get(ctx context.Context, obj *customtypes.ResolverMapQuery, name string) (*models.ResolverMap, error) {
	resolverMap, err := r.ResolverMaps.Read(obj.Namespace, name, clients.ReadOpts{
		Ctx: ctx,
	})
	if err != nil {
		return nil, err
	}
	return r.Converter.ConvertOutputResolverMap(resolverMap), nil
}
