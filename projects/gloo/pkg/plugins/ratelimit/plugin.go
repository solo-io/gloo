package ratelimit

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/api/v1/plugins/ratelimit"

	envoyroute "github.com/envoyproxy/go-control-plane/envoy/api/v2/route"
	envoyhttp "github.com/envoyproxy/go-control-plane/envoy/config/filter/network/http_connection_manager/v2"
	"github.com/solo-io/solo-kit/pkg/utils/protoutils"

	"github.com/gogo/protobuf/types"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins/utils"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/extauth"
)

/*
Background:

Currently the way to created descriptors in envoy is somewhate limited.
Even though we can use the server configuration to express many forms of rate limtits, we
are limited to configurations that we can also express in envoy.

I modeled the desired user configuration of rate limits for authenticated users and anonymous requests in envoy,
as such:

actions:
- header_value_match: {"descriptor_value":"is-authenticated", "expect_match":true, "headers":[{"name":"Authorization", "present_match":true}]}
- request_headers: {"header_name":"Authorization", "descriptor_key":"userid"}
actions:
- header_value_match: {"descriptor_value":"not-authenticated", "expect_match":false, "headers":[{"name":"Authorization", "present_match":true}]}
- remote_address: {}

Two actions, where the first one is the negation of the other. Since a failed entry causes the
whole action to not be generated, only one action (descriptor?) will be sent to the server.

The first action check checks if the Authorization is present. if it is we assume we can trust it, as the request
should have pass an auth filter first.
If the header is present, the second one get the header so we can rate limit on a per user basis.

If not (the second actions\descriptor), then the remote address is retrieved so we can limit per IP.

Given this envoy configuraiton, the appropriate server configuration would be:

constraints:
- key: generic_key
  value: <vhost_name>
  constraints:
  - key: header_match
    value: not-authenticated
    constraints:
    - key: remote_address
      rate_limit:
        unit: MINUTE
        requests_per_unit: 3
  - key: header_match
    value: is-authenticated
    constraints:
    - key: userid
      rate_limit:
        unit: MINUTE
        requests_per_unit: 10
*/

const (
	ExtensionName      = "rate-limit"
	EnvoyExtensionName = "envoy-rate-limit"
	IngressDomain      = "ingress"
	CustomDomain       = "custom"
	requestType        = "external"
	userid             = "userid"

	authenticated = "is-authenticated"
	anonymous     = "not-authenticated"

	stage       = 0
	customStage = 1
	timeout     = 20 * time.Millisecond

	headerMatch   = "header_match"
	genericKey    = "generic_key"
	remoteAddress = "remote_address"
)

const (
	filterName = "envoy.rate_limit"
	// rate limiting should happen after auth
	filterStage = plugins.PostInAuth
)

type Plugin struct {
	upstreamRef      *core.ResourceRef
	authUserIdHeader string
}

func NewPlugin() plugins.Plugin {
	return &Plugin{}
}

//// TODO(yuval-k): Copied from ext auth. the real solution is to add it to upstream Gloo
type tmpPluginContainer struct {
	params plugins.InitParams
}

func (t *tmpPluginContainer) GetExtensions() *v1.Extensions {
	return t.params.ExtensionsSettings
}

func (p *Plugin) Init(params plugins.InitParams) error {

	var settings ratelimit.Settings
	p.upstreamRef = nil
	err := utils.UnmarshalExtension(&tmpPluginContainer{params}, ExtensionName, &settings)
	if err != nil {
		if err == utils.NotFoundError {
			return nil
		}
		return err
	}

	p.upstreamRef = settings.RatelimitServerRef

	authSettings, _ := extauth.GetSettings(params)
	p.authUserIdHeader = extauth.GetAuthHeader(authSettings)

	return nil
}

func (p *Plugin) ProcessVirtualHost(params plugins.Params, in *v1.VirtualHost, out *envoyroute.VirtualHost) error {

	err := p.ProcessVirtualHostSimple(params, in, out)
	if err != nil {
		return err
	}
	return p.ProcessVirtualHostCustom(params, in, out)
}

func (p *Plugin) ProcessVirtualHostSimple(params plugins.Params, in *v1.VirtualHost, out *envoyroute.VirtualHost) error {
	var rateLimit ratelimit.IngressRateLimit
	err := utils.UnmarshalExtension(in.VirtualHostPlugins, ExtensionName, &rateLimit)
	if err != nil {
		if err == utils.NotFoundError {
			return nil
		}
		return errors.Wrapf(err, "Error converting proto to ingress rate limit plugin")
	}
	_, err = TranslateUserConfigToRateLimitServerConfig(in.Name, rateLimit)
	if err != nil {
		return err
	}

	vhost := generateEnvoyConfigForVhost(in.Name, p.authUserIdHeader)
	out.RateLimits = vhost

	return nil
}
func (p *Plugin) ProcessVirtualHostCustom(params plugins.Params, in *v1.VirtualHost, out *envoyroute.VirtualHost) error {
	var rateLimit ratelimit.RateLimitVhostExtension
	err := utils.UnmarshalExtension(in.VirtualHostPlugins, EnvoyExtensionName, &rateLimit)
	if err != nil {
		if err == utils.NotFoundError {
			return nil
		}
		return errors.Wrapf(err, "Error converting proto to vhost rate limit plugin")
	}

	out.RateLimits = generateCustomEnvoyConfigForVhost(rateLimit.RateLimits)

	return nil
}
func (p *Plugin) ProcessRoute(params plugins.Params, in *v1.Route, out *envoyroute.Route) error {
	var rateLimit ratelimit.RateLimitRouteExtension
	err := utils.UnmarshalExtension(in.RoutePlugins, EnvoyExtensionName, &rateLimit)
	if err != nil {
		if err == utils.NotFoundError {
			return nil
		}
		return errors.Wrapf(err, "Error converting proto any to vhost rate limit plugin")
	}
	ra := out.GetRoute()
	if ra != nil {
		ra.RateLimits = generateCustomEnvoyConfigForVhost(rateLimit.RateLimits)
		ra.IncludeVhRateLimits = &types.BoolValue{Value: rateLimit.IncludeVhRateLimits}
	} else {
		// TODO(yuval-k): maybe reaturn nil here instread and just log a warning?
		return fmt.Errorf("cannot apply rate limits without a route action")
	}

	return nil
}

func (p *Plugin) HttpFilters(params plugins.Params, listener *v1.HttpListener) ([]plugins.StagedHttpFilter, error) {
	if p.upstreamRef == nil {
		return nil, nil
	}

	conf, err := protoutils.MarshalStruct(generateEnvoyConfigForFilter(*p.upstreamRef))
	if err != nil {
		return nil, err
	}

	customConf, err := protoutils.MarshalStruct(generateEnvoyConfigForCustomFilter(*p.upstreamRef))
	if err != nil {
		return nil, err
	}

	return []plugins.StagedHttpFilter{
		{
			HttpFilter: &envoyhttp.HttpFilter{
				Name: filterName,
				ConfigType: &envoyhttp.HttpFilter_Config{
					Config: customConf,
				},
			},
			Stage: filterStage,
		},
		{
			HttpFilter: &envoyhttp.HttpFilter{
				Name: filterName,
				ConfigType: &envoyhttp.HttpFilter_Config{
					Config: conf,
				},
			},
			Stage: filterStage,
		},
	}, nil
}
