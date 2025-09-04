package translator

import (
	"context"
	"fmt"
	"hash/fnv"
	"os"
	"strconv"
	"sync"

	"github.com/solo-io/gloo/pkg/utils/api_conversion"
	"github.com/solo-io/gloo/pkg/utils/envutils"
	"github.com/solo-io/gloo/projects/gloo/constants"

	envoy_config_cluster_v3 "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	envoy_config_endpoint_v3 "github.com/envoyproxy/go-control-plane/envoy/config/endpoint/v3"
	envoy_config_listener_v3 "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	envoy_config_route_v3 "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	"github.com/golang/protobuf/proto"
	errors "github.com/rotisserie/eris"
	envoyvalidation "github.com/solo-io/gloo/pkg/utils/envoyutils/validation"
	validationapi "github.com/solo-io/gloo/projects/gloo/pkg/api/grpc/validation"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins"
	"github.com/solo-io/gloo/projects/gloo/pkg/utils"
	"github.com/solo-io/gloo/projects/gloo/pkg/utils/validation"
	"github.com/solo-io/gloo/projects/gloo/pkg/xds"
	"github.com/solo-io/go-utils/contextutils"
	"github.com/solo-io/protoc-gen-ext/pkg/hasher/hashstructure"
	envoycache "github.com/solo-io/solo-kit/pkg/api/v1/control-plane/cache"
	"github.com/solo-io/solo-kit/pkg/api/v1/control-plane/resource"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-kit/pkg/api/v2/reporter"
	"go.opencensus.io/trace"
	"go.uber.org/zap"
	proto2 "google.golang.org/protobuf/proto"
)

// logComputeTranslator is a helper function that logs translator messages only when COMPUTE_TRANSLATOR_LOGS is enabled
func logComputeTranslator(logger *zap.SugaredLogger, msg string, keysAndValues ...interface{}) {
	if envutils.IsEnvTruthy(constants.ComputeTranslatorLogsEnv) {
		// Add the issue label to all gated logs
		keysAndValues = append([]interface{}{"issue", "8539"}, keysAndValues...)
		logger.Infow(msg, keysAndValues...)
	}
}

type Translator interface {
	// Translate converts a Proxy CR into an xDS Snapshot
	// Any errors that are encountered during translation are appended to the ResourceReports
	// It is invalid for us to return an error here, since translation of resources always needs
	// to results in an xDS Snapshot so we are resilient to pod restarts
	Translate(
		params plugins.Params,
		proxy *v1.Proxy,
	) (envoycache.Snapshot, reporter.ResourceReports, *validationapi.ProxyReport)
}
type ClusterTranslator interface {
	// Translate converts a Upstream CR into an xDS Snapshot
	// Any errors or warnings that are encountered during translation are returned, along with the
	// envoy cluster.
	TranslateCluster(
		params plugins.Params,
		upstream *v1.Upstream,
	) (*envoy_config_cluster_v3.Cluster, []error)
}

var (
	_ Translator = new(translatorInstance)
)

// translatorInstance is the implementation for a Translator used during Gloo translation
type translatorInstance struct {
	lock                        sync.Mutex
	pluginRegistry              plugins.PluginRegistry
	settings                    *v1.Settings
	hasher                      func(resources []envoycache.Resource) (uint64, error)
	listenerTranslatorFactory   *ListenerSubsystemTranslatorFactory
	shouldEnforceNamespaceMatch bool
}

func NewDefaultTranslator(settings *v1.Settings, pluginRegistry plugins.PluginRegistry) *translatorInstance {
	return NewTranslatorWithHasher(utils.NewSslConfigTranslator(), settings, pluginRegistry, EnvoyCacheResourcesListToFnvHash)
}

