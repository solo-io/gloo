package routeoptions

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/rotisserie/eris"
	"github.com/solo-io/gloo/projects/gloo/pkg/xds"
	"github.com/solo-io/go-utils/contextutils"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-kit/pkg/api/v2/reporter"
	"go.uber.org/zap"
	"istio.io/istio/pkg/kube/krt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/sets"
	"sigs.k8s.io/controller-runtime/pkg/client"
	gwv1 "sigs.k8s.io/gateway-api/apis/v1"

	sologatewayv1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	solokubev1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1/kube/apis/gateway.solo.io/v1"
	gwquery "github.com/solo-io/gloo/projects/gateway2/query"
	"github.com/solo-io/gloo/projects/gateway2/reports"
	"github.com/solo-io/gloo/projects/gateway2/translator/plugins"
	rtoptquery "github.com/solo-io/gloo/projects/gateway2/translator/plugins/routeoptions/query"
	"github.com/solo-io/gloo/projects/gateway2/translator/plugins/utils"
	"github.com/solo-io/gloo/projects/gateway2/translator/routeutils"
	"github.com/solo-io/gloo/projects/gateway2/wellknown"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/grpc/validation"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	glooutils "github.com/solo-io/gloo/projects/gloo/pkg/utils"
)

const (
	// wildcardField is used to enable overriding all fields in RouteOptions inherited from the parent route.
	wildcardField = "*"
)

var (
	_ plugins.RoutePlugin  = &plugin{}
	_ plugins.StatusPlugin = &plugin{}

	ReadingRouteOptionErrStr = "error reading RouteOption"
)

// holds the data structures needed to derive and report a classic GE status
type legacyStatus struct {
	// maps proxyName -> proxyStatus
	subresourceStatus map[string]*core.Status
	// *All* of the route errors encountered during processing for gloov1.Routes which receive their
	// options for this RouteOption
	routeErrors []*validation.RouteReport_Error
}

func newLegacyStatus() *legacyStatus {
	return &legacyStatus{
		subresourceStatus: map[string]*core.Status{},
		routeErrors:       []*validation.RouteReport_Error{},
	}
}

// holds status structure for each RouteOption we have processed and attached
// this is used because a RouteOption is attached to a Route, but a Route may be
// attached to multiple Gateways/Listeners, so we need a single status object
// to contain the subresourceStatus for each Proxy it was translated too, but also
// all the errors specifically encountered
type legacyStatusCache = map[types.NamespacedName]*legacyStatus

type plugin struct {
	gwQueries             gwquery.GatewayQueries
	rtOptQueries          rtoptquery.RouteOptionQueries
	legacyStatusCache     legacyStatusCache
	routeOptionCollection krt.Collection[*solokubev1.RouteOption]
	statusReporter        reporter.StatusReporter
}

func NewPlugin(
	gwQueries gwquery.GatewayQueries,
	client client.Client,
	routeOptionCollection krt.Collection[*solokubev1.RouteOption],
	statusReporter reporter.StatusReporter,
) *plugin {
	legacyStatusCache := make(legacyStatusCache)
	return &plugin{
		gwQueries:             gwQueries,
		rtOptQueries:          rtoptquery.NewQuery(client),
		legacyStatusCache:     legacyStatusCache,
		routeOptionCollection: routeOptionCollection,
		statusReporter:        statusReporter,
	}
}

func (p *plugin) ApplyRoutePlugin(
	ctx context.Context,
	routeCtx *plugins.RouteContext,
	outputRoute *gloov1.Route,
) error {
	// check for RouteOptions applied to the given routeCtx
	routeOptions, _, sources, err := p.handleAttachment(ctx, routeCtx)
	if err != nil {
		return err
	}
	if routeOptions == nil {
		return nil
	}

	merged, OptionsMergeResult := mergeOptionsForRoute(ctx, routeCtx.HTTPRoute, routeOptions, outputRoute.GetOptions())
	if OptionsMergeResult == glooutils.OptionsMergedNone {
		// No existing options merged into 'sources', so set the 'sources' on the outputRoute
		routeutils.SetRouteSources(outputRoute, sources)
	} else if OptionsMergeResult == glooutils.OptionsMergedPartial {
		// Some existing options merged into 'sources', so append the 'sources' on the outputRoute
		routeutils.AppendRouteSources(outputRoute, sources)
	} // In case OptionsMergedFull, the correct sources are already set on the outputRoute

	// Set the merged RouteOptions on the outputRoute
	outputRoute.Options = merged

	// Track that we used this RouteOption is our status cache
	// we do this so we can persist status later for all attached RouteOptions
	p.trackAcceptedRouteOptions(outputRoute.GetMetadataStatic().GetSources())

	return nil
}

