package jwt_test

import (
	"crypto/x509"
	"encoding/json"
	"fmt"

	v3 "github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/extensions/filters/http/jwt_authn/v3"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/core/matchers"

	envoy_config_route_v3 "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/duration"
	"github.com/golang/protobuf/ptypes/empty"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	envoycore "github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/config/core/v3"
	. "github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/extensions/jwt"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/jwt"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins"
	"github.com/solo-io/gloo/projects/gloo/pkg/translator"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	test_matchers "github.com/solo-io/solo-kit/test/matchers"
	. "github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/jwt"
	"gopkg.in/square/go-jose.v2"
)

var _ = Describe("JWT Plugin", func() {

	var (
		plugin      *Plugin
		params      plugins.Params
		vhostParams plugins.VirtualHostParams
		routeParams plugins.RouteParams
		virtualHost *v1.VirtualHost
		route       *v1.Route
	)

	const (
		publicKey = `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA4XbzUpqbgKbDLngsLp4b
pjf04WkMzXx8QsZAorkuGprIc2BYVwAmWD2tZvez4769QfXsohu85NRviYsrqbyC
w/NTs3fMlcgld+ayfb/1X3+6u4f1Q8JsDm4fkSWoBUlTkWO7Mcts2hF8OJ8LlGSw
zUDj3TJLQXwtfM0Ty1VzGJQMJELeBuOYHl/jaTdGogI8zbhDZ986CaIfO+q/UM5u
kDA3NJ7oBQEH78N6BTsFpjDUKeTae883CCsRDbsytWgfKT8oA7C4BFkvRqVMSek7
FYkg7AesknSyCIVMObSaf6ZO3T2jVGrWc0iKfrR3Oo7WpiMH84SdBYXPaS1VdLC1
7QIDAQAB
-----END PUBLIC KEY-----
	`
		jwk  = `{"use":"sig","kty":"RSA","alg":"RS256","n":"4XbzUpqbgKbDLngsLp4bpjf04WkMzXx8QsZAorkuGprIc2BYVwAmWD2tZvez4769QfXsohu85NRviYsrqbyCw_NTs3fMlcgld-ayfb_1X3-6u4f1Q8JsDm4fkSWoBUlTkWO7Mcts2hF8OJ8LlGSwzUDj3TJLQXwtfM0Ty1VzGJQMJELeBuOYHl_jaTdGogI8zbhDZ986CaIfO-q_UM5ukDA3NJ7oBQEH78N6BTsFpjDUKeTae883CCsRDbsytWgfKT8oA7C4BFkvRqVMSek7FYkg7AesknSyCIVMObSaf6ZO3T2jVGrWc0iKfrR3Oo7WpiMH84SdBYXPaS1VdLC17Q","e":"AQAB"}`
		jwks = `{"keys":[` + jwk + `]}`
	)
	var (
		keySet jose.JSONWebKeySet
	)

	BeforeEach(func() {
		err := json.Unmarshal([]byte(jwks), &keySet)
		Expect(err).NotTo(HaveOccurred())
	})

	BeforeEach(func() {
		plugin = NewPlugin()
		plugin.Init(plugins.InitParams{})

	})

	Context("With deprecated JWT api", func() {
		var (
			jwtVhost *jwt.VhostExtension
		)
		loadSnapshot := func() {
			virtualHost = &v1.VirtualHost{
				Name:    "virt1",
				Domains: []string{"*"},
				Options: &v1.VirtualHostOptions{
					// Use deprecated Jwt config
					JwtConfig: &v1.VirtualHostOptions_Jwt{
						Jwt: jwtVhost,
					},
				},
				Routes: []*v1.Route{route},
			}

			proxy := &v1.Proxy{
				Metadata: &core.Metadata{
					Name:      "secret",
					Namespace: "default",
				},
				Listeners: []*v1.Listener{{
					Name: "default",
					ListenerType: &v1.Listener_HttpListener{
						HttpListener: &v1.HttpListener{
							VirtualHosts: []*v1.VirtualHost{virtualHost},
						},
					},
				}},
			}

			params.Snapshot = &v1.ApiSnapshot{
				Proxies: v1.ProxyList{proxy},
			}
			vhostParams = plugins.VirtualHostParams{
				Params:   params,
				Proxy:    proxy,
				Listener: proxy.Listeners[0],
			}
			routeParams = plugins.RouteParams{
				VirtualHostParams: vhostParams,
				VirtualHost:       virtualHost,
			}
		}

		BeforeEach(func() {
			jwtDisabledRoute := &jwt.RouteExtension{
				Disable: true,
			}

			route = &v1.Route{
				Name: "some-route",
				Matchers: []*matchers.Matcher{{
					PathSpecifier: &matchers.Matcher_Prefix{
						Prefix: "/",
					},
				}},
				Action: &v1.Route_DirectResponseAction{
					DirectResponseAction: &v1.DirectResponseAction{
						Status: 200,
						Body:   "test",
					},
				},
				Options: &v1.RouteOptions{
					JwtConfig: &v1.RouteOptions_Jwt{
						Jwt: jwtDisabledRoute,
					},
				},
			}

			jwtVhost = &jwt.VhostExtension{
				Providers: map[string]*jwt.Provider{
					"provider1": {
						Jwks: &jwt.Jwks{
							Jwks: &jwt.Jwks_Remote{
								Remote: &jwt.RemoteJwks{
									CacheDuration: &duration.Duration{Seconds: 5},
									Url:           "testium",
									UpstreamRef:   &core.ResourceRef{Name: "test", Namespace: "testns"},
								},
							},
						},
						Audiences: []string{"testaud"},
						Issuer:    "testiss",
					},
				},
			}
		})

		JustBeforeEach(func() {
			loadSnapshot()
		})

		Context("Process snapshot", func() {
			var (
				outRoute     envoy_config_route_v3.Route
				outVhost     envoy_config_route_v3.VirtualHost
				outFilters   []plugins.StagedHttpFilter
				keySetString []byte
				cfg          *JwtWithStage
			)
			JustBeforeEach(func() {
				outVhost = envoy_config_route_v3.VirtualHost{
					Name: "test",
				}
				outRoute = envoy_config_route_v3.Route{}

				// run it like the translator:
				err := plugin.ProcessRoute(routeParams, route, &outRoute)
				Expect(err).NotTo(HaveOccurred())
				err = plugin.ProcessVirtualHost(vhostParams, virtualHost, &outVhost)
				Expect(err).NotTo(HaveOccurred())
				outFilters, err = plugin.HttpFilters(params, nil)
				Expect(err).NotTo(HaveOccurred())

				Expect(outFilters).To(HaveLen(1))
				filter := outFilters[0]
				cfgSt := filter.HttpFilter.GetTypedConfig()
				cfg = &JwtWithStage{}
				err = ptypes.UnmarshalAny(cfgSt, cfg)
				Expect(err).NotTo(HaveOccurred())
			})

			BeforeEach(func() {
				keySetString, _ = json.Marshal(&keySet)
			})

			It("plugin initializes/clears previous config", func() {
				oldConfig := virtualHost.Options.JwtConfig
				virtualHost.Options.JwtConfig = &v1.VirtualHostOptions_JwtStaged{
					JwtStaged: &jwt.JwtStagedVhostExtension{
						BeforeExtAuth: jwtVhost,
					},
				}
				err := plugin.ProcessVirtualHost(vhostParams, virtualHost, &outVhost)
				Expect(err).NotTo(HaveOccurred())
				outFilters, err := plugin.HttpFilters(params, nil)
				Expect(err).NotTo(HaveOccurred())

				Expect(outFilters).To(HaveLen(2))
				// re-initialize plugin
				virtualHost.Options.JwtConfig = oldConfig
				plugin.Init(plugins.InitParams{})
				err = plugin.ProcessVirtualHost(vhostParams, virtualHost, &outVhost)
				Expect(err).NotTo(HaveOccurred())
				outFilters, err = plugin.HttpFilters(params, nil)
				Expect(err).NotTo(HaveOccurred())

				Expect(outFilters).To(HaveLen(1))
			})

			It("should process virtual host", func() {

				goTpfc := outVhost.TypedPerFilterConfig[SoloJwtFilterName]
				Expect(goTpfc).NotTo(BeNil())
				var routeCfg StagedJwtAuthnPerRoute
				err := ptypes.UnmarshalAny(goTpfc, &routeCfg)
				Expect(err).NotTo(HaveOccurred())
				Expect(routeCfg.GetJwtConfigs()[AfterExtAuthStage].Requirement).To(Equal(virtualHost.Name))
			})

			Context("allow missing or failed JWT VirtualHost", func() {
				BeforeEach(func() {
					loadSnapshot()
					virtualHost.GetOptions().GetJwt().AllowMissingOrFailedJwt = true
				})

				It("should process virtual host that will allow JWT fail or missing JWT", func() {

					req := cfg.JwtAuthn.FilterStateRules.GetRequires()[virtualHost.Name]
					expectedReq := new(v3.JwtRequirement_RequiresAny)
					requiresAny := req.RequiresType
					Expect(requiresAny).To(BeAssignableToTypeOf(expectedReq))
					ORreqs := requiresAny.(*v3.JwtRequirement_RequiresAny).RequiresAny.GetRequirements()
					// Make sure that the virtual host has the allow_missing_or_failed req with at least one other req
					Expect(ORreqs).To(HaveLen(2))
					missingOrFailedReq := &v3.JwtRequirement{
						RequiresType: &v3.JwtRequirement_AllowMissingOrFailed{
							AllowMissingOrFailed: &empty.Empty{},
						},
					}
					Expect(ORreqs).To(ContainElement(missingOrFailedReq))
				})
			})
			Context("allow missing or failed JWT Route", func() {

				BeforeEach(func() {
					loadSnapshot()
					route.GetOptions().GetJwt().Disable = false
				})
			})

			It("should process regular route", func() {
				goTpfc := outRoute.TypedPerFilterConfig[SoloJwtFilterName]
				Expect(goTpfc).NotTo(BeNil())
				var routeCfg StagedJwtAuthnPerRoute
				err := ptypes.UnmarshalAny(goTpfc, &routeCfg)
				Expect(err).NotTo(HaveOccurred())
				Expect(routeCfg.JwtConfigs[AfterExtAuthStage].Requirement).To(Equal(DisableName))
			})

			It("should process filters", func() {
				// Expect a requirement
				providerName := ProviderName(virtualHost.Name, "provider1")
				expectedCfg := &JwtWithStage{
					JwtAuthn: &v3.JwtAuthentication{
						Providers: map[string]*v3.JwtProvider{
							providerName: {
								Issuer:            jwtVhost.Providers["provider1"].Issuer,
								Audiences:         jwtVhost.Providers["provider1"].Audiences,
								PayloadInMetadata: providerName,
								JwksSourceSpecifier: &v3.JwtProvider_RemoteJwks{
									RemoteJwks: &v3.RemoteJwks{
										CacheDuration: &duration.Duration{Seconds: 5},
										HttpUri: &envoycore.HttpUri{
											Timeout: &duration.Duration{Seconds: RemoteJwksTimeoutSecs},
											Uri:     jwtVhost.Providers["provider1"].GetJwks().GetRemote().Url,
											HttpUpstreamType: &envoycore.HttpUri_Cluster{
												Cluster: translator.UpstreamToClusterName(jwtVhost.Providers["provider1"].GetJwks().GetRemote().UpstreamRef),
											},
										},
									},
								},
							},
						},
						FilterStateRules: &v3.FilterStateRule{
							Name: getFilterStateName(0),
							Requires: map[string]*v3.JwtRequirement{
								virtualHost.Name: {
									RequiresType: &v3.JwtRequirement_ProviderName{
										ProviderName: providerName,
									},
								},
							},
						},
					},
				}
				Expect(expectedCfg).To(test_matchers.MatchProto(cfg))
			})

			Context("local jwks", func() {
				BeforeEach(func() {
					jwtVhost = &jwt.VhostExtension{
						Providers: map[string]*jwt.Provider{
							"default": {
								Jwks: &jwt.Jwks{
									Jwks: &jwt.Jwks_Local{
										Local: &jwt.LocalJwks{
											Key: jwks,
										},
									},
								},
								Audiences: []string{"testaud"},
								Issuer:    "testiss",
							},
						},
					}
				})

				It("should process filters", func() {

					// Expect a requirement
					providerName := ProviderName(virtualHost.Name, "default")
					expectedCfg := &JwtWithStage{
						JwtAuthn: &v3.JwtAuthentication{
							Providers: map[string]*v3.JwtProvider{
								providerName: {
									Issuer:            jwtVhost.Providers["default"].Issuer,
									Audiences:         jwtVhost.Providers["default"].Audiences,
									PayloadInMetadata: providerName,
									JwksSourceSpecifier: &v3.JwtProvider_LocalJwks{
										LocalJwks: &envoycore.DataSource{
											Specifier: &envoycore.DataSource_InlineString{
												InlineString: string(keySetString),
											},
										},
									},
								},
							},
							FilterStateRules: &v3.FilterStateRule{
								Name: getFilterStateName(0),
								Requires: map[string]*v3.JwtRequirement{
									virtualHost.Name: {
										RequiresType: &v3.JwtRequirement_ProviderName{
											ProviderName: providerName,
										},
									},
								},
							},
						},
					}
					Expect(expectedCfg).To(test_matchers.MatchProto(cfg))
				})
			})

			Context("claims to headers", func() {
				BeforeEach(func() {
					route.Options = nil
					jwks := &jwt.Jwks{
						Jwks: &jwt.Jwks_Local{
							Local: &jwt.LocalJwks{
								Key: jwks,
							},
						},
					}
					jwtVhost = &jwt.VhostExtension{
						Providers: map[string]*jwt.Provider{
							"provider1": {
								Jwks:      jwks,
								Audiences: []string{"testaud1"},
								Issuer:    "testiss1",
								ClaimsToHeaders: []*jwt.ClaimToHeader{{
									Claim:  "sub",
									Header: "x-sub",
									Append: true,
								}},
							},
						},
					}
				})

				It("should translate claims to headers", func() {

					goTpfc := outVhost.TypedPerFilterConfig[SoloJwtFilterName]
					Expect(goTpfc).NotTo(BeNil())
					var routeCfg StagedJwtAuthnPerRoute
					err := ptypes.UnmarshalAny(goTpfc, &routeCfg)
					Expect(err).NotTo(HaveOccurred())

					provider1Name := ProviderName(virtualHost.Name, "provider1")
					expectedCfg := &StagedJwtAuthnPerRoute{
						JwtConfigs: map[uint32]*SoloJwtAuthnPerRoute{
							AfterExtAuthStage: {
								Requirement: virtualHost.Name,
								ClaimsToHeaders: map[string]*SoloJwtAuthnPerRoute_ClaimToHeaders{
									provider1Name: {
										Claims: []*SoloJwtAuthnPerRoute_ClaimToHeader{{
											Claim:  "sub",
											Header: "x-sub",
											Append: true,
										}},
									},
								},
								ClearRouteCache:   true,
								PayloadInMetadata: PayloadInMetadata,
							},
						},
					}
					Expect(expectedCfg).To(test_matchers.MatchProto(&routeCfg))
				})
			})

			Context("claims token source", func() {
				BeforeEach(func() {
					route.Options = nil
					jwks := &jwt.Jwks{
						Jwks: &jwt.Jwks_Local{
							Local: &jwt.LocalJwks{
								Key: jwks,
							},
						},
					}
					jwtVhost = &jwt.VhostExtension{
						Providers: map[string]*jwt.Provider{
							"provider1": {
								Jwks:      jwks,
								Audiences: []string{"testaud1"},
								Issuer:    "testiss1",
								TokenSource: &jwt.TokenSource{
									Headers: []*jwt.TokenSource_HeaderSource{{
										Header: "header",
										Prefix: "prefix",
									}},
									QueryParams: []string{"query"},
								},
							},
						},
					}
					virtualHost = &v1.VirtualHost{
						Name:    "virt1",
						Domains: []string{"*"},
						Options: &v1.VirtualHostOptions{
							JwtConfig: &v1.VirtualHostOptions_Jwt{
								Jwt: jwtVhost,
							},
						},
						Routes: []*v1.Route{route},
					}
				})

				It("should translate token source", func() {
					provider1Name := ProviderName(virtualHost.Name, "provider1")
					expectedCfg := &JwtWithStage{
						JwtAuthn: &v3.JwtAuthentication{
							Providers: map[string]*v3.JwtProvider{
								provider1Name: {
									Issuer:            "testiss1",
									Audiences:         []string{"testaud1"},
									PayloadInMetadata: provider1Name,
									FromHeaders: []*v3.JwtHeader{{
										Name:        "header",
										ValuePrefix: "prefix",
									}},
									FromParams: []string{"query"},
									JwksSourceSpecifier: &v3.JwtProvider_LocalJwks{
										LocalJwks: &envoycore.DataSource{
											Specifier: &envoycore.DataSource_InlineString{
												InlineString: string(keySetString),
											},
										},
									},
								},
							},
							FilterStateRules: &v3.FilterStateRule{
								Name: getFilterStateName(0),
								Requires: map[string]*v3.JwtRequirement{
									virtualHost.Name: {
										RequiresType: &v3.JwtRequirement_ProviderName{
											ProviderName: provider1Name,
										},
									},
								},
							},
						},
					}
					Expect(expectedCfg).To(test_matchers.MatchProto(cfg))
				})
			})

			Context("multiple providers", func() {
				BeforeEach(func() {
					jwks := &jwt.Jwks{
						Jwks: &jwt.Jwks_Local{
							Local: &jwt.LocalJwks{
								Key: jwks,
							},
						},
					}
					jwtVhost = &jwt.VhostExtension{
						Providers: map[string]*jwt.Provider{
							"provider1": {
								Jwks:      jwks,
								Audiences: []string{"testaud1"},
								Issuer:    "testiss1",
							},
							"provider2": {
								Jwks:      jwks,
								Audiences: []string{"testaud2"},
								Issuer:    "testiss2",
							},
						},
					}
				})

				It("should translate multiple providers", func() {
					// Expect a requirement
					provider1Name := ProviderName(virtualHost.Name, "provider1")
					provider2Name := ProviderName(virtualHost.Name, "provider2")
					expectedCfg := &JwtWithStage{
						JwtAuthn: &v3.JwtAuthentication{
							Providers: map[string]*v3.JwtProvider{
								provider1Name: {
									Issuer:            "testiss1",
									Audiences:         []string{"testaud1"},
									PayloadInMetadata: provider1Name,
									JwksSourceSpecifier: &v3.JwtProvider_LocalJwks{
										LocalJwks: &envoycore.DataSource{
											Specifier: &envoycore.DataSource_InlineString{
												InlineString: string(keySetString),
											},
										},
									},
								},
								provider2Name: {
									Issuer:            "testiss2",
									Audiences:         []string{"testaud2"},
									PayloadInMetadata: provider2Name,
									JwksSourceSpecifier: &v3.JwtProvider_LocalJwks{
										LocalJwks: &envoycore.DataSource{
											Specifier: &envoycore.DataSource_InlineString{
												InlineString: string(keySetString),
											},
										},
									},
								},
							},
							FilterStateRules: &v3.FilterStateRule{
								Name: getFilterStateName(0),
								Requires: map[string]*v3.JwtRequirement{
									virtualHost.Name: {
										RequiresType: &v3.JwtRequirement_RequiresAny{
											RequiresAny: &v3.JwtRequirementOrList{
												Requirements: []*v3.JwtRequirement{
													{
														RequiresType: &v3.JwtRequirement_ProviderName{
															ProviderName: provider1Name,
														},
													}, {
														RequiresType: &v3.JwtRequirement_ProviderName{
															ProviderName: provider2Name,
														},
													},
												},
											},
										},
									},
								},
							},
						},
					}
					Expect(expectedCfg).To(test_matchers.MatchProto(cfg))
				})

			})

		})

		Context("translate key", func() {

			It("should translate PEM", func() {
				jwks, err := TranslateKey(publicKey)
				Expect(err).NotTo(HaveOccurred())
				// make certs empty and not nil for comparison
				jwks.Keys[0].Certificates = make([]*x509.Certificate, 0)
				// set to empty array instead of nil, so assert below will work
				jwks.Keys[0].CertificateThumbprintSHA1 = []byte{}
				jwks.Keys[0].CertificateThumbprintSHA256 = []byte{}
				Expect(*jwks).To(BeEquivalentTo(keySet))
			})
			It("should translate key", func() {
				jwks, err := TranslateKey(jwk)
				Expect(err).NotTo(HaveOccurred())
				Expect(*jwks).To(BeEquivalentTo(keySet))
			})
			It("should translate key set", func() {
				jwks, err := TranslateKey(jwks)
				Expect(err).NotTo(HaveOccurred())
				Expect(*jwks).To(BeEquivalentTo(keySet))
			})
		})
	})

	Context("With staged JWT api", func() {
		var beforeExtAuthJwtVhost, afterExtauthJwtVhost *jwt.VhostExtension

		BeforeEach(func() {

			beforeExtAuthJwtVhost = &jwt.VhostExtension{
				Providers: map[string]*jwt.Provider{
					"before_provider": {
						Jwks: &jwt.Jwks{
							Jwks: &jwt.Jwks_Remote{
								Remote: &jwt.RemoteJwks{
									CacheDuration: &duration.Duration{Seconds: 5},
									Url:           "before_testium",
									UpstreamRef:   &core.ResourceRef{Name: "test", Namespace: "testns"},
									AsyncFetch: &v3.JwksAsyncFetch{
										FastListener: true,
									},
								},
							},
						},
						Audiences: []string{"before_testaud"},
						Issuer:    "before_testiss",
					},
				},
			}

			afterExtauthJwtVhost = &jwt.VhostExtension{
				Providers: map[string]*jwt.Provider{
					"after_provider": {
						Jwks: &jwt.Jwks{
							Jwks: &jwt.Jwks_Remote{
								Remote: &jwt.RemoteJwks{
									CacheDuration: &duration.Duration{Seconds: 5},
									Url:           "after_testium",
									UpstreamRef:   &core.ResourceRef{Name: "test", Namespace: "testns"},
								},
							},
						},
						Audiences: []string{"after_testaud"},
						Issuer:    "after_testiss",
					},
				},
			}

			route = &v1.Route{
				Name: "some-route",
				Matchers: []*matchers.Matcher{{
					PathSpecifier: &matchers.Matcher_Prefix{
						Prefix: "/",
					},
				}},
				Action: &v1.Route_DirectResponseAction{
					DirectResponseAction: &v1.DirectResponseAction{
						Status: 200,
						Body:   "test",
					},
				},
				Options: &v1.RouteOptions{
					// Use staged Jwt config
					JwtConfig: &v1.RouteOptions_JwtStaged{
						JwtStaged: &jwt.JwtStagedRouteExtension{
							AfterExtAuth:  &jwt.RouteExtension{Disable: true},
							BeforeExtAuth: &jwt.RouteExtension{Disable: true},
						},
					},
				},
			}

			virtualHost = &v1.VirtualHost{
				Name:    "virt1",
				Domains: []string{"*"},
				Options: &v1.VirtualHostOptions{
					// Use staged Jwt config
					JwtConfig: &v1.VirtualHostOptions_JwtStaged{
						JwtStaged: &jwt.JwtStagedVhostExtension{
							AfterExtAuth:  afterExtauthJwtVhost,
							BeforeExtAuth: beforeExtAuthJwtVhost,
						},
					},
				},
				Routes: []*v1.Route{route},
			}

		})
		loadSnapshot := func() {
			proxy := &v1.Proxy{
				Metadata: &core.Metadata{
					Name:      "secret",
					Namespace: "default",
				},
				Listeners: []*v1.Listener{{
					Name: "default",
					ListenerType: &v1.Listener_HttpListener{
						HttpListener: &v1.HttpListener{
							VirtualHosts: []*v1.VirtualHost{virtualHost},
						},
					},
				}},
			}

			params.Snapshot = &v1.ApiSnapshot{
				Proxies: v1.ProxyList{proxy},
			}
			vhostParams = plugins.VirtualHostParams{
				Params:   params,
				Proxy:    proxy,
				Listener: proxy.Listeners[0],
			}
			routeParams = plugins.RouteParams{
				VirtualHostParams: vhostParams,
				VirtualHost:       virtualHost,
			}
		}

		JustBeforeEach(func() {
			loadSnapshot()
		})

		Context("Process snapshot", func() {

			var (
				outRoute            envoy_config_route_v3.Route
				outVhost            envoy_config_route_v3.VirtualHost
				outFilters          []plugins.StagedHttpFilter
				beforeCfg, afterCfg *JwtWithStage
			)
			JustBeforeEach(func() {
				outVhost = envoy_config_route_v3.VirtualHost{
					Name: "test",
				}
				outRoute = envoy_config_route_v3.Route{}

				// run it like the translator:
				err := plugin.ProcessRoute(routeParams, route, &outRoute)
				Expect(err).NotTo(HaveOccurred())
				err = plugin.ProcessVirtualHost(vhostParams, virtualHost, &outVhost)
				Expect(err).NotTo(HaveOccurred())
				outFilters, err = plugin.HttpFilters(params, nil)
				Expect(err).NotTo(HaveOccurred())
				Expect(outFilters).To(HaveLen(2))
				beforeCfg = &JwtWithStage{}
				err = ptypes.UnmarshalAny(outFilters[1].HttpFilter.GetTypedConfig(), beforeCfg)
				Expect(err).NotTo(HaveOccurred())
				afterCfg = &JwtWithStage{}
				err = ptypes.UnmarshalAny(outFilters[0].HttpFilter.GetTypedConfig(), afterCfg)
				Expect(err).NotTo(HaveOccurred())
			})

			It("should process virtual host with jwt configs for both stages", func() {
				goTpfc := outVhost.TypedPerFilterConfig[SoloJwtFilterName]
				Expect(goTpfc).NotTo(BeNil())
				var routeCfg StagedJwtAuthnPerRoute
				err := ptypes.UnmarshalAny(goTpfc, &routeCfg)
				Expect(err).NotTo(HaveOccurred())
				Expect(routeCfg.GetJwtConfigs()[AfterExtAuthStage].Requirement).To(Equal(virtualHost.Name))
				Expect(routeCfg.GetJwtConfigs()[BeforeExtAuthStage].Requirement).To(Equal(virtualHost.Name))

			})

			Context("allow missing or failed JWT VirtualHost", func() {
				BeforeEach(func() {
					loadSnapshot()
					virtualHost.GetOptions().GetJwtStaged().AfterExtAuth.AllowMissingOrFailedJwt = true
					virtualHost.GetOptions().GetJwtStaged().BeforeExtAuth.AllowMissingOrFailedJwt = true
				})

				It("should process missing or failed jwt requirements on staged jwt", func() {

					checkAllowMissingOrFailed := func(stage *JwtWithStage) {
						req := stage.JwtAuthn.FilterStateRules.GetRequires()[virtualHost.Name]
						expectedReq := new(v3.JwtRequirement_RequiresAny)
						requiresAny := req.RequiresType
						ExpectWithOffset(1, requiresAny).To(BeAssignableToTypeOf(expectedReq))
						ORreqs := requiresAny.(*v3.JwtRequirement_RequiresAny).RequiresAny.GetRequirements()
						// Make sure that the virtual host has the allow_missing_or_failed req with at least one other req
						ExpectWithOffset(1, ORreqs).To(HaveLen(2))
						missingOrFailedReq := &v3.JwtRequirement{
							RequiresType: &v3.JwtRequirement_AllowMissingOrFailed{
								AllowMissingOrFailed: &empty.Empty{},
							},
						}
						ExpectWithOffset(1, ORreqs).To(ContainElement(missingOrFailedReq))
					}

					checkAllowMissingOrFailed(beforeCfg)
					checkAllowMissingOrFailed(afterCfg)
				})
			})

			It("should process regular route with jwt configs for both stages", func() {
				goTpfc := outRoute.TypedPerFilterConfig[SoloJwtFilterName]
				Expect(goTpfc).NotTo(BeNil())
				var routeCfg StagedJwtAuthnPerRoute
				err := ptypes.UnmarshalAny(goTpfc, &routeCfg)
				Expect(err).NotTo(HaveOccurred())
				Expect(routeCfg.JwtConfigs[AfterExtAuthStage].Requirement).To(Equal(DisableName))
				Expect(routeCfg.JwtConfigs[BeforeExtAuthStage].Requirement).To(Equal(DisableName))

			})

			It("should process filters", func() {
				// Expect a requirement
				providerName := ProviderName(virtualHost.Name, "after_provider")
				expectedCfg := &JwtWithStage{
					JwtAuthn: &v3.JwtAuthentication{
						Providers: map[string]*v3.JwtProvider{
							providerName: {
								Issuer:            afterExtauthJwtVhost.Providers["after_provider"].Issuer,
								Audiences:         afterExtauthJwtVhost.Providers["after_provider"].Audiences,
								PayloadInMetadata: providerName,
								JwksSourceSpecifier: &v3.JwtProvider_RemoteJwks{
									RemoteJwks: &v3.RemoteJwks{
										CacheDuration: &duration.Duration{Seconds: 5},
										HttpUri: &envoycore.HttpUri{
											Timeout: &duration.Duration{Seconds: RemoteJwksTimeoutSecs},
											Uri:     afterExtauthJwtVhost.Providers["after_provider"].GetJwks().GetRemote().Url,
											HttpUpstreamType: &envoycore.HttpUri_Cluster{
												Cluster: translator.UpstreamToClusterName(afterExtauthJwtVhost.Providers["after_provider"].GetJwks().GetRemote().UpstreamRef),
											},
										},
									},
								},
							},
						},
						FilterStateRules: &v3.FilterStateRule{
							Name: getFilterStateName(0),
							Requires: map[string]*v3.JwtRequirement{
								virtualHost.Name: {
									RequiresType: &v3.JwtRequirement_ProviderName{
										ProviderName: providerName,
									},
								},
							},
						},
					},
				}
				Expect(expectedCfg).To(test_matchers.MatchProto(afterCfg))
				providerName = ProviderName(virtualHost.Name, "before_provider")
				expectedCfg = &JwtWithStage{
					JwtAuthn: &v3.JwtAuthentication{
						Providers: map[string]*v3.JwtProvider{
							providerName: {
								Issuer:            beforeExtAuthJwtVhost.Providers["before_provider"].Issuer,
								Audiences:         beforeExtAuthJwtVhost.Providers["before_provider"].Audiences,
								PayloadInMetadata: providerName,
								JwksSourceSpecifier: &v3.JwtProvider_RemoteJwks{
									RemoteJwks: &v3.RemoteJwks{
										CacheDuration: &duration.Duration{Seconds: 5},
										HttpUri: &envoycore.HttpUri{
											Timeout: &duration.Duration{Seconds: RemoteJwksTimeoutSecs},
											Uri:     beforeExtAuthJwtVhost.Providers["before_provider"].GetJwks().GetRemote().Url,
											HttpUpstreamType: &envoycore.HttpUri_Cluster{
												Cluster: translator.UpstreamToClusterName(beforeExtAuthJwtVhost.Providers["before_provider"].GetJwks().GetRemote().UpstreamRef),
											},
										},
										AsyncFetch: &v3.JwksAsyncFetch{
											FastListener: true,
										},
									},
								},
							},
						},
						FilterStateRules: &v3.FilterStateRule{
							Name: getFilterStateName(1),
							Requires: map[string]*v3.JwtRequirement{
								virtualHost.Name: {
									RequiresType: &v3.JwtRequirement_ProviderName{
										ProviderName: providerName,
									},
								},
							},
						},
					},
					Stage: 1,
				}
				Expect(expectedCfg).To(test_matchers.MatchProto(beforeCfg))
			})
		})
	})
})

func getFilterStateName(stage uint32) string {
	return fmt.Sprintf("stage%d-%s", stage, StateName)
}
