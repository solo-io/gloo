package routeoptions

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"

	transformation1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/transformation"

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
	wildcardField                 = "*"
	PortalMetadataNamespace       = "io.solo.gloo.portal"
	PortalCustomMetadataNamespace = "io.solo.gloo.portal.custom_metadata"
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

// holds status structure for each RouteOption we have processed and attached
// this is used because a RouteOption is attached to a Route, but a Route may be
// attached to multiple Gateways/Listeners, so we need a single status object
// to contain the subresourceStatus for each Proxy it was translated too, but also
// all the errors specifically encountered
type legacyStatusCache = map[types.NamespacedName]legacyStatus

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

	// merge portal specific transformations into the destination route options post merge
	err = mergePortalTransformations(merged, outputRoute.GetOptions())
	if err != nil {
		contextutils.LoggerFrom(ctx).Errorf("error merging portal transformations: %v", err)
	}

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

			var newStatus legacyStatus
			newStatus.subresourceStatus = make(map[string]*core.Status)

			// update the cache
			p.legacyStatusCache[roKey] = newStatus
		}
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
		mayberoObj := p.routeOptionCollection.GetKey(krt.Key[*solokubev1.RouteOption](krt.Named{Namespace: roKey.Namespace, Name: roKey.Name}.ResourceName()))
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
		var newStatus legacyStatus
		newStatus.subresourceStatus = make(map[string]*core.Status)
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

func mergePortalTransformations(dst, src *gloov1.RouteOptions) error {
	if src == nil || src.GetStagedTransformations().GetEarly() == nil {
		return nil // nothing to merge
	}

	var portalDynamicMetadataTransformations []*transformation1.TransformationTemplate_DynamicMetadataValue
	// portal transformations are applied in the early stage so extract the appropriate transformations from the source route options
	if src.GetStagedTransformations().GetEarly() != nil {
		for _, transformation := range src.GetStagedTransformations().GetEarly().GetRequestTransforms() {
			if transformation.GetRequestTransformation().GetTransformationTemplate().GetDynamicMetadataValues() != nil {
				// extract portal specific transformations and add them to the destination route options
				srcDynamicMetadataTransformations := transformation.GetRequestTransformation().GetTransformationTemplate().GetDynamicMetadataValues()
				for _, srcDynamicMetadataTransformation := range srcDynamicMetadataTransformations {
					if srcDynamicMetadataTransformation.GetMetadataNamespace() == PortalMetadataNamespace || srcDynamicMetadataTransformation.GetMetadataNamespace() == PortalCustomMetadataNamespace {
						portalDynamicMetadataTransformations = append(portalDynamicMetadataTransformations, srcDynamicMetadataTransformation)
					}
				}
			}
		}
	}

	var portalTransformation *transformation1.Transformation
	if len(portalDynamicMetadataTransformations) == 0 {
		return nil // nothing to merge
	}

	portalTransformation = &transformation1.Transformation{
		TransformationType: &transformation1.Transformation_TransformationTemplate{
			TransformationTemplate: &transformation1.TransformationTemplate{
				ParseBodyBehavior:     transformation1.TransformationTemplate_DontParse,
				DynamicMetadataValues: portalDynamicMetadataTransformations,
			},
		},
	}

	// if there are no staged early request transformations in the destination route option, we can add the portal metadata transformation as-is
	if dst.GetStagedTransformations().GetEarly().GetRequestTransforms() == nil {
		if dst.GetStagedTransformations() == nil {
			dst.StagedTransformations = &transformation1.TransformationStages{}
		}
		if dst.GetStagedTransformations().GetEarly() == nil {
			dst.GetStagedTransformations().Early = &transformation1.RequestResponseTransformations{}
		}

		dst.GetStagedTransformations().GetEarly().RequestTransforms = []*transformation1.RequestMatch{{
			RequestTransformation: portalTransformation,
		}}
		return nil
	}

	// if there are early transforms, merge the portal metadata with any existing transformation templates
	for _, requestMatch := range dst.GetStagedTransformations().GetEarly().GetRequestTransforms() {
		// an error should only occur if the request match is nil
		err := setMetadataOnMatch(requestMatch, portalTransformation)
		if err != nil {
			return err
		}
	}
	return nil
}

// setMetadataOnMatch merges the portal transformation metadata into a user's transformation in-place.
// If a user redefines a portal metadata, meaning they reuse a key and namespace, the user's value is prioritized over the portal's value.
func setMetadataOnMatch(match *transformation1.RequestMatch, transform *transformation1.Transformation) error {
	// we can't modify a nil request match, this is done for defensive programming, although it should never happen
	if match == nil {
		return eris.New("request match is nil")
	}

	// add the portal transformation to the existing request transformation, if it exists
	if requestTransform := match.GetRequestTransformation(); requestTransform != nil {
		err := mergePortalTransformation(requestTransform, transform)
		if err != nil {
			return err
		}
	} else {
		// if the request transformation does not exist, set it as the portal transformation
		match.RequestTransformation = transform
	}

	return nil
}

// mergePortalTransformation merges the src transformation (portal) into a dest transformation (user) in-place, only merging the values needed for portal metadata.
// It sorts the dynamic metadata values for deterministic output.
// If the dest transformation is nil, an error is returned.
func mergePortalTransformation(dest, src *transformation1.Transformation) error {
	// we can't modify a nil transformation
	if dest == nil {
		return eris.New("cannot merge into nil transformation")
	}
	// do not merge if the src is nil
	if src == nil {
		return nil
	}

	// we can only merge onto transformation templates
	if userTemplate, ok := dest.GetTransformationType().(*transformation1.Transformation_TransformationTemplate); ok {
		// set the parse body behavior from the new transform (which should be DontParse for the portal metadata)
		userTemplate.TransformationTemplate.ParseBodyBehavior = src.GetTransformationTemplate().GetParseBodyBehavior()

		metadataMap := make(map[string]*transformation1.TransformationTemplate_DynamicMetadataValue, len(userTemplate.TransformationTemplate.GetDynamicMetadataValues()))
		for _, value := range userTemplate.TransformationTemplate.GetDynamicMetadataValues() {
			key := getDynamicMetadataKey(value)
			metadataMap[key] = value
		}

		for _, value := range src.GetTransformationTemplate().GetDynamicMetadataValues() {
			key := getDynamicMetadataKey(value)
			// if the key does not exist in the dynamic metadata map, add it so that we don't overwrite user-defined values or have duplicates
			if _, ok := metadataMap[key]; !ok {
				metadataMap[key] = value
			}
		}

		// sort existing and new metadata for deterministic output
		sortedKeys := make([]string, 0, len(metadataMap))
		for key := range metadataMap {
			sortedKeys = append(sortedKeys, key)
		}
		sort.Slice(sortedKeys, func(i, j int) bool { return sortedKeys[i] < sortedKeys[j] })

		dynamicMetadataValues := make([]*transformation1.TransformationTemplate_DynamicMetadataValue, 0, len(sortedKeys))
		for _, key := range sortedKeys {
			dynamicMetadataValues = append(dynamicMetadataValues, metadataMap[key])
		}
		userTemplate.TransformationTemplate.DynamicMetadataValues = dynamicMetadataValues
	}
	return nil
}

// getDynamicMetadataKey returns a unique key for a dynamic metadata value.
func getDynamicMetadataKey(v *transformation1.TransformationTemplate_DynamicMetadataValue) string {
	return v.GetMetadataNamespace() + "." + v.GetKey()
}
