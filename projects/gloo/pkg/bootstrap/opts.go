package bootstrap

import (
	"github.com/solo-io/gloo/projects/gloo/pkg/bootstrap"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/factory"
)

type EnterpriseOpts struct {
	Schemas         factory.ResourceClientFactory
	*bootstrap.Opts
}
