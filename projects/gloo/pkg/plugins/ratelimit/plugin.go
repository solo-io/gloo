package ratelimit

import (
	"errors"
	"time"

	envoyvhostratelimit "github.com/envoyproxy/go-control-plane/envoy/api/v2/route"
	envoyratelimit "github.com/envoyproxy/go-control-plane/envoy/config/filter/http/rate_limit/v2"
	types "github.com/gogo/protobuf/types"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1/plugins/ratelimit"
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

func translateUserConfigToRateLimitServerConfig(userRl ratelimit.UserRateLimit) (*ratelimit.RateLimitConfig, error) {

	rl := &ratelimit.RateLimitConfig{
		Domain: domain,
	}
	if userRl.Anonymous != nil {

		if userRl.Anonymous.Unit == ratelimit.RateLimit_UNKNOWN {
			return nil, errors.New("unknown unit for anonymous config")
		}

		c := &ratelimit.Constraint{
			Key:   headerMatch,
			Value: anonymous,
			Constraints: []*ratelimit.Constraint{
				{
					Key:       remoteAddress,
					RateLimit: userRl.Anonymous,
				},
			},
		}
		rl.Constraints = append(rl.Constraints, c)
	}

	if userRl.Authenticated != nil {

		if userRl.Authenticated.Unit == ratelimit.RateLimit_UNKNOWN {
			return nil, errors.New("unknown unit for authenticated config")
		}

		c := &ratelimit.Constraint{
			Key:   headerMatch,
			Value: authenticated,
			Constraints: []*ratelimit.Constraint{
				{
					Key:       userid,
					RateLimit: userRl.Authenticated,
				},
			},
		}
		rl.Constraints = append(rl.Constraints, c)
	}

	return rl, nil

}

func generateEnvoyConfig(headername string) (*envoyratelimit.RateLimit, []*envoyvhostratelimit.RateLimit) {
	// the filter config, virtual host config are always the same:
	timeout := timeout
	envoyrl := envoyratelimit.RateLimit{
		Domain:      domain,
		Stage:       stage,
		RequestType: requestType,
		Timeout:     &timeout,
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
	return &envoyrl, vhostrl
}