func mergeOptionsForRoute(
	ctx context.Context,
	route *gwv1.HTTPRoute,
	dst, src *gloov1.RouteOptions,
) (*gloov1.RouteOptions, glooutils.OptionsMergeResult) {
	// By default, lower priority options cannot override higher priority ones
	// and can only augment them during a merge such that fields unset in the higher
	// priority options can be merged in from the lower priority options.
	// In the case of delegated routes, a parent route can enable child routes to override
	// all (wildcard *) or specific fields using the wellknown.PolicyOverrideAnnotation.
	fieldsAllowedToOverride := sets.New[string]()

	// If the route already has options set, we should override/augment them.
	// This is important because for delegated routes, the plugin will
	// be invoked on the child routes multiple times for each parent route
	// that may override/augment them.
	//
	// By default, parent options (routeOptions) are preferred, unless the parent explicitly
	// enabled child routes (outputRoute.Options) to override parent options.
	fieldsStr, delegatedPolicyOverride := route.Annotations[wellknown.PolicyOverrideAnnotation]
	if delegatedPolicyOverride {
		delegatedFieldsToOverride := parseDelegationFieldOverrides(fieldsStr)
		if delegatedFieldsToOverride.Len() == 0 {
			// Invalid annotation value, so log an error but enforce the default behavior of preferring the parent options.
			contextutils.LoggerFrom(ctx).Errorf("invalid value %q for annotation %s on route %s; must be %s or a comma-separated list of field names",
				fieldsStr, wellknown.PolicyOverrideAnnotation, client.ObjectKeyFromObject(route), wildcardField)
		} else {
			fieldsAllowedToOverride = delegatedFieldsToOverride
		}
	}

	return glooutils.MergeRouteOptionsWithOverrides(dst, src, fieldsAllowedToOverride)
}

func (p *plugin) InitStatusPlugin(ctx context.Context, statusCtx *plugins.StatusContext) error {
	for _, proxyWithReport := range statusCtx.ProxiesWithReports {
		// now that we translate proxies one by one, we can't assume ApplyRoutePlugin is called before ApplyStatusPlugin for all proxies
		// ApplyStatusPlugin should be come idempotent, as also now it gets applied outside of translation context.
		// we need to track ownership separately. TODO: re-think this on monday

		// for this specific proxy, get all the route errors and their associated RouteOption sources
		routeErrors := extractRouteErrors(proxyWithReport.Reports.ProxyReport)
		for roKey := range routeErrors {
			p.legacyStatusCache[roKey] = newLegacyStatus()
		}
	}
	return nil
}

// MergeStatusPlugin merges the status of the source plugin in to the current plugin.
// This is used late in the proxy sync process to report the status of the RO when
// it references multiple proxies/snapshots.
func (p *plugin) MergeStatusPlugin(ctx context.Context, source any) error {
	sourceStatusPlugin, ok := source.(*plugin)
	if !ok {
		return nil
	}

	for cacheKey, status := range sourceStatusPlugin.legacyStatusCache {
		destStatus, ok := p.legacyStatusCache[cacheKey]
		if !ok {
			destStatus = newLegacyStatus()
		}

		destStatus.routeErrors = append(destStatus.routeErrors, status.routeErrors...)
		for subresourceKey, subresourceStatus := range status.subresourceStatus {
			destStatus.subresourceStatus[subresourceKey] = subresourceStatus
		}

		p.legacyStatusCache[cacheKey] = destStatus
	}

	return nil
}