func NewTranslatorWithHasher(
	sslConfigTranslator utils.SslConfigTranslator,
	settings *v1.Settings,
	pluginRegistry plugins.PluginRegistry,
	hasher func(resources []envoycache.Resource) (uint64, error),
) *translatorInstance {
	shouldEnforceStr := os.Getenv(api_conversion.MatchingNamespaceEnv)
	shouldEnforceNamespaceMatch := false
	if shouldEnforceStr != "" {
		var err error
		shouldEnforceNamespaceMatch, err = strconv.ParseBool(shouldEnforceStr)
		if err != nil {
			// TODO: what to do here?
		}
	}
	return &translatorInstance{
		lock:                        sync.Mutex{},
		pluginRegistry:              pluginRegistry,
		settings:                    settings,
		hasher:                      hasher,
		listenerTranslatorFactory:   NewListenerSubsystemTranslatorFactory(pluginRegistry, sslConfigTranslator, settings),
		shouldEnforceNamespaceMatch: shouldEnforceNamespaceMatch,
	}
}

func (t *translatorInstance) Translate(
	params plugins.Params,
	proxy *v1.Proxy,
) (envoycache.Snapshot, reporter.ResourceReports, *validationapi.ProxyReport) {
	ctx, span := trace.StartSpan(params.Ctx, "gloo.translator.Translate")
	defer span.End()
	params.Ctx = contextutils.WithLogger(ctx, "translator")
	logger := contextutils.LoggerFrom(params.Ctx)

	// ADD TRANSLATION START LOGGING
	logComputeTranslator(logger, "Starting proxy translation",
		"proxyName", proxy.GetMetadata().GetName(),
		"proxyNamespace", proxy.GetMetadata().GetNamespace(),
		"upstreamCount", len(params.Snapshot.Upstreams),
		"upstreamGroupCount", len(params.Snapshot.UpstreamGroups),
		"secretCount", len(params.Snapshot.Secrets),
		"artifactCount", len(params.Snapshot.Artifacts),
		"endpointCount", len(params.Snapshot.Endpoints))

	// Plugins may modify the snapshot, so we need to make sure we are operating on a copy
	//     We make a shallow copy here, which means that the top level slices are copied, but the
	//     objects in the slices are not. This is sufficient for our use case, since we don't
	//     expect plugins to modify the objects themselves, only to add/remove objects from the slices.
	//     If we ever need to support plugins that modify the objects themselves, we would need to
	//     make a deep copy here. However, that would be a significant performance hit, so we should
	//     avoid it if possible.
	//
	//     NOTE: We should be careful about this. If we ever add a plugin that modifies the objects
	//     themselves, we would need to make a deep copy here. Otherwise, the plugin would modify
	//     the original objects, which could cause issues if the snapshot is used elsewhere after
	//     reset.
	logComputeTranslator(logger, "Initializing translation plugins",
		"pluginCount", len(t.pluginRegistry.GetPlugins()))

	t.lock.Lock()
	defer t.lock.Unlock()
	for _, p := range t.pluginRegistry.GetPlugins() {
		p.Init(plugins.InitParams{
			Ctx:      params.Ctx,
			Settings: t.settings,
		})
	}

	// prepare reports used to aggregate Warnings/Errors encountered during translation
	reports := make(reporter.ResourceReports)
	proxyReport := validation.MakeReport(proxy)
	logComputeTranslator(logger, "Starting cluster subsystem translation",
		"proxyName", proxy.GetMetadata().GetName())

	var clusters []*envoy_config_cluster_v3.Cluster
	var endpoints []*envoy_config_endpoint_v3.ClusterLoadAssignment
	clusters, endpoints = t.translateClusterSubsystemComponents(params, proxy, reports)

	logComputeTranslator(logger, "Completed cluster subsystem translation",
		"proxyName", proxy.GetMetadata().GetName(),
		"clusterCount", len(clusters))

	logComputeTranslator(logger, "Starting listener subsystem translation",
		"proxyName", proxy.GetMetadata().GetName())

	routeConfigs, listeners := t.translateListenerSubsystemComponents(params, proxy, proxyReport)

	logComputeTranslator(logger, "Completed listener subsystem translation",
		"proxyName", proxy.GetMetadata().GetName(),
		"listenerCount", len(listeners))

	// run Resource Generator Plugins
	logComputeTranslator(logger, "Running resource generator plugins",
		"pluginCount", len(t.pluginRegistry.GetResourceGeneratorPlugins()))

	for _, plugin := range t.pluginRegistry.GetResourceGeneratorPlugins() {
		generatedClusters, generatedEndpoints, generatedRouteConfigs, generatedListeners, err := plugin.GeneratedResources(params, clusters, endpoints, routeConfigs, listeners)
		if err != nil {
			logger.Warnw("Resource generator plugin failed",
				"issue", "8539",
				"plugin", fmt.Sprintf("%T", plugin),
				"error", err.Error())
			reports.AddError(proxy, err)
		}
		clusters = append(clusters, generatedClusters...)
		endpoints = append(endpoints, generatedEndpoints...)
		routeConfigs = append(routeConfigs, generatedRouteConfigs...)
		listeners = append(listeners, generatedListeners...)
	}

	logger.Infow("Generating final xDS snapshot",
		"issue", "8539",
		"finalClusters", len(clusters),
		"finalEndpoints", len(endpoints),
		"finalRouteConfigs", len(routeConfigs),
		"finalListeners", len(listeners))

	xdsSnapshot := t.generateXDSSnapshot(params, clusters, endpoints, routeConfigs, listeners)

	if err := validation.GetProxyError(proxyReport); err != nil {
		logComputeTranslator(logger, "Proxy translation validation error detected",
			"proxyName", proxy.GetMetadata().GetName(),
			"error", err.Error())
		reports.AddError(proxy, err)
	}

	if warnings := validation.GetProxyWarning(proxyReport); len(warnings) > 0 {
		logComputeTranslator(logger, "Proxy translation validation warnings detected",
			"proxyName", proxy.GetMetadata().GetName(),
			"warningCount", len(warnings))
		for _, warning := range warnings {
			reports.AddWarning(proxy, warning)
		}
	}

	// Validate the xDS snapshot
	ctx = contextutils.WithLogger(ctx, "envoy_validation")
	logger = contextutils.LoggerFrom(ctx)

	// If full envoy validation is disabled, skip validation
	if !t.settings.GetGateway().GetValidation().GetFullEnvoyValidation().GetValue() {
		logComputeTranslator(logger, "Skipping full Envoy validation",
			"proxyName", proxy.GetMetadata().GetName())
		return xdsSnapshot, reports, proxyReport
	}

	logComputeTranslator(logger, "Running full Envoy validation",
		"proxyName", proxy.GetMetadata().GetName())

	if err := envoyvalidation.ValidateSnapshot(ctx, xdsSnapshot); err != nil {
		logComputeTranslator(logger, "Full Envoy validation failed",
			"proxyName", proxy.GetMetadata().GetName(),
			"error", err.Error())
		reports.AddError(proxy, err)
	}

	logComputeTranslator(logger, "Proxy translation completed successfully",
		"proxyName", proxy.GetMetadata().GetName(),
		"hasErrors", reports.ValidateStrict() != nil,
		"hasWarnings", reports.Validate() != nil)

	return xdsSnapshot, reports, proxyReport
}

