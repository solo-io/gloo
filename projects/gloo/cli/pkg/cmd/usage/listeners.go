package usage

import (
	"time"

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

func (f *FeatureCalculator) processSettings(proxyNames []string) error {
	defaultMetadata := UsageMetadata{
		ProxyNames: proxyNames,
		Category:   settingsCategory,
		API:        GlooEdgeAPI,
	}
	for _, settings := range f.Configs.Settings() {
		defaultMetadata.Name = settings.Name
		defaultMetadata.Namespace = settings.Namespace
		if settings.Spec.GetDiscoveryNamespace() != "" {
			f.AddUsageStat(&UsageStat{
				Type:     DISCOVERY_NAMESPACE,
				Metadata: defaultMetadata,
			})
		}
		if settings.Spec.GetWatchNamespaces() != nil && len(settings.Spec.GetWatchNamespaces()) > 0 {
			f.AddUsageStat(&UsageStat{
				Type:     WATCH_NAMESPACES,
				Metadata: defaultMetadata,
			})
		}
		if settings.Spec.GetKubernetesConfigSource() != nil {
			f.AddUsageStat(&UsageStat{
				Type:     KUBERNETES_CONFIG_SOURCE,
				Metadata: defaultMetadata,
			})
		}
		if settings.Spec.GetDirectoryConfigSource() != nil {
			f.AddUsageStat(&UsageStat{
				Type:     DIRECTORY_CONFIG_SOURCE,
				Metadata: defaultMetadata,
			})
		}
		if settings.Spec.GetConsulKvSource() != nil {
			f.AddUsageStat(&UsageStat{
				Type:     CONSUL_KV_SOURCE,
				Metadata: defaultMetadata,
			})
		}
		if settings.Spec.GetKubernetesSecretSource() != nil {
			f.AddUsageStat(&UsageStat{
				Type:     KUBERNETES_SECRET_SOURCE,
				Metadata: defaultMetadata,
			})
		}
		if settings.Spec.GetVaultSecretSource() != nil {
			f.AddUsageStat(&UsageStat{
				Type:     VAULT_SECRET_SOURCE,
				Metadata: defaultMetadata,
			})
		}
		if settings.Spec.GetDirectorySecretSource() != nil {
			f.AddUsageStat(&UsageStat{
				Type:     DIRECTORY_SECRET_SOURCE,
				Metadata: defaultMetadata,
			})
		}
		if settings.Spec.GetSecretOptions() != nil && len(settings.Spec.GetSecretOptions().Sources) > 0 {
			for _, source := range settings.Spec.GetSecretOptions().Sources {
				switch source.Source.(type) {
				case *gloov1.Settings_SecretOptions_Source_Kubernetes:
					f.AddUsageStat(&UsageStat{
						Type: KUBERNETES_SECRET_SOURCE_OPTIONS,
					})
				case *gloov1.Settings_SecretOptions_Source_Vault:
					f.AddUsageStat(&UsageStat{
						Type: VAULT_SECRET_SOURCE_OPTIONS,
					})
				case *gloov1.Settings_SecretOptions_Source_Directory:
					f.AddUsageStat(&UsageStat{
						Type: DIRECTORY_SECRET_SOURCE_OPTIONS,
					})
				}
			}
		}
		if settings.Spec.GetKubernetesArtifactSource() != nil {
			f.AddUsageStat(&UsageStat{
				Type:     KUBERNETES_ARTIFACT_SOURCE,
				Metadata: defaultMetadata,
			})
		}
		if settings.Spec.GetDirectoryArtifactSource() != nil {
			f.AddUsageStat(&UsageStat{
				Type:     DIRECTORY_ARTIFACT_SOURCE,
				Metadata: defaultMetadata,
			})
		}
		if settings.Spec.GetConsulKvArtifactSource() != nil {
			f.AddUsageStat(&UsageStat{
				Type:     CONSUL_KV_ARTIFACT_SOURCE,
				Metadata: defaultMetadata,
			})
		}
		if settings.Spec.GetRefreshRate() != nil {
			f.AddUsageStat(&UsageStat{
				Type:     REFRESH_RATE,
				Metadata: defaultMetadata,
			})
		}
		if settings.Spec.GetDevMode() == true {
			f.AddUsageStat(&UsageStat{
				Type:     DEV_MODE,
				Metadata: defaultMetadata,
			})
		}
		if settings.Spec.GetLinkerd() == true {
			f.AddUsageStat(&UsageStat{
				Type:     LINKERD,
				Metadata: defaultMetadata,
			})
		}
		if settings.Spec.GetKnative() != nil {
			f.AddUsageStat(&UsageStat{
				Type:     KNATIVE,
				Metadata: defaultMetadata,
			})
		}
		if settings.Spec.GetDiscovery() != nil {
			if settings.Spec.GetDiscovery().GetFdsMode() != gloov1.Settings_DiscoveryOptions_BLACKLIST {
				if settings.Spec.GetDiscovery().GetFdsMode() != gloov1.Settings_DiscoveryOptions_WHITELIST {
					f.AddUsageStat(&UsageStat{
						Type:     FDS_WHITELIST,
						Metadata: defaultMetadata,
					})
				}
				if settings.Spec.GetDiscovery().GetFdsMode() != gloov1.Settings_DiscoveryOptions_DISABLED {
					f.AddUsageStat(&UsageStat{
						Type:     FDS_DISABLED,
						Metadata: defaultMetadata,
					})
				}
			}
			if settings.Spec.GetDiscovery().GetUdsOptions() != nil && settings.Spec.GetDiscovery().GetUdsOptions().GetEnabled() != nil && settings.Spec.GetDiscovery().GetUdsOptions().GetEnabled().Value == true {
				f.AddUsageStat(&UsageStat{
					Type:     UDS_DISCOVERY,
					Metadata: defaultMetadata,
				})
			}
			if settings.Spec.GetDiscovery().GetFdsOptions() != nil && settings.Spec.GetDiscovery().GetFdsOptions().GraphqlEnabled != nil && settings.Spec.GetDiscovery().GetFdsOptions().GraphqlEnabled.Value == true {
				f.AddUsageStat(&UsageStat{
					Type:     GRAPHQL_FDS,
					Metadata: defaultMetadata,
				})
			}
		}
		if settings.Spec.GetGloo() != nil {
			if settings.Spec.GetGloo().GetXdsBindAddr() != "" {
				f.AddUsageStat(&UsageStat{
					Type:     XDS_BIND_ADDR,
					Metadata: defaultMetadata,
				})
			}
			if settings.Spec.GetGloo().GetValidationBindAddr() != "" {
				f.AddUsageStat(&UsageStat{
					Type:     VALIDATION_BIND_ADDR,
					Metadata: defaultMetadata,
				})
			}
			if settings.Spec.GetGloo().GetCircuitBreakers() != nil {
				f.AddUsageStat(&UsageStat{
					Type:     GLOBAL_CIRCUIT_BREAKER,
					Metadata: defaultMetadata,
				})
			}
			if settings.Spec.GetGloo().GetEndpointsWarmingTimeout() != nil && settings.Spec.GetGloo().GetEndpointsWarmingTimeout().AsDuration() != (5*time.Minute) {
				f.AddUsageStat(&UsageStat{
					Type:     ENDPOINTS_WARMING_TIMEOUT,
					Metadata: defaultMetadata,
				})
			}
			if settings.Spec.GetGloo().GetAwsOptions() != nil {
				f.AddUsageStat(&UsageStat{
					Type:     AWS_OPTIONS,
					Metadata: defaultMetadata,
				})
			}
			if settings.Spec.GetGloo().GetInvalidConfigPolicy() != nil {
				f.AddUsageStat(&UsageStat{
					Type:     INVALID_CONFIG_POLICY,
					Metadata: defaultMetadata,
				})
			}
			if settings.Spec.GetGloo().GetDisableKubernetesDestinations() == true {
				f.AddUsageStat(&UsageStat{
					Type:     DISABLE_KUBERNETES_DESTINATIONS,
					Metadata: defaultMetadata,
				})
			}
			if settings.Spec.GetGloo().GetDisableGrpcWeb() != nil && settings.Spec.GetGloo().GetDisableGrpcWeb().Value == true {
				f.AddUsageStat(&UsageStat{
					Type:     DISABLE_GRPC_WEB,
					Metadata: defaultMetadata,
				})
			}
			if settings.Spec.GetGloo().GetDisableProxyGarbageCollection() != nil && settings.Spec.GetGloo().GetDisableProxyGarbageCollection().Value == true {
				f.AddUsageStat(&UsageStat{
					Type:     DISABLE_PROXY_GARBAGE_COLLECTION,
					Metadata: defaultMetadata,
				})
			}
			if settings.Spec.GetGloo().GetRegexMaxProgramSize() != nil && settings.Spec.GetGloo().GetRegexMaxProgramSize().Value != 100 {
				f.AddUsageStat(&UsageStat{
					Type:     REGEX_MAX_PROGRAM_SIZE,
					Metadata: defaultMetadata,
				})
			}
			if settings.Spec.GetGloo().GetEnableRestEds() != nil && settings.Spec.GetGloo().GetEnableRestEds().Value == true {
				f.AddUsageStat(&UsageStat{
					Type:     REST_XDS_BIND_ADDR,
					Metadata: defaultMetadata,
				})
			}
			if settings.Spec.GetGloo().GetFailoverUpstreamDnsPollingInterval() != nil && settings.Spec.GetGloo().GetFailoverUpstreamDnsPollingInterval().AsDuration() != 10*time.Second {
				f.AddUsageStat(&UsageStat{
					Type:     FAIL_OVER_UPSTREAM_DNS_POLLING_INTERVAL,
					Metadata: defaultMetadata,
				})
			}
			if settings.Spec.GetGloo().GetRemoveUnusedFilters() != nil && settings.Spec.GetGloo().GetRemoveUnusedFilters().Value == true {
				f.AddUsageStat(&UsageStat{
					Type:     REMOVE_UNUSED_FILTERS,
					Metadata: defaultMetadata,
				})
			}
			if settings.Spec.GetGloo().GetProxyDebugBindAddr() != "" {
				f.AddUsageStat(&UsageStat{
					Type:     PROXY_DEBUG_BIND_ADDR,
					Metadata: defaultMetadata,
				})
			}
			if settings.Spec.GetGloo().GetLogTransformationRequestResponseInfo() != nil && settings.Spec.GetGloo().GetLogTransformationRequestResponseInfo().Value == true {
				f.AddUsageStat(&UsageStat{
					Type:     LOG_TRANSFORMATION_REQUEST_RESPONSE_INFO,
					Metadata: defaultMetadata,
				})
			}
			if settings.Spec.GetGloo().GetTransformationEscapeCharacters() != nil && settings.Spec.GetGloo().GetTransformationEscapeCharacters().Value == true {
				f.AddUsageStat(&UsageStat{
					Type:     TRANSFORMATION_ESCAPE_CHARACTERS,
					Metadata: defaultMetadata,
				})
			}
			if settings.Spec.GetGloo().GetIstioOptions() != nil {
				f.AddUsageStat(&UsageStat{
					Type:     ISTIO_OPTIONS,
					Metadata: defaultMetadata,
				})
			}
		} // end of gloo

		if settings.Spec.GetGateway() != nil {
			if settings.Spec.GetGateway().GetValidation() != nil {
				f.AddUsageStat(&UsageStat{
					Type:     VALIDATION_OPTIONS,
					Metadata: defaultMetadata,
				})
			}
			if settings.Spec.GetGateway().GetReadGatewaysFromAllNamespaces() != true {
				f.AddUsageStat(&UsageStat{
					Type:     READ_GATEWAYS_FROM_ALL_NAMESPACES,
					Metadata: defaultMetadata,
				})
			}
			if settings.Spec.GetGateway().GetCompressedProxySpec() == true {
				f.AddUsageStat(&UsageStat{
					Type:     COMPRESSED_PROXY_SPEC,
					Metadata: defaultMetadata,
				})
			}
			if settings.Spec.GetGateway().GetVirtualServiceOptions() != nil && settings.Spec.GetGateway().GetVirtualServiceOptions().OneWayTls != nil && settings.Spec.GetGateway().GetVirtualServiceOptions().OneWayTls.Value == true {
				f.AddUsageStat(&UsageStat{
					Type:     GLOBAL_ONE_WAY_TLS,
					Metadata: defaultMetadata,
				})
			}
			if settings.Spec.GetGateway().GetPersistProxySpec() != nil && settings.Spec.GetGateway().GetPersistProxySpec().Value == true {
				f.AddUsageStat(&UsageStat{
					Type:     PERSIST_PROXY_SPEC,
					Metadata: defaultMetadata,
				})
			}
			if settings.Spec.GetGateway().GetEnableGatewayController() != nil && settings.Spec.GetGateway().GetEnableGatewayController().Value == false {
				f.AddUsageStat(&UsageStat{
					Type:     DISABLE_GATEWAY_CONTROLLER,
					Metadata: defaultMetadata,
				})
			}
			if settings.Spec.GetGateway().GetIsolateVirtualHostsBySslConfig() != nil && settings.Spec.GetGateway().GetIsolateVirtualHostsBySslConfig().Value == true {
				f.AddUsageStat(&UsageStat{
					Type:     ISOLATE_VIRTUAL_HOSTS_BY_SSL_CONFIG,
					Metadata: defaultMetadata,
				})
			}
			if settings.Spec.GetGateway().GetTranslateEmptyGateways() != nil && settings.Spec.GetGateway().GetTranslateEmptyGateways().Value == true {
				f.AddUsageStat(&UsageStat{
					Type:     TRANSLATE_EMPTY_GATEWAYS,
					Metadata: defaultMetadata,
				})
			}
			if settings.Spec.GetConsul() != nil {
				f.AddUsageStat(&UsageStat{
					Type:     CONSUL_CONFIGURATION,
					Metadata: defaultMetadata,
				})
			}
			if settings.Spec.GetConsulDiscovery() != nil {
				f.AddUsageStat(&UsageStat{
					Type:     CONSUL_DISCOVERY,
					Metadata: defaultMetadata,
				})
			}
			if settings.Spec.GetKubernetes() != nil && settings.Spec.GetKubernetes().GetRateLimits() != nil {
				f.AddUsageStat(&UsageStat{
					Type:     KUBERNETES_API_RATE_LIMITS,
					Metadata: defaultMetadata,
				})
			}
			if settings.Spec.GetExtensions() != nil && len(settings.Spec.GetExtensions().GetConfigs()) > 0 {
				f.AddUsageStat(&UsageStat{
					Type:     GLOBAL_EXTENSIONS,
					Metadata: defaultMetadata,
				})
			}
			if settings.Spec.GetRbac() != nil && settings.Spec.GetRbac().RequireRbac == true {
				f.AddUsageStat(&UsageStat{
					Type:     GLOBAL_RBAC,
					Metadata: defaultMetadata,
				})
			}
			if settings.Spec.GetExtauth() != nil {
				f.AddUsageStat(&UsageStat{
					Type:     GLOBAL_EXTAUTH_SETTINGS,
					Metadata: defaultMetadata,
				})
			}
			if settings.Spec.GetNamedExtauth() != nil && len(settings.Spec.GetNamedExtauth()) > 0 {
				f.AddUsageStat(&UsageStat{
					Type:     GLOBAL_NAMED_EXTAUTH_SETTINGS,
					Metadata: defaultMetadata,
				})
			}
			if settings.Spec.GetCachingServer() != nil {
				f.AddUsageStat(&UsageStat{
					Type:     GLOBAL_CACHING_SERVER_SETTINGS,
					Metadata: defaultMetadata,
				})
			}
			if settings.Spec.GetMetadata() != nil {
				f.AddUsageStat(&UsageStat{
					Type:     GLOBAL_METADATA_SETTINGS,
					Metadata: defaultMetadata,
				})
			}
			if settings.Spec.GetObservabilityOptions() != nil {
				f.AddUsageStat(&UsageStat{
					Type:     GLOBAL_OBSERVABILITY_SETTINGS,
					Metadata: defaultMetadata,
				})
			}
			if settings.Spec.GetUpstreamOptions() != nil {
				f.AddUsageStat(&UsageStat{
					Type:     GLOBAL_UPSTREAM_SETTINGS,
					Metadata: defaultMetadata,
				})
			}
			if settings.Spec.GetGraphqlOptions() != nil {
				f.AddUsageStat(&UsageStat{
					Type:     GLOBAL_GRAPHQL_SETTINGS,
					Metadata: defaultMetadata,
				})
			}
			if settings.Spec.GetExtProc() != nil {
				f.AddUsageStat(&UsageStat{
					Type:     GLOBAL_EXTPROC_SETTINGS,
					Metadata: defaultMetadata,
				})
			}
			if settings.Spec.GetWatchNamespaceSelectors() != nil {
				f.AddUsageStat(&UsageStat{
					Type:     GLOBAL_WATCH_NAMESPACE_SELECTORS,
					Metadata: defaultMetadata,
				})
			}
		}
	} // end of settings
	return nil
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
