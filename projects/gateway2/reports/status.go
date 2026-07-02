package reports

import (
	"context"
	"fmt"
	"reflect"
	"slices"

	"github.com/solo-io/go-utils/contextutils"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	gwv1 "sigs.k8s.io/gateway-api/apis/v1"
	gwv1a2 "sigs.k8s.io/gateway-api/apis/v1alpha2"
)

// TODO: refactor this struct + methods to better reflect the usage now in proxy_syncer

func (r *ReportMap) BuildGWStatus(ctx context.Context, gw gwv1.Gateway) *gwv1.GatewayStatus {
	gwReport := r.Gateway(&gw)
	if gwReport == nil {
		return nil
	}

	finalListeners := make([]gwv1.ListenerStatus, 0, len(gw.Spec.Listeners))
	for _, lis := range gw.Spec.Listeners {
		listenerStatus := listenerStatusWithDefaults(listenerReport(gwReport.listeners, lis.Name), lis.Name)

		finalConditions := make([]metav1.Condition, 0, len(listenerStatus.Conditions))
		oldLisStatusIndex := slices.IndexFunc(gw.Status.Listeners, func(l gwv1.ListenerStatus) bool {
			return l.Name == lis.Name
		})
		for _, lisCondition := range listenerStatus.Conditions {
			lisCondition.ObservedGeneration = gwReport.observedGeneration

			// copy old condition from gw so LastTransitionTime is set correctly below by SetStatusCondition()
			if oldLisStatusIndex != -1 {
				if cond := meta.FindStatusCondition(gw.Status.Listeners[oldLisStatusIndex].Conditions, lisCondition.Type); cond != nil {
					meta.SetStatusCondition(&finalConditions, *cond)
				}
			}
			meta.SetStatusCondition(&finalConditions, lisCondition)
		}
		listenerStatus.Conditions = finalConditions
		finalListeners = append(finalListeners, listenerStatus)
	}

	finalConditions := make([]metav1.Condition, 0)
	for _, gwCondition := range gatewayConditionsWithDefaults(gwReport.GetConditions()) {
		gwCondition.ObservedGeneration = gwReport.observedGeneration

		// copy old condition from gw so LastTransitionTime is set correctly below by SetStatusCondition()
		if cond := meta.FindStatusCondition(gw.Status.Conditions, gwCondition.Type); cond != nil {
			meta.SetStatusCondition(&finalConditions, *cond)
		}
		meta.SetStatusCondition(&finalConditions, gwCondition)
	}
	// If there are conditions on the Gateway that are not owned by our reporter, include
	// them in the final list of conditions to preseve conditions we do not own
	for _, condition := range gw.Status.Conditions {
		if meta.FindStatusCondition(finalConditions, condition.Type) == nil {
			meta.SetStatusCondition(&finalConditions, condition)
		}
	}

	finalGwStatus := gwv1.GatewayStatus{}
	finalGwStatus.Conditions = finalConditions
	finalGwStatus.Listeners = finalListeners
	return &finalGwStatus
}

