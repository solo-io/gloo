package reports

import (
	"context"

	"github.com/solo-io/go-utils/contextutils"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	gwv1 "sigs.k8s.io/gateway-api/apis/v1"
	gwv1alpha2 "sigs.k8s.io/gateway-api/apis/v1alpha2"
	gwxv1alpha1 "sigs.k8s.io/gateway-api/apisx/v1alpha1"
)

type ReportMap struct {
	Gateways     map[types.NamespacedName]*GatewayReport
	ListenerSets map[types.NamespacedName]*ListenerSetReport
	HTTPRoutes   map[types.NamespacedName]*RouteReport
	TCPRoutes    map[types.NamespacedName]*RouteReport
	TLSRoutes    map[types.NamespacedName]*RouteReport
}

type GatewayReport struct {
	conditions         []metav1.Condition
	listeners          map[string]*ListenerReport
	observedGeneration int64
}

type ListenerSetReport struct {
	conditions         []metav1.Condition
	listeners          map[string]*ListenerReport
	observedGeneration int64
}

type ListenerReport struct {
	Status gwv1.ListenerStatus
}

type RouteReport struct {
	Parents            map[ParentRefKey]*ParentRefReport
	observedGeneration int64
}

// TODO: rename to e.g. RouteParentRefReport
type ParentRefReport struct {
	Conditions []metav1.Condition
}

type ParentRefKey struct {
	Group string
	Kind  string
	types.NamespacedName
}

func NewReportMap() ReportMap {
	gateways := make(map[types.NamespacedName]*GatewayReport)
	listenerSets := make(map[types.NamespacedName]*ListenerSetReport)
	httpRoutes := make(map[types.NamespacedName]*RouteReport)
	tcpRoutes := make(map[types.NamespacedName]*RouteReport)
	tlsRoutes := make(map[types.NamespacedName]*RouteReport)
	return ReportMap{
		Gateways:     gateways,
		ListenerSets: listenerSets,
		HTTPRoutes:   httpRoutes,
		TCPRoutes:    tcpRoutes,
		TLSRoutes:    tlsRoutes,
	}
}

// Returns a GatewayReport for the provided Gateway, nil if there is not a report present.
// This is different than the Reporter.Gateway() method, as we need to understand when
// reports are not generated for a Gateway that has been translated.
//
// NOTE: Exported for unit testing, validation_test.go should be refactored to reduce this visibility
func (r *ReportMap) Gateway(gateway *gwv1.Gateway) *GatewayReport {
	key := client.ObjectKeyFromObject(gateway)
	return r.Gateways[key]
}

func (r *ReportMap) newGatewayReport(gateway *gwv1.Gateway) *GatewayReport {
	gr := &GatewayReport{}
	gr.observedGeneration = gateway.Generation
	key := client.ObjectKeyFromObject(gateway)
	r.Gateways[key] = gr
	return gr
}

// Returns a ListenerSetReport for the provided ListenerSet, nil if there is not a report present.
// This is different than the Reporter.ListenerSet() method, as we need to understand when
// reports are not generated for a ListenerSet that has been translated.
//
// NOTE: Exported for unit testing, validation_test.go should be refactored to reduce this visibility
func (r *ReportMap) ListenerSet(listenerSet *gwxv1alpha1.XListenerSet) *ListenerSetReport {
	key := client.ObjectKeyFromObject(listenerSet)
	return r.ListenerSets[key]
}

func (r *ReportMap) newListenerSetReport(listenerSet *gwxv1alpha1.XListenerSet) *ListenerSetReport {
	lsr := &ListenerSetReport{}
	lsr.observedGeneration = listenerSet.Generation
	key := client.ObjectKeyFromObject(listenerSet)
	r.ListenerSets[key] = lsr
	return lsr
}

// route returns a RouteReport for the provided route object, nil if a report is not present.
// This is different than the Reporter.Route() method, as we need to understand when
// reports are not generated for a route that has been translated. Supported object types are:
//
// * HTTPRoute
// * TCPRoute
func (r *ReportMap) route(obj client.Object) *RouteReport {
	key := client.ObjectKeyFromObject(obj)

	switch obj.(type) {
	case *gwv1.HTTPRoute:
		return r.HTTPRoutes[key]
	case *gwv1alpha2.TCPRoute:
		return r.TCPRoutes[key]
	case *gwv1alpha2.TLSRoute:
		return r.TLSRoutes[key]
	default:
		contextutils.LoggerFrom(context.TODO()).Warnf("Unsupported route type: %T", obj)
		return nil
	}
}

