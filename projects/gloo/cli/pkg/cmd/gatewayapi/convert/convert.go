package convert

import (
	"fmt"
	corev3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	"github.com/solo-io/cue/pkg/time"
	ext_core_v3 "github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/config/core/v3"
	v4 "github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/extensions/filters/http/csrf/v3"
	gloo_type_matcher "github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/type/matcher/v3"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/jwt"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/rbac"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/cors"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/protocol"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/ssl"
	"k8s.io/apimachinery/pkg/api/resource"
	"strconv"
	"strings"

	v3 "github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/type/v3"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/ai"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/extauth/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/extproc"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/ratelimit"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/aws"
	transformation2 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/transformation"
	v1alpha2 "github.com/solo-io/solo-apis/pkg/api/ratelimit.solo.io/v1alpha1"
	corev1 "k8s.io/api/core/v1"

	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/als"
	"k8s.io/apimachinery/pkg/runtime"

	"k8s.io/apimachinery/pkg/types"

	"github.com/solo-io/gloo/projects/gloo/cli/pkg/snapshot"

	"encoding/json"

	kgateway "github.com/kgateway-dev/kgateway/v2/api/v1alpha1"
	gloogateway "github.com/solo-io/gloo-gateway/api/v1alpha1"
	gloogwv1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gateway2/api/v1alpha1"
	"github.com/solo-io/gloo/projects/gateway2/wellknown"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/core/matchers"
	glookube "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/kube/apis/gloo.solo.io/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
	gwv1 "sigs.k8s.io/gateway-api/apis/v1"
	apixv1a1 "sigs.k8s.io/gateway-api/apisx/v1alpha1"
)

const (
	perConnectionBufferLimit = "kgateway.dev/per-connection-buffer-limit"
	routeWeight              = "kgateway.dev/route-weight"
)

func (o *GatewayAPIOutput) Convert() error {

	// Convert upstreams to backends first so that we can reference them in the Settings and Routes
	for _, upstream := range o.edgeCache.Upstreams() {
		o.convertUpstreamToBackend(upstream)
	}

	for _, settings := range o.edgeCache.Settings() {
		// We only translate virtual services for ones that match a gateway selector
		// TODO in the future we could blindly convert VS and not attach it to anything
		err := o.convertSettings(settings)
		if err != nil {
			return err
		}
	}

	for _, gateway := range o.edgeCache.GlooGateways() {
		// We only translate virtual services for ones that match a gateway selector
		// TODO(nick) - in the future we could blindly convert VS and not attach it to anything
		err := o.convertGatewayAndVirtualServices(gateway)
		if err != nil {
			return err
		}
	}

	for _, routeTable := range o.edgeCache.RouteTables() {
		err := o.convertRouteTableToHTTPRoute(routeTable)
		if err != nil {
			return err
		}
	}

	return nil
}
func (o *GatewayAPIOutput) convertSettings(settings *snapshot.SettingsWrapper) error {
	if settings == nil {
		return nil
	}
	spec := settings.Settings.Spec

	if spec.GetDiscoveryNamespace() != "" {
		//TODO(nick): how is this set now?
	}
	if spec.GetWatchNamespaces() != nil {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, settings, "watchNamespaces is not supported")
	}
	if spec.GetKubernetesConfigSource() != nil {
		// this is a default
	}
	if spec.GetDirectoryConfigSource() != nil {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, settings, "directoryConfigSource is not supported")
	}
	if spec.GetConsulKvSource() != nil {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, settings, "consulKvSource is not supported")
	}
	if spec.GetKubernetesSecretSource() != nil {
		// This is the default
	}
	if spec.GetVaultSecretSource() != nil {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, settings, "vaultSecretSource is not supported")
	}
	if spec.GetDirectorySecretSource() != nil {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, settings, "directorySecretSource is not supported")
	}
	if spec.GetSecretOptions() != nil {
		// This is no longer needed and helps configure the secret source
	}
	if spec.GetKubernetesArtifactSource() != nil {
		// This is the default
	}
	if spec.GetDirectoryArtifactSource() != nil {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, settings, "directoryArtifactSource is not supported")
	}
	if spec.GetConsulKvArtifactSource() != nil {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, settings, "consulKvArtifactSouce is not supported")
	}
	if spec.GetRefreshRate() != nil {
		//this is not needed in kgateway, all users should start with the default
	}
	if spec.GetLinkerd() == true {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, settings, "linkerd no longer supported")
	}
	if spec.GetKnative() != nil {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, settings, "knative no longer supported")
	}
	if spec.GetDiscovery() != nil {
		// This is the default
		//  discovery:
		//    fdsMode: WHITELIST
		//TODO(nick) - not sure if there is anything to convert here, for now we will just ignore
	}
	if spec.GetGloo() != nil {
		gloo := spec.GetGloo()
		// here are the defaults
		//  gloo:
		//    disableKubernetesDestinations: false
		//    disableProxyGarbageCollection: false
		//    enableRestEds: false
		//    invalidConfigPolicy:
		//      invalidRouteResponseBody: Gloo Gateway has invalid configuration. Administrators
		//        should run `glooctl check` to find and fix config errors.
		//      invalidRouteResponseCode: 404
		//      replaceInvalidRoutes: false
		//    istioOptions:
		//      appendXForwardedHost: true
		//      enableAutoMtls: false
		//      enableIntegration: false
		//    proxyDebugBindAddr: 0.0.0.0:9966
		//    regexMaxProgramSize: 1024
		//    restXdsBindAddr: 0.0.0.0:9976
		//    xdsBindAddr: 0.0.0.0:9977

		//TODO(nick) - verify AWS support
		if gloo.GetAwsOptions() != nil && spec.GetGloo().GetAwsOptions().GetEnableCredentialsDiscovey() == true {
			o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, settings, "AWS credentialDiscovery is not supported")
		}
		if gloo.GetAwsOptions() != nil && spec.GetGloo().GetAwsOptions().GetFallbackToFirstFunction() != nil && spec.GetGloo().GetAwsOptions().GetFallbackToFirstFunction().GetValue() == true {
			o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, settings, "AWS fallbackToFirstFunction is not supported")
		}
		if gloo.GetAwsOptions() != nil && spec.GetGloo().GetAwsOptions().GetPropagateOriginalRouting() != nil && spec.GetGloo().GetAwsOptions().GetPropagateOriginalRouting().GetValue() == true {
			o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, settings, "AWS propagateOriginalRouting is not supported")
		}
		if gloo.GetAwsOptions() != nil && spec.GetGloo().GetAwsOptions().GetServiceAccountCredentials() != nil {
			o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, settings, "AWS serviceAccountCredentials is not supported")
		}
		if gloo.GetCircuitBreakers() != nil {
			o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, settings, "gloo.circuitBreakers is not supported")
		}
	}
	if spec.GetGateway() != nil {
		// These are the defaults
		//  gateway:
		//    enableGatewayController: true
		//    isolateVirtualHostsBySslConfig: false
		//    readGatewaysFromAllNamespaces: false
		//    validation:
		//      allowWarnings: true
		//      alwaysAccept: true
		//      disableTransformationValidation: false
		//      fullEnvoyValidation: false
		//      proxyValidationServerAddr: gloo:9988
		//      serverEnabled: true
		//      validationServerGrpcMaxSizeBytes: 104857600
		//      warnMissingTlsSecret: true
		//      warnRouteShortCircuiting: false
		if spec.GetGateway().GetValidation() != nil {
			// nothing to warn about here as we are hoping validation is better in kgateway
		}
	}

	if spec.GetConsul() != nil {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, settings, "consul is not supported")
	}
	if spec.GetConsulDiscovery() != nil {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, settings, "consul discovery is not supported")
	}
	if spec.GetKubernetes() != nil {
		// nothing to warn on here
	}
	if spec.GetExtensions() != nil {
		// TODO(nick): what are config extensions?
	}

	if spec.GetRatelimit() != nil {
		// descriptors are set here
		//"descriptors": []ratelimit.api.solo.io.Descriptor
		//"setDescriptors": []ratelimit.api.solo.io.SetDescriptor

		// TODO(nick): convert the descriptors to a RateLimitConfig that does not yet exist in gloo-gateway
	}

	if spec.GetRatelimitServer() != nil {
		extension := o.generateGatewayExtensionForRateLimit(spec.GetRatelimitServer(), "rate-limit", settings)
		if extension != nil {
			o.gatewayAPICache.AddGatewayExtension(snapshot.NewGatewayExtensionWrapper(extension, settings.FileOrigin()))
		}
	}
	if spec.GetRbac() != nil && spec.GetRbac().GetRequireRbac() == true {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, settings, "requireRbac is not supported")
	}
	if spec.GetExtauth() != nil {
		extension := o.generateGatewayExtensionForExtAuth(spec.GetExtauth(), "ext-authz", settings)
		if extension != nil {
			o.gatewayAPICache.AddGatewayExtension(snapshot.NewGatewayExtensionWrapper(extension, settings.FileOrigin()))
		}
	}
	if spec.GetNamedExtauth() != nil {
		for name, extauth := range spec.GetNamedExtauth() {
			extension := o.generateGatewayExtensionForExtAuth(extauth, name, settings)
			if extension != nil {
				o.gatewayAPICache.AddGatewayExtension(snapshot.NewGatewayExtensionWrapper(extension, settings.FileOrigin()))
			}
		}
	}
	if spec.GetCachingServer() != nil {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, settings, "cachingServer is not supported")
	}
	if spec.GetMetadata() != nil {
		// not sure if we should do anything here
	}
	if spec.GetNamespacedStatuses() != nil {
		// no need to do anything here
	}
	if spec.GetObservabilityOptions() != nil {
		// no need to do anything here
	}
	if spec.GetUpstreamOptions() != nil {
		if spec.GetUpstreamOptions().GetSslParameters() != nil {
			o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, settings, "global upstream SSLParameters  is not supported")
		}
	}
	if spec.GetConsoleOptions() != nil {
		// no need to do anything here
	}
	//if spec.GetGraphqlOptions() != nil {
	//	o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, settings, "graphql is not supported")
	//}
	if spec.GetExtProc() != nil {
		o.generateGatewayExtensionForExtProc(spec.GetExtProc(), "global-ext-proc", settings)
	}
	if spec.GetKnative() != nil {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, settings, "knative is not supported")
	}
	if spec.GetWatchNamespaceSelectors() != nil {
		// nothing to do here
	}

	return nil
}

// TODO(nick): does aws backend support awsec2 upstream?
func (o *GatewayAPIOutput) convertUpstreamToBackend(upstream *snapshot.UpstreamWrapper) {
	// Add all existing upstreams except for kube services which will be referenced directly
	if upstream.Spec.GetKube() != nil {
		// do nothing, let it continue in case there were other policies attached to the kube that we can warn about
	}
	if upstream.Spec.GetAi() != nil {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, upstream, "gcp AI is not supported")
	}

	// determine if we need to create an upstream policy which will apply to the upstream
	// TODO we might need to see if the upstream is a kube service and then apply it to that?
	o.convertUpstreamPolicy(upstream)

	if upstream.Spec.GetAws() != nil {
		if len(upstream.Spec.GetAws().GetLambdaFunctions()) > 0 {
			backend := o.convertAWSBackend(upstream, nil)
			o.gatewayAPICache.AddBackend(backend)
		} else {
			//TODO (nick): we create multiple backends here but we need to fix naming and backendPolicies
			for _, lambda := range upstream.Spec.GetAws().GetLambdaFunctions() {
				backend := o.convertAWSBackend(upstream, lambda)
				o.gatewayAPICache.AddBackend(backend)
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
				o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, upstream, "static upstream healthCheckConfig is not supported")
			}
			if hosts.GetLoadBalancingWeight() != nil {
				o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, upstream, "static upstream loadBalancingWeight is not supported")
			}
			if hosts.GetSniAddr() != "" {
				o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, upstream, "static upstream sni is not supported")
			}

			if hosts.GetAddr() != "" {
				host := kgateway.Host{
					Host: hosts.GetAddr(),
				}
				if hosts.GetPort() != 0 {
					host.Port = gwv1.PortNumber(hosts.GetPort())
				}
				o.convertUpstreamPolicy(upstream)
				backend.Spec.Static.Hosts = append(backend.Spec.Static.Hosts, host)
			}
		}
		o.gatewayAPICache.AddBackend(backend)
	}
	if upstream.Spec.GetAwsEc2() != nil {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, upstream, "awsec2 upstream is not supported")
	}
	if upstream.Spec.GetConsul() != nil {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, upstream, "consul upstream is not supported")
	}
	if upstream.Spec.GetAzure() != nil {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, upstream, "azure upstream is not supported")
	}
	if upstream.Spec.GetGcp() != nil {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, upstream, "gcp upstream is not supported")
	}

}
func (o *GatewayAPIOutput) convertUpstreamPolicy(upstream *snapshot.UpstreamWrapper) {
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
			o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, upstream, "static upstream autoSniRewrite is not supported")
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
				o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, upstream, "sslConfig.secretRef.namespace %s is not the same as the backendConfig's %s", upstream.GetNamespace(), upstream.Spec.GetSslConfig().GetSecretRef().GetNamespace())
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
				o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, upstream, "sslConfig.sslFiles.ocspStaple is not supported")
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
			o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, upstream, "sslConfig.sds is not supported")
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
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, upstream, "circuitBreakers is not supported")
	}
	if upstream.Spec.GetConnectionConfig() != nil {
		if upstream.Spec.GetConnectionConfig().GetPerConnectionBufferLimitBytes() != nil {
			backendConfig.Spec.PerConnectionBufferLimitBytes = ptr.To(int(upstream.Spec.GetConnectionConfig().GetPerConnectionBufferLimitBytes().Value))
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
				o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, upstream, "commonHTTPProtocolOptions.headersWithUndercoresAction is not supported")
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
				o.AddErrorFromWrapper(ERROR_TYPE_IGNORED, upstream, "commonHealthCheck.healthChecks only using first health check")
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
						o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, upstream, "healthCheck.httpHealthCheck.useHTTP2 is not supported")
					}
					if hc.GetHttpHealthCheck().GetServiceName() != "" {
						o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, upstream, "healthCheck.httpHealthCheck.serviceName is not supported")
					}
					if hc.GetHttpHealthCheck().GetExpectedStatuses() != nil {
						o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, upstream, "healthCheck.httpHealthCheck.expectedStatuses is not supported")
					}
					if hc.GetHttpHealthCheck().GetRequestHeadersToAdd() != nil {
						o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, upstream, "healthCheck.httpHealthCheck.requestHeadersToAdd is not supported")
					}
					if hc.GetHttpHealthCheck().GetRequestHeadersToRemove() != nil {
						o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, upstream, "healthCheck.httpHealthCheck.requestHeadersToRemove is not supported")
					}
					if hc.GetHttpHealthCheck().GetResponseAssertions() != nil {
						o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, upstream, "healthCheck.httpHealthCheck.responseAssertions is not supported")
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
						o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, upstream, "healthCheck.grpcHealthCheck.initialMetadata is not supported")
					}
					healthCheck.Grpc = grpc
				}
				if hc.GetAlwaysLogHealthCheckFailures() {
					o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, upstream, "healthCheck.alwaysLogHealthCheckFailures is not supported")
				}
				if hc.GetCustomHealthCheck() != nil {
					o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, upstream, "healthCheck.customHealthCheck is not supported")
				}
				if hc.GetEventLogPath() != "" {
					o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, upstream, "healthCheck.eventLogPath is not supported")
				}
				if hc.GetHealthyEdgeInterval() != nil {
					o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, upstream, "healthCheck.healthyEdgeInterval is not supported")
				}
				if hc.GetInitialJitter() != nil {
					o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, upstream, "healthCheck.initialJitter is not supported")
				}
				if hc.GetIntervalJitterPercent() != 0 {
					o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, upstream, "healthCheck.initialJitterPercent is not supported")
				}
				if hc.GetNoTrafficInterval() != nil {
					o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, upstream, "healthCheck.noTrafficInterval is not supported")
				}
				if hc.GetReuseConnection() != nil {
					o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, upstream, "healthCheck.reuseConnection is not supported")
				}
				if hc.GetTcpHealthCheck() != nil {
					o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, upstream, "healthCheck.tcpHealthCheck is not supported")
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
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, upstream, "dnsRefreshRate is not supported")
	}
	if upstream.Spec.GetFailover() != nil {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, upstream, "failover is not supported")
	}

	if upstream.Spec.GetHttpConnectHeaders() != nil {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, upstream, "httpConnectHeaders is not supported")
	}
	if upstream.Spec.GetHttpConnectSslConfig() != nil {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, upstream, "httpConnectSslConfig is not supported")
	}
	if upstream.Spec.GetIgnoreHealthOnHostRemoval() != nil {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, upstream, "ignoreHealthOnHostRemoval is not supported")
	}
	if upstream.Spec.GetInitialConnectionWindowSize() != nil {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, upstream, "initialConnectionWindowSize is not supported")
	}
	if upstream.Spec.GetInitialStreamWindowSize() != nil {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, upstream, "initialStreamWindowSize is not supported")
	}
	if upstream.Spec.GetHttpProxyHostname() != nil {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, upstream, "httpProxyHostname is not supported")
	}
	if upstream.Spec.GetOutlierDetection() != nil {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, upstream, "outlierDetection is not supported")
	}
	if upstream.Spec.GetOverrideStreamErrorOnInvalidHttpMessage() != nil {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, upstream, "overrideStreamErrorOnInvalidHttpMessage is not supported")
	}
	if upstream.Spec.GetMaxConcurrentStreams() != nil {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, upstream, "maxConcurrentStreams is not supported")
	}
	if upstream.Spec.GetPipe() != nil {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, upstream, "pipe is not supported")
	}
	if upstream.Spec.GetPreconnectPolicy() != nil {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, upstream, "preconnectPolicy is not supported")
	}
	if upstream.Spec.GetProxyProtocolVersion() != nil {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, upstream, "proxy protocol version is not supported")
	}
	if upstream.Spec.GetRespectDnsTtl() != nil {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, upstream, "respectDnsTtl is not supported")
	}

	if configExists {
		bcpw := snapshot.NewBackendConfigPolicyWrapper(backendConfig, upstream.FileOrigin())

		if upstream.Spec.GetKube() != nil {
			if upstream.GetNamespace() != upstream.Spec.GetKube().GetServiceNamespace() {
				o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, bcpw, "BackendConfigPolicy cannot apply to kube services in other namespaces.")
			}
		}
		o.gatewayAPICache.AddBackendConfigPolicy(bcpw)
	}
	return
}

