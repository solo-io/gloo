package gateway_test

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/solo-io/gloo/test/helpers"
	"github.com/solo-io/go-utils/testutils"
	"gopkg.in/yaml.v2"
	"k8s.io/client-go/kubernetes"

	"github.com/solo-io/go-utils/testutils/helper"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	v1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gateway/pkg/defaults"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/go-utils/kubeutils"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/factory"
	"github.com/solo-io/solo-kit/pkg/errors"
	"github.com/solo-io/solo-projects/projects/apiserver/test/harness"
	qm "github.com/solo-io/solo-projects/projects/apiserver/test/queries"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"

	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/kube"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"k8s.io/client-go/rest"
)

var _ = Describe("Installing gloo in gateway mode", func() {

	var (
		ctx        context.Context
		cancel     context.CancelFunc
		cfg        *rest.Config
		kubeClient kubernetes.Interface

		gatewayClient        v1.GatewayClient
		virtualServiceClient v1.VirtualServiceClient
	)

	BeforeEach(func() {
		ctx, cancel = context.WithCancel(context.Background())

		var err error
		cfg, err = kubeutils.GetConfig("", "")
		Expect(err).NotTo(HaveOccurred())

		kubeClient, err = kubernetes.NewForConfig(cfg)
		Expect(err).NotTo(HaveOccurred())

		cache := kube.NewKubeCache(ctx)
		gatewayClientFactory := &factory.KubeResourceClientFactory{
			Crd:         v1.GatewayCrd,
			Cfg:         cfg,
			SharedCache: cache,
		}
		virtualServiceClientFactory := &factory.KubeResourceClientFactory{
			Crd:         v1.VirtualServiceCrd,
			Cfg:         cfg,
			SharedCache: cache,
		}

		gatewayClient, err = v1.NewGatewayClient(gatewayClientFactory)
		Expect(err).NotTo(HaveOccurred())

		virtualServiceClient, err = v1.NewVirtualServiceClient(virtualServiceClientFactory)
		Expect(err).NotTo(HaveOccurred())

	})

	AfterEach(func() {
		cancel()
		err := virtualServiceClient.Delete(testHelper.InstallNamespace, "vs", clients.DeleteOpts{})
		if err != nil && !errors.IsNotExist(err) {
			Expect(err).NotTo(HaveOccurred())
		}
	})

	It("can route request to upstream", func() {

		_, err := virtualServiceClient.Write(&v1.VirtualService{

			Metadata: core.Metadata{
				Name:      "vs",
				Namespace: testHelper.InstallNamespace,
			},
			VirtualHost: &gloov1.VirtualHost{
				Name:    "default",
				Domains: []string{"*"},
				Routes: []*gloov1.Route{{
					Matcher: &gloov1.Matcher{
						PathSpecifier: &gloov1.Matcher_Prefix{
							Prefix: "/",
						},
					},
					Action: &gloov1.Route_RouteAction{
						RouteAction: &gloov1.RouteAction{
							Destination: &gloov1.RouteAction_Single{
								Single: &gloov1.Destination{
									Upstream: core.ResourceRef{
										Namespace: testHelper.InstallNamespace,
										Name:      fmt.Sprintf("%s-%s-%v", testHelper.InstallNamespace, "testrunner", helper.TestRunnerPort)},
								},
							},
						},
					},
				}},
			},
		}, clients.WriteOpts{})
		Expect(err).NotTo(HaveOccurred())

		defaultGateway := defaults.DefaultGateway(testHelper.InstallNamespace)
		// wait for default gateway to be created
		Eventually(func() (*v1.Gateway, error) {
			return gatewayClient.Read(testHelper.InstallNamespace, defaultGateway.Metadata.Name, clients.ReadOpts{})
		}, "15s", "0.5s").Should(Not(BeNil()))

		gatewayProxy := "gateway-proxy"
		gatewayPort := int(80)
		testHelper.CurlEventuallyShouldRespond(helper.CurlOpts{
			Protocol:          "http",
			Path:              "/",
			Method:            "GET",
			Host:              gatewayProxy,
			Service:           gatewayProxy,
			Port:              gatewayPort,
			ConnectionTimeout: 10, // this is important, as the first curl call sometimes hangs indefinitely
		}, helper.SimpleHttpResponse, 1, time.Minute*2)
	})

	Context("virtual service in configured with SSL", func() {

		BeforeEach(func() {
			// get the certificate so it is generated in the background
			go helpers.Certificate()
		})

		AfterEach(func() {
			err := kubeClient.CoreV1().Secrets(testHelper.InstallNamespace).Delete("secret", nil)
			Expect(err).NotTo(HaveOccurred())
		})

		It("can route https request to upstream", func() {
			createdSecret, err := kubeClient.CoreV1().Secrets(testHelper.InstallNamespace).Create(helpers.GetKubeSecret("secret", testHelper.InstallNamespace))
			Expect(err).NotTo(HaveOccurred())

			_, err = virtualServiceClient.Write(&v1.VirtualService{

				Metadata: core.Metadata{
					Name:      "vs",
					Namespace: testHelper.InstallNamespace,
				},
				SslConfig: &gloov1.SslConfig{
					SslSecrets: &gloov1.SslConfig_SecretRef{
						SecretRef: &core.ResourceRef{
							Name:      createdSecret.ObjectMeta.Name,
							Namespace: createdSecret.ObjectMeta.Namespace,
						},
					},
				},
				VirtualHost: &gloov1.VirtualHost{
					Name:    "default",
					Domains: []string{"*"},
					Routes: []*gloov1.Route{{
						Matcher: &gloov1.Matcher{
							PathSpecifier: &gloov1.Matcher_Prefix{
								Prefix: "/",
							},
						},
						Action: &gloov1.Route_RouteAction{
							RouteAction: &gloov1.RouteAction{
								Destination: &gloov1.RouteAction_Single{
									Single: &gloov1.Destination{
										Upstream: core.ResourceRef{
											Namespace: testHelper.InstallNamespace,
											Name:      fmt.Sprintf("%s-%s-%v", testHelper.InstallNamespace, "testrunner", helper.TestRunnerPort)},
									},
								},
							},
						},
					}},
				},
			}, clients.WriteOpts{})
			Expect(err).NotTo(HaveOccurred())

			defaultGateway := defaults.DefaultGateway(testHelper.InstallNamespace)
			// wait for default gateway to be created
			Eventually(func() (*v1.Gateway, error) {
				return gatewayClient.Read(testHelper.InstallNamespace, defaultGateway.Metadata.Name, clients.ReadOpts{})
			}, "15s", "0.5s").Should(Not(BeNil()))

			gatewayProxy := "gateway-proxy"
			gatewayPort := int(443)
			caFile := ToFile(helpers.Certificate())
			//noinspection GoUnhandledErrorResult
			defer os.Remove(caFile)

			err = testutils.Kubectl("cp", caFile, testHelper.InstallNamespace+"/testrunner:/tmp/ca.crt")
			Expect(err).NotTo(HaveOccurred())

			testHelper.CurlEventuallyShouldRespond(helper.CurlOpts{
				Protocol:          "https",
				Path:              "/",
				Method:            "GET",
				Host:              gatewayProxy,
				Service:           gatewayProxy,
				Port:              gatewayPort,
				CaFile:            "/tmp/ca.crt",
				ConnectionTimeout: 10, // this is important, as the first curl call sometimes hangs indefinitely
			}, helper.SimpleHttpResponse, 1, time.Minute*2)
		})
	})

	Context("Call Apiserver", func() {
		apiServer := harness.ApiServer{
			Origin: "http://apiserver-ui:8088",
			// Origin: "http://localhost:8082",
			// Token:  "any_string",
		}
		addr := os.Getenv("APISERVER_ADDR")
		if addr != "" {
			apiServer.Origin = addr
		}

		It("should call apiserver correctly", func() {
			apiserverCurlOptions := func(body string) helper.CurlOpts {
				serviceName := "apiserver-ui"
				return helper.CurlOpts{
					Protocol:          "http",
					Path:              "/query",
					Method:            "POST",
					Host:              serviceName,
					Service:           serviceName,
					Port:              8088,
					ConnectionTimeout: 10, // this is important, as the first curl call sometimes hangs indefinitely
					Body:              body,
					WithoutStats:      true,
				}
			}
			unmarshalGqlResponse := func(r string) harness.ResponseErrors {
				re := harness.ResponseErrors{}
				err := yaml.Unmarshal([]byte(r), &re)
				ExpectWithOffset(1, err).NotTo(HaveOccurred())
				return re
			}

			By("Basic Query")
			testHelper.CurlEventuallyShouldRespond(
				apiserverCurlOptions(qm.BasicUpstreamsQuery),
				qm.BasicUpstreamsQueryMatch,
				1,
				time.Minute*2)

			// After the first call has succeeded we can trust that the server is ready
			By("Root Query")
			response, err := testHelper.Curl(apiserverCurlOptions(qm.RootQuery))
			Expect(err).NotTo(HaveOccurred())
			re := unmarshalGqlResponse(response)
			Expect(re.Errors).NotTo(HaveOccurred())
			nsInterfaceList, ok := re.Data["allNamespaces"]
			Expect(ok).To(BeTrue())
			nsList, ok := nsInterfaceList.([]interface{})
			Expect(ok).To(BeTrue())
			Expect(len(nsList)).To(BeNumerically(">", 0))

			// Run the recorded queries
			recordedQueries := qm.Queries_6362475465567435311
			recordedQueries = append(recordedQueries, qm.ReplaceNamespaces(qm.Queries_1731087844762345556, testHelper.InstallNamespace)...)
			// recordedQueries = append(recordedQueries, qm.Queries_to_be_recorded_later...)
			By(fmt.Sprintf("Run %v recorded queries", len(recordedQueries)))
			for i, q := range recordedQueries {
				By(fmt.Sprintf("Recorded query %v", i))
				response, err := testHelper.Curl(apiserverCurlOptions(q))
				Expect(err).NotTo(HaveOccurred())
				re := unmarshalGqlResponse(response)
				Expect(re.Errors).NotTo(HaveOccurred())
			}

		})
	})
})

func ToFile(content string) string {
	f, err := ioutil.TempFile("", "")
	ExpectWithOffset(1, err).NotTo(HaveOccurred())
	n, err := f.WriteString(content)
	ExpectWithOffset(1, err).NotTo(HaveOccurred())
	ExpectWithOffset(1, n).To(Equal(len(content)))
	_ = f.Close()
	return f.Name()
}
