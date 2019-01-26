package graphql

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	gatewayv1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/defaults"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-projects/projects/apiserver/pkg/graphql/customtypes"
	"github.com/solo-io/solo-projects/projects/apiserver/pkg/graphql/graph"
	"github.com/solo-io/solo-projects/projects/apiserver/pkg/graphql/models"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

type ApiResolver struct {
	UpstreamClient       v1.UpstreamClient
	SecretClient         v1.SecretClient
	ArtifactClient       v1.ArtifactClient
	SettingsClient       v1.SettingsClient
	VirtualServiceClient gatewayv1.VirtualServiceClient
	KubeClient           corev1.CoreV1Interface
}

func NewResolvers(upstreams v1.UpstreamClient,
	artifacts v1.ArtifactClient,
	settings v1.SettingsClient,
	secrets v1.SecretClient,
	virtualServices gatewayv1.VirtualServiceClient,
	kubeClient corev1.CoreV1Interface,
) *ApiResolver {
	return &ApiResolver{
		UpstreamClient:       upstreams,
		VirtualServiceClient: virtualServices,
		ArtifactClient:       artifacts,
		SettingsClient:       settings,
		SecretClient:         secrets,
		KubeClient:           kubeClient,
		// TODO(ilackarms): just make these private functions, remove converter
	}
}

func (r *ApiResolver) Query() graph.QueryResolver {
	return &queryResolver{r}
}
func (r *ApiResolver) Mutation() graph.MutationResolver {
	return &mutationResolver{r}
}
func (r *ApiResolver) Subscription() graph.SubscriptionResolver {
	return &subscriptionResolver{r}
}
func (r *ApiResolver) Namespace() graph.NamespaceResolver {
	return &namespaceResolver{r}
}
func (r *ApiResolver) UpstreamMutation() graph.UpstreamMutationResolver {
	return &upstreamMutationResolver{r}
}
func (r *ApiResolver) VirtualServiceMutation() graph.VirtualServiceMutationResolver {
	return &virtualServiceMutationResolver{r}
}
func (r *ApiResolver) ArtifactMutation() graph.ArtifactMutationResolver {
	return &artifactMutationResolver{r}
}
func (r *ApiResolver) SettingsMutation() graph.SettingsMutationResolver {
	return &settingsMutationResolver{r}
}
func (r *ApiResolver) SecretMutation() graph.SecretMutationResolver {
	return &secretMutationResolver{r}
}

type queryResolver struct{ *ApiResolver }

func (r *queryResolver) GetOAuthEndpoint(ctx context.Context) (models.OAuthEndpoint, error) {
	return getOAuthEndpoint()
}

func (r *queryResolver) Version(ctx context.Context) (string, error) {
	return getAPIVersion(), nil
}

func (r *queryResolver) Namespace(ctx context.Context, name string) (customtypes.Namespace, error) {
	return customtypes.Namespace{Name: name}, nil
}

// This method causes the namespaceResolver (and the Kubernetes clients it uses) to be invoked once for every namespace.
// The Kubernetes go-client could issue a single API call for all namespaces (by passing an empty namespace), but the
// way that gqlgen generates resolvers prevents us from taking advantage of this without creating additional custom
// types that would pollute the GraphQL schema.
// Fortunately, we can live with this, since the client caching mechanism we have in place mitigates the problem.
func (r *queryResolver) AllNamespaces(ctx context.Context) ([]customtypes.Namespace, error) {
	nsList, err := r.KubeClient.Namespaces().List(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	var namespaces []customtypes.Namespace
	for _, ns := range nsList.Items {
		namespaces = append(namespaces, customtypes.Namespace{Name: ns.Name})
	}
	return namespaces, nil
}

func (r *queryResolver) Settings(ctx context.Context) (*models.Settings, error) {
	namespace := defaults.GlooSystem
	name := defaults.SettingsName
	settings, err := r.SettingsClient.Read(namespace, name, clients.ReadOpts{
		Ctx: ctx,
	})
	if err != nil {
		return nil, err
	}
	return NewConverter(r.ApiResolver, ctx).ConvertOutputSettings(settings), nil
}

type mutationResolver struct{ *ApiResolver }

func (r *mutationResolver) Upstreams(ctx context.Context) (customtypes.UpstreamMutation, error) {
	return customtypes.UpstreamMutation{}, nil
}
func (r *mutationResolver) VirtualServices(ctx context.Context) (customtypes.VirtualServiceMutation, error) {
	return customtypes.VirtualServiceMutation{}, nil
}
func (r *mutationResolver) Secrets(ctx context.Context) (customtypes.SecretMutation, error) {
	return customtypes.SecretMutation{}, nil
}
func (r *mutationResolver) Artifacts(ctx context.Context) (customtypes.ArtifactMutation, error) {
	return customtypes.ArtifactMutation{}, nil
}
func (r *mutationResolver) Settings(ctx context.Context) (customtypes.SettingsMutation, error) {
	return customtypes.SettingsMutation{}, nil
}

type namespaceResolver struct{ *ApiResolver }
