package virtualhost_options

import (
	"context"
	"net/http"

	"github.com/onsi/gomega"
	"github.com/onsi/gomega/gstruct"
	"github.com/stretchr/testify/suite"

	"github.com/solo-io/gloo/pkg/utils/kubeutils"
	"github.com/solo-io/gloo/pkg/utils/requestutils/curl"
	"github.com/solo-io/gloo/projects/gloo/pkg/defaults"
	"github.com/solo-io/gloo/test/gomega/matchers"
	"github.com/solo-io/gloo/test/helpers"
	"github.com/solo-io/gloo/test/kubernetes/e2e"
	testdefaults "github.com/solo-io/gloo/test/kubernetes/e2e/defaults"
	"github.com/solo-io/gloo/test/kubernetes/e2e/features/listenerset"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ e2e.NewSuiteFunc = NewTestingSuite

// testingSuite is the entire Suite of tests for the "VirtualHostOptions" feature
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
	// Check that the common setup manifest is applied
	for _, manifest := range setupManifests {
		err := s.testInstallation.Actions.Kubectl().ApplyFile(s.ctx, manifest)
		s.NoError(err, "can apply "+manifest)
	}

	if listenerset.RequiredCrdExists(s.testInstallation) {
		err := s.testInstallation.Actions.Kubectl().ApplyFile(s.ctx, manifestListenerSetup)
		s.NoError(err, "can apply "+manifestListenerSetup)
	}

	s.testInstallation.AssertionsT(s.T()).EventuallyObjectsExist(s.ctx, proxyService1,
		proxyDeployment1, proxyService2, proxyDeployment2, exampleSvc, nginxPod, testdefaults.CurlPod)

	// Check that test resources are running
	s.testInstallation.AssertionsT(s.T()).EventuallyPodsRunning(s.ctx, nginxPod.ObjectMeta.GetNamespace(),
		metav1.ListOptions{
			LabelSelector: "app.kubernetes.io/name=nginx",
		})
	s.testInstallation.AssertionsT(s.T()).EventuallyPodsRunning(s.ctx, proxyDeployment1.ObjectMeta.GetNamespace(),
		metav1.ListOptions{
			LabelSelector: "app.kubernetes.io/name=gloo-proxy-gw-1",
		})
	s.testInstallation.AssertionsT(s.T()).EventuallyPodsRunning(s.ctx, proxyDeployment2.ObjectMeta.GetNamespace(),
		metav1.ListOptions{
			LabelSelector: "app.kubernetes.io/name=gloo-proxy-gw-2",
		})
	s.testInstallation.AssertionsT(s.T()).EventuallyPodsRunning(s.ctx, testdefaults.CurlPod.GetNamespace(),
		metav1.ListOptions{
			LabelSelector: "app.kubernetes.io/name=curl",
		})
}

func (s *testingSuite) TearDownSuite() {
	// Check that the common setup manifest is deleted
	for _, manifest := range setupManifests {
		output, err := s.testInstallation.Actions.Kubectl().DeleteFileWithOutput(s.ctx, manifest)
		s.NoError(err, "can delete "+manifest)
		s.testInstallation.AssertionsT(s.T()).ExpectObjectDeleted(manifest, err, output)
	}

	if listenerset.RequiredCrdExists(s.testInstallation) {
		output, err := s.testInstallation.Actions.Kubectl().DeleteFileWithOutput(s.ctx, manifestListenerSetup)
		s.NoError(err, "can delete "+manifestListenerSetup)
		s.testInstallation.AssertionsT(s.T()).ExpectObjectDeleted(manifestListenerSetup, err, output)
	}
}

// TestConfirmSetup tests that the setup is correct
//
// The default state should have two gateways, each with two listeners.
// The first gateway, gw-1, on ports 8080 and 8081, and the headers x-bar and x-baz should be added to the response
// A ListenerSet which has listeners on ports 8085 and 8086, is attched to gw-1 and should have the same headers as gw-1
// The second gateway, gw-2, on ports 8083 and 8084, and the headers x-bar-2 and x-baz-2 should be added to the response
func (s *testingSuite) TestConfirmSetup() {

	matchersForListeners := map[string]map[int]*matchers.HttpResponse{
		proxyService1Fqdn: {
			gw1port1: defaultResponseGw1,
			gw1port2: defaultResponseGw1,
		},
		proxyService2Fqdn: {
			gw2port1: defaultResponseGw2,
			gw2port2: defaultResponseGw2,
		},
	}

	if listenerset.RequiredCrdExists(s.testInstallation) {
		matchersForListeners[proxyService1Fqdn][lsPort1] = defaultResponseGw1
		matchersForListeners[proxyService1Fqdn][lsPort2] = defaultResponseGw1
	}

	s.testExpectedResponsesForManifests(nil, matchersForListeners, true)
}

