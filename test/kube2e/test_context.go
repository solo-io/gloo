package kube2e

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/avast/retry-go"
	"github.com/google/uuid"
	. "github.com/onsi/gomega"
	v1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	gatewaydefaults "github.com/solo-io/gloo/projects/gateway/pkg/defaults"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/gloosnapshot"
	gloodefaults "github.com/solo-io/gloo/projects/gloo/pkg/defaults"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins/kubernetes"
	"github.com/solo-io/gloo/test/helpers"
	"github.com/solo-io/gloo/test/kube2e"
	"github.com/solo-io/gloo/test/testutils"
	"github.com/solo-io/k8s-utils/testutils/helper"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-projects/test/kubeutils"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	// uniqueTestResourceLabel is assigned to the default VirtualService used by kube2e tests
	// This unique label per test run ensures that the generated snapshot is different on subsequent runs
	// We have previously seen flakes where a resource is deleted and re-created with the same hash and thus
	// the emitter can miss the update
	uniqueTestResourceLabel = "gloo-kube2e-test-id"

	DefaultSettingsName = gloodefaults.SettingsName
)

type TestContextFactory struct {
	TestHelper        *helper.SoloTestHelper
	snapshotWriter    helpers.SnapshotWriter
	resourceClientSet *kube2e.KubeResourceClientSet
}

func (f *TestContextFactory) SetupSnapshotAndClientSet(ctx context.Context) {
	resourceClientSet, err := kube2e.NewDefaultKubeResourceClientSet(ctx)
	ExpectWithOffset(1, err).NotTo(HaveOccurred(), "can create kube resource client set")

	snapshotWriter := helpers.NewSnapshotWriter(resourceClientSet).
		WithWriteNamespace(f.InstallNamespace()).
		WithRetryOptions([]retry.Option{
			// Some tests modify multiple resources so giving a lenient retry amount when writing snapshots
			retry.Attempts(6),
		})
	f.snapshotWriter = snapshotWriter
	f.resourceClientSet = resourceClientSet
}

func (f *TestContextFactory) NewTestContext() *TestContext {
	return &TestContext{
		testHelper:        f.TestHelper,
		snapshotWriter:    f.snapshotWriter,
		resourceClientSet: f.resourceClientSet,
	}
}

// InstallGloo installs Gloo if the "SKIP_INSTALL" environment variable is not true.
func (f *TestContextFactory) InstallGloo(ctx context.Context, valuesFileName string) {
	if testutils.ShouldSkipInstall() {
		return
	}

	cwd, err := os.Getwd()
	ExpectWithOffset(1, err).NotTo(HaveOccurred(), "working dir should be retrieved while installing gloo")
	helmValuesFile := filepath.Join(cwd, "artifacts", valuesFileName)

	err = f.TestHelper.InstallGloo(ctx, helper.GATEWAY, 5*time.Minute, helper.ExtraArgs("--values", helmValuesFile))
	ExpectWithOffset(1, err).NotTo(HaveOccurred(), "expect gloo to install successfully")

	f.waitForGlooHealthy()
}

func (f *TestContextFactory) waitForGlooHealthy() {
	// Check that everything is OK
	kube2e.GlooctlCheckEventuallyHealthy(1, f.TestHelper, "90s")

	// Ensure gloo reaches valid state and doesn't continually resync
	kube2e.EventuallyReachesConsistentState(f.TestHelper.InstallNamespace)
}

func (f *TestContextFactory) UninstallGloo(ctx context.Context) {
	if !testutils.ShouldTearDown() {
		return
	}

	err := f.TestHelper.UninstallGlooAll()
	ExpectWithOffset(1, err).NotTo(HaveOccurred())
	_, err = kube2e.MustKubeClient().CoreV1().Namespaces().Get(ctx, f.TestHelper.InstallNamespace, metav1.GetOptions{})
	ExpectWithOffset(1, apierrors.IsNotFound(err)).To(BeTrue())
}

