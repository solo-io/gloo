package mocks

import (
	"sync"

	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/errors"
	"github.com/solo-io/solo-kit/pkg/utils/errutils"
)

type FestingEmitter interface {
	Register() error
	MockResource() MockResourceClient
	Snapshots(watchNamespaces []string, opts clients.WatchOpts) (<-chan *FestingSnapshot, <-chan error, error)
}

func NewFestingEmitter(mockResourceClient MockResourceClient) FestingEmitter {
	return NewFestingEmitterWithEmit(mockResourceClient, make(chan struct{}))
}

func NewFestingEmitterWithEmit(upstreamClient UpstreamClient, emit <-chan struct{}) FestingEmitter {
	return &festingEmitter{
		mockResource: mockResourceClient,
		forceEmit:    emit,
	}
}

type festingEmitter struct {
	forceEmit    <-chan struct{}
	mockResource MockResourceClient
}

func (c *festingEmitter) Register() error {
	if err := c.mockResource.Register(); err != nil {
		return err
	}
	return nil
}

func (c *festingEmitter) MockResource() MockResourceClient {
	return c.mockResource
}

func (c *festingEmitter) Snapshots(watchNamespaces []string, opts clients.WatchOpts) (<-chan *FestingSnapshot, <-chan error, error) {
	errs := make(chan error)
	var done sync.WaitGroup
	/* Create channel for MockResource */
	type mockResourceListWithNamespace struct {
		list      MockResourceList
		namespace string
	}
	mockResourceChan := make(chan mockResourceListWithNamespace)

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
				}
			}
		}(namespace)
	}

	snapshots := make(chan *FestingSnapshot)
	go func() {
		currentSnapshot := FestingSnapshot{}
		sync := func(newSnapshot FestingSnapshot) {
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
			}
		}
	}()
	return snapshots, errs, nil
}
