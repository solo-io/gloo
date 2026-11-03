package syncer_test

import (
	"context"
	"errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	gatewayv1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	gwreconciler "github.com/solo-io/gloo/projects/gateway/pkg/reconciler"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	v1snap "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/gloosnapshot"
	. "github.com/solo-io/gloo/projects/gloo/pkg/syncer"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	envoycache "github.com/solo-io/solo-kit/pkg/api/v1/control-plane/cache"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-kit/pkg/api/v2/reporter"
)

// The store below starts empty and the stub reconciler writes nothing, so a reconcile that fails
// lists no Proxies.
var _ = Describe("Translator syncer after a failed proxy reconcile", func() {
	var (
		ctx       context.Context
		cancel    context.CancelFunc
		xdsCache  *keyRecordingXdsCache
		extension *recordingExtension
		reconcile *stubProxyReconciler
		snap      *v1snap.ApiSnapshot
	)

	BeforeEach(func() {
		ctx, cancel = context.WithCancel(context.Background())
		xdsCache = &keyRecordingXdsCache{statusKeys: []string{"connected-node"}}
		extension = &recordingExtension{}
		reconcile = &stubProxyReconciler{}
		snap = &v1snap.ApiSnapshot{Gateways: gatewayv1.GatewayList{{
			Metadata: &core.Metadata{Namespace: cancelTestNamespace, Name: "gw"},
		}}}
	})

	AfterEach(func() {
		cancel()
	})

	sync := func() error {
		syncer := newTestTranslatorSyncerWithGateway(ctx, &cancellingTranslator{}, xdsCache,
			&testGateway{translator: &stubGatewayTranslator{}, reconciler: reconcile},
			[]TranslatorSyncerExtension{extension})
		return syncer.Sync(ctx, snap)
	}

	It("does not garbage collect xDS snapshots or run extension syncs", func() {
		reconcile.err = errors.New("validation process was killed")

		err := sync()

		Expect(err).To(MatchError(ContainSubstring("translating proxies")))
		Expect(xdsCache.setKeys).NotTo(HaveKey("connected-node"), "a node missing from an incomplete Proxy list must keep its snapshot")
		Expect(extension.called).To(BeFalse(), "extension config is built from the Proxy list and must not be rebuilt from an incomplete one")
	})

	It("garbage collects xDS snapshots and runs extension syncs when the reconcile succeeds (sanity check of this setup)", func() {
		Expect(sync()).To(Succeed())

		Expect(xdsCache.setKeys).To(HaveKey("connected-node"))
		Expect(extension.called).To(BeTrue())
	})
})

// stubGatewayTranslator produces an empty Proxy for each proxy name.
type stubGatewayTranslator struct{}

func (t *stubGatewayTranslator) Translate(_ context.Context, proxyName string, _ *v1snap.ApiSnapshot, _ gatewayv1.GatewayList) (*v1.Proxy, reporter.ResourceReports) {
	return &v1.Proxy{Metadata: &core.Metadata{Namespace: cancelTestNamespace, Name: proxyName}}, reporter.ResourceReports{}
}

// stubProxyReconciler writes nothing and returns err.
type stubProxyReconciler struct {
	err error
}

func (r *stubProxyReconciler) ReconcileProxies(context.Context, gwreconciler.GeneratedProxies, string, clients.ListOpts) error {
	return r.err
}

type recordingExtension struct {
	called bool
}

func (e *recordingExtension) ID() string {
	return "test-extension"
}

func (e *recordingExtension) Sync(context.Context, *v1snap.ApiSnapshot, *v1.Settings, SnapshotSetter, reporter.ResourceReports) {
	e.called = true
}

// keyRecordingXdsCache reports statusKeys as connected nodes and records each key SetSnapshot is called with.
type keyRecordingXdsCache struct {
	MockXdsCache
	statusKeys []string
	setKeys    map[string]envoycache.Snapshot
}

func (c *keyRecordingXdsCache) GetStatusKeys() []string {
	return c.statusKeys
}

func (c *keyRecordingXdsCache) SetSnapshot(node string, snapshot envoycache.Snapshot) {
	if c.setKeys == nil {
		c.setKeys = map[string]envoycache.Snapshot{}
	}
	c.setKeys[node] = snapshot
}
