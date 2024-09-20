package proxy_syncer

import (
	"context"
	"errors"
	"strconv"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"

	sologatewayv1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1/kube/apis/gateway.solo.io/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/gloosnapshot"
	glookubev1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/kube/apis/gloo.solo.io/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins/kubernetes"

	// solokubeclient "github.com/solo-io/solo-kit/pkg/api/v1/clients/kube"
	// solokubecrd "github.com/solo-io/solo-kit/pkg/api/v1/clients/kube/crd/solo.io/v1"
	// "github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"istio.io/istio/pkg/kube"
	"istio.io/istio/pkg/kube/kclient"
	"istio.io/istio/pkg/kube/krt"
	"istio.io/istio/pkg/kube/kubetypes"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/golang/protobuf/proto"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/solo-io/gloo/pkg/utils/statsutils"
	gatewayv1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gateway2/extensions"
	"github.com/solo-io/gloo/projects/gateway2/reports"
	gwplugins "github.com/solo-io/gloo/projects/gateway2/translator/plugins"
	"github.com/solo-io/gloo/projects/gateway2/translator/plugins/registry"
	"github.com/solo-io/gloo/projects/gateway2/translator/translatorutils"
	kubeconverters "github.com/solo-io/gloo/projects/gloo/pkg/api/converters/kube"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	kubeplugin "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/kubernetes"
	"github.com/solo-io/gloo/projects/gloo/pkg/syncer"
	"github.com/solo-io/gloo/projects/gloo/pkg/translator"
	kubeupstreams "github.com/solo-io/gloo/projects/gloo/pkg/upstreams/kubernetes"
	"github.com/solo-io/gloo/projects/gloo/pkg/utils"
	"github.com/solo-io/go-utils/contextutils"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	envoycache "github.com/solo-io/solo-kit/pkg/api/v1/control-plane/cache"

	// "github.com/solo-io/solo-kit/pkg/utils/kubeutils"
	// "github.com/solo-io/solo-kit/pkg/utils/protoutils"
	"github.com/solo-io/solo-kit/pkg/utils/statusutils"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	gwv1 "sigs.k8s.io/gateway-api/apis/v1"
)

// ProxySyncer is responsible for translating Kubernetes Gateway CRs into Gloo Proxies
// and syncing the proxyClient with the newly translated proxies.
type ProxySyncer struct {
	controllerName string
	writeNamespace string

	inputs          *GatewayInputChannels
	mgr             manager.Manager
	k8sGwExtensions extensions.K8sGatewayExtensions

	// proxyReconciler wraps the client that writes Proxy resources into an in-memory cache
	// This cache is utilized by RateLimit and the debug.ProxyEndpointServer
	proxyReconciler gloov1.ProxyReconciler

	proxyTranslator ProxyTranslator
	istioClient     kube.Client

	// used for converting from kube type to gloo type
	// TODO: abstract away the need for this by refactoring convert func()
	// legacyClients map[reflect.Type]*solokubeclient.ResourceClient
}

type GatewayInputChannels struct {
	genericEvent AsyncQueue[struct{}]
	secretEvent  AsyncQueue[SecretInputs]
}

func (x *GatewayInputChannels) Kick(ctx context.Context) {
	x.genericEvent.Enqueue(struct{}{})
}

func (x *GatewayInputChannels) UpdateSecretInputs(ctx context.Context, inputs SecretInputs) {
	x.secretEvent.Enqueue(inputs)
}

func NewGatewayInputChannels() *GatewayInputChannels {
	return &GatewayInputChannels{
		genericEvent: NewAsyncQueue[struct{}](),
		secretEvent:  NewAsyncQueue[SecretInputs](),
	}
}

// labels used to uniquely identify Proxies that are managed by the kube gateway controller
var kubeGatewayProxyLabels = map[string]string{
	// the proxy type key/value must stay in sync with the one defined in projects/gateway2/translator/gateway_translator.go
	utils.ProxyTypeKey: utils.GatewayApiProxyValue,
}

