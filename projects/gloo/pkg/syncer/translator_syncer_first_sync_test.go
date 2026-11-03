package syncer_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/solo-io/gloo/pkg/bootstrap/leaderelector"
	"github.com/solo-io/gloo/pkg/bootstrap/leaderelector/singlereplica"
	"github.com/solo-io/gloo/pkg/utils/setuputils"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/grpc/validation"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	v1snap "github.com/solo-io/gloo/projects/gloo/pkg/api/v1/gloosnapshot"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins"
	. "github.com/solo-io/gloo/projects/gloo/pkg/syncer"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/kube"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/memory"
	envoycache "github.com/solo-io/solo-kit/pkg/api/v1/control-plane/cache"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-kit/pkg/api/v2/reporter"
)

var _ = Describe("Translator syncer in a setup run", func() {
	It("holds setup restarts until its first Sync", func() {
		ctx, cancel := context.WithCancel(context.Background())
		DeferCleanup(cancel)

		settings := &v1.Settings{Metadata: &core.Metadata{Name: "default", Namespace: "gloo-system"}}
		snapshot := &v1.SetupSnapshot{Settings: v1.SettingsList{settings}}

		translatorSyncers := make(chan v1snap.ApiSyncer, 1)
		probe := &probeTranslator{}
		setupSyncer := setuputils.NewSetupSyncer(
			settings.GetMetadata().Ref(),
			func(
				ctx context.Context,
				kubeCache kube.SharedCache,
				inMemoryCache memory.InMemoryResourceCache,
				settings *v1.Settings,
				identity leaderelector.Identity) error {
				translatorSyncers <- newTestTranslatorSyncer(ctx, probe, &MockXdsCache{})
				return nil
			},
			singlereplica.Identity())

		// The setup event loop cannot restart setup until Sync returns.
		setupDone := make(chan error, 1)
		go func() { setupDone <- setupSyncer.Sync(ctx, snapshot) }()
		var translatorSyncer v1snap.ApiSyncer
		Eventually(translatorSyncers, "2s").Should(Receive(&translatorSyncer))
		Consistently(setupDone, "500ms").ShouldNot(Receive(), "a restart before the first translator Sync must be held")

		var translated, heldDuringTranslate bool
		probe.onTranslate = func() {
			translated = true
			heldDuringTranslate = len(setupDone) == 0
		}
		proxy := &v1.Proxy{Metadata: &core.Metadata{Namespace: cancelTestNamespace, Name: "first-sync-proxy"}}
		Expect(translatorSyncer.Sync(ctx, &v1snap.ApiSnapshot{Proxies: v1.ProxyList{proxy}})).To(Succeed())
		Expect(translated).To(BeTrue())
		Expect(heldDuringTranslate).To(BeTrue(), "a restart while the first Sync translates, before its xDS write, must be held")
		Eventually(setupDone, "2s").Should(Receive(BeNil()))
	})
})

// probeTranslator runs onTranslate while the envoy syncer translates a proxy, before the xDS write.
type probeTranslator struct {
	onTranslate func()
}

func (t *probeTranslator) Translate(params plugins.Params, proxy *v1.Proxy) (envoycache.Snapshot, reporter.ResourceReports, *validation.ProxyReport) {
	t.onTranslate()
	return envoycache.NilSnapshot{}, reporter.ResourceReports{}, &validation.ProxyReport{}
}