func (p *plugin) ApplyStatusPlugin(ctx context.Context, statusCtx *plugins.StatusContext) error {
	logger := contextutils.LoggerFrom(ctx).Desugar()
	// gather all RouteOptions we need to report status for
	for _, proxyWithReport := range statusCtx.ProxiesWithReports {
		// now that we translate proxies one by one, we can't assume ApplyRoutePlugin is called before ApplyStatusPlugin for all proxies
		// ApplyStatusPlugin should be come idempotent, as also now it gets applied outside of translation context.
		// we need to track ownership separately. TODO: re-think this on monday

		// get proxy status to use for RouteOption status
		proxyStatus := p.statusReporter.StatusFromReport(proxyWithReport.Reports.ResourceReports[proxyWithReport.Proxy], nil)

		// for this specific proxy, get all the route errors and their associated RouteOption sources
		routeErrors := extractRouteErrors(proxyWithReport.Reports.ProxyReport)
		for roKey, rerrs := range routeErrors {
			// grab the existing status object for this RouteOption
			statusForRO, ok := p.legacyStatusCache[roKey]
			if !ok {
				// we are processing an error that has a RouteOption source that we hadn't encountered until now
				// this shouldn't happen
				logger.DPanic("while trying to apply status for RouteOptions, we found a Route error sourced by an unknown RouteOption", zap.Stringer("RouteOption", roKey))
			}

			// set the subresource status for this specific proxy on the RO
			thisSubresourceStatus := statusForRO.subresourceStatus
			thisSubresourceStatus[xds.SnapshotCacheKey(proxyWithReport.Proxy)] = proxyStatus
			statusForRO.subresourceStatus = thisSubresourceStatus

			// add any routeErrors from this Proxy translation
			statusForRO.routeErrors = append(statusForRO.routeErrors, rerrs...)

			// update the cache
			p.legacyStatusCache[roKey] = statusForRO
		}
	}
	routeOptionReport := make(reporter.ResourceReports)
	var multierr *multierror.Error
	for roKey, status := range p.legacyStatusCache {
		// get the obj by namespacedName
		mayberoObj := p.routeOptionCollection.GetKey(krt.Named{Namespace: roKey.Namespace, Name: roKey.Name}.ResourceName())
		if mayberoObj == nil {
			err := errors.New("RouteOption not found")
			multierr = multierror.Append(multierr, eris.Wrapf(err, "%s %s in namespace %s", ReadingRouteOptionErrStr, roKey.Name, roKey.Namespace))
			continue
		}
		roObj := **mayberoObj
		roObj.Spec.Metadata = &core.Metadata{}
		roObj.Spec.GetMetadata().Name = roObj.GetName()
		roObj.Spec.GetMetadata().Namespace = roObj.GetNamespace()
		roObjSk := &roObj.Spec
		roObjSk.NamespacedStatuses = &roObj.Status

		// mark this object to be processed
		routeOptionReport.Accept(roObjSk)

		// add any route errors for this obj
		for i, rerr := range status.routeErrors {
			rErr := errors.New(rerr.GetReason())
			logger.Debug("adding error to RouteOption status", zap.Stringer("RouteOption", roKey), zap.Error(rErr), zap.Int("routeErrorIndex", i))
			routeOptionReport.AddError(roObjSk, rErr)
		}

		// actually write out the reports!
		err := p.statusReporter.WriteReports(ctx, routeOptionReport, status.subresourceStatus)
		if err != nil {
			multierr = multierror.Append(multierr, fmt.Errorf("error writing status report from RouteOptionPlugin: %w", err))
			continue
		}
	}
	return multierr.ErrorOrNil()
}

// tracks the attachment of a RouteOption so we know which RouteOptions to report status for
func (p *plugin) trackAcceptedRouteOptions(
	sources []*gloov1.SourceMetadata_SourceRef,
) {
	for _, source := range sources {
		newStatus := newLegacyStatus()
		p.legacyStatusCache[client.ObjectKey{
			Namespace: source.GetResourceRef().GetNamespace(),
			Name:      source.GetResourceRef().GetName(),
		}] = newStatus
	}
}

