package plugin

import (
	envoyapi "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	"github.com/solo-io/glue/pkg/api/types/v1"
	"github.com/solo-io/glue/pkg/secretwatcher"
)

type Dependencies struct {
	SecretRefs []string
}

type Inputs struct {
	Cfg     v1.Config
	Secrets secretwatcher.SecretMap
}

type TranslatorPlugin interface {
	GetDependencies(cfg v1.Config) *Dependencies
}

type UpstreamPlugin interface {
	TranslatorPlugin
	Process(out *envoyapi.Cluster)
}

type FunctionPlugin interface {
	UpstreamPlugin
}

type RoutePlugin interface {
	TranslatorPlugin
}