// setupCollectionDynamic uses the dynamic client to setup an informer for a resource
// and then uses an intermediate krt collection to type the unstructured resource.
// This is a temporary workaround until we update to the latest istio version and can
// uncomment the code below for registering types.
// HACK: we don't want to use this long term, but it's letting me push forward with deveopment
func setupCollectionDynamic[T any](ctx context.Context, client kube.Client, gvr schema.GroupVersionResource, opts ...krt.CollectionOption) krt.Collection[*T] {
	gatewayClient := kclient.NewDelayedInformer[*unstructured.Unstructured](client, gvr, kubetypes.DynamicInformer, kclient.Filter{})
	GatewayMapper := krt.WrapClient(gatewayClient, opts...)
	return krt.NewCollection(GatewayMapper, func(krtctx krt.HandlerContext, i *unstructured.Unstructured) **T {
		var empty T
		out := &empty
		err := runtime.DefaultUnstructuredConverter.FromUnstructured(i.UnstructuredContent(), out)
		if err != nil {
			contextutils.LoggerFrom(ctx).DPanic("failed converting unstructured into %T: %v", empty, i)
			return nil
		}
		return &out
	})
}

// NewProxySyncer returns an implementation of the ProxySyncer
// The provided GatewayInputChannels are used to trigger syncs.
// The proxy sync is triggered by the `genericEvent` which is kicked when
// we reconcile gateway in the gateway controller. The `secretEvent` is kicked when a secret is created, updated,
func NewProxySyncer(
	controllerName string,
	writeNamespace string,
	inputs *GatewayInputChannels,
	mgr manager.Manager,
	k8sGwExtensions extensions.K8sGatewayExtensions,
	proxyClient gloov1.ProxyClient,
	upstreamClient gloov1.UpstreamClient,
	translator translator.Translator,
	xdsCache envoycache.SnapshotCache,
	settings *gloov1.Settings,
	syncerExtensions []syncer.TranslatorSyncerExtension,
) *ProxySyncer {
	restCfg := kube.NewClientConfigForRestConfig(mgr.GetConfig())
	client, err := kube.NewClient(restCfg, "")
	if err != nil {
		// TODO move this init somewhere we can handle the err
		panic(err)
	}
	kube.EnableCrdWatcher(client)

	// legacyClients := map[reflect.Type]*solokubeclient.ResourceClient{}
	// usc := upstreamClient.BaseClient()
	// if kusc, ok := usc.(*solokubeclient.ResourceClient); !ok {
	// 	panic("upstream base client isn't a kube client!!!")
	// } else {
	// 	legacyClients[reflect.TypeOf(gloov1.Upstream{})] = kusc
	// }
	return &ProxySyncer{
		controllerName:  controllerName,
		writeNamespace:  writeNamespace,
		inputs:          inputs,
		mgr:             mgr,
		k8sGwExtensions: k8sGwExtensions,
		proxyReconciler: gloov1.NewProxyReconciler(proxyClient, statusutils.NewNoOpStatusClient()),
		proxyTranslator: NewProxyTranslator(translator, xdsCache, settings, syncerExtensions),
		istioClient:     client,
		// legacyClients:   legacyClients,
	}
}

type krtCtxKey struct{}

func krtFromCtx(ctx context.Context) (krt.HandlerContext, error) {
	krtctx, ok := ctx.Value(krtCtxKey{}).(krt.HandlerContext)
	if !ok {
		return nil, errors.New("ctx does not wrap a krt.HandlerContext")
	}
	return krtctx, nil
}

func WithKrtCtx(ctx context.Context, krtctx krt.HandlerContext) context.Context {
	return context.WithValue(ctx, krtCtxKey{}, krtctx)
}

type ProxyTranslator struct {
	translator       translator.Translator
	settings         *gloov1.Settings
	syncerExtensions []syncer.TranslatorSyncerExtension
	xdsCache         envoycache.SnapshotCache
}

func NewProxyTranslator(translator translator.Translator,
	xdsCache envoycache.SnapshotCache,
	settings *gloov1.Settings,
	syncerExtensions []syncer.TranslatorSyncerExtension,
) ProxyTranslator {
	return ProxyTranslator{
		translator:       translator,
		xdsCache:         xdsCache,
		settings:         settings,
		syncerExtensions: syncerExtensions,
	}
}

var _ krt.ResourceNamer = &glooEndpoint{}

