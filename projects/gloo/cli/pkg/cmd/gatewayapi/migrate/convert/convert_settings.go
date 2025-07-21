package convert

import (
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/snapshot"
)

func (g *GatewayAPIOutput) convertSettings(settings *snapshot.SettingsWrapper) error {
	if settings == nil {
		return nil
	}
	spec := settings.Settings.Spec

	if spec.GetDiscoveryNamespace() != "" {
		//TODO(nick): how is this set now?
	}
	if spec.GetWatchNamespaces() != nil {
		g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, settings, "watchNamespaces is not supported")
	}
	if spec.GetKubernetesConfigSource() != nil {
		// this is a default
	}
	if spec.GetDirectoryConfigSource() != nil {
		g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, settings, "directoryConfigSource is not supported")
	}
	if spec.GetConsulKvSource() != nil {
		g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, settings, "consulKvSource is not supported")
	}
	if spec.GetKubernetesSecretSource() != nil {
		// This is the default
	}
	if spec.GetVaultSecretSource() != nil {
		g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, settings, "vaultSecretSource is not supported")
	}
	if spec.GetDirectorySecretSource() != nil {
		g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, settings, "directorySecretSource is not supported")
	}
	if spec.GetSecretOptions() != nil {
		// This is no longer needed and helps configure the secret source
	}
	if spec.GetKubernetesArtifactSource() != nil {
		// This is the default
	}
	if spec.GetDirectoryArtifactSource() != nil {
		g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, settings, "directoryArtifactSource is not supported")
	}
	if spec.GetConsulKvArtifactSource() != nil {
		g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, settings, "consulKvArtifactSouce is not supported")
	}
	if spec.GetRefreshRate() != nil {
		//this is not needed in kgateway, all users should start with the default
	}
	if spec.GetLinkerd() == true {
		g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, settings, "linkerd no longer supported")
	}
	if spec.GetKnative() != nil {
		g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, settings, "knative no longer supported")
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
			g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, settings, "AWS credentialDiscovery is not supported")
		}
		if gloo.GetAwsOptions() != nil && spec.GetGloo().GetAwsOptions().GetFallbackToFirstFunction() != nil && spec.GetGloo().GetAwsOptions().GetFallbackToFirstFunction().GetValue() == true {
			g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, settings, "AWS fallbackToFirstFunction is not supported")
		}
		if gloo.GetAwsOptions() != nil && spec.GetGloo().GetAwsOptions().GetPropagateOriginalRouting() != nil && spec.GetGloo().GetAwsOptions().GetPropagateOriginalRouting().GetValue() == true {
			g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, settings, "AWS propagateOriginalRouting is not supported")
		}
		if gloo.GetAwsOptions() != nil && spec.GetGloo().GetAwsOptions().GetServiceAccountCredentials() != nil {
			g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, settings, "AWS serviceAccountCredentials is not supported")
		}
		if gloo.GetCircuitBreakers() != nil {
			g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, settings, "gloo.circuitBreakers is not supported")
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
		g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, settings, "consul is not supported")
	}
	if spec.GetConsulDiscovery() != nil {
		g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, settings, "consul discovery is not supported")
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
		extension := g.generateGatewayExtensionForRateLimit(spec.GetRatelimitServer(), "rate-limit", settings)
		if extension != nil {
			g.gatewayAPICache.AddGatewayExtension(snapshot.NewGatewayExtensionWrapper(extension, settings.FileOrigin()))
		}
	}
	if spec.GetRbac() != nil && spec.GetRbac().GetRequireRbac() == true {
		g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, settings, "requireRbac is not supported")
	}
	if spec.GetExtauth() != nil {
		extension := g.generateGatewayExtensionForExtAuth(spec.GetExtauth(), "ext-authz", settings)
		if extension != nil {
			g.gatewayAPICache.AddGatewayExtension(snapshot.NewGatewayExtensionWrapper(extension, settings.FileOrigin()))
		}
	}
	if spec.GetNamedExtauth() != nil {
		for name, extauth := range spec.GetNamedExtauth() {
			extension := g.generateGatewayExtensionForExtAuth(extauth, name, settings)
			if extension != nil {
				g.gatewayAPICache.AddGatewayExtension(snapshot.NewGatewayExtensionWrapper(extension, settings.FileOrigin()))
			}
		}
	}
	if spec.GetCachingServer() != nil {
		g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, settings, "cachingServer is not supported")
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
			g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, settings, "global upstream SSLParameters  is not supported")
		}
	}
	if spec.GetConsoleOptions() != nil {
		// no need to do anything here
	}
	//if spec.GetGraphqlOptions() != nil {
	//	o.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, settings, "graphql is not supported")
	//}
	if spec.GetExtProc() != nil {
		g.generateGatewayExtensionForExtProc(spec.GetExtProc(), "global-ext-proc", settings)
	}
	if spec.GetKnative() != nil {
		g.AddErrorFromWrapper(ERROR_TYPE_NOT_SUPPORTED, settings, "knative is not supported")
	}
	if spec.GetWatchNamespaceSelectors() != nil {
		// nothing to do here
	}

	return nil
}
