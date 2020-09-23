package translation

import (
	envoyratelimit "github.com/envoyproxy/go-control-plane/envoy/api/v2/route"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/rotisserie/eris"
	rl_opts "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/ratelimit"
	solo_api_rl_types "github.com/solo-io/solo-apis/pkg/api/ratelimit.solo.io/v1alpha1"
	"github.com/solo-io/solo-projects/projects/rate-limit/pkg/translation/internal"
)

//go:generate mockgen -destination ./mocks/mock_interfaces.go -source ./basic.go

// Original code: https://github.com/solo-io/solo-projects/blob/cf2e4e3a9c33b0189adb9a0cfae154e1980c1e77/projects/gloo/pkg/plugins/ratelimit/plugin.go#L17

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
  value: <resource_name>
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
type BasicRateLimitTranslator interface {
	GenerateServerConfig(resourceName string, ingressRl rl_opts.IngressRateLimit) (*solo_api_rl_types.Descriptor, error)
	GenerateResourceConfig(resourceName, headerName string, stage uint32) []*envoyratelimit.RateLimit
}

type translator struct{}

func NewBasicRateLimitTranslator() BasicRateLimitTranslator {
	return translator{}
}

func (translator) GenerateServerConfig(resourceName string, ingressRl rl_opts.IngressRateLimit) (*solo_api_rl_types.Descriptor, error) {
	rootDescriptor := &solo_api_rl_types.Descriptor{
		Key:         internal.GenericKey,
		Value:       resourceName,
		Descriptors: []*solo_api_rl_types.Descriptor{},
	}

	if ingressRl.AnonymousLimits != nil {

		if ingressRl.AnonymousLimits.Unit == solo_api_rl_types.RateLimit_UNKNOWN {
			return nil, eris.New("unknown unit for anonymous config")
		}

		c := &solo_api_rl_types.Descriptor{
			Key:   internal.HeaderMatch,
			Value: internal.Anonymous,
			Descriptors: []*solo_api_rl_types.Descriptor{
				{
					Key:       internal.RemoteAddress,
					RateLimit: ingressRl.AnonymousLimits,
				},
			},
		}

		rootDescriptor.Descriptors = append(rootDescriptor.Descriptors, c)
	}

	if ingressRl.AuthorizedLimits != nil {

		if ingressRl.AuthorizedLimits.Unit == solo_api_rl_types.RateLimit_UNKNOWN {
			return nil, eris.New("unknown unit for authenticated config")
		}

		c := &solo_api_rl_types.Descriptor{
			Key:   internal.HeaderMatch,
			Value: internal.Authenticated,
			Descriptors: []*solo_api_rl_types.Descriptor{
				{
					Key:       internal.UserId,
					RateLimit: ingressRl.AuthorizedLimits,
				},
			},
		}
		rootDescriptor.Descriptors = append(rootDescriptor.Descriptors, c)
	}

	return rootDescriptor, nil
}

func (translator) GenerateResourceConfig(resourceName, headerName string, stage uint32) []*envoyratelimit.RateLimit {
	// the filter config, virtual host config are always the same:

	if headerName == "" {
		// TODO(yuval-k): fix this hack
		headerName = "not-a-header"
	}

	action := getPerResourceRateLimit(resourceName)

	getAuthRateLimits := func(b bool) *envoyratelimit.RateLimit_Action { return getAuthHeaderRateLimit(headerName, b) }

	rateLimits := []*envoyratelimit.RateLimit{
		{
			Stage: &wrappers.UInt32Value{Value: stage},
			Actions: []*envoyratelimit.RateLimit_Action{
				action,
				getAuthRateLimits(true),
				getUserIdRateLimit(headerName),
			},
		},
		{
			Stage: &wrappers.UInt32Value{Value: stage},
			Actions: []*envoyratelimit.RateLimit_Action{
				action,
				getAuthRateLimits(false),
				getPerIpRateLimit(),
			},
		},
	}
	return rateLimits
}

func getPerResourceRateLimit(resourceName string) *envoyratelimit.RateLimit_Action {
	return &envoyratelimit.RateLimit_Action{
		ActionSpecifier: &envoyratelimit.RateLimit_Action_GenericKey_{
			GenericKey: &envoyratelimit.RateLimit_Action_GenericKey{
				DescriptorValue: resourceName,
			},
		},
	}
}

func getAuthHeaderRateLimit(headername string, match bool) *envoyratelimit.RateLimit_Action {

	headersmatcher := []*envoyratelimit.HeaderMatcher{{
		Name:                 headername,
		HeaderMatchSpecifier: &envoyratelimit.HeaderMatcher_PresentMatch{PresentMatch: true},
	}}

	var value string
	if match {
		value = internal.Authenticated
	} else {
		value = internal.Anonymous
	}

	return &envoyratelimit.RateLimit_Action{
		ActionSpecifier: &envoyratelimit.RateLimit_Action_HeaderValueMatch_{
			HeaderValueMatch: &envoyratelimit.RateLimit_Action_HeaderValueMatch{
				DescriptorValue: value,
				ExpectMatch:     &wrappers.BoolValue{Value: match},
				Headers:         headersmatcher,
			},
		},
	}
}

func getUserIdRateLimit(headername string) *envoyratelimit.RateLimit_Action {
	return &envoyratelimit.RateLimit_Action{
		ActionSpecifier: &envoyratelimit.RateLimit_Action_RequestHeaders_{
			RequestHeaders: &envoyratelimit.RateLimit_Action_RequestHeaders{
				DescriptorKey: internal.UserId,
				HeaderName:    headername,
			},
		},
	}
}

func getPerIpRateLimit() *envoyratelimit.RateLimit_Action {
	return &envoyratelimit.RateLimit_Action{
		ActionSpecifier: &envoyratelimit.RateLimit_Action_RemoteAddress_{
			RemoteAddress: &envoyratelimit.RateLimit_Action_RemoteAddress{},
		},
	}
}
