package e2e_test

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/solo-io/gloo/test/helpers"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"

	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/core/matchers"

	"github.com/fgrosse/zaptest"
	"github.com/golang/protobuf/ptypes/wrappers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/extensions/transformers/xslt"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/transformation"
	"github.com/solo-io/gloo/projects/gloo/pkg/defaults"
	"github.com/solo-io/go-utils/contextutils"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/memory"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	. "github.com/solo-io/solo-projects/test/e2e/transformation_helpers"
	"github.com/solo-io/solo-projects/test/services"
	"github.com/solo-io/solo-projects/test/v1helpers"
)

var _ = Describe("XSLT Transformer E2E", func() {

	var (
		ctx           context.Context
		cancel        context.CancelFunc
		testClients   services.TestClients
		envoyInstance *services.EnvoyInstance
		testUpstream  *v1helpers.TestUpstream
		envoyPort     uint32
		transform     *transformation.TransformationStages
	)

	BeforeEach(func() {

		logger := zaptest.LoggerWriter(GinkgoWriter)
		contextutils.SetFallbackLogger(logger.Sugar())

		ctx, cancel = context.WithCancel(context.Background())
		cache := memory.NewInMemoryResourceCache()

		testClients = services.GetTestClients(ctx, cache)
		testClients.GlooPort = int(services.AllocateGlooPort())

		what := services.What{
			DisableGateway: true,
			DisableFds:     true,
			DisableUds:     true,
		}

		services.RunGlooGatewayUdsFdsOnPort(ctx, cache, int32(testClients.GlooPort), what, "gloo-system", nil, nil, &gloov1.Settings{
			Gateway: &gloov1.GatewayOptions{
				Validation: &gloov1.GatewayOptions_ValidationOptions{
					DisableTransformationValidation: &wrappers.BoolValue{Value: true},
				},
			},
		})
	})

	setupProxy := func() {
		var err error
		envoyInstance, err = envoyFactory.NewEnvoyInstance()
		Expect(err).NotTo(HaveOccurred())
		err = envoyInstance.Run(testClients.GlooPort)
		Expect(err).NotTo(HaveOccurred())

		testUpstream = v1helpers.NewTestHttpUpstream(ctx, envoyInstance.LocalAddr())

		_, err = testClients.UpstreamClient.Write(testUpstream.Upstream, clients.WriteOpts{})
		Expect(err).NotTo(HaveOccurred())

		helpers.EventuallyResourceAccepted(func() (resources.InputResource, error) {
			return testClients.UpstreamClient.Read(testUpstream.Upstream.Metadata.Namespace, testUpstream.Upstream.Metadata.Name, clients.ReadOpts{})
		})

		envoyPort = defaults.HttpPort
		proxy := getProxyXsltTransform(envoyPort, transform, testUpstream.Upstream.Metadata.Ref())

		_, err = testClients.ProxyClient.Write(proxy, clients.WriteOpts{})
		Expect(err).NotTo(HaveOccurred())

		helpers.EventuallyResourceAccepted(func() (resources.InputResource, error) {
			return testClients.ProxyClient.Read(proxy.Metadata.Namespace, proxy.Metadata.Name, clients.ReadOpts{})
		})

		request, err := http.NewRequest("GET", fmt.Sprintf("http://localhost:%d", envoyInstance.AdminPort), nil)
		Expect(err).NotTo(HaveOccurred())
		client := &http.Client{}
		Eventually(func() (int, error) {
			response, err := client.Do(request)
			if response == nil {
				return 0, err
			}
			return response.StatusCode, err
		}, 20*time.Second, 1*time.Second).Should(Equal(200))

	}

	AfterEach(func() {
		if envoyInstance != nil {
			_ = envoyInstance.Clean()
		}
		cancel()
	})

	Context("With envoy", func() {

		getXsltTransform := func(transform, setContentType string, nonXmlTransform bool) *transformation.TransformationStages {
			return &transformation.TransformationStages{
				Regular: &transformation.RequestResponseTransformations{
					RequestTransforms: []*transformation.RequestMatch{
						{
							Matcher: &matchers.Matcher{
								PathSpecifier: &matchers.Matcher_Prefix{Prefix: "/"},
							},
							RequestTransformation: &transformation.Transformation{
								TransformationType: &transformation.Transformation_XsltTransformation{
									XsltTransformation: &xslt.XsltTransformation{
										Xslt:            transform,
										SetContentType:  setContentType,
										NonXmlTransform: nonXmlTransform,
									},
								},
							},
						},
					},
				},
			}
		}

		DescribeTable("with envoy",
			func(xsltTransformation, setContentType string, nonXmlTransform bool, body, expectedBody string) {
				transform = getXsltTransform(xsltTransformation, setContentType, nonXmlTransform)
				setupProxy()
				request, err := http.NewRequest("POST", fmt.Sprintf("http://localhost:%d/", envoyPort), strings.NewReader(body))
				ExpectWithOffset(1, err).NotTo(HaveOccurred())
				var response *http.Response
				EventuallyWithOffset(1, func() (int, error) {
					response, err = http.DefaultClient.Do(request)
					if response == nil {
						return 0, err
					}
					return response.StatusCode, err
				}, 10*time.Second, 1*time.Second).Should(Equal(200))
				bodyBytes, err := ioutil.ReadAll(response.Body)
				ExpectWithOffset(1, err).NotTo(HaveOccurred())
				ExpectWithOffset(1, string(bodyBytes)).To(Equal(expectedBody))
			},
			Entry("json -> xml", JsonToXmlTransform, "application/xml", true, AutomobileJson, AutomobileXml),
			Entry("xml -> json", XmlToJsonTransform, "application/json", false, CarsXml, CarsJson),
			Entry("soap xml -> xml", XmlToXmlTransform, "applications/xml", false, SoapMessageXml, TransformedSoapXml),
		)

		It("rejects invalid input correctly", func() {
			expectBadRequest := func(body string) {
				request, err := http.NewRequest("POST", fmt.Sprintf("http://localhost:%d/", envoyPort), strings.NewReader(body))
				ExpectWithOffset(1, err).NotTo(HaveOccurred())
				var response *http.Response
				EventuallyWithOffset(1, func() (int, error) {
					response, _ = http.DefaultClient.Do(request)
					if response == nil {
						return 0, err
					}
					return response.StatusCode, err
				}, 15*time.Second, 1*time.Second).Should(Equal(400))
				bodyBytes, err := ioutil.ReadAll(response.Body)
				ExpectWithOffset(1, string(bodyBytes)).To(ContainSubstring("bad request"))
			}

			transform = getXsltTransform(XmlToXmlTransform, "application/xml", false)
			setupProxy()
			// Invalid body
			expectBadRequest(`<This is invalid xml />`)
			expectBadRequest("")
		})
	})

})