// TestConfigureVirtualHostOptions tests the basic functionality of VirtualHostOptions using a single VHO
func (s *testingSuite) TestConfigureVirtualHostOptions() {
	manifests := map[string]*metav1.ObjectMeta{
		manifestVhoRemoveXBar: &vhoRemoveXBar,
	}

	matchersForListeners := map[string]map[int]*matchers.HttpResponse{
		proxyService1Fqdn: {
			gw1port1: expectedResponseWithoutXBar,
		},
	}

	s.testExpectedResponsesForManifests(manifests, matchersForListeners, true)
}

// TestConfigureVirtualHostOptions tests the basic functionality of VirtualHostOptions using a single VHO
// and multiple target refs. This test also indirectly validates targetRefs with sectionName.
func (s *testingSuite) TestConfigureVirtualHostOptionsMultipleTargetRefs() {

	manifests := map[string]*metav1.ObjectMeta{
		manifestVhoMultipleTargetRefs: &vhoMultipleTargetRefs,
	}

	matchersForListeners := map[string]map[int]*matchers.HttpResponse{
		proxyService1Fqdn: {
			gw1port1: expectedResponseWithXFoo("foo"),
			gw1port2: expectedResponseWithoutXFoo,
		},
		proxyService2Fqdn: {
			gw2port1: expectedResponseWithoutXFoo,
			gw2port2: expectedResponseWithXFoo("foo"),
		},
	}

	if listenerset.RequiredCrdExists(s.testInstallation) {
		matchersForListeners[proxyService1Fqdn][lsPort1] = expectedResponseWithoutXFoo
		matchersForListeners[proxyService1Fqdn][lsPort2] = expectedResponseWithoutXFoo
	}

	s.testExpectedResponsesForManifests(manifests, matchersForListeners, true)
}

// TestConfigureVirtualHostOptions tests the basic functionality of VirtualHostOptions using a single VHO
// and multiple target refs. This test also indirectly validates targetRefs with sectionName.
func (s *testingSuite) TestConfigureVirtualHostListenerSetTargetRef() {
	if !listenerset.RequiredCrdExists(s.testInstallation) {
		s.T().Skip("Skipping as the XListenerSet CRD is not installed")
	}

	manifests := map[string]*metav1.ObjectMeta{
		manifestVhoListenerSetTargetRef: &vhoListenerSetTargetRef,
	}

	matchersForListeners := map[string]map[int]*matchers.HttpResponse{
		proxyService1Fqdn: {
			gw1port1: expectedResponseWithoutXFoo,
			gw1port2: expectedResponseWithoutXFoo,
			lsPort1:  expectedResponseWithXFoo("foo-lis"),
			lsPort2:  expectedResponseWithXFoo("foo-lis"),
		},
		proxyService2Fqdn: {
			gw2port1: expectedResponseWithoutXFoo,
			gw2port2: expectedResponseWithoutXFoo,
		},
	}

	s.testExpectedResponsesForManifests(manifests, matchersForListeners, true)
}

// TestConfigureVirtualHostOptions tests the basic functionality of VirtualHostOptions using a single VHO
// and multiple target refs. This test also indirectly validates targetRefs with sectionName.
func (s *testingSuite) TestConfigureVirtualHostListenerSetSectionedTargetRef() {
	if !listenerset.RequiredCrdExists(s.testInstallation) {
		s.T().Skip("Skipping as the XListenerSet CRD is not installed")
	}

	manifests := map[string]*metav1.ObjectMeta{
		manifestVhoListenerSetSectionedTargetRef: &vhoListenerSetSectionedTargetRef,
	}

	// Setup the matchers for requests to the different listeners
	matchersForListeners := map[string]map[int]*matchers.HttpResponse{
		proxyService1Fqdn: {
			gw1port1: expectedResponseWithoutXFoo,
			gw1port2: expectedResponseWithoutXFoo,
			lsPort1:  expectedResponseWithXFoo("foo-lis-sec"),
			lsPort2:  expectedResponseWithoutXFoo,
		},
		proxyService2Fqdn: {
			gw2port1: expectedResponseWithoutXFoo,
			gw2port2: expectedResponseWithoutXFoo,
		},
	}
	s.testExpectedResponsesForManifests(manifests, matchersForListeners, true)
}