// stolen from projects/gloo/pkg/upstreams/serviceentry/krtwrappers.go
// TODO: consolidate this stuff
func UnwrapEps(geps []*glooEndpoint) gloov1.EndpointList {
	out := make(gloov1.EndpointList, 0, len(geps))
	for _, ep := range geps {
		out = append(out, ep.Endpoint)
	}
	return out
}

// glooEndpoint provides a krt keying function for Gloo's `v1.Endpoint`
type glooEndpoint struct {
	*gloov1.Endpoint
}

func (ep *glooEndpoint) Equals(in *glooEndpoint) bool {
	return proto.Equal(ep, in)
}

func (ep *glooEndpoint) ResourceName() string {
	return ep.Metadata.GetName() + "/" + ep.Metadata.GetNamespace()
}

var _ krt.ResourceNamer = &upstream{}

// upstream provides a keying function for Gloo's `v1.Upstream`
type upstream struct {
	*gloov1.Upstream
}

func (us *upstream) ResourceName() string {
	return us.Metadata.GetName() + "/" + us.Metadata.GetNamespace()
}
func (us *upstream) Equals(in *upstream) bool {
	return proto.Equal(us, in)
}

func (s *ProxySyncer) Start(ctx context.Context) error {
	ctx = contextutils.WithLogger(ctx, "k8s-gw-syncer")
	contextutils.LoggerFrom(ctx).Debug("starting syncer for k8s gateway proxies")

	// create krt collections needed for building ApiSnapshot
	RouteOptions := setupCollectionDynamic[sologatewayv1.RouteOption](
		ctx,
		s.istioClient,
		sologatewayv1.SchemeGroupVersion.WithResource("routeoptions"),
		krt.WithName("RouteOptions"),
	)
	// VHostOptions := setupCollectionDynamic[sologatewayv1.VirtualHostOption](
	// 	ctx,
	// 	s.istioClient,
	// 	sologatewayv1.SchemeGroupVersion.WithResource("virtualhostoptions"),
	// 	krt.WithName("VirtualHostOptions"),
	// )

	// TODO: handle cfgmap noisiness: https://github.com/solo-io/gloo/blob/main/projects/gloo/pkg/api/converters/kube/artifact_converter.go#L31
	configMapClient := kclient.New[*corev1.ConfigMap](s.istioClient)
	ConfigMaps := krt.WrapClient(configMapClient, krt.WithName("ConfigMaps"))

	KubeUpstreams := setupCollectionDynamic[glookubev1.Upstream](
		ctx,
		s.istioClient,
		glookubev1.SchemeGroupVersion.WithResource("upstreams"),
		krt.WithName("Upstreams"),
	)
	GlooUpstreams := krt.NewCollection(KubeUpstreams, func(kctx krt.HandlerContext, u *glookubev1.Upstream) **upstream {
		// TODO: not cloning, this is already a copy from the underlying cache, right?!
		glooUs := &u.Spec
		glooUs.Metadata = &core.Metadata{}
		glooUs.Metadata.Name = u.GetName()
		glooUs.Metadata.Namespace = u.GetNamespace()
		us := &upstream{glooUs}
		return &us
	}, krt.WithName("InMemoryUpstreams"))

	serviceClient := kclient.New[*corev1.Service](s.istioClient)
	Services := krt.WrapClient(serviceClient, krt.WithName("Services"))
	InMemUpstreams := krt.NewManyCollection(Services, func(kctx krt.HandlerContext, svc *corev1.Service) []*upstream {
		uss := []*upstream{}
		for _, port := range svc.Spec.Ports {
			us := kubeupstreams.ServiceToUpstream(ctx, svc, port)
			uss = append(uss, &upstream{us})
		}
		return uss
	})
	// TODO: get upstream collections from extensions
	FinalUpstreams := krt.JoinCollection([]krt.Collection[*upstream]{GlooUpstreams, InMemUpstreams})

	podClient := kclient.New[*corev1.Pod](s.istioClient)
	Pods := krt.WrapClient(podClient, krt.WithName("Pods"))
	epClient := kclient.New[*corev1.Endpoints](s.istioClient)
	KubeEndpoints := krt.WrapClient(epClient, krt.WithName("Endpoints"))

	GlooEndpoints := krt.NewManyCollection(FinalUpstreams, func(kctx krt.HandlerContext, us *upstream) []*glooEndpoint {
		// ripped from: projects/gloo/pkg/plugins/kubernetes/eds.go#newEndpointsWatcher(...)
		upstreamSpecs := make(map[*core.ResourceRef]*kubeplugin.UpstreamSpec)
		for _, us := range FinalUpstreams.List() {
			kubeUpstream, ok := us.Upstream.GetUpstreamType().(*gloov1.Upstream_Kube)
			// only care about kube upstreams
			if !ok {
				continue
			}
			upstreamSpecs[us.GetMetadata().Ref()] = kubeUpstream.Kube
		}
		keps := krt.Fetch(kctx, KubeEndpoints)
		svcs := krt.Fetch(kctx, Services)
		pods := krt.Fetch(kctx, Pods)
		endpoints, warns, errs := kubernetes.FilterEndpoints(
			ctx,
			"gloo-system",
			keps,
			svcs,
			pods,
			upstreamSpecs,
		)
		if len(warns) > 0 || len(errs) > 0 {
			// do something
		}
		out := make([]*glooEndpoint, 0, len(endpoints))
		for _, gep := range endpoints {
			out = append(out, &glooEndpoint{gep})
		}
		return out
	}, krt.WithName("GlooEndpoints"))

	var (
		// totalResyncs is used to track the number of times the proxy syncer has been triggered
		totalResyncs int
	)
	resyncProxies := func() {
		totalResyncs++
		contextutils.LoggerFrom(ctx).Debugf("resyncing k8s gateway proxies [%v]", totalResyncs)
		stopwatch := statsutils.NewTranslatorStopWatch("ProxySyncer")
		stopwatch.Start()
		var (
			proxies gloov1.ProxyList
		)
		defer func() {
			duration := stopwatch.Stop(ctx)
			contextutils.LoggerFrom(ctx).Debugf("translated and wrote %d proxies in %s", len(proxies), duration.String())
		}()

		var gwl gwv1.GatewayList
		err := s.mgr.GetClient().List(ctx, &gwl)
		if err != nil {
			// This should never happen, try again?
			return
		}

		pluginRegistry := s.k8sGwExtensions.CreatePluginRegistry(ctx)
		rm := reports.NewReportMap()
		r := reports.NewReporter(&rm)

		var (
			translatedGateways []gwplugins.TranslatedGateway
		)
		for _, gw := range gwl.Items {
			gatewayTranslator := s.k8sGwExtensions.GetTranslator(ctx, &gw, pluginRegistry)
			if gatewayTranslator == nil {
				contextutils.LoggerFrom(ctx).Errorf("no translator found for Gateway %s (gatewayClass %s)", gw.Name, gw.Spec.GatewayClassName)
				continue
			}
			proxy := gatewayTranslator.TranslateProxy(ctx, &gw, s.writeNamespace, r)
			if proxy != nil {
				// Add proxy id to the proxy metadata to track proxies for status reporting
				proxyAnnotations := proxy.GetMetadata().GetAnnotations()
				if proxyAnnotations == nil {
					proxyAnnotations = make(map[string]string)
				}
				proxyAnnotations[utils.ProxySyncId] = strconv.Itoa(totalResyncs)
				proxy.GetMetadata().Annotations = proxyAnnotations

				proxies = append(proxies, proxy)
				translatedGateways = append(translatedGateways, gwplugins.TranslatedGateway{
					Gateway: gw,
				})
			}
		}

		applyPostTranslationPlugins(ctx, pluginRegistry, &gwplugins.PostTranslationContext{
			TranslatedGateways: translatedGateways,
		})

		s.reconcileProxies(ctx, proxies)

		latestSnap := gloosnapshot.ApiSnapshot{}
		latestSnap.Proxies = proxies

		krtCfgMaps := ConfigMaps.List()
		as := make([]*gloov1.Artifact, 0, len(krtCfgMaps))
		for _, u := range krtCfgMaps {
			a := kubeconverters.KubeConfigMapToArtifact(u)
			as = append(as, a)
		}
		latestSnap.Artifacts = as

		kus := FinalUpstreams.List()
		upstreams := make([]*gloov1.Upstream, 0, len(kus))
		for _, u := range kus {
			upstreams = append(upstreams, u.Upstream)
		}
		latestSnap.Upstreams = upstreams

		geps := GlooEndpoints.List()
		eps := UnwrapEps(geps)
		latestSnap.Endpoints = eps

		krtRouteOpts := RouteOptions.List()
		glooRtOpts := make([]*gatewayv1.RouteOption, 0, len(krtRouteOpts))
		for _, u := range krtRouteOpts {
			glooUs := proto.Clone(&u.Spec).(*gatewayv1.RouteOption)
			glooUs.Metadata = &core.Metadata{}
			glooUs.Metadata.Name = u.GetName()
			glooUs.Metadata.Namespace = u.GetNamespace()
			glooRtOpts = append(glooRtOpts, glooUs)
		}
		latestSnap.RouteOptions = glooRtOpts

		proxiesWithReports := s.proxyTranslator.glooSync(ctx, &latestSnap)
		applyStatusPlugins(ctx, proxiesWithReports, pluginRegistry)
		s.syncStatus(ctx, rm, gwl)
		s.syncRouteStatus(ctx, rm)
	}

	go s.istioClient.RunAndWait(ctx.Done())

	// wait for caches to sync before accepting events and syncing xds
	if !s.mgr.GetCache().WaitForCacheSync(ctx) {
		return errors.New("kube gateway sync loop waiting for all caches to sync failed")
	}

	for {
		select {
		case <-ctx.Done():
			contextutils.LoggerFrom(ctx).Debug("context done, stopping proxy syncer")
			return nil
		case <-s.inputs.genericEvent.Next():
			resyncProxies()
		case <-s.inputs.secretEvent.Next():
			resyncProxies()
		}
	}
}

