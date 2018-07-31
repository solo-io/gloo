package apiclient_test

import (
	"testing"

	"google.golang.org/grpc"

	"fmt"
	"net"

	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"github.com/grpc-ecosystem/go-grpc-middleware/tags"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/solo-io/gloo/pkg/log"
	"github.com/solo-io/gloo/test/helpers"
	"github.com/solo-io/solo-kit/pkg/api/v1/apiserver"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/memory"
	"github.com/solo-io/solo-kit/test/mocks"
	"go.uber.org/zap"
)

func TestApiclient(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Apiclient Suite")
}

var (
	resourceClient = memory.NewResourceClient(&mocks.MockData{})
	port           = 1234
	server         *grpc.Server
)

var _ = BeforeSuite(func() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	helpers.Must(err)
	server = grpc.NewServer(grpc.StreamInterceptor(
		grpc_middleware.ChainStreamServer(
			grpc_ctxtags.StreamServerInterceptor(),
			grpc_zap.StreamServerInterceptor(zap.NewNop()),
			func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
				log.Printf("%v", info.FullMethod)
				return handler(srv, ss)
			},
		)))
	apiserver.NewApiServer(server, resourceClient)
	log.Printf("grpc listening on %v", port)
	go server.Serve(lis)
})

var _ = AfterSuite(func() {
	server.Stop()
})
