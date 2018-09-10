package mocks

import (
	"sync"

	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/errors"
	"github.com/solo-io/solo-kit/pkg/utils/errutils"
)

type BlestingEmitter interface {
	Register() error
	FakeResource() FakeResourceClient
	Snapshots(watchNamespaces []string, opts clients.WatchOpts) (<-chan *BlestingSnapshot, <-chan error, error)
}

func NewBlestingEmitter(fakeResourceClient FakeResourceClient) BlestingEmitter {
	return NewBlestingEmitterWithEmit(fakeResourceClient, make(chan struct{}))
}

func NewBlestingEmitterWithEmit(fakeResourceClient FakeResourceClient, emit <-chan struct{}) BlestingEmitter {
	return &blestingEmitter{
		fakeResource: fakeResourceClient,
		forceEmit:    emit,
	}
}

type blestingEmitter struct {
	forceEmit    <-chan struct{}
	fakeResource FakeResourceClient
}

func (c *blestingEmitter) Register() error {
	if err := c.fakeResource.Register(); err != nil {
		return err
	}
	return nil
}

func (c *blestingEmitter) FakeResource() FakeResourceClient {
	return c.fakeResource
}

func (c *blestingEmitter) Snapshots(watchNamespaces []string, opts clients.WatchOpts) (<-chan *BlestingSnapshot, <-chan error, error) {
	errs := make(chan error)
	var done sync.WaitGroup
	/* Create channel for FakeResource */
	type fakeResourceListWithNamespace struct {
		list      FakeResourceList
		namespace string
	}
	fakeResourceChan := make(chan fakeResourceListWithNamespace)

	for _, namespace := range watchNamespaces {
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

	snapshots := make(chan *BlestingSnapshot)
	go func() {
		currentSnapshot := BlestingSnapshot{}
		sync := func(newSnapshot BlestingSnapshot) {
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