// func (p *ProxySyncer) convertCrdToResource(typ reflect.Type, resourceCrd *solokubecrd.Resource) (resources.Resource, error) {
// 	rc := p.legacyClients[typ]
// 	// this is the only use at this point of rc, find a better way to do this
// 	resource := rc.NewResource()
// 	resource.SetMetadata(kubeutils.FromKubeMeta(resourceCrd.ObjectMeta, true))

// 	// have to recreate, original one is private
// 	statusUnmarshaler := statusutils.NewNamespacedStatusesUnmarshaler(protoutils.UnmarshalMapToProto)
// 	if withStatus, ok := resource.(resources.InputResource); ok {
// 		statusUnmarshaler.UnmarshalStatus(resourceCrd.Status, withStatus)
// 	}
// 	// if resourceCrd.Spec != nil {
// 	// 	if err := specutils.UnmarshalSpecMapToResource(*resourceCrd.Spec, resource); err != nil {
// 	// 		// copy/paste as resourceName is private as well
// 	// 		resourceName := strings.Replace(typ.String(), "*", "", -1)
// 	// 		resourceName = strings.Replace(resourceName, ".", "", -1)
// 	// 		return nil, fmt.Errorf("unmarshal err: '%w' reading crd spec on resource %v in namespace %v into %v", err, resourceCrd.Name, resourceCrd.Namespace, resourceName)
// 	// 	}
// 	// }
// 	return resource, nil
// }

