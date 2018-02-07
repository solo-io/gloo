package translator

import (
	envoycache "github.com/envoyproxy/go-control-plane/pkg/cache"
	"github.com/solo-io/glue/pkg/api/types/v1"
	"github.com/solo-io/glue/pkg/endpointdiscovery"
	"github.com/solo-io/glue/pkg/secretwatcher"

	"github.com/hashicorp/go-multierror"
)

type UserResource interface {
	GetDependencies(cfg *v1.Config) DependenciesDescription
}

type DependenciesDescription interface {
	SecretRefs() []string
}

type Dependencies interface {
	Secrets() secretwatcher.SecretMap
}

type EnvoyNameConverter interface {
	ToEnvoyClusterName(upstreamName string) string
	ToEnvoyVhostName(*v1.VirtualHost) string
}

type ConfigStatus struct {
	Cfg v1.ConfigObject
	Err *multierror.Error
}

func NewConfigMultiError(Cfg v1.ConfigObject, err *multierror.Error) ConfigStatus {
	return ConfigStatus{Cfg: Cfg, Err: err}
}

func NewConfigOk(Cfg v1.ConfigObject) ConfigStatus {
	return ConfigStatus{Cfg: Cfg, Err: nil}
}

type Interface interface {
	UserResource

	Translate(cfg *v1.Config, secretMap secretwatcher.SecretMap, endpoints endpointdiscovery.EndpointGroups) (*envoycache.Snapshot, []ConfigStatus)
}
