package config

import (
	"sync"
)

type UpdateHandler func(blobs []byte) error

type Config struct {
	raw           []byte
	handlers      []UpdateHandler
	resources     map[string]EnvoyResources
	resourcesLock sync.RWMutex
}

func NewConfig() *Config {
	return &Config{
		resources: make(map[string]EnvoyResources),
	}
}

func (c *Config) RegisterHandler(handler UpdateHandler) {
	c.handlers = append(c.handlers, handler)
}

func (c *Config) Update(data []byte) error {
	c.raw = data
	for _, cb := range c.handlers {
		if err := cb(c.raw); err != nil {
			return err
		}
	}
	return nil
}

func (c *Config) UpdateResources(moduleIdentifier string, resources EnvoyResources) {
	c.resourcesLock.Lock()
	defer c.resourcesLock.Unlock()
	c.resources[moduleIdentifier] = resources
}

func (c *Config) GetResources() []EnvoyResources {
	var resources []EnvoyResources
	c.resourcesLock.RLock()
	defer c.resourcesLock.RUnlock()
	for _, resource := range c.resources {
		resources = append(resources, resource)
	}
	return resources
}

func (c *Config) Raw() []byte {
	return c.raw
}
