package e2e_test

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/config"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"

	"github.com/solo-io/ext-auth-service/pkg/config/oidc"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	extauthrunner "github.com/solo-io/solo-projects/projects/extauth/pkg/runner"
	"github.com/solo-io/solo-projects/test/services"

	extauth "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/extauth/v1"

	"github.com/dgrijalva/jwt-go"
	"github.com/fgrosse/zaptest"
	"github.com/solo-io/gloo/pkg/utils"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	gloov1static "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/static"
	"github.com/solo-io/gloo/projects/gloo/pkg/defaults"
	"github.com/solo-io/go-utils/contextutils"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/memory"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-projects/test/v1helpers"
)

var (
	baseExtauthPort = uint32(27000)
)

var _ = Describe("External auth", func() {

	var (
		ctx         context.Context
		cancel      context.CancelFunc
		testClients services.TestClients
		settings    extauthrunner.Settings
		cache       memory.InMemoryResourceCache
	)

	BeforeEach(func() {
		extAuthPort := atomic.AddUint32(&baseExtauthPort, 1) + uint32(config.GinkgoConfig.ParallelNode*1000)

		logger := zaptest.LoggerWriter(GinkgoWriter)
		contextutils.SetFallbackLogger(logger.Sugar())

		ctx, cancel = context.WithCancel(context.Background())
		cache = memory.NewInMemoryResourceCache()

		testClients = services.GetTestClients(cache)
		testClients.GlooPort = int(services.AllocateGlooPort())

		extauthAddr := "localhost"
		if runtime.GOOS == "darwin" {
			extauthAddr = "host.docker.internal"
		}

		extAuthServer := &gloov1.Upstream{
			Metadata: core.Metadata{
				Name:      "extauth-server",
				Namespace: "default",
			},
			UseHttp2: true,
			UpstreamType: &gloov1.Upstream_Static{
				Static: &gloov1static.UpstreamSpec{
					Hosts: []*gloov1static.Host{{
						Addr: extauthAddr,
						Port: extAuthPort,
					}},
				},
			},
		}

		_, err := testClients.UpstreamClient.Write(extAuthServer, clients.WriteOpts{})
		Expect(err).NotTo(HaveOccurred())

		ref := extAuthServer.Metadata.Ref()
		extauthSettings := &extauth.Settings{
			ExtauthzServerRef: &ref,
		}

		settings = extauthrunner.Settings{
			GlooAddress:  fmt.Sprintf("localhost:%d", testClients.GlooPort),
			DebugPort:    0,
			ServerPort:   int(extAuthPort),
			SigningKey:   "hello",
			UserIdHeader: "X-User-Id",
		}

		glooSettings := &gloov1.Settings{Extauth: extauthSettings}

		what := services.What{
			DisableGateway: true,
			DisableUds:     true,
			DisableFds:     true,
		}

		services.RunGlooGatewayUdsFdsOnPort(ctx, cache, int32(testClients.GlooPort), what, defaults.GlooSystem, nil, nil, glooSettings)
		go func(testCtx context.Context) {
			defer GinkgoRecover()
			err := extauthrunner.RunWithSettings(testCtx, settings)
			if testCtx.Err() == nil {
				Expect(err).NotTo(HaveOccurred())
			}
		}(ctx)
	})

	AfterEach(func() {
		cancel()
	})

	Context("With envoy", func() {

		var (
			envoyInstance *services.EnvoyInstance
			testUpstream  *v1helpers.TestUpstream
			envoyPort     = uint32(8080)
		)

		BeforeEach(func() {
			var err error
			envoyInstance, err = envoyFactory.NewEnvoyInstance()
			Expect(err).NotTo(HaveOccurred())

			err = envoyInstance.Run(testClients.GlooPort)
			Expect(err).NotTo(HaveOccurred())

			testUpstream = v1helpers.NewTestHttpUpstream(ctx, envoyInstance.LocalAddr())

			var opts clients.WriteOpts
			up := testUpstream.Upstream
			_, err = testClients.UpstreamClient.Write(up, opts)
			Expect(err).NotTo(HaveOccurred())

		})

		AfterEach(func() {
			if envoyInstance != nil {
				_ = envoyInstance.Clean()
			}
		})

		Context("using new config format", func() {

			Context("basic auth sanity tests", func() {

				BeforeEach(func() {

					// drain channel as we dont care about it
					go func() {
						for range testUpstream.C {
						}
					}()

					_, err := testClients.AuthConfigClient.Write(&extauth.AuthConfig{
						Metadata: core.Metadata{
							Name:      GetBasicAuthExtension().GetConfigRef().Name,
							Namespace: GetBasicAuthExtension().GetConfigRef().Namespace,
						},
						Configs: []*extauth.AuthConfig_Config{{
							AuthConfig: &extauth.AuthConfig_Config_BasicAuth{
								BasicAuth: getBasicAuthConfig(),
							},
						}},
					}, clients.WriteOpts{Ctx: ctx})
					Expect(err).NotTo(HaveOccurred())

					proxy := getProxyExtAuthBasicAuth(envoyPort, testUpstream.Upstream.Metadata.Ref())

					_, err = testClients.ProxyClient.Write(proxy, clients.WriteOpts{Ctx: ctx})
					Expect(err).NotTo(HaveOccurred())

					Eventually(func() (core.Status, error) {
						proxy, err := testClients.ProxyClient.Read(proxy.Metadata.Namespace, proxy.Metadata.Name, clients.ReadOpts{})
						if err != nil {
							return core.Status{}, err
						}

						return proxy.Status, nil
					}, "5s", "0.1s").Should(MatchFields(IgnoreExtras, Fields{
						"Reason": BeEmpty(),
						"State":  Equal(core.Status_Accepted),
					}))
				})

				It("should deny ext auth envoy", func() {
					Eventually(func() (int, error) {
						resp, err := http.Get(fmt.Sprintf("http://%s:%d/1", "localhost", envoyPort))
						if err != nil {
							return 0, err
						}
						return resp.StatusCode, nil
					}, "5s", "0.5s").Should(Equal(http.StatusUnauthorized))
				})

				It("should allow ext auth envoy", func() {
					Eventually(func() (int, error) {
						resp, err := http.Get(fmt.Sprintf("http://user:password@%s:%d/1", "localhost", envoyPort))
						if err != nil {
							return 0, err
						}
						return resp.StatusCode, nil
					}, "5s", "0.5s").Should(Equal(http.StatusOK))
				})

				It("should deny ext auth with wrong password", func() {
					Eventually(func() (int, error) {
						resp, err := http.Get(fmt.Sprintf("http://user:password2@%s:%d/1", "localhost", envoyPort))
						if err != nil {
							return 0, err
						}
						return resp.StatusCode, nil
					}, "5s", "0.5s").Should(Equal(http.StatusUnauthorized))
				})
			})

			Context("oidc sanity", func() {
				var (
					privateKey      *rsa.PrivateKey
					discoveryServer fakeDiscoveryServer
					secret          *gloov1.Secret
					proxy           *gloov1.Proxy
					token           string
				)
				BeforeEach(func() {
					discoveryServer = fakeDiscoveryServer{}

					privateKey = discoveryServer.Start()

					clientSecret := &extauth.OauthSecret{
						ClientSecret: "test",
					}

					secret = &gloov1.Secret{
						Metadata: core.Metadata{
							Name:      "secret",
							Namespace: "default",
						},
						Kind: &gloov1.Secret_Oauth{
							Oauth: clientSecret,
						},
					}
					_, err := testClients.SecretClient.Write(secret, clients.WriteOpts{})
					Expect(err).NotTo(HaveOccurred())

					_, err = testClients.AuthConfigClient.Write(&extauth.AuthConfig{
						Metadata: core.Metadata{
							Name:      getOidcExtAuthExtension().GetConfigRef().Name,
							Namespace: getOidcExtAuthExtension().GetConfigRef().Namespace,
						},
						Configs: []*extauth.AuthConfig_Config{{
							AuthConfig: &extauth.AuthConfig_Config_Oauth{
								Oauth: getOauthConfig(secret.Metadata.Ref()),
							},
						}},
					}, clients.WriteOpts{Ctx: ctx})
					Expect(err).NotTo(HaveOccurred())

					proxy = getProxyExtAuthOIDC(envoyPort, testUpstream.Upstream.Metadata.Ref())

					// create an id token
					tokenToSign := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
						"foo": "bar",
						"aud": "test-clientid",
						"sub": "user",
						"iss": "http://localhost:5556",
					})
					tokenToSign.Header["kid"] = "test-123"
					token, err = tokenToSign.SignedString(privateKey)
					Expect(err).NotTo(HaveOccurred())
				})

				JustBeforeEach(func() {
					_, err := testClients.ProxyClient.Write(proxy, clients.WriteOpts{})
					Expect(err).NotTo(HaveOccurred())

					Eventually(func() (core.Status, error) {
						proxy, err := testClients.ProxyClient.Read(proxy.Metadata.Namespace, proxy.Metadata.Name, clients.ReadOpts{})
						if err != nil {
							return core.Status{}, err
						}

						return proxy.Status, nil
					}, "5s", "0.1s").Should(MatchFields(IgnoreExtras, Fields{
						"Reason": BeEmpty(),
						"State":  Equal(core.Status_Accepted),
					}))
				})

				AfterEach(discoveryServer.Stop)

				Context("Oidc tests that don't forward to upstream", func() {
					BeforeEach(func() {
						// drain channel as we dont care about it
						go func() {
							for range testUpstream.C {
							}
						}()
					})

					It("should redirect to auth page", func() {
						Eventually(func() (string, error) {
							resp, err := http.Get(fmt.Sprintf("http://%s:%d/1", "localhost", envoyPort))
							if err != nil {
								return "", err
							}
							body, err := ioutil.ReadAll(resp.Body)
							if err != nil {
								return "", err
							}
							return string(body), nil
						}, "5s", "0.5s").Should(Equal("auth"))
					})

					It("should include email scope in url", func() {
						client := &http.Client{
							CheckRedirect: func(req *http.Request, via []*http.Request) error {
								return http.ErrUseLastResponse
							},
						}
						req, err := http.NewRequest("GET", fmt.Sprintf("http://%s:%d/1", "localhost", envoyPort), nil)
						Expect(err).NotTo(HaveOccurred())

						Eventually(func() (http.Response, error) {
							r, err := client.Do(req)
							if err != nil {
								return http.Response{}, err
							}
							return *r, err
						}, "5s", "0.5s").Should(MatchFields(IgnoreExtras, Fields{
							"StatusCode": Equal(http.StatusFound),
							"Header":     HaveKeyWithValue("Location", ContainElement(ContainSubstring("email"))),
						}))
					})

					It("should exchange token", func() {
						finalpage := fmt.Sprintf("http://%s:%d/success", "localhost", envoyPort)
						client := &http.Client{
							CheckRedirect: func(req *http.Request, via []*http.Request) error {
								return http.ErrUseLastResponse
							},
						}

						st := oidc.NewStateSigner([]byte(settings.SigningKey))
						signedState, err := st.Sign(finalpage)
						Expect(err).NotTo(HaveOccurred())
						req, err := http.NewRequest("GET", fmt.Sprintf("http://%s:%d/callback?code=1234&state="+string(signedState), "localhost", envoyPort), nil)
						Expect(err).NotTo(HaveOccurred())

						Eventually(func() (http.Response, error) {
							r, err := client.Do(req)
							if err != nil {
								return http.Response{}, err
							}
							return *r, err
						}, "5s", "0.5s").Should(MatchFields(IgnoreExtras, Fields{
							"StatusCode": Equal(http.StatusFound),
							"Header":     HaveKeyWithValue("Location", []string{finalpage}),
						}))
					})

					Context("oidc + opa sanity", func() {
						BeforeEach(func() {
							policy := &gloov1.Artifact{
								Metadata: core.Metadata{
									Name:      "jwt",
									Namespace: "default",
									Labels:    map[string]string{"team": "infrastructure"},
								},
								Data: map[string]string{
									"jwt.rego": `package test
	
				default allow = false
				allow {
					[header, payload, signature] = io.jwt.decode(input.state.jwt)
					payload["foo"] = "not-bar"
				}
				`}}
							modules := []*core.ResourceRef{{Name: policy.Metadata.Name}}

							_, err := testClients.AuthConfigClient.Write(&extauth.AuthConfig{
								Metadata: core.Metadata{
									Name:      getOidcAndOpaExtAuthExtension().GetConfigRef().Name,
									Namespace: getOidcAndOpaExtAuthExtension().GetConfigRef().Namespace,
								},
								Configs: []*extauth.AuthConfig_Config{
									{
										AuthConfig: &extauth.AuthConfig_Config_Oauth{
											Oauth: getOauthConfig(secret.Metadata.Ref()),
										},
									},
									{
										AuthConfig: &extauth.AuthConfig_Config_OpaAuth{
											OpaAuth: getOpaConfig(modules),
										},
									},
								},
							}, clients.WriteOpts{Ctx: ctx})
							Expect(err).NotTo(HaveOccurred())

							proxy = getProxyExtAuthOIDCAndOpa(envoyPort, secret.Metadata.Ref(), testUpstream.Upstream.Metadata.Ref(), modules)

							_, err = testClients.ArtifactClient.Write(policy, clients.WriteOpts{})
							Expect(err).ToNot(HaveOccurred())
						})

						It("should NOT allow access", func() {
							EventuallyWithOffset(1, func() (int, error) {
								req, err := http.NewRequest("GET", fmt.Sprintf("http://%s:%d/1", "localhost", envoyPort), nil)
								req.Header.Add("Authorization", "Bearer "+token)

								resp, err := http.DefaultClient.Do(req)
								if err != nil {
									return 0, err
								}
								return resp.StatusCode, nil
							}, "5s", "0.5s").Should(Equal(http.StatusForbidden))

						})

					})
				})

				ExpectUpstreamRequest := func() {
					EventuallyWithOffset(1, func() (int, error) {
						req, err := http.NewRequest("GET", fmt.Sprintf("http://%s:%d/1", "localhost", envoyPort), nil)
						req.Header.Add("Authorization", "Bearer "+token)

						resp, err := http.DefaultClient.Do(req)
						if err != nil {
							return 0, err
						}
						return resp.StatusCode, nil
					}, "5s", "0.5s").Should(Equal(http.StatusOK))

					select {
					case r := <-testUpstream.C:
						ExpectWithOffset(1, r.Headers["X-User-Id"]).To(HaveLen(1))
						ExpectWithOffset(1, r.Headers["X-User-Id"][0]).To(Equal("http://localhost:5556;user"))
					case <-time.After(time.Second):
						Fail("expected a message to be received")
					}
				}

				Context("Oidc tests that do forward to upstream", func() {
					It("should allow access with proper jwt token", func() {
						ExpectUpstreamRequest()
					})
				})

				Context("oidc + opa sanity", func() {
					BeforeEach(func() {
						policy := &gloov1.Artifact{
							Metadata: core.Metadata{
								Name:      "jwt",
								Namespace: "default",
								Labels:    map[string]string{"team": "infrastructure"},
							},
							Data: map[string]string{
								"jwt.rego": `package test

			default allow = false
			allow {
				[header, payload, signature] = io.jwt.decode(input.state.jwt)
				payload["foo"] = "bar"
			}
			`}}
						modules := []*core.ResourceRef{{Name: policy.Metadata.Name}}
						_, err := testClients.AuthConfigClient.Write(&extauth.AuthConfig{
							Metadata: core.Metadata{
								Name:      getOidcAndOpaExtAuthExtension().GetConfigRef().Name,
								Namespace: getOidcAndOpaExtAuthExtension().GetConfigRef().Namespace,
							},
							Configs: []*extauth.AuthConfig_Config{
								{
									AuthConfig: &extauth.AuthConfig_Config_Oauth{
										Oauth: getOauthConfig(secret.Metadata.Ref()),
									},
								},
								{
									AuthConfig: &extauth.AuthConfig_Config_OpaAuth{
										OpaAuth: getOpaConfig(modules),
									},
								},
							},
						}, clients.WriteOpts{Ctx: ctx})
						Expect(err).NotTo(HaveOccurred())
						proxy = getProxyExtAuthOIDCAndOpa(envoyPort, secret.Metadata.Ref(), testUpstream.Upstream.Metadata.Ref(), modules)

						_, err = testClients.ArtifactClient.Write(policy, clients.WriteOpts{})
						Expect(err).ToNot(HaveOccurred())
					})
					It("should allow access", func() {
						ExpectUpstreamRequest()
					})
				})

			})

			Context("api key sanity tests", func() {
				BeforeEach(func() {

					// drain channel as we dont care about it
					go func() {
						for range testUpstream.C {
						}
					}()

					_, err := testClients.AuthConfigClient.Write(&extauth.AuthConfig{
						Metadata: core.Metadata{
							Name:      getApiKeyExtAuthExtension().GetConfigRef().Name,
							Namespace: getApiKeyExtAuthExtension().GetConfigRef().Namespace,
						},
						Configs: []*extauth.AuthConfig_Config{{
							AuthConfig: &extauth.AuthConfig_Config_ApiKeyAuth{
								ApiKeyAuth: getApiKeyAuthConfig(),
							},
						}},
					}, clients.WriteOpts{Ctx: ctx})
					Expect(err).NotTo(HaveOccurred())

					apiKeySecret1 := &extauth.ApiKeySecret{
						ApiKey: "secretApiKey1",
					}

					secret1 := &gloov1.Secret{
						Metadata: core.Metadata{
							Name:      "secret1",
							Namespace: "default",
						},
						Kind: &gloov1.Secret_ApiKey{
							ApiKey: apiKeySecret1,
						},
					}

					apiKeySecret2 := &extauth.ApiKeySecret{
						ApiKey: "secretApiKey2",
					}

					secret2 := &gloov1.Secret{
						Metadata: core.Metadata{
							Name:      "secret2",
							Namespace: "default",
							Labels:    map[string]string{"team": "infrastructure"},
						},
						Kind: &gloov1.Secret_ApiKey{
							ApiKey: apiKeySecret2,
						},
					}

					_, err = testClients.SecretClient.Write(secret1, clients.WriteOpts{})
					Expect(err).ToNot(HaveOccurred())

					_, err = testClients.SecretClient.Write(secret2, clients.WriteOpts{})
					Expect(err).ToNot(HaveOccurred())

					proxy := getProxyExtAuthApiKeyAuth(envoyPort, testUpstream.Upstream.Metadata.Ref())

					_, err = testClients.ProxyClient.Write(proxy, clients.WriteOpts{})
					Expect(err).NotTo(HaveOccurred())

					Eventually(func() (core.Status, error) {
						proxy, err := testClients.ProxyClient.Read(proxy.Metadata.Namespace, proxy.Metadata.Name, clients.ReadOpts{})
						if err != nil {
							return core.Status{}, err
						}

						return proxy.Status, nil
					}, "5s", "0.1s").Should(MatchFields(IgnoreExtras, Fields{
						"Reason": BeEmpty(),
						"State":  Equal(core.Status_Accepted),
					}))
				})

				It("should deny ext auth envoy without apikey", func() {
					Eventually(func() (int, error) {
						resp, err := http.Get(fmt.Sprintf("http://%s:%d/1", "localhost", envoyPort))
						if err != nil {
							return 0, err
						}
						return resp.StatusCode, nil
					}, "5s", "0.5s").Should(Equal(http.StatusUnauthorized))
				})

				It("should deny ext auth envoy with incorrect apikey", func() {
					Eventually(func() (int, error) {
						req, err := http.NewRequest("GET", fmt.Sprintf("http://%s:%d/1", "localhost", envoyPort), nil)
						req.Header.Add("api-key", "badApiKey")
						resp, err := http.DefaultClient.Do(req)

						if err != nil {
							return 0, err
						}
						return resp.StatusCode, nil
					}, "5s", "0.5s").Should(Equal(http.StatusUnauthorized))
				})

				It("should accept ext auth envoy with correct apikey -- secret ref match", func() {
					Eventually(func() (int, error) {
						req, err := http.NewRequest("GET", fmt.Sprintf("http://%s:%d/1", "localhost", envoyPort), nil)
						req.Header.Add("api-key", "secretApiKey1")
						resp, err := http.DefaultClient.Do(req)

						if err != nil {
							return 0, err
						}
						return resp.StatusCode, nil
					}, "5s", "0.5s").Should(Equal(http.StatusOK))
				})

				It("should accept ext auth envoy with correct apikey -- label match", func() {
					Eventually(func() (int, error) {
						req, err := http.NewRequest("GET", fmt.Sprintf("http://%s:%d/1", "localhost", envoyPort), nil)
						req.Header.Add("api-key", "secretApiKey2")
						resp, err := http.DefaultClient.Do(req)

						if err != nil {
							return 0, err
						}
						return resp.StatusCode, nil
					}, "5s", "0.5s").Should(Equal(http.StatusOK))
				})
			})
		})

	})
})

