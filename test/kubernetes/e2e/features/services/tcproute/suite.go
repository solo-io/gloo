package tcproute

import (
	"context"
	"fmt"
	"os"

	"github.com/stretchr/testify/suite"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "sigs.k8s.io/gateway-api/apis/v1"

	"github.com/solo-io/gloo/pkg/utils/kubeutils"
	"github.com/solo-io/gloo/pkg/utils/kubeutils/kubectl"
	"github.com/solo-io/gloo/pkg/utils/requestutils/curl"
	"github.com/solo-io/gloo/test/gomega/matchers"
	"github.com/solo-io/gloo/test/kubernetes/e2e"
	"github.com/solo-io/gloo/test/kubernetes/e2e/defaults"
)

// testingSuite is the entire suite of tests for testing K8s Service-specific features/fixes
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

func (s *testingSuite) SetupSuite() {
	var cancel context.CancelFunc
	s.ctx, cancel = context.WithTimeout(context.Background(), ctxTimeout)
	s.T().Cleanup(cancel)

	manifests := []string{
		singleSvcNsManifest,
		singleSvcGatewayAndClientManifest,
		singleSvcBackendManifest,
		singleSvcTcpRouteManifest,
		multiSvcNsManifest,
		multiSvcGatewayAndClientManifest,
		multiSvcBackendManifest,
		multiSvcTcpRouteManifest,
	}
	for _, file := range manifests {
		s.Require().NoError(validateManifestFile(file), "Invalid manifest file: %s", file)
	}
}

type tcpRouteTestCase struct {
	name                string
	nsManifest          string
	gtwName             string
	gtwNs               string
	gtwManifest         string
	svcManifest         string
	tcpRouteManifest    string
	proxyService        *corev1.Service
	proxyDeployment     *appsv1.Deployment
	expectedResponses   []*matchers.HttpResponse
	ports               []int
	listenerNames       []v1.SectionName
	expectedRouteCounts []int32
	tcpRouteNames       []string
}