// This test should be updated to confirm statuses on conflicting VHOs once statuses are fixed
// this may involve updating testExpectedResponsesForManifests to allow either
func (s *testingSuite) TestConfigureVirtualHostOptionsWithConflictingVHO() {
	if !listenerset.RequiredCrdExists(s.testInstallation) {
		s.T().Skip("Skipping as the XListenerSet CRD is not installed")
	}

	manifests := map[string]*metav1.ObjectMeta{
		manifestVhoSectionAddXFoo:                &vhoSectionAddXFoo,
		manifestVhoGwAddXFoo:                     &vhoGwAddXFoo,
		manifestVhoListenerSetTargetRef:          &vhoListenerSetTargetRef,
		manifestVhoListenerSetSectionedTargetRef: &vhoListenerSetSectionedTargetRef,
	}

	matchersForListeners := map[string]map[int]*matchers.HttpResponse{
		proxyService1Fqdn: {
			gw1port1: expectedResponseWithXFoo("foo-http"),
			gw1port2: expectedResponseWithXFoo("foo"),
			lsPort1:  expectedResponseWithXFoo("foo-lis-sec"),
			lsPort2:  expectedResponseWithXFoo("foo-lis"),
		},
		proxyService2Fqdn: {
			gw2port1: expectedResponseWithoutXFoo,
			gw2port2: expectedResponseWithoutXFoo,
		},
	}

	// Skip validation for this test because the status is not accepted
	s.testExpectedResponsesForManifests(manifests, matchersForListeners, false)
}

// TestConfigureInvalidVirtualHostOptions confirms that an invalid VirtualHostOption is rejected
func (s *testingSuite) TestConfigureInvalidVirtualHostOptions() {
	s.T().Cleanup(func() {
		output, err := s.testInstallation.Actions.Kubectl().DeleteFileWithOutput(s.ctx, manifestVhoRemoveXBar)
		s.testInstallation.AssertionsT(s.T()).ExpectObjectDeleted(manifestVhoRemoveXBar, err, output)

		output, err = s.testInstallation.Actions.Kubectl().DeleteFileWithOutput(s.ctx, manifestVhoWebhookReject)
		s.testInstallation.AssertionsT(s.T()).ExpectObjectDeleted(manifestVhoWebhookReject, err, output)
	})

	err := s.testInstallation.Actions.Kubectl().ApplyFile(s.ctx, manifestVhoRemoveXBar)
	s.NoError(err, "can apply "+manifestVhoRemoveXBar)

	// Check status is accepted on VirtualHostOption
	s.testInstallation.AssertionsT(s.T()).EventuallyResourceStatusMatchesState(
		s.getterForMeta(&vhoRemoveXBar),
		core.Status_Accepted,
		defaults.KubeGatewayReporter,
	)

	output, err := s.testInstallation.Actions.Kubectl().ApplyFileWithOutput(s.ctx, manifestVhoWebhookReject)
	s.testInstallation.AssertionsT(s.T()).ExpectObjectAdmitted(manifestVhoWebhookReject, err, output,
		"Validating *v1.VirtualHostOption failed")

	if !s.testInstallation.Metadata.ValidationAlwaysAccept {
		s.testInstallation.AssertionsT(s.T()).ExpectGlooObjectNotExist(
			s.ctx,
			s.getterForMeta(&vhoWebhookReject),
			&vhoWebhookReject,
		)
	} else {
		// Check status is rejected on bad VirtualHostOption
		s.testInstallation.AssertionsT(s.T()).EventuallyResourceStatusMatchesState(
			s.getterForMeta(&vhoWebhookReject),
			core.Status_Rejected,
			defaults.KubeGatewayReporter,
		)
	}

	// Check healthy response with no x-bar header
	s.testInstallation.AssertionsT(s.T()).AssertEventualCurlResponse(
		s.ctx,
		testdefaults.CurlPodExecOpt,
		[]curl.Option{
			curl.WithHost(kubeutils.ServiceFQDN(proxyService1.ObjectMeta)),
			curl.WithHostHeader("example.com"),
		},
		expectedResponseWithoutXBar)
}

