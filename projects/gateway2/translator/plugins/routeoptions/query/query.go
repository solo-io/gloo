package query

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/rotisserie/eris"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	gwv1 "sigs.k8s.io/gateway-api/apis/v1"

	sologatewayv1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	solokubev1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1/kube/apis/gateway.solo.io/v1"
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
	//
	// MEMORY/MUTABILITY CONTRACT: to avoid deep-copying potentially large RouteOptions for every
	// route rule on every translation (solo-io/solo-projects#8802), the lookups pass
	// client.UnsafeDisableDeepCopy and the merged RouteOption shares its options sub-messages with
	// the objects in the client's cache. The merged options are a distinct top-level message, so
	// callers may reassign its top-level fields, but must NEVER mutate nested messages, slices, or
	// maps reachable from it: those are shared with the informer cache and with every other route
	// referencing the same RouteOption.
	GetRouteOptionForRouteRule(
		ctx context.Context,
		route types.NamespacedName,
		rule *gwv1.HTTPRouteRule,
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
) (*solokubev1.RouteOption, []*gloov1.SourceMetadata_SourceRef, error) {
	var sources []*gloov1.SourceMetadata_SourceRef
	merged := &solokubev1.RouteOption{}

	// mergeAttachment folds a single RouteOption attachment into the accumulated `merged` result,
	// recording it as a source if any of its fields were used.
	//
	// The first attachment seeds `merged` with a shallow copy (sharing the attachment's immutable
	// sub-messages by pointer) rather than a deep clone. Deep-cloning the first attachment per route
	// is what dominated translation heap, since every route referencing the same RouteOption received
	// its own deep copy of identical (and often large) transformation templates. `merged.Spec.Options`
	// is a distinct top-level message per route, so downstream route plugins can still reassign its
	// top-level fields safely; they must not mutate the shared sub-messages in place.
	mergeAttachment := func(opt *solokubev1.RouteOption) {
		optionUsed := false
		if merged.Spec.GetOptions() == nil {
			if src := opt.Spec.GetOptions(); src != nil {
				merged.Spec.Options = glooutils.ShallowCopyRouteOptions(src)
				optionUsed = true
			}
		} else {
			merged.Spec.Options, optionUsed = glooutils.ShallowMergeRouteOptions(merged.Spec.GetOptions(), opt.Spec.GetOptions())
		}
		if optionUsed {
			sources = append(sources, routeOptionToSourceRef(opt))
		}
	}

	filterAttachments, err := r.lookupFilterAttachments(ctx, route, rule)
	if err != nil {
		return nil, nil, err
	}
	for _, opt := range filterAttachments {
		mergeAttachment(opt)
	}

	var list solokubev1.RouteOptionList
	if err := r.c.List(
		ctx,
		&list,
		client.MatchingFieldsSelector{Selector: fields.OneTermEqualSelector(RouteOptionTargetField, route.String())},
		client.InNamespace(route.Namespace),
		// Do not deep-copy the matching RouteOptions out of the cache on every call: this query
		// runs for every route rule on every translation, and per-call copies defeat the
		// sub-message sharing that keeps translation heap bounded when many routes reference the
		// same RouteOption (solo-io/solo-projects#8802). The returned objects are shared with the
		// cache and must be treated as read-only, per the contract on RouteOptionQueries.
		client.UnsafeDisableDeepCopy,
	); err != nil {
		return nil, nil, err
	}

	if len(list.Items) == 0 {
		return nilOptionIfEmpty(merged), sources, nil
	}

	out := make([]*solokubev1.RouteOption, len(list.Items))
	for i := range list.Items {
		out[i] = &list.Items[i]
	}
	gwutils.SortByCreationTime(out)
	for _, opt := range out {
		mergeAttachment(opt)
	}

	return nilOptionIfEmpty(merged), sources, nil
}

func nilOptionIfEmpty(opt *solokubev1.RouteOption) *solokubev1.RouteOption {
	if opt == nil || opt.Spec.GetOptions() == nil {
		return nil
	}
	return opt
}

// lookupFilterAttachments returns the RouteOptions attached to the route via ExtensionRef filters on the route's rule.
// ExtensionRefs are local object references, so the lookup is always in the route's namespace.
func (r *routeOptionQueries) lookupFilterAttachments(
	ctx context.Context,
	route types.NamespacedName,
	rule *gwv1.HTTPRouteRule,
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
	for _, filter := range filters {
		routeOption := &solokubev1.RouteOption{}
		err := r.c.Get(
			ctx,
			types.NamespacedName{Namespace: route.Namespace, Name: string(filter.ExtensionRef.Name)},
			routeOption,
			// Shared with the cache and read-only, same as the List below; see the contract on
			// RouteOptionQueries.
			client.UnsafeDisableDeepCopy,
		)
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
