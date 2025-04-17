package usage

import (
	gatewayv1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/snapshot"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
)

type FeatureCalculator struct {
	Configs  *snapshot.Instance
	Features map[API][]*UsageStat `json:"glooFeatureUsage"`
}

type UsageMetadata struct {
	ProxyNames      []string `json:"proxyNames" yaml:"proxyNames,flow"`
	Name            string   `json:"name" yaml:"name"`
	Namespace       string   `json:"namespace" yaml:"namespace"`
	ParentName      string   `json:"parentName" yaml:"parentName"`
	ParentNamespace string   `json:"parentNamespace" yaml:"parentNamespace"`
	Kind            string   `json:"kind" yaml:"kind"`
	Category        Category `json:"category" yaml:"category,double"`
	API             API      `json:"api" yaml:"api,double"`
}

type UsageStat struct {
	Type     FeatureType   `json:"type" yaml:"type,double"`
	Metadata UsageMetadata `json:"metadata" yaml:"metadata"`
}

func (f *FeatureCalculator) AddUsageStat(stat *UsageStat) {
	if stat.Metadata.API == "" {
		stat.Metadata.API = GlooEdgeAPI
	}
	if f.Features == nil {
		f.Features = make(map[API][]*UsageStat)
	}
	if f.Features[stat.Metadata.API] == nil {
		f.Features[stat.Metadata.API] = []*UsageStat{}
	}
	f.Features[stat.Metadata.API] = append(f.Features[stat.Metadata.API], stat)
}

