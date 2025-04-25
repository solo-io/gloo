package convert

import (
	"fmt"
	"os"
)

func printMetrics(output *GatewayAPIOutput, filesEvaluated int) {

	_, _ = fmt.Fprintf(os.Stdout, "-------------------------------------")

	totalEdgeAPIs := len(output.edgeCache.Gateways()) + len(output.edgeCache.VirtualServices()) + len(output.edgeCache.RouteTables()) + len(output.edgeCache.Upstreams()) + len(output.edgeCache.VirtualHostOptions()) + len(output.edgeCache.ListenerOptions()) + len(output.edgeCache.HTTPListenerOptions()) + len(output.edgeCache.AuthConfigs()) + len(output.edgeCache.Settings())

	fmt.Printf("\nGloo Edge APIs:")
	fmt.Printf("\n\tGloo Gateways: %d", len(output.edgeCache.Gateways()))
	fmt.Printf("\n\tVirtualService: %d", len(output.edgeCache.VirtualServices()))
	fmt.Printf("\n\tRouteTables: %d", len(output.edgeCache.RouteTables()))
	fmt.Printf("\n\tUpstreams: %d", len(output.edgeCache.Upstreams()))
	fmt.Printf("\n\tVirtualHostOptions: %d", len(output.edgeCache.VirtualHostOptions()))
	fmt.Printf("\n\tListenerOptions: %d", len(output.edgeCache.ListenerOptions()))
	fmt.Printf("\n\tHTTPListenerOptions: %d", len(output.edgeCache.HTTPListenerOptions()))
	fmt.Printf("\n\tAuthConfigs: %d", len(output.edgeCache.AuthConfigs()))
	fmt.Printf("\n\tSettings: %d", len(output.edgeCache.Settings()))
	_, _ = fmt.Fprintf(os.Stdout, "\nTotal Edge API Config: %v", totalEdgeAPIs)

	_, _ = fmt.Fprintf(os.Stdout, "\n-------------------------------------")

	numGatewayObjects := len(output.gatewayAPICache.Gateways)
	numListenerSets := len(output.gatewayAPICache.ListenerSets)
	numHTTPRoutes := len(output.gatewayAPICache.HTTPRoutes)
	numDirectResponses := len(output.gatewayAPICache.DirectResponses)
	numHTTPListenerOptions := len(output.gatewayAPICache.HTTPListenerOptions)
	numAuthConfigs := len(output.gatewayAPICache.AuthConfigs)
	numSettings := len(output.gatewayAPICache.Settings)
	numRouteOptions := len(output.gatewayAPICache.RouteOptions)
	numUpstreams := len(output.gatewayAPICache.Upstreams)

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
	for _, err := range output.errors {
		totalErrors += len(err)
	}
	_, _ = fmt.Fprintf(os.Stdout, "\nTotal Errors: %v", totalErrors)
	_, _ = fmt.Fprintf(os.Stdout, "\nFiles evaluated: %v", filesEvaluated)

}
