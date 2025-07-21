package convert

import (
	corev3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	kgateway "github.com/kgateway-dev/kgateway/v2/api/v1alpha1"
	gloogwv1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/snapshot"
	ext_core_v3 "github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/config/core/v3"
	glookube "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/kube/apis/gloo.solo.io/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/aws"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/protocol"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/ssl"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/utils/ptr"
	gwv1 "sigs.k8s.io/gateway-api/apis/v1"
)

// TODO(nick): does aws backend support awsec2 upstream?
func (g *GatewayAPIOutput) convertUpstreamToBackend(upstream *snapshot.UpstreamWrapper) {
	// Add all existing upstreams except for kube services which will be referenced directly
	if upstream.Spec.GetKube() != nil {
		// do nothing, let it continue in case there were other policies attached to the kube that we can warn about
	}
	if upstream.Spec.GetAi() != nil {
		g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, upstream, "gcp AI is not supported")
	}

	// determine if we need to create an upstream policy which will apply to the upstream
	// TODO we might need to see if the upstream is a kube service and then apply it to that?
	g.convertUpstreamPolicy(upstream)

	if upstream.Spec.GetAws() != nil {
		if len(upstream.Spec.GetAws().GetLambdaFunctions()) > 0 {
			backend := g.convertAWSBackend(upstream, nil)
			g.gatewayAPICache.AddBackend(backend)
		} else {
			//TODO (nick): we create multiple backends here but we need to fix naming and backendPolicies
			for _, lambda := range upstream.Spec.GetAws().GetLambdaFunctions() {
				backend := g.convertAWSBackend(upstream, lambda)
				g.gatewayAPICache.AddBackend(backend)
			}
		}
	}
	if upstream.Spec.GetStatic() != nil {

		backend := &snapshot.BackendWrapper{
			Backend: &kgateway.Backend{
				TypeMeta: metav1.TypeMeta{
					Kind:       "Backend",
					APIVersion: kgateway.SchemeGroupVersion.String(),
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      upstream.GetName(),
					Namespace: upstream.GetNamespace(),
				},
				Spec: kgateway.BackendSpec{
					Type: kgateway.BackendTypeStatic,
					AI:   nil, // existing
					Aws:  nil, // existing
					Static: &kgateway.StaticBackend{
						Hosts: []kgateway.Host{},
					}, // existing
					DynamicForwardProxy: nil,
				},
			},
		}
		for _, hosts := range upstream.Spec.GetStatic().GetHosts() {

			if hosts.GetHealthCheckConfig() != nil {
				g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, upstream, "static upstream healthCheckConfig is not supported")
			}
			if hosts.GetLoadBalancingWeight() != nil {
				g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, upstream, "static upstream loadBalancingWeight is not supported")
			}
			if hosts.GetSniAddr() != "" {
				g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, upstream, "static upstream sni is not supported")
			}

			if hosts.GetAddr() != "" {
				host := kgateway.Host{
					Host: hosts.GetAddr(),
				}
				if hosts.GetPort() != 0 {
					host.Port = gwv1.PortNumber(hosts.GetPort())
				}
				g.convertUpstreamPolicy(upstream)
				backend.Spec.Static.Hosts = append(backend.Spec.Static.Hosts, host)
			}
		}
		g.gatewayAPICache.AddBackend(backend)
	}
	if upstream.Spec.GetAwsEc2() != nil {
		g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, upstream, "awsec2 upstream is not supported")
	}
	if upstream.Spec.GetConsul() != nil {
		g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, upstream, "consul upstream is not supported")
	}
	if upstream.Spec.GetAzure() != nil {
		g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, upstream, "azure upstream is not supported")
	}
	if upstream.Spec.GetGcp() != nil {
		g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, upstream, "gcp upstream is not supported")
	}

}
func (g *GatewayAPIOutput) convertUpstreamPolicy(upstream *snapshot.UpstreamWrapper) {
	configExists := false

	targetRef := kgateway.LocalPolicyTargetReference{
		Group: kgateway.GroupName,
		Kind:  "Backend",
		Name:  gwv1.ObjectName(upstream.GetName()),
	}

	if upstream.Spec.GetKube() != nil {
		targetRef = kgateway.LocalPolicyTargetReference{
			Group: "v1",
			Kind:  "Service",
			Name:  gwv1.ObjectName(upstream.GetName()),
		}
	}

	backendConfig := &kgateway.BackendConfigPolicy{
		TypeMeta: metav1.TypeMeta{
			Kind:       "BackendConfigPolicy",
			APIVersion: kgateway.GroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      upstream.GetName(),
			Namespace: upstream.GetNamespace(),
		},
		Spec: kgateway.BackendConfigPolicySpec{
			TargetRefs: []kgateway.LocalPolicyTargetReference{
				targetRef,
			},
			ConnectTimeout:                nil, // existing
			PerConnectionBufferLimitBytes: nil, // existing
			TCPKeepalive:                  nil, // existing
			CommonHttpProtocolOptions:     nil, // existing
			Http1ProtocolOptions:          nil, // existing
			Http2ProtocolOptions:          nil, // existing
			TLS:                           nil, // existing
			LoadBalancer:                  nil, // existing
			HealthCheck:                   nil, // existing
			TargetSelectors:               nil, //existing
		},
	}
	if upstream.Spec.GetStatic() != nil {
		if upstream.Spec.GetStatic().GetUseTls() != nil || upstream.Spec.GetStatic().GetUseTls().GetValue() {
			// TODO(nick): This currently does not do anything. Need support for TLS
			backendConfig.Spec.TLS = &kgateway.TLS{}
		}
		if upstream.Spec.GetStatic().GetAutoSniRewrite() != nil && upstream.Spec.GetStatic().GetAutoSniRewrite().GetValue() {
			g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, upstream, "static upstream autoSniRewrite is not supported")
		}
	}

	if upstream.Spec.GetSslConfig() != nil {
		tls := kgateway.TLS{
			SecretRef:            nil,                                   // existing
			TLSFiles:             nil,                                   // existing
			Sni:                  upstream.Spec.GetSslConfig().GetSni(), // existing
			VerifySubjectAltName: nil,                                   // existing
			Parameters:           nil,                                   // existing
			AlpnProtocols:        nil,                                   // existing
			AllowRenegotiation:   nil,                                   // existing
			OneWayTLS:            nil,                                   // existing
		}
		if upstream.Spec.GetSslConfig().GetSecretRef() != nil {
			if upstream.Spec.GetSslConfig().GetSecretRef().GetNamespace() != upstream.GetNamespace() {
				g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, upstream, "sslConfig.secretRef.namespace %s is not the same as the backendConfig's %s", upstream.GetNamespace(), upstream.Spec.GetSslConfig().GetSecretRef().GetNamespace())
			}
			tls.SecretRef = &corev1.LocalObjectReference{
				Name: upstream.Spec.GetSslConfig().GetSecretRef().GetName(),
			}
		}
		if upstream.Spec.GetSslConfig().GetSslFiles() != nil {
			tls.TLSFiles = &kgateway.TLSFiles{
				TLSCertificate: upstream.Spec.GetSslConfig().GetSslFiles().GetTlsCert(),
				TLSKey:         upstream.Spec.GetSslConfig().GetSslFiles().GetTlsKey(),
				RootCA:         upstream.Spec.GetSslConfig().GetSslFiles().GetRootCa(),
			}
			if upstream.Spec.GetSslConfig().GetSslFiles().GetOcspStaple() != "" {
				g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, upstream, "sslConfig.sslFiles.ocspStaple is not supported")
			}
		}

		if upstream.Spec.GetSslConfig().GetOneWayTls() != nil {
			tls.OneWayTLS = ptr.To(upstream.Spec.GetSslConfig().GetOneWayTls().GetValue())
		}
		if len(upstream.Spec.GetSslConfig().GetAlpnProtocols()) > 0 {
			tls.AlpnProtocols = upstream.Spec.GetSslConfig().GetAlpnProtocols()
		}
		if upstream.Spec.GetSslConfig().GetAllowRenegotiation() != nil {
			tls.AllowRenegotiation = ptr.To(upstream.Spec.GetSslConfig().GetAllowRenegotiation().GetValue())
		}
		if upstream.Spec.GetSslConfig().GetParameters() != nil {
			params := &kgateway.Parameters{
				TLSMinVersion: nil, // existing
				TLSMaxVersion: nil, // existing
				CipherSuites:  nil, // existing
				EcdhCurves:    nil, // existing
			}
			if len(upstream.Spec.GetSslConfig().GetParameters().GetEcdhCurves()) > 0 {
				params.EcdhCurves = upstream.Spec.GetSslConfig().GetParameters().GetEcdhCurves()
			}
			if len(upstream.Spec.GetSslConfig().GetParameters().GetCipherSuites()) > 0 {
				params.CipherSuites = upstream.Spec.GetSslConfig().GetParameters().GetCipherSuites()
			}
			switch upstream.Spec.GetSslConfig().GetParameters().GetMaximumProtocolVersion() {
			case ssl.SslParameters_TLS_AUTO:
				params.TLSMaxVersion = ptr.To(kgateway.TLSVersionAUTO)
			case ssl.SslParameters_TLSv1_0:
				params.TLSMaxVersion = ptr.To(kgateway.TLSVersion1_0)
			case ssl.SslParameters_TLSv1_1:
				params.TLSMaxVersion = ptr.To(kgateway.TLSVersion1_1)
			case ssl.SslParameters_TLSv1_2:
				params.TLSMaxVersion = ptr.To(kgateway.TLSVersion1_2)
			case ssl.SslParameters_TLSv1_3:
				params.TLSMaxVersion = ptr.To(kgateway.TLSVersion1_3)
			}
			switch upstream.Spec.GetSslConfig().GetParameters().GetMinimumProtocolVersion() {
			case ssl.SslParameters_TLS_AUTO:
				params.TLSMinVersion = ptr.To(kgateway.TLSVersionAUTO)
			case ssl.SslParameters_TLSv1_0:
				params.TLSMinVersion = ptr.To(kgateway.TLSVersion1_0)
			case ssl.SslParameters_TLSv1_1:
				params.TLSMinVersion = ptr.To(kgateway.TLSVersion1_1)
			case ssl.SslParameters_TLSv1_2:
				params.TLSMinVersion = ptr.To(kgateway.TLSVersion1_2)
			case ssl.SslParameters_TLSv1_3:
				params.TLSMinVersion = ptr.To(kgateway.TLSVersion1_3)
			}
			tls.Parameters = params
		}
		if upstream.Spec.GetSslConfig().GetSds() != nil {
			g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, upstream, "sslConfig.sds is not supported")
		}
		if upstream.Spec.GetSslConfig().GetOneWayTls() != nil {
			tls.OneWayTLS = ptr.To(upstream.Spec.GetSslConfig().GetOneWayTls().GetValue())
		}
		if upstream.Spec.GetSslConfig().GetVerifySubjectAltName() != nil {
			tls.VerifySubjectAltName = upstream.Spec.GetSslConfig().GetVerifySubjectAltName()
		}

		backendConfig.Spec.TLS = &tls
		configExists = true
	}

	if upstream.Spec.GetCircuitBreakers() != nil {
		g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, upstream, "circuitBreakers is not supported")
	}
	if upstream.Spec.GetConnectionConfig() != nil {
		if upstream.Spec.GetConnectionConfig().GetPerConnectionBufferLimitBytes() != nil {
			backendConfig.Spec.PerConnectionBufferLimitBytes = ptr.To(int(upstream.Spec.GetConnectionConfig().GetPerConnectionBufferLimitBytes().GetValue()))
		}
		if upstream.Spec.GetConnectionConfig().GetConnectTimeout() != nil {
			backendConfig.Spec.ConnectTimeout = ptr.To(metav1.Duration{Duration: upstream.Spec.GetConnectionConfig().GetConnectTimeout().AsDuration()})
		}
		if upstream.Spec.GetConnectionConfig().GetTcpKeepalive() != nil {
			keepAlive := &kgateway.TCPKeepalive{
				KeepAliveProbes:   ptr.To(int(upstream.Spec.GetConnectionConfig().GetTcpKeepalive().GetKeepaliveProbes())),
				KeepAliveTime:     nil, // existing
				KeepAliveInterval: nil, // existing
			}
			if upstream.Spec.GetConnectionConfig().GetTcpKeepalive().GetKeepaliveInterval() != nil {
				keepAlive.KeepAliveTime = ptr.To(metav1.Duration{Duration: upstream.Spec.GetConnectionConfig().GetTcpKeepalive().GetKeepaliveTime().AsDuration()})
			}
			if upstream.Spec.GetConnectionConfig().GetTcpKeepalive().GetKeepaliveInterval() != nil {
				keepAlive.KeepAliveInterval = ptr.To(metav1.Duration{Duration: upstream.Spec.GetConnectionConfig().GetTcpKeepalive().GetKeepaliveInterval().AsDuration()})
			}

			backendConfig.Spec.TCPKeepalive = keepAlive
			configExists = true
		}
		if upstream.Spec.GetConnectionConfig().GetCommonHttpProtocolOptions() != nil || upstream.Spec.GetConnectionConfig().GetMaxRequestsPerConnection() > 0 {

			options := &kgateway.CommonHttpProtocolOptions{
				IdleTimeout:              nil, // existing
				MaxHeadersCount:          nil, // existing
				MaxStreamDuration:        nil, // existing
				MaxRequestsPerConnection: nil, // existing
			}
			if upstream.Spec.GetConnectionConfig().GetCommonHttpProtocolOptions().GetIdleTimeout() != nil {
				options.IdleTimeout = ptr.To(metav1.Duration{Duration: upstream.Spec.GetConnectionConfig().GetCommonHttpProtocolOptions().GetIdleTimeout().AsDuration()})
			}
			if upstream.Spec.GetConnectionConfig().GetCommonHttpProtocolOptions().GetMaxStreamDuration() != nil {
				options.MaxStreamDuration = ptr.To(metav1.Duration{Duration: upstream.Spec.GetConnectionConfig().GetCommonHttpProtocolOptions().GetMaxStreamDuration().AsDuration()})
			}
			if upstream.Spec.GetConnectionConfig().GetCommonHttpProtocolOptions().GetMaxHeadersCount() > 0 {
				options.MaxHeadersCount = ptr.To(int(upstream.Spec.GetConnectionConfig().GetCommonHttpProtocolOptions().GetMaxHeadersCount()))
			}
			if upstream.Spec.GetConnectionConfig().GetCommonHttpProtocolOptions().GetHeadersWithUnderscoresAction() > 0 {
				g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, upstream, "commonHTTPProtocolOptions.headersWithUndercoresAction is not supported")
			}
			if upstream.Spec.GetConnectionConfig().GetMaxRequestsPerConnection() > 0 {
				options.MaxRequestsPerConnection = ptr.To(int(upstream.Spec.GetConnectionConfig().GetMaxRequestsPerConnection()))
			}
			backendConfig.Spec.CommonHttpProtocolOptions = options
			configExists = true
		}
		if upstream.Spec.GetConnectionConfig().GetHttp1ProtocolOptions() != nil {
			options := &kgateway.Http1ProtocolOptions{
				EnableTrailers:                          nil,
				HeaderFormat:                            nil,
				OverrideStreamErrorOnInvalidHttpMessage: nil,
			}
			if upstream.Spec.GetConnectionConfig().GetHttp1ProtocolOptions().GetOverrideStreamErrorOnInvalidHttpMessage() != nil {
				options.OverrideStreamErrorOnInvalidHttpMessage = ptr.To(upstream.Spec.GetConnectionConfig().GetHttp1ProtocolOptions().GetOverrideStreamErrorOnInvalidHttpMessage().GetValue())
			}
			if upstream.Spec.GetConnectionConfig().GetHttp1ProtocolOptions().GetEnableTrailers() {
				options.EnableTrailers = ptr.To(upstream.Spec.GetConnectionConfig().GetHttp1ProtocolOptions().GetEnableTrailers())
			}
			if upstream.Spec.GetConnectionConfig().GetHttp1ProtocolOptions().GetHeaderFormat() != nil {
				switch upstream.Spec.GetConnectionConfig().GetHttp1ProtocolOptions().GetHeaderFormat().(type) {
				case *protocol.Http1ProtocolOptions_ProperCaseHeaderKeyFormat:
					options.HeaderFormat = ptr.To(kgateway.ProperCaseHeaderKeyFormat)
				case *protocol.Http1ProtocolOptions_PreserveCaseHeaderKeyFormat:
					options.HeaderFormat = ptr.To(kgateway.PreserveCaseHeaderKeyFormat)
				}

				backendConfig.Spec.Http1ProtocolOptions = options
			}
			configExists = true
		}
		if upstream.Spec.GetHealthChecks() != nil {
			healthCheck := &kgateway.HealthCheck{
				Timeout:            nil, // existing
				Interval:           nil, // existing
				UnhealthyThreshold: nil, // existing
				HealthyThreshold:   nil, // existing
				Http:               nil, // existing
				Grpc:               nil, // existing
			}
			// going to just take the first health check
			if len(upstream.Spec.GetHealthChecks()) > 0 {
				g.AddErrorFromWrapper(ERROR_TYPE_IGNORED, upstream, "commonHealthCheck.healthChecks only using first health check")
				hc := upstream.Spec.GetHealthChecks()[0]
				if hc.GetTimeout() != nil {
					healthCheck.Timeout = ptr.To(metav1.Duration{Duration: hc.GetTimeout().AsDuration()})
				}
				if hc.GetInterval() != nil {
					healthCheck.Interval = ptr.To(metav1.Duration{Duration: hc.GetInterval().AsDuration()})
				}
				if hc.GetUnhealthyThreshold() != nil {
					healthCheck.UnhealthyThreshold = ptr.To(hc.GetUnhealthyThreshold().GetValue())
				}
				if hc.GetHealthyThreshold() != nil {
					healthCheck.HealthyThreshold = ptr.To(hc.GetHealthyThreshold().GetValue())
				}
				if hc.GetHttpHealthCheck() != nil {
					http := &kgateway.HealthCheckHttp{
						Host:   ptr.To(hc.GetHttpHealthCheck().GetHost()),
						Path:   hc.GetHttpHealthCheck().GetPath(), // existing
						Method: nil,                               // existing
					}
					if hc.GetHttpHealthCheck().GetUseHttp2() {
						g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, upstream, "healthCheck.httpHealthCheck.useHTTP2 is not supported")
					}
					if hc.GetHttpHealthCheck().GetServiceName() != "" {
						g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, upstream, "healthCheck.httpHealthCheck.serviceName is not supported")
					}
					if hc.GetHttpHealthCheck().GetExpectedStatuses() != nil {
						g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, upstream, "healthCheck.httpHealthCheck.expectedStatuses is not supported")
					}
					if hc.GetHttpHealthCheck().GetRequestHeadersToAdd() != nil {
						g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, upstream, "healthCheck.httpHealthCheck.requestHeadersToAdd is not supported")
					}
					if hc.GetHttpHealthCheck().GetRequestHeadersToRemove() != nil {
						g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, upstream, "healthCheck.httpHealthCheck.requestHeadersToRemove is not supported")
					}
					if hc.GetHttpHealthCheck().GetResponseAssertions() != nil {
						g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, upstream, "healthCheck.httpHealthCheck.responseAssertions is not supported")
					}
					switch hc.GetHttpHealthCheck().GetMethod() {
					case ext_core_v3.RequestMethod(corev3.RequestMethod_METHOD_UNSPECIFIED):
						http.Method = ptr.To("GET")
					case ext_core_v3.RequestMethod(corev3.RequestMethod_GET):
						http.Method = ptr.To("GET")
					case ext_core_v3.RequestMethod(corev3.RequestMethod_HEAD):
						http.Method = ptr.To("HEAD")
					case ext_core_v3.RequestMethod(corev3.RequestMethod_POST):
						http.Method = ptr.To("POST")
					case ext_core_v3.RequestMethod(corev3.RequestMethod_PUT):
						http.Method = ptr.To("PUT")
					case ext_core_v3.RequestMethod(corev3.RequestMethod_DELETE):
						http.Method = ptr.To("DELETE")
					case ext_core_v3.RequestMethod(corev3.RequestMethod_CONNECT):
						http.Method = ptr.To("CONNECT")
					case ext_core_v3.RequestMethod(corev3.RequestMethod_OPTIONS):
						http.Method = ptr.To("OPTIONS")
					case ext_core_v3.RequestMethod(corev3.RequestMethod_TRACE):
						http.Method = ptr.To("TRACE")
					case ext_core_v3.RequestMethod(corev3.RequestMethod_PATCH):
						http.Method = ptr.To("PATCH")
					}
					healthCheck.Http = http
				}
				if hc.GetGrpcHealthCheck() != nil {
					grpc := &kgateway.HealthCheckGrpc{
						ServiceName: ptr.To(hc.GetGrpcHealthCheck().GetServiceName()),
						Authority:   ptr.To(hc.GetGrpcHealthCheck().GetAuthority()),
					}
					if len(hc.GetGrpcHealthCheck().GetInitialMetadata()) > 0 {
						g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, upstream, "healthCheck.grpcHealthCheck.initialMetadata is not supported")
					}
					healthCheck.Grpc = grpc
				}
				if hc.GetAlwaysLogHealthCheckFailures() {
					g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, upstream, "healthCheck.alwaysLogHealthCheckFailures is not supported")
				}
				if hc.GetCustomHealthCheck() != nil {
					g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, upstream, "healthCheck.customHealthCheck is not supported")
				}
				if hc.GetEventLogPath() != "" {
					g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, upstream, "healthCheck.eventLogPath is not supported")
				}
				if hc.GetHealthyEdgeInterval() != nil {
					g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, upstream, "healthCheck.healthyEdgeInterval is not supported")
				}
				if hc.GetInitialJitter() != nil {
					g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, upstream, "healthCheck.initialJitter is not supported")
				}
				if hc.GetIntervalJitterPercent() != 0 {
					g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, upstream, "healthCheck.initialJitterPercent is not supported")
				}
				if hc.GetNoTrafficInterval() != nil {
					g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, upstream, "healthCheck.noTrafficInterval is not supported")
				}
				if hc.GetReuseConnection() != nil {
					g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, upstream, "healthCheck.reuseConnection is not supported")
				}
				if hc.GetTcpHealthCheck() != nil {
					g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, upstream, "healthCheck.tcpHealthCheck is not supported")
				}
				backendConfig.Spec.HealthCheck = healthCheck
				configExists = true
			}
		}
	}
	//if upstream.Spec.GetDiscoveryMetadata() != nil {
	//	o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, upstream, "discoveryMetadata is not supported")
	//}
	if upstream.Spec.GetDnsRefreshRate() != nil {
		g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, upstream, "dnsRefreshRate is not supported")
	}
	if upstream.Spec.GetFailover() != nil {
		g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, upstream, "failover is not supported")
	}

	if upstream.Spec.GetHttpConnectHeaders() != nil {
		g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, upstream, "httpConnectHeaders is not supported")
	}
	if upstream.Spec.GetHttpConnectSslConfig() != nil {
		g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, upstream, "httpConnectSslConfig is not supported")
	}
	if upstream.Spec.GetIgnoreHealthOnHostRemoval() != nil {
		g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, upstream, "ignoreHealthOnHostRemoval is not supported")
	}
	if upstream.Spec.GetInitialConnectionWindowSize() != nil {
		g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, upstream, "initialConnectionWindowSize is not supported")
	}
	if upstream.Spec.GetInitialStreamWindowSize() != nil {
		g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, upstream, "initialStreamWindowSize is not supported")
	}
	if upstream.Spec.GetHttpProxyHostname() != nil {
		g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, upstream, "httpProxyHostname is not supported")
	}
	if upstream.Spec.GetOutlierDetection() != nil {
		g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, upstream, "outlierDetection is not supported")
	}
	if upstream.Spec.GetOverrideStreamErrorOnInvalidHttpMessage() != nil {
		g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, upstream, "overrideStreamErrorOnInvalidHttpMessage is not supported")
	}
	if upstream.Spec.GetMaxConcurrentStreams() != nil {
		g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, upstream, "maxConcurrentStreams is not supported")
	}
	if upstream.Spec.GetPipe() != nil {
		g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, upstream, "pipe is not supported")
	}
	if upstream.Spec.GetPreconnectPolicy() != nil {
		g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, upstream, "preconnectPolicy is not supported")
	}
	if upstream.Spec.GetProxyProtocolVersion() != nil {
		g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, upstream, "proxy protocol version is not supported")
	}
	if upstream.Spec.GetRespectDnsTtl() != nil {
		g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, upstream, "respectDnsTtl is not supported")
	}

	if configExists {
		bcpw := snapshot.NewBackendConfigPolicyWrapper(backendConfig, upstream.FileOrigin())

		if upstream.Spec.GetKube() != nil {
			if upstream.GetNamespace() != upstream.Spec.GetKube().GetServiceNamespace() {
				g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, bcpw, "BackendConfigPolicy cannot apply to kube services in other namespaces.")
			}
		}
		g.gatewayAPICache.AddBackendConfigPolicy(bcpw)
	}
	return
}

