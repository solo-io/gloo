package extauth_test

import (
	"context"
	"net/http"
	"sync/atomic"

	envoy_service_auth_v3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	structpb "github.com/golang/protobuf/ptypes/struct"
	errors "github.com/rotisserie/eris"
	grpcPassthrough "github.com/solo-io/ext-auth-service/pkg/config/passthrough/grpc"
	passthrough_utils "github.com/solo-io/ext-auth-service/pkg/config/passthrough/test_utils"
	gatewayv1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	v3 "github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/config/trace/v3"
	extauth "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/extauth/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/hcm"
	gloov1static "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/static"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/tracing"
	"github.com/solo-io/gloo/test/ginkgo/parallel"
	"github.com/solo-io/gloo/test/helpers"
	"github.com/solo-io/gloo/test/testutils"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/solo-projects/test/e2e"
)

/*
	TODO Move chaining tests into their own file(s).
*/

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
			authServer     *passthrough_utils.GrpcAuthServer
			authServerPort uint32
		)
		BeforeEach(func() {
			authServerPort = atomic.AddUint32(&baseExtauthPort, 1) + uint32(parallel.GetPortOffset())
		})

		setupAuthServerWithPassthroughConfig := func(failureModeAllow bool) {
			// start auth server
			err := authServer.Start(int(authServerPort))
			ExpectWithOffset(1, err).NotTo(HaveOccurred())

			// write auth configuration
			ac := &extauth.AuthConfig{
				Metadata: &core.Metadata{
					Name:      GetPassThroughExtAuthExtension().GetConfigRef().Name,
					Namespace: GetPassThroughExtAuthExtension().GetConfigRef().Namespace,
				},
				Configs: []*extauth.AuthConfig_Config{{
					AuthConfig: &extauth.AuthConfig_Config_PassThroughAuth{
						PassThroughAuth: getPassThroughAuthConfig(authServer.GetAddress(), failureModeAllow),
					},
				}},
			}
			_, err = testContext.TestClients().AuthConfigClient.Write(ac, clients.WriteOpts{})
			ExpectWithOffset(1, err).NotTo(HaveOccurred())

			// get proxy with pass through auth extension
			testContext.PatchDefaultVirtualService(func(vs *gatewayv1.VirtualService) *gatewayv1.VirtualService {
				vsBuilder := helpers.BuilderFromVirtualService(testContext.ResourcesToCreate().VirtualServices[0])
				vsBuilder.WithVirtualHostOptions(&gloov1.VirtualHostOptions{
					Extauth: GetPassThroughExtAuthExtension(),
				})
				return vsBuilder.Build()
			})
		}

		Context("failure_mode_allow=false (default)", func() {
			const (
				zipkinPort = 9411
			)

			JustBeforeEach(func() {
				setupAuthServerWithPassthroughConfig(false)
			})

			AfterEach(func() {
				authServer.Stop()
			})

			Context("when auth server returns ok response", func() {

				BeforeEach(func() {
					authServerResponse := passthrough_utils.OkResponse()
					authServer = passthrough_utils.NewGrpcAuthServerWithResponse(authServerResponse, nil)
				})

				It("should accept extauth passthrough", func() {
					expectRequestEventuallyReturnsResponseCode(http.StatusOK)
				})

			})

			Context("when auth server returns denied response", func() {

				BeforeEach(func() {
					authServerResponse := passthrough_utils.DeniedResponse()
					authServer = passthrough_utils.NewGrpcAuthServerWithResponse(authServerResponse, nil)
				})

				It("should deny extauth passthrough", func() {
					expectRequestEventuallyReturnsResponseCode(http.StatusUnauthorized)
				})

			})

			Context("when auth server errors", func() {

				BeforeEach(func() {
					authServerError := errors.New("auth server internal server error")
					authServer = passthrough_utils.NewGrpcAuthServerWithResponse(nil, authServerError)
				})

				It("should deny extauth passthrough", func() {
					expectRequestEventuallyReturnsResponseCode(http.StatusForbidden)
				})

			})

			Context("when auth server returns ok response with valid dynamic metadata properties", func() {

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
					authServer = passthrough_utils.NewGrpcAuthServerWithResponse(authServerResponse, nil)
				})

				It("should accept extauth passthrough", func() {
					expectRequestEventuallyReturnsResponseCode(http.StatusOK)
				})

			})

			Context("with tracing enabled", func() {
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

					authServer = passthrough_utils.NewGrpcAuthServerWithTracingRequired()
				})

				Context("when auth server returns ok response", func() {
					It("should accept extauth passthrough", func() {
						expectRequestEventuallyReturnsResponseCode(http.StatusOK)
					})
				})
			})
		})

		Context("failure_mode_allow=true", func() {
			JustBeforeEach(func() {
				setupAuthServerWithPassthroughConfig(true)
			})

			AfterEach(func() {
				authServer.Stop()
			})

			Context("when auth server returns denied response", func() {
				BeforeEach(func() {
					authServerResponse := passthrough_utils.DeniedResponse()
					authServer = passthrough_utils.NewGrpcAuthServerWithResponse(authServerResponse, nil)
				})

				It("should deny extauth passthrough", func() {
					expectRequestEventuallyReturnsResponseCode(http.StatusUnauthorized)
				})
			})

			Context("when auth server returns a 5XX server error", func() {
				BeforeEach(func() {
					resp := passthrough_utils.ServerErrorResponse()
					authServer = passthrough_utils.NewGrpcAuthServerWithResponse(resp, nil)
				})

				It("should allow extauth passthrough", func() {
					expectRequestEventuallyReturnsResponseCode(http.StatusOK)
				})
			})

			Context("when auth server has errors", func() {
				BeforeEach(func() {
					authServerError := errors.New("lorem ipsum, this causes a Check to return an err")
					resp := passthrough_utils.ServerErrorResponse()
					authServer = passthrough_utils.NewGrpcAuthServerWithResponse(resp, authServerError)
				})

				It("should allow extauth passthrough", func() {
					expectRequestEventuallyReturnsResponseCode(http.StatusOK)
				})
			})
		})

	})

	// These tests are used to validate that custom config is passed properly to the passthrough service
	Context("passthrough auth config sanity", func() {
		var (
			authServer     *passthrough_utils.GrpcAuthServer
			authServerPort uint32
		)

		newGrpcAuthServerWithRequiredConfig := func() *passthrough_utils.GrpcAuthServer {
			return &passthrough_utils.GrpcAuthServer{
				AuthChecker: func(ctx context.Context, req *envoy_service_auth_v3.CheckRequest) (*envoy_service_auth_v3.CheckResponse, error) {
					// Check if config exists in the FilterMetadata under the MetadataConfigKey.
					if passThroughFilterMetadata, ok := req.GetAttributes().GetMetadataContext().GetFilterMetadata()[grpcPassthrough.MetadataConfigKey]; ok {
						passThroughFields := passThroughFilterMetadata.GetFields()
						if value, ok := passThroughFields["customConfig1"]; ok && value.GetBoolValue() == true {
							// Required key was in FilterMetadata, succeed request
							return passthrough_utils.OkResponse(), nil
						}
						// Required key was not in FilterMetadata, deny fail request
						return passthrough_utils.DeniedResponse(), nil
					}
					// No passthrough properties were sent in FilterMetadata, fail request
					return passthrough_utils.DeniedResponse(), nil
				},
			}
		}

		BeforeEach(func() {
			authServerPort = atomic.AddUint32(&baseExtauthPort, 1) + uint32(parallel.GetPortOffset())
			authServer = newGrpcAuthServerWithRequiredConfig()

			// start auth server
			err := authServer.Start(int(authServerPort))
			Expect(err).NotTo(HaveOccurred())

			authConfig := &extauth.AuthConfig{
				Metadata: &core.Metadata{
					Name:      GetPassThroughExtAuthExtension().GetConfigRef().Name,
					Namespace: GetPassThroughExtAuthExtension().GetConfigRef().Namespace,
				},
				Configs: []*extauth.AuthConfig_Config{{
					AuthConfig: &extauth.AuthConfig_Config_PassThroughAuth{
						PassThroughAuth: getPassThroughAuthWithCustomConfig(authServer.GetAddress(), false),
					},
				}},
			}
			vsBuilder := helpers.BuilderFromVirtualService(testContext.ResourcesToCreate().VirtualServices[0])
			vsBuilder.WithVirtualHostOptions(&gloov1.VirtualHostOptions{
				Extauth: GetPassThroughExtAuthExtension(),
			})

			testContext.ResourcesToCreate().AuthConfigs = extauth.AuthConfigList{
				authConfig,
			}
			testContext.ResourcesToCreate().VirtualServices = gatewayv1.VirtualServiceList{
				vsBuilder.Build(),
			}
		})

		AfterEach(func() {
			authServer.Stop()
		})

		Context("passes config block to passthrough auth service", func() {
			It("correctly", func() {
				expectRequestEventuallyReturnsResponseCode(http.StatusOK)
			})
		})
	})

	Context("passthrough chaining sanity", func() {
		var (
			authServerA     *passthrough_utils.GrpcAuthServer
			authServerAPort uint32

			authServerB     *passthrough_utils.GrpcAuthServer
			authServerBPort uint32
		)

		BeforeEach(func() {
			authServerAPort = atomic.AddUint32(&baseExtauthPort, 1) + uint32(parallel.GetPortOffset())
			authServerBPort = atomic.AddUint32(&baseExtauthPort, 1) + uint32(parallel.GetPortOffset())
		})

		authConfigWithFailureModeAllow := func(address string, failureModeAllow bool) *extauth.AuthConfig_Config {
			return &extauth.AuthConfig_Config{
				AuthConfig: &extauth.AuthConfig_Config_PassThroughAuth{
					PassThroughAuth: getPassThroughAuthConfig(address, failureModeAllow),
				},
			}
		}

		// This is called on test's runtime (AFTER `testContext.JustBeforeEach`) to allow modifying the failureModeAllow values,
		// so we need to use the test clients to manually write the desired configurations
		setupAuthServersWithFailureModeAllow := func(failureModeAllowA, failureModeAllowB bool) {
			// start auth servers
			err := authServerA.Start(int(authServerAPort))
			ExpectWithOffset(1, err).NotTo(HaveOccurred())

			err = authServerB.Start(int(authServerBPort))
			ExpectWithOffset(1, err).NotTo(HaveOccurred())

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
			_, err = testContext.TestClients().AuthConfigClient.Write(ac, clients.WriteOpts{})
			ExpectWithOffset(1, err).NotTo(HaveOccurred())

			testContext.PatchDefaultVirtualService(func(service *gatewayv1.VirtualService) *gatewayv1.VirtualService {
				vsBuilder := helpers.BuilderFromVirtualService(service)
				vsBuilder.WithVirtualHostOptions(&gloov1.VirtualHostOptions{
					Extauth: GetPassThroughExtAuthExtension(),
				})
				return vsBuilder.Build()
			})
		}

		AfterEach(func() {
			authServerA.Stop()
			authServerB.Stop()
		})

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
				authServerA = passthrough_utils.NewGrpcAuthServerWithResponse(authServerAResponse, nil)

				// Configure AuthServerB (second in chain) to expect those dynamic metadata keys
				authServerB = passthrough_utils.NewGrpcAuthServerWithRequiredMetadata([]string{
					"key",
					"non-string-value",
				})
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
				authServerAResponse := passthrough_utils.OkResponse()
				authServerA = passthrough_utils.NewGrpcAuthServerWithResponse(authServerAResponse, nil)

				// Configure AuthServerB (second in chain) to expect dynamic metadata keys
				authServerB = passthrough_utils.NewGrpcAuthServerWithRequiredMetadata([]string{
					"key",
					"non-string-value",
				})
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
				authServerAResponse := passthrough_utils.ServerErrorResponse()
				authServerA = passthrough_utils.NewGrpcAuthServerWithResponse(authServerAResponse, nil)

				// Configure AuthServerB (second in chain) to not need any metadata to pass
				authServerBResponse := passthrough_utils.DeniedResponse()
				authServerB = passthrough_utils.NewGrpcAuthServerWithResponse(authServerBResponse, nil)
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

func getPassThroughAuthConfig(address string, failureModeAllow bool) *extauth.PassThroughAuth {
	return &extauth.PassThroughAuth{
		Protocol: &extauth.PassThroughAuth_Grpc{
			Grpc: &extauth.PassThroughGrpc{
				Address: address,
				// use default connection timeout
			},
		},
		FailureModeAllow: failureModeAllow,
	}
}

// This provides PassThroughAuth AuthConfig with Custom Config
func getPassThroughAuthWithCustomConfig(address string, failureModeAllow bool) *extauth.PassThroughAuth {
	passThroughAuth := getPassThroughAuthConfig(address, failureModeAllow)
	passThroughAuth.Config = &structpb.Struct{
		Fields: map[string]*structpb.Value{
			"customConfig1": {
				Kind: &structpb.Value_BoolValue{
					BoolValue: true,
				},
			},
		},
	}
	return passThroughAuth
}