func (t *translatorInstance) translateClusterSubsystemComponents(params plugins.Params, proxy *v1.Proxy, reports reporter.ResourceReports) (
	[]*envoy_config_cluster_v3.Cluster,
	[]*envoy_config_endpoint_v3.ClusterLoadAssignment,
) {
	logger := contextutils.LoggerFrom(params.Ctx)

	logger.Debugf("verifying upstream groups: %v", proxy.GetMetadata().GetName())
	t.verifyUpstreamGroups(params, reports)

	upstreamRefKeyToEndpoints := createUpstreamToEndpointsMap(params.Snapshot.Upstreams, params.Snapshot.Endpoints)

	// endpoints and clusters are shared between listeners
	logger.Debugf("computing envoy clusters for proxy: %v", proxy.GetMetadata().GetName())
	clusters, clusterToUpstreamMap := t.computeClusters(params, reports, upstreamRefKeyToEndpoints, proxy)
	logger.Debugf("computing envoy endpoints for proxy: %v", proxy.GetMetadata().GetName())

	endpoints := t.computeClusterEndpoints(params, upstreamRefKeyToEndpoints, reports)

	upstreamMap := make(map[string]struct{}, len(params.Snapshot.Upstreams))
	// make sure to call EndpointPlugin with empty endpoint
	for _, upstream := range params.Snapshot.Upstreams {
		key := UpstreamToClusterName(&core.ResourceRef{
			Name:      upstream.GetMetadata().GetName(),
			Namespace: upstream.GetMetadata().GetNamespace(),
		})
		upstreamMap[key] = struct{}{}
	}
	endpointMap := make(map[string][]*envoy_config_endpoint_v3.ClusterLoadAssignment, len(endpoints))
	for _, ep := range endpoints {
		if _, ok := endpointMap[ep.GetClusterName()]; !ok {
			endpointMap[ep.GetClusterName()] = []*envoy_config_endpoint_v3.ClusterLoadAssignment{ep}
		} else {
			// TODO: should check why has duplicated upstream
			endpointMap[ep.GetClusterName()] = append(endpointMap[ep.GetClusterName()], ep)
		}
	}
	// Find all the EDS clusters without endpoints (can happen with kube service that have no endpoints), and create a zero sized load assignment
	// this is important as otherwise envoy will wait for them forever wondering their fate and not doing much else.
ClusterLoop:
	for _, c := range clusters {
		if c.GetType() != envoy_config_cluster_v3.Cluster_EDS {
			continue
		}
		// get upstream that generated this cluster
		upstream := clusterToUpstreamMap[c]
		endpointClusterName, err := GetEndpointClusterName(c.GetName(), upstream)
		if err != nil {
			reports.AddError(upstream, errors.Wrapf(err, "could not marshal upstream to JSON"))
		}
		// Workaround for envoy bug: https://github.com/envoyproxy/envoy/issues/13009
		// Change the cluster eds config, forcing envoy to re-request latest EDS config
		c.GetEdsClusterConfig().ServiceName = endpointClusterName
		if eList, ok := endpointMap[c.GetName()]; ok {
			for _, ep := range eList {
				// the endpoint ClusterName needs to match the cluster's EdsClusterConfig ServiceName
				ep.ClusterName = endpointClusterName
			}
			continue ClusterLoop
		}
		emptyEndpointList := &envoy_config_endpoint_v3.ClusterLoadAssignment{
			ClusterName: endpointClusterName,
		}
		// make sure to call EndpointPlugin with empty endpoint
		if _, ok := upstreamMap[c.GetName()]; ok {
			for _, plugin := range t.pluginRegistry.GetEndpointPlugins() {
				if err := plugin.ProcessEndpoints(params, upstream, emptyEndpointList); err != nil {
					reports.AddError(upstream, err)
				}
			}
		}
		if _, ok := endpointMap[emptyEndpointList.GetClusterName()]; !ok {
			endpointMap[emptyEndpointList.GetClusterName()] = []*envoy_config_endpoint_v3.ClusterLoadAssignment{emptyEndpointList}
		} else {
			endpointMap[emptyEndpointList.GetClusterName()] = append(endpointMap[emptyEndpointList.GetClusterName()], emptyEndpointList)
		}
		endpoints = append(endpoints, emptyEndpointList)
	}

	return clusters, endpoints
}

