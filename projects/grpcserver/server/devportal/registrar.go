package devportal

import (
	"github.com/solo-io/dev-portal/pkg/api/grpc/admin"
	"google.golang.org/grpc"
)

type Registrar interface {
	// Registers the dev portal services to the given gRPC server.
	RegisterTo(server *grpc.Server)
}

func NewRegistrar(portalService admin.PortalApiServer) Registrar {
	return &devPortalRegistrar{
		portalService: portalService,
	}
}

// TODO(marco): add remaining dev portal services here
type devPortalRegistrar struct {
	portalService admin.PortalApiServer
}

func (r *devPortalRegistrar) RegisterTo(server *grpc.Server) {
	admin.RegisterPortalApiServer(server, r.portalService)
}

// This registrar is used when the portal is not enabled.
type noOpRegistrar struct{}

func NewNoOpRegistrar() Registrar {
	return &noOpRegistrar{}
}

func (*noOpRegistrar) RegisterTo(server *grpc.Server) {
}