func (f *TestContextFactory) InstallNamespace() string {
	return f.TestHelper.InstallNamespace
}

type TestContext struct {
	ctx               context.Context
	cancel            context.CancelFunc
	testHelper        *helper.SoloTestHelper
	resourceClientSet *kube2e.KubeResourceClientSet
	snapshotWriter    helpers.SnapshotWriter
	resourcesToWrite  *gloosnapshot.ApiSnapshot
}

const (
	DefaultRouteName          = "testrunner-route"
	DefaultVirtualServiceName = "vs"
)

func (t *TestContext) BeforeEach() {
	t.ctx, t.cancel = context.WithCancel(context.Background())

	defaultVs := helpers.NewVirtualServiceBuilder().
		WithNamespace(t.InstallNamespace()).
		WithName(DefaultVirtualServiceName).
		WithDomain(gatewaydefaults.GatewayProxyName).
		WithRoutePrefixMatcher(DefaultRouteName, TestMatcherPrefix).
		WithRouteActionToUpstreamRef(DefaultRouteName, t.TestRunnerUpstreamRef()).
		Build()
	defaultVs.Metadata.Labels = map[string]string{
		uniqueTestResourceLabel: uuid.New().String(),
	}

	t.resourcesToWrite = &gloosnapshot.ApiSnapshot{
		VirtualServices: v1.VirtualServiceList{defaultVs},
	}
}

func (t *TestContext) JustBeforeEach() {
	err := t.snapshotWriter.WriteSnapshot(t.resourcesToWrite, clients.WriteOpts{
		Ctx:               t.Ctx(),
		OverwriteExisting: false,
	})
	ExpectWithOffset(1, err).NotTo(HaveOccurred())

	t.EventuallyProxyAccepted()
}

func (t *TestContext) AfterEach() {
	t.CancelContext()
}

func (t *TestContext) JustAfterEach() {
	err := t.snapshotWriter.DeleteSnapshot(t.resourcesToWrite, clients.DeleteOpts{
		Ctx:            t.Ctx(),
		IgnoreNotExist: true,
	})
	ExpectWithOffset(1, err).NotTo(HaveOccurred())
}

// EventuallyProxyAccepted is useful for tests that rely on changing an existing configuration.
func (t *TestContext) EventuallyProxyAccepted() {
	// Wait for a proxy to be accepted.
	// Some tests (gateway_test) configure multiple resources which locally have taken ~80 seconds to be reconciled.
	helpers.EventuallyResourceAccepted(func() (resources.InputResource, error) {
		return t.ResourceClientSet().ProxyClient().Read(t.InstallNamespace(), gatewaydefaults.GatewayProxyName, clients.ReadOpts{Ctx: t.ctx})
	}, 2*time.Minute)
}

// PatchDefaultVirtualService mutates the existing VirtualService generated by the TestContext
func (t *TestContext) PatchDefaultVirtualService(mutator func(*v1.VirtualService) *v1.VirtualService) {
	err := helpers.PatchResourceWithOffset(
		1,
		t.Ctx(),
		&core.ResourceRef{
			Name:      DefaultVirtualServiceName,
			Namespace: t.InstallNamespace(),
		},
		func(resource resources.Resource) resources.Resource {
			return mutator(resource.(*v1.VirtualService))
		},
		t.ResourceClientSet().VirtualServiceClient().BaseClient(),
	)
	ExpectWithOffset(1, err).NotTo(HaveOccurred())
}

// PatchDefaultGateway mutates the existing Gateway generated by the TestContext
func (t *TestContext) PatchDefaultGateway(mutator func(*v1.Gateway) *v1.Gateway) {
	err := helpers.PatchResourceWithOffset(
		1,
		t.Ctx(),
		&core.ResourceRef{
			Name:      gatewaydefaults.GatewayProxyName,
			Namespace: t.InstallNamespace(),
		},
		func(resource resources.Resource) resources.Resource {
			return mutator(resource.(*v1.Gateway))
		},
		t.ResourceClientSet().GatewayClient().BaseClient(),
	)
	ExpectWithOffset(1, err).NotTo(HaveOccurred())
}

