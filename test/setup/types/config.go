package types

import (
	"encoding/json"
	"os"
)

// Config is the top-level configuration for the application.
type Config struct {
	// Clusters is a list of clusters to create.
	Clusters []*Cluster `yaml:"clusters,omitempty" json:"clusters,omitempty"`
}

func (c *Config) GetClusters() []*Cluster {
	return c.Clusters
}

func (c *Config) Print() error {
	enc := json.NewEncoder(os.Stdout)

	enc.SetIndent("", "  ")

	return enc.Encode(c)
}
