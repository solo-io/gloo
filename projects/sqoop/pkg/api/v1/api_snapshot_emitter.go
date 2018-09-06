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
	return &apiEmitter{
		resolverMap: resolverMapClient,
		schema:      schemaClient,
	}
}

type apiEmitter struct {
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
	snapshots := make(chan *ApiSnapshot)
	errs := make(chan error)
	var done sync.WaitGroup

	currentSnapshot := ApiSnapshot{}

	sync := func(newSnapshot ApiSnapshot) {
		if currentSnapshot.Hash() == newSnapshot.Hash() {
			return
		}
		currentSnapshot = newSnapshot
		snapshots <- &currentSnapshot
	}

	for _, namespace := range watchNamespaces {
		resolverMapChan, resolverMapErrs, err := c.resolverMap.Watch(namespace, opts)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "starting ResolverMap watch")
		}
		done.Add(1)
		go func() {
			defer done.Done()
			errutils.AggregateErrs(opts.Ctx, errs, resolverMapErrs, namespace+"-resolverMaps")
		}()

		done.Add(1)
		go func(namespace string, resolverMapChan <-chan ResolverMapList) {
			defer done.Done()
			for {
				select {
				case <-opts.Ctx.Done():
					return
				case resolverMapList := <-resolverMapChan:
					newSnapshot := currentSnapshot.Clone()
					newSnapshot.ResolverMaps.Clear(namespace)
					newSnapshot.ResolverMaps.Add(resolverMapList...)
					sync(newSnapshot)
				}
			}
		}(namespace, resolverMapChan)
		schemaChan, schemaErrs, err := c.schema.Watch(namespace, opts)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "starting Schema watch")
		}
		done.Add(1)
		go func() {
			defer done.Done()
			errutils.AggregateErrs(opts.Ctx, errs, schemaErrs, namespace+"-schemas")
		}()

		done.Add(1)
		go func(namespace string, schemaChan <-chan SchemaList) {
			defer done.Done()
			for {
				select {
				case <-opts.Ctx.Done():
					return
				case schemaList := <-schemaChan:
					newSnapshot := currentSnapshot.Clone()
					newSnapshot.Schemas.Clear(namespace)
					newSnapshot.Schemas.Add(schemaList...)
					sync(newSnapshot)
				}
			}
		}(namespace, schemaChan)
	}

	go func() {
		select {
		case <-opts.Ctx.Done():
			done.Wait()
			close(snapshots)
			close(errs)
		}
	}()
	return snapshots, errs, nil
}