// TestConfigureVirtualHostOptionsWithSectionName tests a complex scenario where multiple
// VirtualHostOptions conflicting across multiple listeners are applied to a single gateway
//
// The goal here is to test the behavior when multiple VHOs target a gateway with multiple listeners and only some
// conflict. This will generate a warning on the conflicted resource, but the VHO should be attached properly and
// options propagated for the listener.
func (s *testingSuite) TestConfigureVirtualHostOptionsWithSectionNameManualSetup() {
	s.T().Cleanup(func() {
		output, err := s.testInstallation.Actions.Kubectl().DeleteFileWithOutput(s.ctx, manifestVhoRemoveXBar)
		s.testInstallation.AssertionsT(s.T()).ExpectObjectDeleted(manifestVhoRemoveXBar, err, output)

		output, err = s.testInstallation.Actions.Kubectl().DeleteFileWithOutput(s.ctx, manifestVhoRemoveXBaz)
		s.testInstallation.AssertionsT(s.T()).ExpectObjectDeleted(manifestVhoRemoveXBaz, err, output)

		output, err = s.testInstallation.Actions.Kubectl().DeleteFileWithOutput(s.ctx, manifestVhoSectionAddXFoo)
		s.testInstallation.AssertionsT(s.T()).ExpectObjectDeleted(manifestVhoSectionAddXFoo, err, output)
	})

	// Apply our manifests so we can assert that basic vho exists before applying conflicting VHOs.
	// This is needed because our solo-kit clients currently do not return creationTimestamp
	err := s.testInstallation.Actions.Kubectl().ApplyFile(s.ctx, manifestVhoRemoveXBar)
	s.NoError(err, "can apply "+manifestVhoRemoveXBar)

	// Check status is accepted before moving on to apply conflicting vho
	s.testInstallation.AssertionsT(s.T()).EventuallyResourceStatusMatchesState(
		s.getterForMeta(&vhoRemoveXBar),
		core.Status_Accepted,
		defaults.KubeGatewayReporter,
	)

	err = s.testInstallation.Actions.Kubectl().ApplyFile(s.ctx, manifestVhoRemoveXBaz)
	s.NoError(err, "can apply "+manifestVhoRemoveXBaz)

	err = s.testInstallation.Actions.Kubectl().ApplyFile(s.ctx, manifestVhoSectionAddXFoo)
	s.NoError(err, "can apply "+manifestVhoSectionAddXFoo)

	// Check status is accepted on VirtualHostOption with section name
	s.testInstallation.AssertionsT(s.T()).EventuallyResourceStatusMatchesState(
		s.getterForMeta(&vhoSectionAddXFoo),
		core.Status_Accepted,
		defaults.KubeGatewayReporter,
	)

	// Check status is warning on VirtualHostOption not selected for attachment
	// to either of the listeners
	s.testInstallation.AssertionsT(s.T()).EventuallyResourceStatusMatchesWarningReasons(
		s.getterForMeta(&vhoRemoveXBaz),
		[]string{
			"VirtualHostOption 'default/remove-x-baz-header' not attached to listener 'other' on Gateway 'default/gw-1' due to conflict with more specific or older VirtualHostOptions 'default/remove-x-bar-header'",
			"VirtualHostOption 'default/remove-x-baz-header' not attached to listener 'http' on Gateway 'default/gw-1' due to conflict with more specific or older VirtualHostOptions 'default/add-x-foo-header-section'",
		},
		defaults.KubeGatewayReporter,
	)

	// Check status is warning on VirtualHostOption with conflicting attachment,
	// despite being properly attached to 8081 listener
	s.testInstallation.AssertionsT(s.T()).EventuallyResourceStatusMatchesWarningReasons(
		s.getterForMeta(&vhoRemoveXBar),
		[]string{"conflict with more specific or older VirtualHostOptions"},
		defaults.KubeGatewayReporter,
	)

	// Check healthy response with added foo header to listener targeted by sectionName
	s.testInstallation.AssertionsT(s.T()).AssertEventualCurlResponse(
		s.ctx,
		testdefaults.CurlPodExecOpt,
		[]curl.Option{
			curl.WithHost(kubeutils.ServiceFQDN(proxyService1.ObjectMeta)),
			curl.WithHostHeader("example.com"),
			curl.WithPort(gw1port1),
		},
		&matchers.HttpResponse{
			StatusCode: http.StatusOK,
			Custom: gomega.And(
				// attached to this listener
				matchers.ContainHeaderKeys([]string{"x-foo"}),
				// not removed because conflicts with earlier VHO
				matchers.ContainHeaderKeys([]string{"x-bar"}),
				// not removed because conflicts with earlier VHO
				matchers.ContainHeaderKeys([]string{"x-baz"}),
			),
			Body: gstruct.Ignore(),
		})

	// Check healthy response with x-bar removed to listener NOT targeted by sectionName
	s.testInstallation.AssertionsT(s.T()).AssertEventualCurlResponse(
		s.ctx,
		testdefaults.CurlPodExecOpt,
		[]curl.Option{
			curl.WithHost(kubeutils.ServiceFQDN(proxyService1.ObjectMeta)),
			curl.WithHostHeader("example.com"),
			curl.WithPort(gw1port2),
		},
		&matchers.HttpResponse{
			StatusCode: http.StatusOK,
			Custom: gomega.And(
				// not attached to this listener
				gomega.Not(matchers.ContainHeaderKeys([]string{"x-foo"})),
				// removed by the earliest VHO
				gomega.Not(matchers.ContainHeaderKeys([]string{"x-bar"})),
				// not removed because conflicts with earlier VHO
				matchers.ContainHeaderKeys([]string{"x-baz"}),
			),
			Body: gstruct.Ignore(),
		})
}

