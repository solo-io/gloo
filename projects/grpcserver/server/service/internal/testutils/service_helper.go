package testutils

import (
	"fmt"
	"net"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc"
)

func TestVirtualservice(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Virtualservice Suite")
}

type registerServices func(s *grpc.Server)

func MustRunGrpcServer(f registerServices) (*grpc.Server, *grpc.ClientConn) {
	lis, err := net.Listen("tcp", ":0")
	Expect(err).NotTo(HaveOccurred())
	port := lis.Addr().(*net.TCPAddr).Port
	s := grpc.NewServer()
	f(s)
	go func() {
		defer GinkgoRecover()
		err = s.Serve(lis)
		Expect(err).NotTo(HaveOccurred())
	}()
	conn, err := grpc.Dial(fmt.Sprintf("localhost:%d", port), grpc.WithInsecure())
	Expect(err).NotTo(HaveOccurred())
	return s, conn
}
