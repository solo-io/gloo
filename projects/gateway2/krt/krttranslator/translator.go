package krttranslator

import (
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/kube/apis/gloo.solo.io/v1"

	gwapi "sigs.k8s.io/gateway-api/apis/v1beta1"
)

type ProxyGateway struct {
	Proxy   *v1.Proxy
	Gateway *gwapi.Gateway
}



	// rm := reports.NewReportMap()
	// r := reports.NewReporter(&rm)
	// applyPostTranslationPlugins(ctx, pluginRegistry, &gwplugins.PostTranslationContext{
	// 	TranslatedGateways: translatedGateways,
	// })
