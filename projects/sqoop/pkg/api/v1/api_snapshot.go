package v1

import (
	"github.com/mitchellh/hashstructure"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources"
	"github.com/solo-io/solo-kit/pkg/api/v1/resources/core"
)

type ApiSnapshot struct {
	ResolverMaps ResolverMapsByNamespace
	Schemas      SchemasByNamespace
}

func (s ApiSnapshot) Clone() ApiSnapshot {
	return ApiSnapshot{
		ResolverMaps: s.ResolverMaps.Clone(),
		Schemas:      s.Schemas.Clone(),
	}
}

func (s ApiSnapshot) Hash() uint64 {
	snapshotForHashing := s.Clone()
	for _, resolverMap := range snapshotForHashing.ResolverMaps.List() {
		resources.UpdateMetadata(resolverMap, func(meta *core.Metadata) {
			meta.ResourceVersion = ""
		})
		resolverMap.SetStatus(core.Status{})
	}
	for _, schema := range snapshotForHashing.Schemas.List() {
		resources.UpdateMetadata(schema, func(meta *core.Metadata) {
			meta.ResourceVersion = ""
		})
		schema.SetStatus(core.Status{})
	}
	h, err := hashstructure.Hash(snapshotForHashing, nil)
	if err != nil {
		panic(err)
	}
	return h
}
