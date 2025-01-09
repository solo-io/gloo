package deployer

import (
	"context"
	"fmt"
	"time"

	tlsv3 "github.com/envoyproxy/go-control-plane/envoy/extensions/transport_sockets/tls/v3"
	"github.com/onsi/gomega"
	"github.com/onsi/gomega/gstruct"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	gwv1 "sigs.k8s.io/gateway-api/apis/v1"

	"github.com/solo-io/gloo/pkg/utils/envoyutils/admincli"
	"github.com/solo-io/gloo/pkg/utils/kubeutils"
	"github.com/solo-io/gloo/projects/gateway2/api/v1alpha1"
	"github.com/solo-io/gloo/projects/gateway2/wellknown"
	"github.com/solo-io/gloo/projects/gloo/pkg/syncer/setup"
	"github.com/solo-io/gloo/projects/gloo/pkg/utils"
	"github.com/solo-io/gloo/test/kubernetes/e2e"
	testdefaults "github.com/solo-io/gloo/test/kubernetes/e2e/defaults"
)

var _ e2e.NewSuiteFunc = NewTestingSuite

// testingSuite is the entire Suite of tests for the "deployer" feature
// The "deployer" code can be found here: /projects/gateway2/deployer
type testingSuite struct {
	suite.Suite

	ctx context.Context

	// testInstallation contains all the metadata/utilities necessary to execute a series of tests
	// against an installation of Gloo Gateway
	testInstallation *e2e.TestInstallation

	// manifests maps test name to a list of manifests to apply before the test
	manifests map[string][]string

	// manifestObjects maps a manifest file to a list of objects that are contained in that file
	manifestObjects map[string][]client.Object
}

func NewTestingSuite(ctx context.Context, testInst *e2e.TestInstallation) suite.TestingSuite {
	return &testingSuite{
		ctx:              ctx,
		testInstallation: testInst,
	}
}

func (s *testingSuite) SetupSuite() {
	s.manifests = map[string][]string{
		"TestProvisionDeploymentAndService": {
			testdefaults.NginxPodManifest,
			gatewayWithoutParameters,
		},
		"TestConfigureProxiesFromGatewayParameters": {
			testdefaults.NginxPodManifest,
			gatewayParametersCustom,
			gatewayWithParameters,
		},
		"TestProvisionResourcesUpdatedWithValidParameters": {
			testdefaults.NginxPodManifest,
			gatewayWithParameters,
		},
		"TestProvisionResourcesNotUpdatedWithInvalidParameters": {
			testdefaults.NginxPodManifest,
			gatewayWithParameters,
		},
		"TestSelfManagedGateway": {
			selfManagedGateway,
		},
		"TestConfigureAwsLambda": {
			testdefaults.NginxPodManifest,
			gatewayWithoutParameters,
		},
	}
	s.manifestObjects = map[string][]client.Object{
		testdefaults.NginxPodManifest: {testdefaults.NginxPod, testdefaults.NginxSvc},
		gatewayWithoutParameters:      {proxyService, proxyServiceAccount, proxyDeployment},
		gatewayWithParameters:         {proxyService, proxyServiceAccount, proxyDeployment, gwParamsDefault},
		gatewayParametersCustom:       {gwParamsCustom},
		selfManagedGateway:            {gwParamsDefault},
	}
}

func (s *testingSuite) TearDownSuite() {
	// nothing at the moment
}

func (s *testingSuite) BeforeTest(suiteName, testName string) {
	manifests := s.manifests[testName]
	for _, manifest := range manifests {
		err := s.testInstallation.Actions.Kubectl().ApplyFile(s.ctx, manifest)
		s.Require().NoError(err)
		s.testInstallation.Assertions.EventuallyObjectsExist(s.ctx, s.manifestObjects[manifest]...)
	}
}

func (s *testingSuite) AfterTest(suiteName, testName string) {
	manifests := s.manifests[testName]
	for _, manifest := range manifests {
		err := s.testInstallation.Actions.Kubectl().DeleteFileSafe(s.ctx, manifest)
		s.Require().NoError(err)
		s.testInstallation.Assertions.EventuallyObjectsNotExist(s.ctx, s.manifestObjects[manifest]...)
	}
}

