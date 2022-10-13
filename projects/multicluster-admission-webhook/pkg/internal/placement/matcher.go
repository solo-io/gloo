package placement

import (
	"context"

	"github.com/solo-io/go-utils/contextutils"
	"github.com/solo-io/go-utils/stringutils"
	multicluster_types "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/multicluster.solo.io/v1alpha1/types"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/util/sets"
)

//go:generate mockgen -source ./matcher.go -destination ./mocks/mock_placement.go

/*
	Compute whether the placement of a resource is allowed by a given rule.

	Checks whether the clusters defined by the rule are a superset of the ones on the resource. Followed by the
	same process for the namespaces.

	This function will return an error if the format of any of the lists is incorrect
*/
type Matcher interface {
	Matches(ctx context.Context, resource, rule *multicluster_types.Placement) bool
}

func NewMatcher() Matcher {
	return &placementMatcher{}
}

type placementMatcher struct{}

func (p *placementMatcher) Matches(ctx context.Context, resource, rule *multicluster_types.Placement) bool {
	logger := contextutils.LoggerFrom(contextutils.WithLoggerValues(ctx,
		zap.Strings("rule_clusters", rule.GetClusters()),
		zap.Strings("resource_clusters", resource.GetClusters()),
		zap.Strings("rule_namespaces", rule.GetNamespaces()),
		zap.Strings("resource_namespaces", resource.GetNamespaces()),
	))

	if !IsValid(rule) {
		logger.Debug("rule placement is invalid")
		return false
	}
	if !IsValid(resource) {
		logger.Debug("resource placement is invalid")
		return false
	}

	if !IsClusterWildcard(rule) {

		// If the resource has a wildcard, and the rule does not, reject
		if IsClusterWildcard(resource) {
			logger.Debug("resource has a wildcard cluster list, but rule does not")
			return false
		}

		// If the resource intends to modify clusters not included in the rule, disallow
		if !sets.NewString(rule.GetClusters()...).IsSuperset(sets.NewString(resource.GetClusters()...)) {
			logger.Debug("rule cluster list is not a superset of the resource cluster list")
			return false
		}
	}

	if !IsNamespaceWildcard(rule) {

		// If the resource has a wildcard, and the rule does not, reject
		if IsNamespaceWildcard(resource) {
			logger.Debug("resource has a wildcard namespace list, but resource does not")
			return false
		}

		// If the resource intends to modify namespaces not included in the rule, disallow
		if !sets.NewString(rule.GetNamespaces()...).IsSuperset(sets.NewString(resource.GetNamespaces()...)) {
			logger.Debug("rule namespace list is not a superset of the resource namespace list")
			return false
		}
	}

	return true
}

func IsValid(placement *multicluster_types.Placement) bool {
	return len(placement.GetClusters()) != 0 && len(placement.GetNamespaces()) != 0
}

func IsClusterWildcard(placement *multicluster_types.Placement) bool {
	return stringutils.ContainsString("*", placement.GetClusters())
}

func IsNamespaceWildcard(placement *multicluster_types.Placement) bool {
	return stringutils.ContainsString("*", placement.GetNamespaces())
}
