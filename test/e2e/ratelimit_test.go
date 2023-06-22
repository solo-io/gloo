package e2e_test

import (
	"fmt"
	"net"
	"os"

	v1alpha1 "github.com/solo-io/gloo/projects/gloo/pkg/api/external/solo/ratelimit"

	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/core/matchers"
	extauthpb "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/extauth/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/ratelimit"
	rlv1alpha1 "github.com/solo-io/solo-apis/pkg/api/ratelimit.solo.io/v1alpha1"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
)

// This file is a remnant of the old rate limiting e2e tests implementation.
// Only utilities exist in this file, and we should migrate them to either:
// 1. A shared place for all tests (test/gomega)
// 2. The file that it is used in

var (
	baseRateLimitPort = uint32(18081)
)

func getRedisPath() string {
	binaryPath := os.Getenv("REDIS_BINARY")
	if binaryPath != "" {
		return binaryPath
	}
	return "redis-server"
}

type rateLimitingProxyBuilder struct {
	port              uint32
	virtualHostConfig map[string]virtualHostConfig
	// Will be used for all routes
	routeAction *gloov1.Route_RouteAction
}

type routeConfig struct {
	prefix                         string
	extAuth                        *extauthpb.ExtAuthExtension
	ingressRateLimit               *ratelimit.IngressRateLimit
	rateLimitConfigRef             *core.ResourceRef
	earlyStageRateLimitConfigRef   *core.ResourceRef
	regularStageRateLimitConfigRef *core.ResourceRef
}

type virtualHostConfig struct {
	// A simple catch-all route to the target upstream will always be appended to this slice
	routes  []routeConfig
	extAuth *extauthpb.ExtAuthExtension
	// Check the builder implementation to see the supported config types
	rateLimitConfig interface{}
}

func newRateLimitingProxyBuilder(port uint32, targetUpstream *core.ResourceRef) *rateLimitingProxyBuilder {
	return &rateLimitingProxyBuilder{
		port: port,
		routeAction: &gloov1.Route_RouteAction{
			RouteAction: &gloov1.RouteAction{
				Destination: &gloov1.RouteAction_Single{
					Single: &gloov1.Destination{
						DestinationType: &gloov1.Destination_Upstream{
							Upstream: targetUpstream,
						},
					},
				},
			},
		},
		virtualHostConfig: make(map[string]virtualHostConfig),
	}
}

func (b *rateLimitingProxyBuilder) withVirtualHost(domain string, config virtualHostConfig) *rateLimitingProxyBuilder {
	if _, ok := b.virtualHostConfig[domain]; ok {
		panic("already have a virtual host with domain: " + domain)
	}

	b.virtualHostConfig[domain] = config
	return b
}

