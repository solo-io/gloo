package translator

import (
	"fmt"

	envoycache "github.com/envoyproxy/go-control-plane/pkg/cache"
	"github.com/solo-io/glue/pkg/api/types"
)

type Translator struct{}

func NewTranslator() *Translator {
	return &Translator{}
}

func (t Translator) Translate(cfg types.Config) (envoycache.Snapshot, error) {
	return envoycache.Snapshot{}, fmt.Errorf("not implemented")
}
