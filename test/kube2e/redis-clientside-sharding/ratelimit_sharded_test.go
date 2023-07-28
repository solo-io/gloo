package clientside_sharding_test

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/solo-io/gloo/test/helpers"

	"github.com/solo-io/solo-projects/test/kube2e"
	"github.com/solo-io/solo-projects/test/services"

	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/core/matchers"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/ratelimit"
	rlv1alpha1 "github.com/solo-io/solo-apis/pkg/api/ratelimit.solo.io/v1alpha1"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	v1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"

	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"

	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
)

var _ = Describe("RateLimit tests", func() {

	var (
		testContext  *kube2e.TestContext
		origSettings *gloov1.Settings // used to capture & restore initial Settings so each test can modify them
	)

	const (
		response200 = "HTTP/1.1 200 OK"
		response429 = "HTTP/1.1 429 Too Many Requests"
	)

	BeforeEach(func() {
		testContext = testContextFactory.NewTestContext()
		testContext.BeforeEach()
		var err error

		origSettings, err = testContext.ResourceClientSet().SettingsClient().Read(testContext.InstallNamespace(), "default", clients.ReadOpts{Ctx: testContext.Ctx()})
		Expect(err).NotTo(HaveOccurred(), "Should be able to read initial settings")

		// Enable RateLimit with DenyOnFail
		currentSettings, err := testContext.ResourceClientSet().SettingsClient().Read(testContext.InstallNamespace(), "default", clients.ReadOpts{Ctx: testContext.Ctx()})
		Expect(err).NotTo(HaveOccurred(), "Should be able to read current settings")
		currentSettings.RatelimitServer.RatelimitServerRef = &core.ResourceRef{
			Name:      "rate-limit",
			Namespace: testContext.InstallNamespace(),
		}
		currentSettings.RatelimitServer.DenyOnFail = true
		_, err = testContext.ResourceClientSet().SettingsClient().Write(currentSettings, clients.WriteOpts{Ctx: testContext.Ctx(), OverwriteExisting: true})
		Expect(err).ToNot(HaveOccurred())

		// bounce envoy, get a clean state (draining listener can break this test). see https://github.com/solo-io/solo-projects/issues/2921 for more.
		out, err := services.KubectlOut(strings.Split("rollout restart -n "+testContext.InstallNamespace()+" deploy/gateway-proxy", " ")...)
		fmt.Println(out)
		Expect(err).ToNot(HaveOccurred())
		out, err = services.KubectlOut(strings.Split("rollout status -n "+testContext.InstallNamespace()+" deploy/gateway-proxy", " ")...)
		fmt.Println(out)
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		// deferring this so that the testContext's `ctx` doesn't get cancelled before the following uses of it
		defer testContext.AfterEach()

		currentSettings, err := testContext.ResourceClientSet().SettingsClient().Read(testContext.InstallNamespace(), "default", clients.ReadOpts{Ctx: testContext.Ctx()})
		Expect(err).NotTo(HaveOccurred(), "Should be able to read current settings")

		if origSettings.Metadata.ResourceVersion != currentSettings.Metadata.ResourceVersion {
			origSettings.Metadata.ResourceVersion = currentSettings.Metadata.ResourceVersion // so we can overwrite settings
			origSettings.RatelimitServer.DenyOnFail = true
			_, err = testContext.ResourceClientSet().SettingsClient().Write(origSettings, clients.WriteOpts{Ctx: testContext.Ctx(), OverwriteExisting: true})
			Expect(err).ToNot(HaveOccurred())
		}
	})

	JustBeforeEach(func() {
		testContext.JustBeforeEach()
	})

	JustAfterEach(func() {
		testContext.JustAfterEach()
	})

	requestOnPath := func(path string) (string, error) {
		curlOpts := testContext.DefaultCurlOptsBuilder().
			WithPath(path).WithHost("rate-limit.example.com").
			WithConnectionTimeout(10).WithVerbose(true).Build()
		res, err := testContext.TestHelper().Curl(curlOpts)
		if err != nil {
			return "", err
		}
		return res, nil
	}

	Context("rate limit at specific rate, even with scaled redis", func() {

		It("can rate limit to upstream", func() {
			testContext.PatchDefaultVirtualService(func(service *v1.VirtualService) *v1.VirtualService {
				return helpers.BuilderFromVirtualService(service).
					WithDomain("rate-limit.example.com").
					WithRoute("health-check", generateHealthCheckRoute(testContext.InstallNamespace())).
					// patch default route to be rate limited to 10 requests per day
					WithRouteOptions(kube2e.DefaultRouteName, &gloov1.RouteOptions{
						PrefixRewrite: &wrappers.StringValue{
							Value: "/",
						},
						RatelimitBasic: &ratelimit.IngressRateLimit{
							AnonymousLimits: &rlv1alpha1.RateLimit{
								Unit:            rlv1alpha1.RateLimit_DAY,
								RequestsPerUnit: 10,
							},
						},
					}).
					Build()
			})
			testContext.EventuallyProxyAccepted()

			// If we're about to roll over to the next day and
			// risk resetting the daily rate-limit mid-test,
			// just wait a minute before starting.
			if time.Now().Format("1504") == "2359" {
				fmt.Fprintln(GinkgoWriter, "the current time is 11.59pm, there is a risk of the daily rate limit resetting mid-test, so we will wait 1 minute before starting")
				time.Sleep(1 * time.Minute)
			}

			// Wait for the service to be happily responding
			Eventually(func() (string, error) {
				return requestOnPath("/HealthCheck")
			}, "120s", "1s").Should(ContainSubstring(response200))

			totalRequestsSent := 0
			for i := 0; i < 10; i++ {
				// First 10 requests should work (all should go to the same Redis)
				res, err := requestOnPath(kube2e.TestMatcherPrefix)
				totalRequestsSent++
				Expect(err).NotTo(HaveOccurred())
				Expect(res).To(ContainSubstring(response200), "request #"+strconv.Itoa(totalRequestsSent)+" should not be rate limited")
			}

			// After we've hit our limit, we should consistently be rate limited:
			Consistently(func() error {
				res, err := requestOnPath(kube2e.TestMatcherPrefix)
				totalRequestsSent++
				if err != nil {
					return err
				}
				if !strings.Contains(res, response429) {
					return fmt.Errorf("request #"+strconv.Itoa(totalRequestsSent)+" was not rate limited, expected %v to be found within %v", response429, res)
				}
				return nil
			}, "5s", ".5s").Should(BeNil())

		})
	})
})

func generateHealthCheckRoute(namespace string) *v1.Route {
	return &v1.Route{
		Action: &v1.Route_RouteAction{
			RouteAction: &gloov1.RouteAction{
				Destination: &gloov1.RouteAction_Single{
					Single: &gloov1.Destination{
						DestinationType: &gloov1.Destination_Upstream{
							Upstream: &core.ResourceRef{
								Name:      "gloo-system-testrunner-1234",
								Namespace: namespace,
							},
						},
					},
				},
			},
		},
		Options: &gloov1.RouteOptions{
			PrefixRewrite: &wrappers.StringValue{
				Value: "/",
			},
		},
		Matchers: []*matchers.Matcher{
			{
				PathSpecifier: &matchers.Matcher_Exact{
					Exact: "/HealthCheck",
				},
			},
		},
	}
}