// BuildRouteStatus returns a newly constructed and fully defined RouteStatus for the supplied route object
// according to the state of the ReportMap. If the ReportMap does not have a RouteReport for the given route,
// e.g. because it did not encounter the route during translation, or the object is an unsupported route kind,
// nil is returned. Supported object types are:
//
// * HTTPRoute
// * TCPRoute
func (r *ReportMap) BuildRouteStatus(ctx context.Context, obj client.Object, cName string) *gwv1.RouteStatus {
	routeReport := r.route(obj)
	if routeReport == nil {
		contextutils.LoggerFrom(ctx).Infof("missing route report for %T %s/%s", obj, obj.GetName(), obj.GetNamespace())
		return nil
	}

	contextutils.LoggerFrom(ctx).Debugf("building status for %s %s/%s",
		obj.GetObjectKind().GroupVersionKind().Kind, obj.GetNamespace(),
		obj.GetName())

	// Default to using spec.ParentRefs when building the parent statuses for a route.
	// However, for delegatee (child) routes, the parentRefs field is optional and such routes
	// may not specify it. In this case, we infer the parentRefs form the RouteReport
	// corresponding to the delegatee (child) route as the route's report is associated to a parentRef.
	var existingStatus gwv1.RouteStatus
	var parentRefs []gwv1.ParentReference
	switch route := obj.(type) {
	case *gwv1.HTTPRoute:
		existingStatus = route.Status.RouteStatus
		parentRefs = append(parentRefs, route.Spec.ParentRefs...)
		if len(parentRefs) == 0 {
			parentRefs = append(parentRefs, routeReport.parentRefs()...)
		}
	case *gwv1a2.TCPRoute:
		existingStatus = route.Status.RouteStatus
		parentRefs = append(parentRefs, route.Spec.ParentRefs...)
		if len(parentRefs) == 0 {
			parentRefs = append(parentRefs, routeReport.parentRefs()...)
		}
	default:
		contextutils.LoggerFrom(ctx).Error(fmt.Errorf("unsupported route type %T", obj), "failed to build route status")
		return nil
	}

	// Process the parent references to build the RouteParentStatus
	routeStatus := gwv1.RouteStatus{}
	for _, parentRef := range parentRefs {
		parentConditions := parentRefConditionsWithDefaults(parentRefReport(routeReport, &parentRef))

		// Get the status of the current parentRef conditions if they exist
		var currentParentRefConditions []metav1.Condition
		currentParentRefIdx := slices.IndexFunc(existingStatus.Parents, func(s gwv1.RouteParentStatus) bool {
			return reflect.DeepEqual(s.ParentRef, parentRef)
		})
		if currentParentRefIdx != -1 {
			currentParentRefConditions = existingStatus.Parents[currentParentRefIdx].Conditions
		}

		finalConditions := make([]metav1.Condition, 0, len(parentConditions))
		for _, pCondition := range parentConditions {
			pCondition.ObservedGeneration = routeReport.observedGeneration

			// Copy old condition to preserve LastTransitionTime, if it exists
			if cond := meta.FindStatusCondition(currentParentRefConditions, pCondition.Type); cond != nil {
				meta.SetStatusCondition(&finalConditions, *cond)
			}
			meta.SetStatusCondition(&finalConditions, pCondition)
		}
		// If there are conditions on the HTTPRoute that are not owned by our reporter, include
		// them in the final list of conditions to preseve conditions we do not own
		for _, condition := range currentParentRefConditions {
			if meta.FindStatusCondition(finalConditions, condition.Type) == nil {
				meta.SetStatusCondition(&finalConditions, condition)
			}
		}

		routeParentStatus := gwv1.RouteParentStatus{
			ParentRef:      parentRef,
			ControllerName: gwv1.GatewayController(cName),
			Conditions:     finalConditions,
		}
		routeStatus.Parents = append(routeStatus.Parents, routeParentStatus)
	}

	return &routeStatus
}

// Reports will initially only contain negative conditions found during translation,
// so all missing conditions are assumed to be positive. The helpers below add
// those defaults to clones so status rendering does not mutate the ReportMap
// used by krt equality.
func gatewayConditionsWithDefaults(reportConditions []metav1.Condition) []metav1.Condition {
	conditions := slices.Clone(reportConditions)
	if cond := meta.FindStatusCondition(conditions, string(gwv1.GatewayConditionAccepted)); cond == nil {
		meta.SetStatusCondition(&conditions, metav1.Condition{
			Type:   string(gwv1.GatewayConditionAccepted),
			Status: metav1.ConditionTrue,
			Reason: string(gwv1.GatewayReasonAccepted),
		})
	}
	if cond := meta.FindStatusCondition(conditions, string(gwv1.GatewayConditionProgrammed)); cond == nil {
		meta.SetStatusCondition(&conditions, metav1.Condition{
			Type:   string(gwv1.GatewayConditionProgrammed),
			Status: metav1.ConditionTrue,
			Reason: string(gwv1.GatewayReasonProgrammed),
		})
	}
	return conditions
}

