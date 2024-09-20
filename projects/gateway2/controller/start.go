package controller

import (
	"context"

	glooschemes "github.com/solo-io/gloo/pkg/schemes"

	"k8s.io/apimachinery/pkg/util/sets"
	ctrl "sigs.k8s.io/controller-runtime"

	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	metricsserver "sigs.k8s.io/controller-runtime/pkg/metrics/server"

	gatewayv1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gateway2/extensions"
	"github.com/solo-io/gloo/projects/gateway2/proxy_syncer"
	"github.com/solo-io/gloo/projects/gateway2/secrets"
	"github.com/solo-io/gloo/projects/gateway2/wellknown"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	api "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/extauth/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/bootstrap"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins"
	"github.com/solo-io/gloo/projects/gloo/pkg/syncer"
	"github.com/solo-io/gloo/projects/gloo/pkg/translator"
	envoycache "github.com/solo-io/solo-kit/pkg/api/v1/control-plane/cache"
	"github.com/solo-io/solo-kit/pkg/api/v2/reporter"
)

const (
	// AutoProvision controls whether the controller will be responsible for provisioning dynamic
	// infrastructure for the Gateway API.
	AutoProvision = true
)

var setupLog = ctrl.Log.WithName("setup")

type StartConfig struct {
	Dev  bool
	Opts bootstrap.Opts

	Mgr manager.Manager

	// ExtensionsFactory is the factory function which will return an extensions.K8sGatewayExtensions
	// This is responsible for producing the extension points that this controller requires
	ExtensionsFactory extensions.K8sGatewayExtensionsFactory

	// GlooPluginRegistryFactory is the factory function to produce a PluginRegistry
	// The plugins in this registry are used during the conversion of a Proxy resource into an xDS Snapshot
	GlooPluginRegistryFactory plugins.PluginRegistryFactory

	// ProxyClient is the client that writes Proxy resources into an in-memory cache
	// This cache ultimately populates the ApiSnapshot (important for extensions such as Rate Limit to work correctly)
	// and is also utilized by the debug.ProxyEndpointServer
	ProxyClient v1.ProxyClient

	// AuthConfigClient is the client used for retrieving AuthConfig objects within the Portal Plugin
	AuthConfigClient api.AuthConfigClient

	// RouteOptionClient is the client used for retrieving RouteOption objects within the RouteOptionsPlugin
	RouteOptionClient gatewayv1.RouteOptionClient

	// VirtualHostOptionClient is the client used for retrieving VirtualHostOption objects within the VirtualHostOptionsPlugin
	VirtualHostOptionClient gatewayv1.VirtualHostOptionClient

	// StatusReporter is used within any StatusPlugins that must persist a GE-classic style status
	StatusReporter reporter.StatusReporter

	Translator       translator.Translator
	XdsCache         envoycache.SnapshotCache
	Settings         *v1.Settings
	SyncerExtensions []syncer.TranslatorSyncerExtension

	ProxySyncer *proxy_syncer.ProxySyncer

	K8sGwExtensions extensions.K8sGatewayExtensions
	InputChannels   *proxy_syncer.GatewayInputChannels
}

func BuildMgr(devMode bool) (manager.Manager, error) {
	var opts []zap.Opts
	if devMode {
		setupLog.Info("starting log in dev mode")
		opts = append(opts, zap.UseDevMode(true))
	}
	ctrl.SetLogger(zap.New(opts...))

	mgrOpts := ctrl.Options{
		Scheme:           glooschemes.DefaultScheme(),
		PprofBindAddress: "127.0.0.1:9099",
		// if you change the port here, also change the port "health" in the helmchart.
		HealthProbeBindAddress: ":9093",
		Metrics: metricsserver.Options{
			BindAddress: ":9092",
		},
	}
	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), mgrOpts)
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		return nil, err
	}
	return mgr, nil
}

