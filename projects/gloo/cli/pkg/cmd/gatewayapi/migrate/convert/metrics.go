package convert

import (
	"fmt"
	"os"
)

func (g *GatewayAPIOutput) PrintMetrics(filesEvaluated int) {

	_, _ = fmt.Fprintf(os.Stdout, "-------------------------------------")

	totalEdgeAPIs := len(g.edgeCache.Gateways()) + len(g.edgeCache.VirtualServices()) + len(g.edgeCache.RouteTables()) + len(g.edgeCache.Upstreams()) + len(g.edgeCache.VirtualHostOptions()) + len(g.edgeCache.ListenerOptions()) + len(g.edgeCache.HTTPListenerOptions()) + len(g.edgeCache.AuthConfigs()) + len(g.edgeCache.Settings())

	fmt.Printf("\nGloo Edge APIs:")
	fmt.Printf("\n\tGloo Gateways: %d", len(g.edgeCache.Gateways()))
	fmt.Printf("\n\tVirtualService: %d", len(g.edgeCache.VirtualServices()))
	fmt.Printf("\n\tRouteTables: %d", len(g.edgeCache.RouteTables()))
	fmt.Printf("\n\tUpstreams: %d", len(g.edgeCache.Upstreams()))
	fmt.Printf("\n\tVirtualHostOptions: %d", len(g.edgeCache.VirtualHostOptions()))
	fmt.Printf("\n\tListenerOptions: %d", len(g.edgeCache.ListenerOptions()))
	fmt.Printf("\n\tHTTPListenerOptions: %d", len(g.edgeCache.HTTPListenerOptions()))
	fmt.Printf("\n\tAuthConfigs: %d", len(g.edgeCache.AuthConfigs()))
	fmt.Printf("\n\tSettings: %d", len(g.edgeCache.Settings()))
	_, _ = fmt.Fprintf(os.Stdout, "\nTotal Edge API Config: %v", totalEdgeAPIs)

	_, _ = fmt.Fprintf(os.Stdout, "\n-------------------------------------")

	numGatewayObjects := len(g.gatewayAPICache.Gateways)
	numListenerSets := len(g.gatewayAPICache.ListenerSets)
	numHTTPRoutes := len(g.gatewayAPICache.HTTPRoutes)
	numDirectResponses := len(g.gatewayAPICache.DirectResponses)
	numGlooTrafficPolicies := len(g.gatewayAPICache.GlooTrafficPolicies)
	numAuthConfigs := len(g.gatewayAPICache.AuthConfigs)
	numBackends := len(g.gatewayAPICache.Backends)
	numBackendConfigPolicies := len(g.gatewayAPICache.BackendConfigPolicy)
	numGatewayParameters := len(g.gatewayAPICache.KGatewayParameters)
	numHTTPListenerPolicies := len(g.gatewayAPICache.HTTPListenerPolicies)
	numGatewayExtensions := len(g.gatewayAPICache.GatewayExtensions)

	fmt.Printf("\nGateway API APIs:")
	fmt.Printf("\n\tGateways: %d", numGatewayObjects)
	fmt.Printf("\n\tGatewayParameters: %d", numGatewayParameters)
	fmt.Printf("\n\tGatewayExtensions: %d", numGatewayExtensions)
	fmt.Printf("\n\tListenerSets: %d", numListenerSets)
	fmt.Printf("\n\tHTTPRoutes: %d", numHTTPRoutes)
	fmt.Printf("\n\tDirectResponses: %d", numDirectResponses)
	fmt.Printf("\n\tHTTPListenerPolicies: %d", numHTTPListenerPolicies)
	fmt.Printf("\n\tGlooTrafficPolicies: %d", numGlooTrafficPolicies)
	fmt.Printf("\n\tAuthConfigs: %d", numAuthConfigs)
	fmt.Printf("\n\tBackends: %d", numBackends)
	fmt.Printf("\n\tBackendConfigPolicies: %d", numBackendConfigPolicies)

	numGatewayAPIObjects := numGatewayObjects + numListenerSets + numHTTPRoutes + numDirectResponses + numGlooTrafficPolicies + numAuthConfigs + numBackends + numBackendConfigPolicies + numGatewayParameters + numHTTPListenerPolicies + numGatewayExtensions

	_, _ = fmt.Fprintf(os.Stdout, "\nTotal Gateway API Config: %v", numGatewayAPIObjects)
	_, _ = fmt.Fprintf(os.Stdout, "\n-------------------------------------")
	totalErrors := 0
	for _, err := range g.errors {
		totalErrors += len(err)
	}
	_, _ = fmt.Fprintf(os.Stdout, "\nTotal Errors: %v", totalErrors)
	_, _ = fmt.Fprintf(os.Stdout, "\nFiles evaluated: %v", filesEvaluated)

}
