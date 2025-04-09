package utils

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/solo-io/gloo/projects/gateway2/query"
	"github.com/solo-io/gloo/projects/gateway2/translator/plugins"
	"github.com/solo-io/gloo/projects/gateway2/utils"
	"github.com/solo-io/gloo/projects/gateway2/wellknown"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	skv2corev1 "github.com/solo-io/skv2/pkg/api/core.skv2.solo.io/v1"
	wrapperspb "google.golang.org/protobuf/types/known/wrapperspb"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	gwv1 "sigs.k8s.io/gateway-api/apis/v1"
	apixv1a1 "sigs.k8s.io/gateway-api/apisx/v1alpha1"
)

var (
	ErrUnexpectedListenerType = errors.New("unexpected listener type")
	ErrUnexpectedListener     = func(l *gloov1.Listener) error {
		return fmt.Errorf("%w: expected AggregateListener, got %T", ErrUnexpectedListenerType, l.GetListenerType())
	}

	ListenerTargetRefGVKs = []schema.GroupVersionKind{
		wellknown.GatewayGVK,
		wellknown.XListenerSetGVK,
	}
)

// FindAppliedRouteFilters finds all instances of the supplied filterTypes for the Rule supplied in the RouteContext.
// Should only be used for plugins that support multiple filters as part of a single Rule
func FindAppliedRouteFilters(
	routeCtx *plugins.RouteContext,
	filterTypes ...gwv1.HTTPRouteFilterType,
) []gwv1.HTTPRouteFilter {
	if routeCtx.Rule == nil {
		return nil
	}
	var appliedFilters []gwv1.HTTPRouteFilter
	for _, filter := range routeCtx.Rule.Filters {
		for _, filterType := range filterTypes {
			if filter.Type == filterType {
				appliedFilters = append(appliedFilters, filter)
			}
		}
	}
	return appliedFilters
}

// FindAppliedRouteFilter finds the first instance of the filterType supplied in the Rule being processed.
// Returns nil if the Rule doesn't contain a filter of the provided Type
func FindAppliedRouteFilter(
	routeCtx *plugins.RouteContext,
	filterType gwv1.HTTPRouteFilterType,
) *gwv1.HTTPRouteFilter {
	if routeCtx.Rule == nil {
		return nil
	}
	// TODO: check full Filter list for duplicates and error?
	for _, filter := range routeCtx.Rule.Filters {
		if filter.Type == filterType {
			return &filter
		}
	}
	return nil
}

// FindExtensionRefFilters returns a list ExtensionRef filters that
// references the supplied GroupKind in the Rule being processed.
// Returns nil if the Rule doesn't contain a matching ExtensionRef filter
func FindExtensionRefFilters(
	rule *gwv1.HTTPRouteRule,
	gk schema.GroupKind,
) []gwv1.HTTPRouteFilter {
	if rule == nil {
		return nil
	}

	var filters []gwv1.HTTPRouteFilter
	// TODO: check full Filter list for duplicates and error?
	for _, filter := range rule.Filters {
		if filter.Type == gwv1.HTTPRouteFilterExtensionRef {
			if filter.ExtensionRef.Group == gwv1.Group(gk.Group) && filter.ExtensionRef.Kind == gwv1.Kind(gk.Kind) {
				filters = append(filters, filter)
			}
		}
	}
	return filters
}

var ErrTypesNotEqual = fmt.Errorf("types not equal")

// GetExtensionRefObj uses the provided query engine to retrieve an ExtensionRef object
// and return the object of the same type as the type parameter.
// An error will be returned if the Get was unsuccessful or if the type parameter was not correct.
// A nil error indicates success and `obj` should be usable as normal.
func GetExtensionRefObj[T client.Object](
	ctx context.Context,
	route *gwv1.HTTPRoute,
	queries query.GatewayQueries,
	extensionRef *gwv1.LocalObjectReference,
) (T, error) {
	return GetExtensionRefObjFrom[T](ctx, queries.ObjToFrom(route), queries, extensionRef)
}