func (o *GatewayAPIOutput) convertAWSBackend(upstream *snapshot.UpstreamWrapper, lambda *aws.LambdaFunctionSpec) *snapshot.BackendWrapper {
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
			o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, upstream, "aws upstream references secret that exists in another namespace %s which is not supported", upstream.Spec.GetAws().GetSecretRef().GetNamespace())
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

func (o *GatewayAPIOutput) convertGatewayAndVirtualServices(glooGateway *snapshot.GlooGatewayWrapper) error {

	// we first need to generate Gateway objects with the correct names based on proxy Names
	// spec.proxyNames
	o.generateGatewaysFromProxyNames(glooGateway)

	gatewayVs, err := o.edgeCache.GlooGatewayVirtualServices(glooGateway)
	if err != nil {
		return err
	}
	if len(gatewayVs) == 0 {
		o.AddErrorFromWrapper(ERROR_TYPE_NO_REFERENCES, glooGateway, "gateway does not contain virtual services")
	}
	for _, vs := range gatewayVs {
		proxyNames := glooGateway.Spec.GetProxyNames()
		if len(proxyNames) == 0 {
			proxyNames = append(proxyNames, "gateway-proxy")
		}
		for _, proxyName := range proxyNames {
			listenerName := fmt.Sprintf("%s-%d-%s-%s", proxyName, glooGateway.Spec.GetBindPort(), vs.Name, vs.Namespace)
			// convert the listener portion of the virtual service
			if err := o.convertVirtualServiceListener(vs, glooGateway, listenerName, proxyName); err != nil {
				return err
			}
			// convert the routing portion of the virtual service
			err := o.convertVirtualServiceHTTPRoutes(vs, glooGateway, listenerName)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (o *GatewayAPIOutput) convertVirtualServiceListener(vs *snapshot.VirtualServiceWrapper, glooGateway *snapshot.GlooGatewayWrapper, listenerName string, gatewayName string) error {

	// for each VirtualService generate a listener set given the gateway port
	listenerSet := &apixv1a1.XListenerSet{
		TypeMeta: metav1.TypeMeta{
			Kind:       "XListenerSet",
			APIVersion: apixv1a1.GroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      listenerName,
			Namespace: vs.GetNamespace(),
			Labels:    vs.GetLabels(),
		},
		Spec: apixv1a1.ListenerSetSpec{
			ParentRef: apixv1a1.ParentGatewayReference{
				Group:     ptr.To(gwv1.Group(wellknown.GatewayGroup)),
				Kind:      ptr.To(gwv1.Kind(wellknown.GatewayKind)),
				Namespace: ptr.To(gwv1.Namespace(glooGateway.GetNamespace())),
				Name:      gwv1.ObjectName(gatewayName),
			},
			Listeners: []apixv1a1.ListenerEntry{},
		},
	}

	// we only create the listener part, not the http matchers
	for _, hostname := range vs.Spec.GetVirtualHost().GetDomains() {
		if strings.Contains(hostname, ":") {
			o.AddErrorFromWrapper(ERROR_TYPE_IGNORED, vs, "contains port in hostname %s, its being ignored for ListenerSet %s/%s", hostname, listenerSet.Namespace, listenerSet.Name)
			continue
		}

		// listener entry does not support wildcard
		listenerEntryName := strings.ReplaceAll(fmt.Sprintf("%s-%s", vs.Name, hostname), "*", "star")
		entry := apixv1a1.ListenerEntry{
			Name:     gwv1.SectionName(listenerEntryName),
			Hostname: ptr.To(gwv1.Hostname(hostname)),
			Port:     apixv1a1.PortNumber(glooGateway.Spec.GetBindPort()),
			Protocol: gwv1.HTTPProtocolType,
		}
		if vs.Spec.GetSslConfig() != nil {
			tlsConfig := o.generateTLSConfiguration(vs)
			if tlsConfig != nil {
				entry.TLS = tlsConfig
				entry.Protocol = gwv1.HTTPSProtocolType
			}
		}

		// allowed routes
		entry.AllowedRoutes = &gwv1.AllowedRoutes{
			Namespaces: &gwv1.RouteNamespaces{
				From: ptr.To(gwv1.NamespacesFromAll),
			},
		}
		listenerSet.Spec.Listeners = append(listenerSet.Spec.Listeners, entry)
	}

	if vs.Spec.GetVirtualHost().GetOptionsConfigRefs() != nil && len(vs.Spec.GetVirtualHost().GetOptionsConfigRefs().GetDelegateOptions()) > 0 {
		delegateOptions := vs.Spec.GetVirtualHost().GetOptionsConfigRefs().GetDelegateOptions()
		for _, delegateOption := range delegateOptions {
			// check to see if this already exists in gatewayAPI cache, if not move it over from edge cache
			gtp, exists := o.gatewayAPICache.GlooTrafficPolicies[types.NamespacedName{Name: delegateOption.GetName(), Namespace: delegateOption.GetNamespace()}]
			if !exists {
				vho, exists := o.edgeCache.VirtualHostOptions()[types.NamespacedName{Name: delegateOption.GetName(), Namespace: delegateOption.GetNamespace()}]
				if !exists {
					o.AddErrorFromWrapper(ERROR_TYPE_UNKNOWN_REFERENCE, vs, "references VirtualHostOption %s that does not exist", types.NamespacedName{Name: delegateOption.GetName(), Namespace: delegateOption.GetNamespace()})
					continue
				}
				gtp = o.convertVirtualHostOptionToGlooTrafficPolicy(vho)
				if gtp == nil {
					o.AddErrorFromWrapper(ERROR_TYPE_IGNORED, vs, "references VirtualHostOption %s - No options converted", types.NamespacedName{Name: delegateOption.GetName(), Namespace: delegateOption.GetNamespace()})
				}
			}
			if listenerSet.Namespace != gtp.GlooTrafficPolicy.GetNamespace() {
				o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, vs, "VirtualHostOption %s references a listener set in a different namespace %s which is not supported", types.NamespacedName{Name: vs.GetName(), Namespace: vs.GetNamespace()}, types.NamespacedName{Name: listenerSet.GetName(), Namespace: listenerSet.GetNamespace()})
			}
			// add the target ref to the listener
			gtp.GlooTrafficPolicy.Spec.TargetRefs = append(gtp.GlooTrafficPolicy.Spec.TargetRefs, kgateway.LocalPolicyTargetReferenceWithSectionName{
				LocalPolicyTargetReference: kgateway.LocalPolicyTargetReference{
					Group: apixv1a1.GroupName,
					Kind:  "XListenerSet",
					Name:  gwv1.ObjectName(listenerSet.Name),
				},
			})
			o.gatewayAPICache.AddGlooTrafficPolicy(gtp)
		}
	}

	// we need to get the virtualhostoptions and update their references
	if vs.Spec.GetVirtualHost().GetOptions() != nil {
		// create a separate virtualhost option and link it
		gtp := &gloogateway.GlooTrafficPolicy{
			TypeMeta: metav1.TypeMeta{
				Kind:       "GlooTrafficPolicy",
				APIVersion: gloogateway.GroupVersion.String(),
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      listenerSet.Name,
				Namespace: listenerSet.Namespace,
			},
		}
		// go through each option and add it to traffic policy
		spec := o.convertVHOOptionsToTrafficPolicySpec(vs.Spec.GetVirtualHost().GetOptions(), vs)

		// attach the xListenerSet to the GlooTrafficPolicy
		if spec.TrafficPolicySpec.TargetRefs == nil {
			spec.TrafficPolicySpec.TargetRefs = []kgateway.LocalPolicyTargetReferenceWithSectionName{}
		}
		spec.TrafficPolicySpec.TargetRefs = append(spec.TrafficPolicySpec.TargetRefs, kgateway.LocalPolicyTargetReferenceWithSectionName{
			LocalPolicyTargetReference: kgateway.LocalPolicyTargetReference{
				Group: apixv1a1.GroupName,
				Kind:  "XListenerSet",
				Name:  gwv1.ObjectName(listenerSet.Name),
			},
		})

		gtp.Spec = spec

		o.gatewayAPICache.AddGlooTrafficPolicy(snapshot.NewGlooTrafficPolicyWrapper(gtp, vs.FileOrigin()))
	}
	o.gatewayAPICache.AddListenerSet(snapshot.NewListenerSetWrapper(listenerSet, vs.FileOrigin()))

	return nil
}

func (o *GatewayAPIOutput) convertVirtualHostOptionToGlooTrafficPolicy(vho *snapshot.VirtualHostOptionWrapper) *snapshot.GlooTrafficPolicyWrapper {

	policy := &gloogateway.GlooTrafficPolicy{
		TypeMeta: metav1.TypeMeta{
			Kind:       "GlooTrafficPolicy",
			APIVersion: gloogateway.GroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      vho.GetName(),
			Namespace: vho.GetNamespace(),
		},
	}

	policy.Spec = o.convertVHOOptionsToTrafficPolicySpec(vho.VirtualHostOption.Spec.Options, vho)

	wrapper := snapshot.NewGlooTrafficPolicyWrapper(policy, vho.FileOrigin())
	return wrapper
}

func (o *GatewayAPIOutput) convertVHOOptionsToTrafficPolicySpec(vho *gloov1.VirtualHostOptions, wrapper snapshot.Wrapper) gloogateway.GlooTrafficPolicySpec {

	spec := gloogateway.GlooTrafficPolicySpec{
		TrafficPolicySpec: kgateway.TrafficPolicySpec{
			TargetRefs:      nil, // existing
			TargetSelectors: nil, // existing
			AI:              nil, // existing
			Transformation:  nil, // existing
			ExtProc:         nil, // existing
			ExtAuth:         nil, // existing
			RateLimit:       nil, // existing
			Cors:            nil, // existing
			Csrf:            nil, // existing
			Buffer:          nil, // existing
		},
		Waf:                      nil, // existing
		Retry:                    nil, // existing
		Timeouts:                 nil, // existing
		RateLimitEnterprise:      nil, // existing
		ExtAuthEnterprise:        nil, // existing
		TransformationEnterprise: nil, // existing
		JWTEnterprise:            nil, // existing
		RBACEnterprise:           nil, // existing
	}
	if vho != nil {
		if vho.GetExtauth() != nil {
			// we need to copy over the auth config ref if it exists
			ref := vho.GetExtauth().GetConfigRef()
			ac, exists := o.edgeCache.AuthConfigs()[types.NamespacedName{Name: ref.GetName(), Namespace: ref.GetNamespace()}]
			if !exists {
				o.AddErrorFromWrapper(ERROR_TYPE_UNKNOWN_REFERENCE, wrapper, "did not find AuthConfig %s/%s for delegated route option reference", ref.GetName(), ref.GetNamespace())
			} else {
				o.gatewayAPICache.AddAuthConfig(ac)

				spec.ExtAuthEnterprise = &gloogateway.ExtAuthEnterprise{
					ExtensionRef: &corev1.LocalObjectReference{
						Name: "ext-authz",
					},
					AuthConfigRef: gloogateway.AuthConfigRef{
						Name:      vho.GetExtauth().GetConfigRef().GetName(),
						Namespace: ptr.To(vho.GetExtauth().GetConfigRef().GetNamespace()),
					},
				}
			}
		}
		if vho.GetExtProc() != nil {
			// TODO(nick) the extproc on the VHO allows a user to disable the global
			// one but idk if there is an equivalent in gateway api?
			//extProc := &kgateway.ExtProcPolicy{
			//	ExtensionRef:   &corev1.LocalObjectReference{Name: "global-ext-proc"},
			//	ProcessingMode: &kgateway.ProcessingMode{
			//		RequestHeaderMode:   nil,
			//		ResponseHeaderMode:  nil,
			//		RequestBodyMode:     nil,
			//		ResponseBodyMode:    nil,
			//		RequestTrailerMode:  nil,
			//		ResponseTrailerMode: nil,
			//	},
			//}
			//
			//if vho.GetExtProc().GetDisabled() != nil {
			//
			//}
			//if vho.GetExtProc().GetOverride() != nil {}
			//
			//
			//spec.TrafficPolicySpec.ExtProc = extProc
		}
		if vho.GetWaf() != nil {
			waf := &gloogateway.Waf{
				Disabled:      ptr.To(vho.GetWaf().GetDisabled()),
				CustomMessage: ptr.To(vho.GetWaf().GetCustomInterventionMessage()),
				Rules:         []gloogateway.WafRule{},
			}
			for _, r := range vho.GetWaf().RuleSets {
				waf.Rules = append(waf.Rules, gloogateway.WafRule{
					RuleStr: ptr.To(r.RuleStr),
				})
				if r.Files != nil && len(r.Files) > 0 {
					o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "WAF files is not supported")
				}
			}
			spec.Waf = waf
		}
		if vho.GetRatelimitBasic() != nil {
			o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "rateLimitBasic is not supported")
		}
		if vho.GetRatelimitEarly() != nil {
			o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "rateLimitEarly is not supported, defaulting to regular rate limiting")

			rle := &gloogateway.RateLimitEnterprise{
				Global: &gloogateway.GlobalRateLimit{
					// Need to find the Gateway Extension for Global Rate Limit Server
					ExtensionRef: &corev1.LocalObjectReference{
						Name: "rate-limit",
					},

					RateLimits: []gloogateway.RateLimitActions{},
					// RateLimitConfig for the policy, not sure how it works for rate limit basic
					// TODO(nick) grab the global rate limit config ref
					RateLimitConfigRef: nil,
				},
			}
			for _, rl := range vho.GetRatelimitEarly().GetRateLimits() {
				rateLimit := &gloogateway.RateLimitActions{
					Actions:    []gloogateway.Action{},
					SetActions: []gloogateway.Action{},
				}
				for _, action := range rl.GetActions() {
					rateLimitAction := o.convertRateLimitAction(action)
					rateLimit.Actions = append(rateLimit.Actions, rateLimitAction)
				}
				for _, action := range rl.GetSetActions() {
					rateLimitAction := o.convertRateLimitAction(action)
					rateLimit.SetActions = append(rateLimit.SetActions, rateLimitAction)
				}
				if rl.GetLimit() != nil {
					o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "rateLimit action limit is not supported")
				}
			}
			spec.RateLimitEnterprise = rle
		}
		if vho.GetRatelimitRegular() != nil {
			rle := &gloogateway.RateLimitEnterprise{
				Global: &gloogateway.GlobalRateLimit{
					// Need to find the Gateway Extension for Global Rate Limit Server
					ExtensionRef: &corev1.LocalObjectReference{
						Name: "rate-limit",
					},

					RateLimits: []gloogateway.RateLimitActions{},
					// RateLimitConfig for the policy, not sure how it works for rate limit basic
					// TODO(nick) grab the global rate limit config ref
					RateLimitConfigRef: nil,
				},
			}
			for _, rl := range vho.GetRatelimitRegular().GetRateLimits() {
				rateLimit := &gloogateway.RateLimitActions{
					Actions:    []gloogateway.Action{},
					SetActions: []gloogateway.Action{},
				}
				for _, action := range rl.GetActions() {
					rateLimitAction := o.convertRateLimitAction(action)
					rateLimit.Actions = append(rateLimit.Actions, rateLimitAction)
				}
				for _, action := range rl.GetSetActions() {
					rateLimitAction := o.convertRateLimitAction(action)
					rateLimit.SetActions = append(rateLimit.SetActions, rateLimitAction)
				}
				if rl.GetLimit() != nil {
					o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "rateLimit action limit is not supported")
				}
			}
			spec.RateLimitEnterprise = rle
		}
		if vho.GetHeaderManipulation() != nil {
			// this is natively supported on the HTTPRoute
		}
		if vho.GetCors() != nil {
			policy := o.convertCORS(vho.GetCors(), wrapper)
			spec.Cors = policy
		}
		if vho.GetTransformations() != nil {
			// TODO(nick) should we try to translate this or require the end user to migrate to staged?
		}
		if vho.GetStagedTransformations() != nil {
			transformation := o.convertStagedTransformation(vho.GetStagedTransformations(), wrapper)
			spec.TransformationEnterprise = transformation
		}
		if vho.GetJwt() != nil {
			// TODO(nick) should we try to translate this or require the end user to migrate to staged?
			o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "jwt is deprecated in edge and supported")
		}
		if vho.GetJwtStaged() != nil {
			spec.JWTEnterprise = &gloogateway.StagedJWT{
				AfterExtAuth:  nil, // existing
				BeforeExtAuth: nil, // existing
			}
			if vho.GetJwtStaged().GetBeforeExtAuth() != nil {
				jwte := o.convertJWTStagedExtAuth(vho.GetJwtStaged().GetBeforeExtAuth(), wrapper)
				spec.JWTEnterprise.BeforeExtAuth = jwte
			}
			if vho.GetJwtStaged().GetAfterExtAuth() != nil {
				jwte := o.convertJWTStagedExtAuth(vho.GetJwtStaged().GetAfterExtAuth(), wrapper)
				spec.JWTEnterprise.BeforeExtAuth = jwte
			}
		}
		if vho.GetRbac() != nil {
			rbe := o.convertRBAC(vho.GetRbac())
			spec.RBACEnterprise = rbe
		}
		if vho.GetBufferPerRoute() != nil && vho.GetBufferPerRoute().GetBuffer() != nil && vho.GetBufferPerRoute().GetBuffer().GetMaxRequestBytes() != nil {
			spec.Buffer = &kgateway.Buffer{
				MaxRequestSize: resource.NewQuantity(int64(vho.GetBufferPerRoute().GetBuffer().GetMaxRequestBytes().GetValue()), resource.BinarySI),
			}
			if vho.GetBufferPerRoute().GetDisabled() {
				o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "bufferPerRoute.disabled is not supported")
			}
		}
		if vho.GetIncludeRequestAttemptCount() != nil {
			o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "includeRequestAttemptCount is not supported")
		}
		if vho.GetIncludeAttemptCountInResponse() != nil {
			o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "includeRequestAttemptCountInResponse is not supported")
		}
		if vho.GetCorsPolicyMergeSettings() != nil {
			o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "corsPolicyMergeSettings is not supported")
		}
		if vho.GetDlp() != nil {
			o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "dlp is not supported")
		}
		if vho.GetCsrf() != nil {
			csrf := o.convertCSRF(vho.GetCsrf())
			spec.TrafficPolicySpec.Csrf = csrf
		}
		if vho.GetExtensions() != nil {
			o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "gloo edge extensions is not supported")
		}
	}
	return spec
}

