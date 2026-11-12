package proxy_syncer

import (
	"context"
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	envoy_config_route_v3 "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	"github.com/onsi/gomega"
	"github.com/rotisserie/eris"
	"go.opencensus.io/stats/view"
	"go.uber.org/zap"
	"istio.io/istio/pkg/kube/krt"
	corev1 "k8s.io/api/core/v1"

	"github.com/solo-io/gloo/projects/envoyinit/pkg/runner"
	"github.com/solo-io/gloo/projects/gateway2/krtcollections"
	gloov1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	extauthkubev1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/extauth/v1/kube/apis/enterprise.gloo.solo.io/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/gloosnapshot"
	glookubev1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/kube/apis/gloo.solo.io/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins"
	gloopluginregistry "github.com/solo-io/gloo/projects/gloo/pkg/plugins/registry"
	"github.com/solo-io/gloo/projects/gloo/pkg/syncer/setup"
	"github.com/solo-io/gloo/projects/gloo/pkg/syncer/stats"
	"github.com/solo-io/gloo/projects/gloo/pkg/xds"
	rlkubev1a1 "github.com/solo-io/solo-apis/pkg/api/ratelimit.solo.io/v1alpha1"
	envoycache "github.com/solo-io/solo-kit/pkg/api/v1/control-plane/cache"
	"github.com/solo-io/solo-kit/pkg/api/v1/control-plane/types"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"istio.io/istio/pkg/kube/controllers"
)

// Regression tests for interrupted validation on the kube gateway path.
// An interrupted translation must not reach xDS or proxy status.

const testProxyName = "test-proxy"

func TestBuildXdsSnapshotAbortsOnInterruptedValidation(t *testing.T) {
	g := gomega.NewWithT(t)

	pt := newTestProxyTranslator(interruptedValidationErr())

	snapshot, reports, proxyReport, err := pt.buildXdsSnapshot(
		krt.TestingDummyContext{}, context.Background(), testProxy(), &gloosnapshot.ApiSnapshot{})

	g.Expect(err).To(gomega.MatchError(runner.ErrValidationInterrupted))
	g.Expect(err.Error()).To(gomega.ContainSubstring("envoy config validation was interrupted"))

	g.Expect(snapshot).To(gomega.BeNil())
	g.Expect(reports).To(gomega.BeNil())
	g.Expect(proxyReport).To(gomega.BeNil())
}

// A plugin error that is not an interruption belongs on the reports, so the proxy is rejected as usual.
func TestBuildXdsSnapshotReportsPluginErrors(t *testing.T) {
	g := gomega.NewWithT(t)

	proxy := testProxy()
	pt := newTestProxyTranslator(eris.New("bad config"))

	snapshot, reports, proxyReport, err := pt.buildXdsSnapshot(
		krt.TestingDummyContext{}, context.Background(), proxy, &gloosnapshot.ApiSnapshot{})

	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(snapshot).NotTo(gomega.BeNil())
	g.Expect(proxyReport).NotTo(gomega.BeNil())
	g.Expect(reports.ValidateStrict()).To(gomega.MatchError(gomega.ContainSubstring("bad config")))
}

func TestBuildXdsSnapshotSucceedsWithoutInterruption(t *testing.T) {
	g := gomega.NewWithT(t)

	proxy := testProxy()
	pt := newTestProxyTranslator(nil)

	snapshot, reports, proxyReport, err := pt.buildXdsSnapshot(
		krt.TestingDummyContext{}, context.Background(), proxy, &gloosnapshot.ApiSnapshot{})

	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(snapshot).NotTo(gomega.BeNil())
	g.Expect(proxyReport).NotTo(gomega.BeNil())
	g.Expect(reports.ValidateStrict()).NotTo(gomega.HaveOccurred())
}

// translateProxy has no error return, so a nil wrapper is the only signal that the proxy was skipped.
func TestTranslateProxySkipsInterruptedValidation(t *testing.T) {
	g := gomega.NewWithT(t)

	before := interruptedValidationSkips(g, testProxyResourceName())

	s := &ProxySyncer{proxyTranslator: newTestProxyTranslator(interruptedValidationErr())}
	wrapper := translateTestProxy(s)

	g.Expect(wrapper).To(gomega.BeNil())
	g.Eventually(func() float64 {
		return interruptedValidationSkips(g, testProxyResourceName())
	}).Should(gomega.Equal(before+1),
		"the skip must be counted so a proxy whose xDS updates are stalled is visible")
}