// TestMultipleVirtualHostOptionsSetup tests a complex scenario where multiple VirtualHostOptions conflict
//
// The goal here is to test the behavior when multiple VHOs are targeting a gateway without sectionName.
// The expected behavior is that the oldest resource is used
func (s *testingSuite) TestMultipleVirtualHostOptionsSetup() {
	s.T().Cleanup(func() {
		output, err := s.testInstallation.Actions.Kubectl().DeleteFileWithOutput(s.ctx, manifestVhoRemoveXBar)
		s.testInstallation.AssertionsT(s.T()).ExpectObjectDeleted(manifestVhoRemoveXBar, err, output)

		output, err = s.testInstallation.Actions.Kubectl().DeleteFileWithOutput(s.ctx, manifestVhoRemoveXBaz)
		s.testInstallation.AssertionsT(s.T()).ExpectObjectDeleted(manifestVhoRemoveXBaz, err, output)
	})

	// Manually apply our manifests so we can assert that basic vho exists before applying extra vho.
	// This is needed because our solo-kit clients currently do not return creationTimestamp
	err := s.testInstallation.Actions.Kubectl().ApplyFile(s.ctx, manifestVhoRemoveXBar)
	s.NoError(err, "can apply "+manifestVhoRemoveXBar)

	err = s.testInstallation.Actions.Kubectl().ApplyFile(s.ctx, manifestVhoRemoveXBaz)
	s.NoError(err, "can apply "+manifestVhoRemoveXBaz)

	// Check status is warning on newer VirtualHostOption not selected for attachment
	s.testInstallation.AssertionsT(s.T()).EventuallyResourceStatusMatchesWarningReasons(
		s.getterForMeta(&vhoRemoveXBaz),
		[]string{"conflict with more specific or older VirtualHostOptions"},
		defaults.KubeGatewayReporter,
	)

	// Check status is accepted on older VirtualHostOption
	s.testInstallation.AssertionsT(s.T()).EventuallyResourceStatusMatchesState(
		s.getterForMeta(&vhoRemoveXBar),
		core.Status_Accepted,
		defaults.KubeGatewayReporter,
	)

	// Check healthy response with no x-bar header
	s.testInstallation.AssertionsT(s.T()).AssertEventualCurlResponse(
		s.ctx,
		testdefaults.CurlPodExecOpt,
		[]curl.Option{
			curl.WithHost(kubeutils.ServiceFQDN(proxyService1.ObjectMeta)),
			curl.WithHostHeader("example.com"),
		},
		expectedResponseWithoutXBar)
}