var startDiscoveryServerOnce sync.Once
var cachedPrivateKey *rsa.PrivateKey

type fakeDiscoveryServer struct {
	s http.Server
}

func (f *fakeDiscoveryServer) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_ = f.s.Shutdown(ctx)
}

func (f *fakeDiscoveryServer) Start() *rsa.PrivateKey {
	startDiscoveryServerOnce.Do(func() {
		var err error
		cachedPrivateKey, err = rsa.GenerateKey(rand.Reader, 512)
		Expect(err).NotTo(HaveOccurred())
	})

	n := base64.RawURLEncoding.EncodeToString(cachedPrivateKey.N.Bytes())
	e := base64.RawURLEncoding.EncodeToString(big.NewInt(0).SetUint64(uint64(cachedPrivateKey.E)).Bytes())

	tokenToSign := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"foo": "bar",
		"aud": "test-clientid",
		"sub": "user",
		"iss": "http://localhost:5556",
	})
	tokenToSign.Header["kid"] = "test-123"
	token, err := tokenToSign.SignedString(cachedPrivateKey)
	Expect(err).NotTo(HaveOccurred())

	f.s = http.Server{
		Addr: ":5556",
	}

	f.s.Handler = http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("content-type", "application/json")

		switch r.URL.Path {
		case "/auth":
			_, _ = rw.Write([]byte(`auth`))
		case "/.well-known/openid-configuration":
			_, _ = rw.Write([]byte(`
		{
			"issuer": "http://localhost:5556",
			"authorization_endpoint": "http://localhost:5556/auth",
			"token_endpoint": "http://localhost:5556/token",
			"jwks_uri": "http://localhost:5556/keys",
			"response_types_supported": [
			  "code"
			],
			"subject_types_supported": [
			  "public"
			],
			"id_token_signing_alg_values_supported": [
			  "RS256"
			],
			"scopes_supported": [
			  "openid",
			  "email",
			  "profile"
			]
		  }
		`))
		case "/token":
			_, _ = rw.Write([]byte(`
			{
				"access_token": "SlAV32hkKG",
				"token_type": "Bearer",
				"refresh_token": "8xLOxBtZp8",
				"expires_in": 3600,
				"id_token": "` + token + `"
			 }
	`))
		case "/keys":
			_, _ = rw.Write([]byte(`
		{
			"keys": [
			  {
				"use": "sig",
				"kty": "RSA",
				"kid": "test-123",
				"alg": "RS256",
				"n": "` + n + `",
				"e": "` + e + `"
			  }
			]
		  }
		`))

		}
	})

	go func() {
		defer GinkgoRecover()
		err := f.s.ListenAndServe()
		if err != http.ErrServerClosed {
			Expect(err).NotTo(HaveOccurred())
		}
	}()

	return cachedPrivateKey
}

