package reconciler

import (
	"context"

	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/reconcile"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-kit/pkg/errors"
	"github.com/solo-io/solo-kit/pkg/utils/contextutils"
	"github.com/solo-io/solo-kit/projects/gloo/pkg/api/v1"
)
// Option to copy anything from the original to the desired before writing
type TransitionUpstreamsFunc func(original, desired *v1.Upstream)

type UpstreamReconciler interface {
	Reconcile(namespace string, desiredResources []*v1.Upstream, opts clients.ListOpts) error
}

func upstreamsToResources(list ... *v1.Upstream) []resources.Resource {
	var resourceList []resources.Resource
	for _, upstream := range list {
		resourceList = append(resourceList, upstream)
	}
	return resourceList
}

func resourcesToUpstreams(list ... resources.Resource) []*v1.Upstream {
	var upstreamList []*v1.Upstream
	for _, resource := range list {
		upstreamList = append(upstreamList, resource.(*v1.Upstream))
	}
	return upstreamList
}

func NewUpstreamReconciler(client v1.UpstreamClient) {
	edsReconciler := reconcile.NewReconciler(endpointsResourceClient)
	edsReconciler.Reconcile()
}

type ListResources func() ([]*v1.Upstream, error)

type Reconciler interface {
	Reconcile(namespace string, opts clients.ListOpts, kind string, desiredResources []*v1.Upstream) error
}

type reconciler struct {
	resourceClient v1.UpstreamClient
}

func NewReconciler(resourceClient v1.UpstreamClient) Reconciler {
	return &reconciler{
		resourceClient: resourceClient,
	}
}

func (r *reconciler) Reconcile(namespace string, opts clients.ListOpts, kind string, desiredResources []*v1.Upstream) error {
	opts = opts.WithDefaults()
	opts.Ctx = contextutils.WithLogger(opts.Ctx, "reconciler")
	originalResources, err := r.resourceClient.List(namespace, opts)
	if err != nil {
		return err
	}
	for _, desired := range desiredResources {
		if err := r.syncResource(opts.Ctx, desired, originalResources); err != nil {
			return errors.Wrapf(err, "reconciling resource %v", desired.GetMetadata().Name)
		}
	}
	// delete unused
	for _, original := range originalResources {
		unused := findResource(original.GetMetadata().Namespace, original.GetMetadata().Name, desiredResources) == nil
		if unused {
			if err := r.deleteStaleResource(opts.Ctx, original); err != nil {
				return errors.Wrapf(err, "deleting stale resource %v", original.GetMetadata().Name)
			}
		}
	}

	return nil
}

func (r *reconciler) syncResource(ctx context.Context, desired *v1.Upstream, originalResources []*v1.Upstream) error {
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
	_, err := r.resourceClient.Write(desired, clients.WriteOpts{Ctx: ctx, OverwriteExisting: overwriteExisting})
	return err
}

func (r *reconciler) deleteStaleResource(ctx context.Context, original *v1.Upstream) error {
	return r.resourceClient.Delete(original.GetMetadata().Namespace, original.GetMetadata().Name, clients.DeleteOpts{
		Ctx:            ctx,
		IgnoreNotExist: true,
	})
}

func findResource(namespace, name string, rss []*v1.Upstream) *v1.Upstream {
	for _, resource := range rss {
		if resource.GetMetadata().Namespace == namespace && resource.GetMetadata().Name == name {
			return resource
		}
	}
	return nil
}