// TestConfigureVirtualHostOptionsWarningMultipleGatewaysSetup tests a complex scenario where multiple
// VirtualHostOptions conflicting across multiple listeners and multiple gateways are applied
//
// The goal here is to confirm that that if multiple gateways have warnings on the same resource,
// the warnings are aggregated and presented in the status of the resource.
func (s *testingSuite) TestConfigureVirtualHostOptionsWarningMultipleGatewaysSetup() {
	s.T().Cleanup(func() {
		output, err := s.testInstallation.Actions.Kubectl().DeleteFileWithOutput(s.ctx, manifestVhoMultipleGatewayWarnings)
		s.testInstallation.AssertionsT(s.T()).ExpectObjectDeleted(manifestVhoMultipleGatewayWarnings, err, output)
	})

	err := s.testInstallation.Actions.Kubectl().ApplyFile(s.ctx, manifestVhoMultipleGatewayWarnings)
	s.NoError(err, "can apply "+manifestVhoMultipleGatewayWarnings)

	// Check status is accepted on VirtualHostOption with section name
	s.testInstallation.AssertionsT(s.T()).EventuallyResourceStatusMatchesState(
		s.getterForMeta(&metav1.ObjectMeta{
			Name:      "add-x-foo-header-2",
			Namespace: "default",
		}),
		core.Status_Accepted,
		defaults.KubeGatewayReporter,
	)

	// Check status is warning on VirtualHostOption not selected for attachment
	// to either of the gateways and listeners
	s.testInstallation.AssertionsT(s.T()).EventuallyResourceStatusMatchesWarningReasons(
		s.getterForMeta(&metav1.ObjectMeta{
			Name:      "remove-x-baz-header-2",
			Namespace: "default",
		}),
		[]string{
			"VirtualHostOption 'default/remove-x-baz-header-2' not attached to listener 'http' on Gateway 'default/gw-1' due to conflict with more specific or older VirtualHostOptions 'default/add-x-foo-header-2'",
			"VirtualHostOption 'default/remove-x-baz-header-2' not attached to listener 'other' on Gateway 'default/gw-1' due to conflict with more specific or older VirtualHostOptions 'default/remove-x-bar-header-2'",
			"VirtualHostOption 'default/remove-x-baz-header-2' not attached to listener 'http' on Gateway 'default/gw-2' due to conflict with more specific or older VirtualHostOptions 'default/add-x-foo-header-2'",
			"VirtualHostOption 'default/remove-x-baz-header-2' not attached to listener 'other' on Gateway 'default/gw-2' due to conflict with more specific or older VirtualHostOptions 'default/remove-x-bar-header-2'",
		},

		defaults.KubeGatewayReporter,
	)

	// Check status is warning on VirtualHostOption with conflicting attachment,
	// despite being properly attached to 8081 listener
	s.testInstallation.AssertionsT(s.T()).EventuallyResourceStatusMatchesWarningReasons(
		s.getterForMeta(&metav1.ObjectMeta{
			Name:      "remove-x-bar-header-2",
			Namespace: "default",
		}),
		[]string{"conflict with more specific or older VirtualHostOptions"},
		defaults.KubeGatewayReporter,
	)

	// Check healthy response with added foo header to listener targeted by sectionName
	s.testInstallation.AssertionsT(s.T()).AssertEventualCurlResponse(
		s.ctx,
		testdefaults.CurlPodExecOpt,
		[]curl.Option{
			curl.WithHost(kubeutils.ServiceFQDN(proxyService1.ObjectMeta)),
			curl.WithHostHeader("example.com"),
			curl.WithPort(gw1port1),
		},
		&matchers.HttpResponse{
			StatusCode: http.StatusOK,
			Custom: gomega.And(
				// attached to this listener
				matchers.ContainHeaderKeys([]string{"x-foo"}),
				// not removed because conflicts with earlier VHO
				matchers.ContainHeaderKeys([]string{"x-bar"}),
				// not removed because conflicts with earlier VHO
				matchers.ContainHeaderKeys([]string{"x-baz"}),
			),
			Body: gstruct.Ignore(),
		})

	// Check healthy response with x-bar removed to listener NOT targeted by sectionName
	s.testInstallation.AssertionsT(s.T()).AssertEventualCurlResponse(
		s.ctx,
		testdefaults.CurlPodExecOpt,
		[]curl.Option{
			curl.WithHost(kubeutils.ServiceFQDN(proxyService1.ObjectMeta)),
			curl.WithHostHeader("example.com"),
			curl.WithPort(gw1port2),
		},
		&matchers.HttpResponse{
			StatusCode: http.StatusOK,
			Custom: gomega.And(
				// not attached to this listener
				gomega.Not(matchers.ContainHeaderKeys([]string{"x-foo"})),
				// removed by the earliest VHO
				gomega.Not(matchers.ContainHeaderKeys([]string{"x-bar"})),
				// not removed because conflicts with earlier VHO
				matchers.ContainHeaderKeys([]string{"x-baz"}),
			),
			Body: gstruct.Ignore(),
		})
}