func getOauthConfig(secretRef core.ResourceRef) *extauth.OAuth {
	return &extauth.OAuth{
		ClientId:        "test-clientid",
		ClientSecretRef: &secretRef,
		IssuerUrl:       "http://localhost:5556/",
		AppUrl:          "http://example.com",
		CallbackPath:    "/callback",
		Scopes:          []string{"email"},
	}
}

func getProxyExtAuthOIDC(envoyPort uint32, upstream core.ResourceRef) *gloov1.Proxy {
	return getProxyExtAuth(envoyPort, upstream, getOidcExtAuthExtension())
}

func getOidcExtAuthExtension() *extauth.ExtAuthExtension {
	return &extauth.ExtAuthExtension{
		Spec: &extauth.ExtAuthExtension_ConfigRef{
			ConfigRef: &core.ResourceRef{
				Name:      "oidc-auth",
				Namespace: defaults.GlooSystem,
			},
		},
	}
}

func getProxyExtAuthOIDCAndOpa(envoyPort uint32, secretRef, upstream core.ResourceRef, modules []*core.ResourceRef) *gloov1.Proxy {
	return getProxyExtAuth(envoyPort, upstream, getOidcAndOpaExtAuthExtension())
}

func getOidcAndOpaExtAuthExtension() *extauth.ExtAuthExtension {
	return &extauth.ExtAuthExtension{
		Spec: &extauth.ExtAuthExtension_ConfigRef{
			ConfigRef: &core.ResourceRef{
				Name:      "oidcand-opa-auth",
				Namespace: defaults.GlooSystem,
			},
		},
	}
}