func GetExtensionRefObjFrom[T client.Object](
	ctx context.Context,
	from query.From,
	queries query.GatewayQueries,
	extensionRef *gwv1.LocalObjectReference,
) (T, error) {
	var t T
	localObj, err := queries.GetLocalObjRef(ctx, from, *extensionRef)
	if err != nil {
		return t, err
	}

	typed, ok := localObj.(T)
	if !ok {
		return t, fmt.Errorf(
			"%w: generic object typeOf: '%T' localObj typeOf: '%T'",
			ErrTypesNotEqual, t, localObj,
		)
	}
	return typed, nil
}

// PolicyWithSectionedTargetRefs is a wrapper type to represent policy objects
// that attach via TargetRefWtihSectionName
type PolicyWithSectionedTargetRefs[T client.Object] interface {
	GetTargetRefs() []*skv2corev1.PolicyTargetReferenceWithSectionName
	GetObject() T
}

// GetPrioritizedListenerPolicies accepts a slice of Gateway-attached policies (that may explicitly
// target a specific Listener and returns a slice of these policies (or a subset) resources.
// The returned policy list is sorted by specificity in the order of
//
// ListenerSet targetRefs:
//  1. older with section name
//  2. newer with section name
//  3. older without section name
//  4. newer without section name
//
// Gateway targetRefs:
//  1. older with section name
//  2. newer with section name
//  3. older without section name
//  4. newer without section name
//
// This function will process all targetRefs for each policy
func GetPrioritizedListenerPolicies[T client.Object](
	items []PolicyWithSectionedTargetRefs[T],
	listener *gwv1.Listener,
	parentGwName string,
	listenerSet *apixv1a1.XListenerSet,
) []T {
	listenerSetName := ""
	if listenerSet != nil {
		listenerSetName = listenerSet.GetName()
	}

	// gw - gateway, ls - listener set
	var gwOptsWithSectionName, gwOptsWithoutSectionName, lsOptsWithSectionName, lsOptsWithoutSectionName []T
	for i := range items {
		item := items[i]

		// Loop over all targetRefs and check if any have a section name
		appendOptsWithoutSectionName := false
		for _, targetRef := range item.GetTargetRefs() {
			// Check that this is the right gw
			gwMatch := (targetRef.GetGroup() == gwv1.GroupName && targetRef.GetKind() == wellknown.GatewayKind && targetRef.GetName() == parentGwName)
			lsMatch := (targetRef.GetGroup() == apixv1a1.GroupName && targetRef.GetKind() == wellknown.XListenerSetKind && targetRef.GetName() == listenerSetName)
			targetRefMatch := gwMatch || lsMatch

			if !targetRefMatch {
				continue
			}

			sectionName := targetRef.GetSectionName()

			if sectionName != nil && sectionName.GetValue() != "" {
				// we have a section name, now check if it matches the specific listener provided
				if sectionName.GetValue() == string(listener.Name) {
					switch {
					case gwMatch:
						gwOptsWithSectionName = append(gwOptsWithSectionName, item.GetObject())
					case lsMatch:
						lsOptsWithSectionName = append(lsOptsWithSectionName, item.GetObject())
					default:
						// Given the current implementation of the targetRef, this can never happen
						// panic to alert the developer if they are changing the targetRef match logic and forget to update this switch
						panic(fmt.Sprintf("unhandled case when matching targetRef: %v", targetRef))
					}
				}
			} else {
				appendOptsWithoutSectionName = true
			}

			if appendOptsWithoutSectionName {
				// attach all matched items that do not have a section name and let the caller be discerning
				switch {
				case gwMatch:
					gwOptsWithoutSectionName = append(gwOptsWithoutSectionName, item.GetObject())
				case lsMatch:
					lsOptsWithoutSectionName = append(lsOptsWithoutSectionName, item.GetObject())
				default:
					// Given the current implementation of the targetRef, this can never happen
					// panic to alert the developer if they are changing the targetRef match logic and forget to update this switch
					panic(fmt.Sprintf("unhandled case when matching targetRef: %v", targetRef))
				}
			}
		}
	}

	// this can happen if the policy list only contains items targeting other Listeners by section name
	if len(gwOptsWithoutSectionName)+len(gwOptsWithSectionName)+len(lsOptsWithSectionName)+len(lsOptsWithoutSectionName) == 0 {
		return nil
	}

	utils.SortByCreationTime(gwOptsWithSectionName)
	utils.SortByCreationTime(gwOptsWithoutSectionName)
	utils.SortByCreationTime(lsOptsWithSectionName)
	utils.SortByCreationTime(lsOptsWithoutSectionName)
	sortedPolicies := append(append(append(lsOptsWithSectionName, lsOptsWithoutSectionName...), gwOptsWithSectionName...), gwOptsWithoutSectionName...)
	return sortedPolicies
}

