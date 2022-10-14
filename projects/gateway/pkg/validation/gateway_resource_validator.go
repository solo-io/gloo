package validation

import (
	v1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const GatewayGroup = "gateway.solo.io"

// GvkSupportedValidationGatewayResources the current group of resources that can be validated
var GvkSupportedValidationGatewayResources = map[schema.GroupVersionKind]bool{
	v1.GatewayGVK:        true,
	v1.VirtualServiceGVK: true,
	v1.RouteTableGVK:     true,
}

// GvkSupportedDeleteGatewayResources the current group of resources that can be validated
var GvkSupportedDeleteGatewayResources = map[schema.GroupVersionKind]bool{
	v1.VirtualServiceGVK: true,
	v1.RouteTableGVK:     true,
}
