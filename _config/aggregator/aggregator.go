package aggregator

import (
	"github.com/envoyproxy/go-control-plane/pkg/cache"
	"github.com/solo-io/glue/config"
	"github.com/solo-io/glue/pkg/model"
)

type Aggregator interface {
	CreateSnapshot(cfg *config.Config, modules []Module) (cache.Snapshot, error)
}

type Module interface {
	Parse(blobs [][]byte) (model.Config, error)
}

type aggregator struct{}

func (agr *aggregator) CreateSnapshot(cfg *config.Config, modules []Module) (cache.Snapshot, error) {

}
