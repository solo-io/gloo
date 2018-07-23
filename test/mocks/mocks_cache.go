package mocks

import (
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/factory"
	"github.com/solo-io/solo-kit/pkg/errors"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/mitchellh/hashstructure"
	"github.com/gogo/protobuf/proto"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
)

type Snapshot struct {
	MockResourceList []*MockResource
	FakeResourceList []*FakeResource
}

func (s Snapshot) Clone() Snapshot {
	var mockResourceList []*MockResource
	var fakeResourceList []*FakeResource
	for _, mockResource := range s.MockResourceList {
		mockResourceList = append(mockResourceList, proto.Clone(mockResource).(*MockResource))
	}
	for _, fakeResource := range s.FakeResourceList {
		fakeResourceList = append(fakeResourceList, proto.Clone(fakeResource).(*FakeResource))
	}
	return Snapshot{
		MockResourceList: mockResourceList,
		FakeResourceList: fakeResourceList,
	}
}

func (s Snapshot) Hash() uint64 {
	snapshotForHashing := s.Clone()
	for _, mockResource := range s.MockResourceList {
		resources.UpdateMetadata(mockResource, func(meta *core.Metadata) {
			meta.ResourceVersion = ""
		})
		mockResource.SetStatus(core.Status{})
	}
	for _, fakeResource := range s.FakeResourceList {
		resources.UpdateMetadata(fakeResource, func(meta *core.Metadata) {
			meta.ResourceVersion = ""
		})
		fakeResource.SetStatus(core.Status{})
	}
	h, err := hashstructure.Hash(snapshotForHashing, nil)
	if err != nil {
		panic(err)
	}
	return h
}

type Cache interface {
	MockResource() MockResourceClient
	FakeResource() FakeResourceClient
	Snapshots(opts clients.WatchOpts) (<-chan *Snapshot, <-chan error, error)
}

func NewCache(factory *factory.ResourceClientFactory) Cache {
	return &cache{
		mockResource: NewMockResourceClient(factory),
		fakeResource: NewFakeResourceClient(factory),
	}
}

type cache struct {
	mockResource MockResourceClient
	fakeResource FakeResourceClient
}

func (c *cache) MockResource() MockResourceClient {
	return c.mockResource
}

func (c *cache) FakeResource() FakeResourceClient {
	return c.fakeResource
}


func (c *cache) Snapshots(opts clients.WatchOpts) (<-chan *Snapshot, <-chan error, error) {
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

	mockResourceChan, mockResourceErrs, err := c.mockResource.Watch(opts)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "starting MockResource watch")
	}

	fakeResourceChan, fakeResourceErrs, err := c.fakeResource.Watch(opts)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "starting FakeResource watch")
	}
	go func() {
		for {
			select {
			case mockResourceList := <-mockResourceChan:
				newSnapshot := currentSnapshot.Clone()
				newSnapshot.MockResourceList = mockResourceList
				sync(newSnapshot)
			case fakeResourceList := <-fakeResourceChan:
				newSnapshot := currentSnapshot.Clone()
				newSnapshot.FakeResourceList = fakeResourceList
			case err := <- mockResourceErrs:
				errs <- err
			case err := <- fakeResourceErrs:
				errs <- err
			}
		}
	}()
	return snapshots, errs, nil
}