func getProxyXsltTransform(envoyPort uint32, transform *transformation.TransformationStages, upstream *core.ResourceRef) *gloov1.Proxy {
	var vhosts []*gloov1.VirtualHost

	vhost := &gloov1.VirtualHost{
		Name:    "virt1",
		Domains: []string{"*"},
		Options: &gloov1.VirtualHostOptions{
			StagedTransformations: transform,
		},
		Routes: []*gloov1.Route{{
			Action: &gloov1.Route_RouteAction{
				RouteAction: &gloov1.RouteAction{
					Destination: &gloov1.RouteAction_Single{
						Single: &gloov1.Destination{
							DestinationType: &gloov1.Destination_Upstream{
								Upstream: upstream,
							},
						},
					},
				},
			}}},
	}

	vhosts = append(vhosts, vhost)

	p := &gloov1.Proxy{
		Metadata: &core.Metadata{
			Name:      "proxy",
			Namespace: "default",
		},
		Listeners: []*gloov1.Listener{{
			Name:        "listener",
			BindAddress: "0.0.0.0",
			BindPort:    envoyPort,
			ListenerType: &gloov1.Listener_HttpListener{
				HttpListener: &gloov1.HttpListener{
					VirtualHosts: vhosts,
				},
			},
		}},
	}

	return p
}
