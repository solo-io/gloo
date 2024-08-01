package deployer

import (
	"context"

	"github.com/stretchr/testify/suite"

	"github.com/solo-io/gloo/test/kubernetes/e2e"
)

var _ e2e.NewSuiteFunc = NewMinimalDefaultGatewayParametersTestingSuite

// minimalDefaultGatewayParametersDeployerSuite tests the "deployer" feature in situations where users have applied `null` values
// to as many of the helm values controlling the default GatewayParameters for the gloo-gateway GatewayClass as possible.
// The "deployer" code can be found here: /projects/gateway2/deployer
type minimalDefaultGatewayParametersDeployerSuite struct {
	suite.Suite

	ctx context.Context

	// testInstallation contains all the metadata/utilities necessary to execute a series of tests
	// against an installation of Gloo Gateway
	testInstallation *e2e.TestInstallation
}

func NewMinimalDefaultGatewayParametersTestingSuite(ctx context.Context, testInst *e2e.TestInstallation) suite.TestingSuite {
	return &minimalDefaultGatewayParametersDeployerSuite{
		ctx:              ctx,
		testInstallation: testInst,
	}
}

// This test has been commented out as a bug in helm prevents the runAsUser in the OSS sub-chart to be removed
// by setting it as null in the enterprise chart.
// Ref: https://github.com/helm/helm/issues/12637
// TODO (davidjumani): Add this back once the bug has been fixed or remove it once an alternative mechanism has been provied
// func (s *minimalDefaultGatewayParametersDeployerSuite) TestConfigureProxiesFromGatewayParameters() {
// 	s.T().Cleanup(func() {
// 		err := s.testInstallation.Actions.Kubectl().DeleteFile(s.ctx, gwParametersManifestFile)
// 		s.NoError(err, "can delete basic gateway manifest")
// 		s.testInstallation.Assertions.EventuallyObjectsNotExist(s.ctx, gwParams, proxyService, proxyDeployment)
// 	})

// 	err := s.testInstallation.Actions.Kubectl().ApplyFile(s.ctx, gwParametersManifestFile)
// 	s.Require().NoError(err, "can apply basic gateway manifest")
// 	s.testInstallation.Assertions.EventuallyObjectsExist(s.ctx, gwParams, proxyService, proxyDeployment)

// 	deployment, err := s.testInstallation.ClusterContext.Clientset.AppsV1().Deployments(proxyDeployment.GetNamespace()).Get(s.ctx, proxyDeployment.GetName(), metav1.GetOptions{})
// 	s.Require().NoError(err, "can get deployment")
// 	s.Require().Len(deployment.Spec.Template.Spec.Containers, 1)
// 	secCtx := deployment.Spec.Template.Spec.Containers[0].SecurityContext
// 	s.Require().NotNil(secCtx)
// 	s.Require().Nil(secCtx.RunAsUser)
// 	s.Require().NotNil(secCtx.RunAsNonRoot)
// 	s.Require().False(*secCtx.RunAsNonRoot)
// 	s.Require().NotNil(secCtx.AllowPrivilegeEscalation)
// 	s.Require().True(*secCtx.AllowPrivilegeEscalation)
// }
