package jwt_test

import (
	"crypto/x509"
	"encoding/json"

	envoycore "github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	envoyroute "github.com/envoyproxy/go-control-plane/envoy/api/v2/route"
	envoyauth "github.com/envoyproxy/go-control-plane/envoy/config/filter/http/jwt_authn/v2alpha"
	"github.com/envoyproxy/go-control-plane/pkg/util"
	"github.com/gogo/protobuf/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/square/go-jose.v2"

	. "github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/extensions/jwt"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/plugins/jwt"

	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins"
	"github.com/solo-io/gloo/projects/gloo/pkg/translator"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	. "github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/jwt"
)

var _ = Describe("JWT Plugin", func() {
	var (
		plugin      *Plugin
		params      plugins.Params
		vhostParams plugins.VirtualHostParams
		routeParams plugins.RouteParams
		virtualHost *v1.VirtualHost
		route       *v1.Route
		jwtVhost    *jwt.VhostExtension
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

		jwtRoute := &jwt.RouteExtension{
			Disable: true,
		}
		jwtRouteSt, err := util.MessageToStruct(jwtRoute)
		Expect(err).NotTo(HaveOccurred())

		route = &v1.Route{
			Matcher: &v1.Matcher{
				PathSpecifier: &v1.Matcher_Prefix{
					Prefix: "/",
				},
			},
			Action: &v1.Route_DirectResponseAction{
				DirectResponseAction: &v1.DirectResponseAction{
					Status: 200,
					Body:   "test",
				},
			},
			RoutePlugins: &v1.RoutePlugins{
				Extensions: &v1.Extensions{
					Configs: map[string]*types.Struct{
						ExtensionName: jwtRouteSt,
					},
				},
			},
		}

		jwtVhost = &jwt.VhostExtension{
			Jwks: &jwt.Jwks{
				Jwks: &jwt.Jwks_Remote{
					Remote: &jwt.RemoteJwks{
						Url:         "testium",
						UpstreamRef: &core.ResourceRef{Name: "test", Namespace: "testns"},
					},
				},
			},
			Audiences: []string{"testaud"},
			Issuer:    "testiss",
		}

	})
	JustBeforeEach(func() {
		jwtVhostSt, err := util.MessageToStruct(jwtVhost)
		Expect(err).NotTo(HaveOccurred())

		virtualHost = &v1.VirtualHost{
			Name:    "virt1",
			Domains: []string{"*"},
			VirtualHostPlugins: &v1.VirtualHostPlugins{
				Extensions: &v1.Extensions{
					Configs: map[string]*types.Struct{
						ExtensionName: jwtVhostSt,
					},
				},
			},
			Routes: []*v1.Route{route},
		}

		proxy := &v1.Proxy{
			Metadata: core.Metadata{
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
	})

	Context("Process snapshot", func() {
		var (
			outRoute     envoyroute.Route
			outVhost     envoyroute.VirtualHost
			outFilters   []plugins.StagedHttpFilter
			keySetString []byte
			cfg          *envoyauth.JwtAuthentication
		)
		JustBeforeEach(func() {
			outVhost = envoyroute.VirtualHost{
				Name: "test",
			}
			outRoute = envoyroute.Route{}

			// run it like the translator:
			err := plugin.ProcessRoute(routeParams, route, &outRoute)
			Expect(err).NotTo(HaveOccurred())
			err = plugin.ProcessVirtualHost(vhostParams, virtualHost, &outVhost)
			Expect(err).NotTo(HaveOccurred())
			outFilters, err = plugin.HttpFilters(params, nil)
			Expect(err).NotTo(HaveOccurred())

			Expect(outFilters).To(HaveLen(1))
			filter := outFilters[0]
			cfgSt := filter.HttpFilter.GetConfig()
			cfg = &envoyauth.JwtAuthentication{}
			err = util.StructToMessage(cfgSt, cfg)
			Expect(err).NotTo(HaveOccurred())
		})

		BeforeEach(func() {

			// in theory i could have used jwks instead of serializing,
			// but i want to make sure that this will work in future go versions, if
			// the serialization order changes.
			keySetString, _ = json.Marshal(&keySet)
		})

		It("should process virtual host", func() {
			pfc := outVhost.PerFilterConfig[JwtFilterName]
			Expect(pfc).NotTo(BeNil())

			var routeCfg SoloJwtAuthnPerRoute
			err := util.StructToMessage(pfc, &routeCfg)
			Expect(err).NotTo(HaveOccurred())
			Expect(routeCfg.Requirement).To(Equal(virtualHost.Name))
		})

		It("should process route", func() {
			pfc := outRoute.PerFilterConfig[JwtFilterName]
			Expect(pfc).NotTo(BeNil())

			var routeCfg SoloJwtAuthnPerRoute
			err := util.StructToMessage(pfc, &routeCfg)
			Expect(err).NotTo(HaveOccurred())
			Expect(routeCfg.Requirement).To(Equal(DisableName))
		})

		It("should process filters", func() {
			// Expect a requirement
			providerName := ProviderName(virtualHost.Name, "default")
			expectedCfg := envoyauth.JwtAuthentication{
				Providers: map[string]*envoyauth.JwtProvider{
					providerName: {
						Issuer:            jwtVhost.Issuer,
						Audiences:         jwtVhost.Audiences,
						PayloadInMetadata: PayloadInMetadata,
						JwksSourceSpecifier: &envoyauth.JwtProvider_RemoteJwks{
							RemoteJwks: &envoyauth.RemoteJwks{
								HttpUri: &envoycore.HttpUri{
									Uri: jwtVhost.GetJwks().GetRemote().Url,
									HttpUpstreamType: &envoycore.HttpUri_Cluster{
										Cluster: translator.UpstreamToClusterName(*jwtVhost.GetJwks().GetRemote().UpstreamRef),
									},
								},
							},
						},
					},
				},
				FilterStateRules: &envoyauth.FilterStateRule{
					Name: StateName,
					Requires: map[string]*envoyauth.JwtRequirement{
						virtualHost.Name: {
							RequiresType: &envoyauth.JwtRequirement_ProviderName{
								ProviderName: providerName,
							},
						},
					},
				},
			}
			Expect(expectedCfg).To(Equal(*cfg))
		})

		Context("local jwks", func() {
			BeforeEach(func() {
				jwtVhost = &jwt.VhostExtension{
					Jwks: &jwt.Jwks{
						Jwks: &jwt.Jwks_Local{
							Local: &jwt.LocalJwks{
								Key: jwks,
							},
						},
					},
					Audiences: []string{"testaud"},
					Issuer:    "testiss",
				}
			})

			It("should process filters", func() {

				// Expect a requirement
				providerName := ProviderName(virtualHost.Name, "default")
				expectedCfg := envoyauth.JwtAuthentication{
					Providers: map[string]*envoyauth.JwtProvider{
						providerName: {
							Issuer:            jwtVhost.Issuer,
							Audiences:         jwtVhost.Audiences,
							PayloadInMetadata: PayloadInMetadata,
							JwksSourceSpecifier: &envoyauth.JwtProvider_LocalJwks{
								LocalJwks: &envoycore.DataSource{
									Specifier: &envoycore.DataSource_InlineString{
										InlineString: string(keySetString),
									},
								},
							},
						},
					},
					FilterStateRules: &envoyauth.FilterStateRule{
						Name: StateName,
						Requires: map[string]*envoyauth.JwtRequirement{
							virtualHost.Name: {
								RequiresType: &envoyauth.JwtRequirement_ProviderName{
									ProviderName: providerName,
								},
							},
						},
					},
				}
				Expect(expectedCfg).To(Equal(*cfg))
			})
		})

		Context("claims to headers", func() {
			BeforeEach(func() {
				route.RoutePlugins = nil
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

				pfc := outVhost.PerFilterConfig[JwtFilterName]
				Expect(pfc).NotTo(BeNil())

				var routeCfg SoloJwtAuthnPerRoute
				err := util.StructToMessage(pfc, &routeCfg)
				Expect(err).NotTo(HaveOccurred())

				provider1Name := ProviderName(virtualHost.Name, "provider1")
				expectedCfg := SoloJwtAuthnPerRoute{
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
				}
				Expect(expectedCfg).To(Equal(routeCfg))
			})
		})

		Context("claims token source", func() {
			BeforeEach(func() {
				route.RoutePlugins = nil
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
			})

			It("should translate token source", func() {
				provider1Name := ProviderName(virtualHost.Name, "provider1")
				expectedCfg := envoyauth.JwtAuthentication{
					Providers: map[string]*envoyauth.JwtProvider{
						provider1Name: {
							Issuer:            "testiss1",
							Audiences:         []string{"testaud1"},
							PayloadInMetadata: provider1Name,
							FromHeaders: []*envoyauth.JwtHeader{{
								Name:        "header",
								ValuePrefix: "prefix",
							}},
							FromParams: []string{"query"},
							JwksSourceSpecifier: &envoyauth.JwtProvider_LocalJwks{
								LocalJwks: &envoycore.DataSource{
									Specifier: &envoycore.DataSource_InlineString{
										InlineString: string(keySetString),
									},
								},
							},
						},
					},
					FilterStateRules: &envoyauth.FilterStateRule{
						Name: StateName,
						Requires: map[string]*envoyauth.JwtRequirement{
							virtualHost.Name: {
								RequiresType: &envoyauth.JwtRequirement_ProviderName{
									ProviderName: provider1Name,
								},
							},
						},
					},
				}
				Expect(expectedCfg).To(Equal(*cfg))
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
				expectedCfg := envoyauth.JwtAuthentication{
					Providers: map[string]*envoyauth.JwtProvider{
						provider1Name: {
							Issuer:            "testiss1",
							Audiences:         []string{"testaud1"},
							PayloadInMetadata: provider1Name,
							JwksSourceSpecifier: &envoyauth.JwtProvider_LocalJwks{
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
							JwksSourceSpecifier: &envoyauth.JwtProvider_LocalJwks{
								LocalJwks: &envoycore.DataSource{
									Specifier: &envoycore.DataSource_InlineString{
										InlineString: string(keySetString),
									},
								},
							},
						},
					},
					FilterStateRules: &envoyauth.FilterStateRule{
						Name: StateName,
						Requires: map[string]*envoyauth.JwtRequirement{
							virtualHost.Name: {
								RequiresType: &envoyauth.JwtRequirement_RequiresAny{
									RequiresAny: &envoyauth.JwtRequirementOrList{
										Requirements: []*envoyauth.JwtRequirement{
											{
												RequiresType: &envoyauth.JwtRequirement_ProviderName{
													ProviderName: provider1Name,
												},
											}, {
												RequiresType: &envoyauth.JwtRequirement_ProviderName{
													ProviderName: provider2Name,
												},
											},
										},
									},
								},
							},
						},
					},
				}
				Expect(expectedCfg).To(Equal(*cfg))
			})

		})

	})

	Context("translate key", func() {

		It("should translate PEM", func() {
			jwks, err := TranslateKey(publicKey)
			Expect(err).NotTo(HaveOccurred())
			// make certs empty and not nil for comparison
			jwks.Keys[0].Certificates = make([]*x509.Certificate, 0)
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