func listenerStatusWithDefaults(lisReport *ListenerReport, name gwv1.SectionName) gwv1.ListenerStatus {
	status := gwv1.ListenerStatus{Name: name}
	if lisReport != nil {
		status = lisReport.Status
		// The rendered status name is the spec listener name the caller asked for,
		// independent of how the report happens to be keyed.
		status.Name = name
		status.Conditions = slices.Clone(lisReport.Status.Conditions)
		status.SupportedKinds = slices.Clone(lisReport.Status.SupportedKinds)
	}
	// set healthy conditions for Condition Types not set yet (i.e. no negative status yet, we can assume positive)
	if cond := meta.FindStatusCondition(status.Conditions, string(gwv1.ListenerConditionAccepted)); cond == nil {
		meta.SetStatusCondition(&status.Conditions, metav1.Condition{
			Type:   string(gwv1.ListenerConditionAccepted),
			Status: metav1.ConditionTrue,
			Reason: string(gwv1.ListenerReasonAccepted),
		})
	}
	if cond := meta.FindStatusCondition(status.Conditions, string(gwv1.ListenerConditionConflicted)); cond == nil {
		meta.SetStatusCondition(&status.Conditions, metav1.Condition{
			Type:   string(gwv1.ListenerConditionConflicted),
			Status: metav1.ConditionFalse,
			Reason: string(gwv1.ListenerReasonNoConflicts),
		})
	}
	if cond := meta.FindStatusCondition(status.Conditions, string(gwv1.ListenerConditionResolvedRefs)); cond == nil {
		meta.SetStatusCondition(&status.Conditions, metav1.Condition{
			Type:   string(gwv1.ListenerConditionResolvedRefs),
			Status: metav1.ConditionTrue,
			Reason: string(gwv1.ListenerReasonResolvedRefs),
		})
	}
	if cond := meta.FindStatusCondition(status.Conditions, string(gwv1.ListenerConditionProgrammed)); cond == nil {
		meta.SetStatusCondition(&status.Conditions, metav1.Condition{
			Type:   string(gwv1.ListenerConditionProgrammed),
			Status: metav1.ConditionTrue,
			Reason: string(gwv1.ListenerReasonProgrammed),
		})
	}
	return status
}

func parentRefConditionsWithDefaults(report *ParentRefReport) []metav1.Condition {
	var conditions []metav1.Condition
	if report != nil {
		conditions = slices.Clone(report.Conditions)
	}
	if cond := meta.FindStatusCondition(conditions, string(gwv1.RouteConditionAccepted)); cond == nil {
		meta.SetStatusCondition(&conditions, metav1.Condition{
			Type:   string(gwv1.RouteConditionAccepted),
			Status: metav1.ConditionTrue,
			Reason: string(gwv1.RouteReasonAccepted),
		})
	}
	if cond := meta.FindStatusCondition(conditions, string(gwv1.RouteConditionResolvedRefs)); cond == nil {
		meta.SetStatusCondition(&conditions, metav1.Condition{
			Type:   string(gwv1.RouteConditionResolvedRefs),
			Status: metav1.ConditionTrue,
			Reason: string(gwv1.RouteReasonResolvedRefs),
		})
	}
	return conditions
}

// listenerReport returns the ListenerReport for the named listener, or nil if
// none was recorded during translation. Indexing a nil map is safe, so callers
// may pass a nil listeners map.
func listenerReport(listeners map[string]*ListenerReport, name gwv1.SectionName) *ListenerReport {
	return listeners[string(name)]
}

func parentRefReport(routeReport *RouteReport, parentRef *gwv1.ParentReference) *ParentRefReport {
	if routeReport == nil || routeReport.Parents == nil {
		return nil
	}
	return routeReport.Parents[getParentRefKey(parentRef)]
}
