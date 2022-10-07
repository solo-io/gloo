package validation

import (
	"context"

	v1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	gloov1snap "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/gloosnapshot"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// DeleteGatewayResourceValidator are resources that have a Group of gateway.solo.io
// this interface helps to delete when the resource is being validated
type DeleteGatewayResourceValidator interface {
	// DeleteResource will validate the process of deleting the resource
	DeleteResource(ctx context.Context, ref *core.ResourceRef, v Validator, dryRun bool) error
}

// GatewayResourceValidator are resources that have a Group of gateway.solo.io
// this interface helps to retrieve the proxies when the resource is being validated
type GatewayResourceValidator interface {
	// GetProxies will retrieve the proxies based off the resource type
	GetProxies(ctx context.Context, resource resources.HashableInputResource, snap *gloov1snap.ApiSnapshot) ([]string, error)
}

// GvkToGatewayValidator the current group of resources that can be validated, that implement the GatewayResoruceValidation interface
var GvkToGatewayResourceValidator = map[schema.GroupVersionKind]GatewayResourceValidator{
	v1.GatewayGVK:        &GatewayValidator{},
	v1.VirtualServiceGVK: &VirtualServiceValidation{},
	v1.RouteTableGVK:     &RouteTableValidator{},
}

// GvkToDeleteGatewayResourceValidator the current group of resources that can be validated, that implement the GvkToDeleteGatewayResourceValidator interface
var GvkToDeleteGatewayResourceValidator = map[schema.GroupVersionKind]DeleteGatewayResourceValidator{
	v1.VirtualServiceGVK: &VirtualServiceValidation{},
	v1.RouteTableGVK:     &RouteTableValidator{},
}
