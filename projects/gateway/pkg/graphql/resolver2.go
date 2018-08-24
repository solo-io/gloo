//go:generate gorunpkg github.com/99designs/gqlgen

package graphql

import (
	context "context"

	customtypes "github.com/solo-io/solo-kit/projects/gateway/pkg/graphql/customtypes"
	graph "github.com/solo-io/solo-kit/projects/gateway/pkg/graphql/graph"
	models "github.com/solo-io/solo-kit/projects/gateway/pkg/graphql/models"
)

type Resolver2 struct{}

func (r *Resolver2) Mutation() graph.MutationResolver {
	return &mutationResolver{r}
}
func (r *Resolver2) Query() graph.QueryResolver {
	return &queryResolver{r}
}
func (r *Resolver2) ResolverMapMutation() graph.ResolverMapMutationResolver {
	return &resolverMapMutationResolver{r}
}
func (r *Resolver2) ResolverMapQuery() graph.ResolverMapQueryResolver {
	return &resolverMapQueryResolver{r}
}
func (r *Resolver2) UpstreamMutation() graph.UpstreamMutationResolver {
	return &upstreamMutationResolver{r}
}
func (r *Resolver2) UpstreamQuery() graph.UpstreamQueryResolver {
	return &upstreamQueryResolver{r}
}
func (r *Resolver2) VirtualServiceMutation() graph.VirtualServiceMutationResolver {
	return &virtualServiceMutationResolver{r}
}
func (r *Resolver2) VirtualServiceQuery() graph.VirtualServiceQueryResolver {
	return &virtualServiceQueryResolver{r}
}

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) Upstreams(ctx context.Context, namespace string) (customtypes.UpstreamMutation, error) {
	panic("not implemented")
}
func (r *mutationResolver) VirtualServices(ctx context.Context, namespace string) (customtypes.VirtualServiceMutation, error) {
	panic("not implemented")
}
func (r *mutationResolver) ResolverMaps(ctx context.Context, namespace string) (customtypes.ResolverMapMutation, error) {
	panic("not implemented")
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) Upstreams(ctx context.Context, namespace string) (customtypes.UpstreamQuery, error) {
	panic("not implemented")
}
func (r *queryResolver) VirtualServices(ctx context.Context, namespace string) (customtypes.VirtualServiceQuery, error) {
	panic("not implemented")
}
func (r *queryResolver) ResolverMaps(ctx context.Context, namespace string) (customtypes.ResolverMapQuery, error) {
	panic("not implemented")
}

type resolverMapMutationResolver struct{ *Resolver }

func (r *resolverMapMutationResolver) Create(ctx context.Context, obj *customtypes.ResolverMapMutation, resolverMap models.InputResolverMap) (*models.ResolverMap, error) {
	panic("not implemented")
}
func (r *resolverMapMutationResolver) Update(ctx context.Context, obj *customtypes.ResolverMapMutation, resolverMap models.InputResolverMap) (*models.ResolverMap, error) {
	panic("not implemented")
}
func (r *resolverMapMutationResolver) Delete(ctx context.Context, obj *customtypes.ResolverMapMutation, name string) (*models.ResolverMap, error) {
	panic("not implemented")
}

type resolverMapQueryResolver struct{ *Resolver }

func (r *resolverMapQueryResolver) List(ctx context.Context, obj *customtypes.ResolverMapQuery, selector *customtypes.MapStringString) ([]*models.ResolverMap, error) {
	panic("not implemented")
}
func (r *resolverMapQueryResolver) Get(ctx context.Context, obj *customtypes.ResolverMapQuery, name string) (*models.ResolverMap, error) {
	panic("not implemented")
}

type upstreamMutationResolver struct{ *Resolver }

func (r *upstreamMutationResolver) Create(ctx context.Context, obj *customtypes.UpstreamMutation, upstream models.InputUpstream) (*models.Upstream, error) {
	panic("not implemented")
}
func (r *upstreamMutationResolver) Update(ctx context.Context, obj *customtypes.UpstreamMutation, upstream models.InputUpstream) (*models.Upstream, error) {
	panic("not implemented")
}
func (r *upstreamMutationResolver) Delete(ctx context.Context, obj *customtypes.UpstreamMutation, name string) (*models.Upstream, error) {
	panic("not implemented")
}

type upstreamQueryResolver struct{ *Resolver }

func (r *upstreamQueryResolver) List(ctx context.Context, obj *customtypes.UpstreamQuery, selector *customtypes.MapStringString) ([]*models.Upstream, error) {
	panic("not implemented")
}
func (r *upstreamQueryResolver) Get(ctx context.Context, obj *customtypes.UpstreamQuery, name string) (*models.Upstream, error) {
	panic("not implemented")
}

type virtualServiceMutationResolver struct{ *Resolver }

func (r *virtualServiceMutationResolver) Create(ctx context.Context, obj *customtypes.VirtualServiceMutation, virtualService models.InputVirtualService) (*models.VirtualService, error) {
	panic("not implemented")
}
func (r *virtualServiceMutationResolver) Update(ctx context.Context, obj *customtypes.VirtualServiceMutation, virtualService models.InputVirtualService) (*models.VirtualService, error) {
	panic("not implemented")
}
func (r *virtualServiceMutationResolver) Delete(ctx context.Context, obj *customtypes.VirtualServiceMutation, name string) (*models.VirtualService, error) {
	panic("not implemented")
}
func (r *virtualServiceMutationResolver) AddRoute(ctx context.Context, obj *customtypes.VirtualServiceMutation, virtualServiceName string, index int, route models.InputRoute) (*models.VirtualService, error) {
	panic("not implemented")
}
func (r *virtualServiceMutationResolver) UpdateRoute(ctx context.Context, obj *customtypes.VirtualServiceMutation, virtualServiceName string, index int, route models.InputRoute) (*models.VirtualService, error) {
	panic("not implemented")
}
func (r *virtualServiceMutationResolver) DeleteRoute(ctx context.Context, obj *customtypes.VirtualServiceMutation, virtualServiceName string, index int) (*models.VirtualService, error) {
	panic("not implemented")
}
func (r *virtualServiceMutationResolver) SwapRoutes(ctx context.Context, obj *customtypes.VirtualServiceMutation, virtualServiceName string, index1 int, index2 int) (*models.VirtualService, error) {
	panic("not implemented")
}

type virtualServiceQueryResolver struct{ *Resolver }

func (r *virtualServiceQueryResolver) List(ctx context.Context, obj *customtypes.VirtualServiceQuery, selector *customtypes.MapStringString) ([]*models.VirtualService, error) {
	panic("not implemented")
}
func (r *virtualServiceQueryResolver) Get(ctx context.Context, obj *customtypes.VirtualServiceQuery, name string) (*models.VirtualService, error) {
	panic("not implemented")
}
