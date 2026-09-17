package proxy_syncer

import (
	"context"
	"testing"

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
	rlkubev1a1 "github.com/solo-io/solo-apis/pkg/api/ratelimit.solo.io/v1alpha1"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
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
	return s.translateProxy(
		ctx,
		krt.TestingDummyContext{},
		zap.NewNop().Sugar(),
		&glooProxy{Proxy: testProxy()},
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
