package options

type DiscoveryOptions struct {
	AutoDiscoverSwagger bool
	SwaggerUrisToTry    []string

	AutoDiscoverNATS    bool
	AutoDiscoverFaaS    bool
	AutoDiscoverFission bool
	ClusterIDsToTry     []string

	AutoDiscoverGRPC bool
}
