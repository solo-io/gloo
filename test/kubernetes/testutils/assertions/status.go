package assertions

import (
	"context"
	"fmt"
	"time"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"github.com/onsi/gomega/gstruct"
	errors "github.com/rotisserie/eris"
	"github.com/solo-io/gloo/test/gomega/matchers"
	"github.com/solo-io/gloo/test/helpers"
	"github.com/solo-io/gloo/test/kube2e/helper"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	gwv1 "sigs.k8s.io/gateway-api/apis/v1"
	gwv1a2 "sigs.k8s.io/gateway-api/apis/v1alpha2"
)

// Checks GetNamespacedStatuses status for gloo installation namespace
func (p *Provider) EventuallyResourceStatusMatchesWarningReasons(getter helpers.InputResourceGetter, desiredStatusReasons []string, desiredReporter string, timeout ...time.Duration) {
	ginkgo.GinkgoHelper()

	currentTimeout, pollingInterval := helper.GetTimeouts(timeout...)
	gomega.Eventually(func(g gomega.Gomega) {
		statusWarningsMatcher := matchers.MatchStatusInNamespace(
			p.glooGatewayContext.InstallNamespace,
			gomega.And(matchers.HaveWarningStateWithReasonSubstrings(desiredStatusReasons...), matchers.HaveReportedBy(desiredReporter)),
		)

		status, err := getResourceNamespacedStatus(getter)
		g.Expect(err).NotTo(gomega.HaveOccurred(), "failed to get resource namespaced status")
		g.Expect(status).ToNot(gomega.BeNil())
		g.Expect(status).To(gomega.HaveValue(statusWarningsMatcher))
	}, currentTimeout, pollingInterval).Should(gomega.Succeed())
}

func (p *Provider) EventuallyResourceStatusMatchesRejectedReasons(getter helpers.InputResourceGetter, desiredStatusReasons []string, desiredReporter string, timeout ...time.Duration) {
	ginkgo.GinkgoHelper()

	currentTimeout, pollingInterval := helper.GetTimeouts(timeout...)
	gomega.Eventually(func(g gomega.Gomega) {
		statusRejectionsMatcher := matchers.MatchStatusInNamespace(
			p.glooGatewayContext.InstallNamespace,
			gomega.And(matchers.HaveRejectedStateWithReasonSubstrings(desiredStatusReasons...), matchers.HaveReportedBy(desiredReporter)),
		)

		status, err := getResourceNamespacedStatus(getter)
		g.Expect(err).NotTo(gomega.HaveOccurred(), "failed to get resource namespaced status")
		g.Expect(status).ToNot(gomega.BeNil())
		g.Expect(status).To(gomega.HaveValue(statusRejectionsMatcher))
	}, currentTimeout, pollingInterval).Should(gomega.Succeed())
}

func (p *Provider) EventuallyResourceStatusMatchesState(
	getter helpers.InputResourceGetter,
	desiredState core.Status_State,
	desiredReporter string,
	timeout ...time.Duration,
) {
	currentTimeout, pollingInterval := helper.GetTimeouts(timeout...)
	p.Gomega.Eventually(func(g gomega.Gomega) {
		statusStateMatcher := matchers.MatchStatusInNamespace(
			p.glooGatewayContext.InstallNamespace,
			gomega.And(matchers.HaveState(desiredState), matchers.HaveReportedBy(desiredReporter)),
		)
		status, err := getResourceNamespacedStatus(getter)
		g.Expect(err).NotTo(gomega.HaveOccurred(), "failed to get resource namespaced status")
		g.Expect(status).ToNot(gomega.BeNil())
		g.Expect(status).To(gomega.HaveValue(statusStateMatcher))
	}, currentTimeout, pollingInterval).Should(gomega.Succeed())
}

func (p *Provider) EventuallyResourceStatusMatchesSubResource(
	getter helpers.InputResourceGetter,
	desiredSubresourceName string,
	desiredSubresource matchers.SoloKitSubresourceStatus,
	timeout ...time.Duration,
) {
	currentTimeout, pollingInterval := helper.GetTimeouts(timeout...)
	p.Gomega.Eventually(func(g gomega.Gomega) {
		subResourceStatusMatcher := matchers.HaveSubResourceStatusState(desiredSubresourceName, desiredSubresource)
		status, err := getResourceNamespacedStatus(getter)
		g.Expect(err).NotTo(gomega.HaveOccurred(), "failed to get resource namespaced status")
		g.Expect(status).ToNot(gomega.BeNil())
		g.Expect(status).To(gomega.HaveValue(subResourceStatusMatcher))
	}, currentTimeout, pollingInterval).Should(gomega.Succeed())
}