func (t *translatorInstance) translateListenerSubsystemComponents(params plugins.Params, proxy *v1.Proxy, proxyReport *validationapi.ProxyReport) (
	[]*envoy_config_route_v3.RouteConfiguration,
	[]*envoy_config_listener_v3.Listener,
) {
	var (
		routeConfigs []*envoy_config_route_v3.RouteConfiguration
		listeners    []*envoy_config_listener_v3.Listener
	)

	logger := contextutils.LoggerFrom(params.Ctx)

	logger.Infow("Starting listener subsystem translation",
		"issue", "8539",
		"proxy_name", proxy.GetMetadata().GetName(),
		"proxy_namespace", proxy.GetMetadata().GetNamespace(),
		"listener_count", len(proxy.GetListeners()))

	for i, listener := range proxy.GetListeners() {
		logger.Infof("computing envoy resources for listener: %v", listener.GetName())
		logger.Infow("Processing listener for translation",
			"issue", "8539",
			"listener_name", listener.GetName(),
			"listener_index", i,
			"listener_type", fmt.Sprintf("%T", listener.GetListenerType()),
			"bind_address", listener.GetBindAddress(),
			"bind_port", listener.GetBindPort())

		listenerReport := proxyReport.GetListenerReports()[i]

		// TODO: This only needs to happen once, we should move it out of the loop
		validateListenerPorts(proxy, listenerReport)

		// Select a ListenerTranslator and RouteConfigurationTranslator, based on the type of listener (ie TCP, HTTP, Hybrid, or Aggregate)
		listenerTranslator, routeConfigurationTranslator := t.listenerTranslatorFactory.GetTranslators(params.Ctx, proxy, listener, listenerReport)

		logComputeTranslator(logger, "Selected translators for listener",
			"listener_name", listener.GetName(),
			"listener_translator_type", fmt.Sprintf("%T", listenerTranslator),
			"route_config_translator_type", fmt.Sprintf("%T", routeConfigurationTranslator))

		// 1. Compute RouteConfiguration
		// This way we call ProcessVirtualHost / ProcessRoute first
		logComputeTranslator(logger, "Computing route configuration",
			"issue", "8539",
			"listener_name", listener.GetName())
		envoyRouteConfiguration := routeConfigurationTranslator.ComputeRouteConfiguration(params)

		// 2. Compute Listener
		// This way we evaluate HttpFilters second, which allows us to avoid appending an HttpFilter
		// that is not used by any Route / VirtualHost
		logComputeTranslator(logger, "Computing envoy listener",
			"issue", "8539",
			"listener_name", listener.GetName())
		envoyListener := listenerTranslator.ComputeListener(params)

		if envoyListener != nil {
			listeners = append(listeners, envoyListener)
			if len(envoyRouteConfiguration) > 0 {
				routeConfigs = append(routeConfigs, envoyRouteConfiguration...)
			}
			logComputeTranslator(logger, "Successfully computed listener and routes",
				"issue", "8539",
				"listener_name", listener.GetName(),
				"envoy_listener_name", envoyListener.GetName(),
				"route_config_count", len(envoyRouteConfiguration))
		} else {
			logComputeTranslator(logger, "Listener translation returned nil",
				"issue", "8539",
				"listener_name", listener.GetName())
		}
	}

	logger.Infow("Completed listener subsystem translation",
		"issue", "8539",
		"total_listeners_created", len(listeners),
		"total_route_configs_created", len(routeConfigs))

	return routeConfigs, listeners
}

