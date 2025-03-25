package envoy

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "sigs.k8s.io/gateway-api/apis/v1"
)

// ListenerNamespaces indicate which namespaces ListenerSets should be selected from.
type ListenerNamespaces struct {
	// From indicates where ListenerSets can attach to this Gateway. Possible
	// values are:
	//
	// * Same: Only ListenerSets in the same namespace may be attached to this Gateway.
	// * Selector: ListenerSets in namespaces selected by the selector may be attached to this Gateway.:w
	// * All: ListenerSets in all namespaces may be attached to this Gateway.
	// * None: Only listeners defined in the Gateway's spec are allowed
	//
	// +optional
	// +kubebuilder:default=None
	// +kubebuilder:validation:Enum=Same;None;Selector;All
	From *v1.FromNamespaces `json:"from,omitempty"`

	// Selector must be specified when From is set to "Selector". In that case,
	// only ListenerSets in Namespaces matching this Selector will be selected by this
	// Gateway. This field is ignored for other values of "From".
	//
	// +optional
	Selector *metav1.LabelSelector `json:"selector,omitempty"`
}

// ListenerSet defines a set of additional listeners to attach to an existing Gateway.
type ListenerSet struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Spec defines the desired state of ListenerSet.
	Spec ListenerSetSpec `json:"spec"`

	// Status defines the current state of ListenerSet.
	Status ListenerSetStatus `json:"status,omitempty"`
}

// ListenerSetSpec defines the desired state of a ListenerSet.
type ListenerSetSpec struct {
	// ParentRef references the Gateway that the listeners are attached to.
	ParentRef ParentGatewayReference `json:"parentRef"`

	// Listeners associated with this ListenerSet. Listeners define
	// logical endpoints that are bound on this referenced parent Gateway's addresses.
	//
	// Listeners in a `Gateway` and their attached `ListenerSets` are concatenated
	// as a list when programming the underlying infrastructure.
	//
	// Listeners should be merged using the following precedence:
	//
	// 1. "parent" Gateway
	// 2. ListenerSet ordered by creation time (oldest first)
	// 3. ListenerSet ordered alphabetically by “{namespace}/{name}”.
	//
	// +listType=map
	// +listMapKey=name
	// +kubebuilder:validation:MinItems=1
	// +kubebuilder:validation:MaxItems=64
	Listeners []ListenerEntry `json:"listeners"`
}

// ListenerEntry embodies the concept of a logical endpoint where a Gateway accepts
// network connections.
type ListenerEntry struct {
	// Name is the name of the Listener. This name MUST be unique within a
	// ListenerSet.
	//
	// Name is not required to be unique across a Gateway and ListenerSets.
	// Routes can attach to a Listener by having a ListenerSet as a parentRef
	// and setting the SectionName
	Name v1.SectionName `json:"name"`

	// Hostname specifies the virtual hostname to match for protocol types that
	// define this concept. When unspecified, all hostnames are matched. This
	// field is ignored for protocols that don't require hostname based
	// matching.
	//
	// Implementations MUST apply Hostname matching appropriately for each of
	// the following protocols:
	//
	// * TLS: The Listener Hostname MUST match the SNI.
	// * HTTP: The Listener Hostname MUST match the Host header of the request.
	// * HTTPS: The Listener Hostname SHOULD match at both the TLS and HTTP
	//   protocol layers as described above. If an implementation does not
	//   ensure that both the SNI and Host header match the Listener hostname,
	//   it MUST clearly document that.
	//
	// For HTTPRoute and TLSRoute resources, there is an interaction with the
	// `spec.hostnames` array. When both listener and route specify hostnames,
	// there MUST be an intersection between the values for a Route to be
	// accepted. For more information, refer to the Route specific Hostnames
	// documentation.
	//
	// Hostnames that are prefixed with a wildcard label (`*.`) are interpreted
	// as a suffix match. That means that a match for `*.example.com` would match
	// both `test.example.com`, and `foo.test.example.com`, but not `example.com`.
	//
	// +optional
	Hostname *v1.Hostname `json:"hostname,omitempty"`

	// Port is the network port. Multiple listeners may use the
	// same port, subject to the Listener compatibility rules.
	//
	// If the port is specified as zero, the implementation will assign
	// a unique port. If the implementation does not support dynamic port
	// assignment, it MUST set `Accepted` condition to `False` with the
	// `UnsupportedPort` reason.
	//
	// Support: Core
	//
	// +optional
	Port *v1.PortNumber `json:"port,omitempty"`

	// Protocol specifies the network protocol this listener expects to receive.
	//
	// Support: Core
	Protocol v1.ProtocolType `json:"protocol"`

	// TLS is the TLS configuration for the Listener. This field is required if
	// the Protocol field is "HTTPS" or "TLS". It is invalid to set this field
	// if the Protocol field is "HTTP", "TCP", or "UDP".
	//
	// The association of SNIs to Certificate defined in GatewayTLSConfig is
	// defined based on the Hostname field for this listener.
	//
	// The GatewayClass MUST use the longest matching SNI out of all
	// available certificates for any TLS handshake.
	//
	// Support: Core
	//
	// +optional
	TLS *v1.GatewayTLSConfig `json:"tls,omitempty"`

	// AllowedRoutes defines the types of routes that MAY be attached to a
	// Listener and the trusted namespaces where those Route resources MAY be
	// present.
	//
	// Although a client request may match multiple route rules, only one rule
	// may ultimately receive the request. Matching precedence MUST be
	// determined in order of the following criteria:
	//
	// * The most specific match as defined by the Route type.
	// * The oldest Route based on creation timestamp. For example, a Route with
	//   a creation timestamp of "2020-09-08 01:02:03" is given precedence over
	//   a Route with a creation timestamp of "2020-09-08 01:02:04".
	// * If everything else is equivalent, the Route appearing first in
	//   alphabetical order (namespace/name) should be given precedence. For
	//   example, foo/bar is given precedence over foo/baz.
	//
	// All valid rules within a Route attached to this Listener should be
	// implemented. Invalid Route rules can be ignored (sometimes that will mean
	// the full Route). If a Route rule transitions from valid to invalid,
	// support for that Route rule should be dropped to ensure consistency. For
	// example, even if a filter specified by a Route rule is invalid, the rest
	// of the rules within that Route should still be supported.
	//
	// Support: Core
	// +kubebuilder:default={namespaces:{from: Same}}
	// +optional
	AllowedRoutes *v1.AllowedRoutes `json:"allowedRoutes,omitempty"`
}