func (r *ReportMap) newRouteReport(obj client.Object) *RouteReport {
	rr := &RouteReport{
		observedGeneration: obj.GetGeneration(),
	}

	key := client.ObjectKeyFromObject(obj)

	switch obj.(type) {
	case *gwv1.HTTPRoute:
		r.HTTPRoutes[key] = rr
	case *gwv1alpha2.TCPRoute:
		r.TCPRoutes[key] = rr
	case *gwv1alpha2.TLSRoute:
		r.TLSRoutes[key] = rr
	default:
		contextutils.LoggerFrom(context.TODO()).Warnf("Unsupported route type: %T", obj)
		return nil
	}

	return rr
}

func (g *GatewayReport) Listener(listener *gwv1.Listener) ListenerReporter {
	return g.listener(listener)
}

func (g *GatewayReport) listener(listener *gwv1.Listener) *ListenerReport {
	if g.listeners == nil {
		g.listeners = make(map[string]*ListenerReport)
	}

	// Return the ListenerReport if it already exists
	if lr, exists := g.listeners[string(listener.Name)]; exists {
		return lr
	}

	// Create and add the new ListenerReport if it doesn't exist
	lr := NewListenerReport(string(listener.Name))
	g.listeners[string(listener.Name)] = lr
	return lr
}

func (g *GatewayReport) GetConditions() []metav1.Condition {
	if g == nil {
		return []metav1.Condition{}
	}
	return g.conditions
}

func (g *GatewayReport) SetCondition(gc GatewayCondition) {
	condition := metav1.Condition{
		Type:    string(gc.Type),
		Status:  gc.Status,
		Reason:  string(gc.Reason),
		Message: gc.Message,
	}
	meta.SetStatusCondition(&g.conditions, condition)
}

func (g *ListenerSetReport) Listener(listener *gwv1.Listener) ListenerReporter {
	return g.listener(listener)
}

func (g *ListenerSetReport) listener(listener *gwv1.Listener) *ListenerReport {
	if g.listeners == nil {
		g.listeners = make(map[string]*ListenerReport)
	}

	// Return the ListenerReport if it already exists
	if lr, exists := g.listeners[string(listener.Name)]; exists {
		return lr
	}

	// Create and add the new ListenerReport if it doesn't exist
	lr := NewListenerReport(string(listener.Name))
	g.listeners[string(listener.Name)] = lr
	return lr
}

func (g *ListenerSetReport) GetConditions() []metav1.Condition {
	if g == nil {
		return []metav1.Condition{}
	}
	return g.conditions
}

func (g *ListenerSetReport) SetCondition(gc GatewayCondition) {
	condition := metav1.Condition{
		Type:    string(gc.Type),
		Status:  gc.Status,
		Reason:  string(gc.Reason),
		Message: gc.Message,
	}
	g.conditions = append(g.conditions, condition)
}

func NewListenerReport(name string) *ListenerReport {
	lr := ListenerReport{}
	lr.Status.Name = gwv1.SectionName(name)
	return &lr
}

func (l *ListenerReport) SetCondition(lc ListenerCondition) {
	condition := metav1.Condition{
		Type:    string(lc.Type),
		Status:  lc.Status,
		Reason:  string(lc.Reason),
		Message: lc.Message,
	}
	meta.SetStatusCondition(&l.Status.Conditions, condition)
}

func (l *ListenerReport) SetSupportedKinds(rgks []gwv1.RouteGroupKind) {
	l.Status.SupportedKinds = rgks
}

func (l *ListenerReport) SetAttachedRoutes(n uint) {
	l.Status.AttachedRoutes = int32(n)
}

type reporter struct {
	report *ReportMap
}

func (r *reporter) Gateway(gateway *gwv1.Gateway) GatewayReporter {
	gr := r.report.Gateway(gateway)
	if gr == nil {
		gr = r.report.newGatewayReport(gateway)
	}
	return gr
}

