package setup

import (
	"testing"

	ggv2utils "github.com/solo-io/gloo/projects/gateway2/utils"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
)

func TestProxyReplayAcrossSetupChanges(t *testing.T) {
	settings := baseSettings("gloo-system")
	r := newProxyCacheRun(t, settings)
	latest := ggv2utils.NewLatest[v1.ProxyList]()
	latest.Publish(v1.ProxyList{revisionProxy("current")})
	for _, change := range []func(){func() {}, func() { settings.DevMode = true }, func() { settings.WatchNamespaces = []string{"gloo-system", "tenant-b"} }} {
		change()
		old := r.client
		r.run()
		r.assertRetained()
		startHandoffConsumer(t, latest, r.client)
		awaitProxyRevision(t, r.client, "current")
		if old != nil {
			stale, err := old.Read("gloo-system", "gateway", clients.ReadOpts{})
			if err != nil {
				t.Fatal(err)
			}
			stale.Metadata.Labels["revision"] = "stale"
			if _, err := old.Write(stale, clients.WriteOpts{OverwriteExisting: true}); err != nil {
				t.Fatal(err)
			}
			awaitProxyRevision(t, r.client, "current")
		}
	}
}
