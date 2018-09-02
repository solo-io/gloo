package mocks

import (
	"github.com/mitchellh/hashstructure"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
)

type BlestingSnapshot struct {
	Mocks MocksByNamespace
	Fakes FakesByNamespace
}

func (s BlestingSnapshot) Clone() BlestingSnapshot {
	return BlestingSnapshot{
		Mocks: s.Mocks.Clone(),
		Fakes: s.Fakes.Clone(),
	}
}

func (s BlestingSnapshot) Hash() uint64 {
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
	h, err := hashstructure.Hash(snapshotForHashing, nil)
	if err != nil {
		panic(err)
	}
	return h
}
