package translator

import (
	"errors"
	"fmt"
	"strings"
	"time"

	envoy_config_cluster_v3 "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	envoy_config_core_v3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	"github.com/envoyproxy/go-control-plane/pkg/wellknown"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/duration"
	_struct "github.com/golang/protobuf/ptypes/struct"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/rotisserie/eris"
	"github.com/solo-io/gloo/pkg/utils/api_conversion"
	"github.com/solo-io/gloo/pkg/utils/envutils"
	"github.com/solo-io/gloo/projects/gloo/constants"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	v1_options "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/ssl"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins"
	upstream_proxy_protocol "github.com/solo-io/gloo/projects/gloo/pkg/plugins/utils/upstreamproxyprotocol"
	"github.com/solo-io/gloo/projects/gloo/pkg/utils"
	"github.com/solo-io/gloo/projects/gloo/pkg/xds"
	"github.com/solo-io/go-utils/contextutils"
	"github.com/solo-io/solo-kit/pkg/api/v2/reporter"
	"github.com/solo-io/solo-kit/pkg/utils/prototime"
	"go.opencensus.io/trace"
	"go.uber.org/zap"
	_structpb "google.golang.org/protobuf/types/known/structpb"
)

func (t *translatorInstance) computeClusters(
	params plugins.Params,
	reports reporter.ResourceReports,
	upstreamRefKeyToEndpoints map[string][]*v1.Endpoint,
	proxy *v1.Proxy,
) ([]*envoy_config_cluster_v3.Cluster, map[*envoy_config_cluster_v3.Cluster]*v1.Upstream) {
	ctx, span := trace.StartSpan(params.Ctx, "gloo.translator.computeClusters")
	defer span.End()
	params.Ctx = contextutils.WithLogger(ctx, "compute_clusters")

	// snapshot contains both real and service-derived upstreams
	upstreamGroups := params.Snapshot.UpstreamGroups
	upstreams := params.Snapshot.Upstreams
	clusters := make([]*envoy_config_cluster_v3.Cluster, 0, len(upstreams))
	validateUpstreamLambdaFunctions(proxy, upstreams, upstreamGroups, reports)

	clusterToUpstreamMap := make(map[*envoy_config_cluster_v3.Cluster]*v1.Upstream)
	for _, upstream := range upstreams {
		eds := false
		if eps, ok := upstreamRefKeyToEndpoints[upstream.GetMetadata().Ref().Key()]; ok && len(eps) > 0 {
			eds = true
		}
		cluster, errs := t.computeCluster(params, upstream, eds)
		for _, err := range errs {
			var warning *Warning
			if errors.As(err, &warning) {
				reports.AddWarning(upstream, err.Error())
			} else {
				reports.AddError(upstream, err)
			}
		}
		clusterToUpstreamMap[cluster] = upstream
		clusters = append(clusters, cluster)
	}

	return clusters, clusterToUpstreamMap
}

// This function is intented to be used when translating a single upstream outside of the context of a full snapshot.
// This happens in the kube gateway krt implementation.
func (t *translatorInstance) TranslateCluster(
	params plugins.Params,
	upstream *v1.Upstream,
) (*envoy_config_cluster_v3.Cluster, []error) {
	t.lock.Lock()
	defer t.lock.Unlock()
	for _, p := range t.pluginRegistry.GetUpstreamPlugins() {
		p.Init(plugins.InitParams{Ctx: params.Ctx, Settings: t.settings})
	}
	// as we don't know if we have endpoints for this upstream,
	// we will let the upstream plugins will set the cluster type
	eds := false
	c, err := t.computeCluster(params, upstream, eds)
	if c != nil && c.GetEdsClusterConfig() != nil {
		endpointClusterName, err2 := GetEndpointClusterName(c.GetName(), upstream)
		if err2 == nil {
			c.GetEdsClusterConfig().ServiceName = endpointClusterName
		}
	}
	return c, err
}