func (o *GatewayAPIOutput) convertJWTStagedExtAuth(auth *jwt.VhostExtension, wrapper snapshot.Wrapper) *gloogateway.JWTEnterprise {
	jwte := &gloogateway.JWTEnterprise{
		Providers:        nil, // existing
		ValidationPolicy: nil, // existing
		Disable:          nil, // existing
	}

	switch auth.GetValidationPolicy() {
	case jwt.VhostExtension_REQUIRE_VALID:
		jwte.ValidationPolicy = ptr.To(gloogateway.ValidationPolicyRequireValid)
	case jwt.VhostExtension_ALLOW_MISSING:
		jwte.ValidationPolicy = ptr.To(gloogateway.ValidationPolicyAllowMissing)
	case jwt.VhostExtension_ALLOW_MISSING_OR_FAILED:
		jwte.ValidationPolicy = ptr.To(gloogateway.ValidationPolicyAllowMissingOrFailed)
	}

	if auth.GetProviders() != nil {
		jwte.Providers = make(map[string]gloogateway.JWTProvider)
		for k, provider := range auth.GetProviders() {
			p := gloogateway.JWTProvider{
				JWKS:                         nil, // existing
				Audiences:                    nil, // existing
				Issuer:                       ptr.To(provider.Issuer),
				TokenSource:                  nil, // existing
				KeepToken:                    ptr.To(provider.KeepToken),
				ClaimsToHeaders:              nil, // existing
				ClockSkewSeconds:             nil, // existing
				AttachFailedStatusToMetadata: ptr.To(provider.AttachFailedStatusToMetadata),
			}
			if provider.GetClockSkewSeconds() != nil {
				p.ClockSkewSeconds = ptr.To(provider.ClockSkewSeconds.Value)
			}
			if len(provider.GetAudiences()) > 0 {
				p.Audiences = provider.GetAudiences()
			}
			if len(provider.GetClaimsToHeaders()) > 0 {
				p.ClaimsToHeaders = make([]gloogateway.ClaimToHeader, 0)
				for _, h := range provider.GetClaimsToHeaders() {
					p.ClaimsToHeaders = append(p.ClaimsToHeaders, gloogateway.ClaimToHeader{
						Claim:  h.GetClaim(),
						Header: h.GetHeader(),
						Append: ptr.To(h.GetAppend()),
					})
				}
			}

			if provider.GetTokenSource() != nil {
				p.TokenSource = &gloogateway.TokenSource{
					Headers:     make([]gloogateway.TokenSourceHeaderSource, 0),
					QueryParams: provider.GetTokenSource().GetQueryParams(),
				}
				for _, h := range provider.GetTokenSource().GetHeaders() {
					p.TokenSource.Headers = append(p.TokenSource.Headers, gloogateway.TokenSourceHeaderSource{
						Header: h.GetHeader(),
						Prefix: ptr.To(h.GetPrefix()),
					})
				}
			}
			if provider.GetJwks() != nil {
				jwks := &gloogateway.JWKS{
					Local:  nil,
					Remote: nil,
				}
				if provider.GetJwks().GetLocal() != nil {
					jwks.Local = &gloogateway.LocalJWKS{Key: provider.GetJwks().GetLocal().GetKey()}
				}
				if provider.GetJwks().GetRemote() != nil {
					jwks.Remote = &gloogateway.RemoteJWKS{
						Url:           provider.GetJwks().GetRemote().GetUrl(),
						BackendRef:    nil, // existing
						CacheDuration: nil, // existing
						AsyncFetch:    nil, // existing
					}

					if provider.GetJwks().GetRemote().GetCacheDuration() != nil && provider.GetJwks().GetRemote().GetCacheDuration().Nanos != 0 {
						jwks.Remote.CacheDuration = &metav1.Duration{Duration: provider.GetJwks().GetRemote().CacheDuration.AsDuration()}
					}
					if provider.GetJwks().GetRemote().GetAsyncFetch() != nil {
						jwks.Remote.AsyncFetch = &gloogateway.JwksAsyncFetch{FastListener: ptr.To(provider.GetJwks().GetRemote().GetAsyncFetch().GetFastListener())}
					}

					if provider.GetJwks().GetRemote().GetUpstreamRef() != nil {

						backendRef := &gwv1.BackendRef{
							BackendObjectReference: gwv1.BackendObjectReference{
								Group:     nil,
								Kind:      nil,
								Namespace: nil,
								Port:      nil,
							},
							Weight: nil,
						}
						// need to look up the upstream to see if its kube or not
						upstream := o.GetEdgeCache().GetUpstream(types.NamespacedName{Name: provider.GetJwks().GetRemote().GetUpstreamRef().GetName(), Namespace: provider.GetJwks().GetRemote().GetUpstreamRef().GetNamespace()})
						if upstream == nil {
							// just treat it as a kube service because we dont know what it might be
							o.AddErrorFromWrapper(ERROR_TYPE_UNKNOWN_REFERENCE, wrapper, "jwtStaged remote jwks references upstream %s/%s which was not found", provider.GetJwks().GetRemote().GetUpstreamRef().GetNamespace(), provider.GetJwks().GetRemote().GetUpstreamRef().GetName())
							backendRef.Name = gwv1.ObjectName(provider.GetJwks().GetRemote().GetUpstreamRef().GetName())
							backendRef.Namespace = ptr.To(gwv1.Namespace(provider.GetJwks().GetRemote().GetUpstreamRef().GetNamespace()))
						} else {
							if upstream.Upstream.Spec.GetKube() != nil {
								// references a kubernetes service
								backendRef.Name = gwv1.ObjectName(upstream.Upstream.Spec.GetKube().GetServiceName())
								backendRef.Namespace = ptr.To(gwv1.Namespace(upstream.Upstream.Spec.GetKube().GetServiceNamespace()))
								backendRef.Port = ptr.To(gwv1.PortNumber(upstream.Upstream.Spec.GetKube().GetServicePort()))
							} else {
								// it needs to reference a backend
								backendRef.Name = gwv1.ObjectName(upstream.Name)
								backendRef.Namespace = ptr.To(gwv1.Namespace(upstream.Namespace))
								backendRef.Kind = (*gwv1.Kind)(ptr.To("Backend"))
								backendRef.Group = (*gwv1.Group)(ptr.To(glookube.GroupName))
							}
						}
						jwks.Remote.BackendRef = backendRef
					}
				}
				p.JWKS = jwks
			}

			jwte.Providers[k] = p
		}
	}

	return jwte
}

func (o *GatewayAPIOutput) convertCORS(policy *cors.CorsPolicy, wrapper snapshot.Wrapper) *kgateway.CorsPolicy {
	filter := &gwv1.HTTPCORSFilter{
		AllowOrigins:     []gwv1.AbsoluteURI{},
		AllowCredentials: gwv1.TrueField(policy.GetAllowCredentials()),
		AllowMethods:     []gwv1.HTTPMethodWithWildcard{},
		AllowHeaders:     []gwv1.HTTPHeaderName{},
		ExposeHeaders:    []gwv1.HTTPHeaderName{},
		MaxAge:           0,
	}
	if policy.GetAllowOrigin() != nil {
		for _, origin := range policy.GetAllowOrigin() {
			filter.AllowOrigins = append(filter.AllowOrigins, gwv1.AbsoluteURI(origin))
		}
	}
	if policy.GetAllowMethods() != nil {
		for _, method := range policy.GetAllowMethods() {
			filter.AllowMethods = append(filter.AllowMethods, gwv1.HTTPMethodWithWildcard(method))
		}
	}
	if policy.GetAllowHeaders() != nil {
		for _, header := range policy.GetAllowHeaders() {
			filter.AllowHeaders = append(filter.AllowHeaders, gwv1.HTTPHeaderName(header))
		}
	}
	if policy.GetExposeHeaders() != nil {
		for _, header := range policy.GetExposeHeaders() {
			filter.ExposeHeaders = append(filter.ExposeHeaders, gwv1.HTTPHeaderName(header))
		}
	}
	if policy.GetMaxAge() != "" {
		age, err := strconv.Atoi(policy.GetMaxAge())
		if err != nil {
			// try to parse duration
			duration, err := time.ParseDuration(policy.GetMaxAge())
			if err != nil {
				o.AddErrorFromWrapper(ERROR_TYPE_IGNORED, wrapper, "invalid max age %s", policy.GetMaxAge())
			} else {
				filter.MaxAge = int32(duration / time.Second)
			}
		} else {
			filter.MaxAge = int32(age)
		}
	}
	if policy.GetAllowOriginRegex() != nil {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "allowOriginRegex not supported")

	}
	if policy.GetDisableForRoute() != true {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "cors disabledForRoute not supported")
	}
	return &kgateway.CorsPolicy{
		HTTPCORSFilter: filter,
	}
}

func (o *GatewayAPIOutput) convertStagedTransformation(transformation *transformation2.TransformationStages, wrapper snapshot.Wrapper) *gloogateway.TransformationEnterprise {
	stagedTransformations := &gloogateway.TransformationEnterprise{
		Stages:    &gloogateway.StagedTransformations{}, // existing
		AWSLambda: nil,                                  // existing
	}

	if transformation.GetInheritTransformation() == true {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "transformation inherit transformation is not supported")
	}
	if transformation.GetEarly() != nil {
		routing := o.convertRequestTransformation(transformation.GetEarly(), wrapper)
		stagedTransformations.Stages.Early = routing
	}
	if transformation.GetRegular() != nil {
		routing := o.convertRequestTransformation(transformation.GetRegular(), wrapper)
		stagedTransformations.Stages.Regular = routing
	}
	if transformation.GetPostRouting() != nil {
		routing := o.convertRequestTransformation(transformation.GetPostRouting(), wrapper)
		stagedTransformations.Stages.PostRouting = routing
	}

	if transformation.GetLogRequestResponseInfo() != nil && transformation.GetLogRequestResponseInfo().GetValue() == true {
		stagedTransformations.Stages.LogRequestResponseInfo = ptr.To(true)
	}

	if transformation.GetEscapeCharacters() != nil {
		if transformation.GetEscapeCharacters().GetValue() {
			stagedTransformations.Stages.EscapeCharacters = ptr.To(gloogateway.EscapeCharactersEscape)
		} else {
			stagedTransformations.Stages.EscapeCharacters = ptr.To(gloogateway.EscapeCharactersDontEscape)
		}
	}
	return stagedTransformations
}

func (o *GatewayAPIOutput) generateTLSConfiguration(vs *snapshot.VirtualServiceWrapper) *gwv1.GatewayTLSConfig {
	tlsConfig := &gwv1.GatewayTLSConfig{
		Mode: ptr.To(gwv1.TLSModeTerminate),
		//FrontendValidation: nil, // TODO do we need to set this?
		//Options:            nil, // TODO do we need to set this?
	}
	if vs.Spec.GetSslConfig().GetSecretRef() != nil {
		tlsConfig.CertificateRefs = []gwv1.SecretObjectReference{
			{
				Group:     ptr.To(gwv1.Group("")),
				Kind:      ptr.To(gwv1.Kind("Secret")),
				Name:      gwv1.ObjectName(vs.Spec.GetSslConfig().GetSecretRef().GetName()),
				Namespace: ptr.To(gwv1.Namespace(vs.Spec.GetSslConfig().GetSecretRef().GetNamespace())),
			},
		}
	}
	// TODO There is a situation where a SSLSecret contains a ca.crt which triggers mTLS in Gloo Edge we have no way to determine this
	if vs.Spec.GetSslConfig().GetSslFiles() != nil {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, vs, "has SSLFiles but its not supported in Gateway API")
	}
	if vs.Spec.GetSslConfig().GetSds() != nil {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, vs, "has SDS Certificates but its not supported in Gateway API")
	}
	if len(vs.Spec.GetSslConfig().GetVerifySubjectAltName()) > 0 {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, vs, "has VerifySubjectAltName but its not supported in Gateway API")
	}
	if len(vs.Spec.GetSslConfig().GetAlpnProtocols()) > 0 {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, vs, "has AlpnProtocols but its not supported in Gateway API")
	}
	if vs.Spec.GetSslConfig().GetOcspStaplePolicy() > 0 {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, vs, "has OcspStaplePolicy %d but its not supported in Gateway API", vs.Spec.GetSslConfig().GetOcspStaplePolicy())
	}

	return tlsConfig
}

func (o *GatewayAPIOutput) generateGatewaysFromProxyNames(glooGateway *snapshot.GlooGatewayWrapper) {

	proxyNames := glooGateway.Gateway.Spec.GetProxyNames()

	if len(proxyNames) == 0 {
		proxyNames = append(proxyNames, "gateway-proxy")
	}

	for _, proxyName := range proxyNames {
		// check to see if we already created the Gateway, if we did then just move on
		existingGw := o.gatewayAPICache.GetGateway(types.NamespacedName{Name: proxyName, Namespace: glooGateway.Gateway.Namespace})
		if existingGw == nil {
			// create a new gateway
			existingGw = snapshot.NewGatewayWrapper(&gwv1.Gateway{
				TypeMeta: metav1.TypeMeta{
					Kind:       "Gateway",
					APIVersion: gwv1.GroupVersion.String(),
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      proxyName,
					Namespace: glooGateway.Gateway.Namespace,
					Labels:    glooGateway.Gateway.Labels,
				},
				Spec: gwv1.GatewaySpec{
					AllowedListeners: &gwv1.AllowedListeners{
						Namespaces: &gwv1.ListenerNamespaces{
							From: ptr.To(gwv1.NamespacesFromAll),
						},
					},
					Listeners: []gwv1.Listener{
						{
							Name:     "dummy",
							Port:     8888,
							Protocol: "HTTP",
							AllowedRoutes: &gwv1.AllowedRoutes{
								Namespaces: &gwv1.RouteNamespaces{
									From: ptr.To(gwv1.NamespacesFromSelector),
									Selector: &metav1.LabelSelector{
										MatchLabels: map[string]string{
											"dummy": "dummy",
										},
									},
								},
							},
						},
					},
					GatewayClassName: "gloo-gateway",
				},
			}, glooGateway.FileOrigin())
		}
		// special case for per connection buffer limits to apply to the gateway as an annotation - https://github.com/kgateway-dev/kgateway/pull/11505
		if glooGateway.Spec.GetOptions() != nil && glooGateway.Spec.GetOptions().GetPerConnectionBufferLimitBytes() != nil && glooGateway.Spec.GetOptions().GetPerConnectionBufferLimitBytes().GetValue() != 0 {
			if existingGw.Annotations == nil {
				existingGw.Annotations = make(map[string]string)
			}
			existingGw.Annotations[perConnectionBufferLimit] = glooGateway.Spec.GetOptions().GetPerConnectionBufferLimitBytes().String()
		}

		o.gatewayAPICache.AddGateway(existingGw)

		if glooGateway.Spec.GetHttpGateway() != nil && glooGateway.Spec.GetHttpGateway().GetOptions() != nil {
			o.convertHTTPListenerOptions(glooGateway.Spec.GetHttpGateway().Options, glooGateway, proxyName)
		}
		if glooGateway.Spec.GetOptions() != nil && glooGateway.Spec.GetOptions() != nil {
			o.convertListenerOptions(glooGateway, proxyName)
		}
	}
}

func (o *GatewayAPIOutput) convertListenerOptions(glooGateway *snapshot.GlooGatewayWrapper, proxyName string) {
	options := glooGateway.Spec.GetOptions()
	if options == nil {
		return
	}
	listenerPolicy := &kgateway.HTTPListenerPolicy{
		TypeMeta: metav1.TypeMeta{
			Kind:       "HTTPListenerPolicy",
			APIVersion: kgateway.GroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      glooGateway.GetName(),
			Namespace: glooGateway.GetNamespace(),
			Labels:    glooGateway.Gateway.Labels,
		},
		Spec: kgateway.HTTPListenerPolicySpec{
			TargetRefs: []kgateway.LocalPolicyTargetReference{
				{
					Group: gwv1.Group(gwv1.GroupVersion.Group),
					Kind:  "Gateway",
					Name:  gwv1.ObjectName(proxyName),
				},
			},
		},
	}
	if options.GetExtensions() != nil {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, glooGateway, "gloo edge listener option extensions are not supported for HTTPTrafficPolicy")
	}
	if options.GetSocketOptions() != nil {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, glooGateway, "gloo edge listener option socket options are not supported for HTTPTrafficPolicy")
	}
	if options.GetAccessLoggingService() != nil {
		o.convertListenerOptionAccessLogging(glooGateway, options, listenerPolicy)
	}
	if options.GetListenerAccessLoggingService() != nil {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, glooGateway, "gloo edge listener option listenerAccessLoggingService is not supported for HTTPTrafficPolicy")
	}
	if options.GetConnectionBalanceConfig() != nil {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, glooGateway, "gloo edge listener option connectionBalanceConfig is not supported for HTTPTrafficPolicy")
	}
	//if options.GetPerConnectionBufferLimitBytes() != nil {
	// This is now set as an annotation on Gateway
	//	o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, glooGateway, "gloo edge listener option perConnectionBufferLimitBytes is not supported for HTTPTrafficPolicy")
	//}
	if options.GetProxyProtocol() != nil {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, glooGateway, "gloo edge listener option proxyProtocol is not supported for HTTPTrafficPolicy")
	}
	if options.GetTcpStats() != nil {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, glooGateway, "gloo edge listener option tcpStats is not supported for HTTPTrafficPolicy")
	}

	o.gatewayAPICache.AddHTTPListenerPolicy(snapshot.NewHTTPListenerPolicyWrapper(listenerPolicy, glooGateway.FileOrigin()))
}

