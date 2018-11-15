package constants

const (
	AppName = "vcs"

	GatewayRootDir        = "gateways"
	VirtualServiceRootDir = "virtual-services"
	ProxyRootDir          = "proxies"
	SchemaRootDir         = "schemas"
	ResolverMapRootDir    = "resolver-maps"
	UpstreamRootDir       = "upstreams"
	SettingsRootDir       = "settings"

	// Name of the environment variable that holds the token used to authenticate with the git remote
	AuthTokenEnvVariableName = "SOLO_GITHUB_TOKEN"
)