// Start runs the controllers responsible for processing the K8s Gateway API objects
// It is intended to be run in a goroutine as the function will block until the supplied
// context is cancelled
func Start(ctx context.Context, cfg StartConfig) error {

	// TODO: replace this with something that checks that we have xds snapshot ready (or that we don't need one).
	cfg.Mgr.AddReadyzCheck("ready-ping", healthz.Ping)

	// Create the proxy syncer for the Gateway API resources
	proxySyncer := cfg.ProxySyncer
	if err := cfg.Mgr.Add(proxySyncer); err != nil {
		setupLog.Error(err, "unable to add proxySyncer runnable")
		return err
	}

	gwCfg := GatewayConfig{
		Mgr:            cfg.Mgr,
		GWClasses:      sets.New(append(cfg.Opts.ExtraGatewayClasses, wellknown.GatewayClassName)...),
		ControllerName: wellknown.GatewayControllerName,
		AutoProvision:  AutoProvision,
		ControlPlane:   cfg.Opts.ControlPlane, // TODO(Law) type seems FAR too broad, only used to provide the xds addr
		IstioValues:    cfg.Opts.GlooGateway.IstioValues,
		Kick:           cfg.InputChannels.Kick,
		Extensions:     cfg.K8sGwExtensions,
	}
	if err := NewBaseGatewayController(ctx, gwCfg); err != nil {
		setupLog.Error(err, "unable to create controller")
		return err
	}

	if err := secrets.NewSecretsController(ctx, cfg.Mgr, cfg.InputChannels); err != nil {
		setupLog.Error(err, "unable to create controller")
		return err
	}

	return cfg.Mgr.Start(ctx)
}

// 1. krt collection of gloov1.Proxy
//    a. derived from krt collection for every resource in the current API snapshot (minus the GE specific things)
//    b. Endpoint krt collection that DOES NOT need an EDS Discovery Plugin minus channel machinery
//    c. scale benefit: per-dependent resource translation
//       i. first step could be per Upstream translation
// 2. possible option: convert ApiSnapshot to interface to allow pluggable impls
// notes: deprecate UpstreamGroups
//
// RouteOptions collection; index by targetRef
// HTTPRoute collection; index by parentRef
// Gateway collection
// Proxy/Gateway collection (fetch HTTPRoutes indexed by parentRef
// proxy <- Gateway <- HTTPRoutes <- RouteOptions <- AuthConfig
//   |
//    --> Listener (xds)

// this theoretically should work, however there's no easy way to do a watch via dynamic client that would return the typed objects
// we need. we could use a separate client than the ClientGetter, such as:
// - a typed client
// - kube client with correct scheme (and Watch() support)
//   - the challenge here is that we also want this to be backed by a cache so we're in a weird place and would need
//     to create the client to share the cache etc. from the controller-runtime client. but one of those goals in krt migration
//     is to not have to use controller runtime directly anymore?

// kubeclient.Register[*sologatewayv1.RouteOption](
// 	sologatewayv1.SchemeGroupVersion.WithResource("routeoptions"),
// 	sologatewayv1.SchemeGroupVersion.WithKind("RouteOption"),
// 	func(c kubeclient.ClientGetter, ns string, o metav1.ListOptions) (runtime.Object, error) {
// 		rtopts := sologatewayv1.RouteOptionList{}
// 		l, err := c.Dynamic().Resource(sologatewayv1.SchemeGroupVersion.WithResource("routeoptions")).Namespace(ns).List(ctx, o)
// 		err = runtime.DefaultUnstructuredConverter.FromUnstructured(l.UnstructuredContent(), &rtopts)
// 		if err != nil {
// 			return nil, err
// 		}
// 		return &rtopts, err
// 	},
// 	func(c kubeclient.ClientGetter, ns string, o metav1.ListOptions) (watch.Interface, error) {
// 		cfg.Mgr.
// 		return c.Dynamic().Resource(sologatewayv1.SchemeGroupVersion.WithResource("routeoptions")).Namespace(ns).Watch(ctx, o)
// 	},
// )
