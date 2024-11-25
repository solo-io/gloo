package e2e

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/solo-io/gloo/pkg/utils/glooadminutils/admincli"
	"github.com/solo-io/gloo/pkg/utils/kubeutils/portforward"
	"github.com/solo-io/gloo/pkg/utils/requestutils/curl"
	"github.com/solo-io/gloo/projects/gloo/pkg/servers/admin"
	"github.com/solo-io/gloo/test/kubernetes/testutils/actions"
	"github.com/solo-io/gloo/test/kubernetes/testutils/assertions"
	"github.com/solo-io/gloo/test/kubernetes/testutils/cluster"
	"github.com/solo-io/gloo/test/kubernetes/testutils/gloogateway"
	"github.com/solo-io/gloo/test/kubernetes/testutils/helper"
	testruntime "github.com/solo-io/gloo/test/kubernetes/testutils/runtime"
	"github.com/solo-io/gloo/test/testutils"
)

// MustTestHelper returns the SoloTestHelper used for e2e tests
// The SoloTestHelper is a wrapper around `glooctl` and we should eventually phase it out
// in favor of using the exact tool that users rely on
func MustTestHelper(ctx context.Context, installation *TestInstallation) *helper.SoloTestHelper {
	testHelper, err := helper.GetTestHelperForRootDir(ctx, testutils.GitRootDirectory(), installation.Metadata.InstallNamespace)
	if err != nil {
		panic(err)
	}

	testHelper.SetKubeCli(installation.ClusterContext.Cli)

	// TODO(npolshak) We should remove the test helper and have the test installation own this
	// set installation metadata
	installation.Metadata.TestAssetDir = testHelper.TestAssetDir
	installation.Metadata.ChartVersion = testHelper.ChartVersion()
	installation.Metadata.ReleasedVersion = testHelper.ReleasedVersion
	installation.Metadata.HelmChartName = testHelper.HelmChartName
	installation.Metadata.HelmRepoIndexFileName = testHelper.HelmRepoIndexFileName
	installation.Metadata.ChartUri = filepath.Join(testutils.GitRootDirectory(), installation.Metadata.TestAssetDir, installation.Metadata.HelmChartName+"-"+installation.Metadata.ChartVersion+".tgz")

	return testHelper
}

// CreateTestInstallation is the simplest way to construct a TestInstallation in Gloo Gateway OSS
// It is syntactic sugar on top of CreateTestInstallationForCluster
func CreateTestInstallation(
	t *testing.T,
	glooGatewayContext *gloogateway.Context,
) *TestInstallation {
	runtimeContext := testruntime.NewContext()
	clusterContext := cluster.MustKindContext(runtimeContext.ClusterName)

	if err := gloogateway.ValidateGlooGatewayContext(glooGatewayContext); err != nil {
		// We error loudly if the context is misconfigured
		panic(err)
	}

	return CreateTestInstallationForCluster(t, runtimeContext, clusterContext, glooGatewayContext)
}

// CreateTestInstallationForCluster is the standard way to construct a TestInstallation
// It accepts context objects from 3 relevant sources:
//
//	runtime - These are properties that are supplied at runtime and will impact how tests are executed
//	cluster - These are properties that are used to connect to the Kubernetes cluster
//	glooGateway - These are properties that are relevant to how Gloo Gateway will be configured
func CreateTestInstallationForCluster(
	t *testing.T,
	runtimeContext testruntime.Context,
	clusterContext *cluster.Context,
	glooGatewayContext *gloogateway.Context,
) *TestInstallation {
	installation := &TestInstallation{
		// RuntimeContext contains the set of properties that are defined at runtime by whoever is invoking tests
		RuntimeContext: runtimeContext,

		// ClusterContext contains the metadata about the Kubernetes Cluster that is used for this TestCluster
		ClusterContext: clusterContext,

		// Maintain a reference to the Metadata used for this installation
		Metadata: glooGatewayContext,

		// ResourceClients are only available _after_ installing Gloo Gateway
		ResourceClients: nil,

		// Create an actions provider, and point it to the running installation
		Actions: actions.NewActionsProvider().
			WithClusterContext(clusterContext).
			WithGlooGatewayContext(glooGatewayContext),

		// Create an assertions provider, and point it to the running installation
		Assertions: assertions.NewProvider(t).
			WithClusterContext(clusterContext).
			WithGlooGatewayContext(glooGatewayContext),

		// GeneratedFiles contains the unique location where files generated during the execution
		// of tests against this installation will be stored
		// By creating a unique location, per TestInstallation and per Cluster.Name we guarantee isolation
		// between TestInstallation outputs per CI run
		GeneratedFiles: MustGeneratedFiles(glooGatewayContext.InstallNamespace, clusterContext.Name),
	}
	runtime.SetFinalizer(installation, func(i *TestInstallation) { i.finalize() })
	return installation
}

