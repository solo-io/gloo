package regressions

import (
	"context"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/golang/protobuf/ptypes/wrappers"
	. "github.com/onsi/gomega"
	"github.com/solo-io/gloo/pkg/cliutil"
	v1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/core/matchers"
	"github.com/solo-io/go-utils/contextutils"
	"github.com/solo-io/k8s-utils/kubeutils"
	"github.com/solo-io/k8s-utils/testutils/helper"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/factory"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/kube"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"go.uber.org/zap"
	"k8s.io/client-go/kubernetes"
)

const (
	TestMatcherPrefix = "/test"
)

func WriteVirtualService(
	ctx context.Context,
	testHelper *helper.SoloTestHelper,
	vsClient v1.VirtualServiceClient,
	virtualHostOptions *gloov1.VirtualHostOptions,
	routeOptions *gloov1.RouteOptions,
	sslConfig *gloov1.SslConfig,
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
	sslConfig *gloov1.SslConfig,
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

func MustKubeClient() kubernetes.Interface {
	restConfig, err := kubeutils.GetConfig("", "")
	ExpectWithOffset(1, err).NotTo(HaveOccurred())
	kubeClient, err := kubernetes.NewForConfig(restConfig)
	ExpectWithOffset(1, err).NotTo(HaveOccurred())
	return kubeClient
}

func PrintGlooDebugLogs() {
	logs, _ := ioutil.ReadFile(cliutil.GetLogsPath())
	fmt.Println("*** Gloo debug logs ***")
	fmt.Println(string(logs))
	fmt.Println("*** End Gloo debug logs ***")
}
