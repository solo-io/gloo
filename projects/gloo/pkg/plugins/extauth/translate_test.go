package extauth_test

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/solo-io/solo-projects/projects/gloo/pkg/api/v1/plugins/extauth"
	. "github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/extauth"

	"github.com/envoyproxy/go-control-plane/pkg/util"
	"github.com/gogo/protobuf/types"
	"github.com/solo-io/gloo/pkg/utils"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	static_plugin_gloo "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/plugins/static"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
)

var _ = Describe("Translate", func() {
	var (
		params       plugins.Params
		virtualHost  *v1.VirtualHost
		upstream     *v1.Upstream
		secret       *v1.Secret
		route        *v1.Route
		extAuthVhost *extauth.VhostExtension
		clientSecret *extauth.OauthSecret
		apiKeySecret *extauth.ApiKeySecret
	)

	BeforeEach(func() {

		upstream = &v1.Upstream{
			Metadata: core.Metadata{
				Name:      "extauth",
				Namespace: "default",
			},
			UpstreamSpec: &v1.UpstreamSpec{
				UpstreamType: &v1.UpstreamSpec_Static{
					Static: &static_plugin_gloo.UpstreamSpec{
						Hosts: []*static_plugin_gloo.Host{{
							Addr: "test",
							Port: 1234,
						}},
					},
				},
			},
		}
		route = &v1.Route{
			Matcher: &v1.Matcher{
				PathSpecifier: &v1.Matcher_Prefix{
					Prefix: "/",
				},
			},
			Action: &v1.Route_RouteAction{
				RouteAction: &v1.RouteAction{
					Destination: &v1.RouteAction_Single{
						Single: &v1.Destination{
							DestinationType: &v1.Destination_Upstream{
								Upstream: utils.ResourceRefPtr(upstream.Metadata.Ref()),
							},
						},
					},
				},
			},
		}

		apiKeySecret = &extauth.ApiKeySecret{
			ApiKey: "apiKey1",
		}

		clientSecret = &extauth.OauthSecret{
			ClientSecret: "1234",
		}

		st, err := util.MessageToStruct(clientSecret)
		Expect(err).NotTo(HaveOccurred())

		secret = &v1.Secret{
			Metadata: core.Metadata{
				Name:      "secret",
				Namespace: "default",
			},
			Kind: &v1.Secret_Extension{
				Extension: &v1.Extension{
					Config: st,
				},
			},
		}
		secretRef := secret.Metadata.Ref()
		extAuthVhost = &extauth.VhostExtension{
			AuthConfig: &extauth.VhostExtension_Oauth{
				Oauth: &extauth.OAuth{
					ClientSecretRef: &secretRef,
					ClientId:        "ClientId",
					IssuerUrl:       "IssuerUrl",
					AppUrl:          "AppUrl",
					CallbackPath:    "CallbackPath",
				},
			},
		}

		params.Snapshot = &v1.ApiSnapshot{
			Upstreams: v1.UpstreamList{upstream},
		}
	})
	JustBeforeEach(func() {

		extAuthSt, err := util.MessageToStruct(extAuthVhost)
		Expect(err).NotTo(HaveOccurred())

		virtualHost = &v1.VirtualHost{
			Name:    "virt1",
			Domains: []string{"*"},
			VirtualHostPlugins: &v1.VirtualHostPlugins{
				Extensions: &v1.Extensions{
					Configs: map[string]*types.Struct{
						ExtensionName: extAuthSt,
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

		params.Snapshot.Proxies = v1.ProxyList{proxy}
		params.Snapshot.Secrets = v1.SecretList{secret}
	})

	It("should translate oauth config for extauth server", func() {
		cfg, err := TranslateUserConfigToExtAuthServerConfig(context.TODO(), params.Snapshot.Proxies[0], params.Snapshot.Proxies[0].Listeners[0], virtualHost, params.Snapshot, *extAuthVhost)
		Expect(err).NotTo(HaveOccurred())
		Expect(cfg.Vhost).To(Equal(GetResourceName(params.Snapshot.Proxies[0], params.Snapshot.Proxies[0].Listeners[0], virtualHost)))
		authCfg := cfg.AuthConfig.(*extauth.ExtAuthConfig_Oauth).Oauth
		expectAuthCfg := extAuthVhost.AuthConfig.(*extauth.VhostExtension_Oauth).Oauth
		Expect(authCfg.IssuerUrl).To(Equal(expectAuthCfg.IssuerUrl))
		Expect(authCfg.ClientId).To(Equal(expectAuthCfg.ClientId))
		Expect(authCfg.ClientSecret).To(Equal(clientSecret.ClientSecret))
		Expect(authCfg.AppUrl).To(Equal(expectAuthCfg.AppUrl))
		Expect(authCfg.CallbackPath).To(Equal(expectAuthCfg.CallbackPath))
	})

	Context("with api key extauth", func() {
		BeforeEach(func() {
			st, err := util.MessageToStruct(apiKeySecret)
			Expect(err).NotTo(HaveOccurred())

			secret = &v1.Secret{
				Metadata: core.Metadata{
					Name:      "secretName",
					Namespace: "default",
					Labels:    map[string]string{"team": "infrastructure"},
				},
				Kind: &v1.Secret_Extension{
					Extension: &v1.Extension{
						Config: st,
					},
				},
			}
			secretRef := secret.Metadata.Ref()

			extAuthVhost = &extauth.VhostExtension{
				AuthConfig: &extauth.VhostExtension_ApiKeyAuth{
					ApiKeyAuth: &extauth.ApiKeyAuth{
						ApiKeySecretRefs: []*core.ResourceRef{&secretRef},
					},
				},
			}
		})

		Context("with api key extauth, secret ref matching", func() {
			It("should translate api keys config for extauth server - matching secret ref", func() {
				cfg, err := TranslateUserConfigToExtAuthServerConfig(context.TODO(), params.Snapshot.Proxies[0], params.Snapshot.Proxies[0].Listeners[0], virtualHost, params.Snapshot, *extAuthVhost)
				Expect(err).NotTo(HaveOccurred())
				Expect(cfg.Vhost).To(Equal(GetResourceName(params.Snapshot.Proxies[0], params.Snapshot.Proxies[0].Listeners[0], virtualHost)))
				authCfg := cfg.AuthConfig.(*extauth.ExtAuthConfig_ApiKeyAuth).ApiKeyAuth
				Expect(authCfg.ValidApiKeyAndUser).To(Equal(map[string]string{"apiKey1": "secretName"}))
			})

			It("should translate api keys config for extauth server - mismatching secret ref", func() {
				secret.Metadata.Name = "mismatchName"
				_, err := TranslateUserConfigToExtAuthServerConfig(context.TODO(), params.Snapshot.Proxies[0], params.Snapshot.Proxies[0].Listeners[0], virtualHost, params.Snapshot, *extAuthVhost)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("list did not find secret"))
			})
		})

		Context("with api key ext auth, label matching", func() {
			BeforeEach(func() {
				extAuthVhost = &extauth.VhostExtension{
					AuthConfig: &extauth.VhostExtension_ApiKeyAuth{
						ApiKeyAuth: &extauth.ApiKeyAuth{
							LabelSelector: map[string]string{"team": "infrastructure"},
						},
					},
				}
			})

			It("should translate api keys config for extauth server - matching label", func() {
				cfg, err := TranslateUserConfigToExtAuthServerConfig(context.TODO(), params.Snapshot.Proxies[0], params.Snapshot.Proxies[0].Listeners[0], virtualHost, params.Snapshot, *extAuthVhost)
				Expect(err).NotTo(HaveOccurred())
				Expect(cfg.Vhost).To(Equal(GetResourceName(params.Snapshot.Proxies[0], params.Snapshot.Proxies[0].Listeners[0], virtualHost)))
				authCfg := cfg.AuthConfig.(*extauth.ExtAuthConfig_ApiKeyAuth).ApiKeyAuth
				Expect(authCfg.ValidApiKeyAndUser).To(Equal(map[string]string{"apiKey1": "secretName"}))
			})

			It("should translate api keys config for extauth server - mismatched labels", func() {
				secret.Metadata.Labels = map[string]string{"missingLabel": "missingValue"}
				_, err := TranslateUserConfigToExtAuthServerConfig(context.TODO(), params.Snapshot.Proxies[0], params.Snapshot.Proxies[0].Listeners[0], virtualHost, params.Snapshot, *extAuthVhost)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal(NoMatchesForGroupError(map[string]string{"team": "infrastructure"}).Error()))
			})

			It("should translate api keys config for extauth server - mismatched labels", func() {
				secret.Metadata.Labels = map[string]string{}
				_, err := TranslateUserConfigToExtAuthServerConfig(context.TODO(), params.Snapshot.Proxies[0], params.Snapshot.Proxies[0].Listeners[0], virtualHost, params.Snapshot, *extAuthVhost)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal(NoMatchesForGroupError(map[string]string{"team": "infrastructure"}).Error()))
			})

			It("should translate api keys config for extauth server - mismatched labels", func() {
				secret.Metadata.Labels = nil
				_, err := TranslateUserConfigToExtAuthServerConfig(context.TODO(), params.Snapshot.Proxies[0], params.Snapshot.Proxies[0].Listeners[0], virtualHost, params.Snapshot, *extAuthVhost)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal(NoMatchesForGroupError(map[string]string{"team": "infrastructure"}).Error()))
			})
		})
	})
	Context("with OPA extauth", func() {
		BeforeEach(func() {
			extAuthVhost = &extauth.VhostExtension{
				Configs: []*extauth.AuthConfig{{
					AuthConfig: &extauth.AuthConfig_OpaAuth{
						OpaAuth: &extauth.OpaAuth{
							Modules: []*core.ResourceRef{{Namespace: "namespace", Name: "name"}},
							Query:   "true",
						},
					},
				}},
			}

			params.Snapshot.Artifacts = v1.ArtifactList{
				{

					Metadata: core.Metadata{
						Name:      "name",
						Namespace: "namespace",
					},
					Data: map[string]string{"module.rego": "package foo"},
				},
			}
		})

		It("should OPA config", func() {
			cfg, err := TranslateUserConfigToExtAuthServerConfig(context.TODO(), params.Snapshot.Proxies[0], params.Snapshot.Proxies[0].Listeners[0], virtualHost, params.Snapshot, *extAuthVhost)
			Expect(err).NotTo(HaveOccurred())
			Expect(cfg.Vhost).To(Equal(GetResourceName(params.Snapshot.Proxies[0], params.Snapshot.Proxies[0].Listeners[0], virtualHost)))
			authCfg := cfg.Configs[0].GetOpaAuth()
			expectAuthCfg := extAuthVhost.Configs[0].GetOpaAuth()
			Expect(authCfg.Query).To(Equal(expectAuthCfg.Query))
			data := params.Snapshot.Artifacts[0].Data
			Expect(authCfg.Modules).To(Equal(data))
		})

	})
})