func (t *translatorInstance) computeCluster(
	params plugins.Params,
	upstream *v1.Upstream,
	eds bool,
) (*envoy_config_cluster_v3.Cluster, []error) {
	logger := contextutils.LoggerFrom(params.Ctx)
	params.Ctx = contextutils.WithLogger(params.Ctx, upstream.GetMetadata().GetName())
	out, errs := t.initializeCluster(upstream, eds, &params.Snapshot.Secrets)

	for _, plugin := range t.pluginRegistry.GetUpstreamPlugins() {
		if err := plugin.ProcessUpstream(params, upstream, out); err != nil {
			logger.Debug("Error processing upstream", zap.String("upstream", upstream.GetMetadata().Ref().String()), zap.Error(err), zap.String("plugin", plugin.Name()))
			errs = append(errs, err)
		}
	}
	if err := validateCluster(out); err != nil {
		logger.Debug("Error validating cluster ", zap.String("upstream", upstream.GetMetadata().Ref().String()), zap.Error(err))
		errs = append(errs, eris.Wrap(err, "cluster was configured improperly by one or more plugins"))
	}
	return out, errs
}

func (t *translatorInstance) initializeCluster(
	upstream *v1.Upstream,
	eds bool,
	secrets *v1.SecretList,
) (*envoy_config_cluster_v3.Cluster, []error) {
	var errorList []error
	hcConfig, err := createHealthCheckConfig(upstream, secrets, t.shouldEnforceNamespaceMatch)
	if err != nil {
		errorList = append(errorList, err)
	}
	detectCfg, err := createOutlierDetectionConfig(upstream)
	if err != nil {
		errorList = append(errorList, err)
	}

	preconnect, err := getPreconnectPolicy(upstream.GetPreconnectPolicy())
	if err != nil {
		errorList = append(errorList, err)
	}
	dnsRefreshRate, err := getDnsRefreshRate(upstream)
	if err != nil {
		errorList = append(errorList, err)
	}

	circuitBreakers := t.settings.GetGloo().GetCircuitBreakers()
	out := &envoy_config_cluster_v3.Cluster{
		Name:             UpstreamToClusterName(upstream.GetMetadata().Ref()),
		Metadata:         new(envoy_config_core_v3.Metadata),
		CircuitBreakers:  getCircuitBreakers(upstream.GetCircuitBreakers(), circuitBreakers),
		LbSubsetConfig:   createLbConfig(upstream),
		HealthChecks:     hcConfig,
		OutlierDetection: detectCfg,
		// defaults to Cluster_USE_CONFIGURED_PROTOCOL
		ProtocolSelection: envoy_config_cluster_v3.Cluster_ClusterProtocolSelection(upstream.GetProtocolSelection()),
		// this field can be overridden by plugins
		ConnectTimeout:            ptypes.DurationProto(ClusterConnectionTimeout),
		Http2ProtocolOptions:      getHttp2options(upstream),
		IgnoreHealthOnHostRemoval: upstream.GetIgnoreHealthOnHostRemoval().GetValue(),
		RespectDnsTtl:             upstream.GetRespectDnsTtl().GetValue(),
		DnsRefreshRate:            dnsRefreshRate,
		PreconnectPolicy:          preconnect,
	}
	// for kube gateway, use new stats name format
	if envutils.IsEnvTruthy(constants.GlooGatewayEnableK8sGwControllerEnv) {
		out.AltStatName = UpstreamToClusterStatsName(upstream)
	}

	if sslConfig := upstream.GetSslConfig(); sslConfig != nil {
		applyDefaultsToUpstreamSslConfig(sslConfig, t.settings.GetUpstreamOptions())
		cfg, err := utils.NewSslConfigTranslator().ResolveUpstreamSslConfig(*secrets, sslConfig)
		if err != nil {
			// if we are configured to warn on missing tls secret and we match that error, add a
			// warning instead of error to the report.
			if t.settings.GetGateway().GetValidation().GetWarnMissingTlsSecret().GetValue() &&
				errors.Is(err, utils.SslSecretNotFoundError) {
				errorList = append(errorList, &Warning{
					Message: err.Error(),
				})
			} else {
				errorList = append(errorList, err)
			}
		} else {
			typedConfig, err := utils.MessageToAny(cfg)
			if err != nil {
				// TODO: Need to change the upstream to use a direct response action instead of leaving the upstream untouched
				// Difficult because direct response is not on the upsrtream but on the virtual host
				// The fallback listener would take much more piping as well
				panic(err)
			} else {
				out.TransportSocket = &envoy_config_core_v3.TransportSocket{
					Name:       wellknown.TransportSocketTls,
					ConfigType: &envoy_config_core_v3.TransportSocket_TypedConfig{TypedConfig: typedConfig},
				}
			}
		}
	}
	// proxyprotocol may be wiped by some plugins that transform transport sockets
	// see static and failover at time of writing.
	if upstream.GetProxyProtocolVersion() != nil {

		tp, err := upstream_proxy_protocol.WrapWithPProtocol(out.GetTransportSocket(), upstream.GetProxyProtocolVersion().GetValue())
		if err != nil {
			errorList = append(errorList, err)
		} else {
			out.TransportSocket = tp
		}
	}

	// set Type = EDS if we have endpoints for the upstream
	if eds {
		xds.SetEdsOnCluster(out, t.settings)
	}
	return out, errorList
}

