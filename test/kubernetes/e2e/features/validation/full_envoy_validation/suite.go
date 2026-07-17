package full_envoy_validation

import (
	"context"
	"fmt"

	"github.com/solo-io/gloo/test/kubernetes/e2e"
	testdefaults "github.com/solo-io/gloo/test/kubernetes/e2e/defaults"
	"github.com/solo-io/gloo/test/kubernetes/e2e/features/validation"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/stretchr/testify/suite"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ e2e.NewSuiteFunc = NewTestingSuite

// testingSuite is the entire Suite of tests for the webhook validation fullEnvoyValidation=true feature
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

// TestRejectInvalidTransformation checks webhook rejects invalid transformation when fullEnvoyValidation=true
func (s *testingSuite) TestRejectInvalidTransformation() {

	s.T().Cleanup(func() {
		err := s.testInstallation.Actions.Kubectl().DeleteFileSafe(s.ctx, validation.ExampleUpstream, "-n", s.testInstallation.Metadata.InstallNamespace)
		s.Assertions.NoError(err, "can delete upstream")

		err = s.testInstallation.Actions.Kubectl().DeleteFileSafe(s.ctx, testdefaults.NginxPodManifest)
		s.Assertions.NoError(err, "can delete nginx pod")
	})

	err := s.testInstallation.Actions.Kubectl().ApplyFile(s.ctx, testdefaults.NginxPodManifest)
	s.Assert().NoError(err)
	// Check that test resources are running
	s.testInstallation.Assertions.EventuallyPodsRunning(s.ctx, testdefaults.NginxPod.ObjectMeta.GetNamespace(), metav1.ListOptions{
		LabelSelector: "app.kubernetes.io/name=nginx",
	})

	err = s.testInstallation.Actions.Kubectl().ApplyFile(s.ctx, validation.ExampleUpstream, "-n", s.testInstallation.Metadata.InstallNamespace)
	s.Assert().NoError(err)
	// Upstreams no longer report status if they have not been translated at all to avoid conflicting with
	// other syncers that have translated them, so we can only detect that the objects exist here
	s.testInstallation.Assertions.EventuallyResourceExists(
		func() (resources.Resource, error) {
			return s.testInstallation.ResourceClients.UpstreamClient().Read(s.testInstallation.Metadata.InstallNamespace, "nginx-upstream", clients.ReadOpts{Ctx: s.ctx})
		},
	)

	// we need to make sure Gloo has had a chance to process it
	s.testInstallation.Assertions.ConsistentlyResourceExists(
		s.ctx,
		func() (resources.Resource, error) {
			return s.testInstallation.ResourceClients.UpstreamClient().Read(s.testInstallation.Metadata.InstallNamespace, "nginx-upstream", clients.ReadOpts{Ctx: s.ctx})
		},
	)
	s.T().Cleanup(func() {
		err := s.testInstallation.Actions.Kubectl().DeleteFileSafe(s.ctx, validation.VSTransformationHeaderText, "-n", s.testInstallation.Metadata.InstallNamespace)
		s.Assert().NoError(err)
	})
	// rejects invalid inja template in transformation
	output, err := s.testInstallation.Actions.Kubectl().ApplyFileWithOutput(s.ctx, validation.VSTransformationHeaderText, "-n", s.testInstallation.Metadata.InstallNamespace)
	s.Assert().Error(err)
	s.Assert().Contains(output, "Failed to parse response template: Failed to parse "+
		"header template ':status': [inja.exception.parser_error] (at 1:92) expected statement close, got '%'")
}

// TestRejectsDeleteOfInUseUpstream verifies that in a hybrid install (Edge + K8s Gateway API), deleting an
// Upstream that a VirtualService still routes to is rejected by the webhook. Upstream admission in K8s
// Gateway mode also validates against the real Edge proxies, so this in-use check must keep working.
func (s *testingSuite) TestRejectsDeleteOfInUseUpstream() {
	s.T().Cleanup(func() {
		// Delete the VirtualService first so the Upstream is no longer in use and can be removed.
		err := s.testInstallation.Actions.Kubectl().DeleteFileSafe(s.ctx, validation.ExampleVS, "-n", s.testInstallation.Metadata.InstallNamespace)
		s.Assert().NoError(err)
		err = s.testInstallation.Actions.Kubectl().DeleteFileSafe(s.ctx, validation.ExampleUpstream, "-n", s.testInstallation.Metadata.InstallNamespace)
		s.Assert().NoError(err)
		err = s.testInstallation.Actions.Kubectl().DeleteFileSafe(s.ctx, testdefaults.NginxPodManifest)
		s.Assert().NoError(err)
	})

	// nginx backs the example Upstream.
	err := s.testInstallation.Actions.Kubectl().ApplyFile(s.ctx, testdefaults.NginxPodManifest)
	s.Assert().NoError(err)
	s.testInstallation.Assertions.EventuallyPodsRunning(s.ctx, testdefaults.NginxPod.ObjectMeta.GetNamespace(), metav1.ListOptions{
		LabelSelector: "app.kubernetes.io/name=nginx",
	})

	err = s.testInstallation.Actions.Kubectl().ApplyFile(s.ctx, validation.ExampleUpstream, "-n", s.testInstallation.Metadata.InstallNamespace)
	s.Assert().NoError(err, "can apply the upstream")

	err = s.testInstallation.Actions.Kubectl().ApplyFile(s.ctx, validation.ExampleVS, "-n", s.testInstallation.Metadata.InstallNamespace)
	s.Assert().NoError(err, "can apply the virtual service routing to the upstream")

	// Deleting the Upstream while the VirtualService still routes to it must be rejected.
	output, err := s.testInstallation.Actions.Kubectl().DeleteFileWithOutput(s.ctx, validation.ExampleUpstream, "-n", s.testInstallation.Metadata.InstallNamespace)
	s.Assert().Error(err, "deleting an in-use upstream should be rejected")
	s.Assert().Contains(output, "admission webhook")
	s.Assert().Contains(output, validation.ExampleUpstreamName)
}

// TestLargeConfiguration checks webhook accepts large configuration when fullEnvoyValidation=true
func (s *testingSuite) TestLargeConfiguration() {
	s.T().Skip("we need to make sure we have all formats working")
	s.T().Cleanup(func() {
		err := s.testInstallation.Actions.Kubectl().DeleteFileSafe(s.ctx, validation.LargeConfiguration, "-n",
			s.testInstallation.Metadata.InstallNamespace)
		s.Assertions.NoError(err, "can delete large configuration")

		err = s.testInstallation.Actions.Kubectl().DeleteFileSafe(s.ctx, validation.ExampleUpstream)
		s.Assertions.NoError(err, "can delete example upstream")
	})

	err := s.testInstallation.Actions.Kubectl().ApplyFile(s.ctx, validation.ExampleUpstream, "-n", s.testInstallation.Metadata.InstallNamespace)
	s.Assert().NoError(err)

	err = s.testInstallation.Actions.Kubectl().ApplyFile(s.ctx, validation.LargeConfiguration, "-n",
		s.testInstallation.Metadata.InstallNamespace)
	fmt.Println(err)
}
