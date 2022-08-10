package kube2e

import (
	"context"
	"fmt"
	"io/ioutil"
	"runtime"
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

func GetHttpEchoImage() string {
	httpEchoImageRegistry := "hashicorp/http-echo"
	if runtime.GOARCH == "arm64" {
		httpEchoImageRegistry = "gcr.io/solo-test-236622/http-echo"
	}
	return httpEchoImageRegistry
}

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

// https://github.com/solo-io/gloo/issues/4043#issuecomment-772706604
// We should move tests away from using the testrunner, and instead depend on EphemeralContainers.
// The default response changed in later kube versions, which caused this value to change.
// Ideally the test utilities used by Gloo are maintained in the Gloo repo, so I opted to move
// this constant here.
// This response is given by the testrunner when the SimpleServer is started
func GetSimpleTestRunnerHttpResponse() string {
	if runtime.GOARCH == "arm64" {
		return SimpleTestRunnerHttpResponseArm
	} else {
		return SimpleTestRunnerHttpResponse
	}
}

const SimpleTestRunnerHttpResponse = `<!DOCTYPE html PUBLIC "-//W3C//DTD HTML 3.2 Final//EN"><html>
<title>Directory listing for /</title>
<body>
<h2>Directory listing for /</h2>
<hr>
<ul>
<li><a href="bin/">bin/</a>
<li><a href="boot/">boot/</a>
<li><a href="dev/">dev/</a>
<li><a href="etc/">etc/</a>
<li><a href="home/">home/</a>
<li><a href="lib/">lib/</a>
<li><a href="lib64/">lib64/</a>
<li><a href="media/">media/</a>
<li><a href="mnt/">mnt/</a>
<li><a href="opt/">opt/</a>
<li><a href="proc/">proc/</a>
<li><a href="product_name">product_name</a>
<li><a href="product_uuid">product_uuid</a>
<li><a href="root/">root/</a>
<li><a href="root.crt">root.crt</a>
<li><a href="run/">run/</a>
<li><a href="sbin/">sbin/</a>
<li><a href="srv/">srv/</a>
<li><a href="sys/">sys/</a>
<li><a href="tmp/">tmp/</a>
<li><a href="usr/">usr/</a>
<li><a href="var/">var/</a>
</ul>
<hr>
</body>
</html>`

// I think this is for any local request
const SimpleTestRunnerHttpResponseArm = `<!DOCTYPE html PUBLIC "-//W3C//DTD HTML 3.2 Final//EN"><html>
<title>Directory listing for /</title>
<body>
<h2>Directory listing for /</h2>
<hr>
<ul>
<li><a href="bin/">bin/</a>
<li><a href="boot/">boot/</a>
<li><a href="dev/">dev/</a>
<li><a href="etc/">etc/</a>
<li><a href="home/">home/</a>
<li><a href="lib/">lib/</a>
<li><a href="lib64/">lib64/</a>
<li><a href="media/">media/</a>
<li><a href="mnt/">mnt/</a>
<li><a href="opt/">opt/</a>
<li><a href="proc/">proc/</a>
<li><a href="product_uuid">product_uuid</a>
<li><a href="root/">root/</a>
<li><a href="root.crt">root.crt</a>
<li><a href="run/">run/</a>
<li><a href="sbin/">sbin/</a>
<li><a href="srv/">srv/</a>
<li><a href="sys/">sys/</a>
<li><a href="tmp/">tmp/</a>
<li><a href="usr/">usr/</a>
<li><a href="var/">var/</a>
</ul>
<hr>
</body>
</html>`
