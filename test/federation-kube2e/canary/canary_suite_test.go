package canary_test

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/golang/protobuf/ptypes/wrappers"
	skv2v1 "github.com/solo-io/skv2/pkg/api/core.skv2.solo.io/v1"
	multicluster_v1alpha1 "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/multicluster.solo.io/v1alpha1"
	multicluster_types "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/multicluster.solo.io/v1alpha1/types"
	"github.com/solo-io/solo-projects/test/kubeutils"
	"github.com/solo-io/solo-projects/test/kubeutils/installer"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	skhelpers "github.com/solo-io/solo-kit/test/helpers"

	"github.com/solo-io/gloo/projects/gloo/cli/pkg/testutils"
	"github.com/solo-io/gloo/projects/gloo/pkg/defaults"
	"github.com/solo-io/solo-projects/test/services"

	"os"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/solo-io/go-utils/log"
)

const (
	managementClusterEnvName    = "MANAGEMENT_CLUSTER"
	remoteReleaseClusterEnvName = "REMOTE_RELEASE_CLUSTER"
	remoteCanaryClusterEnvName  = "REMOTE_CANARY_CLUSTER"

	releaseNamespace    = "gloo-system-release"
	canaryNamespace     = "gloo-system-canary"
	federationNamespace = "gloo-fed"
	remoteClusterName   = "remote-cluster"
)

func TestCanary(t *testing.T) {
	if !kubeutils.IsKubeTestType("canary") {
		log.Warnf("This test is disabled. To enable, set KUBE2E_TESTS to 'canary' in your env.")
		return
	}

	requiredEnvForTest := []string{
		kubeutils.GlooLicenseKey,
		managementClusterEnvName,
		remoteReleaseClusterEnvName,
		remoteCanaryClusterEnvName,
	}
	if !kubeutils.IsEnvDefined(requiredEnvForTest) {
		log.Warnf("This test is disabled. To enable, set %v in your env.", requiredEnvForTest)
		return
	}
	RegisterFailHandler(Fail)

	RunSpecs(t, "Federation Canary Suite")
}

var (
	ctx    context.Context
	cancel context.CancelFunc

	resourcesFolder string

	managementClusterConfig    *kubeutils.ClusterConfig
	remoteReleaseClusterConfig *kubeutils.ClusterConfig
	remoteCanaryClusterConfig  *kubeutils.ClusterConfig

	orchestrator          kubeutils.Orchestrator
	multiClusterInstaller *installer.MultiClusterInstaller
)

var _ = SynchronizedBeforeSuite(func() []byte {
	managementClusterConfig = kubeutils.CreateClusterConfigFromKubeClusterNameEnv(managementClusterEnvName)
	remoteReleaseClusterConfig = kubeutils.CreateClusterConfigFromKubeClusterNameEnv(remoteReleaseClusterEnvName)
	remoteCanaryClusterConfig = kubeutils.CreateClusterConfigFromKubeClusterNameEnv(remoteCanaryClusterEnvName)

	ctx, cancel = context.WithCancel(context.Background())

	cwd, err := os.Getwd()
	Expect(err).NotTo(HaveOccurred())
	resourcesFolder = filepath.Join(cwd, "resources")

	// Configure the MultiClusterInstaller Installer with the set of actions required at install time
	orchestrator = kubeutils.NewKindOrchestrator()
	multiClusterInstaller = installer.NewMultiCluster(orchestrator)
	multiClusterInstaller.RegisterInstaller(getManagementReleaseInstaller())
	multiClusterInstaller.RegisterInstaller(getManagementCanaryInstaller())
	multiClusterInstaller.RegisterInstaller(getRemoteReleaseInstaller())
	multiClusterInstaller.RegisterInstaller(getRemoteCanaryInstaller())

	// Configure Pre-Fail Handlers so that we output debug information on fails
	skhelpers.RegisterCommonFailHandlers()
	skhelpers.RegisterPreFailHandler(
		kubeutils.GetClusteredPreFailHandler(ctx, orchestrator, GinkgoWriter, []kubeutils.InstallRef{
			{
				ClusterName: managementClusterConfig.ClusterName,
				Namespace:   releaseNamespace,
			},
			{
				ClusterName: remoteReleaseClusterConfig.ClusterName,
				Namespace:   defaults.GlooSystem,
			},
		}))

	By("Install Gloo Edge in management and remote clusters")
	installErr := multiClusterInstaller.Install(ctx)
	Expect(installErr).NotTo(HaveOccurred())

	By("Register remote clusters")
	err = orchestrator.SetClusterContext(ctx, managementClusterConfig.ClusterName)
	Expect(err).NotTo(HaveOccurred())

	registerErr := registerCluster(remoteReleaseClusterConfig, releaseNamespace)
	Expect(registerErr).NotTo(HaveOccurred())

	registerErr = registerCluster(remoteCanaryClusterConfig, canaryNamespace)
	Expect(registerErr).NotTo(HaveOccurred())

	By("Apply multicluster role bindings")
	createErr := createClusterRoleBinding(ctx, managementClusterConfig, []string{
		releaseNamespace,
		canaryNamespace,
	})
	Expect(createErr).NotTo(HaveOccurred())

	By("Create gloo-fed namespace")
	err = services.Kubectl("create", "namespace", federationNamespace)
	Expect(err).NotTo(HaveOccurred())

	return nil
}, func([]byte) {})

