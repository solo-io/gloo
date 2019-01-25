package graphql

import (
	"context"

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

func (r *UnregisteredResolver) Mutation() graph.MutationResolver {
	panic("client must present a token to access this feature")
}

func (r *UnregisteredResolver) Query() graph.QueryResolver {
	panic("client must present a token to access this feature")
}

func (r *UnregisteredResolver) Namespace() graph.NamespaceResolver {
	panic("client must present a token to access this feature")
}

func (r *UnregisteredResolver) SecretMutation() graph.SecretMutationResolver {
	panic("client must present a token to access this feature")
}

func (r *UnregisteredResolver) SettingsMutation() graph.SettingsMutationResolver {
	panic("client must present a token to access this feature")
}

func (r *UnregisteredResolver) Subscription() graph.SubscriptionResolver {
	panic("client must present a token to access this feature")
}

func (r *UnregisteredResolver) UpstreamMutation() graph.UpstreamMutationResolver {
	panic("client must present a token to access this feature")
}

func (r *UnregisteredResolver) VirtualServiceMutation() graph.VirtualServiceMutationResolver {
	panic("client must present a token to access this feature")
}
