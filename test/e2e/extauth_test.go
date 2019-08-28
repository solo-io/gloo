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

	envoyutil "github.com/envoyproxy/go-control-plane/pkg/util"
	extauth2 "github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/extauth"

	"github.com/gogo/protobuf/types"

	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/config"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"

	"github.com/solo-io/ext-auth-service/pkg/config/oidc"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	extauthrunner "github.com/solo-io/solo-projects/projects/extauth/pkg/runner"
	"github.com/solo-io/solo-projects/test/services"

	"github.com/solo-io/solo-projects/projects/gloo/pkg/api/v1/plugins/extauth"

	"github.com/dgrijalva/jwt-go"
	"github.com/fgrosse/zaptest"
	"github.com/solo-io/gloo/pkg/utils"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	gloov1static "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/plugins/static"
	"github.com/solo-io/gloo/projects/gloo/pkg/defaults"
	"github.com/solo-io/go-utils/contextutils"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/memory"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-projects/test/v1helpers"
)

var (
	baseExtauthPort = uint32(27000)
)

var _ = Describe("External ", func() {

	var (
		ctx         context.Context
		cancel      context.CancelFunc
		testClients services.TestClients
		settings    extauthrunner.Settings
	)

	BeforeEach(func() {
		extauthport := atomic.AddUint32(&baseExtauthPort, 1) + uint32(config.GinkgoConfig.ParallelNode*1000)

		logger := zaptest.LoggerWriter(GinkgoWriter)
		contextutils.SetFallbackLogger(logger.Sugar())

		ctx, cancel = context.WithCancel(context.Background())
		cache := memory.NewInMemoryResourceCache()

		testClients = services.GetTestClients(cache)
		testClients.GlooPort = int(services.AllocateGlooPort())

		extauthAddr := "localhost"
		if runtime.GOOS == "darwin" {
			extauthAddr = "host.docker.internal"
		}

		extauthserver := &gloov1.Upstream{
			Metadata: core.Metadata{
				Name:      "extauth-server",
				Namespace: "default",
			},
			UpstreamSpec: &gloov1.UpstreamSpec{
				UseHttp2: true,
				UpstreamType: &gloov1.UpstreamSpec_Static{
					Static: &gloov1static.UpstreamSpec{
						Hosts: []*gloov1static.Host{{
							Addr: extauthAddr,
							Port: extauthport,
						}},
					},
				},
			},
		}

		testClients.UpstreamClient.Write(extauthserver, clients.WriteOpts{})
		ref := extauthserver.Metadata.Ref()
		extauthSettings := &extauth.Settings{
			ExtauthzServerRef: &ref,
		}
		settingsStruct, err := envoyutil.MessageToStruct(extauthSettings)
		Expect(err).NotTo(HaveOccurred())

		extensions := &gloov1.Extensions{
			Configs: map[string]*types.Struct{
				extauth2.ExtensionName: settingsStruct,
			},
		}

		settings = extauthrunner.Settings{
			GlooAddress:  fmt.Sprintf("localhost:%d", testClients.GlooPort),
			DebugPort:    0,
			ServerPort:   int(extauthport),
			SigningKey:   "hello",
			UserIdHeader: "X-User-Id",
		}

		what := services.What{
			DisableGateway: true,
			DisableUds:     true,
			DisableFds:     true,
		}

		services.RunGlooGatewayUdsFdsOnPort(ctx, cache, int32(testClients.GlooPort), what, defaults.GlooSystem, nil, extensions)
		go func(testctx context.Context) {
			defer GinkgoRecover()
			err := extauthrunner.RunWithSettings(testctx, settings)
			if testctx.Err() == nil {
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
				envoyInstance.Clean()
			}
		})

		Context("basic auth sanity tests", func() {
			BeforeEach(func() {

				// drain channel as we dont care about it
				go func() {
					for range testUpstream.C {
					}
				}()

				proxy := getProxyExtAuthBasicAuth(envoyPort, testUpstream.Upstream.Metadata.Ref())

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
				privatekey      *rsa.PrivateKey
				discoveryServer fakeDiscoveryServer
				secret          *gloov1.Secret
				proxy           *gloov1.Proxy
				token           string
			)
			BeforeEach(func() {
				discoveryServer = fakeDiscoveryServer{}

				privatekey = discoveryServer.Start()

				clientSecret := &extauth.OauthSecret{
					ClientSecret: "test",
				}
				secretStruct, err := envoyutil.MessageToStruct(clientSecret)
				Expect(err).NotTo(HaveOccurred())

				secret = &gloov1.Secret{
					Metadata: core.Metadata{
						Name:      "secret",
						Namespace: "default",
					},
					Kind: &gloov1.Secret_Extension{
						Extension: &gloov1.Extension{
							Config: secretStruct,
						},
					},
				}
				testClients.SecretClient.Write(secret, clients.WriteOpts{})

				proxy = getProxyExtAuthOIDC(envoyPort, secret.Metadata.Ref(), testUpstream.Upstream.Metadata.Ref())

				// create an id token
				tokentosign := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
					"foo": "bar",
					"aud": "test-clientid",
					"sub": "user",
					"iss": "http://127.0.0.1:5556",
				})
				tokentosign.Header["kid"] = "test-123"
				token, err = tokentosign.SignedString(privatekey)
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
						proxy = getProxyExtAuthOIDCAndOpa(envoyPort, secret.Metadata.Ref(), testUpstream.Upstream.Metadata.Ref(), modules)

						_, err := testClients.ArtifactClient.Write(policy, clients.WriteOpts{})
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
					ExpectWithOffset(1, r.Headers["X-User-Id"][0]).To(Equal("http://127.0.0.1:5556;user"))
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
					proxy = getProxyExtAuthOIDCAndOpa(envoyPort, secret.Metadata.Ref(), testUpstream.Upstream.Metadata.Ref(), modules)

					_, err := testClients.ArtifactClient.Write(policy, clients.WriteOpts{})
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

				apiKeySecret1 := &extauth.ApiKeySecret{
					ApiKey: "secretApiKey1",
				}
				secretStruct1, err := envoyutil.MessageToStruct(apiKeySecret1)
				Expect(err).NotTo(HaveOccurred())

				secret1 := &gloov1.Secret{
					Metadata: core.Metadata{
						Name:      "secret1",
						Namespace: "default",
					},
					Kind: &gloov1.Secret_Extension{
						Extension: &gloov1.Extension{
							Config: secretStruct1,
						},
					},
				}

				apiKeySecret2 := &extauth.ApiKeySecret{
					ApiKey: "secretApiKey2",
				}
				secretStruct2, err := envoyutil.MessageToStruct(apiKeySecret2)
				Expect(err).NotTo(HaveOccurred())

				secret2 := &gloov1.Secret{
					Metadata: core.Metadata{
						Name:      "secret2",
						Namespace: "default",
						Labels:    map[string]string{"team": "infrastructure"},
					},
					Kind: &gloov1.Secret_Extension{
						Extension: &gloov1.Extension{
							Config: secretStruct2,
						},
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

var startDiscoveryServerOnce sync.Once
var cachedPrivateKey *rsa.PrivateKey

type fakeDiscoveryServer struct {
	s http.Server
}

func (f *fakeDiscoveryServer) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	f.s.Shutdown(ctx)
}

func (f *fakeDiscoveryServer) Start() *rsa.PrivateKey {
	startDiscoveryServerOnce.Do(func() {
		var err error
		cachedPrivateKey, err = rsa.GenerateKey(rand.Reader, 512)
		Expect(err).NotTo(HaveOccurred())
	})

	n := base64.RawURLEncoding.EncodeToString(cachedPrivateKey.N.Bytes())
	e := base64.RawURLEncoding.EncodeToString(big.NewInt(0).SetUint64(uint64(cachedPrivateKey.E)).Bytes())

	tokentosign := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"foo": "bar",
		"aud": "test-clientid",
		"sub": "user",
		"iss": "http://127.0.0.1:5556",
	})
	tokentosign.Header["kid"] = "test-123"
	token, err := tokentosign.SignedString(cachedPrivateKey)
	Expect(err).NotTo(HaveOccurred())

	f.s = http.Server{
		Addr: ":5556",
	}

	f.s.Handler = http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

		rw.Header().Set("content-type", "application/json")

		switch r.URL.Path {
		case "/auth":
			rw.Write([]byte(`auth`))
		case "/.well-known/openid-configuration":
			rw.Write([]byte(`
		{
			"issuer": "http://127.0.0.1:5556",
			"authorization_endpoint": "http://127.0.0.1:5556/auth",
			"token_endpoint": "http://127.0.0.1:5556/token",
			"jwks_uri": "http://127.0.0.1:5556/keys",
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
			rw.Write([]byte(`
			{
				"access_token": "SlAV32hkKG",
				"token_type": "Bearer",
				"refresh_token": "8xLOxBtZp8",
				"expires_in": 3600,
				"id_token": "` + token + `"
			 }
	`))
		case "/keys":
			rw.Write([]byte(`
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

func getProxyExtAuthOIDC(envoyPort uint32, secretRef, upstream core.ResourceRef) *gloov1.Proxy {
	extauthCfg := &extauth.VhostExtension{
		AuthConfig: &extauth.VhostExtension_Oauth{
			Oauth: getOauthConfig(secretRef),
		},
	}
	return getProxyExtAuth(envoyPort, upstream, extauthCfg)
}

func getProxyExtAuthOIDCAndOpa(envoyPort uint32, secretRef, upstream core.ResourceRef, modules []*core.ResourceRef) *gloov1.Proxy {
	extauthCfg := &extauth.VhostExtension{
		Configs: []*extauth.AuthConfig{
			{
				AuthConfig: &extauth.AuthConfig_Oauth{
					Oauth: getOauthConfig(secretRef),
				},
			},
			{
				AuthConfig: &extauth.AuthConfig_OpaAuth{
					OpaAuth: &extauth.OpaAuth{
						Modules: modules,
						Query:   "data.test.allow == true",
					},
				},
			},
		},
	}
	return getProxyExtAuth(envoyPort, upstream, extauthCfg)
}

func getProxyExtAuthBasicAuth(envoyPort uint32, upstream core.ResourceRef) *gloov1.Proxy {
	extauthCfg := GetBasicAuthExtension()
	return getProxyExtAuth(envoyPort, upstream, extauthCfg)
}

func GetBasicAuthExtension() *extauth.VhostExtension {
	return &extauth.VhostExtension{
		AuthConfig: &extauth.VhostExtension_BasicAuth{
			BasicAuth: &extauth.BasicAuth{
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
			},
		},
	}
}

func getProxyExtAuthApiKeyAuth(envoyPort uint32, upstream core.ResourceRef) *gloov1.Proxy {
	extauthCfg := GetApiKeyAuthExtension()
	return getProxyExtAuth(envoyPort, upstream, extauthCfg)
}

func GetApiKeyAuthExtension() *extauth.VhostExtension {
	return &extauth.VhostExtension{
		AuthConfig: &extauth.VhostExtension_ApiKeyAuth{
			ApiKeyAuth: &extauth.ApiKeyAuth{
				ApiKeySecretRefs: []*core.ResourceRef{
					{
						Namespace: "default",
						Name:      "secret1",
					},
				},
				LabelSelector: map[string]string{"team": "infrastructure"},
			},
		},
	}
}

func getProxyExtAuth(envoyPort uint32, upstream core.ResourceRef, extauthCfg *extauth.VhostExtension) *gloov1.Proxy {
	var extensions *gloov1.Extensions

	extauthStruct, err := envoyutil.MessageToStruct(extauthCfg)
	Expect(err).NotTo(HaveOccurred())
	protos := map[string]*types.Struct{
		extauth2.ExtensionName: extauthStruct,
	}

	extensions = &gloov1.Extensions{
		Configs: protos,
	}

	var vhosts []*gloov1.VirtualHost

	vhost := &gloov1.VirtualHost{
		Name:               "gloo-system.virt1",
		Domains:            []string{"*"},
		VirtualHostPlugins: &gloov1.VirtualHostPlugins{},
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
							DestinationType: &gloov1.Destination_Upstream{
								Upstream: utils.ResourceRefPtr(upstream),
							},
						},
					},
				},
			},
		}},
	}

	vhost.VirtualHostPlugins.Extensions = extensions

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
