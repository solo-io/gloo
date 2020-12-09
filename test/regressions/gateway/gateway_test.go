package gateway_test

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/extensions/transformation"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/core/matchers"
	transformation2 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/transformation"
	"github.com/solo-io/solo-projects/test/regressions"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	v1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"

	"github.com/solo-io/gloo/test/helpers"
	"github.com/solo-io/go-utils/testutils"
	"k8s.io/client-go/kubernetes"

	"github.com/solo-io/k8s-utils/testutils/helper"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	v2 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gateway/pkg/defaults"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/k8s-utils/kubeutils"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/factory"
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

		gatewayClient        v2.GatewayClient
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
			Crd:         v2.GatewayCrd,
			Cfg:         cfg,
			SharedCache: cache,
		}
		virtualServiceClientFactory := &factory.KubeResourceClientFactory{
			Crd:         v1.VirtualServiceCrd,
			Cfg:         cfg,
			SharedCache: cache,
		}

		gatewayClient, err = v2.NewGatewayClient(ctx, gatewayClientFactory)
		Expect(err).NotTo(HaveOccurred())

		virtualServiceClient, err = v1.NewVirtualServiceClient(ctx, virtualServiceClientFactory)
		Expect(err).NotTo(HaveOccurred())

	})

	AfterEach(func() {
		regressions.DeleteVirtualService(virtualServiceClient, testHelper.InstallNamespace, "vs", clients.DeleteOpts{Ctx: ctx, IgnoreNotExist: true})
		cancel()
	})

	It("can route request to upstream", func() {

		regressions.WriteVirtualService(ctx, testHelper, virtualServiceClient, nil, nil, nil)

		defaultGateway := defaults.DefaultGateway(testHelper.InstallNamespace)
		// wait for default gateway to be created
		Eventually(func() (*v2.Gateway, error) {
			return gatewayClient.Read(testHelper.InstallNamespace, defaultGateway.Metadata.Name, clients.ReadOpts{Ctx: ctx})
		}, "15s", "0.5s").Should(Not(BeNil()))

		gatewayPort := 80
		testHelper.CurlEventuallyShouldRespond(helper.CurlOpts{
			Protocol:          "http",
			Path:              regressions.TestMatcherPrefix,
			Method:            "GET",
			Host:              defaults.GatewayProxyName,
			Service:           defaults.GatewayProxyName,
			Port:              gatewayPort,
			ConnectionTimeout: 10, // this is important, as the first curl call sometimes hangs indefinitely
		}, helper.SimpleHttpResponse, 1, time.Minute*5)
	})

	Context("virtual service in configured with SSL", func() {

		BeforeEach(func() {
			// get the certificate so it is generated in the background
			go helpers.Certificate()
		})

		AfterEach(func() {
			err := kubeClient.CoreV1().Secrets(testHelper.InstallNamespace).Delete(ctx, "secret", metav1.DeleteOptions{})
			Expect(err).NotTo(HaveOccurred())
		})

		It("can route https request to upstream", func() {
			sslSecret := helpers.GetKubeSecret("secret", testHelper.InstallNamespace)
			createdSecret, err := kubeClient.CoreV1().Secrets(testHelper.InstallNamespace).Create(ctx, sslSecret, metav1.CreateOptions{})
			Expect(err).NotTo(HaveOccurred())

			Eventually(func() error {
				_, err := kubeClient.CoreV1().Secrets(sslSecret.Namespace).Get(ctx, sslSecret.Name, metav1.GetOptions{})
				return err
			}, "10s", "0.5s").Should(BeNil())
			time.Sleep(3 * time.Second) // Wait a few seconds so Gloo can pick up the secret, otherwise the webhook validation might fail

			sslConfig := &gloov1.SslConfig{
				SslSecrets: &gloov1.SslConfig_SecretRef{
					SecretRef: &core.ResourceRef{
						Name:      createdSecret.ObjectMeta.Name,
						Namespace: createdSecret.ObjectMeta.Namespace,
					},
				},
			}

			regressions.WriteVirtualService(ctx, testHelper, virtualServiceClient, nil, nil, sslConfig)

			defaultGateway := defaults.DefaultGateway(testHelper.InstallNamespace)
			// wait for default gateway to be created
			Eventually(func() (*v2.Gateway, error) {
				return gatewayClient.Read(testHelper.InstallNamespace, defaultGateway.Metadata.Name, clients.ReadOpts{})
			}, "15s", "0.5s").Should(Not(BeNil()))

			gatewayPort := 443
			caFile := ToFile(helpers.Certificate())
			//noinspection GoUnhandledErrorResult
			defer os.Remove(caFile)

			err = testutils.Kubectl("cp", caFile, testHelper.InstallNamespace+"/testrunner:/tmp/ca.crt")
			Expect(err).NotTo(HaveOccurred())

			testHelper.CurlEventuallyShouldRespond(helper.CurlOpts{
				Protocol:          "https",
				Path:              regressions.TestMatcherPrefix,
				Method:            "GET",
				Host:              defaults.GatewayProxyName,
				Service:           defaults.GatewayProxyName,
				Port:              gatewayPort,
				CaFile:            "/tmp/ca.crt",
				ConnectionTimeout: 10, // this is important, as the first curl call sometimes hangs indefinitely
			}, helper.SimpleHttpResponse, 1, time.Minute*2)
		})
	})

	It("rejects invalid inja template in transformation", func() {
		injaTransform := `{% if default(data.error.message, "") != "" %}400{% else %}{{ header(":status") }}{% endif %}`
		t := &transformation2.Transformations{
			ClearRouteCache: true,
			ResponseTransformation: &transformation.Transformation{
				TransformationType: &transformation.Transformation_TransformationTemplate{
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

		_, err := virtualServiceClient.Write(vs, clients.WriteOpts{Ctx: ctx})
		Expect(err).ToNot(HaveOccurred())

		err = virtualServiceClient.Delete(vs.Metadata.Namespace, vs.Metadata.Name, clients.DeleteOpts{Ctx: ctx})
		Expect(err).ToNot(HaveOccurred())

		// trim trailing "}", which should invalidate our inja template
		t.ResponseTransformation.GetTransformationTemplate().Headers[":status"].Text = strings.TrimSuffix(injaTransform, "}")

		_, err = virtualServiceClient.Write(vs, clients.WriteOpts{Ctx: ctx})
		Expect(err).To(MatchError(ContainSubstring("Failed to parse response template: Failed to parse " +
			"header template ':status': [inja.exception.parser_error] expected statement close, got '%'")))
	})

})

func getVirtualService(dest *gloov1.Destination, sslConfig *gloov1.SslConfig) *v1.VirtualService {
	return getVirtualServiceWithRoute(getRouteWithDest(dest, "/"), sslConfig)
}

func getVirtualServiceWithRoute(route *v1.Route, sslConfig *gloov1.SslConfig) *v1.VirtualService {
	return &v1.VirtualService{
		Metadata: core.Metadata{
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

func ToFile(content string) string {
	f, err := ioutil.TempFile("", "")
	ExpectWithOffset(1, err).NotTo(HaveOccurred())
	n, err := f.WriteString(content)
	ExpectWithOffset(1, err).NotTo(HaveOccurred())
	ExpectWithOffset(1, n).To(Equal(len(content)))
	_ = f.Close()
	return f.Name()
}