func applyStatusPlugins(
	ctx context.Context,
	proxiesWithReports []translatorutils.ProxyWithReports,
	registry registry.PluginRegistry,
) {
	ctx = contextutils.WithLogger(ctx, "k8sGatewayStatusPlugins")
	logger := contextutils.LoggerFrom(ctx)

	statusCtx := &gwplugins.StatusContext{
		ProxiesWithReports: proxiesWithReports,
	}
	for _, plugin := range registry.GetStatusPlugins() {
		err := plugin.ApplyStatusPlugin(ctx, statusCtx)
		if err != nil {
			logger.Errorf("Error applying status plugin: %v", err)
			continue
		}
	}
}

func (s *ProxySyncer) syncRouteStatus(ctx context.Context, rm reports.ReportMap) {
	ctx = contextutils.WithLogger(ctx, "routeStatusSyncer")
	logger := contextutils.LoggerFrom(ctx)
	logger.Debugf("syncing k8s gateway route status")
	stopwatch := statsutils.NewTranslatorStopWatch("HTTPRouteStatusSyncer")
	stopwatch.Start()
	defer stopwatch.Stop(ctx)

	rl := gwv1.HTTPRouteList{}
	err := s.mgr.GetClient().List(ctx, &rl)
	if err != nil {
		logger.Error(err)
		return
	}

	for _, route := range rl.Items {
		route := route // pike
		if status := rm.BuildRouteStatus(ctx, route, s.controllerName); status != nil {
			if !isHTTPRouteStatusEqual(&route.Status, status) {
				route.Status = *status
				if err := s.mgr.GetClient().Status().Update(ctx, &route); err != nil {
					logger.Error(err)
				}
			}
		}
	}
}

