package syncer_test

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/solo-io/gloo/pkg/utils/statsutils/metrics"
	"github.com/solo-io/gloo/projects/gloo/pkg/servers/iosnapshot"

	"github.com/solo-io/gloo/pkg/bootstrap/leaderelector/singlereplica"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/solo-io/gloo/pkg/utils/statusutils"
	"github.com/solo-io/gloo/projects/envoyinit/pkg/runner"
	gwreconciler "github.com/solo-io/gloo/projects/gateway/pkg/reconciler"
	gwsyncer "github.com/solo-io/gloo/projects/gateway/pkg/syncer"
	gwtranslator "github.com/solo-io/gloo/projects/gateway/pkg/translator"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/grpc/validation"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	v1snap "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/gloosnapshot"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins"
	. "github.com/solo-io/gloo/projects/gloo/pkg/syncer"
	glootranslator "github.com/solo-io/gloo/projects/gloo/pkg/translator"
	"github.com/solo-io/gloo/projects/gloo/pkg/xds"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/factory"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/memory"
	envoycache "github.com/solo-io/solo-kit/pkg/api/v1/control-plane/cache"
	"github.com/solo-io/solo-kit/pkg/api/v1/control-plane/types"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-kit/pkg/api/v2/reporter"
)

// Regression tests for interrupted validation
//
// A sync whose context ended, or whose envoy validation was interrupted, must not write to the xDS
// cache, through SetSnapshot or through garbage collection.
var _ = Describe("Translator syncer with a cancelled context", func() {

	var (
		ctx    context.Context
		cancel context.CancelFunc

		proxy *v1.Proxy
		snap  *v1snap.ApiSnapshot
	)

	newSyncer := func(translator glootranslator.Translator, xdsCache envoycache.SnapshotCache, emit ...chan<- struct{}) v1snap.ApiSyncer {
		return newTestTranslatorSyncer(ctx, translator, xdsCache, emit...)
	}

	BeforeEach(func() {
		ctx, cancel = context.WithCancel(context.Background())
		proxy = &v1.Proxy{
			Metadata: &core.Metadata{Namespace: cancelTestNamespace, Name: "cancel-proxy"},
		}
		snap = &v1snap.ApiSnapshot{Proxies: v1.ProxyList{proxy}}
	})

	AfterEach(func() {
		cancel()
	})

	It("does not set xDS snapshots for a translation the context ended during", func() {
		xdsCache := &MockXdsCache{}
		// The cancellation lands mid-translation.
		syncer := newSyncer(&cancellingTranslator{cancel: cancel}, xdsCache)

		_ = syncer.Sync(ctx, snap)

		Expect(xdsCache.Called).To(BeFalse(), "a translation the context ended during must not reach the xDS cache")
	})

	It("does not garbage collect xDS snapshots when the context has already ended", func() {
		xdsCache := &staleKeyXdsCache{}
		syncer := newSyncer(&cancellingTranslator{}, xdsCache)

		cancel()
		err := syncer.Sync(ctx, snap)

		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("context ended before or during proxy translation"))
		Expect(xdsCache.Called).To(BeFalse(), "a dying sync must not garbage collect the xDS cache")
	})

	It("sets xDS snapshots when the context is live (sanity check of this setup)", func() {
		xdsCache := &MockXdsCache{}
		syncer := newSyncer(&cancellingTranslator{}, xdsCache)

		err := syncer.Sync(ctx, snap)

		Expect(err).NotTo(HaveOccurred())
		Expect(xdsCache.Called).To(BeTrue())
	})

	DescribeTable("retries interrupted translation without changing resources until it succeeds", func(existing bool) {
		emit := make(chan struct{})
		tx := &recoveringTranslator{remaining: 2}
		xdsCache := xds.NewAdsSnapshotCache(ctx)
		key := xds.SnapshotCacheKey(proxy)
		old := &xds.EnvoySnapshot{Clusters: envoycache.Resources{Version: "last-good"}}
		if existing {
			xdsCache.SetSnapshot(key, old)
		}
		syncer := newSyncer(tx, xdsCache, emit)
		Expect(syncer.Sync(ctx, snap)).To(Succeed())
		current, err := xdsCache.GetSnapshot(key)
		if existing {
			Expect(err).NotTo(HaveOccurred())
			Expect(current.GetResources(types.ClusterTypeV3).Version).To(Equal("last-good"))
		} else {
			Expect(err).To(HaveOccurred())
		}
		for range 2 {
			select {
			case <-emit:
				// Like the force-emitter path, use a fresh copy of the latest inputs.
				latest := snap.Clone()
				Expect(syncer.Sync(ctx, &latest)).To(Succeed())
			case <-time.After(5 * time.Second):
				Fail("interrupted proxy was never retried")
			}
		}
		current, err = xdsCache.GetSnapshot(key)
		Expect(err).NotTo(HaveOccurred())
		Expect(current.GetResources(types.ClusterTypeV3).Version).NotTo(Equal("last-good"))
		Consistently(emit, 1200*time.Millisecond).ShouldNot(Receive())
	}, Entry("new proxy", false), Entry("existing snapshot is preserved until recovery", true))

	It("cancels a pending retry when setup shuts down", func() {
		emit := make(chan struct{}, 1)
		syncer := newSyncer(&interruptingTranslator{}, &MockXdsCache{}, emit)
		Expect(syncer.Sync(ctx, snap)).To(Succeed())
		cancel()
		Consistently(emit, 1200*time.Millisecond).ShouldNot(Receive())
	})

	// A fork killed from outside leaves the context live, so the ctx.Err() checks above do not apply.
	It("does not set xDS snapshots for a translation whose envoy validation was interrupted, even with a live context", func() {
		xdsCache := &MockXdsCache{}
		syncer := newSyncer(&interruptingTranslator{}, xdsCache)

		err := syncer.Sync(ctx, snap)

		Expect(err).NotTo(HaveOccurred())
		Expect(ctx.Err()).To(Succeed(), "this spec is only meaningful while the context is live")
		Expect(xdsCache.Called).To(BeFalse(),
			"a translation that could not validate its config must not reach the xDS cache")
	})
})

