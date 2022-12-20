package discovery

import (
	"context"

	"github.com/solo-io/skv2/pkg/ezkube"
	fedv1 "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.solo.io/v1"
	"github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.solo.io/v1/input"
	hub_v1sets "github.com/solo-io/solo-projects/projects/gloo-fed/pkg/api/fed.solo.io/v1/sets"
	"github.com/solo-io/solo-projects/projects/gloo-fed/pkg/discovery/translator"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	SnapshotName = "discovery"
)

type GlooResourceReconciler interface {
	ReconcileAll(id ezkube.ClusterResourceId) (bool, error)
}

type glooResourceReconcilerImpl struct {
	ctx                context.Context
	glooInstanceClient fedv1.GlooInstanceClient
	fedInputSnapshot   input.Builder
	translator         translator.GlooInstanceTranslator
}

func (g *glooResourceReconcilerImpl) ReconcileAll(
	_ ezkube.ClusterResourceId,
) (bool, error) { // The first bool indicates whether we should retry
	// Build a snapshot that gets everything across all clusters.
	// Translate that snapshot into a list of Gloo Instances.

	discoverySnapshot, err := g.fedInputSnapshot.BuildSnapshot(g.ctx, SnapshotName, input.BuildOptions{})
	if err != nil {
		return true, err
	}
	instances := g.translator.FromSnapshot(g.ctx, discoverySnapshot)
	newSet := hub_v1sets.NewGlooInstanceSet(instances...)

	existingInstances, err := g.glooInstanceClient.ListGlooInstance(g.ctx)
	if err != nil {
		return true, err
	}

	existingSet := hub_v1sets.NewGlooInstanceSet()
	for _, existingIter := range existingInstances.Items {
		existing := existingIter
		existingSet.Insert(&existing)
	}

	for _, instance := range newSet.List() {
		if err = g.glooInstanceClient.UpsertGlooInstance(g.ctx, instance); IgnoreAlreadyExists(err) != nil {
			return true, err
		}
	}

	for _, deletedIter := range existingSet.Difference(newSet).List() {
		deleted := deletedIter
		if err = g.glooInstanceClient.DeleteGlooInstance(g.ctx, client.ObjectKey{
			Namespace: deleted.GetNamespace(),
			Name:      deleted.GetName(),
		}); err != nil {
			return true, err
		}
	}
	return true, nil
}

// Visible for testing
func NewGlooResourceReconciler(
	ctx context.Context,
	glooInstanceClient fedv1.GlooInstanceClient,
	snapshotBuilder input.Builder,
	translator translator.GlooInstanceTranslator,
) GlooResourceReconciler {
	return &glooResourceReconcilerImpl{
		ctx:                ctx,
		glooInstanceClient: glooInstanceClient,
		fedInputSnapshot:   snapshotBuilder,
		translator:         translator,
	}
}
