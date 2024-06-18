package query

import (
	"context"
	"fmt"

	"github.com/solo-io/gloo/projects/gateway2/translator/backendref"
	"github.com/solo-io/gloo/projects/gateway2/wellknown"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/controller-runtime/pkg/client"

	gwv1 "sigs.k8s.io/gateway-api/apis/v1"
	v1 "sigs.k8s.io/gateway-api/apis/v1"
)

// HTTPRouteInfo contains pre-resolved backends (Services, Upstreams and delegate HTTPRoutes)
// This allows all querying to happen upfront, and detailed logic for delegation to can be
// as part of translation.
type HTTPRouteInfo struct {
	// HTTPRoute with rules and matches filtered to only those
	// that are valid based on the parent rule.
	gwv1.HTTPRoute

	// ParentRef points to the Gateway (and optionally Listener) or HTTPRoute
	ParentRef *gwv1.ParentReference

	// hostnameOverrides can replace the HTTPRoute hostnames with those that intersect
	// the attached listener's hostname(s).
	// R
	HostnameOverrides []string

	// Backends are pre-resolved here. This list will not contain delegates.
	// Map values are either client.Object or error (errors can be passed to ProcessBackendRef).
	// TODO should we ProcessBackendRef early and put cluster names here?)
	Backends BackendMap[client.Object]
	// Children contains all delegate HTTPRoutes referenced in any rule of this
	// HTTPRoute, keyed by the backend ref for easy lookup.
	// This tree structure can have cyclic references. Check them when recursing through the tree.
	Children BackendMap[[]*HTTPRouteInfo]
}

func (h HTTPRouteInfo) GetName() string {
	return h.HTTPRoute.GetName()
}

func (h HTTPRouteInfo) GetNamespace() string {
	return h.HTTPRoute.GetNamespace()
}

type BackendMap[T any] struct {
	items  map[backendRefKey]T
	errors map[backendRefKey]error
}

func NewBackendMap[T any]() BackendMap[T] {
	return BackendMap[T]{
		items:  make(map[backendRefKey]T),
		errors: make(map[backendRefKey]error),
	}
}

type backendRefKey string

func backendToRefKey(ref gwv1.BackendObjectReference) backendRefKey {
	const delim = "~"
	return backendRefKey(
		string(ptrOrDefault(ref.Group, "")) + delim +
			string(ptrOrDefault(ref.Kind, "")) + delim +
			string(ref.Name) + delim +
			string(ptrOrDefault(ref.Namespace, "")),
	)
}

func ptrOrDefault[T comparable](p *T, fallback T) T {
	if p == nil {
		return fallback
	}
	return *p
}

func (bm BackendMap[T]) Add(backendRef gwv1.BackendObjectReference, value T) {
	key := backendToRefKey(backendRef)
	bm.items[key] = value
}

func (bm BackendMap[T]) AddError(backendRef gwv1.BackendObjectReference, err error) {
	key := backendToRefKey(backendRef)
	bm.errors[key] = err
}

func (bm BackendMap[T]) get(backendRef gwv1.BackendObjectReference, def T) (T, error) {
	key := backendToRefKey(backendRef)
	if err, ok := bm.errors[key]; ok {
		return def, err
	}
	if res, ok := bm.items[key]; ok {
		return res, nil
	}
	return def, ErrUnresolvedReference
}

func (hr *HTTPRouteInfo) Hostnames() []string {
	if len(hr.HostnameOverrides) > 0 {
		return hr.HostnameOverrides
	}
	strs := make([]string, 0, len(hr.Spec.Hostnames))
	for _, v := range hr.Spec.Hostnames {
		strs = append(strs, string(v))
	}
	return strs
}

func (hr *HTTPRouteInfo) GetBackendForRef(backendRef gwv1.BackendObjectReference) (client.Object, error) {
	return hr.Backends.get(backendRef, nil)
}

func (hr *HTTPRouteInfo) GetChildrenForRef(backendRef gwv1.BackendObjectReference) ([]*HTTPRouteInfo, error) {
	return hr.Children.get(backendRef, nil)
}

func (hr *HTTPRouteInfo) Clone() *HTTPRouteInfo {
	if hr == nil {
		return nil
	}
	return &HTTPRouteInfo{
		HTTPRoute: hr.HTTPRoute, // TODO DeepCopy here too?
		ParentRef: hr.ParentRef.DeepCopy(),
		Backends:  hr.Backends,
		Children:  hr.Children,
	}
}

type GatewayHTTPRouteInfo struct {
	ListenerResults map[string][]*HTTPRouteInfo
	ListenerErrors  map[string]error
	RouteErrors     []*RouteError
}

