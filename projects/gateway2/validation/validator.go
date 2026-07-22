package validation

import (
	"context"

	sologatewayv1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/gloosnapshot"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/static"
	gloovalidation "github.com/solo-io/gloo/projects/gloo/pkg/validation"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
)

// TranslateK8sGatewayProxies builds an extremely basic "dummy" gloov1.Proxy object with the provided `res`
// attached in the correct location, i.e. if a RouteOption is passed it will be on the Route or if a VirtualHostOption is
// passed it will be on the VirtualHost
func TranslateK8sGatewayProxies(ctx context.Context, snap *gloosnapshot.ApiSnapshot, res resources.Resource) ([]*gloov1.Proxy, error) {
	us := gloov1.NewUpstream("default", "zzz_fake-upstream-for-gloo-validation")
	us.UpstreamType = &gloov1.Upstream_Static{
		Static: &static.UpstreamSpec{
			Hosts: []*static.Host{
				{Addr: "solo.io", Port: 80},
			},
		},
	}
	snap.UpsertToResourceList(us)

	proxy, route, vhost := buildValidationProxy(us.GetMetadata().Ref())

	// add the policy object we are validating to the correct location
	switch policy := res.(type) {
	case *sologatewayv1.RouteOption:
		route.Options = policy.GetOptions()
	case *sologatewayv1.VirtualHostOption:
		vhost.Options = policy.GetOptions()
	}

	return []*gloov1.Proxy{proxy}, nil
}

// TranslateK8sGatewayProxiesForUpstream builds the same "dummy" gloov1.Proxy as TranslateK8sGatewayProxies,
// but routes to the given `upstream` instead of the fake one. This is how an Upstream is validated in K8s
// Gateway API mode, where the Upstream is otherwise never translated because HTTPRoutes are not in the ApiSnapshot.
// The real Upstream must be the route destination because the plugin chain (for example the grpcjson plugin's
// ProcessUpstream) only runs for Upstreams referenced by a route in the proxy. ValidateGloo adds the Upstream to
// the snapshot before translation, so the reference resolves.
func TranslateK8sGatewayProxiesForUpstream(ctx context.Context, snap *gloosnapshot.ApiSnapshot, upstream *gloov1.Upstream) ([]*gloov1.Proxy, error) {
	proxy, _, _ := buildValidationProxy(upstream.GetMetadata().Ref())
	return []*gloov1.Proxy{proxy}, nil
}

// buildValidationProxy builds the "dummy" gloov1.Proxy shared by the K8s Gateway API validation paths: a single
// catch-all virtual host routing to `usRef`. It returns the route and virtual host so callers can attach a policy.
func buildValidationProxy(usRef *core.ResourceRef) (*gloov1.Proxy, *gloov1.Route, *gloov1.VirtualHost) {
	route := &gloov1.Route{
		Name: "route",
		Action: &gloov1.Route_RouteAction{
			RouteAction: &gloov1.RouteAction{
				Destination: &gloov1.RouteAction_Single{
					Single: &gloov1.Destination{
						DestinationType: &gloov1.Destination_Upstream{
							Upstream: usRef,
						},
					},
				},
			},
		},
	}

	vhost := &gloov1.VirtualHost{
		Name:    "vhost",
		Domains: []string{"*"},
		Routes:  []*gloov1.Route{route},
	}

	aggregateListener := &gloov1.Listener{
		Name:        "aggregate-listener",
		BindAddress: "127.0.0.1",
		BindPort:    8082,
		ListenerType: &gloov1.Listener_AggregateListener{
			AggregateListener: &gloov1.AggregateListener{
				HttpResources: &gloov1.AggregateListener_HttpResources{
					VirtualHosts: map[string]*gloov1.VirtualHost{
						"vhost": vhost,
					},
				},
				HttpFilterChains: []*gloov1.AggregateListener_HttpFilterChain{{
					HttpOptionsRef:  "opts1",
					VirtualHostRefs: []string{"vhost"},
				}},
			},
		},
	}

	proxy := &gloov1.Proxy{
		Metadata: &core.Metadata{
			Name:      "zzz-fake-proxy-for-validation",
			Namespace: "gloo-system",
		},
		Listeners: []*gloov1.Listener{
			aggregateListener,
		},
	}

	return proxy, route, vhost
}

// GetSimpleErrorFromGlooValidation will get the Errors for the provided `proxy` that exist in the provided `reports`
// This function will return with the first error present if multiple reports exist
func GetSimpleErrorFromGlooValidation(
	reports []*gloovalidation.GlooValidationReport,
	proxy *gloov1.Proxy,
) error {
	var errs error

	for _, report := range reports {
		proxyResReport := report.ResourceReports[proxy]
		if proxyResReport.Errors != nil {
			return proxyResReport.Errors
		}
	}

	return errs
}

// GetSimpleErrorFromGlooValidationForUpstream returns the errors reported against the given `upstream` in the
// provided `reports`. Unlike route and virtual host errors, ProcessUpstream errors such as an invalid grpcjson
// protoDescriptorBin are attributed to the Upstream resource rather than the proxy, so they are looked up by the
// Upstream's kind and ref. The pre-sanitization reports are read because the xDS sanitizers demote Upstream
// errors to warnings before the final reports are produced. The first error found is returned if multiple
// reports exist.
func GetSimpleErrorFromGlooValidationForUpstream(
	reports []*gloovalidation.GlooValidationReport,
	upstream *gloov1.Upstream,
) error {
	for _, report := range reports {
		_, usResReport := report.PreSanitizationReports.Find(resources.Kind(upstream), upstream.GetMetadata().Ref())
		if usResReport.Errors != nil {
			return usResReport.Errors
		}
	}

	return nil
}
