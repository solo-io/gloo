package v1

import (
	"github.com/mitchellh/hashstructure"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-kit/pkg/errors"
	"github.com/solo-io/solo-kit/pkg/utils/errutils"
)

type Snapshot struct {
	ResolverMaps ResolverMapsByNamespace
}

func (s Snapshot) Clone() Snapshot {
	return Snapshot{
		ResolverMaps: s.ResolverMaps.Clone(),
	}
}

func (s Snapshot) Hash() uint64 {
	snapshotForHashing := s.Clone()
	for _, resolverMap := range snapshotForHashing.ResolverMaps.List() {
		resources.UpdateMetadata(resolverMap, func(meta *core.Metadata) {
			meta.ResourceVersion = ""
		})
		resolverMap.SetStatus(core.Status{})
	}
	h, err := hashstructure.Hash(snapshotForHashing, nil)
	if err != nil {
		panic(err)
	}
	return h
}

type Cache interface {
	Register() error
	ResolverMap() ResolverMapClient
	Snapshots(watchNamespaces []string, opts clients.WatchOpts) (<-chan *Snapshot, <-chan error, error)
}

func NewCache(resolverMapClient ResolverMapClient) Cache {
	return &cache{
		resolverMap: resolverMapClient,
	}
}

type cache struct {
	resolverMap ResolverMapClient
}

func (c *cache) Register() error {
	if err := c.resolverMap.Register(); err != nil {
		return err
	}
	return nil
}

func (c *cache) ResolverMap() ResolverMapClient {
	return c.resolverMap
}

func (c *cache) Snapshots(watchNamespaces []string, opts clients.WatchOpts) (<-chan *Snapshot, <-chan error, error) {
	snapshots := make(chan *Snapshot)
	errs := make(chan error)

	currentSnapshot := Snapshot{}

	sync := func(newSnapshot Snapshot) {
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
		go errutils.AggregateErrs(opts.Ctx, errs, resolverMapErrs, namespace+"-resolverMaps")
		go func(namespace string, resolverMapChan <-chan ResolverMapList) {
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
	}

	go func() {
		select {
		case <-opts.Ctx.Done():
			close(snapshots)
			close(errs)
		}
	}()
	return snapshots, errs, nil
}
