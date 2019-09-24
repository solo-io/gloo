package nackdetector

import (
	"context"

	"log"

	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
)

var (
	nodeIdKey, _   = tag.NewKey("nodeid")
	resourceKey, _ = tag.NewKey("resource")

	mGlooeXdsTotalEntities = stats.Int64("glooe.solo.io/xds/total_entities", "Total number of entities", "1")
	GlooeTotalEntities     = &view.View{
		Name:        "glooe.solo.io/xds/total_entities",
		Measure:     mGlooeXdsTotalEntities,
		Description: "The total number of XDS streams",
		Aggregation: view.Sum(),
		TagKeys:     []tag.Key{resourceKey},
	}

	mGlooeXdsOutOfSync = stats.Int64("glooe.solo.io/xds/outofsync", "The of envoys out of sync", "1")
	GlooeOutOfSync     = &view.View{
		Name:        "glooe.solo.io/xds/outofsync",
		Measure:     mGlooeXdsOutOfSync,
		Description: "The number of envoys out of sync",
		Aggregation: view.Sum(),
		TagKeys:     []tag.Key{resourceKey},
	}

	mGlooeXdsNack = stats.Int64("glooe.solo.io/xds/nack", "The of envoys reported a nack", "1")
	GlooeNack     = &view.View{
		Name:        "glooe.solo.io/xds/nack",
		Measure:     mGlooeXdsNack,
		Description: "The number of envoys that reported NACK",
		Aggregation: view.Sum(),
		TagKeys:     []tag.Key{resourceKey},
	}

	mGlooeXdsInSync = stats.Int64("glooe.solo.io/xds/insync", "The of envoys in sync", "1")
	GlooeInSync     = &view.View{
		Name:        "glooe.solo.io/xds/insync",
		Measure:     mGlooeXdsInSync,
		Description: "The envoys that are in sync",
		Aggregation: view.Sum(),
		TagKeys:     []tag.Key{resourceKey},
	}
)

func init() {
	if err := view.Register(GlooeTotalEntities, GlooeOutOfSync, GlooeNack, GlooeInSync); err != nil {
		log.Printf("failed to register stats views [%v,%v,%v,%v]", GlooeTotalEntities, GlooeOutOfSync, GlooeNack, GlooeInSync)

	}
}

type StatGen struct {
	ctx context.Context
}

var _ StateChangedCallback = new(StatGen).Stat

func NewStatsGen(ctx context.Context) *StatGen {
	return &StatGen{
		ctx: ctx,
	}

}

func (s *StatGen) Stat(id EnvoyStatusId, oldst, st State) {
	ctx := s.ctx

	// if ctxWithTags, err := tag.New(ctx, tag.Insert(nodeIdKey, id.NodeId)); err == nil {
	// 	ctx = ctxWithTags
	// }
	record := func(metric *stats.Int64Measure, v int64) {
		stats.RecordWithTags(ctx, tags(id), metric.M(v))
	}
	switch st {
	case New:
		record(mGlooeXdsTotalEntities, 1)
	case InSync:
		record(mGlooeXdsInSync, 1)
	case OutOfSync:
		record(mGlooeXdsOutOfSync, 1)
	case OutOfSyncNack:
		record(mGlooeXdsNack, 1)
	case Gone:
		record(mGlooeXdsTotalEntities, -1)
	}

	switch oldst {
	case New:
		// this case is handled above
	case InSync:
		record(mGlooeXdsInSync, -1)
	case OutOfSync:
		record(mGlooeXdsOutOfSync, -1)
	case OutOfSyncNack:
		record(mGlooeXdsNack, -1)
	case Gone:
		// this case should never happen
	}

}

func tags(id EnvoyStatusId) []tag.Mutator {
	return []tag.Mutator{tag.Insert(resourceKey, id.StreamId.TypeUrl)}
}
