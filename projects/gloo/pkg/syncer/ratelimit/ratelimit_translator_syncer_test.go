package ratelimit_test

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise"
	"github.com/solo-io/gloo/projects/gloo/pkg/syncer"

	rlsyncer "github.com/solo-io/solo-projects/projects/gloo/pkg/syncer/ratelimit"

	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	envoycache "github.com/solo-io/solo-kit/pkg/api/v1/control-plane/cache"
	skcore "github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
)

var _ = Describe("RateLimitTranslatorSyncer", func() {
	var (
		proxy       *gloov1.Proxy
		params      syncer.TranslatorSyncerExtensionParams
		translator  syncer.TranslatorSyncerExtension
		apiSnapshot *gloov1.ApiSnapshot
		snapCache   *mockSetSnapshot
	)

	BeforeEach(func() {
		params = syncer.TranslatorSyncerExtensionParams{}
	})

	JustBeforeEach(func() {
		var err error
		translator, err = rlsyncer.NewTranslatorSyncerExtension(context.TODO(), params)
		Expect(err).ToNot(HaveOccurred())
		apiSnapshot = &gloov1.ApiSnapshot{
			Proxies: []*gloov1.Proxy{proxy},
		}
		snapCache = &mockSetSnapshot{}
	})

	translate := func() envoycache.Snapshot {
		err := translator.Sync(context.Background(), apiSnapshot, snapCache)
		Expect(err).NotTo(HaveOccurred())
		Expect(snapCache.Snapshots).To(HaveKey("ratelimit"))
		return snapCache.Snapshots["ratelimit"]
	}

	Context("config that needs to be translated (non-custom)", func() {

		BeforeEach(func() {
			proxy = getProxy()
		})

		It("should work with one listener", func() {
			snap := translate()
			res := snap.GetResources(enterprise.RateLimitConfigType)
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
			res := snap.GetResources(enterprise.RateLimitConfigType)
			Expect(res.Items).To(HaveLen(1))
		})

	})

})

func getProxy() *gloov1.Proxy {
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
					}},
				},
			},
		}},
	}
}

type mockSetSnapshot struct {
	Snapshots map[string]envoycache.Snapshot
}

func (m *mockSetSnapshot) CreateWatch(envoycache.Request) (value chan envoycache.Response, cancel func()) {
	panic("implement me")
}

func (m *mockSetSnapshot) Fetch(context.Context, envoycache.Request) (*envoycache.Response, error) {
	panic("implement me")
}

func (m *mockSetSnapshot) GetStatusInfo(string) envoycache.StatusInfo {
	panic("implement me")
}

func (m *mockSetSnapshot) GetStatusKeys() []string {
	panic("implement me")
}

func (m *mockSetSnapshot) GetSnapshot(node string) (envoycache.Snapshot, error) {
	panic("implement me")
}

func (m *mockSetSnapshot) ClearSnapshot(node string) {
	panic("implement me")
}

func (m *mockSetSnapshot) SetSnapshot(node string, snapshot envoycache.Snapshot) error {
	if m.Snapshots == nil {
		m.Snapshots = make(map[string]envoycache.Snapshot)
	}

	m.Snapshots[node] = snapshot
	return nil
}