func (p *plugin) handleAttachment(
	ctx context.Context,
	routeCtx *plugins.RouteContext,
) (*gloov1.RouteOptions, *solokubev1.RouteOption, []*gloov1.SourceMetadata_SourceRef, error) {
	// TODO: This is far too naive and we should optimize the amount of querying we do.
	// Route plugins run on every match for every Rule in a Route but the attached options are
	// the same each time; i.e. HTTPRoute <-1:1-> RouteOptions.
	// We should only make this query once per HTTPRoute.
	attachedOption, sources, err := p.rtOptQueries.GetRouteOptionForRouteRule(
		ctx,
		types.NamespacedName{Name: routeCtx.HTTPRoute.Name, Namespace: routeCtx.HTTPRoute.Namespace},
		routeCtx.Rule,
		p.gwQueries,
	)
	if err != nil {
		contextutils.LoggerFrom(ctx).Errorf("error getting RouteOptions for Route: %v", err)
		switch {
		case errors.Is(err, utils.ErrTypesNotEqual):
		default:
			routeCtx.Reporter.SetCondition(reports.RouteCondition{
				Type:    gwv1.RouteConditionResolvedRefs,
				Status:  metav1.ConditionFalse,
				Reason:  gwv1.RouteReasonBackendNotFound,
				Message: err.Error(),
			})
		}
		return nil, nil, nil, err
	}
	if attachedOption == nil || attachedOption.Spec.GetOptions() == nil {
		return nil, nil, nil, nil
	}

	return attachedOption.Spec.GetOptions(), attachedOption, sources, nil
}

// given a ProxyReport, extract and aggregate all Route errors that have RouteOption source metadata
// and key them by the source RouteOption NamespacedName
func extractRouteErrors(proxyReport *validation.ProxyReport) map[types.NamespacedName][]*validation.RouteReport_Error {
	routeErrors := make(map[types.NamespacedName][]*validation.RouteReport_Error)
	routeReports := getAllRouteReports(proxyReport.GetListenerReports())
	for _, rr := range routeReports {
		for _, rerr := range rr.GetErrors() {
			// if we've found a RouteReport with an Error, let's check if it has a sourced RouteOption
			// if so, we will add that error to the list of errors associated to that RouteOption
			if roKey, ok := extractRouteOptionSourceKeys(rerr); ok {
				errors := routeErrors[roKey]
				errors = append(errors, rerr)
				routeErrors[roKey] = errors
			}
		}
	}
	return routeErrors
}

// given a list of ListenerReports, iterate all HttpListeners to find and return all RouteReports
func getAllRouteReports(listenerReports []*validation.ListenerReport) []*validation.RouteReport {
	routeReports := []*validation.RouteReport{}
	for _, lr := range listenerReports {
		for _, hlr := range lr.GetAggregateListenerReport().GetHttpListenerReports() {
			for _, vhr := range hlr.GetVirtualHostReports() {
				routeReports = append(routeReports, vhr.GetRouteReports()...)
			}
		}
	}
	return routeReports
}

// if the Route error has a RouteOption source associated with it, extract the source and return it
func extractRouteOptionSourceKeys(routeErr *validation.RouteReport_Error) (types.NamespacedName, bool) {
	metadata := routeErr.GetMetadata()
	if metadata == nil {
		return types.NamespacedName{}, false
	}

	for _, src := range metadata.GetSources() {
		if src.GetResourceKind() == sologatewayv1.RouteOptionGVK.Kind {
			key := types.NamespacedName{
				Namespace: src.GetResourceRef().GetNamespace(),
				Name:      src.GetResourceRef().GetName(),
			}
			return key, true
		}
	}

	return types.NamespacedName{}, false
}

func parseDelegationFieldOverrides(val string) sets.Set[string] {
	if val == wildcardField {
		return sets.New(wildcardField)
	}

	set := sets.New[string]()
	parts := strings.Split(val, ",")
	for _, part := range parts {
		set.Insert(strings.ToLower(strings.TrimSpace(part)))
	}
	return set
}
