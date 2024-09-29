package kubernetes

import (
	"encoding/binary"
	"fmt"
	"hash/fnv"
	"os"
	"sort"
	"strings"

	errors "github.com/rotisserie/eris"
	corev1 "k8s.io/api/core/v1"
	discoveryv1 "k8s.io/api/discovery/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"

	"github.com/solo-io/go-utils/contextutils"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	corecache "github.com/solo-io/solo-kit/pkg/api/v1/clients/kube/cache"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"

	"github.com/solo-io/gloo/pkg/utils/settingsutil"
	"github.com/solo-io/gloo/projects/gloo/constants"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	kubeplugin "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/kubernetes"
)

type podMap struct {
	// ipMap will record pods which
	// - pod.Status.Phase is corev1.PodRunning
	// - pod.Status.PodIP != ""
	// - !pod.Spec.HostNetwork
	// and use pod.Status.PodIP as map key
	ipMap map[string]*corev1.Pod
	// metaMap will record all pods and use
	// `pod.GetName()+"/"+pod.GetNamespace()` as map key
	metaMap map[string]*corev1.Pod
}

func newPodMap(pods []*corev1.Pod) *podMap {
	podsIPMap := make(map[string]*corev1.Pod, len(pods))
	podsMetaMap := make(map[string]*corev1.Pod, len(pods))
	for _, pod := range pods {
		if pod == nil {
			continue
		}
		podsMetaMap[pod.GetName()+"/"+pod.GetNamespace()] = pod
		if pod.Status.Phase != corev1.PodRunning {
			continue
		}
		if pod.Status.PodIP == "" {
			continue
		}
		if pod.Spec.HostNetwork {
			// we can't tell pods apart if they are all on the host network.
			continue
		}
		podsIPMap[pod.Status.PodIP] = pod
	}
	return &podMap{podsIPMap, podsMetaMap}
}

func (pm *podMap) GetLabelsForIP(ip string, podName, podNamespace string) (map[string]string, error) {
	if podName != "" && podNamespace != "" {
		if p, ok := pm.metaMap[podName+"/"+podNamespace]; ok {
			return p.GetLabels(), nil
		}
	}
	if ip != "" {
		if p, ok := pm.ipMap[ip]; ok {
			return p.GetLabels(), nil
		}
	}
	return nil, errors.Errorf("running pod not found with IP %v", ip)
}

func (p *plugin) WatchEndpoints(writeNamespace string, upstreamsToTrack v1.UpstreamList, opts clients.WatchOpts) (<-chan v1.EndpointList, <-chan error, error) {
	kubeFactory := func(namespaces []string) KubePluginSharedFactory {
		return getInformerFactory(opts.Ctx, p.kube, namespaces)
	}
	watcher, err := newEndpointWatcherForUpstreams(kubeFactory, p.kubeCoreCache, upstreamsToTrack, opts, p.settings)
	if err != nil {
		return nil, nil, err
	}
	return watcher.watch(writeNamespace, opts)
}

func newEndpointWatcherForUpstreams(kubeFactoryFactory func(ns []string) KubePluginSharedFactory, kubeCoreCache corecache.KubeCoreCache, upstreamsToTrack v1.UpstreamList, opts clients.WatchOpts, settings *v1.Settings) (*edsWatcher, error) {
	var namespaces []string
	if settingsutil.IsAllNamespacesFromSettings(settings) {
		namespaces = []string{metav1.NamespaceAll}
	} else {
		nsSet := map[string]bool{}
		for _, upstream := range upstreamsToTrack {
			svcNs := upstream.GetKube().GetServiceNamespace()
			// only care about kube upstreams
			if svcNs == "" {
				continue
			}
			nsSet[svcNs] = true
		}
		for ns := range nsSet {
			namespaces = append(namespaces, ns)
		}
	}

	// If there are no upstreams to watch (eg: if discovery is disabled), namespaces remains an empty list.
	// When creating the InformerFactory, by convention, an empty namespace list means watch all namespaces.
	// To ensure that we only watch what we are supposed to, fallback to WatchNamespaces if namespaces is an empty list.
	if len(namespaces) == 0 {
		namespaces = settingsutil.GetNamespacesToWatch(settings)
	}
	kubeFactory := kubeFactoryFactory(namespaces)
	// this can take a bit of time some make sure we are still in business
	if opts.Ctx.Err() != nil {
		return nil, opts.Ctx.Err()
	}
	opts = opts.WithDefaults()

	return newEndpointsWatcher(kubeCoreCache, namespaces, kubeFactory, upstreamsToTrack), nil
}

