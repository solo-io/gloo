package route_delegation

import (
	"context"
	"fmt"
	"net/http"
	"time"

	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/suite"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	gwv1 "sigs.k8s.io/gateway-api/apis/v1"

	"github.com/solo-io/gloo/pkg/utils/requestutils/curl"
	testmatchers "github.com/solo-io/gloo/test/gomega/matchers"
	"github.com/solo-io/gloo/test/kubernetes/e2e"
	"github.com/solo-io/gloo/test/kubernetes/e2e/defaults"
)

var _ e2e.NewSuiteFunc = NewTestingSuite

type tsuite struct {
	suite.Suite

	ctx context.Context
	// ti contains all the metadata/utilities necessary to execute a series of tests
	// against an installation of Gloo Gateway
	ti *e2e.TestInstallation

	// maps test name to a list of manifests to apply before the test
	manifests map[string][]string

	manifestObjects map[string][]client.Object

	// resources from common manifest
	commonResources []client.Object
}

func NewTestingSuite(ctx context.Context, testInst *e2e.TestInstallation) suite.TestingSuite {
	return &tsuite{
		ctx: ctx,
		ti:  testInst,
	}
}

func (s *tsuite) SetupSuite() {
	s.manifests = map[string][]string{
		"TestBasic":                       {basicRoutesManifest},
		"TestRecursive":                   {recursiveRoutesManifest},
		"TestCyclic":                      {cyclicRoutesManifest},
		"TestInvalidChild":                {invalidChildRoutesManifest},
		"TestHeaderQueryMatch":            {headerQueryMatchRoutesManifest},
		"TestMultipleParents":             {multipleParentsManifest},
		"TestInvalidChildValidStandalone": {invalidChildValidStandaloneManifest},
		"TestUnresolvedChild":             {unresolvedChildManifest},
		"TestRouteOptions":                {routeOptionsManifest},
		"TestMatcherInheritance":          {matcherInheritanceManifest},
		"TestLabelSelector":               {labelSelectorManifest},
	}
	// Not every resource that is applied needs to be verified. We are not testing `kubectl apply`,
	// but the below code demonstrates how it can be done if necessary
	s.manifestObjects = map[string][]client.Object{
		basicRoutesManifest:                 {routeRoot, routeTeam1, routeTeam2},
		cyclicRoutesManifest:                {routeRoot, routeTeam1, routeTeam2},
		recursiveRoutesManifest:             {routeRoot, routeTeam1, routeTeam2},
		invalidChildRoutesManifest:          {routeRoot, routeTeam1, routeTeam2},
		headerQueryMatchRoutesManifest:      {routeRoot, routeTeam1, routeTeam2},
		multipleParentsManifest:             {routeParent1, routeParent2, routeTeam1, routeTeam2},
		invalidChildValidStandaloneManifest: {proxyTestService, proxyTestDeployment, routeRoot, routeTeam1, routeTeam2},
		unresolvedChildManifest:             {routeRoot},
		routeOptionsManifest:                {routeRoot, routeTeam1, routeTeam2},
		matcherInheritanceManifest:          {routeParent1, routeParent2, routeTeam1},
		labelSelectorManifest:               {routeRoot, routeTeam1, routeTeam2},
	}

	s.commonResources = []client.Object{
		// resources from manifest
		httpbinTeam1Service, httpbinTeam1Deployment, httpbinTeam2Service, httpbinTeam2Deployment, gateway,
		// deployer-generated resources
		proxyDeployment, proxyService,
	}

	// set up common resources once
	err := s.ti.Actions.Kubectl().ApplyFile(s.ctx, commonManifest)
	s.Require().NoError(err, "can apply common manifest")

	a := s.ti.AssertionsT(s.T())

	a.EventuallyObjectsExist(s.ctx, s.commonResources...)
	// make sure pods are running
	a.EventuallyPodsRunning(s.ctx, httpbinTeam1Deployment.GetNamespace(), metav1.ListOptions{
		LabelSelector: "app=httpbin,version=v1",
	})
	a.EventuallyPodsRunning(s.ctx, httpbinTeam2Deployment.GetNamespace(), metav1.ListOptions{
		LabelSelector: "app=httpbin,version=v2",
	})
	a.EventuallyPodsRunning(s.ctx, proxyMeta.GetNamespace(), metav1.ListOptions{
		LabelSelector: fmt.Sprintf("app.kubernetes.io/name=%s", proxyMeta.GetName()),
	})

	// set up curl once
	err = s.ti.Actions.Kubectl().ApplyFile(s.ctx, defaults.CurlPodManifest)
	s.Require().NoError(err, "can apply curl pod manifest")
	a.EventuallyPodsRunning(s.ctx, defaults.CurlPod.GetNamespace(), metav1.ListOptions{
		LabelSelector: defaults.CurlPodLabelSelector,
	})
}