func (o *GatewayAPIOutput) convertListenerOptionAccessLogging(glooGateway *snapshot.GlooGatewayWrapper, options *gloov1.ListenerOptions, listenerPolicy *kgateway.HTTPListenerPolicy) {
	accessLoggingService := options.GetAccessLoggingService()

	for _, edgeAccessLog := range accessLoggingService.GetAccessLog() {
		if listenerPolicy.Spec.AccessLog == nil {
			listenerPolicy.Spec.AccessLog = []kgateway.AccessLog{}
		}
		accessLog := kgateway.AccessLog{
			FileSink:    nil, // existing
			GrpcService: nil, // existing
			Filter:      nil, // existing
		}
		if edgeAccessLog.GetFileSink() != nil {
			fileSink := &kgateway.FileSink{
				Path: edgeAccessLog.GetFileSink().Path,
			}
			if jsonFormat := edgeAccessLog.GetFileSink().GetJsonFormat(); jsonFormat != nil {
				jsonBytes, err := json.Marshal(jsonFormat.AsMap())
				if err != nil {
					o.AddErrorFromWrapper(ERROR_TYPE_IGNORED, glooGateway, "unable to marshal json format for accessLoggingService %v", err)
				} else {
					fileSink.JsonFormat = &runtime.RawExtension{Raw: jsonBytes}
				}
			}
			if edgeAccessLog.GetFileSink().GetStringFormat() != "" {
				fileSink.StringFormat = edgeAccessLog.GetFileSink().GetStringFormat()
			}
			accessLog.FileSink = fileSink
		}
		if edgeAccessLog.GetGrpcService() != nil {
			accessLog.GrpcService = &kgateway.AccessLogGrpcService{
				CommonAccessLogGrpcService: kgateway.CommonAccessLogGrpcService{
					//CommonGrpcService: nil,// TODO(nick) what do we need to set here?
					LogName: edgeAccessLog.GetGrpcService().LogName,
				},
				AdditionalRequestHeadersToLog:   edgeAccessLog.GetGrpcService().AdditionalRequestHeadersToLog,
				AdditionalResponseHeadersToLog:  edgeAccessLog.GetGrpcService().AdditionalResponseHeadersToLog,
				AdditionalResponseTrailersToLog: edgeAccessLog.GetGrpcService().AdditionalResponseTrailersToLog,
			}

			// backend Ref
			switch edgeAccessLog.GetGrpcService().GetServiceRef().(type) {
			case *als.GrpcService_StaticClusterName:
				accessLog.GrpcService.BackendRef = &gwv1.BackendRef{
					BackendObjectReference: gwv1.BackendObjectReference{
						Name:      gwv1.ObjectName(edgeAccessLog.GetGrpcService().GetStaticClusterName()),
						Namespace: nil,
						Port:      nil,
					},
				}
				o.AddErrorFromWrapper(ERROR_TYPE_UNKNOWN_REFERENCE, glooGateway, "", edgeAccessLog.GetGrpcService().GetStaticClusterName())
			}
			if edgeAccessLog.GetGrpcService().GetFilterStateObjectsToLog() != nil {
				o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, glooGateway, "accessLoggingService grpcService has filterStateObjectsToLog but its not supported")
			}
		}
		if edgeAccessLog.GetFilter() != nil {
			if accessLog.Filter == nil {
				accessLog.Filter = &kgateway.AccessLogFilter{}
			}
			if edgeAccessLog.GetFilter().GetOrFilter() != nil {

				if accessLog.Filter.OrFilter == nil {
					accessLog.Filter.OrFilter = []kgateway.FilterType{}
				}
				for _, filter := range edgeAccessLog.GetFilter().GetOrFilter().GetFilters() {
					accessLog.Filter.AndFilter = append(accessLog.Filter.AndFilter, *o.convertAccessLogFitler(filter, glooGateway))
				}
			} else if edgeAccessLog.GetFilter().GetAndFilter() != nil {
				if accessLog.Filter.AndFilter == nil {
					accessLog.Filter.AndFilter = []kgateway.FilterType{}
				}
				for _, filter := range edgeAccessLog.GetFilter().GetAndFilter().GetFilters() {
					accessLog.Filter.AndFilter = append(accessLog.Filter.AndFilter, *o.convertAccessLogFitler(filter, glooGateway))
				}
			} else {
				// just and inline filter
				accessLog.Filter.FilterType = o.convertAccessLogFitler(edgeAccessLog.GetFilter(), glooGateway)
			}
		}
		listenerPolicy.Spec.AccessLog = append(listenerPolicy.Spec.AccessLog, accessLog)
	}
}

func (o *GatewayAPIOutput) convertAccessLogFitler(filter *als.AccessLogFilter, wrapper snapshot.Wrapper) *kgateway.FilterType {

	filterType := &kgateway.FilterType{
		StatusCodeFilter:     nil,   // existing
		DurationFilter:       nil,   // existing
		NotHealthCheckFilter: false, // existing
		TraceableFilter:      false, // existing
		HeaderFilter:         nil,   // existing
		ResponseFlagFilter:   nil,   // existing
		GrpcStatusFilter:     nil,   // existing
		CELFilter:            nil,   // existing
	}

	if filter.GetDurationFilter() != nil {
		filterType.DurationFilter = &kgateway.DurationFilter{
			Op:    kgateway.Op(filter.GetDurationFilter().GetComparison().Op),
			Value: filter.GetDurationFilter().GetComparison().GetValue().GetDefaultValue(),
		}
	}
	if filter.GetHeaderFilter() != nil && filter.GetHeaderFilter().GetHeader() != nil {
		headerMatch := gwv1.HTTPHeaderMatch{
			Name: gwv1.HTTPHeaderName(filter.GetHeaderFilter().GetHeader().Name),
		}

		if filter.GetHeaderFilter().GetHeader().GetExactMatch() != "" {
			o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "accessLoggingService grpcService has filterStateObjectsToLog but its not supported")
		}

		if filter.GetHeaderFilter().GetHeader().GetInvertMatch() == true {
			o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "accessLoggingService filter header invert match is not supported")
		}

		if filter.GetHeaderFilter().GetHeader().GetPresentMatch() == true {
			// TODO(nick): is this supported in Gateway API?
			o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "accessLoggingService filter header range match match is not supported")
		}

		if filter.GetHeaderFilter().GetHeader().GetPrefixMatch() != "" {
			//	HeaderMatchExact             HeaderMatchType = "Exact"
			//	HeaderMatchRegularExpression HeaderMatchType = "RegularExpression"
			//TODO(nick): can someone verify this is the equivalent?
			headerMatch.Type = ptr.To(gwv1.HeaderMatchRegularExpression)
			headerMatch.Value = filter.GetHeaderFilter().GetHeader().GetPrefixMatch() + ".*"
		}

		if filter.GetHeaderFilter().GetHeader().GetRangeMatch() != nil {
			o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "accessLoggingService filter header range match match is not supported")
		}

		if filter.GetHeaderFilter().GetHeader().GetSafeRegexMatch() != nil {
			// Edge only supported Googles Regex (RE2) which might not be compatible with Gateway API regex
			headerMatch.Type = ptr.To(gwv1.HeaderMatchRegularExpression)
			headerMatch.Value = filter.GetHeaderFilter().GetHeader().GetSafeRegexMatch().Regex
		}

		if filter.GetHeaderFilter().GetHeader().GetSuffixMatch() != "" {
			//TODO(nick): can someone verify this is the equivalent?
			headerMatch.Type = ptr.To(gwv1.HeaderMatchRegularExpression)
			headerMatch.Value = ".*" + filter.GetHeaderFilter().GetHeader().GetPrefixMatch()
		}

		filterType.HeaderFilter = &kgateway.HeaderFilter{
			Header: headerMatch,
		}
	}

	if filter.GetGrpcStatusFilter() != nil {
		grpcFilter := &kgateway.GrpcStatusFilter{
			Statuses: []kgateway.GrpcStatus{},
			Exclude:  filter.GetGrpcStatusFilter().Exclude,
		}
		for _, status := range filter.GetGrpcStatusFilter().Statuses {
			grpcFilter.Statuses = append(grpcFilter.Statuses, kgateway.GrpcStatus(status))
		}
		filterType.GrpcStatusFilter = grpcFilter
	}

	if filter.GetNotHealthCheckFilter() != nil {
		//unsure if this is correct. it appears this just needs to exist to function?
		filterType.NotHealthCheckFilter = true
	}

	if filter.GetResponseFlagFilter() != nil {
		filterType.ResponseFlagFilter = &kgateway.ResponseFlagFilter{
			Flags: filter.GetResponseFlagFilter().GetFlags(),
		}
	}

	if filter.GetRuntimeFilter() != nil {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "accessLoggingService runtimeFilter is not supported")
	}

	if filter.GetTraceableFilter() != nil {
		//unsure if this is correct. it appears this just needs to exist to function?
		filterType.TraceableFilter = true
	}

	if filter.GetStatusCodeFilter() != nil {
		filterType.StatusCodeFilter = &kgateway.StatusCodeFilter{
			Op:    kgateway.Op(filter.GetStatusCodeFilter().GetComparison().GetOp()),
			Value: filter.GetStatusCodeFilter().GetComparison().GetValue().GetDefaultValue(),
		}
	}

	return nil
}

// convertHTTPListenerOptions - generates GlooTrafficPolicy applied to the Gateway
// TODO(nick) - need to figure out which fields go to which policy. For example: httpConnectionManagerSettings.streamIdleTimeout: 3600s
func (o *GatewayAPIOutput) convertHTTPListenerOptions(options *gloov1.HttpListenerOptions, wrapper snapshot.Wrapper, proxyName string) {
	if options == nil {
		return
	}

	trafficPolicy := &gloogateway.GlooTrafficPolicy{
		TypeMeta: metav1.TypeMeta{
			Kind:       "GlooTrafficPolicy",
			APIVersion: gloogateway.GroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      wrapper.GetName(),
			Namespace: wrapper.GetNamespace(),
			Labels:    wrapper.GetLabels(),
		},
	}

	tps := gloogateway.GlooTrafficPolicySpec{
		TrafficPolicySpec: kgateway.TrafficPolicySpec{
			TargetRefs: []kgateway.LocalPolicyTargetReferenceWithSectionName{
				{
					LocalPolicyTargetReference: kgateway.LocalPolicyTargetReference{
						Group: gwv1.Group(gwv1.GroupVersion.Group),
						Kind:  "Gateway",
						Name:  gwv1.ObjectName(proxyName),
					},
				},
			},
			TargetSelectors: nil, // existing
			AI:              nil, // existing
			Transformation:  nil, // existing
			ExtProc:         nil, // existing
			ExtAuth:         nil, // existing
			RateLimit:       nil, // existing
		},
		Waf:                      nil, // existing
		Retry:                    nil, // existing
		Timeouts:                 nil, // existing
		RateLimitEnterprise:      nil, // existing
		ExtAuthEnterprise:        nil, // existing
		TransformationEnterprise: nil, // existing
	}

	// go through each option in Gateway Options and convert to listener policy

	// inline extAuth settings
	if options.GetExtauth() != nil {
		// These are global extAuthSettings that are also on the Settings Object.
		// If this exists we need to generate a GatewayExtensionObject for this
		gatewayExtensions := o.generateGatewayExtensionForExtAuth(options.Extauth, wrapper.GetName(), wrapper)
		o.gatewayAPICache.AddGatewayExtension(snapshot.NewGatewayExtensionWrapper(gatewayExtensions, wrapper.FileOrigin()))
	}

	// inline extProc settings
	if options.GetExtProc() != nil {
		gatewayExtensions := o.generateGatewayExtensionForExtProc(options.GetExtProc(), wrapper.GetName(), wrapper)
		o.gatewayAPICache.AddGatewayExtension(snapshot.NewGatewayExtensionWrapper(gatewayExtensions, wrapper.FileOrigin()))
	}

	// inline rate limit settings
	if options.GetRatelimitServer() != nil {
		gatewayExtensions := o.generateGatewayExtensionForRateLimit(options.GetRatelimitServer(), wrapper.GetName(), wrapper)
		o.gatewayAPICache.AddGatewayExtension(snapshot.NewGatewayExtensionWrapper(gatewayExtensions, wrapper.FileOrigin()))
	}

	if options.GetHttpLocalRatelimit() != nil {
		if options.GetHttpLocalRatelimit().GetEnableXRatelimitHeaders() != nil && options.GetHttpLocalRatelimit().GetEnableXRatelimitHeaders().GetValue() {
			o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "httpLocalRateLimit enableXRateLimitHeaders is not supported")
		}
		if options.GetHttpLocalRatelimit().GetLocalRateLimitPerDownstreamConnection() != nil && options.GetHttpLocalRatelimit().GetLocalRateLimitPerDownstreamConnection().GetValue() {
			o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "httpLocalRateLimit localRateLimitPerDownstreamConnection is not supported")
		}
		if options.GetHttpLocalRatelimit().GetDefaultLimit() != nil {
			rl := &kgateway.RateLimit{
				Local: &kgateway.LocalRateLimitPolicy{
					TokenBucket: &kgateway.TokenBucket{
						MaxTokens: options.GetHttpLocalRatelimit().GetDefaultLimit().GetMaxTokens(),
					},
				},
			}
			if options.GetHttpLocalRatelimit().GetDefaultLimit().GetTokensPerFill() != nil {
				rl.Local.TokenBucket.TokensPerFill = ptr.To(options.GetHttpLocalRatelimit().GetDefaultLimit().GetTokensPerFill().GetValue())
			}
			if options.GetHttpLocalRatelimit().GetDefaultLimit().GetFillInterval() != nil {
				rl.Local.TokenBucket.FillInterval = gwv1.Duration(options.GetHttpLocalRatelimit().GetDefaultLimit().GetFillInterval().AsDuration().String())
			}
			tps.TrafficPolicySpec.RateLimit = rl
		}
	}

	if options.GetWaf() != nil {
		waf := &gloogateway.Waf{
			Disabled:      ptr.To(options.GetWaf().GetDisabled()),
			CustomMessage: ptr.To(options.GetWaf().GetCustomInterventionMessage()),
			Rules:         []gloogateway.WafRule{},
		}
		for _, r := range options.GetWaf().RuleSets {
			waf.Rules = append(waf.Rules, gloogateway.WafRule{
				RuleStr: ptr.To(r.RuleStr),
			})
			if r.Files != nil && len(r.Files) > 0 {
				o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "WAF files is not supported")
			}
		}
		tps.Waf = waf
	}
	if options.GetDisableExtProc() != nil {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "disableExtProc is not supported")
	}
	if options.GetNetworkLocalRatelimit() != nil {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "networkLocalRateLimit is not supported")
	}
	if options.GetDlp() != nil {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "dlp is not supported")
	}
	if options.GetCsrf() != nil {
		csrf := o.convertCSRF(options.GetCsrf())
		tps.TrafficPolicySpec.Csrf = csrf
	}
	if options.GetBuffer() != nil {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "buffer is not supported")
	}
	if options.GetCaching() != nil {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "caching is not supported")
	}
	if options.GetConnectionLimit() != nil {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "connectionlimit is not supported")
	}
	if options.GetDynamicForwardProxy() != nil {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "dynamicForwardProxy (DFP) is not supported")
	}
	if options.GetExtensions() != nil {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "gloo edge extensions is not supported")
	}
	if options.GetGrpcJsonTranscoder() != nil {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "grpcToJson is not supported")
	}
	if options.GetGrpcWeb() != nil {
		//TODO(nick) : GRPCWeb is enabled by default in edge. we need to verify the same.
		//o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, glooGateway, "grpcWeb is not supported")
	}
	if options.GetGzip() != nil {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "gzip is not supported")
	}
	if options.GetHeaderValidationSettings() != nil {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "header validation is not supported")
	}
	if options.GetHealthCheck() != nil {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "health check is not supported")
	}
	if options.GetHttpConnectionManagerSettings() != nil {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "httpConnectionManagerSettings is not supported")
	}
	if options.GetProxyLatency() != nil {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "proxy latency is not supported")
	}
	if options.GetRouter() != nil {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "router (envoy filter maps) is not supported")
	}
	if options.GetSanitizeClusterHeader() != nil {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "sanitize cluster header is not supported")
	}
	if options.GetStatefulSession() != nil {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "statefulSession is not supported")
	}
	if options.GetTap() != nil {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "Tap filter is not supported")
	}
	if options.GetWasm() != nil {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "WASM is not supported")
	}
	trafficPolicy.Spec = tps

	o.gatewayAPICache.AddGlooTrafficPolicy(snapshot.NewGlooTrafficPolicyWrapper(trafficPolicy, wrapper.FileOrigin()))
}