type edsWatcher struct {
	upstreams         map[*core.ResourceRef]*kubeplugin.UpstreamSpec
	kubeShareFactory  KubePluginSharedFactory
	kubeCoreCache     corecache.KubeCoreCache
	namespaces        []string
	lastEndpointsHash uint64
}

func newEndpointsWatcher(kubeCoreCache corecache.KubeCoreCache, namespaces []string, kubeShareFactory KubePluginSharedFactory, upstreams v1.UpstreamList) *edsWatcher {
	upstreamSpecs := make(map[*core.ResourceRef]*kubeplugin.UpstreamSpec)
	for _, us := range upstreams {
		kubeUpstream, ok := us.GetUpstreamType().(*v1.Upstream_Kube)
		// only care about kube upstreams
		if !ok {
			continue
		}
		upstreamSpecs[us.GetMetadata().Ref()] = kubeUpstream.Kube
	}
	return &edsWatcher{
		upstreams:        upstreamSpecs,
		kubeShareFactory: kubeShareFactory,
		kubeCoreCache:    kubeCoreCache,
		namespaces:       namespaces,
	}
}

func (c *edsWatcher) List(writeNamespace string, opts clients.ListOpts) (v1.EndpointList, error) {
	ctx := contextutils.WithLogger(opts.Ctx, "kubernetes_eds")

	var (
		endpointList      []*corev1.Endpoints
		endpointSliceList []*discoveryv1.EndpointSlice
		serviceList       []*corev1.Service
		podList           []*corev1.Pod
		warnsToLog        []string
	)
	for _, ns := range c.namespaces {
		if c.kubeCoreCache.NamespacedServiceLister(ns) == nil {
			// this namespace is not watched, ignore it.
			warnsToLog = append(warnsToLog, fmt.Sprintf("namespace %v is not watched, and has upstreams pointing to it", ns))
			continue
		}
		services, err := c.kubeCoreCache.NamespacedServiceLister(ns).List(labels.SelectorFromSet(opts.Selector))
		if err != nil {
			return nil, err
		}
		serviceList = append(serviceList, services...)

		pods, err := c.kubeCoreCache.NamespacedPodLister(ns).List(labels.SelectorFromSet(opts.Selector))
		if err != nil {
			return nil, err
		}
		podList = append(podList, pods...)

		endpoints, err := c.kubeShareFactory.EndpointsLister(ns).List(labels.SelectorFromSet(opts.Selector))
		if err != nil {
			return nil, err
		}
		endpointList = append(endpointList, endpoints...)

		endpointSlices, err := c.kubeShareFactory.EndpointSlicesLister(ns).List(labels.SelectorFromSet(opts.Selector))
		if err != nil {
			return nil, err
		}
		endpointSliceList = append(endpointSliceList, endpointSlices...)
	}

	// TODO(tim): introduce `type endpointBuilder struct` to encapsulate the logic of building endpoints?
	// This will help cut down on the number of arguments passed to filterEndpoints.

	eps, warns, errsToLog := filterEndpoints(
		writeNamespace,
		endpointList,
		endpointSliceList,
		serviceList,
		podList,
		c.upstreams,
	)
	if len(warns) > 0 {
		warnsToLog = append(warnsToLog, warns...)
	}

	hasher := fnv.New64()
	for _, ep := range eps {
		// The "Hash" method is deprecated. See https://github.com/solo-io/solo-projects/issues/6976
		// for using the "HashUnique" method instead.
		result, err := ep.Hash(hasher)
		if err != nil {
			return nil, err
		}
		if err := binary.Write(hasher, binary.LittleEndian, result); err != nil {
			return nil, err
		}
	}
	hash := hasher.Sum64()

	if c.lastEndpointsHash == hash {
		return eps, nil
	}
	c.lastEndpointsHash = hash

	// We return and log any messages here later because endpoints are constantly being updated / triggering the
	// EDS loop by the endpoint shared informer. By ensuring the hash is different, we know the endpoints are
	// functionally different and logging any new information might mean different / new information.
	//
	// This change helps avoid disk filling up from endless logs and cluttering the logs with the same message.
	logger := contextutils.LoggerFrom(ctx)
	for _, warnToLog := range warnsToLog {
		logger.Warn(warnToLog)
	}
	for _, errToLog := range errsToLog {
		logger.Error(errToLog)
	}

	return eps, nil
}

