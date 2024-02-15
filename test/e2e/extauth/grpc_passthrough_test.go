package extauth_test

import (
	"context"
	"net/http"
	"time"

	envoy_service_auth_v3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	grpcPassthrough "github.com/solo-io/ext-auth-service/pkg/config/passthrough/grpc"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/wrapperspb"

	errors "github.com/rotisserie/eris"
	v3 "github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/config/trace/v3"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/hcm"
	gloov1static "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/static"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/tracing"
	extauth_test_server "github.com/solo-io/solo-projects/test/services/extauth/servers"

	structpb "github.com/golang/protobuf/ptypes/struct"
	passthrough_utils "github.com/solo-io/ext-auth-service/pkg/config/passthrough/test_utils"
	gatewayv1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	extauth "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/extauth/v1"
	"github.com/solo-io/gloo/test/helpers"
	"github.com/solo-io/gloo/test/testutils"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/solo-projects/test/e2e"
)

var _ = Describe("GRPC Passthrough", func() {

	var (
		testContext *e2e.TestContextWithExtensions
	)

	BeforeEach(func() {
		testContext = testContextFactory.NewTestContextWithExtensions(e2e.TestContextExtensions{
			ExtAuth: true,
		})
		testContext.BeforeEach()

		vsBuilder := helpers.BuilderFromVirtualService(testContext.ResourcesToCreate().VirtualServices[0])
		vsBuilder.WithVirtualHostOptions(&gloov1.VirtualHostOptions{
			Extauth: GetPassThroughExtAuthExtension(),
		})
		testContext.ResourcesToCreate().VirtualServices = gatewayv1.VirtualServiceList{
			vsBuilder.Build(),
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

	expectRequestEventuallyReturnsResponseCode := func(responseCode int) {
		httpReqBuilder := testContext.GetHttpRequestBuilder()
		EventuallyWithOffset(1, func(g Gomega) *http.Response {
			resp, err := testutils.DefaultHttpClient.Do(httpReqBuilder.Build())
			g.Expect(err).NotTo(HaveOccurred())
			return resp
		}, "5s", "0.5s").Should(HaveHTTPStatus(responseCode))
	}

	Context("passthrough sanity", func() {
		var (
			authServer       *extauth_test_server.GrpcServer
			failureModeAllow bool
			authConfigCfg    *structpb.Struct
			retryPolicy      *extauth.RetryPolicy
		)

		BeforeEach(func() {
			failureModeAllow = false
			authConfigCfg = nil
			retryPolicy = nil
		})

		JustBeforeEach(func() {
			// start auth server
			// auth servers are configured in the BeforeEach per context, so it should be started JustBeforeEach
			authServer.Start(testContext.Ctx())

			// write auth configuration
			ac := &extauth.AuthConfig{
				Metadata: &core.Metadata{
					Name:      GetPassThroughExtAuthExtension().GetConfigRef().Name,
					Namespace: GetPassThroughExtAuthExtension().GetConfigRef().Namespace,
				},
				Configs: []*extauth.AuthConfig_Config{{
					AuthConfig: &extauth.AuthConfig_Config_PassThroughAuth{
						PassThroughAuth: getPassThroughAuthConfig(authServer.GetAddress(), failureModeAllow, authConfigCfg, retryPolicy),
					},
				}},
			}
			_, err := testContext.TestClients().AuthConfigClient.Write(ac, clients.WriteOpts{
				OverwriteExisting: true,
			})
			Expect(err).NotTo(HaveOccurred())
		})

		Context("passthrough sanity", func() {
			When("auth server returns OK response", func() {
				BeforeEach(func() {
					server := passthrough_utils.NewGrpcAuthServerWithResponse(passthrough_utils.OkResponse(), nil)
					authServer = extauth_test_server.NewGrpcServer(server)
				})

				It("should accept extauth passthrough", func() {
					expectRequestEventuallyReturnsResponseCode(http.StatusOK)
				})
			})

			Context("failure_mode_allow=false (default)", func() {
				BeforeEach(func() {
					failureModeAllow = false
				})

				When("auth server returns denied response", func() {
					BeforeEach(func() {
						server := passthrough_utils.NewGrpcAuthServerWithResponse(passthrough_utils.DeniedResponse(), nil)
						authServer = extauth_test_server.NewGrpcServer(server)
					})

					It("should deny extauth passthrough", func() {
						expectRequestEventuallyReturnsResponseCode(http.StatusUnauthorized)
					})

				})

				When("auth server errors", func() {
					BeforeEach(func() {
						server := passthrough_utils.NewGrpcAuthServerWithResponse(nil, errors.New("auth server internal server error"))
						authServer = extauth_test_server.NewGrpcServer(server)
					})

					It("should deny extauth passthrough", func() {
						expectRequestEventuallyReturnsResponseCode(http.StatusForbidden)
					})

				})

				When("auth server returns ok response with valid dynamic metadata properties", func() {
					BeforeEach(func() {
						authServerResponse := passthrough_utils.OkResponseWithDynamicMetadata(&structpb.Struct{
							Fields: map[string]*structpb.Value{
								"current-state-key": {
									Kind: &structpb.Value_StringValue{
										StringValue: "new-state-value",
									},
								},
								"new-state-key": {
									Kind: &structpb.Value_StringValue{
										StringValue: "new-state-value",
									},
								},
							},
						})
						server := passthrough_utils.NewGrpcAuthServerWithResponse(authServerResponse, nil)
						authServer = extauth_test_server.NewGrpcServer(server)
					})

					It("should accept extauth passthrough", func() {
						expectRequestEventuallyReturnsResponseCode(http.StatusOK)
					})

				})

				// The zipkin tests rely on a hardcoded port (9411), so we shouldn't run them in parallel
				Context("with tracing enabled", Serial, func() {
					const (
						zipkinPort = 9411
					)

					BeforeEach(func() {
						zipkinUs := &gloov1.Upstream{
							Metadata: &core.Metadata{
								Name:      "zipkin",
								Namespace: "default",
							},
							UpstreamType: &gloov1.Upstream_Static{
								Static: &gloov1static.UpstreamSpec{
									Hosts: []*gloov1static.Host{
										{
											Addr: testContext.EnvoyInstance().LocalAddr(),
											Port: zipkinPort,
										},
									},
								},
							},
						}
						// "patch" the default gateway before it is written to include zipkin tracing
						testContext.ResourcesToCreate().Gateways[0].GetHttpGateway().Options = &gloov1.HttpListenerOptions{
							HttpConnectionManagerSettings: &hcm.HttpConnectionManagerSettings{
								Tracing: &tracing.ListenerTracingSettings{
									ProviderConfig: &tracing.ListenerTracingSettings_ZipkinConfig{
										ZipkinConfig: &v3.ZipkinConfig{
											CollectorCluster: &v3.ZipkinConfig_CollectorUpstreamRef{
												CollectorUpstreamRef: zipkinUs.Metadata.Ref(),
											},
											CollectorEndpoint:        "/api/v2/spans",
											CollectorEndpointVersion: v3.ZipkinConfig_HTTP_JSON,
										},
									},
								},
							},
						}
						testContext.ResourcesToCreate().Upstreams = append(testContext.ResourcesToCreate().Upstreams, zipkinUs)

						server := passthrough_utils.NewGrpcAuthServerWithTracingRequired()
						authServer = extauth_test_server.NewGrpcServer(server)
					})

					When("auth server returns ok response", func() {
						It("should accept extauth passthrough", func() {
							expectRequestEventuallyReturnsResponseCode(http.StatusOK)
						})
					})
				})
			})

			Context("failure_mode_allow=true", func() {
				BeforeEach(func() {
					failureModeAllow = true
				})

				When("auth server returns denied response", func() {
					BeforeEach(func() {
						server := passthrough_utils.NewGrpcAuthServerWithResponse(passthrough_utils.DeniedResponse(), nil)
						authServer = extauth_test_server.NewGrpcServer(server)
					})

					It("should deny extauth passthrough", func() {
						expectRequestEventuallyReturnsResponseCode(http.StatusUnauthorized)
					})
				})

				Context("when auth server returns a 5XX server error", func() {
					BeforeEach(func() {
						server := passthrough_utils.NewGrpcAuthServerWithResponse(passthrough_utils.ServerErrorResponse(), nil)
						authServer = extauth_test_server.NewGrpcServer(server)
					})

					It("should allow extauth passthrough", func() {
						expectRequestEventuallyReturnsResponseCode(http.StatusOK)
					})
				})

				When("auth server has errors", func() {
					BeforeEach(func() {
						authServerError := errors.New("lorem ipsum, this causes a Check to return an err")
						server := passthrough_utils.NewGrpcAuthServerWithResponse(passthrough_utils.ServerErrorResponse(), authServerError)
						authServer = extauth_test_server.NewGrpcServer(server)
					})

					It("should allow extauth passthrough", func() {
						expectRequestEventuallyReturnsResponseCode(http.StatusOK)
					})
				})
			})

		})

		// These tests are used to validate that custom config is passed properly to the passthrough service
		Context("with custom auth config", func() {
			const configFieldKey = "customConfig1"

			BeforeEach(func() {
				// Sets up an auth server which requires a key in the passthrough filter metadata
				authHandler := func(ctx context.Context, req *envoy_service_auth_v3.CheckRequest) (*envoy_service_auth_v3.CheckResponse, error) {
					// Check if config exists in the FilterMetadata under the MetadataConfigKey.
					if passThroughFilterMetadata, ok := req.GetAttributes().GetMetadataContext().GetFilterMetadata()[grpcPassthrough.MetadataConfigKey]; ok {
						passThroughFields := passThroughFilterMetadata.GetFields()
						if value, ok := passThroughFields[configFieldKey]; ok && value.GetBoolValue() == true {
							// Required key was in FilterMetadata, succeed request
							return passthrough_utils.OkResponse(), nil
						}
						// Required key was not in FilterMetadata, deny fail request
						return passthrough_utils.DeniedResponse(), nil
					}
					// No passthrough properties were sent in FilterMetadata, fail request
					return passthrough_utils.DeniedResponse(), nil
				}
				auth := &passthrough_utils.GrpcAuthServer{
					AuthChecker: authHandler,
				}
				authServer = extauth_test_server.NewGrpcServer(auth)
			})

			authConfigCfgWithFieldValue := func(fieldValue bool) *structpb.Struct {
				return &structpb.Struct{
					Fields: map[string]*structpb.Value{
						configFieldKey: {
							Kind: &structpb.Value_BoolValue{
								BoolValue: fieldValue,
							},
						},
					},
				}
			}

			When("auth config is set with the expected fields", func() {
				BeforeEach(func() {
					authConfigCfg = authConfigCfgWithFieldValue(true)
				})

				It("returns a 200 OK response", func() {
					expectRequestEventuallyReturnsResponseCode(http.StatusOK)
				})
			})

			When("auth config is set with the wrong fields", func() {
				BeforeEach(func() {
					authConfigCfg = authConfigCfgWithFieldValue(false)
				})

				It("returns a 401 Unauthorized response", func() {
					expectRequestEventuallyReturnsResponseCode(http.StatusUnauthorized)
				})
			})

			When("auth config is missing expected fields", func() {
				BeforeEach(func() {
					authConfigCfg = nil
				})

				It("returns a 401 Unauthorized response", func() {
					expectRequestEventuallyReturnsResponseCode(http.StatusUnauthorized)
				})
			})

			When("auth config is set with retry", func() {
				simulateGlitchOnAuthServer := func() chan bool {
					authServer.Stop()
					glitch := make(chan bool)
					go func() {
						defer GinkgoRecover()
						// Wait until the client is ready to attempt an RPC to ensure the authServer is down when called
						<-glitch
						// Wait to ensure the first call fails
						time.Sleep(50 * time.Millisecond)
						// Restart the server on the same port
						server := passthrough_utils.NewGrpcAuthServerWithResponse(passthrough_utils.OkResponse(), nil)
						authServer = extauth_test_server.NewGrpcServerOnPort(server, authServer.GetPort())
						authServer.Start(testContext.Ctx())
						// Ensure we're done so the cleanup can begin
						glitch <- true
					}()
					return glitch
				}

				BeforeEach(func() {
					authConfigCfg = authConfigCfgWithFieldValue(true)

					retryPolicy = &extauth.RetryPolicy{
						NumRetries: &wrapperspb.UInt32Value{
							Value: 100,
						},
						Strategy: &extauth.RetryPolicy_RetryBackOff{
							RetryBackOff: &extauth.BackoffStrategy{
								BaseInterval: &durationpb.Duration{
									Nanos: 10000000, // 10ms
								},
								MaxInterval: &durationpb.Duration{
									Nanos: 10000000, // 10ms
								},
							},
						},
					}
				})

				AfterEach(func() {
					authServer.Stop()
				})

				It("retries and succeeds", func() {
					// Let it succeed the first time so we know the authconfig has been accepted
					expectRequestEventuallyReturnsResponseCode(http.StatusOK)

					glitch := simulateGlitchOnAuthServer()
					httpReqBuilder := testContext.GetHttpRequestBuilder()

					// Sync the glitch to ensure we attempt an RPC when the authServer is down
					glitch <- true
					resp, err := testutils.DefaultHttpClient.Do(httpReqBuilder.Build())
					Expect(err).NotTo(HaveOccurred())
					Expect(resp).To(HaveHTTPStatus(http.StatusOK))

					// Ensure we're done so the cleanup can begin
					<-glitch
				})
			})
		})
	})

	// TODO: move chaining tests to their own file
	Context("passthrough chaining sanity", func() {
		var (
			authServerA *extauth_test_server.GrpcServer
			authServerB *extauth_test_server.GrpcServer
		)

		authConfigWithFailureModeAllow := func(address string, failureModeAllow bool) *extauth.AuthConfig_Config {
			return &extauth.AuthConfig_Config{
				AuthConfig: &extauth.AuthConfig_Config_PassThroughAuth{
					PassThroughAuth: getPassThroughAuthConfig(address, failureModeAllow, nil, nil),
				},
			}
		}

		// This is called on test's runtime (AFTER `testContext.JustBeforeEach`) to allow modifying the failureModeAllow values,
		// so we need to use the test clients to manually write the desired configurations
		setupAuthServersWithFailureModeAllow := func(failureModeAllowA, failureModeAllowB bool) {
			// start auth servers
			authServerA.Start(testContext.Ctx())
			authServerB.Start(testContext.Ctx())

			// write auth configuration
			ac := &extauth.AuthConfig{
				Metadata: &core.Metadata{
					Name:      GetPassThroughExtAuthExtension().GetConfigRef().Name,
					Namespace: GetPassThroughExtAuthExtension().GetConfigRef().Namespace,
				},
				Configs: []*extauth.AuthConfig_Config{
					// ServerA is the initial chain, we only care about how it's final status affects Server B
					// This can apply to N+1 chains, where ServerA is the accumulated response from (0..N), and B is N+1
					authConfigWithFailureModeAllow(authServerA.GetAddress(), failureModeAllowA),
					authConfigWithFailureModeAllow(authServerB.GetAddress(), failureModeAllowB),
				},
			}
			_, err := testContext.TestClients().AuthConfigClient.Write(ac, clients.WriteOpts{})
			ExpectWithOffset(1, err).NotTo(HaveOccurred())

			testContext.PatchDefaultVirtualService(func(service *gatewayv1.VirtualService) *gatewayv1.VirtualService {
				vsBuilder := helpers.BuilderFromVirtualService(service)
				vsBuilder.WithVirtualHostOptions(&gloov1.VirtualHostOptions{
					Extauth: GetPassThroughExtAuthExtension(),
				})
				return vsBuilder.Build()
			})
		}

		Context("first auth server writes metadata, second requires it", func() {
			BeforeEach(func() {
				// Configure AuthServerA (first in chain) to return DynamicMetadata.
				authServerAResponse := passthrough_utils.OkResponseWithDynamicMetadata(&structpb.Struct{
					Fields: map[string]*structpb.Value{
						"key": {
							Kind: &structpb.Value_StringValue{
								StringValue: "value",
							},
						},
						"non-string-value": {
							Kind: &structpb.Value_StructValue{
								StructValue: &structpb.Struct{
									Fields: map[string]*structpb.Value{
										"nested-key": {
											Kind: &structpb.Value_StringValue{
												StringValue: "nested-value",
											},
										},
									},
								},
							},
						},
					},
				})
				serverA := passthrough_utils.NewGrpcAuthServerWithResponse(authServerAResponse, nil)
				authServerA = extauth_test_server.NewGrpcServer(serverA)

				// Configure AuthServerB (second in chain) to expect those dynamic metadata keys
				serverB := passthrough_utils.NewGrpcAuthServerWithRequiredMetadata([]string{
					"key",
					"non-string-value",
				})
				authServerB = extauth_test_server.NewGrpcServer(serverB)
			})

			It("should accept extauth passthrough", func() {
				// This will pass only if the following events occur:
				//		1. AuthServerA returns DynamicMetadata under PassThrough Key and that data is stored on AuthorizationRequest
				//		2. State on AuthorizationRequest is parsed and sent on subsequent request to AuthServerB
				//		3. AuthServerB receives the Metadata and returns ok if all keys are present.
				setupAuthServersWithFailureModeAllow(false, false)
				expectRequestEventuallyReturnsResponseCode(http.StatusOK)
			})
		})

		Context("first auth server does not write metadata, second requires it", func() {
			BeforeEach(func() {
				// Configure AuthServerA (first in chain) to NOT return DynamicMetadata.
				serverA := passthrough_utils.NewGrpcAuthServerWithResponse(passthrough_utils.OkResponse(), nil)
				authServerA = extauth_test_server.NewGrpcServer(serverA)

				// Configure AuthServerB (second in chain) to expect dynamic metadata keys
				serverB := passthrough_utils.NewGrpcAuthServerWithRequiredMetadata([]string{
					"key",
					"non-string-value",
				})
				authServerB = extauth_test_server.NewGrpcServer(serverB)
			})

			It("should deny extauth passthrough", func() {
				setupAuthServersWithFailureModeAllow(false, false)
				// This will deny the request because:
				//		1. AuthServerA does not return DynamicMetadata under PassThrough Key. So there is not AuthorizationRequest.State
				//		2. Since there is no AuthorizationRequest.State, no Metadata is sent in request to AuthServerB
				//		3. AuthServerB receives no Metadata, but requires certain fields and returns 401 since there are missing properties
				expectRequestEventuallyReturnsResponseCode(http.StatusUnauthorized)
			})
		})

		Context("first auth server has server issues", func() {
			BeforeEach(func() {
				serverA := passthrough_utils.NewGrpcAuthServerWithResponse(passthrough_utils.ServerErrorResponse(), nil)
				authServerA = extauth_test_server.NewGrpcServer(serverA)

				// Configure AuthServerB (second in chain) to not need any metadata to pass
				serverB := passthrough_utils.NewGrpcAuthServerWithResponse(passthrough_utils.DeniedResponse(), nil)
				authServerB = extauth_test_server.NewGrpcServer(serverB)
			})

			DescribeTable("should return expected response, depending on failure_mode_allow", func(failureModeAllowA, failureModeAllowB bool, expectedResponse int) {
				setupAuthServersWithFailureModeAllow(failureModeAllowA, failureModeAllowB)
				expectRequestEventuallyReturnsResponseCode(expectedResponse)
			},
				Entry("neither service has failure_mode_allow enabled", false, false, http.StatusServiceUnavailable),
				Entry("service B has failure_mode_allow enabled", false, true, http.StatusServiceUnavailable),
				Entry("service A has failure_mode_allow enabled", true, false, http.StatusUnauthorized),
				Entry("both services have failure_mode_allow enabled", true, true, http.StatusUnauthorized),
			)
		})
	})
})

func getPassThroughAuthConfig(address string, failureModeAllow bool, config *structpb.Struct, retryPolicy *extauth.RetryPolicy) *extauth.PassThroughAuth {
	return &extauth.PassThroughAuth{
		Protocol: &extauth.PassThroughAuth_Grpc{
			Grpc: &extauth.PassThroughGrpc{
				Address:     address,
				RetryPolicy: retryPolicy,
				// use default connection timeout
			},
		},
		FailureModeAllow: failureModeAllow,
		Config:           config,
	}
}
