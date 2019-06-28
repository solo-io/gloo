package testutils

import (
	"net"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

type registerServices func(s *grpc.Server)

func MustRunGrpcServer(f registerServices) (*grpc.Server, *grpc.ClientConn) {
	lis := bufconn.Listen(1024 * 1024)
	mockDialer := func(string, time.Duration) (net.Conn, error) {
		return lis.Dial()
	}
	s := grpc.NewServer()
	f(s)
	go func() {
		defer GinkgoRecover()
		err := s.Serve(lis)
		Expect(err).NotTo(HaveOccurred())
	}()
	conn, err := grpc.Dial("mock", grpc.WithDialer(mockDialer), grpc.WithInsecure())
	Expect(err).NotTo(HaveOccurred())
	return s, conn
}
