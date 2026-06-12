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
	// client.UnsafeDisableDeepCopy and the query interns one deep copy per unique RouteOption
	// for its own lifetime (a fresh query is constructed per translation pass and retained with
	// that pass's output; see NewQuery). The merged RouteOption shares its options sub-messages
	// with that interned copy, never with the informer cache, so the cache cannot be corrupted
	// through the merged result. The merged options are a distinct top-level message, so callers
	// may reassign its top-level fields; they still must not mutate nested messages, slices, or
	// maps reachable from it, since those are shared with every other route referencing the same
	// RouteOption in this pass.
	GetRouteOptionForRouteRule(
		ctx context.Context,
		route types.NamespacedName,
		rule *gwv1.HTTPRouteRule,
	) (*solokubev1.RouteOption, []*gloov1.SourceMetadata_SourceRef, error)
}

type routeOptionQueries struct {
	c client.Client

	// interned holds this query's one deep copy per unique RouteOption. An entry is replaced
	// when the cached RouteOption's resourceVersion moves, so a lookup can never be served stale
	// options and the map stays bounded by the number of RouteOptions even if the query outlives
	// the translation pass it was built for. The interned copies are what keep the informer
	// cache unreachable from translation output while still sharing one copy across all routes
	// that reference the same RouteOption (solo-io/solo-projects#8802). Not safe for concurrent
	// use: route plugins run sequentially within a pass.
	interned map[types.NamespacedName]internedRouteOption
}

// internedRouteOption is one RouteOption's deep-copied options plus the resourceVersion they
// were copied at.
type internedRouteOption struct {
	resourceVersion string
	options         *gloov1.RouteOptions
}

// NewQuery returns a RouteOptionQueries meant to live for a single translation pass: the proxy
// syncer builds a fresh plugin registry — and with it a fresh query — per pass, and retains it
// with that pass's output for status syncing. The query's interned RouteOption copies are
// retained along with it, which is what bounds translation memory at one copy per unique
// RouteOption per pass instead of one per route (solo-io/solo-projects#8802).
func NewQuery(c client.Client) RouteOptionQueries {
	return &routeOptionQueries{
		c:        c,
		interned: map[types.NamespacedName]internedRouteOption{},
	}
}

// internedOptions returns this query's private deep copy of the RouteOption's options, cloning
// on first sight or when the cached object's resourceVersion has moved since the copy was
// taken. opt is shared with the informer cache (the lookups disable deep copies) and is only
// ever read, never written to.
func (r *routeOptionQueries) internedOptions(opt *solokubev1.RouteOption) *gloov1.RouteOptions {
	src := opt.Spec.GetOptions()
	if src == nil {
		return nil
	}
	key := types.NamespacedName{Namespace: opt.GetNamespace(), Name: opt.GetName()}
	if entry, ok := r.interned[key]; ok && entry.resourceVersion == opt.GetResourceVersion() {
		return entry.options
	}
	copied := src.Clone().(*gloov1.RouteOptions)
	r.interned[key] = internedRouteOption{
		resourceVersion: opt.GetResourceVersion(),
		options:         copied,
	}
	return copied
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
	// The merge reads from the query's interned copy of each attachment, never from the
	// cache-shared object itself, and the first attachment seeds `merged` with a shallow copy
	// (sharing the interned copy's sub-messages by pointer) rather than a deep clone.
	// Deep-cloning the first attachment per route is what dominated translation heap, since every
	// route referencing the same RouteOption received its own deep copy of identical (and often
	// large) transformation templates. `merged.Spec.Options` is a distinct top-level message per
	// route, so downstream route plugins can still reassign its top-level fields safely; they
	// must not mutate the shared sub-messages in place.
	mergeAttachment := func(opt *solokubev1.RouteOption) {
		options := r.internedOptions(opt)
		if options == nil {
			return
		}
		optionUsed := false
		if merged.Spec.GetOptions() == nil {
			merged.Spec.Options = glooutils.ShallowCopyRouteOptions(options)
			optionUsed = true
		} else {
			merged.Spec.Options, optionUsed = glooutils.ShallowMergeRouteOptions(merged.Spec.GetOptions(), options)
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
		// Skip the client's per-call deep copy out of the cache: the merge only ever reads this
		// query's interned copies (see internedOptions), so a per-call copy would be pure
		// allocation churn thrown away after each lookup of a query that runs for every route
		// rule on every translation (solo-io/solo-projects#8802). The returned objects are
		// shared with the cache and are only read, never written.
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
			// Shared with the cache and only read, never written, same as the List in
			// GetRouteOptionForRouteRule; the merge reads only this query's interned copies
			// (see internedOptions).
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