// GetHTTPRouteChains queries for HTTPRoutes that attach directly to a Gateway.
// The returned []HTTPRouteInfo items are each top-level HTTPRoutes, children
// are delgated routes. Delegated children may not be fully-valid based on
// matchers, and are given unmodified. This decouples the querying from the
// more nuanced parts of translation. BackendRefs are followed blindly.
func (r *gatewayQueries) GetHTTPRouteChains(
	ctx context.Context,
	gw *v1.Gateway,
) (GatewayHTTPRouteInfo, error) {
	// TODO inline this or make it private?
	topLevel, err := r.getRoutesForGw(ctx, gw)
	if err != nil {
		return GatewayHTTPRouteInfo{}, err
	}
	res := GatewayHTTPRouteInfo{
		ListenerResults: make(map[string][]*HTTPRouteInfo, len(topLevel.ListenerResults)),
		ListenerErrors:  make(map[string]error, len(topLevel.ListenerResults)),
		RouteErrors:     topLevel.RouteErrors,
	}
	for listener, listenerRes := range topLevel.ListenerResults {
		if listenerRes.Error != nil {
			res.ListenerErrors[listener] = listenerRes.Error
			continue
		}
		for _, route := range listenerRes.Routes {
			route := route // pike: ptr to iterator var
			gwRef := route.ParentRef
			info := &HTTPRouteInfo{
				HTTPRoute:         route.Route,
				HostnameOverrides: route.Hostnames,
				ParentRef:         &gwRef,
				Backends:          r.resolveRouteBackends(ctx, &route.Route),
				Children:          r.getDelegatedChildren(ctx, &route.Route, nil),
			}
			res.ListenerResults[listener] = append(res.ListenerResults[listener], info)
		}
	}
	return res, nil
}

func (r *gatewayQueries) getRoutesForGw(ctx context.Context, gw *gwv1.Gateway) (RoutesForGwResult, error) {
	ret := RoutesForGwResult{
		ListenerResults: map[string]*ListenerResult{},
	}

	nns := types.NamespacedName{
		Namespace: gw.Namespace,
		Name:      gw.Name,
	}

	var hrlist gwv1.HTTPRouteList
	err := r.client.List(ctx, &hrlist, client.MatchingFieldsSelector{Selector: fields.OneTermEqualSelector(HttpRouteTargetField, nns.String())})
	if err != nil {
		return ret, err
	}

	for _, hr := range hrlist.Items {
		refs := getParentRefsForGw(gw, &hr)
		for _, ref := range refs {
			anyRoutesAllowed := false
			anyListenerMatched := false
			anyHostsMatch := false
			for _, l := range gw.Spec.Listeners {
				lr := ret.ListenerResults[string(l.Name)]

				if lr == nil {
					lr = &ListenerResult{}
					ret.ListenerResults[string(l.Name)] = lr
				}

				allowedNs, allowedKinds, err := r.allowedRoutes(gw, &l)
				if err != nil {
					lr.Error = err
					continue
				}

				if !isHttpRouteAllowed(allowedKinds) {
					continue
				}
				if !allowedNs(hr.Namespace) {
					continue
				}
				anyRoutesAllowed = true

				if !parentRefMatchListener(ref, &l) {
					continue
				}
				anyListenerMatched = true

				ok, hostnames := hostnameIntersect(&l, &hr)
				if !ok {
					continue
				}
				lrr := &ListenerRouteResult{
					Route:     hr,
					Hostnames: hostnames,
					ParentRef: ref,
				}
				anyHostsMatch = true
				lr.Routes = append(lr.Routes, lrr)
			}

			if !anyRoutesAllowed {
				ret.RouteErrors = append(ret.RouteErrors, &RouteError{
					Route:     hr,
					ParentRef: ref,
					Error:     Error{E: ErrNotAllowedByListeners, Reason: gwv1.RouteReasonNotAllowedByListeners},
				})
			} else if !anyListenerMatched {
				ret.RouteErrors = append(ret.RouteErrors, &RouteError{
					Route:     hr,
					ParentRef: ref,
					Error:     Error{E: ErrNoMatchingParent, Reason: gwv1.RouteReasonNoMatchingParent},
				})
			} else if !anyHostsMatch {
				ret.RouteErrors = append(ret.RouteErrors, &RouteError{
					Route:     hr,
					ParentRef: ref,
					Error:     Error{E: ErrNoMatchingListenerHostname, Reason: gwv1.RouteReasonNoMatchingListenerHostname},
				})
			}
		}
	}
	return ret, nil
}