func (o *GatewayAPIOutput) convertCSRF(policy *v4.CsrfPolicy) *kgateway.CSRFPolicy {
	csrf := &kgateway.CSRFPolicy{
		PercentageEnabled:  nil,
		PercentageShadowed: nil,
		AdditionalOrigins:  nil,
	}
	if policy.GetFilterEnabled() != nil {
		filterEnabled := policy.GetFilterEnabled()

		// Convert FractionalPercent to numerical percentage
		var percentage float64
		switch filterEnabled.GetDefaultValue().GetDenominator() {
		case v3.FractionalPercent_HUNDRED:
			percentage = float64(filterEnabled.GetDefaultValue().GetNumerator())
		case v3.FractionalPercent_TEN_THOUSAND:
			percentage = float64(filterEnabled.GetDefaultValue().GetNumerator()) / 100.0
		case v3.FractionalPercent_MILLION:
			percentage = float64(filterEnabled.GetDefaultValue().GetNumerator()) / 10000.0
		default:
			// Default to HUNDRED if denominator is not set
			percentage = float64(filterEnabled.GetDefaultValue().GetNumerator())
		}
		csrf.PercentageEnabled = ptr.To(uint32(percentage))
	}
	if policy.GetAdditionalOrigins() != nil {
		// Convert the additional origins from Gloo Edge format to kgateway format
		additionalOrigins := []*kgateway.StringMatcher{}
		for _, origin := range policy.GetAdditionalOrigins() {
			switch typed := origin.GetMatchPattern().(type) {
			case *gloo_type_matcher.StringMatcher_Exact:
				additionalOrigins = append(additionalOrigins, &kgateway.StringMatcher{Exact: ptr.To(typed.Exact)})
			case *gloo_type_matcher.StringMatcher_Prefix:
				additionalOrigins = append(additionalOrigins, &kgateway.StringMatcher{Prefix: ptr.To(typed.Prefix)})
			case *gloo_type_matcher.StringMatcher_Suffix:
				additionalOrigins = append(additionalOrigins, &kgateway.StringMatcher{Suffix: ptr.To(typed.Suffix)})
			case *gloo_type_matcher.StringMatcher_SafeRegex:
				additionalOrigins = append(additionalOrigins, &kgateway.StringMatcher{SafeRegex: ptr.To(typed.SafeRegex.GetRegex())})
			}
		}
		csrf.AdditionalOrigins = additionalOrigins
	}
	if policy.GetShadowEnabled() != nil {
		shadowEnabled := policy.GetShadowEnabled()

		// Convert FractionalPercent to numerical percentage
		var percentage float64
		switch shadowEnabled.GetDefaultValue().GetDenominator() {
		case v3.FractionalPercent_HUNDRED:
			percentage = float64(shadowEnabled.GetDefaultValue().GetNumerator())
		case v3.FractionalPercent_TEN_THOUSAND:
			percentage = float64(shadowEnabled.GetDefaultValue().GetNumerator()) / 100.0
		case v3.FractionalPercent_MILLION:
			percentage = float64(shadowEnabled.GetDefaultValue().GetNumerator()) / 10000.0
		default:
			// Default to HUNDRED if denominator is not set
			percentage = float64(shadowEnabled.GetDefaultValue().GetNumerator())
		}

		csrf.PercentageShadowed = ptr.To(uint32(percentage))
	}
	return csrf
}
func (o *GatewayAPIOutput) generateGatewayExtensionForExtProc(extProc *extproc.Settings, name string, wrapper snapshot.Wrapper) *kgateway.GatewayExtension {

	gatewayExtension := &kgateway.GatewayExtension{
		TypeMeta: metav1.TypeMeta{
			Kind:       "GatewayExtension",
			APIVersion: kgateway.GroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: wrapper.GetNamespace(),
			Labels:    wrapper.GetLabels(),
		},
		Spec: kgateway.GatewayExtensionSpec{
			Type:    kgateway.GatewayExtensionTypeExtProc,
			ExtProc: &kgateway.ExtProcProvider{},
		},
	}

	//TODO(nick): Implement ExtProc - https://github.com/kgateway-dev/kgateway/issues/11424
	if extProc.GetStatPrefix() != nil {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "extProc statPrefix is not supported")
	}
	if extProc.GetFailureModeAllow() != nil && extProc.GetFailureModeAllow().GetValue() == true {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "extProc failureModeAllow is not supported")
	}
	if extProc.GetAllowModeOverride() != nil && extProc.GetAllowModeOverride().GetValue() == true {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "extProc allowModeOverride is not supported")
	}
	if extProc.GetAsyncMode() != nil && extProc.GetAsyncMode().GetValue() == true {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "extProc asyncMode is not supported")
	}
	if extProc.GetDisableClearRouteCache() != nil && extProc.GetDisableClearRouteCache().GetValue() == true {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "extProc disableClearRouteCache is not supported")
	}
	if extProc.GetFilterMetadata() != nil {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "extProc filterMetadata is not supported")
	}
	if extProc.GetFilterStage() != nil {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "extProc filterStage is not supported")
	}
	if extProc.GetForwardRules() != nil {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "extProc forwardRules is not supported")
	}
	if extProc.GetMaxMessageTimeout() != nil {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "extProc maxMessageTimeout is not supported")
	}
	if extProc.GetMessageTimeout() != nil {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "extProc messageTimeout is not supported")
	}
	if extProc.GetMetadataContextNamespaces() != nil {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "extProc metadataContextNamespaces is not supported")
	}
	if extProc.GetMutationRules() != nil {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "extProc mutationRules is not supported")
	}
	if extProc.GetRequestAttributes() != nil {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "extProc requestAttributes is not supported")
	}
	if extProc.GetResponseAttributes() != nil {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "extProc responseAttributes is not supported")
	}
	if extProc.GetTypedMetadataContextNamespaces() != nil {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "extProc typedMetadataContextNamespaces is not supported")
	}
	if extProc.GetGrpcService() != nil {

		backend := o.GetGatewayAPICache().GetBackend(types.NamespacedName{Namespace: extProc.GetGrpcService().GetExtProcServerRef().GetNamespace(), Name: extProc.GetGrpcService().GetExtProcServerRef().GetName()})
		if backend == nil {
			o.AddErrorFromWrapper(ERROR_TYPE_UNKNOWN_REFERENCE, wrapper, "extProc grpcService backend not found")
		}
		grpcService := &kgateway.ExtGrpcService{
			BackendRef: &gwv1.BackendRef{
				BackendObjectReference: gwv1.BackendObjectReference{
					Name:      gwv1.ObjectName(backend.GetName()),
					Namespace: ptr.To(gwv1.Namespace(backend.GetNamespace())),
					// using the default port here
					Port: ptr.To(gwv1.PortNumber(4444)),
				},
			},
		}
		gatewayExtension.Spec.ExtProc.GrpcService = grpcService
	}
	return gatewayExtension
}
func (o *GatewayAPIOutput) generateGatewayExtensionForRateLimit(rateLimitSettings *ratelimit.Settings, name string, wrapper snapshot.Wrapper) *kgateway.GatewayExtension {

	gatewayExtension := &kgateway.GatewayExtension{
		TypeMeta: metav1.TypeMeta{
			Kind:       "GatewayExtension",
			APIVersion: kgateway.GroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: wrapper.GetNamespace(),
			Labels:    wrapper.GetLabels(),
		},
		Spec: kgateway.GatewayExtensionSpec{
			Type:      kgateway.GatewayExtensionTypeRateLimit,
			RateLimit: &kgateway.RateLimitProvider{},
		},
	}

	//TODO(nick): Implement RateLimitSettings - https://github.com/kgateway-dev/kgateway/issues/11424
	if rateLimitSettings.GetRateLimitBeforeAuth() == true {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "RateLimitSettings rateLimitBeforeAuth is not supported")

	}
	if rateLimitSettings.GetEnableXRatelimitHeaders() == true {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "RateLimitSettings enableXRatelimitHeaders is not supported")
	}
	if rateLimitSettings.GetDenyOnFail() == false {
		gatewayExtension.Spec.RateLimit.FailOpen = !rateLimitSettings.GetDenyOnFail()
	}
	if rateLimitSettings.GetRequestTimeout() != nil {
		gatewayExtension.Spec.RateLimit.Timeout = gwv1.Duration(rateLimitSettings.GetRequestTimeout().String())
	}
	if rateLimitSettings.GetRatelimitServerRef() != nil {
		backend := o.GetGatewayAPICache().GetBackend(types.NamespacedName{Namespace: rateLimitSettings.GetRatelimitServerRef().GetNamespace(), Name: rateLimitSettings.GetRatelimitServerRef().GetName()})
		if backend == nil {
			o.AddErrorFromWrapper(ERROR_TYPE_UNKNOWN_REFERENCE, wrapper, "rateLimitSettings ratelimitServerRef backend not found")
			return nil
		}

		grpcService := &kgateway.ExtGrpcService{
			BackendRef: &gwv1.BackendRef{
				BackendObjectReference: gwv1.BackendObjectReference{
					Name:      gwv1.ObjectName(backend.GetName()),
					Namespace: ptr.To(gwv1.Namespace(backend.GetNamespace())),
					// using the default port here
					Port: ptr.To(gwv1.PortNumber(18081)),
				},
			},
		}
		if rateLimitSettings.GetGrpcService() != nil {
			grpcService.Authority = ptr.To(rateLimitSettings.GetGrpcService().GetAuthority())
		}
		gatewayExtension.Spec.RateLimit.GrpcService = grpcService
	}

	return gatewayExtension
}
func (o *GatewayAPIOutput) generateGatewayExtensionForExtAuth(extauth *v1.Settings, name string, wrapper snapshot.Wrapper) *kgateway.GatewayExtension {
	if extauth == nil {
		return nil
	}
	var grpcService *kgateway.ExtGrpcService

	if extauth.GetExtauthzServerRef() != nil {
		backend := o.GetGatewayAPICache().GetBackend(types.NamespacedName{Namespace: extauth.GetExtauthzServerRef().GetNamespace(), Name: extauth.GetExtauthzServerRef().GetName()})
		if backend == nil {
			o.AddErrorFromWrapper(ERROR_TYPE_UNKNOWN_REFERENCE, wrapper, "extauth extauthzServerRef backend not found")
			return nil
		}
		grpcService = &kgateway.ExtGrpcService{
			BackendRef: &gwv1.BackendRef{
				BackendObjectReference: gwv1.BackendObjectReference{
					Name:      gwv1.ObjectName(backend.GetName()),
					Namespace: ptr.To(gwv1.Namespace(backend.GetNamespace())),
					// using the default port here
					Port: ptr.To(gwv1.PortNumber(8083)),
				},
			},
		}

		if extauth.GetGrpcService() != nil {
			grpcService.Authority = ptr.To(extauth.GetGrpcService().GetAuthority())
		}
	}
	//TODO(nick): Implement ExtAuthSettings - https://github.com/kgateway-dev/kgateway/issues/11424
	if extauth.GetClearRouteCache() == true {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "extAuth settings clearRouteCache is not supported")
	}
	if extauth.GetFailureModeAllow() == true {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "extAuth settings clearRouteCache is not supported")
	}
	if extauth.GetHttpService() != nil {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "extAuth settings httpService is not supported")
	}
	if extauth.GetRequestBody() != nil {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "extAuth settings requestBody is not supported")
	}
	if extauth.GetRequestTimeout() != nil {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "extAuth settings requestTimeout is not supported")
	}
	if extauth.GetStatPrefix() != "" {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "extAuth settings statPrefix is not supported")
	}
	if extauth.GetStatusOnError() != 0 {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "extAuth settings statusOnError is not supported")
	}
	if extauth.GetTransportApiVersion() != 0 {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "extAuth settings transportApiVersion is not supported")
	}
	if extauth.GetUserIdHeader() != "" {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "extAuth settings userIdHeader is not supported")
	}

	gatewayExtension := &kgateway.GatewayExtension{
		TypeMeta: metav1.TypeMeta{
			Kind:       "GatewayExtension",
			APIVersion: kgateway.GroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: wrapper.GetNamespace(),
			Labels:    wrapper.GetLabels(),
		},
		Spec: kgateway.GatewayExtensionSpec{
			Type: kgateway.GatewayExtensionTypeExtAuth,
			ExtAuth: &kgateway.ExtAuthProvider{
				GrpcService: grpcService,
			},
		},
	}
	return gatewayExtension
}

func (o *GatewayAPIOutput) convertVirtualServiceHTTPRoutes(vs *snapshot.VirtualServiceWrapper, glooGateway *snapshot.GlooGatewayWrapper, listenerName string) error {

	hr := &gwv1.HTTPRoute{
		TypeMeta: metav1.TypeMeta{
			Kind:       "HTTPRoute",
			APIVersion: gwv1.GroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      vs.GetName(),
			Namespace: vs.GetNamespace(),
			Labels:    vs.GetLabels(),
		},
		Spec: gwv1.HTTPRouteSpec{
			CommonRouteSpec: gwv1.CommonRouteSpec{
				ParentRefs: []gwv1.ParentReference{
					{
						Name:      gwv1.ObjectName(listenerName),
						Namespace: ptr.To(gwv1.Namespace(glooGateway.GetNamespace())),
						Kind:      ptr.To(gwv1.Kind("XListenerSet")),
						Group:     ptr.To(gwv1.Group(apixv1a1.GroupVersion.Group)),
					},
				},
			},
			Hostnames: convertDomains(vs.Spec.GetVirtualHost().GetDomains()),
			Rules:     []gwv1.HTTPRouteRule{},
		},
	}

	for _, route := range vs.Spec.GetVirtualHost().GetRoutes() {
		rule, err := o.convertRouteToRule(route, vs)
		if err != nil {
			return err
		}
		hr.Spec.Rules = append(hr.Spec.Rules, rule)
	}

	o.gatewayAPICache.AddHTTPRoute(snapshot.NewHTTPRouteWrapper(hr, vs.FileOrigin()))

	return nil
}