var (
	DefaultHealthCheckTimeout  = &duration.Duration{Seconds: 5}
	DefaultHealthCheckInterval = prototime.DurationToProto(time.Millisecond * 100)
	DefaultThreshold           = &wrappers.UInt32Value{
		Value: 5,
	}

	NilFieldError = func(fieldName string) error {
		return eris.Errorf("The field %s cannot be nil", fieldName)
	}

	minimumDnsRefreshRate = prototime.DurationToProto(time.Millisecond * 1)
)

func createHealthCheckConfig(upstream *v1.Upstream, secrets *v1.SecretList, shouldEnforceNamespaceMatch bool) ([]*envoy_config_core_v3.HealthCheck, error) {
	if upstream == nil {
		return nil, nil
	}
	result := make([]*envoy_config_core_v3.HealthCheck, 0, len(upstream.GetHealthChecks()))
	for i, hc := range upstream.GetHealthChecks() {
		// These values are required by envoy, but not explicitly
		if hc.GetHealthyThreshold() == nil {
			return nil, NilFieldError(fmt.Sprintf("HealthCheck[%d].HealthyThreshold", i))
		}
		if hc.GetUnhealthyThreshold() == nil {
			return nil, NilFieldError(fmt.Sprintf("HealthCheck[%d].UnhealthyThreshold", i))
		}
		if hc.GetHealthChecker() == nil {
			return nil, NilFieldError(fmt.Sprintf("HealthCheck[%d].HealthChecker", i))
		}
		options := api_conversion.HeaderSecretOptions{UpstreamNamespace: upstream.GetMetadata().GetNamespace(), EnforceNamespaceMatch: shouldEnforceNamespaceMatch}
		converted, err := api_conversion.ToEnvoyHealthCheck(hc, secrets, options)
		if err != nil {
			return nil, err
		}
		result = append(result, converted)
	}
	return result, nil
}

func createOutlierDetectionConfig(upstream *v1.Upstream) (*envoy_config_cluster_v3.OutlierDetection, error) {
	if upstream.GetOutlierDetection() == nil {
		return nil, nil
	}
	if upstream.GetOutlierDetection().GetInterval() == nil {
		return nil, NilFieldError("OutlierDetection.HealthChecker.Interval")
	}
	return api_conversion.ToEnvoyOutlierDetection(upstream.GetOutlierDetection()), nil
}

func convertDefaultSubset(defaultSubset *v1_options.Subset) *_struct.Struct {
	if defaultSubset == nil {
		return nil
	}
	subsetVals := make(map[string]interface{}, len(defaultSubset.GetValues()))
	for k, v := range defaultSubset.GetValues() {
		subsetVals[k] = v
	}
	converted, err := _structpb.NewStruct(subsetVals)
	if err != nil {
		return nil
	}
	return converted
}

