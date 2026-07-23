package query

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/rotisserie/eris"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	gwv1 "sigs.k8s.io/gateway-api/apis/v1"

	sologatewayv1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	solokubev1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1/kube/apis/gateway.solo.io/v1"
	gwquery "github.com/solo-io/gloo/projects/gateway2/query"
	utils "github.com/solo-io/gloo/projects/gateway2/translator/plugins/utils"
	gwutils "github.com/solo-io/gloo/projects/gateway2/utils"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	glooutils "github.com/solo-io/gloo/projects/gloo/pkg/utils"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
)

var routeOptionGK = schema.GroupKind{
	Group: sologatewayv1.RouteOptionGVK.Group,
	Kind:  sologatewayv1.RouteOptionGVK.Kind,
}

type RouteOptionQueries interface {
	// GetRouteOptionForRouteRule returns the RouteOption attached to the given route and rule.
	//
	// It performs a merge of the ExtensionRef filter attachment and targetRef based attachment
	// while giving priority to the ExtensionRef filter attachment.
	// A lower priority RouteOption may augment the top-level fields in a higher priority RouteOption,
	// but can never override them.
	//
	// When multiple RouteOptions are attached to the route via targetRef, only the earliest created
	// resource is considered for the merge.
	//
	// It returns the merged RouteOption, a list of sources corresponding to the merge, and an error if one occurs.
	GetRouteOptionForRouteRule(
		ctx context.Context,
		route types.NamespacedName,
		rule *gwv1.HTTPRouteRule,
		gwQueries gwquery.GatewayQueries,
	) (*solokubev1.RouteOption, []*gloov1.SourceMetadata_SourceRef, error)
}

type routeOptionQueries struct {
	c client.Client
}

func NewQuery(c client.Client) RouteOptionQueries {
	return &routeOptionQueries{c}
}

func (r *routeOptionQueries) GetRouteOptionForRouteRule(
	ctx context.Context,
	route types.NamespacedName,
	rule *gwv1.HTTPRouteRule,
	gwQueries gwquery.GatewayQueries,
) (*solokubev1.RouteOption, []*gloov1.SourceMetadata_SourceRef, error) {
	filterAttachments, err := lookupFilterAttachments(ctx, route, rule, gwQueries)
	if err != nil {
		return nil, nil, err
	}

	var list solokubev1.RouteOptionList
	if err := r.c.List(
		ctx,
		&list,
		client.MatchingFieldsSelector{Selector: fields.OneTermEqualSelector(RouteOptionTargetField, route.String())},
		client.InNamespace(route.Namespace),
	); err != nil {
		return nil, nil, err
	}

	// Build the ordered list of candidate RouteOptions: ExtensionRef filter
	// attachments take priority over targetRef attachments, and the earliest
	// created targetRef resource wins over later ones (handled by SortByCreationTime).
	candidates := make([]*solokubev1.RouteOption, 0, len(filterAttachments)+len(list.Items))
	candidates = append(candidates, filterAttachments...)
	targetRefs := make([]*solokubev1.RouteOption, len(list.Items))
	for i := range list.Items {
		targetRefs[i] = &list.Items[i]
	}
	gwutils.SortByCreationTime(targetRefs)
	candidates = append(candidates, targetRefs...)

	sources, merged := mergeCandidateRouteOptions(candidates)
	return nilOptionIfEmpty(merged), sources, nil
}