// syncStatus updates the status of the Gateway CRs
func (s *ProxySyncer) syncStatus(ctx context.Context, rm reports.ReportMap, gwl gwv1.GatewayList) {
	ctx = contextutils.WithLogger(ctx, "statusSyncer")
	logger := contextutils.LoggerFrom(ctx)
	stopwatch := statsutils.NewTranslatorStopWatch("GatewayStatusSyncer")
	stopwatch.Start()
	defer stopwatch.Stop(ctx)

	for _, gw := range gwl.Items {
		gw := gw // pike
		if status := rm.BuildGWStatus(ctx, gw); status != nil {
			if !isGatewayStatusEqual(&gw.Status, status) {
				gw.Status = *status
				if err := s.mgr.GetClient().Status().Patch(ctx, &gw, client.Merge); err != nil {
					logger.Error(err)
				}
			}
		}
	}
}

// reconcileProxies persists the proxies that were generated during translations and stores them in an in-memory cache
// The Gloo Xds Translator will receive these proxies via List() using a MultiResourceClient; two reasons it is needed there:
// 1. To allow Rate Limit extensions to work, as it only syncs RL configs it finds used on Proxies in the snapshots
// 2. This cache is utilized by the debug.ProxyEndpointServer
func (s *ProxySyncer) reconcileProxies(ctx context.Context, proxyList gloov1.ProxyList) {
	ctx = contextutils.WithLogger(ctx, "proxyCache")
	logger := contextutils.LoggerFrom(ctx)

	// Proxy CR is located in the writeNamespace, which may be different from the originating Gateway CR
	err := s.proxyReconciler.Reconcile(
		s.writeNamespace,
		proxyList,
		func(original, desired *gloov1.Proxy) (bool, error) {
			// only reconcile if proxies are equal
			// we reconcile so ggv2 proxies can be used in extension syncing and debug snap storage
			// but if we reconcile every time, we end in a loop where the proxies being synced here
			// trigger an apisnapshot Sync (as Proxies are in the ApiSnapshot) which will cause us to
			// regen and reconcile Proxies, in an endless loop
			// also for now, we need to translate and sync on ApiSnapshot syncs because that is the
			// source-of-truth for Gloo related translation and syncing (e.g. Endpoints)
			// if we didn't, we would need to watch for e.g. Endpoint events but there's no guarantee the
			// latest ApiSnapshot we stored would contain that latest Endpoint event
			return proto.Equal(original, desired), nil
		},
		clients.ListOpts{
			Ctx:      ctx,
			Selector: kubeGatewayProxyLabels,
		})
	if err != nil {
		// A write error to our cache should not impact translation
		// We will emit a message, and continue
		logger.Error(err)
	}
}

func applyPostTranslationPlugins(ctx context.Context, pluginRegistry registry.PluginRegistry, translationContext *gwplugins.PostTranslationContext) {
	ctx = contextutils.WithLogger(ctx, "postTranslation")
	logger := contextutils.LoggerFrom(ctx)

	for _, postTranslationPlugin := range pluginRegistry.GetPostTranslationPlugins() {
		err := postTranslationPlugin.ApplyPostTranslationPlugin(ctx, translationContext)
		if err != nil {
			logger.Errorf("Error applying post-translation plugin: %v", err)
			continue
		}
	}
}

var opts = cmp.Options{
	cmpopts.IgnoreFields(metav1.Condition{}, "LastTransitionTime"),
	cmpopts.IgnoreMapEntries(func(k string, _ any) bool {
		return k == "lastTransitionTime"
	}),
}

func isGatewayStatusEqual(objA, objB *gwv1.GatewayStatus) bool {
	return cmp.Equal(objA, objB, opts)
}

func isHTTPRouteStatusEqual(objA, objB *gwv1.HTTPRouteStatus) bool {
	return cmp.Equal(objA, objB, opts)
}
