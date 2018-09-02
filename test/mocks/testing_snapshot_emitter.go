package mocks

import (
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/errors"
	"github.com/solo-io/solo-kit/pkg/utils/errutils"
)

type TestingEmitter interface {
	Register() error
	MockResource() MockResourceClient
	FakeResource() FakeResourceClient
	Snapshots(watchNamespaces []string, opts clients.WatchOpts) (<-chan *Snapshot, <-chan error, error)
}

func NewTestingEmitter(mockResourceClient MockResourceClient, fakeResourceClient FakeResourceClient) TestingEmitter {
	return &testingEmitter{
		mockResource: mockResourceClient,
		fakeResource: fakeResourceClient,
	}
}

type testingEmitter struct {
	mockResource MockResourceClient
	fakeResource FakeResourceClient
}

func (c *testingEmitter) Register() error {
	if err := c.mockResource.Register(); err != nil {
		return err
	}
	if err := c.fakeResource.Register(); err != nil {
		return err
	}
	return nil
}

func (c *testingEmitter) MockResource() MockResourceClient {
	return c.mockResource
}

func (c *testingEmitter) FakeResource() FakeResourceClient {
	return c.fakeResource
}

func (c *testingEmitter) Snapshots(watchNamespaces []string, opts clients.WatchOpts) (<-chan *Snapshot, <-chan error, error) {
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
		go func(namespace string, mockResourceChan  <- chan MockResourceList) {
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
		}(namespace, mockResourceChan)
		fakeResourceChan, fakeResourceErrs, err := c.fakeResource.Watch(namespace, opts)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "starting FakeResource watch")
		}
		go errutils.AggregateErrs(opts.Ctx, errs, fakeResourceErrs, namespace+"-fakes")
		go func(namespace string, fakeResourceChan  <- chan FakeResourceList) {
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
		}(namespace, fakeResourceChan)
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