// TestInstallation is the structure around a set of tests that validate behavior for an installation
// of Gloo Gateway.
type TestInstallation struct {
	fmt.Stringer

	// RuntimeContext contains the set of properties that are defined at runtime by whoever is invoking tests
	RuntimeContext testruntime.Context

	// ClusterContext contains the metadata about the Kubernetes Cluster that is used for this TestCluster
	ClusterContext *cluster.Context

	// Metadata contains the properties used to install Gloo Gateway
	Metadata *gloogateway.Context

	// ResourceClients is a set of clients that can manipulate resources owned by Gloo Gateway
	ResourceClients gloogateway.ResourceClients

	// Actions is the entity that creates actions that can be executed by the Operator
	Actions *actions.Provider

	// Assertions is the entity that creates assertions that can be executed by the Operator
	Assertions *assertions.Provider

	// GeneratedFiles is the collection of directories and files that this test installation _may_ create
	GeneratedFiles GeneratedFiles

	// IstioctlBinary is the path to the istioctl binary that can be used to interact with Istio
	IstioctlBinary string
}

func (i *TestInstallation) String() string {
	return i.Metadata.InstallNamespace
}

func (i *TestInstallation) finalize() {
	if err := os.RemoveAll(i.GeneratedFiles.TempDir); err != nil {
		panic(fmt.Sprintf("Failed to remove temporary directory: %s", i.GeneratedFiles.TempDir))
	}
}

func (i *TestInstallation) AddIstioctl(ctx context.Context) error {
	istioctl, err := cluster.GetIstioctl(ctx)
	if err != nil {
		return fmt.Errorf("failed to download istio: %w", err)
	}
	i.IstioctlBinary = istioctl
	return nil
}

func (i *TestInstallation) InstallMinimalIstio(ctx context.Context) error {
	return cluster.InstallMinimalIstio(ctx, i.IstioctlBinary, i.ClusterContext.KubeContext)
}

func (i *TestInstallation) InstallRevisionedIstio(ctx context.Context, rev, profile string) error {
	return cluster.InstallRevisionedIstio(ctx, i.IstioctlBinary, i.ClusterContext.KubeContext, rev, profile)
}

func (i *TestInstallation) UninstallIstio() error {
	return cluster.UninstallIstio(i.IstioctlBinary, i.ClusterContext.KubeContext)
}

func (i *TestInstallation) CreateIstioBugReport(ctx context.Context) {
	cluster.CreateIstioBugReport(ctx, i.IstioctlBinary, i.ClusterContext.KubeContext, i.GeneratedFiles.FailureDir)
}

// InstallGlooGatewayWithTestHelper is the common way to install Gloo Gateway.
// However, it relies on a SoloTestHelper which is an artifact of the legacy e2e tests that we hope to deprecate
func (i *TestInstallation) InstallGlooGatewayWithTestHelper(ctx context.Context, testHelper *helper.SoloTestHelper, timeout time.Duration) {
	installFn := func(ctx context.Context) error {
		return testHelper.InstallGloo(
			ctx,
			timeout,
			// We layer the values that are provided to a given TestInstallation, as a way to reduce the complexity that developers
			// have to concern themselves with.
			//	1. First, we define the common recommendations. These are values which we want all tests to inherit by default.
			//	2. Next, we define a "profile", which is a standard collection of values for a type of installation
			//	3. Finally, we define test specific values, so that a test can have the final say over which values to define.
			helper.WithExtraArgs("--values", CommonRecommendationManifest),
			helper.WithExtraArgs("--values", i.Metadata.ProfileValuesManifestFile),
			helper.WithExtraArgs("--values", i.Metadata.ValuesManifestFile),
		)
	}

	i.InstallGlooGateway(ctx, installFn)
}

