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

type Resolver struct {
	Upstreams       v1.UpstreamClient
	VirtualServices gatewayv1.VirtualServiceClient
	Converter       *Converter
}

func (r *Resolver) Mutation() graph.MutationResolver {
	return &mutationResolver{r}
}
func (r *Resolver) Query() graph.QueryResolver {
	return &queryResolver{r}
}
func (r *Resolver) UpstreamMutation() graph.UpstreamMutationResolver {
	return &upstreamMutationResolver{r}
}
func (r *Resolver) UpstreamQuery() graph.UpstreamQueryResolver {
	return &upstreamQueryResolver{r}
}
func (r *Resolver) VirtualServiceMutation() graph.VirtualServiceMutationResolver {
	return &virtualServiceMutationResolver{r}
}
func (r *Resolver) VirtualServiceQuery() graph.VirtualServiceQueryResolver {
	return &virtualServiceQueryResolver{r}
}

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) Upstreams(ctx context.Context, namespace string) (customtypes.UpstreamMutation, error) {
	return customtypes.UpstreamMutation{Namespace: namespace}, nil
}
func (r *mutationResolver) VirtualServices(ctx context.Context, namespace string) (customtypes.VirtualServiceMutation, error) {
	return customtypes.VirtualServiceMutation{Namespace: namespace}, nil
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) Upstreams(ctx context.Context, namespace string) (customtypes.UpstreamQuery, error) {
	return customtypes.UpstreamQuery{Namespace: namespace}, nil
}
func (r *queryResolver) VirtualServices(ctx context.Context, namespace string) (customtypes.VirtualServiceQuery, error) {
	return customtypes.VirtualServiceQuery{Namespace: namespace}, nil
}

type upstreamMutationResolver struct{ *Resolver }

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

type upstreamQueryResolver struct{ *Resolver }

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

type virtualServiceMutationResolver struct{ *Resolver }

func (r *virtualServiceMutationResolver) Create(ctx context.Context, obj *customtypes.VirtualServiceMutation, upstream models.InputVirtualService) (*models.VirtualService, error) {
	panic("not implemented")
}
func (r *virtualServiceMutationResolver) Update(ctx context.Context, obj *customtypes.VirtualServiceMutation, upstream models.InputVirtualService) (*models.VirtualService, error) {
	panic("not implemented")
}
func (r *virtualServiceMutationResolver) Delete(ctx context.Context, obj *customtypes.VirtualServiceMutation, name string) (*models.VirtualService, error) {
	panic("not implemented")
}

type virtualServiceQueryResolver struct{ *Resolver }

func (r *virtualServiceQueryResolver) List(ctx context.Context, obj *customtypes.VirtualServiceQuery, selector *customtypes.MapStringString) ([]*models.VirtualService, error) {
	panic("not implemented")
}
func (r *virtualServiceQueryResolver) Get(ctx context.Context, obj *customtypes.VirtualServiceQuery, name string) (*models.VirtualService, error) {
	panic("not implemented")
}
