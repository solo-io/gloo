package utils

import (
	"context"
	"fmt"
	"strconv"
	"time"

	gatewayv1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/core/matchers"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
)

// creates a virtual service with the given selector, that routes to the specified upstream
func CreateVirtualServiceForUpstream(
	ctx context.Context,
	ns string,
	virtualServiceClient gatewayv1.VirtualServiceClient,
	selector map[string]string,
	upstreamRefs ...*core.ResourceRef,
) (*gatewayv1.VirtualService, error) {

	timestamp := strconv.FormatInt(time.Now().UnixNano(), 10)

	var routes []*gatewayv1.Route
	for idx, ref := range upstreamRefs {
		var routeMatchers []*matchers.Matcher
		if idx == len(upstreamRefs)-1 {
			routeMatchers = []*matchers.Matcher{
				{
					PathSpecifier: &matchers.Matcher_Prefix{
						Prefix: "/",
					},
				},
			}
		} else {
			routeMatchers = []*matchers.Matcher{
				{
					PathSpecifier: &matchers.Matcher_Prefix{
						Prefix: fmt.Sprintf("/%d", idx),
					},
				},
			}
		}

		routes = append(routes, &gatewayv1.Route{
			Options: &gloov1.RouteOptions{
				HostRewriteType: &gloov1.RouteOptions_HostRewrite{
					HostRewrite: "httpbin.org",
				},
			},
			Matchers: routeMatchers,
			Action: &gatewayv1.Route_RouteAction{
				RouteAction: &gloov1.RouteAction{
					Destination: &gloov1.RouteAction_Single{
						Single: &gloov1.Destination{
							DestinationType: &gloov1.Destination_Upstream{
								Upstream: ref,
							},
						},
					},
				},
			},
		})
	}

	vs := &gatewayv1.VirtualService{
		Metadata: &core.Metadata{
			Name:      "vs-" + timestamp,
			Namespace: ns,
			Labels:    selector,
		},
		VirtualHost: &gatewayv1.VirtualHost{
			Domains: []string{"domain-" + timestamp},
			Routes:  routes,
		},
	}

	return virtualServiceClient.Write(vs, clients.WriteOpts{Ctx: ctx})
}
