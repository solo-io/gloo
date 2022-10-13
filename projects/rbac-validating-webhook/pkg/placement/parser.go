package placement

import (
	"context"

	fed_enterprise_v1 "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.enterprise.gloo.solo.io/v1"
	fed_gateway_v1 "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.gateway.solo.io/v1"
	fed_gloo_v1 "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.gloo.solo.io/v1"
	fed_ratelimit_v1alpha1 "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.ratelimit.solo.io/v1alpha1"
	fedv1 "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.solo.io/v1"
	multicluster_types "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/multicluster.solo.io/v1alpha1/types"
)

func NewTypedParser() TypedParser {
	return &typedParserImpl{}
}

type typedParserImpl struct {
}

// Failover procesor only actually modifies primary upstream
// If READ permissions are added later, then we can address the targets
func (t *typedParserImpl) ParseFailoverScheme(_ context.Context, obj *fedv1.FailoverScheme) ([]*multicluster_types.Placement, error) {
	return []*multicluster_types.Placement{&multicluster_types.Placement{
		Namespaces: []string{obj.Spec.GetPrimary().GetNamespace()},
		Clusters:   []string{obj.Spec.GetPrimary().GetClusterName()},
	}}, nil
}

func (t *typedParserImpl) ParseFederatedAuthConfig(_ context.Context, obj *fed_enterprise_v1.FederatedAuthConfig) ([]*multicluster_types.Placement, error) {
	return []*multicluster_types.Placement{obj.Spec.GetPlacement()}, nil
}

func (t *typedParserImpl) ParseFederatedRateLimitConfig(_ context.Context, obj *fed_ratelimit_v1alpha1.FederatedRateLimitConfig) ([]*multicluster_types.Placement, error) {
	return []*multicluster_types.Placement{obj.Spec.GetPlacement()}, nil
}

func (t *typedParserImpl) ParseFederatedUpstream(_ context.Context, obj *fed_gloo_v1.FederatedUpstream) ([]*multicluster_types.Placement, error) {
	return []*multicluster_types.Placement{obj.Spec.GetPlacement()}, nil
}

func (t *typedParserImpl) ParseFederatedSettings(_ context.Context, obj *fed_gloo_v1.FederatedSettings) ([]*multicluster_types.Placement, error) {
	return []*multicluster_types.Placement{obj.Spec.GetPlacement()}, nil
}

func (t *typedParserImpl) ParseFederatedUpstreamGroup(_ context.Context, obj *fed_gloo_v1.FederatedUpstreamGroup) ([]*multicluster_types.Placement, error) {
	return []*multicluster_types.Placement{obj.Spec.GetPlacement()}, nil
}

func (t *typedParserImpl) ParseFederatedGateway(_ context.Context, obj *fed_gateway_v1.FederatedGateway) ([]*multicluster_types.Placement, error) {
	return []*multicluster_types.Placement{obj.Spec.GetPlacement()}, nil
}

func (t *typedParserImpl) ParseFederatedMatchableHttpGateway(ctx context.Context, obj *fed_gateway_v1.FederatedMatchableHttpGateway) ([]*multicluster_types.Placement, error) {
	return []*multicluster_types.Placement{obj.Spec.GetPlacement()}, nil
}

func (t *typedParserImpl) ParseFederatedVirtualService(_ context.Context, obj *fed_gateway_v1.FederatedVirtualService) ([]*multicluster_types.Placement, error) {
	return []*multicluster_types.Placement{obj.Spec.GetPlacement()}, nil
}

func (t *typedParserImpl) ParseFederatedRouteTable(_ context.Context, obj *fed_gateway_v1.FederatedRouteTable) ([]*multicluster_types.Placement, error) {
	return []*multicluster_types.Placement{obj.Spec.GetPlacement()}, nil
}
