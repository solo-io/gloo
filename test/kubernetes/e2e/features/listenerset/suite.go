package listenerset

import (
	"context"

	glooschemes "github.com/solo-io/gloo/pkg/schemes"
	"github.com/solo-io/gloo/pkg/utils/kubeutils"
	"github.com/solo-io/gloo/pkg/utils/requestutils/curl"
	"github.com/solo-io/gloo/projects/gateway2/translator/listener"
	"github.com/solo-io/gloo/projects/gateway2/wellknown"
	"github.com/solo-io/gloo/test/kubernetes/e2e"
	"github.com/solo-io/gloo/test/kubernetes/e2e/defaults"
	"github.com/solo-io/gloo/test/kubernetes/e2e/tests/base"
	"github.com/stretchr/testify/suite"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	gwv1 "sigs.k8s.io/gateway-api/apis/v1"
	gwxv1a1 "sigs.k8s.io/gateway-api/apisx/v1alpha1"
)

var _ e2e.NewSuiteFunc = NewTestingSuite

type testingSuite struct {
	*base.BaseTestingSuite
}

func NewTestingSuite(ctx context.Context, testInst *e2e.TestInstallation) suite.TestingSuite {
	return &testingSuite{
		base.NewBaseTestingSuite(ctx, testInst, setup, testCases),
	}
}

func (s *testingSuite) SetupSuite() {
	if !RequiredCrdExists(s.TestInstallation) {
		s.T().Skip("Skipping as the XListenerSet CRD is not installed")
	}

	s.BaseTestingSuite.SetupSuite()
}

func (s *testingSuite) TestValidListenerSet() {
	s.expectListenerSetAccepted(validListenerSet)

	// The route attached to the gateway should work on the listener defined on the gateway
	s.TestInstallation.AssertionsT(s.T()).AssertEventualCurlResponse(
		s.Ctx,
		defaults.CurlPodExecOpt,
		[]curl.Option{
			curl.WithHost(kubeutils.ServiceFQDN(proxyService.ObjectMeta)),
			curl.WithPort(8080),
			curl.WithHostHeader("example.com"),
		},
		expectOK)

	// The route attached to the listenerset should NOT work on the listener defined on the gateway
	s.TestInstallation.AssertionsT(s.T()).AssertEventualCurlResponse(
		s.Ctx,
		defaults.CurlPodExecOpt,
		[]curl.Option{
			curl.WithHost(kubeutils.ServiceFQDN(proxyService.ObjectMeta)),
			curl.WithPort(8080),
			curl.WithHostHeader("listenerset.com"),
		},
		expectNotFound)

	// The route attached to the gateway should work on the listener defined on the listener set
	s.TestInstallation.AssertionsT(s.T()).AssertEventualCurlResponse(
		s.Ctx,
		defaults.CurlPodExecOpt,
		[]curl.Option{
			curl.WithHost(kubeutils.ServiceFQDN(proxyService.ObjectMeta)),
			curl.WithPort(8081),
			curl.WithHostHeader("example.com"),
		},
		expectOK)

	// The route attached to the listenerset should work on the listener defined on the listener set
	s.TestInstallation.AssertionsT(s.T()).AssertEventualCurlResponse(
		s.Ctx,
		defaults.CurlPodExecOpt,
		[]curl.Option{
			curl.WithHost(kubeutils.ServiceFQDN(proxyService.ObjectMeta)),
			curl.WithPort(8081),
			curl.WithHostHeader("listenerset.com"),
		},
		expectOK)
}

func (s *testingSuite) TestInvalidListenerSetNotAllowed() {
	s.expectListenerSetNotAllowed(invalidListenerSetNotAllowed)

	// The route attached to the gateway should work on the listener defined on the gateway
	s.TestInstallation.AssertionsT(s.T()).AssertEventualCurlResponse(
		s.Ctx,
		defaults.CurlPodExecOpt,
		[]curl.Option{
			curl.WithHost(kubeutils.ServiceFQDN(proxyService.ObjectMeta)),
			curl.WithPort(8080),
			curl.WithHostHeader("example.com"),
		},
		expectOK)

	// The listener defined on the invalid listenerset should not work
	s.TestInstallation.AssertionsT(s.T()).AssertEventualCurlError(
		s.Ctx,
		defaults.CurlPodExecOpt,
		[]curl.Option{
			curl.WithHost(kubeutils.ServiceFQDN(proxyService.ObjectMeta)),
			curl.WithPort(8081),
			curl.WithHostHeader("example.com"),
		},
		curlExitErrorCode)

}