func (f *FeatureCalculator) processGlooGateways() error {

	for _, gateway := range f.Configs.GlooGateways() {
		spec := gateway.Gateway.Spec
		proxyNames := gateway.Spec.ProxyNames
		if len(proxyNames) == 0 {
			proxyNames = append(proxyNames, "gateway-proxy")
		}
		if spec.UseProxyProto != nil && spec.UseProxyProto.Value == true {
			f.AddUsageStat(&UsageStat{
				Type: PROXY_PROTOCOL,
				Metadata: UsageMetadata{
					ProxyNames: proxyNames,
					Name:       gateway.Name,
					Kind:       "GlooGateway",
					Category:   listenerCatagory,
					API:        GlooEdgeAPI,
				},
			})
		}

		if spec.GetTcpGateway() != nil {
			f.processTCPGateway(spec.GetTcpGateway(), gateway.GetName(), proxyNames)
		}
		if spec.GetOptions() != nil {
			f.processListenerOptions(spec.GetOptions(), gateway.GetName(), "Gateway", GlooEdgeAPI, proxyNames)
		}
		if spec.GetRouteOptions() != nil {
			f.processRouteConfigurationOptions(spec.GetRouteOptions(), gateway.GetName(), "Gateway", proxyNames)
		}
		if spec.GetHttpGateway() != nil {
			err := f.processHTTPGateway(gateway, proxyNames)
			if err != nil {
				return err
			}
		}
		if spec.GetHybridGateway() != nil {
			f.AddUsageStat(&UsageStat{
				Type: HYBRID_GATEWAY,
				Metadata: UsageMetadata{
					ProxyNames: proxyNames,
					Name:       gateway.Name,
					Kind:       "GlooGateway",
					Category:   listenerCatagory,
					API:        GlooEdgeAPI,
				},
			})
		}
	}

	return nil
}
func (f *FeatureCalculator) processRouteConfigurationOptions(options *gloov1.RouteConfigurationOptions, parentName string, parentKind string, proxyNames []string) {
	defaultMetadata := UsageMetadata{
		ProxyNames: proxyNames,
		Kind:       parentKind,
		Name:       parentName,
		Category:   listenerCatagory,
		API:        GlooEdgeAPI,
	}
	if options.GetMostSpecificHeaderMutationsWins() != nil && options.GetMostSpecificHeaderMutationsWins().Value == true {
		f.AddUsageStat(&UsageStat{
			Type:     MOST_SPECIAL_HEADER_MUTATIONS_WINS,
			Metadata: defaultMetadata,
		})
	}
	if options.GetMaxDirectResponseBodySizeBytes() != nil && options.GetMaxDirectResponseBodySizeBytes().Value > 0 {
		f.AddUsageStat(&UsageStat{
			Type:     MAX_DIRECT_RESPONSE_BODY_SIZE,
			Metadata: defaultMetadata,
		})
	}
}
func (f *FeatureCalculator) processListenerOptions(options *gloov1.ListenerOptions, parentName, parentKind string, parentAPI API, proxyNames []string) {

	defaultMetadata := UsageMetadata{
		ProxyNames: proxyNames,
		Kind:       parentKind,
		Name:       parentName,
		Category:   listenerCatagory,
		API:        parentAPI,
	}
	if options.SocketOptions != nil {
		f.AddUsageStat(&UsageStat{
			Type:     SOCKET_OPTIONS,
			Metadata: defaultMetadata,
		})
	}
	if options.AccessLoggingService != nil {
		f.AddUsageStat(&UsageStat{
			Type:     ACCESS_LOGGING,
			Metadata: defaultMetadata,
		})
	}
	if options.ConnectionBalanceConfig != nil {
		f.AddUsageStat(&UsageStat{
			Type:     CONNECTION_BALANCING,
			Metadata: defaultMetadata,
		})
	}
	if options.ListenerAccessLoggingService != nil {
		f.AddUsageStat(&UsageStat{
			Type:     EARLY_ACCESS_LOGGING,
			Metadata: defaultMetadata,
		})
	}
	if options.ProxyProtocol != nil {
		f.AddUsageStat(&UsageStat{
			Type:     PROXY_PROTOCOL,
			Metadata: defaultMetadata,
		})
	}
}
func (f *FeatureCalculator) processTCPGateway(tcpGateway *gatewayv1.TcpGateway, gatewayName string, proxyNames []string) {
	defaultMetadata := UsageMetadata{
		ProxyNames: proxyNames,
		Kind:       "Gateway",
		Name:       gatewayName,
		Category:   listenerCatagory,
		API:        GlooEdgeAPI,
	}
	if tcpGateway != nil {
		for _, host := range tcpGateway.GetTcpHosts() {
			if host.SslConfig != nil {
				f.AddUsageStat(&UsageStat{
					Type:     TLS_ROUTING,
					Metadata: defaultMetadata,
				})
			} else {
				f.AddUsageStat(&UsageStat{
					Type:     TCP_ROUTING,
					Metadata: defaultMetadata,
				})
			}
		}
		if tcpGateway.GetOptions().LocalRatelimit != nil {
			f.AddUsageStat(&UsageStat{
				Type:     LOCAL_RATE_LIMITING,
				Metadata: defaultMetadata,
			})
		}
		if tcpGateway.GetOptions().ConnectionLimit != nil {
			f.AddUsageStat(&UsageStat{
				Type:     CONNECTION_LIMIT,
				Metadata: defaultMetadata,
			})
		}
	}
}

// The most common type, it selects virtual services
func (f *FeatureCalculator) processHTTPGateway(wrapper *snapshot.GlooGatewayWrapper, proxyNames []string) error {

	if wrapper.Spec.GetHttpGateway().GetOptions() != nil {
		f.ProcessHTTPListenerOptions(wrapper.Spec.GetHttpGateway().GetOptions(), wrapper.GetName(), "HTTPGateway", GlooEdgeAPI, proxyNames)
	}

	// get the virtual services associated with the gateway
	services, err := f.Configs.GlooGatewayVirtualServices(wrapper)
	if err != nil {
		return err
	}

	for _, vs := range services {
		err := f.processVirtualService(vs.VirtualService, proxyNames)
		if err != nil {
			return err
		}
	}
	return nil
}

