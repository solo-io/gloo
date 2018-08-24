package v1

import (
	"github.com/gogo/protobuf/proto"
	"github.com/mitchellh/hashstructure"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-kit/pkg/errors"
)

type Snapshot struct {
	ResolverMapList ResolverMapList
}

func (s Snapshot) Clone() Snapshot {
	var resolverMapList []*ResolverMap
	for _, resolverMap := range s.ResolverMapList {
		resolverMapList = append(resolverMapList, proto.Clone(resolverMap).(*ResolverMap))
	}
	return Snapshot{
		ResolverMapList: resolverMapList,
	}
}

func (s Snapshot) Hash() uint64 {
	snapshotForHashing := s.Clone()
	for _, resolverMap := range snapshotForHashing.ResolverMapList {
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
	Snapshots(namespace string, opts clients.WatchOpts) (<-chan *Snapshot, <-chan error, error)
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

func (c *cache) Snapshots(namespace string, opts clients.WatchOpts) (<-chan *Snapshot, <-chan error, error) {
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
	resolverMapChan, resolverMapErrs, err := c.resolverMap.Watch(namespace, opts)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "starting ResolverMap watch")
	}

	go func() {
		for {
			select {
			case resolverMapList := <-resolverMapChan:
				newSnapshot := currentSnapshot.Clone()
				newSnapshot.ResolverMapList = resolverMapList
				sync(newSnapshot)
			case err := <-resolverMapErrs:
				errs <- err
			}
		}
	}()
	return snapshots, errs, nil
}
