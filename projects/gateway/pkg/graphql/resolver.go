//go:generate gorunpkg github.com/99designs/gqlgen

package graphql

import (
	"context"

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

func (r *mutationResolver) Upstreams(ctx context.Context, namespace string) (*models.UpstreamMutation, error) {
	panic("not implemented")
}
func (r *mutationResolver) VirtualServices(ctx context.Context, namespace string) (*models.VirtualServiceMutation, error) {
	panic("not implemented")
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) Upstreams(ctx context.Context, namespace string) (*models.UpstreamQuery, error) {
	panic("not implemented")
}
func (r *queryResolver) VirtualServices(ctx context.Context, namespace string) (*models.VirtualServiceQuery, error) {
	panic("not implemented")
}

type upstreamMutationResolver struct {
	*Resolver
	namespace string
}

func (r *upstreamMutationResolver) Create(ctx context.Context, obj *models.UpstreamMutation, upstream models.InputUpstream) (*models.Upstream, error) {
	panic("not implemented")
}
func (r *upstreamMutationResolver) Update(ctx context.Context, obj *models.UpstreamMutation, upstream models.InputUpstream) (*models.Upstream, error) {
	panic("not implemented")
}
func (r *upstreamMutationResolver) Delete(ctx context.Context, obj *models.UpstreamMutation, name string) (*models.Upstream, error) {
	panic("not implemented")
}

type upstreamQueryResolver struct {
	*Resolver
	namespace string
}

func (r *upstreamQueryResolver) List(ctx context.Context, obj *models.UpstreamQuery, selector *customtypes.MapStringString) ([]*models.Upstream, error) {
	panic("not implemented")
}
func (r *upstreamQueryResolver) Get(ctx context.Context, obj *models.UpstreamQuery, name string) (*models.Upstream, error) {
	panic("not implemented")
}

type virtualServiceMutationResolver struct {
	*Resolver
	namespace string
}

func (r *virtualServiceMutationResolver) Create(ctx context.Context, obj *models.VirtualServiceMutation, upstream models.InputVirtualService) (*models.VirtualService, error) {
	panic("not implemented")
}
func (r *virtualServiceMutationResolver) Update(ctx context.Context, obj *models.VirtualServiceMutation, upstream models.InputVirtualService) (*models.VirtualService, error) {
	panic("not implemented")
}
func (r *virtualServiceMutationResolver) Delete(ctx context.Context, obj *models.VirtualServiceMutation, name string) (*models.VirtualService, error) {
	panic("not implemented")
}

type virtualServiceQueryResolver struct {
	*Resolver
	namespace string
}

func (r *virtualServiceQueryResolver) List(ctx context.Context, obj *models.VirtualServiceQuery, selector *customtypes.MapStringString) ([]*models.VirtualService, error) {
	panic("not implemented")
}
func (r *virtualServiceQueryResolver) Get(ctx context.Context, obj *models.VirtualServiceQuery, name string) (*models.VirtualService, error) {
	panic("not implemented")
}
