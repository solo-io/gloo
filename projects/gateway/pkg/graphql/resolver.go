//go:generate gorunpkg github.com/99designs/gqlgen

package graphql

import (
	context "context"

	customtypes "github.com/solo-io/solo-kit/projects/gateway/pkg/graphql/customtypes"
	graph "github.com/solo-io/solo-kit/projects/gateway/pkg/graphql/graph"
	models "github.com/solo-io/solo-kit/projects/gateway/pkg/graphql/models"
)

type Resolver struct{}

func (r *Resolver) Mutation() graph.MutationResolver {
	return &mutationResolver{r}
}
func (r *Resolver) Query() graph.QueryResolver {
	return &queryResolver{r}
}
func (r *Resolver) UpstreamQuery() graph.UpstreamQueryResolver {
	return &upstreamQueryResolver{r}
}

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) Upstreams(ctx context.Context, namespace string) (*models.UpstreamMutation, error) {
	panic("not implemented")
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) Upstreams(ctx context.Context, namespace string) (*customtypes.UpstreamQuery, error) {
	panic("not implemented")
}

type upstreamQueryResolver struct{ *Resolver }

func (r *upstreamQueryResolver) List(ctx context.Context, obj *customtypes.UpstreamQuery, selector *customtypes.MapStringString) ([]*models.Upstream, error) {
	panic("not implemented")
}
func (r *upstreamQueryResolver) Get(ctx context.Context, obj *customtypes.UpstreamQuery, name string) (*models.Upstream, error) {
	panic("not implemented")
}