func convertFallbackPolicy(fallbackPolicy v1_options.FallbackPolicy) envoy_config_cluster_v3.Cluster_LbSubsetConfig_LbSubsetFallbackPolicy {
	if fallbackPolicy == v1_options.FallbackPolicy_NO_FALLBACK {
		return envoy_config_cluster_v3.Cluster_LbSubsetConfig_NO_FALLBACK
	} else if fallbackPolicy == v1_options.FallbackPolicy_ANY_ENDPOINT {
		return envoy_config_cluster_v3.Cluster_LbSubsetConfig_ANY_ENDPOINT
	} else if fallbackPolicy == v1_options.FallbackPolicy_DEFAULT_SUBSET {
		return envoy_config_cluster_v3.Cluster_LbSubsetConfig_DEFAULT_SUBSET
	}
	// this should not happen, return the desired default
	return envoy_config_cluster_v3.Cluster_LbSubsetConfig_ANY_ENDPOINT
}

func createLbConfig(upstream *v1.Upstream) *envoy_config_cluster_v3.Cluster_LbSubsetConfig {
	specGetter, ok := upstream.GetUpstreamType().(v1.SubsetSpecGetter)
	if !ok {
		return nil
	}
	glooSubsetConfig := specGetter.GetSubsetSpec()
	if glooSubsetConfig == nil {
		return nil
	}

	subsetConfig := &envoy_config_cluster_v3.Cluster_LbSubsetConfig{
		// when omitted, fallback policy defaults to ANY_ENDPOINT
		FallbackPolicy: convertFallbackPolicy(glooSubsetConfig.GetFallbackPolicy()),
		DefaultSubset:  convertDefaultSubset(glooSubsetConfig.GetDefaultSubset()),
	}
	for _, selector := range glooSubsetConfig.GetSelectors() {
		subsetConfig.SubsetSelectors = append(subsetConfig.GetSubsetSelectors(), &envoy_config_cluster_v3.Cluster_LbSubsetConfig_LbSubsetSelector{
			Keys:                selector.GetKeys(),
			SingleHostPerSubset: selector.GetSingleHostPerSubset(),
		})
	}

	return subsetConfig
}

// TODO: add more validation here
func validateCluster(c *envoy_config_cluster_v3.Cluster) error {
	if c.GetClusterType() != nil {
		// TODO(yuval-k): this is a custom cluster, we cant validate it for now.
		return nil
	}
	clusterType := c.GetType()
	if clusterType == envoy_config_cluster_v3.Cluster_STATIC ||
		clusterType == envoy_config_cluster_v3.Cluster_STRICT_DNS ||
		clusterType == envoy_config_cluster_v3.Cluster_LOGICAL_DNS {
		if len(c.GetLoadAssignment().GetEndpoints()) == 0 {
			return eris.Errorf("cluster type %v specified but LoadAssignment was empty", clusterType.String())
		}
	}

	if c.GetDnsRefreshRate() != nil {
		if clusterType != envoy_config_cluster_v3.Cluster_STRICT_DNS && clusterType != envoy_config_cluster_v3.Cluster_LOGICAL_DNS {
			return eris.Errorf("DnsRefreshRate is only valid with STRICT_DNS or LOGICAL_DNS cluster type, found %v", clusterType)
		}
	}
	return nil
}

// Convert the first non nil circuit breaker.
func getCircuitBreakers(cfgs ...*v1.CircuitBreakerConfig) *envoy_config_cluster_v3.CircuitBreakers {
	for _, cfg := range cfgs {
		if cfg != nil {
			envoyCfg := &envoy_config_cluster_v3.CircuitBreakers{}
			envoyCfg.Thresholds = []*envoy_config_cluster_v3.CircuitBreakers_Thresholds{{
				MaxConnections:     cfg.GetMaxConnections(),
				MaxPendingRequests: cfg.GetMaxPendingRequests(),
				MaxRequests:        cfg.GetMaxRequests(),
				MaxRetries:         cfg.GetMaxRetries(),
			}}
			return envoyCfg
		}
	}
	return nil
}