// translateProxy runs inside a krt collection, as it does in the syncer, and its first two
// translations are interrupted. No input changes after the collection starts, so krt runs
// translateProxy again only when a retry fires the proxy's recompute trigger. The subtests cover a
// new proxy (no xDS snapshot yet) and a proxy with a last good snapshot.
func TestTranslateProxyRetriesWithoutInputChanges(t *testing.T) {
	for _, existing := range []bool{false, true} {
		t.Run(fmt.Sprintf("existing=%v", existing), func(t *testing.T) {
			g := gomega.NewWithT(t)
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			var attempts atomic.Int32
			pt := newTestProxyTranslator(nil)
			pt.xdsCache = xds.NewAdsSnapshotCache(ctx)
			key := testProxyResourceName()
			old := &xds.EnvoySnapshot{Clusters: envoycache.Resources{Version: "last-good"}}
			if existing {
				pt.xdsCache.SetSnapshot(key, old)
			}
			// The first two translations are interrupted; the third succeeds.
			pt.translator.PluginRegistry = func(context.Context) plugins.PluginRegistry {
				var err error
				if attempts.Add(1) <= 2 {
					err = interruptedValidationErr()
				}
				return gloopluginregistry.NewPluginRegistry([]plugins.Plugin{&failingPlugin{err: err}})
			}
			s := &ProxySyncer{proxyTranslator: pt}
			proxies := krt.NewStaticCollection[glooProxy](nil, []glooProxy{{Proxy: testProxy()}})
			s.initTranslationRetries(ctx, proxies)
			snapshots := krt.NewCollection(proxies, func(kctx krt.HandlerContext, p glooProxy) *XdsSnapWrapper {
				return translateTestProxyWithHandler(ctx, kctx, s, &p)
			}, krt.WithStop(ctx.Done()))
			snapshots.Register(func(e krt.Event[XdsSnapWrapper]) {
				if e.Event != controllers.EventDelete {
					p := e.Latest()
					s.proxyTranslator.syncXds(ctx, p.snap, p.proxyKey)
				}
			})
			// After the first, interrupted translation: the last good snapshot is kept, or there is none.
			g.Eventually(attempts.Load).Should(gomega.Equal(int32(1)))
			current, err := pt.xdsCache.GetSnapshot(key)
			if existing {
				g.Expect(err).NotTo(gomega.HaveOccurred())
				g.Expect(current.GetResources(types.ClusterTypeV3).Version).To(gomega.Equal("last-good"))
			} else {
				g.Expect(err).To(gomega.HaveOccurred())
			}

			// Retries alone (1s, then 2s) reach the successful third translation, which is published.
			g.Eventually(func() bool {
				current, err := pt.xdsCache.GetSnapshot(key)
				return err == nil && current.GetResources(types.ClusterTypeV3).Version != "last-good"
			}, 8*time.Second).Should(gomega.BeTrue())
			g.Expect(attempts.Load()).To(gomega.Equal(int32(3)))

			// Retries stop after the success.
			g.Consistently(attempts.Load, 1200*time.Millisecond).Should(gomega.Equal(int32(3)))
		})
	}
}

// Every translation is interrupted, so only deleting the proxy or shutting down stops the retries.
func TestTranslateProxyStopsRetryingAfterDeletionOrShutdown(t *testing.T) {
	for _, shutdown := range []bool{false, true} {
		t.Run(fmt.Sprintf("shutdown=%v", shutdown), func(t *testing.T) {
			g := gomega.NewWithT(t)
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			var attempts atomic.Int32
			pt := newTestProxyTranslator(nil)
			pt.translator.PluginRegistry = func(context.Context) plugins.PluginRegistry {
				attempts.Add(1)
				return gloopluginregistry.NewPluginRegistry([]plugins.Plugin{&failingPlugin{err: interruptedValidationErr()}})
			}
			s := &ProxySyncer{proxyTranslator: pt}
			proxies := krt.NewStaticCollection[glooProxy](nil, []glooProxy{{Proxy: testProxy()}})
			s.initTranslationRetries(ctx, proxies)
			snapshots := krt.NewCollection(proxies, func(kctx krt.HandlerContext, p glooProxy) *XdsSnapWrapper {
				return translateTestProxyWithHandler(ctx, kctx, s, &p)
			}, krt.WithStop(ctx.Done()))
			// The first translation is interrupted and schedules a retry for 1s later.
			g.Eventually(snapshots.HasSynced).Should(gomega.BeTrue())
			g.Expect(attempts.Load()).To(gomega.Equal(int32(1)))

			// Shut down or delete the proxy before that retry runs; it must never translate again.
			if shutdown {
				cancel()
			} else {
				proxies.DeleteObjects(func(glooProxy) bool { return true })
			}
			g.Consistently(attempts.Load, 1500*time.Millisecond).Should(gomega.Equal(int32(1)))
		})
	}
}

