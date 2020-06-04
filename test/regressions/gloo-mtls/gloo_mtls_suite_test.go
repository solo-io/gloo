package gloo_mtls_test

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/avast/retry-go"

	v1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/core/matchers"
	"github.com/solo-io/go-utils/contextutils"
	"go.uber.org/zap"

	"github.com/gogo/protobuf/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	errors "github.com/rotisserie/eris"
	"github.com/solo-io/gloo/pkg/cliutil"
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/cmd/check"
	"github.com/solo-io/gloo/projects/gloo/cli/pkg/cmd/options"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/test/helpers"
	"github.com/solo-io/go-utils/kubeutils"
	"github.com/solo-io/go-utils/testutils"
	"github.com/solo-io/go-utils/testutils/clusterlock"
	"github.com/solo-io/go-utils/testutils/exec"
	"github.com/solo-io/go-utils/testutils/helper"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/factory"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/kube"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	skhelpers "github.com/solo-io/solo-kit/test/helpers"
	"k8s.io/client-go/kubernetes"
)

// This file is largely copied from test/regressions/gateway/gateway_suite_test.go (May 2020)

func TestGateway(t *testing.T) {
	if testutils.AreTestsDisabled() {
		return
	}
	skhelpers.RegisterCommonFailHandlers()
	skhelpers.SetupLog()
	_ = os.Remove(cliutil.GetLogsPath())
	skhelpers.RegisterPreFailHandler(printGlooDebugLogs)
	RunSpecs(t, "Gateway Suite")
}

func printGlooDebugLogs() {
	logs, _ := ioutil.ReadFile(cliutil.GetLogsPath())
	fmt.Println("*** Gloo debug logs ***")
	fmt.Println(string(logs))
	fmt.Println("*** End Gloo debug logs ***")
}

const (
	testMatcherPrefix = "/test"
)

var (
	testHelper *helper.SoloTestHelper
	locker     *clusterlock.TestClusterLocker
)

var _ = BeforeSuite(func() {
	cwd, err := os.Getwd()
	Expect(err).NotTo(HaveOccurred())

	testHelper, err = helper.NewSoloTestHelper(func(defaults helper.TestConfig) helper.TestConfig {
		defaults.RootDir = filepath.Join(cwd, "../../..")
		defaults.HelmChartName = "gloo-ee"
		defaults.LicenseKey = "eyJleHAiOjM4Nzk1MTY3ODYsImlhdCI6MTU1NDk0MDM0OCwiayI6IkJ3ZXZQQSJ9.tbJ9I9AUltZ-iMmHBertugI2YIg1Z8Q0v6anRjc66Jo"
		defaults.InstallNamespace = "gloo-system"
		return defaults
	})
	Expect(err).NotTo(HaveOccurred())

	skhelpers.RegisterPreFailHandler(helpers.KubeDumpOnFail(GinkgoWriter, testHelper.InstallNamespace))

	locker, err = clusterlock.NewKubeClusterLocker(mustKubeClient(), clusterlock.Options{})
	Expect(err).NotTo(HaveOccurred())
	Expect(locker.AcquireLock(retry.Attempts(40))).NotTo(HaveOccurred())

	// Install Gloo
	values, cleanup := getHelmOverrides()
	defer cleanup()

	err = testHelper.InstallGloo(helper.GATEWAY, 5*time.Minute, helper.ExtraArgs("--values", values, "-v"))
	Expect(err).NotTo(HaveOccurred())
	Eventually(func() error {
		opts := &options.Options{
			Top: options.Top{
				Ctx: context.Background(),
			},
			Metadata: core.Metadata{
				Namespace: testHelper.InstallNamespace,
			},
		}
		ok, err := check.CheckResources(opts)
		if err != nil {
			return errors.Wrapf(err, "unable to run glooctl check")
		}
		if ok {
			return nil
		}
		return errors.New("glooctl check detected a problem with the installation")
	}, 2*time.Minute, "5s").Should(BeNil())

	// Print out the versions of CLI and server components
	glooctlVersionCommand := []string{
		filepath.Join(testHelper.BuildAssetDir, testHelper.GlooctlExecName),
		"version", "-n", testHelper.InstallNamespace}
	output, err := exec.RunCommandOutput(testHelper.RootDir, true, glooctlVersionCommand...)
	Expect(err).NotTo(HaveOccurred())
	fmt.Println(output)

	// TODO(marco): explicitly enable strict validation, this can be removed once we enable validation by default
	// See https://github.com/solo-io/gloo/issues/1374
	enableStrictValidation()

})

var _ = AfterSuite(func() {
	defer locker.ReleaseLock()

	err := testHelper.UninstallGloo()
	Expect(err).NotTo(HaveOccurred())

	// glooctl should delete the namespace. we do it again just in case it failed
	// ignore errors
	_ = testutils.Kubectl("delete", "namespace", testHelper.InstallNamespace)

	EventuallyWithOffset(1, func() error {
		return testutils.Kubectl("get", "namespace", testHelper.InstallNamespace)
	}, "60s", "1s").Should(HaveOccurred())
})

