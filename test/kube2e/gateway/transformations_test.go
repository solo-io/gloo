package gateway_test

import (
	"fmt"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/golang/protobuf/ptypes/wrappers"
	glooKube2e "github.com/solo-io/gloo/test/kube2e"

	"github.com/solo-io/gloo/projects/gloo/pkg/plugins/kubernetes"

	"github.com/solo-io/gloo/test/helpers"

	"github.com/solo-io/solo-projects/test/kube2e"
	"github.com/solo-io/solo-projects/test/services"

	"github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/extensions/transformers/xslt"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/core/matchers"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/dlp"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/transformation"
	"github.com/solo-io/go-utils/testutils"
	envoy_type "github.com/solo-io/solo-kit/pkg/api/external/envoy/type"
	"github.com/solo-io/solo-projects/test/e2e/transformation_helpers"

	"github.com/solo-io/k8s-utils/testutils/helper"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	v1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gateway/pkg/defaults"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"

	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
)

var _ = Describe("dlp tests", func() {

	var (
		testContext      *kube2e.TestContext
		httpEcho         helper.TestRunner
		preservedGateway *v1.Gateway
	)

	BeforeEach(func() {
		testContext = testContextFactory.NewTestContext()
		testContext.BeforeEach()

		var err error
		httpEcho, err = helper.NewEchoHttp(testContext.InstallNamespace())
		Expect(err).NotTo(HaveOccurred())

		err = httpEcho.Deploy(2 * time.Minute)
		Expect(err).NotTo(HaveOccurred())

		// the gateway may be modified during the test, so we retrieve it beforehand and return it to the previous state after
		preservedGateway, err = testContext.ResourceClientSet().GatewayClient().Read(testContext.InstallNamespace(), defaults.DefaultGateway(testContext.InstallNamespace()).GetMetadata().GetName(), clients.ReadOpts{})
		Expect(err).NotTo(HaveOccurred())
		preservedGateway.Metadata.ResourceVersion = ""

		// bounce envoy, get a clean state (draining listener can break this test). see https://github.com/solo-io/solo-projects/issues/2921 for more.
		out, err := services.KubectlOut(strings.Split("rollout restart -n "+testContext.InstallNamespace()+" deploy/gateway-proxy", " ")...)
		fmt.Println(out)
		Expect(err).ToNot(HaveOccurred())
		out, err = services.KubectlOut(strings.Split("rollout status -n "+testContext.InstallNamespace()+" deploy/gateway-proxy", " ")...)
		fmt.Println(out)
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		testContext.AfterEach()
		err := httpEcho.Terminate()
		Expect(err).NotTo(HaveOccurred())

		// return the gateway to previous state, deleting previous gateway to avoid resource version conflict
		err = testContext.ResourceClientSet().GatewayClient().Delete(testContext.InstallNamespace(), preservedGateway.Metadata.GetName(), clients.DeleteOpts{IgnoreNotExist: true})
		Expect(err).NotTo(HaveOccurred())
		_, err = testContext.ResourceClientSet().GatewayClient().Write(preservedGateway, clients.WriteOpts{OverwriteExisting: true})
		Expect(err).NotTo(HaveOccurred())

		// Delete http echo service
		err = testutils.Kubectl("delete", "service", "-n", testContext.InstallNamespace(), helper.HttpEchoName, "--grace-period=0")
		Expect(err).NotTo(HaveOccurred())
	})

	JustBeforeEach(func() {
		testContext.JustBeforeEach()
	})

	JustAfterEach(func() {
		testContext.JustAfterEach()
	})

	checkConnection := func(path, body string) {
		testContext.EventuallyProxyAccepted()

		curlOpts := testContext.DefaultCurlOptsBuilder().
			WithPath(path).WithHeaders(map[string]string{"hello": "world"}).
			WithConnectionTimeout(10).WithVerbose(true).Build()
		testContext.TestHelper().CurlEventuallyShouldRespond(curlOpts, body, 1, time.Minute*5)
	}

	Context("data loss prevention", func() {

		It("will mask regex matches", func() {

			dlpVhost := &dlp.Config{
				Actions: []*dlp.Action{
					{
						ActionType: dlp.Action_CUSTOM,
						CustomAction: &dlp.CustomAction{
							Name:     "test",
							Regex:    []string{"hello", "world"},
							MaskChar: "Y",
							Percent: &envoy_type.Percent{
								Value: 60,
							},
						},
					},
				},
			}

			httpEchoRef := &core.ResourceRef{
				Namespace: testContext.InstallNamespace(),
				Name:      kubernetes.UpstreamName(testContext.InstallNamespace(), helper.HttpEchoName, helper.HttpEchoPort),
			}
			testContext.PatchDefaultVirtualService(func(service *v1.VirtualService) *v1.VirtualService {
				virtualHostPlugins := &gloov1.VirtualHostOptions{
					Dlp: dlpVhost,
				}

				return helpers.BuilderFromVirtualService(service).WithVirtualHostOptions(virtualHostPlugins).WithRouteActionToUpstreamRef(kube2e.DefaultRouteName, httpEchoRef).Build()
			})
			checkConnection(kube2e.TestMatcherPrefix, `"YYYlo":"YYYld"`)
		})

		It("will only apply to matched routes on gateway-level dlp option", func() {
			unmatchedPath := "/unmatched"
			matcher := &matchers.Matcher{
				PathSpecifier: &matchers.Matcher_Regex{
					Regex: "\\/t.*t",
				},
			}
			action := &dlp.Action{
				ActionType: dlp.Action_CUSTOM,
				CustomAction: &dlp.CustomAction{
					Name:     "test",
					Regex:    []string{"hello", "world"},
					MaskChar: "Y",
					Percent: &envoy_type.Percent{
						Value: 60,
					},
				},
			}
			testContext.PatchDefaultGateway(func(gateway *v1.Gateway) *v1.Gateway {
				gateway.GatewayType = &v1.Gateway_HttpGateway{
					HttpGateway: &v1.HttpGateway{
						Options: &gloov1.HttpListenerOptions{
							Dlp: &dlp.FilterConfig{
								DlpRules: []*dlp.DlpRule{
									{
										Matcher: matcher,
										Actions: []*dlp.Action{
											action,
										},
									},
								},
							},
						},
					},
				}
				return gateway
			})
			testContext.EventuallyProxyAccepted()

			httpEchoRef := &core.ResourceRef{
				Namespace: testContext.InstallNamespace(),
				Name:      kubernetes.UpstreamName(testContext.InstallNamespace(), helper.HttpEchoName, helper.HttpEchoPort),
			}

			// we expect DLP to apply when requesting a path, /test, that matches the matcher in the DlpRule
			testContext.PatchDefaultVirtualService(func(service *v1.VirtualService) *v1.VirtualService {
				return helpers.BuilderFromVirtualService(service).WithRouteActionToUpstreamRef(kube2e.DefaultRouteName, httpEchoRef).Build()
			})
			checkConnection(kube2e.TestMatcherPrefix, `"YYYlo":"YYYld"`)

			// we do not expect DLP to apply when requesting a path, /unmatched, that does not match the matcher in the DlpRule
			testContext.PatchDefaultVirtualService(func(service *v1.VirtualService) *v1.VirtualService {
				return helpers.BuilderFromVirtualService(service).WithRouteActionToUpstreamRef(kube2e.DefaultRouteName, httpEchoRef).WithRoutePrefixMatcher(kube2e.DefaultRouteName, unmatchedPath).Build()
			})
			checkConnection(unmatchedPath, `"hello":"world"`)
		})
	})

	// These tests are `Ordered` so the settings are only updated once, and not per-test.
	Context("xslt transformer", Ordered, func() {
		BeforeAll(func() {
			// Strict validation is enabled by default, but we need to disable it for these tests.
			glooKube2e.UpdateSettings(testContext.Ctx(), func(settings *gloov1.Settings) {
				settings.Gateway.Validation.AlwaysAccept = &wrappers.BoolValue{Value: true}
				settings.Gateway.Validation.AllowWarnings = &wrappers.BoolValue{Value: true}
			}, testContext.InstallNamespace())
		})

		AfterAll(func() {
			// Reset settings to default values
			glooKube2e.UpdateSettings(testContext.Ctx(), func(settings *gloov1.Settings) {
				settings.Gateway.Validation.AlwaysAccept = &wrappers.BoolValue{Value: false}
				settings.Gateway.Validation.AllowWarnings = &wrappers.BoolValue{Value: true}
			}, testContext.InstallNamespace())
		})

		expectBody := func(body, expectedBody string) {
			testContext.EventuallyProxyAccepted()
			expectedString := regexp.MustCompile("[\\r\\n\\s]+").ReplaceAllString(expectedBody, "")
			curlOpts := testContext.DefaultCurlOptsBuilder().
				WithHeaders(map[string]string{"hello": "world"}).WithConnectionTimeout(10).
				WithVerbose(true).WithBody(body).Build()
			testContext.TestHelper().CurlEventuallyShouldRespond(curlOpts, expectedString, 1, time.Second*20)
		}

		It("will transform xml -> json", func() {
			if runtime.GOARCH == "arm64" {
				Skip("Fails on arm64")
			}
			virtualHostPlugins := &gloov1.VirtualHostOptions{
				StagedTransformations: &transformation.TransformationStages{
					Early: &transformation.RequestResponseTransformations{
						RequestTransforms: []*transformation.RequestMatch{
							{
								Matcher: &matchers.Matcher{
									PathSpecifier: &matchers.Matcher_Prefix{Prefix: "/"},
								},
								RequestTransformation: &transformation.Transformation{
									TransformationType: &transformation.Transformation_XsltTransformation{
										XsltTransformation: &xslt.XsltTransformation{
											Xslt: transformation_helpers.XmlToJsonTransform,
										},
									},
								},
							},
						},
					},
				},
			}

			httpEchoRef := &core.ResourceRef{
				Namespace: testContext.InstallNamespace(),
				Name:      kubernetes.UpstreamName(testContext.InstallNamespace(), helper.HttpEchoName, helper.HttpEchoPort),
			}
			testContext.PatchDefaultVirtualService(func(service *v1.VirtualService) *v1.VirtualService {
				return helpers.BuilderFromVirtualService(service).WithVirtualHostOptions(virtualHostPlugins).WithRouteActionToUpstreamRef(kube2e.DefaultRouteName, httpEchoRef).Build()
			})
			expectBody(transformation_helpers.CarsXml, transformation_helpers.CarsJson)
		})
	})
})
