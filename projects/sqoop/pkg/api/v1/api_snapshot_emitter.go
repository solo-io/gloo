package v1

import (
	"sync"

	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/errors"
	"github.com/solo-io/solo-kit/pkg/utils/errutils"
)

type ApiEmitter interface {
	Register() error
	ResolverMap() ResolverMapClient
	Schema() SchemaClient
	Snapshots(watchNamespaces []string, opts clients.WatchOpts) (<-chan *ApiSnapshot, <-chan error, error)
}

func NewApiEmitter(resolverMapClient ResolverMapClient, schemaClient SchemaClient) ApiEmitter {
	return NewApiEmitterWithEmit(resolverMapClient, schemaClient, make(chan struct{}))
}

func NewApiEmitterWithEmit(resolverMapClient ResolverMapClient, schemaClient SchemaClient, emit <-chan struct{}) ApiEmitter {
	return &apiEmitter{
		resolverMap: resolverMapClient,
		schema:      schemaClient,
		forceEmit:   emit,
	}
}

type apiEmitter struct {
	forceEmit   <-chan struct{}
	resolverMap ResolverMapClient
	schema      SchemaClient
}

func (c *apiEmitter) Register() error {
	if err := c.resolverMap.Register(); err != nil {
		return err
	}
	if err := c.schema.Register(); err != nil {
		return err
	}
	return nil
}

func (c *apiEmitter) ResolverMap() ResolverMapClient {
	return c.resolverMap
}

func (c *apiEmitter) Schema() SchemaClient {
	return c.schema
}

func (c *apiEmitter) Snapshots(watchNamespaces []string, opts clients.WatchOpts) (<-chan *ApiSnapshot, <-chan error, error) {
	errs := make(chan error)
	var done sync.WaitGroup
	/* Create channel for ResolverMap */
	type resolverMapListWithNamespace struct {
		list      ResolverMapList
		namespace string
	}
	resolverMapChan := make(chan resolverMapListWithNamespace)
	/* Create channel for Schema */
	type schemaListWithNamespace struct {
		list      SchemaList
		namespace string
	}
	schemaChan := make(chan schemaListWithNamespace)

	for _, namespace := range watchNamespaces {
		/* Setup watch for ResolverMap */
		resolverMapNamespacesChan, resolverMapErrs, err := c.resolverMap.Watch(namespace, opts)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "starting ResolverMap watch")
		}

		done.Add(1)
		go func(namespace string) {
			defer done.Done()
			errutils.AggregateErrs(opts.Ctx, errs, resolverMapErrs, namespace+"-resolverMaps")
		}(namespace)
		/* Setup watch for Schema */
		schemaNamespacesChan, schemaErrs, err := c.schema.Watch(namespace, opts)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "starting Schema watch")
		}

		done.Add(1)
		go func(namespace string) {
			defer done.Done()
			errutils.AggregateErrs(opts.Ctx, errs, schemaErrs, namespace+"-schemas")
		}(namespace)

		/* Watch for changes and update snapshot */
		go func(namespace string) {
			for {
				select {
				case <-opts.Ctx.Done():
					return
				case resolverMapList := <-resolverMapNamespacesChan:
					select {
					case <-opts.Ctx.Done():
						return
					case resolverMapChan <- resolverMapListWithNamespace{list: resolverMapList, namespace: namespace}:
					}
				case schemaList := <-schemaNamespacesChan:
					select {
					case <-opts.Ctx.Done():
						return
					case schemaChan <- schemaListWithNamespace{list: schemaList, namespace: namespace}:
					}
				}
			}
		}(namespace)
	}

	snapshots := make(chan *ApiSnapshot)
	go func() {
		currentSnapshot := ApiSnapshot{}
		sync := func(newSnapshot ApiSnapshot) {
			if currentSnapshot.Hash() == newSnapshot.Hash() {
				return
			}
			currentSnapshot = newSnapshot
			sentSnapshot := currentSnapshot.Clone()
			snapshots <- &sentSnapshot
		}
		for {
			select {
			case <-opts.Ctx.Done():
				close(snapshots)
				done.Wait()
				close(errs)
				return
			case <-c.forceEmit:
				sentSnapshot := currentSnapshot.Clone()
				snapshots <- &sentSnapshot
			case resolverMapNamespacedList := <-resolverMapChan:
				namespace := resolverMapNamespacedList.namespace
				resolverMapList := resolverMapNamespacedList.list

				newSnapshot := currentSnapshot.Clone()
				newSnapshot.ResolverMaps.Clear(namespace)
				newSnapshot.ResolverMaps.Add(resolverMapList...)
				sync(newSnapshot)
			case schemaNamespacedList := <-schemaChan:
				namespace := schemaNamespacedList.namespace
				schemaList := schemaNamespacedList.list

				newSnapshot := currentSnapshot.Clone()
				newSnapshot.Schemas.Clear(namespace)
				newSnapshot.Schemas.Add(schemaList...)
				sync(newSnapshot)
			}
		}
	}()
	return snapshots, errs, nil
}
