package plugin

import (
	"github.com/solo-io/glue/pkg/api/types/v1"
	"github.com/solo-io/glue/secrets"
)

type Plugin interface {
	// validate
	Validate(cfg v1.Config)
	// get dependenceies
	GetDependencies(cfg v1.Config) []string
	// create
	Translate(cfg v1.Config, secretMap secrets.SecretMap)
}
