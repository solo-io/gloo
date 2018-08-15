package v1

import (
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/reconcile"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/solo-io/solo-kit/pkg/utils/contextutils"
)

// Option to copy anything from the original to the desired before writing
type TransitionArtifactFunc func(original, desired *Artifact)

type ArtifactReconciler interface {
	Reconcile(namespace string, desiredResources []*Artifact, transition TransitionArtifactFunc, opts clients.ListOpts) error
}

func artifactsToResources(list ArtifactList) []resources.Resource {
	var resourceList []resources.Resource
	for _, artifact := range list {
		resourceList = append(resourceList, artifact)
	}
	return resourceList
}

func NewArtifactReconciler(client ArtifactClient) ArtifactReconciler {
	return &artifactReconciler{
		base: reconcile.NewReconciler(client.BaseClient()),
	}
}

type artifactReconciler struct {
	base reconcile.Reconciler
}

func (r *artifactReconciler) Reconcile(namespace string, desiredResources []*Artifact, transition TransitionArtifactFunc, opts clients.ListOpts) error {
	opts = opts.WithDefaults()
	opts.Ctx = contextutils.WithLogger(opts.Ctx, "artifact_reconciler")
	var transitionResources reconcile.TransitionResourcesFunc
	if transition != nil {
		transitionResources = func(original, desired resources.Resource) {
			transition(original.(*Artifact), desired.(*Artifact))
		}
	}
	return r.base.Reconcile(namespace, artifactsToResources(desiredResources), transitionResources, opts)
}
