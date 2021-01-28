package placement

import (
	"context"

	"github.com/solo-io/skv2-enterprise/multicluster-admission-webhook/pkg/api/multicluster.solo.io/v1alpha1"
	fed_enterprise_v1 "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.enterprise.gloo.solo.io/v1"
	fed_gateway_v1 "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.gateway.solo.io/v1"
	fed_gloo_v1 "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.gloo.solo.io/v1"
	fed_ratelimit_v1alpha1 "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.ratelimit.solo.io/v1alpha1"
	fedv1 "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.solo.io/v1"
)

func NewTypedParser() TypedParser {
	return &typedParserImpl{}
}

type typedParserImpl struct {
}

// Failover procesor only actually modifies primary upstream
// If READ permissions are added later, then we can address the targets
func (t *typedParserImpl) ParseFailoverScheme(_ context.Context, obj *fedv1.FailoverScheme) ([]*v1alpha1.Placement, error) {
	return []*v1alpha1.Placement{&v1alpha1.Placement{
		Namespaces: []string{obj.Spec.GetPrimary().GetNamespace()},
		Clusters:   []string{obj.Spec.GetPrimary().GetClusterName()},
	}}, nil
}

func (t *typedParserImpl) ParseFederatedAuthConfig(_ context.Context, obj *fed_enterprise_v1.FederatedAuthConfig) ([]*v1alpha1.Placement, error) {
	return []*v1alpha1.Placement{obj.Spec.GetPlacement()}, nil
}

func (t *typedParserImpl) ParseFederatedRateLimitConfig(_ context.Context, obj *fed_ratelimit_v1alpha1.FederatedRateLimitConfig) ([]*v1alpha1.Placement, error) {
	return []*v1alpha1.Placement{obj.Spec.GetPlacement()}, nil
}

func (t *typedParserImpl) ParseFederatedUpstream(_ context.Context, obj *fed_gloo_v1.FederatedUpstream) ([]*v1alpha1.Placement, error) {
	return []*v1alpha1.Placement{obj.Spec.GetPlacement()}, nil
}

func (t *typedParserImpl) ParseFederatedSettings(_ context.Context, obj *fed_gloo_v1.FederatedSettings) ([]*v1alpha1.Placement, error) {
	return []*v1alpha1.Placement{obj.Spec.GetPlacement()}, nil
}

func (t *typedParserImpl) ParseFederatedUpstreamGroup(_ context.Context, obj *fed_gloo_v1.FederatedUpstreamGroup) ([]*v1alpha1.Placement, error) {
	return []*v1alpha1.Placement{obj.Spec.GetPlacement()}, nil
}

func (t *typedParserImpl) ParseFederatedGateway(_ context.Context, obj *fed_gateway_v1.FederatedGateway) ([]*v1alpha1.Placement, error) {
	return []*v1alpha1.Placement{obj.Spec.GetPlacement()}, nil
}

func (t *typedParserImpl) ParseFederatedVirtualService(_ context.Context, obj *fed_gateway_v1.FederatedVirtualService) ([]*v1alpha1.Placement, error) {
	return []*v1alpha1.Placement{obj.Spec.GetPlacement()}, nil
}

func (t *typedParserImpl) ParseFederatedRouteTable(_ context.Context, obj *fed_gateway_v1.FederatedRouteTable) ([]*v1alpha1.Placement, error) {
	return []*v1alpha1.Placement{obj.Spec.GetPlacement()}, nil
}