const cancelTestNamespace = "cancel-ns"

// newTestTranslatorSyncer builds a translator syncer on in-memory clients.
func newTestTranslatorSyncer(ctx context.Context, translator glootranslator.Translator, xdsCache envoycache.SnapshotCache, emit ...chan<- struct{}) v1snap.ApiSyncer {
	return newTestTranslatorSyncerWithGateway(ctx, translator, xdsCache, nil, nil, emit...)
}

// testGateway is the gateway translation a test translator syncer runs before its own.
type testGateway struct {
	translator gwtranslator.Translator
	reconciler gwreconciler.ProxyReconciler
}

// newTestTranslatorSyncerWithGateway builds a translator syncer on in-memory clients. With gw set, it
// runs gateway translation first and lists Proxies from its in-memory store. With emit set, the
// syncer sends its retries on it.
func newTestTranslatorSyncerWithGateway(
	ctx context.Context,
	translator glootranslator.Translator,
	xdsCache envoycache.SnapshotCache,
	gw *testGateway,
	extensions []TranslatorSyncerExtension,
	emit ...chan<- struct{},
) v1snap.ApiSyncer {
	resourceClientFactory := &factory.MemoryResourceClientFactory{
		Cache: memory.NewInMemoryResourceCache(),
	}
	proxyClient, err := v1.NewProxyClient(ctx, resourceClientFactory)
	Expect(err).NotTo(HaveOccurred())
	upstreamClient, err := resourceClientFactory.NewResourceClient(ctx, factory.NewResourceClientParams{ResourceType: &v1.Upstream{}})
	Expect(err).NotTo(HaveOccurred())

	settings := &v1.Settings{}
	statusClient := statusutils.GetStatusClientFromEnvOrDefault(cancelTestNamespace)
	statusMetrics, err := metrics.NewConfigStatusMetrics(metrics.GetDefaultConfigStatusOptions())
	Expect(err).NotTo(HaveOccurred())
	rep := reporter.NewReporter("cancel-test", statusClient, proxyClient.BaseClient(), upstreamClient)
	history := iosnapshot.GetHistoryFactory()(iosnapshot.HistoryFactoryParameters{
		Settings: settings,
		Cache:    xdsCache,
	})

	var gatewaySyncer *gwsyncer.TranslatorSyncer
	if gw != nil {
		gatewaySyncer = gwsyncer.NewTranslatorSyncer(ctx, cancelTestNamespace, proxyClient, gw.reconciler, rep,
			gw.translator, statusClient, statusMetrics, singlereplica.Identity())
	}

	return NewTranslatorSyncer(ctx, translator, xdsCache, &MockXdsSanitizer{}, rep, false, extensions, settings,
		statusMetrics, gatewaySyncer, proxyClient, cancelTestNamespace, singlereplica.Identity(), history, emit...)
}

