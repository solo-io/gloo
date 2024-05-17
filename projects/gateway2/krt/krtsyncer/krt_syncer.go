package krttranslator

import (
	"context"
	"strconv"

	"github.com/solo-io/gloo/projects/gateway2/extensions"
	"github.com/solo-io/gloo/projects/gateway2/krt/krtquery"
	"github.com/solo-io/gloo/projects/gateway2/reports"
	"github.com/solo-io/gloo/projects/gateway2/translator"
	gwplugins "github.com/solo-io/gloo/projects/gateway2/translator/plugins"
	"github.com/solo-io/gloo/projects/gateway2/translator/plugins/registry"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/utils"
	"github.com/solo-io/go-utils/contextutils"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"istio.io/istio/pkg/kube"
	"istio.io/istio/pkg/kube/controllers"
	"istio.io/istio/pkg/kube/krt"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	gwapi "sigs.k8s.io/gateway-api/apis/v1"
)

// QueueStatusForProxiesFn queues a list of proxies to be synced and the plugin registry that produced them for a given sync iteration
type QueueStatusForProxiesFn func(proxies gloov1.ProxyList, pluginRegistry *registry.PluginRegistry, totalSyncCount int)

type Syncer struct {
	// regular controller runtime mgr
	mgr manager.Manager
	// krt stuff
	client  kube.Client
	queries krtquery.Queries
	proxies krt.Collection[ProxyGateway]

	// settings
	controllerName  string
	writeNamespace  string
	k8sGwExtensions extensions.K8sGatewayExtensions

	totalResyncs int

	// proxyReconciler wraps the client that writes Proxy resources into an in-memory cache
	// This cache is utilized by the debug.ProxyEndpointServer
	proxyReconciler gloov1.ProxyReconciler

	// queueStatusForProxies stores a list of proxies that need the proxy status synced and the plugin registry
	// that produced them for a given sync iteration
	queueStatusForProxies QueueStatusForProxiesFn
}

type ProxyGateway struct {
	Proxy   *gloov1.Proxy
	Gateway *gwapi.Gateway
}

func New(
  ctx context.Context, 
	controllerName, writeNamespace string,
	// inputs *GatewayInputChannels,
	mgr manager.Manager,
	k8sGwExtensions extensions.K8sGatewayExtensions,
	proxyClient gloov1.ProxyClient,
	queueStatusForProxies QueueStatusForProxiesFn,
) (*Syncer, error) {
  // currently krt requires istio's kube client
  cfg := kube.NewClientConfigForRestConfig(mgr.GetConfig())
	client, err := kube.NewClient(cfg, "gloo-gateway")
	if err != nil {
		return nil, err
	}

  // all external data access for the translator is done through queries
	queries, err := krtquery.New(client)
	if err != nil {
		return nil, err
	}

	Proxies := krt.NewCollection(
		queries.Gateways,
		func(krtctx krt.HandlerContext, i *gwapi.Gateway) *ProxyGateway {
      // TODO get translator from extensions when that PR is in
      translator.NewTranslator(queries, k8sGwExtensions.CreatePluginRegistry(ctx))
			return nil
		},
	)

	return &Syncer{
		mgr:             nil,
		client:          client,
		queries:         queries,
		proxies:         Proxies,
		controllerName:  "",
		writeNamespace:  "",
		k8sGwExtensions: k8sGwExtensions,
		totalResyncs:    0,
		proxyReconciler: nil,
		queueStatusForProxies: func(proxies gloov1.ProxyList, pluginRegistry *registry.PluginRegistry, totalSyncCount int) {
		},
	}, nil
}

func (s *Syncer) Run(ctx context.Context) {
	s.proxies.RegisterBatch(func(events []krt.Event[ProxyGateway], initialSync bool) {
		s.resyncProxies(ctx, events)
	}, false)

	s.client.RunAndWait(ctx.Done())
}

