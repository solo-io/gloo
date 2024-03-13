package wellknown

const (
	// GatewayClassName represents the name of the GatewayClass to watch for
	GatewayClassName         = "gloo-gateway"
	WaypointGatewayClassName = "gloo-waypoint"

	// GatewayControllerName is the name of the controller that has implemented the Gateway API
	// It is configured to manage GatewayClasses with the name GatewayClassName
	GatewayControllerName = "solo.io/gloo-gateway"

	// PROXY protocol in Envoy can be used to set the source/destination addresses
	// and provide arbitrary additional metadata via TLV headers.
	PROXYProtocol = "PROXY"
)
