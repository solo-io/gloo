package kube2e

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	. "github.com/solo-io/gloo/test/kube2e"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"os/exec"

	"github.com/golang/protobuf/ptypes/wrappers"
	. "github.com/onsi/gomega"
	"github.com/solo-io/gloo/pkg/cliutil"
	v1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/core/matchers"
	gloossl "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/ssl"
	"github.com/solo-io/go-utils/contextutils"
	"github.com/solo-io/k8s-utils/kubeutils"
	"github.com/solo-io/k8s-utils/testutils/helper"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/factory"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/kube"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"go.uber.org/zap"

	kubetestutils "github.com/solo-io/solo-projects/test/kubeutils"
)

const (
	TestMatcherPrefix = "/test"
	GlooeRepoName     = "https://storage.googleapis.com/gloo-ee-helm"
)

func WriteVirtualService(
	ctx context.Context,
	testHelper *helper.SoloTestHelper,
	vsClient v1.VirtualServiceClient,
	virtualHostOptions *gloov1.VirtualHostOptions,
	routeOptions *gloov1.RouteOptions,
	sslConfig *gloossl.SslConfig,
) {

	upstreamRef := &core.ResourceRef{
		Namespace: testHelper.InstallNamespace,
		Name:      fmt.Sprintf("%s-%s-%v", testHelper.InstallNamespace, "testrunner", helper.TestRunnerPort),
	}
	WriteCustomVirtualService(
		ctx,
		2,
		testHelper,
		vsClient,
		virtualHostOptions,
		routeOptions,
		sslConfig,
		upstreamRef,
		TestMatcherPrefix,
	)
}

func WriteCustomVirtualService(
	ctx context.Context,
	offset int,
	testHelper *helper.SoloTestHelper,
	vsClient v1.VirtualServiceClient,
	virtualHostOptions *gloov1.VirtualHostOptions,
	routeOptions *gloov1.RouteOptions,
	sslConfig *gloossl.SslConfig,
	upstreamRef *core.ResourceRef,
	matcherPrefix string,
) {

	if routeOptions.GetPrefixRewrite() == nil {
		if routeOptions == nil {
			routeOptions = &gloov1.RouteOptions{}
		}
		routeOptions.PrefixRewrite = &wrappers.StringValue{
			Value: "/",
		}
	}

	// We wrap this in a eventually because the validating webhook may reject the virtual service if one of the
	// resources the VS depends on is not yet available.
	EventuallyWithOffset(offset, func() error {
		_, err := vsClient.Write(&v1.VirtualService{

			Metadata: &core.Metadata{
				Name:      "vs",
				Namespace: testHelper.InstallNamespace,
			},
			SslConfig: sslConfig,
			VirtualHost: &v1.VirtualHost{
				Options: virtualHostOptions,
				Domains: []string{"*"},
				Routes: []*v1.Route{{
					Options: routeOptions,
					Matchers: []*matchers.Matcher{{
						PathSpecifier: &matchers.Matcher_Prefix{
							Prefix: matcherPrefix,
						},
					}},
					Action: &v1.Route_RouteAction{
						RouteAction: &gloov1.RouteAction{
							Destination: &gloov1.RouteAction_Single{
								Single: &gloov1.Destination{
									DestinationType: &gloov1.Destination_Upstream{
										Upstream: upstreamRef,
									},
								},
							},
						},
					},
				}},
			},
		}, clients.WriteOpts{Ctx: ctx})

		if err != nil {
			contextutils.LoggerFrom(ctx).Warnw("failed to create virtual service", zap.Error(err))
		}

		return err
	}, time.Minute, "5s").Should(BeNil())
}

func DeleteVirtualService(vsClient v1.VirtualServiceClient, ns, name string, opts clients.DeleteOpts) {
	// We wrap this in a eventually because the validating webhook may reject the virtual service if one of the
	// resources the VS depends on is not yet available.
	EventuallyWithOffset(1, func() error {
		return vsClient.Delete(ns, name, opts)
	}, time.Minute, "5s").Should(BeNil())
}

func EnableStrictValidation(testHelper *helper.SoloTestHelper) {
	// enable strict validation
	// this can be removed once we enable validation by default
	// set projects/gateway/pkg/syncer.AcceptAllResourcesByDefault is set to false
	cfg, err := kubeutils.GetConfig("", "")
	Expect(err).NotTo(HaveOccurred())

	var ctx = context.Background()
	kubeCache := kube.NewKubeCache(ctx)
	settingsClientFactory := &factory.KubeResourceClientFactory{
		Crd:         gloov1.SettingsCrd,
		Cfg:         cfg,
		SharedCache: kubeCache,
	}

	settingsClient, err := gloov1.NewSettingsClient(ctx, settingsClientFactory)
	ExpectWithOffset(1, err).NotTo(HaveOccurred())

	settings, err := settingsClient.Read(testHelper.InstallNamespace, "default", clients.ReadOpts{})
	ExpectWithOffset(1, err).NotTo(HaveOccurred())

	ExpectWithOffset(1, settings.Gateway).NotTo(BeNil())
	ExpectWithOffset(1, settings.Gateway.Validation).NotTo(BeNil())
	settings.Gateway.Validation.AlwaysAccept = &wrappers.BoolValue{Value: false}

	_, err = settingsClient.Write(settings, clients.WriteOpts{OverwriteExisting: true})
	ExpectWithOffset(1, err).NotTo(HaveOccurred())
}