func (o *GatewayAPIOutput) convertRouteOptions(
	options *gloov1.RouteOptions,
	routeName string,
	wrapper snapshot.Wrapper,
) (*gloogateway.GlooTrafficPolicy, *gwv1.HTTPRouteFilter) {

	var trafficPolicy *gloogateway.GlooTrafficPolicy
	var filter *gwv1.HTTPRouteFilter
	associationID := RandStringRunes(RandomSuffix)
	if routeName == "" {
		routeName = "route-association"
	}
	associationName := fmt.Sprintf("%s-%s", routeName, associationID)

	if !isRouteOptionsSet(options) {
		return nil, nil
	}
	// converts options to RouteOptions but we need to this for everything except prefixrewrite and a few others now
	gtpSpec := gloogateway.GlooTrafficPolicySpec{
		TrafficPolicySpec: kgateway.TrafficPolicySpec{
			TargetRefs:      nil, // existing
			TargetSelectors: nil, // existing
			AI:              nil, // existing
			Transformation:  nil, // existing
			ExtProc:         nil, // existing
			ExtAuth:         nil, // existing
			RateLimit:       nil, // existing
		},
		Waf:                      nil, // existing
		Retry:                    nil, // existing
		Timeouts:                 nil, // existing
		RateLimitEnterprise:      nil, // existing
		ExtAuthEnterprise:        nil, // existing
		TransformationEnterprise: nil, // existing
	}

	//Features Supported By GatewayAPI
	// - RequestHeaderModifier
	// - ResponseHeaderModifier
	// - RequestRedirect
	// - URLRewrite
	// - Request Mirror
	// - CORS
	// - ExtensionRef
	// - Timeout (done)
	// - Retry (done)
	// - Session

	//// Because we move rewrites to a filter we need to remove it from RouteOptions
	// TODO(nick): delete this because this was for RouteOption and not needed for GlooTrafficPolicy we still need to add it to the HTTPRouteThough
	//if options.GetPrefixRewrite() != nil {
	//	trafficPolicy.Spec.GetOptions().PrefixRewrite = nil
	//}

	filter = &gwv1.HTTPRouteFilter{
		Type: gwv1.HTTPRouteFilterExtensionRef,
		ExtensionRef: &gwv1.LocalObjectReference{
			Group: glookube.GroupName,
			Kind:  "GlooTrafficPolicy",
			Name:  gwv1.ObjectName(associationName),
		},
	}
	if options.GetExtauth() != nil && options.GetExtauth().GetConfigRef() != nil {
		// we need to copy over the auth config ref if it exists
		ref := options.GetExtauth().GetConfigRef()
		ac, exists := o.edgeCache.AuthConfigs()[types.NamespacedName{Name: ref.GetName(), Namespace: ref.GetNamespace()}]
		if !exists {
			o.AddErrorFromWrapper(ERROR_TYPE_UNKNOWN_REFERENCE, wrapper, "did not find AuthConfig %s/%s for delegated route option reference", ref.GetName(), ref.GetNamespace())
		} else {
			o.gatewayAPICache.AddAuthConfig(ac)

			gtpSpec.ExtAuthEnterprise = &gloogateway.ExtAuthEnterprise{
				ExtensionRef: &corev1.LocalObjectReference{
					Name: "ext-authz",
				},
				AuthConfigRef: gloogateway.AuthConfigRef{
					Name:      ac.GetName(),
					Namespace: ptr.To(ac.GetNamespace()),
				},
			}
		}
	}
	if options.GetAi() != nil {
		aip := &kgateway.AIPolicy{
			PromptEnrichment: nil,
			PromptGuard:      nil,
			Defaults:         []kgateway.FieldDefault{},
		}
		switch options.GetAi().GetRouteType() {
		case ai.RouteSettings_CHAT:
			aip.RouteType = ptr.To(kgateway.CHAT)
		case ai.RouteSettings_CHAT_STREAMING:
			aip.RouteType = ptr.To(kgateway.CHAT_STREAMING)
		}
		for _, d := range options.GetAi().GetDefaults() {
			aip.Defaults = append(aip.Defaults, kgateway.FieldDefault{
				Field:    d.Field,
				Value:    d.Value.String(),
				Override: ptr.To(d.Override),
			})
		}
		if options.GetAi().GetPromptEnrichment() != nil {
			enrichment := &kgateway.AIPromptEnrichment{}

			for _, prepend := range options.GetAi().GetPromptEnrichment().GetPrepend() {
				enrichment.Prepend = append(enrichment.Prepend, kgateway.Message{
					Role:    prepend.GetRole(),
					Content: prepend.GetContent(),
				})
			}
			aip.PromptEnrichment = enrichment
		}
		if options.GetAi().GetRag() != nil {
			o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "ai RAG is not supported")
		}
		if options.GetAi().GetPromptGuard() != nil {
			guard := o.generateAIPromptGuard(options, wrapper)
			aip.PromptGuard = guard
		}
		if options.GetAi().GetSemanticCache() != nil {
			o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "ai SemanticCache is not supported")
		}
		gtpSpec.AI = aip
	}

	if options.GetWaf() != nil {
		// TODO(nick): Finish Implementing WAF -https://github.com/solo-io/gloo-gateway/issues/32
		gtpSpec.Waf = &gloogateway.Waf{
			Disabled:      ptr.To(options.GetWaf().GetDisabled()),
			Rules:         []gloogateway.WafRule{},
			CustomMessage: ptr.To(options.GetWaf().GetCustomInterventionMessage()),
		}
		for _, rule := range options.GetWaf().GetRuleSets() {
			gtpSpec.Waf.Rules = append(gtpSpec.Waf.Rules, gloogateway.WafRule{RuleStr: ptr.To(rule.GetRuleStr())})
			if rule.GetFiles() != nil {
				o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "WAF rule files is not supported")
			}
			if rule.GetDirectory() != "" {
				o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "WAF rule directory %s is not supported", rule.GetDirectory())
			}
		}
		if len(options.GetWaf().GetConfigMapRuleSets()) > 0 {
			o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "WAF configMapRuleSets is not supported")
		}
		if options.GetWaf().GetCoreRuleSet() != nil {
			o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "WAF coreRuleSets is not supported")
		}
		if options.GetWaf().GetRequestHeadersOnly() == true {
			o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "WAF requestHeadersOnly is not supported")
		}
		if options.GetWaf().GetResponseHeadersOnly() == true {
			o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "WAF responseHeadersOnly is not supported")
		}
		if options.GetWaf().GetAuditLogging() != nil {
			o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "WAF auditLogging is not supported")
		}
	}
	if options.GetCors() != nil {
		policy := o.convertCORS(options.GetCors(), wrapper)
		gtpSpec.Cors = policy
	}

	if options.GetRatelimit() != nil && len(options.GetRatelimit().GetRateLimits()) > 0 {

		rle := &gloogateway.RateLimitEnterprise{
			Global: &gloogateway.GlobalRateLimit{
				// Need to find the Gateway Extension for Global Rate Limit Server
				ExtensionRef: &corev1.LocalObjectReference{
					Name: "rate-limit",
				},

				RateLimits: []gloogateway.RateLimitActions{},
				// RateLimitConfig for the policy, not sure how it works for rate limit basic
				// TODO(nick) grab the global rate limit config ref
				RateLimitConfigRef: nil,
			},
		}
		for _, rl := range options.GetRatelimit().GetRateLimits() {
			rateLimit := &gloogateway.RateLimitActions{
				Actions:    []gloogateway.Action{},
				SetActions: []gloogateway.Action{},
			}
			for _, action := range rl.GetActions() {
				rateLimitAction := o.convertRateLimitAction(action)
				rateLimit.Actions = append(rateLimit.Actions, rateLimitAction)
			}
			for _, action := range rl.GetSetActions() {
				rateLimitAction := o.convertRateLimitAction(action)
				rateLimit.SetActions = append(rateLimit.SetActions, rateLimitAction)
			}
			if rl.GetLimit() != nil {
				o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "rateLimit action limit is not supported")
			}
		}
		gtpSpec.RateLimitEnterprise = rle
	}
	if options.GetRatelimitBasic() != nil {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "rateLimitBasic is not supported")

		//TODO (nick) : How do we translate rateLimitBasic to kgateway?
		//if options.GetRatelimitBasic().GetAuthorizedLimits() != nil {
		//
		//}
		//if options.GetRatelimitBasic().GetAnonymousLimits() != nil {
		//
		//}
		//gtpSpec.RateLimitEnterprise = &gloogateway.RateLimitEnterprise{
		//	Global: gloogateway.GlobalRateLimit{
		//		// Need to find the Gateway Extension for Global Rate Limit Server
		//		ExtensionRef: nil,
		//
		//		RateLimits: []gloogateway.RateLimitActions{},
		//		// RateLimitConfig for the policy, not sure how it works for rate limit basic
		//		RateLimitConfigRef: nil,
		//	},
		//}
	}
	if options.GetTransformations() != nil {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "legacy style transformation is not supported")
	}
	if options.GetStagedTransformations() != nil {
		transformation := o.convertStagedTransformation(options.GetStagedTransformations(), wrapper)
		gtpSpec.TransformationEnterprise = transformation
	}
	if options.GetDlp() != nil {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "dlp is not supported")
	}
	if options.GetCsrf() != nil {
		csrf := o.convertCSRF(options.GetCsrf())
		gtpSpec.TrafficPolicySpec.Csrf = csrf
	}
	if options.GetExtensions() != nil {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "gloo edge extensions is not supported")
	}
	if options.GetBufferPerRoute() != nil {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "bufferPerRoute is not supported")
	}
	if options.GetAppendXForwardedHost() != nil && options.GetAppendXForwardedHost().GetValue() == true {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "appendXForwardedHost is not supported")
	}
	if options.GetAutoHostRewrite() != nil {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "autoHostRewrite is not supported")
	}
	if options.GetEnvoyMetadata() != nil {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "envoyMetadata is not supported")
	}
	if options.GetFaults() != nil {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "faults is not supported")
	}
	if options.GetHostRewriteHeader() != nil {
		// TODO (nick): not sure how this is supported?
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "hostRewriteHeader is not supported")
	}
	if options.GetHostRewritePathRegex() != nil {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "hostRewritePathRegex is not supported")
	}
	if options.GetIdleTimeout() != nil {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "idleTimeout is not supported")
	}
	if options.GetJwtProvidersStaged() != nil {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "jwtProvidersStaged is not supported")
	}
	if options.GetJwtStaged() != nil {
		gtpSpec.JWTEnterprise = &gloogateway.StagedJWT{
			AfterExtAuth:  nil, // existing
			BeforeExtAuth: nil, // existing
		}
		if options.GetJwtStaged().GetBeforeExtAuth() != nil && options.GetJwtStaged().GetBeforeExtAuth().GetDisable() {
			gtpSpec.JWTEnterprise.BeforeExtAuth = &gloogateway.JWTEnterprise{
				Providers:        nil,                                                            // existing
				ValidationPolicy: nil,                                                            // existing
				Disable:          ptr.To(options.GetJwtStaged().GetBeforeExtAuth().GetDisable()), // existing
			}
		}
		if options.GetJwtStaged().GetAfterExtAuth() != nil && options.GetJwtStaged().GetAfterExtAuth().GetDisable() {
			gtpSpec.JWTEnterprise.AfterExtAuth = &gloogateway.JWTEnterprise{
				Providers:        nil,                                                           // existing
				ValidationPolicy: nil,                                                           // existing
				Disable:          ptr.To(options.GetJwtStaged().GetAfterExtAuth().GetDisable()), // existing
			}
		}
	}
	if options.GetLbHash() != nil {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "lbHash is not supported")
	}
	if options.GetMaxStreamDuration() != nil {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "maxStreamDuration is not supported")
	}
	if options.GetRbac() != nil {
		rbe := o.convertRBAC(options.GetRbac())
		gtpSpec.RBACEnterprise = rbe
	}
	if options.GetShadowing() != nil {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "shadowing is not supported")
	}
	if options.GetUpgrades() != nil {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "upgrades is not supported")
	}

	trafficPolicy = &gloogateway.GlooTrafficPolicy{
		TypeMeta: metav1.TypeMeta{
			Kind:       "GlooTrafficPolicy",
			APIVersion: gloogateway.GroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      associationName,
			Namespace: wrapper.GetNamespace(),
		},
		Spec: gtpSpec,
	}

	return trafficPolicy, filter
}

func (o *GatewayAPIOutput) convertRBAC(extension *rbac.ExtensionSettings) *gloogateway.RBACEnterprise {
	rbe := &gloogateway.RBACEnterprise{
		Disable:  ptr.To(extension.GetDisable()),
		Policies: map[string]gloogateway.RBACPolicy{},
	}
	for k, policy := range extension.GetPolicies() {
		rp := gloogateway.RBACPolicy{
			Principals:           make([]gloogateway.RBACPrincipal, 0),
			Permissions:          nil,
			NestedClaimDelimiter: ptr.To(policy.GetNestedClaimDelimiter()),
		}
		if policy.GetPermissions() != nil {
			rp.Permissions = &gloogateway.RBACPermissions{
				PathPrefix: ptr.To(policy.GetPermissions().GetPathPrefix()),
				Methods:    policy.GetPermissions().GetMethods(),
			}
		}

		for _, principle := range policy.Principals {
			if principle.GetJwtPrincipal() != nil {
				p := gloogateway.RBACPrincipal{
					JWTPrincipal: gloogateway.RBACJWTPrincipal{
						Claims:   principle.GetJwtPrincipal().GetClaims(),
						Provider: ptr.To(principle.GetJwtPrincipal().GetProvider()),
						Matcher:  nil,
					},
				}
				switch principle.GetJwtPrincipal().GetMatcher() {
				case rbac.JWTPrincipal_EXACT_STRING:
					p.JWTPrincipal.Matcher = ptr.To(gloogateway.JwtPrincipalClaimMatcherExactString)
				case rbac.JWTPrincipal_BOOLEAN:
					p.JWTPrincipal.Matcher = ptr.To(gloogateway.JwtPrincipalClaimMatcherBoolean)
				case rbac.JWTPrincipal_LIST_CONTAINS:
					p.JWTPrincipal.Matcher = ptr.To(gloogateway.JwtPrincipalClaimMatcherListContains)
				}
				rp.Principals = append(rp.Principals, p)
			}
		}
		rbe.Policies[k] = rp
	}
	return rbe
}

func (o *GatewayAPIOutput) convertRateLimitAction(action *v1alpha2.Action) gloogateway.Action {

	ggAction := gloogateway.Action{
		SourceCluster:      nil,
		DestinationCluster: nil,
		RequestHeaders:     nil,
		RemoteAddress:      nil,
		GenericKey:         nil,
		HeaderValueMatch:   nil,
		Metadata:           nil,
	}
	if action.GetSourceCluster() != nil {
		ggAction.SourceCluster = &gloogateway.SourceClusterAction{}
	}
	if action.GetDestinationCluster() != nil {
		ggAction.DestinationCluster = &gloogateway.DestinationClusterAction{}
	}
	if action.GetGenericKey() != nil {
		ggAction.GenericKey = &gloogateway.GenericKeyAction{
			DescriptorValue: action.GetGenericKey().GetDescriptorValue(),
		}
	}
	if action.GetHeaderValueMatch() != nil {
		hvm := &gloogateway.HeaderValueMatchAction{
			DescriptorValue: action.GetHeaderValueMatch().GetDescriptorValue(),
			ExpectMatch:     nil,
			Headers:         []gloogateway.HeaderMatcher{},
		}
		if action.GetHeaderValueMatch().GetExpectMatch() != nil {
			hvm.ExpectMatch = ptr.To(action.GetHeaderValueMatch().GetExpectMatch().GetValue())
		}
		for _, header := range action.GetHeaderValueMatch().GetHeaders() {
			var rangeMatch *gloogateway.Int64Range
			if header.GetRangeMatch() != nil {
				rangeMatch = &gloogateway.Int64Range{
					Start: header.GetRangeMatch().GetStart(),
					End:   header.GetRangeMatch().GetEnd(),
				}
			}
			//TODO(nick) this might set them all instead of the ones that exist
			hvm.Headers = append(hvm.Headers, gloogateway.HeaderMatcher{
				Name:         header.GetName(),
				ExactMatch:   ptr.To(header.GetExactMatch()),
				RegexMatch:   ptr.To(header.GetRegexMatch()),
				PresentMatch: ptr.To(header.GetPresentMatch()),
				PrefixMatch:  ptr.To(header.GetPrefixMatch()),
				SuffixMatch:  ptr.To(header.GetSuffixMatch()),
				InvertMatch:  ptr.To(header.GetInvertMatch()),
				RangeMatch:   rangeMatch,
			})
		}
		ggAction.HeaderValueMatch = hvm
	}
	return ggAction
}

func (o *GatewayAPIOutput) convertRequestTransformation(transformationRouting *transformation2.RequestResponseTransformations, wrapper snapshot.Wrapper) *gloogateway.RequestResponseTransformations {
	routing := &gloogateway.RequestResponseTransformations{}
	requestMatchers := o.convertRequestTransforms(transformationRouting.GetRequestTransforms(), wrapper)
	routing.Requests = requestMatchers

	responseMatchers := o.convertResponseTranforms(transformationRouting.GetResponseTransforms(), wrapper)
	routing.Responses = responseMatchers

	return routing
}
func (o *GatewayAPIOutput) convertResponseTranforms(responseTransform []*transformation2.ResponseMatch, wrapper snapshot.Wrapper) []gloogateway.ResponseMatcher {
	responseMatchers := []gloogateway.ResponseMatcher{}
	for _, rule := range responseTransform {
		match := gloogateway.ResponseMatcher{
			Headers:             []gloogateway.TransformationHeaderMatcher{},
			ResponseCodeDetails: ptr.To(rule.ResponseCodeDetails),
		}
		if rule.GetMatchers() != nil {
			for _, header := range rule.GetMatchers() {
				match.Headers = append(match.Headers, gloogateway.TransformationHeaderMatcher{
					Name:        header.GetName(),
					Value:       header.GetValue(),
					Regex:       header.GetRegex(),
					InvertMatch: header.GetInvertMatch(),
				})
			}
		}
		if rule.GetResponseTransformation() != nil {
			transformation := o.convertTransformationMatch(rule.GetResponseTransformation())
			match.Transformation = transformation
		}

		responseMatchers = append(responseMatchers, match)
	}
	return responseMatchers
}
func (o *GatewayAPIOutput) convertRequestTransforms(requestTranforms []*transformation2.RequestMatch, wrapper snapshot.Wrapper) []gloogateway.RequestMatcher {
	requestMatchers := []gloogateway.RequestMatcher{}
	for _, rule := range requestTranforms {
		match := gloogateway.RequestMatcher{}
		if rule.GetClearRouteCache() == true {
			o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "transformation rule clearRouteCache is not supported")
		}
		if rule.GetMatcher() != nil {
			match.Matcher = &gloogateway.TransformationRequestMatcher{
				Headers: []gloogateway.TransformationHeaderMatcher{},
			}
			for _, header := range rule.GetMatcher().GetHeaders() {
				match.Matcher.Headers = append(match.Matcher.Headers, gloogateway.TransformationHeaderMatcher{
					Name:        header.GetName(),
					Value:       header.GetValue(),
					Regex:       header.GetRegex(),
					InvertMatch: header.GetInvertMatch(),
				})
			}
			if rule.GetMatcher().GetConnectMatcher() != nil {
				o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "transformation rule connect match is not supported")
			}
			if rule.GetMatcher().GetCaseSensitive() != nil {
				o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "transformation rule caseSensitive match is not supported")
			}
			if rule.GetMatcher().GetExact() != "" {
				o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "transformation rule exact match is not supported")
			}
			if rule.GetMatcher().GetMethods() != nil {
				o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "transformation rule methods match is not supported")
			}
			if rule.GetMatcher().GetPrefix() != "" {
				o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "transformation rule prefix match is not supported")
			}
			if rule.GetMatcher().GetQueryParameters() != nil {
				o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "transformation rule queryParameters match is not supported")
			}
		}
		if rule.GetRequestTransformation() != nil {
			transformation := o.convertTransformationMatch(rule.GetRequestTransformation())
			match.Transformation = transformation
		}

		requestMatchers = append(requestMatchers, match)
	}
	return requestMatchers
}

func (o *GatewayAPIOutput) convertTransformationMatch(rule *transformation2.Transformation) gloogateway.Transformation {
	transformation := gloogateway.Transformation{
		Template:   nil, // existing
		HeaderBody: nil, // existing
	}

	// TODO fill this out and look for more options on Gloo edge transformation template
	if rule.GetTransformationTemplate() != nil {
		tt := rule.GetTransformationTemplate()
		template := &gloogateway.TransformationTemplate{
			AdvancedTemplates:     ptr.To(rule.GetTransformationTemplate().AdvancedTemplates),
			Extractors:            nil,
			Headers:               nil,
			HeadersToAppend:       []gloogateway.HeaderToAppend{},
			HeadersToRemove:       []string{},
			BodyTransformation:    nil,
			ParseBodyBehavior:     nil,
			IgnoreErrorOnParse:    ptr.To(tt.GetIgnoreErrorOnParse()),
			DynamicMetadataValues: []gloogateway.DynamicMetadataValue{},
			EscapeCharacters:      nil,
			SpanTransformer:       nil,
		}
		for name, ext := range tt.GetExtractors() {
			if template.Extractors == nil {
				template.Extractors = map[string]*gloogateway.Extraction{}
			}
			extraction := &gloogateway.Extraction{
				ExtractionHeader: ptr.To(ext.GetHeader()),
				Regex:            ext.GetRegex(),
				Subgroup:         ptr.To(ext.GetSubgroup()),
			}
			if ext.GetBody() != nil {
				extraction.ExtractionBody = ptr.To(true)
			}
			if ext.GetReplacementText() != nil {
				extraction.ReplacementText = ptr.To(ext.GetReplacementText().Value)
			}
			switch ext.GetMode() {
			case transformation2.Extraction_EXTRACT:
				extraction.Mode = ptr.To(gloogateway.ModeExtract)
			case transformation2.Extraction_SINGLE_REPLACE:
				extraction.Mode = ptr.To(gloogateway.ModeSingleReplace)
			case transformation2.Extraction_REPLACE_ALL:
				extraction.Mode = ptr.To(gloogateway.ModeReplaceAll)
			}
			template.Extractors[name] = extraction
		}
		for name, header := range tt.GetHeaders() {
			if template.Headers == nil {
				template.Headers = make(map[string]gloogateway.InjaTemplate)
			}
			template.Headers[name] = gloogateway.InjaTemplate(header.GetText())
		}
		for _, hta := range tt.GetHeadersToAppend() {
			h := gloogateway.HeaderToAppend{
				Key: hta.Key,
			}
			if hta.Value != nil {
				h.Value = gloogateway.InjaTemplate(hta.Value.String())
			}
			template.HeadersToAppend = append(template.HeadersToAppend, h)
		}
		for _, htr := range tt.GetHeadersToRemove() {
			template.HeadersToRemove = append(template.HeadersToRemove, htr)
		}
		if tt.GetBody() != nil {
			template.BodyTransformation = &gloogateway.BodyTransformation{
				Type: gloogateway.BodyTransformationTypeBody,
				Body: ptr.To(gloogateway.InjaTemplate(tt.GetBody().String())),
			}
		}
		if tt.GetPassthrough() != nil {
			template.BodyTransformation = &gloogateway.BodyTransformation{
				Type: gloogateway.BodyTransformationTypePassthrough,
				Body: ptr.To(gloogateway.InjaTemplate(tt.GetPassthrough().String())),
			}
		}
		if tt.GetMergeExtractorsToBody() != nil {
			template.BodyTransformation = &gloogateway.BodyTransformation{
				Type: gloogateway.BodyTransformationTypeMergeExtractorsToBody,
				Body: ptr.To(gloogateway.InjaTemplate(tt.GetMergeExtractorsToBody().String())),
			}
		}
		if tt.GetMergeJsonKeys() != nil {
			template.BodyTransformation = &gloogateway.BodyTransformation{
				Type: gloogateway.BodyTransformationTypeMergeJsonKeys,
				Body: ptr.To(gloogateway.InjaTemplate(tt.GetMergeJsonKeys().String())),
			}
		}
		if tt.GetParseBodyBehavior() == transformation2.TransformationTemplate_ParseAsJson {
			template.ParseBodyBehavior = ptr.To(gloogateway.ParseAsJson)
		}
		if tt.GetParseBodyBehavior() == transformation2.TransformationTemplate_DontParse {
			template.ParseBodyBehavior = ptr.To(gloogateway.DontParse)
		}
		for _, m := range tt.GetDynamicMetadataValues() {
			dm := gloogateway.DynamicMetadataValue{
				MetadataNamespace: ptr.To(m.GetMetadataNamespace()),
				Key:               m.GetKey(),
				Value:             gloogateway.InjaTemplate(m.GetValue().String()),
				JsonToProto:       ptr.To(m.JsonToProto),
			}
			if m.GetValue() != nil {
				dm.Value = gloogateway.InjaTemplate(m.GetValue().String())
			}
			template.DynamicMetadataValues = append(template.DynamicMetadataValues, dm)
		}
		if tt.GetEscapeCharacters() != nil {
			if tt.GetEscapeCharacters().GetValue() {
				template.EscapeCharacters = ptr.To(gloogateway.EscapeCharactersEscape)
			}
		}
		if tt.GetSpanTransformer() != nil && tt.GetSpanTransformer().GetName() != nil {
			template.SpanTransformer = &gloogateway.SpanTransformer{
				Name: gloogateway.InjaTemplate(tt.GetSpanTransformer().GetName().GetText()),
			}
		}

		transformation.Template = template
	}
	return transformation
}