func (i *TestInstallation) InstallGlooGateway(ctx context.Context, installFn func(ctx context.Context) error) {
	if !testutils.ShouldSkipInstall() {
		err := installFn(ctx)
		i.Assertions.Require.NoError(err)
		i.Assertions.EventuallyInstallationSucceeded(ctx)
		i.Assertions.EventuallyGlooReachesConsistentState(i.Metadata.InstallNamespace)
	}

	// We can only create the ResourceClients after the CRDs exist in the Cluster
	clients, err := gloogateway.NewResourceClients(ctx, i.ClusterContext)
	i.Assertions.Require.NoError(err)
	i.ResourceClients = clients
}

// UninstallGlooGatewayWithTestHelper is the common way to uninstall Gloo Gateway.
// However, it relies on a SoloTestHelper which is an artifact of the legacy e2e tests that we hope to deprecate
func (i *TestInstallation) UninstallGlooGatewayWithTestHelper(ctx context.Context, testHelper *helper.SoloTestHelper) {
	uninstallFn := func(ctx context.Context) error {
		return testHelper.UninstallGlooAll()
	}

	i.UninstallGlooGateway(ctx, uninstallFn)
}

func (i *TestInstallation) UninstallGlooGateway(ctx context.Context, uninstallFn func(ctx context.Context) error) {
	if testutils.ShouldSkipInstall() {
		return
	}
	err := uninstallFn(ctx)
	i.Assertions.Require.NoError(err)
	i.Assertions.EventuallyUninstallationSucceeded(ctx)
}