// getPreconnectPolicy follows the naming of the rest of functions here
// it MAY someday deal with inheritance like the other functions
// it consumes an ordered list of preconnect policies
// it returns the first non-nil policy
func getPreconnectPolicy(cfgs ...*v1.PreconnectPolicy) (*envoy_config_cluster_v3.Cluster_PreconnectPolicy, error) {
	// since we dont want strict reliance on envoy's current api
	// but still able to map as closely as possible
	// if not nil then convert the gloo configurations to envoy
	for _, curConfig := range cfgs {
		if curConfig == nil {
			continue
		}

		// we have a policy. However we dont respect proto constraints
		// therefore we need to check here to see if the applied configuration is
		// SILLY and warn the user.
		perUpstream := curConfig.GetPerUpstreamPreconnectRatio()
		predictive := curConfig.GetPredictivePreconnectRatio()

		eStrings := []string{}
		if perUpstream != nil && (perUpstream.GetValue() < 1 || perUpstream.GetValue() > 3) {
			eStrings = append(eStrings, fmt.Sprintf("perupstream must be between 1 and 3 if set, it was: %v", perUpstream.GetValue()))
		}
		if predictive != nil && (predictive.GetValue() < 1 || predictive.GetValue() > 3) {
			eStrings = append(eStrings, fmt.Sprintf("predictive must be between 1 and 3 if set, it was: %v", perUpstream.GetValue()))
		}
		if len(eStrings) > 0 {
			return nil, errors.New("invalid preconnect policy: " + strings.Join(eStrings, "; "))
		}
		return &envoy_config_cluster_v3.Cluster_PreconnectPolicy{
			PerUpstreamPreconnectRatio: curConfig.GetPerUpstreamPreconnectRatio(),
			PredictivePreconnectRatio:  curConfig.GetPredictivePreconnectRatio(),
		}, nil
	}
	return nil, nil
}

func getHttp2options(us *v1.Upstream) *envoy_config_core_v3.Http2ProtocolOptions {
	if us.GetUseHttp2().GetValue() {
		return &envoy_config_core_v3.Http2ProtocolOptions{}
	}
	return nil
}

// getDnsRefreshRate returns the DnsRefreshRate for an Upstream:
// - defined and valid: returns the duration
// - defined and invalid: adds a warning and returns nil
// - undefined: returns nil
func getDnsRefreshRate(us *v1.Upstream) (*duration.Duration, error) {
	refreshRate := us.GetDnsRefreshRate()
	if refreshRate == nil {
		return nil, nil
	}

	if refreshRate.AsDuration() < minimumDnsRefreshRate.AsDuration() {
		return nil, &Warning{
			Message: fmt.Sprintf("dnsRefreshRate was set below minimum requirement (%s), ignoring configuration", minimumDnsRefreshRate),
		}
	}

	return refreshRate, nil
}

// Validates routes that point to the current AWS lambda upstream
// Checks that the function the route is pointing to is available on the upstream
// else it adds an error to the upstream, so that invalid route replacement can be used.
func validateUpstreamLambdaFunctions(proxy *v1.Proxy, upstreams v1.UpstreamList, upstreamGroups v1.UpstreamGroupList, reports reporter.ResourceReports) {
	// Create a set of the lambda functions in each upstream
	upstreamLambdas := make(map[string]map[string]bool)
	for _, upstream := range upstreams {
		lambdaFuncs := upstream.GetAws().GetLambdaFunctions()
		for _, lambda := range lambdaFuncs {
			upstreamRef := UpstreamToClusterName(upstream.GetMetadata().Ref())
			if upstreamLambdas[upstreamRef] == nil {
				upstreamLambdas[upstreamRef] = make(map[string]bool)
			}
			upstreamLambdas[upstreamRef][lambda.GetLogicalName()] = true
		}
	}

	for _, listener := range proxy.GetListeners() {
		virtualHosts := utils.GetVirtualHostsForListener(listener)

		for _, virtualHost := range virtualHosts {
			// Validate all routes to make sure that if they point to a lambda, it exists.
			for _, route := range virtualHost.GetRoutes() {
				validateRouteDestinationForValidLambdas(proxy, route, upstreamGroups, reports, upstreamLambdas)
			}
		}
	}
}

