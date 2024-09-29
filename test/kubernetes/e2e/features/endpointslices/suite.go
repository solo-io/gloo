package endpointslices

import (
	"context"
	"net/http"
	"path/filepath"
	"time"

	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/suite"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	discoveryv1 "k8s.io/api/discovery/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/solo-io/gloo/pkg/utils/kubeutils"
	"github.com/solo-io/gloo/pkg/utils/requestutils/curl"
	"github.com/solo-io/gloo/test/gomega/matchers"
	"github.com/solo-io/gloo/test/kubernetes/e2e"
	"github.com/solo-io/gloo/test/kubernetes/e2e/defaults"
	testdefaults "github.com/solo-io/gloo/test/kubernetes/e2e/defaults"
	"github.com/solo-io/skv2/codegen/util"
)

var (
	setupManifest               = filepath.Join(util.MustGetThisDir(), "testdata", "setup.yaml")
	gatewayManifest             = filepath.Join(util.MustGetThisDir(), "testdata", "gateway.yaml")
	manualEndpointSliceManifest = filepath.Join(util.MustGetThisDir(), "testdata", "manual-endpointslice.yaml")

	glooProxyObjectMeta = metav1.ObjectMeta{
		Name:      "gloo-proxy-gw",
		Namespace: "default",
	}
	proxyDeployment   = &appsv1.Deployment{ObjectMeta: glooProxyObjectMeta}
	proxyService      = &corev1.Service{ObjectMeta: glooProxyObjectMeta}
	httpbinDeployment = &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "httpbin",
			Namespace: "httpbin",
		},
	}
	epsMeta = &discoveryv1.EndpointSlice{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "httpbin-svc",
			Namespace: "httpbin",
		},
	}
)

type testingSuite struct {
	suite.Suite
	ctx context.Context
	ti  *e2e.TestInstallation
	// maps test name to a list of manifests to apply before the test
	manifests map[string][]string
}

func NewTestingSuite(
	ctx context.Context,
	testInst *e2e.TestInstallation,
) suite.TestingSuite {
	return &testingSuite{
		ctx: ctx,
		ti:  testInst,
	}
}

func (s *testingSuite) SetupSuite() {
	// Check that the common setup manifest is applied
	err := s.ti.Actions.Kubectl().ApplyFile(s.ctx, setupManifest)
	s.NoError(err, "can apply "+setupManifest)
	err = s.ti.Actions.Kubectl().ApplyFile(s.ctx, testdefaults.CurlPodManifest)
	s.NoError(err, "can apply curl pod manifest")

	// Check that istio injection is successful and httpbin is running
	s.ti.Assertions.EventuallyObjectsExist(s.ctx, httpbinDeployment)
	// httpbin can take a while to start up with Istio sidecar
	s.ti.Assertions.EventuallyPodsRunning(s.ctx, httpbinDeployment.GetNamespace(), metav1.ListOptions{
		LabelSelector: "app=httpbin",
	})
	s.ti.Assertions.EventuallyPodsRunning(s.ctx, testdefaults.CurlPod.GetNamespace(), metav1.ListOptions{
		LabelSelector: "app.kubernetes.io/name=curl",
	})

	// TODO: wrap in eventually statement since Update can/will fail.
	pods := &corev1.PodList{}
	err = s.ti.ClusterContext.Client.List(s.ctx, pods, []client.ListOption{
		client.InNamespace(httpbinDeployment.GetNamespace()),
		client.MatchingLabels{"app": "httpbin"},
	}...)
	s.NoError(err, "can list pods")
	s.Greater(len(pods.Items), 0, "httpbin pod exists")

	podIP := pods.Items[0].Status.PodIP

	// TODO: wrap in eventually statement since Update can/will fail.
	eps := &discoveryv1.EndpointSlice{}
	err = s.ti.ClusterContext.Client.Get(s.ctx, client.ObjectKeyFromObject(epsMeta), eps)
	s.NoError(err, "can get endpoint slice")

	eps.Endpoints = []discoveryv1.Endpoint{{
		Addresses: []string{podIP},
		Conditions: discoveryv1.EndpointConditions{
			Ready: ptr.To(true),
		},
	}}
	err = s.ti.ClusterContext.Client.Update(s.ctx, eps)
	s.NoError(err, "can update endpoint slice")

	// TODO: add test case that modifies the Service in some way to verify everything is still gucci.
	// TODO: targetRef for Pod?

	// include gateway manifests for the tests, so we recreate it for each test run
	s.manifests = map[string][]string{
		"TestBasicEndpointSlice": {gatewayManifest, manualEndpointSliceManifest},
	}
}

func (s *testingSuite) TearDownSuite() {
	err := s.ti.Actions.Kubectl().DeleteFileSafe(s.ctx, setupManifest)
	s.NoError(err, "can delete setup manifest")
	err = s.ti.Actions.Kubectl().DeleteFileSafe(s.ctx, testdefaults.CurlPodManifest)
	s.NoError(err, "can delete curl pod manifest")
	s.ti.Assertions.EventuallyObjectsNotExist(s.ctx, httpbinDeployment)
}

func (s *testingSuite) BeforeTest(suiteName, testName string) {
	manifests, ok := s.manifests[testName]
	if !ok {
		s.FailNow("no manifests found for %s, manifest map contents: %v", testName, s.manifests)
	}
	for _, manifest := range manifests {
		err := s.ti.Actions.Kubectl().ApplyFile(s.ctx, manifest)
		s.Assert().NoError(err, "can apply manifest "+manifest)
	}

	// we recreate the `Gateway` resource (and thus dynamically provision the proxy pod) for each test run
	// so let's assert the proxy svc and pod is ready before moving on
	s.ti.Assertions.EventuallyObjectsExist(s.ctx, proxyService, proxyDeployment)
	s.ti.Assertions.EventuallyPodsRunning(s.ctx, proxyDeployment.ObjectMeta.GetNamespace(), metav1.ListOptions{
		LabelSelector: "app.kubernetes.io/name=gloo-proxy-gw",
	})
}

func (s *testingSuite) AfterTest(suiteName, testName string) {
	manifests, ok := s.manifests[testName]
	if !ok {
		s.FailNow("no manifests found for " + testName)
	}

	for _, manifest := range manifests {
		output, err := s.ti.Actions.Kubectl().DeleteFileWithOutput(s.ctx, manifest)
		s.ti.Assertions.ExpectObjectDeleted(manifest, err, output)
	}
}

func (s *testingSuite) TestBasicEndpointSlice() {
	// verify that a direct response route works as expected
	s.ti.Assertions.AssertEventualCurlResponse(
		s.ctx,
		defaults.CurlPodExecOpt,
		[]curl.Option{
			curl.WithHost(kubeutils.ServiceFQDN(proxyService.ObjectMeta)),
			curl.WithHostHeader("www.example.com"),
			curl.WithPath("/manual-endpointslices"),
		},
		&matchers.HttpResponse{
			StatusCode: http.StatusOK,
			Body:       ContainSubstring("X-Envoy-Original-Path"),
		},
		time.Minute,
	)
}
