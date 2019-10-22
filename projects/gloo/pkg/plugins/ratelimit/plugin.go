package ratelimit

import (
	"time"

	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/plugins/ratelimit"
	"github.com/solo-io/go-utils/errors"

	envoyroute "github.com/envoyproxy/go-control-plane/envoy/api/v2/route"

	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins"
	rlplugin "github.com/solo-io/gloo/projects/gloo/pkg/plugins/ratelimit"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins/utils"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/extauth"
)

/*
Background:

Currently the way to create descriptors in envoy is somewhat limited.
Even though we can use the server configuration to express many forms of rate limits, we
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

The first action check checks if the Authorization header is present. If it is we assume we can trust it, as the request
should have passed an auth filter first.

If the header is present, the second one gets the header so we can rate limit on a per user basis.

If not (the second action/generated descriptor), then the remote address is retrieved so we can limit per IP.

Given this envoy configuration, the appropriate server configuration would be:

descriptors:
- key: generic_key
  value: <vhost_name>
  descriptors:
  - key: header_match
    value: not-authenticated
    descriptors:
    - key: remote_address
      rate_limit:
        unit: MINUTE
        requests_per_unit: 3
  - key: header_match
    value: is-authenticated
    descriptors:
    - key: userid
      rate_limit:
        unit: MINUTE
        requests_per_unit: 10
*/

const (
	IngressDomain = "ingress"
	userid        = "userid"

	authenticated = "is-authenticated"
	anonymous     = "not-authenticated"

	stage = 0

	headerMatch   = "header_match"
	genericKey    = "generic_key"
	remoteAddress = "remote_address"
)

type Plugin struct {
	upstreamRef *core.ResourceRef
	timeout     *time.Duration
	denyOnFail  bool

	authUserIdHeader string
}

func NewPlugin() *Plugin {
	return &Plugin{}
}

//// TODO(yuval-k): Copied from ext auth. the real solution is to add it to upstream Gloo
type tmpPluginContainer struct {
	params plugins.InitParams
}

func (t *tmpPluginContainer) GetExtensions() *v1.Extensions {
	return t.params.ExtensionsSettings
}

//TODO(kdorosh) delete once support for old config is dropped
func (p *Plugin) handleDeprecatedPluginConfig(params plugins.InitParams) error {
	var settings ratelimit.Settings
	p.upstreamRef = nil
	err := utils.UnmarshalExtension(&tmpPluginContainer{params}, rlplugin.ExtensionName, &settings)
	if err != nil {
		if err == utils.NotFoundError {
			return nil
		}
		return err
	}

	p.upstreamRef = settings.RatelimitServerRef
	p.timeout = settings.RequestTimeout
	p.denyOnFail = settings.DenyOnFail
	return nil
}

func (p *Plugin) Init(params plugins.InitParams) error {
	authSettings, _, _ := extauth.GetSettings(params)
	p.authUserIdHeader = extauth.GetAuthHeader(authSettings)

	if err := p.handleDeprecatedPluginConfig(params); err != nil {
		return err
	}

	if rlServer := params.Settings.GetRatelimitServer(); rlServer != nil {
		p.upstreamRef = rlServer.RatelimitServerRef
		p.timeout = rlServer.RequestTimeout
		p.denyOnFail = rlServer.DenyOnFail
	}

	return nil
}

func (p *Plugin) ProcessVirtualHost(params plugins.VirtualHostParams, in *v1.VirtualHost, out *envoyroute.VirtualHost) error {
	return p.ProcessVirtualHostSimple(params, in, out)
}

//TODO(kdorosh) delete once support for old config is dropped
func (p *Plugin) handleDeprecatedVirtualHostSimple(params plugins.VirtualHostParams, in *v1.VirtualHost) (*ratelimit.IngressRateLimit, error) {
	var rateLimit ratelimit.IngressRateLimit
	err := utils.UnmarshalExtension(in.VirtualHostPlugins, rlplugin.ExtensionName, &rateLimit)
	if err != nil {
		if err == utils.NotFoundError {
			return nil, nil
		}
		return nil, errors.Wrapf(err, "Error converting proto to ingress rate limit plugin")
	}
	return &rateLimit, nil
}

func (p *Plugin) ProcessVirtualHostSimple(params plugins.VirtualHostParams, in *v1.VirtualHost, out *envoyroute.VirtualHost) error {
	rateLimit, err := p.handleDeprecatedVirtualHostSimple(params, in)
	if err != nil {
		return err
	}

	if rl := in.GetVirtualHostPlugins().GetRatelimitBasic(); rl != nil {
		rateLimit = rl
	}

	if rateLimit == nil {
		// no rate limit virtual host config found, nothing to do here
		return nil
	}

	_, err = TranslateUserConfigToRateLimitServerConfig(in.Name, *rateLimit)
	if err != nil {
		return err
	}

	out.RateLimits = generateEnvoyConfigForVhost(in.Name, p.authUserIdHeader)
	return nil
}

func (p *Plugin) HttpFilters(params plugins.Params, listener *v1.HttpListener) ([]plugins.StagedHttpFilter, error) {
	if p.upstreamRef == nil {
		return nil, nil
	}

	conf := generateEnvoyConfigForFilter(*p.upstreamRef, p.timeout, p.denyOnFail)
	stagedFilter, err := plugins.NewStagedFilterWithConfig(rlplugin.FilterName, conf, rlplugin.FilterStage)
	if err != nil {
		return nil, err
	}

	return []plugins.StagedHttpFilter{
		stagedFilter,
	}, nil
}
