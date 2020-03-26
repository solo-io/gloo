package devportal

import (
	"github.com/solo-io/dev-portal/pkg/api/grpc/admin"
	"google.golang.org/grpc"
)

type Registrar interface {
	// Registers the dev portal services to the given gRPC server.
	RegisterTo(server *grpc.Server)
}

func NewRegistrar(
	portalService admin.PortalApiServer,
	apiDocService admin.ApiDocApiServer,
	userService admin.UserApiServer,
	groupService admin.GroupApiServer,
) Registrar {
	return &devPortalRegistrar{
		portalService: portalService,
		apiDocService: apiDocService,
		userService:   userService,
		groupService:  groupService,
	}
}

type devPortalRegistrar struct {
	portalService admin.PortalApiServer
	apiDocService admin.ApiDocApiServer
	userService   admin.UserApiServer
	groupService  admin.GroupApiServer
}

func (r *devPortalRegistrar) RegisterTo(server *grpc.Server) {
	admin.RegisterPortalApiServer(server, r.portalService)
	admin.RegisterApiDocApiServer(server, r.apiDocService)
	admin.RegisterUserApiServer(server, r.userService)
	admin.RegisterGroupApiServer(server, r.groupService)
}

// This registrar is used when the portal is not enabled.
type noOpRegistrar struct{}

func NewNoOpRegistrar() Registrar {
	return &noOpRegistrar{}
}

func (*noOpRegistrar) RegisterTo(_ *grpc.Server) {}