func (f *FeatureCalculator) ProcessHTTPListenerOptions(options *gloov1.HttpListenerOptions, parentName string, parentKind string, parentAPI API, proxyNames []string) {
	defaultMetadata := UsageMetadata{
		ProxyNames: proxyNames,
		Kind:       parentKind,
		Name:       parentName,
		Category:   listenerCatagory,
		API:        parentAPI,
	}
	if options.GetHttpConnectionManagerSettings() != nil {
		f.AddUsageStat(&UsageStat{
			Type:     HTTP_CONNECTION_MANAGEMENT,
			Metadata: defaultMetadata,
		})
	}
	if options.GetHealthCheck() != nil {
		f.AddUsageStat(&UsageStat{
			Type:     INBOUND_HEALTH_CHECK,
			Metadata: defaultMetadata,
		})
	}
	if options.GetWaf() != nil {
		f.AddUsageStat(&UsageStat{
			Type:     WEB_APPLICATION_FIREWALL,
			Metadata: defaultMetadata,
		})
	}
	if options.GetDlp() != nil {
		f.AddUsageStat(&UsageStat{
			Type:     DATA_LOSS_PREVENTION,
			Metadata: defaultMetadata,
		})
	}
	if options.GetWasm() != nil {
		f.AddUsageStat(&UsageStat{
			Type:     WASM,
			Metadata: defaultMetadata,
		})
	}
	if options.GetCaching() != nil {
		f.AddUsageStat(&UsageStat{
			Type:     RESPONSE_CACHING,
			Metadata: defaultMetadata,
		})
	}
	if options.GetGzip() != nil {
		f.AddUsageStat(&UsageStat{
			Type:     GZIP_COMPRESSION,
			Metadata: defaultMetadata,
		})
	}
	if options.GetExtProc() != nil {
		f.AddUsageStat(&UsageStat{
			Type:     EXTERNAL_PROCESSING,
			Metadata: defaultMetadata,
		})
	}
	if options.GetBuffer() != nil {
		f.AddUsageStat(&UsageStat{
			Type:     BUFFER_LIMIT,
			Metadata: defaultMetadata,
		})
	}
	if options.GetCsrf() != nil {
		f.AddUsageStat(&UsageStat{
			Type:     CSRF,
			Metadata: defaultMetadata,
		})
	}
	if options.GetGrpcJsonTranscoder() != nil {
		f.AddUsageStat(&UsageStat{
			Type:     GRPC_JSON_TRANSCODING,
			Metadata: defaultMetadata,
		})
	}
	if options.GetDynamicForwardProxy() != nil {
		f.AddUsageStat(&UsageStat{
			Type:     DYNAMIC_FORWAD_PROXY,
			Metadata: defaultMetadata,
		})
	}
	if options.GetConnectionLimit() != nil {
		f.AddUsageStat(&UsageStat{
			Type:     CONNECTION_LIMIT,
			Metadata: defaultMetadata,
		})
	}
	if options.GetNetworkLocalRatelimit() != nil {
		f.AddUsageStat(&UsageStat{
			Type:     LOCAL_RATE_LIMITING,
			Metadata: defaultMetadata,
		})
	}
	if options.GetHttpLocalRatelimit() != nil {
		f.AddUsageStat(&UsageStat{
			Type:     LOCAL_RATE_LIMITING,
			Metadata: defaultMetadata,
		})
	}
	if options.GetTap() != nil {
		f.AddUsageStat(&UsageStat{
			Type:     TAP_FILTER,
			Metadata: defaultMetadata,
		})
	}
	if options.GetStatefulSession() != nil {
		f.AddUsageStat(&UsageStat{
			Type:     STATEFUL_SESSION,
			Metadata: defaultMetadata,
		})
	}
	if options.GetHeaderValidationSettings() != nil {
		f.AddUsageStat(&UsageStat{
			Type:     HEADER_VALIDATION,
			Metadata: defaultMetadata,
		})
	}

}