func (s *Syncer) resyncProxies(ctx context.Context, events []krt.Event[ProxyGateway]) {
	// if !secretsWarmed {
	// 	return
	// }
	s.totalResyncs++
	contextutils.LoggerFrom(ctx).Debugf("resyncing k8s gateway proxies [%v]", s.totalResyncs)

	var (
		proxies            gloov1.ProxyList
		translatedGateways []gwplugins.TranslatedGateway
		// TODO merge reports from events
		rm = reports.NewReportMap()
	)

	for _, ev := range events {
		if ev.Event == controllers.EventDelete {
			// TODO
			continue
		}
		proxy := ev.New.Proxy
		gw := ev.New.Gateway
		// Add proxy id to the proxy metadata to track proxies for status reporting
		proxyAnnotations := proxy.GetMetadata().GetAnnotations()
		if proxyAnnotations == nil {
			proxyAnnotations = make(map[string]string)
		}
		proxyAnnotations[utils.ProxySyncId] = strconv.Itoa(s.totalResyncs)
		proxy.GetMetadata().Annotations = proxyAnnotations

		proxies = append(proxies, proxy)
		translatedGateways = append(translatedGateways, gwplugins.TranslatedGateway{
			Gateway: *gw,
		})
	}

	pluginRegistry := s.k8sGwExtensions.CreatePluginRegistry(ctx)
	applyPostTranslationPlugins(ctx, pluginRegistry, &gwplugins.PostTranslationContext{
		TranslatedGateways: translatedGateways,
	})

	s.queueStatusForProxies(proxies, &pluginRegistry, s.totalResyncs)
	s.syncStatus(ctx, rm, translatedGateways)
	s.syncRouteStatus(ctx, rm)
	s.reconcileProxies(ctx, proxies)
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

// TODO krt-ify the entire reporting setup
func (s *Syncer) syncRouteStatus(ctx context.Context, rm reports.ReportMap) {
	ctx = contextutils.WithLogger(ctx, "routeStatusSyncer")
	logger := contextutils.LoggerFrom(ctx)
	rl := gwapi.HTTPRouteList{}
	err := s.mgr.GetClient().List(ctx, &rl)
	if err != nil {
		logger.Error(err)
		return
	}

	for _, route := range rl.Items {
		route := route // pike
		if status := rm.BuildRouteStatus(ctx, route, s.controllerName); status != nil {
			route.Status = *status
			if err := s.mgr.GetClient().Status().Update(ctx, &route); err != nil {
				logger.Error(err)
			}
		}
	}
}

// syncStatus updates the status of the Gateway CRs
func (s *Syncer) syncStatus(ctx context.Context, rm reports.ReportMap, gws []gwplugins.TranslatedGateway) {
	ctx = contextutils.WithLogger(ctx, "statusSyncer")
	logger := contextutils.LoggerFrom(ctx)
	for _, gw := range gws {
		gw := gw.Gateway // pike
		if status := rm.BuildGWStatus(ctx, gw); status != nil {
			gw.Status = *status
			if err := s.mgr.GetClient().Status().Patch(ctx, &gw, client.Merge); err != nil {
				logger.Error(err)
			}
		}
	}
}

// reconcileProxies persists the proxies that were generated during translations and stores them in an in-memory cache
// This cache is utilized by the debug.ProxyEndpointServer
// As well as to resync the Gloo Xds Translator (when it receives new proxies using a MultiResourceClient)
func (s *Syncer) reconcileProxies(ctx context.Context, proxyList gloov1.ProxyList) {
	ctx = contextutils.WithLogger(ctx, "proxyCache")
	logger := contextutils.LoggerFrom(ctx)

	// Proxy CR is located in the writeNamespace, which may be different from the originating Gateway CR
	err := s.proxyReconciler.Reconcile(
		s.writeNamespace,
		proxyList,
		func(original, desired *gloov1.Proxy) (bool, error) {
			// ignore proxies that do not have our owner label
			if original.GetMetadata().GetLabels() == nil || original.GetMetadata().GetLabels()[utils.ProxyTypeKey] != utils.GatewayApiProxyValue {
				// TODO(npolshak): Currently we update all Gloo Gateway proxies. We should create a new label and ignore proxies that are not owned by Gloo control plane running in a specific namespace via POD_NAMESPACE
				logger.Debugf("ignoring proxy %v in namespace %v, does not have owner label %v", original.GetMetadata().GetName(), original.GetMetadata().GetNamespace(), utils.GatewayApiProxyValue)
				return false, nil
			}
			// otherwise always update
			return true, nil
		},
		clients.ListOpts{
			Ctx: ctx,
		})
	if err != nil {
		// A write error to our cache should not impact translation
		// We will emit a message, and continue
		logger.Error(err)
	}
}
