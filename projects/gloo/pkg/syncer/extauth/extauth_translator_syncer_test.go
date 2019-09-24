package extauth_test

import (
	"context"

	"github.com/gogo/protobuf/proto"
	"github.com/solo-io/gloo/projects/gloo/pkg/defaults"

	"github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	"github.com/gogo/protobuf/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/solo-io/solo-projects/projects/gloo/pkg/syncer/extauth"

	envoyutil "github.com/envoyproxy/go-control-plane/pkg/util"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	extauth "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/plugins/extauth/v1"
	envoycache "github.com/solo-io/solo-kit/pkg/api/v1/control-plane/cache"
	skcore "github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	extauth2 "github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/extauth"
)

var _ = Describe("ExtauthTranslatorSyncer", func() {
	var (
		proxy      *gloov1.Proxy
		translator *TranslatorSyncerExtension
		secret     *gloov1.Secret
		validAuthConfig,
		customAuthConfig *extauth.AuthConfig
		apiSnapshot *gloov1.ApiSnapshot
		snapCache   *mockSetSnapshot
	)

	JustBeforeEach(func() {
		translator = NewTranslatorSyncerExtension()
		secret = &gloov1.Secret{
			Metadata: skcore.Metadata{
				Name:      "secret",
				Namespace: "gloo-system",
			},

			Kind: &gloov1.Secret_Extension{
				Extension: oidcSecret(),
			},
		}
		apiSnapshot = &gloov1.ApiSnapshot{
			Proxies:     []*gloov1.Proxy{proxy},
			Secrets:     []*gloov1.Secret{secret},
			AuthConfigs: extauth.AuthConfigList{validAuthConfig, customAuthConfig},
		}
		snapCache = &mockSetSnapshot{}
	})

	translate := func() envoycache.Snapshot {
		err := translator.SyncAndSet(context.Background(), apiSnapshot, snapCache)
		Expect(err).NotTo(HaveOccurred())
		Expect(snapCache.Snapshots).To(HaveKey("extauth"))
		return snapCache.Snapshots["extauth"]
	}

	// TODO(marco): just remove this whole block when removing deprecated API
	Context("deprecated config format", func() {
		BeforeEach(func() {
			proxy = getProxy(true, skcore.ResourceRef{})
		})

		It("should work with one listener", func() {
			snap := translate()
			res := snap.GetResources(extauth.ExtAuthConfigType)
			Expect(res.Items).To(HaveLen(1))
		})

		It("should work with two listeners", func() {
			proxy.Listeners = append(proxy.Listeners, &gloov1.Listener{
				Name: "listener-::-8080",
				ListenerType: &gloov1.Listener_HttpListener{
					HttpListener: &gloov1.HttpListener{
						VirtualHosts: []*gloov1.VirtualHost{{
							Name: "gloo-system.default",
						}},
					},
				},
			})

			snap := translate()
			res := snap.GetResources(extauth.ExtAuthConfigType)
			Expect(res.Items).To(HaveLen(1))
		})

		It("generates a single snapshot resource if two listeners use the same auth config", func() {
			newListener := *proxy.Listeners[0]
			newListener.Name = "listener2"
			proxy.Listeners = append(proxy.Listeners, &newListener)

			snap := translate()
			res := snap.GetResources(extauth.ExtAuthConfigType)
			Expect(res.Items).To(HaveLen(1))
		})
	})

	Context("new config format", func() {

		BeforeEach(func() {
			validAuthConfig = &extauth.AuthConfig{
				Metadata: skcore.Metadata{
					Name:      "auth",
					Namespace: defaults.GlooSystem,
				},
				Configs: []*extauth.AuthConfig_Config{{
					AuthConfig: &extauth.AuthConfig_Config_Oauth{
						Oauth: &extauth.OAuth{
							AppUrl:       "https://blah.example.com",
							CallbackPath: "/CallbackPath",
							ClientId:     "oidc.ClientId",
							ClientSecretRef: &skcore.ResourceRef{
								Name:      "secret",
								Namespace: "gloo-system",
							},
							IssuerUrl: "https://issuer.example.com",
						},
					},
				}},
			}

			customAuthConfig = &extauth.AuthConfig{
				Metadata: skcore.Metadata{
					Name:      "custom-auth",
					Namespace: defaults.GlooSystem,
				},
				Configs: []*extauth.AuthConfig_Config{{
					AuthConfig: &extauth.AuthConfig_Config_CustomAuth{CustomAuth: &extauth.CustomAuth{}},
				}},
			}

		})

		Context("config that needs to be translated (non-custom)", func() {

			BeforeEach(func() {
				proxy = getProxy(false, validAuthConfig.Metadata.Ref())
			})

			It("should work with one listener", func() {
				snap := translate()
				res := snap.GetResources(extauth.ExtAuthConfigType)
				Expect(res.Items).To(HaveLen(1))
			})

			It("should work with two listeners", func() {
				proxy.Listeners = append(proxy.Listeners, &gloov1.Listener{
					Name: "listener-::-8080",
					ListenerType: &gloov1.Listener_HttpListener{
						HttpListener: &gloov1.HttpListener{
							VirtualHosts: []*gloov1.VirtualHost{{
								Name: "gloo-system.default",
							}},
						},
					},
				})

				snap := translate()
				res := snap.GetResources(extauth.ExtAuthConfigType)
				Expect(res.Items).To(HaveLen(1))
			})

			It("generates a single snapshot resource if two listeners use the same auth config", func() {
				newListener := *proxy.Listeners[0]
				newListener.Name = "listener2"
				proxy.Listeners = append(proxy.Listeners, &newListener)

				snap := translate()
				res := snap.GetResources(extauth.ExtAuthConfigType)
				Expect(res.Items).To(HaveLen(1))
			})
		})

		Context("config that needs to be translated (non-custom)", func() {

			BeforeEach(func() {
				proxy = getProxy(false, customAuthConfig.Metadata.Ref())
			})

			It("should work with one listener", func() {
				snap := translate()
				res := snap.GetResources(extauth.ExtAuthConfigType)
				Expect(res.Items).To(HaveLen(0))
			})
		})
	})

})

