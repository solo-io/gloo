package pkg

import (
	"context"
	"fmt"
	"strings"

	"github.com/solo-io/ext-auth-plugins/api"
	"github.com/solo-io/go-utils/contextutils"
	"github.com/solo-io/go-utils/errors"
	"go.uber.org/zap"
)

const MembershipAttribute = "memberOf"

var (
	MissingRequiredConfigError = func(missingAttrs []string) error {
		return errors.Errorf("plugin configuration is missing required attributes: %v", missingAttrs)
	}
	MalformedTemplate = func(placeholderNum int) error {
		return errors.Errorf("user DN template must contain exactly one [%%s] placeholder. Found: %d", placeholderNum)
	}
	UnexpectedConfigError = func(typ interface{}) error {
		return errors.New(fmt.Sprintf("unexpected config type %T", typ))
	}
)

type LdapFactory struct{}

type Config struct {
	// URL of the LDAP server
	ServerUrl string
	// Template to build user entry distinguished names (DN)
	// E.g. "uid=%s,ou=people,dc=solo,dc=io"
	UserDnTemplate string
	// User must be member of one of these groups for the request to be authenticated
	// E.g. []string{ "cn=developers,ou=groups,dc=solo,dc=io" }
	AllowedGroups []string
}

func (l LdapFactory) NewConfigInstance(ctx context.Context) (configInstance interface{}, err error) {
	return &Config{}, nil
}

func (l LdapFactory) GetAuthService(ctx context.Context, configInstance interface{}) (api.AuthService, error) {
	config, ok := configInstance.(*Config)
	if !ok {
		return nil, UnexpectedConfigError(configInstance)
	}

	logger(ctx).Debugw("Parsed config",
		zap.Any("serverUrl", config.ServerUrl),
		zap.Any("userDnTemplate", config.UserDnTemplate),
		zap.Any("allowedGroups", config.AllowedGroups),
	)

	if err := validateConfig(config); err != nil {
		return nil, err
	}

	return NewLdapAuthService(NewLdapClientBuilder(), config), nil
}

func validateConfig(config *Config) error {
	var missing []string
	if config.ServerUrl == "" {
		missing = append(missing, "ServerUrl")
	}
	if config.UserDnTemplate == "" {
		missing = append(missing, "UserDnTemplate")
	}
	if len(config.AllowedGroups) == 0 {
		missing = append(missing, "AllowedGroups")
	}

	if len(missing) > 0 {
		return MissingRequiredConfigError(missing)
	}

	placeholderOccurrences := strings.Count(config.UserDnTemplate, "%s")
	if placeholderOccurrences != 1 {
		return MalformedTemplate(placeholderOccurrences)
	}

	return nil
}

func logger(ctx context.Context) *zap.SugaredLogger {
	return contextutils.LoggerFrom(contextutils.WithLogger(ctx, "ldap_plugin"))
}