const enableIstioSidecarOnGatewayDeprecatedWarning = "istioIntegration.enableIstioSidecarOnGateway is deprecated. Use istioSDS.enabled instead."

// isIstioInjectionEnabled returns true if Istio integration is enabled, where endpoints
// must be defined by IP address rather than hostnames. For more details, see:
// https://github.com/solo-io/gloo/issues/6195
func isIstioInjectionEnabled() (bool, []string) {
	lookupResult, found := os.LookupEnv(constants.IstioInjectionEnabled)
	istioIsEnabled := found && strings.EqualFold(lookupResult, "true")
	if istioIsEnabled {
		return true, []string{enableIstioSidecarOnGatewayDeprecatedWarning}
	}
	return false, nil
}

func (c *edsWatcher) watch(writeNamespace string, opts clients.WatchOpts) (<-chan v1.EndpointList, <-chan error, error) {
	watch := c.kubeShareFactory.Subscribe()

	endpointsChan := make(chan v1.EndpointList)
	errs := make(chan error)

	updateResourceList := func() {
		list, err := c.List(writeNamespace, clients.ListOpts{
			Ctx:      opts.Ctx,
			Selector: opts.Selector,
		})
		if err != nil {
			errs <- err
			return
		}
		select {
		case <-opts.Ctx.Done():
			return
		case endpointsChan <- list:
		}
	}

	go func() {
		defer c.kubeShareFactory.Unsubscribe(watch)
		defer close(endpointsChan)
		defer close(errs)

		updateResourceList()
		for {
			select {
			case _, ok := <-watch:
				if !ok {
					return
				}
				updateResourceList()
			case <-opts.Ctx.Done():
				return
			}
		}
	}()

	return endpointsChan, errs, nil
}

// epKey represents a unique key for an endpoint. This key is used to group
// endpoints by address and port with additional metadata to identify the
// upstream and service.
type epKey struct {
	address     string
	port        uint32
	name        string
	namespace   string
	upstreamRef *core.ResourceRef
	// While we can use the upstream ref to get the upstream and service, if there
	// are too many upstreams that could be slow.
	isHeadless bool
}

// findPortFromServices returns first matching port by searching through the list
// of services according to the given upstream spec. When we fail to find a matching
// port, we return nil and false.
func findPortFromServices(services []*corev1.Service, spec *kubeplugin.UpstreamSpec) (*corev1.ServicePort, bool) {
	for _, svc := range services {
		if svc.Namespace != spec.GetServiceNamespace() || svc.Name != spec.GetServiceName() {
			continue
		}
		for _, port := range svc.Spec.Ports {
			if spec.GetServicePort() == uint32(port.Port) {
				return &port, len(svc.Spec.Ports) == 1
			}
		}
	}
	return nil, false
}