func getOpaConfig(modules []*core.ResourceRef) *extauth.OpaAuth {
	return &extauth.OpaAuth{
		Modules: modules,
		Query:   "data.test.allow == true",
	}
}

func getProxyExtAuthBasicAuth(envoyPort uint32, upstream core.ResourceRef) *gloov1.Proxy {
	return getProxyExtAuth(envoyPort, upstream, GetBasicAuthExtension())
}

//TODO(kdorosh) make sure no flakes here and that order doesn't matter
func GetBasicAuthExtension() *extauth.ExtAuthExtension {
	return &extauth.ExtAuthExtension{
		Spec: &extauth.ExtAuthExtension_ConfigRef{
			ConfigRef: &core.ResourceRef{
				Name:      "basic-auth",
				Namespace: defaults.GlooSystem,
			},
		},
	}
}

func getBasicAuthConfig() *extauth.BasicAuth {
	return &extauth.BasicAuth{
		Realm: "gloo",
		Apr: &extauth.BasicAuth_Apr{
			Users: map[string]*extauth.BasicAuth_Apr_SaltedHashedPassword{
				"user": {
					// Password is password
					Salt:           "0adzfifo",
					HashedPassword: "14o4fMw/Pm2L34SvyyA2r.",
				},
			},
		},
	}
}

