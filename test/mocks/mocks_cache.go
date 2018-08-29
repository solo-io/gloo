package mocks

import (
	"github.com/mitchellh/hashstructure"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
	"github.com/solo-io/solo-kit/pkg/errors"
	"github.com/solo-io/solo-kit/pkg/utils/errutils"
)

type Snapshot struct {
	Mocks     MocksByNamespace
	Fakes     FakesByNamespace
	MockDatas MockDatasByNamespace
}

func (s Snapshot) Clone() Snapshot {
	return Snapshot{
		Mocks:     s.Mocks.Clone(),
		Fakes:     s.Fakes.Clone(),
		MockDatas: s.MockDatas.Clone(),
	}
}

func (s Snapshot) Hash() uint64 {
	snapshotForHashing := s.Clone()
	for _, mockResource := range snapshotForHashing.Mocks.List() {
		resources.UpdateMetadata(mockResource, func(meta *core.Metadata) {
			meta.ResourceVersion = ""
		})
		mockResource.SetStatus(core.Status{})
	}
	for _, fakeResource := range snapshotForHashing.Fakes.List() {
		resources.UpdateMetadata(fakeResource, func(meta *core.Metadata) {
			meta.ResourceVersion = ""
		})
		fakeResource.SetStatus(core.Status{})
	}
	for _, mockData := range snapshotForHashing.MockDatas.List() {
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
	Snapshots(watchNamespaces []string, opts clients.WatchOpts) (<-chan *Snapshot, <-chan error, error)
}

func NewCache(mockResourceClient MockResourceClient, fakeResourceClient FakeResourceClient, mockDataClient MockDataClient) Cache {
	return &cache{
		mockResource: mockResourceClient,
		fakeResource: fakeResourceClient,
		mockData:     mockDataClient,
	}
}

type cache struct {
	mockResource MockResourceClient
	fakeResource FakeResourceClient
	mockData     MockDataClient
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
		mockResourceChan, mockResourceErrs, err := c.mockResource.Watch(namespace, opts)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "starting MockResource watch")
		}
		go errutils.AggregateErrs(opts.Ctx, errs, mockResourceErrs, namespace+"-mocks")
		go func() {
			for {
				select {
				case <-opts.Ctx.Done():
					return
				case mockResourceList := <-mockResourceChan:
					newSnapshot := currentSnapshot.Clone()
					newSnapshot.Mocks.Clear(namespace)
					newSnapshot.Mocks.Add(mockResourceList...)
					sync(newSnapshot)
				}
			}
		}()
		fakeResourceChan, fakeResourceErrs, err := c.fakeResource.Watch(namespace, opts)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "starting FakeResource watch")
		}
		go errutils.AggregateErrs(opts.Ctx, errs, fakeResourceErrs, namespace+"-fakes")
		go func() {
			for {
				select {
				case <-opts.Ctx.Done():
					return
				case fakeResourceList := <-fakeResourceChan:
					newSnapshot := currentSnapshot.Clone()
					newSnapshot.Fakes.Clear(namespace)
					newSnapshot.Fakes.Add(fakeResourceList...)
					sync(newSnapshot)
				}
			}
		}()
		mockDataChan, mockDataErrs, err := c.mockData.Watch(namespace, opts)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "starting MockData watch")
		}
		go errutils.AggregateErrs(opts.Ctx, errs, mockDataErrs, namespace+"-mockDatas")
		go func() {
			for {
				select {
				case <-opts.Ctx.Done():
					return
				case mockDataList := <-mockDataChan:
					newSnapshot := currentSnapshot.Clone()
					newSnapshot.MockDatas.Clear(namespace)
					newSnapshot.MockDatas.Add(mockDataList...)
					sync(newSnapshot)
				}
			}
		}()
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