// ListenerSetStatus defines the observed state of a ListenerSet
type ListenerSetStatus struct {
	// Listeners provide status for each unique listener port defined in the Spec.
	//
	// +optional
	// +listType=map
	// +listMapKey=name
	// +kubebuilder:validation:MaxItems=64
	Listeners []ListenerEntryStatus `json:"listeners,omitempty"`

	// Conditions describe the current conditions of the ListenerSet.
	//
	// Implementations should prefer to express ListenerSet conditions
	// using the `GatewayConditionType` and `GatewayConditionReason`
	// constants so that operators and tools can converge on a common
	// vocabulary to describe Gateway state.
	//
	// Known condition types are:
	//
	// * "Accepted"
	// * "Programmed"
	//
	// +optional
	// +listType=map
	// +listMapKey=type
	// +kubebuilder:validation:MaxItems=8
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// ListenerEntryStatus is the status associated with a ListenerEntry.
type ListenerEntryStatus struct {
	// Name is the name of the Listener that this status corresponds to.
	Name v1.SectionName `json:"name"`

	// Port is the network port that this listener is listening on.
	Port v1.PortNumber `json:"port"`

	// SupportedKinds is the list indicating the Kinds supported by this
	// listener. This MUST represent the kinds an implementation supports for
	// that Listener configuration.
	//
	// If kinds are specified in Spec that are not supported, they MUST NOT
	// appear in this list and an implementation MUST set the "ResolvedRefs"
	// condition to "False" with the "InvalidRouteKinds" reason. If both valid
	// and invalid Route kinds are specified, the implementation MUST
	// reference the valid Route kinds that have been specified.
	//
	// +kubebuilder:validation:MaxItems=8
	SupportedKinds []v1.RouteGroupKind `json:"supportedKinds"`

	// AttachedRoutes represents the total number of Routes that have been
	// successfully attached to this Listener.
	//
	// Successful attachment of a Route to a Listener is based solely on the
	// combination of the AllowedRoutes field on the corresponding Listener
	// and the Route's ParentRefs field. A Route is successfully attached to
	// a Listener when it is selected by the Listener's AllowedRoutes field
	// AND the Route has a valid ParentRef selecting the whole Gateway
	// resource or a specific Listener as a parent resource (more detail on
	// attachment semantics can be found in the documentation on the various
	// Route kinds ParentRefs fields). Listener or Route status does not impact
	// successful attachment, i.e. the AttachedRoutes field count MUST be set
	// for Listeners with condition Accepted: false and MUST count successfully
	// attached Routes that may themselves have Accepted: false conditions.
	//
	// Uses for this field include troubleshooting Route attachment and
	// measuring blast radius/impact of changes to a Listener.
	AttachedRoutes int32 `json:"attachedRoutes"`

	// Conditions describe the current condition of this listener.
	//
	// +listType=map
	// +listMapKey=type
	// +kubebuilder:validation:MaxItems=8
	Conditions []metav1.Condition `json:"conditions"`
}

// ParentGatewayReference identifies an API object including its namespace,
// defaulting to Gateway.
type ParentGatewayReference struct {
	// Group is the group of the referent.
	//
	// +optional
	// +kubebuilder:default="gateway.networking.k8s.io"
	Group *v1.Group `json:"group"`

	// Kind is kind of the referent. For example "Gateway".
	//
	// +optional
	// +kubebuilder:default=Gateway
	Kind *v1.Kind `json:"kind"`

	// Name is the name of the referent.
	Name v1.ObjectName `json:"name"`

	// Namespace is the name of the referent.
	// +optional
	Namespace *v1.ObjectName `json:"namespace"`
}
