// +build wireinject

package devportal

import (
	"context"

	"github.com/google/wire"
	"github.com/solo-io/go-utils/envutils"
	"github.com/solo-io/solo-projects/projects/grpcserver/server/setup"
)

func InitDevPortal(ctx context.Context) (Registrar, error) {
	wire.Build(
		envutils.MustGetPodNamespace,
		setup.NewKubeConfig,
		ProviderSet)
	return &devPortalRegistrar{}, nil
}
