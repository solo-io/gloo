package mocks

import (
	"github.com/gogo/protobuf/proto"
	"github.com/mitchellh/hashstructure"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-kit/pkg/errors"
)

type Snapshot struct {
	MockResourceList []*MockResource
	FakeResourceList []*FakeResource
	MockDataList []*MockData
}

func (s Snapshot) Clone() Snapshot {
	var mockResourceList []*MockResource
	for _, mockResource := range s.MockResourceList {
		mockResourceList = append(mockResourceList, proto.Clone(mockResource).(*MockResource))
	}
	var fakeResourceList []*FakeResource
	for _, fakeResource := range s.FakeResourceList {
		fakeResourceList = append(fakeResourceList, proto.Clone(fakeResource).(*FakeResource))
	}
	var mockDataList []*MockData
	for _, mockData := range s.MockDataList {
		mockDataList = append(mockDataList, proto.Clone(mockData).(*MockData))
	}
	return Snapshot{
		MockResourceList: mockResourceList,
		FakeResourceList: fakeResourceList,
		MockDataList: mockDataList,
	}
}

func (s Snapshot) Hash() uint64 {
	snapshotForHashing := s.Clone()
	for _, mockResource := range snapshotForHashing.MockResourceList {
		resources.UpdateMetadata(mockResource, func(meta *core.Metadata) {
			meta.ResourceVersion = ""
		})
		mockResource.SetStatus(core.Status{})
	}
	for _, fakeResource := range snapshotForHashing.FakeResourceList {
		resources.UpdateMetadata(fakeResource, func(meta *core.Metadata) {
			meta.ResourceVersion = ""
		})
		fakeResource.SetStatus(core.Status{})
	}
	for _, mockData := range snapshotForHashing.MockDataList {
		resources.UpdateMetadata(mockData, func(meta *core.Metadata) {
			meta.ResourceVersion = ""
		})
		mockData.SetStatus(core.Status{})
	}
	h, err := hashstructure.Hash(snapshotForHashing, nil)
	if err != nil {
		panic(err)
	}
	return h
}

type Cache interface {
	Register() error
	MockResource() MockResourceClient
	FakeResource() FakeResourceClient
	MockData() MockDataClient
	Snapshots(namespace string, opts clients.WatchOpts) (<-chan *Snapshot, <-chan error, error)
}

func NewCache(mockResourceClient MockResourceClient, fakeResourceClient FakeResourceClient, mockDataClient MockDataClient) Cache {
	return &cache{
		mockResource: mockResourceClient,
		fakeResource: fakeResourceClient,
		mockData: mockDataClient,
	}
}

type cache struct {
	mockResource MockResourceClient
	fakeResource FakeResourceClient
	mockData MockDataClient
}

func (c *cache) Register() error {
	if err := c.mockResource.Register(); err != nil {
		return err
	}
	if err := c.fakeResource.Register(); err != nil {
		return err
	}
	if err := c.mockData.Register(); err != nil {
		return err
	}
	return nil
}

func (c *cache) MockResource() MockResourceClient {
	return c.mockResource
}

func (c *cache) FakeResource() FakeResourceClient {
	return c.fakeResource
}

func (c *cache) MockData() MockDataClient {
	return c.mockData
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
	mockResourceChan, mockResourceErrs, err := c.mockResource.Watch(namespace, opts)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "starting MockResource watch")
	}
	fakeResourceChan, fakeResourceErrs, err := c.fakeResource.Watch(namespace, opts)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "starting FakeResource watch")
	}
	mockDataChan, mockDataErrs, err := c.mockData.Watch(namespace, opts)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "starting MockData watch")
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
				sync(newSnapshot)
			case mockDataList := <-mockDataChan:
				newSnapshot := currentSnapshot.Clone()
				newSnapshot.MockDataList = mockDataList
				sync(newSnapshot)
			case err := <-mockResourceErrs:
				errs <- err
			case err := <-fakeResourceErrs:
				errs <- err
			case err := <-mockDataErrs:
				errs <- err
			}
		}
	}()
	return snapshots, errs, nil
}