func (s *testingSuite) TestProvisionDeploymentAndService() {
	s.testInstallation.Assertions.EventuallyRunningReplicas(s.ctx, proxyDeployment.ObjectMeta, gomega.Equal(1))
}

func (s *testingSuite) TestConfigureProxiesFromGatewayParameters() {
	s.testInstallation.Assertions.EventuallyRunningReplicas(s.ctx, proxyDeployment.ObjectMeta, gomega.Equal(1))

	// check that the labels and annotations got passed through from GatewayParameters to the ServiceAccount
	sa := &corev1.ServiceAccount{}
	err := s.testInstallation.ClusterContext.Client.Get(s.ctx,
		types.NamespacedName{Name: glooProxyObjectMeta.Name, Namespace: glooProxyObjectMeta.Namespace},
		sa)
	s.Require().NoError(err)
	s.testInstallation.Assertions.Gomega.Expect(sa.GetLabels()).To(
		gomega.HaveKeyWithValue("sa-label-key", "sa-label-val"))
	s.testInstallation.Assertions.Gomega.Expect(sa.GetAnnotations()).To(
		gomega.HaveKeyWithValue("sa-anno-key", "sa-anno-val"))

	// check that the labels and annotations got passed through from GatewayParameters to the Service
	svc := &corev1.Service{}
	err = s.testInstallation.ClusterContext.Client.Get(s.ctx,
		types.NamespacedName{Name: glooProxyObjectMeta.Name, Namespace: glooProxyObjectMeta.Namespace},
		svc)
	s.Require().NoError(err)
	s.testInstallation.Assertions.Gomega.Expect(svc.GetLabels()).To(
		gomega.HaveKeyWithValue("svc-label-key", "svc-label-val"))
	s.testInstallation.Assertions.Gomega.Expect(svc.GetAnnotations()).To(
		gomega.HaveKeyWithValue("svc-anno-key", "svc-anno-val"))

	// check that external traffic policy got passwed through from GatewayParameters to the Service
	s.testInstallation.Assertions.Gomega.Expect(svc.Spec.ExternalTrafficPolicy).To(
		gomega.Equal(corev1.ServiceExternalTrafficPolicyCluster))

	// Update the Gateway to use the custom GatewayParameters
	gwName := types.NamespacedName{Name: gw.Name, Namespace: gw.Namespace}
	err = s.testInstallation.ClusterContext.Client.Get(s.ctx, gwName, gw)
	s.Require().NoError(err)
	s.patchGateway(gw.ObjectMeta, func(gw *gwv1.Gateway) {
		gw.Annotations[wellknown.GatewayParametersAnnotationName] = gwParamsCustom.Name
	})

	// Assert that the expected custom configuration exists.
	s.testInstallation.Assertions.EventuallyRunningReplicas(s.ctx, proxyDeployment.ObjectMeta, gomega.Equal(2))

	s.testInstallation.Assertions.AssertEnvoyAdminApi(
		s.ctx,
		proxyDeployment.ObjectMeta,
		serverInfoLogLevelAssertion(s.testInstallation, "debug", "connection:trace,upstream:debug"),
		xdsClusterAssertion(s.testInstallation),
	)
}

func (s *testingSuite) TestConfigureAwsLambda() {
	s.testInstallation.Assertions.EventuallyRunningReplicas(s.ctx, proxyDeployment.ObjectMeta, gomega.Equal(1))

	s.testInstallation.Assertions.AssertEnvoyAdminApi(
		s.ctx,
		proxyDeployment.ObjectMeta,
		awsStsClusterAssertion(s.testInstallation),
	)
}

