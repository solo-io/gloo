package extauth_test

import (
	"context"

	"github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	"github.com/gogo/protobuf/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/solo-io/solo-projects/projects/gloo/pkg/syncer/extauth"

	envoyutil "github.com/envoyproxy/go-control-plane/pkg/util"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	envoycache "github.com/solo-io/solo-kit/pkg/api/v1/control-plane/cache"
	skcore "github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/api/v1/plugins/extauth"
	extauth2 "github.com/solo-io/solo-projects/projects/gloo/pkg/plugins/extauth"
)

var _ = Describe("ExtauthTranslatorSyncer", func() {
	var (
		proxy       *gloov1.Proxy
		translator  *ExtAuthTranslatorSyncerExtension
		secret      *gloov1.Secret
		apiSnapshot *gloov1.ApiSnapshot
		snapcache   *mockSetSnapshot
	)
	BeforeEach(func() {
		proxy = getProxy()
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
		apiSnapshot = &gloov1.ApiSnapshot{Proxies: []*gloov1.Proxy{proxy},
			Secrets: []*gloov1.Secret{secret}}
		snapcache = &mockSetSnapshot{}
	})

	translate := func() envoycache.Snapshot {
		translator.SyncAndSet(context.Background(), apiSnapshot, snapcache)
		Expect(snapcache.Snapshots).To(HaveKey("extauth"))
		return snapcache.Snapshots["extauth"]
	}

	It("should work with one listeners", func() {
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

	It("should work with two listeners with auth config", func() {
		newListener := *proxy.Listeners[0]
		newListener.Name = "listener2"
		proxy.Listeners = append(proxy.Listeners, &newListener)

		snap := translate()
		res := snap.GetResources(extauth.ExtAuthConfigType)
		Expect(res.Items).To(HaveLen(2))
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

func getProxy() *gloov1.Proxy {

	vhostAuth := &extauth.VhostExtension{
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