// mergeCandidateRouteOptions merges the given RouteOptions in priority order
// (earlier candidates win) into a single RouteOption.
//
// It uses copy-on-write to avoid deep-cloning the (potentially large)
// RouteOptions tree on every route rule: when only a single candidate
// contributes options, that candidate's Options are returned by reference
// rather than cloned. This is the dominant case (one RouteOption attached to a
// route) and previously accounted for a large share of translation heap because
// the full transformation template tree was deep-copied per route rule.
// See https://github.com/solo-io/solo-projects/issues/8802.
//
// MUTATION SAFETY: when a single candidate is returned by reference, the result
// aliases a resource owned by the informer cache and MUST NOT be mutated in
// place. The translation pipeline upholds this: the routeoptions plugin only
// ever passes the result as the read-only src of MergeRouteOptionsWithOverrides,
// which clones it (when the output route has no options yet) or copies top-level
// field pointers into a separate destination, before assigning it to the output
// route. As soon as a second candidate must be merged, the base is cloned once
// so the in-place field merge never touches a cached resource.
func mergeCandidateRouteOptions(candidates []*solokubev1.RouteOption) ([]*gloov1.SourceMetadata_SourceRef, *solokubev1.RouteOption) {
	var sources []*gloov1.SourceMetadata_SourceRef
	merged := &solokubev1.RouteOption{}

	// owned tracks whether merged.Spec.Options is a copy we own (and may mutate)
	// or a reference to a cached resource that must remain read-only.
	owned := false
	for _, opt := range candidates {
		src := opt.Spec.GetOptions()
		if src == nil {
			continue
		}

		if merged.Spec.GetOptions() == nil {
			// First contributing candidate: share by reference, do not clone.
			merged.Spec.Options = src
			sources = append(sources, routeOptionToSourceRef(opt))
			continue
		}

		// A second candidate needs to be merged. ShallowMergeRouteOptions merges
		// src into dst in place, so clone the aliased base exactly once before the
		// first such merge to avoid mutating the cached resource.
		if !owned {
			merged.Spec.Options = merged.Spec.GetOptions().Clone().(*gloov1.RouteOptions)
			owned = true
		}

		optionUsed := false
		merged.Spec.Options, optionUsed = glooutils.ShallowMergeRouteOptions(merged.Spec.GetOptions(), src)
		if optionUsed {
			sources = append(sources, routeOptionToSourceRef(opt))
		}
	}

	return sources, merged
}

func nilOptionIfEmpty(opt *solokubev1.RouteOption) *solokubev1.RouteOption {
	if opt == nil || opt.Spec.GetOptions() == nil {
		return nil
	}
	return opt
}

// lookupFilterAttachments returns the RouteOptions attached to the route via ExtensionRef filters on the route's rule
func lookupFilterAttachments(
	ctx context.Context,
	route types.NamespacedName,
	rule *gwv1.HTTPRouteRule,
	gwQueries gwquery.GatewayQueries,
) ([]*solokubev1.RouteOption, error) {
	if rule == nil {
		return nil, nil
	}

	filters := utils.FindExtensionRefFilters(rule, routeOptionGK)
	if filters == nil {
		return nil, nil
	}

	var out []*solokubev1.RouteOption
	var multiErr *multierror.Error
	extLookup := extensionRefLookup{namespace: route.Namespace}
	for _, filter := range filters {
		routeOption, err := utils.GetExtensionRefObjFrom[*solokubev1.RouteOption](ctx, extLookup, gwQueries, filter.ExtensionRef)
		if err != nil {
			// If the filter is not found, report a specific error so that it can reflect more
			// clearly on the status of the HTTPRoute.
			if apierrors.IsNotFound(err) {
				multiErr = multierror.Append(multiErr, errFilterNotFound(route.Namespace, &filter))
			} else {
				multiErr = multierror.Append(multiErr, err)
			}
			continue
		}
		out = append(out, routeOption)
	}

	return out, multiErr.ErrorOrNil()
}

type extensionRefLookup struct {
	namespace string
}

func (e extensionRefLookup) GroupKind() (metav1.GroupKind, error) {
	return metav1.GroupKind{
		Group: routeOptionGK.Group,
		Kind:  routeOptionGK.Kind,
	}, nil
}

func (e extensionRefLookup) Namespace() string {
	return e.namespace
}

func errFilterNotFound(namespace string, filter *gwv1.HTTPRouteFilter) error {
	return eris.Errorf(
		"extensionRef '%s' of type %s.%s in namespace '%s' not found",
		filter.ExtensionRef.Group,
		filter.ExtensionRef.Kind,
		filter.ExtensionRef.Name,
		namespace,
	)
}

func routeOptionToSourceRef(opt *solokubev1.RouteOption) *gloov1.SourceMetadata_SourceRef {
	return &gloov1.SourceMetadata_SourceRef{
		ResourceRef: &core.ResourceRef{
			Name:      opt.GetName(),
			Namespace: opt.GetNamespace(),
		},
		ResourceKind:       routeOptionGK.Kind,
		ObservedGeneration: opt.GetGeneration(),
	}
}
