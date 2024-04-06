package routeoptions

import (
	"context"
	"errors"
	"fmt"

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
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/go-utils/contextutils"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	gwv1 "sigs.k8s.io/gateway-api/apis/v1"
	gwv1alpha2 "sigs.k8s.io/gateway-api/apis/v1alpha2"
)

var _ plugins.RoutePlugin = &plugin{}
var _ plugins.StatusPlugin = &plugin{}

var routeOptionGK = schema.GroupKind{
	Group: sologatewayv1.RouteOptionGVK.Group,
	Kind:  sologatewayv1.RouteOptionGVK.Kind,
}

type PolicyReport struct {
	Ancestors          map[types.NamespacedName]*PolicyAncestorReport
	ObservedGeneration int64
}

type PolicyAncestorReport struct {
	Condition metav1.Condition
}

// RouteOption resource ->  PolicyReport
// PolictReport: AncestorRef(HTTPRoute) -> Actual condition
// NOTE: This assumes the only ancestor type of a RouteOption is HTTPRoute; if this changes, we need
// to track Group & Kind along with the types.NN
type policyStatusCache = map[types.NamespacedName]*PolicyReport

type plugin struct {
	gwQueries    gwquery.GatewayQueries
	rtOptQueries rtoptquery.RouteOptionQueries
	statusCache  policyStatusCache
}

func NewPlugin(gwQueries gwquery.GatewayQueries, client client.Client) *plugin {
	policyStatusCache := make(policyStatusCache)
	return &plugin{
		gwQueries,
		rtoptquery.NewQuery(client),
		policyStatusCache,
	}
}

func (p *plugin) getPolicyReport(key types.NamespacedName) *PolicyReport {
	return p.statusCache[key]
}

func (p *plugin) getOrCreatePolicyReport(routeOption *solokubev1.RouteOption) *PolicyReport {
	pr := p.getPolicyReport(client.ObjectKeyFromObject(routeOption))
	if pr == nil {
		pr = &PolicyReport{}
		pr.ObservedGeneration = routeOption.GetGeneration()
		pr.Ancestors = make(map[types.NamespacedName]*PolicyAncestorReport)
		p.statusCache[client.ObjectKeyFromObject(routeOption)] = pr
	}
	return pr
}

func (pr *PolicyReport) getAncestorReport(ancestorKey types.NamespacedName) *PolicyAncestorReport {
	return pr.Ancestors[ancestorKey]
}

func (pr *PolicyReport) upsertAncestorCondition(ancestorKey types.NamespacedName, condition metav1.Condition) {
	ar := pr.getAncestorReport(ancestorKey)
	if ar == nil {
		ar = &PolicyAncestorReport{}
	}
	ar.Condition = condition
	pr.Ancestors[ancestorKey] = ar
}

func (p *plugin) ApplyStatusPlugin(ctx context.Context, statusCtx plugins.StatusContext) {
	routeErrors := extractRouteErrors(statusCtx.ProxyReport)
	// we can coalesce route errors here to be keyed by the HTTPRoute/RouteOption object
	// as each HTTPRoute can only have a single attached RouteOption, we know that all gloov1.Routes
	// from a given HTTPRoute *should* have the same routeError

	for _, rerr := range routeErrors {
		route, ro := extractSourceKeys(rerr.GetMetadata())
		pr := p.getPolicyReport(ro)
		if pr == nil {
			// TODO: we got a route error for a routeoption that we weren't tracking during the plugin run; what happened?
		}

		ar := pr.getAncestorReport(route)
		if ar == nil {
			// TODO: this route error was sourced from an HTTPRoute that wasn't originally tracked; what happened?
		}
		pr.upsertAncestorCondition(route, metav1.Condition{
			Type:    string(gwv1alpha2.PolicyConditionAccepted),
			Status:  metav1.ConditionFalse,
			Reason:  string(gwv1alpha2.PolicyReasonInvalid),
			Message: rerr.GetReason(),
		})
	}
}

func (p *plugin) setPolicyStatusAccepted(
	routeOption *solokubev1.RouteOption,
	route *gwv1.HTTPRoute,
) {
	pr := p.getOrCreatePolicyReport(routeOption)
	pr.upsertAncestorCondition(client.ObjectKeyFromObject(route), metav1.Condition{
		Type:    string(gwv1alpha2.PolicyConditionAccepted),
		Status:  metav1.ConditionTrue,
		Reason:  string(gwv1alpha2.PolicyReasonAccepted),
		Message: "Attached successfully",
	})
}

// RouteError metadata should have a single HTTPRoute & possibly RouteOption resource
// TODO: add error handling, this always assumes happy path
func extractSourceKeys(
	metadata *v1.SourceMetadata,
) (route, routeOption types.NamespacedName) {
	var routeRef, routeOptionRef *v1.SourceMetadata_SourceRef
	for _, src := range metadata.GetSources() {
		if src.GetResourceKind() == wellknown.HTTPRouteKind {
			routeRef = src
		} else if src.GetResourceKind() == sologatewayv1.RouteOptionGVK.Kind {
			routeOptionRef = src
		}
	}
	route = types.NamespacedName{
		Namespace: routeRef.ResourceRef.GetNamespace(),
		Name:      routeRef.ResourceRef.GetName(),
	}
	routeOption = types.NamespacedName{
		Namespace: routeOptionRef.ResourceRef.GetNamespace(),
		Name:      routeOptionRef.ResourceRef.GetName(),
	}
	return route, routeOption
}

