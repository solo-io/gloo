package mocks

import (
	"github.com/mitchellh/hashstructure"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
)

type FestingSnapshot struct {
	Mocks MocksByNamespace
}

func (s FestingSnapshot) Clone() FestingSnapshot {
	return FestingSnapshot{
		Mocks: s.Mocks.Clone(),
	}
}

func (s FestingSnapshot) Hash() uint64 {
	snapshotForHashing := s.Clone()
	for _, mockResource := range snapshotForHashing.Mocks.List() {
		resources.UpdateMetadata(mockResource, func(meta *core.Metadata) {
			meta.ResourceVersion = ""
		})
		mockResource.SetStatus(core.Status{})
	}
	h, err := hashstructure.Hash(snapshotForHashing, nil)
	if err != nil {
		panic(err)
	}
	return h
}
