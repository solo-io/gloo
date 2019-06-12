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
	jose "gopkg.in/square/go-jose.v2"

	"github.com/solo-io/solo-projects/projects/gloo/pkg/api/v1/plugins/jwt"

	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins"
	"github.com/solo-io/gloo/projects/gloo/pkg/translator"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	. "github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/jwt"
)

var _ = Describe("Plugin", func() {
	var (
		plugin      *Plugin
		params      plugins.Params
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
			Jwks: &jwt.VhostExtension_Jwks{
				Jwks: &jwt.VhostExtension_Jwks_Remote{
					Remote: &jwt.VhostExtension_RemoteJwks{
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
	})

	Context("Process snapshot", func() {
		var (
			outRoute   envoyroute.Route
			outVhost   envoyroute.VirtualHost
			outFilters []plugins.StagedHttpFilter
		)
		JustBeforeEach(func() {
			outVhost = envoyroute.VirtualHost{
				Name: "test",
			}
			outRoute = envoyroute.Route{}

			// run it like the translator:
			err := plugin.ProcessRoute(params, route, &outRoute)
			Expect(err).NotTo(HaveOccurred())
			err = plugin.ProcessVirtualHost(params, virtualHost, &outVhost)
			Expect(err).NotTo(HaveOccurred())
			outFilters, err = plugin.HttpFilters(params, nil)
			Expect(err).NotTo(HaveOccurred())
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
			Expect(outFilters).To(HaveLen(1))
			filter := outFilters[0]
			Expect(filter.Stage).To(Equal(plugins.InAuth))
			Expect(filter.HttpFilter.Name).To(Equal(JwtFilterName))

			cfgSt := filter.HttpFilter.GetConfig()
			cfg := envoyauth.JwtAuthentication{}
			err := util.StructToMessage(cfgSt, &cfg)
			Expect(err).NotTo(HaveOccurred())

			// Expect a requirement
			expectedCfg := envoyauth.JwtAuthentication{
				Providers: map[string]*envoyauth.JwtProvider{
					virtualHost.Name: &envoyauth.JwtProvider{
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
						virtualHost.Name: &envoyauth.JwtRequirement{
							RequiresType: &envoyauth.JwtRequirement_ProviderName{
								ProviderName: virtualHost.Name,
							},
						},
					},
				},
			}
			Expect(expectedCfg).To(Equal(cfg))
		})

		Context("local jwks", func() {
			BeforeEach(func() {
				jwtVhost = &jwt.VhostExtension{
					Jwks: &jwt.VhostExtension_Jwks{
						Jwks: &jwt.VhostExtension_Jwks_Local{
							Local: &jwt.VhostExtension_LocalJwks{
								Key: jwks,
							},
						},
					},
					Audiences: []string{"testaud"},
					Issuer:    "testiss",
				}
			})

			It("should process filters", func() {
				Expect(outFilters).To(HaveLen(1))
				filter := outFilters[0]
				Expect(filter.Stage).To(Equal(plugins.InAuth))
				Expect(filter.HttpFilter.Name).To(Equal(JwtFilterName))

				cfgSt := filter.HttpFilter.GetConfig()
				cfg := envoyauth.JwtAuthentication{}
				err := util.StructToMessage(cfgSt, &cfg)
				Expect(err).NotTo(HaveOccurred())

				// in theory i could have used jwks instead of serializing,
				// but i want to make sure that this will work in future go versions, if
				// the serialization order changes.
				keySetString, _ := json.Marshal(&keySet)

				// Expect a requirement
				expectedCfg := envoyauth.JwtAuthentication{
					Providers: map[string]*envoyauth.JwtProvider{
						virtualHost.Name: &envoyauth.JwtProvider{
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
							virtualHost.Name: &envoyauth.JwtRequirement{
								RequiresType: &envoyauth.JwtRequirement_ProviderName{
									ProviderName: virtualHost.Name,
								},
							},
						},
					},
				}
				Expect(expectedCfg).To(Equal(cfg))
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
