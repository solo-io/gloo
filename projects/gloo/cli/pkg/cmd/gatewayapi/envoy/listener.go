package envoy

import (
	"fmt"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/extensions/jwt"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/extensions/waf"
	extv1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/extauth/v1"
	ratelimitv1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/ratelimit"
	"k8s.io/utils/ptr"

	listenerv3 "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	ext_authzv3 "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/ext_authz/v3"
	grpcstats "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/grpc_stats/v3"
	ratelimitv3 "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/ratelimit/v3"
	rbacv3 "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/rbac/v3"
	envoy_hcm "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/network/http_connection_manager/v3"
	envoytcp "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/network/tcp_proxy/v3"
	aws "github.com/solo-io/envoy-gloo/go/config/filter/http/aws_lambda/v2"
	transformation "github.com/solo-io/envoy-gloo/go/config/filter/http/transformation/v2"
	glooaws "github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/extensions/aws"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/wrapperspb"
	gwv1 "sigs.k8s.io/gateway-api/apis/v1"
)

func (o *Outputs) convertListener(l *listenerv3.Listener) ([]gwv1.Listener, []RouteInfo) {
	// translate each filter chain

	baseListener := gwv1.Listener{
		Name: gwv1.SectionName(fmt.Sprintf("listener-%d", l.GetAddress().GetSocketAddress().GetPortValue())),
		Port: gwv1.PortNumber(l.GetAddress().GetSocketAddress().GetPortValue()),
		AllowedRoutes: &gwv1.AllowedRoutes{
			Namespaces: &gwv1.RouteNamespaces{
				From: ptr.To(gwv1.FromNamespaces("gloo-system")),
			},
			Kinds: []gwv1.RouteGroupKind{
				{
					Kind: "HTTPRoute",
				},
			},
		},
	}
	var rds []RouteInfo
	var lds []gwv1.Listener
	for _, fc := range l.FilterChains {
		var snis []string
		if m := fc.GetFilterChainMatch(); m != nil {
			snis = m.ServerNames
			m.ServerNames = nil
			if !isEmpty(m) {
				o.Errors = append(o.Errors, fmt.Errorf("filter chain match is not empty. %s", l.Name))
			}
		}

		if len(snis) == 0 {
			snis = append(snis, "")
		}

		filters := fc.Filters
		if len(filters) != 1 {
			o.Errors = append(o.Errors, fmt.Errorf("filter chain has more than one filter. %s", l.Name))
		}
		lastFilter := filters[len(filters)-1]
		filter, err := convertAny(lastFilter.GetTypedConfig())
		if err != nil {
			o.Errors = append(o.Errors, fmt.Errorf("error unmarshalling last filter: %v", err))
		}

		isTls := fc.TransportSocket != nil // TODO: check if it's a TLS transport socket

		for _, sni := range snis {
			listener := baseListener.DeepCopy()
			if sni != "" {
				s := gwv1.Hostname(sni)
				listener.Hostname = &s
			}
			switch filter := filter.(type) {
			case *envoy_hcm.HttpConnectionManager:
				if isTls {
					listener.Protocol = gwv1.HTTPSProtocolType
				} else {
					listener.Protocol = gwv1.HTTPProtocolType
				}

				// TODO: add the hcm filter here too; we need to make sure
				// it has the options we expect, or panic
				policies := o.convertListenerPolicy(filter.HttpFilters)
				rds = append(rds, RouteInfo{Rds: filter.GetRds().RouteConfigName, FiltersOnChain: policies})

			case *envoytcp.TcpProxy:
				if isTls {
					listener.Protocol = gwv1.TLSProtocolType
				} else {
					listener.Protocol = gwv1.TCPProtocolType
				}
			}
			lds = append(lds, *listener)
		}

	}

	return lds, rds
}