func (o *GatewayAPIOutput) generateAIPromptGuard(options *gloov1.RouteOptions, wrapper snapshot.Wrapper) *kgateway.AIPromptGuard {
	guard := &kgateway.AIPromptGuard{
		Request:  nil,
		Response: nil,
	}
	if options.GetAi().GetPromptGuard().GetRequest() != nil {
		request := o.convertPromptGuardRequest(options, wrapper)
		guard.Request = request
	}
	if options.GetAi().GetPromptGuard().GetResponse() != nil {
		response := o.convertPromptGuardResponse(options, wrapper)
		guard.Response = response
	}
	return guard
}

func (o *GatewayAPIOutput) convertPromptGuardResponse(options *gloov1.RouteOptions, wrapper snapshot.Wrapper) *kgateway.PromptguardResponse {
	response := &kgateway.PromptguardResponse{
		Regex:   nil, // existing
		Webhook: nil, // existing
	}

	if options.GetAi().GetPromptGuard().GetResponse().GetWebhook() != nil {
		webhook := &kgateway.Webhook{
			Host: kgateway.Host{
				Host: options.GetAi().GetPromptGuard().GetResponse().GetWebhook().GetHost(),
				Port: gwv1.PortNumber(options.GetAi().GetPromptGuard().GetResponse().GetWebhook().GetPort()),
				//InsecureSkipVerify: nil,
			},
			ForwardHeaders: []gwv1.HTTPHeaderMatch{},
		}
		for _, h := range options.GetAi().GetPromptGuard().GetResponse().GetWebhook().GetForwardHeaders() {
			match := gwv1.HTTPHeaderMatch{
				Name: gwv1.HTTPHeaderName(h.GetKey()),
				//Value: nil,
			}
			// TODO(nick) - We have a lot of options but gateway API only has exact or regex....
			switch h.GetMatchType() {
			case ai.AIPromptGuard_Webhook_HeaderMatch_CONTAINS:
				match.Type = ptr.To(gwv1.HeaderMatchExact)
				o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "ai headerMatch 'contains' is not supported")
			case ai.AIPromptGuard_Webhook_HeaderMatch_EXACT:
				match.Type = ptr.To(gwv1.HeaderMatchExact)
			case ai.AIPromptGuard_Webhook_HeaderMatch_PREFIX:
				o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "ai headerMatch 'prefix' is not supported")
				match.Type = ptr.To(gwv1.HeaderMatchExact)
			case ai.AIPromptGuard_Webhook_HeaderMatch_REGEX:
				match.Type = ptr.To(gwv1.HeaderMatchRegularExpression)
			case ai.AIPromptGuard_Webhook_HeaderMatch_SUFFIX:
				o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "ai headerMatch 'suffix' is not supported")
				match.Type = ptr.To(gwv1.HeaderMatchExact)
			}
			webhook.ForwardHeaders = append(webhook.ForwardHeaders, match)
		}
		response.Webhook = webhook
	}

	if options.GetAi().GetPromptGuard().GetResponse().GetRegex() != nil {
		response.Regex = &kgateway.Regex{
			Matches:  []kgateway.RegexMatch{},
			Builtins: []kgateway.BuiltIn{},
		}
		switch options.GetAi().GetPromptGuard().GetResponse().GetRegex().GetAction() {
		case ai.AIPromptGuard_Regex_MASK:
			response.Regex.Action = ptr.To(kgateway.MASK)
		case ai.AIPromptGuard_Regex_REJECT:
			response.Regex.Action = ptr.To(kgateway.REJECT)
		}

		for _, match := range options.GetAi().GetPromptGuard().GetResponse().GetRegex().GetMatches() {
			response.Regex.Matches = append(response.Regex.Matches, kgateway.RegexMatch{
				Pattern: ptr.To(match.GetPattern()),
				Name:    ptr.To(match.GetName()),
			})
		}
		response.Regex.Builtins = []kgateway.BuiltIn{}
		for _, builtIns := range options.GetAi().GetPromptGuard().GetResponse().GetRegex().GetBuiltins() {
			switch builtIns {
			case ai.AIPromptGuard_Regex_SSN:
				response.Regex.Builtins = append(response.Regex.Builtins, kgateway.SSN)
			case ai.AIPromptGuard_Regex_CREDIT_CARD:
				response.Regex.Builtins = append(response.Regex.Builtins, kgateway.CREDIT_CARD)
			case ai.AIPromptGuard_Regex_PHONE_NUMBER:
				response.Regex.Builtins = append(response.Regex.Builtins, kgateway.PHONE_NUMBER)
			}
		}
	}
	return response
}

func (o *GatewayAPIOutput) convertPromptGuardRequest(options *gloov1.RouteOptions, wrapper snapshot.Wrapper) *kgateway.PromptguardRequest {
	request := &kgateway.PromptguardRequest{
		CustomResponse: &kgateway.CustomResponse{
			Message:    ptr.To(options.GetAi().GetPromptGuard().GetRequest().GetCustomResponse().GetMessage()),
			StatusCode: ptr.To(options.GetAi().GetPromptGuard().GetRequest().GetCustomResponse().GetStatusCode()),
		},
	}
	if options.GetAi().GetPromptGuard().GetRequest().GetModeration() != nil && options.GetAi().GetPromptGuard().GetRequest().GetModeration().GetOpenai() != nil {
		request.Moderation = &kgateway.Moderation{
			OpenAIModeration: &kgateway.OpenAIConfig{
				Model: ptr.To(options.GetAi().GetPromptGuard().GetRequest().GetModeration().GetOpenai().GetModel()),
			},
		}
		if options.GetAi().GetPromptGuard().GetRequest().GetModeration().GetOpenai().GetAuthToken() != nil {
			authToken := kgateway.SingleAuthToken{}
			if options.GetAi().GetPromptGuard().GetRequest().GetModeration().GetOpenai().GetAuthToken().GetInline() != "" {
				authToken.Kind = kgateway.Inline
				authToken.Inline = ptr.To(options.GetAi().GetPromptGuard().GetRequest().GetModeration().GetOpenai().GetAuthToken().GetInline())
			}
			if options.GetAi().GetPromptGuard().GetRequest().GetModeration().GetOpenai().GetAuthToken().GetSecretRef() != nil {
				authToken.Kind = kgateway.SecretRef
				authToken.SecretRef = &corev1.LocalObjectReference{
					Name: options.GetAi().GetPromptGuard().GetRequest().GetModeration().GetOpenai().GetAuthToken().GetSecretRef().GetName(),
				}
				if options.GetAi().GetPromptGuard().GetRequest().GetModeration().GetOpenai().GetAuthToken().GetSecretRef().GetNamespace() != wrapper.GetNamespace() {
					o.AddErrorFromWrapper(ERROR_TYPE_UNKNOWN_REFERENCE, wrapper, "AI AuthToken secretRef may be referencing secret outside configuration namespace")
				}
			}
			request.Moderation.OpenAIModeration.AuthToken = authToken
		}
	}

	if options.GetAi().GetPromptGuard().GetRequest().GetWebhook() != nil {
		webhook := &kgateway.Webhook{
			Host: kgateway.Host{
				Host: options.GetAi().GetPromptGuard().GetRequest().GetWebhook().GetHost(),
				Port: gwv1.PortNumber(options.GetAi().GetPromptGuard().GetRequest().GetWebhook().GetPort()),
				//InsecureSkipVerify: nil,
			},
			ForwardHeaders: []gwv1.HTTPHeaderMatch{},
		}
		for _, h := range options.GetAi().GetPromptGuard().GetRequest().GetWebhook().GetForwardHeaders() {
			match := gwv1.HTTPHeaderMatch{
				Name: gwv1.HTTPHeaderName(h.GetKey()),
				//Value: nil,
			}
			// TODO(nick) - We have a lot of options but gateway API only has exact or regex....
			switch h.GetMatchType() {
			case ai.AIPromptGuard_Webhook_HeaderMatch_CONTAINS:
				match.Type = ptr.To(gwv1.HeaderMatchExact)
				o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "ai headerMatch 'contains' is not supported")
			case ai.AIPromptGuard_Webhook_HeaderMatch_EXACT:
				match.Type = ptr.To(gwv1.HeaderMatchExact)
			case ai.AIPromptGuard_Webhook_HeaderMatch_PREFIX:
				o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "ai headerMatch 'prefix' is not supported")
				match.Type = ptr.To(gwv1.HeaderMatchExact)
			case ai.AIPromptGuard_Webhook_HeaderMatch_REGEX:
				match.Type = ptr.To(gwv1.HeaderMatchRegularExpression)
			case ai.AIPromptGuard_Webhook_HeaderMatch_SUFFIX:
				o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "ai headerMatch 'suffix' is not supported")
				match.Type = ptr.To(gwv1.HeaderMatchExact)
			}
			webhook.ForwardHeaders = append(webhook.ForwardHeaders, match)
		}
		request.Webhook = webhook
	}

	if options.GetAi().GetPromptGuard().GetRequest().GetCustomResponse() != nil {
		request.CustomResponse = &kgateway.CustomResponse{
			Message:    ptr.To(options.GetAi().GetPromptGuard().GetRequest().GetCustomResponse().GetMessage()),
			StatusCode: ptr.To(options.GetAi().GetPromptGuard().GetRequest().GetCustomResponse().GetStatusCode()),
		}
	}
	if options.GetAi().GetPromptGuard().GetRequest().GetRegex() != nil {
		request.Regex = &kgateway.Regex{
			Matches:  []kgateway.RegexMatch{},
			Builtins: []kgateway.BuiltIn{},
		}
		switch options.GetAi().GetPromptGuard().GetRequest().GetRegex().GetAction() {
		case ai.AIPromptGuard_Regex_MASK:
			request.Regex.Action = ptr.To(kgateway.MASK)
		case ai.AIPromptGuard_Regex_REJECT:
			request.Regex.Action = ptr.To(kgateway.REJECT)
		}

		for _, match := range options.GetAi().GetPromptGuard().GetRequest().GetRegex().GetMatches() {
			request.Regex.Matches = append(request.Regex.Matches, kgateway.RegexMatch{
				Pattern: ptr.To(match.GetPattern()),
				Name:    ptr.To(match.GetName()),
			})
		}
		request.Regex.Builtins = []kgateway.BuiltIn{}
		for _, builtIns := range options.GetAi().GetPromptGuard().GetRequest().GetRegex().GetBuiltins() {
			switch builtIns {
			case ai.AIPromptGuard_Regex_SSN:
				request.Regex.Builtins = append(request.Regex.Builtins, kgateway.SSN)
			case ai.AIPromptGuard_Regex_CREDIT_CARD:
				request.Regex.Builtins = append(request.Regex.Builtins, kgateway.CREDIT_CARD)
			case ai.AIPromptGuard_Regex_PHONE_NUMBER:
				request.Regex.Builtins = append(request.Regex.Builtins, kgateway.PHONE_NUMBER)
			}
		}
	}
	return request
}

func (o *GatewayAPIOutput) convertRouteToRule(r *gloogwv1.Route, wrapper snapshot.Wrapper) (gwv1.HTTPRouteRule, error) {

	rr := gwv1.HTTPRouteRule{
		Name:               nil, //existing
		Matches:            []gwv1.HTTPRouteMatch{},
		Filters:            []gwv1.HTTPRouteFilter{},
		BackendRefs:        []gwv1.HTTPBackendRef{},
		Timeouts:           nil, //existing
		Retry:              nil, //existing
		SessionPersistence: nil, //existing
	}

	// unused fields
	if r.GetInheritablePathMatchers() != nil && r.GetInheritablePathMatchers().GetValue() == true {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "has inheritable path matchers but there is not equivalent in Gateway API")
	}
	if r.GetInheritableMatchers() != nil && r.GetInheritableMatchers().GetValue() == true {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "has inheritable matchers but there is not equivalent in Gateway API")
	}

	for _, m := range r.GetMatchers() {
		match, err := o.convertMatch(m, wrapper)
		if err != nil {
			return rr, err
		}
		rr.Matches = append(rr.Matches, match)
	}
	if r.GetOptions() != nil {
		options := r.GetOptions()

		// TODO we might still want to do all of these in TP or GTP due to them potentially applying at the listener or gateway level and having more features.
		// Features Supported By GatewayAPI
		// - RequestHeaderModifier
		// - ResponseHeaderModifier
		// - RequestRedirect
		// - URLRewrite
		// - Request Mirror
		// - CORS
		// - ExtensionRef
		// - Timeout (done)
		// - Retry (done)
		// - Session

		// prefix rewrite, sets it on HTTPRoute
		if options.GetPrefixRewrite() != nil {
			rf := o.generateFilterForURLRewrite(r, wrapper)
			if rf != nil {
				rr.Filters = append(rr.Filters, *rf)
			}
		}

		if options.GetTimeout() != nil {
			rr.Timeouts = &gwv1.HTTPRouteTimeouts{
				Request: ptr.To(gwv1.Duration(options.GetTimeout().AsDuration().String())),
			}
		}
		if options.GetRetries() != nil {
			retry := &gwv1.HTTPRouteRetry{
				Codes:    []gwv1.HTTPRouteRetryStatusCode{},
				Attempts: ptr.To(int(options.GetRetries().GetNumRetries())),
				Backoff:  nil,
			}
			if options.GetRetries().GetRetryOn() != "" {
				// TODO need to convert envoy x-envoy-retry-on to HTTPRouteRetry
				o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "retry does not support x-envoy-retry-on")
			}
			if options.GetRetries().GetPreviousPriorities() != nil {
				o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "retry does not support envoy previous priorities retry selector")
			}

			if options.GetRetries().GetPerTryTimeout() != nil {
				retry.Backoff = ptr.To(gwv1.Duration(options.GetRetries().GetPerTryTimeout().String()))
			}
			if options.GetRetries().GetRetriableStatusCodes() != nil {
				for _, code := range options.GetRetries().GetRetriableStatusCodes() {
					retry.Codes = append(retry.Codes, gwv1.HTTPRouteRetryStatusCode(code))
				}
			}

			rr.Retry = retry
		}

		glooTrafficPolicy, filter := o.convertRouteOptions(options, r.GetName(), wrapper)
		if filter != nil {
			rr.Filters = append(rr.Filters, *filter)
		}
		if glooTrafficPolicy != nil {
			o.gatewayAPICache.AddGlooTrafficPolicy(snapshot.NewGlooTrafficPolicyWrapper(glooTrafficPolicy, wrapper.FileOrigin()))
		}
	}
	// Process Route_Actions
	if r.GetRouteAction() != nil {
		// Route_Route Action
		if r.GetRouteAction().GetClusterHeader() != "" {
			o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "has cluster header action set but there is not equivalent in Gateway API")
		}
		if r.GetRouteAction().GetDynamicForwardProxy() != nil {
			o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "has dynamic forward proxy action set but there is not equivalent in Gateway API")
		}
		if r.GetRouteAction().GetMulti() != nil {
			o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "has multi detination action set but there is not equivalent in Gateway API")
		}

		if r.GetRouteAction().GetSingle() != nil {
			// single static upstream
			if r.GetRouteAction().GetSingle().GetUpstream() != nil {
				backendRef := o.generateBackendRefForSingleUpstream(r, wrapper)

				rr.BackendRefs = append(rr.BackendRefs, backendRef)
			}
		}
		if r.GetRouteAction().GetUpstreamGroup() != nil {
			o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "has upstream group action set but there is not equivalent in Gateway API")
		}

	} else if r.GetRedirectAction() != nil {
		rdf := o.convertRedirect(r, wrapper)

		rr.Filters = append(rr.Filters, gwv1.HTTPRouteFilter{
			Type:            "RequestRedirect",
			RequestRedirect: rdf,
		})

	} else if r.GetDirectResponseAction() != nil {

		dr := convertDirectResponse(r.GetDirectResponseAction())
		if dr != nil {
			// TODO(nick): what if route name is nil?
			rName := r.GetName()
			if rName == "" {
				rName = RandStringRunes(6)
			}
			drName := fmt.Sprintf("directresponse-%s-%s", wrapper.GetName(), rName)
			dr.Name = drName
			dr.Namespace = wrapper.GetNamespace()
			o.gatewayAPICache.AddDirectResponse(snapshot.NewDirectResponseWrapper(dr, wrapper.FileOrigin()))

			rr.Filters = append(rr.Filters, gwv1.HTTPRouteFilter{
				Type: "ExtensionRef",
				ExtensionRef: &gwv1.LocalObjectReference{
					Group: v1alpha1.Group,
					Kind:  "DirectResponse",
					Name:  gwv1.ObjectName(drName),
				},
			})
		}

	} else if r.GetDelegateAction() != nil {
		// delegate action
		// intermediate delegation step. This is a placeholder for the next path to do delegation
		backendRef := o.generateBackendRefForDelegateAction(r, wrapper)

		if len(backendRef) > 0 {
			for _, b := range backendRef {
				rr.BackendRefs = append(rr.BackendRefs, *b)
			}
		}
	}

	if r.GetOptionsConfigRefs() != nil && len(r.GetOptionsConfigRefs().GetDelegateOptions()) > 0 {
		// these are references to other RouteOptions, we need to add them
		for _, delegateOptions := range r.GetOptionsConfigRefs().GetDelegateOptions() {
			if delegateOptions.GetNamespace() != wrapper.GetNamespace() {
				o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "delegates to route options not in same namespace (this does not work in Gateway API)")
			}

			// grab that route option and convert it to GlooTrafficPolicy
			ro, exists := o.edgeCache.RouteOptions()[types.NamespacedName{Name: delegateOptions.GetName(), Namespace: delegateOptions.GetNamespace()}]
			if !exists {
				o.AddErrorFromWrapper(ERROR_TYPE_UNKNOWN_REFERENCE, wrapper, "did not find RouteOption %s/%s for delegated route option reference", delegateOptions.GetNamespace(), delegateOptions.GetName())
			}

			if ro.Spec.GetOptions() != nil && ro.Spec.GetOptions().GetExtauth() != nil && ro.Spec.GetOptions().GetExtauth().GetConfigRef() != nil {
				// we need to copy over the auth config ref if it exists
				ref := ro.Spec.GetOptions().GetExtauth().GetConfigRef()
				ac, exists := o.edgeCache.AuthConfigs()[types.NamespacedName{Name: ref.GetName(), Namespace: ref.GetNamespace()}]
				if !exists {
					o.AddErrorFromWrapper(ERROR_TYPE_UNKNOWN_REFERENCE, ro, "did not find AuthConfig %s/%s for delegated route option reference", ref.GetName(), ref.GetNamespace())
				}
				o.gatewayAPICache.AddAuthConfig(ac)
			}

			gtp, filter := o.convertRouteOptions(ro.RouteOption.Spec.GetOptions(), delegateOptions.GetName(), ro)
			if gtp != nil {
				o.gatewayAPICache.AddGlooTrafficPolicy(snapshot.NewGlooTrafficPolicyWrapper(gtp, ro.FileOrigin()))
			}
			if filter != nil {
				rr.Filters = append(rr.Filters, *filter)
			}
		}
	}

	return rr, nil
}

