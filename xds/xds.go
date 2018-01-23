package xds

import (
	"fmt"
	"net"

	envoyapi "github.com/envoyproxy/go-control-plane/api"
	envoycache "github.com/envoyproxy/go-control-plane/pkg/cache"
	xds "github.com/envoyproxy/go-control-plane/pkg/server"
	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/solo-io/glue/config"
	"github.com/solo-io/glue/pkg/log"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

const nodeKey = envoycache.Key("mock-node")

type hasher struct {
}

func (h hasher) Hash(node *envoyapi.Node) (envoycache.Key, error) {
	return nodeKey, nil
}

func RunXDS(gatewayConfig *config.Config, port int, configChanged <-chan bool) error {
	envoyConfig := envoycache.NewSimpleCache(hasher{}, func(key envoycache.Key) {
		log.Printf("CACHE: Key Updated: %s", key)
	})
	opts := []grpc_zap.Option{
		grpc_zap.WithDecider(func(fullMethodName string, err error) bool {
			// by default everything will be logged
			return true
		}),
	}
	grpcServer := grpc.NewServer(grpc.StreamInterceptor(
		grpc_middleware.ChainStreamServer(
			grpc_ctxtags.StreamServerInterceptor(),
			grpc_zap.StreamServerInterceptor(zap.NewNop(), opts...),
			func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
				log.Printf("STREAM: info: %v\n", info)
				log.Debugf("STREAM: serverStream: %v\n", ss)
				return handler(srv, ss)
			},
		)),
	)
	xdsSerevr := xds.NewServer(envoyConfig)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}
	envoyapi.RegisterAggregatedDiscoveryServiceServer(grpcServer, xdsSerevr)
	envoyapi.RegisterEndpointDiscoveryServiceServer(grpcServer, xdsSerevr)
	envoyapi.RegisterClusterDiscoveryServiceServer(grpcServer, xdsSerevr)
	envoyapi.RegisterRouteDiscoveryServiceServer(grpcServer, xdsSerevr)
	envoyapi.RegisterListenerDiscoveryServiceServer(grpcServer, xdsSerevr)
	log.Printf("xDS server listening on %d", port)
	go func() {
		if err = grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve grpc: %v", err)
		}
	}()
	for {
		select {
		case <-configChanged:
			configureCache(gatewayConfig.GetResources(), envoyConfig)
		}
	}
}

func configureCache(resources []config.EnvoyResources, config envoycache.Cache) {
	snapshot, err := createSnapshot(resources)
	if err != nil {
		log.Printf("ERROR:failed to create snapshot: %v", err)
	}
	config.SetSnapshot(nodeKey, snapshot)
	if err != nil {
		log.Printf("ERROR:failed to set snapshot: %v", err)
	}
	log.Debugf("set new snapshot: %v", snapshot)
}
