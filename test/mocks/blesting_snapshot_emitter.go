package mocks

import (
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
	return &blestingEmitter{
		fakeResource: fakeResourceClient,
	}
}

type blestingEmitter struct {
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
	snapshots := make(chan *BlestingSnapshot)
	errs := make(chan error)

	currentSnapshot := BlestingSnapshot{}

	sync := func(newSnapshot BlestingSnapshot) {
		if currentSnapshot.Hash() == newSnapshot.Hash() {
			return
		}
		currentSnapshot = newSnapshot
		snapshots <- &currentSnapshot
	}

	for _, namespace := range watchNamespaces {
		fakeResourceChan, fakeResourceErrs, err := c.fakeResource.Watch(namespace, opts)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "starting FakeResource watch")
		}
		go errutils.AggregateErrs(opts.Ctx, errs, fakeResourceErrs, namespace+"-fakes")
		go func(namespace string, fakeResourceChan <-chan FakeResourceList) {
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
