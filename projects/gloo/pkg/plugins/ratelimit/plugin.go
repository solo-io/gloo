package ratelimit

import (
	"time"

	envoyroute "github.com/envoyproxy/go-control-plane/envoy/api/v2/route"
	envoyhttp "github.com/envoyproxy/go-control-plane/envoy/config/filter/network/http_connection_manager/v2"
	"github.com/solo-io/solo-kit/pkg/utils/protoutils"

	"github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/plugins"
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
	domain      = "ingress"
	requestType = "external"
	userid      = "userid"

	authenticated = "is-authenticated"
	anonymous     = "not-authenticated"

	stage   = 0
	timeout = 20 * time.Millisecond

	headerMatch   = "header_match"
	remoteAddress = "remote_address"
)

const (
	filterName = "envoy.rate_limit"
	// rate limiting should happen after auth
	filterStage = plugins.PostInAuth
)

type Plugin struct {
}

func NewPlugin() plugins.Plugin {
	return &Plugin{}
}

func (p *Plugin) Init(params plugins.InitParams) error {

	return nil
}

func (p *Plugin) ProcessVirtualHost(params plugins.Params, in *v1.VirtualHost, out *envoyroute.VirtualHost) error {
	if in.VirtualHostPlugins == nil {
		return nil
	}
	if in.VirtualHostPlugins.RateLimits == nil {
		return nil
	}
	_, err := TranslateUserConfigToRateLimitServerConfig(*in.VirtualHostPlugins.RateLimits)
	if err != nil {
		return err
	}

	vhost := generateEnvoyConfigForVhost(in.VirtualHostPlugins.RateLimits.AuthrorizedHeader)
	out.RateLimits = vhost

	return nil
}

func (p *Plugin) HttpFilters(params plugins.Params, listener *v1.HttpListener) ([]plugins.StagedHttpFilter, error) {
	conf, err := protoutils.MarshalStruct(generateEnvoyConfigForFilter())
	if err != nil {
		return nil, err
	}
	return []plugins.StagedHttpFilter{
		{
			HttpFilter: &envoyhttp.HttpFilter{Name: filterName,
				Config: conf},
			Stage: filterStage,
		},
	}, nil
}