func (g *GatewayAPIOutput) convertAWSBackend(upstream *snapshot.UpstreamWrapper, lambda *aws.LambdaFunctionSpec) *snapshot.BackendWrapper {
	backend := &snapshot.BackendWrapper{
		Backend: &kgateway.Backend{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Backend",
				APIVersion: kgateway.GroupVersion.String(),
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      upstream.GetName(),
				Namespace: upstream.GetNamespace(),
			},
			Spec: kgateway.BackendSpec{
				Type: kgateway.BackendTypeAWS,
				Aws: &kgateway.AwsBackend{
					AccountId: upstream.Spec.GetAws().GetAwsAccountId(),
					Region:    ptr.To(upstream.Spec.GetAws().GetRegion()),
				},
			},
		},
	}
	if upstream.Spec.GetAws().GetSecretRef() != nil {
		// if the upstream doesnt have the same namespace as the ARN secret we might have problems
		if upstream.Spec.GetAws().GetSecretRef().GetNamespace() != "" && upstream.Spec.GetAws().GetSecretRef().GetNamespace() != upstream.GetNamespace() {
			g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, upstream, "aws upstream references secret that exists in another namespace %s which is not supported", upstream.Spec.GetAws().GetSecretRef().GetNamespace())
		}
		backend.Spec.Aws.Auth = &kgateway.AwsAuth{
			Type: kgateway.AwsAuthTypeSecret,
			SecretRef: &corev1.LocalObjectReference{
				Name: upstream.Spec.GetAws().GetSecretRef().GetName(),
			},
		}
	}
	if lambda != nil {
		backend.Spec.Aws.Lambda = kgateway.AwsLambda{
			EndpointURL:          "",                             // existing
			FunctionName:         lambda.GetLambdaFunctionName(), // existing
			InvocationMode:       "",                             // existing
			Qualifier:            lambda.GetQualifier(),          // existing
			PayloadTransformMode: "",                             // existing
		}
	}
	return backend
}

