package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	bookstore "github.com/solo-io/gloo-plugins/grpc/grpc-test-service/bookstore/protos"
	"google.golang.org/grpc"
)

// requires https://github.com/googleapis/googleapis to be in /tmp/googleapis

//go:generate mkdir -p bookstore
//go:generate mkdir -p descriptors
//go:generate protoc -I${HOME}/workspace/googleapis -I. --include_source_info --go_out=plugins=grpc:bookstore   --descriptor_set_out=descriptors/proto.pb protos/bookstore.proto

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