// delayedCleanupCollection holds each delete event for initTranslationRetries' delete handler until
// release is closed. krt runs that handler on its own goroutine, so it can lag behind translation.
type delayedCleanupCollection struct {
	krt.Collection[glooProxy]
	entered chan struct{}
	release chan struct{}
	done    chan struct{}
	stop    <-chan struct{}
}

func (c delayedCleanupCollection) Register(f func(krt.Event[glooProxy])) krt.HandlerRegistration {
	return c.Collection.Register(func(e krt.Event[glooProxy]) {
		if e.Event == controllers.EventDelete {
			close(c.entered)
			select {
			case <-c.release:
			case <-c.stop:
				return
			}
			f(e)
			close(c.done)
			return
		}
		f(e)
	})
}

// A proxy is deleted and recreated under the same key. The delete handler runs only after the
// replacement has translated, been interrupted, and scheduled a retry. That retry must still
// translate the replacement.
func TestTranslateProxyRetriesAfterRecreationWithDelayedCleanup(t *testing.T) {
	g := gomega.NewWithT(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Only the second translation is interrupted: the first and third succeed.
	var attempts atomic.Int32
	pt := newTestProxyTranslator(nil)
	pt.translator.PluginRegistry = func(context.Context) plugins.PluginRegistry {
		var err error
		if attempts.Add(1) == 2 {
			err = interruptedValidationErr()
		}
		return gloopluginregistry.NewPluginRegistry([]plugins.Plugin{&failingPlugin{err: err}})
	}
	s := &ProxySyncer{proxyTranslator: pt}

	// The delete handler sees proxies through delayed, so it waits for release. Translation reads
	// proxies directly and runs without delay.
	p := glooProxy{Proxy: testProxy()}
	proxies := krt.NewStaticCollection[glooProxy](nil, []glooProxy{p}, krt.WithStop(ctx.Done()))
	delayed := delayedCleanupCollection{
		Collection: proxies,
		entered:    make(chan struct{}),
		release:    make(chan struct{}),
		done:       make(chan struct{}),
		stop:       ctx.Done(),
	}
	s.initTranslationRetries(ctx, delayed)
	translated := make(chan struct{}, 10)
	snapshots := krt.NewCollection(proxies, func(kctx krt.HandlerContext, p glooProxy) *XdsSnapWrapper {
		result := translateTestProxyWithHandler(ctx, kctx, s, &p)
		translated <- struct{}{}
		return result
	}, krt.WithStop(ctx.Done()))

	// The first translation succeeds and publishes a snapshot.
	g.Eventually(translated).Should(gomega.Receive())
	g.Eventually(func() int { return len(snapshots.List()) }).Should(gomega.Equal(1))

	// Delete the proxy. Its snapshot is removed, but the delete handler is held before cleanup.
	proxies.DeleteObjects(func(glooProxy) bool { return true })
	g.Eventually(delayed.entered).Should(gomega.BeClosed())
	g.Eventually(func() int { return len(snapshots.List()) }).Should(gomega.Equal(0))

	// Recreate the proxy. Its translation is interrupted and schedules a retry.
	proxies.Reset([]glooProxy{p})
	g.Eventually(translated).Should(gomega.Receive())
	g.Expect(attempts.Load()).To(gomega.Equal(int32(2)))

	// Let the held delete handler run now that the replacement has scheduled its retry.
	close(delayed.release)
	g.Eventually(delayed.done).Should(gomega.BeClosed())

	// The recreated proxy must retry its interrupted validation without further input changes.
	g.Eventually(attempts.Load, 3*time.Second).Should(gomega.Equal(int32(3)))
	g.Eventually(func() int { return len(snapshots.List()) }).Should(gomega.Equal(1))
}

func TestTranslateProxyDoesNotCountSkipsFromAnEndedContext(t *testing.T) {
	g := gomega.NewWithT(t)

	before := interruptedValidationSkips(g, testProxyResourceName())

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	s := &ProxySyncer{proxyTranslator: newTestProxyTranslator(interruptedValidationErr())}
	wrapper := translateTestProxyWithCtx(ctx, s)

	g.Expect(wrapper).To(gomega.BeNil())
	g.Consistently(func() float64 {
		return interruptedValidationSkips(g, testProxyResourceName())
	}).Should(gomega.Equal(before),
		"a skip caused by the context ending says nothing about this proxy's validation")
}

func TestTranslateProxyEmitsSnapshotWithoutInterruption(t *testing.T) {
	g := gomega.NewWithT(t)

	s := &ProxySyncer{proxyTranslator: newTestProxyTranslator(nil)}
	wrapper := translateTestProxy(s)

	g.Expect(wrapper).NotTo(gomega.BeNil())
	g.Expect(wrapper.proxyKey).To(gomega.Equal(testProxyResourceName()))
}

func testProxyResourceName() string {
	return glooProxy{Proxy: testProxy()}.ResourceName()
}

func translateTestProxy(s *ProxySyncer) *XdsSnapWrapper {
	return translateTestProxyWithCtx(context.Background(), s)
}

func translateTestProxyWithCtx(ctx context.Context, s *ProxySyncer) *XdsSnapWrapper {
	return translateTestProxyWithHandler(ctx, krt.TestingDummyContext{}, s, &glooProxy{Proxy: testProxy()})
}

func translateTestProxyWithHandler(ctx context.Context, kctx krt.HandlerContext, s *ProxySyncer, p *glooProxy) *XdsSnapWrapper {
	return s.translateProxy(
		ctx,
		kctx,
		zap.NewNop().Sugar(),
		p,
		krt.NewStaticCollection[*corev1.ConfigMap](nil, nil),
		krt.NewStaticCollection[EndpointResources](nil, nil),
		krt.NewStaticCollection[RedactedSecret](nil, nil),
		krt.NewStaticCollection[krtcollections.UpstreamWrapper](nil, nil),
		krt.NewStaticCollection[*extauthkubev1.AuthConfig](nil, nil),
		krt.NewStaticCollection[*rlkubev1a1.RateLimitConfig](nil, nil),
	)
}

// newTestProxyTranslator builds a ProxyTranslator whose only plugin fails translation with err. A nil
// err translates cleanly.
func newTestProxyTranslator(err error) ProxyTranslator {
	settings := &glookubev1.Settings{}
	return ProxyTranslator{
		translator: setup.TranslatorFactory{
			PluginRegistry: func(context.Context) plugins.PluginRegistry {
				return gloopluginregistry.NewPluginRegistry([]plugins.Plugin{&failingPlugin{err: err}})
			},
		},
		settings: krt.NewStatic(&settings, true),
	}
}

// interruptedValidationErr is the error a plugin returns when its envoy validation fork was killed.
func interruptedValidationErr() error {
	return eris.Wrap(runner.ErrValidationInterrupted,
		"envoy validation of envoy.filters.http.waf config was interrupted")
}

func testProxy() *gloov1.Proxy {
	return &gloov1.Proxy{
		Metadata: &core.Metadata{Namespace: "gloo-system", Name: testProxyName},
		Listeners: []*gloov1.Listener{{
			Name:        "http",
			BindAddress: "::",
			BindPort:    8080,
			ListenerType: &gloov1.Listener_HttpListener{
				HttpListener: &gloov1.HttpListener{
					VirtualHosts: []*gloov1.VirtualHost{{
						Name:    "vhost",
						Domains: []string{"*"},
					}},
				},
			},
		}},
	}
}

// interruptedValidationSkips reads the InterruptedValidationSkips count for a proxy. The count
// persists across tests, so tests compare against a baseline.
func interruptedValidationSkips(g *gomega.WithT, resourceName string) float64 {
	rows, err := view.RetrieveData(stats.InterruptedValidationSkips.Name())
	g.Expect(err).NotTo(gomega.HaveOccurred())

	for _, row := range rows {
		for _, t := range row.Tags {
			if t.Key == stats.ProxyNameKey && t.Value == resourceName {
				return row.Data.(*view.SumData).Value
			}
		}
	}
	return 0
}

// failingPlugin returns err from every virtual host it processes. A nil err makes it a no-op.
type failingPlugin struct {
	err error
}

func (p *failingPlugin) Name() string { return "failing-plugin" }

func (p *failingPlugin) Init(plugins.InitParams) {}

func (p *failingPlugin) ProcessVirtualHost(
	plugins.VirtualHostParams,
	*gloov1.VirtualHost,
	*envoy_config_route_v3.VirtualHost,
) error {
	return p.err
}

var _ plugins.VirtualHostPlugin = &failingPlugin{}