func (s *tsuite) TearDownSuite() {
	a := s.ti.AssertionsT(s.T())

	// clean up curl
	err := s.ti.Actions.Kubectl().DeleteFileSafe(s.ctx, defaults.CurlPodManifest)
	s.Require().NoError(err, "can delete curl pod manifest")
	a.EventuallyObjectsNotExist(s.ctx, defaults.CurlPod)

	// clean up common resources
	err = s.ti.Actions.Kubectl().DeleteFileSafe(s.ctx, commonManifest)
	s.Require().NoError(err, "can delete common manifest")
	a.EventuallyObjectsNotExist(s.ctx, s.commonResources...)
	a.EventuallyPodsNotExist(s.ctx, httpbinTeam1Deployment.GetNamespace(), metav1.ListOptions{
		LabelSelector: "app=httpbin,version=v1",
	})
	a.EventuallyPodsNotExist(s.ctx, httpbinTeam2Deployment.GetNamespace(), metav1.ListOptions{
		LabelSelector: "app=httpbin,version=v2",
	})
	a.EventuallyPodsNotExist(s.ctx, proxyMeta.GetNamespace(), metav1.ListOptions{
		LabelSelector: fmt.Sprintf("app.kubernetes.io/name=%s", proxyMeta.GetName()),
	})
}

func (s *tsuite) BeforeTest(suiteName, testName string) {
	a := s.ti.AssertionsT(s.T())

	manifests := s.manifests[testName]
	for _, manifest := range manifests {
		err := s.ti.Actions.Kubectl().ApplyFile(s.ctx, manifest)
		s.Require().NoError(err)
		a.EventuallyObjectsExist(s.ctx, s.manifestObjects[manifest]...)
	}
}

func (s *tsuite) AfterTest(suiteName, testName string) {
	a := s.ti.AssertionsT(s.T())
	manifests := s.manifests[testName]
	for _, manifest := range manifests {
		err := s.ti.Actions.Kubectl().DeleteFileSafe(s.ctx, manifest)
		s.Require().NoError(err)
		a.EventuallyObjectsNotExist(s.ctx, s.manifestObjects[manifest]...)
	}
}

func (s *tsuite) TestBasic() {
	a := s.ti.AssertionsT(s.T())
	// Assert traffic to team1 route
	a.AssertEventuallyConsistentCurlResponse(s.ctx, defaults.CurlPodExecOpt, []curl.Option{curl.WithHostPort(proxyHostPort), curl.WithPath(pathTeam1)},
		&testmatchers.HttpResponse{StatusCode: http.StatusOK, Body: ContainSubstring(pathTeam1)})

	// Assert traffic to team2 route
	a.AssertEventuallyConsistentCurlResponse(s.ctx, defaults.CurlPodExecOpt, []curl.Option{curl.WithHostPort(proxyHostPort), curl.WithPath(pathTeam2)},
		&testmatchers.HttpResponse{StatusCode: http.StatusOK, Body: ContainSubstring(pathTeam2)})
}

func (s *tsuite) TestRecursive() {
	a := s.ti.AssertionsT(s.T())
	// Assert traffic to team1 route
	a.AssertEventuallyConsistentCurlResponse(s.ctx, defaults.CurlPodExecOpt, []curl.Option{curl.WithHostPort(proxyHostPort), curl.WithPath(pathTeam1)},
		&testmatchers.HttpResponse{StatusCode: http.StatusOK, Body: ContainSubstring(pathTeam1)})

	// Assert traffic to team2 route
	a.AssertEventuallyConsistentCurlResponse(s.ctx, defaults.CurlPodExecOpt, []curl.Option{curl.WithHostPort(proxyHostPort), curl.WithPath(pathTeam2)},
		&testmatchers.HttpResponse{StatusCode: http.StatusOK, Body: ContainSubstring(pathTeam2)})
}

func (s *tsuite) TestCyclic() {
	a := s.ti.AssertionsT(s.T())
	// Assert traffic to team1 route
	a.AssertEventuallyConsistentCurlResponse(s.ctx, defaults.CurlPodExecOpt, []curl.Option{curl.WithHostPort(proxyHostPort), curl.WithPath(pathTeam1)},
		&testmatchers.HttpResponse{StatusCode: http.StatusOK, Body: ContainSubstring(pathTeam1)})

	// Assert traffic to team2 route fails with HTTP 404 as it is a cyclic route
	a.AssertEventuallyConsistentCurlResponse(s.ctx, defaults.CurlPodExecOpt, []curl.Option{curl.WithHostPort(proxyHostPort), curl.WithPath(pathTeam2)},
		&testmatchers.HttpResponse{StatusCode: http.StatusNotFound})

	a.EventuallyHTTPRouteStatusContainsMessage(s.ctx, routeTeam2.Name, routeTeam2.Namespace,
		"cyclic reference detected", 10*time.Second, 1*time.Second)
}

