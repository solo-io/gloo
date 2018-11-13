package graphql

import (
	"context"

	"github.com/solo-io/solo-projects/projects/apiserver/pkg/graphql/customtypes"
	"github.com/solo-io/solo-projects/projects/apiserver/pkg/graphql/graph"
	"github.com/solo-io/solo-projects/projects/apiserver/pkg/graphql/models"
)

// If client does not present a token, we allow them to still query certain functions
type UnregisteredResolver struct{}

func NewUnregisteredResolver() *UnregisteredResolver {
	return &UnregisteredResolver{}
}

func (r *UnregisteredResolver) GetOAuthEndpoint(ctx context.Context) (models.OAuthEndpoint, error) {
	return getOAuthEndpoint()
}

func (r *UnregisteredResolver) Version(ctx context.Context) (string, error) {
	return getAPIVersion(), nil
}

func (r *UnregisteredResolver) ArtifactMutation() graph.ArtifactMutationResolver {
	panic("client must present a token to access this feature")
}

func (r *UnregisteredResolver) ArtifactQuery() graph.ArtifactQueryResolver {
	panic("client must present a token to access this feature")
}

func (r *UnregisteredResolver) Mutation() graph.MutationResolver {
	panic("client must present a token to access this feature")
}

func (r *UnregisteredResolver) Query() graph.QueryResolver {
	return r
}

func (r *UnregisteredResolver) ResolverMapMutation() graph.ResolverMapMutationResolver {
	panic("client must present a token to access this feature")
}

func (r *UnregisteredResolver) ResolverMapQuery() graph.ResolverMapQueryResolver {
	panic("client must present a token to access this feature")
}

func (r *UnregisteredResolver) SchemaMutation() graph.SchemaMutationResolver {
	panic("client must present a token to access this feature")
}

func (r *UnregisteredResolver) SchemaQuery() graph.SchemaQueryResolver {
	panic("client must present a token to access this feature")
}

func (r *UnregisteredResolver) SecretMutation() graph.SecretMutationResolver {
	panic("client must present a token to access this feature")
}

func (r *UnregisteredResolver) SecretQuery() graph.SecretQueryResolver {
	panic("client must present a token to access this feature")
}

func (r *UnregisteredResolver) SettingsMutation() graph.SettingsMutationResolver {
	panic("client must present a token to access this feature")
}

func (r *UnregisteredResolver) SettingsQuery() graph.SettingsQueryResolver {
	panic("client must present a token to access this feature")
}

func (r *UnregisteredResolver) Subscription() graph.SubscriptionResolver {
	panic("client must present a token to access this feature")
}

func (r *UnregisteredResolver) UpstreamMutation() graph.UpstreamMutationResolver {
	panic("client must present a token to access this feature")
}

func (r *UnregisteredResolver) UpstreamQuery() graph.UpstreamQueryResolver {
	panic("client must present a token to access this feature")
}

func (r *UnregisteredResolver) VirtualServiceMutation() graph.VirtualServiceMutationResolver {
	panic("client must present a token to access this feature")
}

func (r *UnregisteredResolver) VirtualServiceQuery() graph.VirtualServiceQueryResolver {
	panic("client must present a token to access this feature")
}

func (r *UnregisteredResolver) Resource(ctx context.Context, guid string) (models.Resource, error) {
	panic("client must present a token to access this feature")
}

func (r *UnregisteredResolver) Upstreams(ctx context.Context, namespace string) (customtypes.UpstreamQuery, error) {
	panic("client must present a token to access this feature")
}

func (r *UnregisteredResolver) VirtualServices(ctx context.Context, namespace string) (customtypes.VirtualServiceQuery, error) {
	panic("client must present a token to access this feature")
}

func (r *UnregisteredResolver) ResolverMaps(ctx context.Context, namespace string) (customtypes.ResolverMapQuery, error) {
	panic("client must present a token to access this feature")
}

func (r *UnregisteredResolver) Schemas(ctx context.Context, namespace string) (customtypes.SchemaQuery, error) {
	panic("client must present a token to access this feature")
}

func (r *UnregisteredResolver) Secrets(ctx context.Context, namespace string) (customtypes.SecretQuery, error) {
	panic("client must present a token to access this feature")
}

func (r *UnregisteredResolver) Artifacts(ctx context.Context, namespace string) (customtypes.ArtifactQuery, error) {
	panic("client must present a token to access this feature")
}

func (r *UnregisteredResolver) Settings(ctx context.Context) (customtypes.SettingsQuery, error) {
	panic("client must present a token to access this feature")
}

func (r *UnregisteredResolver) VcsMutation() graph.VcsMutationResolver {
	panic("client must present a token to access this feature")
}
