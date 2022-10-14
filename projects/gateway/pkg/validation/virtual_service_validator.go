package validation

import (
	"context"

	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
)

var _ DeleteGatewayResourceValidator = VirtualServiceValidation{}

type VirtualServiceValidation struct {
}

func (vsv VirtualServiceValidation) DeleteResource(ctx context.Context, ref *core.ResourceRef, v Validator, dryRun bool) error {
	return v.ValidateDeleteVirtualService(ctx, ref, dryRun)
}
