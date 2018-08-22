//go:generate gorunpkg github.com/99designs/gqlgen

package graphql

import (
	"context"

	"github.com/solo-io/solo-kit/projects/gateway/pkg/graphql/customtypes"
	"github.com/solo-io/solo-kit/projects/gateway/pkg/graphql/graph"
	"github.com/solo-io/solo-kit/projects/gateway/pkg/graphql/models"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1"
	gatewayv1 "github.com/solo-io/solo-kit/projects/gateway/pkg/api/v1"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
)

type Resolver struct {
	Upstreams       v1.UpstreamClient
	VirtualServices gatewayv1.VirtualServiceClient
	Converter       *models.Converter
}

func (r *Resolver) Mutation() graph.MutationResolver {
	return &mutationResolver{r}
}
func (r *Resolver) Query() graph.QueryResolver {
	return &queryResolver{r}
}
func (r *Resolver) UpstreamMutation() graph.UpstreamMutationResolver {
	return &upstreamMutationResolver{Resolver: r}
}
func (r *Resolver) UpstreamQuery() graph.UpstreamQueryResolver {
	return &upstreamQueryResolver{Resolver: r}
}

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) Upstreams(ctx context.Context, namespace string) (*customtypes.UpstreamMutation, error) {
	return &customtypes.UpstreamMutation{Namespace: namespace}, nil
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) Upstreams(ctx context.Context, namespace string) (*customtypes.UpstreamQuery, error) {
	return &customtypes.UpstreamQuery{Namespace: namespace}, nil
}

type upstreamMutationResolver struct {
	*Resolver
	namespace string
}

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
	panic("not implemented")
}

type upstreamQueryResolver struct {
	*Resolver
	namespace string
}

func (r *upstreamQueryResolver) List(ctx context.Context, obj *customtypes.UpstreamQuery, selector *customtypes.MapStringString) ([]*models.Upstream, error) {
	panic("not implemented")
}
func (r *upstreamQueryResolver) Get(ctx context.Context, obj *customtypes.UpstreamQuery, name string) (*models.Upstream, error) {
	panic("not implemented")
}