func (t *translatorInstance) generateXDSSnapshot(
	params plugins.Params,
	clusters []*envoy_config_cluster_v3.Cluster,
	endpoints []*envoy_config_endpoint_v3.ClusterLoadAssignment,
	routeConfigs []*envoy_config_route_v3.RouteConfiguration,
	listeners []*envoy_config_listener_v3.Listener,
) envoycache.Snapshot {
	var endpointsProto, clustersProto, listenersProto []envoycache.Resource

	for _, ep := range endpoints {
		endpointsProto = append(endpointsProto, resource.NewEnvoyResource(ep))
	}
	for _, cluster := range clusters {
		clustersProto = append(clustersProto, resource.NewEnvoyResource(cluster))
	}
	for _, listener := range listeners {
		// don't add empty listeners, envoy will complain
		if len(listener.GetFilterChains()) < 1 {
			continue
		}
		listenersProto = append(listenersProto, resource.NewEnvoyResource(listener))
	}
	// construct version
	// TODO: investigate whether we need a more sophisticated versioning algorithm
	endpointsVersion, endpointsErr := t.hasher(endpointsProto)
	if endpointsErr != nil {
		contextutils.LoggerFrom(params.Ctx).DPanic(fmt.Sprintf("error trying to hash endpointsProto: %v", endpointsErr))
	}
	clustersVersion, clustersErr := t.hasher(clustersProto)
	if clustersErr != nil {
		contextutils.LoggerFrom(params.Ctx).DPanic(fmt.Sprintf("error trying to hash clustersProto: %v", clustersErr))
	}
	listenersVersion, listenersErr := t.hasher(listenersProto)
	if listenersErr != nil {
		contextutils.LoggerFrom(params.Ctx).DPanic(fmt.Sprintf("error trying to hash listenersProto: %v", listenersErr))
	}

	// if clusters are updated, provider a new version of the endpoints,
	// so the clusters are warm
	endpointsNew := envoycache.NewResources(fmt.Sprintf("%v-%v", clustersVersion, endpointsVersion), endpointsProto)
	if endpointsErr != nil || clustersErr != nil {
		endpointsNew = envoycache.NewResources("endpoints-hashErr", endpointsProto)
	}
	clustersNew := envoycache.NewResources(fmt.Sprintf("%v", clustersVersion), clustersProto)
	if clustersErr != nil {
		clustersNew = envoycache.NewResources("clusters-hashErr", endpointsProto)
	}
	listenersNew := envoycache.NewResources(fmt.Sprintf("%v", listenersVersion), listenersProto)
	if listenersErr != nil {
		listenersNew = envoycache.NewResources("listeners-hashErr", listenersProto)
	}
	return xds.NewSnapshotFromResources(
		endpointsNew,
		clustersNew,
		MakeRdsResources(routeConfigs),
		listenersNew)
}