func oidcSecret() *gloov1.Extension {
	input := extauth.OauthSecret{
		ClientSecret: "123",
	}
	secretStruct, err := envoyutil.MessageToStruct(&input)
	Expect(err).NotTo(HaveOccurred())

	return &gloov1.Extension{
		Config: secretStruct,
	}
}

func getProxy(useDeprecatedVersion bool, authConfigRef skcore.ResourceRef) *gloov1.Proxy {
	var vhostAuth proto.Message
	if useDeprecatedVersion {
		vhostAuth = &extauth.VhostExtension{
			AuthConfig: &extauth.VhostExtension_Oauth{
				Oauth: &extauth.OAuth{
					AppUrl:       "https://blah.example.com",
					CallbackPath: "/CallbackPath",
					ClientId:     "oidc.ClientId",
					ClientSecretRef: &skcore.ResourceRef{
						Name:      "secret",
						Namespace: "gloo-system",
					},
					IssuerUrl: "https://issuer.example.com",
				},
			},
		}
	} else {
		vhostAuth = &extauth.ExtAuthExtension{
			Spec: &extauth.ExtAuthExtension_ConfigRef{
				ConfigRef: &authConfigRef,
			},
		}
	}

	vhostAuthStruct, err := envoyutil.MessageToStruct(vhostAuth)
	Expect(err).NotTo(HaveOccurred())
	return &gloov1.Proxy{
		Metadata: skcore.Metadata{
			Name:      "proxy",
			Namespace: "gloo-system",
		},
		Listeners: []*gloov1.Listener{{
			Name: "listener-::-8443",
			ListenerType: &gloov1.Listener_HttpListener{
				HttpListener: &gloov1.HttpListener{
					VirtualHosts: []*gloov1.VirtualHost{{
						Name: "gloo-system.default",
						VirtualHostPlugins: &gloov1.VirtualHostPlugins{
							Extensions: &gloov1.Extensions{
								Configs: map[string]*types.Struct{extauth2.ExtensionName: vhostAuthStruct},
							},
						},
					}},
				},
			},
		}},
	}
}

type nodeHash struct{}

// ID function defines a unique string identifier for the remote Envoy node.
func (nodeHash) ID(node *core.Node) string { return "foo" }

type mockSetSnapshot struct {
	Snapshots map[string]envoycache.Snapshot
}

func (m *mockSetSnapshot) SetSnapshot(node string, snapshot envoycache.Snapshot) error {
	if m.Snapshots == nil {
		m.Snapshots = make(map[string]envoycache.Snapshot)
	}

	m.Snapshots[node] = snapshot
	return nil
}
