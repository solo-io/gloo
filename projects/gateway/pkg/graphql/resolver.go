//go:generate gorunpkg github.com/99designs/gqlgen

package graphql

import (
	"context"

	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/errors"
	gatewayv1 "github.com/solo-io/solo-kit/projects/gateway/pkg/api/v1"
	"github.com/solo-io/solo-kit/projects/gateway/pkg/graphql/customtypes"
	"github.com/solo-io/solo-kit/projects/gateway/pkg/graphql/graph"
	"github.com/solo-io/solo-kit/projects/gateway/pkg/graphql/models"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1"
)

type ApiResolver struct {
	Upstreams       v1.UpstreamClient
	VirtualServices gatewayv1.VirtualServiceClient
	Converter       *Converter
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
	panic("not implemented")
}

type queryResolver struct{ *ApiResolver }

func (r *queryResolver) Upstreams(ctx context.Context, namespace string) (customtypes.UpstreamQuery, error) {
	return customtypes.UpstreamQuery{Namespace: namespace}, nil
}
func (r *queryResolver) VirtualServices(ctx context.Context, namespace string) (customtypes.VirtualServiceQuery, error) {
	return customtypes.VirtualServiceQuery{Namespace: namespace}, nil
}
func (r *queryResolver) ResolverMaps(ctx context.Context, namespace string) (customtypes.ResolverMapQuery, error) {
	panic("not implemented")
}

type upstreamMutationResolver struct{ *ApiResolver }

func (r *upstreamMutationResolver) write(overwrite bool, ctx context.Context, obj *customtypes.UpstreamMutation, upstream models.InputUpstream) (*models.Upstream, error) {
	ups := r.Converter.ConvertInputUpstream(upstream)
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

func (r *upstreamQueryResolver) List(ctx context.Context, obj *customtypes.UpstreamQuery, selector *customtypes.MapStringString) ([]*models.Upstream, error) {
	var convertedSelector map[string]string
	if selector != nil {
		convertedSelector = selector.GetMap()
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
	ups, err := r.Converter.ConvertInputVirtualService(virtualService)
	if err != nil {
		return nil, err
	}
	out, err := r.VirtualServices.Write(ups, clients.WriteOpts{
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

type virtualServiceQueryResolver struct{ *ApiResolver }

func (r *virtualServiceQueryResolver) List(ctx context.Context, obj *customtypes.VirtualServiceQuery, selector *customtypes.MapStringString) ([]*models.VirtualService, error) {
	var convertedSelector map[string]string
	if selector != nil {
		convertedSelector = selector.GetMap()
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
	ups := r.Converter.ConvertInputResolverMap(resolverMap)
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

func (r *resolverMapQueryResolver) List(ctx context.Context, obj *customtypes.ResolverMapQuery, selector *customtypes.MapStringString) ([]*models.ResolverMap, error) {
	var convertedSelector map[string]string
	if selector != nil {
		convertedSelector = selector.GetMap()
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
