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
		fmt.Printf("status: %v\n", status)
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
		fmt.Printf("status: %v\n", status)
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
