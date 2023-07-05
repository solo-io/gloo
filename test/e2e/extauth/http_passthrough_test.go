package extauth_test

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	gatewayv1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	"github.com/solo-io/gloo/test/helpers"

	"github.com/solo-io/ext-auth-service/pkg/controller/translation"

	"github.com/golang/protobuf/ptypes/duration"
	structpb "github.com/golang/protobuf/ptypes/struct"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	extauth "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/extauth/v1"
	"github.com/solo-io/gloo/test/gomega/matchers"
	"github.com/solo-io/gloo/test/gomega/transforms"
	"github.com/solo-io/gloo/test/testutils"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	skmatchers "github.com/solo-io/solo-kit/test/matchers"
	"github.com/solo-io/solo-projects/test/e2e"
	"github.com/solo-io/solo-projects/test/v1helpers"
)

/*
	TODO Move chaining tests into their own file(s).
*/

var _ = Describe("HTTP Passthrough", func() {

	const (
		defaultUpstreamRequestTimeout = 5 * time.Second
	)

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

	expectStatusCodeWithHeaders := func(responseCode int, reqHeadersToAdd map[string][]string, responseHeadersToExpect map[string]interface{}) {
		httpReqBuilder := testContext.GetHttpRequestBuilder()
		req := httpReqBuilder.Build()
		req.Header = reqHeadersToAdd

		EventuallyWithOffset(1, func(g Gomega) {
			// setting the expected matcher inside the Eventually to avoid using the same matcher multiple times
			expectedResponse := &matchers.HttpResponse{
				StatusCode: responseCode,
				Headers:    responseHeadersToExpect,
			}

			resp, err := testutils.DefaultHttpClient.Do(req)
			g.Expect(err).NotTo(HaveOccurred())
			g.Expect(resp).To(matchers.HaveHttpResponse(expectedResponse))
		}, "5s", "0.5s").Should(Succeed())
	}

	Context("http passthrough sanity", func() {
		var (
			httpPassthroughConfig *extauth.PassThroughHttp
			httpAuthServer        *v1helpers.TestUpstream
			authConfig            *extauth.AuthConfig
			authServerUpstream    *gloov1.Upstream
			authconfigCfg         *structpb.Struct
			authConfigRequestPath string
		)

		BeforeEach(func() {
			httpPassthroughConfig = &extauth.PassThroughHttp{
				Request:  &extauth.PassThroughHttp_Request{},
				Response: &extauth.PassThroughHttp_Response{},
				ConnectionTimeout: &duration.Duration{
					Seconds: 10,
				},
			}
			authconfigCfg = nil
			authConfigRequestPath = ""
		})

		// This should be run AFTER the parent's testContext.JustBeforeEach()
		setupServerWithFailureModeAllow := func(failureModeAllow bool, handler v1helpers.ExtraHandlerFunc) {
			// Create the auth server and upstream with handler functions
			httpAuthServer = v1helpers.NewTestHttpUpstreamWithHandler(testContext.Ctx(), testContext.EnvoyInstance().LocalAddr(), handler)
			authServerUpstream = httpAuthServer.Upstream

			httpPassthroughConfig.Url = fmt.Sprintf("http://%s%s", httpAuthServer.Address, authConfigRequestPath)
			authConfig = &extauth.AuthConfig{
				Metadata: &core.Metadata{
					Name:      GetPassThroughExtAuthExtension().GetConfigRef().Name,
					Namespace: GetPassThroughExtAuthExtension().GetConfigRef().Namespace,
				},
				Configs: []*extauth.AuthConfig_Config{{
					AuthConfig: &extauth.AuthConfig_Config_PassThroughAuth{
						PassThroughAuth: &extauth.PassThroughAuth{
							Protocol: &extauth.PassThroughAuth_Http{
								Http: httpPassthroughConfig,
							},
							Config:           authconfigCfg,
							FailureModeAllow: failureModeAllow,
						},
					},
				}},
			}

			// These resources are dynamically created and modified in the tests, so we need to manually write them AFTER the JustBeforeEach.
			_, err := testContext.TestClients().AuthConfigClient.Write(authConfig, clients.WriteOpts{})
			Expect(err).NotTo(HaveOccurred())
			_, err = testContext.TestClients().UpstreamClient.Write(authServerUpstream, clients.WriteOpts{})
			Expect(err).NotTo(HaveOccurred())

			testContext.EventuallyProxyAccepted()
		}

		It("works", func() {
			setupServerWithFailureModeAllow(false, nil)
			expectStatusCodeWithHeaders(http.StatusOK, nil, nil)

			select {
			case received := <-httpAuthServer.C:
				Expect(received.Method).To(Equal(http.MethodPost))
			case <-time.After(defaultUpstreamRequestTimeout):
				Fail("request didn't make it upstream")
			}
		})

		// The failure_mode_allow feature allows for a user to be authorized even if the auth service is unavailable, BUT NOT when they are unauthorized.
		DescribeTable("failure_mode_allow sanity tests", func(failureModeAllow bool, clientResponse, expectedServiceResponse int) {
			handler := func(rw http.ResponseWriter, r *http.Request) bool {
				rw.WriteHeader(clientResponse)
				return true
			}
			setupServerWithFailureModeAllow(failureModeAllow, handler)
			// Any non-200 response from the client results in the passthrough service to return a 401
			expectStatusCodeWithHeaders(expectedServiceResponse, nil, nil)

			select {
			case received := <-httpAuthServer.C:
				Expect(received.Method).To(Equal(http.MethodPost))
			case <-time.After(defaultUpstreamRequestTimeout):
				Fail("request didn't make it upstream")
			}
		},
			Entry("unauthorized when failure_mode_allow=false and the service is unavailable", false, http.StatusServiceUnavailable, http.StatusUnauthorized),
			Entry("authorized when failure_mode_allow=true, but service is unavailable", true, http.StatusServiceUnavailable, http.StatusOK),
			Entry("unauthorized when failure_mode_allow=false and the user is unauthorized", false, http.StatusUnauthorized, http.StatusUnauthorized),
			Entry("unauthorized when failure_mode_allow=true, but the user is unauthorized", true, http.StatusUnauthorized, http.StatusUnauthorized),
		)

		Context("setting path on URL", func() {
			// We are running the server setup in the JustBeforeEach because we are writing configs, onto the clients which are created in the testContext.JustBeforeEach.
			JustBeforeEach(func() {
				authConfigRequestPath += "/auth"
				setupServerWithFailureModeAllow(false, nil)
			})

			It("has correct path to auth server", func() {
				expectStatusCodeWithHeaders(http.StatusOK, nil, nil)

				select {
				case received := <-httpAuthServer.C:
					Expect(received.Method).To(Equal(http.MethodPost))
					Expect(received.URL.Path).To(Equal("/auth"))
				case <-time.After(defaultUpstreamRequestTimeout):
					Fail("request didn't make it upstream")
				}
			})
		})

		Context("request", func() {
			JustBeforeEach(func() {
				httpPassthroughConfig.Request = &extauth.PassThroughHttp_Request{
					AllowedHeaders: []string{"x-passthrough-1", "x-passthrough-2"},
					HeadersToAdd: map[string]string{
						"x-added-header-1": "net new header",
					},
				}

				setupServerWithFailureModeAllow(false, nil)
			})

			It("copies `allowed_headers` request headers and adds `headers_to_add` headers to auth request", func() {
				expectStatusCodeWithHeaders(http.StatusOK, map[string][]string{
					"x-passthrough-1":    {"some header from request"},
					"x-passthrough-2":    {"some header from request 2"},
					"x-dont-passthrough": {"some header from request that shouldn't be passed through to auth server"},
				}, nil)

				select {
				case received := <-httpAuthServer.C:
					Expect(received.Method).To(Equal("POST"))
					Expect(received.Headers).To(And(
						HaveKeyWithValue("X-Passthrough-1", ContainElements("some header from request")),
						HaveKeyWithValue("X-Passthrough-2", ContainElements("some header from request 2")),
						HaveKeyWithValue("X-Added-Header-1", ContainElements("net new header")),
						Not(HaveKey("X-Dont-Passthrough")),
					))
				case <-time.After(defaultUpstreamRequestTimeout):
					Fail("request didn't make it upstream")
				}
			})
		})

		Context("response", func() {
			var (
				handler v1helpers.ExtraHandlerFunc
			)

			JustBeforeEach(func() {
				httpPassthroughConfig.Response = &extauth.PassThroughHttp_Response{
					AllowedUpstreamHeaders:       []string{"x-auth-header-1", "x-auth-header-2"},
					AllowedClientHeadersOnDenied: []string{"x-auth-header-1"},
				}

				setupServerWithFailureModeAllow(false, handler)
			})

			Context("On authorized response", func() {
				BeforeEach(func() {
					handler = func(rw http.ResponseWriter, r *http.Request) bool {
						rw.Header().Set("x-auth-header-1", "some value")
						rw.Header().Set("x-auth-header-2", "some value 2")
						rw.Header().Set("x-shouldnt-upstream", "shouldn't upstream")
						return true
					}
				})

				It("copies `allowed_headers` request headers and adds `headers_to_add` headers to auth request", func() {
					expectStatusCodeWithHeaders(http.StatusOK, map[string][]string{
						"x-auth-header-1": {"hello"},
						"x-passthrough-1": {"some header from request that should go to upstream"},
					}, nil)

					select {
					case received := <-testContext.TestUpstream().C:
						Expect(received.Method).To(Equal(http.MethodGet))
						Expect(received.Headers).To(And(
							// This header should have an appended value since it exists on the original request
							HaveKeyWithValue("X-Auth-Header-1", ContainElements("hello,some value")),
							HaveKeyWithValue("X-Auth-Header-2", ContainElements("some value 2")),
							Not(HaveKey("X-Shouldnt-Upstream")),
						))

					case <-time.After(defaultUpstreamRequestTimeout):
						Fail("request didn't make it upstream")
					}
				})
			})

			Context("on authorized response", func() {
				BeforeEach(func() {
					handler = func(rw http.ResponseWriter, r *http.Request) bool {
						rw.Header().Set("x-auth-header-1", "some value")
						rw.WriteHeader(http.StatusUnauthorized)
						return true
					}
				})

				It("sends allowed authorization headers back to downstream", func() {
					expectStatusCodeWithHeaders(http.StatusUnauthorized, nil, map[string]interface{}{"x-auth-header-1": "some value"})
				})
			})
		})

		Context("request to Auth Server Body", func() {
			BeforeEach(func() {
				// We need these settings so envoy buffers the request body and sends it to the ext-auth-service
				testContext.UpdateRunSettings(func(settings *gloov1.Settings) {
					settings.Extauth.RequestBody = &extauth.BufferSettings{
						MaxRequestBytes:     uint32(1024),
						AllowPartialMessage: true,
					}
				})
			})
			JustBeforeEach(func() {
				setupServerWithFailureModeAllow(false, nil)
			})

			expectStatusCodeWithBody := func(responseCode int, body string) {
				httpReqBuilder := testContext.GetHttpRequestBuilder().WithContentType("text/plain").WithPostBody(body)
				EventuallyWithOffset(1, func(g Gomega) *http.Response {
					resp, err := testutils.DefaultHttpClient.Do(httpReqBuilder.Build())
					g.Expect(err).NotTo(HaveOccurred())
					return resp
				}, "15s", "0.5s").Should(HaveHTTPStatus(responseCode))
			}

			Context("with PassThroughBody=true", func() {
				BeforeEach(func() {
					httpPassthroughConfig.Request.PassThroughBody = true
				})

				It("passes the request body correctly", func() {
					expectStatusCodeWithBody(http.StatusOK, "some body")

					select {
					case received := <-httpAuthServer.C:
						// The received body is in json format, so transforming it into a map to assert a key/value on it.
						Expect(received.Body).To(WithTransform(transforms.WithJsonBody(), HaveKeyWithValue("body", "some body")))
					case <-time.After(defaultUpstreamRequestTimeout):
						Fail("request didn't make it upstream")
					}
				})
			})

			Context("with PassThroughBody=false", func() {
				BeforeEach(func() {
					httpPassthroughConfig.Request.PassThroughBody = false
				})

				It("body is empty if no body, config, state, or filtermetadata passthrough is set", func() {
					expectStatusCodeWithBody(http.StatusOK, "some body")

					select {
					case received := <-httpAuthServer.C:
						Expect(string(received.Body)).To(BeEmpty())
					case <-time.After(defaultUpstreamRequestTimeout):
						Fail("request didn't make it upstream")
					}
				})
			})
		})

		Context("pass config specified on auth config in auth request body", func() {
			JustBeforeEach(func() {
				authconfigCfg = &structpb.Struct{
					Fields: map[string]*structpb.Value{
						"nestedStruct": {
							Kind: &structpb.Value_StructValue{
								StructValue: &structpb.Struct{
									Fields: map[string]*structpb.Value{
										"list": {
											Kind: &structpb.Value_ListValue{
												ListValue: &structpb.ListValue{
													Values: []*structpb.Value{
														{
															Kind: &structpb.Value_StringValue{
																StringValue: "some string",
															},
														},
														{
															Kind: &structpb.Value_NumberValue{
																NumberValue: float64(23),
															},
														},
													},
												},
											},
										},
										"string": {
											Kind: &structpb.Value_StringValue{
												StringValue: "some string",
											},
										},
										"int": {
											Kind: &structpb.Value_NumberValue{
												NumberValue: float64(23),
											},
										},
										"bool": {
											Kind: &structpb.Value_BoolValue{
												BoolValue: true,
											},
										},
									},
								},
							},
						},
					},
				}
				setupServerWithFailureModeAllow(false, nil)
			})
			It("passes through config from auth config to passthrough server", func() {
				type ConfigStruct struct {
					Config *structpb.Struct `json:"config"`
				}
				expectStatusCodeWithHeaders(http.StatusOK, nil, nil)

				select {
				case received := <-httpAuthServer.C:
					cfgStruct := &ConfigStruct{}
					err := json.Unmarshal(received.Body, cfgStruct)
					Expect(err).NotTo(HaveOccurred())
					Expect(cfgStruct.Config).To(skmatchers.MatchProto(authconfigCfg))
				case <-time.After(defaultUpstreamRequestTimeout):
					Fail("request didn't make it upstream")
				}
			})
		})

	})

	Context("http passthrough chaining sanity", func() {
		var (
			authConfig *extauth.AuthConfig
			httpAuthServerA,
			httpAuthServerB *v1helpers.TestUpstream
			httpPassthroughConfigA,
			httpPassthroughConfigB *extauth.PassThroughHttp
		)

		authConfigWithFailureModeAllow := func(httpPassthrough *extauth.PassThroughHttp, failureModeAllow bool) *extauth.AuthConfig_Config {
			return &extauth.AuthConfig_Config{
				AuthConfig: &extauth.AuthConfig_Config_PassThroughAuth{
					PassThroughAuth: &extauth.PassThroughAuth{
						Protocol: &extauth.PassThroughAuth_Http{
							Http: httpPassthrough,
						},
						FailureModeAllow: failureModeAllow,
					},
				},
			}
		}

		setupAuthConfigs := func(failureModeAllowA, failureModeAllowB bool) {
			httpPassthroughConfigA = &extauth.PassThroughHttp{
				Request:  &extauth.PassThroughHttp_Request{},
				Response: &extauth.PassThroughHttp_Response{},
				ConnectionTimeout: &duration.Duration{
					Seconds: 10,
				},
			}
			httpPassthroughConfigB = &extauth.PassThroughHttp{
				Request:  &extauth.PassThroughHttp_Request{},
				Response: &extauth.PassThroughHttp_Response{},
				ConnectionTimeout: &duration.Duration{
					Seconds: 10,
				},
			}
			authConfig = &extauth.AuthConfig{
				Metadata: &core.Metadata{
					Name:      GetPassThroughExtAuthExtension().GetConfigRef().Name,
					Namespace: GetPassThroughExtAuthExtension().GetConfigRef().Namespace,
				},
				Configs: []*extauth.AuthConfig_Config{
					authConfigWithFailureModeAllow(httpPassthroughConfigA, failureModeAllowA),
					authConfigWithFailureModeAllow(httpPassthroughConfigB, failureModeAllowB),
				},
			}
		}

		setupServers := func(handlerA, handlerB v1helpers.ExtraHandlerFunc) {
			httpAuthServerA = v1helpers.NewTestHttpUpstreamWithHandler(testContext.Ctx(), "127.0.0.1", handlerA)
			httpPassthroughConfigA.Url = fmt.Sprintf("http://%s", httpAuthServerA.Address)

			httpAuthServerB = v1helpers.NewTestHttpUpstreamWithHandler(testContext.Ctx(), "127.0.0.1", handlerB)
			httpPassthroughConfigB.Url = fmt.Sprintf("http://%s", httpAuthServerB.Address)

			// These resources are dynamically created and modified in the tests, so we need to manually write them AFTER the JustBeforeEach.
			_, err := testContext.TestClients().AuthConfigClient.Write(authConfig, clients.WriteOpts{})
			Expect(err).NotTo(HaveOccurred())
			_, err = testContext.TestClients().UpstreamClient.Write(httpAuthServerA.Upstream, clients.WriteOpts{})
			Expect(err).NotTo(HaveOccurred())
			_, err = testContext.TestClients().UpstreamClient.Write(httpAuthServerB.Upstream, clients.WriteOpts{})
			Expect(err).NotTo(HaveOccurred())
		}

		_ = func(responseCode int, body string) {
			httpReqBuilder := testContext.GetHttpRequestBuilder().WithContentType("text/plain").WithPostBody(body)
			EventuallyWithOffset(1, func(g Gomega) *http.Response {
				resp, err := testutils.DefaultHttpClient.Do(httpReqBuilder.Build())
				g.Expect(err).NotTo(HaveOccurred())
				return resp
			}).Should(HaveHTTPStatus(responseCode))
		}

		It("works", func() {
			setupAuthConfigs(false, false)
			setupServers(nil, nil)
			expectStatusCodeWithHeaders(http.StatusOK, nil, nil)
		})

		Context("server A has server error ", func() {
			DescribeTable("it returns expected response, depending on failure mode allow", func(failureModeAllowA, failureModeAllowB bool, expectedStatus int) {
				setupAuthConfigs(failureModeAllowA, failureModeAllowB)
				// note: HTTP PassThrough services return either a 200 (accepted) or 401 (denied) status,
				// regardless of the specific status gotten from the internal request
				handlerA := func(rw http.ResponseWriter, r *http.Request) bool {
					rw.WriteHeader(http.StatusServiceUnavailable)
					return true
				}
				handlerB := func(rw http.ResponseWriter, r *http.Request) bool {
					rw.WriteHeader(http.StatusOK)
					return true
				}
				setupServers(handlerA, handlerB)

				expectStatusCodeWithHeaders(expectedStatus, nil, nil)
			},
				Entry("neither service has failure_mode_allow enabled", false, false, http.StatusUnauthorized),
				Entry("service B has failure_mode_allow enabled", false, true, http.StatusUnauthorized),
				Entry("service A has failure_mode_allow enabled", true, false, http.StatusOK),
				Entry("both services have failure_mode_allow enabled", true, true, http.StatusOK),
			)
		})

		Context("can modify state", func() {
			BeforeEach(func() {
				setupAuthConfigs(false, false)
				httpPassthroughConfigA.Request.PassThroughState = true
				httpPassthroughConfigA.Response.ReadStateFromResponse = true
				httpPassthroughConfigB.Request.PassThroughState = true
			})

			JustBeforeEach(func() {
				handlerA := func(rw http.ResponseWriter, r *http.Request) bool {
					rw.Write([]byte(`{"state":{"list": ["item1", "item2", 3, {"item4":""}], "string": "hello", "integer": 9, "nestedObject":{"key":"value"}}}`))
					return false
				}
				// Setting up servers JustBeforeEach, since it has to happen AFTER the testContext.JustBeforeEach() where test clients are created
				setupServers(handlerA, nil)
			})
			It("modifies state in authServerA and authServerB can see the new state", func() {
				expectStatusCodeWithHeaders(200, nil, nil)

				select {
				case received := <-httpAuthServerB.C:
					Expect(string(received.Body)).To(Equal(`{"state":{"integer":9,"list":["item1","item2",3,{"item4":""}],"nestedObject":{"key":"value"},"string":"hello"}}`))
				case <-time.After(time.Second * 5):
					Fail("request didn't make it upstream")
				}
			})
		})
	})

	Context("https", func() {
		var (
			httpAuthServer *v1helpers.TestUpstream
			rootCaBytes    []byte
		)

		BeforeEach(func() {
			httpPassthroughConfig := &extauth.PassThroughHttp{
				Request:  &extauth.PassThroughHttp_Request{},
				Response: &extauth.PassThroughHttp_Response{},
				ConnectionTimeout: &duration.Duration{
					Seconds: 10,
				},
			}
			rootCaBytes, httpAuthServer = v1helpers.NewTestHttpsUpstreamWithHandler(testContext.Ctx(), testContext.EnvoyInstance().LocalAddr(), nil)
			// set environment variable for ext auth server https passthrough
			err := os.Setenv(translation.HttpsPassthroughCaCert, base64.StdEncoding.EncodeToString(rootCaBytes))
			Expect(err).NotTo(HaveOccurred())
			httpPassthroughConfig.Url = fmt.Sprintf("https://%s", httpAuthServer.Address)
			authConfig := &extauth.AuthConfig{
				Metadata: &core.Metadata{
					Name:      GetPassThroughExtAuthExtension().GetConfigRef().Name,
					Namespace: GetPassThroughExtAuthExtension().GetConfigRef().Namespace,
				},
				Configs: []*extauth.AuthConfig_Config{{
					AuthConfig: &extauth.AuthConfig_Config_PassThroughAuth{
						PassThroughAuth: &extauth.PassThroughAuth{
							Protocol: &extauth.PassThroughAuth_Http{
								Http: httpPassthroughConfig,
							},
							Config: nil,
						},
					},
				}},
			}

			testContext.ResourcesToCreate().AuthConfigs = extauth.AuthConfigList{
				authConfig,
			}
			testContext.ResourcesToCreate().Upstreams = append(testContext.ResourcesToCreate().Upstreams, httpAuthServer.Upstream)
		})

		AfterEach(func() {
			err := os.Unsetenv(translation.HttpsPassthroughCaCert)
			Expect(err).NotTo(HaveOccurred())
		})

		It("works", func() {
			expectStatusCodeWithHeaders(http.StatusOK, nil, nil)

			select {
			case received := <-httpAuthServer.C:
				Expect(received.Method).To(Equal(http.MethodPost))
			case <-time.After(defaultUpstreamRequestTimeout):
				Fail("request didn't make it upstream")
			}
		})
	})
})

func GetPassThroughExtAuthExtension() *extauth.ExtAuthExtension {
	return &extauth.ExtAuthExtension{
		Spec: &extauth.ExtAuthExtension_ConfigRef{
			ConfigRef: &core.ResourceRef{
				Name:      "passthrough-auth",
				Namespace: e2e.WriteNamespace,
			},
		},
	}
}