func getProxyExtAuthApiKeyAuth(envoyPort uint32, upstream core.ResourceRef) *gloov1.Proxy {
	return getProxyExtAuth(envoyPort, upstream, getApiKeyExtAuthExtension())
}

func getApiKeyAuthConfig() *extauth.ApiKeyAuth {
	return &extauth.ApiKeyAuth{
		ApiKeySecretRefs: []*core.ResourceRef{
			{
				Namespace: "default",
				Name:      "secret1",
			},
		},
		LabelSelector: map[string]string{"team": "infrastructure"},
	}
}

func getApiKeyExtAuthExtension() *extauth.ExtAuthExtension {
	return &extauth.ExtAuthExtension{
		Spec: &extauth.ExtAuthExtension_ConfigRef{
			ConfigRef: &core.ResourceRef{
				Name:      "apikey-auth",
				Namespace: defaults.GlooSystem,
			},
		},
	}
}

func getProxyExtAuth(envoyPort uint32, upstream core.ResourceRef, extauthCfg *extauth.ExtAuthExtension) *gloov1.Proxy {
	var vhosts []*gloov1.VirtualHost

	vhost := &gloov1.VirtualHost{
		Name:    "gloo-system.virt1",
		Domains: []string{"*"},
		Options: &gloov1.VirtualHostOptions{
			Extauth: extauthCfg,
		},
		Routes: []*gloov1.Route{{
			Action: &gloov1.Route_RouteAction{
				RouteAction: &gloov1.RouteAction{
					Destination: &gloov1.RouteAction_Single{
						Single: &gloov1.Destination{
							DestinationType: &gloov1.Destination_Upstream{
								Upstream: utils.ResourceRefPtr(upstream),
							},
						},
					},
				},
			},
		}},
	}

	vhosts = append(vhosts, vhost)

	p := &gloov1.Proxy{
		Metadata: core.Metadata{
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
