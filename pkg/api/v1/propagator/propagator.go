package propagator

import (
	"fmt"
	"strings"

	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-kit/pkg/errors"
	"github.com/solo-io/solo-kit/pkg/utils/errutils"
)

type ResourcesByType map[string]resources.ResourceList

type Propagator struct {
	forController     string
	children, parents resources.InputResourceList
	resourceClients   clients.ResourceClients
}

func NewPropagator(forController string, parents, children resources.InputResourceList, ResourceClients clients.ResourceClients) *Propagator {
	return &Propagator{
		forController:   forController,
		children:        children,
		parents:         parents,
		resourceClients: ResourceClients,
	}
}

// sources can be multiple types
func (p *Propagator) PropagateStatuses(writeErrs chan error, opts clients.WatchOpts) error {
	// each ressource by kind, then namespace
	childrenByClientAndNamespace, err := byKindByNamespace(p.resourceClients, p.children)
	if err != nil {
		return err
	}

	childrenChannel := make(chan resources.ResourceList)

	if err := createWatchForResources(childrenByClientAndNamespace, childrenChannel, writeErrs, opts); err != nil {
		return errors.Wrapf(err, "creating watch for child resources")
	}

	parentsByClientAndNamespace, err := byKindByNamespace(p.resourceClients, p.parents)
	if err != nil {
		return err
	}
	parentsChannel := make(chan resources.ResourceList)

	if err := createWatchForResources(parentsByClientAndNamespace, parentsChannel, writeErrs, opts); err != nil {
		return errors.Wrapf(err, "creating watch for child resources")
	}

	// aggregate all the different watches, perform sync
	go func() {
		uniqueChildren := make(resources.ResourcesById)
		uniqueParents := make(resources.ResourcesById)
		var lastParents, lastChildren resources.ResourceList
		for {
			select {
			case children := <-childrenChannel:
				if children.Equal(lastChildren) {
					continue
				}
				for _, child := range children {
					uniqueChildren[resources.Key(child)] = child.(resources.InputResource)
				}
				lastChildren = uniqueChildren.List()
				if err := p.syncStatuses(lastParents, lastChildren, opts); err != nil {
					writeErrs <- errors.Wrapf(err, "syncing statuses from children to parents")
				}
			case parents := <-parentsChannel:
				if parents.Equal(lastParents) {
					continue
				}
				for _, parent := range parents {
					uniqueParents[resources.Key(parent)] = parent.(resources.InputResource)
				}
				lastParents = uniqueParents.List()
				if err := p.syncStatuses(lastParents, lastChildren, opts); err != nil {
					writeErrs <- errors.Wrapf(err, "syncing statuses from children to parents")
				}
			case <-opts.Ctx.Done():
				return
			}
		}
	}()
	return nil
}

func createWatchForResources(resByKindAndNamespace map[clients.ResourceClient]map[string]resources.InputResourceList, destinationChannel chan resources.ResourceList, writeErrs chan error, opts clients.WatchOpts) error {
	for clientForKind, resourcesByNamespace := range resByKindAndNamespace {
		for namespace, resourcesToWatch := range resourcesByNamespace {
			watch, errs, err := clientForKind.Watch(namespace, opts)
			if err != nil {
				return err
			}
			go errutils.AggregateErrs(opts.Ctx, writeErrs, errs)
			go receiveResources(watch, resourcesToWatch, destinationChannel, opts)
		}
	}
	return nil
}

func receiveResources(watch <-chan resources.ResourceList, resourcesToWatch resources.InputResourceList, destinationChannel chan resources.ResourceList, opts clients.WatchOpts) {
	for {
		select {
		case resourceList := <-watch:
			// filter only the resources we want
			// TODO(ilackarms): move this abstraction down the stack, see if we can get it into the
			// storage layer api request for max efficiency
			resourceList = resourceList.FilterByNames(resourcesToWatch.Names())
			destinationChannel <- resourceList
		case <-opts.Ctx.Done():
			return
		}
	}
}
func byKindByNamespace(resourceClients clients.ResourceClients, ress resources.InputResourceList) (map[clients.ResourceClient]map[string]resources.InputResourceList, error) {
	resByKindAndNamespace := make(map[clients.ResourceClient]map[string]resources.InputResourceList)
	for _, r := range ress {
		client, err := resourceClients.ForResource(r)
		if err != nil {
			return nil, err
		}
		namespace := r.GetMetadata().Namespace
		if resByKindAndNamespace[client] == nil {
			resByKindAndNamespace[client] = make(map[string]resources.InputResourceList)
		}
		resByKindAndNamespace[client][namespace] = append(resByKindAndNamespace[client][namespace], r)
	}
	return resByKindAndNamespace, nil
}

func containsStatus(st1, st2 core.Status) bool {
	return st1.State == st2.State && st1.ReportedBy == st2.ReportedBy && strings.Contains(st1.Reason, st2.Reason)
}

func (p *Propagator) syncStatuses(parents, children resources.ResourceList, opts clients.WatchOpts) error {
	if !parents.Contains(p.parents.AsResourceList()) {
		return errors.Errorf("updated list of parent resource(s) was missing a resource to update")
	}
	if !children.Contains(p.children.AsResourceList()) {
		return errors.Errorf("updated list of child resource(s) was missing a resource to read status from")
	}
	status, err := createCombinedStatus(p.forController, children)
	if err != nil {
		return err
	}
	for _, parentRes := range parents {
		parent, ok := parentRes.(resources.InputResource)
		if !ok {
			return errors.Errorf("internal error: %v.%v is not an input resource", parentRes.GetMetadata().Namespace, parentRes.GetMetadata().Name)
		}
		if containsStatus(parent.GetStatus(), status) {
			// no-op
			continue
		}
		mergedStatus := mergeStatuses(parent.GetStatus(), status)
		mergedStatus.ReportedBy = p.forController
		parent.SetStatus(mergedStatus)
		rc, err := p.resourceClients.ForResource(parent)
		if err != nil {
			return errors.Wrapf(err, "resource client for parent not found")
		}
		_, err = rc.Write(parent, clients.WriteOpts{
			Ctx:               opts.Ctx,
			OverwriteExisting: true,
		})
		if err != nil {
			return errors.Wrapf(err, "updating status on parent resource")
		}
	}
	return nil
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
			return core.Status{}, errors.Errorf("internal error: %v.%v is not an input resource", baseRes.GetMetadata().Namespace, baseRes.GetMetadata().Name)
		}
		stat := res.GetStatus()
		switch stat.State {
		case core.Status_Rejected:
			state = core.Status_Rejected
			reason += fmt.Sprintf("child resource %v.%v has an error: %v\n", res.GetMetadata().Namespace, res.GetMetadata().Name, stat.Reason)
		case core.Status_Pending:
			// accepteds should be pending
			// errors should still be error
			if state == core.Status_Accepted {
				state = core.Status_Pending
			}
			reason += fmt.Sprintf("child resource %v.%v is still pending\n", res.GetMetadata().Namespace, res.GetMetadata().Name)
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
