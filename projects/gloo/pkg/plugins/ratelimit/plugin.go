package ratelimit

import (
	"errors"
	"time"

	envoyroute "github.com/envoyproxy/go-control-plane/envoy/api/v2/route"
	envoyvhostratelimit "github.com/envoyproxy/go-control-plane/envoy/api/v2/route"
	envoyratelimit "github.com/envoyproxy/go-control-plane/envoy/config/filter/http/rate_limit/v2"
	envoyhttp "github.com/envoyproxy/go-control-plane/envoy/config/filter/network/http_connection_manager/v2"
	"github.com/solo-io/solo-kit/pkg/utils/protoutils"

	types "github.com/gogo/protobuf/types"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1/plugins/ratelimit"
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
	rlconfig *v1.RateLimitConfig
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
	cfg, err := translateUserConfigToRateLimitServerConfig(*in.VirtualHostPlugins.RateLimits)
	if err != nil {
		return err
	}

	vhost := generateEnvoyConfigForVhost(in.VirtualHostPlugins.RateLimits.AuthrorizedHeader)
	out.RateLimits = vhost

	p.rlconfig = cfg
	return nil
}

func (p *Plugin) HttpFilters(params plugins.Params, listener *v1.HttpListener) ([]plugins.StagedHttpFilter, error) {
	conf, err := protoutils.MarshalPbStruct(generateEnvoyConfigForFilter())
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

/*
translate virtual hosts
save them
then translate get rate limit configs
*/

func translateUserConfigToRateLimitServerConfig(ingressRl ratelimit.IngressRateLimit) (*v1.RateLimitConfig, error) {
	rl := &v1.RateLimitConfig{
		Domain: domain,
	}
	if ingressRl.AnonymousLimits != nil {

		if ingressRl.AnonymousLimits.Unit == ratelimit.RateLimit_UNKNOWN {
			return nil, errors.New("unknown unit for anonymous config")
		}

		c := &v1.Constraint{
			Key:   headerMatch,
			Value: anonymous,
			Constraints: []*v1.Constraint{
				{
					Key:       remoteAddress,
					RateLimit: ingressRl.AnonymousLimits,
				},
			},
		}
		rl.Constraints = append(rl.Constraints, c)
	}

	if ingressRl.AuthorizedLimits != nil {

		if ingressRl.AuthorizedLimits.Unit == ratelimit.RateLimit_UNKNOWN {
			return nil, errors.New("unknown unit for authenticated config")
		}

		c := &v1.Constraint{
			Key:   headerMatch,
			Value: authenticated,
			Constraints: []*v1.Constraint{
				{
					Key:       userid,
					RateLimit: ingressRl.AuthorizedLimits,
				},
			},
		}
		rl.Constraints = append(rl.Constraints, c)
	}

	return rl, nil
}

func generateEnvoyConfigForFilter() *envoyratelimit.RateLimit {
	timeout := timeout
	envoyrl := envoyratelimit.RateLimit{
		Domain:      domain,
		Stage:       stage,
		RequestType: requestType,
		Timeout:     &timeout,
	}
	return &envoyrl
}

func generateEnvoyConfigForVhost(headername string) []*envoyvhostratelimit.RateLimit {
	// the filter config, virtual host config are always the same:

	empty := headername == ""

	if empty {
		// TODO(yuval-k): fix this hack
		headername = "not-a-header"
	}

	headersmatcher := []*envoyvhostratelimit.HeaderMatcher{{
		Name:                 headername,
		HeaderMatchSpecifier: &envoyvhostratelimit.HeaderMatcher_PresentMatch{PresentMatch: true},
	}}
	//:[{"name":"Authorization", "present_match":true}]
	vhostrl := []*envoyvhostratelimit.RateLimit{
		{
			Stage: &types.UInt32Value{Value: stage},
			Actions: []*envoyvhostratelimit.RateLimit_Action{
				{
					ActionSpecifier: &envoyvhostratelimit.RateLimit_Action_HeaderValueMatch_{
						HeaderValueMatch: &envoyvhostratelimit.RateLimit_Action_HeaderValueMatch{

							DescriptorValue: authenticated,
							ExpectMatch:     &types.BoolValue{Value: true},
							Headers:         headersmatcher,
						},
					},
				},
				{
					ActionSpecifier: &envoyvhostratelimit.RateLimit_Action_RequestHeaders_{
						RequestHeaders: &envoyvhostratelimit.RateLimit_Action_RequestHeaders{
							DescriptorKey: userid,
							HeaderName:    headername,
						},
					},
				},
			},
		},
		{
			Stage: &types.UInt32Value{Value: stage},
			Actions: []*envoyvhostratelimit.RateLimit_Action{
				{
					ActionSpecifier: &envoyvhostratelimit.RateLimit_Action_HeaderValueMatch_{
						HeaderValueMatch: &envoyvhostratelimit.RateLimit_Action_HeaderValueMatch{
							DescriptorValue: anonymous,
							ExpectMatch:     &types.BoolValue{Value: false},
							Headers:         headersmatcher,
						},
					},
				},
				{
					ActionSpecifier: &envoyvhostratelimit.RateLimit_Action_RemoteAddress_{
						RemoteAddress: &envoyvhostratelimit.RateLimit_Action_RemoteAddress{},
					},
				},
			},
		},
	}
	return vhostrl
}