func (s *testingSuite) TestProvisionResourcesUpdatedWithValidParameters() {
	s.testInstallation.Assertions.EventuallyRunningReplicas(s.ctx, proxyDeployment.ObjectMeta, gomega.Equal(1))

	// modify the number of replicas in the GatewayParameters
	s.patchGatewayParameters(gwParamsDefault.ObjectMeta, func(parameters *v1alpha1.GatewayParameters) {
		parameters.Spec.Kube.Deployment.Replicas = ptr.To(uint32(2))
	})

	// the GatewayParameters modification should cause the deployer to re-run and update the
	// deployment to have 2 replicas
	s.testInstallation.Assertions.EventuallyRunningReplicas(s.ctx, proxyDeployment.ObjectMeta, gomega.Equal(2))

	// modify the external traffic policy in the GatewayParameters
	s.patchGatewayParameters(gwParamsDefault.ObjectMeta, func(parameters *v1alpha1.GatewayParameters) {
		parameters.Spec.Kube.Service.ExternalTrafficPolicy = ptr.To(corev1.ServiceExternalTrafficPolicyLocal)
	})

	// the GatewayParameters modification should cause the deployer to re-run and update the
	// service to have ExternalTrafficPolicy = Local
	s.testInstallation.Assertions.EventuallyExternalTrafficPolicy(s.ctx, *proxyService, gomega.Equal(corev1.ServiceExternalTrafficPolicyLocal))
}

func (s *testingSuite) TestProvisionResourcesNotUpdatedWithInvalidParameters() {
	s.testInstallation.Assertions.EventuallyRunningReplicas(s.ctx, proxyDeployment.ObjectMeta, gomega.Equal(1))

	var (
		// initially, allowPrivilegeEscalation should be true and privileged should not be set
		origAllowPrivilegeEscalation = gstruct.PointTo(gomega.BeTrue())
		origPrivileged               = gomega.BeNil()
	)

	s.patchGatewayParameters(gwParamsDefault.ObjectMeta, func(parameters *v1alpha1.GatewayParameters) {
		gomega.Expect(proxyDeployment.Spec.Template.Spec.Containers).To(gomega.HaveLen(1))
		envoyContainer := proxyDeployment.Spec.Template.Spec.Containers[0]
		gomega.Expect(envoyContainer.SecurityContext.AllowPrivilegeEscalation).To(origAllowPrivilegeEscalation)
		gomega.Expect(envoyContainer.SecurityContext.Privileged).To(origPrivileged)

		// try to modify GatewayParameters with invalid values
		// K8s won't allow setting both allowPrivilegeEscalation=false and privileged=true,
		// so the proposed patch should fail and the original values should be retained.
		parameters.Spec.Kube.EnvoyContainer = &v1alpha1.EnvoyContainer{
			SecurityContext: &corev1.SecurityContext{
				Privileged:               ptr.To(true),
				AllowPrivilegeEscalation: ptr.To(false),
			},
		}

		// This is valid, but should be ignored, because another part of this patch is invalid
		parameters.Spec.Kube.Deployment.Replicas = ptr.To(uint32(2))
	})

	// We keep checking for some amount of time (30s) to account for the time it might take for
	// the deployer to run and re-provision resources. If the original values are consistently
	// retained after that amount of time, we can be confident that the deployer has had time to
	// consume the new values and fail to apply them.
	s.testInstallation.Assertions.Gomega.Consistently(func(g gomega.Gomega) {
		err := s.testInstallation.ClusterContext.Client.Get(s.ctx, client.ObjectKeyFromObject(proxyDeployment), proxyDeployment)
		g.Expect(err).NotTo(gomega.HaveOccurred())
		g.Expect(proxyDeployment.Spec.Template.Spec.Containers[0].SecurityContext.AllowPrivilegeEscalation).To(origAllowPrivilegeEscalation)
		g.Expect(proxyDeployment.Spec.Template.Spec.Containers[0].SecurityContext.Privileged).To(origPrivileged)
		g.Expect(proxyDeployment.Spec.Replicas).To(gstruct.PointTo(gomega.Equal(int32(1))))
	}, "30s", "1s").Should(gomega.Succeed())

}

