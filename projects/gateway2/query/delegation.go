package query

import (
	"context"
	"path"
	"reflect"
	"slices"
	"strings"

	"github.com/solo-io/gloo/projects/gateway2/wellknown"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	gwv1 "sigs.k8s.io/gateway-api/apis/v1"
)

// inheritMatcherAnnotation is the annotation used on an child HTTPRoute that
// participates in a delegation chain to indicate that child route should inherit
// the route matcher from the parent route.
const inheritMatcherAnnotation = "delegation.gateway.solo.io/inherit-parent-matcher"

// GetDelegatedRoutes returns the child routes that match the specified selector using
// backend ref and the parent route matcher.
//
// A child route may other match a parent rule matcher or may opt to use inherit the parent's matcher.
// In the case of rule inheritance, a child route does not need to match the parent matcher. By default.
// a child route must match the parent matcher.
//
// A child route matches a parent if all of the following conditions are met:
//
// 1. The child route matches the route selector specified by the parent's backendRef
//
// 2.The child route has a rule that matches the parent's route matcher:
//   - The child route's path must contain the parent's path as a prefix
//   - The child route's headers must be a superset of the parent's headers
//   - The child route's query parameters must be a superset of the parent's query parameters
//
// If a child route's rule does not match the given parent match, it is not included in the route returned.
func (r *gatewayQueries) GetDelegatedRoutes(
	ctx context.Context,
	backendRef gwv1.BackendObjectReference,
	parentMatch gwv1.HTTPRouteMatch,
	parentRef types.NamespacedName,
) ([]gwv1.HTTPRoute, error) {
	if parentMatch.Path == nil || parentMatch.Path.Type == nil || *parentMatch.Path.Type != gwv1.PathMatchPathPrefix {
		return nil, nil
	}

	// Lookup the delegated routes; default to parent namespace
	delegatedNs := parentRef.Namespace
	if backendRef.Namespace != nil {
		delegatedNs = string(*backendRef.Namespace)
	}

	var children []gwv1.HTTPRoute
	if string(backendRef.Name) == "" || string(backendRef.Name) == "*" {
		// consider this to be a wildcard
		var hrlist gwv1.HTTPRouteList
		err := r.client.List(ctx, &hrlist, client.InNamespace(delegatedNs))
		if err != nil {
			return nil, err
		}
		children = append(children, hrlist.Items...)
	} else {
		delegatedRef := types.NamespacedName{
			Namespace: delegatedNs,
			Name:      string(backendRef.Name),
		}
		// Lookup the child route
		child := &gwv1.HTTPRoute{}
		err := r.client.Get(ctx, types.NamespacedName{Namespace: delegatedRef.Namespace, Name: delegatedRef.Name}, child)
		if err != nil {
			return nil, err
		}
		children = append(children, *child)
	}

	// Select the child routes that match the parent
	var selected []gwv1.HTTPRoute
	for _, child := range children {
		// Check if the child route is allowed to be delegated to by the parent
		if !isAllowedParent(parentRef, child.Namespace, child.Spec.ParentRefs) {
			continue
		}

		inheritMatcher := shouldInheritMatcher(child)

		// Check if the child route has a prefix that matches the parent.
		// Only rules matching the parent prefix are considered.
		//
		// We use validRules to store the rules in the child route that are valid
		// (matches in the rule match the parent route matcher). If a specific rule
		// in the child is not valid, then we discard it in the final child route
		// returned by this function.
		var validRules []gwv1.HTTPRouteRule
		for i, rule := range child.Spec.Rules {
			var validMatches []gwv1.HTTPRouteMatch

			// If the child route opts to inherit the parent's matcher and it does not specify its own matcher,
			// simply inherit the parent's matcher.
			if inheritMatcher && len(rule.Matches) == 0 {
				validMatches = append(validMatches, parentMatch)
			}

			for _, match := range rule.Matches {
				// When inheriting the parent's matcher, all matches are valid.
				// In this case, the child inherits the parents matcher so we merge
				// the parent's matcher with the child's.
				if inheritMatcher {
					mergeParentChildRouteMatch(&parentMatch, &match)
					validMatches = append(validMatches, match)
					continue
				}

				// Non-inherited matcher delegation requires matching child matcher to parent matcher
				// to delegate from the parent route to the child.
				if ok := isDelegatedRouteMatch(parentMatch, match); ok {
					validMatches = append(validMatches, match)
				}
			}

			// Matchers in this rule match the parent route matcher, so consider the valid matchers on the child,
			// and discard rules on the child that do not match the parent route matcher.
			if len(validMatches) > 0 {
				validRule := child.Spec.Rules[i]
				validRule.Matches = validMatches
				validRules = append(validRules, validRule)
			}
		}
		if len(validRules) > 0 {
			child.Spec.Rules = validRules
			selected = append(selected, child)
		}
	}

	return selected, nil
}