// PreFailHandler is the function that is invoked if a test in the given TestInstallation fails
func (i *TestInstallation) PreFailHandler(ctx context.Context) {
	// This is a work in progress
	// The idea here is we want to accumulate ALL information about this TestInstallation into a single directory
	// That way we can upload it in CI, or inspect it locally

	failureDir := i.GeneratedFiles.FailureDir
	err := os.Mkdir(failureDir, os.ModePerm)
	// We don't want to fail on the output directory already existing. This could occur
	// if multiple tests running in the same cluster from the same installation namespace
	// fail.
	if err != nil && !errors.Is(err, fs.ErrExist) {
		i.Assertions.Require.NoError(err)
	}

	glooLogFilePath := filepath.Join(failureDir, "gloo.log")
	glooLogFile, err := os.OpenFile(glooLogFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModePerm)
	i.Assertions.Require.NoError(err)
	defer glooLogFile.Close()

	glooLogsCmd := i.Actions.Kubectl().Command(ctx, "logs", "-n", i.Metadata.InstallNamespace, "deployments/gloo")
	_ = glooLogsCmd.WithStdout(glooLogFile).WithStderr(glooLogFile).Run()

	edgeGatewayLogFilePath := filepath.Join(failureDir, "edge_gateway.log")
	edgeGatewayLogFile, err := os.OpenFile(edgeGatewayLogFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModePerm)
	i.Assertions.Require.NoError(err)
	defer edgeGatewayLogFile.Close()

	kubeGatewayLogFilePath := filepath.Join(failureDir, "kube_gateway.log")
	kubeGatewayLogFile, err := os.OpenFile(kubeGatewayLogFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModePerm)
	i.Assertions.Require.NoError(err)
	defer kubeGatewayLogFile.Close()

	namespaces, err := i.Actions.Kubectl().Namespaces(ctx)
	i.Assertions.Require.NoError(err)
	for _, n := range namespaces {
		edgeGatewayLogFile.WriteString(fmt.Sprintf("Logs for edge gateway proxies in namespace %s\n", n))
		edgeGatewayLogsCmd := i.Actions.Kubectl().Command(ctx, "logs", "--all-containers", "--namespace", n, "--prefix", "-l", "gloo=gateway-proxy")
		_ = edgeGatewayLogsCmd.WithStdout(edgeGatewayLogFile).WithStderr(edgeGatewayLogFile).Run()
		edgeGatewayLogFile.WriteString("----------------------------------------------------------------------------------------------------------\n")

		kubeGatewayLogFile.WriteString(fmt.Sprintf("Logs for kube gateway proxies in namespace %s\n", n))
		kubeGatewayLogsCmd := i.Actions.Kubectl().Command(ctx, "logs", "--all-containers", "--namespace", n, "--prefix", "-l", "gloo=kube-gateway")
		_ = kubeGatewayLogsCmd.WithStdout(kubeGatewayLogFile).WithStderr(kubeGatewayLogFile).Run()
		kubeGatewayLogFile.WriteString("----------------------------------------------------------------------------------------------------------\n")
	}

	clusterStateFilePath := filepath.Join(failureDir, "cluster_state.log")
	clusterStateFile, err := os.OpenFile(clusterStateFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModePerm)
	i.Assertions.Require.NoError(err)
	defer clusterStateFile.Close()

	kubectlGetAllCmd := i.Actions.Kubectl().Command(ctx, "get", "all", "-A", "-owide")
	_ = kubectlGetAllCmd.WithStdout(clusterStateFile).WithStderr(clusterStateFile).Run()
	clusterStateFile.WriteString("\n")

	resourcesToGet := []string{
		// Kubernetes resources
		"secrets",
		// Kube GW API resources
		"gateways.gateway.networking.k8s.io",
		"gatewayclasses.gateway.networking.k8s.io",
		"httproutes.gateway.networking.k8s.io",
		"referencegrants.gateway.networking.k8s.io",
		// GG Kube GW resources
		"gatewayparameters.gateway.gloo.solo.io",
		"listeneroptions.gateway.solo.io",     // only implemented for kube gw as of now
		"httplisteneroptions.gateway.solo.io", // only implemented for kube gw as of now
		// GG Gloo resources
		"graphqlapis.graphql.gloo.solo.io",
		"proxies.gloo.solo.io",
		"settings.gloo.solo.io",
		"upstreamgroups.gloo.solo.io",
		"upstreams.gloo.solo.io",
		// GG Edge GW resources
		"gateways.gateway.solo.io",
		"httpgateways.gateway.solo.io",
		"tcpgateways.gateway.solo.io",
		"virtualservices.gateway.solo.io",
		// Shared GW resources
		"routeoptions.gateway.solo.io",
		"virtualhostoptions.gateway.solo.io",
		// Dataplane extensions resources
		"authconfigs.enterprise.gloo.solo.io",
		"ratelimitconfigs.ratelimit.solo.io",
	}
	kubectlGetResourcesCmd := i.Actions.Kubectl().Command(ctx, "get", strings.Join(resourcesToGet, ","), "-A", "-owide")
	_ = kubectlGetResourcesCmd.WithStdout(clusterStateFile).WithStderr(clusterStateFile).Run()
	clusterStateFile.WriteString("\n")

	podStdOut := bytes.NewBuffer(nil)
	podStdErr := bytes.NewBuffer(nil)

	// Fetch the name of the Gloo Gateway controller pod
	getGlooPodNameCmd := i.Actions.Kubectl().Command(ctx, "get", "pod", "-n", i.Metadata.InstallNamespace,
		"--selector", "gloo=gloo", "--output", "jsonpath='{.items[0].metadata.name}'")
	_ = getGlooPodNameCmd.WithStdout(podStdOut).WithStderr(podStdErr).Run()

	// Clean up and check the output
	glooPodName := strings.Trim(podStdOut.String(), "'")
	if glooPodName == "" {
		fmt.Printf("Failed to get the name of the Gloo Gateway controller pod: %s\n", podStdErr.String())
		return
	}

	// Get the metrics from the Gloo Gateway controller pod and write them to a file
	metricsFilePath := filepath.Join(failureDir, "metrics.log")
	metricsFile, err := os.OpenFile(metricsFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModePerm)
	i.Assertions.Require.NoError(err)

	// Using an ephemeral debug pod fetch the metrics from the Gloo Gateway controller
	metricsCmd := i.Actions.Kubectl().Command(ctx, "debug", "-n", i.Metadata.InstallNamespace,
		"-it", "--image=curlimages/curl:7.83.1", glooPodName, "--",
		"curl", "http://localhost:9091/metrics")
	_ = metricsCmd.WithStdout(metricsFile).WithStderr(metricsFile).Run()
	metricsFile.Close()

	// Open a port-forward to the Gloo Gateway controller pod's admin port
	portForwarder, err := i.Actions.Kubectl().StartPortForward(ctx,
		portforward.WithDeployment("gloo", i.Metadata.InstallNamespace),
		portforward.WithPorts(int(admin.AdminPort), int(admin.AdminPort)),
	)
	i.Assertions.Require.NoError(err)

	defer func() {
		portForwarder.Close()
		portForwarder.WaitForStop()
	}()

	adminClient := admincli.NewClient().
		WithReceiver(io.Discard). // adminAssertion can overwrite this
		WithCurlOptions(
			curl.WithRetries(3, 0, 10),
			curl.WithPort(int(admin.AdminPort)),
		)

	// Get krt snapshot from the Gloo Gateway controller pod and write it to a file
	krtSnapshotFilePath := filepath.Join(failureDir, "krt_snapshot.log")
	krtSnapshotFile, err := os.OpenFile(krtSnapshotFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModePerm)
	i.Assertions.Require.NoError(err)

	cmdErr := adminClient.KrtSnapshotCmd(ctx).
		WithStdout(krtSnapshotFile).
		WithStderr(krtSnapshotFile).
		Run()
	i.Assertions.Require.NoError(cmdErr)

	// Get xds snapshot from the Gloo Gateway controller pod and write it to a file
	xdsSnapshotFilePath := filepath.Join(failureDir, "xds_snapshot.log")
	xdsSnapshotFile, err := os.OpenFile(xdsSnapshotFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModePerm)
	i.Assertions.Require.NoError(err)

	cmdErr = adminClient.XdsSnapshotCmd(ctx).
		WithStdout(xdsSnapshotFile).
		WithStderr(xdsSnapshotFile).
		Run()
	i.Assertions.Require.NoError(cmdErr)

	fmt.Printf("Test failed. Logs and cluster state are available in %s\n", failureDir)
}