func (o *GatewayAPIOutput) convertRedirect(r *gloogwv1.Route, wrapper snapshot.Wrapper) *gwv1.HTTPRequestRedirectFilter {
	rdf := &gwv1.HTTPRequestRedirectFilter{}

	action := r.GetRedirectAction()
	if action.GetHttpsRedirect() {
		rdf.Scheme = ptr.To("https")
	}
	if action.GetStripQuery() {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "has stripQuery redirect action but there is not equivalent in Gateway API")
	}
	if action.GetRegexRewrite() != nil {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "has regexRewrite redirect action but there is not equivalent in Gateway API")
	}
	if action.GetPrefixRewrite() != "" {
		match, err := isPrefixMatch(r.GetMatchers())
		if err != nil {
			o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "has multiple matchers in same route")
		}

		if match {
			// full path rewrite
			rdf.Path = &gwv1.HTTPPathModifier{
				Type:               gwv1.PrefixMatchHTTPPathModifier,
				ReplacePrefixMatch: ptr.To(action.GetPrefixRewrite()),
			}
		}

	}
	if action.GetHostRedirect() != "" {
		rdf.Hostname = ptr.To(gwv1.PreciseHostname(action.GetHostRedirect()))
	}
	if action.GetPathRedirect() != "" {
		match, err := isExactMatch(r.GetMatchers())
		if err != nil {
			o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "has multiple matchers in same route")
		}
		if match {
			// full path rewrite
			rdf.Path = &gwv1.HTTPPathModifier{
				Type:            gwv1.FullPathHTTPPathModifier,
				ReplaceFullPath: ptr.To(action.GetPathRedirect()),
			}
		}
	}

	if action.GetPortRedirect() != nil {
		rdf.Port = ptr.To(gwv1.PortNumber(action.GetPortRedirect().GetValue()))
	}

	switch action.GetResponseCode() {
	case gloov1.RedirectAction_MOVED_PERMANENTLY:
		rdf.StatusCode = ptr.To(301)
	case gloov1.RedirectAction_FOUND:
		rdf.StatusCode = ptr.To(302)
	case gloov1.RedirectAction_SEE_OTHER:
		rdf.StatusCode = ptr.To(303)
	case gloov1.RedirectAction_TEMPORARY_REDIRECT:
		rdf.StatusCode = ptr.To(307)
	case gloov1.RedirectAction_PERMANENT_REDIRECT:
		rdf.StatusCode = ptr.To(308)
	default:
		rdf.StatusCode = ptr.To(301)
	}
	return rdf
}
func convertDirectResponse(action *gloov1.DirectResponseAction) *kgateway.DirectResponse {
	if action == nil {
		return nil
	}
	dr := &kgateway.DirectResponse{
		TypeMeta: metav1.TypeMeta{
			Kind:       "DirectResponse",
			APIVersion: kgateway.GroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{},
		Spec: kgateway.DirectResponseSpec{
			StatusCode: action.GetStatus(),
			Body:       action.GetBody(),
		},
	}

	return dr
}

func (o *GatewayAPIOutput) generateBackendRefForDelegateAction(
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
				o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "has delegate action with more than one label selector which is not supported in Gateway API")
			}
			// create a backend ref for each label
			for _, v := range selector.GetLabels() {
				// just grab the first label
				backendRef := &gwv1.HTTPBackendRef{
					BackendRef: gwv1.BackendRef{
						BackendObjectReference: gwv1.BackendObjectReference{
							Name:      gwv1.ObjectName(v),                               // label value
							Namespace: ptr.To(gwv1.Namespace(namespace)),                // defaults to parent namespace if unset
							Kind:      ptr.To(gwv1.Kind("label")),                       // label is the only value
							Group:     ptr.To(gwv1.Group("delegation.gateway.solo.io")), // custom group for delegation
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

func (o *GatewayAPIOutput) generateFilterForURLRewrite(r *gloogwv1.Route, wrapper snapshot.Wrapper) *gwv1.HTTPRouteFilter {

	rf := &gwv1.HTTPRouteFilter{
		Type: gwv1.HTTPRouteFilterURLRewrite,
		URLRewrite: &gwv1.HTTPURLRewriteFilter{
			Path: &gwv1.HTTPPathModifier{},
		},
	}
	match, err := isExactMatch(r.GetMatchers())
	if err != nil {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "has multiple matchers with different types in same route")
	}
	if match {
		rf.URLRewrite.Path.Type = gwv1.FullPathHTTPPathModifier
		rf.URLRewrite.Path.ReplaceFullPath = ptr.To(r.GetOptions().GetPrefixRewrite().GetValue())
		rf.URLRewrite.Path.ReplacePrefixMatch = nil
	}
	match, err = isPrefixMatch(r.GetMatchers())
	if err != nil {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "has multiple matchers in same route")
	}

	if match {
		rf.URLRewrite.Path.Type = gwv1.PrefixMatchHTTPPathModifier
		rf.URLRewrite.Path.ReplacePrefixMatch = ptr.To(r.GetOptions().GetPrefixRewrite().GetValue())
		rf.URLRewrite.Path.ReplaceFullPath = nil
	}

	match, err = isRegexMatch(r.GetMatchers())
	if err != nil {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "has multiple matchers with different types in same route")
	}
	if match {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "has regex matchers and cannot be used with path rewrites in Gateway API")
		return nil
	}
	// regex rewrite, NOT SUPPORTED IN GATEWAY API
	if r.GetOptions().GetRegexRewrite() != nil {
		o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "regex rewrite not supported in Gateway API")
	}

	return rf
}

// Converts a single upstream to a GatewayAPI backend ref
func (o *GatewayAPIOutput) generateBackendRefForSingleUpstream(r *gloogwv1.Route, wrapper snapshot.Wrapper) gwv1.HTTPBackendRef {
	upstream := r.GetRouteAction().GetSingle().GetUpstream()
	var backendRef gwv1.HTTPBackendRef

	//TODO we need to lookup the upstream to see if its kube and then just reference kube directly
	var up *snapshot.UpstreamWrapper
	//if it is not a kube service or does not need http2
	var upstreamNs = upstream.GetNamespace()
	if upstreamNs == "" {
		upstreamNs = wrapper.GetNamespace()
	}

	up = o.edgeCache.GetUpstream(types.NamespacedName{Name: upstream.GetName(), Namespace: upstreamNs})

	if up == nil {
		// unknown reference to backend
		o.AddErrorFromWrapper(ERROR_TYPE_UNKNOWN_REFERENCE, wrapper, "upstream %s not found, referencing unknown upstream backend ref", upstream.GetName())

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
		o.AddErrorFromWrapper(ERROR_TYPE_UPDATE_OBJECT, wrapper, "service %s/%s uses http2, update its k8s service appProtocol=http2", up.Spec.GetKube().GetServiceNamespace(), up.Spec.GetKube().GetServiceName())
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

func (o *GatewayAPIOutput) convertMatch(m *matchers.Matcher, wrapper snapshot.Wrapper) (gwv1.HTTPRouteMatch, error) {
	hrm := gwv1.HTTPRouteMatch{
		QueryParams: []gwv1.HTTPQueryParamMatch{},
	}

	// header matching
	if len(m.GetHeaders()) > 0 {
		hrm.Headers = []gwv1.HTTPHeaderMatch{}
		for _, h := range m.GetHeaders() {
			// support invert header match https://github.com/solo-io/gloo/blob/main/projects/gateway2/translator/httproute/gateway_http_route_translator.go#L274
			if h.GetInvertMatch() == true {
				o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "invert match not currently supported")
			}
			if h.GetRegex() {
				hrm.Headers = append(hrm.Headers, gwv1.HTTPHeaderMatch{
					Type:  ptr.To(gwv1.HeaderMatchRegularExpression),
					Value: h.GetValue(),
					Name:  gwv1.HTTPHeaderName(h.GetName()),
				})
			} else {
				if h.GetValue() == "" {
					// no header value set so any value is good
					hrm.Headers = append(hrm.Headers, gwv1.HTTPHeaderMatch{
						Type:  ptr.To(gwv1.HeaderMatchRegularExpression),
						Value: "*",
						Name:  gwv1.HTTPHeaderName(h.GetName()),
					})
				} else {
					hrm.Headers = append(hrm.Headers, gwv1.HTTPHeaderMatch{
						Type:  ptr.To(gwv1.HeaderMatchExact),
						Value: h.GetValue(),
						Name:  gwv1.HTTPHeaderName(h.GetName()),
					})
				}
			}
		}

	}

	// method matching
	if len(m.GetMethods()) > 0 {
		if len(m.GetMethods()) > 1 {
			o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, wrapper, "gateway API only supports 1 method match per rule and %d were detected", len(m.GetMethods()))
		}
		hrm.Method = (*gwv1.HTTPMethod)(ptr.To(strings.ToUpper(m.GetMethods()[0])))
	}

	// query param matching
	if len(m.GetQueryParameters()) > 0 {
		for _, m := range m.GetQueryParameters() {
			if m.GetRegex() {
				hrm.QueryParams = append(hrm.QueryParams, gwv1.HTTPQueryParamMatch{
					Type:  ptr.To(gwv1.QueryParamMatchRegularExpression),
					Name:  (gwv1.HTTPHeaderName)(m.GetName()),
					Value: m.GetValue(),
				})
			} else {
				hrm.QueryParams = append(hrm.QueryParams, gwv1.HTTPQueryParamMatch{
					Type:  ptr.To(gwv1.QueryParamMatchExact),
					Name:  (gwv1.HTTPHeaderName)(m.GetName()),
					Value: m.GetValue(),
				})
			}
		}
	}

	// Path matching
	if m.GetPathSpecifier() != nil {
		if m.GetPrefix() != "" {
			hrm.Path = &gwv1.HTTPPathMatch{
				Type:  ptr.To(gwv1.PathMatchPathPrefix),
				Value: ptr.To(m.GetPrefix()),
			}
		}
		if m.GetExact() != "" {
			hrm.Path = &gwv1.HTTPPathMatch{
				Type:  ptr.To(gwv1.PathMatchExact),
				Value: ptr.To(m.GetExact()),
			}
		}
		if m.GetRegex() != "" {
			hrm.Path = &gwv1.HTTPPathMatch{
				Type:  ptr.To(gwv1.PathMatchRegularExpression),
				Value: ptr.To(m.GetRegex()),
			}
		}
	}
	return hrm, nil
}

func (o *GatewayAPIOutput) convertRouteTableToHTTPRoute(rt *snapshot.RouteTableWrapper) error {

	hr := &gwv1.HTTPRoute{
		TypeMeta: metav1.TypeMeta{
			Kind:       "HTTPRoute",
			APIVersion: gwv1.GroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      rt.Name,
			Namespace: rt.Namespace,
			Labels:    rt.Labels,
		},
		Spec: gwv1.HTTPRouteSpec{
			// CommonRouteSpec: gwv1.CommonRouteSpec{},
			// Hostnames: [],
			Rules: []gwv1.HTTPRouteRule{},
		},
	}
	if rt.Spec.GetWeight() != nil {
		if hr.ObjectMeta.GetLabels() == nil {
			hr.ObjectMeta.Labels = map[string]string{}
		}
		o.AddErrorFromWrapper(ERROR_TYPE_UPDATE_OBJECT, rt, "route weights are being set, enable KGW_WEIGHTED_ROUTE_PRECEDENCE=true environment variable.")

		hr.ObjectMeta.Labels[routeWeight] = rt.Spec.GetWeight().String()
	}

	for _, route := range rt.Spec.GetRoutes() {
		rule, err := o.convertRouteToRule(route, rt)
		if err != nil {
			return err
		}
		hr.Spec.Rules = append(hr.Spec.Rules, rule)
	}
	o.gatewayAPICache.AddHTTPRoute(snapshot.NewHTTPRouteWrapper(hr, rt.FileOrigin()))

	return nil
}

// This function validates that the RouteRable matchers are the same match type prefix or exact
// The reason being is that if you are doing a rewrite you can only have one type of filter applied
func validateMatchersAreSame(matches []*matchers.Matcher) error {

	var foundExact, foundPrefix, foundRegex bool
	for _, m := range matches {
		if m.GetExact() != "" {
			if foundPrefix || foundRegex {
				return fmt.Errorf("multiple matchers found")
			}
			foundExact = true
		}
		if m.GetPrefix() != "" {
			if foundExact || foundRegex {
				return fmt.Errorf("multiple matchers found")
			}
			foundPrefix = true
		}
		if m.GetRegex() != "" {
			if foundExact || foundPrefix {
				return fmt.Errorf("multiple matchers found")
			}
			foundRegex = true
		}
	}
	return nil
}

// tests to see if all matchers are exact
func isExactMatch(matches []*matchers.Matcher) (bool, error) {
	if err := validateMatchersAreSame(matches); err != nil {
		return false, err
	}
	for _, m := range matches {
		if m.GetExact() != "" {
			return true, nil
		}
	}
	return false, nil
}

// tests to see if all matchers are exact
func isPrefixMatch(matches []*matchers.Matcher) (bool, error) {
	if err := validateMatchersAreSame(matches); err != nil {
		return false, err
	}
	for _, m := range matches {
		if m.GetPrefix() != "" {
			return true, nil
		}
	}
	return false, nil
}

// tests to see if all matchers are regex
func isRegexMatch(matches []*matchers.Matcher) (bool, error) {
	if err := validateMatchersAreSame(matches); err != nil {
		return false, err
	}
	for _, m := range matches {
		if m.GetRegex() != "" {
			return true, nil
		}
	}
	return false, nil
}

func doHttpRouteLabelsMatch(matches map[string]string, labels map[string]string) bool {
	for k, v := range matches {
		if labels[k] != v {
			return false
		}
	}
	return true
}

// This checks to see if any of the route options are set for ones that are not supported in Gateway API
func isRouteOptionsSet(options *gloov1.RouteOptions) bool {
	//Features Supported By GatewayAPI
	// - RequestHeaderModifier
	// - ResponseHeaderModifier
	// - RequestRedirect
	// - URLRewrite
	// - Request Mirror
	// - CORS
	// - ExtensionRef
	// - Timeout (done)
	// - Retry (done)
	// - Session
	return options.GetExtProc() != nil ||
		options.GetStagedTransformations() != nil ||
		options.GetAutoHostRewrite() != nil ||
		options.GetFaults() != nil ||
		options.GetExtensions() != nil ||
		options.GetTracing() != nil ||
		options.GetAppendXForwardedHost() != nil ||
		options.GetLbHash() != nil ||
		options.GetUpgrades() != nil ||
		options.GetRatelimit() != nil ||
		options.GetRatelimitBasic() != nil ||
		options.GetWaf() != nil ||
		options.GetJwtConfig() != nil ||
		options.GetRbac() != nil ||
		options.GetDlp() != nil ||
		options.GetStagedTransformations() != nil ||
		options.GetEnvoyMetadata() != nil ||
		options.GetMaxStreamDuration() != nil ||
		options.GetIdleTimeout() != nil ||
		options.GetRegexRewrite() != nil ||
		options.GetExtauth() != nil ||
		options.GetAi() != nil ||
		options.GetCors() != nil
}
