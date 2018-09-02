package mocks

import (
	"github.com/mitchellh/hashstructure"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
)

type TestingSnapshot struct {
	Mocks MocksByNamespace
	Fakes FakesByNamespace
}

func (s TestingSnapshot) Clone() TestingSnapshot {
	return TestingSnapshot{
		Mocks: s.Mocks.Clone(),
		Fakes: s.Fakes.Clone(),
	}
}

func (s TestingSnapshot) Hash() uint64 {
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