var _ = SynchronizedAfterSuite(
	func() {
		// This runs on all nodes
	},
	func() {
		// This runs only on 1 node

		if !kubeutils.ShouldTearDown() {
			return
		}

		err := orchestrator.SetClusterContext(ctx, managementClusterConfig.ClusterName)
		Expect(err).NotTo(HaveOccurred())

		By("Uninstall Gloo Edge from management and remote clusters")
		uninstallErr := multiClusterInstaller.Uninstall(ctx)
		Expect(uninstallErr).NotTo(HaveOccurred())

		cancel()
	},
)

func getManagementReleaseInstaller() installer.Installer {
	managementReleaseInstaller, err := installer.NewGlooInstaller(installer.InstallConfig{
		ClusterName:      managementClusterConfig.ClusterName,
		InstallNamespace: releaseNamespace,
		HelmValuesFile:   filepath.Join(resourcesFolder, "management-cluster-values.yaml"),
	})
	ExpectWithOffset(1, err).NotTo(HaveOccurred())

	return managementReleaseInstaller
}

func getManagementCanaryInstaller() installer.Installer {
	managementReleaseInstaller, err := installer.NewGlooInstaller(installer.InstallConfig{
		ClusterName:      managementClusterConfig.ClusterName,
		InstallNamespace: canaryNamespace,
		HelmValuesFile:   filepath.Join(resourcesFolder, "management-cluster-values.yaml"),
	})
	ExpectWithOffset(1, err).NotTo(HaveOccurred())

	return managementReleaseInstaller
}

func getRemoteReleaseInstaller() installer.Installer {
	remoteReleaseInstaller, err := installer.NewGlooInstaller(installer.InstallConfig{
		ClusterName:      remoteReleaseClusterConfig.ClusterName,
		InstallNamespace: defaults.GlooSystem,
		HelmValuesFile:   filepath.Join(resourcesFolder, "remote-cluster-values.yaml"),
	})
	ExpectWithOffset(1, err).NotTo(HaveOccurred())

	return remoteReleaseInstaller
}

func getRemoteCanaryInstaller() installer.Installer {
	remoteReleaseInstaller, err := installer.NewGlooInstaller(installer.InstallConfig{
		ClusterName:      remoteCanaryClusterConfig.ClusterName,
		InstallNamespace: defaults.GlooSystem,
		HelmValuesFile:   filepath.Join(resourcesFolder, "remote-cluster-values.yaml"),
	})
	ExpectWithOffset(1, err).NotTo(HaveOccurred())

	return remoteReleaseInstaller
}

func registerCluster(clusterConfig *kubeutils.ClusterConfig, ns string) error {
	args := []string{
		"cluster",
		"register",
		"--federation-namespace", // federation-namespace is the ns that gloo-fed is installed in
		ns,
		"--remote-namespace", // remote-namespace is the ns that the remote gloo-ee is installed in
		defaults.GlooSystem,
		"--cluster-name", // cluster-name is the user-provided name for the KubernetesCluster CR that will be created
		remoteClusterName,
		"--remote-context",
		clusterConfig.KubeContext,
		" --local-cluster-domain-override",
		getApiServerDomainForCluster(clusterConfig),
	}
	return testutils.Glooctl(strings.Join(args, " "))
}

func createClusterRoleBinding(ctx context.Context, clusterConfig *kubeutils.ClusterConfig, namespaces []string) error {

	for _, namespace := range namespaces {
		clusterRoleBinding := &multicluster_v1alpha1.MultiClusterRoleBinding{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "kind-admin",
				Namespace: namespace,
			},
			Spec: multicluster_types.MultiClusterRoleBindingSpec{
				Subjects: []*skv2v1.TypedObjectRef{
					{
						Kind: &wrappers.StringValue{
							Value: "User",
						},
						Name: "kubernetes-admin",
					},
				},
				RoleRef: &skv2v1.ObjectRef{
					Name:      "gloo-fed",
					Namespace: namespace,
				},
			},
		}
		if err := clusterConfig.MulticlusterClientset.MultiClusterRoleBindings().CreateMultiClusterRoleBinding(ctx, clusterRoleBinding); err != nil {
			return err
		}
	}
	return nil
}

func getApiServerDomainForCluster(clusterConfig *kubeutils.ClusterConfig) string {
	apiServerPort := 6443 // The default value used for Kind clusters
	return fmt.Sprintf("%s-control-plane:%d", clusterConfig.ClusterName, apiServerPort)
}