func mustKubeClient() kubernetes.Interface {
	restConfig, err := kubeutils.GetConfig("", "")
	ExpectWithOffset(1, err).NotTo(HaveOccurred())
	kubeClient, err := kubernetes.NewForConfig(restConfig)
	ExpectWithOffset(1, err).NotTo(HaveOccurred())
	return kubeClient
}

func getHelmOverrides() (filename string, cleanup func()) {
	values, err := ioutil.TempFile("", "*.yaml")
	Expect(err).NotTo(HaveOccurred())
	// Set global.glooMtls.enabled = true, and make sure to pull the quay.io/solo-io
	_, err = values.Write([]byte(`
gloo:
  rbac:    
    namespaced: true
    nameSuffix: e2e-test-rbac-suffix
settings:
  singleNamespace: true
  create: true
prometheus:
  podSecurityPolicy:
    enabled: true
grafana:
  testFramework:
    enabled: false
global:
  glooMtls:
    enabled: true
    sds:
      image:
        registry: quay.io/solo-io
  extensions:
    extAuth:
      # we want to deploy extauth as both a standalone deployment (the default) and as a sidecar in the envoy pod, so we can test both
      envoySidecar: true
`)) // need to override registry because we use gcr and quay confusingly https://github.com/solo-io/solo-projects/issues/1733
	Expect(err).NotTo(HaveOccurred())
	err = values.Close()
	Expect(err).NotTo(HaveOccurred())

	return values.Name(), func() {
		_ = os.Remove(values.Name())
	}
}

func enableStrictValidation() {
	// enable strict validation
	// this can be removed once we enable validation by default
	// set projects/gateway/pkg/syncer.AcceptAllResourcesByDefault is set to false
	cfg, err := kubeutils.GetConfig("", "")
	Expect(err).NotTo(HaveOccurred())

	kubeCache := kube.NewKubeCache(context.Background())
	settingsClientFactory := &factory.KubeResourceClientFactory{
		Crd:         gloov1.SettingsCrd,
		Cfg:         cfg,
		SharedCache: kubeCache,
	}

	settingsClient, err := gloov1.NewSettingsClient(settingsClientFactory)
	Expect(err).NotTo(HaveOccurred())

	settings, err := settingsClient.Read(testHelper.InstallNamespace, "default", clients.ReadOpts{})
	Expect(err).NotTo(HaveOccurred())

	Expect(settings.Gateway).NotTo(BeNil())
	Expect(settings.Gateway.Validation).NotTo(BeNil())
	settings.Gateway.Validation.AlwaysAccept = &types.BoolValue{Value: false}

	_, err = settingsClient.Write(settings, clients.WriteOpts{OverwriteExisting: true})
	Expect(err).NotTo(HaveOccurred())
}

func writeVirtualService(ctx context.Context, vsClient v1.VirtualServiceClient,
	virtualHostOptions *gloov1.VirtualHostOptions, routeOptions *gloov1.RouteOptions,
	sslConfig *gloov1.SslConfig) {

	upstreamRef := &core.ResourceRef{
		Namespace: testHelper.InstallNamespace,
		Name:      fmt.Sprintf("%s-%s-%v", testHelper.InstallNamespace, "testrunner", helper.TestRunnerPort),
	}
	writeCustomVirtualService(ctx, vsClient, virtualHostOptions, routeOptions, sslConfig, upstreamRef)
}

func writeCustomVirtualService(ctx context.Context, vsClient v1.VirtualServiceClient,
	virtualHostOptions *gloov1.VirtualHostOptions, routeOptions *gloov1.RouteOptions,
	sslConfig *gloov1.SslConfig, upstreamRef *core.ResourceRef) {

	if routeOptions.GetPrefixRewrite() == nil {
		if routeOptions == nil {
			routeOptions = &gloov1.RouteOptions{}
		}
		routeOptions.PrefixRewrite = &types.StringValue{
			Value: "/",
		}
	}

	// We wrap this in a eventually because the validating webhook may reject the virtual service if one of the
	// resources the VS depends on is not yet available.
	EventuallyWithOffset(1, func() error {
		_, err := vsClient.Write(&v1.VirtualService{

			Metadata: core.Metadata{
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
							Prefix: testMatcherPrefix,
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

func deleteVirtualService(vsClient v1.VirtualServiceClient, ns, name string, opts clients.DeleteOpts) {
	// We wrap this in a eventually because the validating webhook may reject the virtual service if one of the
	// resources the VS depends on is not yet available.
	EventuallyWithOffset(1, func() error {
		return vsClient.Delete(ns, name, opts)
	}, time.Minute, "5s").Should(BeNil())
}