// Validates a route that may have a single or multi upstream destinations to make sure that any lambda upstreams are referencing valid lambdas
func validateRouteDestinationForValidLambdas(
	proxy *v1.Proxy,
	route *v1.Route,
	upstreamGroups v1.UpstreamGroupList,
	reports reporter.ResourceReports,
	upstreamLambdas map[string]map[string]bool,
) {
	// Append destinations to a destination list to process all of them in one go
	var destinations []*v1.Destination

	_, ok := route.GetAction().(*v1.Route_RouteAction)
	if !ok {
		// If this is not a Route_RouteAction (e.g. Route_DirectResponseAction, Route_RedirectAction), there is no destination to validate
		return
	}
	routeAction := route.GetRouteAction()
	switch typedRoute := routeAction.GetDestination().(type) {
	case *v1.RouteAction_Single:
		{
			destinations = append(destinations, routeAction.GetSingle())
		}
	case *v1.RouteAction_Multi:
		{
			multiDest := routeAction.GetMulti()
			for _, weightedDest := range multiDest.GetDestinations() {
				destinations = append(destinations, weightedDest.GetDestination())
			}
		}
	case *v1.RouteAction_UpstreamGroup:
		{
			ugRef := typedRoute.UpstreamGroup
			ug, err := upstreamGroups.Find(ugRef.GetNamespace(), ugRef.GetName())
			if err != nil {
				reports.AddError(proxy, fmt.Errorf("upstream group not found, (Name: %s, Namespace: %s)", ugRef.GetName(), ugRef.GetNamespace()))
				return
			}
			for _, weightedDest := range ug.GetDestinations() {
				destinations = append(destinations, weightedDest.GetDestination())
			}
		}
	case *v1.RouteAction_ClusterHeader:
		{
			// no upstream configuration is provided in this case; can't validate the route
			return
		}
	case *v1.RouteAction_DynamicForwardProxy:
		{
			// no upstream configuration is provided in this case; can't validate the route
			return
		}
	default:
		{
			reports.AddError(proxy, fmt.Errorf("route destination type %T not supported with AWS Lambda", typedRoute))
			return
		}
	}

	// Process destinations (upstreams)
	for _, dest := range destinations {
		routeUpstream := dest.GetUpstream()
		// Check that route is pointing to current upstream
		if routeUpstream != nil {
			// Get the lambda functions that this upstream has
			lambdaFuncSet := upstreamLambdas[UpstreamToClusterName(routeUpstream)]
			routeLambda := dest.GetDestinationSpec().GetAws()
			routeLambdaName := routeLambda.GetLogicalName()
			// If route is pointing to a lambda that does not exist on this upstream, report error on the upstream
			if routeLambda != nil && lambdaFuncSet[routeLambdaName] == false {
				// Add error to the proxy which has the faulty route pointing to a non-existent lambda
				reports.AddError(proxy, fmt.Errorf("a route references %s AWS lambda which does not exist on the route's upstream", routeLambdaName))
			}
		}
	}
}

// Apply defaults to UpstreamSslConfig
func applyDefaultsToUpstreamSslConfig(sslConfig *ssl.UpstreamSslConfig, options *v1.UpstreamOptions) {
	if options == nil {
		return
	}

	// Apply default SslParameters if none are defined on upstream
	if sslConfig.GetParameters() == nil {
		sslConfig.Parameters = options.GetSslParameters()
	}
}
