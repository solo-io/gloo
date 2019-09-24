package options

import (
	"sort"

	extauth "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/plugins/extauth/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/plugins/ratelimit"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
)

type ExtraOptions struct {
	Install    InstallExtended
	RateLimit  RateLimit
	OIDCAuth   OIDCAuth
	ApiKeyAuth ApiKeyAuth
	OpaAuth    OpaAuth
}

var RateLimit_TimeUnits = func() []string {
	var vals []string
	for _, name := range ratelimit.RateLimit_Unit_name {
		vals = append(vals, name)
	}
	sort.Strings(vals)
	return vals
}()

type RateLimit struct {
	Enable              bool
	TimeUnit            string
	RequestsPerTimeUnit uint32
}

type Dashboard struct {
}

type OIDCAuth struct {
	Enable bool

	// Include all options from the vhost extension
	extauth.OAuth
}

type ApiKeyAuth struct {
	Enable          bool
	Labels          []string
	SecretNamespace string
	SecretName      string
}

type OIDCSettings struct {
	ExtAtuhServerUpstreamRef core.ResourceRef
}

type InstallExtended struct {
	LicenseKey string
}

type OpaAuth struct {
	Enable bool

	Query   string
	Modules []string
}
