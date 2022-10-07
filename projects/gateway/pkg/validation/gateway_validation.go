package validation

import (
	"context"

	v1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gateway/pkg/utils"
	gloov1snap "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/gloosnapshot"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
)

var _ GatewayResourceValidator = GatewayValidator{}

type GatewayValidator struct {
}

func (rtv GatewayValidator) GetProxies(ctx context.Context, resource resources.HashableInputResource, snap *gloov1snap.ApiSnapshot) ([]string, error) {
	return utils.GetProxyNamesForGateway(resource.(*v1.Gateway)), nil
}