// policyTargetReference is an interface that represents a policy target reference, and is used to consolidate
// code for indexing policy target references with and without a section name.
type policyTargetReference interface {
	GetGroup() string
	GetKind() string
	GetNamespace() *wrapperspb.StringValue
	GetName() string
}

type NamespacedNameKind struct {
	Name      string
	Namespace string
	Kind      string
}

const NamespacedNameKindSeparator = '/'

func (n NamespacedNameKind) String() string {
	return strings.Join([]string{n.Namespace, n.Name, n.Kind}, string(NamespacedNameKindSeparator))
}

// Different index functions for different types of target references
type indexFunction[T policyTargetReference] func(targetRef T, namespace string) string

func indexTargetRefsNnk[T policyTargetReference](targetRef T, defaultNamespace string) string {
	ns := targetRef.GetNamespace().GetValue()
	if ns == "" {
		ns = defaultNamespace
	}

	targetNnk := NamespacedNameKind{
		Namespace: ns,
		Name:      targetRef.GetName(),
		Kind:      targetRef.GetKind(),
	}
	return targetNnk.String()
}

func indexTargetRefsNns[T policyTargetReference](targetRef T, defaultNamespace string) string {
	ns := targetRef.GetNamespace().GetValue()
	if ns == "" {
		ns = defaultNamespace
	}

	targetNn := types.NamespacedName{
		Namespace: ns,
		Name:      targetRef.GetName(),
	}
	return targetNn.String()
}

// indexTargetRefs indexes a list of policy target references by namespace and name using the provided indexer function
func indexTargetRefs[T policyTargetReference](targetRefs []T, namespace string, gvks []schema.GroupVersionKind, indexer indexFunction[T]) []string {
	var res []string

	if len(targetRefs) == 0 {
		return res
	}

	foundNns := map[string]any{}

	for _, targetRef := range targetRefs {
		matchGroupKind := false
		for _, kind := range gvks {
			if targetRef.GetGroup() == kind.Group && targetRef.GetKind() == kind.Kind {
				matchGroupKind = true
				break
			}
		}

		if !matchGroupKind {
			continue
		}

		foundNns[indexer(targetRef, namespace)] = struct{}{}
	}

	for k := range foundNns {
		res = append(res, k)
	}

	return res

}

// IndexTargetRefsNnk indexes a list of policy target references by namespace, name, and kind.
// The kinds parameter is a list of GroupVersionKinds that are allowed to be indexed, though version is ignored.
func IndexTargetRefsNnk[T policyTargetReference](targetRefs []T, namespace string, gvks []schema.GroupVersionKind) []string {
	return indexTargetRefs(targetRefs, namespace, gvks, indexTargetRefsNnk[T])
}

// IndexTargetRefsNns indexes a list of policy target references by namespace and name.
// The kinds parameter is a list of GroupVersionKinds that are allowed to be indexed, though version is ignored.
func IndexTargetRefsNns[T policyTargetReference](targetRefs []T, namespace string, gvks []schema.GroupVersionKind) []string {
	return indexTargetRefs(targetRefs, namespace, gvks, indexTargetRefsNns[T])
}
