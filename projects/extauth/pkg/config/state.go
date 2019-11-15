package config

import (
	"context"

	"github.com/solo-io/ext-auth-plugins/api"
	extauth "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/plugins/extauth/v1"
)

// The state of the extauth server
type serverState struct {
	// The set of currently running AuthService instances.
	configs map[string]*configState
	// The name of the header used  by Gloo to identify authenticated users.
	userAuthHeader string
}

// Represents a running AuthService instance.
type configState struct {
	// The xDS resource that was used to generate the AuthService
	config *extauth.ExtAuthConfig
	// The AuthService resulting from the xDS resource
	authService api.AuthService
	// A hash of the xDS resource. We use this to determine if a particular resource has changed between updates.
	hash uint64
	// The context that determines the lifecycle of the above AuthService.
	ctx context.Context
	// Used to terminate previous instances of AuthServices.
	cancel context.CancelFunc
}

func (c *serverState) AuthService(configId string) api.AuthService {
	configState, ok := c.configs[configId]
	if !ok {
		// ext-auth-service can handle this case
		return nil
	}
	return configState.authService
}

func (c *serverState) UserHeader() string {
	return c.userAuthHeader
}

func (c *serverState) GetConfigCount() int {
	return len(c.configs)
}