func (g *GatewayAPIOutput) generateBackendRefForDelegateAction(
	r *gloogwv1.Route,
	wrapper snapshot.Wrapper,
) []*gwv1.HTTPBackendRef {
	var backends []*gwv1.HTTPBackendRef
	if r.GetDelegateAction().GetRef() != nil {
		delegate := r.GetDelegateAction().GetRef()
		backendRef := &gwv1.HTTPBackendRef{
			BackendRef: gwv1.BackendRef{
				BackendObjectReference: gwv1.BackendObjectReference{
					Name:      gwv1.ObjectName(delegate.GetName()),
					Namespace: (*gwv1.Namespace)(ptr.To(delegate.GetNamespace())),
					Kind:      (*gwv1.Kind)(ptr.To("HTTPRoute")),
					Group:     (*gwv1.Group)(ptr.To(gwv1.GroupVersion.Group)),
				},
			},
		}
		backends = append(backends, backendRef)
	} else if r.GetDelegateAction().GetSelector() != nil {

		selector := r.GetDelegateAction().GetSelector()
		namespaces := selector.GetNamespaces()
		if namespaces != nil || len(selector.GetNamespaces()) == 0 {
			// default namespace is gloo-system
			namespaces = []string{"gloo-system"}
		}

		for _, namespace := range selector.GetNamespaces() {
			if namespace == "*" {
				namespace = "all"
			}

			if len(selector.GetLabels()) > 1 {
				g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "has delegate action with more than one label selector which is not supported in Gateway API")
			}
			// create a backend ref for each label
			for _, v := range selector.GetLabels() {
				// just grab the first label
				backendRef := &gwv1.HTTPBackendRef{
					BackendRef: gwv1.BackendRef{
						BackendObjectReference: gwv1.BackendObjectReference{
							Name:      gwv1.ObjectName(v),                       // label value
							Namespace: ptr.To(gwv1.Namespace(namespace)),        // defaults to parent namespace if unset
							Kind:      ptr.To(gwv1.Kind("label")),               // label is the only value
							Group:     ptr.To(gwv1.Group(delegationLabelGroup)), // custom group for delegation
						},
					},
				}
				backends = append(backends, backendRef)
				break
			}
		}
	}

	return backends
}

