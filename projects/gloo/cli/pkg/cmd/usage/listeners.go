package usage

import (
	gatewayv1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
)

func (u *UsageStats) processGlooGateways() error {

	for _, gateway := range u.GlooEdgeConfigs.GlooGateways() {
		spec := gateway.Gateway.Spec
		if spec.UseProxyProto != nil && spec.UseProxyProto.Value == true {
			u.AddUsageStat(&UsageStat{
				Type: PROXY_PROTOCOL,
				Metadata: UsageMetadata{
					Name:     gateway.Name,
					Kind:     "GlooGateway",
					Category: listenerCatagory,
					API:      GlooEdgeAPI,
				},
			})
		}
		if spec.GetTcpGateway() != nil {
			u.processTCPGateway(spec.GetTcpGateway(), gateway.GetName())
		}
		if spec.GetOptions() != nil {
			u.processListenerOptions(spec.GetOptions(), gateway.GetName(), "Gateway", GlooEdgeAPI)
		}
		if spec.GetRouteOptions() != nil {
			u.processRouteConfigurationOptions(spec.GetRouteOptions(), gateway.GetName(), "Gateway")
		}
		if spec.GetHttpGateway() != nil {
			u.processHTTPGateway(spec.GetHttpGateway(), gateway.GetName())
		}
		if spec.GetHybridGateway() != nil {
			u.AddUsageStat(&UsageStat{
				Type: HYBRID_GATEWAY,
				Metadata: UsageMetadata{
					Name:     gateway.Name,
					Kind:     "GlooGateway",
					Category: listenerCatagory,
					API:      GlooEdgeAPI,
				},
			})
		}
	}

	return nil
}
func (u *UsageStats) processRouteConfigurationOptions(options *gloov1.RouteConfigurationOptions, parentName, parentKind string) {
	defaultMetadata := UsageMetadata{
		Kind:     parentKind,
		Name:     parentName,
		Category: listenerCatagory,
		API:      GlooEdgeAPI,
	}
	if options.GetMostSpecificHeaderMutationsWins() != nil && options.GetMostSpecificHeaderMutationsWins().Value == true {
		u.AddUsageStat(&UsageStat{
			Type:     MOST_SPECIAL_HEADER_MUTATIONS_WINS,
			Metadata: defaultMetadata,
		})
	}
	if options.GetMaxDirectResponseBodySizeBytes() != nil && options.GetMaxDirectResponseBodySizeBytes().Value > 0 {
		u.AddUsageStat(&UsageStat{
			Type:     MAX_DIRECT_RESPONSE_BODY_SIZE,
			Metadata: defaultMetadata,
		})
	}
}
func (u *UsageStats) processListenerOptions(options *gloov1.ListenerOptions, parentName, parentKind string, parentAPI API) {

	defaultMetadata := UsageMetadata{
		Kind:     parentKind,
		Name:     parentName,
		Category: listenerCatagory,
		API:      parentAPI,
	}
	if options.SocketOptions != nil {
		u.AddUsageStat(&UsageStat{
			Type:     SOCKET_OPTIONS,
			Metadata: defaultMetadata,
		})
	}
	if options.AccessLoggingService != nil {
		u.AddUsageStat(&UsageStat{
			Type:     ACCESS_LOGGING,
			Metadata: defaultMetadata,
		})
	}
	if options.ConnectionBalanceConfig != nil {
		u.AddUsageStat(&UsageStat{
			Type:     CONNECTION_BALANCING,
			Metadata: defaultMetadata,
		})
	}
	if options.ListenerAccessLoggingService != nil {
		u.AddUsageStat(&UsageStat{
			Type:     EARLY_ACCESS_LOGGING,
			Metadata: defaultMetadata,
		})
	}
	if options.ProxyProtocol != nil {
		u.AddUsageStat(&UsageStat{
			Type:     PROXY_PROTOCOL,
			Metadata: defaultMetadata,
		})
	}

}
func (u *UsageStats) processTCPGateway(tcpGateway *gatewayv1.TcpGateway, gatewayName string) {
	defaultMetadata := UsageMetadata{
		Kind:     "Gateway",
		Name:     gatewayName,
		Category: listenerCatagory,
		API:      GlooEdgeAPI,
	}
	if tcpGateway != nil {
		for _, host := range tcpGateway.GetTcpHosts() {
			if host.SslConfig != nil {
				u.AddUsageStat(&UsageStat{
					Type:     TLS_ROUTING,
					Metadata: defaultMetadata,
				})
			} else {
				u.AddUsageStat(&UsageStat{
					Type:     TCP_ROUTING,
					Metadata: defaultMetadata,
				})
			}
		}
		if tcpGateway.GetOptions().LocalRatelimit != nil {
			u.AddUsageStat(&UsageStat{
				Type:     LOCAL_RATE_LIMITING,
				Metadata: defaultMetadata,
			})
		}
		if tcpGateway.GetOptions().ConnectionLimit != nil {
			u.AddUsageStat(&UsageStat{
				Type:     CONNECTION_LIMIT,
				Metadata: defaultMetadata,
			})
		}
	}
}

