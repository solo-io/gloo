package gateway_test

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"time"

	"github.com/avast/retry-go"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/solo-io/gloo/pkg/cliutil/install"
	v1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/extensions/transformation"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/core/matchers"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/gloosnapshot"
	glooStatic "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/static"
	glootransformation "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/transformation"

	defaults2 "github.com/solo-io/gloo/projects/gloo/pkg/defaults"
	"github.com/solo-io/gloo/projects/gloo/pkg/translator"
	glooKube2e "github.com/solo-io/gloo/test/kube2e"

	"github.com/solo-io/solo-projects/test/kube2e"
	"k8s.io/apimachinery/pkg/util/intstr"

	kubernetes2 "github.com/solo-io/gloo/projects/gloo/pkg/plugins/kubernetes"
	"github.com/solo-io/gloo/test/helpers"
	"github.com/solo-io/go-utils/testutils"
	"github.com/solo-io/k8s-utils/testutils/helper"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/solo-io/gloo/projects/gateway/pkg/defaults"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	gloossl "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/ssl"
	"github.com/solo-io/k8s-utils/kubeutils"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"

	osskube2e "github.com/solo-io/gloo/test/kube2e"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-kit/test/setup"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("Installing gloo in gateway mode", func() {

	var (
		ctx               context.Context
		cancel            context.CancelFunc
		resourceClientset *glooKube2e.KubeResourceClientSet
		snapshotWriter    helpers.SnapshotWriter
		glooResources     *gloosnapshot.ApiSnapshot
	)

	BeforeEach(func() {
		ctx, cancel = context.WithCancel(context.Background())
		cfg, err := kubeutils.GetConfig("", "")
		Expect(err).NotTo(HaveOccurred())

		resourceClientset, err = glooKube2e.NewKubeResourceClientSet(ctx, cfg)
		Expect(err).NotTo(HaveOccurred())

		// Create a VirtualService routing directly to the testrunner kubernetes service
		// A virtual service has to be created to test the gloo validations
		testRunnerDestination := &gloov1.Destination{
			DestinationType: &gloov1.Destination_Kube{
				Kube: &gloov1.KubernetesServiceDestination{
					Ref: &core.ResourceRef{
						Namespace: testHelper.InstallNamespace,
						Name:      helper.TestrunnerName,
					},
					Port: uint32(helper.TestRunnerPort),
				},
			},
		}
		testRunnerVs := helpers.NewVirtualServiceBuilder().
			WithName(helper.TestrunnerName).
			WithNamespace(testHelper.InstallNamespace).
			WithDomain(helper.TestrunnerName).
			WithRoutePrefixMatcher(helper.TestrunnerName, "/").
			WithRouteActionToSingleDestination(helper.TestrunnerName, testRunnerDestination).
			Build()

		// The set of resources that these tests will generate
		glooResources = &gloosnapshot.ApiSnapshot{
			VirtualServices: v1.VirtualServiceList{
				// many tests route to the TestRunner Service so it makes sense to just
				// always create it
				// the other benefit is this ensures that all tests start with a valid Proxy CR
				testRunnerVs,
			},
		}
		snapshotWriter = helpers.NewSnapshotWriter(resourceClientset, []retry.Option{})
	})

	AfterEach(func() {
		kube2e.DeleteVirtualService(resourceClientset.VirtualServiceClient(), testHelper.InstallNamespace, "vs", clients.DeleteOpts{Ctx: ctx, IgnoreNotExist: true})
		cancel()
	})

	JustBeforeEach(func() {
		err := snapshotWriter.WriteSnapshot(glooResources, clients.WriteOpts{
			Ctx:               ctx,
			OverwriteExisting: false,
		})
		Expect(err).NotTo(HaveOccurred())
	})

	JustAfterEach(func() {
		err := snapshotWriter.DeleteSnapshot(glooResources, clients.DeleteOpts{
			Ctx:            ctx,
			IgnoreNotExist: true,
		})
		Expect(err).NotTo(HaveOccurred())
	})

	It("can route request to upstream", func() {

		kube2e.WriteVirtualService(ctx, testHelper, resourceClientset.VirtualServiceClient(), nil, nil, nil)

		defaultGateway := defaults.DefaultGateway(testHelper.InstallNamespace)
		// wait for default gateway to be created
		Eventually(func() (*v1.Gateway, error) {
			return resourceClientset.GatewayClient().Read(testHelper.InstallNamespace, defaultGateway.Metadata.Name, clients.ReadOpts{Ctx: ctx})
		}, "15s", "0.5s").Should(Not(BeNil()))

		gatewayPort := 80
		testHelper.CurlEventuallyShouldRespond(helper.CurlOpts{
			Protocol:          "http",
			Path:              kube2e.TestMatcherPrefix,
			Method:            "GET",
			Host:              defaults.GatewayProxyName,
			Service:           defaults.GatewayProxyName,
			Port:              gatewayPort,
			ConnectionTimeout: 10, // this is important, as the first curl call sometimes hangs indefinitely
		}, osskube2e.GetSimpleTestRunnerHttpResponse(), 1, time.Minute*5)
	})

	Context("virtual service in configured with SSL", func() {

		BeforeEach(func() {
			// get the certificate so it is generated in the background
			go helpers.Certificate()
		})

		AfterEach(func() {
			err := resourceClientset.KubeClients().CoreV1().Secrets(testHelper.InstallNamespace).Delete(ctx, "secret", metav1.DeleteOptions{})
			Expect(err).NotTo(HaveOccurred())
		})

		It("can route https request to upstream", func() {
			sslSecret := helpers.GetKubeSecret("secret", testHelper.InstallNamespace)
			createdSecret, err := resourceClientset.KubeClients().CoreV1().Secrets(testHelper.InstallNamespace).Create(ctx, sslSecret, metav1.CreateOptions{})
			Expect(err).NotTo(HaveOccurred())

			Eventually(func() error {
				_, err := resourceClientset.KubeClients().CoreV1().Secrets(sslSecret.Namespace).Get(ctx, sslSecret.Name, metav1.GetOptions{})
				return err
			}, "10s", "0.5s").Should(BeNil())
			time.Sleep(3 * time.Second) // Wait a few seconds so Gloo can pick up the secret, otherwise the webhook validation might fail

			sslConfig := &gloossl.SslConfig{
				SslSecrets: &gloossl.SslConfig_SecretRef{
					SecretRef: &core.ResourceRef{
						Name:      createdSecret.ObjectMeta.Name,
						Namespace: createdSecret.ObjectMeta.Namespace,
					},
				},
			}

			kube2e.WriteVirtualService(ctx, testHelper, resourceClientset.VirtualServiceClient(), nil, nil, sslConfig)

			defaultGateway := defaults.DefaultGateway(testHelper.InstallNamespace)
			// wait for default gateway to be created
			Eventually(func() (*v1.Gateway, error) {
				return resourceClientset.GatewayClient().Read(testHelper.InstallNamespace, defaultGateway.Metadata.Name, clients.ReadOpts{})
			}, "15s", "0.5s").Should(Not(BeNil()))

			gatewayPort := 443
			caFile := osskube2e.ToFile(helpers.Certificate())
			//goland:noinspection  GoUnhandledErrorResult
			defer os.Remove(caFile)

			err = testutils.Kubectl("cp", caFile, testHelper.InstallNamespace+"/testrunner:/tmp/ca.crt")
			Expect(err).NotTo(HaveOccurred())

			testHelper.CurlEventuallyShouldRespond(helper.CurlOpts{
				Protocol:          "https",
				Path:              kube2e.TestMatcherPrefix,
				Method:            "GET",
				Host:              defaults.GatewayProxyName,
				Service:           defaults.GatewayProxyName,
				Port:              gatewayPort,
				CaFile:            "/tmp/ca.crt",
				ConnectionTimeout: 10, // this is important, as the first curl call sometimes hangs indefinitely
			}, osskube2e.GetSimpleTestRunnerHttpResponse(), 1, time.Minute*2)
		})
	})

	It("rejects invalid inja template in transformation", func() {
		injaTransform := `{% if default(data.error.message, "") != "" %}400{% else %}{{ header(":status") }}{% endif %}`
		t := &glootransformation.Transformations{
			ClearRouteCache: true,
			ResponseTransformation: &glootransformation.Transformation{
				TransformationType: &glootransformation.Transformation_TransformationTemplate{
					TransformationTemplate: &transformation.TransformationTemplate{
						Headers: map[string]*transformation.InjaTemplate{
							":status": {Text: injaTransform},
						},
					},
				},
			},
		}

		dest := &gloov1.Destination{
			DestinationType: &gloov1.Destination_Upstream{
				Upstream: &core.ResourceRef{
					Namespace: testHelper.InstallNamespace,
					Name:      fmt.Sprintf("%s-%s-%v", testHelper.InstallNamespace, helper.TestrunnerName, helper.TestRunnerPort),
				},
			},
		}

		vs := getVirtualService(dest, nil)
		vs.VirtualHost.Options = &gloov1.VirtualHostOptions{Transformations: t}

		_, err := resourceClientset.VirtualServiceClient().Write(vs, clients.WriteOpts{Ctx: ctx})
		Expect(err).ToNot(HaveOccurred())

		err = resourceClientset.VirtualServiceClient().Delete(vs.Metadata.Namespace, vs.Metadata.Name, clients.DeleteOpts{Ctx: ctx})
		Expect(err).ToNot(HaveOccurred())

		// trim trailing "}", which should invalidate our inja template
		t.ResponseTransformation.GetTransformationTemplate().Headers[":status"].Text = strings.TrimSuffix(injaTransform, "}")

		_, err = resourceClientset.VirtualServiceClient().Write(vs, clients.WriteOpts{Ctx: ctx})
		Expect(err).To(MatchError(ContainSubstring("Failed to parse response template: Failed to parse " +
			"header template ':status': [inja.exception.parser_error] expected statement close, got '%'")))
	})

	Context("tests for validation", func() {
		BeforeEach(func() {
			glooKube2e.UpdateAlwaysAcceptSetting(ctx, false, testHelper.InstallNamespace)
		})
		testValidation := func(yaml, expectedErr string) {
			out, err := install.KubectlApplyOut([]byte(yaml))

			testValidationDidError := func() {
				ExpectWithOffset(1, err).To(HaveOccurred())
				ExpectWithOffset(1, string(out)).To(ContainSubstring(expectedErr))
			}

			testValidationDidSucceed := func() {
				ExpectWithOffset(1, err).NotTo(HaveOccurred())
				// To ensure that we do not leave artifacts between tests
				// we cleanup the resource after it is accepted
				err = install.KubectlDelete([]byte(yaml))
				ExpectWithOffset(1, err).NotTo(HaveOccurred())
			}

			if expectedErr == "" {
				testValidationDidSucceed()
			} else {
				testValidationDidError()
			}
		}

		Context("extension resources", func() {

			type testCase struct {
				resourceYaml, expectedErr string
			}

			BeforeEach(func() {
				// Set the validation settings to be as strict as possible so that we can trigger
				// rejections by just producing a warning on the resource
				glooKube2e.UpdateSettings(ctx, func(settings *gloov1.Settings) {
					Expect(settings.GetGateway().GetValidation()).NotTo(BeNil())
					settings.GetGateway().GetValidation().AllowWarnings = &wrappers.BoolValue{Value: false}
				}, testHelper.InstallNamespace)
			})

			AfterEach(func() {
				glooKube2e.UpdateSettings(ctx, func(settings *gloov1.Settings) {
					Expect(settings.GetGateway().GetValidation()).NotTo(BeNil())
					settings.GetGateway().GetValidation().AllowWarnings = &wrappers.BoolValue{Value: true}
				}, testHelper.InstallNamespace)
			})

			JustBeforeEach(func() {
				// Validation of Gloo resources requires that a Proxy resource exist
				// Therefore, before the tests start, we must attempt updates that should be rejected
				// They will only be rejected once a Proxy exists in the ApiSnapshot

				// the action value is not equal to the Descriptor value, so this should always fail
				upstream := &gloov1.Upstream{
					Metadata: &core.Metadata{
						Name:      "",
						Namespace: testHelper.InstallNamespace,
					},
					UpstreamType: &gloov1.Upstream_Static{
						Static: &glooStatic.UpstreamSpec{
							Hosts: []*glooStatic.Host{{
								Addr: "~",
							}},
						},
					},
				}
				attempt := 0
				Eventually(func(g Gomega) bool {
					upstream.Metadata.Name = fmt.Sprintf("invalid-placeholder-upstream-%d", attempt)
					_, err := resourceClientset.UpstreamClient().Write(upstream, clients.WriteOpts{Ctx: ctx})
					if err != nil {
						serr := err.Error()
						g.Expect(serr).Should(ContainSubstring("admission webhook"))
						g.Expect(serr).Should(ContainSubstring("port cannot be empty for host"))
						// We have successfully rejected an invalid upstream
						// This means that the webhook is fully warmed, and contains a Snapshot with a Proxy
						return true
					}

					err = resourceClientset.UpstreamClient().Delete(
						upstream.GetMetadata().GetNamespace(),
						upstream.GetMetadata().GetName(),
						clients.DeleteOpts{Ctx: ctx, IgnoreNotExist: true})
					g.Expect(err).NotTo(HaveOccurred())

					attempt += 1
					return false
				}, time.Second*15, time.Second*1).Should(BeTrue())
			})

			It("rejects bad resources", func() {
				testCases := []testCase{{
					resourceYaml: `
apiVersion: gateway.solo.io/v1
kind: VirtualService
metadata:
  name: missing-rlc-vs
  namespace: gloo-system
spec:
  virtualHost:
    domains:
      - "my-invalid-rate-limit-domain"
    options:
      rateLimitConfigs:
        refs:
          - name: invalid-rlc-name
            namespace: gloo-system
`,
					expectedErr: "could not find RateLimitConfig resource with name",
				},
				}
				for _, tc := range testCases {
					testValidation(tc.resourceYaml, tc.expectedErr)
				}
			})

			checkThatVSRefCantBeDeleted := func(resourceYaml, vsYaml string) string {
				err := install.KubectlApply([]byte(resourceYaml))
				Expect(err).ToNot(HaveOccurred())
				Eventually(func() error {
					// eventually the resource will be applied and we can apply the virtual service
					err = install.KubectlApply([]byte(vsYaml))
					return err
				}, "5s", "1s").Should(BeNil())
				var out []byte
				// we should get an error saying that the admission webhook can not find the resource this is because the VS
				// references the resource, and the allowWarnings property is not set.  The warning for a resource missing should
				// error in the reports
				// adding a sleep here, because it seems that the snapshot take time to pick up the new VS, and RLC
				time.Sleep(5 * time.Second)
				Eventually(func() error {
					out, err = install.KubectlOut(bytes.NewBuffer([]byte(resourceYaml)), []string{"delete", "-f", "-"}...)
					return err
				}, "5s", "1s").Should(Not(BeNil()))

				// delete the VS and the resource that the VS references have to wait for the snapshot to sync in the gateway
				// validator for the resource to be deleted
				Eventually(func(g Gomega) {
					err = install.KubectlDelete([]byte(vsYaml))
					g.Expect(err).ToNot(HaveOccurred())
				}, "5s", "1s")
				Eventually(func(g Gomega) {
					err = install.KubectlDelete([]byte(resourceYaml))
					g.Expect(err).ToNot(HaveOccurred())
				}, "5s", "1s")
				return string(out)
			}

			It("rejects deleting rate limit config referenced on a Virtual Service", func() {
				rateLimitYaml := `
apiVersion: ratelimit.solo.io/v1alpha1
kind: RateLimitConfig
metadata:
  name: rlc
  namespace: gloo-system
spec:
  raw:
    descriptors:
      - key: foo
        value: foo
        rateLimit:
          requestsPerUnit: 1
          unit: MINUTE
    rateLimits:
      - actions:
        - genericKey:
            descriptorValue: bar
`
				vsYaml := `
apiVersion: gateway.solo.io/v1
kind: VirtualService
metadata:
  name: vs
  namespace: gloo-system
spec:
  virtualHost:
    domains:
      - "valid-domain"
    options:
      rateLimitConfigs:
        refs:
          - name: rlc
            namespace: gloo-system
`
				out := checkThatVSRefCantBeDeleted(rateLimitYaml, vsYaml)
				Expect(out).To(ContainSubstring("Error from server"))
				Expect(out).To(ContainSubstring("admission webhook"))
				Expect(out).To(ContainSubstring("could not find RateLimitConfig resource with name"))
			})

		})
	})

	Context("matchable hybrid gateway", func() {

		var (
			hybridProxyServicePort = corev1.ServicePort{
				Name:       "hybrid-proxy",
				Port:       int32(defaults2.HybridPort),
				TargetPort: intstr.FromInt(int(defaults2.HybridPort)),
				Protocol:   "TCP",
			}
			tcpEchoClusterName string
		)

		exposePortOnGwProxyService := func(servicePort corev1.ServicePort) {
			gwSvc, err := resourceClientset.KubeClients().CoreV1().Services(testHelper.InstallNamespace).Get(ctx, defaults.GatewayProxyName, metav1.GetOptions{})
			Expect(err).NotTo(HaveOccurred())

			// Append servicePort if not found already
			found := false
			for _, v := range gwSvc.Spec.Ports {
				if v.Name == hybridProxyServicePort.Name || v.Port == hybridProxyServicePort.Port {
					found = true
					break
				}
			}
			if !found {
				gwSvc.Spec.Ports = append(gwSvc.Spec.Ports, hybridProxyServicePort)
			}

			_, err = resourceClientset.KubeClients().CoreV1().Services(testHelper.InstallNamespace).Update(ctx, gwSvc, metav1.UpdateOptions{})
			Expect(err).NotTo(HaveOccurred())
		}
		var httpEcho helper.TestRunner
		var err error
		BeforeEach(func() {
			Skip("to merge other 1.15 code")

			caFile := glooKube2e.ToFile(helpers.Certificate())
			//goland:noinspection GoUnhandledErrorResult
			defer os.Remove(caFile)
			err = setup.Kubectl("cp", caFile, testHelper.InstallNamespace+"/testrunner:/tmp/ca.crt")
			Expect(err).NotTo(HaveOccurred())
			exposePortOnGwProxyService(hybridProxyServicePort)

			httpEcho, err = helper.NewEchoHttp(testHelper.InstallNamespace)
			Expect(err).NotTo(HaveOccurred())

			err = httpEcho.Deploy(2 * time.Minute) // mimick transformations test
			Expect(err).NotTo(HaveOccurred())

			tcpEchoClusterName = translator.UpstreamToClusterName(&core.ResourceRef{
				Namespace: testHelper.InstallNamespace,
				Name:      kubernetes2.UpstreamName(testHelper.InstallNamespace, helper.HttpEchoName, helper.HttpEchoPort),
			})

			// Create a MatchableHttpGateway
			matchableHttpGateway := &v1.MatchableHttpGateway{
				Metadata: &core.Metadata{
					Name:      "matchable-http-gateway",
					Namespace: testHelper.InstallNamespace,
				},
				HttpGateway: &v1.HttpGateway{
					// match all virtual services
				},
			}

			// Create a MatchableTcpGateway
			matchableTcpGateway := &v1.MatchableTcpGateway{
				Metadata: &core.Metadata{
					Name:      "matchable-tcp-gateway",
					Namespace: testHelper.InstallNamespace,
				},
				Matcher: &v1.MatchableTcpGateway_Matcher{

					PassthroughCipherSuites: []string{"AES128-SHA256"},
				},
				TcpGateway: &v1.TcpGateway{
					TcpHosts: []*gloov1.TcpHost{{
						Name: tcpEchoClusterName,
						Destination: &gloov1.TcpHost_TcpAction{
							Destination: &gloov1.TcpHost_TcpAction_Single{
								Single: &gloov1.Destination{
									DestinationType: &gloov1.Destination_Upstream{
										Upstream: &core.ResourceRef{
											Namespace: testHelper.InstallNamespace,
											Name:      fmt.Sprintf("%s-%s-%v", testHelper.InstallNamespace, helper.HttpEchoName, helper.HttpEchoPort),
										},
									},
								},
							},
						},
					}},
				},
			}

			// Create a HybridGateway that references that MatchableHttpGateway
			hybridGateway := &v1.Gateway{
				Metadata: &core.Metadata{
					Name:      fmt.Sprintf("%s-hybrid", defaults.GatewayProxyName),
					Namespace: testHelper.InstallNamespace,
				},
				GatewayType: &v1.Gateway_HybridGateway{
					HybridGateway: &v1.HybridGateway{
						DelegatedHttpGateways: &v1.DelegatedHttpGateway{
							SelectionType: &v1.DelegatedHttpGateway_Ref{
								Ref: &core.ResourceRef{
									Name:      matchableHttpGateway.GetMetadata().GetName(),
									Namespace: matchableHttpGateway.GetMetadata().GetNamespace(),
								},
							},
						},
						DelegatedTcpGateways: &v1.DelegatedTcpGateway{
							SelectionType: &v1.DelegatedTcpGateway_Ref{
								Ref: &core.ResourceRef{
									Name:      matchableTcpGateway.GetMetadata().GetName(),
									Namespace: matchableTcpGateway.GetMetadata().GetNamespace(),
								},
							},
						},
					},
				},
				ProxyNames:    []string{defaults.GatewayProxyName},
				BindAddress:   defaults.GatewayBindAddress,
				BindPort:      defaults2.HybridPort,
				UseProxyProto: &wrappers.BoolValue{Value: false},
			}

			glooResources.HttpGateways = v1.MatchableHttpGatewayList{matchableHttpGateway}
			glooResources.TcpGateways = v1.MatchableTcpGatewayList{matchableTcpGateway}
			glooResources.Gateways = v1.GatewayList{hybridGateway}
		})
		AfterEach(func() {
			httpEcho.Terminate()

		})
		It("works", func() {

			checkCmd := exec.Command("curl", "--http0.9", "-sv", "--request", "POST", "--ciphers", "AES128-SHA256", tcpEchoClusterName, "-d", "something ")
			err = checkCmd.Run()
			checkCmd.Stdout = GinkgoWriter
			checkCmd.Stderr = GinkgoWriter
			Expect(err).NotTo(HaveOccurred())

		})

	})

})

func getVirtualService(dest *gloov1.Destination, sslConfig *gloossl.SslConfig) *v1.VirtualService {
	return getVirtualServiceWithRoute(getRouteWithDest(dest, "/"), sslConfig)
}

func getVirtualServiceWithRoute(route *v1.Route, sslConfig *gloossl.SslConfig) *v1.VirtualService {
	return &v1.VirtualService{
		Metadata: &core.Metadata{
			Name:      "vs",
			Namespace: testHelper.InstallNamespace,
		},
		SslConfig: sslConfig,
		VirtualHost: &v1.VirtualHost{
			Domains: []string{"*"},

			Routes: []*v1.Route{route},
		},
	}
}

func getRouteWithDest(dest *gloov1.Destination, path string) *v1.Route {
	return &v1.Route{
		Matchers: []*matchers.Matcher{{
			PathSpecifier: &matchers.Matcher_Prefix{
				Prefix: path,
			},
		}},
		Action: &v1.Route_RouteAction{
			RouteAction: &gloov1.RouteAction{
				Destination: &gloov1.RouteAction_Single{
					Single: dest,
				},
			},
		},
	}
}