// Converts a single upstream to a GatewayAPI backend ref
func (g *GatewayAPIOutput) generateBackendRefForSingleUpstream(r *gloogwv1.Route, wrapper snapshot.Wrapper) gwv1.HTTPBackendRef {
	upstream := r.GetRouteAction().GetSingle().GetUpstream()
	var backendRef gwv1.HTTPBackendRef

	//TODO we need to lookup the upstream to see if its kube and then just reference kube directly
	var up *snapshot.UpstreamWrapper
	//if it is not a kube service or does not need http2
	var upstreamNs = upstream.GetNamespace()
	if upstreamNs == "" {
		upstreamNs = wrapper.GetNamespace()
	}

	up = g.edgeCache.GetUpstream(types.NamespacedName{Name: upstream.GetName(), Namespace: upstreamNs})

	if up == nil {
		// unknown reference to backend
		g.AddErrorFromWrapper(ERROR_TYPE_UNKNOWN_REFERENCE, wrapper, "upstream %s not found, referencing unknown upstream backend ref", upstream.GetName())

		backendRef = gwv1.HTTPBackendRef{
			BackendRef: gwv1.BackendRef{
				BackendObjectReference: gwv1.BackendObjectReference{
					Name:      gwv1.ObjectName(upstream.GetName()),
					Namespace: (*gwv1.Namespace)(ptr.To(upstreamNs)),
					Kind:      (*gwv1.Kind)(ptr.To("Backend")),
					Group:     (*gwv1.Group)(ptr.To(glookube.GroupName)),
				},
			},
		}
	} else if up.Spec.GetKube() == nil {
		// non kubernetes upstream
		backendRef = gwv1.HTTPBackendRef{
			BackendRef: gwv1.BackendRef{
				BackendObjectReference: gwv1.BackendObjectReference{
					Name:      gwv1.ObjectName(upstream.GetName()),
					Namespace: (*gwv1.Namespace)(ptr.To(upstreamNs)),
					Kind:      (*gwv1.Kind)(ptr.To("Backend")),
					Group:     (*gwv1.Group)(ptr.To(glookube.GroupName)),
				},
			},
		}
	} else if up.Spec.GetKube() != nil && up.Spec.GetUseHttp2() != nil && up.Spec.GetUseHttp2().GetValue() == true {
		g.AddErrorFromWrapper(ERROR_TYPE_UPDATE_OBJECT, wrapper, "service %s/%s uses http2, update its k8s service appProtocol=http2", up.Spec.GetKube().GetServiceNamespace(), up.Spec.GetKube().GetServiceName())
		// normal backend ref but let the user know htey need to annotate their service
		backendRef = gwv1.HTTPBackendRef{
			BackendRef: gwv1.BackendRef{
				BackendObjectReference: gwv1.BackendObjectReference{
					Name:      gwv1.ObjectName(up.Spec.GetKube().GetServiceName()),
					Namespace: (*gwv1.Namespace)(ptr.To(up.Spec.GetKube().GetServiceNamespace())),
					Port:      ptr.To(gwv1.PortNumber(int32(up.Spec.GetKube().GetServicePort()))),
				},
			},
		}
	} else {
		//use kube backendref
		backendRef = gwv1.HTTPBackendRef{
			BackendRef: gwv1.BackendRef{
				BackendObjectReference: gwv1.BackendObjectReference{
					Name:      gwv1.ObjectName(up.Spec.GetKube().GetServiceName()),
					Namespace: (*gwv1.Namespace)(ptr.To(up.Spec.GetKube().GetServiceNamespace())),
					Port:      ptr.To(gwv1.PortNumber(int32(up.Spec.GetKube().GetServicePort()))),
				},
			},
		}
	}

	// AWS lambda integration
	if r.GetRouteAction().GetSingle().GetDestinationSpec() != nil && r.GetRouteAction().GetSingle().GetDestinationSpec().GetAws() != nil {
		// we need to add a parameter for the lambda name reference
		backendRef.Filters = append(backendRef.Filters, gwv1.HTTPRouteFilter{
			Type: gwv1.HTTPRouteFilterExtensionRef,
			ExtensionRef: &gwv1.LocalObjectReference{
				Kind:  "Parameter",
				Group: glookube.GroupName,
				Name:  (gwv1.ObjectName)(r.GetRouteAction().GetSingle().GetDestinationSpec().GetAws().GetLogicalName()),
			},
		})
	}
	return backendRef
}
