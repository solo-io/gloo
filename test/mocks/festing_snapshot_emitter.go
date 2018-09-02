package mocks

import (
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
	return &festingEmitter{
		mockResource: mockResourceClient,
	}
}

type festingEmitter struct {
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
	snapshots := make(chan *FestingSnapshot)
	errs := make(chan error)

	currentSnapshot := FestingSnapshot{}

	sync := func(newSnapshot FestingSnapshot) {
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
