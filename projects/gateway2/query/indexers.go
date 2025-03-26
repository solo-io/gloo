package query

import (
	"errors"

	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	gwv1 "sigs.k8s.io/gateway-api/apis/v1"
	gwv1a2 "sigs.k8s.io/gateway-api/apis/v1alpha2"
	gwv1b1 "sigs.k8s.io/gateway-api/apis/v1beta1"
	apixv1a1 "sigs.k8s.io/gateway-api/apisx/v1alpha1"
	gwxv1a1 "sigs.k8s.io/gateway-api/apisx/v1alpha1"

	"github.com/solo-io/gloo/projects/gateway2/wellknown"
)

const (
	HttpRouteTargetField            = "http-route-target"
	HttpRouteDelegatedLabelSelector = "http-route-delegated-label-selector"
	TcpRouteTargetField             = "tcp-route-target"
	TlsRouteTargetField             = "tls-route-target"
	ReferenceGrantFromField         = "ref-grant-from"
	ListenerSetTargetField          = "listener-set-target"
)

// IterateIndices calls the provided function for each indexable object with the appropriate indexer function.
func IterateIndices(f func(client.Object, string, client.IndexerFunc) error) error {
	return errors.Join(
		f(&gwv1.HTTPRoute{}, HttpRouteTargetField, IndexerByObjType),
		f(&gwv1.HTTPRoute{}, HttpRouteDelegatedLabelSelector, IndexByHTTPRouteDelegationLabelSelector),
		f(&gwv1a2.TCPRoute{}, TcpRouteTargetField, IndexerByObjType),
		f(&gwv1a2.TLSRoute{}, TlsRouteTargetField, IndexerByObjType),
		f(&gwv1b1.ReferenceGrant{}, ReferenceGrantFromField, IndexerByObjType),
		f(&gwxv1a1.XListenerSet{}, ListenerSetTargetField, IndexerByObjType),
	)
}

func matchesGK(pRef gwv1.ParentReference, parentGroup gwv1.Group, parentKind gwv1.Kind) bool {
	if pRef.Group != nil && *pRef.Group != parentGroup {
		return false
	}
	if pRef.Kind != nil && *pRef.Kind != parentKind {
		return false
	}
	return true
}

func isGateway(pRef gwv1.ParentReference) bool {
	return matchesGK(pRef, gwv1.GroupName, wellknown.GatewayKind)
}

func isListenerSet(pRef gwv1.ParentReference) bool {
	return matchesGK(pRef, apixv1a1.GroupName, wellknown.XListenerSetKind)
}

// IndexerByObjType indexes objects based on the provided object type. The following object types are supported:
//
//   - HTTPRoute
//   - TCPRoute
//   - TLSRoute
//   - ReferenceGrant
func IndexerByObjType(obj client.Object) []string {
	var results []string

	switch resource := obj.(type) {
	case *gwv1.HTTPRoute:
		for _, pRef := range resource.Spec.ParentRefs {
			if !(isGateway(pRef) || isListenerSet(pRef)) {
				continue
			}
			ns := resolveNs(pRef.Namespace)
			if ns == "" {
				ns = resource.Namespace
			}
			nns := types.NamespacedName{
				Namespace: ns,
				Name:      string(pRef.Name),
			}
			results = append(results, nns.String())
		}
	case *gwv1a2.TCPRoute:
		for _, pRef := range resource.Spec.ParentRefs {
			if !(isGateway(pRef) || isListenerSet(pRef)) {
				continue
			}
			ns := resolveNs(pRef.Namespace)
			if ns == "" {
				ns = resource.Namespace
			}
			nns := types.NamespacedName{
				Namespace: ns,
				Name:      string(pRef.Name),
			}
			results = append(results, nns.String())
		}
	case *gwv1a2.TLSRoute:
		for _, pRef := range resource.Spec.ParentRefs {
			if !(isGateway(pRef) || isListenerSet(pRef)) {
				continue
			}
			ns := resolveNs(pRef.Namespace)
			if ns == "" {
				ns = resource.Namespace
			}
			nns := types.NamespacedName{
				Namespace: ns,
				Name:      string(pRef.Name),
			}
			results = append(results, nns.String())
		}
	case *gwxv1a1.XListenerSet:
		if resource.Spec.ParentRef.Group != nil && *resource.Spec.ParentRef.Group != gwv1a2.GroupName {
			break
		}
		if resource.Spec.ParentRef.Kind != nil && *resource.Spec.ParentRef.Kind != wellknown.GatewayKind {
			break
		}
		ns := resolveNs(resource.Spec.ParentRef.Namespace)
		if ns == "" {
			ns = resource.Namespace
		}
		nns := types.NamespacedName{
			Namespace: ns,
			Name:      string(resource.Spec.ParentRef.Name),
		}
		results = append(results, nns.String())

	case *gwv1b1.ReferenceGrant:
		for _, from := range resource.Spec.From {
			if from.Namespace != "" {
				results = append(results, string(from.Namespace))
			}
		}
	default:
		// Unsupported route type
		return results
	}

	return results
}

func IndexByHTTPRouteDelegationLabelSelector(obj client.Object) []string {
	route := obj.(*gwv1.HTTPRoute)
	value, ok := route.Labels[wellknown.RouteDelegationLabelSelector]
	if !ok {
		return nil
	}
	return []string{value}
}

// resolveNs resolves the namespace from an optional Namespace field.
func resolveNs(ns *gwv1.Namespace) string {
	if ns == nil {
		return ""
	}
	return string(*ns)
}