func (s *tsuite) TestInvalidChild() {
	a := s.ti.AssertionsT(s.T())
	// Assert traffic to team1 route
	a.AssertEventuallyConsistentCurlResponse(s.ctx, defaults.CurlPodExecOpt, []curl.Option{curl.WithHostPort(proxyHostPort), curl.WithPath(pathTeam1)},
		&testmatchers.HttpResponse{StatusCode: http.StatusOK, Body: ContainSubstring(pathTeam1)})

	// Assert traffic to team2 route fails with HTTP 404 as the route is invalid due to specifying a hostname on the child route
	a.AssertEventuallyConsistentCurlResponse(s.ctx, defaults.CurlPodExecOpt, []curl.Option{curl.WithHostPort(proxyHostPort), curl.WithPath(pathTeam2)},
		&testmatchers.HttpResponse{StatusCode: http.StatusNotFound})

	a.EventuallyHTTPRouteStatusContainsMessage(s.ctx, routeTeam2.Name, routeTeam2.Namespace,
		"spec.hostnames must be unset", 10*time.Second, 1*time.Second)
}

func (s *tsuite) TestHeaderQueryMatch() {
	a := s.ti.AssertionsT(s.T())
	// Assert traffic to team1 route with matching header and query parameters
	a.AssertEventuallyConsistentCurlResponse(s.ctx, defaults.CurlPodExecOpt,
		[]curl.Option{
			curl.WithHostPort(proxyHostPort), curl.WithPath(pathTeam1),
			curl.WithHeader("header1", "val1"),
			curl.WithHeader("headerX", "valX"),
			curl.WithQueryParameters(map[string]string{"query1": "val1", "queryX": "valX"}),
		},
		&testmatchers.HttpResponse{StatusCode: http.StatusOK, Body: ContainSubstring(pathTeam1)})

	// Assert traffic to team2 route fails with HTTP 404 as it does not match the parent's header and query parameters
	a.AssertEventuallyConsistentCurlResponse(s.ctx, defaults.CurlPodExecOpt,
		[]curl.Option{
			curl.WithHostPort(proxyHostPort), curl.WithPath(pathTeam2),
			curl.WithHeader("headerX", "valX"),
			curl.WithQueryParameters(map[string]string{"queryX": "valX"}),
		},
		&testmatchers.HttpResponse{StatusCode: http.StatusNotFound})
}

func (s *tsuite) TestMultipleParents() {
	a := s.ti.AssertionsT(s.T())
	// Assert traffic to parent1.com/anything/team1
	a.AssertEventuallyConsistentCurlResponse(s.ctx, defaults.CurlPodExecOpt,
		[]curl.Option{
			curl.WithHostPort(proxyHostPort),
			curl.WithPath(pathTeam1),
			curl.WithHostHeader(routeParent1Host),
		},
		&testmatchers.HttpResponse{StatusCode: http.StatusOK, Body: ContainSubstring(pathTeam1)})

	// Assert traffic to parent1.com/anything/team2
	a.AssertEventuallyConsistentCurlResponse(s.ctx, defaults.CurlPodExecOpt,
		[]curl.Option{
			curl.WithHostPort(proxyHostPort),
			curl.WithPath(pathTeam2),
			curl.WithHostHeader(routeParent1Host),
		},
		&testmatchers.HttpResponse{StatusCode: http.StatusOK, Body: ContainSubstring(pathTeam2)})

	// Assert traffic to parent2.com/anything/team1
	a.AssertEventuallyConsistentCurlResponse(s.ctx, defaults.CurlPodExecOpt,
		[]curl.Option{
			curl.WithHostPort(proxyHostPort),
			curl.WithPath(pathTeam1),
			curl.WithHostHeader(routeParent2Host),
		},
		&testmatchers.HttpResponse{StatusCode: http.StatusOK, Body: ContainSubstring(pathTeam1)})

	// Assert traffic to parent2.com/anything/team2 fails as it is not selected by parent2 route
	a.AssertEventuallyConsistentCurlResponse(s.ctx, defaults.CurlPodExecOpt,
		[]curl.Option{
			curl.WithHostPort(proxyHostPort),
			curl.WithPath(pathTeam2),
			curl.WithHostHeader(routeParent2Host),
		},
		&testmatchers.HttpResponse{StatusCode: http.StatusNotFound})
}