func (r *gatewayQueries) allowedRoutes(gw *gwv1.Gateway, l *gwv1.Listener) (func(string) bool, []metav1.GroupKind, error) {
	var allowedKinds []metav1.GroupKind

	switch l.Protocol {
	case gwv1.HTTPSProtocolType:
		fallthrough
	case gwv1.HTTPProtocolType:
		allowedKinds = []metav1.GroupKind{{Kind: wellknown.HTTPRouteKind, Group: gwv1.GroupName}}
	case gwv1.TLSProtocolType:
		fallthrough
	case gwv1.TCPProtocolType:
		allowedKinds = []metav1.GroupKind{{}}
	case gwv1.UDPProtocolType:
		allowedKinds = []metav1.GroupKind{{}}
	default:
		// allow custom protocols to work
		allowedKinds = []metav1.GroupKind{{Kind: wellknown.HTTPRouteKind, Group: gwv1.GroupName}}
	}

	allowedNs := SameNamespace(gw.Namespace)
	if ar := l.AllowedRoutes; ar != nil {
		if ar.Kinds != nil {
			allowedKinds = nil
			for _, k := range ar.Kinds {
				gk := metav1.GroupKind{Kind: string(k.Kind)}
				if k.Group != nil {
					gk.Group = string(*k.Group)
				} else {
					gk.Group = gwv1.GroupName
				}
				allowedKinds = append(allowedKinds, gk)
			}
		}
		if ar.Namespaces != nil && ar.Namespaces.From != nil {
			switch *ar.Namespaces.From {
			case gwv1.NamespacesFromAll:
				allowedNs = AllNamespace()
			case gwv1.NamespacesFromSelector:
				if ar.Namespaces.Selector == nil {
					return nil, nil, fmt.Errorf("selector must be set")
				}
				selector, err := metav1.LabelSelectorAsSelector(ar.Namespaces.Selector)
				if err != nil {
					return nil, nil, err
				}
				allowedNs = r.NamespaceSelector(selector)
			}
		}
	}
	return allowedNs, allowedKinds, nil
}

func (r *gatewayQueries) resolveRouteBackends(ctx context.Context, hr *gwv1.HTTPRoute) BackendMap[client.Object] {
	out := NewBackendMap[client.Object]()
	for _, rule := range hr.Spec.Rules {
		for _, backendRef := range rule.BackendRefs {
			obj, err := r.GetBackendForRef(ctx, r.ObjToFrom(hr), &backendRef.BackendObjectReference)
			if err != nil {
				out.AddError(backendRef.BackendObjectReference, err)
				continue
			}
			out.Add(backendRef.BackendObjectReference, obj)
		}
	}
	return out
}

func (r *gatewayQueries) getDelegatedChildren(
	ctx context.Context,
	parent *gwv1.HTTPRoute,
	visited sets.Set[types.NamespacedName],
) BackendMap[[]*HTTPRouteInfo] {
	if visited == nil {
		visited = sets.New[types.NamespacedName]()
	}
	parentRef := namespacedName(parent)
	visited.Insert(parentRef)

	children := NewBackendMap[[]*HTTPRouteInfo]()
	for _, parentRule := range parent.Spec.Rules {
		var refChildren []*HTTPRouteInfo
		for _, backendRef := range parentRule.BackendRefs {
			if !backendref.RefIsHTTPRoute(backendRef.BackendObjectReference) {
				continue
			}
			referencedRoutes, err := r.fetchChildRoutes(ctx, parent.Namespace, backendRef)
			if err != nil {
				children.AddError(backendRef.BackendObjectReference, err)
				continue
			}
			for _, childRoute := range referencedRoutes {
				childRoute := childRoute // pike: ptr to loop item
				childRef := namespacedName(&childRoute)
				if visited.Has(childRef) {
					err := fmt.Errorf("ignoring child route %s for parent %s: %w", parentRef, childRef, ErrCyclicReference)
					children.AddError(backendRef.BackendObjectReference, err)
					// don't resolve child routes; the entire backendRef is invalid
					break
				}
				routeInfo := &HTTPRouteInfo{
					HTTPRoute: childRoute,
					ParentRef: &gwv1.ParentReference{
						Group:     ptr.To(gwv1.Group(wellknown.GatewayGroup)),
						Kind:      ptr.To(gwv1.Kind(wellknown.HTTPRouteKind)),
						Namespace: ptr.To(gwv1.Namespace(parent.Namespace)),
						Name:      v1.ObjectName(parent.Name),
					},
					Backends: r.resolveRouteBackends(ctx, &childRoute),
					Children: r.getDelegatedChildren(ctx, &childRoute, visited),
				}
				refChildren = append(refChildren, routeInfo)
			}
			children.Add(backendRef.BackendObjectReference, refChildren)
		}
	}
	return children
}

func (r *gatewayQueries) fetchChildRoutes(
	ctx context.Context,
	parentNamespace string,
	backendRef gwv1.HTTPBackendRef,
) ([]gwv1.HTTPRoute, error) {
	delegatedNs := parentNamespace
	if !backendref.RefIsHTTPRoute(backendRef.BackendObjectReference) {
		return nil, nil
	}
	if backendRef.Namespace != nil {
		delegatedNs = string(*backendRef.Namespace)
	}
	var refChildren []gwv1.HTTPRoute
	if string(backendRef.Name) == "" || string(backendRef.Name) == "*" {
		// consider this to be a wildcard
		var hrlist gwv1.HTTPRouteList
		err := r.client.List(ctx, &hrlist, client.InNamespace(delegatedNs))
		if err != nil {
			return nil, err
		}
		refChildren = append(refChildren, hrlist.Items...)
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
		refChildren = append(refChildren, *child)
	}
	if len(refChildren) == 0 {
		return nil, ErrUnresolvedReference
	}
	return refChildren, nil
}

type Namespaced interface {
	GetName() string
	GetNamespace() string
}

func namespacedName(o Namespaced) types.NamespacedName {
	return types.NamespacedName{Name: o.GetName(), Namespace: o.GetNamespace()}
}
