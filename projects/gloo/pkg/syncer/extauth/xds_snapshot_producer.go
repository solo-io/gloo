package extauth

import (
	"context"
	"sort"

	"github.com/hashicorp/go-multierror"

	"github.com/golang/protobuf/proto"
	errors "github.com/rotisserie/eris"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	extauth "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/extauth/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/gloosnapshot"
	glooutils "github.com/solo-io/gloo/projects/gloo/pkg/utils"
	"github.com/solo-io/go-utils/contextutils"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-kit/pkg/api/v2/reporter"
)

var (
	_ XdsSnapshotProducer = new(proxySourcedXdsSnapshotProducer)
	_ XdsSnapshotProducer = new(snapshotSourcedXdsSnapshotProducer)
)

// XdsSnapshotProducer produces the slice of AuthConfigs which will be sent over xDS to the ExtAuth Service
// Any errors that are encountered are added to the report
type XdsSnapshotProducer interface {
	// ProduceXdsSnapshot produces a list of AuthConfigs for a given API Snapshot
	// NOT Thread-safe
	// This can be invoked by either a translator or validation request, and we want this to be thread-safe
	// eventually. In the meantime, instantiate a new instance per Sync
	ProduceXdsSnapshot(
		ctx context.Context,
		settings *v1.Settings,
		snapshot *gloosnapshot.ApiSnapshot,
		reports reporter.ResourceReports,
	) []*extauth.ExtAuthConfig
}

// proxySourcedXdsSnapshotProducer is the previous implementation of our xdsSnapshotProducer
// It works by walking the Proxy object, and using any AuthConfig references defined on that object
// as the source of truth for which AuthConfigs the Control Plane needs to translate and send to the ext-auth-service
// over xDS.
// The downside to this implementation is that:
//
//		A: Each sync is distinct, so we reprocess the entire set of configuration each time
//	 B: If the Proxy CR is deleted (sometimes done to Debug), we will identify 0 AuthConfigs and send
//			that to the ext-auth-service. The danger of this is outlined in https://github.com/solo-io/solo-projects/issues/3558
//		C: Errors that exist on AuthConfig objects are reported on the objects themselves, but not the Proxy that references
//			it. This means that an invalid AuthConfig can be referenced by a VirtualService and we will accept the VirtualService
//			instead of rejecting it.
type proxySourcedXdsSnapshotProducer struct {
	translatedConfigs map[string]*extauth.ExtAuthConfig
}

func NewProxySourcedXdsSnapshotProducer() *proxySourcedXdsSnapshotProducer {
	return &proxySourcedXdsSnapshotProducer{
		translatedConfigs: make(map[string]*extauth.ExtAuthConfig),
	}
}

func (i *proxySourcedXdsSnapshotProducer) reset() {
	i.translatedConfigs = make(map[string]*extauth.ExtAuthConfig)
}

func (i *proxySourcedXdsSnapshotProducer) ProduceXdsSnapshot(
	ctx context.Context,
	settings *v1.Settings,
	snapshot *gloosnapshot.ApiSnapshot,
	reports reporter.ResourceReports,
) []*extauth.ExtAuthConfig {
	for _, proxy := range snapshot.Proxies {
		for _, listener := range proxy.Listeners {
			virtualHosts := glooutils.GetVirtualHostsForListener(listener)

			for _, virtualHost := range virtualHosts {
				virtualHost = proto.Clone(virtualHost).(*v1.VirtualHost)
				virtualHost.Name = glooutils.SanitizeForEnvoy(ctx, virtualHost.Name, "virtual host")

				if err := i.processAuthExtension(ctx, snapshot, settings, virtualHost.GetOptions().GetExtauth(), reports, proxy); err != nil {
					// Continue to next virtualHost, error has been added to the report.
					continue
				}

				for _, route := range virtualHost.Routes {
					if err := i.processAuthExtension(ctx, snapshot, settings, route.GetOptions().GetExtauth(), reports, proxy); err != nil {
						// Continue to next route, error has been added to the report.
						continue
					}

					for _, weightedDestination := range route.GetRouteAction().GetMulti().GetDestinations() {
						if err := i.processAuthExtension(ctx, snapshot, settings, weightedDestination.GetOptions().GetExtauth(),
							reports, proxy); err != nil {
							// Continue to next weighted destination, error has been added to the report.
							continue
						}
					}
				}
			}
		}
	}

	return convertConfigMapToSortedList(i.translatedConfigs)
}