func (s *testingSuite) TestConfigureTCPRouteBackingDestinations() {
	testCases := []tcpRouteTestCase{
		{
			name:             "SingleServiceTCPRoute",
			nsManifest:       singleSvcNsManifest,
			gtwName:          singleSvcGatewayName,
			gtwNs:            singleSvcNsName,
			gtwManifest:      singleSvcGatewayAndClientManifest,
			svcManifest:      singleSvcBackendManifest,
			tcpRouteManifest: singleSvcTcpRouteManifest,
			proxyService:     singleSvcProxyService,
			proxyDeployment:  singleSvcProxyDeployment,
			expectedResponses: []*matchers.HttpResponse{
				expectedSingleSvcResp,
			},
			ports: []int{8087},
			listenerNames: []v1.SectionName{
				v1.SectionName(singleSvcListenerName8087),
			},
			expectedRouteCounts: []int32{1},
			tcpRouteNames:       []string{singleSvcTCPRouteName},
		},
		{
			name:             "MultiServicesTCPRoute",
			nsManifest:       multiSvcNsManifest,
			gtwName:          multiSvcGatewayName,
			gtwNs:            multiSvcNsName,
			gtwManifest:      multiSvcGatewayAndClientManifest,
			svcManifest:      multiSvcBackendManifest,
			tcpRouteManifest: multiSvcTcpRouteManifest,
			proxyService:     multiProxyService,
			proxyDeployment:  multiProxyDeployment,
			expectedResponses: []*matchers.HttpResponse{
				expectedMultiSvc1Resp,
				expectedMultiSvc2Resp,
			},
			ports: []int{8088, 8089},
			listenerNames: []v1.SectionName{
				v1.SectionName(multiSvcListenerName8088),
				v1.SectionName(multiSvcListenerName8089),
			},
			expectedRouteCounts: []int32{1, 1},
			tcpRouteNames:       []string{multiSvcTCPRouteName1, multiSvcTCPRouteName2},
		},
	}

	for _, tc := range testCases {
		tc := tc // capture range variable
		s.Run(tc.name, func() {
			// Cleanup function
			s.T().Cleanup(func() {
				err := s.deleteManifests(tc.nsManifest)
				s.Require().NoError(err, fmt.Sprintf("Failed to delete manifest %s", tc.nsManifest))

				s.testInstallation.Assertions.EventuallyObjectsNotExist(s.ctx, &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: tc.gtwNs}})
			})

			// Setup environment
			s.setupTestEnvironment(
				tc.nsManifest,
				tc.gtwName,
				tc.gtwNs,
				tc.gtwManifest,
				tc.svcManifest,
				tc.proxyService,
				tc.proxyDeployment,
			)

			// Apply TCPRoute manifest
			err := s.testInstallation.Actions.Kubectl().ApplyFile(s.ctx, tc.tcpRouteManifest)
			s.Require().NoError(err, fmt.Sprintf("Failed to apply manifest %s", tc.tcpRouteManifest))

			// Assert gateway conditions
			s.testInstallation.Assertions.EventuallyGatewayCondition(s.ctx, tc.gtwName, tc.gtwNs, v1.GatewayConditionAccepted, metav1.ConditionTrue, timeout)

			// Assert TCPRoute conditions
			for _, tcpRouteName := range tc.tcpRouteNames {
				s.testInstallation.Assertions.EventuallyTCPRouteCondition(s.ctx, tcpRouteName, tc.gtwNs, v1.RouteConditionAccepted, metav1.ConditionTrue, timeout)
				s.testInstallation.Assertions.EventuallyTCPRouteCondition(s.ctx, tcpRouteName, tc.gtwNs, v1.RouteConditionResolvedRefs, metav1.ConditionTrue, timeout)
			}

			// Assert gateway programmed condition
			s.testInstallation.Assertions.EventuallyGatewayCondition(s.ctx, tc.gtwName, tc.gtwNs, v1.GatewayConditionProgrammed, metav1.ConditionTrue, timeout)

			// Assert listener attached routes
			for i, listenerName := range tc.listenerNames {
				expectedRouteCount := tc.expectedRouteCounts[i]
				s.testInstallation.Assertions.EventuallyGatewayListenerAttachedRoutes(s.ctx, tc.gtwName, tc.gtwNs, listenerName, expectedRouteCount, timeout)
			}

			// Assert expected responses
			for i, port := range tc.ports {
				expectedResponse := tc.expectedResponses[i]
				s.testInstallation.Assertions.AssertEventualCurlResponse(
					s.ctx,
					s.execOpts(tc.gtwNs),
					[]curl.Option{
						curl.WithHost(kubeutils.ServiceFQDN(tc.proxyService.ObjectMeta)),
						curl.WithPort(port),
					},
					expectedResponse)
			}
		})
	}
}

func validateManifestFile(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("Manifest file not found: %s", path)
	}
	return nil
}

func (s *testingSuite) setupTestEnvironment(nsManifest, gtwName, gtwNs, gtwManifest, svcManifest string, proxySvc *corev1.Service, proxyDeploy *appsv1.Deployment) {
	s.applyManifests(nsManifest)

	s.applyManifests(gtwManifest)
	s.testInstallation.Assertions.EventuallyGatewayCondition(s.ctx, gtwName, gtwNs, v1.GatewayConditionAccepted, metav1.ConditionTrue, timeout)

	s.applyManifests(svcManifest)
	s.testInstallation.Assertions.EventuallyObjectsExist(s.ctx, proxySvc, proxyDeploy)
}

func (s *testingSuite) applyManifests(manifests ...string) {
	for _, manifest := range manifests {
		s.Eventually(func() bool {
			err := s.testInstallation.Actions.Kubectl().ApplyFile(s.ctx, manifest)
			if err != nil {
				s.T().Logf("Retrying apply manifest: %s, error: %v", manifest, err)
				return false
			}
			return true
		}, waitTime, tickTime, fmt.Sprintf("Can apply manifest %s", manifest))
	}
}

func (s *testingSuite) deleteManifests(manifests ...string) error {
	for _, manifest := range manifests {
		if err := s.testInstallation.Actions.Kubectl().DeleteFile(s.ctx, manifest); err != nil {
			return err
		}
	}
	return nil
}

func (s *testingSuite) execOpts(ns string) kubectl.PodExecOptions {
	opts := defaults.CurlPodExecOpt
	opts.Namespace = ns
	return opts
}
