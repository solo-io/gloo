package utils

import (
	"reflect"
	"strings"

	"k8s.io/apimachinery/pkg/types"
	gwv1 "sigs.k8s.io/gateway-api/apis/v1"

	"github.com/solo-io/gloo/projects/gateway2/wellknown"
)

// ChildRouteCanAttachToParentRef returns a boolean indicating whether the given delegatee/child
// route can attach to a parent referenced by its NamespacedName.
//
// A delegatee route can attach to a parent if either of the following conditions are true:
//   - the child does not specify ParentRefs (implicit attachment)
//   - the child has an HTTPRoute ParentReference that matches parentRef
func ChildRouteCanAttachToParentRef(
	route *gwv1.HTTPRoute,
	parentRef types.NamespacedName,
) bool {
	if route == nil {
		return false
	}

	childParentRefs := route.Spec.ParentRefs

	// no explicit parentRefs, so any parent is allowed
	if len(childParentRefs) == 0 {
		return true
	}

	// validate that the child's parentRefs contains the specified parentRef
	for _, ref := range childParentRefs {
		// default to the child's namespace if not specified
		refNs := route.Namespace
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

// ShouldChildRouteInheritParentMatcher returns true if the child HTTPRoute should inherit the parent HTTPRoute's matcher
func ShouldChildRouteInheritParentMatcher(route *gwv1.HTTPRoute) bool {
	if route == nil {
		return false
	}
	val, ok := route.Annotations[wellknown.InheritMatcherAnnotation]
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

// IsDelegatedRouteMatch returns true if the child HTTPRouteMatch matches (is a subset) of the parent HTTPRouteMatch.
func IsDelegatedRouteMatch(
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