func GetEnterpriseTestHelper(ctx context.Context, namespace string) (*helper.SoloTestHelper, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	if useVersion := GetTestReleasedVersion(ctx, "solo-projects"); useVersion != "" {
		return helper.NewSoloTestHelper(func(defaults helper.TestConfig) helper.TestConfig {
			defaults.RootDir = filepath.Join(cwd, "../../..")
			defaults.HelmChartName = "gloo-ee"
			defaults.ReleasedVersion = useVersion
			defaults.LicenseKey = kubetestutils.LicenseKey()
			defaults.InstallNamespace = namespace
			return defaults
		})
	} else {
		return helper.NewSoloTestHelper(func(defaults helper.TestConfig) helper.TestConfig {
			defaults.RootDir = filepath.Join(cwd, "../../..")
			defaults.HelmChartName = "gloo-ee"
			defaults.LicenseKey = kubetestutils.LicenseKey()
			defaults.InstallNamespace = namespace
			return defaults
		})
	}
}
func PrintGlooDebugLogs() {
	logs, _ := ioutil.ReadFile(cliutil.GetLogsPath())
	fmt.Println("*** Gloo debug logs ***")
	fmt.Println(string(logs))
	fmt.Println("*** End Gloo debug logs ***")
}

func RunAndCleanCommand(name string, arg ...string) []byte {
	cmd := exec.Command(name, arg...)
	b, err := cmd.Output()
	if err != nil {
		if v, ok := err.(*exec.ExitError); ok {
			fmt.Println("ExitError: ", string(v.Stderr))
		}
	}
	ExpectWithOffset(1, err).To(BeNil())
	cmd.Process.Kill()
	cmd.Process.Release()
	return b
}

func CheckGlooHealthy(testHelper *helper.SoloTestHelper) {
	GlooctlCheckEventuallyHealthy(2, testHelper, "180s")
}

func InstallGloo(testHelper *helper.SoloTestHelper, fromRelease string, strictValidation bool, helmOverrideFilePath string) {
	fmt.Printf("\n=============== Installing Gloo : %s ===============\n", fromRelease)
	// construct helm args
	var args = []string{"install", testHelper.HelmChartName}

	RunAndCleanCommand("helm", "repo", "add", testHelper.HelmChartName, GlooeRepoName,
		"--force-update")
	args = append(args, testHelper.HelmChartName+"/gloo-ee",
		"--version", fromRelease)

	args = append(args, "-n", testHelper.InstallNamespace,
		"--create-namespace",
		"--set-string", "license_key="+testHelper.LicenseKey,
		"--values", helmOverrideFilePath)

	fmt.Printf("running helm with args: %v\n", args)

	RunAndCleanCommand("helm", args...)

	if err := testHelper.Deploy(5 * time.Minute); err != nil {
		Expect(err).NotTo(HaveOccurred())
	}

	// Check that everything is OK
	CheckGlooHealthy(testHelper)
}

func InstallGlooWithArgs(testHelper *helper.SoloTestHelper, fromRelease string, additionalArgs []string, helmOverrideFilePath string) {
	fmt.Printf("\n=============== Installing Gloo : %s ===============\n", fromRelease)
	// construct helm args
	var args = []string{"install", testHelper.HelmChartName}

	RunAndCleanCommand("helm", "repo", "add", testHelper.HelmChartName, GlooeRepoName,
		"--force-update")
	args = append(args, testHelper.HelmChartName+"/gloo-ee",
		"--version", fromRelease)

	args = append(args, "-n", testHelper.InstallNamespace,
		"--create-namespace",
		"--set-string", "license_key="+testHelper.LicenseKey,
		"--values", helmOverrideFilePath)

	args = append(args, additionalArgs...)

	fmt.Printf("running helm with args: %v\n", args)

	RunAndCleanCommand("helm", args...)

	if err := testHelper.Deploy(5 * time.Minute); err != nil {
		Expect(err).NotTo(HaveOccurred())
	}

	// Check that everything is OK
	CheckGlooHealthy(testHelper)
}

func UninstallGloo(testHelper *helper.SoloTestHelper, ctx context.Context, cancel context.CancelFunc) {
	Expect(testHelper).ToNot(BeNil())
	err := testHelper.UninstallGloo()
	Expect(err).NotTo(HaveOccurred())
	_, err = MustKubeClient().CoreV1().Namespaces().Get(ctx, testHelper.InstallNamespace, metav1.GetOptions{})
	Expect(apierrors.IsNotFound(err)).To(BeTrue())
	cancel()
}
