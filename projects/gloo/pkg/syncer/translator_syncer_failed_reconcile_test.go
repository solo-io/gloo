package syncer_test

import (
	"context"
	"errors"
	"fmt"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/solo-io/gloo/projects/envoyinit/pkg/runner"
	gatewayv1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	gwreconciler "github.com/solo-io/gloo/projects/gateway/pkg/reconciler"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	v1snap "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/gloosnapshot"
	. "github.com/solo-io/gloo/projects/gloo/pkg/syncer"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	envoycache "github.com/solo-io/solo-kit/pkg/api/v1/control-plane/cache"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-kit/pkg/api/v2/reporter"
	skerrors "github.com/solo-io/solo-kit/pkg/errors"
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
		emit      chan struct{}
	)

	BeforeEach(func() {
		ctx, cancel = context.WithCancel(context.Background())
		emit = make(chan struct{}, 1)
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
			[]TranslatorSyncerExtension{extension}, emit)
		return syncer.Sync(ctx, snap)
	}

	It("does not garbage collect xDS snapshots, but still runs extension syncs", func() {
		reconcile.err = errors.New("validation process was killed")

		err := sync()

		Expect(err).To(MatchError(ContainSubstring("translating proxies")))
		Expect(xdsCache.setKeys).NotTo(HaveKey("connected-node"), "a node missing from an incomplete Proxy list must keep its snapshot")
		Expect(extension.called).To(BeTrue(), "ext-auth config does not depend on the Proxy list, so AuthConfig and Secret changes must still be published")
	})

	It("requests another sync without a resource change when envoy validation was interrupted", func() {
		// Wrapped the way the gloo validator and the gateway proxy reconciler wrap it.
		reconcile.err = skerrors.Wrapf(
			fmt.Errorf("cannot validate proxy: envoy config validation was interrupted: %w", runner.ErrValidationInterrupted),
			"failed to add proxy validation results to reports")

		Expect(sync()).To(MatchError(ContainSubstring("translating proxies")))

		Eventually(emit, 3*time.Second).Should(Receive())
	})

	It("does not request another sync for other reconcile errors", func() {
		reconcile.err = errors.New("proxy write failed")

		Expect(sync()).To(MatchError(ContainSubstring("translating proxies")))

		Consistently(emit, 1500*time.Millisecond).ShouldNot(Receive())
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
