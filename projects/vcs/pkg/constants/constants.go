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

	// This is the name of the git branch that is considered to be the master
	MasterBranchName = "master"

	// Branches must match this regular expression
	BranchRegExp = "^[a-zA-Z0-9_-]+$"

	// Name of the environment variable that holds the token used to authenticate with the git remote
	AuthTokenEnvVariableName = "SOLO_GITHUB_TOKEN"

	// Name of the environment variable that holds the git remote URI
	RemoteUriEnvVariableName = "SOLO_REMOTE_URI"
)
