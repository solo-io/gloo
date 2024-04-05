package kubeutils

const (
	GlooDeploymentName = "gloo"
	GlooServiceName    = "gloo"

	GatewayProxyDeploymentName = "gateway-proxy"

	// The name of the port in the gloo control plane Kubernetes Service that serves xDS config.
	GlooXdsPortName = "grpc-xds"
)
