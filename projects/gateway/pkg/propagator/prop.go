package propagator

import (
	"fmt"

	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-kit/pkg/errors"
	"github.com/solo-io/solo-kit/pkg/utils/errutils"
)

type ResourcesByType map[string]resources.ResourceList

type Propagator struct {
	ForController         string
	ResourceClients       clients.ResourceClients
	Sources, Destinations resources.InputResourceList
}

func NewPropagator(forController string, destinations, sources resources.InputResourceList, ResourceClients clients.ResourceClients) *Propagator {
	return &Propagator{
		ForController: forController,
		Sources:       sources,
		Destinations:  destinations,
	}
}

// sources can be multiple types
func (p *Propagator) PropagateStatus(writeErrs chan error, opts clients.WatchOpts) error {
	// each ressource by kind
	sourcesByClient := make(map[clients.ResourceClient]resources.InputResourceList)
	for _, r := range p.Sources {
		client, err := p.ResourceClients.ForResource(r)
		if err != nil {
			return err
		}
		sourcesByClient[client] = append(sourcesByClient[client], r)
	}
	destinationsByClient := make(map[clients.ResourceClient]resources.InputResourceList)
	for _, r := range p.Destinations {
		client, err := p.ResourceClients.ForResource(r)
		if err != nil {
			return err
		}
		destinationsByClient[client] = append(destinationsByClient[client], r)
	}

	syncStatuses := func(destinations, sources resources.InputResourceList) error {
		for _, requiredDests := range destinationsByClient {
			requiredDests
		}
		for _, dest := range destinations {

		}
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
					// ensure all the resiources we expect to be there are there
					// filter only the resources we want
					// TODO(ilackarms): move this abstraction down the stack, see if we can get it into the
					// kube api request for max efficiency
					resourceList = resourceList.FilterByNamespaces(sourcesToWatch.Namespaces()).FilterByNames(sourcesToWatch.Names())
					status, err := createCombinedStatus(p.ForController, resourceList)
					if err != nil {
						writeErrs <- err
						continue
					}
					for _, dest := range p.Destinations {
						status = mergeStatuses(dest.GetStatus(), status)
						if dest.GetStatus().Equal(status) {
							// no-op
							continue
						}
						dest.SetStatus(status)
					}
				case <-opts.Ctx.Done():
					return
				}
			}
		}()
	}
}

func mergeStatuses(dest, src core.Status) core.Status {
	switch src.State {
	case core.Status_Accepted:
	case core.Status_Pending:
		if dest.State == core.Status_Accepted {
			dest.State = core.Status_Pending
		}
		dest.Reason += src.Reason
	case core.Status_Rejected:
		dest.State = core.Status_Rejected
		dest.Reason += src.Reason
	}
	return dest
}

func createCombinedStatus(forController string, fromResources resources.ResourceList) (core.Status, error) {
	state := core.Status_Accepted
	reason := ""

	for _, baseRes := range fromResources {
		res, ok := baseRes.(resources.InputResource)
		if !ok {
			return core.Status{}, errors.Errorf("internal error: %v.%v is not an input resource", baseRes.GetMetadata().ObjectRef())
		}
		stat := res.GetStatus()
		switch stat.State {
		case core.Status_Rejected:
			state = core.Status_Rejected
			reason += fmt.Sprintf("child resource %v.%v has an error\n", res.GetMetadata().ObjectRef())
		case core.Status_Pending:
			// accepteds should be pending
			// errors should still be error
			if state == core.Status_Accepted {
				state = core.Status_Pending
			}
			reason += fmt.Sprintf("child resource %v.%v is still pending\n", res.GetMetadata().ObjectRef())
		case core.Status_Accepted:
			continue
		}
	}
	return core.Status{
		State:      state,
		Reason:     reason,
		ReportedBy: forController,
	}, nil
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
