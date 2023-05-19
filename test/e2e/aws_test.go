package e2e_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/memory"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/solo-io/gloo/projects/gloo/pkg/defaults"

	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/aws"

	"github.com/solo-io/solo-projects/test/services"
)

const (
	region               = "us-east-1"
	webIdentityTokenFile = "AWS_WEB_IDENTITY_TOKEN_FILE"
	jwtPrivateKey        = "JWT_PRIVATE_KEY"
	awsRoleArnSts        = "AWS_ROLE_ARN_STS"
	awsRoleArn           = "AWS_ROLE_ARN"
)

var _ = Describe("AWS Lambda ", func() {

	var (
		ctx    context.Context
		cancel context.CancelFunc

		testClients   services.TestClients
		envoyInstance *services.EnvoyInstance
		secret        *gloov1.Secret
		upstream      *gloov1.Upstream
		envoyPort     uint32
	)

	setupEnvoy := func() {
		ctx, cancel = context.WithCancel(context.Background())
		cache := memory.NewInMemoryResourceCache()

		testClients = services.GetTestClients(ctx, cache)
		testClients.GlooPort = int(services.AllocateGlooPort())

		var err error
		envoyInstance, err = envoyFactory.NewEnvoyInstance()
		Expect(err).NotTo(HaveOccurred())

		settings := &gloov1.Settings{}

		what := services.What{
			DisableGateway: true,
			DisableUds:     true,
			DisableFds:     false,
		}

		services.RunGlooGatewayUdsFdsOnPort(services.RunGlooGatewayOpts{Ctx: ctx, Cache: cache, LocalGlooPort: int32(testClients.GlooPort), What: what, Namespace: defaults.GlooSystem, Settings: settings})

		err = envoyInstance.Run(testClients.GlooPort)
		Expect(err).NotTo(HaveOccurred())

		envoyPort = defaults.HttpPort
	}

	addBasicCredentials := func() {

		// Look in ~/.aws/credentials for local AWS credentials
		// see the gloo OSS e2e test README for more information about configuring credentials for AWS e2e tests
		// https://github.com/solo-io/gloo/blob/main/test/e2e/README.md
		localAwsCredentials := credentials.NewSharedCredentials("", "")
		v, err := localAwsCredentials.Get()
		if err != nil {
			Fail("no AWS creds available")
		}

		var opts clients.WriteOpts

		accesskey := v.AccessKeyID
		secretkey := v.SecretAccessKey

		secret = &gloov1.Secret{
			Metadata: &core.Metadata{
				Namespace: "default",
				Name:      region,
			},
			Kind: &gloov1.Secret_Aws{
				Aws: &gloov1.AwsSecret{
					AccessKey: accesskey,
					SecretKey: secretkey,
				},
			},
		}

		_, err = testClients.SecretClient.Write(secret, opts)
		Expect(err).NotTo(HaveOccurred())
	}

	validateLambda := func(offset int, envoyPort uint32, substring string) {
		body := []byte("\"solo.io\"")

		Eventually(func(g Gomega) {
			// send a request with a body
			var buf bytes.Buffer
			buf.Write(body)

			url := fmt.Sprintf("http://%s:%d/1?param_a=value_1&param_b=value_b", "localhost", envoyPort)
			res, err := http.Post(url, "application/octet-stream", &buf)
			g.Expect(err).NotTo(HaveOccurred())
			defer res.Body.Close()
			g.Expect(res.StatusCode).To(Equal(http.StatusOK))

			body, err := io.ReadAll(res.Body)
			g.Expect(err).NotTo(HaveOccurred())

			g.Expect(string(body)).To(ContainSubstring(substring))
		}, "10s", "1s").Should(Succeed())
	}

	validateLambdaUppercase := func(envoyPort uint32) {
		validateLambda(2, envoyPort, "SOLO.IO")
	}

	addUpstream := func() {
		upstream = &gloov1.Upstream{
			Metadata: &core.Metadata{
				Namespace: "default",
				Name:      region,
			},
			UpstreamType: &gloov1.Upstream_Aws{
				Aws: &aws.UpstreamSpec{
					Region:    region,
					SecretRef: secret.Metadata.Ref(),
				},
			},
		}

		var opts clients.WriteOpts
		_, err := testClients.UpstreamClient.Write(upstream, opts)
		Expect(err).NotTo(HaveOccurred())

		Eventually(func() []*aws.LambdaFunctionSpec {
			us, err := testClients.UpstreamClient.Read(
				upstream.GetMetadata().Namespace,
				upstream.GetMetadata().Name,
				clients.ReadOpts{},
			)
			if err != nil {
				return nil
			}
			return us.GetAws().GetLambdaFunctions()
		}, "15s", "1s").ShouldNot(BeEmpty())
	}

	testProxy := func() {
		proxy := &gloov1.Proxy{
			Metadata: &core.Metadata{
				Name:      "proxy",
				Namespace: "default",
			},
			Listeners: []*gloov1.Listener{{
				Name:        "listener",
				BindAddress: "::",
				BindPort:    envoyPort,
				ListenerType: &gloov1.Listener_HttpListener{
					HttpListener: &gloov1.HttpListener{
						VirtualHosts: []*gloov1.VirtualHost{{
							Name:    "virt1",
							Domains: []string{"*"},
							Routes: []*gloov1.Route{{
								Action: &gloov1.Route_RouteAction{
									RouteAction: &gloov1.RouteAction{
										Destination: &gloov1.RouteAction_Single{
											Single: &gloov1.Destination{
												DestinationType: &gloov1.Destination_Upstream{
													Upstream: upstream.Metadata.Ref(),
												},
												DestinationSpec: &gloov1.DestinationSpec{
													DestinationType: &gloov1.DestinationSpec_Aws{
														Aws: &aws.DestinationSpec{
															LogicalName: "uppercase",
														},
													},
												},
											},
										},
									},
								},
							}},
						}},
					},
				},
			}},
		}

		var opts clients.WriteOpts
		_, err := testClients.ProxyClient.Write(proxy, opts)
		Expect(err).NotTo(HaveOccurred())

		validateLambdaUppercase(defaults.HttpPort)
	}

	createEchoProxy := func(unwrapAsApiGateway bool, wrapAsApiGateway bool) {
		proxy := &gloov1.Proxy{
			Metadata: &core.Metadata{
				Name:      "proxy",
				Namespace: "default",
			},
			Listeners: []*gloov1.Listener{{
				Name:        "listener",
				BindAddress: "::",
				BindPort:    envoyPort,
				ListenerType: &gloov1.Listener_HttpListener{
					HttpListener: &gloov1.HttpListener{
						VirtualHosts: []*gloov1.VirtualHost{{
							Name:    "virt1",
							Domains: []string{"*"},
							Routes: []*gloov1.Route{{
								Action: &gloov1.Route_RouteAction{
									RouteAction: &gloov1.RouteAction{
										Destination: &gloov1.RouteAction_Single{
											Single: &gloov1.Destination{
												DestinationType: &gloov1.Destination_Upstream{
													Upstream: upstream.Metadata.Ref(),
												},
												DestinationSpec: &gloov1.DestinationSpec{
													DestinationType: &gloov1.DestinationSpec_Aws{
														Aws: &aws.DestinationSpec{
															LogicalName:        "echo",
															UnwrapAsApiGateway: unwrapAsApiGateway,
															WrapAsApiGateway:   wrapAsApiGateway,
														},
													},
												},
											},
										},
									},
								},
							}},
						},
						},
					},
				},
			}},
		}

		var opts clients.WriteOpts
		_, err := testClients.ProxyClient.Write(proxy, opts)
		Expect(err).NotTo(HaveOccurred())
	}

	Context("Enterprise Lambda plugin", func() {

		BeforeEach(func() {
			setupEnvoy()
			addBasicCredentials()
			addUpstream()
		})

		AfterEach(func() {
			envoyInstance.Clean()
			cancel()
		})

		It("Can configure simple Lambda upstream", func() {
			testProxy()
		})

		Context("Enterprise-specific functionality", func() {
			It("Can configure unwrapAsApiGateway", func() {
				createEchoProxy(true, false)

				// format API gateway response.
				bodyString := "test"
				statusCode := 200
				headers := map[string]string{"Foo": "bar"}
				jsonHeaderStr, err := json.Marshal(headers)
				Expect(err).NotTo(HaveOccurred())
				apiGatewayResponse := fmt.Sprintf("{\"body\": \"%s\", \"statusCode\": %d, \"headers\":%s}", bodyString, statusCode, string(jsonHeaderStr))

				body := []byte(apiGatewayResponse)
				expectResponse := func(response *http.Response, body string, statusCode int, headers map[string]string) {
					Expect(response.StatusCode).To(Equal(statusCode))

					defer response.Body.Close()
					responseBody, err := io.ReadAll(response.Body)
					Expect(err).NotTo(HaveOccurred())
					Expect(string(responseBody)).To(Equal(body))
					for k, v := range headers {
						Expect(response.Header.Get(k)).To(Equal(v))
					}
				}

				Eventually(func(g Gomega) {
					// send a request with a body
					var err error
					var buf bytes.Buffer
					buf.Write(body)

					// send request to echo lambda, mimicking a service that generates an API gateway response.
					url := fmt.Sprintf("http://%s:%d/1?param_a=value_1&param_b=value_b", "localhost", defaults.HttpPort)
					res, err := http.Post(url, "application/octet-stream", &buf)
					g.Expect(err).NotTo(HaveOccurred())
					g.Expect(res.StatusCode).To(Equal(200))
					expectResponse(res, bodyString, statusCode, headers)
				}, "5s", "1s").Should(Succeed())
			})

			Context("Can configure wrapAsApiGateway", func() {
				It("works with a simple request", func() {
					createEchoProxy(false, true)

					// format API gateway request.
					bodyString := "test"
					body := []byte(bodyString)
					headers := map[string][]string{
						"single-value-header": {"value"},
						"multi-value-header":  {"value1", "value2"},
					}

					// expect that the response (which is identical to the payload sent to the lambda) is transformed appropriately.
					expectResponse := func(response *http.Response, body string, statusCode int, headers map[string][]string) {
						Expect(response.StatusCode).To(Equal(statusCode))

						defer response.Body.Close()
						responseBody, err := io.ReadAll(response.Body)
						Expect(err).NotTo(HaveOccurred())

						jsonResponseBody := make(map[string]interface{})
						err = json.Unmarshal(responseBody, &jsonResponseBody)
						Expect(err).NotTo(HaveOccurred())

						Expect(jsonResponseBody["body"]).To(Equal(bodyString))
						Expect(jsonResponseBody["httpMethod"]).To(Equal("GET"))
						Expect(jsonResponseBody["path"]).To(Equal("/1"))

						responseHeaders := jsonResponseBody["headers"].(map[string]interface{})
						Expect(responseHeaders["single-value-header"]).To(Equal(headers["single-value-header"][0]))

						responseMultiValueHeaders := jsonResponseBody["multiValueHeaders"].(map[string]interface{})
						responseMultiValueHeader := responseMultiValueHeaders["multi-value-header"].([]interface{})
						Expect(len(responseMultiValueHeader)).To(Equal(len(headers["multi-value-header"])))
						for i, v := range responseMultiValueHeader {
							Expect(v).To(Equal(headers["multi-value-header"][i]))
						}

						responseQueryStringParameters := jsonResponseBody["queryStringParameters"].(map[string]interface{})
						Expect(responseQueryStringParameters["param_a"]).To(Equal("value_1"))
						// note that we expect the value of param_b to be the last value assigned to it in the query string.
						Expect(responseQueryStringParameters["param_b"]).To(Equal("value_2"))

						responseMultiValueQueryStringParameters := jsonResponseBody["multiValueQueryStringParameters"].(map[string]interface{})
						responseMultiValueQueryStringParameter := responseMultiValueQueryStringParameters["param_b"].([]interface{})
						Expect(len(responseMultiValueQueryStringParameter)).To(Equal(2))
						Expect(responseMultiValueQueryStringParameter[0]).To(Equal("value_b"))
						Expect(responseMultiValueQueryStringParameter[1]).To(Equal("value_2"))
					}

					Eventually(func(g Gomega) {
						var buf bytes.Buffer
						buf.Write(body)

						client := http.DefaultClient
						var err error
						// send request to echo lambda, mimicking a service that generates an API gateway response.
						res, err := client.Do(&http.Request{
							Method: "GET",
							URL: &url.URL{
								Scheme:   "http",
								Host:     fmt.Sprintf("localhost:%d", defaults.HttpPort),
								Path:     "/1",
								RawQuery: "param_a=value_1&param_b=value_b&param_b=value_2",
							},
							Body:   io.NopCloser(&buf),
							Header: headers,
						})
						g.Expect(err).NotTo(HaveOccurred())
						expectResponse(res, bodyString, 200, headers)
					}, "5s", "1s").Should(Succeed())
				})

				It("works with a request with no body", func() {
					createEchoProxy(false, true)

					// expect that the response (which is identical to the payload sent to the lambda) is transformed appropriately.
					expectResponse := func(response *http.Response, statusCode int) {
						Expect(response.StatusCode).To(Equal(statusCode))

						defer response.Body.Close()
						responseBody, err := io.ReadAll(response.Body)
						Expect(err).NotTo(HaveOccurred())

						jsonResponseBody := make(map[string]interface{})
						err = json.Unmarshal(responseBody, &jsonResponseBody)
						Expect(err).NotTo(HaveOccurred())

						Expect(jsonResponseBody["body"]).To(Equal(""))
						Expect(jsonResponseBody["httpMethod"]).To(Equal("GET"))
						Expect(jsonResponseBody["path"]).To(Equal("/1"))

						responseQueryStringParameters := jsonResponseBody["queryStringParameters"].(map[string]interface{})
						Expect(responseQueryStringParameters["param_a"]).To(Equal("value_1"))
						// note that we expect the value of param_b to be the last value assigned to it in the query string.
						Expect(responseQueryStringParameters["param_b"]).To(Equal("value_2"))

						responseMultiValueQueryStringParameters := jsonResponseBody["multiValueQueryStringParameters"].(map[string]interface{})
						responseMultiValueQueryStringParameter := responseMultiValueQueryStringParameters["param_b"].([]interface{})
						Expect(len(responseMultiValueQueryStringParameter)).To(Equal(2))
						Expect(responseMultiValueQueryStringParameter[0]).To(Equal("value_b"))
						Expect(responseMultiValueQueryStringParameter[1]).To(Equal("value_2"))
					}

					Eventually(func(g Gomega) {
						client := http.DefaultClient
						var err error

						res, err := client.Do(&http.Request{
							Method: "GET",
							URL: &url.URL{
								Scheme:   "http",
								Host:     fmt.Sprintf("localhost:%d", defaults.HttpPort),
								Path:     "/1",
								RawQuery: "param_a=value_1&param_b=value_b&param_b=value_2",
							},
						})

						g.Expect(err).NotTo(HaveOccurred())
						expectResponse(res, 200)
					}, "5s", "1s").Should(Succeed())
				})
			})
		})

	})
})
