package validation

import (
	"context"
	"errors"

	v1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gateway/pkg/utils"
	gloov1snap "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/gloosnapshot"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
)

var _ GatewayResourceValidation = GatewayValidation{}

type GatewayValidation struct {
}

func (vsv GatewayValidation) DeleteResource(ctx context.Context, ref *core.ResourceRef, v Validator, dryRun bool) error {
	// Not Implemented
	return errors.New("cannot validate a deletion of a gateway resource")
}

func (rtv GatewayValidation) GetProxies(ctx context.Context, resource resources.HashableInputResource, snap *gloov1snap.ApiSnapshot) ([]string, error) {
	return utils.GetProxyNamesForGateway(resource.(*v1.Gateway)), nil
}
