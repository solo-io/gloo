package propagator

import (
	"fmt"

	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-kit/pkg/errors"
	"github.com/solo-io/solo-kit/pkg/utils/errutils"
)

type ResourcesByType map[string]ResourceList

type Propagator struct {
	ResourceClients clients.ResourceClients
}

// sources can be multiple types
func (p *Propagator) PropagateStatus(destinations, sources resources.InputResourceList, writeErrs chan error, opts clients.WatchOpts) error {
	// each ressource by kind
	sourcesByClient := make(map[clients.ResourceClient]resources.InputResourceList)
	for _, r := range sources {
		client, err := p.ResourceClients.ForResource(r)
		if err != nil {
			return err
		}
		sourcesByClient[client] = append(sourcesByClient[client], r)
	}
	destinationsByClient := make(map[clients.ResourceClient]resources.InputResourceList)
	for _, r := range destinations {
		client, err := p.ResourceClients.ForResource(r)
		if err != nil {
			return err
		}
		destinationsByClient[client] = append(destinationsByClient[client], r)
	}

	for sourceRc, sourcesToWatch := range sourcesByClient {
		watch, errs, err := sourceRc.Watch("TODO", opts)
		if err != nil {
			return err
		}
		go errutils.AggregateErrs(opts.Ctx, writeErrs, errs)
		go func() {
			for {
				select {
				case resourceList := <-watch:
					resourceList = resourceList.FilterByList(sourcesToWatch)
				case <-opts.Ctx.Done():
					return
				}
			}
		}()
	}
}

func updateDestStatusIfNeeded(dest, source resources.InputResource, destRc clients.ResourceClient, opts clients.WatchOpts) error {
	destinationRes, err := destRc.Read(dest.GetMetadata().Namespace, dest.GetMetadata().Name, clients.ReadOpts{Ctx: opts.Ctx})
	if err != nil {
		return errors.Wrapf(err, "dependent resource %v no longer found", dest.GetMetadata().Name)
	}
	dest, ok := destinationRes.(resources.InputResource)
	if !ok {
		return errors.Errorf("internal error: bad type assertion")
	}

	updatedSourceRes, err := resourceList.Find(source.GetMetadata().Namespace, source.GetMetadata().Name)
	if err != nil {
		return errors.Wrapf(err, "could not find %v which is a dependency of %v", source.GetMetadata().Name,
			dest.GetMetadata().Name)
	}
	src, ok := updatedSourceRes.(resources.InputResource)
	if !ok {
		return errors.Errorf("internal error: bad type assertion")
	}

	updated := propagateStatus(src.GetMetadata().Name, src.GetStatus(), dest.GetStatus())
	if updated.Equal(dest.GetStatus()) {
		return nil
	}
	dest.SetStatus(updated)

	if _, err := destRc.Write(dest, clients.WriteOpts{
		Ctx:               opts.Ctx,
		OverwriteExisting: true,
	}); err != nil {
		return errors.Wrapf(err, "failed to update %v with new status", dest.GetMetadata().Name)
	}
	return nil
}

func propagateStatus(dependencyName string, primary, secondary core.Status) core.Status {
	switch primary.State {
	case core.Status_Rejected:
		secondary.State = core.Status_Rejected
		secondary.Reason += fmt.Sprintf("\ndependent on resource %v which failed with err: %v",
			dependencyName, primary.Reason)
	case core.Status_Pending:
		secondary.State = core.Status_Pending
		secondary.Reason += fmt.Sprintf("\ndependent on resource %v which is still pending", dependencyName)
	case core.Status_Accepted:
		if primary.State == core.Status_Pending {
			secondary.State = core.Status_Accepted
			secondary.Reason = ""
		}
	}
	return secondary
}

func PropagateStatuses(namespace string, opts clients.WatchOpts,
	errs chan error,
	resourceClients clients.ResourceClients,
	from resources.InputResourceList,
	to resources.InputResourceList) error {
	sourceResourcesByKind := make(resources.InputResourcesByKind)
	for _, source := range from {
		sourceResourcesByKind.Add(source)
	}
	for kind, sourceResources := range sourceResourcesByKind {
		rc, ok := resourceClients.ForKind(kind)
		if !ok {
			return errors.Errorf("no resource client registered for %v", kind)
		}
		w, watchErrs, err := rc.Watch(namespace, opts)
		if err != nil {
			return err
		}
		go errutils.AggregateErrs(opts.Ctx, errs, watchErrs)
		go func() {
			for {
				select {
				case resourceList := <-w:
					// select only the source resources
					resourceList = resourceList.FilterByList(sourceResources)
				case <-opts.Ctx.Done():
					return
				}
			}
		}()

		for _, dest := range to {

		}
	}
}
