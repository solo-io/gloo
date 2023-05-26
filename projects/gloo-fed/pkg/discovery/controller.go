package discovery

import (
	"context"
	"time"

	v1 "github.com/solo-io/external-apis/pkg/api/k8s/core/v1"
	"github.com/solo-io/skv2/pkg/multicluster"
	"github.com/solo-io/skv2/pkg/reconcile"
	fedv1 "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.solo.io/v1"
	"github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.solo.io/v1/input"
	"github.com/solo-io/solo-projects/projects/gloo-fed/pkg/discovery/translator"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

func InitializeDiscovery(
	ctx context.Context,
	writeNamespace string,
	localManager manager.Manager,
	mcClient multicluster.Client,
	cw multicluster.ClusterWatcher,
) {
	glooResourceReconciler := NewGlooResourceReconciler(
		ctx,
		fedv1.NewClientset(localManager.GetClient()).GlooInstances(),
		input.NewMultiClusterBuilder(cw.(multicluster.Interface), mcClient),
		translator.NewTranslator(writeNamespace, v1.NewMulticlusterClientset(mcClient)),
	)

	// RegisterMultiClusterReconciler starts a reconcile loop for each type of watched resource on each remote cluster.
	// When an event occurs for a watched resource, the ReconcileAll function is called to reconcile the state of the world
	// (in this case, creates an updated input snapshot to produce a list of GlooInstances).
	// The verifier allows us to ignore missing CRD errors on remote clusters during their reconcile loops.
	verifier := buildVerifier(ctx)
	input.RegisterMultiClusterReconciler(ctx, cw, glooResourceReconciler.ReconcileAll, time.Second/2, input.ReconcileOptions{
		Gateways:              reconcile.Options{Verifier: verifier},
		MatchableHttpGateways: reconcile.Options{Verifier: verifier},
		MatchableTcpGateways:  reconcile.Options{Verifier: verifier},
		VirtualServices:       reconcile.Options{Verifier: verifier},
		RouteTables:           reconcile.Options{Verifier: verifier},
		Upstreams:             reconcile.Options{Verifier: verifier},
		UpstreamGroups:        reconcile.Options{Verifier: verifier},
		Proxies:               reconcile.Options{Verifier: verifier},
		AuthConfigs:           reconcile.Options{Verifier: verifier},
		RateLimitConfigs:      reconcile.Options{Verifier: verifier},
	})
}