func (i *proxySourcedXdsSnapshotProducer) processAuthExtension(ctx context.Context, snap *gloosnapshot.ApiSnapshot, settings *v1.Settings, config *extauth.ExtAuthExtension,
	reports reporter.ResourceReports, parentProxy resources.InputResource) error {
	if config.GetConfigRef() != nil {
		return i.processAuthExtensionConfigRef(ctx, snap, config.GetConfigRef(), reports, parentProxy)
	}

	if config.GetCustomAuth() != nil {
		return i.processAuthExtensionCustomAuth(ctx, settings, config.GetCustomAuth(), reports, parentProxy)
	}

	// Just return if there is nothing to process
	return nil
}

func (i *proxySourcedXdsSnapshotProducer) processAuthExtensionConfigRef(ctx context.Context, snap *gloosnapshot.ApiSnapshot, configRef *core.ResourceRef,
	reports reporter.ResourceReports, parentProxy resources.InputResource) error {

	if configRef == nil {
		// Just return if there is nothing to translate
		return nil
	}

	// Don't perform duplicate work if we already have translated this resource
	if _, ok := i.translatedConfigs[configRef.Key()]; ok {
		return nil
	}

	cfg, err := snap.AuthConfigs.Find(configRef.GetNamespace(), configRef.GetName())
	if err != nil {
		contextutils.LoggerFrom(ctx).Warnf("Unable to find referenced auth config %v in snapshot", configRef)
		reports.AddError(parentProxy, err)
		return err
	}

	// do validation
	ValidateAuthConfig(cfg, reports)

	translatedConfig, err := TranslateExtAuthConfig(ctx, snap, configRef)
	if err != nil {
		reports.AddError(cfg, err)
		return err
	} else if translatedConfig == nil {
		// Do nothing if config is empty or consists only of custom auth configs
		return nil
	}

	i.translatedConfigs[configRef.Key()] = translatedConfig
	return nil
}

func (i *proxySourcedXdsSnapshotProducer) processAuthExtensionCustomAuth(ctx context.Context, settings *v1.Settings, customAuth *extauth.CustomAuth,
	reports reporter.ResourceReports, parentProxy resources.InputResource) error {

	customAuthServerName := customAuth.GetName()
	if customAuthServerName == "" {
		// If name is not specified, there is nothing to validate
		return nil
	}

	namedExtAuthSettings := settings.GetNamedExtauth()
	if namedExtAuthSettings == nil {
		// A name is specified, but no settings are configured
		err := errors.New("Unable to find named_extauth in Settings")
		contextutils.LoggerFrom(ctx).Warnf("%v", err)
		reports.AddError(parentProxy, err)
		return err
	}

	if _, ok := namedExtAuthSettings[customAuthServerName]; !ok {
		// A name is specified, but it isn't one of the settings that are configured
		err := errors.Errorf("Unable to find custom auth server [%s] in named_extauth in Settings", customAuthServerName)
		contextutils.LoggerFrom(ctx).Warnf("%v", err)
		reports.AddError(parentProxy, err)
		return err
	}

	return nil
}

// convertConfigMapToSortedList converts a map of AuthConfigs into a slice with a stable order.
// It doesn't really matter what order as long as it's consistent.
// If we don't do this, different orders of configs will produce different
// hashes, which will trick the system into unnecessarily thinking that we need to update the extauth service.
func convertConfigMapToSortedList(configMap map[string]*extauth.ExtAuthConfig) []*extauth.ExtAuthConfig {
	// extract values for sorting
	var configs []*extauth.ExtAuthConfig
	var configKeys []string
	for key, cfg := range configMap {
		configs = append(configs, cfg)
		configKeys = append(configKeys, key)
	}

	sort.SliceStable(configs, func(i, j int) bool {
		return configKeys[i] < configKeys[j]
	})

	return configs
}

