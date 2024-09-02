package proxy_syncer

import (
	"context"
	"errors"
	"strconv"

	"github.com/golang/protobuf/proto"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/solo-io/gloo/pkg/utils/statsutils"
	"github.com/solo-io/gloo/projects/gateway2/extensions"
	"github.com/solo-io/gloo/projects/gateway2/reports"
	gwplugins "github.com/solo-io/gloo/projects/gateway2/translator/plugins"
	"github.com/solo-io/gloo/projects/gateway2/translator/plugins/registry"
	"github.com/solo-io/gloo/projects/gateway2/translator/translatorutils"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	v1snap "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/gloosnapshot"
	"github.com/solo-io/gloo/projects/gloo/pkg/syncer"
	"github.com/solo-io/gloo/projects/gloo/pkg/translator"
	"github.com/solo-io/gloo/projects/gloo/pkg/utils"
	"github.com/solo-io/go-utils/contextutils"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	envoycache "github.com/solo-io/solo-kit/pkg/api/v1/control-plane/cache"
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
}

type GatewayInputChannels struct {
	genericEvent     AsyncQueue[struct{}]
	secretEvent      AsyncQueue[SecretInputs]
	apiSnapshotEvent AsyncQueue[*v1snap.ApiSnapshot]
}

func (x *GatewayInputChannels) Kick(ctx context.Context) {
	x.genericEvent.Enqueue(struct{}{})
}

func (x *GatewayInputChannels) NewSnap(snap *v1snap.ApiSnapshot) {
	x.apiSnapshotEvent.Enqueue(snap)
}

func (x *GatewayInputChannels) UpdateSecretInputs(ctx context.Context, inputs SecretInputs) {
	x.secretEvent.Enqueue(inputs)
}

func NewGatewayInputChannels() *GatewayInputChannels {
	return &GatewayInputChannels{
		genericEvent:     NewAsyncQueue[struct{}](),
		secretEvent:      NewAsyncQueue[SecretInputs](),
		apiSnapshotEvent: NewAsyncQueue[*v1snap.ApiSnapshot](),
	}
}

// labels used to uniquely identify Proxies that are managed by the kube gateway controller
var kubeGatewayProxyLabels = map[string]string{
	// the proxy type key/value must stay in sync with the one defined in projects/gateway2/translator/gateway_translator.go
	utils.ProxyTypeKey: utils.GatewayApiProxyValue,
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
	translator translator.Translator,
	xdsCache envoycache.SnapshotCache,
	settings *gloov1.Settings,
	syncerExtensions []syncer.TranslatorSyncerExtension,
) *ProxySyncer {
	return &ProxySyncer{
		controllerName:  controllerName,
		writeNamespace:  writeNamespace,
		inputs:          inputs,
		mgr:             mgr,
		k8sGwExtensions: k8sGwExtensions,
		proxyReconciler: gloov1.NewProxyReconciler(proxyClient, statusutils.NewNoOpStatusClient()),
		proxyTranslator: NewProxyTranslator(translator, xdsCache, settings, syncerExtensions),
	}
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

func (s *ProxySyncer) Sync(ctx context.Context, snap *v1snap.ApiSnapshot) error {
	s.inputs.apiSnapshotEvent.Enqueue(snap)
	return nil
}

func (s *ProxySyncer) Start(ctx context.Context) error {
	ctx = contextutils.WithLogger(ctx, "k8s-gw-syncer")
	contextutils.LoggerFrom(ctx).Debug("starting syncer for k8s gateway proxies")

	var (
		// totalResyncs is used to track the number of times the proxy syncer has been triggered
		totalResyncs int
		latestSnap   *v1snap.ApiSnapshot
	)
	resyncProxies := func() {
		if latestSnap == nil {
			return
		}
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
		// TODO should clone?
		latestSnap.Proxies = proxies
		proxiesWithReports := s.proxyTranslator.glooSync(ctx, latestSnap)
		contextutils.LoggerFrom(ctx).Info("LAW before status plugins")
		applyStatusPlugins(ctx, proxiesWithReports, pluginRegistry)
		contextutils.LoggerFrom(ctx).Info("LAW after status plugins")
		s.syncStatus(ctx, rm, gwl)
		s.syncRouteStatus(ctx, rm)
		contextutils.LoggerFrom(ctx).Info("LAW after regular status")
	}

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
		case snap := <-s.inputs.apiSnapshotEvent.Next():
			latestSnap = snap
			resyncProxies()
		}
	}
}

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