func (s *testingSuite) TestSelfManagedGateway() {
	s.Require().EventuallyWithT(func(c *assert.CollectT) {
		gw := &gwv1.Gateway{}
		err := s.testInstallation.ClusterContext.Client.Get(s.ctx,
			types.NamespacedName{Name: glooProxyObjectMeta.Name, Namespace: glooProxyObjectMeta.Namespace},
			gw)
		assert.NoError(c, err, "gateway not found")

		accepted := false
		for _, conditions := range gw.Status.Conditions {
			if conditions.Type == string(gwv1.GatewayConditionAccepted) && conditions.Status == metav1.ConditionTrue {
				accepted = true
				break
			}
		}
		assert.True(c, accepted, "gateway status not accepted")
	}, 10*time.Second, 1*time.Second)

	s.testInstallation.Assertions.ConsistentlyObjectsNotExist(s.ctx, proxyService, proxyServiceAccount, proxyDeployment)
}

// patchGateway accepts a reference to an object, and a patch function. It then queries the object,
// performs the patch in memory, and writes the object back to the cluster.
func (s *testingSuite) patchGateway(objectMeta metav1.ObjectMeta, patchFn func(*gwv1.Gateway)) {
	gw := new(gwv1.Gateway)
	gwName := types.NamespacedName{
		Namespace: objectMeta.GetNamespace(),
		Name:      objectMeta.GetName(),
	}
	err := s.testInstallation.ClusterContext.Client.Get(s.ctx, gwName, gw)
	s.Assert().NoError(err, "can get the Gateway object")
	updated := gw.DeepCopy()

	patchFn(updated)

	err = s.testInstallation.ClusterContext.Client.Patch(s.ctx, updated, client.MergeFrom(gw))
	s.Assert().NoError(err, "can update the Gateway object")
}

// patchGatewayParameters accepts a reference to an object, and a patch function
// It then queries the object, performs the patch in memory, and writes the object back to the cluster
func (s *testingSuite) patchGatewayParameters(objectMeta metav1.ObjectMeta, patchFn func(*v1alpha1.GatewayParameters)) {
	gatewayParameters := &v1alpha1.GatewayParameters{}
	err := s.testInstallation.ClusterContext.Client.Get(s.ctx, client.ObjectKey{
		Name:      objectMeta.GetName(),
		Namespace: objectMeta.GetNamespace(),
	}, gatewayParameters)
	s.Assert().NoError(err, "can query the GatewayParameters object")
	modifiedGatewayParameters := gatewayParameters.DeepCopy()

	patchFn(modifiedGatewayParameters)

	err = s.testInstallation.ClusterContext.Client.Patch(s.ctx, modifiedGatewayParameters, client.MergeFrom(gatewayParameters))
	s.Assert().NoError(err, "can update the GatewayParameters object")
}

func serverInfoLogLevelAssertion(testInstallation *e2e.TestInstallation, expectedLogLevel, expectedComponentLogLevel string) func(ctx context.Context, adminClient *admincli.Client) {
	return func(ctx context.Context, adminClient *admincli.Client) {
		testInstallation.Assertions.Gomega.Eventually(func(g gomega.Gomega) {
			serverInfo, err := adminClient.GetServerInfo(ctx)
			g.Expect(err).NotTo(gomega.HaveOccurred(), "can get server info")
			g.Expect(serverInfo.GetCommandLineOptions().GetLogLevel()).To(
				gomega.Equal(expectedLogLevel), "defined on the GatewayParameters CR")
			g.Expect(serverInfo.GetCommandLineOptions().GetComponentLogLevel()).To(
				gomega.Equal(expectedComponentLogLevel), "defined on the GatewayParameters CR")
		}).
			WithContext(ctx).
			WithTimeout(time.Second * 10).
			WithPolling(time.Millisecond * 200).
			Should(gomega.Succeed())
	}
}