// PatchDefaultSettings mutates the default Settings
func (t *TestContext) PatchDefaultSettings(mutator func(*gloov1.Settings) *gloov1.Settings) {
	err := helpers.PatchResourceWithOffset(
		1,
		t.Ctx(),
		&core.ResourceRef{
			Name:      DefaultSettingsName,
			Namespace: t.InstallNamespace(),
		},
		func(resource resources.Resource) resources.Resource {
			return mutator(resource.(*gloov1.Settings))
		},
		t.ResourceClientSet().SettingsClient().BaseClient(),
	)
	ExpectWithOffset(1, err).NotTo(HaveOccurred())
}

func (t *TestContext) GetDefaultSettings() *gloov1.Settings {
	settings, err := t.resourceClientSet.SettingsClient().Read(t.InstallNamespace(), DefaultSettingsName, clients.ReadOpts{Ctx: t.ctx})
	ExpectWithOffset(1, err).NotTo(HaveOccurred(), "read default settings")
	return settings
}

// ModifyDeploymentReplicas sets the given deployment's number of replicas, and then waits for the deployment
// to be rolled out successfully (or times out).
func (t *TestContext) ModifyDeploymentReplicas(deploymentName string, replicas int32) {
	deploymentClient := t.ResourceClientSet().KubeClients().AppsV1().Deployments(t.InstallNamespace())

	d, err := deploymentClient.Get(t.Ctx(), deploymentName, metav1.GetOptions{})
	ExpectWithOffset(1, err).NotTo(HaveOccurred())

	d.Spec.Replicas = &replicas
	_, err = deploymentClient.Update(t.Ctx(), d, metav1.UpdateOptions{})
	ExpectWithOffset(1, err).NotTo(HaveOccurred())
	fmt.Printf("modified %s replicas to %v, now waiting for rollout\n", deploymentName, replicas)

	kubeutils.WaitForRolloutWithOffset(1, deploymentName, t.InstallNamespace(), "60s", "1s")
}

// WaitForDeploymentReplicas waits for the given deployment's replicas to equal the given number of replicas
// (or times out).
func (t *TestContext) WaitForDeploymentReplicas(deploymentName string, replicas int32) {
	EventuallyWithOffset(1, func() (int32, error) {
		deploymentClient := t.ResourceClientSet().KubeClients().AppsV1().Deployments(t.InstallNamespace())

		d, err := deploymentClient.Get(t.Ctx(), deploymentName, metav1.GetOptions{})
		if err != nil {
			return 0, err
		}

		if d.Status.UpdatedReplicas != d.Status.ReadyReplicas {
			return 0, errors.New("updated and ready replicas do not match")
		}

		return d.Status.ReadyReplicas, nil
	}, "60s", "1s").Should(Equal(replicas))
}

// ModifyDeploymentEnv adds (if it does not exist) or modifies (if it already exists) the specified env variable
// on the specified container of the deployment's pod spec, and then waits for the deployment to be rolled out
// successfully (or times out).
func (t *TestContext) ModifyDeploymentEnv(deploymentName string, containerIndex int, envVar corev1.EnvVar) {
	deploymentClient := t.ResourceClientSet().KubeClients().AppsV1().Deployments(t.InstallNamespace())

	d, err := deploymentClient.Get(t.Ctx(), deploymentName, metav1.GetOptions{})
	ExpectWithOffset(1, err).NotTo(HaveOccurred())

	// make sure we are referencing a valid container
	ExpectWithOffset(1, len(d.Spec.Template.Spec.Containers)).To(BeNumerically(">", containerIndex))

	// if an env var with the given name already exists, modify it
	exists := false
	for i, env := range d.Spec.Template.Spec.Containers[containerIndex].Env {
		if env.Name == envVar.Name {
			d.Spec.Template.Spec.Containers[containerIndex].Env[i].Value = envVar.Value
			exists = true
			break
		}
	}
	// otherwise add a new env var
	if !exists {
		d.Spec.Template.Spec.Containers[containerIndex].Env = append(d.Spec.Template.Spec.Containers[containerIndex].Env, envVar)
	}
	_, err = deploymentClient.Update(t.Ctx(), d, metav1.UpdateOptions{})
	ExpectWithOffset(1, err).NotTo(HaveOccurred())

	kubeutils.WaitForRolloutWithOffset(1, deploymentName, t.InstallNamespace(), "60s", "1s")
}