// TestDeletingNonConflictingVirtualHostOptions tests the behavior when a VHO that was blocking
// another VHO is deleted
//
// The expected behavior is that the previously blocked VHO is now attached and the
// headers are mutated as expected
func (s *testingSuite) TestDeletingConflictingVirtualHostOptions() {
	s.T().Cleanup(func() {
		// this should already be deleted, confirm
		s.testInstallation.AssertionsT(s.T()).ExpectGlooObjectNotExist(
			s.ctx,
			s.getterForMeta(&vhoRemoveXBar),
			&vhoRemoveXBar,
		)

		output, err := s.testInstallation.Actions.Kubectl().DeleteFileWithOutput(s.ctx, manifestVhoRemoveXBaz)
		s.testInstallation.AssertionsT(s.T()).ExpectObjectDeleted(manifestVhoRemoveXBaz, err, output)
	})

	err := s.testInstallation.Actions.Kubectl().ApplyFile(s.ctx, manifestVhoRemoveXBar)
	s.NoError(err, "can apply "+manifestVhoRemoveXBar)

	err = s.testInstallation.Actions.Kubectl().ApplyFile(s.ctx, manifestVhoRemoveXBaz)
	s.NoError(err, "can apply "+manifestVhoRemoveXBaz)

	// Check status is warning on newer VirtualHostOption not selected for attachment
	s.testInstallation.AssertionsT(s.T()).EventuallyResourceStatusMatchesWarningReasons(
		s.getterForMeta(&vhoRemoveXBaz),
		[]string{"conflict with more specific or older VirtualHostOptions"},
		defaults.KubeGatewayReporter,
	)

	// Check status is accepted on older VirtualHostOption
	s.testInstallation.AssertionsT(s.T()).EventuallyResourceStatusMatchesState(
		s.getterForMeta(&vhoRemoveXBar),
		core.Status_Accepted,
		defaults.KubeGatewayReporter,
	)

	// Check healthy response with no x-bar header
	s.testInstallation.AssertionsT(s.T()).AssertEventualCurlResponse(
		s.ctx,
		testdefaults.CurlPodExecOpt,
		[]curl.Option{
			curl.WithHost(kubeutils.ServiceFQDN(proxyService1.ObjectMeta)),
			curl.WithHostHeader("example.com"),
		},
		expectedResponseWithoutXBar)

	// Delete the VHO that was blocking the other VHO
	output, err := s.testInstallation.Actions.Kubectl().DeleteFileWithOutput(s.ctx, manifestVhoRemoveXBar)
	s.testInstallation.AssertionsT(s.T()).ExpectObjectDeleted(manifestVhoRemoveXBar, err, output)

	// Check status is accepted on VirtualHostOption
	s.testInstallation.AssertionsT(s.T()).EventuallyResourceStatusMatchesState(
		s.getterForMeta(&vhoRemoveXBaz),
		core.Status_Accepted,
		defaults.KubeGatewayReporter,
	)

	// Check healthy response with no x-bar header
	s.testInstallation.AssertionsT(s.T()).AssertEventualCurlResponse(
		s.ctx,
		testdefaults.CurlPodExecOpt,
		[]curl.Option{
			curl.WithHost(kubeutils.ServiceFQDN(proxyService1.ObjectMeta)),
			curl.WithHostHeader("example.com"),
		},
		expectedResponseWithoutXBaz)
}

