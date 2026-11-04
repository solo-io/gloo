package setup

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/bootstrap"
	"github.com/solo-io/gloo/projects/gloo/pkg/defaults"
	"github.com/solo-io/gloo/projects/gloo/pkg/utils"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/memory"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-kit/pkg/utils/statusutils"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

// A watchNamespaceSelectors change re-runs setup in-process with a new context but the same
// process-level cache. The gateway reconciler only preserves a rejected listener when it can read
// the Proxy the previous run wrote, so the Proxy client must survive the re-run.
var _ = Describe("Proxy client across setup runs", func() {
	var (
		memCache         memory.InMemoryResourceCache
		edgeLabels       = map[string]string{utils.ProxyTypeKey: utils.GlooEdgeProxyValue}
		gatewayApiLabels = map[string]string{utils.ProxyTypeKey: utils.GatewayApiProxyValue}
	)

	BeforeEach(func() {
		memCache = memory.NewInMemoryResourceCache()
	})

	// newRunProxyClient builds a setup run's Proxy client. A nil settings has no config source, so every
	// client is in-memory and no kube is needed.
	newRunProxyClient := func(settings *v1.Settings, writeNamespace string) v1.ProxyClient {
		ctx, cancel := context.WithCancel(context.Background())
		DeferCleanup(cancel)
		if settings == nil {
			settings = &v1.Settings{}
		}
		opts, err := constructOpts(ctx, &bootstrap.SetupOpts{}, constructOptsParams{
			memCache:       memCache,
			settings:       settings,
			writeNamespace: writeNamespace,
		})
		Expect(err).NotTo(HaveOccurred())
		client, err := v1.NewProxyClient(ctx, opts.Proxies)
		Expect(err).NotTo(HaveOccurred())
		return client
	}
	proxy := func(name string, labels map[string]string) *v1.Proxy {
		return &v1.Proxy{
			Metadata: &core.Metadata{Name: name, Namespace: defaults.GlooSystem, Labels: labels},
		}
	}
	writeFirstRunProxies := func() {
		run1 := newRunProxyClient(nil, defaults.GlooSystem)
		for _, p := range []*v1.Proxy{
			proxy("proxy-a", edgeLabels),
			proxy("proxy-b", edgeLabels),
			proxy("gateway-api-proxy", gatewayApiLabels),
		} {
			_, err := run1.Write(p, clients.WriteOpts{})
			Expect(err).NotTo(HaveOccurred())
		}
	}
	listAll := func(client v1.ProxyClient) []string {
		proxies, err := client.List("", clients.ListOpts{})
		Expect(err).NotTo(HaveOccurred())
		return proxies.Names()
	}

	It("keeps in-memory Proxies written by a previous run until the next reconcile removes them", func() {
		writeFirstRunProxies()

		run2 := newRunProxyClient(nil, defaults.GlooSystem)
		Expect(listAll(run2)).To(ConsistOf("proxy-a", "proxy-b", "gateway-api-proxy"))

		// run2's Edge translation produces only proxy-a: the reconcile deletes run1's proxy-b
		// and leaves the Proxy owned by the Gateway API controller alone.
		err := v1.NewProxyReconciler(run2, statusutils.NewNoOpStatusClient()).Reconcile(
			defaults.GlooSystem,
			v1.ProxyList{proxy("proxy-a", edgeLabels)},
			nil,
			clients.ListOpts{Selector: edgeLabels},
		)
		Expect(err).NotTo(HaveOccurred())
		Expect(listAll(run2)).To(ConsistOf("proxy-a", "gateway-api-proxy"))
	})

	It("deletes the previous run's Proxies when the write namespace changes", func() {
		writeFirstRunProxies()

		run2 := newRunProxyClient(nil, "new-write-namespace")
		Expect(listAll(run2)).To(BeEmpty(), "no reconcile lists Proxies outside the write namespace")
	})

	It("deletes the previous run's Proxies when persistProxySpec is set without a config source", func() {
		writeFirstRunProxies()

		run2 := newRunProxyClient(&v1.Settings{
			Gateway: &v1.GatewayOptions{PersistProxySpec: &wrapperspb.BoolValue{Value: true}},
		}, "new-write-namespace")
		Expect(listAll(run2)).To(BeEmpty(), "without a config source the Proxy store is still in memory")
	})

	It("deletes the previous run's Edge Proxies when the gateway controller is disabled", func() {
		writeFirstRunProxies()

		run2 := newRunProxyClient(&v1.Settings{
			Gateway: &v1.GatewayOptions{EnableGatewayController: &wrapperspb.BoolValue{Value: false}},
		}, defaults.GlooSystem)
		Expect(listAll(run2)).To(ConsistOf("gateway-api-proxy"), "the Gateway API controller still reconciles its Proxy")
	})
})