// interruptedReason is the validator's wrap of runner.ErrValidationInterrupted as a string. The real
// translator keeps interruptions off reports. The stubs below put this text on theirs to prove the
// report content is not what protects the xDS cache.
const interruptedReason = "envoy validation of envoy.filters.http.waf config was interrupted: " +
	"envoy validation process was interrupted before it completed (ctxErr=context canceled): " +
	"command \"envoy --mode validate\" failed with error: signal: killed"

// cancellingTranslator cancels the sync context mid-translation and errors the proxy report with the
// interruption text. Without a cancel func it translates cleanly.
type cancellingTranslator struct {
	cancel context.CancelFunc
}

func (t *cancellingTranslator) Translate(params plugins.Params, proxy *v1.Proxy) (envoycache.Snapshot, reporter.ResourceReports, *validation.ProxyReport) {
	rpts := reporter.ResourceReports{}
	if t.cancel != nil {
		t.cancel()
		rpts.AddError(proxy, errors.New(interruptedReason))
	}
	return envoycache.NilSnapshot{}, rpts, &validation.ProxyReport{}
}

// interruptingTranslator simulates an envoy validation fork killed from outside the process. It
// records the interruption on params like reporting.go does, with a live context. Unlike the real
// translator it also errors the proxy report, to prove the syncer reads the params signal and not
// the report.
type interruptingTranslator struct{}

func (t *interruptingTranslator) Translate(params plugins.Params, proxy *v1.Proxy) (envoycache.Snapshot, reporter.ResourceReports, *validation.ProxyReport) {
	params.ValidationInterruptions.Add(fmt.Errorf("envoy validation of waf config was interrupted: %w", runner.ErrValidationInterrupted))

	rpts := reporter.ResourceReports{}
	rpts.AddError(proxy, errors.New(interruptedReason))

	return envoycache.NilSnapshot{}, rpts, &validation.ProxyReport{}
}

// staleKeyXdsCache reports a status key absent from the API snapshot, so it is eligible for garbage
// collection.
type staleKeyXdsCache struct {
	MockXdsCache
}

func (c *staleKeyXdsCache) GetStatusKeys() []string {
	return []string{"stale-node"}
}

// Fail twice with unchanged inputs, then produce a usable snapshot.
type recoveringTranslator struct{ remaining int }

func (t *recoveringTranslator) Translate(params plugins.Params, proxy *v1.Proxy) (envoycache.Snapshot, reporter.ResourceReports, *validation.ProxyReport) {
	if t.remaining > 0 {
		t.remaining--
		params.ValidationInterruptions.Add(runner.ErrValidationInterrupted)
	}
	return envoycache.NilSnapshot{}, reporter.ResourceReports{}, &validation.ProxyReport{}
}