func (u *UsageStats) processHTTPGateway(gateway *gatewayv1.HttpGateway, parentName string) {

	if gateway.GetOptions() != nil {
		u.ProcessHTTPListenerOptions(gateway.GetOptions(), parentName, "HTTPGateway", GlooEdgeAPI)
	}

}

func (u *UsageStats) ProcessHTTPListenerOptions(options *gloov1.HttpListenerOptions, parentName string, parentKind string, parentAPI API) {
	defaultMetadata := UsageMetadata{
		Kind:     parentKind,
		Name:     parentName,
		Category: listenerCatagory,
		API:      parentAPI,
	}
	if options.GetHttpConnectionManagerSettings() != nil {
		u.AddUsageStat(&UsageStat{
			Type:     HTTP_CONNECTION_MANAGEMENT,
			Metadata: defaultMetadata,
		})
	}
	if options.GetHealthCheck() != nil {
		u.AddUsageStat(&UsageStat{
			Type:     INBOUND_HEALTH_CHECK,
			Metadata: defaultMetadata,
		})
	}
	if options.GetWaf() != nil {
		u.AddUsageStat(&UsageStat{
			Type:     WEB_APPLICATION_FIREWALL,
			Metadata: defaultMetadata,
		})
	}
	if options.GetDlp() != nil {
		u.AddUsageStat(&UsageStat{
			Type:     DATA_LOSS_PREVENTION,
			Metadata: defaultMetadata,
		})
	}
	if options.GetWasm() != nil {
		u.AddUsageStat(&UsageStat{
			Type:     WASM,
			Metadata: defaultMetadata,
		})
	}
	if options.GetCaching() != nil {
		u.AddUsageStat(&UsageStat{
			Type:     RESPONSE_CACHING,
			Metadata: defaultMetadata,
		})
	}
	if options.GetGzip() != nil {
		u.AddUsageStat(&UsageStat{
			Type:     GZIP_COMPRESSION,
			Metadata: defaultMetadata,
		})
	}
	if options.GetExtProc() != nil {
		u.AddUsageStat(&UsageStat{
			Type:     EXTERNAL_PROCESSING,
			Metadata: defaultMetadata,
		})
	}
	if options.GetBuffer() != nil {
		u.AddUsageStat(&UsageStat{
			Type:     BUFFER_LIMIT,
			Metadata: defaultMetadata,
		})
	}
	if options.GetCsrf() != nil {
		u.AddUsageStat(&UsageStat{
			Type:     CSRF,
			Metadata: defaultMetadata,
		})
	}
	if options.GetGrpcJsonTranscoder() != nil {
		u.AddUsageStat(&UsageStat{
			Type:     GRPC_JSON_TRANSCODING,
			Metadata: defaultMetadata,
		})
	}
	if options.GetDynamicForwardProxy() != nil {
		u.AddUsageStat(&UsageStat{
			Type:     DYNAMIC_FORWAD_PROXY,
			Metadata: defaultMetadata,
		})
	}
	if options.GetConnectionLimit() != nil {
		u.AddUsageStat(&UsageStat{
			Type:     CONNECTION_LIMIT,
			Metadata: defaultMetadata,
		})
	}
	if options.GetNetworkLocalRatelimit() != nil {
		u.AddUsageStat(&UsageStat{
			Type:     LOCAL_RATE_LIMITING,
			Metadata: defaultMetadata,
		})
	}
	if options.GetHttpLocalRatelimit() != nil {
		u.AddUsageStat(&UsageStat{
			Type:     LOCAL_RATE_LIMITING,
			Metadata: defaultMetadata,
		})
	}
	if options.GetTap() != nil {
		u.AddUsageStat(&UsageStat{
			Type:     TAP_FILTER,
			Metadata: defaultMetadata,
		})
	}
	if options.GetStatefulSession() != nil {
		u.AddUsageStat(&UsageStat{
			Type:     STATEFUL_SESSION,
			Metadata: defaultMetadata,
		})
	}
	if options.GetHeaderValidationSettings() != nil {
		u.AddUsageStat(&UsageStat{
			Type:     HEADER_VALIDATION,
			Metadata: defaultMetadata,
		})
	}

}
