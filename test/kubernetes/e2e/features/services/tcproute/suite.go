package tcproute

import (
	"context"

	"github.com/stretchr/testify/suite"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "sigs.k8s.io/gateway-api/apis/v1"

	"github.com/solo-io/gloo/pkg/utils/kubeutils"
	"github.com/solo-io/gloo/pkg/utils/requestutils/curl"
	"github.com/solo-io/gloo/test/kubernetes/e2e"
	"github.com/solo-io/gloo/test/kubernetes/e2e/defaults"
)

// testingSuite is the entire Suite of tests for testing K8s Service-specific features/fixes
type testingSuite struct {
	suite.Suite

	ctx context.Context

	// testInstallation contains all the metadata/utilities necessary to execute a series of tests
	// against an installation of Gloo Gateway
	testInstallation *e2e.TestInstallation
}

func NewTestingSuite(ctx context.Context, testInst *e2e.TestInstallation) suite.TestingSuite {
	return &testingSuite{
		ctx:              ctx,
		testInstallation: testInst,
	}
}

func (s *testingSuite) TestConfigureTCPRouteBackingDestinationsWithSingleService() {
	s.T().Cleanup(func() {
		err := s.testInstallation.Actions.Kubectl().DeleteFile(s.ctx, singleTcpRouteManifest)
		s.NoError(err, "can delete manifest")
		err = s.testInstallation.Actions.Kubectl().DeleteFile(s.ctx, multiBackendServiceManifest)
		s.NoError(err, "can delete manifest")
		err = s.testInstallation.Actions.Kubectl().DeleteFile(s.ctx, singleListenerGatewayAndClientManifest)
		s.NoError(err, "can delete manifest")
		s.testInstallation.Assertions.EventuallyObjectsNotExist(s.ctx, proxyService, proxyDeployment)
	})

	err := s.testInstallation.Actions.Kubectl().ApplyFile(s.ctx, singleListenerGatewayAndClientManifest)
	s.Assert().NoError(err, "can apply gateway and client manifest")
	s.testInstallation.Assertions.EventuallyGatewayCondition(s.ctx, "tcp-gateway", "default", v1.GatewayConditionAccepted, metav1.ConditionTrue, timeout)

	err = s.testInstallation.Actions.Kubectl().ApplyFile(s.ctx, multiBackendServiceManifest)
	s.Assert().NoError(err, "can apply backend service manifest")
	s.testInstallation.Assertions.EventuallyObjectsExist(s.ctx, proxyService, proxyDeployment)

	err = s.testInstallation.Actions.Kubectl().ApplyFile(s.ctx, singleTcpRouteManifest)
	s.Assert().NoError(err, "can apply tcproute manifest")
	s.testInstallation.Assertions.EventuallyTCPRouteCondition(s.ctx, "tcp-app-1", "default", v1.RouteConditionAccepted, metav1.ConditionTrue, timeout)
	s.testInstallation.Assertions.EventuallyTCPRouteCondition(s.ctx, "tcp-app-1", "default", v1.RouteConditionResolvedRefs, metav1.ConditionTrue, timeout)
	s.testInstallation.Assertions.EventuallyGatewayCondition(s.ctx, "tcp-gateway", "default", v1.GatewayConditionProgrammed, metav1.ConditionTrue, timeout)
	s.testInstallation.Assertions.EventuallyGatewayListenerAttachedRoutes(s.ctx, "tcp-gateway", "default", v1.SectionName("foo"), int32(1), timeout)

	s.testInstallation.Assertions.AssertEventualCurlResponse(
		s.ctx,
		defaults.CurlPodExecOpt,
		[]curl.Option{
			curl.WithHost(kubeutils.ServiceFQDN(proxyService.ObjectMeta)),
			curl.WithPort(8088),
		},
		expectedTcpFooSvcResp)
	s.testInstallation.Assertions.AssertEventualCurlResponse(
		s.ctx,
		defaults.CurlPodExecOpt,
		[]curl.Option{
			curl.WithHost(kubeutils.ServiceFQDN(proxyService.ObjectMeta)),
			curl.WithPort(8088),
		},
		expectedTcpBarSvcResp)
}

