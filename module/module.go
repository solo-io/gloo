package module

import "github.com/solo-io/mockway/config"

type Module interface {
	// Identifier
	// identifier represents the root key of the config tree
	// that we pass to the module
	Identifier() string
	// the module should parse the config and tell us the names
	// of any kubernetes secrets we need to watch
	SecretsToWatch(configBlob []byte) (SecretNames []string)
	// translate should take all the current values of secrets
	// as well as the user config
	// and return the corresponding envoy configuration resources
	Translate(secrets map[string]string, configBlob []byte) (config.EnvoyResources, error)
}
