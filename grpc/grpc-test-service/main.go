package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/solo-io/gloo-plugins/grpc/grpc-test-service/bookstore"
	"google.golang.org/grpc"
)

func main() {
	port := flag.Int("p", 8080, "port")
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	bookstore.RegisterBookstoreServer(grpcServer, NewServer())
	log.Printf("listening on %v", *port)
	grpcServer.Serve(lis)
}