// deprecated, use EnvoyCacheResourcesListToFnvHash
func MustEnvoyCacheResourcesListToFnvHash(resources []envoycache.Resource) uint64 {
	out, err := EnvoyCacheResourcesListToFnvHash(resources)
	if err != nil {
		contextutils.LoggerFrom(context.Background()).DPanic(err)
	}
	return out
}

func EnvoyCacheResourcesListToFnvHash(resources []envoycache.Resource) (uint64, error) {
	hasher := fnv.New64()
	// 8kb capacity, consider raising if we find the buffer is frequently being
	// re-allocated by MarshalAppend to fit larger protos.
	// the goal is to keep allocations constant for GC, without allocating an
	// unnecessarily large buffer.
	buffer := make([]byte, 0, 8*1024)
	mo := proto2.MarshalOptions{Deterministic: true}
	for _, r := range resources {
		buf := buffer[:0]
		// proto.MessageV2 will create another allocation, updating solo-kit
		// to use google protos (rather than github protos, i.e. use v2) is
		// another path to further improve performance here.
		out, err := mo.MarshalAppend(buf, proto.MessageV2(r.ResourceProto()))
		if err != nil {
			contextutils.LoggerFrom(context.Background()).DPanic(errors.Wrap(err, "marshalling envoy snapshot components"))
			return 0, errors.Wrap(err, "marshalling envoy snapshot components")
		}
		_, err = hasher.Write(out)
		if err != nil {
			contextutils.LoggerFrom(context.Background()).DPanic(errors.Wrap(err, "constructing hash for envoy snapshot components"))
			return 0, errors.Wrap(err, "constructing hash for envoy snapshot components")
		}
	}
	return hasher.Sum64(), nil
}

// deprecated, slower than EnvoyCacheResourcesListToFnvHash
func EnvoyCacheResourcesListToHash(resources []envoycache.Resource) (uint64, error) {
	return hashstructure.Hash(resources, nil)
}

func MakeRdsResources(routeConfigs []*envoy_config_route_v3.RouteConfiguration) envoycache.Resources {
	var routesProto []envoycache.Resource

	for _, routeCfg := range routeConfigs {
		// don't add empty route configs, envoy will complain
		if len(routeCfg.GetVirtualHosts()) < 1 {
			continue
		}
		routesProto = append(routesProto, resource.NewEnvoyResource(routeCfg))

	}

	routesVersion, err := EnvoyCacheResourcesListToFnvHash(routesProto)
	if err != nil {
		contextutils.LoggerFrom(context.Background()).DPanic(fmt.Sprintf("error trying to hash routesProto: %v", err))
		return envoycache.NewResources("routes-hashErr", routesProto)
	}
	return envoycache.NewResources(fmt.Sprintf("%v", routesVersion), routesProto)
}

func GetEndpointClusterName(clusterName string, upstream *v1.Upstream) (string, error) {
	hash, err := upstream.Hash(nil)
	if err != nil {
		return "", err
	}
	//note: we add the upstream hash here because of
	// https://github.com/envoyproxy/envoy/issues/13009
	endpointClusterName := fmt.Sprintf("%s-%d", clusterName, hash)
	return endpointClusterName, nil
}
