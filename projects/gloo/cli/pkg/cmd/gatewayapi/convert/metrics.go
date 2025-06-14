package convert

import (
	"fmt"
	"os"
)

func (o *GatewayAPIOutput) PrintMetrics(filesEvaluated int) {

	_, _ = fmt.Fprintf(os.Stdout, "-------------------------------------")

	totalEdgeAPIs := len(o.edgeCache.Gateways()) + len(o.edgeCache.VirtualServices()) + len(o.edgeCache.RouteTables()) + len(o.edgeCache.Upstreams()) + len(o.edgeCache.VirtualHostOptions()) + len(o.edgeCache.ListenerOptions()) + len(o.edgeCache.HTTPListenerOptions()) + len(o.edgeCache.AuthConfigs()) + len(o.edgeCache.Settings())

	fmt.Printf("\nGloo Edge APIs:")
	fmt.Printf("\n\tGloo Gateways: %d", len(o.edgeCache.Gateways()))
	fmt.Printf("\n\tVirtualService: %d", len(o.edgeCache.VirtualServices()))
	fmt.Printf("\n\tRouteTables: %d", len(o.edgeCache.RouteTables()))
	fmt.Printf("\n\tUpstreams: %d", len(o.edgeCache.Upstreams()))
	fmt.Printf("\n\tVirtualHostOptions: %d", len(o.edgeCache.VirtualHostOptions()))
	fmt.Printf("\n\tListenerOptions: %d", len(o.edgeCache.ListenerOptions()))
	fmt.Printf("\n\tHTTPListenerOptions: %d", len(o.edgeCache.HTTPListenerOptions()))
	fmt.Printf("\n\tAuthConfigs: %d", len(o.edgeCache.AuthConfigs()))
	fmt.Printf("\n\tSettings: %d", len(o.edgeCache.Settings()))
	_, _ = fmt.Fprintf(os.Stdout, "\nTotal Edge API Config: %v", totalEdgeAPIs)

	_, _ = fmt.Fprintf(os.Stdout, "\n-------------------------------------")

	numGatewayObjects := len(o.gatewayAPICache.Gateways)
	numListenerSets := len(o.gatewayAPICache.ListenerSets)
	numHTTPRoutes := len(o.gatewayAPICache.HTTPRoutes)
	numDirectResponses := len(o.gatewayAPICache.DirectResponses)
	numHTTPListenerOptions := len(o.gatewayAPICache.HTTPListenerOptions)
	numAuthConfigs := len(o.gatewayAPICache.AuthConfigs)
	numSettings := len(o.gatewayAPICache.Settings)
	numRouteOptions := len(o.gatewayAPICache.RouteOptions)
	numUpstreams := len(o.gatewayAPICache.Upstreams)

	fmt.Printf("\nGateway API APIs:")
	fmt.Printf("\n\tGateways: %d", numGatewayObjects)
	fmt.Printf("\n\tListenerSets: %d", numListenerSets)
	fmt.Printf("\n\tHTTPRoutes: %d", numHTTPRoutes)
	fmt.Printf("\n\tDirectResponses: %d", numDirectResponses)
	fmt.Printf("\n\tHTTPListenerOptions: %d", numHTTPListenerOptions)
	fmt.Printf("\n\tAuthConfigs: %d", numAuthConfigs)
	fmt.Printf("\n\tSettings: %d", numSettings)
	fmt.Printf("\n\tRouteOptions: %d", numRouteOptions)
	fmt.Printf("\n\tUpstreams: %d", numUpstreams)

	numGatewayAPIObjects := numGatewayObjects + numListenerSets + numHTTPRoutes + numDirectResponses + numHTTPListenerOptions + numAuthConfigs + numSettings + numRouteOptions + numUpstreams

	_, _ = fmt.Fprintf(os.Stdout, "\nTotal Gateway API Config: %v", numGatewayAPIObjects)
	_, _ = fmt.Fprintf(os.Stdout, "\n-------------------------------------")
	totalErrors := 0
	for _, err := range o.errors {
		totalErrors += len(err)
	}
	_, _ = fmt.Fprintf(os.Stdout, "\nTotal Errors: %v", totalErrors)
	_, _ = fmt.Fprintf(os.Stdout, "\nFiles evaluated: %v", filesEvaluated)

}
