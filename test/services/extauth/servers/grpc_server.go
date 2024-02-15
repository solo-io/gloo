package extauth_test_server

import (
	"context"

	. "github.com/onsi/gomega"
	passthrough_utils "github.com/solo-io/ext-auth-service/pkg/config/passthrough/test_utils"
	"github.com/solo-io/gloo/test/ginkgo/parallel"
)

var (
	baseExtauthPort = uint32(27000)
)

type GrpcServer struct {
	*passthrough_utils.GrpcAuthServer
	Port int
}

// NewGrpcServer Creates a wrapper around the passthrough_utils.GrpcAuthServer using an available port.
func NewGrpcServer(authServer *passthrough_utils.GrpcAuthServer) *GrpcServer {
	port := parallel.AdvancePortSafeListen(&baseExtauthPort)
	return &GrpcServer{
		GrpcAuthServer: authServer,
		Port:           int(port),
	}
}

// NewGrpcServer Creates a wrapper around the passthrough_utils.GrpcAuthServer on the given port.
func NewGrpcServerOnPort(authServer *passthrough_utils.GrpcAuthServer, port int) *GrpcServer {
	return &GrpcServer{
		GrpcAuthServer: authServer,
		Port:           port,
	}
}

func (g *GrpcServer) Start(ctx context.Context) {
	err := g.GrpcAuthServer.Start(g.Port)
	ExpectWithOffset(1, err).NotTo(HaveOccurred())

	go func(testContext context.Context) {
		for {
			select {
			case <-testContext.Done():
				g.Stop()
				return
			}
		}
	}(ctx)
}

func (g *GrpcServer) Stop() {
	g.GrpcAuthServer.Stop()
}

func (g *GrpcServer) GetAddress() string {
	return g.GrpcAuthServer.GetAddress()
}

func (g *GrpcServer) GetPort() int {
	return g.Port
}
