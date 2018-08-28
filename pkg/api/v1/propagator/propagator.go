package propagator

import (
	"fmt"
	"reflect"

	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-kit/pkg/errors"
	"github.com/solo-io/solo-kit/pkg/utils/contextutils"
	"github.com/solo-io/solo-kit/pkg/utils/errutils"
)

type ResourcesByType map[string]resources.ResourceList

type Propagator struct {
	forController     string
	children, parents resources.InputResourceList
	resourceClients   clients.ResourceClients
	writeErrs         chan error
}

func NewPropagator(forController string, parents, children resources.InputResourceList, ResourceClients clients.ResourceClients, writeErrs chan error) *Propagator {
	return &Propagator{
		forController:   forController,
		children:        children,
		parents:         parents,
		resourceClients: ResourceClients,
		writeErrs:       writeErrs,
	}
}

// sources can be multiple types
func (p *Propagator) PropagateStatuses(opts clients.WatchOpts) error {
	// each ressource by kind, then namespace
	childrenByClientAndNamespace, err := byKindByNamespace(p.resourceClients, p.children)
	if err != nil {
		return err
	}

	childrenChannel := make(chan resources.ResourceList)

	if err := createWatchForResources(childrenByClientAndNamespace, childrenChannel, p.writeErrs, opts); err != nil {
		return errors.Wrapf(err, "creating watch for child resources")
	}

	parentsByClientAndNamespace, err := byKindByNamespace(p.resourceClients, p.parents)
	if err != nil {
		return err
	}
	parentsChannel := make(chan resources.ResourceList)

	if err := createWatchForResources(parentsByClientAndNamespace, parentsChannel, p.writeErrs, opts); err != nil {
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
					p.writeErrs <- errors.Wrapf(err, "syncing statuses from children to parents1")
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
					p.writeErrs <- errors.Wrapf(err, "syncing statuses from children to parents2")
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
			go func(namespace string, clientForKind clients.ResourceClient) {
				errutils.AggregateErrs(opts.Ctx, writeErrs, errs, "resource watch on "+fmt.Sprintf("%v.%v", namespace, clientForKind.Kind()))
			}(namespace, clientForKind)
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

func (p *Propagator) syncStatuses(parents, children resources.ResourceList, opts clients.WatchOpts) error {
	if !parents.Contains(p.parents.AsResourceList()) {
		return errors.Errorf("updated list of parent resource(s) was missing a resource to update")
	}
	if !children.Contains(p.children.AsResourceList()) {
		return errors.Errorf("updated list of child resource(s) was missing a resource to read status from")
	}
	childStatuses, err := makeChildStatusMap(children)
	if err != nil {
		return err
	}
	for _, parentRes := range parents {
		parent, ok := parentRes.(resources.InputResource)
		if !ok {
			return errors.Errorf("internal error: %v.%v is not an input resource", parentRes.GetMetadata().Namespace, parentRes.GetMetadata().Name)
		}
		if containsStatuses(parent.GetStatus(), childStatuses) {
			// no-op
			continue
		}
		resources.UpdateStatus(parent, func(status *core.Status) {
			status.SubresourceStatuses = childStatuses
		})
		rc, err := p.resourceClients.ForResource(parent)
		if err != nil {
			return errors.Wrapf(err, "resource client for parent not found")
		}
		_, err = rc.Write(parent, clients.WriteOpts{
			Ctx:               opts.Ctx,
			OverwriteExisting: true,
		})
		if err != nil {
			// ignore rv errors, we should read a new one
			if !errors.IsResourceVersion(err) {
				return errors.Wrapf(err, "updating status on parent resource %v", resources.Key(parent))
			}
			contextutils.LoggerFrom(opts.Ctx).Debugf("received an invalid resource version err on write: %v", err)
		}
	}
	return nil
}

func makeChildStatusMap(children resources.ResourceList) (map[string]*core.Status, error) {
	statuses := make(map[string]*core.Status)
	for _, childRes := range children {
		child, ok := childRes.(resources.InputResource)
		if !ok {
			return nil, errors.Errorf("internal error: %v.%v is not an input resource", childRes.GetMetadata().Namespace, childRes.GetMetadata().Name)
		}
		stat := child.GetStatus()
		statuses[resources.Key(child)] = &stat
	}
	return statuses, nil
}

func containsStatuses(parent core.Status, statuses map[string]*core.Status) bool {
	return reflect.DeepEqual(parent.SubresourceStatuses, statuses)
}
