package config

import (
	"context"
	"time"

	"github.com/solo-io/go-utils/hashutils"

	"github.com/hashicorp/go-multierror"
	errors "github.com/rotisserie/eris"
	"github.com/solo-io/go-utils/contextutils"
	"go.uber.org/zap"

	extauthv1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/extauth/v1"
)

const (
	DefaultCallback      = "/oauth-gloo-callback"
	DefaultOAuthCacheTtl = time.Minute * 10
	// Default to 30 days (in seconds)
	defaultMaxAge = 30 * 24 * 60 * 60
)

var (
	MissingAuthConfigRefError = errors.New("missing required [authConfigRefName] field")
	GetAuthServiceError       = func(err error, id string, keepingPreviousConfig bool) error {
		additionalContext := "this configuration will be ignored"
		if keepingPreviousConfig {
			additionalContext = "server will continue using previous configuration for this id"
		}
		return errors.Wrapf(err, "failed to get auth service for auth config with id [%s]; %s", id, additionalContext)
	}
)

type Generator interface {
	GenerateConfig(resources []*extauthv1.ExtAuthConfig) (*serverState, error)
}

func NewGenerator(
	ctx context.Context,
	userIdHeader string,
	configTranslator ExtAuthConfigTranslator,
) *configGenerator {
	return &configGenerator{
		originalCtx:      ctx,
		userIdHeader:     userIdHeader,
		configTranslator: configTranslator,

		// Initial state will be an empty config
		currentState: newServerState(ctx, userIdHeader, nil),
	}
}

type configGenerator struct {
	originalCtx      context.Context
	userIdHeader     string
	configTranslator ExtAuthConfigTranslator

	cancel       context.CancelFunc
	currentState *serverState
}

func (c *configGenerator) GenerateConfig(resources []*extauthv1.ExtAuthConfig) (*serverState, error) {
	errs := &multierror.Error{}

	// Initialize new server state
	newState := newServerState(c.originalCtx, c.userIdHeader, resources)

	var authConfigsToStart []string
	for configId, newConfig := range newState.configs {

		currentConfig, currentlyExists := c.currentState.configs[configId]

		// If the config has not changed, just use the current one in the new state.
		// We do NOT want to cancel the context and restart the service in this case.
		if currentlyExists && currentConfig.hash == newConfig.hash {
			newState.configs[configId] = currentConfig
			continue
		}

		// Create context for new config
		newConfig.ctx, newConfig.cancel = context.WithCancel(c.originalCtx)

		// Create an AuthService from the new config
		authService, err := c.configTranslator.Translate(newConfig.ctx, newConfig.config)
		if err != nil {

			// Cancel context to be safe
			newConfig.cancel()

			if currentlyExists {
				// If the current state contains a valid config with this id (i.e. a previously valid AuthConfig),
				// then keep the current config running.
				errs = multierror.Append(errs, GetAuthServiceError(err, newConfig.config.AuthConfigRefName, true))
				newState.configs[configId] = currentConfig
			} else {
				// If this configuration is new, just drop it.
				errs = multierror.Append(errs, GetAuthServiceError(err, newConfig.config.AuthConfigRefName, false))
				delete(newState.configs, configId)
			}

			continue
		}

		newConfig.authService = authService

		authConfigsToStart = append(authConfigsToStart, configId)
	}

	// Log errors, if any
	if err := errs.ErrorOrNil(); err != nil {
		contextutils.LoggerFrom(c.originalCtx).
			Errorw("Errors encountered while processing new server configuration", zap.Error(err))
	}

	// Check for current configurations that are orphaned and cancel their context to avoid leaks
	for id, currentConfig := range c.currentState.configs {
		if _, exists := newState.configs[id]; !exists {
			currentConfig.cancel()
		}
	}

	// For each of the AuthServices that are either new or have changed:
	// - if an instance is already running, terminate it by cancelling its context
	// - call the Start function
	for _, id := range authConfigsToStart {

		if currentConfig, exists := c.currentState.configs[id]; exists {
			currentConfig.cancel()
		}

		newConfig := newState.configs[id]
		go func() {
			if err := newConfig.authService.Start(newConfig.ctx); err != nil {
				contextutils.LoggerFrom(c.originalCtx).Errorw("Error calling Start function",
					zap.Error(err), zap.String("authConfig", newConfig.config.AuthConfigRefName))
			}
		}()
	}

	// Store the new state so that it is available when the next config update is received.
	c.currentState = newState

	return newState, nil
}

func newServerState(ctx context.Context, userIdHeader string, resources []*extauthv1.ExtAuthConfig) *serverState {
	state := &serverState{
		userAuthHeader: userIdHeader,
		configs:        map[string]*configState{},
	}

	for _, resource := range resources {

		if resource.AuthConfigRefName == "" {
			// this should never happen
			contextutils.LoggerFrom(ctx).DPanicw("Invalid ExtAuthConfig resource will be ignored",
				zap.Error(MissingAuthConfigRefError), zap.Any("resource", resource))
			continue
		}

		state.configs[resource.AuthConfigRefName] = &configState{
			config: resource,
			hash:   hashutils.HashAll(resource),
		}
	}

	return state
}
