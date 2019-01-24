package configproto

import (
	extauthservice "github.com/solo-io/ext-auth-service/pkg/service"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/api/v1/plugins/extauth"
)

type ConfigGenerator interface {
	GenerateConfig(resources []*extauth.ExtAuthConfig) (*extauthservice.Config, error)
}
