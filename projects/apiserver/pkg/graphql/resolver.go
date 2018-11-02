package graphql

import (
	"context"

	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
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
	Secrets         v1.SecretClient
	Artifacts       v1.ArtifactClient
	Settings        v1.SettingsClient
	VirtualServices gatewayv1.VirtualServiceClient
	ResolverMaps    sqoopv1.ResolverMapClient
	Schemas         sqoopv1.SchemaClient
}

func NewResolvers(upstreams v1.UpstreamClient,
	schemas sqoopv1.SchemaClient,
	artifacts v1.ArtifactClient,
	settings v1.SettingsClient,
	secrets v1.SecretClient,
	virtualServices gatewayv1.VirtualServiceClient,
	resolverMaps sqoopv1.ResolverMapClient) *ApiResolver {
	return &ApiResolver{
		Upstreams:       upstreams,
		VirtualServices: virtualServices,
		ResolverMaps:    resolverMaps,
		Schemas:         schemas,
		Artifacts:       artifacts,
		Settings:        settings,
		Secrets:         secrets,
		// TODO(ilackarms): just make these private functions, remove converter
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

func (r *ApiResolver) SchemaMutation() graph.SchemaMutationResolver {
	return &schemaMutationResolver{r}
}
func (r *ApiResolver) SchemaQuery() graph.SchemaQueryResolver {
	return &schemaQueryResolver{r}
}

func (r *ApiResolver) ArtifactMutation() graph.ArtifactMutationResolver {
	return &artifactMutationResolver{r}
}

func (r *ApiResolver) ArtifactQuery() graph.ArtifactQueryResolver {
	return &artifactQueryResolver{r}
}

func (r *ApiResolver) SettingsMutation() graph.SettingsMutationResolver {
	return &settingsMutationResolver{r}
}

func (r *ApiResolver) SettingsQuery() graph.SettingsQueryResolver {
	return &settingsQueryResolver{r}
}

func (r *ApiResolver) SecretMutation() graph.SecretMutationResolver {
	return &secretMutationResolver{r}
}

func (r *ApiResolver) SecretQuery() graph.SecretQueryResolver {
	return &secretQueryResolver{r}
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
func (r *mutationResolver) Schemas(ctx context.Context, namespace string) (customtypes.SchemaMutation, error) {
	return customtypes.SchemaMutation{Namespace: namespace}, nil
}
func (r *mutationResolver) Secrets(ctx context.Context, namespace string) (customtypes.SecretMutation, error) {
	return customtypes.SecretMutation{Namespace: namespace}, nil
}
func (r *mutationResolver) Artifacts(ctx context.Context, namespace string) (customtypes.ArtifactMutation, error) {
	return customtypes.ArtifactMutation{Namespace: namespace}, nil
}
func (r *mutationResolver) Settings(ctx context.Context) (customtypes.SettingsMutation, error) {
	return customtypes.SettingsMutation{}, nil
}

type queryResolver struct{ *ApiResolver }

func (r *queryResolver) Resource(ctx context.Context, guid string) (models.Resource, error) {
	kind, namespace, name, err := resources.SplitKey(guid)
	if err != nil {
		return nil, err
	}
	switch kind {
	case resources.Kind(&v1.Upstream{}):
		return r.UpstreamQuery().Get(ctx, &customtypes.UpstreamQuery{Namespace: namespace}, name)
	case resources.Kind(&gatewayv1.VirtualService{}):
		return r.VirtualServiceQuery().Get(ctx, &customtypes.VirtualServiceQuery{Namespace: namespace}, name)
	case resources.Kind(&sqoopv1.ResolverMap{}):
		return r.ResolverMapQuery().Get(ctx, &customtypes.ResolverMapQuery{Namespace: namespace}, name)
	case resources.Kind(&sqoopv1.Schema{}):
		return r.SchemaQuery().Get(ctx, &customtypes.SchemaQuery{Namespace: namespace}, name)
	case resources.Kind(&v1.Secret{}):
		return r.SecretQuery().Get(ctx, &customtypes.SecretQuery{Namespace: namespace}, name)
	case resources.Kind(&v1.Artifact{}):
		return r.ArtifactQuery().Get(ctx, &customtypes.ArtifactQuery{Namespace: namespace}, name)
	case resources.Kind(&v1.Settings{}):
		return r.SettingsQuery().Get(ctx, &customtypes.SettingsQuery{})
	}
	return nil, errors.Errorf("unknown kind %v", kind)
}

func (r *queryResolver) GetOAuthEndpoint(ctx context.Context) (models.OAuthEndpoint, error) {
	return getOAuthEndpoint()
}

func (r *queryResolver) Version(ctx context.Context) (string, error) {
	return getAPIVersion(), nil
}

func (r *queryResolver) Upstreams(ctx context.Context, namespace string) (customtypes.UpstreamQuery, error) {
	return customtypes.UpstreamQuery{Namespace: namespace}, nil
}
func (r *queryResolver) VirtualServices(ctx context.Context, namespace string) (customtypes.VirtualServiceQuery, error) {
	return customtypes.VirtualServiceQuery{Namespace: namespace}, nil
}
func (r *queryResolver) ResolverMaps(ctx context.Context, namespace string) (customtypes.ResolverMapQuery, error) {
	return customtypes.ResolverMapQuery{Namespace: namespace}, nil
}
func (r *queryResolver) Schemas(ctx context.Context, namespace string) (customtypes.SchemaQuery, error) {
	return customtypes.SchemaQuery{Namespace: namespace}, nil
}
func (r *queryResolver) Secrets(ctx context.Context, namespace string) (customtypes.SecretQuery, error) {
	return customtypes.SecretQuery{Namespace: namespace}, nil
}
func (r *queryResolver) Artifacts(ctx context.Context, namespace string) (customtypes.ArtifactQuery, error) {
	return customtypes.ArtifactQuery{Namespace: namespace}, nil
}
func (r *queryResolver) Settings(ctx context.Context) (customtypes.SettingsQuery, error) {
	return customtypes.SettingsQuery{}, nil
}
func (r *ApiResolver) Subscription() graph.SubscriptionResolver {
	return &subscriptionResolver{r}
}
