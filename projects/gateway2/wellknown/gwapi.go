package wellknown

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
	apiv1 "sigs.k8s.io/gateway-api/apis/v1"
)

const (
	// Group string for Gateway API resources
	GatewayGroup = apiv1.GroupName

	// Kind string for k8s service
	ServiceKind = "Service"

	// Kind string for HTTPRoute resource
	HTTPRouteKind = "HTTPRoute"

	// Kind string for Gateway resource
	GatewayKind = "Gateway"

	// Kind string for GatewayClass resource
	GatewayClassKind = "GatewayClass"

	// Kind string for ReferenceGrant resource
	ReferenceGrantKind = "ReferenceGrant"

	// Kind strings for Gateway API list types
	HTTPRouteListKind      = "HTTPRouteList"
	GatewayListKind        = "GatewayList"
	GatewayClassListKind   = "GatewayClassList"
	ReferenceGrantListKind = "ReferenceGrantList"
)

var (
	// SupportedVersions specifies the supported versions of Gateway API.
	// [TODO] Remove "v1.0.0-rc1" when https://github.com/solo-io/gloo/issues/10115 is fixed.
	SupportedVersions = []string{"v1.1.0", "v1.0.0", "v1.0.0-rc1"}

	GatewayGVK = schema.GroupVersionKind{
		Group:   GatewayGroup,
		Version: "v1",
		Kind:    GatewayKind,
	}
	GatewayClassGVK = schema.GroupVersionKind{
		Group:   GatewayGroup,
		Version: "v1",
		Kind:    GatewayClassKind,
	}
	HTTPRouteGVK = schema.GroupVersionKind{
		Group:   GatewayGroup,
		Version: "v1",
		Kind:    HTTPRouteKind,
	}
	ReferenceGrantGVK = schema.GroupVersionKind{
		Group:   GatewayGroup,
		Version: "v1beta1",
		Kind:    ReferenceGrantKind,
	}

	GatewayListGVK = schema.GroupVersionKind{
		Group:   GatewayGroup,
		Version: "v1",
		Kind:    GatewayListKind,
	}
	GatewayClassListGVK = schema.GroupVersionKind{
		Group:   GatewayGroup,
		Version: "v1",
		Kind:    GatewayClassListKind,
	}
	HTTPRouteListGVK = schema.GroupVersionKind{
		Group:   GatewayGroup,
		Version: "v1",
		Kind:    HTTPRouteListKind,
	}
	ReferenceGrantListGVK = schema.GroupVersionKind{
		Group:   GatewayGroup,
		Version: "v1beta1",
		Kind:    ReferenceGrantListKind,
	}
)