func (s *tsuite) TestInvalidChildValidStandalone() {
	a := s.ti.AssertionsT(s.T())
	// Assert traffic to team1 route
	a.AssertEventuallyConsistentCurlResponse(s.ctx, defaults.CurlPodExecOpt,
		[]curl.Option{
			curl.WithHostPort(proxyTestHostPort),
			curl.WithPath(pathTeam1),
			curl.WithHostHeader(routeParentHost),
		},
		&testmatchers.HttpResponse{StatusCode: http.StatusOK, Body: ContainSubstring(pathTeam1)})

	// Assert traffic to team2 route on parent hostname fails with HTTP 404 as the route is invalid due to specifying a hostname on the child route
	a.AssertEventuallyConsistentCurlResponse(s.ctx, defaults.CurlPodExecOpt,
		[]curl.Option{
			curl.WithHostPort(proxyTestHostPort),
			curl.WithPath(pathTeam2),
			curl.WithHostHeader(routeParentHost),
		},
		&testmatchers.HttpResponse{StatusCode: http.StatusNotFound})

	// Assert traffic to team2 route on standalone host succeeds
	a.AssertEventuallyConsistentCurlResponse(s.ctx, defaults.CurlPodExecOpt,
		[]curl.Option{
			curl.WithHostPort(proxyTestHostPort),
			curl.WithPath(pathTeam2),
			curl.WithHostHeader(routeTeam2Host),
		},
		&testmatchers.HttpResponse{StatusCode: http.StatusOK, Body: ContainSubstring(pathTeam2)})

	a.EventuallyHTTPRouteStatusContainsMessage(s.ctx, routeTeam2.Name, routeTeam2.Namespace,
		"spec.hostnames must be unset", 10*time.Second, 1*time.Second)
}

func (s *tsuite) TestUnresolvedChild() {
	a := s.ti.AssertionsT(s.T())
	a.EventuallyHTTPRouteStatusContainsReason(s.ctx, routeRoot.Name, routeRoot.Namespace,
		string(gwv1.RouteReasonBackendNotFound), 10*time.Second, 1*time.Second)
}

func (s *tsuite) TestRouteOptions() {
	a := s.ti.AssertionsT(s.T())
	// Assert traffic to team1 route experiences the injected fault using RouteOption
	a.AssertEventuallyConsistentCurlResponse(s.ctx, defaults.CurlPodExecOpt,
		[]curl.Option{
			curl.WithHostPort(proxyHostPort),
			curl.WithPath(pathTeam1),
		},
		&testmatchers.HttpResponse{
			StatusCode: http.StatusTeapot,
			Body:       ContainSubstring("fault filter abort"),
		})

	// Assert traffic to team2 route succeeds with path rewrite using RouteOption
	// while also containing the response header set by the root RouteOption
	a.AssertEventuallyConsistentCurlResponse(s.ctx, defaults.CurlPodExecOpt,
		[]curl.Option{
			curl.WithHostPort(proxyHostPort),
			curl.WithPath(pathTeam2),
		},
		&testmatchers.HttpResponse{
			StatusCode: http.StatusOK,
			Body:       ContainSubstring("/anything/rewrite"),
			Headers:    map[string]interface{}{"x-foo": Equal("baz")},
		})
}

func (s *tsuite) TestMatcherInheritance() {
	a := s.ti.AssertionsT(s.T())
	// Assert traffic on parent1's prefix
	a.AssertEventuallyConsistentCurlResponse(s.ctx, defaults.CurlPodExecOpt,
		[]curl.Option{curl.WithHostPort(proxyHostPort), curl.WithPath("/anything/foo/child")},
		&testmatchers.HttpResponse{StatusCode: http.StatusOK, Body: ContainSubstring("/anything/foo/child")})

	// Assert traffic on parent2's prefix
	a.AssertEventuallyConsistentCurlResponse(s.ctx, defaults.CurlPodExecOpt,
		[]curl.Option{curl.WithHostPort(proxyHostPort), curl.WithPath("/anything/baz/child")},
		&testmatchers.HttpResponse{StatusCode: http.StatusOK, Body: ContainSubstring("/anything/baz/child")})
}

func (s *tsuite) TestLabelSelector() {
	a := s.ti.AssertionsT(s.T())

	// Assert traffic to team1 route
	a.AssertEventuallyConsistentCurlResponse(s.ctx, defaults.CurlPodExecOpt, []curl.Option{curl.WithHostPort(proxyHostPort), curl.WithPath(pathTeam1)},
		&testmatchers.HttpResponse{StatusCode: http.StatusOK, Body: ContainSubstring(pathTeam1)})

	// Assert traffic to team2 route
	a.AssertEventuallyConsistentCurlResponse(s.ctx, defaults.CurlPodExecOpt, []curl.Option{curl.WithHostPort(proxyHostPort), curl.WithPath(pathTeam2)},
		&testmatchers.HttpResponse{StatusCode: http.StatusOK, Body: ContainSubstring(pathTeam2)})
}