func extractRouteErrors(proxyReport *validation.ProxyReport) []*validation.RouteReport_Error {
	routeErrors := []*validation.RouteReport_Error{}
	for _, lr := range proxyReport.GetListenerReports() {
		for _, hlr := range lr.GetAggregateListenerReport().GetHttpListenerReports() {
			for _, vhr := range hlr.GetVirtualHostReports() {
				for _, rr := range vhr.GetRouteReports() {
					for _, rerr := range rr.GetErrors() {
						routeErrors = append(routeErrors, rerr)
					}
				}
			}
		}
	}
	return routeErrors
}

func (p *plugin) ApplyRoutePlugin(
	ctx context.Context,
	routeCtx *plugins.RouteContext,
	outputRoute *v1.Route,
) error {
	var routeOptions *v1.RouteOptions
	// check for RouteOptions applied to full Route
	routeOptions, routeOptionKube := p.handleAttachment(ctx, routeCtx)

	// allow for ExtensionRef filter override
	filterRouteOptions, err := p.handleFilter(ctx, routeCtx)
	if err != nil {
		return err
	}
	if filterRouteOptions != nil {
		routeOptions = filterRouteOptions
	}

	if routeOptions != nil {
		// clobber the existing RouteOptions; merge semantics may be desired later
		outputRoute.Options = routeOptions

		if filterRouteOptions == nil {
			// if we didn't use a filter to derive the RouteOptions for this v1.Route, let's track the
			// RouteOption policy object that was used so we can report status on it
			routeutils.AppendSourceToRoute(outputRoute, routeOptionKube)

			// report success on attaching this policy, although it may be overturned
			// later if the Proxy translation results in a processing error for these RouteOptions
			// we will track the ancestor as the targeted HTTPRoute; we may want change to Gateway later
			p.setPolicyStatusAccepted(routeOptionKube, routeCtx.Route)
		}
	}
	return nil
}

func (p *plugin) handleAttachment(
	ctx context.Context,
	routeCtx *plugins.RouteContext,
) (*v1.RouteOptions, *solokubev1.RouteOption) {
	// TODO: This is far too naive and we should optimize the amount of querying we do.
	// Route plugins run on every match for every Rule in a Route but the attached options are
	// the same each time; i.e. HTTPRoute <-1:1-> RouteOptions.
	// We should only make this query once per HTTPRoute.
	attachedOptions := getAttachedRouteOptions(ctx, routeCtx.Route, p.rtOptQueries)
	if len(attachedOptions) == 0 {
		return nil, nil
	}

	// sort attached options and apply only the earliest
	utils.SortByCreationTime(attachedOptions)
	earliestOption := attachedOptions[0]

	if earliestOption.Spec.GetOptions() != nil {
		return earliestOption.Spec.GetOptions(), earliestOption
	}
	return nil, nil
}

func (p *plugin) handleFilter(
	ctx context.Context,
	routeCtx *plugins.RouteContext,
) (*v1.RouteOptions, error) {
	filter := utils.FindExtensionRefFilter(routeCtx, routeOptionGK)
	if filter == nil {
		return nil, nil
	}

	routeOption := &solokubev1.RouteOption{}
	err := utils.GetExtensionRefObj(context.Background(), routeCtx, p.gwQueries, filter.ExtensionRef, routeOption)
	if err != nil {
		switch {
		case apierrors.IsNotFound(err):
			notFoundMsg := formatNotFoundMessage(routeCtx, filter)
			routeCtx.ParentRefReporter.SetCondition(reports.HTTPRouteCondition{
				Type:    gwv1.RouteConditionResolvedRefs,
				Status:  metav1.ConditionFalse,
				Reason:  gwv1.RouteReasonBackendNotFound,
				Message: notFoundMsg,
			})
			return nil, errors.New(notFoundMsg)
		case errors.Is(err, utils.ErrTypesNotEqual):
		case errors.Is(err, utils.ErrNotSettable):
			devErr := fmt.Errorf("developer error while getting RouteOptions as ExtensionRef: %w", err)
			contextutils.LoggerFrom(ctx).DPanic(devErr)
			return nil, devErr
		default:
			return nil, fmt.Errorf("error while getting RouteOptions as ExtensionRef route filter: %w", err)
		}
	}

	if routeOption.Spec.GetOptions() != nil {
		return routeOption.Spec.GetOptions(), nil
	}
	return nil, nil
}

func getAttachedRouteOptions(ctx context.Context, route *gwv1.HTTPRoute, queries rtoptquery.RouteOptionQueries) []*solokubev1.RouteOption {
	var routeOptionList solokubev1.RouteOptionList
	err := queries.GetRouteOptionsForRoute(ctx, route, &routeOptionList)
	if err != nil {
		contextutils.LoggerFrom(ctx).Errorf("error while Listing RouteOptions: %v", err)
		// TODO: add status to policy on error
		return nil
	}

	// as the RouteOptionList does not contain pointers, and RouteOption is a concrete proto message,
	// we need to turn it into a pointer slice to avoid copying proto message state around, copying locks, etc.
	// while we perform operations on the RouteOptionList
	ptrSlice := []*solokubev1.RouteOption{}
	items := routeOptionList.Items
	for i := range items {
		ptrSlice = append(ptrSlice, &items[i])
	}
	return ptrSlice
}

func formatNotFoundMessage(routeCtx *plugins.RouteContext, filter *gwv1.HTTPRouteFilter) string {
	return fmt.Sprintf(
		"extensionRef '%s' of type %s.%s in namespace '%s' not found",
		filter.ExtensionRef.Group,
		filter.ExtensionRef.Kind,
		filter.ExtensionRef.Name,
		routeCtx.Route.GetNamespace(),
	)
}
