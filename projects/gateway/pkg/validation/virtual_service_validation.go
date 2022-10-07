package validation

import (
	"context"
	"sort"

	v1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gateway/pkg/translator"
	"github.com/solo-io/gloo/projects/gateway/pkg/utils"
	gloov1snap "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/gloosnapshot"
	"github.com/solo-io/go-utils/contextutils"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
)

var _ GatewayResourceValidator = VirtualServiceValidation{}
var _ DeleteGatewayResourceValidator = VirtualServiceValidation{}

type VirtualServiceValidation struct {
}

func (vsv VirtualServiceValidation) DeleteResource(ctx context.Context, ref *core.ResourceRef, v Validator, dryRun bool) error {
	return v.ValidateDeleteVirtualService(ctx, ref, dryRun)
}

func (rtv VirtualServiceValidation) GetProxies(ctx context.Context, resource resources.HashableInputResource, snap *gloov1snap.ApiSnapshot) ([]string, error) {
	return proxiesForVirtualService(ctx, snap.Gateways, snap.HttpGateways, resource.(*v1.VirtualService)), nil
}

func proxiesForVirtualService(ctx context.Context, gwList v1.GatewayList, httpGwList v1.MatchableHttpGatewayList, vs *v1.VirtualService) []string {
	gatewaysByProxy := utils.GatewaysByProxyName(gwList)

	var proxiesToConsider []string

	for proxyName, gatewayList := range gatewaysByProxy {
		if gatewayListContainsVirtualService(ctx, gatewayList, httpGwList, vs) {
			// we only care about validating this proxy if it contains this virtual service
			proxiesToConsider = append(proxiesToConsider, proxyName)
		}
	}

	sort.Strings(proxiesToConsider)

	return proxiesToConsider
}

func gatewayListContainsVirtualService(ctx context.Context, gwList v1.GatewayList, httpGwList v1.MatchableHttpGatewayList, vs *v1.VirtualService) bool {
	for _, gw := range gwList {
		if gatewayContainsVirtualService(ctx, httpGwList, gw, vs) {
			return true
		}
	}

	return false
}

func gatewayContainsVirtualService(ctx context.Context, httpGwList v1.MatchableHttpGatewayList, gw *v1.Gateway, vs *v1.VirtualService) bool {
	if gw.GetTcpGateway() != nil {
		return false
	}

	if httpGateway := gw.GetHttpGateway(); httpGateway != nil {
		return httpGatewayContainsVirtualService(httpGateway, vs, gw.GetSsl())
	}

	if hybridGateway := gw.GetHybridGateway(); hybridGateway != nil {
		matchedGateways := hybridGateway.GetMatchedGateways()
		if matchedGateways != nil {
			for _, mg := range hybridGateway.GetMatchedGateways() {
				if httpGateway := mg.GetHttpGateway(); httpGateway != nil {
					if httpGatewayContainsVirtualService(httpGateway, vs, mg.GetMatcher().GetSslConfig() != nil) {
						return true
					}
				}
			}
		} else {
			delegatedGateway := hybridGateway.GetDelegatedHttpGateways()
			selectedGatewayList := translator.NewHttpGatewaySelector(httpGwList).SelectMatchableHttpGateways(delegatedGateway, func(err error) {
				logger := contextutils.LoggerFrom(ctx)
				logger.Warnf("failed to select matchable http gateways on gateway: %v", err.Error())
			})
			for _, selectedHttpGw := range selectedGatewayList {
				if httpGatewayContainsVirtualService(selectedHttpGw.GetHttpGateway(), vs, selectedHttpGw.GetMatcher().GetSslConfig() != nil) {
					return true
				}
			}
		}

	}

	return false
}

func httpGatewayContainsVirtualService(httpGateway *v1.HttpGateway, vs *v1.VirtualService, requireSsl bool) bool {
	contains, err := translator.HttpGatewayContainsVirtualService(httpGateway, vs, requireSsl)
	if err != nil {
		return false
	}
	return contains
}