// getServiceFromUpstreamSpec returns the Service object that matches the given
// upstream spec. If no matching service is found, we return nil, so the caller
// needs to handle this case.
func getServiceFromUpstreamSpec(spec *kubeplugin.UpstreamSpec, services []*corev1.Service) *corev1.Service {
	for _, svc := range services {
		if svc.Namespace == spec.GetServiceNamespace() && svc.Name == spec.GetServiceName() {
			return svc
		}
	}
	return nil
}

// filterEndpoints computes all the required endpoints for Gloo from the given Kubernetes
// endpoints/endpointslices, services, and Gloo upstreams that are being tracked. It returns
// the computed endpoints, and any warnings or errors encountered during the computation.
func filterEndpoints(
	writeNamespace string,
	endpoints []*corev1.Endpoints,
	endpointSlices []*discoveryv1.EndpointSlice,
	services []*corev1.Service,
	pods []*corev1.Pod,
	upstreams map[*core.ResourceRef]*kubeplugin.UpstreamSpec,
) (v1.EndpointList, []string, []string) {
	var warnsToLog, errorsToLog []string
	istioInjectionEnabled, warnings := isIstioInjectionEnabled()
	if len(warnings) > 0 {
		warnsToLog = append(warnsToLog, warnings...)
	}

	podLabelSource := newPodMap(pods)
	endpointsMap := make(map[epKey][]*core.ResourceRef)

	for usRef, spec := range upstreams {
		kubeServicePort, singlePortService := findPortFromServices(services, spec)
		if kubeServicePort == nil {
			errorsToLog = append(errorsToLog, fmt.Sprintf("upstream %v: port %v not found for service %v",
				usRef.Key(),
				spec.GetServicePort(),
				spec.GetServiceName(),
			))
			continue
		}

		svc := getServiceFromUpstreamSpec(spec, services)
		if svc == nil {
			errorsToLog = append(errorsToLog, fmt.Sprintf("upstream %v: service %v not found",
				usRef.Key(),
				spec.GetServiceName(),
			))
			continue
		}

		isHeadlessSvc := svc.Spec.ClusterIP == "None"
		// Istio uses the service's port for routing requests
		// Headless services don't have a cluster IP, so we'll resort to pod IP endpoints
		if istioInjectionEnabled && !isHeadlessSvc {
			hostname := fmt.Sprintf("%v.%v", spec.GetServiceName(), spec.GetServiceNamespace())
			copyRef := *usRef
			key := epKey{
				address:     hostname,
				port:        uint32(kubeServicePort.Port),
				name:        spec.GetServiceName(),
				namespace:   spec.GetServiceNamespace(),
				upstreamRef: &copyRef,
				isHeadless:  isHeadlessSvc,
			}
			endpointsMap[key] = append(endpointsMap[key], &copyRef)
			continue
		}

		// TODO(tim): can we introduce a feature flag that allows users to opt-in to
		// using endpoint slices? we'd graduate this to a stable feature in the future.

		// find each matching endpoint slice. we'll then use the address and port
		// fields from the endpoint slice to create the endpoint(s) needs for EDS.
		for _, eps := range endpointSlices {
			label, ok := eps.GetLabels()[discoveryv1.LabelServiceName]
			if !ok || label != spec.GetServiceName() {
				// filter out endpoint slices that don't match the service name.
				// this label is used to link an EndpointSlices to a Service and
				// is either set manually by the user or by a controller.
				continue
			}

			// we need to find the port that matches the configured Upstream service port
			port := findFirstPortInEndpointSlice(eps.Ports, singlePortService, kubeServicePort)
			if port == 0 {
				warnsToLog = append(warnsToLog, fmt.Sprintf("upstream %v: port %v not found for service %v in endpoint slice %v",
					usRef.Key(),
					spec.GetServicePort(),
					spec.GetServiceName(),
					eps.GetName(),
				))
				continue
			}

			for _, endpoint := range eps.Endpoints {
				if endpoint.Conditions.Ready == nil || !*endpoint.Conditions.Ready {
					// filter out non-ready endpoints from EDS
					continue
				}
				// TODO: handle targetRef pointing to a Pod like Endpoint logic does.
				for _, addr := range endpoint.Addresses {
					key := epKey{
						address:     addr,
						port:        port,
						name:        spec.GetServiceName(),
						namespace:   spec.GetServiceNamespace(),
						upstreamRef: usRef,
						isHeadless:  isHeadlessSvc,
					}
					copyRef := *usRef
					endpointsMap[key] = append(endpointsMap[key], &copyRef)
				}
			}
		}

		// find each matching endpoint
		//
		// TODO(tim): introduce additional filtering logic that filters out Endpoints that
		// are managed by kube-controller-manager (or third-party controllers) if we don't
		// pursue the feature flag approach? There doesn't seem to easily identify Endpoint
		// resources that are manually managed (i.e. not managed via an EndpointSlice) without
		// introducing more code complexity. We could possibly infer this via the skip-mirror
		// label, but kube has a bunch of other caveats on what it considers a managed resource.
		for _, eps := range endpoints {
			if eps.Namespace != spec.GetServiceNamespace() || eps.Name != spec.GetServiceName() {
				continue
			}
			for _, subset := range eps.Subsets {
				port := findFirstPortInEndpointSubsets(subset.Ports, singlePortService, kubeServicePort)
				if port == 0 {
					warnsToLog = append(warnsToLog, fmt.Sprintf("upstream %v: port %v not found for service %v in endpoint %v",
						usRef.Key(),
						spec.GetServicePort(),
						spec.GetServiceName(),
						subset,
					))
					continue
				}

				warnings := processSubsetAddresses(
					subset.Addresses,
					spec,
					podLabelSource,
					usRef,
					port,
					endpointsMap,
					isHeadlessSvc,
				)
				if len(warnings) > 0 {
					warnsToLog = append(warnsToLog, warnings...)
				}
			}
		}
	}

	endpointList := generateFilteredEndpointList(
		endpointsMap,
		services,
		podLabelSource,
		writeNamespace,
		istioInjectionEnabled,
	)

	return endpointList, warnsToLog, errorsToLog
}

