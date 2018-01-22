package module

import (
	"sync"

	"github.com/solo-io/glue/config"
	"github.com/solo-io/glue/pkg/log"
)

var ready sync.WaitGroup
var globalRegistry *registry

func init() {
	ready.Add(1)
}

func Init(cfg *config.Config) {
	globalRegistry = rewRegistry(cfg)
	ready.Done()
}

func Register(module Module) {
	go func() {
		ready.Wait()
		globalRegistry.register(module)
	}()
}

type registry struct {
	cfg *config.Config
}

func rewRegistry(cfg *config.Config) *registry {
	return &registry{
		cfg: cfg,
	}
}

func (r *registry) register(module Module) {
	log.Printf("registering module %v\n", module)
	r.cfg.RegisterHandler(func(raw []byte) error {
		blobs, err := getBlobsFromYml(r.cfg.Raw, module.Identifier())
		if err != nil {
			return err
		}
		r.cfg.ResourcesLock.Lock()
		r.cfg.Resources[module.Identifier()], err = module.Translate(nil, blobs)
		r.cfg.ResourcesLock.Unlock()
		if err != nil {
			return err
		}
		return nil
	})
}
