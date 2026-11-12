package syncer_test

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	gatewayv1 "github.com/solo-io/gloo/projects/gateway/pkg/api/v1"
	ratelimit "github.com/solo-io/gloo/projects/gloo/pkg/api/external/solo/ratelimit"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	extauthv1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/extauth/v1"
	v1snap "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/gloosnapshot"
	"github.com/solo-io/gloo/projects/gloo/pkg/xds"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/factory"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/memory"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
)

// The retry only asks the API emitter for another snapshot, so it depends on the emitter and the event
// loop delivering that snapshot to Sync even though no resource changed.
var _ = Describe("Translator syncer retries through the API event loop", func() {
	It("publishes an interrupted proxy without a resource change", func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		emit := make(chan struct{})

		proxy := &v1.Proxy{Metadata: &core.Metadata{Namespace: cancelTestNamespace, Name: "retry-proxy"}}
		emitter := newMemoryApiEmitter(ctx, emit, proxy)

		xdsCache := xds.NewAdsSnapshotCache(ctx)
		syncer := newTestTranslatorSyncerWithGateway(ctx, &recoveringTranslator{remaining: 2}, xdsCache, nil, nil, emit)
		// Run the translator syncer inside ApiSyncers, as setup does.
		_, err := v1snap.NewApiEventLoop(emitter, v1snap.ApiSyncers{syncer}).
			Run([]string{cancelTestNamespace}, clients.WatchOpts{Ctx: ctx})
		Expect(err).NotTo(HaveOccurred())

		// The first two translations are interrupted, so the third Sync publishes the snapshot.
		Eventually(func() error {
			_, err := xdsCache.GetSnapshot(xds.SnapshotCacheKey(proxy))
			return err
		}, 8*time.Second).Should(Succeed())
	})
})

// newMemoryApiEmitter returns an API emitter on in-memory clients that holds proxy and sends a snapshot
// on each receive from emit.
func newMemoryApiEmitter(ctx context.Context, emit <-chan struct{}, proxy *v1.Proxy) v1snap.ApiEmitter {
	f := &factory.MemoryResourceClientFactory{Cache: memory.NewInMemoryResourceCache()}
	must := func(err error) {
		GinkgoHelper()
		Expect(err).NotTo(HaveOccurred())
	}
	artifacts, err := v1.NewArtifactClient(ctx, f)
	must(err)
	endpoints, err := v1.NewEndpointClient(ctx, f)
	must(err)
	proxies, err := v1.NewProxyClient(ctx, f)
	must(err)
	upstreamGroups, err := v1.NewUpstreamGroupClient(ctx, f)
	must(err)
	secrets, err := v1.NewSecretClient(ctx, f)
	must(err)
	upstreams, err := v1.NewUpstreamClient(ctx, f)
	must(err)
	authConfigs, err := extauthv1.NewAuthConfigClient(ctx, f)
	must(err)
	rateLimitConfigs, err := ratelimit.NewRateLimitConfigClient(ctx, f)
	must(err)
	virtualServices, err := gatewayv1.NewVirtualServiceClient(ctx, f)
	must(err)
	routeTables, err := gatewayv1.NewRouteTableClient(ctx, f)
	must(err)
	gateways, err := gatewayv1.NewGatewayClient(ctx, f)
	must(err)
	virtualHostOptions, err := gatewayv1.NewVirtualHostOptionClient(ctx, f)
	must(err)
	routeOptions, err := gatewayv1.NewRouteOptionClient(ctx, f)
	must(err)
	httpGateways, err := gatewayv1.NewMatchableHttpGatewayClient(ctx, f)
	must(err)
	tcpGateways, err := gatewayv1.NewMatchableTcpGatewayClient(ctx, f)
	must(err)

	_, err = proxies.Write(proxy, clients.WriteOpts{Ctx: ctx})
	must(err)

	return v1snap.NewApiEmitterWithEmit(artifacts, endpoints, proxies, upstreamGroups, secrets, upstreams,
		authConfigs, rateLimitConfigs, virtualServices, routeTables, gateways, virtualHostOptions, routeOptions,
		httpGateways, tcpGateways, emit)
}