func getResourceNamespacedStatus(getter helpers.InputResourceGetter) (*core.NamespacedStatuses, error) {
	resource, err := getter()
	if err != nil {
		return &core.NamespacedStatuses{}, errors.Wrapf(err, "failed to get resource")
	}

	namespacedStatuses := resource.GetNamespacedStatuses()

	// In newer versions of Gloo Edge we provide a default "empty" status, which allows us to patch it to perform updates
	// As a result, a nil check isn't enough to determine that that status hasn't been reported
	if namespacedStatuses == nil || namespacedStatuses.GetStatuses() == nil {
		return &core.NamespacedStatuses{}, errors.Wrapf(err, "waiting for %v status to be non-empty", resource.GetMetadata().GetName())
	}

	return namespacedStatuses, nil
}

// EventuallyHTTPRouteStatusContainsMessage asserts that eventually at least one of the HTTPRoute's route parent statuses contains
// the given message substring.
func (p *Provider) EventuallyHTTPRouteStatusContainsMessage(
	ctx context.Context,
	routeName string,
	routeNamespace string,
	message string,
	timeout ...time.Duration) {
	currentTimeout, pollingInterval := helper.GetTimeouts(timeout...)
	p.Gomega.Eventually(func(g gomega.Gomega) {
		matcher := matchers.HaveKubeGatewayRouteStatus(&matchers.KubeGatewayRouteStatus{
			Custom: gstruct.MatchFields(gstruct.IgnoreExtras, gstruct.Fields{
				"Parents": gomega.ContainElement(gstruct.MatchFields(gstruct.IgnoreExtras, gstruct.Fields{
					"Conditions": gomega.ContainElement(gstruct.MatchFields(gstruct.IgnoreExtras, gstruct.Fields{
						"Message": matchers.ContainSubstrings([]string{message}),
					})),
				})),
			}),
		})

		route := &gwv1.HTTPRoute{}
		err := p.clusterContext.Client.Get(ctx, types.NamespacedName{Name: routeName, Namespace: routeNamespace}, route)
		g.Expect(err).NotTo(gomega.HaveOccurred(), "can get httproute")
		g.Expect(route.Status.RouteStatus).To(gomega.HaveValue(matcher))
	}, currentTimeout, pollingInterval).Should(gomega.Succeed())
}

// EventuallyHTTPRouteStatusContainsReason asserts that eventually at least one of the HTTPRoute's route parent statuses contains
// the given reason substring.
func (p *Provider) EventuallyHTTPRouteStatusContainsReason(
	ctx context.Context,
	routeName string,
	routeNamespace string,
	reason string,
	timeout ...time.Duration) {
	currentTimeout, pollingInterval := helper.GetTimeouts(timeout...)
	p.Gomega.Eventually(func(g gomega.Gomega) {
		matcher := matchers.HaveKubeGatewayRouteStatus(&matchers.KubeGatewayRouteStatus{
			Custom: gstruct.MatchFields(gstruct.IgnoreExtras, gstruct.Fields{
				"Parents": gomega.ContainElement(gstruct.MatchFields(gstruct.IgnoreExtras, gstruct.Fields{
					"Conditions": gomega.ContainElement(gstruct.MatchFields(gstruct.IgnoreExtras, gstruct.Fields{
						"Reason": matchers.ContainSubstrings([]string{reason}),
					})),
				})),
			}),
		})

		route := &gwv1.HTTPRoute{
			ObjectMeta: metav1.ObjectMeta{
				Name:      routeName,
				Namespace: routeNamespace,
			},
		}
		err := p.clusterContext.Client.Get(ctx, types.NamespacedName{Name: routeName, Namespace: routeNamespace}, route)
		g.Expect(err).NotTo(gomega.HaveOccurred(), "can get httproute")
		g.Expect(route.Status.RouteStatus).To(gomega.HaveValue(matcher))
	}, currentTimeout, pollingInterval).Should(gomega.Succeed())
}

