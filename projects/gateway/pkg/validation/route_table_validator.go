package validation

import (
	"context"

	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
)

var _ DeleteGatewayResourceValidator = RouteTableValidator{}

type RouteTableValidator struct {
}

func (vsv RouteTableValidator) DeleteResource(ctx context.Context, ref *core.ResourceRef, v Validator, dryRun bool) error {
	return v.ValidateDeleteRouteTable(ctx, ref, dryRun)
}
