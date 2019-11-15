package extauth_test

import (
	"context"

	"github.com/solo-io/gloo/projects/gloo/pkg/defaults"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/solo-io/solo-projects/projects/gloo/pkg/syncer/extauth"
	. "github.com/solo-io/solo-projects/test/extauth/helpers"

	envoyutil "github.com/envoyproxy/go-control-plane/pkg/util"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	extauth "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/plugins/extauth/v1"
	envoycache "github.com/solo-io/solo-kit/pkg/api/v1/control-plane/cache"
	skcore "github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
)

var _ = Describe("ExtauthTranslatorSyncer", func() {
	var (
		proxy           *gloov1.Proxy
		translator      *TranslatorSyncerExtension
		secret          *gloov1.Secret
		oauthAuthConfig *extauth.AuthConfig
		apiSnapshot     *gloov1.ApiSnapshot
		snapCache       *mockSetSnapshot
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
			AuthConfigs: extauth.AuthConfigList{oauthAuthConfig},
		}
		snapCache = &mockSetSnapshot{}
	})

	translate := func() envoycache.Snapshot {
		err := translator.SyncAndSet(context.Background(), apiSnapshot, snapCache)
		Expect(err).NotTo(HaveOccurred())
		Expect(snapCache.Snapshots).To(HaveKey("extauth"))
		return snapCache.Snapshots["extauth"]
	}

	// TODO(kdorosh) remove outer context right before merge -- leave around for PR review for easy diff
	Context("strongly typed config", func() {

		BeforeEach(func() {
			oauthAuthConfig = &extauth.AuthConfig{
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

		})

		Context("config that needs to be translated (non-custom)", func() {

			BeforeEach(func() {
				proxy = getProxy(StronglyTyped, oauthAuthConfig.Metadata.Ref())
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

func getProxy(configFormat ConfigFormatType, authConfigRef skcore.ResourceRef) *gloov1.Proxy {
	proxy := &gloov1.Proxy{
		Metadata: skcore.Metadata{
			Name:      "proxy",
			Namespace: "gloo-system",
		},
		Listeners: []*gloov1.Listener{{
			Name: "listener-::-8443",
			ListenerType: &gloov1.Listener_HttpListener{
				HttpListener: &gloov1.HttpListener{
					VirtualHosts: []*gloov1.VirtualHost{{
						Name:               "gloo-system.default",
						VirtualHostPlugins: nil,
					}},
				},
			},
		}},
	}

	var plugins *gloov1.VirtualHostPlugins
	if configFormat == StronglyTyped {
		plugins = &gloov1.VirtualHostPlugins{
			Extauth: &extauth.ExtAuthExtension{
				Spec: &extauth.ExtAuthExtension_ConfigRef{
					ConfigRef: &authConfigRef,
				},
			},
		}

	}

	proxy.Listeners[0].GetHttpListener().VirtualHosts[0].VirtualHostPlugins = plugins

	return proxy
}

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