// TestOptionsMerge tests shallow merging of VirtualHostOptions larger in the precedence chain
func (s *testingSuite) TestOptionsMerge() {
	s.T().Cleanup(func() {
		output, err := s.testInstallation.Actions.Kubectl().DeleteFileWithOutput(s.ctx, manifestVhoRemoveXBar)
		s.testInstallation.AssertionsT(s.T()).ExpectObjectDeleted(manifestVhoRemoveXBar, err, output)

		output, err = s.testInstallation.Actions.Kubectl().DeleteFileWithOutput(s.ctx, manifestVhoMergeRemoveXBaz)
		s.testInstallation.AssertionsT(s.T()).ExpectObjectDeleted(manifestVhoMergeRemoveXBaz, err, output)
	})

	_, err := s.testInstallation.Actions.Kubectl().ApplyFileWithOutput(s.ctx, manifestVhoRemoveXBar)
	s.Require().NoError(err)

	_, err = s.testInstallation.Actions.Kubectl().ApplyFileWithOutput(s.ctx, manifestVhoMergeRemoveXBaz)
	s.Require().NoError(err)

	s.testInstallation.AssertionsT(s.T()).EventuallyResourceStatusMatchesState(
		s.getterForMeta(&vhoMergeRemoveXBaz),
		core.Status_Accepted,
		defaults.KubeGatewayReporter,
	)

	s.testInstallation.AssertionsT(s.T()).EventuallyResourceStatusMatchesState(
		s.getterForMeta(&vhoRemoveXBar),
		core.Status_Accepted,
		defaults.KubeGatewayReporter,
	)

	s.testInstallation.AssertionsT(s.T()).AssertEventualCurlResponse(
		s.ctx,
		testdefaults.CurlPodExecOpt,
		[]curl.Option{
			curl.WithHost(kubeutils.ServiceFQDN(proxyService1.ObjectMeta)),
			curl.WithHostHeader("example.com"),
		},
		// Expect:
		// - x-bar header to be removed by vho-remove-x-bar
		// - x-baz header to not be removed as the option conflicts with vho-remove-x-bar and is not merged
		// - x-envoy-attempt-count header to be added by vho-merge-remove-x-baz.yaml
		&matchers.HttpResponse{
			StatusCode: http.StatusOK,
			Custom: gomega.And(
				gomega.Not(matchers.ContainHeaderKeys([]string{"x-bar"})),
				matchers.ContainHeaderKeys([]string{"x-baz"}),
				matchers.ContainHeaderKeys([]string{"x-envoy-attempt-count"}),
			),
			Body: gstruct.Ignore(),
		},
	)
}

func (s *testingSuite) getterForMeta(meta *metav1.ObjectMeta) helpers.InputResourceGetter {
	return func() (resources.InputResource, error) {
		return s.testInstallation.ResourceClients.VirtualHostOptionClient().Read(meta.GetNamespace(),
			meta.GetName(), clients.ReadOpts{})
	}
}

// testExpectedResponsesForManifests tests is a utility function that applies a set of manifests and runs a set of curls with defined matchers
// matchersForListeners is map of service fqdn to map of port to matcher
func (s *testingSuite) testExpectedResponsesForManifests(manifests map[string]*metav1.ObjectMeta, matchersForListeners map[string]map[int]*matchers.HttpResponse, validateAcceptedStatus bool) {
	s.T().Cleanup(func() {
		for manifestFile := range manifests {
			output, err := s.testInstallation.Actions.Kubectl().DeleteFileWithOutput(s.ctx, manifestFile)
			s.testInstallation.AssertionsT(s.T()).ExpectObjectDeleted(manifestFile, err, output)
		}
	})

	// Setup and validate the manifests
	for manifestFile, manifestMeta := range manifests {
		err := s.testInstallation.Actions.Kubectl().ApplyFile(s.ctx, manifestFile)
		s.NoError(err, "can apply "+manifestFile)

		if validateAcceptedStatus {
			s.testInstallation.AssertionsT(s.T()).EventuallyResourceStatusMatchesState(
				s.getterForMeta(manifestMeta),
				core.Status_Accepted,
				defaults.KubeGatewayReporter,
			)
		}
	}

	for host, ports := range matchersForListeners {
		for port, matcher := range ports {
			s.testInstallation.AssertionsT(s.T()).AssertEventualCurlResponse(
				s.ctx,
				testdefaults.CurlPodExecOpt,
				[]curl.Option{
					curl.WithHost(host),
					curl.WithHostHeader("example.com"),
					curl.WithPort(port),
				},
				matcher,
			)
		}
	}
}
