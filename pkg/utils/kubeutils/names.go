package kubeutils

const (
	GlooDeploymentName = "gloo"
	GlooServiceName    = "gloo"

	// The name of the port in the gloo control plane Kubernetes Service that serves xDS config.
	GlooXdsPortName = "grpc-xds"

	// Image defaults
	DefaultGlooImageRegistry     = "quay.io/solo-io"
	DefaultGlooImageName         = "gloo"
	DefaultDiscoveryImageName    = "discovery"
	DefaultIgressImageName       = "ingress"
	DefaultSdsImageName          = "sds"
	DefaultCertGenJobImageName   = "certgen"
	DefaultCleanupJobImageName   = "kubectl"
	DefaultRolloutJobImageName   = "kubectl"
	DefaultAccessLoggerImageName = "access-logger"
	DefaultGatewayImageName      = "gloo-envoy-wrapper"

	// Istio defaults
	DefaultIstioTag            = "1.18.2"
	DefaultIstioImageRegistry  = "docker.io/istio"
	DefaultIstioProxyImageName = "proxyv2"
)