func processSubsetAddresses(
	addresses []corev1.EndpointAddress,
	spec *kubeplugin.UpstreamSpec,
	pods *podMap,
	usRef *core.ResourceRef,
	port uint32,
	endpointsMap map[epKey][]*core.ResourceRef,
	isHeadlessService bool,
) []string {
	var warnings []string
	for _, addr := range addresses {
		var (
			podName      string
			podNamespace string
		)
		if addr.TargetRef != nil && addr.TargetRef.Kind == "Pod" {
			podName = addr.TargetRef.Name
			podNamespace = addr.TargetRef.Namespace
		}
		if len(spec.GetSelector()) != 0 {
			// determine whether labels for the owner of this ip (pod) matches the spec
			podLabels, err := pods.GetLabelsForIP(addr.IP, podName, podNamespace)
			if err != nil {
				warnings = append(warnings, fmt.Sprintf("error for upstream %v service %v: %v",
					usRef.Key(),
					spec.GetServiceName(),
					err,
				))
				continue
			}
			if !labels.SelectorFromSet(spec.GetSelector()).Matches(labels.Set(podLabels)) {
				continue
			}
			// pod hasn't been assigned address yet
			if addr.IP == "" {
				continue
			}
		}
		key := epKey{
			address:     addr.IP,
			port:        port,
			name:        podName,
			namespace:   podNamespace,
			upstreamRef: usRef,
			isHeadless:  isHeadlessService,
		}
		copyRef := *usRef
		endpointsMap[key] = append(endpointsMap[key], &copyRef)
	}
	return warnings
}

func findFirstPortInEndpointSubsets(
	ports []corev1.EndpointPort,
	singlePortService bool,
	kubeServicePort *corev1.ServicePort,
) uint32 {
	var port uint32
	for _, p := range ports {
		// if the endpoint port is not named, it implies that
		// the kube service only has a single unnamed port as well.
		switch {
		case singlePortService:
			port = uint32(p.Port)
		case p.Name == kubeServicePort.Name:
			port = uint32(p.Port)
			break
		}
	}
	return port
}

