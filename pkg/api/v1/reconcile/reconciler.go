package reconcile

import (
	"context"

	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-kit/pkg/errors"
	"github.com/solo-io/solo-kit/pkg/utils/contextutils"
)

// Option to copy anything from the original to the desired before writing
type TransitionResourcesFunc func(original, desired resources.Resource)

type Reconciler interface {
	Reconcile(namespace string, desiredResources []resources.Resource, opts clients.ListOpts) error
}

type reconciler struct {
	rc         clients.ResourceClient
	transition TransitionResourcesFunc
}

func NewReconciler(resourceClient clients.ResourceClient, transitionFunc TransitionResourcesFunc) Reconciler {
	return &reconciler{
		rc:         resourceClient,
		transition: transitionFunc,
	}
}

func (r *reconciler) Reconcile(namespace string, desiredResources []resources.Resource, opts clients.ListOpts) error {
	opts = opts.WithDefaults()
	opts.Ctx = contextutils.WithLogger(opts.Ctx, "reconciler")
	originalResources, err := r.rc.List(namespace, opts)
	if err != nil {
		return err
	}
	for _, desired := range desiredResources {
		if err := syncResource(opts.Ctx, r.rc, desired, originalResources); err != nil {
			return errors.Wrapf(err, "reconciling resource %v", desired.GetMetadata().Name)
		}
	}
	// delete unused
	for _, original := range originalResources {
		unused := findResource(original.GetMetadata().Namespace, original.GetMetadata().Name, desiredResources) == nil
		if unused {
			if err := deleteStaleResource(opts.Ctx, r.rc, original); err != nil {
				return errors.Wrapf(err, "deleting stale resource %v", original.GetMetadata().Name)
			}
		}
	}

	return nil
}

func syncResource(ctx context.Context, rc clients.ResourceClient, desired resources.Resource, originalResources []resources.Resource) error {
	var overwriteExisting bool
	original := findResource(desired.GetMetadata().Namespace, desired.GetMetadata().Name, originalResources)
	if original != nil {
		// if this is an update,
		// update resource version
		// set status to 0, needs to be re-processed
		overwriteExisting = true
		resources.UpdateMetadata(desired, func(meta *core.Metadata) {
			meta.ResourceVersion = original.GetMetadata().ResourceVersion
		})
		if desiredInput, ok := desired.(resources.InputResource); ok {
			desiredInput.SetStatus(core.Status{})
		}
	}
	_, err := rc.Write(desired, clients.WriteOpts{Ctx: ctx, OverwriteExisting: overwriteExisting})
	return err
}

func deleteStaleResource(ctx context.Context, rc clients.ResourceClient, original resources.Resource) error {
	return rc.Delete(original.GetMetadata().Namespace, original.GetMetadata().Name, clients.DeleteOpts{
		Ctx:            ctx,
		IgnoreNotExist: true,
	})
}

func findResource(namespace, name string, rss []resources.Resource) resources.Resource {
	for _, resource := range rss {
		if resource.GetMetadata().Namespace == namespace && resource.GetMetadata().Name == name {
			return resource
		}
	}
	return nil
}