// EventuallyGatewayCondition checks the provided Gateway condition is set to expect.
func (p *Provider) EventuallyGatewayCondition(
	ctx context.Context,
	gatewayName string,
	gatewayNamespace string,
	cond gwv1.GatewayConditionType,
	expect metav1.ConditionStatus,
	timeout ...time.Duration,
) {
	ginkgo.GinkgoHelper()
	currentTimeout, pollingInterval := helper.GetTimeouts(timeout...)
	p.Gomega.Eventually(func(g gomega.Gomega) {
		gateway := &gwv1.Gateway{}
		err := p.clusterContext.Client.Get(ctx, types.NamespacedName{Name: gatewayName, Namespace: gatewayNamespace}, gateway)
		g.Expect(err).NotTo(gomega.HaveOccurred(), fmt.Sprintf("failed to get Gateway %s/%s", gatewayNamespace, gatewayName))

		condition := getConditionByType(gateway.Status.Conditions, string(cond))
		g.Expect(condition).NotTo(gomega.BeNil(), fmt.Sprintf("%v condition not found for Gateway %s/%s", cond, gatewayNamespace, gatewayName))
		g.Expect(condition.Status).To(gomega.Equal(expect), fmt.Sprintf("%v condition is not %v for Gateway %s/%s",
			cond, expect, gatewayNamespace, gatewayName))
	}, currentTimeout, pollingInterval).Should(gomega.Succeed())
}

// EventuallyGatewayListenerAttachedRoutes checks the provided Gateway contains the expected attached routes for the listener.
func (p *Provider) EventuallyGatewayListenerAttachedRoutes(
	ctx context.Context,
	gatewayName string,
	gatewayNamespace string,
	listener gwv1.SectionName,
	routes int32,
	timeout ...time.Duration,
) {
	ginkgo.GinkgoHelper()
	currentTimeout, pollingInterval := helper.GetTimeouts(timeout...)
	p.Gomega.Eventually(func(g gomega.Gomega) {
		gateway := &gwv1.Gateway{}
		err := p.clusterContext.Client.Get(ctx, types.NamespacedName{Name: gatewayName, Namespace: gatewayNamespace}, gateway)
		g.Expect(err).NotTo(gomega.HaveOccurred(), fmt.Sprintf("failed to get Gateway %s/%s", gatewayNamespace, gatewayName))

		found := false
		for _, l := range gateway.Status.Listeners {
			if l.Name == listener {
				found = true
				g.Expect(l.AttachedRoutes).To(gomega.Equal(routes), fmt.Sprintf("%v listener does not contain %d attached routes for Gateway %s/%s",
					l, routes, gatewayNamespace, gatewayName))
			}
		}
		g.Expect(found).To(gomega.BeTrue(), fmt.Sprintf("%v listener not found for Gateway %s/%s", listener, gatewayNamespace, gatewayName))
	}, currentTimeout, pollingInterval).Should(gomega.Succeed())
}

// EventuallyTCPRouteCondition checks that provided TCPRoute condition is set to expect.
func (p *Provider) EventuallyTCPRouteCondition(
	ctx context.Context,
	routeName string,
	routeNamespace string,
	cond gwv1.RouteConditionType,
	expect metav1.ConditionStatus,
	timeout ...time.Duration,
) {
	ginkgo.GinkgoHelper()
	currentTimeout, pollingInterval := helper.GetTimeouts(timeout...)
	p.Gomega.Eventually(func(g gomega.Gomega) {
		route := &gwv1a2.TCPRoute{}
		err := p.clusterContext.Client.Get(ctx, types.NamespacedName{Name: routeName, Namespace: routeNamespace}, route)
		g.Expect(err).NotTo(gomega.HaveOccurred(), "failed to get TCPRoute %s/%s", routeNamespace, routeName)

		var conditionFound bool
		for _, parentStatus := range route.Status.Parents {
			condition := getConditionByType(parentStatus.Conditions, string(cond))
			if condition != nil && condition.Status == expect {
				conditionFound = true
				break
			}
		}
		g.Expect(conditionFound).To(gomega.BeTrue(), fmt.Sprintf("%v condition is not %v for any parent of TCPRoute %s/%s",
			cond, expect, routeNamespace, routeName))
	}, currentTimeout, pollingInterval).Should(gomega.Succeed())
}

// Helper function to retrieve a condition by type from a list of conditions.
func getConditionByType(conditions []metav1.Condition, conditionType string) *metav1.Condition {
	for _, condition := range conditions {
		if condition.Type == conditionType {
			return &condition
		}
	}
	return nil
}