// snapshotSourcedXdsSnapshotProducer is the new implementation of our xdsSnapshotProducer
// It works by walking the API Snapshot, and using any AuthConfig defined there
// as the source of truth for which AuthConfigs the Control Plane needs to translate and send to the ext-auth-service
// over xDS.
// The upside to this implementation is that:
//
//	A: We consistently send the entire set of AuthConfigs to the ext-auth-service, even if the Proxy object is deleted
//	B: We report errors from AuthConfigs as errors on the Proxy. This will make validating AuthConfigs (https://github.com/solo-io/gloo/issues/7272)
//		easier to implement.
//
// The downside to this implementation is that:
//
//	A: Each sync is distinct, so we reprocess the entire set of configuration each time. We can more easily add intelligence
//		to avoid processing if the set of AuthConfigs are the same
type snapshotSourcedXdsSnapshotProducer struct {
	authConfigCache map[string]authConfigCacheValue
}

// authConfigCacheValue represents the Errors and Warnings with an AuthConfig
type authConfigCacheValue struct {
	Warnings, Errors error
}

// NewSnapshotSourcedXdsSnapshotProducer returns an instances of a snapshotSourcedXdsSnapshotProducer
func NewSnapshotSourcedXdsSnapshotProducer() *snapshotSourcedXdsSnapshotProducer {
	return &snapshotSourcedXdsSnapshotProducer{
		authConfigCache: map[string]authConfigCacheValue{},
	}
}

// ProduceXdsSnapshot produces a list of AuthConfigs for a given API Snapshot
// NOT Thread-safe
func (e *snapshotSourcedXdsSnapshotProducer) ProduceXdsSnapshot(
	ctx context.Context,
	settings *v1.Settings,
	snapshot *gloosnapshot.ApiSnapshot,
	reports reporter.ResourceReports,
) []*extauth.ExtAuthConfig {
	// 1. Process all the AuthConfigs in the Snapshot
	xdsAuthConfigs := e.processAuthConfigs(ctx, snapshot, reports)

	// 2. Process the Proxies
	for _, proxy := range snapshot.Proxies {
		e.processProxy(settings, proxy, reports)
	}

	// 3. Return a sorted list of ALL VALID AuthConfigs in the cache
	return xdsAuthConfigs
}

func (e *snapshotSourcedXdsSnapshotProducer) processAuthConfigs(
	ctx context.Context,
	snapshot *gloosnapshot.ApiSnapshot,
	reports reporter.ResourceReports,
) []*extauth.ExtAuthConfig {
	xdsAuthConfigs := make([]*extauth.ExtAuthConfig, 0, len(snapshot.AuthConfigs))

	// Reset the cache
	// In the future, this can be a long-lived cache, but for now it is not
	e.authConfigCache = make(map[string]authConfigCacheValue, len(snapshot.AuthConfigs))

	snapshot.AuthConfigs.Each(func(authConfig *extauth.AuthConfig) {
		// Convert the user facing (external) AuthConfig into the internal API that we use to send resources over xDS
		validationErr := ErrorIfInvalidAuthConfig(authConfig)
		xdsAuthConfig, conversionErr := ConvertExternalAuthConfigToXdsAuthConfig(ctx, snapshot, authConfig)

		authConfigErr := multierror.Append(validationErr, conversionErr)
		if authConfigErr.ErrorOrNil() != nil {
			// We encountered an error during the translation of this AuthConfig
			// Write that error to a report on the AuthConfig
			reports.AddErrors(authConfig, authConfigErr.Errors...)
		} else {
			// We did not encounter an error, accumulate all valid AuthConfigs in a list
			// These AuthConfigs are unsorted at the moment, but sorted after we aggregate them all
			xdsAuthConfigs = append(xdsAuthConfigs, xdsAuthConfig)
		}

		// Update an internal cache with the validation results of each AuthConfig
		// This will be used when we traverse the proxy to determine if resources (routes, virtual services)
		// reference valid AuthConfigs
		e.authConfigCache[authConfig.GetMetadata().Ref().Key()] = authConfigCacheValue{
			Errors: authConfigErr.ErrorOrNil(),
			// We don't yet support placing warnings on AuthConfigs
			// https://github.com/solo-io/solo-projects/issues/2950
			Warnings: nil,
		}
	})

	// Sort the AuthConfig slice alphabetically
	sort.SliceStable(xdsAuthConfigs, func(i, j int) bool {
		return xdsAuthConfigs[i].GetAuthConfigRefName() < xdsAuthConfigs[j].GetAuthConfigRefName()
	})
	return xdsAuthConfigs
}

