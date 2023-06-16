package e2e_test

import (
	"fmt"
	"net/http"
	"os"

	errors "github.com/rotisserie/eris"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/solo-io/solo-projects/test/services"

	"github.com/solo-io/gloo/test/gomega/transforms"

	"github.com/aws/aws-sdk-go/aws/credentials"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	v1gw "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/aws"
	"github.com/solo-io/gloo/test/gomega/matchers"
	"github.com/solo-io/gloo/test/helpers"
	"github.com/solo-io/gloo/test/testutils"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-projects/test/e2e"
)

const (
	region = "us-east-1"
)

var _ = Describe("AWS Lambda", FlakeAttempts(5), func() {
	var (
		testContext *e2e.TestContext
		upstream    *gloov1.Upstream
	)

	BeforeEach(func() {
		testContext = testContextFactory.NewTestContext(testutils.AwsCredentials())
		testContext.BeforeEach()
		testContext.SetRunServices(services.What{
			DisableGateway: false,
			DisableUds:     true,
			// Enabling FDS to discover the lambda functions
			DisableFds: false,
		})

		// Look in ~/.aws/credentials for local AWS credentials
		// see the gloo OSS e2e test README for more information about configuring credentials for AWS e2e tests
		// https://github.com/solo-io/gloo/blob/main/test/e2e/README.md
		localAwsCredentials := credentials.NewSharedCredentials("", "")
		v, err := localAwsCredentials.Get()
		Expect(err).NotTo(HaveOccurred(), "No AWS credentials available")

		secret := &gloov1.Secret{
			Metadata: &core.Metadata{
				Namespace: "default",
				Name:      region,
			},
			Kind: &gloov1.Secret_Aws{
				Aws: &gloov1.AwsSecret{
					AccessKey: v.AccessKeyID,
					SecretKey: v.SecretAccessKey,
				},
			},
		}

		// Some users may require a session token as part of their credentials, so we set the session token (if any) when running locally.
		if os.Getenv("GCLOUD_BUILD_ID") == "" {
			secret.Kind.(*gloov1.Secret_Aws).Aws.SessionToken = v.SessionToken
		}

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

		// appending because setting it to just the upstream breaks since the test_context expects the default test upstream to exist
		testContext.ResourcesToCreate().Upstreams = append(testContext.ResourcesToCreate().Upstreams, upstream)
		testContext.ResourcesToCreate().Secrets = v1.SecretList{
			secret,
		}
	})

	JustBeforeEach(func() {
		testContext.JustBeforeEach()
	})

	AfterEach(func() {
		testContext.AfterEach()
	})

	JustAfterEach(func() {
		testContext.JustAfterEach()
	})

	lambdaDestination := func(functionName string, unwrapAsApiGateway, wrapAsApiGateway bool) *gloov1.Destination {
		return &gloov1.Destination{
			DestinationType: &gloov1.Destination_Upstream{
				Upstream: upstream.Metadata.Ref(),
			},
			DestinationSpec: &gloov1.DestinationSpec{
				DestinationType: &gloov1.DestinationSpec_Aws{
					Aws: &aws.DestinationSpec{
						LogicalName:        functionName,
						UnwrapAsApiGateway: unwrapAsApiGateway,
						WrapAsApiGateway:   wrapAsApiGateway,
					},
				},
			},
		}
	}

	Context("Enterprise Lambda plugin", func() {
		JustBeforeEach(func() {
			// Ensure that the upstream has been accepted and the lambda functions have been discovered
			helpers.EventuallyResourceAccepted(func() (resources.InputResource, error) {
				us, err := testContext.TestClients().UpstreamClient.Read(
					upstream.GetMetadata().GetNamespace(),
					upstream.GetMetadata().GetName(),
					clients.ReadOpts{})
				if err != nil {
					return nil, err
				}
				if len(us.GetAws().GetLambdaFunctions()) == 0 {
					return nil, errors.New("no lambda functions discovered")
				}

				return us, nil
			})
		})

		It("Can configure simple Lambda upstream", func() {
			// Update the virtual service to route to the lambda upstream "uppercase" function
			testContext.PatchDefaultVirtualService(func(vs *v1gw.VirtualService) *v1gw.VirtualService {
				vsBuilder := helpers.BuilderFromVirtualService(vs)
				vsBuilder.WithRouteActionToSingleDestination(e2e.DefaultRouteName, lambdaDestination("uppercase", false, false))
				return vsBuilder.Build()
			})
			httpRequestBuilder := testContext.GetHttpRequestBuilder().WithPostBody(`"solo.io"`)
			expectedResponse := &matchers.HttpResponse{
				StatusCode: http.StatusOK,
				Body: WithTransform(transforms.WithJsonBody(),
					And(
						HaveKeyWithValue("body", "SOLO.IO"),
						HaveKeyWithValue("statusCode", float64(http.StatusOK)),
					),
				),
			}
			Eventually(func(g Gomega) {
				response, err := testutils.DefaultHttpClient.Do(httpRequestBuilder.Build())
				Expect(err).NotTo(HaveOccurred())
				g.Expect(response).To(matchers.HaveHttpResponse(expectedResponse))
			}, "10s", "1s").Should(Succeed())
		})

		Context("Enterprise-specific functionality", func() {
			It("Can configure unwrapAsApiGateway", func() {
				testContext.PatchDefaultVirtualService(func(vs *v1gw.VirtualService) *v1gw.VirtualService {
					vsBuilder := helpers.BuilderFromVirtualService(vs)
					vsBuilder.WithRouteActionToSingleDestination(e2e.DefaultRouteName, lambdaDestination("echo", true, false))
					return vsBuilder.Build()
				})

				bodyString := "test"
				httpRequestBuilder := testContext.GetHttpRequestBuilder().
					WithPath("1?param_a=value_1&param_b=value_b").WithHeader("Foo", "bar").
					WithPostBody(fmt.Sprintf(`{"statusCode": %v, "body": "%v"}`, http.StatusOK, bodyString))
				expectedResponse := &matchers.HttpResponse{
					StatusCode: http.StatusOK,
					Body: WithTransform(transforms.WithJsonBody(),
						And(
							HaveKeyWithValue("body", bodyString),
							HaveKeyWithValue("statusCode", float64(http.StatusOK)),
							Not(HaveKey("headers")),
							Not(HaveKey("queryStringParameters")),
						),
					),
				}
				Eventually(func(g Gomega) {
					g.Expect(testutils.DefaultHttpClient.Do(httpRequestBuilder.Build())).Should(matchers.HaveHttpResponse(expectedResponse))
				}, "10s", "1s").Should(Succeed())
			})

			Context("Can configure wrapAsApiGateway", func() {
				JustBeforeEach(func() {
					testContext.PatchDefaultVirtualService(func(vs *v1gw.VirtualService) *v1gw.VirtualService {
						vsBuilder := helpers.BuilderFromVirtualService(vs)
						vsBuilder.WithRouteActionToSingleDestination(e2e.DefaultRouteName, lambdaDestination("echo", false, true))
						return vsBuilder.Build()
					})
				})

				It("works with a simple request", func() {
					body := "test"
					headers := map[string][]string{
						"single-value-header": {"value"},
						"multi-value-header":  {"value1", "value2"},
					}
					httpRequestBuilder := testContext.GetHttpRequestBuilder().
						WithPath("1?param_a=value_1&param_b=value_b&param_b=value_2").WithBody(body).
						WithHeader("single-value-header", "value").WithHeader("multi-value-header", "value1,value2")
					expectedResponse := &matchers.HttpResponse{
						StatusCode: http.StatusOK,
						Body: WithTransform(transforms.WithJsonBody(),
							And(
								HaveKeyWithValue("body", body),
								HaveKeyWithValue("httpMethod", "GET"),
								HaveKeyWithValue("path", "/1"),
								HaveKeyWithValue("headers", HaveKeyWithValue("single-value-header", headers["single-value-header"][0])),
								HaveKeyWithValue("multiValueHeaders", HaveKeyWithValue("multi-value-header", HaveExactElements(headers["multi-value-header"][0], headers["multi-value-header"][1]))),
								HaveKeyWithValue("multiValueQueryStringParameters", HaveKeyWithValue("param_b", HaveExactElements("value_b", "value_2"))),
								// note that we expect the value of param_b to be the last value assigned to it in the query string.
								HaveKeyWithValue("queryStringParameters", SatisfyAll(HaveKeyWithValue("param_a", "value_1"), HaveKeyWithValue("param_b", "value_2"))),
							),
						),
					}
					Eventually(func(g Gomega) {
						response, err := testutils.DefaultHttpClient.Do(httpRequestBuilder.Build())
						g.Expect(err).NotTo(HaveOccurred())
						g.Expect(response).To(matchers.HaveHttpResponse(expectedResponse))
					}, "10s", "1s").Should(Succeed())
				})

				It("works with a request with no body", func() {
					// expect that the response (which is identical to the payload sent to the lambda) is transformed appropriately.
					httpRequestBuilder := testContext.GetHttpRequestBuilder().
						WithPath("1?param_a=value_1&param_b=value_b&param_b=value_2")
					expectedResponse := &matchers.HttpResponse{
						StatusCode: http.StatusOK,
						Body: WithTransform(transforms.WithJsonBody(),
							And(
								HaveKeyWithValue("body", ""),
								HaveKeyWithValue("httpMethod", "GET"),
								HaveKeyWithValue("path", "/1"),
								HaveKeyWithValue("multiValueQueryStringParameters", HaveKeyWithValue("param_b", HaveExactElements("value_b", "value_2"))),
								// note that we expect the value of param_b to be the last value assigned to it in the query string.
								HaveKeyWithValue("queryStringParameters", SatisfyAll(HaveKeyWithValue("param_a", "value_1"), HaveKeyWithValue("param_b", "value_2"))),
							),
						),
					}
					Eventually(func(g Gomega) {
						response, err := testutils.DefaultHttpClient.Do(httpRequestBuilder.Build())
						g.Expect(err).NotTo(HaveOccurred())
						g.Expect(response).To(matchers.HaveHttpResponse(expectedResponse))
					}, "5s", "1s").Should(Succeed())
				})
			})
		})
	})
})
