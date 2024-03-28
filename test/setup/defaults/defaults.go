package defaults

const (
	DefaultIstioTag           = "1.18.2"
	DefaultIstioImageRegistry = "docker.io/istio"
	IstioctlVersionEnv        = "ISTIOCTL_VERSION"

	DefaultGlooImageRegistry = "quay.io/solo-io"

	DefaultGlooImageName         = "gloo"
	DefaultDiscoveryImageName    = "discovery"
	DefaultIgressImageName       = "ingress"
	DefaultSdsImageName          = "sds"
	DefaultIstioProxyImageName   = "proxyv2"
	DefaultCertGenJobImageName   = "certgen"
	DefaultCleanupJobImageName   = "kubectl"
	DefaultRolloutJobImageName   = "kubectl"
	DefaultAccessLoggerImageName = "access-logger"
	DefaultGatewayImageName      = "gloo-envoy-wrapper"
)
