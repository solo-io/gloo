package mocks

import (
	"sync"
	"time"

	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"

	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/errors"
	"github.com/solo-io/solo-kit/pkg/utils/errutils"
)

var (
	mFestingSnapshotIn  = stats.Int64("festing_snap_emitter/snap_in", "The number of snapshots in", "1")
	mFestingSnapshotOut = stats.Int64("festing_snap_emitter/snap_out", "The number of snapshots out", "1")

	festingsnapshotInView = &view.View{
		Name:        "festing_snap_emitter/snap_in",
		Measure:     mFestingSnapshotIn,
		Description: "The number of snapshots updates coming in",
		Aggregation: view.Count(),
		TagKeys:     []tag.Key{},
	}
	festingsnapshotOutView = &view.View{
		Name:        "festing_snap_emitter/snap_out",
		Measure:     mFestingSnapshotOut,
		Description: "The number of snapshots updates going out",
		Aggregation: view.Count(),
		TagKeys:     []tag.Key{},
	}
)

func init() {
	view.Register(festingsnapshotInView, festingsnapshotOutView)
}

type FestingEmitter interface {
	Register() error
	MockResource() MockResourceClient
	Snapshots(watchNamespaces []string, opts clients.WatchOpts) (<-chan *FestingSnapshot, <-chan error, error)
}

func NewFestingEmitter(mockResourceClient MockResourceClient) FestingEmitter {
	return NewFestingEmitterWithEmit(mockResourceClient, make(chan struct{}))
}

func NewFestingEmitterWithEmit(mockResourceClient MockResourceClient, emit <-chan struct{}) FestingEmitter {
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
	ctx := opts.Ctx
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
			errutils.AggregateErrs(ctx, errs, mockResourceErrs, namespace+"-mocks")
		}(namespace)

		/* Watch for changes and update snapshot */
		go func(namespace string) {
			for {
				select {
				case <-ctx.Done():
					return
				case mockResourceList := <-mockResourceNamespacesChan:
					select {
					case <-ctx.Done():
						return
					case mockResourceChan <- mockResourceListWithNamespace{list: mockResourceList, namespace: namespace}:
					}
				}
			}
		}(namespace)
	}

	snapshots := make(chan *FestingSnapshot)
	go func() {
		originalSnapshot := FestingSnapshot{}
		currentSnapshot := originalSnapshot.Clone()
		timer := time.NewTicker(time.Second * 5)
		sync := func() {
			if originalSnapshot.Hash() == currentSnapshot.Hash() {
				return
			}
			originalSnapshot = currentSnapshot.Clone()
			sentSnapshot := currentSnapshot.Clone()
			snapshots <- &sentSnapshot
		}

		/* TODO (yuval-k): figure out how to make this work to avoid a stale snapshot.
		   		// construct the first snapshot from all the configs that are currently there
		   		// that guarantees that the first snapshot contains all the data.
		   		for range watchNamespaces {
		      mockResourceNamespacedList := <- mockResourceChan:
		   	namespace := mockResourceNamespacedList.namespace
		   	mockResourceList := mockResourceNamespacedList.list

		   	currentSnapshot.Mocks.Clear(namespace)
		   	currentSnapshot.Mocks.Add(mockResourceList...)
		   		}
		*/

		for {
			select {
			case <-timer.C:
				sync()
			case <-ctx.Done():
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

				currentSnapshot.Mocks.Clear(namespace)
				currentSnapshot.Mocks.Add(mockResourceList...)
			}

			// if we got here its because a new entry in the channel
			stats.Record(ctx, mFestingSnapshotIn.M(1))
		}
	}()
	return snapshots, errs, nil
}