func (s *testingSuite) TestInvalidListenerSetNonExistingGW() {
	s.expectListenerSetUnknown(invalidListenerSetNonExistingGW)

	// The route attached to the gateway should work on the listener defined on the gateway
	s.TestInstallation.AssertionsT(s.T()).AssertEventualCurlResponse(
		s.Ctx,
		defaults.CurlPodExecOpt,
		[]curl.Option{
			curl.WithHost(kubeutils.ServiceFQDN(proxyService.ObjectMeta)),
			curl.WithPort(8080),
			curl.WithHostHeader("example.com"),
		},
		expectOK)

	// The listener defined on the invalid listenerset should not work
	s.TestInstallation.AssertionsT(s.T()).AssertEventualCurlError(
		s.Ctx,
		defaults.CurlPodExecOpt,
		[]curl.Option{
			curl.WithHost(kubeutils.ServiceFQDN(proxyService.ObjectMeta)),
			curl.WithPort(8081),
			curl.WithHostHeader("example.com"),
		},
		curlExitErrorCode)
}

func (s *testingSuite) expectListenerSetAccepted(namespacedName types.NamespacedName) {
	s.TestInstallation.AssertionsT(s.T()).EventuallyGatewayCondition(s.Ctx, gatewayObjectMeta.Name, glooProxyObjectMeta.Namespace, listener.AttachedListenerSetsConditionType, metav1.ConditionTrue)

	s.TestInstallation.AssertionsT(s.T()).EventuallyListenerSetStatus(s.Ctx, namespacedName.Name, namespacedName.Namespace,
		gwxv1a1.ListenerSetStatus{
			Conditions: []metav1.Condition{
				{
					Type:   string(gwxv1a1.ListenerSetConditionAccepted),
					Status: metav1.ConditionTrue,
					Reason: string(gwxv1a1.ListenerSetReasonAccepted),
				},
				{
					Type:   string(gwxv1a1.ListenerSetConditionProgrammed),
					Status: metav1.ConditionTrue,
					Reason: string(gwxv1a1.ListenerSetReasonProgrammed),
				},
			},
			Listeners: []gwxv1a1.ListenerEntryStatus{
				{
					Name:           "http",
					Port:           8081,
					AttachedRoutes: 2,
					Conditions: []metav1.Condition{
						{
							Type:   string(gwxv1a1.ListenerEntryConditionAccepted),
							Status: metav1.ConditionTrue,
							Reason: string(gwxv1a1.ListenerEntryReasonAccepted),
						},
						{
							Type:   string(gwxv1a1.ListenerEntryConditionConflicted),
							Status: metav1.ConditionFalse,
							Reason: string(gwv1.ListenerReasonNoConflicts),
						},
						{
							Type:   string(gwxv1a1.ListenerEntryConditionResolvedRefs),
							Status: metav1.ConditionTrue,
							Reason: string(gwxv1a1.ListenerEntryReasonResolvedRefs),
						},
						{
							Type:   string(gwxv1a1.ListenerEntryConditionProgrammed),
							Status: metav1.ConditionTrue,
							Reason: string(gwxv1a1.ListenerEntryReasonProgrammed),
						},
					},
				},
			},
		})
}

func (s *testingSuite) expectListenerSetNotAllowed(namespacedName types.NamespacedName) {
	s.TestInstallation.AssertionsT(s.T()).EventuallyGatewayCondition(s.Ctx, gatewayObjectMeta.Name, glooProxyObjectMeta.Namespace, listener.AttachedListenerSetsConditionType, metav1.ConditionFalse)

	s.TestInstallation.AssertionsT(s.T()).EventuallyListenerSetStatus(s.Ctx, namespacedName.Name, namespacedName.Namespace,
		gwxv1a1.ListenerSetStatus{
			Conditions: []metav1.Condition{
				{
					Type:   string(gwxv1a1.ListenerSetConditionAccepted),
					Status: metav1.ConditionFalse,
					Reason: string(gwxv1a1.ListenerSetReasonNotAllowed),
				},
				{
					Type:   string(gwxv1a1.ListenerSetConditionProgrammed),
					Status: metav1.ConditionFalse,
					Reason: string(gwxv1a1.ListenerSetReasonNotAllowed),
				},
			},
		})
}

func (s *testingSuite) expectListenerSetUnknown(namespacedName types.NamespacedName) {
	s.TestInstallation.AssertionsT(s.T()).EventuallyGatewayCondition(s.Ctx, gatewayObjectMeta.Name, glooProxyObjectMeta.Namespace, listener.AttachedListenerSetsConditionType, metav1.ConditionFalse)

	s.TestInstallation.AssertionsT(s.T()).EventuallyListenerSetStatus(s.Ctx, namespacedName.Name, namespacedName.Namespace,
		gwxv1a1.ListenerSetStatus{
			Conditions: []metav1.Condition{
				{
					Type:   string(gwxv1a1.ListenerSetConditionAccepted),
					Status: metav1.ConditionUnknown,
				},
				{
					Type:   string(gwxv1a1.ListenerSetConditionProgrammed),
					Status: metav1.ConditionUnknown,
				},
			},
		})
}

func RequiredCrdExists(testInstallation *e2e.TestInstallation) bool {
	xListenerSetExists, err := glooschemes.CRDExists(testInstallation.ClusterContext.RestConfig, gwxv1a1.GroupVersion.Group, gwxv1a1.GroupVersion.Version, wellknown.XListenerSetKind)
	testInstallation.Assertions.Assert.NoError(err)
	return xListenerSetExists
}