func (s *testingSuite) TestConfigureTCPRouteBackingDestinationsWithMultiServices() {
	s.T().Skip("skipping test until we resolve the reason it is flaky")

	s.T().Cleanup(func() {
		err := s.testInstallation.Actions.Kubectl().DeleteFile(s.ctx, multiTcpRouteManifest)
		s.NoError(err, "can delete manifest")
		err = s.testInstallation.Actions.Kubectl().DeleteFile(s.ctx, multiBackendServiceManifest)
		s.NoError(err, "can delete manifest")
		err = s.testInstallation.Actions.Kubectl().DeleteFile(s.ctx, multiListenerGatewayAndClientManifest)
		s.NoError(err, "can delete manifest")
		s.testInstallation.Assertions.EventuallyObjectsNotExist(s.ctx, proxyService, proxyDeployment)
	})

	err := s.testInstallation.Actions.Kubectl().ApplyFile(s.ctx, multiListenerGatewayAndClientManifest)
	s.Assert().NoError(err, "can apply gateway and client manifest")
	s.testInstallation.Assertions.EventuallyGatewayCondition(s.ctx, "tcp-gateway", "default", v1.GatewayConditionAccepted, metav1.ConditionTrue, timeout)

	err = s.testInstallation.Actions.Kubectl().ApplyFile(s.ctx, multiBackendServiceManifest)
	s.Assert().NoError(err, "can apply backend service manifest")
	s.testInstallation.Assertions.EventuallyObjectsExist(s.ctx, proxyService, proxyDeployment)

	err = s.testInstallation.Actions.Kubectl().ApplyFile(s.ctx, multiTcpRouteManifest)
	s.Assert().NoError(err, "can apply tcproute manifest")
	s.testInstallation.Assertions.EventuallyTCPRouteCondition(s.ctx, "tcp-app-1", "default", v1.RouteConditionAccepted, metav1.ConditionTrue, timeout)
	s.testInstallation.Assertions.EventuallyTCPRouteCondition(s.ctx, "tcp-app-1", "default", v1.RouteConditionResolvedRefs, metav1.ConditionTrue, timeout)
	s.testInstallation.Assertions.EventuallyTCPRouteCondition(s.ctx, "tcp-app-2", "default", v1.RouteConditionAccepted, metav1.ConditionTrue, timeout)
	s.testInstallation.Assertions.EventuallyTCPRouteCondition(s.ctx, "tcp-app-2", "default", v1.RouteConditionResolvedRefs, metav1.ConditionTrue, timeout)
	s.testInstallation.Assertions.EventuallyGatewayCondition(s.ctx, "tcp-gateway", "default", v1.GatewayConditionProgrammed, metav1.ConditionTrue, timeout)
	s.testInstallation.Assertions.EventuallyGatewayListenerAttachedRoutes(s.ctx, "tcp-gateway", "default", v1.SectionName("foo"), int32(1), timeout)
	s.testInstallation.Assertions.EventuallyGatewayListenerAttachedRoutes(s.ctx, "tcp-gateway", "default", v1.SectionName("bar"), int32(1), timeout)

	s.testInstallation.Assertions.AssertEventualCurlResponse(
		s.ctx,
		defaults.CurlPodExecOpt,
		[]curl.Option{
			curl.WithHost(kubeutils.ServiceFQDN(proxyService.ObjectMeta)),
			curl.WithPort(8088),
		},
		expectedTcpFooSvcResp)
	s.testInstallation.Assertions.AssertEventualCurlResponse(
		s.ctx,
		defaults.CurlPodExecOpt,
		[]curl.Option{
			curl.WithHost(kubeutils.ServiceFQDN(proxyService.ObjectMeta)),
			curl.WithPort(8089),
		},
		expectedTcpBarSvcResp)
}