func (r *reporter) ListenerSet(listenerSet *gwxv1alpha1.XListenerSet) ListenerSetReporter {
	lsr := r.report.ListenerSet(listenerSet)
	if lsr == nil {
		lsr = r.report.newListenerSetReport(listenerSet)
	}
	return lsr
}

func (r *reporter) Route(obj client.Object) RouteReporter {
	rr := r.report.route(obj)
	if rr == nil {
		rr = r.report.newRouteReport(obj)
	}
	return rr
}

// TODO: flesh out
func getParentRefKey(parentRef *gwv1.ParentReference) ParentRefKey {
	var group string
	if parentRef.Group != nil {
		group = string(*parentRef.Group)
	}
	var kind string
	if parentRef.Kind != nil {
		kind = string(*parentRef.Kind)
	}
	var ns string
	if parentRef.Namespace != nil {
		ns = string(*parentRef.Namespace)
	}
	return ParentRefKey{
		Group: group,
		Kind:  kind,
		NamespacedName: types.NamespacedName{
			Namespace: ns,
			Name:      string(parentRef.Name),
		},
	}
}

func (r *RouteReport) parentRef(parentRef *gwv1.ParentReference) *ParentRefReport {
	key := getParentRefKey(parentRef)
	if r.Parents == nil {
		r.Parents = make(map[ParentRefKey]*ParentRefReport)
	}
	var prr *ParentRefReport
	prr, ok := r.Parents[key]
	if !ok {
		prr = &ParentRefReport{}
		r.Parents[key] = prr
	}
	return prr
}

// parentRefs returns a list of ParentReferences associated with the RouteReport.
// It is used to update the Status of delegatee routes who may not specify
// the parentRefs field.
func (r *RouteReport) parentRefs() []gwv1.ParentReference {
	var refs []gwv1.ParentReference
	for key := range r.Parents {
		parentRef := gwv1.ParentReference{
			Group:     ptr.To(gwv1.Group(key.Group)),
			Kind:      ptr.To(gwv1.Kind(key.Kind)),
			Name:      gwv1.ObjectName(key.Name),
			Namespace: ptr.To(gwv1.Namespace(key.Namespace)),
		}
		refs = append(refs, parentRef)
	}
	return refs
}

func (r *RouteReport) ParentRef(parentRef *gwv1.ParentReference) ParentRefReporter {
	return r.parentRef(parentRef)
}

func (prr *ParentRefReport) SetCondition(rc RouteCondition) {
	condition := metav1.Condition{
		Type:    string(rc.Type),
		Status:  rc.Status,
		Reason:  string(rc.Reason),
		Message: rc.Message,
	}
	meta.SetStatusCondition(&prr.Conditions, condition)
}

func NewReporter(reportMap *ReportMap) Reporter {
	return &reporter{report: reportMap}
}

type Reporter interface {
	Gateway(gateway *gwv1.Gateway) GatewayReporter
	ListenerSet(listenerSet *gwxv1alpha1.XListenerSet) ListenerSetReporter
	Route(obj client.Object) RouteReporter
}

type GatewayReporter interface {
	Listener(listener *gwv1.Listener) ListenerReporter
	SetCondition(condition GatewayCondition)
}

type ListenerSetReporter interface {
	Listener(listener *gwv1.Listener) ListenerReporter
	SetCondition(condition GatewayCondition)
}

type ListenerReporter interface {
	SetCondition(ListenerCondition)
	SetSupportedKinds([]gwv1.RouteGroupKind)
	SetAttachedRoutes(n uint)
}

type RouteReporter interface {
	ParentRef(parentRef *gwv1.ParentReference) ParentRefReporter
}

type ParentRefReporter interface {
	SetCondition(condition RouteCondition)
}

type GatewayCondition struct {
	Type    gwv1.GatewayConditionType
	Status  metav1.ConditionStatus
	Reason  gwv1.GatewayConditionReason
	Message string
}

type ListenerCondition struct {
	Type    gwv1.ListenerConditionType
	Status  metav1.ConditionStatus
	Reason  gwv1.ListenerConditionReason
	Message string
}

type RouteCondition struct {
	Type    gwv1.RouteConditionType
	Status  metav1.ConditionStatus
	Reason  gwv1.RouteConditionReason
	Message string
}
