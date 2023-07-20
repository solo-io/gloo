package kube2e

import (
	"context"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"

	"github.com/solo-io/gloo/test/testutils"

	"github.com/solo-io/gloo/projects/gloo/pkg/plugins/kubernetes"

	"github.com/solo-io/gloo/projects/gateway/pkg/defaults"

	"github.com/golang/protobuf/ptypes/wrappers"
	v1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/gloosnapshot"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"

	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"

	. "github.com/onsi/gomega"
	"github.com/solo-io/gloo/test/helpers"
	"github.com/solo-io/gloo/test/kube2e"
	"github.com/solo-io/k8s-utils/testutils/helper"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	// uniqueTestResourceLabel is assigned to the default VirtualService used by kube2e tests
	// This unique label per test run ensures that the generated snapshot is different on subsequent runs
	// We have previously seen flakes where a resource is deleted and re-created with the same hash and thus
	// the emitter can miss the update
	uniqueTestResourceLabel = "gloo-kube2e-test-id"
)

type TestContextFactory struct {
	TestHelper        *helper.SoloTestHelper
	snapshotWriter    helpers.SnapshotWriter
	resourceClientSet *kube2e.KubeResourceClientSet
}

func (f *TestContextFactory) SetupSnapshotAndClientset(ctx context.Context) {
	resourceClientSet, err := kube2e.NewDefaultKubeResourceClientSet(ctx)
	ExpectWithOffset(1, err).NotTo(HaveOccurred(), "can create kube resource client set")

	snapshotWriter := helpers.NewSnapshotWriter(resourceClientSet).WithWriteNamespace(f.TestHelper.InstallNamespace)
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
// It uses the helm values stored in each test's `artifacts/helm.yaml` file to configure Gloo.
func (f *TestContextFactory) InstallGloo(ctx context.Context) {
	if testutils.ShouldSkipInstall() {
		return
	}

	cwd, err := os.Getwd()
	ExpectWithOffset(1, err).NotTo(HaveOccurred(), "working dir should be retrieved while installing gloo")
	helmValuesFile := filepath.Join(cwd, "artifacts", "helm.yaml")

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
		WithDomain(defaults.GatewayProxyName).
		WithRoutePrefixMatcher(DefaultRouteName, TestMatcherPrefix).
		WithRouteOptions(DefaultRouteName, &gloov1.RouteOptions{
			PrefixRewrite: &wrappers.StringValue{Value: "/"},
		}).
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
	t.cancel()
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
	// Wait for a proxy to be accepted
	helpers.EventuallyResourceAccepted(func() (resources.InputResource, error) {
		return t.resourceClientSet.ProxyClient().Read(t.InstallNamespace(), defaults.GatewayProxyName, clients.ReadOpts{Ctx: t.ctx})
	})
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
		t.resourceClientSet.VirtualServiceClient().BaseClient(),
	)
	ExpectWithOffset(1, err).NotTo(HaveOccurred())
}

func (t *TestContext) ResourceClientSet() *kube2e.KubeResourceClientSet {
	return t.resourceClientSet
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
			Host:              defaults.GatewayProxyName,
			Service:           defaults.GatewayProxyName,
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