func (t *TestContext) ResourceClientSet() *kube2e.KubeResourceClientSet {
	return t.resourceClientSet
}

func (t *TestContext) SnapshotWriter() helpers.SnapshotWriter {
	return t.snapshotWriter
}

func (t *TestContext) ResourcesToWrite() *gloosnapshot.ApiSnapshot {
	return t.resourcesToWrite
}

func (t *TestContext) TestHelper() *helper.SoloTestHelper {
	return t.testHelper
}

func (t *TestContext) InstallNamespace() string {
	return t.testHelper.InstallNamespace
}

func (t *TestContext) Ctx() context.Context {
	return t.ctx
}

func (t *TestContext) CancelContext() {
	t.cancel()
}

func (t *TestContext) TestRunnerUpstreamRef() *core.ResourceRef {
	return &core.ResourceRef{
		Namespace: t.InstallNamespace(),
		Name:      kubernetes.UpstreamName(t.InstallNamespace(), helper.TestrunnerName, helper.TestRunnerPort),
	}
}

type CurlOptsBuilder struct {
	opts helper.CurlOpts
}

func (t *TestContext) DefaultCurlOptsBuilder() *CurlOptsBuilder {
	return &CurlOptsBuilder{
		opts: helper.CurlOpts{
			Protocol:          "http",
			Method:            http.MethodGet,
			Path:              TestMatcherPrefix,
			Port:              80, // Gateway port
			Host:              gatewaydefaults.GatewayProxyName,
			Service:           gatewaydefaults.GatewayProxyName,
			ConnectionTimeout: 1,
			WithoutStats:      true,
			Verbose:           false,
		},
	}
}

func (c *CurlOptsBuilder) Build() helper.CurlOpts {
	return c.opts
}

func (c *CurlOptsBuilder) WithProtocol(protocol string) *CurlOptsBuilder {
	c.opts.Protocol = protocol
	return c
}

func (c *CurlOptsBuilder) WithCaFile(caFile string) *CurlOptsBuilder {
	c.opts.CaFile = caFile
	return c
}

func (c *CurlOptsBuilder) WithHeaders(headers map[string]string) *CurlOptsBuilder {
	c.opts.Headers = headers
	return c
}

func (c *CurlOptsBuilder) WithBody(body string) *CurlOptsBuilder {
	c.opts.Body = body
	return c
}

func (c *CurlOptsBuilder) WithMethod(method string) *CurlOptsBuilder {
	c.opts.Method = method
	return c
}

func (c *CurlOptsBuilder) WithPath(path string) *CurlOptsBuilder {
	c.opts.Path = path
	return c
}

func (c *CurlOptsBuilder) WithPort(port int) *CurlOptsBuilder {
	c.opts.Port = port
	return c
}

func (c *CurlOptsBuilder) WithHost(host string) *CurlOptsBuilder {
	c.opts.Host = host
	return c
}

func (c *CurlOptsBuilder) WithService(service string) *CurlOptsBuilder {
	c.opts.Service = service
	return c
}

func (c *CurlOptsBuilder) WithConnectionTimeout(timeout int) *CurlOptsBuilder {
	c.opts.ConnectionTimeout = timeout
	return c
}

func (c *CurlOptsBuilder) WithStats(withStats bool) *CurlOptsBuilder {
	c.opts.WithoutStats = !withStats
	return c
}

func (c *CurlOptsBuilder) WithVerbose(verbose bool) *CurlOptsBuilder {
	c.opts.Verbose = verbose
	return c
}
