package configproto

import (
	extauthservice "github.com/solo-io/ext-auth-service/pkg/service"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/plugins/extauth"
)

type ConfigGenerator interface {
	GenerateConfig(resources []*extauth.ExtAuthConfig) (*extauthservice.Config, error)
}
