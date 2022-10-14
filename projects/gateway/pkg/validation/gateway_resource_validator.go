package validation

import (
	"context"

	v1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const GatewayGroup = "gateway.solo.io"

/*
	Resources that are implemented in Validation must implement the DeleteGatewayResourceValidator and/or GatewayResourceValidator interfaces.
	These are used to handle the specific logic used in validation of gateway resources.
	Have to add  to GVK maps once supported GvkToGatewayResourceValidator and GvkToDeleteGatewayResourceValidator.
*/

// DeleteGatewayResourceValidator are resources that have a Group of gateway.solo.io
// this interface helps to delete when the resource is being validated
type DeleteGatewayResourceValidator interface {
	// DeleteResource will validate the process of deleting the resource
	DeleteResource(ctx context.Context, ref *core.ResourceRef, v Validator, dryRun bool) error
}

// GvkSupportedValidationGatewayResources the current group of resources that can be validated
var GvkSupportedValidationGatewayResources = map[schema.GroupVersionKind]bool{
	v1.GatewayGVK:        true,
	v1.VirtualServiceGVK: true,
	v1.RouteTableGVK:     true,
}

// GvkSupportedDeleteGatewayResources the current group of resources that can be validated
var GvkSupportedDeleteGatewayResources = map[schema.GroupVersionKind]DeleteGatewayResourceValidator{
	v1.VirtualServiceGVK: &VirtualServiceValidation{},
	v1.RouteTableGVK:     &RouteTableValidator{},
}
