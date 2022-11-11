package discovery

import (
	"context"
	"time"

	v1 "github.com/solo-io/external-apis/pkg/api/k8s/core/v1"
	"github.com/solo-io/skv2/pkg/multicluster"
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

	input.RegisterMultiClusterReconciler(ctx, cw, glooResourceReconciler.ReconcileAll, time.Second/2, input.ReconcileOptions{})
}
