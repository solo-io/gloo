package nginx

// An Nginx configuration
type IngressNginxConfig struct {
	// TODO(talnordan): support multiple servers.
	Server Server
}
