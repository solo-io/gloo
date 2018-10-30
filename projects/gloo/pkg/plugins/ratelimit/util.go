package ratelimit

import (
	"errors"

	envoyvhostratelimit "github.com/envoyproxy/go-control-plane/envoy/api/v2/route"
	envoyratelimit "github.com/envoyproxy/go-control-plane/envoy/config/filter/http/rate_limit/v2"

	types "github.com/gogo/protobuf/types"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1/plugins/ratelimit"
)

/*
translate virtual hosts
save them
then translate get rate limit configs
*/

func TranslateUserConfigToRateLimitServerConfig(ingressRl ratelimit.IngressRateLimit) (*v1.RateLimitConfig, error) {
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
