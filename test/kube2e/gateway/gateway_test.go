package gateway_test

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"strings"

	"time"

	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/solo-io/gloo/pkg/cliutil/install"
	v1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/extensions/transformation"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/core/matchers"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/gloosnapshot"
	glooStatic "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/static"
	glootransformation "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/transformation"
	kubernetes2 "github.com/solo-io/gloo/projects/gloo/pkg/plugins/kubernetes"
	"google.golang.org/protobuf/types/known/wrapperspb"
	corev1 "k8s.io/api/core/v1"

	static_plugin_gloo "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/static"
	defaults2 "github.com/solo-io/gloo/projects/gloo/pkg/defaults"
	"github.com/solo-io/gloo/projects/gloo/pkg/translator"
	"github.com/solo-io/gloo/test/gomega/transforms"
	glooKube2e "github.com/solo-io/gloo/test/kube2e"

	"github.com/solo-io/solo-projects/test/kube2e"

	"github.com/solo-io/gloo/test/helpers"
	"github.com/solo-io/go-utils/testutils"
	"github.com/solo-io/k8s-utils/testutils/helper"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/solo-io/gloo/projects/gateway/pkg/defaults"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	gloossl "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/ssl"
	"github.com/solo-io/k8s-utils/kubeutils"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"

	osskube2e "github.com/solo-io/gloo/test/kube2e"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-kit/test/setup"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
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
		snapshotWriter = helpers.NewSnapshotWriter(resourceClientset)
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

			Context("rate limit config referenced by virtual service", func() {
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

				getRateLimitYaml := func(requestsPerUnit int) string {
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
          requestsPerUnit: %d
          unit: MINUTE
    rateLimits:
      - actions:
        - genericKey:
            descriptorValue: bar
`
					return fmt.Sprintf(rateLimitYaml, requestsPerUnit)
				}

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

				It("rejects deleting rate limit config referenced on a Virtual Service", func() {
					out := checkThatVSRefCantBeDeleted(getRateLimitYaml(1), vsYaml)
					Expect(out).To(ContainSubstring("Error from server"))
					Expect(out).To(ContainSubstring("admission webhook"))
					Expect(out).To(ContainSubstring("could not find RateLimitConfig resource with name"))
				})

				It("accepts updating rate limit config referenced by a Virtual Service", func() {
					var out []byte
					var err error

					Eventually(func() error {
						out, err = install.KubectlApplyOut([]byte(getRateLimitYaml(100)))
						return err
					}, "5s", "1s").Should(BeNil())
					Expect(string(out)).ToNot(ContainSubstring("Error from server"))

					Eventually(func() error {
						out, err = install.KubectlApplyOut([]byte(vsYaml))
						return err
					}, "5s", "1s").Should(BeNil())
					Expect(string(out)).ToNot(ContainSubstring("Error from server"))

					Eventually(func() error {
						out, err = install.KubectlApplyOut([]byte(getRateLimitYaml(200)))
						return err
					}, "5s", "1s").Should(BeNil())
					Expect(string(out)).ToNot(And(ContainSubstring("Error from server"), ContainSubstring("unchanged")))

				})
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
			tcpEchoClusterName  string
			tcpEchoShutdownFunc func()
			httpEcho            helper.TestRunner

			// define some ciphers which we will use both for configuration in
			// the server and defining test requests that will be sent by the
			// client

			// passthroughCipher: configured in envoy for passthrough - must be
			// supported by the tcp-echo-tls server, but does not need to be
			// supported by envoy
			passthroughCipher = "AES128-SHA"
			// terminatedCipher: configure for TLS termination in envoy
			// (must be supported by envoy)
			terminatedCipher = "ECDHE-RSA-AES256-GCM-SHA384"
			// nonTerminatedCipher: a cipher used for testing that's not
			// configured for passthrough or tls termination in envoy.
			// connections made with this cipher and no others should be
			// terminated if dcp is enabled
			nonTerminatedCipher = "ECDHE-RSA-CHACHA20-POLY1305"
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

		var err error
		var virtualservice *v1.VirtualService
		var hybridGateway *v1.Gateway
		BeforeEach(func() {
			// We are disabling tests temporarily, until we resolve:
			// https://github.com/solo-io/solo-projects/issues/5156
			Skip("Skipping DCP Tests")

			caFile := glooKube2e.ToFile(helpers.Certificate())
			//goland:noinspection GoUnhandledErrorResult
			defer os.Remove(caFile)
			err = setup.Kubectl("cp", caFile, testHelper.InstallNamespace+"/testrunner:/tmp/ca.crt")
			Expect(err).NotTo(HaveOccurred())
			exposePortOnGwProxyService(hybridProxyServicePort)

			tcpEchoClusterName = translator.UpstreamToClusterName(&core.ResourceRef{
				Namespace: testHelper.InstallNamespace,
				Name:      kubernetes2.UpstreamName(testHelper.InstallNamespace, helper.HttpEchoName, helper.HttpEchoPort),
			})

			httpEchoUpstream := &gloov1.Upstream{
				Metadata: &core.Metadata{
					Name:      "http-echo",
					Namespace: testHelper.InstallNamespace,
				},
				UpstreamType: &gloov1.Upstream_Static{
					Static: &static_plugin_gloo.UpstreamSpec{
						Hosts: []*static_plugin_gloo.Host{{
							Addr: helper.HttpEchoName,
							Port: helper.HttpEchoPort,
						}},
						UseTls: wrapperspb.Bool(false),
					},
				},
			}

			sslSecret := helpers.GetKubeSecret("secret", testHelper.InstallNamespace)
			createdSecret, err := resourceClientset.KubeClients().CoreV1().Secrets(testHelper.InstallNamespace).Create(ctx, sslSecret, metav1.CreateOptions{})
			Expect(err).NotTo(HaveOccurred())

			domain := fmt.Sprintf("%s:%d", defaults.GatewayProxyName, defaults2.HybridPort)
			virtualservice = helpers.NewVirtualServiceBuilder().
				WithName(defaults.GatewayProxyName).
				WithNamespace(testHelper.InstallNamespace).
				WithDomain(domain).
				WithRoutePrefixMatcher("route", "/").
				WithRouteActionToUpstream("route", httpEchoUpstream).
				WithSslConfig(&gloossl.SslConfig{
					Parameters: &gloossl.SslParameters{
						// only treat these ciphers as natively supported
						CipherSuites: []string{terminatedCipher},
					},
					SslSecrets: &gloossl.SslConfig_SecretRef{
						SecretRef: &core.ResourceRef{
							Name:      createdSecret.ObjectMeta.Name,
							Namespace: createdSecret.ObjectMeta.Namespace,
						},
					},
				}).
				Build()

			// Create a tcp upstream pointing to tcp-echo-tls
			tcpBinUp := &gloov1.Upstream{
				Metadata: &core.Metadata{
					Name:      "tcpb",
					Namespace: testHelper.InstallNamespace,
				},
				UpstreamType: &gloov1.Upstream_Static{
					Static: &static_plugin_gloo.UpstreamSpec{
						Hosts: []*static_plugin_gloo.Host{{
							Addr: "tcp-echo-tls",
							Port: 443,
						}},
						UseTls: wrapperspb.Bool(false),
					},
				},
			}
			tcpRef := &core.ResourceRef{
				Namespace: testHelper.InstallNamespace,
				Name:      "tcpb",
			}

			// Create a MatchableHttpGateway. Since this gateway will proxy an
			// ssl-enabled virtual serivce, we *must* specify an SslConfig
			// object to indicate that ssl-enabled virtual services can be
			// matched by this gateway. see
			// https://github.com/solo-io/gloo/blob/dd1db73fb4b3bee1dfbd01d685982b3995325ecc/projects/gateway/pkg/translator/gateway_selector.go#L73
			matchableHttpGateway := &v1.MatchableHttpGateway{
				Metadata: &core.Metadata{
					Name:      "matchable-http-gateway",
					Namespace: testHelper.InstallNamespace,
				},
				HttpGateway: &v1.HttpGateway{
					VirtualServices: []*core.ResourceRef{
						{
							Name:      virtualservice.Metadata.Name,
							Namespace: virtualservice.Metadata.Namespace,
						},
					},
				},
				Matcher: &v1.MatchableHttpGateway_Matcher{
					SslConfig: &gloossl.SslConfig{},
				},
			}

			// Create a MatchableTcpGateway
			matchableTcpGateway := &v1.MatchableTcpGateway{
				Metadata: &core.Metadata{
					Name:      "matchable-tcp-gateway",
					Namespace: testHelper.InstallNamespace,
				},
				Matcher: &v1.MatchableTcpGateway_Matcher{

					PassthroughCipherSuites: []string{passthroughCipher},
				},
				TcpGateway: &v1.TcpGateway{
					TcpHosts: []*gloov1.TcpHost{{
						Name: tcpEchoClusterName,
						Destination: &gloov1.TcpHost_TcpAction{
							Destination: &gloov1.TcpHost_TcpAction_Single{
								Single: &gloov1.Destination{
									DestinationType: &gloov1.Destination_Upstream{
										Upstream: tcpRef,
									},
								},
							},
						},
					}},
				},
			}

			// Create a HybridGateway that references that MatchableHttpGateway
			hybridGateway = &v1.Gateway{
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
							// see comment above on the MatchableHttpGateway
							// about why SslConfig is needed here
							SslConfig:             &gloossl.SslConfig{},
							PreventChildOverrides: true,
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

			glooKube2e.UpdateSettings(ctx, func(settings *gloov1.Settings) {
				Expect(settings.GetGateway().IsolateVirtualHostsBySslConfig).NotTo(BeNil())
				settings.GetGateway().IsolateVirtualHostsBySslConfig = wrapperspb.Bool(true)
			}, testHelper.InstallNamespace)
			glooResources.HttpGateways = v1.MatchableHttpGatewayList{matchableHttpGateway}
			glooResources.TcpGateways = v1.MatchableTcpGatewayList{matchableTcpGateway}
			glooResources.Gateways = v1.GatewayList{hybridGateway}
			glooResources.Upstreams = gloov1.UpstreamList{tcpBinUp, httpEchoUpstream}
			glooResources.VirtualServices = append(glooResources.VirtualServices, virtualservice)
			httpEcho, err = helper.NewEchoHttp(testHelper.InstallNamespace)
			Expect(err).NotTo(HaveOccurred())
			err = httpEcho.Deploy(2 * time.Minute)
			Expect(err).NotTo(HaveOccurred())
			tcpEchoShutdownFunc = createTcpEchoTls()
		})

		AfterEach(func() {
			glooKube2e.UpdateSettings(ctx, func(settings *gloov1.Settings) {
				Expect(settings.GetGateway().IsolateVirtualHostsBySslConfig).NotTo(BeNil())
				settings.GetGateway().IsolateVirtualHostsBySslConfig = wrapperspb.Bool(false)
			}, testHelper.InstallNamespace)
			tcpEchoShutdownFunc()
			httpEcho.Terminate()

			// delete the ssl secret
			err := resourceClientset.KubeClients().CoreV1().Secrets(testHelper.InstallNamespace).Delete(ctx, "secret", metav1.DeleteOptions{})
			Expect(err).NotTo(HaveOccurred())

			// Delete http echo service
			err = testutils.Kubectl("delete", "service", "-n", testHelper.InstallNamespace, helper.HttpEchoName, "--grace-period=0")
			Expect(err).NotTo(HaveOccurred())
			cancel()
		})

		executeRequest := func(ciphers []string, expectedResult func([]string)) {
			cipherString := strings.Join(ciphers, ":")
			curlCmd := fmt.Sprintf("curl -sk --tlsv1.2 https://%s:%d --ciphers %s", defaults.GatewayProxyName, defaults2.HybridPort, cipherString)
			expectedResult(strings.Split(curlCmd, " "))
		}
		expectTcpPassthrough := func(curlArgs []string) {
			// testHelper.Curl does not support the `--ciphers` flag so let's
			// just build the command manually
			Eventually(func(g Gomega) {
				result, err := testHelper.Exec(curlArgs...)
				g.Expect(err).NotTo(HaveOccurred())
				/* sample response:
				   {
				     "name": "Service",
				     "uri": "/",
				     "type": "HTTP",
				     "ip_addresses": [
				   	"10.244.0.155"
				     ],
				     "start_time": "2023-05-24T15:34:26.834787",
				     "end_time": "2023-05-24T15:34:26.834813",
				     "duration": "26.615µs",
				     "body": "Hello World",
				     "code": 200
				   }
				*/
				g.Expect(transforms.WithJsonBody()([]byte(result))).To(And(
					HaveKeyWithValue("code", float64(200)),
					HaveKeyWithValue("body", "Hello World"),
				))
			}, "90s", "5s").Should(Succeed())
		}
		expectTlsTermination := func(curlArgs []string) {
			curlArgs = append(curlArgs, "-d", "test_key=test_value")
			Eventually(func(g Gomega) {
				result, err := testHelper.Exec(curlArgs...)
				g.Expect(err).NotTo(HaveOccurred())
				/* sample response:
				   {
				     "path": "/",
				     "headers": {
				   	"host": "gateway-proxy:8087",
				   	"user-agent": "curl/7.58.0",
				   	"accept": "* /*",
				   	"content-length": "15",
				   	"content-type": "application/x-www-form-urlencoded",
				   	"x-forwarded-proto": "https",
				   	"x-request-id": "9bf32163-57ce-4035-a26c-1b622dcbe752",
				   	"x-envoy-expected-rq-timeout-ms": "15000"
				     },
				     "method": "POST",
				     "body": {
				   	"test_key": "test_value"
				     },
				     "fresh": false,
				     "hostname": "gateway-proxy",
				     "ip": "::ffff:10.244.0.121",
				     "ips": [],
				     "protocol": "http",
				     "query": {},
				     "subdomains": [],
				     "xhr": false
				   }
				*/
				g.Expect(transforms.WithJsonBody()([]byte(result))).To(HaveKeyWithValue("body", map[string]any{"test_key": "test_value"}))
			}, "90s", "5s").Should(Succeed())
		}
		expectSslHandshakeFailure := func(curlArgs []string) {
			// wait for TLS termination case to complete - once it does, the
			// test should be set up properly. then expect the test with
			// curlArgs to fail
			executeRequest([]string{terminatedCipher}, expectTlsTermination)
			_, err := testHelper.Exec(curlArgs...)
			Expect(err).To(MatchError(ContainSubstring("command terminated with exit code 35")))
		}

		// helper to add nonTerminatedCipher to the list of terminated ciphers
		addNonTerminatedCipher := func() {
			helpers.PatchResourceWithOffset(1, ctx, virtualservice.GetMetadata().Ref(), func(resource resources.Resource) resources.Resource {
				virtualservice := resource.(*v1.VirtualService)

				virtualservice.SslConfig.Parameters.CipherSuites = append(
					virtualservice.SslConfig.Parameters.CipherSuites, nonTerminatedCipher)
				return virtualservice
			}, resourceClientset.VirtualServiceClient().BaseClient())
		}

		// helper to disable dcp by deleting DelegatedTcpGateway from the hybrid gateway
		disableDcp := func() {
			helpers.PatchResourceWithOffset(1, ctx, hybridGateway.GetMetadata().Ref(), func(resource resources.Resource) resources.Resource {
				gateway := resource.(*v1.Gateway)
				gateway.GetHybridGateway().DelegatedTcpGateways = nil
				return gateway
			}, resourceClientset.GatewayClient().BaseClient())
		}

		DescribeTable("dcp cipher combination tests",
			/* Table of tests for DCP cipher combinations

			   We initialise the server with a DCP configuration and then test
			   combinations of ciphers and how we expect them to be handled.

			   Configured ciphers:
			   * passthroughCipher (AES128-SHA): configured for passthrough
			   * terminatedCipher (ECDHE-RSA-AES256-GCM-SHA384): configured for
			       tls termination in envoy
			   * nonTerminatedCipher (ECDHE-RSA-CHACHA20-POLY1305): not configured
			       (should be rejected by Envoy)

			   Possible results:
			   * expectTcpPassthrough: Envoy should proxy the request to the
			       configured tcp-echo-tls server
			   * expectTlsTermination: Envoy should terminate TLS and send the
			       request to the configured HTTP echo server
			   * expectSslHandshakeFailure: connection should error. This callback
			       first performs the expectTlsTermination test to ensure that the
			       server has been configured and then performs the requested test

			   In the cases of the expectTlsTermination and expectTcpPassthrough,
			   the response bodies of the 2 configured servers are slightly
			   different, so the different assertions on the response bodies are
			   sufficient to determine which server received the request.

			   Then, we test combinations of ciphers as well. When all 3 ciphers
			   are present, we expect TLS termination in Envoy. When
			   passthroughCipher and nonTerminatedCipher are present, we expect
			   passthrough.

			   Finally, we can also change the server with different
			   configurations (namely, disabling DCP and re-enabling
			   nonTerminatedCipher) to ensure the correct behaviour in different
			   server configurations.
			*/
			func(ciphers []string, expectedResult func([]string), beforeTest []func()) {
				// --tlsv1.2 flag is necessary because without it, curl
				// automatically adds all tlsv1.3 ciphers

				// run any additional setup functions before testing the
				// requests. there is no current check to ensure that
				// configurations are fully applied (with the exception of the
				// call to expectTlsTermination in expectSslHandshakeFailure),
				// so the assertions themselves should be formulated such that
				// they only pass once the configuration has been fully applied
				for _, beforeTestFunc := range beforeTest {
					beforeTestFunc()
				}
				executeRequest(ciphers, expectedResult)
			},
			Entry("passthrough cipher should be passed through",
				[]string{passthroughCipher}, expectTcpPassthrough, nil),
			Entry("terminated cipher should be terminated in envoy",
				[]string{terminatedCipher}, expectTlsTermination, nil),
			Entry("nonTerminatedCipher (not configured for passthrough or termination) should result in a ssl handshake failure",
				[]string{nonTerminatedCipher}, expectSslHandshakeFailure, nil),
			Entry("client provides all 3 kinds of ciphers; TLS termination should be preferred over passthrough or failure",
				[]string{passthroughCipher, terminatedCipher, nonTerminatedCipher}, expectTlsTermination, nil),
			Entry("client provides passthrough and non-terminated ciphers; passthrough should occur",
				[]string{passthroughCipher, nonTerminatedCipher}, expectTcpPassthrough, nil),
			Entry("TLS termination occurs for nonTerminatedCipher after it is added to the terminated cipher list",
				[]string{nonTerminatedCipher}, expectTlsTermination,
				[]func(){addNonTerminatedCipher},
			),
			Entry("still cannot connect with nonTerminatedCipher when dcp is disabled",
				[]string{nonTerminatedCipher}, expectSslHandshakeFailure,
				[]func(){disableDcp},
			),
			Entry("can connect with nonTerminatedCipher when dcp is disabled and it is added to the terminated cipher list",
				[]string{nonTerminatedCipher}, expectTlsTermination,
				[]func(){disableDcp, addNonTerminatedCipher},
			),
		)
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

/* Create a TCP-based echo server that operates over TLS. Return a function to
* clean the services up */
func createTcpEchoTls() func() {
	tcpEchoTls := "tcp-echo-tls"
	cfg, err := kubeutils.GetConfig("", "")
	Expect(err).NotTo(HaveOccurred())
	kube, err := kubernetes.NewForConfig(cfg)
	Expect(err).NotTo(HaveOccurred())

	labels := map[string]string{}
	labels["gloo"] = tcpEchoTls
	metadata := metav1.ObjectMeta{
		Name:      tcpEchoTls,
		Namespace: testHelper.InstallNamespace,
		Labels:    labels,
	}
	zero := int64(0)
	pod, err := kube.CoreV1().Pods(testHelper.InstallNamespace).Create(context.TODO(), &corev1.Pod{
		ObjectMeta: metadata,
		Spec: corev1.PodSpec{
			TerminationGracePeriodSeconds: &zero,
			Containers: []corev1.Container{
				{
					Image:           "gcr.io/solo-test-236622/tcp-echo-tls:0.1.0",
					ImagePullPolicy: corev1.PullIfNotPresent,
					Name:            tcpEchoTls,
				},
			},
		},
	}, metav1.CreateOptions{})
	Expect(err).NotTo(HaveOccurred())

	service, err := kube.CoreV1().Services(testHelper.InstallNamespace).Create(context.Background(), &corev1.Service{
		ObjectMeta: metadata,
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				"gloo": tcpEchoTls,
			},
			Ports: []corev1.ServicePort{
				{
					Name:     tcpEchoTls,
					Protocol: corev1.ProtocolTCP,
					Port:     443,
				},
			},
		},
	}, metav1.CreateOptions{})
	Expect(err).NotTo(HaveOccurred())

	timeout := 2 * time.Minute
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	err = testutils.WaitPodsRunning(ctx, time.Second, testHelper.InstallNamespace, "gloo="+tcpEchoTls)
	Expect(err).NotTo(HaveOccurred())

	return func() {
		kube.CoreV1().Pods(testHelper.InstallNamespace).Delete(context.Background(), pod.GetObjectMeta().GetName(), metav1.DeleteOptions{})
		kube.CoreV1().Services(testHelper.InstallNamespace).Delete(context.Background(), service.GetObjectMeta().GetName(), metav1.DeleteOptions{})
	}
}
