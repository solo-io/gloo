package reconcile

import (
	"context"

	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-kit/pkg/errors"
	"github.com/solo-io/solo-kit/pkg/utils/contextutils"
)

type ListResources func() ([]resources.Resource, error)

type Reconciler interface {
	Reconcile(opts clients.ListOpts, kind string, desiredResources []resources.Resource) error
}

type reconciler struct {
	resourceClients map[string]clients.ResourceClient
}

func NewReconciler(resourceClients ...clients.ResourceClient) Reconciler {
	mapped := make(map[string]clients.ResourceClient)
	for _, rc := range resourceClients {
		mapped[rc.Kind()] = rc
	}
	return &reconciler{
		resourceClients: mapped,
	}
}

func (r *reconciler) resourceClient(kind string) (clients.ResourceClient, error) {
	rc, ok := r.resourceClients[kind]
	if !ok {
		return nil, errors.Errorf("no resource client registered for kind %s", kind)
	}
	return rc, nil
}

func (r *reconciler) Reconcile(opts clients.ListOpts, kind string, desiredResources []resources.Resource) error {
	opts = opts.WithDefaults()
	opts.Ctx = contextutils.WithLogger(opts.Ctx, "reconciler")
	rc, err := r.resourceClient(kind)
	if err != nil {
		return err
	}
	originalResources, err := rc.List(opts)
	if err != nil {
		return err
	}
	for _, desired := range desiredResources {
		if err := syncResource(opts.Ctx, rc, desired, originalResources); err != nil {
			return errors.Wrapf(err, "reconciling resource %v", desired.GetMetadata().Name)
		}
	}
	// delete unused
	for _, original := range originalResources {
		unused := findResource(original.GetMetadata().Namespace, original.GetMetadata().Name, desiredResources) == nil
		if unused {
			if err := deleteStaleResource(opts.Ctx, rc, original); err != nil {
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
		desired.SetStatus(core.Status{})
	}
	_, err := rc.Write(desired, clients.WriteOpts{Ctx: ctx, OverwriteExisting: overwriteExisting})
	return err
}

func deleteStaleResource(ctx context.Context, rc clients.ResourceClient, original resources.Resource) error {
	return rc.Delete(original.GetMetadata().Name, clients.DeleteOpts{
		Ctx:            ctx,
		Namespace:      original.GetMetadata().Namespace,
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