// GeneratedFiles is a collection of files that are generated during the execution of a set of tests
type GeneratedFiles struct {
	// TempDir is the directory where any temporary files should be created
	// Tests may create files for any number of reasons:
	// - A: When a test renders objects in a file, and then uses this file to create and delete values
	// - B: When a test invokes a command that produces a file as a side effect (glooctl, for example)
	// Files in this directory are an implementation detail of the test itself.
	// As a result, it is the callers responsibility to clean up the TempDir when the tests complete
	TempDir string

	// FailureDir is the directory where any assets that are produced on failure will be created
	FailureDir string
}

// MustGeneratedFiles returns GeneratedFiles, or panics if there was an error generating the directories
func MustGeneratedFiles(tmpDirId, clusterId string) GeneratedFiles {
	tmpDir, err := os.MkdirTemp("", tmpDirId)
	if err != nil {
		panic(err)
	}

	// output path is in the format of bug_report/cluster_name/tmp_dir_id
	failureDir := filepath.Join(testruntime.PathToBugReport(), clusterId, tmpDirId)
	err = os.MkdirAll(failureDir, os.ModePerm)
	if err != nil {
		panic(err)
	}

	return GeneratedFiles{
		TempDir:    tmpDir,
		FailureDir: failureDir,
	}
}
