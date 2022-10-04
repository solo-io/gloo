package utils

import (
	"github.com/rotisserie/eris"
	"github.com/solo-io/gloo/projects/gateway/pkg/defaults"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"

	v1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
)

func GatewaysByProxyName(gateways v1.GatewayList) map[string]v1.GatewayList {
	result := make(map[string]v1.GatewayList)
	for _, gw := range gateways {
		proxyNames := GetProxyNamesForGateway(gw)
		for _, name := range proxyNames {
			result[name] = append(result[name], gw)
		}
	}
	return result
}

func GetProxyNamesForGateway(gw *v1.Gateway) []string {
	proxyNames := gw.GetProxyNames()
	if len(proxyNames) == 0 {
		proxyNames = []string{defaults.GatewayProxyName}
	}
	return proxyNames
}

func GetProxiesFromHashableInputResource(resource resources.HashableInputResource) ([]string, error) {
	switch typed := resource.(type) {
	case *v1.Gateway:
		return GetProxyNamesForGateway(typed), nil
	default:
		return nil, eris.New("the type for the resource does not exist when getting the proxies")
	}
}
