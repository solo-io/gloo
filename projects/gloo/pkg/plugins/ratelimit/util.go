package ratelimit

import (
	"errors"
	"time"

	envoyvhostratelimit "github.com/envoyproxy/go-control-plane/envoy/api/v2/route"
	envoyratelimit "github.com/envoyproxy/go-control-plane/envoy/config/filter/http/rate_limit/v2"
	"github.com/gogo/protobuf/types"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/plugins/ratelimit"
	rlplugin "github.com/solo-io/gloo/projects/gloo/pkg/plugins/ratelimit"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
)

/*
translate virtual hosts
save them
then translate get rate limit configs
*/

func TranslateUserConfigToRateLimitServerConfig(vhostname string, ingressRl ratelimit.IngressRateLimit) (*ratelimit.Descriptor, error) {

	vhostDescriptor := &ratelimit.Descriptor{
		Key:         genericKey,
		Value:       vhostname,
		Descriptors: []*ratelimit.Descriptor{},
	}

	if ingressRl.AnonymousLimits != nil {

		if ingressRl.AnonymousLimits.Unit == ratelimit.RateLimit_UNKNOWN {
			return nil, errors.New("unknown unit for anonymous config")
		}

		c := &ratelimit.Descriptor{
			Key:   headerMatch,
			Value: anonymous,
			Descriptors: []*ratelimit.Descriptor{
				{
					Key:       remoteAddress,
					RateLimit: ingressRl.AnonymousLimits,
				},
			},
		}

		vhostDescriptor.Descriptors = append(vhostDescriptor.Descriptors, c)
	}

	if ingressRl.AuthorizedLimits != nil {

		if ingressRl.AuthorizedLimits.Unit == ratelimit.RateLimit_UNKNOWN {
			return nil, errors.New("unknown unit for authenticated config")
		}

		c := &ratelimit.Descriptor{
			Key:   headerMatch,
			Value: authenticated,
			Descriptors: []*ratelimit.Descriptor{
				{
					Key:       userid,
					RateLimit: ingressRl.AuthorizedLimits,
				},
			},
		}
		vhostDescriptor.Descriptors = append(vhostDescriptor.Descriptors, c)
	}

	return vhostDescriptor, nil
}

func generateEnvoyConfigForFilter(ref core.ResourceRef, timeout *time.Duration, denyOnFail bool) *envoyratelimit.RateLimit {
	return rlplugin.GenerateEnvoyConfigForFilterWith(ref, IngressDomain, stage, timeout, denyOnFail)
}

func generateEnvoyConfigForVhost(vhostname, headername string) []*envoyvhostratelimit.RateLimit {
	// the filter config, virtual host config are always the same:

	empty := headername == ""
	if empty {
		// TODO(yuval-k): fix this hack
		headername = "not-a-header"
	}

	vhostAction := getPerVhostRateLimit(vhostname)

	getAuthRateLimits := func(b bool) *envoyvhostratelimit.RateLimit_Action { return getAuthHeaderRateLimit(headername, b) }

	vhostrl := []*envoyvhostratelimit.RateLimit{
		{
			Stage: &types.UInt32Value{Value: stage},
			Actions: []*envoyvhostratelimit.RateLimit_Action{
				vhostAction,
				getAuthRateLimits(true),
				getUserIdRateLimit(headername),
			},
		},
		{
			Stage: &types.UInt32Value{Value: stage},
			Actions: []*envoyvhostratelimit.RateLimit_Action{
				vhostAction,
				getAuthRateLimits(false),
				getPerIpRateLimit(),
			},
		},
	}
	return vhostrl
}

func getPerVhostRateLimit(vhostname string) *envoyvhostratelimit.RateLimit_Action {
	return &envoyvhostratelimit.RateLimit_Action{
		ActionSpecifier: &envoyvhostratelimit.RateLimit_Action_GenericKey_{
			GenericKey: &envoyvhostratelimit.RateLimit_Action_GenericKey{
				DescriptorValue: vhostname,
			},
		},
	}
}

func getAuthHeaderRateLimit(headername string, match bool) *envoyvhostratelimit.RateLimit_Action {

	headersmatcher := []*envoyvhostratelimit.HeaderMatcher{{
		Name:                 headername,
		HeaderMatchSpecifier: &envoyvhostratelimit.HeaderMatcher_PresentMatch{PresentMatch: true},
	}}

	var value string
	if match {
		value = authenticated
	} else {
		value = anonymous
	}

	return &envoyvhostratelimit.RateLimit_Action{
		ActionSpecifier: &envoyvhostratelimit.RateLimit_Action_HeaderValueMatch_{
			HeaderValueMatch: &envoyvhostratelimit.RateLimit_Action_HeaderValueMatch{
				DescriptorValue: value,
				ExpectMatch:     &types.BoolValue{Value: match},
				Headers:         headersmatcher,
			},
		},
	}
}

func getUserIdRateLimit(headername string) *envoyvhostratelimit.RateLimit_Action {
	return &envoyvhostratelimit.RateLimit_Action{
		ActionSpecifier: &envoyvhostratelimit.RateLimit_Action_RequestHeaders_{
			RequestHeaders: &envoyvhostratelimit.RateLimit_Action_RequestHeaders{
				DescriptorKey: userid,
				HeaderName:    headername,
			},
		},
	}
}

func getPerIpRateLimit() *envoyvhostratelimit.RateLimit_Action {
	return &envoyvhostratelimit.RateLimit_Action{
		ActionSpecifier: &envoyvhostratelimit.RateLimit_Action_RemoteAddress_{
			RemoteAddress: &envoyvhostratelimit.RateLimit_Action_RemoteAddress{},
		},
	}
}