// isAllowedParent returns whether the parent specified by `parentRef` is allowed to delegate
// to the child.
//   - `childNs` is the namespace of the child route.
//   - `childParentRefs` is the list of parent references on the child route. If this is empty, then
//     there are no restrictions on which parents can delegate to this child. If it is not empty,
//     then `parentRef` must be in this list in order for the parent to delegate to the child.
func isAllowedParent(
	parentRef types.NamespacedName,
	childNs string,
	childParentRefs []gwv1.ParentReference,
) bool {
	// no explicit parentRefs, so any parent is allowed
	if len(childParentRefs) == 0 {
		return true
	}

	// validate that the child's parentRefs contains the specified parentRef
	for _, ref := range childParentRefs {
		// default to the child's namespace if not specified
		refNs := childNs
		if ref.Namespace != nil {
			refNs = string(*ref.Namespace)
		}
		// check if the ref matches the desired parentRef
		if ref.Group != nil && *ref.Group == wellknown.GatewayGroup &&
			ref.Kind != nil && *ref.Kind == wellknown.HTTPRouteKind &&
			string(ref.Name) == parentRef.Name &&
			refNs == parentRef.Namespace {
			return true
		}
	}
	return false
}

func isDelegatedRouteMatch(
	parent gwv1.HTTPRouteMatch,
	child gwv1.HTTPRouteMatch,
) bool {
	// Validate path
	if parent.Path == nil || parent.Path.Type == nil || *parent.Path.Type != gwv1.PathMatchPathPrefix {
		return false
	}
	parentPath := *parent.Path.Value
	if child.Path == nil || child.Path.Type == nil {
		return false
	}
	childPath := *child.Path.Value
	if !strings.HasPrefix(childPath, parentPath) {
		return false
	}

	// Validate that the child headers are a superset of the parent headers
	for _, parentHeader := range parent.Headers {
		found := false
		for _, childHeader := range child.Headers {
			if reflect.DeepEqual(parentHeader, childHeader) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	// Validate that the child query parameters are a superset of the parent headers
	for _, parentQuery := range parent.QueryParams {
		found := false
		for _, childQuery := range child.QueryParams {
			if reflect.DeepEqual(parentQuery, childQuery) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	// Validate that the child method matches the parent method
	if parent.Method != nil && (child.Method == nil || *parent.Method != *child.Method) {
		return false
	}

	return true
}

// shouldInheritMatcher returns true if the route indicates that it should inherit
// its parent's matcher.
func shouldInheritMatcher(route gwv1.HTTPRoute) bool {
	val, ok := route.Annotations[inheritMatcherAnnotation]
	if !ok {
		return false
	}
	switch strings.ToLower(val) {
	case "true", "yes", "enabled":
		return true

	default:
		return false
	}
}

// mergeParentChildRouteMatch merges the parent route match into the child.
func mergeParentChildRouteMatch(
	parent *gwv1.HTTPRouteMatch,
	child *gwv1.HTTPRouteMatch,
) {
	if parent == nil || child == nil {
		return
	}

	if child.Path == nil {
		child.Path = &gwv1.HTTPPathMatch{
			Type:  ptr.To(gwv1.PathMatchPathPrefix),
			Value: ptr.To(""),
		}
	}
	child.Path.Value = ptr.To(path.Join(*parent.Path.Value, *child.Path.Value))

	// Inherit parent and child headers and query parameters while augmenting the merge
	// with additions specified on the child
	child.Headers = mergeHeaders(parent.Headers, child.Headers)
	child.QueryParams = mergeQueries(parent.QueryParams, child.QueryParams)

	// If parent specifies a method, inherit it
	if parent.Method != nil {
		child.Method = ptr.To(*parent.Method)
	}
}

func mergeHeaders(
	parent, child []gwv1.HTTPHeaderMatch,
) []gwv1.HTTPHeaderMatch {
	merged := make(map[gwv1.HTTPHeaderName]gwv1.HTTPHeaderMatch)
	for _, h := range parent {
		merged[h.Name] = h
	}
	for _, h := range child {
		key := h.Name
		// Only add the child if it does not conflict with the parent
		if _, ok := merged[key]; !ok {
			merged[key] = h
		}
	}
	var result []gwv1.HTTPHeaderMatch
	for _, h := range merged {
		result = append(result, h)
	}
	// Sort for deterministic ordering
	slices.SortFunc(result, func(a, b gwv1.HTTPHeaderMatch) int {
		return strings.Compare(string(a.Name), string(b.Name))
	})
	return result
}

func mergeQueries(
	parent, child []gwv1.HTTPQueryParamMatch,
) []gwv1.HTTPQueryParamMatch {
	merged := make(map[gwv1.HTTPHeaderName]gwv1.HTTPQueryParamMatch)
	for _, h := range parent {
		merged[h.Name] = h
	}
	for _, h := range child {
		key := h.Name
		// Only add the child if it does not conflict with the parent
		if _, ok := merged[key]; !ok {
			merged[key] = h
		}
	}
	var result []gwv1.HTTPQueryParamMatch
	for _, h := range merged {
		result = append(result, h)
	}
	// Sort for deterministic ordering
	slices.SortFunc(result, func(a, b gwv1.HTTPQueryParamMatch) int {
		return strings.Compare(string(a.Name), string(b.Name))
	})
	return result
}
