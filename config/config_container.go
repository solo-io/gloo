package config

import (
	"sync"
)

type UpdateHandler func(blobs []byte) error

type Config struct {
	Raw           []byte
	Handlers      []UpdateHandler
	Resources     map[string]EnvoyResources
	ResourcesLock sync.RWMutex
}

func NewConfig() *Config {
	return &Config{
		Resources: make(map[string]EnvoyResources),
	}
}

func (c *Config) RegisterHandler(handler UpdateHandler) {
	c.Handlers = append(c.Handlers, handler)
}

func (c *Config) Update(data []byte) error {
	c.Raw = data
	for _, cb := range c.Handlers {
		if err := cb(c.Raw); err != nil {
			return err
		}
	}
	return nil
}
