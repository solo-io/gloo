package kubeutils

const (
	GlooDeploymentName   = "gloo"
	GlooServiceName      = "gloo"
	GlooServiceAppLabel  = "gloo"
	GlooServiceGlooLabel = "gloo"

	// GlooXdsPortName is the name of the port in the Gloo Gateway control plane Kubernetes Service that serves xDS config.
	// See: install/helm/gloo/templates/2-gloo-service.yaml
	GlooXdsPortName = "grpc-xds"
)