func (b *rateLimitingProxyBuilder) build() *gloov1.Proxy {
	var virtualHosts []*gloov1.VirtualHost
	for domain, vhostConfig := range b.virtualHostConfig {

		vhost := &gloov1.VirtualHost{
			Name:    "gloo-system_" + domain,
			Domains: []string{domain},
			Options: &gloov1.VirtualHostOptions{},
			Routes:  []*gloov1.Route{},
		}

		if vhostConfig.extAuth != nil {
			vhost.Options.Extauth = vhostConfig.extAuth
		}

		switch rateLimitConfig := vhostConfig.rateLimitConfig.(type) {
		case *v1alpha1.RateLimitConfig:
			vhost.Options.RateLimitConfigType = &gloov1.VirtualHostOptions_RateLimitConfigs{
				RateLimitConfigs: &ratelimit.RateLimitConfigRefs{
					Refs: []*ratelimit.RateLimitConfigRef{
						{
							Namespace: rateLimitConfig.GetNamespace(),
							Name:      rateLimitConfig.GetName(),
						},
					},
				},
			}

		case []*rlv1alpha1.RateLimitActions:
			vhost.Options.RateLimitConfigType = &gloov1.VirtualHostOptions_Ratelimit{
				Ratelimit: &ratelimit.RateLimitVhostExtension{
					RateLimits: rateLimitConfig,
				},
			}
		case *gloov1.VirtualHostOptions_Ratelimit:
			vhost.Options.RateLimitConfigType = rateLimitConfig

		case *gloov1.VirtualHostOptions_RatelimitEarly:
			vhost.Options.RateLimitEarlyConfigType = rateLimitConfig

		case *gloov1.VirtualHostOptions_RateLimitEarlyConfigs:
			vhost.Options.RateLimitEarlyConfigType = rateLimitConfig

		case *ratelimit.IngressRateLimit:
			vhost.Options.RatelimitBasic = rateLimitConfig
		case nil:
			break
		default:
			panic("unexpected rate limit config type")
		}

		for i, routeCfg := range vhostConfig.routes {

			var match []*matchers.Matcher
			if routeCfg.prefix != "" {
				match = []*matchers.Matcher{{
					PathSpecifier: &matchers.Matcher_Prefix{
						Prefix: routeCfg.prefix,
					},
				}}
			}

			routeOptions := &gloov1.RouteOptions{}
			if routeCfg.ingressRateLimit != nil {
				routeOptions.RatelimitBasic = routeCfg.ingressRateLimit
			}
			if routeCfg.earlyStageRateLimitConfigRef != nil {
				routeOptions.RateLimitEarlyConfigType = &gloov1.RouteOptions_RateLimitEarlyConfigs{
					RateLimitEarlyConfigs: &ratelimit.RateLimitConfigRefs{
						Refs: []*ratelimit.RateLimitConfigRef{
							{
								Name:      routeCfg.earlyStageRateLimitConfigRef.Name,
								Namespace: routeCfg.earlyStageRateLimitConfigRef.Namespace,
							},
						},
					},
				}
			}
			if routeCfg.regularStageRateLimitConfigRef != nil {
				routeOptions.RateLimitRegularConfigType = &gloov1.RouteOptions_RateLimitRegularConfigs{
					RateLimitRegularConfigs: &ratelimit.RateLimitConfigRefs{
						Refs: []*ratelimit.RateLimitConfigRef{
							{
								Name:      routeCfg.regularStageRateLimitConfigRef.Name,
								Namespace: routeCfg.regularStageRateLimitConfigRef.Namespace,
							},
						},
					},
				}
			}
			if routeCfg.rateLimitConfigRef != nil {
				routeOptions.RateLimitConfigType = &gloov1.RouteOptions_RateLimitConfigs{
					RateLimitConfigs: &ratelimit.RateLimitConfigRefs{
						Refs: []*ratelimit.RateLimitConfigRef{
							{
								Name:      routeCfg.rateLimitConfigRef.Name,
								Namespace: routeCfg.rateLimitConfigRef.Namespace,
							},
						},
					},
				}
			}
			if routeCfg.extAuth != nil {
				routeOptions.Extauth = routeCfg.extAuth
			}

			vhost.Routes = append(vhost.Routes, &gloov1.Route{
				// Name is required for `RateLimitBasic` config to work
				Name:     fmt.Sprintf("gloo-system_route-%s-%d", domain, i),
				Matchers: match,
				Action:   b.routeAction,
				Options:  routeOptions,
			})
		}

		// Add a fallback route to the target upstream
		vhost.Routes = append(vhost.Routes, &gloov1.Route{
			Action: b.routeAction,
		})

		virtualHosts = append(virtualHosts, vhost)
	}

	return &gloov1.Proxy{
		Metadata: &core.Metadata{
			Name:      "proxy",
			Namespace: "default",
		},
		Listeners: []*gloov1.Listener{
			{
				Name:        "e2e-test-listener",
				BindAddress: net.IPv4zero.String(),
				BindPort:    b.port,
				ListenerType: &gloov1.Listener_HttpListener{
					HttpListener: &gloov1.HttpListener{
						VirtualHosts: virtualHosts,
					},
				},
			},
		},
	}
}
