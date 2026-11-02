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
)

// A watchNamespaceSelectors change re-runs setup in-process with a new context but the same
// process-level cache. The gateway reconciler only preserves a rejected listener when it can read
// the Proxy the previous run wrote, so the Proxy client must survive the re-run.
var _ = Describe("Proxy client across setup runs", func() {

	It("keeps in-memory Proxies written by a previous run until the next reconcile removes them", func() {
		memCache := memory.NewInMemoryResourceCache()
		settings := &v1.Settings{} // no config source: every client is in-memory, no kube needed
		edgeLabels := map[string]string{utils.ProxyTypeKey: utils.GlooEdgeProxyValue}
		gatewayApiLabels := map[string]string{utils.ProxyTypeKey: utils.GatewayApiProxyValue}

		newRunProxyClient := func() v1.ProxyClient {
			ctx, cancel := context.WithCancel(context.Background())
			DeferCleanup(cancel)
			opts, err := constructOpts(ctx, &bootstrap.SetupOpts{}, constructOptsParams{
				memCache:       memCache,
				settings:       settings,
				writeNamespace: defaults.GlooSystem,
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

		run1 := newRunProxyClient()
		for _, p := range []*v1.Proxy{
			proxy("proxy-a", edgeLabels),
			proxy("proxy-b", edgeLabels),
			proxy("gateway-api-proxy", gatewayApiLabels),
		} {
			_, err := run1.Write(p, clients.WriteOpts{})
			Expect(err).NotTo(HaveOccurred())
		}

		run2 := newRunProxyClient()
		proxies, err := run2.List(defaults.GlooSystem, clients.ListOpts{})
		Expect(err).NotTo(HaveOccurred())
		Expect(proxies.Names()).To(ConsistOf("proxy-a", "proxy-b", "gateway-api-proxy"))

		// run2's Edge translation produces only proxy-a: the reconcile deletes run1's proxy-b
		// and leaves the Proxy owned by the Gateway API controller alone.
		err = v1.NewProxyReconciler(run2, statusutils.NewNoOpStatusClient()).Reconcile(
			defaults.GlooSystem,
			v1.ProxyList{proxy("proxy-a", edgeLabels)},
			nil,
			clients.ListOpts{Selector: edgeLabels},
		)
		Expect(err).NotTo(HaveOccurred())
		proxies, err = run2.List(defaults.GlooSystem, clients.ListOpts{})
		Expect(err).NotTo(HaveOccurred())
		Expect(proxies.Names()).To(ConsistOf("proxy-a", "gateway-api-proxy"))
	})
})
