package mocks

import (
	"sync"

	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/errors"
	"github.com/solo-io/solo-kit/pkg/utils/errutils"
)

type TestingEmitter interface {
	Register() error
	MockResource() MockResourceClient
	FakeResource() FakeResourceClient
	Snapshots(watchNamespaces []string, opts clients.WatchOpts) (<-chan *TestingSnapshot, <-chan error, error)
}

func NewTestingEmitter(mockResourceClient MockResourceClient, fakeResourceClient FakeResourceClient) TestingEmitter {
	return NewTestingEmitterWithEmit(mockResourceClient, fakeResourceClient, make(chan struct{}))
}

func NewTestingEmitterWithEmit(upstreamClient UpstreamClient, emit <-chan struct{}) TestingEmitter {
	return &testingEmitter{
		mockResource: mockResourceClient,
		fakeResource: fakeResourceClient,
		forceEmit:    emit,
	}
}

type testingEmitter struct {
	forceEmit    <-chan struct{}
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

func (c *testingEmitter) Snapshots(watchNamespaces []string, opts clients.WatchOpts) (<-chan *TestingSnapshot, <-chan error, error) {
	errs := make(chan error)
	var done sync.WaitGroup
	/* Create channel for MockResource */
	type mockResourceListWithNamespace struct {
		list      MockResourceList
		namespace string
	}
	mockResourceChan := make(chan mockResourceListWithNamespace)
	/* Create channel for FakeResource */
	type fakeResourceListWithNamespace struct {
		list      FakeResourceList
		namespace string
	}
	fakeResourceChan := make(chan fakeResourceListWithNamespace)

	for _, namespace := range watchNamespaces {
		/* Setup watch for MockResource */
		mockResourceNamespacesChan, mockResourceErrs, err := c.mockResource.Watch(namespace, opts)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "starting MockResource watch")
		}

		done.Add(1)
		go func(namespace string) {
			defer done.Done()
			errutils.AggregateErrs(opts.Ctx, errs, mockResourceErrs, namespace+"-mocks")
		}(namespace)
		/* Setup watch for FakeResource */
		fakeResourceNamespacesChan, fakeResourceErrs, err := c.fakeResource.Watch(namespace, opts)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "starting FakeResource watch")
		}

		done.Add(1)
		go func(namespace string) {
			defer done.Done()
			errutils.AggregateErrs(opts.Ctx, errs, fakeResourceErrs, namespace+"-fakes")
		}(namespace)

		/* Watch for changes and update snapshot */
		go func(namespace string) {
			for {
				select {
				case <-opts.Ctx.Done():
					return
				case mockResourceList := <-mockResourceNamespacesChan:
					select {
					case <-opts.Ctx.Done():
						return
					case mockResourceChan <- mockResourceListWithNamespace{list: mockResourceList, namespace: namespace}:
					}
				case fakeResourceList := <-fakeResourceNamespacesChan:
					select {
					case <-opts.Ctx.Done():
						return
					case fakeResourceChan <- fakeResourceListWithNamespace{list: fakeResourceList, namespace: namespace}:
					}
				}
			}
		}(namespace)
	}

	snapshots := make(chan *TestingSnapshot)
	go func() {
		currentSnapshot := TestingSnapshot{}
		sync := func(newSnapshot TestingSnapshot) {
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
			case mockResourceNamespacedList := <-mockResourceChan:
				namespace := mockResourceNamespacedList.namespace
				mockResourceList := mockResourceNamespacedList.list

				newSnapshot := currentSnapshot.Clone()
				newSnapshot.Mocks.Clear(namespace)
				newSnapshot.Mocks.Add(mockResourceList...)
				sync(newSnapshot)
			case fakeResourceNamespacedList := <-fakeResourceChan:
				namespace := fakeResourceNamespacedList.namespace
				fakeResourceList := fakeResourceNamespacedList.list

				newSnapshot := currentSnapshot.Clone()
				newSnapshot.Fakes.Clear(namespace)
				newSnapshot.Fakes.Add(fakeResourceList...)
				sync(newSnapshot)
			}
		}
	}()
	return snapshots, errs, nil
}