func xdsClusterAssertion(testInstallation *e2e.TestInstallation) func(ctx context.Context, adminClient *admincli.Client) {
	return func(ctx context.Context, adminClient *admincli.Client) {
		testInstallation.Assertions.Gomega.Eventually(func(g gomega.Gomega) {
			clusters, err := adminClient.GetStaticClusters(ctx)
			g.Expect(err).NotTo(gomega.HaveOccurred(), "can get static clusters from config dump")

			xdsCluster, ok := clusters["xds_cluster"]
			g.Expect(ok).To(gomega.BeTrue(), "xds_cluster in list")

			g.Expect(xdsCluster.GetLoadAssignment().GetEndpoints()).To(gomega.HaveLen(1))
			g.Expect(xdsCluster.GetLoadAssignment().GetEndpoints()[0].GetLbEndpoints()).To(gomega.HaveLen(1))
			xdsSocketAddress := xdsCluster.GetLoadAssignment().GetEndpoints()[0].GetLbEndpoints()[0].GetEndpoint().GetAddress().GetSocketAddress()
			g.Expect(xdsSocketAddress).NotTo(gomega.BeNil())

			g.Expect(xdsSocketAddress.GetAddress()).To(gomega.Equal(kubeutils.ServiceFQDN(metav1.ObjectMeta{
				Name:      kubeutils.GlooServiceName,
				Namespace: testInstallation.Metadata.InstallNamespace,
			})), "xds socket address points to gloo service, in installation namespace")

			xdsPort, err := setup.GetNamespacedControlPlaneXdsPort(ctx, testInstallation.Metadata.InstallNamespace, testInstallation.ResourceClients.ServiceClient())
			g.Expect(err).NotTo(gomega.HaveOccurred())
			g.Expect(xdsSocketAddress.GetPortValue()).To(gomega.Equal(uint32(xdsPort)), "xds socket port points to gloo service, in installation namespace")
		}).
			WithContext(ctx).
			WithTimeout(time.Second * 10).
			WithPolling(time.Millisecond * 200).
			Should(gomega.Succeed())
	}
}

// awsStsClusterAssertion asserts that:
// - if the installation is configured to use aws with service account creds, then the proxy contains an aws sts cluster with the expected sts uri
// - if the installation is NOT configured to use aws with service account creds, then no aws sts cluster should exist on the proxy
func awsStsClusterAssertion(testInstallation *e2e.TestInstallation) func(ctx context.Context, adminClient *admincli.Client) {
	// get aws values from installation
	awsOpts := testInstallation.Metadata.AwsOptions
	shouldHaveStsCluster := awsOpts != nil && awsOpts.EnableServiceAccountCredentials
	var expectedStsUri string
	if shouldHaveStsCluster {
		expectedStsUri = fmt.Sprintf("sts.%s.amazonaws.com", awsOpts.StsCredentialsRegion)
	}

	return func(ctx context.Context, adminClient *admincli.Client) {
		testInstallation.Assertions.Gomega.Eventually(func(g gomega.Gomega) {
			clusters, err := adminClient.GetStaticClusters(ctx)
			g.Expect(err).NotTo(gomega.HaveOccurred(), "can get static clusters from config dump")

			awsStsCluster, ok := clusters["aws_sts_cluster"]
			if shouldHaveStsCluster {
				g.Expect(ok).To(gomega.BeTrue(), "should contain cluster aws_sts_cluster")
			} else {
				g.Expect(ok).To(gomega.BeFalse(), "should not contain cluster aws_sts_cluster")
				// nothing else to test, so return here
				return
			}

			// check that transport socket has expected values
			msg, err := utils.AnyToMessage(awsStsCluster.GetTransportSocket().GetTypedConfig())
			g.Expect(err).NotTo(gomega.HaveOccurred())
			tlsCtx, ok := msg.(*tlsv3.UpstreamTlsContext)
			g.Expect(ok).To(gomega.BeTrue(), "should be able to get UpstreamTlsContext")
			g.Expect(tlsCtx.GetSni()).To(gomega.Equal(expectedStsUri))

			// check that load assignment has expected values
			g.Expect(awsStsCluster.GetLoadAssignment().GetEndpoints()).To(gomega.HaveLen(1))
			g.Expect(awsStsCluster.GetLoadAssignment().GetEndpoints()[0].GetLbEndpoints()).To(gomega.HaveLen(1))
			socketAddr := awsStsCluster.GetLoadAssignment().GetEndpoints()[0].GetLbEndpoints()[0].GetEndpoint().GetAddress().GetSocketAddress()
			g.Expect(socketAddr).NotTo(gomega.BeNil())
			g.Expect(socketAddr.GetAddress()).To(gomega.Equal(expectedStsUri))
		}).
			WithContext(ctx).
			WithTimeout(time.Second * 10).
			WithPolling(time.Millisecond * 200).
			Should(gomega.Succeed())
	}
}
