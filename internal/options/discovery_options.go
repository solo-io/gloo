package options

type DiscoveryOptions struct {
	AutoDiscoverSwagger bool
	SwaggerUrisToTry    []string

	AutoDiscoverNATS bool
	ClusterIDsToTry  []string
}