// processProxy walks the Proxy object and for any ExtAuthExtensions defined on the Proxy
// (this can be set on the VirtualHost, Route or WeightedDestination level) validates the extension.
// If the extension is invalid, or points to an invalid AuthConfig, amends the Proxy resource report to
// include the Error message.
func (e *snapshotSourcedXdsSnapshotProducer) processProxy(
	settings *v1.Settings,
	proxy *v1.Proxy,
	reports reporter.ResourceReports,
) {
	for _, listener := range proxy.GetListeners() {
		virtualHosts := glooutils.GetVirtualHostsForListener(listener)

		for _, virtualHost := range virtualHosts {
			e.processExtAuthExtension(settings, proxy, virtualHost.GetOptions().GetExtauth(), reports)

			for _, route := range virtualHost.GetRoutes() {
				e.processExtAuthExtension(settings, proxy, route.GetOptions().GetExtauth(), reports)

				for _, weightedDestination := range route.GetRouteAction().GetMulti().GetDestinations() {
					e.processExtAuthExtension(settings, proxy, weightedDestination.GetOptions().GetExtauth(), reports)
				}
			}
		}
	}
}

// processExtAuthExtension validates the ExtAuthExtension and amends the Proxy resource report
// to include an Error if one exists
func (e *snapshotSourcedXdsSnapshotProducer) processExtAuthExtension(
	settings *v1.Settings,
	proxy *v1.Proxy,
	extension *extauth.ExtAuthExtension,
	reports reporter.ResourceReports,
) {
	if extension == nil {
		// No extension defined, nothing to do
		return
	}

	if extension.GetConfigRef() != nil {
		acCacheValue, ok := e.authConfigCache[extension.GetConfigRef().Key()]
		if !ok {
			// ConfigRef points to an unknown AuthConfig
			reports.AddError(proxy, errors.Errorf("list did not find authConfig %s", extension.GetConfigRef().Key()))
		} else if acCacheValue.Errors != nil {
			// ConfigRef points to a known AuthConfig which is invalid
			// https://github.com/solo-io/gloo/issues/7272
			// If we add to the Proxy report, we will start validating VirtualServices (and other resources) which reference
			// invalid AuthConfigs. I tried to introduce this, but a number of tests will need to be updated,
			// so I backed out the change for now.
			// reports.AddError(proxy, errors.Errorf("proxy references an authConfig %s which is invalid", extension.GetConfigRef().Key()))
		} else {
			// ConfigRef points to a known AuthConfig which is valid
			// Do nothing
		}

		return
	}

	if extension.GetCustomAuth().GetName() != "" {
		customAuthServerName := extension.GetCustomAuth().GetName()
		namedExtAuthSettings := settings.GetNamedExtauth()
		if namedExtAuthSettings == nil {
			// A name is specified, but no settings are configured
			reports.AddError(proxy, errors.New("Unable to find namedExtauth in Settings"))
			return
		}

		if _, ok := namedExtAuthSettings[customAuthServerName]; !ok {
			// A name is specified, but it isn't one of the settings that are configured
			reports.AddError(proxy, errors.Errorf("Unable to find custom auth server [%s] in namedExtauth in Settings", customAuthServerName))
			return
		}
	}
}
