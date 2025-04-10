package listener_options

import (
	"context"
	"testing"
	"time"

	"github.com/onsi/gomega"
	"github.com/stretchr/testify/suite"

	"github.com/solo-io/gloo/pkg/utils/envoyutils/admincli"
	"github.com/solo-io/gloo/pkg/utils/requestutils/curl"
	"github.com/solo-io/gloo/test/kubernetes/e2e"
	testdefaults "github.com/solo-io/gloo/test/kubernetes/e2e/defaults"
	"github.com/solo-io/gloo/test/kubernetes/e2e/features/listenerset"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ e2e.NewSuiteFunc = NewTestingSuite

// testingSuite is the entire Suite of tests for the "ListenerOptions" feature
type testingSuite struct {
	suite.Suite
	ctx              context.Context
	testInstallation *e2e.TestInstallation
	// maps test name to a list of manifests to apply before the test
	manifests map[string][]string
}

func NewTestingSuite(
	ctx context.Context,
	testInst *e2e.TestInstallation,
) suite.TestingSuite {
	return &testingSuite{
		ctx:              ctx,
		testInstallation: testInst,
	}
}

func (s *testingSuite) SetupSuite() {
	// Check that the common setup manifest is applied
	for _, manifest := range setupManifests {
		err := s.testInstallation.Actions.Kubectl().ApplyFile(s.ctx, manifest)
		s.NoError(err, "can apply "+manifest)
	}

	if listenerset.RequiredCrdExists(s.testInstallation) {
		err := s.testInstallation.Actions.Kubectl().ApplyFile(s.ctx, listenerSetManifest)
		s.NoError(err, "can apply "+listenerSetManifest)
	}

	s.testInstallation.AssertionsT(s.T()).EventuallyPodsRunning(s.ctx, nginxPod.ObjectMeta.GetNamespace(), metav1.ListOptions{
		LabelSelector: "app.kubernetes.io/name=nginx",
	})
	s.testInstallation.AssertionsT(s.T()).EventuallyPodsRunning(s.ctx, testdefaults.CurlPod.GetNamespace(), metav1.ListOptions{
		LabelSelector: "app.kubernetes.io/name=curl",
	})
	s.testInstallation.AssertionsT(s.T()).EventuallyPodsRunning(s.ctx, proxy1Deployment.ObjectMeta.GetNamespace(), metav1.ListOptions{
		LabelSelector: "app.kubernetes.io/name=gloo-proxy-gw-1",
	})
	s.testInstallation.AssertionsT(s.T()).EventuallyPodsRunning(s.ctx, proxy2Deployment.ObjectMeta.GetNamespace(), metav1.ListOptions{
		LabelSelector: "app.kubernetes.io/name=gloo-proxy-gw-2",
	})

	s.manifests = map[string][]string{
		"TestConfigureListenerOptions":                        {basicLisOptManifest},
		"TestConfigureListenerOptionsWithSectionedTargetRefs": {basicLisOptManifest, lisOptWithSectionedTargetRefsManifest, lisOptWithListenerSetRefsManifest},
	}
}

func (s *testingSuite) TearDownSuite() {
	// Check that the common setup manifest is deleted
	for _, manifest := range setupManifests {
		output, err := s.testInstallation.Actions.Kubectl().DeleteFileWithOutput(s.ctx, manifest)
		s.NoError(err, "can delete "+manifest)
		s.testInstallation.AssertionsT(s.T()).ExpectObjectDeleted(manifest, err, output)
	}

	if listenerset.RequiredCrdExists(s.testInstallation) {
		output, err := s.testInstallation.Actions.Kubectl().DeleteFileWithOutput(s.ctx, listenerSetManifest)
		s.NoError(err, "can delete "+listenerSetManifest)
		s.testInstallation.AssertionsT(s.T()).ExpectObjectDeleted(listenerSetManifest, err, output)
	}
}

func (s *testingSuite) BeforeTest(suiteName, testName string) {
	manifests, ok := s.manifests[testName]
	if !ok {
		s.FailNow("no manifests found for %s, manifest map contents: %v", testName, s.manifests)
	}

	for _, manifest := range manifests {
		err := s.testInstallation.Actions.Kubectl().ApplyFile(s.ctx, manifest)
		s.Assert().NoError(err, "can apply manifest "+manifest)
	}
}

func (s *testingSuite) AfterTest(suiteName, testName string) {
	manifests, ok := s.manifests[testName]
	if !ok {
		s.FailNow("no manifests found for " + testName)
	}

	for _, manifest := range manifests {
		output, err := s.testInstallation.Actions.Kubectl().DeleteFileWithOutput(s.ctx, manifest)
		s.testInstallation.AssertionsT(s.T()).ExpectObjectDeleted(manifest, err, output)
	}
}

func (s *testingSuite) TestConfigureListenerOptions() {
	// Check healthy response
	s.testInstallation.AssertionsT(s.T()).AssertEventualCurlResponse(
		s.ctx,
		testdefaults.CurlPodExecOpt,
		[]curl.Option{
			curl.WithHost(proxy1ServiceFqdn),
			curl.WithHostHeader("example.com"),
		},
		expectedHealthyResponse)

	// Check the buffer limit is set on the Listener via Envoy config dump
	s.testInstallation.AssertionsT(s.T()).AssertEnvoyAdminApi(
		s.ctx,
		proxy1Deployment.ObjectMeta,
		listenerBufferLimitAssertion(s.testInstallation, s.T()),
	)
}

func (s *testingSuite) TestConfigureListenerOptionsWithSectionedTargetRefs() {
	type bufferLimitForListener struct {
		sectionName string
		port        int
		limit       int
	}

	// Setup the expected buffer limits for each listener
	bufferLimitsForListeners := map[string][]*bufferLimitForListener{
		proxy1ServiceFqdn: {
			{sectionName: "http", port: gw1port1, limit: 32000},
			{sectionName: "other", port: gw1port2, limit: 42000},
		},
		proxy2ServiceFqdn: {
			{sectionName: "http", port: gw2port1, limit: 0},
			{sectionName: "other", port: gw2port2, limit: 32000},
		},
	}

	if listenerset.RequiredCrdExists(s.testInstallation) {
		bufferLimitsForListeners[proxy1ServiceFqdn] = append(bufferLimitsForListeners[proxy1ServiceFqdn], &bufferLimitForListener{sectionName: "default/gw-1/listener-1", port: ls1port1, limit: 42000})
		bufferLimitsForListeners[proxy1ServiceFqdn] = append(bufferLimitsForListeners[proxy1ServiceFqdn], &bufferLimitForListener{sectionName: "default/gw-1/listener-2", port: ls1port2, limit: 21000})
	}

	objectMetaForListener := map[string]metav1.ObjectMeta{
		proxy1ServiceFqdn: proxy1Deployment.ObjectMeta,
		proxy2ServiceFqdn: proxy2Deployment.ObjectMeta,
	}

	// Curl each listener for which a matcher is defined
	for host, limits := range bufferLimitsForListeners {
		for _, limit := range limits {
			s.testInstallation.AssertionsT(s.T()).AssertEventualCurlResponse(
				s.ctx,
				testdefaults.CurlPodExecOpt,
				[]curl.Option{
					curl.WithHost(host),
					curl.WithHostHeader("example.com"),
					curl.WithPort(limit.port),
				},
				expectedHealthyResponse,
			)

			// Check the buffer limit is set on the Listener via Envoy config dump
			s.testInstallation.AssertionsT(s.T()).AssertEnvoyAdminApi(
				s.ctx,
				objectMetaForListener[host],
				listenerBufferLimitAssertionForSection(s.testInstallation, s.T(), limit.sectionName, limit.limit),
			)
		}
	}
}

func listenerBufferLimitAssertion(testInstallation *e2e.TestInstallation, t *testing.T) func(ctx context.Context, adminClient *admincli.Client) {
	return func(ctx context.Context, adminClient *admincli.Client) {
		testInstallation.AssertionsT(t).Gomega.Eventually(func(g gomega.Gomega) {
			listener, err := adminClient.GetSingleListenerFromDynamicListeners(ctx, "http")
			g.Expect(err).NotTo(gomega.HaveOccurred(), "error getting listener")
			g.Expect(listener.GetPerConnectionBufferLimitBytes().GetValue()).To(gomega.BeEquivalentTo(42000))
		}).
			WithContext(ctx).
			WithTimeout(time.Second * 10).
			WithPolling(time.Millisecond * 200).
			Should(gomega.Succeed())
	}
}

func listenerBufferLimitAssertionForSection(testInstallation *e2e.TestInstallation, t *testing.T, sectionName string, expectedValue int) func(ctx context.Context, adminClient *admincli.Client) {
	return func(ctx context.Context, adminClient *admincli.Client) {
		testInstallation.AssertionsT(t).Gomega.Eventually(func(g gomega.Gomega) {
			listener, err := adminClient.GetSingleListenerFromDynamicListeners(ctx, sectionName)
			g.Expect(err).NotTo(gomega.HaveOccurred(), "error getting listener")
			g.Expect(listener.GetPerConnectionBufferLimitBytes().GetValue()).To(gomega.BeEquivalentTo(expectedValue))
		}).
			WithContext(ctx).
			WithTimeout(time.Second * 10).
			WithPolling(time.Millisecond * 200).
			Should(gomega.Succeed())
	}
}