func findFirstPortInEndpointSlice(
	ports []discoveryv1.EndpointPort,
	singlePortService bool,
	kubeServicePort *corev1.ServicePort,
) uint32 {
	var port uint32
	for _, p := range ports {
		// if the endpoint port is not named, it implies that
		// the kube service only has a single unnamed port as well.
		switch {
		case singlePortService:
			port = uint32(*p.Port)
		case *p.Name == kubeServicePort.Name:
			port = uint32(*p.Port)
			break
		}
	}
	return port
}

func generateFilteredEndpointList(
	endpointsMap map[epKey][]*core.ResourceRef,
	services []*corev1.Service,
	pods *podMap,
	writeNamespace string,
	istioIntegrationEnabled bool,
) v1.EndpointList {
	var endpoints v1.EndpointList
	for addr, refs := range endpointsMap {
		// sort refs for idempotency
		sort.Slice(refs, func(i, j int) bool { return refs[i].Key() < refs[j].Key() })
		hasher := fnv.New64()
		hasher.Write([]byte(fmt.Sprintf("%+v", addr)))
		dnsname := strings.Map(func(r rune) rune {
			if '0' <= r && r <= '9' {
				return r
			}
			if 'a' <= r && r <= 'z' {
				return r
			}
			return '-'
		}, addr.address)
		endpointName := fmt.Sprintf("ep-%v-%v-%x", dnsname, addr.port, hasher.Sum64())

		var ep *v1.Endpoint
		// While istio integration requires the Service VIP, headless services require the pod IP,
		// as there is no Cluster IP. Note: If istioIntegrationEnabled is enabled, Istio automtls
		// will not be able to generate the endpoint metadata from the Pod to match the transport
		// socket match.
		if istioIntegrationEnabled && !addr.isHeadless {
			// Istio integration requires assigning endpoints the Kub service VIP rather than pod address
			service, _ := getServiceForHostname(addr.address, addr.name, addr.namespace, services)
			ep = createEndpoint(writeNamespace, endpointName, refs, service.Spec.ClusterIP, addr.port, service.GetObjectMeta().GetLabels()) // TODO: labels may be nil
		} else {
			podLabels, _ := pods.GetLabelsForIP(addr.address, addr.name, addr.namespace)
			ep = createEndpoint(writeNamespace, endpointName, refs, addr.address, addr.port, podLabels)
		}
		endpoints = append(endpoints, ep)
	}
	// sort refs for idempotency
	sort.Slice(endpoints, func(i, j int) bool {
		return endpoints[i].GetMetadata().GetName() < endpoints[j].GetMetadata().GetName()
	})
	return endpoints
}

func createEndpoint(
	namespace,
	name string,
	upstreams []*core.ResourceRef,
	address string,
	port uint32,
	labels map[string]string,
) *v1.Endpoint {
	return &v1.Endpoint{
		Metadata: &core.Metadata{
			Namespace: namespace,
			Name:      name,
			Labels:    labels,
		},
		Upstreams: upstreams,
		Address:   address,
		Port:      port,
		// TODO: add locality info
	}
}

func getServiceForHostname(hostname string, serviceName, serviceNamespace string, services []*corev1.Service) (*corev1.Service, error) {
	for _, service := range services {
		if serviceName != "" && serviceNamespace != "" {
			// no need for heuristics!
			// Note(tim): what does that comment mean?
			if serviceName == service.Name && serviceNamespace == service.Namespace {
				return service, nil
			}
		}
		serviceHostname := fmt.Sprintf("%v.%v", service.Name, service.Namespace)
		if serviceHostname != hostname {
			continue
		}
		return service, nil
	}
	return nil, errors.Errorf("service not found with hostname %v", hostname)
}