func (o *Outputs) convertListenerPolicy(filters []*envoy_hcm.HttpFilter) map[string][]proto.Message {
	if o.Settings.Gloo == nil {
		o.Settings.Gloo = &v1.GlooOptions{}
	}

	ret := make(map[string][]proto.Message)
	for _, f := range filters {

		if f.GetName() == "istio.stats" || f.GetName() == "istio.alpn" || f.GetName() == "istio.metadata_exchange" {
			continue
		}

		filter, err := convertAny(f.GetTypedConfig())
		if err != nil {
			o.Errors = append(o.Errors, fmt.Errorf("error unmarshalling last filter: %v", err))
			//TODO commented out due to unknown field "normalize_payload_in_metadata"
			//panic(err)
			continue
		}

		if filter == nil || isEmpty(filter) {
			continue
		}
		ret[f.GetName()] = append(ret[f.GetName()], filter)
		switch filter := filter.(type) {
		case *rbacv3.RBAC:

			o.Errors = append(o.Errors, fmt.Errorf("unhandled filter config %T", filter))
		case *grpcstats.FilterConfig:

		case *waf.ModSecurity:
			if filter.GetDisabled() {
				continue
			}

			o.Errors = append(o.Errors, fmt.Errorf("unhandled filter config %T", filter))
		case *jwt.JwtWithStage:
			// Grab all the jwt providers and give them to the per route config so we can cross reference them
			// to create the jwt policy
			// handled in the orute
		case *ratelimitv3.RateLimit:
			originalFilter := proto.Clone(filter).(*ratelimitv3.RateLimit)
			o.Settings.RatelimitServer = &ratelimitv1.Settings{
				RatelimitServerRef: getRef(filter.GetRateLimitService().GetGrpcService().GetEnvoyGrpc().ClusterName),
				RequestTimeout:     filter.Timeout,
				ServiceType: &ratelimitv1.Settings_GrpcService{
					GrpcService: &ratelimitv1.GrpcService{
						Authority: filter.GetRateLimitService().GetGrpcService().GetEnvoyGrpc().Authority,
					},
				},
			}
			if filter.Stage == 1 {
				o.Settings.RatelimitServer.RateLimitBeforeAuth = true
				filter.Stage = 0
			}

			filter.RateLimitService = nil
			filter.Timeout = nil

			// this is not exposed to users
			filter.Domain = ""
			filter.RequestType = ""
			if !isEmpty(filter) {
				o.Errors = append(o.Errors, fmt.Errorf("unhandled filter config %T", originalFilter))
			}
		case *ext_authzv3.ExtAuthz:
			if o.Settings.Extauth == nil {
				o.Settings.Extauth = &extv1.Settings{}
			}
			if gs := filter.GetGrpcService(); gs != nil {
				o.Settings.Extauth.ServiceType = &extv1.Settings_GrpcService{
					GrpcService: &extv1.GrpcService{
						Authority: gs.GetEnvoyGrpc().Authority,
					},
				}

				o.Settings.Extauth.ExtauthzServerRef = getRef(gs.GetEnvoyGrpc().ClusterName)
				o.Settings.Extauth.RequestTimeout = filter.GetGrpcService().Timeout
				filter.Services = nil
			}
			o.Settings.Extauth.StatPrefix = filter.StatPrefix
			filter.StatPrefix = ""
			o.Settings.Extauth.ClearRouteCache = filter.ClearRouteCache
			filter.ClearRouteCache = false

			// this config is not exposed to users.
			filter.MetadataContextNamespaces = nil
			filter.TransportApiVersion = 0

			if !isEmpty(filter) {
				o.Errors = append(o.Errors, fmt.Errorf("unhandled filter config %T", filter))
			}
		case *transformation.FilterTransformations:
			t := proto.Clone(filter).(*transformation.FilterTransformations)
			t.Stage = 0
			if isEmpty(t) {
				continue
			}
			o.Errors = append(o.Errors, fmt.Errorf("unhandled filter config %T", filter))
		case *aws.AWSLambdaConfig:
			if o.Settings.Gloo.AwsOptions == nil {
				o.Settings.Gloo.AwsOptions = &v1.GlooOptions_AWSOptions{}
			}
			o.Settings.Gloo.AwsOptions.PropagateOriginalRouting = &wrapperspb.BoolValue{Value: filter.PropagateOriginalRouting}
			o.Settings.Gloo.AwsOptions.CredentialRefreshDelay = &durationpb.Duration{Seconds: filter.CredentialRefreshDelay.Seconds}
			// Convert the AWS Lambda credentials fetcher to the Gloo options credentials fetcher
			if enableDiscovery, ok := filter.CredentialsFetcher.(*aws.AWSLambdaConfig_UseDefaultCredentials); ok {
				o.Settings.Gloo.AwsOptions.CredentialsFetcher = &v1.GlooOptions_AWSOptions_EnableCredentialsDiscovey{
					EnableCredentialsDiscovey: enableDiscovery.UseDefaultCredentials.GetValue(),
				}
			} else if serviceAccount, ok := filter.CredentialsFetcher.(*aws.AWSLambdaConfig_ServiceAccountCredentials_); ok {
				o.Settings.Gloo.AwsOptions.CredentialsFetcher = &v1.GlooOptions_AWSOptions_ServiceAccountCredentials{
					ServiceAccountCredentials: protoRoundTrip[*glooaws.AWSLambdaConfig_ServiceAccountCredentials](serviceAccount.ServiceAccountCredentials),
				}
			}
			filter.PropagateOriginalRouting = false
			filter.CredentialRefreshDelay = nil
			filter.CredentialsFetcher = nil

			if !isEmpty(filter) {
				o.Errors = append(o.Errors, fmt.Errorf("unhandled filter config %T", filter))
			}

		default:
			o.Errors = append(o.Errors, fmt.Errorf("unhandled filter config %T", filter))
		}
	}
	return ret
}
