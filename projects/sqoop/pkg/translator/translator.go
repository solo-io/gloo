package translator

import (
	"strings"

	"log"

	"github.com/solo-io/solo-kit/pkg/api/v1/reporter"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	gloov1 "github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1"
	"github.com/solo-io/solo-kit/projects/sqoop/pkg/api/v1"
)

// trnslate a snapshot of schemas to:
// - 1 resolvermap (skeleton, only write for create) per schema
// - 1 proxy for the snapshot, assigned to the sqoop sidecar
func Translate(namespace string, snap *v1.ApiSnapshot) (*gloov1.Proxy, []*v1.ResolverMap, reporter.ResourceErrors) {
	resourceErrs := make(reporter.ResourceErrors)
	resourceErrs.Initialize(snap.Gateways.List().AsInputResources()...)
	resourceErrs.Initialize(snap.VirtualServices.List().AsInputResources()...)
	if len(snap.Gateways.List()) == 0 {
		log.Printf("%v had no gateways", snap.Hash())
		return nil, resourceErrs
	}
	if len(snap.VirtualServices.List()) == 0 {
		log.Printf("%v had no virtual services", snap.Hash())
		return nil, resourceErrs
	}
	validateGateways(snap.Gateways.List(), resourceErrs)
	validateVirtualServices(snap.VirtualServices.List(), resourceErrs)
	var listeners []*gloov1.Listener
	for _, gateway := range snap.Gateways.List() {
		listener := desiredListener(gateway, snap.VirtualServices.List(), resourceErrs)
		listeners = append(listeners, listener)
	}
	return &gloov1.Proxy{
		Metadata: core.Metadata{
			Name:        joinGatewayNames(snap.Gateways.List()) + "-proxy",
			Namespace:   namespace,
			Annotations: map[string]string{"owner_ref": "gateway"},
		},
		Listeners: listeners,
	}, resourceErrs
}

func joinGatewayNames(gateways v1.GatewayList) string {
	var names []string
	for _, gw := range gateways {
		names = append(names, gw.Metadata.Name)
	}
	return strings.Join(names, ".")
}

// TODO(ilackarms): implement validation func
func validateGateways(gateways v1.GatewayList, resourceErrs reporter.ResourceErrors) {

}

func validateVirtualServices(virtualServices v1.VirtualServiceList, resourceErrs reporter.ResourceErrors) {

}

func desiredListener(gateway *v1.Gateway, virtualServices v1.VirtualServiceList, resourceErrs reporter.ResourceErrors) *gloov1.Listener {
	// add all virtual services if empty
	if len(gateway.VirtualServices) == 0 {
		for _, virtualService := range virtualServices {
			gateway.VirtualServices = append(gateway.VirtualServices, core.ResourceRef{
				Name:      virtualService.GetMetadata().Name,
				Namespace: virtualService.GetMetadata().Namespace,
			})
		}
	}

	var (
		virtualHosts []*gloov1.VirtualHost
		sslConfigs   []*gloov1.SslConfig
	)

	for _, ref := range gateway.VirtualServices {
		// virtual service must live in the same namespace as gateway
		virtualService, err := virtualServices.Find(ref.Strings())
		if err != nil {
			resourceErrs.AddError(gateway, err)
			continue
		}
		virtualHosts = append(virtualHosts, virtualService.VirtualHost)
	}
	return &gloov1.Listener{
		Name:        gateway.Metadata.Name,
		BindAddress: gateway.BindAddress,
		BindPort:    gateway.BindPort,
		ListenerType: &gloov1.Listener_HttpListener{
			HttpListener: &gloov1.HttpListener{
				VirtualHosts: virtualHosts,
			},
		},
		SslConfiguations: sslConfigs,
	}
}
